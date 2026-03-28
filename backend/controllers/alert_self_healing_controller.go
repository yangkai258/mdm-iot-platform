package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertSelfHealingController 告警自愈控制器
type AlertSelfHealingController struct {
	DB *gorm.DB
}

// NewAlertSelfHealingController 创建控制器
func NewAlertSelfHealingController(db *gorm.DB) *AlertSelfHealingController {
	return &AlertSelfHealingController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *AlertSelfHealingController) RegisterRoutes(rg *gin.RouterGroup) {
	// 自愈方案管理
	rg.GET("/self-healing", ctrl.ListSelfHealing)
	rg.GET("/self-healing/:id", ctrl.GetSelfHealing)
	rg.POST("/self-healing", ctrl.CreateSelfHealing)
	rg.PUT("/self-healing/:id", ctrl.UpdateSelfHealing)
	rg.DELETE("/self-healing/:id", ctrl.DeleteSelfHealing)
	rg.GET("/self-healing/types", ctrl.GetAlertTypes)

	// 自愈执行记录
	rg.GET("/self-healing/records", ctrl.ListRecords)
	rg.GET("/self-healing/records/:id", ctrl.GetRecord)
	rg.POST("/self-healing/:id/execute", ctrl.ExecuteSelfHealing)

	// 根据告警类型获取建议
	rg.GET("/self-healing/advice/:alert_type", ctrl.GetAdviceForAlert)
}

// ListSelfHealing 获取自愈方案列表
func (ctrl *AlertSelfHealingController) ListSelfHealing(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	alertType := c.Query("alert_type")
	severity := c.Query("severity")
	isActive := c.Query("is_active")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := ctrl.DB.Model(&models.AlertSelfHealing{})

	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var total int64
	var list []models.AlertSelfHealing
	query.Count(&total)

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取自愈方案列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析StepsJSON
	for i := range list {
		if list[i].StepsJSON != "" {
			var steps []models.SelfHealingStep
			if err := json.Unmarshal([]byte(list[i].StepsJSON), &steps); err == nil {
				list[i].SelfHealingSteps = steps
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSelfHealing 获取单个自愈方案
func (ctrl *AlertSelfHealingController) GetSelfHealing(c *gin.Context) {
	id := c.Param("id")

	var item models.AlertSelfHealing
	if err := ctrl.DB.Where("id = ?", id).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "自愈方案不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取自愈方案失败",
		})
		return
	}

	if item.StepsJSON != "" {
		var steps []models.SelfHealingStep
		if err := json.Unmarshal([]byte(item.StepsJSON), &steps); err == nil {
			item.SelfHealingSteps = steps
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": item,
	})
}

// CreateSelfHealing 创建自愈方案
func (ctrl *AlertSelfHealingController) CreateSelfHealing(c *gin.Context) {
	var req struct {
		AlertType       string                  `json:"alert_type" binding:"required"`
		AlertSubType    string                  `json:"alert_sub_type"`
		Severity        int                     `json:"severity"`
		Title           string                  `json:"title" binding:"required"`
		RootCause       string                  `json:"root_cause"`
		Recommendation  string                  `json:"recommendation" binding:"required"`
		SelfHealingSteps []models.SelfHealingStep `json:"self_healing_steps"`
		Tags            string                  `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	stepsJSON := ""
	if len(req.SelfHealingSteps) > 0 {
		if data, err := json.Marshal(req.SelfHealingSteps); err == nil {
			stepsJSON = string(data)
		}
	}

	item := models.AlertSelfHealing{
		AlertType:       req.AlertType,
		AlertSubType:    req.AlertSubType,
		Severity:        req.Severity,
		Title:           req.Title,
		RootCause:       req.RootCause,
		Recommendation:  req.Recommendation,
		StepsJSON:       stepsJSON,
		IsActive:        true,
		Tags:            req.Tags,
		CreatedBy:       c.GetString("username"),
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建自愈方案失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析steps返回
	if stepsJSON != "" {
		var steps []models.SelfHealingStep
		json.Unmarshal([]byte(stepsJSON), &steps)
		item.SelfHealingSteps = steps
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    item,
	})
}

// UpdateSelfHealing 更新自愈方案
func (ctrl *AlertSelfHealingController) UpdateSelfHealing(c *gin.Context) {
	id := c.Param("id")

	var item models.AlertSelfHealing
	if err := ctrl.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "自愈方案不存在",
		})
		return
	}

	var req struct {
		AlertType        string                   `json:"alert_type"`
		AlertSubType     string                   `json:"alert_sub_type"`
		Severity         int                      `json:"severity"`
		Title            string                   `json:"title"`
		RootCause        string                   `json:"root_cause"`
		Recommendation   string                   `json:"recommendation"`
		SelfHealingSteps []models.SelfHealingStep `json:"self_healing_steps"`
		Tags             string                   `json:"tags"`
		IsActive         *bool                    `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	updates := map[string]interface{}{}

	if req.AlertType != "" {
		updates["alert_type"] = req.AlertType
	}
	if req.AlertSubType != "" {
		updates["alert_sub_type"] = req.AlertSubType
	}
	if req.Severity > 0 {
		updates["severity"] = req.Severity
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.RootCause != "" {
		updates["root_cause"] = req.RootCause
	}
	if req.Recommendation != "" {
		updates["recommendation"] = req.Recommendation
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(req.SelfHealingSteps) > 0 {
		if data, err := json.Marshal(req.SelfHealingSteps); err == nil {
			updates["steps_json"] = string(data)
		}
	}

	if err := ctrl.DB.Model(&item).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
			"error":   err.Error(),
		})
		return
	}

	// 重新查询获取完整数据
	ctrl.DB.Where("id = ?", id).First(&item)
	if item.StepsJSON != "" {
		var steps []models.SelfHealingStep
		json.Unmarshal([]byte(item.StepsJSON), &steps)
		item.SelfHealingSteps = steps
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    item,
	})
}

