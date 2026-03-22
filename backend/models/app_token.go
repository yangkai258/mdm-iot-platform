package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppToken App应用Token表（用于App扫码登录等场景）
type AppToken struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Token        string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"token"`
	AppID        string     `gorm:"type:varchar(64);index" json:"app_id"`
	UserID       uint       `gorm:"index" json:"user_id"`
	DeviceID     string     `gorm:"type:varchar(36);index" json:"device_id"` // 关联设备ID（可选）
	Platform     string     `gorm:"type:varchar(20)" json:"platform"`        // ios, android, miniapp
	ClientID     string     `gorm:"type:varchar(128)" json:"client_id"`      // 客户端唯一标识
	Scope        string     `gorm:"type:varchar(255)" json:"scope"`          // 权限范围，逗号分隔
	ExpiresAt    time.Time  `json:"expires_at"`                              // 过期时间
	RefreshCount int        `gorm:"default:0" json:"refresh_count"`          // 刷新次数
	RevokedAt    *time.Time `json:"revoked_at"`                              // 撤销时间（nil表示有效）
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (AppToken) TableName() string {
	return "app_tokens"
}

// BeforeCreate 创建前自动生成 Token
func (a *AppToken) BeforeCreate(tx *gorm.DB) error {
	if a.Token == "" {
		a.Token = uuid.New().String()
	}
	return nil
}

// IsExpired 检查Token是否过期
func (a *AppToken) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

// IsValid 检查Token是否有效（未过期且未撤销）
func (a *AppToken) IsValid() bool {
	return a.IsActive && !a.IsExpired() && a.RevokedAt == nil
}

// AppRefreshToken App刷新Token表
type AppRefreshToken struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	RefreshToken string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"refresh_token"`
	AppTokenID   uint       `gorm:"index" json:"app_token_id"`
	UserID       uint       `gorm:"index" json:"user_id"`
	ClientID     string     `gorm:"type:varchar(128);index" json:"client_id"`
	Platform     string     `gorm:"type:varchar(20)" json:"platform"`
	ExpiresAt    time.Time  `json:"expires_at"` // 刷新Token有效期（通常7天）
	RevokedAt    *time.Time `json:"revoked_at"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (AppRefreshToken) TableName() string {
	return "app_refresh_tokens"
}

// BeforeCreate 创建前自动生成刷新Token
func (a *AppRefreshToken) BeforeCreate(tx *gorm.DB) error {
	if a.RefreshToken == "" {
		a.RefreshToken = uuid.New().String() + uuid.New().String()
	}
	return nil
}

// IsExpired 检查刷新Token是否过期
func (a *AppRefreshToken) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

// IsValid 检查刷新Token是否有效
func (a *AppRefreshToken) IsValid() bool {
	return a.IsActive && !a.IsExpired() && a.RevokedAt == nil
}
