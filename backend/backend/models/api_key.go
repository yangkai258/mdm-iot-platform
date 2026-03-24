package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeveloperApp 开发者应用
type DeveloperApp struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	AppUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"app_uuid"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	AppType      string         `gorm:"size:32;default:'standard'" json:"app_type"` // standard, enterprise
	WebsiteURL   string         `gorm:"size:512" json:"website_url"`
	LogoURL      string         `gorm:"size:512" json:"logo_url"`
	Status       string         `gorm:"size:20;default:'active';index" json:"status"` // active, suspended, inactive
	OwnerID      string         `gorm:"size:64;not null;index" json:"owner_id"`
	OwnerName    string         `gorm:"size:128" json:"owner_name"`
	QuotaTier    string         `gorm:"size:20;default:'free'" json:"quota_tier"` // free, basic, pro, enterprise
	RateLimit    int            `gorm:"default:100" json:"rate_limit"`              // 每分钟请求上限
	MonthlyQuota int64          `gorm:"default:10000" json:"monthly_quota"`        // 月度配额
	MonthlyUsage int64          `gorm:"default:0" json:"monthly_usage"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DeveloperApp) BeforeCreate(tx *gorm.DB) error {
	if d.AppUUID == "" {
		d.AppUUID = uuid.New().String()
	}
	return nil
}

// APIKey API密钥
type APIKey struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	KeyUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"key_uuid"`
	AppID        uint           `gorm:"not null;index" json:"app_id"`
	KeyID        string         `gorm:"size:32;uniqueIndex;not null" json:"key_id"` // 公开的Key标识
	KeyHash      string         `gorm:"size:128;not null" json:"-"`                 // 存储hash，不返回原始key
	KeyPrefix    string         `gorm:"size:8;not null" json:"key_prefix"`          // key前8位用于识别
	Name         string         `gorm:"size:128" json:"name"`                      // key名称
	Scopes       []string       `gorm:"type:text[]" json:"scopes"`                 // 权限范围
	Status       string         `gorm:"size:20;default:'active';index" json:"status"` // active, expired, revoked
	ExpiresAt    *time.Time     `json:"expires_at"`
	LastUsedAt   *time.Time     `json:"last_used_at"`
	LastUsedIP   string         `gorm:"size:45" json:"last_used_ip"`
	TotalCalls   int64          `gorm:"default:0" json:"total_calls"`
	TodayCalls   int64          `gorm:"default:0" json:"today_calls"`
	RateLimitHit int64          `gorm:"default:0" json:"rate_limit_hit"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	if a.KeyUUID == "" {
		a.KeyUUID = uuid.New().String()
	}
	// 生成随机key prefix
	bytes := make([]byte, 4)
	rand.Read(bytes)
	a.KeyPrefix = hex.EncodeToString(bytes)[:8]
	return nil
}

// APIKeyUsage API Key使用记录
type APIKeyUsage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	KeyID     uint      `gorm:"not null;index" json:"key_id"`
	AppID     uint      `gorm:"not null;index" json:"app_id"`
	Path      string    `gorm:"size:255;index" json:"path"`
	Method    string    `gorm:"size:10" json:"method"`
	StatusCode int       `json:"status_code"`
	LatencyMs int64     `json:"latency_ms"`
	IP        string    `gorm:"size:45" json:"ip"`
	UA        string    `gorm:"size:512" json:"ua"`
	CallAt    time.Time `gorm:"not null;index" json:"call_at"`
}

// GenerateKey 生成新的API Key（返回原始key，hash后存储）
func GenerateAPIKey() (rawKey string, keyID string) {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	keyID = "sk_" + hex.EncodeToString(bytes)
	rawKey = keyID
	return
}