// DeleteSelfHealing 删除自愈方案
func (ctrl *AlertSelfHealingController) DeleteSelfHealing(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.DB.Where("id = ?", id).Delete(&models.AlertSelfHealing{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetAlertTypes 获取告警类型列表
func (ctrl *AlertSelfHealingController) GetAlertTypes(c *gin.Context) {
	// 从数据库获取所有告警类型
	var types []string
	ctrl.DB.Model(&models.AlertSelfHealing{}).
		Where("is_active = ?", true).
		Pluck("DISTINCT alert_type", &types)

	// 添加常见告警类型
	commonTypes := []string{
		"device_offline", "device_battery_low", "device_temperature_high",
		"network_unstable", "firmware_update_failed", "ota_failed",
		"ai_response_slow", "ai_quality_low", "pet_health_warning",
		"emotion_negative", "behavior_abnormal", "system_error",
	}

	type Result struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	var results []Result
	for _, t := range types {
		var count int64
		ctrl.DB.Model(&models.AlertSelfHealing{}).Where("alert_type = ?", t).Count(&count)
		results = append(results, Result{Type: t, Count: count})
	}
	_ = commonTypes // 可以在前端展示推荐类型

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"types":      results,
			"common":     commonTypes,
		},
	})
}

// ListRecords 获取自愈执行记录
func (ctrl *AlertSelfHealingController) ListRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	alertType := c.Query("alert_type")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}

	query := ctrl.DB.Model(&models.AlertSelfHealingRecord{})

	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	var list []models.AlertSelfHealingRecord
	query.Count(&total)

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRecord 获取单个执行记录
func (ctrl *AlertSelfHealingController) GetRecord(c *gin.Context) {
	id := c.Param("id")

	var record models.AlertSelfHealingRecord
	if err := ctrl.DB.Where("id = ?", id).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": record,
	})
}

// ExecuteSelfHealing 执行自愈方案
func (ctrl *AlertSelfHealingController) ExecuteSelfHealing(c *gin.Context) {
	id := c.Param("id")
	alertType := c.Query("alert_type")
	triggerCondition := c.Query("condition")

	var selfHealing models.AlertSelfHealing
	if err := ctrl.DB.Where("id = ? AND is_active = ?", id, true).First(&selfHealing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "自愈方案不存在或已禁用",
		})
		return
	}

	// 解析步骤
	var steps []models.SelfHealingStep
	if selfHealing.StepsJSON != "" {
		json.Unmarshal([]byte(selfHealing.StepsJSON), &steps)
	}

	// 创建执行记录
	now := time.Now()
	record := models.AlertSelfHealingRecord{
		AlertID:         fmt.Sprintf("ALT-%d-%d", id, now.Unix()),
		SelfHealingID:   selfHealing.ID,
		AlertType:       selfHealing.AlertType,
		TriggerCondition: triggerCondition,
		StepsTotal:      len(steps),
		Status:          "running",
		ExecutedBy:      "system",
		StartedAt:        &now,
	}

	if alertType != "" {
		record.AlertType = alertType
	}

	ctrl.DB.Create(&record)

	// 模拟执行（实际应该按步骤执行）
	go func() {
		time.Sleep(2 * time.Second)
		completedAt := time.Now()
		result := "自愈方案执行完成\n\n建议：" + selfHealing.Recommendation

		updates := map[string]interface{}{
			"status":       "success",
			"steps_executed": len(steps),
			"result":       result,
			"completed_at":  completedAt,
			"duration":     int(completedAt.Sub(now).Seconds()),
		}

		ctrl.DB.Model(&models.AlertSelfHealingRecord{}).Where("id = ?", record.ID).Updates(updates)

		// 更新使用统计
		ctrl.DB.Model(&models.AlertSelfHealing{}).Where("id = ?", id).Updates(map[string]interface{}{
			"used_count":    selfHealing.UsedCount + 1,
			"success_count": selfHealing.SuccessCount + 1,
		})
	}()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "自愈方案已启动执行",
		"data": gin.H{
			"record_id":     record.ID,
			"alert_id":      record.AlertID,
			"steps_total":   len(steps),
			"self_healing":  selfHealing.Recommendation,
		},
	})
}

// GetAdviceForAlert 根据告警类型获取建议
func (ctrl *AlertSelfHealingController) GetAdviceForAlert(c *gin.Context) {
	alertType := c.Param("alert_type")
	severity := c.DefaultQuery("severity", "0")

	query := ctrl.DB.Model(&models.AlertSelfHealing{}).
		Where("alert_type = ? AND is_active = ?", alertType, true)

	if severity != "0" {
		query = query.Where("severity <= ?", severity)
	}

	var list []models.AlertSelfHealing
	query.Order("success_rate DESC, severity ASC").Limit(5).Find(&list)

	if len(list) == 0 {
		// 返回通用建议
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"advice": gin.H{
					"alert_type":    alertType,
					"title":         "未找到专门建议",
					"recommendation": "请检查设备状态，联系技术支持获取帮助。",
				},
				"has_custom": false,
			},
		})
		return
	}

	// 解析steps
	for i := range list {
		if list[i].StepsJSON != "" {
			var steps []models.SelfHealingStep
			json.Unmarshal([]byte(list[i].StepsJSON), &steps)
			list[i].SelfHealingSteps = steps
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"advice":    list[0], // 返回最匹配的
			"alternatives": list[1:],
			"has_custom": true,
		},
	})
}
