package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== 账单记录 =====

// BillingRecord 账单记录
type BillingRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	BillID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"bill_id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"` // subscription/upgrade/quota
	Amount      float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency    string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/paid/refunded
	PayMethod   string         `gorm:"type:varchar(32)" json:"pay_method"`                 // alipay/wechat/card
	PayTime     *time.Time     `json:"pay_time"`
	Description string         `gorm:"type:varchar(256)" json:"description"`
	OrderNo     string         `gorm:"type:varchar(64)" json:"order_no"` // 外部订单号
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (BillingRecord) TableName() string {
	return "billing_records"
}

// 账单类型常量
const (
	BillingTypeSubscription = "subscription" // 订阅
	BillingTypeUpgrade      = "upgrade"       // 升级
	BillingTypeQuota        = "quota"         // 配额购买
)

// 账单状态常量
const (
	BillingStatusPending  = "pending"  // 待支付
	BillingStatusPaid     = "paid"     // 已支付
	BillingStatusRefunded  = "refunded" // 已退款
	BillingStatusCancelled = "cancelled" // 已取消
)

// ===== 发票 =====

// Invoice 发票
type Invoice struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	InvoiceID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"invoice_id"`
	BillID       string         `gorm:"index;type:varchar(64)" json:"bill_id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Title        string         `gorm:"type:varchar(128);not null" json:"title"`   // 发票抬头
	TaxNo        string         `gorm:"type:varchar(32)" json:"tax_no"`           // 税号
	Amount       float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	TaxAmount    float64        `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	TaxRate      float64        `gorm:"type:decimal(5,2);default:6" json:"tax_rate"` // 税率
	Status       string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/issued/rejected
	InvoiceType  string         `gorm:"type:varchar(20);default:'normal'" json:"invoice_type"` // normal/special/electronic
	ReceiverEmail string        `gorm:"type:varchar(128)" json:"receiver_email"`
	ReceiverPhone string        `gorm:"type:varchar(32)" json:"receiver_phone"`
	IssuedAt     *time.Time     `json:"issued_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Invoice) TableName() string {
	return "invoices"
}

// 发票状态常量
const (
	InvoiceStatusPending = "pending" // 待开
	InvoiceStatusIssued = "issued"  // 已开
	InvoiceStatusRejected = "rejected" // 已拒绝
)
