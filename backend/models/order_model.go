package models

import (
	"time"

	"gorm.io/gorm"
)

// Order 会员订单
type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OrderNo      string         `gorm:"size:50;uniqueIndex" json:"order_no"`    // 订单号
	MemberID     uint           `gorm:"index" json:"member_id"`                 // 会员ID
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`        // 租户ID
	TotalAmount  float64        `gorm:"default:0" json:"total_amount"`         // 订单金额
	Status       int            `gorm:"default:1" json:"status"`               // 状态: 1待支付 2已支付 3已完成 4已取消 5已退款
	OrderType    int            `gorm:"default:1" json:"order_type"`           // 订单类型: 1消费 2充值 3退款
	PayType      int            `gorm:"default:1" json:"pay_type"`             // 支付方式: 1微信 2支付宝 3现金 4其他
	PointsEarned int64          `gorm:"default:0" json:"points_earned"`        // 获得积分
	Discount     float64        `gorm:"default:0" json:"discount"`             // 优惠金额
	Remark       string         `gorm:"size:500" json:"remark"`                // 备注
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
