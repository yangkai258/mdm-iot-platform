package models

import (
	"time"

	"gorm.io/gorm"
)

// MemberTagDef 会员标签定义
type MemberTagDef struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TagName     string         `gorm:"size:100" json:"tag_name"`             // 标签名称
	TagColor    string         `gorm:"size:20" json:"tag_color"`              // 标签颜色
	Description string         `gorm:"size:500" json:"description"`          // 描述
	TenantID    string         `gorm:"size:50;index" json:"tenant_id"`      // 租户ID
	Sort        int            `gorm:"default:0" json:"sort"`                // 排序
	Status      int            `gorm:"default:1" json:"status"`              // 状态: 1启用 2禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// MemberTagRelation 会员标签关联
type MemberTagRelation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MemberID  uint      `gorm:"index" json:"member_id"`    // 会员ID
	TagID     uint      `gorm:"index" json:"tag_id"`       // 标签ID
	TenantID  string    `gorm:"size:50;index" json:"tenant_id"` // 租户ID
	CreatedAt time.Time `json:"created_at"`
}
