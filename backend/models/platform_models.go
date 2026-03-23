package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== 开发者应用 =====

// DeveloperApp 开发者应用
type DeveloperApp struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"` // 开发者用户ID
	AppName     string         `gorm:"type:varchar(255);not null" json:"app_name"`
	AppType     string         `gorm:"type:varchar(50)" json:"app_type"`            // personal/enterprise
	Description string         `gorm:"type:text" json:"description"`
	AppIcon     string         `gorm:"type:varchar(500)" json:"app_icon"`
	WebsiteURL  string         `gorm:"type:varchar(500)" json:"website_url"`
	CallbackURL string         `gorm:"type:varchar(500)" json:"callback_url"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/suspended/deleted
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (DeveloperApp) TableName() string {
	return "developer_apps"
}

// ===== API Key =====

// APIKey API Key
type APIKey struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AppID      uint           `gorm:"index;not null" json:"app_id"` // 关联应用
	KeyPrefix  string         `gorm:"type:varchar(20);not null" json:"key_prefix"`
	KeyHash    string         `gorm:"type:varchar(255);not null" json:"-"` // 不返回给前端
	KeyType    string         `gorm:"type:varchar(20);not null" json:"key_type"` // api_key/oauth_client
	Scopes     StringArray    `gorm:"type:text" json:"scopes"`              // 权限范围
	RateLimit  int            `gorm:"default:1000" json:"rate_limit"`       // 每分钟限制
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	LastUsedAt *time.Time     `json:"last_used_at"`
	ExpiresAt  *time.Time     `json:"expires_at"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (APIKey) TableName() string {
	return "api_keys"
}

// ===== Webhook 订阅 =====

// WebhookSubscription Webhook 订阅
type WebhookSubscription struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`            // 订阅用户
	AppID       *uint          `gorm:"index" json:"app_id"`                       // 可选关联应用
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	URL         string         `gorm:"type:varchar(512);not null" json:"url"`
	Secret      string         `gorm:"type:varchar(256)" json:"secret"`           // 签名密钥
	EventTypes  StringArray    `gorm:"type:text" json:"event_types"`             // ["device.alert","subscription.created"]
	Headers     JSON           `gorm:"type:jsonb" json:"headers"`                // 自定义请求头
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive
	RetryCount  int            `gorm:"default:3" json:"retry_count"`            // 重试次数
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (WebhookSubscription) TableName() string {
	return "webhook_subscriptions"
}

// ===== Webhook 事件 =====

// WebhookEventRecord Webhook 事件记录
type WebhookEventRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SubscriptionID uint         `gorm:"index;not null" json:"subscription_id"`
	EventID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"event_id"`
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

// TableName 表名
func (WebhookEventRecord) TableName() string {
	return "webhook_events"
}

// ===== Webhook 投递记录 =====

// WebhookDelivery Webhook 投递记录
type WebhookDelivery struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	EventID        string         `gorm:"type:varchar(64);index;not null" json:"event_id"`
	SubscriptionID uint           `gorm:"index;not null" json:"subscription_id"`
	URL            string         `gorm:"type:varchar(512);not null" json:"url"`
	RequestHeaders JSON           `gorm:"type:jsonb" json:"request_headers"`
	RequestBody    string         `gorm:"type:text" json:"request_body"`
	ResponseCode   int            `gorm:"type:smallint" json:"response_code"`
	ResponseBody   string         `gorm:"type:text" json:"response_body"`
	Status         string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/success/failed
	Attempts       int            `gorm:"default:0" json:"attempts"`
	ErrorMessage   string         `gorm:"type:varchar(512)" json:"error_message"`
	DurationMs     int            `gorm:"default:0" json:"duration_ms"` // 耗时毫秒
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// TableName 表名
func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}

// ===== Webhook 模板 =====

// WebhookTemplate Webhook 模板（预定义的 webhook 类型）
type WebhookTemplate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(128);not null" json:"name"`
	EventType   string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"event_type"`
	Description string    `gorm:"type:text" json:"description"`
	PayloadExample JSON   `gorm:"type:jsonb" json:"payload_example"` // 示例 payload
	Category    string    `gorm:"type:varchar(50)" json:"category"`  // device/subscription/payment/alert
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 表名
func (WebhookTemplate) TableName() string {
	return "webhook_templates"
}

// ===== 开发者统计 =====

// DeveloperStats 开发者统计
type DeveloperStats struct {
	TotalApps        int   `json:"total_apps"`
	TotalAPIKeys     int   `json:"total_api_keys"`
	TotalCalls       int64 `json:"total_calls"`
	QuotaUsagePercent float64 `json:"quota_usage_percent"`
}
