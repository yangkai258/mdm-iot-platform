package services

import (
	"encoding/json"
	"strings"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// AlertDeduplicationService 告警去重服务
type AlertDeduplicationService struct {
	DB *gorm.DB
}

// NewAlertDeduplicationService 创建服务
func NewAlertDeduplicationService(db *gorm.DB) *AlertDeduplicationService {
	return &AlertDeduplicationService{DB: db}
}

// ShouldProcessAlert 判断告警是否应该处理（不在去重窗口内）
func (s *AlertDeduplicationService) ShouldProcessAlert(alertType, deviceID string, severity int) (bool, *models.AlertDeduplicationRule) {
	var rule models.AlertDeduplicationRule
	err := s.DB.Where("is_active = ? AND alert_type = ?", true, alertType).
		Where("severity_min <= ? AND severity_max >= ?", severity, severity).
		First(&rule).Error

	if err != nil {
		return true, nil
	}

	if rule.DevicePattern != "" && !matchDevicePattern(deviceID, rule.DevicePattern) {
		return true, nil
	}

	windowStart := time.Now().Add(-time.Duration(rule.DedupWindowSeconds) * time.Second)

	var recentRecord models.AlertDeduplicationRecord
	err = s.DB.Where("rule_id = ? AND device_id = ? AND window_start >= ?",
		rule.RuleID, deviceID, windowStart).
		Order("window_start DESC").First(&recentRecord).Error

	if err == gorm.ErrRecordNotFound {
		return true, &rule
	}

	if recentRecord.AlertCount >= rule.MaxCountPerWindow {
		return false, &rule
	}

	return true, &rule
}

// RecordAlert 记录告警到去重窗口
func (s *AlertDeduplicationService) RecordAlert(ruleID, alertType, deviceID, alertID string, alertData interface{}) error {
	windowEnd := time.Now().Add(24 * time.Hour)

	var rule models.AlertDeduplicationRule
	if err := s.DB.Where("rule_id = ?", ruleID).First(&rule).Error; err == nil {
		windowEnd = time.Now().Add(time.Duration(rule.DedupWindowSeconds) * time.Second)
	}

	alertJSON, _ := json.Marshal(alertData)

	record := models.AlertDeduplicationRecord{
		RuleID:             ruleID,
		AlertType:          alertType,
		DeviceID:           deviceID,
		WindowStart:        time.Now(),
		WindowEnd:          windowEnd,
		AlertCount:         1,
		AlertSummary:       string(alertJSON),
		FinalAlertID:       alertID,
		FinalAlertSnapshot:  string(alertJSON),
	}

	return s.DB.Create(&record).Error
}

// IncrementAlertCount 增加窗口内的告警计数
func (s *AlertDeduplicationService) IncrementAlertCount(ruleID, deviceID string) error {
	windowStart := time.Now().Add(-time.Hour)

	var rule models.AlertDeduplicationRule
	if err := s.DB.Where("rule_id = ?", ruleID).First(&rule).Error; err == nil {
		windowStart = time.Now().Add(-time.Duration(rule.DedupWindowSeconds) * time.Second)
	}

	return s.DB.Model(&models.AlertDeduplicationRecord{}).
		Where("rule_id = ? AND device_id = ? AND window_start >= ?", ruleID, deviceID, windowStart).
		UpdateColumn("alert_count", gorm.Expr("alert_count + 1")).Error
}

// GetDeduplicationStats 获取去重统计
func (s *AlertDeduplicationService) GetDeduplicationStats(startTime, endTime time.Time) (int64, int64) {
	var total int64
	s.DB.Model(&models.AlertDeduplicationRecord{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Select("COALESCE(SUM(alert_count), 0)").
		Scan(&total)

	var dedup int64
	s.DB.Model(&models.AlertDeduplicationRecord{}).
		Where("created_at BETWEEN ? AND ? AND alert_count > 1", startTime, endTime).
		Count(&dedup)

	return total, dedup
}

// matchDevicePattern 匹配设备ID模式
func matchDevicePattern(deviceID, pattern string) bool {
	if pattern == "" {
		return true
	}
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(deviceID, prefix)
	}
	return deviceID == pattern
}
