package models

import (
	"time"
)

// AlertSettings 全局告警设置
type AlertSettings struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AlertsEnabled     bool      `gorm:"default:true" json:"alerts_enabled"`      // 全局告警开关
	EmailEnabled      bool      `gorm:"default:false" json:"email_enabled"`      // 邮件通知开关
	SMSEnabled        bool      `gorm:"default:false" json:"sms_enabled"`        // 短信通知开关
	WebhookEnabled    bool      `gorm:"default:false" json:"webhook_enabled"`    // Webhook通知开关
	InAppEnabled      bool      `gorm:"default:true" json:"inapp_enabled"`      // 站内通知开关
	NotifyOnCritical  bool      `gorm:"default:true" json:"notify_on_critical"`  // 严重告警通知
	NotifyOnHigh      bool      `gorm:"default:true" json:"notify_on_high"`      // 高优先级通知
	NotifyOnMedium    bool      `gorm:"default:true" json:"notify_on_medium"`    // 中优先级通知
	NotifyOnLow       bool      `gorm:"default:false" json:"notify_on_low"`      // 低优先级通知
	DigestEnabled     bool      `gorm:"default:false" json:"digest_enabled"`     // 告警摘要（定时汇总）
	DigestInterval    int       `gorm:"default:60" json:"digest_interval"`       // 摘要间隔（分钟）
	QuietHoursEnabled bool      `gorm:"default:false" json:"quiet_hours_enabled"` // 静默时段开关
	QuietHoursStart   string    `gorm:"type:varchar(10);default:'22:00'" json:"quiet_hours_start"` // 静默开始时间
	QuietHoursEnd     string    `gorm:"type:varchar(10);default:'08:00'" json:"quiet_hours_end"`   // 静默结束时间
	MaxPerHour        int       `gorm:"default:100" json:"max_per_hour"`        // 每小时最大告警数
	AutoResolveHours  int       `gorm:"default:24" json:"auto_resolve_hours"`   // 自动解决超时（小时）
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (AlertSettings) TableName() string {
	return "alert_settings"
}

// NotificationChannel 通知渠道配置（SMTP/Webhook/SMS）
type NotificationChannel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ChannelType string    `gorm:"type:varchar(20);not null;index" json:"channel_type"` // smtp, webhook, sms
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`              // 配置名称
	Enabled     bool      `gorm:"default:true" json:"enabled"`                          // 是否启用

	// SMTP 配置
	SMTPHost      string `gorm:"type:varchar(255)" json:"smtp_host"`
	SMTPPort      int    `gorm:"default:587" json:"smtp_port"`
	SMTPUser      string `gorm:"type:varchar(255)" json:"smtp_user"`
	SMTPPassword  string `gorm:"type:varchar(255)" json:"smtp_password"` // 加密存储
	SMTPFrom      string `gorm:"type:varchar(255)" json:"smtp_from"`
	SMTPTo        string `gorm:"type:varchar(255)" json:"smtp_to"` // 默认收件人，多个用逗号分隔
	SMTPUseTLS    bool   `gorm:"default:true" json:"smtp_use_tls"`

	// Webhook 配置
	WebhookURL    string `gorm:"type:varchar(500)" json:"webhook_url"`
	WebhookToken  string `gorm:"type:varchar(255)" json:"webhook_token"` // 加密存储
	WebhookMethod string `gorm:"type:varchar(10);default:'POST'" json:"webhook_method"`

	// SMS 配置（预留）
	SMSProvider string `gorm:"type:varchar(50)" json:"sms_provider"` // aliyun, twilio, etc.
	SMSAccount  string `gorm:"type:varchar(255)" json:"sms_account"`
	SMSSecret   string `gorm:"type:varchar(255)" json:"sms_secret"` // 加密存储
	SMSFrom     string `gorm:"type:varchar(20)" json:"sms_from"`

	// Sprint 11 扩展字段
	Priority       int        `gorm:"default:0" json:"priority"`             // 优先级
	HealthStatus   string     `gorm:"type:varchar(20);default:'unknown'" json:"health_status"` // healthy/unhealthy/unknown
	LastCheckedAt *time.Time `json:"last_checked_at"`
	TenantID      string     `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	CreatedBy     uint       `json:"created_by"`

	// 通用配置
	Remark    string `gorm:"type:varchar(500)" json:"remark"`    // 备注
	IsDefault bool   `gorm:"default:false" json:"is_default"`   // 是否为默认渠道

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (NotificationChannel) TableName() string {
	return "notification_channels"
}
