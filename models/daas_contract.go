package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== DaaS 租赁合同 =====

// DaaSContract 租赁合同
type DaaSContract struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ContractNo      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"contract_no"` // 合同编号
	TenantID        uint           `gorm:"index;not null" json:"tenant_id"`                         // 租户ID
	UserID          uint           `gorm:"index;not null" json:"user_id"`                             // 用户ID（签约人）
	DeviceID        uint           `gorm:"index;not null" json:"device_id"`                         // 设备ID
	DeviceSN        string         `gorm:"type:varchar(128)" json:"device_sn"`                      // 设备序列号
	DeviceName      string         `gorm:"type:varchar(128)" json:"device_name"`                    // 设备名称
	PlanName        string         `gorm:"type:varchar(64)" json:"plan_name"`                      // 租赁套餐名称
	DailyRate       float64        `gorm:"type:decimal(10,2);not null" json:"daily_rate"`           // 每日租金
	MonthlyRate     float64        `gorm:"type:decimal(10,2);not null" json:"monthly_rate"`         // 月租金（可选）
	DepositAmount   float64        `gorm:"type:decimal(10,2);default:0" json:"deposit_amount"`       // 押金
	ContractPeriod  int            `gorm:"default:0" json:"contract_period"`                       // 合同周期（天），0表示不定期
	StartDate       *time.Time     `json:"start_date"`                                             // 合同开始日期
	EndDate         *time.Time     `json:"end_date"`                                               // 合同结束日期
	Status          string         `gorm:"type:varchar(20);default:'active'" json:"status"`       // active/paused/terminated/expired
	TerminateReason string         `gorm:"type:varchar(256)" json:"terminate_reason"`               // 终止原因
	TerminatedAt    *time.Time     `json:"terminated_at"`                                           // 终止时间
	TerminatedBy    uint           `json:"terminated_by"`                                          // 终止操作人
	BillingCycle    string         `gorm:"type:varchar(20);default:'monthly'" json:"billing_cycle"` // daily/monthly
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DaaSContract) TableName() string {
	return "daas_contracts"
}

// 合同状态常量
const (
	DaaSContractStatusActive     = "active"     // 生效中
	DaaSContractStatusPaused      = "paused"     // 已暂停
	DaaSContractStatusTerminated  = "terminated" // 已终止
	DaaSContractStatusExpired     = "expired"    // 已到期
)

// ===== DaaS 设备租赁记录 =====

// DaaSDeviceRental 设备租赁记录
type DaaSDeviceRental struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RentalNo     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"rental_no"` // 租赁流水号
	ContractID   uint           `gorm:"index;not null" json:"contract_id"`                     // 关联合同ID
	TenantID     uint           `gorm:"index;not null" json:"tenant_id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	DeviceID     uint           `gorm:"index;not null" json:"device_id"`
	DeviceSN     string         `gorm:"type:varchar(128)" json:"device_sn"`
	DeviceName   string         `gorm:"type:varchar(128)" json:"device_name"`
	Action       string         `gorm:"type:varchar(20);not null" json:"action"` // rent/return
	RentTime     *time.Time     `json:"rent_time"`                              // 租用时间
	ReturnTime   *time.Time     `json:"return_time"`                            // 归还时间
	ExpectedReturn *time.Time   `json:"expected_return"`                        // 预计归还时间
	DepositAmount float64       `gorm:"type:decimal(10,2);default:0" json:"deposit_amount"`
	Status       string         `gorm:"type:varchar(20);default:'rented'" json:"status"` // rented/returned/overdue
	Notes        string         `gorm:"type:varchar(512)" json:"notes"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DaaSDeviceRental) TableName() string {
	return "daas_device_rentals"
}

// 租赁动作常量
const (
	DaaSRentalActionRent   = "rent"   // 租用
	DaaSRentalActionReturn = "return" // 归还
)

// 租赁状态常量
const (
	DaaSRentalStatusRented   = "rented"   // 已租用
	DaaSRentalStatusReturned = "returned" // 已归还
	DaaSRentalStatusOverdue  = "overdue"  // 已逾期
)

// ===== DaaS 租赁账单 =====

// DaaSBilling 租赁账单
type DaaSBilling struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	BillNo         string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"bill_no"` // 账单编号
	TenantID       uint           `gorm:"index;not null" json:"tenant_id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	ContractID     uint           `gorm:"index;not null" json:"contract_id"`
	RentalID       uint           `gorm:"index" json:"rental_id"` // 可选关联租赁记录
	BillType       string         `gorm:"type:varchar(20);not null" json:"bill_type"` // rental_fee/deposit/deposit_refund/penalty
	PeriodStart    *time.Time     `json:"period_start"`  // 计费周期开始
	PeriodEnd      *time.Time     `json:"period_end"`    // 计费周期结束
	Days           int            `gorm:"default:0" json:"days"`        // 计费天数
	DailyRate      float64        `gorm:"type:decimal(10,2);default:0" json:"daily_rate"` // 日租金（快照）
	Amount         float64        `gorm:"type:decimal(10,2);not null" json:"amount"`     // 账单金额
	DepositAmount  float64        `gorm:"type:decimal(10,2);default:0" json:"deposit_amount"` // 押金金额
	Currency       string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	Status         string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/paid/overdue/cancelled
	PaidAt         *time.Time     `json:"paid_at"`
	PayMethod      string         `gorm:"type:varchar(32)" json:"pay_method"`
	Description    string         `gorm:"type:varchar(512)" json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DaaSBilling) TableName() string {
	return "daas_billings"
}

// 账单类型常量
const (
	DaaSBillTypeRentalFee    = "rental_fee"    // 租赁费
	DaaSBillTypeDeposit       = "deposit"        // 押金
	DaaSBillTypeDepositRefund = "deposit_refund" // 押金退还
	DaaSBillTypePenalty       = "penalty"        // 违约金
)

// 账单状态常量
const (
	DaaSBillStatusPending  = "pending"  // 待支付
	DaaSBillStatusPaid     = "paid"     // 已支付
	DaaSBillStatusOverdue  = "overdue"  // 已逾期
	DaaSBillStatusCancelled = "cancelled" // 已取消
)
