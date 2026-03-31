package models

import (
	"time"

	"gorm.io/gorm"
)

// TempMemberRecord 临时会员记录
type TempMemberRecord struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OpenID     string         `gorm:"size:100;index" json:"open_id"`      // 微信OpenID
	Nickname   string         `gorm:"size:100" json:"nickname"`           // 昵称
	Avatar     string         `gorm:"size:500" json:"avatar"`             // 头像
	Phone      string         `gorm:"size:20" json:"phone"`               // 手机号
	TenantID   string         `gorm:"size:50;index" json:"tenant_id"`    // 租户ID
	StoreID    *uint          `gorm:"index" json:"store_id"`              // 店铺ID
	ExpireTime time.Time      `json:"expire_time"`                        // 过期时间
	Remark     string         `gorm:"size:500" json:"remark"`            // 备注
	Status     int            `gorm:"default:1" json:"status"`           // 状态: 1有效 2过期
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
