package models

import (
	"time"

	"gorm.io/gorm"
)

type SecuritySession struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Username     string         `gorm:"type:varchar(100);not null" json:"username"`
	Token        string         `gorm:"type:varchar(500);uniqueIndex;not null" json:"token"`
	IP           string         `gorm:"type:varchar(45)" json:"ip"`
	UserAgent    string         `gorm:"type:varchar(500)" json:"user_agent"`
	DeviceType   string         `gorm:"type:varchar(50)" json:"device_type"`
	Location     string         `gorm:"type:varchar(255)" json:"location"`
	Status       int            `gorm:"default:1" json:"status"`
	ExpiresAt    time.Time      `json:"expires_at"`
	LastActiveAt time.Time      `json:"last_active_at"`
	TenantID     string         `gorm:"index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (SecuritySession) TableName() string {
	return "security_sessions"
}

type TwoFactorAuth struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	Secret          string     `gorm:"type:varchar(255);not null" json:"-"`
	SecretEncrypted string     `gorm:"type:text" json:"-"`
	IsEnabled       bool       `gorm:"default:false" json:"is_enabled"`
	IsVerified      bool       `gorm:"default:false" json:"is_verified"`
	RecoveryCodes   string     `gorm:"type:text" json:"recovery_codes"`
	EnabledAt       *time.Time `json:"enabled_at"`
	DisabledAt      *time.Time `json:"disabled_at"`
	LastUsedAt      *time.Time `json:"last_used_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (TwoFactorAuth) TableName() string {
	return "two_factor_auths"
}

type LoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"type:varchar(100);index" json:"username"`
	IP        string    `gorm:"type:varchar(45);index" json:"ip"`
	Action    string    `gorm:"type:varchar(50)" json:"action"`
	Status    int       `gorm:"default:1" json:"status"`
	Reason    string    `gorm:"type:varchar(255)" json:"reason"`
	UserAgent string    `gorm:"type:varchar(500)" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}
