package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BatchController 批量操作控制器
type BatchController struct {
	DB *gorm.DB
}

// RegisterBatchRoutes 注册批量操作路由
func (ctrl *BatchController) RegisterBatchRoutes(api *gin.RouterGroup) {
	api.POST("/batch/devices/actions", ctrl.DeviceActions)
	api.POST("/batch/devices/shadow", ctrl.BatchShadowUpdate)
	api.GET("/batch/tasks/:task_id", ctrl.GetTask)
	api.GET("/batch/tasks", ctrl.ListTasks)
}

// DeviceActions 批量设备动作（升级/重启/下发配置）
// POST /api/v1/batch/devices/actions
func (ctrl *BatchController) DeviceActions(c *gin.Context) {
	var req models.BatchTaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	if len(req.DeviceIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "device_ids cannot be empty"})
		return
	}
	if len(req.DeviceIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "maximum 100 devices per batch"})
		return
	}

	taskID := "batch-" + uuid.New().String()
	creatorID := batchGetUserIDFromContext(c)

	task := models.BatchTask{
		TaskID:    taskID,
		TaskType:  req.Action,
		Total:     len(req.DeviceIDs),
		Success:   0,
		Failed:    0,
		Pending:   len(req.DeviceIDs),
		Status:    "running",
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		Results:   "[]",
	}

	if err := ctrl.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create batch task: " + err.Error()})
		return
	}

	// 异步执行批量操作
	go ctrl.executeBatchActions(taskID, req.DeviceIDs, req.Action, req.Params)

	c.JSON(http.StatusAccepted, gin.H{
		"code": 0,
		"data": gin.H{
			"task_id": taskID,
			"total":   len(req.DeviceIDs),
			"status":  "running",
			"message": "Batch task started",
		},
	})
}

