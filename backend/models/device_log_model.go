package models

import (
	"time"

	"gorm.io/gorm"
)

// DeviceLog 设备日志
type DeviceLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DeviceID   uint           `gorm:"index" json:"device_id"`              // 设备ID
	TenantID   string         `gorm:"size:50;index" json:"tenant_id"`      // 租户ID
	LogType    string         `gorm:"size:50;index" json:"log_type"`       // 日志类型: online offline command alert error info
	Action     string         `gorm:"size:100" json:"action"`              // 操作动作
	Details    string         `gorm:"type:text" json:"details"`            // 详细信息
	OperatorID *uint          `gorm:"index" json:"operator_id"`             // 操作人ID
	Operator   string         `gorm:"size:50" json:"operator"`             // 操作人
	IPAddress  string         `gorm:"size:50" json:"ip_address"`           // IP地址
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
