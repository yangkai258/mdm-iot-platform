package models

import (
	"time"

	"gorm.io/gorm"
)

// MiniClawFirmware MiniClaw 固件包
type MiniClawFirmware struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(128);not null" json:"name"`
	Version       string         `gorm:"type:varchar(32);not null;index" json:"version"`
	HardwareModel string         `gorm:"type:varchar(64);not null;index" json:"hardware_model"`
	FileSize      int64          `gorm:"default:0" json:"file_size"`
	FileURL       string         `gorm:"type:varchar(512);not null" json:"file_url"`
	FileMD5       string         `gorm:"type:varchar(32)" json:"file_md5"`
	ReleaseNotes  string         `gorm:"type:text" json:"release_notes"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedBy     string         `gorm:"type:varchar(64);not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MiniClawFirmware) TableName() string {
	return "miniclaw_firmwares"
}

// MiniClawDeviceFirmware 设备固件关联表
type MiniClawDeviceFirmware struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DeviceID     string    `gorm:"type:varchar(36);not null;uniqueIndex" json:"device_id"`
	FirmwareID   uint      `gorm:"not null" json:"firmware_id"`
	FirmwareVer  string    `gorm:"type:varchar(32);not null" json:"firmware_version"`
	UpdatedBy    string    `gorm:"type:varchar(64);not null" json:"updated_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (MiniClawDeviceFirmware) TableName() string {
	return "miniclaw_device_firmwares"
}
