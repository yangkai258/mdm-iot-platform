package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

// UpdateRule 更新告警规则
func (c *AlertController) UpdateRule(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.DeviceAlertRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}

	var input models.DeviceAlertRule
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	rule.Name = input.Name
	rule.DeviceID = input.DeviceID
	rule.AlertType = input.AlertType
	rule.Condition = input.Condition
	rule.Threshold = input.Threshold
	rule.Severity = input.Severity
	rule.Enabled = input.Enabled
	rule.NotifyWays = input.NotifyWays
	rule.Remark = input.Remark

	if err := c.DB.Save(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// DeleteRule 删除告警规则
func (c *AlertController) DeleteRule(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.DeviceAlertRule{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
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
	if alertType := ctx.Query("alert_type"); alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if severity := ctx.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	var alerts []models.DeviceAlert
	query.Order("id DESC").Limit(100).Find(&alerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": alerts},
	})
}

// ConfirmAlert 确认告警
func (c *AlertController) ConfirmAlert(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id") // 从JWT中获取用户ID

	var alert models.DeviceAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
		return
	}

	if alert.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只能确认未处理的告警"})
		return
	}

	now := time.Now()
	alert.Status = 2 // 已确认
	alert.ConfirmedAt = &now
	alert.ConfirmedBy = userID

	if err := c.DB.Save(&alert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "确认失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": alert, "message": "告警已确认"})
}

// ResolveAlert 解决告警
func (c *AlertController) ResolveAlert(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	var alert models.DeviceAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
		return
	}

	if alert.Status == 3 || alert.Status == 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "告警已解决或忽略"})
		return
	}

	now := time.Now()
	alert.Status = 3 // 已解决
	alert.ResolvedAt = &now
	alert.ResolvedBy = userID

	if err := c.DB.Save(&alert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解决失败"})
		return
	}

	// 如果是地理围栏告警，同步更新 geofence_alerts 表
	if alert.AlertType == "geofence_violation" {
		c.DB.Model(&models.GeofenceAlert{}).Where("alert_id = ?", alert.ID).Updates(map[string]interface{}{
			"status": 3,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": alert, "message": "告警已解决"})
}

// IgnoreAlert 忽略告警
func (c *AlertController) IgnoreAlert(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	var alert models.DeviceAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
		return
	}

	if alert.Status == 3 || alert.Status == 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "告警已解决或忽略"})
		return
	}

	now := time.Now()
	alert.Status = 4 // 已忽略
	alert.IgnoredAt = &now
	alert.IgnoredBy = userID

	if err := c.DB.Save(&alert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "忽略失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": alert, "message": "告警已忽略"})
}

// BatchConfirmAlerts 批量确认告警
func (c *AlertController) BatchConfirmAlerts(ctx *gin.Context) {
	var req struct {
		AlertIDs []uint `json:"alert_ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.AlertIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供告警ID列表"})
		return
	}

	userID := ctx.GetString("user_id")
	now := time.Now()

	result := c.DB.Model(&models.DeviceAlert{}).
		Where("id IN ? AND status = ?", req.AlertIDs, 1).
		Updates(map[string]interface{}{
			"status":       2,
			"confirmed_at": now,
			"confirmed_by": userID,
		})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": fmt.Sprintf("已确认 %d 条告警", result.RowsAffected),
	})
}

// BatchResolveAlerts 批量解决告警
func (c *AlertController) BatchResolveAlerts(ctx *gin.Context) {
	var req struct {
		AlertIDs []uint `json:"alert_ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.AlertIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供告警ID列表"})
		return
	}

	userID := ctx.GetString("user_id")
	now := time.Now()

	result := c.DB.Model(&models.DeviceAlert{}).
		Where("id IN ? AND status NOT IN ?", req.AlertIDs, []int{3, 4}).
		Updates(map[string]interface{}{
			"status":      3,
			"resolved_at": now,
			"resolved_by": userID,
		})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": fmt.Sprintf("已解决 %d 条告警", result.RowsAffected),
	})
}

