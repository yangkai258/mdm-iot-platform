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

// TenantQuota 租户配额表
type TenantQuota struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"type:uuid;uniqueIndex;not null" json:"tenant_id"`
	UserCount   int       `gorm:"default:0" json:"user_count"`
	DeviceCount int       `gorm:"default:0" json:"device_count"`
	DeptCount   int       `gorm:"default:0" json:"dept_count"`
	StoreCount  int       `gorm:"default:0" json:"store_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (TenantQuota) TableName() string { return "tenant_quotas" }

// Plan 套餐表
type Plan struct {
	ID           uint                   `gorm:"primaryKey" json:"id"`
	PlanName     string                 `gorm:"type:varchar(20);not null" json:"plan_name"`
	PlanCode     string                 `gorm:"type:varchar(20);uniqueIndex;not null" json:"plan_code"`
	PriceMonthly float64                `gorm:"type:decimal(10,2)" json:"price_monthly"`
	PriceYearly  float64                `gorm:"type:decimal(10,2)" json:"price_yearly"`
	UserQuota    int                    `gorm:"default:5" json:"user_quota"`
	DeviceQuota  int                    `gorm:"default:10" json:"device_quota"`
	DeptQuota    int                    `gorm:"default:1" json:"dept_quota"`
	StoreQuota   int                    `gorm:"default:1" json:"store_quota"`
	Features     map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"features"`
	SortOrder    int                    `gorm:"default:0" json:"sort_order"`
	IsActive     bool                   `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (Plan) TableName() string { return "plans" }

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

// GetQuota 获取当前租户的配额
func GetQuota(db *gorm.DB, tenantID string) (*TenantQuota, error) {
	var quota TenantQuota
	err := db.Where("tenant_id = ?", tenantID).First(&quota).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &quota, nil
}

// IncrementQuota 原子增加配额计数
func IncrementQuota(db *gorm.DB, tenantID, quotaType string) error {
	fieldMap := map[string]string{
		"user":   "user_count",
		"device": "device_count",
		"dept":   "dept_count",
		"store":  "store_count",
	}
	field, ok := fieldMap[quotaType]
	if !ok {
		return nil
	}
	return db.Exec(
		"INSERT INTO tenant_quotas (tenant_id, "+field+", updated_at) VALUES (?, 1, NOW()) "+
			"ON CONFLICT (tenant_id) DO UPDATE SET "+field+" = tenant_quotas."+field+" + 1, updated_at = NOW()",
		tenantID,
	).Error
}

// DecrementQuota 原子减少配额计数
func DecrementQuota(db *gorm.DB, tenantID, quotaType string) error {
	fieldMap := map[string]string{
		"user":   "user_count",
		"device": "device_count",
		"dept":   "dept_count",
		"store":  "store_count",
	}
	field, ok := fieldMap[quotaType]
	if !ok {
		return nil
	}
	return db.Exec(
		"UPDATE tenant_quotas SET "+field+" = GREATEST("+field+" - 1, 0), updated_at = NOW() WHERE tenant_id = ?",
		tenantID,
	).Error
}
