package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Conversation 会话表
type Conversation struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"conversation_id"`
	UserID         int64          `gorm:"not null;index" json:"user_id"`
	DeviceID       string         `gorm:"type:varchar(64);index" json:"device_id"`
	Title          string         `gorm:"type:varchar(256);not null" json:"title"`
	LastMessage    string         `gorm:"type:text" json:"last_message"`
	LastMessageAt  *time.Time     `json:"last_message_at"`
	MessageCount   int            `gorm:"type:int;default:0" json:"message_count"`
	Status         int            `gorm:"type:smallint;default:1" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// BeforeCreate 创建前自动生成 UUID
func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ConversationID == "" {
		c.ConversationID = uuid.New().String()
	}
	return nil
}

// ConversationStatus 会话状态常量
const (
	ConversationStatusNormal = 1 // 正常
	ConversationStatusClosed = 2 // 已归档
)

// ConversationResponse 会话响应
type ConversationResponse struct {
	ConversationID string  `json:"conversation_id"`
	UserID         int64   `json:"user_id"`
	DeviceID       string  `json:"device_id"`
	Title          string  `json:"title"`
	LastMessage    string  `json:"last_message"`
	LastMessageAt  *string `json:"last_message_at"`
	MessageCount   int     `json:"message_count"`
	Status         int     `json:"status"`
	CreatedAt      string  `json:"created_at"`
}

// ToResponse 转换为响应格式
func (c *Conversation) ToResponse() *ConversationResponse {
	resp := &ConversationResponse{
		ConversationID: c.ConversationID,
		UserID:         c.UserID,
		DeviceID:       c.DeviceID,
		Title:          c.Title,
		LastMessage:    c.LastMessage,
		MessageCount:   c.MessageCount,
		Status:         c.Status,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
	}
	if c.LastMessageAt != nil {
		t := c.LastMessageAt.Format(time.RFC3339)
		resp.LastMessageAt = &t
	}
	return resp
}

// CreateConversationRequest 创建会话请求
type CreateConversationRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

// ConversationListQuery 会话列表查询参数
type ConversationListQuery struct {
	DeviceID string `form:"device_id"`
	Status   int    `form:"status"`
}
