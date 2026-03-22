package models

import (
	"time"

	"gorm.io/gorm"
)

// DeveloperApp 开发者应用
type DeveloperApp struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"` // 所属用户
	AppName     string         `gorm:"type:varchar(128);not null" json:"app_name"`
	AppKey      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"app_key"` // 应用唯一标识 (public key)
	Description string         `gorm:"type:text" json:"description"`
	WebsiteURL  string         `gorm:"type:varchar(512)" json:"website_url"`
	LogoURL     string         `gorm:"type:varchar(512)" json:"logo_url"`
	Category    string         `gorm:"type:varchar(64)" json:"category"` // 游戏/工具/社交/电商/物联网/其他
	Platform    string         `gorm:"type:varchar(32)" json:"platform"` // ios/android/windows/mac/linux/web/other
	Status      int            `gorm:"type:smallint;default:1" json:"status"` // 1:启用 0:禁用 -1:待审核
	ApprovedAt  *time.Time     `json:"approved_at"`
	ApprovedBy uint           `json:"approved_by"`
	RejectReason string        `gorm:"type:text" json:"reject_reason"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DeveloperApp) TableName() string {
	return "developer_apps"
}

// DeveloperAppStatus 应用状态常量
const (
	DeveloperAppStatusPending  = -1  // 待审核
	DeveloperAppStatusEnabled = 1   // 启用
	DeveloperAppStatusDisabled = 0   // 禁用
)

// DeveloperAppPlatform 应用平台常量
const (
	DeveloperAppPlatformiOS     = "ios"
	DeveloperAppPlatformAndroid = "android"
	DeveloperAppPlatformWeb     = "web"
	DeveloperAppPlatformOther   = "other"
)