// GetAlertNotifications 获取告警的通知记录
func (c *AlertController) GetAlertNotifications(ctx *gin.Context) {
	alertID := ctx.Param("id")

	var notifications []models.AlertNotification
	c.DB.Where("alert_id = ?", alertID).Order("id DESC").Find(&notifications)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": notifications},
	})
}

// ==================== 地理围栏 API ====================

// GetGeofenceRules 获取地理围栏规则列表
func (c *AlertController) GetGeofenceRules(ctx *gin.Context) {
	var rules []models.GeofenceRule
	c.DB.Order("id DESC").Find(&rules)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": rules}})
}

// CreateGeofenceRule 创建地理围栏规则
func (c *AlertController) CreateGeofenceRule(ctx *gin.Context) {
	var rule models.GeofenceRule
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

// UpdateGeofenceRule 更新地理围栏规则
func (c *AlertController) UpdateGeofenceRule(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.GeofenceRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	var input models.GeofenceRule
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	rule.Name = input.Name
	rule.DeviceID = input.DeviceID
	rule.CenterLat = input.CenterLat
	rule.CenterLng = input.CenterLng
	rule.RadiusMeters = input.RadiusMeters
	rule.AlertOn = input.AlertOn
	rule.Severity = input.Severity
	rule.Enabled = input.Enabled
	rule.NotifyWays = input.NotifyWays
	rule.Remark = input.Remark
	if err := c.DB.Save(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// DeleteGeofenceRule 删除地理围栏规则
func (c *AlertController) DeleteGeofenceRule(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.GeofenceRule{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetGeofenceAlerts 获取地理围栏告警记录
func (c *AlertController) GetGeofenceAlerts(ctx *gin.Context) {
	query := c.DB.Model(&models.GeofenceAlert{})

	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if alertType := ctx.Query("alert_type"); alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}

	var alerts []models.GeofenceAlert
	query.Order("id DESC").Limit(100).Find(&alerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": alerts},
	})
}

// ==================== 内部方法 ====================

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
		case "jailbreak_detected":
			if v, ok := data["is_jailbroken"].(bool); ok && v {
				triggerVal = 1
				triggered = true
			}
			if v, ok := data["root_status"].(string); ok && (v == "rooted" || v == "jailbroken") {
				triggerVal = 1
				triggered = true
			}
		}

		if triggered {
			extraData, _ := json.Marshal(data)
			alert := models.DeviceAlert{
				RuleID:     rule.ID,
				DeviceID:   deviceID,
				AlertType:  rule.AlertType,
				Severity:   rule.Severity,
				Message:    rule.Name,
				TriggerVal: triggerVal,
				Threshold:  rule.Threshold,
				Status:     1,
				ExtraData:  string(extraData),
			}
			db.Create(&alert)

			// 触发通知
			SendAlertNotifications(db, &alert, rule.NotifyWays)
		}
	}
}

