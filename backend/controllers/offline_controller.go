package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OfflineController 离线控制器
type OfflineController struct {
	DB          *gorm.DB
	SyncService *services.OfflineSyncService
}

// NewOfflineController 创建离线控制器
func NewOfflineController(db *gorm.DB) *OfflineController {
	return &OfflineController{
		DB:          db,
		SyncService: services.NewOfflineSyncService(db),
	}
}

// RegisterRoutes 注册离线相关路由
func (ctrl *OfflineController) RegisterRoutes(rg *gin.RouterGroup) {
	offline := rg.Group("/offline")
	{
		// 获取离线缓存数据
		// GET /api/v1/offline/cache?device_id=xxx&data_type=xxx
		offline.GET("/cache", ctrl.GetCache)

		// 同步离线数据
		// POST /api/v1/offline/sync
		offline.POST("/sync", ctrl.SyncOfflineData)

		// 获取待同步队列
		// GET /api/v1/offline/queue?device_id=xxx
		offline.GET("/queue", ctrl.GetQueue)

		// 添加到同步队列
		// POST /api/v1/offline/queue
		offline.POST("/queue", ctrl.AddToQueue)

		// 确认同步完成
		// POST /api/v1/offline/queue/:id/confirm
		offline.POST("/queue/:id/confirm", ctrl.ConfirmSync)

		// 处理离线队列（批量同步）
		// POST /api/v1/offline/process
		offline.POST("/process", ctrl.ProcessQueue)
	}
}

// GetCache 获取离线缓存数据
// GET /api/v1/offline/cache?device_id=xxx&data_type=xxx
func (ctrl *OfflineController) GetCache(c *gin.Context) {
	deviceID := c.Query("device_id")
	dataType := c.Query("data_type")

	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "device_id is required",
		})
		return
	}

	caches, err := ctrl.SyncService.GetCachedData(deviceID, dataType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get cache data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      caches,
			"total":      len(caches),
			"device_id":  deviceID,
			"data_type":  dataType,
			"fetch_time": time.Now(),
		},
	})
}

// SyncOfflineData 同步离线数据
// POST /api/v1/offline/sync
// Body: { "device_id": "xxx", "operations": [...] }
func (ctrl *OfflineController) SyncOfflineData(c *gin.Context) {
	var req struct {
		DeviceID   string                     `json:"device_id" binding:"required"`
		Operations []models.OfflineOperation `json:"operations"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 处理离线操作
	results := map[string]interface{}{
		"device_id": req.DeviceID,
		"succeeded": 0,
		"failed":    0,
		"items":     []map[string]interface{}{},
	}

	for _, op := range req.Operations {
		// 创建离线操作记录
		op.DeviceID = req.DeviceID
		op.Status = "pending"

		if err := ctrl.DB.Create(&op).Error; err != nil {
			results["failed"] = results["failed"].(int) + 1
			results["items"] = append(results["items"].([]map[string]interface{}), map[string]interface{}{
				"operation_id": op.ID,
				"status":       "error",
				"error":        err.Error(),
			})
			continue
		}

		results["succeeded"] = results["succeeded"].(int) + 1
		results["items"] = append(results["items"].([]map[string]interface{}), map[string]interface{}{
			"operation_id": op.ID,
			"status":       "queued",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": results,
	})
}

// GetQueue 获取待同步队列
// GET /api/v1/offline/queue?device_id=xxx&status=pending
func (ctrl *OfflineController) GetQueue(c *gin.Context) {
	deviceID := c.Query("device_id")
	status := c.DefaultQuery("status", "pending")

	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "device_id is required",
		})
		return
	}

	var operations []models.OfflineOperation
	query := ctrl.DB.Where("device_id = ?", deviceID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at ASC").Find(&operations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get queue",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      operations,
			"total":      len(operations),
			"device_id":  deviceID,
			"status":     status,
			"fetch_time": time.Now(),
		},
	})
}

// AddToQueue 添加到同步队列
// POST /api/v1/offline/queue
// Body: { "device_id": "xxx", "operation": "control", "payload": "{}" }
func (ctrl *OfflineController) AddToQueue(c *gin.Context) {
	var req struct {
		DeviceID  string `json:"device_id" binding:"required"`
		Operation string `json:"operation" binding:"required"`
		Payload   string `json:"payload" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	op := models.OfflineOperation{
		DeviceID:  req.DeviceID,
		Operation: req.Operation,
		Payload:   req.Payload,
		CreatedAt: time.Now(),
		Status:    "pending",
	}

	if err := ctrl.DB.Create(&op).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to add to queue",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"operation_id": op.ID,
			"device_id":    op.DeviceID,
			"operation":    op.Operation,
			"status":       op.Status,
			"created_at":   op.CreatedAt,
		},
	})
}

// ConfirmSync 确认同步完成
// POST /api/v1/offline/queue/:id/confirm
func (ctrl *OfflineController) ConfirmSync(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid operation id",
		})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&models.OfflineOperation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    "completed",
		"synced_at": &now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to confirm sync",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "sync confirmed",
	})
}

// ProcessQueue 处理离线队列（批量同步）
// POST /api/v1/offline/process
// Body: { "device_id": "xxx" }
func (ctrl *OfflineController) ProcessQueue(c *gin.Context) {
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	results, err := ctrl.SyncService.ProcessOfflineQueue(req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to process queue",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": results,
	})
}
