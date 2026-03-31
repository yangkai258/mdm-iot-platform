package models

import (
	"time"

	"gorm.io/gorm"
)

// MemberService 会员服务记录
type MemberService struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	MemberID    uint           `gorm:"index" json:"member_id"`                // 会员ID
	TenantID    string         `gorm:"size:50;index" json:"tenant_id"`         // 租户ID
	ServiceType int            `gorm:"default:1" json:"service_type"`          // 服务类型: 1咨询 2体验 3售后 4其他
	Content     string         `gorm:"type:text" json:"content"`              // 服务内容
	OperatorID  *uint          `gorm:"index" json:"operator_id"`              // 操作人ID
	Operator    string         `gorm:"size:50" json:"operator"`               // 操作人姓名
	Rating      int            `gorm:"default:0" json:"rating"`               // 评分 1-5
	Status      int            `gorm:"default:1" json:"status"`               // 状态: 1进行中 2已完成 3已取消
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
