package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringArray 字符串数组类型
type StringArray []string

// Scan 实现 sql.Scanner 接口
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// Value 实现 driver.Valuer 接口
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(s)
}

// ===== Webhook 配置 =====

// Webhook Webhook 配置
type Webhook struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	WebhookID  string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"webhook_id"`
	Name       string         `gorm:"type:varchar(128);not null" json:"name"`
	URL        string         `gorm:"type:varchar(512);not null" json:"url"`
	Secret     string         `gorm:"type:varchar(256)" json:"secret"`           // 签名密钥
	EventTypes StringArray    `gorm:"type:jsonb" json:"event_types"`             // ["subscription.created", "payment.success"]
	Status     string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive
	TenantID   uint           `gorm:"index" json:"tenant_id"`
	Headers    JSON           `gorm:"type:jsonb" json:"headers"`                // 自定义请求头
	RetryCount int            `gorm:"default:3" json:"retry_count"`            // 重试次数
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Webhook) TableName() string {
	return "webhooks"
}

// ===== Webhook 事件 =====

// WebhookEvent Webhook 事件
type WebhookEvent struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	EventID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"event_id"`
	WebhookID    string         `gorm:"index;type:varchar(64);not null" json:"webhook_id"`
	EventType    string         `gorm:"type:varchar(64);not null" json:"event_type"`
	Payload      JSON           `gorm:"type:jsonb" json:"payload"`
	Status       string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/success/failed
	Attempts     int            `gorm:"default:0" json:"attempts"`
	MaxAttempts  int            `gorm:"default:3" json:"max_attempts"`
	LastError    string         `gorm:"type:varchar(512)" json:"last_error"`
	ResponseCode int            `gorm:"type:smallint" json:"response_code"`
	ResponseBody string         `gorm:"type:text" json:"response_body"`
	DeliveredAt  *time.Time     `json:"delivered_at"`
	NextRetryAt  *time.Time     `json:"next_retry_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (WebhookEvent) TableName() string {
	return "webhook_events"
}

// Webhook 事件类型常量
const (
	EventTypeSubscriptionCreated = "subscription.created"
	EventTypeSubscriptionExpired = "subscription.expired"
	EventTypeSubscriptionCancelled = "subscription.cancelled"
	EventTypeSubscriptionRenewed  = "subscription.renewed"
	EventTypePaymentSuccess      = "payment.success"
	EventTypePaymentFailed       = "payment.failed"
	EventTypeQuotaExceeded       = "quota.exceeded"
	EventTypeWebhookTest         = "webhook.test"
)

// Webhook 事件状态
const (
	WebhookStatusPending = "pending"
	WebhookStatusSuccess = "success"
	WebhookStatusFailed = "failed"
)
