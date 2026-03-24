package models

import (
	"time"

	"gorm.io/gorm"
)

// DevicePairing 设备配对记录
type DevicePairing struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PairingUUID   string         `gorm:"size:64;uniqueIndex;not null" json:"pairing_uuid"`
	PairingCode   string         `gorm:"size:8;index" json:"pairing_code"` // 配对码（6位数字）
	DeviceID      string         `gorm:"size:36;index" json:"device_id"`
	DeviceName    string         `gorm:"size:128" json:"device_name"`
	DeviceSN      string         `gorm:"size:64" json:"device_sn"`
	FirmwareVer   string         `gorm:"size:32" json:"firmware_ver"`
	HardwareVer   string         `gorm:"size:32" json:"hardware_ver"`
	MACAddress    string         `gorm:"size:32" json:"mac_address"`
	UserID        string         `gorm:"size:64;index" json:"user_id"`
	UserName      string         `gorm:"size:128" json:"user_name"`
	Status        string         `gorm:"size:20;default:'pending';index" json:"status"` // pending, approved, rejected, unbound
	RejectReason  string         `gorm:"type:text" json:"reject_reason"`
	ApproveBy     string         `gorm:"size:64" json:"approve_by"`
	ApproveAt     *time.Time     `json:"approve_at"`
	RejectAt      *time.Time     `json:"reject_at"`
	UnboundAt     *time.Time     `json:"unbound_at"`
	UnboundReason string         `gorm:"type:text" json:"unbound_reason"`
	ExpiresAt     time.Time      `gorm:"not null" json:"expires_at"`
	PairedAt      *time.Time     `json:"paired_at"`
	IPAddress     string         `gorm:"size:45" json:"ip_address"`
	TenantID      string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DevicePairing) TableName() string {
	return "device_pairings"
}

// DeviceOpenClawBinding 设备AI授权绑定
type DeviceOpenClawBinding struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	BindingUUID       string         `gorm:"size:64;uniqueIndex;not null" json:"binding_uuid"`
	DeviceID          string         `gorm:"size:36;uniqueIndex;not null" json:"device_id"`
	UserID            string         `gorm:"size:64;not null;index" json:"user_id"`
	OpenClawVersionID uint           `gorm:"not null" json:"openclaw_version_id"`
	AiModelType       string         `gorm:"size:50;not null" json:"ai_model_type"`
	AuthStatus        string         `gorm:"size:20;default:'pending';index" json:"auth_status"` // pending, authorized, cancelled
	AuthToken         string         `gorm:"size:256" json:"auth_token"`
	ExpiresAt         *time.Time     `json:"expires_at"`
	AuthorizedAt      *time.Time     `json:"authorized_at"`
	CancelledAt       *time.Time     `json:"cancelled_at"`
	CancelReason      string         `gorm:"type:text" json:"cancel_reason"`
	TenantID          string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DeviceOpenClawBinding) TableName() string {
	return "device_openclaw_bindings"
}
