package models

import (
	"time"

	"gorm.io/gorm"
)

// SubscriptionGift 订阅赠送记录
type SubscriptionGift struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	GiftCode      string         `gorm:"type:varchar(64);uniqueIndex" json:"gift_code"`    // 赠送码
	SenderID      string         `gorm:"type:varchar(64);index" json:"sender_id"`        // 赠送者ID
	SenderName    string         `gorm:"type:varchar(100)" json:"sender_name"`          // 赠送者名称
	RecipientID   string         `gorm:"type:varchar(64);index" json:"recipient_id"`   // 接收者ID
	RecipientName string         `gorm:"type:varchar(100)" json:"recipient_name"`       // 接收者名称
	RecipientEmail string        `gorm:"type:varchar(200)" json:"recipient_email"`     // 接收者邮箱
	PlanID        uint           `gorm:"not null" json:"plan_id"`                     // 订阅计划ID
	PlanName      string         `gorm:"type:varchar(100)" json:"plan_name"`          // 计划名称
	Duration      int            `gorm:"default:30" json:"duration"`                  // 赠送天数
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/sent/claimed/expired/cancelled
	ClaimedAt     *time.Time     `json:"claimed_at"`                                  // 领取时间
	ExpiresAt     *time.Time     `json:"expires_at"`                                  // 过期时间
	SentAt        *time.Time     `json:"sent_at"`                                    // 发送时间
	Message       string         `gorm:"type:text" json:"message"`                   // 祝福语
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (SubscriptionGift) TableName() string {
	return "subscription_gifts"
}

// SubscriptionGiftUsage 赠送码使用统计
type SubscriptionGiftUsage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	GiftID         uint           `gorm:"not null;index" json:"gift_id"`
	UsedByDevice   string         `gorm:"type:varchar(64)" json:"used_by_device"`
	UsedByIP       string         `gorm:"type:varchar(50)" json:"used_by_ip"`
	ClaimedFeature string         `gorm:"type:varchar(100)" json:"claimed_feature"` // 使用的功能
	CreatedAt      time.Time      `json:"created_at"`
}

// TableName 表名
func (SubscriptionGiftUsage) TableName() string {
	return "subscription_gift_usage"
}