// SendAlertNotifications 发送告警通知（stub实现）
func SendAlertNotifications(db *gorm.DB, alert *models.DeviceAlert, notifyWays string) {
	// TODO: implement alert notification sending
	_ = notifyWays
	_ = alert
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

// CheckCompliance 检查设备合规性
func CheckCompliance(db *gorm.DB, deviceID string, data map[string]interface{}) {
	// 查询所有启用的合规策略
	var policies []models.CompliancePolicy
	db.Where("enabled = ?", true).Find(&policies)

	for _, policy := range policies {
		violated := false
		var actualValue string

		switch policy.PolicyType {
		case "battery_level":
			if v, ok := data["battery"].(float64); ok {
				actualValue = fmt.Sprintf("%.0f", v)
				threshold := parseThreshold(policy.TargetValue)
				violated = evaluateCondition(policy.Condition, v, threshold)
			}

		case "offline_duration":
			if v, ok := data["elapsed"].(float64); ok {
				actualValue = fmt.Sprintf("%.0f秒", v)
				threshold := parseThreshold(policy.TargetValue)
				violated = evaluateCondition(policy.Condition, v, threshold)
			}

		case "is_online":
			if v, ok := data["is_online"].(bool); ok {
				actualValue = fmt.Sprintf("%t", v)
				// 只有离线才算违规
				if !v && policy.TargetValue == "false" && policy.Condition == "=" {
					violated = true
				}
			}
		}

		if violated {
			log.Printf("[Compliance] 设备 %s 违反策略 %s: 类型=%s, 期望=%s, 实际=%s",
				deviceID, policy.Name, policy.PolicyType, policy.TargetValue, actualValue)

			// 创建违规记录
			violation := models.ComplianceViolation{
				PolicyID:      policy.ID,
				DeviceID:     deviceID,
				PolicyType:   policy.PolicyType,
				ExpectedValue: policy.TargetValue,
				ActualValue:   actualValue,
				Severity:      policy.Severity,
				Status:        1,
			}
			db.Create(&violation)

			// 执行补救措施
			executeRemediation(db, deviceID, policy, violation.ID)
		}
	}
}

// parseThreshold 解析阈值字符串为 float64
func parseThreshold(s string) float64 {
	var val float64
	fmt.Sscanf(s, "%f", &val)
	return val
}

// executeRemediation 执行补救措施
func executeRemediation(db *gorm.DB, deviceID string, policy models.CompliancePolicy, violationID uint) {
	log.Printf("[Compliance] 对设备 %s 执行补救措施: %s", deviceID, policy.RemediationAction)

	switch policy.RemediationAction {
	case "notify":
		// 发送通知告警
		alert := models.DeviceAlert{
			RuleID:     0, // 合规触发，无关联规则
			DeviceID:   deviceID,
			AlertType:  "compliance_violation",
			Severity:   policy.Severity,
			Message:    fmt.Sprintf("合规违规: %s", policy.Name),
			TriggerVal: 0,
			Threshold:  0,
			Status:     1,
		}
		db.Create(&alert)

	case "isolate":
		// 隔离设备 - 创建高优先级告警
		alert := models.DeviceAlert{
			RuleID:     0,
			DeviceID:   deviceID,
			AlertType:  "device_isolated",
			Severity:   4, // 严重
			Message:    fmt.Sprintf("设备已被隔离: %s", policy.Name),
			TriggerVal: 0,
			Threshold:  0,
			Status:     1,
		}
		db.Create(&alert)

		// 更新设备状态
		db.Model(&models.DeviceShadow{}).Where("device_id = ?", deviceID).Update("current_mode", "isolated")

	case "block":
		// 阻止设备 - 创建告警
		alert := models.DeviceAlert{
			RuleID:     0,
			DeviceID:   deviceID,
			AlertType:  "device_blocked",
			Severity:   4,
			Message:    fmt.Sprintf("设备已被阻止: %s", policy.Name),
			TriggerVal: 0,
			Threshold:  0,
			Status:     1,
		}
		db.Create(&alert)

	case "wipe":
		// 擦除设备 - 创建紧急告警（实际擦除需要设备端配合）
		alert := models.DeviceAlert{
			RuleID:     0,
			DeviceID:   deviceID,
			AlertType:  "device_wipe_required",
			Severity:   4,
			Message:    fmt.Sprintf("设备需要擦除: %s", policy.Name),
			TriggerVal: 0,
			Threshold:  0,
			Status:     1,
		}
		db.Create(&alert)

		// 记录补救措施已执行
		db.Model(&models.ComplianceViolation{}).Where("id = ?", violationID).Updates(map[string]interface{}{
			"action_taken": policy.RemediationAction,
			"status":       2, // 处理中
		})
	}
}

// DashboardStats 统计数据
type DashboardStats struct {
	TotalDevices    int64 `json:"total_devices"`
	OnlineDevices   int64 `json:"online_devices"`
	OfflineDevices   int64 `json:"offline_devices"`
	TotalAlerts     int64 `json:"total_alerts"`
	PendingAlerts   int64 `json:"pending_alerts"`
	GeofenceAlerts   int64 `json:"geofence_alerts"`
	JailbreakAlerts  int64 `json:"jailbreak_alerts"`
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
	c.DB.Model(&models.DeviceAlert{}).Where("alert_type = ?", "geofence_violation").Count(&stats.GeofenceAlerts)
	c.DB.Model(&models.DeviceAlert{}).Where("alert_type = ?", "jailbreak_detected").Count(&stats.JailbreakAlerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}
