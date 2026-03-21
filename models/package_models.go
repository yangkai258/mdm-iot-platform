package models

import (
	"time"

	"gorm.io/gorm"
)

// Package 套餐表（多租户版本）
type Package struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PackageCode  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"package_code"`
	PackageName  string         `gorm:"type:varchar(100);not null" json:"package_name"`
	PlanType     string         `gorm:"type:varchar(20);default:'free'" json:"plan_type"` // free, basic, professional, enterprise
	Description  string         `gorm:"type:varchar(500)" json:"description"`
	PriceMonthly float64        `gorm:"type:decimal(10,2);default:0" json:"price_monthly"`
	PriceYearly  float64        `gorm:"type:decimal(10,2);default:0" json:"price_yearly"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	IsDefault    bool           `gorm:"default:false" json:"is_default"` // 是否为默认套餐
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	Features     StringAnyMap   `gorm:"type:jsonb;default:'{}'" json:"features"`      // 功能特性
	QuotaConfig  StringAnyMap   `gorm:"type:jsonb;default:'{}'" json:"quota_config"`  // 配额配置
	Settings     StringAnyMap   `gorm:"type:jsonb;default:'{}'" json:"settings"`      // 附加设置
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Package) TableName() string { return "packages" }

// PackageQuota 套餐配额表（多租户版本）
type PackageQuota struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"type:uuid;not null;uniqueIndex:uq_pkg_quota_tenant_type" json:"tenant_id"`
	PackageID  uint      `gorm:"not null" json:"package_id"`
	QuotaType  string    `gorm:"type:varchar(50);not null;uniqueIndex:uq_pkg_quota_tenant_type" json:"quota_type"` // user, device, store, dept, ota_deployment, app, notification, alert
	QuotaLimit int       `gorm:"default:0" json:"quota_limit"` // 配额上限，0表示无限制，-1表示无限
	QuotaUsed  int       `gorm:"default:0" json:"quota_used"`  // 当前使用量
	QuotaWarnAt *int     `gorm:"default:80" json:"quota_warn_at"` // 警告阈值（百分比）
	UpdatedAt  time.Time `json:"updated_at"`
}

func (PackageQuota) TableName() string { return "package_quotas" }

// StringAnyMap alias for map[string]interface{} used in gorm jsonb
type StringAnyMap map[string]interface{}
