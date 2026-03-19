package models

import (
	"time"
)

// DeviceAlertRule 设备告警规则
type DeviceAlertRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	DeviceID    string    `gorm:"type:varchar(36)" json:"device_id"` // 空表示所有设备
	AlertType   string    `gorm:"type:varchar(50);not null" json:"alert_type"` // battery_low, offline, temperature_high
	Condition   string    `gorm:"type:varchar(100);not null" json:"condition"` // <, >, =, >=
	Threshold   float64   `not null" json:"threshold"`
	Severity    int       `gorm:"default:1" json:"severity"` // 1:低 2:中 3:高 4:严重
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	NotifyWays  string    `gorm:"type:varchar(100)" json:"notify_ways"` // email,sms,webhook
	Remark      string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DeviceAlertRule) TableName() string {
	return "device_alert_rules"
}

// DeviceAlert 设备告警记录
type DeviceAlert struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RuleID     uint      `json:"rule_id"`
	DeviceID   string    `gorm:"type:varchar(36);index" json:"device_id"`
	AlertType  string    `gorm:"type:varchar(50)" json:"alert_type"`
	Severity   int       `json:"severity"`
	Message    string    `gorm:"type:varchar(500)" json:"message"`
	TriggerVal float64   `json:"trigger_val"`
	Threshold  float64   `json:"threshold"`
	Status     int       `gorm:"default:1" json:"status"` // 1:未处理 2:已确认 3:已解决
	CreatedAt  time.Time `json:"created_at"`
}

func (DeviceAlert) TableName() string {
	return "device_alerts"
}
