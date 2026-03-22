package models

import (
	"time"

	"gorm.io/gorm"
)

// App 应用信息
type App struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	BundleID    string         `gorm:"type:varchar(256);uniqueIndex;not null" json:"bundle_id"`
	Description string         `gorm:"type:text" json:"description"`
	IconURL     string         `gorm:"type:varchar(512)" json:"icon_url"`
	Category    string         `gorm:"type:varchar(64)" json:"category"`
	Developer   string         `gorm:"type:varchar(128)" json:"developer"`
	Status      int            `gorm:"type:smallint;default:1" json:"status"` // 1:启用 0:禁用
	Platform    string         `gorm:"type:varchar(32)" json:"platform"`     // ios / android / windows / mac / multi
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (App) TableName() string {
	return "apps"
}

// AppVersion 应用版本
type AppVersion struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	AppID        uint           `gorm:"not null;index" json:"app_id"`
	Version      string         `gorm:"type:varchar(32);not null" json:"version"`
	BuildNumber  string         `gorm:"type:varchar(64)" json:"build_number"`
	FileSize     int64          `gorm:"default:0" json:"file_size"`
	FileURL      string         `gorm:"type:varchar(512);not null" json:"file_url"`
	FileMD5      string         `gorm:"type:varchar(32)" json:"file_md5"`
	MinOSVersion string         `gorm:"type:varchar(32)" json:"min_os_version"`
	ReleaseNotes string         `gorm:"type:text" json:"release_notes"`
	IsMandatory  bool           `gorm:"default:false" json:"is_mandatory"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AppVersion) TableName() string {
	return "app_versions"
}

// AppDistribution 应用分发任务
type AppDistribution struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	TenantID          string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	Name              string         `gorm:"type:varchar(128);not null" json:"name"`
	AppID             uint           `gorm:"not null;index" json:"app_id"`
	VersionID         uint           `gorm:"not null" json:"version_id"`
	DistributionType  string         `gorm:"type:varchar(16);not null" json:"distribution_type"` // device / user / group
	TargetIDs         string         `gorm:"type:jsonb;default:'[]'" json:"target_ids"`          // JSON array of target ids
	TargetCount       int            `gorm:"default:0" json:"target_count"`
	PendingCount      int            `gorm:"default:0" json:"pending_count"`
	SuccessCount      int            `gorm:"default:0" json:"success_count"`
	FailedCount       int            `gorm:"default:0" json:"failed_count"`
	Status            string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending / running / completed / cancelled / failed
	CancelledBy       string         `gorm:"type:varchar(64)" json:"cancelled_by"`
	CancelledAt       *time.Time     `json:"cancelled_at"`
	CompletedAt       *time.Time     `json:"completed_at"`
	CreatedBy         string         `gorm:"type:varchar(64);not null" json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AppDistribution) TableName() string {
	return "app_distributions"
}

// AppInstallRecord 应用安装记录（用于统计）
type AppInstallRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TenantID        string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	DistributionID   *uint          `gorm:"index" json:"distribution_id"`
	DeviceID        string         `gorm:"type:varchar(36);index" json:"device_id"`
	AppID           uint           `gorm:"not null;index" json:"app_id"`
	VersionID       uint           `gorm:"not null" json:"version_id"`
	Status          string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending / downloading / installed / failed
	InstalledAt     *time.Time     `json:"installed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AppInstallRecord) TableName() string {
	return "app_install_records"
}

// AppLicense 应用许可证
type AppLicense struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	UserID       string         `gorm:"type:varchar(36);index" json:"user_id"`
	AppID        uint           `gorm:"not null;index" json:"app_id"`
	LicenseKey   string         `gorm:"type:varchar(256)" json:"license_key"`
	LicenseType  string         `gorm:"type:varchar(32);default:'trial'" json:"license_type"` // trial / standard / enterprise
	PurchaseAt   *time.Time     `json:"purchase_at"`
	ExpiresAt    *time.Time     `json:"expires_at"`
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active / expired / revoked
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AppLicense) TableName() string {
	return "app_licenses"
}
