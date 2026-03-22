package models

import (
	"time"
)

// NotificationLog 通知发送日志
type NotificationLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ChannelID    uint      `gorm:"index" json:"channel_id"`
	ChannelType  string    `gorm:"type:varchar(20);not null" json:"channel_type"` // email/sms/webhook
	AlertID      uint      `gorm:"index" json:"alert_id"`
	Recipient    string    `gorm:"type:varchar(256)" json:"recipient"` // email/phone/url
	Subject      string    `gorm:"type:varchar(256)" json:"subject"`
	Content      string    `gorm:"type:text" json:"content"`
	Status       string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending/success/failed
	RetryCount   int       `gorm:"default:0" json:"retry_count"`
	ErrorCode    string    `gorm:"type:varchar(50)" json:"error_code"`
	ErrorMsg     string    `gorm:"type:varchar(512)" json:"error_msg"`
	SentAt       *time.Time `json:"sent_at"`
	TenantID     string    `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}

// NotificationLogQuery 通知日志查询参数
type NotificationLogQuery struct {
	ChannelID   uint   `form:"channel_id"`
	ChannelType string `form:"channel_type"`
	AlertID     uint   `form:"alert_id"`
	Status      string `form:"status"` // pending/success/failed
	StartTime   string `form:"start_time"`
	EndTime     string `form:"end_time"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
	Keyword     string `form:"keyword"` // 搜索 recipient/subject
}

// NotificationStats 通知统计
type NotificationStats struct {
	TotalSent    int64 `json:"total_sent"`
	TotalSuccess int64 `json:"total_success"`
	TotalFailed  int64 `json:"total_failed"`
	TodaySent    int64 `json:"today_sent"`
	TodaySuccess int64 `json:"today_success"`
	TodayFailed  int64 `json:"today_failed"`
	ByChannel    map[string]int64 `json:"by_channel"`
}
