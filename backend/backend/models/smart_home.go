package models

import (
	"time"
)

// SmartHomeDevice 智能家居设备
type SmartHomeDevice struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TenantID        string    `gorm:"type:uuid;index" json:"tenant_id"`
	UserID          uint      `gorm:"index" json:"user_id"`
	IntegrationID   uint      `gorm:"index" json:"integration_id"`
	Platform        string    `gorm:"size:50" json:"platform"` // homekit, google_home, alexa, tuya, mi
	PlatformDeviceID string   `gorm:"size:200" json:"platform_device_id"`
	DeviceName      string    `gorm:"size:200" json:"device_name"`
	DeviceType      string    `gorm:"size:50" json:"device_type"` // light, switch, thermostat, camera, lock, sensor
	RoomName        string    `gorm:"size:100" json:"room_name"`
	IsOnline        bool      `gorm:"default:false" json:"is_online"`
	LastOnlineAt    *time.Time `json:"last_online_at"`
	Status          JSONMap   `gorm:"type:jsonb" json:"status"` // JSON
	Config          JSONMap   `gorm:"type:jsonb" json:"config"`  // JSON
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName
func (SmartHomeDevice) TableName() string {
	return "smart_home_devices"
}
