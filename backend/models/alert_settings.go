package models

import (
	"time"
)

// AlertSettings 全局告警设置
type AlertSettings struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	AlertsEnabled     bool      `gorm:"default:true" json:"alerts_enabled"`      // 全局告警开关
	EmailEnabled      bool      `gorm:"default:false" json:"email_enabled"`       // 邮件通知开关
	SMSEnabled        bool      `gorm:"default:false" json:"sms_enabled"`         // 短信通知开关
	WebhookEnabled    bool      `gorm:"default:false" json:"webhook_enabled"`    // Webhook通知开关
	InAppEnabled      bool      `gorm:"default:true" json:"inapp_enabled"`        // 站内通知开关
	NotifyOnCritical  bool      `gorm:"default:true" json:"notify_on_critical"`   // 严重告警通知
	NotifyOnHigh      bool      `gorm:"default:true" json:"notify_on_high"`       // 高优先级通知
	NotifyOnMedium    bool      `gorm:"default:true" json:"notify_on_medium"`     // 中优先级通知
	NotifyOnLow       bool      `gorm:"default:false" json:"notify_on_low"`      // 低优先级通知
	DigestEnabled     bool      `gorm:"default:false" json:"digest_enabled"`      // 告警摘要（定时汇总）
	DigestInterval    int       `gorm:"default:60" json:"digest_interval"`       // 摘要间隔（分钟）
	QuietHoursEnabled bool      `gorm:"default:false" json:"quiet_hours_enabled"`  // 静默时段开关
	QuietHoursStart   string    `gorm:"type:varchar(10);default:'22:00'" json:"quiet_hours_start"` // 静默开始时间
	QuietHoursEnd     string    `gorm:"type:varchar(10);default:'08:00'" json:"quiet_hours_end"`   // 静默结束时间
	MaxPerHour        int       `gorm:"default:100" json:"max_per_hour"`         // 每小时最大告警数
	AutoResolveHours  int       `gorm:"default:24" json:"auto_resolve_hours"`    // 自动解决超时（小时）
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (AlertSettings) TableName() string {
	return "alert_settings"
}