// executeBatchActions 异步执行批量设备动作
func (ctrl *BatchController) executeBatchActions(taskID string, deviceIDs []string, action string, params map[string]interface{}) {
	results := make([]models.BatchTaskDeviceResult, 0, len(deviceIDs))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 限制并发数
	semaphore := make(chan struct{}, 20)

	for _, deviceID := range deviceIDs {
		wg.Add(1)
		go func(did string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := models.BatchTaskDeviceResult{DeviceID: did, Status: "pending"}

			var err error
			switch action {
			case "upgrade":
				err = ctrl.executeUpgrade(did, params)
			case "reboot":
				err = ctrl.executeReboot(did)
			case "config":
				err = ctrl.executeConfig(did, params)
			default:
				err = fmt.Errorf("unknown action: %s", action)
			}

			if err != nil {
				result.Status = "failed"
				result.Error = err.Error()
			} else {
				result.Status = "success"
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(deviceID)
	}

	wg.Wait()

	// 统计结果
	success := 0
	failed := 0
	for _, r := range results {
		if r.Status == "success" {
			success++
		} else {
			failed++
		}
	}

	resultsJSON, _ := json.Marshal(results)
	completedAt := time.Now()

	ctrl.DB.Model(&models.BatchTask{}).Where("task_id = ?", taskID).Updates(map[string]interface{}{
		"success":      success,
		"failed":       failed,
		"pending":      0,
		"status":       "completed",
		"results":      string(resultsJSON),
		"completed_at": completedAt,
	})
}

// executeUpgrade 执行固件升级
func (ctrl *BatchController) executeUpgrade(deviceID string, params map[string]interface{}) error {
	client := mqtt.GetGlobalMQTTClient()
	if client == nil {
		return fmt.Errorf("MQTT client not available")
	}

	firmware := ""
	if v, ok := params["firmware_version"]; ok {
		firmware = fmt.Sprintf("%v", v)
	}

	payload := map[string]interface{}{
		"type":             "ota_upgrade",
		"firmware_version": firmware,
		"timestamp":        time.Now().Unix(),
	}

	data, _ := json.Marshal(payload)
	topic := fmt.Sprintf("/device/%s/down/cmd", deviceID)
	token := client.Publish(topic, 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// executeReboot 执行设备重启
func (ctrl *BatchController) executeReboot(deviceID string) error {
	client := mqtt.GetGlobalMQTTClient()
	if client == nil {
		return fmt.Errorf("MQTT client not available")
	}

	payload := map[string]interface{}{
		"type":      "reboot",
		"timestamp": time.Now().Unix(),
	}

	data, _ := json.Marshal(payload)
	topic := fmt.Sprintf("/device/%s/down/cmd", deviceID)
	token := client.Publish(topic, 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// executeConfig 下发配置
func (ctrl *BatchController) executeConfig(deviceID string, params map[string]interface{}) error {
	client := mqtt.GetGlobalMQTTClient()
	if client == nil {
		return fmt.Errorf("MQTT client not available")
	}

	payload := map[string]interface{}{
		"type":      "config_update",
		"config":    params,
		"timestamp": time.Now().Unix(),
	}

	data, _ := json.Marshal(payload)
	topic := fmt.Sprintf("/device/%s/down/cmd", deviceID)
	token := client.Publish(topic, 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// BatchShadowUpdate 批量更新设备影子
// POST /api/v1/batch/devices/shadow
func (ctrl *BatchController) BatchShadowUpdate(c *gin.Context) {
	var req struct {
		DeviceIDs []string               `json:"device_ids" binding:"required"`
		Shadow    map[string]interface{}  `json:"shadow" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	if len(req.DeviceIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "device_ids cannot be empty"})
		return
	}
	if len(req.DeviceIDs) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "maximum 100 devices per batch"})
		return
	}

	taskID := "batch-shadow-" + uuid.New().String()
	creatorID := batchGetUserIDFromContext(c)

	task := models.BatchTask{
		TaskID:    taskID,
		TaskType:  "shadow_update",
		Total:     len(req.DeviceIDs),
		Success:   0,
		Failed:    0,
		Pending:   len(req.DeviceIDs),
		Status:    "running",
		CreatorID: creatorID,
		CreatedAt: time.Now(),
	}

	if err := ctrl.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create batch task: " + err.Error()})
		return
	}

	go ctrl.executeShadowUpdate(taskID, req.DeviceIDs, req.Shadow)

	c.JSON(http.StatusAccepted, gin.H{
		"code": 0,
		"data": gin.H{
			"task_id": taskID,
			"total":   len(req.DeviceIDs),
			"status":  "running",
		},
	})
}

// executeShadowUpdate 异步执行批量设备影子上报
func (ctrl *BatchController) executeShadowUpdate(taskID string, deviceIDs []string, shadow map[string]interface{}) {
	results := make([]models.BatchTaskDeviceResult, 0, len(deviceIDs))
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 20)

	for _, deviceID := range deviceIDs {
		wg.Add(1)
		go func(did string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := models.BatchTaskDeviceResult{DeviceID: did, Status: "pending"}

			var device models.Device
			if err := ctrl.DB.Where("device_id = ?", did).First(&device).Error; err != nil {
				result.Status = "failed"
				result.Error = "device not found"
			} else {
				shadowJSON, _ := json.Marshal(shadow)
				device.DesiredState = string(shadowJSON)
				if err := ctrl.DB.Save(&device).Error; err != nil {
					result.Status = "failed"
					result.Error = err.Error()
				} else {
					result.Status = "success"
				}
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(deviceID)
	}

	wg.Wait()

	success := 0
	failed := 0
	for _, r := range results {
		if r.Status == "success" {
			success++
		} else {
			failed++
		}
	}

	resultsJSON, _ := json.Marshal(results)
	ctrl.DB.Model(&models.BatchTask{}).Where("task_id = ?", taskID).Updates(map[string]interface{}{
		"success":      success,
		"failed":       failed,
		"pending":      0,
		"status":       "completed",
		"results":      string(resultsJSON),
		"completed_at": time.Now(),
	})
}

// GetTask 查询批量任务状态
// GET /api/v1/batch/tasks/:task_id
func (ctrl *BatchController) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")

	var task models.BatchTask
	if err := ctrl.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Batch task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": task.ToResponse(),
	})
}

// ListTasks 批量任务历史
// GET /api/v1/batch/tasks
func (ctrl *BatchController) ListTasks(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	taskType := c.Query("task_type")
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.BatchTask{})

	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	var total int64
	query.Count(&total)

	var tasks []models.BatchTask
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query batch tasks: " + err.Error(),
		})
		return
	}

	responses := make([]*models.BatchTaskResponse, len(tasks))
	for i := range tasks {
		responses[i] = tasks[i].ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      responses,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ============ 辅助函数 ============

func batchGetUserIDFromContext(c *gin.Context) uint {
	if uid, exists := c.Get("user_id"); exists {
		switch v := uid.(type) {
		case uint:
			return v
		case int:
			return uint(v)
		case int64:
			return uint(v)
		case float64:
			return uint(v)
		case string:
			var id uint
			for _, ch := range v {
				if ch >= '0' && ch <= '9' {
					id = id*10 + uint(ch-'0')
				}
			}
			return id
		}
	}
	return 0
}

func parseUint64(s string) uint64 {
	var n uint64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + uint64(c-'0')
		}
	}
	return n
}
