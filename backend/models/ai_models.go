package models

import (
	"time"
)

// AIConversation AI对话会话
type AIConversation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`         // 用户ID
	SessionID  string    `gorm:"type:varchar(100);index" json:"session_id"` // 会话ID
	Title      string    `gorm:"type:varchar(200)" json:"title"`           // 对话标题
	Model      string    `gorm:"type:varchar(50)" json:"model"`             // 使用的模型
	Status     int       `gorm:"default:1" json:"status"`                  // 状态: 0关闭 1活跃
	TenantID   string    `gorm:"index" json:"tenant_id"`                    // 租户ID
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AIConversation) TableName() string { return "ai_conversations" }

// AIMessage AI对话消息
type AIMessage struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID uint      `gorm:"index" json:"conversation_id"` // 会话ID
	Role           string    `gorm:"type:varchar(20);not null" json:"role"` // 角色: user/assistant/system
	Content        string    `gorm:"type:text" json:"content"`              // 消息内容
	TokenCount     int       `gorm:"default:0" json:"token_count"`         // Token数量
	Model          string    `gorm:"type:varchar(50)" json:"model"`         // 使用的模型
	LatencyMs      int       `gorm:"default:0" json:"latency_ms"`           // 响应延迟(毫秒)
	TenantID       string    `gorm:"index" json:"tenant_id"`                // 租户ID
	CreatedAt      time.Time `json:"created_at"`
}

func (AIMessage) TableName() string { return "ai_messages" }

// AIConfig AI配置
type AIConfig struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Provider   string    `gorm:"type:varchar(50);uniqueIndex" json:"provider"` // 提供商: openai/anthropic/custom
	Model      string    `gorm:"type:varchar(100)" json:"model"`               // 模型名称
	APIKey     string    `gorm:"type:varchar(255)" json:"-"`                   // API密钥(不返回)
	BaseURL    string    `gorm:"type:varchar(255)" json:"base_url"`            // API地址
	MaxTokens  int       `gorm:"default:4096" json:"max_tokens"`              // 最大Token
	Temperature float64  `gorm:"default:0.7" json:"temperature"`              // 温度参数
	Status     int       `gorm:"default:1" json:"status"`                    // 状态: 0禁用 1启用
	IsDefault  int       `gorm:"default:0" json:"is_default"`                // 是否默认: 0否 1是
	TenantID   string    `gorm:"index" json:"tenant_id"`                      // 租户ID
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (AIConfig) TableName() string { return "ai_configs" }
