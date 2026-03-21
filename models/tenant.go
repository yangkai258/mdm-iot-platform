package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户实体（对应 tenants 表）
type Tenant struct {
	ID           string                 `gorm:"type:uuid;primaryKey" json:"id"`
	TenantCode   string                 `gorm:"type:varchar(50);uniqueIndex;not null" json:"tenant_code"`
	Name         string                 `gorm:"type:varchar(200);not null" json:"name"`
	ContactName  string                 `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string                 `gorm:"type:varchar(20)" json:"contact_phone"`
	ContactEmail string                 `gorm:"type:varchar(100)" json:"contact_email"`
	Plan         string                 `gorm:"type:varchar(50);default:'free'" json:"plan"`
	Status       string                 `gorm:"type:varchar(20);default:'pending'" json:"status"`
	LogoURL      string                 `gorm:"type:varchar(500)" json:"logo_url"`
	Domain       string                 `gorm:"type:varchar(200)" json:"domain"`
	ExpiresAt    *time.Time             `json:"expires_at"`
	Settings     map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (Tenant) TableName() string { return "tenants" }

// TenantQuota 租户配额表（对应 package_quotas 表）
type TenantQuota struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"type:uuid;not null" json:"tenant_id"`
	PackageID  uint      `gorm:"not null" json:"package_id"`
	QuotaType  string    `gorm:"type:varchar(50);not null" json:"quota_type"` // user, device, store, dept, ota_deployment, app, notification, alert
	QuotaLimit int       `gorm:"default:0" json:"quota_limit"`                 // 配额上限，0表示无限制(-1)
	QuotaUsed  int       `gorm:"default:0" json:"quota_used"`                  // 当前使用量
	QuotaWarnAt int      `gorm:"default:80" json:"quota_warn_at"`              // 警告阈值（百分比）
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TenantQuota) TableName() string { return "package_quotas" }

// Package 套餐表（对应 packages 表）
// 修复：表名从 plans 改为 packages，字段名从 plan_code 改为 package_code，
// 配额从独立字段改为 JSONB quota_config
type Package struct {
	ID           uint                   `gorm:"primaryKey" json:"id"`
	PackageCode  string                 `gorm:"type:varchar(50);uniqueIndex;not null" json:"package_code"`
	PackageName  string                 `gorm:"type:varchar(100);not null" json:"package_name"`
	PlanType     string                 `gorm:"type:varchar(20);default:'free'" json:"plan_type"`
	Description  string                 `gorm:"type:varchar(500)" json:"description"`
	PriceMonthly float64                `gorm:"type:decimal(10,2)" json:"price_monthly"`
	PriceYearly  float64                `gorm:"type:decimal(10,2)" json:"price_yearly"`
	IsActive     bool                   `gorm:"default:true" json:"is_active"`
	IsDefault    bool                   `gorm:"default:false" json:"is_default"`
	SortOrder    int                    `gorm:"default:0" json:"sort_order"`
	Features     map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"features"`
	QuotaConfig  map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"quota_config"` // JSONB: {"devices":10, "users":5, ...}
	Settings     map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (Package) TableName() string { return "packages" }

// TenantApplication 租户申请记录
type TenantApplication struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ApplicationCode string `gorm:"type:varchar(50);uniqueIndex;not null" json:"application_code"`
	CompanyName  string    `gorm:"type:varchar(200);not null" json:"company_name"`
	ContactName  string    `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string    `gorm:"type:varchar(20)" json:"contact_phone"`
	ContactEmail string    `gorm:"type:varchar(100)" json:"contact_email"`
	Industry     string    `gorm:"type:varchar(50)" json:"industry"`
	CompanySize  string    `gorm:"type:varchar(20)" json:"company_size"`
	PlanID       uint      `json:"plan_id"`
	PlanName     string    `gorm:"type:varchar(50)" json:"plan_name"`
	UseCase      string    `gorm:"type:text" json:"use_case"`
	Status       string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	RejectReason string    `gorm:"type:text" json:"reject_reason"`
	AdminNotes   string    `gorm:"type:text" json:"admin_notes"`
	ApprovedBy   string    `gorm:"type:varchar(64)" json:"approved_by"`
	ApprovedAt   string    `json:"approved_at"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func (TenantApplication) TableName() string { return "tenant_applications" }

// ApprovalHistory 审批历史记录
type ApprovalHistory struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	ApplicationID uint  `gorm:"index;not null" json:"application_id"`
	Action       string `gorm:"type:varchar(20);not null" json:"action"` // approve, reject, withdraw
	ActionText   string `gorm:"type:varchar(200)" json:"action_text"`
	Operator     string `gorm:"type:varchar(64)" json:"operator"`
	Comment      string `gorm:"type:text" json:"comment"`
	CreatedAt    string `json:"created_at"`
}

func (ApprovalHistory) TableName() string { return "approval_histories" }

// ==================== 配额辅助函数 ====================

// GetQuota 获取当前租户的指定类型配额
// 修复：添加 quotaType 参数，因为 package_quotas 表是按类型分行的
func GetQuota(db *gorm.DB, tenantID string, quotaType string) (*TenantQuota, error) {
	var quota TenantQuota
	err := db.Where("tenant_id = ? AND quota_type = ?", tenantID, quotaType).First(&quota).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &quota, nil
}

// IncrementQuota 原子增加配额计数
// 修复：使用 package_quotas 表和 quota_used 字段
func IncrementQuota(db *gorm.DB, tenantID, quotaType string) error {
	return db.Exec(
		"UPDATE package_quotas SET quota_used = quota_used + 1, updated_at = NOW() WHERE tenant_id = ? AND quota_type = ?",
		tenantID, quotaType,
	).Error
}

// DecrementQuota 原子减少配额计数
// 修复：使用 package_quotas 表和 quota_used 字段
func DecrementQuota(db *gorm.DB, tenantID, quotaType string) error {
	return db.Exec(
		"UPDATE package_quotas SET quota_used = GREATEST(quota_used - 1, 0), updated_at = NOW() WHERE tenant_id = ? AND quota_type = ?",
		tenantID, quotaType,
	).Error
	}
	return db.Exec(
		"UPDATE tenant_quotas SET "+field+" = GREATEST("+field+" - 1, 0), updated_at = NOW() WHERE tenant_id = ?",
		tenantID,
	).Error
}
