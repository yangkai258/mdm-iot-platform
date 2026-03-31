package models

import (
	"time"
)

// PairingRecord 配对记录
type PairingRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PairingCode   string    `gorm:"type:varchar(8);uniqueIndex;not null" json:"pairing_code"`
	DeviceID      string    `gorm:"type:varchar(64);index" json:"device_id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	Status        int       `gorm:"type:smallint;default:1" json:"status"` // 1=待激活 2=已配对 3=已解绑
	ExpiresAt     time.Time `gorm:"not null" json:"expires_at"`
	PairedAt      *time.Time `json:"paired_at"`
	UnboundAt     *time.Time `json:"unbound_at"`
	UnboundReason string    `gorm:"type:varchar(256)" json:"unbound_reason"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (PairingRecord) TableName() string {
	return "pairing_records"
}

// PairingStatus 配对状态常量
const (
	PairingStatusPending   = 1 // 待激活
	PairingStatusPaired   = 2 // 已配对
	PairingStatusUnbound  = 3 // 已解绑
)

// DeviceOpenClawBinding 设备AI授权绑定
type DeviceOpenClawBinding struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	DeviceID         string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"device_id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	OpenClawVersionID uint    `gorm:"not null" json:"openclaw_version_id"`
	AIModelType      string    `gorm:"type:varchar(64);not null" json:"ai_model_type"`
	AuthStatus       int       `gorm:"type:smallint;default:1" json:"auth_status"` // 1=待授权 2=已授权 3=已取消
	AuthToken        string    `gorm:"type:varchar(256)" json:"auth_token"`
	ExpiresAt        *time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName 指定表名
func (DeviceOpenClawBinding) TableName() string {
	return "device_openclaw_bindings"
}

// AuthStatus 授权状态常量
const (
	AuthStatusPending = 1 // 待授权
	AuthStatusOK     = 2 // 已授权
	AuthStatusCancel = 3 // 已取消
)

// PairingCodeRequest 生成配对码请求
type PairingCodeRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// PairingVerifyRequest 配对验证请求
type PairingVerifyRequest struct {
	PairingCode string         `json:"pairing_code" binding:"required"`
	DeviceID    string         `json:"device_id" binding:"required"`
	DeviceInfo  PairingDeviceInfo `json:"device_info"`
}

// PairingDeviceInfo 设备信息
type PairingDeviceInfo struct {
	FirmwareVersion string `json:"firmware_version"`
	HardwareVersion string `json:"hardware_version"`
	MACAddress      string `json:"mac_address"`
}

// PairingHistoryQuery 配对历史查询
type PairingHistoryQuery struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size"`
	DeviceID string `form:"device_id"`
	Status   int    `form:"status"`
}
