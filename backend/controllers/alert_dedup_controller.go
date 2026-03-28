package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertDedupController 告警去重控制器
type AlertDedupController struct {
	DB           *gorm.DB
	DedupService *services.AlertDeduplicationService
}

// NewAlertDedupController 创建控制器
func NewAlertDedupController(db *gorm.DB) *AlertDedupController {
	return &AlertDedupController{
		DB:           db,
		DedupService: services.NewAlertDeduplicationService(db),
	}
}

// RegisterRoutes 注册路由
func (ctrl *AlertDedupController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/alerts/dedup/rules", ctrl.ListRules)
	rg.GET("/alerts/dedup/rules/:id", ctrl.GetRule)
	rg.POST("/alerts/dedup/rules", ctrl.CreateRule)
	rg.PUT("/alerts/dedup/rules/:id", ctrl.UpdateRule)
	rg.DELETE("/alerts/dedup/rules/:id", ctrl.DeleteRule)
	rg.GET("/alerts/dedup/records", ctrl.ListRecords)
	rg.GET("/alerts/dedup/stats", ctrl.GetStats)
	rg.POST("/alerts/check", ctrl.CheckAlert)
}

// ListRules 获取去重规则列表
func (ctrl *AlertDedupController) ListRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	alertType := c.Query("alert_type")
	isActive := c.Query("is_active")

	query := ctrl.DB.Model(&models.AlertDeduplicationRule{})

	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var total int64
	var list []models.AlertDeduplicationRule
	query.Count(&total)

	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

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

// GetRule 获取单个规则
func (ctrl *AlertDedupController) GetRule(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertDeduplicationRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// CreateRule 创建去重规则
func (ctrl *AlertDedupController) CreateRule(c *gin.Context) {
	var req struct {
		AlertType          string `json:"alert_type" binding:"required"`
		DevicePattern     string `json:"device_pattern"`
		SeverityMin       int    `json:"severity_min"`
		SeverityMax       int    `json:"severity_max"`
		DedupWindowSeconds int    `json:"dedup_window_seconds"`
		DedupStrategy     string `json:"dedup_strategy"`
		MaxCountPerWindow int    `json:"max_count_per_window"`
		SuppressionType   string `json:"suppression_type"`
		Description       string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	ruleID := fmt.Sprintf("DEDUP-%s-%d", req.AlertType, time.Now().Unix())

	if req.SeverityMin == 0 {
		req.SeverityMin = 1
	}
	if req.SeverityMax == 0 {
		req.SeverityMax = 5
	}
	if req.DedupWindowSeconds == 0 {
		req.DedupWindowSeconds = 300
	}
	if req.DedupStrategy == "" {
		req.DedupStrategy = "first"
	}
	if req.MaxCountPerWindow == 0 {
		req.MaxCountPerWindow = 1
	}
	if req.SuppressionType == "" {
		req.SuppressionType = "none"
	}

	rule := models.AlertDeduplicationRule{
		RuleID:            ruleID,
		AlertType:         req.AlertType,
		DevicePattern:    req.DevicePattern,
		SeverityMin:      req.SeverityMin,
		SeverityMax:      req.SeverityMax,
		DedupWindowSeconds: req.DedupWindowSeconds,
		DedupStrategy:    req.DedupStrategy,
		MaxCountPerWindow: req.MaxCountPerWindow,
		SuppressionType:   req.SuppressionType,
		IsActive:         true,
		Description:     req.Description,
		CreatedBy:       c.GetString("username"),
	}

	if err := ctrl.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": rule})
}

// UpdateRule 更新规则
func (ctrl *AlertDedupController) UpdateRule(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertDeduplicationRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}

	var req struct {
		DevicePattern     string `json:"device_pattern"`
		SeverityMin       int    `json:"severity_min"`
		SeverityMax       int    `json:"severity_max"`
		DedupWindowSeconds int    `json:"dedup_window_seconds"`
		DedupStrategy     string `json:"dedup_strategy"`
		MaxCountPerWindow int    `json:"max_count_per_window"`
		SuppressionType   string `json:"suppression_type"`
		IsActive         *bool  `json:"is_active"`
		Description      string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}

	if req.DevicePattern != "" {
		updates["device_pattern"] = req.DevicePattern
	}
	if req.SeverityMin > 0 {
		updates["severity_min"] = req.SeverityMin
	}
	if req.SeverityMax > 0 {
		updates["severity_max"] = req.SeverityMax
	}
	if req.DedupWindowSeconds > 0 {
		updates["dedup_window_seconds"] = req.DedupWindowSeconds
	}
	if req.DedupStrategy != "" {
		updates["dedup_strategy"] = req.DedupStrategy
	}
	if req.MaxCountPerWindow > 0 {
		updates["max_count_per_window"] = req.MaxCountPerWindow
	}
	if req.SuppressionType != "" {
		updates["suppression_type"] = req.SuppressionType
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	ctrl.DB.Model(&rule).Updates(updates)
	ctrl.DB.Where("id = ?", rule.ID).First(&rule)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": rule})
}

// DeleteRule 删除规则
func (ctrl *AlertDedupController) DeleteRule(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertDeduplicationRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}

	ctrl.DB.Delete(&rule)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListRecords 获取去重记录
func (ctrl *AlertDedupController) ListRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	ruleID := c.Query("rule_id")
	deviceID := c.Query("device_id")

	query := ctrl.DB.Model(&models.AlertDeduplicationRecord{})

	if ruleID != "" {
		query = query.Where("rule_id = ?", ruleID)
	}
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	var total int64
	var list []models.AlertDeduplicationRecord
	query.Count(&total)

	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

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

// GetStats 获取去重统计
func (ctrl *AlertDedupController) GetStats(c *gin.Context) {
	startTime := time.Now().AddDate(0, 0, -7) // 最近7天
	endTime := time.Now()

	total, deduplicated := ctrl.DedupService.GetDeduplicationStats(startTime, endTime)

	var totalRules int64
	ctrl.DB.Model(&models.AlertDeduplicationRule{}).Count(&totalRules)

	var activeRules int64
	ctrl.DB.Model(&models.AlertDeduplicationRule{}).Where("is_active = ?", true).Count(&activeRules)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period": gin.H{
				"start": startTime,
				"end":   endTime,
			},
			"total_alerts":    total,
			"deduplicated":    deduplicated,
			"dedup_rate":     float64(deduplicated) / float64(total) * 100,
			"total_rules":     totalRules,
			"active_rules":    activeRules,
		},
	})
}

// CheckAlert 检查告警是否应该处理
func (ctrl *AlertDedupController) CheckAlert(c *gin.Context) {
	var req struct {
		AlertType string `json:"alert_type" binding:"required"`
		DeviceID string `json:"device_id" binding:"required"`
		Severity int    `json:"severity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Severity == 0 {
		req.Severity = 2 // 默认中等严重程度
	}

	shouldProcess, rule := ctrl.DedupService.ShouldProcessAlert(req.AlertType, req.DeviceID, req.Severity)

	result := gin.H{
		"should_process": shouldProcess,
		"alert_type":     req.AlertType,
		"device_id":      req.DeviceID,
		"severity":       req.Severity,
	}

	if rule != nil {
		result["rule"] = gin.H{
			"rule_id":        rule.RuleID,
			"dedup_strategy": rule.DedupStrategy,
			"window_seconds":  rule.DedupWindowSeconds,
			"max_count":      rule.MaxCountPerWindow,
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}
