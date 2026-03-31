package models

import (
	"time"

	"gorm.io/gorm"
)

// MemberCardRecord 会员卡实例（已发放的会员卡）
type MemberCardRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CardNumber  string         `gorm:"size:50;uniqueIndex" json:"card_number"` // 卡号
	CardType    int            `gorm:"default:1" json:"card_type"`             // 卡类型: 1储值卡 2积分卡 3打折卡
	MemberID    *uint          `gorm:"index" json:"member_id"`                 // 会员ID
	TenantID    string         `gorm:"size:50;index" json:"tenant_id"`         // 租户ID
	Status      int            `gorm:"default:1" json:"status"`               // 状态: 1正常 2冻结 3作废
	IssuedAt    time.Time      `json:"issued_at"`                             // 发卡时间
	ExpiredAt   *time.Time     `json:"expired_at"`                            // 过期时间
	Balance     float64        `gorm:"default:0" json:"balance"`               // 卡内余额
	Points      int64          `gorm:"default:0" json:"points"`               // 卡内积分
	Remark      string         `gorm:"size:500" json:"remark"`                // 备注
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
