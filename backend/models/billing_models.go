package models

import "time"

// Invoice 发票
type Invoice struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	InvoiceNo     string     `gorm:"uniqueIndex;not null" json:"invoice_no"` // 发票号
	UserID        uint       `gorm:"index" json:"user_id"`
	SubscriptionID uint       `gorm:"index" json:"subscription_id"`
	OrderID       uint       `gorm:"index" json:"order_id"`
	Amount        float64    `json:"amount"`      // 发票金额
	TaxAmount     float64    `json:"tax_amount"`  // 税额
	TotalAmount   float64    `json:"total_amount"` // 总金额
	Status        string     `gorm:"default:'pending'" json:"status"` // pending, issued, void
	InvoiceType   string     `json:"invoice_type"`                   // normal, VAT
	Title         string     `json:"title"`                          // 发票抬头
	TaxNo         string     `json:"tax_no"`                         // 税号
	BankName      string     `json:"bank_name"`                      // 开户行
	BankAccount   string     `json:"bank_account"`                   // 银行账号
	Address       string     `json:"address"`                        // 地址电话
	Email         string     `json:"email"`                          // 接收邮箱
	IssueDate     *time.Time `json:"issue_date"`                     // 开票日期
	DueDate       *time.Time `json:"due_date"`                      // 到期日期
	PDFURL        string     `json:"pdf_url"`                        // PDF链接
	CreatedAt     time.Time  `json:"created_at"`
}

// BillingRecord 计费记录
type BillingRecord struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	BillingType string     `json:"billing_type"`  // subscription, usage, API_quota
	ReferenceID uint       `json:"reference_id"`  // 关联订单/订阅ID
	Amount      float64    `json:"amount"`
	Currency    string     `gorm:"default:'CNY'" json:"currency"`
	PeriodStart time.Time  `json:"period_start"`
	PeriodEnd   time.Time  `json:"period_end"`
	Status      string     `gorm:"default:'pending'" json:"status"` // pending, paid, overdue
	PaidAt      *time.Time `json:"paid_at"`
	InvoiceID   *uint      `json:"invoice_id"`
	CreatedAt   time.Time  `json:"created_at"`
}
