package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message 消息表
type Message struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MessageID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"message_id"`
	ConversationID string         `gorm:"type:varchar(64);index;not null" json:"conversation_id"`
	SenderType     int            `gorm:"type:smallint;not null" json:"sender_type"`
	SenderID       *int64         `json:"sender_id"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	ContentType    int            `gorm:"type:smallint;default:1" json:"content_type"`
	MediaURL       string         `gorm:"type:varchar(512)" json:"media_url"`
	Intent         string         `gorm:"type:varchar(64)" json:"intent"`
	Confidence     *float64       `json:"confidence"`
	Metadata       string         `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// BeforeCreate 创建前自动生成 UUID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.MessageID == "" {
		m.MessageID = uuid.New().String()
	}
	return nil
}

// SenderType 常量
const (
	SenderTypeUser = 1 // 用户
	SenderTypePet  = 2 // 宠物
)

// ContentType 常量
const (
	ContentTypeText   = 1 // 文本
	ContentTypeVoice  = 2 // 语音
	ContentTypeImage  = 3 // 图片
	ContentTypeAction = 4 // 动作
)

// MessageResponse 消息响应
type MessageResponse struct {
	MessageID      string   `json:"message_id"`
	ConversationID string   `json:"conversation_id"`
	SenderType     int      `json:"sender_type"`
	SenderID       *int64   `json:"sender_id"`
	Content        string   `json:"content"`
	ContentType    int      `json:"content_type"`
	MediaURL       string   `json:"media_url"`
	Intent         string   `json:"intent"`
	Confidence     *float64 `json:"confidence"`
	Metadata       string   `json:"metadata"`
	CreatedAt      string   `json:"created_at"`
}

// ToResponse 转换为响应格式
func (m *Message) ToResponse() *MessageResponse {
	return &MessageResponse{
		MessageID:      m.MessageID,
		ConversationID: m.ConversationID,
		SenderType:     m.SenderType,
		SenderID:       m.SenderID,
		Content:        m.Content,
		ContentType:    m.ContentType,
		MediaURL:       m.MediaURL,
		Intent:         m.Intent,
		Confidence:     m.Confidence,
		Metadata:       m.Metadata,
		CreatedAt:      m.CreatedAt.Format(time.RFC3339),
	}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	ContentType int    `json:"content_type"`
	MediaURL    string `json:"media_url"`
	Metadata    string `json:"metadata"`
}

// MessageListQuery 消息列表查询参数
type MessageListQuery struct {
	Limit  int `form:"limit,default=20"`
	Offset int `form:"offset,default=0"`
}
