package models

import (
	"time"

	"gorm.io/gorm"
)

// PlanType 套餐类型
type PlanType string

const (
	PlanFree       PlanType = "free"
	PlanBasic      PlanType = "basic"
	PlanPro        PlanType = "pro"
	PlanEnterprise PlanType = "enterprise"
)

// APIQuota API配额表（按租户+月度统计）
type APIQuota struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_tenant_month" json:"tenant_id"`
	Plan        string    `gorm:"type:varchar(20);default:'free'" json:"plan"` // free/basic/pro/enterprise
	MonthlyCalls int64    `gorm:"default:0" json:"monthly_calls"`               // 月度配额上限
	UsedCalls   int64     `gorm:"default:0" json:"used_calls"`                  // 已使用次数
	ResetDate   time.Time `gorm:"not null;uniqueIndex:idx_tenant_month" json:"reset_date"` // 重置日期（每月1日）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (APIQuota) TableName() string { return "api_quotas" }

// APIUsageLog API调用明细记录
type APIUsageLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"type:uuid;not null;index" json:"tenant_id"`
	AppID     uint      `gorm:"index" json:"app_id"`
	APIKeyID  uint      `gorm:"index" json:"api_key_id"`
	Path      string    `gorm:"type:varchar(255);index" json:"path"`
	Method    string    `gorm:"type:varchar(10)" json:"method"`
	StatusCode int       `json:"status_code"`
	LatencyMs int64     `json:"latency_ms"`
	IP        string    `gorm:"type:varchar:45" json:"ip"`
	UserAgent string    `gorm:"type:varchar(512)" json:"user_agent"`
	CallAt    time.Time `gorm:"not null;index" json:"call_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (APIUsageLog) TableName() string { return "api_usage_logs" }

// Invoice 发票表
type Invoice struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	InvoiceNo     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"invoice_no"` // 发票号
	TenantID      string    `gorm:"type:uuid;not null;index" json:"tenant_id"`
	TenantName    string    `gorm:"type:varchar(200)" json:"tenant_name"`
	Plan          string    `gorm:"type:varchar(20)" json:"plan"` // 订阅套餐
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"` // 金额
	Currency      string    `gorm:"type:varchar(3);default:'CNY'" json:"currency"`
	Status        string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/paid/cancelled/refunded
	BillingPeriod string    `gorm:"type:varchar(20)" json:"billing_period"` // 账单周期 e.g. "2026-03"
	DueDate       time.Time `json:"due_date"`
	PaidAt        *time.Time `json:"paid_at"`
	PaymentMethod string    `gorm:"type:varchar(50)" json:"payment_method"`
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Invoice) TableName() string { return "invoices" }

// ==================== 辅助函数 ====================

// GetOrCreateAPIQuota 获取或创建当前租户的API配额
func GetOrCreateAPIQuota(db *gorm.DB, tenantID string) (*APIQuota, error) {
	now := time.Now()
	resetDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0) // 下月1日

	var quota APIQuota
	err := db.Where("tenant_id = ? AND reset_date = ?", tenantID, resetDate).First(&quota).Error
	if err == nil {
		return &quota, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 查询租户套餐
	var tenant Tenant
	if err := db.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return nil, err
	}

	planCode := tenant.Plan
	if planCode == "" {
		planCode = "free"
	}

	monthlyCalls := getAPIMonthlyCalls(planCode)

	quota = APIQuota{
		TenantID:     tenantID,
		Plan:         planCode,
		MonthlyCalls: monthlyCalls,
		UsedCalls:    0,
		ResetDate:    resetDate,
	}
	if err := db.Create(&quota).Error; err != nil {
		return nil, err
	}
	return &quota, nil
}

// IncrementAPICalls 原子增加API调用次数
func IncrementAPICalls(db *gorm.DB, tenantID string) error {
	now := time.Now()
	resetDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0)

	return db.Exec(
		"INSERT INTO api_quotas (tenant_id, used_calls, reset_date, created_at, updated_at) VALUES (?, 1, ?, NOW(), NOW()) "+
			"ON CONFLICT (tenant_id, reset_date) DO UPDATE SET used_calls = api_quotas.used_calls + 1, updated_at = NOW()",
		tenantID, resetDate,
	).Error
}

// getAPIMonthlyCalls 根据套餐代码返回月度API配额
func getAPIMonthlyCalls(planCode string) int64 {
	switch planCode {
	case "free":
		return 1000
	case "basic":
		return 10000
	case "pro":
		return 100000
	case "enterprise":
		return -1 // 不限
	default:
		return 1000
	}
}

// CheckAPIRateLimit 检查API调用是否超限（速率限制）
func CheckAPIRateLimit(db *gorm.DB, apiKeyID uint) (bool, int64, error) {
	// 这里使用简单的滑动窗口：检查最近1分钟的调用次数
	// 实际生产环境建议使用Redis
	now := time.Now()
	windowStart := now.Add(-1 * time.Minute)

	var count int64
	err := db.Model(&APIUsageLog{}).
		Where("api_key_id = ? AND call_at >= ?", apiKeyID, windowStart).
		Count(&count).Error
	if err != nil {
		return false, 0, err
	}

	// 免费用户每分钟100次，付费用户更高
	limit := int64(100)
	if count < limit {
		return true, limit - count, nil // 允许调用
	}
	return false, 0, nil // 超限
}
