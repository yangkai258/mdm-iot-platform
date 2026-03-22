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
	ID           uint       `gorm:"primaryKey" json:"id"`
	RuleID       uint       `json:"rule_id"`
	DeviceID     string     `gorm:"type:varchar(36);index" json:"device_id"`
	AlertType    string     `gorm:"type:varchar(50)" json:"alert_type"`
	Severity     int        `json:"severity"`
	Message      string     `gorm:"type:varchar(500)" json:"message"`
	TriggerVal   float64    `json:"trigger_val"`
	Threshold    float64    `json:"threshold"`
	Status       int        `gorm:"default:1" json:"status"` // 1:未处理 2:已确认 3:已解决 4:已忽略
	ConfirmedAt  *time.Time `json:"confirmed_at"`
	ConfirmedBy  string     `gorm:"type:varchar(36)" json:"confirmed_by"`
	ResolvedAt   *time.Time `json:"resolved_at"`
	ResolvedBy   string     `gorm:"type:varchar(36)" json:"resolved_by"`
	IgnoredAt    *time.Time `json:"ignored_at"`
	IgnoredBy    string     `gorm:"type:varchar(36)" json:"ignored_by"`
	ExtraData    string     `gorm:"type:jsonb" json:"extra_data"` // 附加数据（地理位置、越狱详情等）
	OrgID        uint       `gorm:"index" json:"org_id"`          // 组织ID（用于数据权限）
	CreateUserID uint       `gorm:"index" json:"create_user_id"` // 创建人ID（用于数据权限）
	CreatedAt    time.Time  `json:"created_at"`
}

func (DeviceAlert) TableName() string {
	return "device_alerts"
}

// GeofenceRule 地理围栏规则
type GeofenceRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	DeviceID     string    `gorm:"type:varchar(36)" json:"device_id"` // 空表示所有设备
	CenterLat    float64   `gorm:"type:decimal(10,7);not null" json:"center_lat"`
	CenterLng    float64   `gorm:"type:decimal(10,7);not null" json:"center_lng"`
	RadiusMeters float64   `gorm:"type:decimal(10,2);not null" json:"radius_meters"` // 半径（米）
	AlertOn      string    `gorm:"type:varchar(20);default:'enter'" json:"alert_on"`  // enter, exit, both
	Severity     int       `gorm:"default:2" json:"severity"`                         // 1:低 2:中 3:高 4:严重
	Enabled      bool      `gorm:"default:true" json:"enabled"`
	NotifyWays   string    `gorm:"type:varchar(100)" json:"notify_ways"` // email,webhook,inapp
	Remark       string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (GeofenceRule) TableName() string {
	return "geofence_rules"
}

// GeofenceAlert 地理围栏告警记录
type GeofenceAlert struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	RuleID     uint       `gorm:"not null;index" json:"rule_id"`
	DeviceID   string     `gorm:"type:varchar(36);index" json:"device_id"`
	AlertType  string     `gorm:"type:varchar(20)" json:"alert_type"` // enter, exit
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	Severity   int        `json:"severity"`
	Message    string     `gorm:"type:varchar(500)" json:"message"`
	Status     int        `gorm:"default:1" json:"status"` // 1:未处理 2:已确认 3:已解决 4:已忽略
	AlertID    uint       `gorm:"index" json:"alert_id"`   // 关联的 DeviceAlert ID
	CreatedAt  time.Time  `json:"created_at"`
}

func (GeofenceAlert) TableName() string {
	return "geofence_alerts"
}

// AlertNotification 告警通知记录
type AlertNotification struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AlertID    uint           `gorm:"not null;index" json:"alert_id"`
	AlertType  string         `gorm:"type:varchar(50)" json:"alert_type"`  // email, webhook, inapp
	Status     string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, sent, delivered, failed
	Recipient  string         `gorm:"type:varchar(255)" json:"recipient"` // 邮箱地址/webhook URL
	Subject    string         `gorm:"type:varchar(255)" json:"subject"`
	Content    string         `gorm:"type:text" json:"content"`
	ErrorMsg   string         `gorm:"type:varchar(500)" json:"error_msg"`
	SentAt     *time.Time     `json:"sent_at"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (AlertNotification) TableName() string {
	return "alert_notifications"
}
