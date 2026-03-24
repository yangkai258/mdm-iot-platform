package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Webhook Webhook配置
type Webhook struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	WebhookUUID  string         `gorm:"size:64;uniqueIndex;not null" json:"webhook_uuid"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	URL          string         `gorm:"type:varchar(512);not null" json:"url"`
	Method       string         `gorm:"size:10;default:'POST'" json:"method"` // POST, GET, PUT, DELETE
	Secret       string         `gorm:"size:256" json:"secret"`               // 签名密钥
	Headers      json.RawMessage `gorm:"type:jsonb" json:"headers"`          // 自定义请求头
	EventTypes   []string       `gorm:"type:text[]" json:"event_types"`       // 触发事件类型
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	RetryPolicy  json.RawMessage `gorm:"type:jsonb" json:"retry_policy"`     // 重试策略
	TimeoutSec   int            `gorm:"default:30" json:"timeout_sec"`
	TenantID      string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedBy    string         `gorm:"size:64" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) error {
	if w.WebhookUUID == "" {
		w.WebhookUUID = uuid.New().String()
	}
	return nil
}

// WebhookLog Webhook调用日志
type WebhookLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	LogUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"log_uuid"`
	WebhookID    uint           `gorm:"not null;index" json:"webhook_id"`
	EventType    string         `gorm:"size:50;index" json:"event_type"`
	RequestBody  string         `gorm:"type:text" json:"request_body"`
	ResponseBody string         `gorm:"type:text" json:"response_body"`
	StatusCode   int            `gorm:"default:0" json:"status_code"`
	LatencyMs    int64          `json:"latency_ms"`
	Result       string         `gorm:"size:20;index" json:"result"` // success, failed, timeout, error
	ErrorMsg     string         `gorm:"type:text" json:"error_msg"`
	RetryCount   int            `gorm:"default:0" json:"retry_count"`
	RequestAt    time.Time      `gorm:"not null" json:"request_at"`
	ResponseAt   *time.Time     `json:"response_at"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (w *WebhookLog) BeforeCreate(tx *gorm.DB) error {
	if w.LogUUID == "" {
		w.LogUUID = uuid.New().String()
	}
	return nil
}
