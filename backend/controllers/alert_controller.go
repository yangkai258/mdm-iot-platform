package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertController 告警控制器
type AlertController struct {
	DB *gorm.DB
}

// GetRules 获取告警规则列表
func (c *AlertController) GetRules(ctx *gin.Context) {
	var rules []models.DeviceAlertRule
	c.DB.Order("id DESC").Find(&rules)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": rules},
	})
}

// CreateRule 创建告警规则
func (c *AlertController) CreateRule(ctx *gin.Context) {
	var rule models.DeviceAlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// GetAlerts 获取告警记录
func (c *AlertController) GetAlerts(ctx *gin.Context) {
	query := c.DB.Model(&models.DeviceAlert{})

	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var alerts []models.DeviceAlert
	query.Order("id DESC").Limit(100).Find(&alerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": alerts},
	})
}

// CheckAlerts 检查设备是否触发告警
func CheckAlerts(db *gorm.DB, deviceID string, data map[string]interface{}) {
	// 查询所有启用的规则
	var rules []models.DeviceAlertRule
	db.Where("enabled = ? AND (device_id = ? OR device_id = '')", true, deviceID).Find(&rules)

	for _, rule := range rules {
		triggered := false
		var triggerVal float64

		switch rule.AlertType {
		case "battery_low":
			if v, ok := data["battery"].(float64); ok {
				triggerVal = v
				triggered = evaluateCondition(rule.Condition, v, rule.Threshold)
			}
		case "offline":
			if v, ok := data["is_online"].(bool); ok {
				triggerVal = 0
				if !v && rule.Condition == "=" {
					triggered = true
				}
			}
		}

		if triggered {
			alert := models.DeviceAlert{
				RuleID:     rule.ID,
				DeviceID:   deviceID,
				AlertType:  rule.AlertType,
				Severity:   rule.Severity,
				Message:    rule.Name,
				TriggerVal: triggerVal,
				Threshold:  rule.Threshold,
				Status:     1,
			}
			db.Create(&alert)
		}
	}
}

func evaluateCondition(cond string, val, threshold float64) bool {
	switch cond {
	case "<":
		return val < threshold
	case ">":
		return val > threshold
	case "=":
		return val == threshold
	case "<=":
		return val <= threshold
	case ">=":
		return val >= threshold
	}
	return false
}

// DashboardStats 统计数据
type DashboardStats struct {
	TotalDevices    int64 `json:"total_devices"`
	OnlineDevices  int64 `json:"online_devices"`
	OfflineDevices int64 `json:"offline_devices"`
	TotalAlerts    int64 `json:"total_alerts"`
	PendingAlerts  int64 `json:"pending_alerts"`
}

// GetDashboardStats 获取大盘统计
func (c *AlertController) GetDashboardStats(ctx *gin.Context) {
	var stats DashboardStats

	// 设备统计
	c.DB.Model(&models.Device{}).Count(&stats.TotalDevices)
	c.DB.Model(&models.Device{}).Where("is_online = ?", true).Count(&stats.OnlineDevices)
	stats.OfflineDevices = stats.TotalDevices - stats.OnlineDevices

	// 告警统计
	c.DB.Model(&models.DeviceAlert{}).Count(&stats.TotalAlerts)
	c.DB.Model(&models.DeviceAlert{}).Where("status = ?", 1).Count(&stats.PendingAlerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
