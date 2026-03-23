package models

import "time"

// Subscription 订阅
type Subscription struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"index" json:"user_id"`
	PlanName       string     `json:"plan_name"`
	PlanType       string     `json:"plan_type"` // free, basic, pro, enterprise
	Price          float64    `json:"price"`
	Duration       int        `json:"duration"` // 天数
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	Status         string     `json:"status"` // active, expired, cancelled, suspended
	AutoRenew      bool       `gorm:"default:true" json:"auto_renew"`
	LastRenewAt    *time.Time `json:"last_renew_at"`
	RenewCount     int        `gorm:"default:0" json:"renew_count"`
	RetryCount     int        `gorm:"default:0" json:"retry_count"`         // 续费重试次数
	RenewFailReason string    `json:"renew_fail_reason"`                    // 续费失败原因
	SuspendedAt    *time.Time `json:"suspended_at"`                         // 暂停服务时间
	ReminderSentAt *time.Time `json:"reminder_sent_at"`                    // 提醒发送时间
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// SubscriptionRenewalLog 续费记录
type SubscriptionRenewalLog struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	SubscriptionID  uint       `gorm:"index" json:"subscription_id"`
	Action          string     `json:"action"` // reminder, renewal_success, renewal_failed, suspended
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"` // pending, success, failed
	FailReason      string     `json:"fail_reason"`
	RetryCount      int        `json:"retry_count"`
	CreatedAt       time.Time  `json:"created_at"`
}
