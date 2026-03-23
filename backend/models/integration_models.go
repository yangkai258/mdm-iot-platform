package models

import (
	"time"
)

// ============ 第三方集成相关模型 ============

// Integration 第三方集成记录
type Integration struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	IntegrationType  string    `gorm:"type:varchar(50);not null" json:"integration_type"` // mi_home/tmall_genie/homekit/google_home/vet/insurance/ecommerce
	Status           string    `gorm:"type:varchar(20);default:'disconnected'" json:"status"` // connected/disconnected/error
	Config           JSONMap   `gorm:"type:jsonb" json:"config"`
	ConnectedAt      *time.Time `json:"connected_at"`
	LastSyncAt       *time.Time `json:"last_sync_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Integration) TableName() string { return "integrations" }

// SmartHomeDevice 智能家居设备
type SmartHomeDevice struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"not null;index" json:"user_id"`
	IntegrationID      uint      `json:"integration_id"`
	Platform           string    `gorm:"type:varchar(50);not null" json:"platform"` // mi_home/tmall_genie/homekit/google_home
	PlatformDeviceID   string    `gorm:"type:varchar(100);not null" json:"platform_device_id"`
	DeviceName         string    `gorm:"type:varchar(255)" json:"device_name"`
	DeviceType         string    `gorm:"type:varchar(50)" json:"device_type"` // light/switch/sensor/camera
	Status             JSONMap   `gorm:"type:jsonb" json:"status"`
	IsOnline           bool      `gorm:"default:true" json:"is_online"`
	LastControlAt      *time.Time `json:"last_control_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (SmartHomeDevice) TableName() string { return "smart_home_devices" }

// SmartHomeTrigger 智能联动规则
type SmartHomeTrigger struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	TriggerName  string    `gorm:"type:varchar(255);not null" json:"trigger_name"`
	TriggerType  string    `gorm:"type:varchar(50);not null" json:"trigger_type"` // pet_action/schedule/sensor/voice
	Condition    JSONMap   `gorm:"type:jsonb;not null" json:"condition"`
	Action       JSONMap   `gorm:"type:jsonb;not null" json:"action"`
	DeviceID     uint      `json:"device_id"`
	IsEnabled    bool      `gorm:"default:true" json:"is_enabled"`
	RunCount     int       `gorm:"default:0" json:"run_count"`
	LastRunAt    *time.Time `json:"last_run_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SmartHomeTrigger) TableName() string { return "smart_home_triggers" }

// PetLostReport 寻宠报告
type PetLostReport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PetID       uint      `gorm:"not null;index" json:"pet_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	ReportType  string    `gorm:"type:varchar(20);not null" json:"report_type"` // lost/found
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	LastSeenLat float64   `gorm:"type:decimal(10,6)" json:"last_seen_lat"`
	LastSeenLng float64   `gorm:"type:decimal(10,6)" json:"last_seen_lng"`
	LastSeenAddr string   `gorm:"type:varchar(500)" json:"last_seen_addr"`
	ContactName string    `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string   `gorm:"type:varchar(50)" json:"contact_phone"`
	Reward      float64   `gorm:"type:decimal(10,2)" json:"reward"`
	Photos      StringArray `gorm:"type:varchar(500)[]" json:"photos"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active/found/closed
	ResolvedAt  *time.Time `json:"resolved_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PetLostReport) TableName() string { return "pet_lost_reports" }

// ThirdPartyMapConfig 第三方地图配置
type ThirdPartyMapConfig struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Provider   string    `gorm:"type:varchar(50);not null" json:"provider"` // google_maps/amap/tencent_map
	APIKey     string    `gorm:"type:varchar(500)" json:"api_key"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	QuotaUsed  int       `gorm:"default:0" json:"quota_used"`
	QuotaLimit int       `gorm:"default:10000" json:"quota_limit"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ThirdPartyMapConfig) TableName() string { return "map_configs" }


