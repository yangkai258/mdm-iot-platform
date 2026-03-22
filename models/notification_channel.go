package models

import (
	"time"
)

// NotificationChannel 通知渠道配置
type NotificationChannel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ChannelType string    `gorm:"type:varchar(20);not null" json:"channel_type"` // smtp, webhook, sms
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`        // 配置名称
	Enabled     bool      `gorm:"default:false" json:"enabled"`                  // 是否启用

	// SMTP 配置
	SMTPHost     string `gorm:"type:varchar(255)" json:"smtp_host"`
	SMTPPort     int    `gorm:"default:587" json:"smtp_port"`
	SMTPUser     string `gorm:"type:varchar(255)" json:"smtp_user"`
	SMTPPassword string `gorm:"type:varchar(255)" json:"smtp_password"` // 加密存储
	SMTPFrom     string `gorm:"type:varchar(255)" json:"smtp_from"`
	SMTPTo       string `gorm:"type:varchar(255)" json:"smtp_to"`       // 默认收件人，多个用逗号分隔
	SMTPUseTLS   bool   `gorm:"default:true" json:"smtp_use_tls"`

	// Webhook 配置
	WebhookURL    string `gorm:"type:varchar(500)" json:"webhook_url"`
	WebhookToken  string `gorm:"type:varchar(255)" json:"webhook_token"`  // 加密存储
	WebhookMethod string `gorm:"type:varchar(10);default:'POST'" json:"webhook_method"`

	// SMS 配置
	SMSProvider string `gorm:"type:varchar(50)" json:"sms_provider"`
	SMSAccount  string `gorm:"type:varchar(255)" json:"sms_account"`
	SMSSecret   string `gorm:"type:varchar(255)" json:"sms_secret"`
	SMSFrom     string `gorm:"type:varchar(20)" json:"sms_from"`

	// 通用字段
	Priority      int        `gorm:"default:0" json:"priority"`                  // 优先级
	HealthStatus string     `gorm:"type:varchar(20);default:'unknown'" json:"health_status"` // healthy/unhealthy/unknown
	LastCheckedAt *time.Time `json:"last_checked_at"`
	Remark       string     `gorm:"type:text" json:"remark"`
	IsDefault    bool       `gorm:"default:false" json:"is_default"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (NotificationChannel) TableName() string {
	return "notification_channels"
}
