package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSON 自定义 JSON 类型（用于 GORM 存储）
type JSON map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// NotificationChannel 通知渠道配置
type NotificationChannel struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ChannelType  string    `gorm:"type:varchar(20);not null;index" json:"channel_type"` // email/sms/webhook
	ChannelName  string    `gorm:"type:varchar(100);not null" json:"channel_name"`
	Config       JSON      `gorm:"type:jsonb" json:"config"` // 渠道配置（敏感信息加密存储）
	Enabled      bool      `gorm:"default:true" json:"enabled"`
	IsDefault    bool      `gorm:"default:false" json:"is_default"`
	Priority     int       `gorm:"default:0" json:"priority"` // 优先级，数字越大优先级越高
	HealthStatus string    `gorm:"type:varchar(20);default:'unknown'" json:"health_status"` // healthy/unhealthy/unknown
	LastCheckedAt *time.Time `json:"last_checked_at"`
	TenantID     string    `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	CreatedBy    uint      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (NotificationChannel) TableName() string {
	return "notification_channels"
}

// NotificationChannelConfig 渠道配置（用于 API 请求/响应，加密字段脱敏）
type NotificationChannelConfig struct {
	// Email 配置
	SMTPHost     string `json:"smtp_host,omitempty"`
	SMTPPort     int    `json:"smtp_port,omitempty"`
	SMTPUsername string `json:"smtp_username,omitempty"`
	SMTPPassword string `json:"smtp_password,omitempty"` // 仅在创建/更新时传递，不返回
	SMTPFrom     string `json:"smtp_from,omitempty"`
	UseTLS       bool   `json:"use_tls,omitempty"`

	// SMS 配置
	SMSProvider  string `json:"sms_provider,omitempty"`  // aliyun/tencent
	SMSAccessKey string `json:"sms_access_key,omitempty"`
	SMSSecretKey string `json:"sms_secret_key,omitempty"` // 仅在创建/更新时传递，不返回
	SMSSignName  string `json:"sms_sign_name,omitempty"`
	SMSRegion    string `json:"sms_region,omitempty"`

	// Webhook 配置
	WebhookURL    string `json:"webhook_url,omitempty"`
	WebhookSecret string `json:"webhook_secret,omitempty"` // 仅在创建/更新时传递，不返回
}

// CreateNotificationChannelRequest 创建通知渠道请求
type CreateNotificationChannelRequest struct {
	ChannelType string                   `json:"channel_type" binding:"required"` // email/sms/webhook
	ChannelName string                   `json:"channel_name" binding:"required"`
	Config      NotificationChannelConfig `json:"config" binding:"required"`
	Enabled     bool                     `json:"enabled"`
	IsDefault   bool                     `json:"is_default"`
	Priority    int                      `json:"priority"`
}

// UpdateNotificationChannelRequest 更新通知渠道请求
type UpdateNotificationChannelRequest struct {
	ChannelName string                   `json:"channel_name"`
	Config      NotificationChannelConfig `json:"config"`
	Enabled     *bool                    `json:"enabled"`
	IsDefault   *bool                    `json:"is_default"`
	Priority    *int                     `json:"priority"`
}
