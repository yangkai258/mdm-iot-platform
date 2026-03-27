package models

import "time"

// PaymentMethod 支付方式
type PaymentMethod struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	MethodType string    `json:"method_type"` // credit_card, alipay, wechat
	CardBrand  string    `json:"card_brand"`  // visa, mastercard, amex (for credit_card)
	Last4      string    `json:"last4"`       // 信用卡后4位
	ExpireMonth int      `json:"expire_month"`
	ExpireYear  int      `json:"expire_year"`
	IsDefault  bool      `gorm:"default:false" json:"is_default"`
	ExtData    string    `json:"ext_data"`    // 第三方支付额外信息(JSON)
	Status     string    `json:"status"`      // active, expired, removed
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// SubscriptionRenewal 续费记录
type SubscriptionRenewal struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SubscriptionID uint      `gorm:"index" json:"subscription_id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	PaymentMethodID *uint    `gorm:"index" json:"payment_method_id"`
	Amount         float64   `json:"amount"`
	PaymentMethod  string    `json:"payment_method"` // credit_card, alipay, wechat
	PaymentStatus  string    `json:"payment_status"` // pending, processing, success, failed, refunded
	PaymentTransID string    `json:"payment_trans_id"` // 第三方交易流水号
	FailReason     string    `json:"fail_reason"`
	RetryCount     int       `gorm:"default:0" json:"retry_count"`
	RenewalDate    time.Time `json:"renewal_date"`
	ExpiredDate    time.Time `json:"expired_date"`    // 续费后的到期日期
	NextRetryAt    *time.Time `json:"next_retry_at"` // 下次重试时间
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	Subscription  Subscription   `gorm:"foreignKey:SubscriptionID" json:"-"`
	PaymentMethodRef *PaymentMethod `gorm:"foreignKey:PaymentMethodID" json:"-"`
}

// AutoRenewalSetting 自动续费设置
type AutoRenewalSetting struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex" json:"user_id"`
	SubscriptionID uint      `gorm:"index" json:"subscription_id"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	PaymentMethodID *uint    `gorm:"index" json:"payment_method_id"`
	ReminderDays   string    `json:"reminder_days"` // 逗号分隔: "7,3,1"
	LastRemindAt   *time.Time `json:"last_remind_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
