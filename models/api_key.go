package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// APIKey API 密钥
type APIKey struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	AppID        uint           `gorm:"index;not null" json:"app_id"`          // 关联应用
	KeyID        string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"key_id"` // 密钥标识 (显示给用户)
	KeyPrefix    string         `gorm:"type:varchar(16);not null" json:"key_prefix"`         // 密钥前缀 (mdm_sk_xxx)
	KeyHash      string         `gorm:"type:varchar(128);not null" json:"-"`                // 密钥哈希 (存储用)
	Name         string         `gorm:"type:varchar(128)" json:"name"`         // 密钥名称
	Scopes       JSON           `gorm:"type:jsonb;default:'[]'" json:"scopes"` // 权限范围 ["device:read","device:write"]
	RateLimit    int            `gorm:"default:1000" json:"rate_limit"`        // 速率限制 (次/分钟)
	Status       int            `gorm:"type:smallint;default:1" json:"status"` // 1:启用 0:禁用
	LastUsedAt   *time.Time     `json:"last_used_at"`
	ExpiresAt    *time.Time     `json:"expires_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (APIKey) TableName() string {
	return "api_keys"
}

// APIKeyStatus 密钥状态常量
const (
	APIKeyStatusEnabled  = 1
	APIKeyStatusDisabled = 0
)

// GenerateKeyPair 生成新的密钥对 (返回原始密钥和存储用哈希)
func GenerateKeyPair() (secret string, hash string, err error) {
	// 生成32字节随机数作为密钥
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	secret = "mdm_sk_" + hex.EncodeToString(bytes)

	// SHA256哈希存储
	h := sha256.New()
	h.Write([]byte(secret))
	hash = hex.EncodeToString(h.Sum(nil))
	return secret, hash, nil
}

// GetKeyPrefix 从密钥提取前缀
func GetKeyPrefix(key string) string {
	if len(key) > 16 {
		return key[:16]
	}
	return key
}
