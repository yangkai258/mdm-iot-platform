package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetShortTermMemory 短期记忆表
type PetShortTermMemory struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MemoryID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"memory_id"`
	DeviceID       string         `gorm:"type:varchar(64);index;not null" json:"device_id"`
	UserID         int64          `gorm:"not null" json:"user_id"`
	SessionID      string         `gorm:"type:varchar(64);index;not null" json:"session_id"`
	MessageID      string         `gorm:"type:varchar(64)" json:"message_id"`
	MemoryType     string         `gorm:"type:varchar(32);not null" json:"memory_type"`
	Content        string         `gorm:"type:jsonb;not null" json:"content"`
	Importance     float64        `gorm:"type:float;default:0.5" json:"importance"`
	AccessCount    int            `gorm:"type:int;default:0" json:"access_count"`
	LastAccessedAt *time.Time     `json:"last_accessed_at"`
	ExpiresAt      *time.Time     `json:"expires_at" index:"expires_at"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetShortTermMemory) TableName() string {
	return "short_term_memory"
}

// BeforeCreate 创建前自动生成 UUID
func (s *PetShortTermMemory) BeforeCreate(tx *gorm.DB) error {
	if s.MemoryID == "" {
		s.MemoryID = uuid.New().String()
	}
	return nil
}

// MemoryType 常量
const (
	MemoryTypeInteraction  = "interaction"  // 交互记忆
	MemoryTypeConversation = "conversation" // 对话记忆
	MemoryTypePreference   = "preference"   // 偏好记忆
	MemoryTypeContext      = "context"      // 上下文记忆
)

// ShortTermMemoryResponse 短期记忆响应
type ShortTermMemoryResponse struct {
	MemoryID       string  `json:"memory_id"`
	DeviceID       string  `json:"device_id"`
	UserID         int64   `json:"user_id"`
	SessionID      string  `json:"session_id"`
	MessageID      string  `json:"message_id"`
	MemoryType     string  `json:"memory_type"`
	Content        string  `json:"content"`
	Importance     float64 `json:"importance"`
	AccessCount    int     `json:"access_count"`
	LastAccessedAt *string `json:"last_accessed_at"`
	ExpiresAt      *string `json:"expires_at"`
	CreatedAt      string  `json:"created_at"`
}

// ToResponse 转换为响应格式
func (s *PetShortTermMemory) ToResponse() *ShortTermMemoryResponse {
	resp := &ShortTermMemoryResponse{
		MemoryID:    s.MemoryID,
		DeviceID:    s.DeviceID,
		UserID:      s.UserID,
		SessionID:   s.SessionID,
		MessageID:   s.MessageID,
		MemoryType:  s.MemoryType,
		Content:     s.Content,
		Importance:  s.Importance,
		AccessCount: s.AccessCount,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
	}
	if s.LastAccessedAt != nil {
		t := s.LastAccessedAt.Format(time.RFC3339)
		resp.LastAccessedAt = &t
	}
	if s.ExpiresAt != nil {
		t := s.ExpiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &t
	}
	return resp
}

// CreateMemoryRequest 创建记忆请求
type CreateMemoryRequest struct {
	SessionID  string                 `json:"session_id" binding:"required"`
	MessageID  string                 `json:"message_id"`
	MemoryType string                 `json:"memory_type" binding:"required"`
	Content    map[string]interface{} `json:"content" binding:"required"`
	Importance float64                `json:"importance"`
	ExpiresIn  int                    `json:"expires_in"` // 秒
}

// MemoryListQuery 记忆列表查询参数
type MemoryListQuery struct {
	SessionID string `form:"session_id"`
	Type      string `form:"type"`
	Limit     int    `form:"limit,default=20"`
	Offset    int    `form:"offset,default=0"`
}
