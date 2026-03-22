package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== 订阅计划 =====

// SubscriptionPlan 订阅计划
type SubscriptionPlan struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PlanID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"plan_id"`
	PlanName     string         `gorm:"type:varchar(64);not null" json:"plan_name"`
	PlanType     string         `gorm:"type:varchar(32);not null" json:"plan_type"` // free/basic/pro/enterprise
	Price        float64        `gorm:"type:decimal(10,2);default:0" json:"price"` // 月费
	Currency     string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	DurationDays int            `gorm:"default:30" json:"duration_days"` // 订阅周期（天）
	Features     JSON           `gorm:"type:jsonb" json:"features"`     // 功能列表
	Quotas       JSON           `gorm:"type:jsonb" json:"quotas"`        // 配额限制 {"devices": 10, "users": 5, "storage": 100}
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive
	SortOrder    int            `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

// ===== 用户订阅 =====

// UserSubscription 用户订阅
type UserSubscription struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SubID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"sub_id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	PlanID     string         `gorm:"index;type:varchar(64);not null" json:"plan_id"`
	Status     string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/expired/cancelled/pending
	StartTime  time.Time      `json:"start_time"`
	ExpireTime time.Time      `json:"expire_time"`
	AutoRenew  bool           `gorm:"default:true" json:"auto_renew"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

// ===== 订阅变更记录 =====

// SubscriptionChange 订阅变更记录
type SubscriptionChange struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ChangeID     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"change_id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	SubID        string         `gorm:"index;type:varchar(64)" json:"sub_id"`
	ChangeType   string         `gorm:"type:varchar(20);not null" json:"change_type"` // upgrade/downgrade/renew/cancel/create
	FromPlanID   string         `gorm:"type:varchar(64)" json:"from_plan_id"`
	ToPlanID     string         `gorm:"type:varchar(64);not null" json:"to_plan_id"`
	Amount       float64        `gorm:"type:decimal(10,2);default:0" json:"amount"`
	ChangeReason string         `gorm:"type:text" json:"change_reason"`
	EffectiveAt  time.Time      `json:"effective_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SubscriptionChange) TableName() string {
	return "subscription_changes"
}
