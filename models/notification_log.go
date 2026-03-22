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
	TotalSent     int64            `json:"total_sent"`
	TotalSuccess  int64            `json:"total_success"`
	TotalFailed   int64            `json:"total_failed"`
	TodaySent     int64            `json:"today_sent"`
	TodaySuccess  int64            `json:"today_success"`
	TodayFailed   int64            `json:"today_failed"`
	ByChannel     map[string]int64 `json:"by_channel"`
	RecentLogs    []NotificationLog `json:"recent_logs"`
}

// NotificationTemplate 通知模板（扩展自 NotificationTemplate，定义在 notification.go）
// 本文件定义模板相关的数据结构
type NotificationTemplateV2 struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TenantID      string    `gorm:"type:uuid;index" json:"tenant_id"`
	TemplateName  string    `gorm:"type:varchar(100);not null" json:"template_name"`
	ChannelType   string    `gorm:"type:varchar(20);not null" json:"channel_type"` // email/sms/webhook
	TemplateType  string    `gorm:"type:varchar(50);not null" json:"template_type"` // alert/reminder/system
	SubjectTpl    string    `gorm:"type:varchar(255)" json:"subject_tpl"`         // 标题模板，支持 {{variable}}
	BodyTpl       string    `gorm:"type:text;not null" json:"body_tpl"`            // 内容模板，支持 {{variable}}
	Variables     JSON      `gorm:"type:jsonb" json:"variables"`                  // 支持的变量列表
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedBy     uint      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (NotificationTemplateV2) TableName() string {
	return "notification_templates_v2"
}

// CreateNotificationTemplateRequest 创建模板请求
type CreateNotificationTemplateRequest struct {
	TemplateName string `json:"template_name" binding:"required"`
	ChannelType  string `json:"channel_type" binding:"required"`
	TemplateType string `json:"template_type" binding:"required"`
	SubjectTpl   string `json:"subject_tpl"`
	BodyTpl      string `json:"body_tpl" binding:"required"`
	Variables    []string `json:"variables"`
	IsActive     bool   `json:"is_active"`
}

// UpdateNotificationTemplateRequest 更新模板请求
type UpdateNotificationTemplateRequest struct {
	TemplateName string   `json:"template_name"`
	ChannelType  string   `json:"channel_type"`
	TemplateType string   `json:"template_type"`
	SubjectTpl   string   `json:"subject_tpl"`
	BodyTpl      string   `json:"body_tpl"`
	Variables    []string `json:"variables"`
	IsActive     *bool   `json:"is_active"`
}
