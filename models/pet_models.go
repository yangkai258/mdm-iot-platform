package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetConversation 宠物对话会话
type PetConversation struct {
	ConversationID string         `gorm:"primaryKey;type:varchar(36)" json:"conversation_id"`
	UserID         string         `gorm:"type:varchar(36);index" json:"user_id"`
	DeviceID       string         `gorm:"type:varchar(36);index" json:"device_id"`
	Title          string         `gorm:"type:varchar(128)" json:"title"`
	Status         string         `gorm:"type:varchar(16);default:'active'" json:"status"` // active, closed
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (p *PetConversation) BeforeCreate(tx *gorm.DB) error {
	if p.ConversationID == "" {
		p.ConversationID = uuid.New().String()
	}
	return nil
}

// PetMessage 宠物对话消息
type PetMessage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID string         `gorm:"type:varchar(36);index" json:"conversation_id"`
	UserID         string         `gorm:"type:varchar(36);index" json:"user_id"`
	DeviceID       string         `gorm:"type:varchar(36);index" json:"device_id"`
	SenderType     string         `gorm:"type:varchar(16);not null" json:"sender_type"` // user, pet
	SenderID       string         `gorm:"type:varchar(36)" json:"sender_id"`
	Content        string         `gorm:"type:text" json:"content"`
	MessageType    string         `gorm:"type:varchar(16);default:'text'" json:"message_type"` // text, voice, image, action
	Metadata       map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// PetStatus 宠物状态
type PetStatus struct {
	DeviceID         string    `gorm:"primaryKey;type:varchar(36)" json:"device_id"`
	Mood             string    `gorm:"type:varchar(32);default:'happy'" json:"mood"` // happy, sad, hungry, sleepy, playful, angry
	Energy           int       `gorm:"type:smallint;default:80" json:"energy"`       // 0-100
	Hunger           int       `gorm:"type:smallint;default:30" json:"hunger"`       // 0-100
	Health           int       `gorm:"type:smallint;default:100" json:"health"`     // 0-100
	Experience       int       `gorm:"type:int;default:0" json:"experience"`        // 经验值
	Level            int       `gorm:"type:int;default:1" json:"level"`             // 等级
	CurrentAction    string    `gorm:"type:varchar(64)" json:"current_action"`       // 当前动作
	CurrentEmotion   string    `gorm:"type:varchar(32)" json:"current_emotion"`      // 当前情绪
	LastInteraction  *time.Time `json:"last_interaction"`
	LastFed          *time.Time `json:"last_fed"`
	LastSlept        *time.Time `json:"last_slept"`
	TotalInteractions int      `gorm:"type:int;default:0" json:"total_interactions"`
	TotalConversations int     `gorm:"type:int;default:0" json:"total_conversations"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ShortTermMemory 短期记忆（Redis）
type ShortTermMemory struct {
	DeviceID    string                 `json:"device_id"`
	SessionID   string                 `json:"session_id"`
	Events      []MemoryEvent          `json:"events"`
	LastUpdated time.Time              `json:"last_updated"`
	TTL         int                    `json:"ttl"` // 秒
}

// MemoryEvent 记忆事件
type MemoryEvent struct {
	ID        uint                   `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"` // interaction, conversation, action, emotion, sensor
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Importance int                   `json:"importance"` // 1-10
}

// LongTermMemory 长期记忆（数据库持久化）
type LongTermMemory struct {
	ID         uint                   `gorm:"primaryKey" json:"id"`
	DeviceID   string                 `gorm:"type:varchar(36);index" json:"device_id"`
	EventType  string                 `gorm:"type:varchar(32);index" json:"event_type"`
	Content    string                 `gorm:"type:text" json:"content"`
	Summary    string                 `gorm:"type:varchar(256)" json:"summary"`
	Metadata   map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	Importance int                    `gorm:"type:smallint;default:5" json:"importance"` // 1-10
	Tags       string                 `gorm:"type:varchar(256)" json:"tags"`              // 逗号分隔标签
	Embedding  string                 `gorm:"type:text" json:"embedding"`                 // 向量嵌入（预留）
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// PetBehaviorAction 宠物行为动作记录
type PetBehaviorAction struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DeviceID     string    `gorm:"type:varchar(36);index" json:"device_id"`
	ActionType   string    `gorm:"type:varchar(64)" json:"action_type"`
	ActionName   string    `gorm:"type:varchar(128)" json:"action_name"`
	Sequence     int       `gorm:"type:int" json:"sequence"`         // 动作序号
	Duration     int       `gorm:"type:int" json:"duration"`        // 持续时间（毫秒）
	Trigger      string    `gorm:"type:varchar(32)" json:"trigger"` // manual, auto, sensor, schedule
	DecisionPath string    `gorm:"type:varchar(128)" json:"decision_path"`
	Parameters   map[string]interface{} `gorm:"type:jsonb" json:"parameters"`
	Result       string    `gorm:"type:varchar(32)" json:"result"` // success, failed, interrupted
	CreatedAt    time.Time `json:"created_at"`
}
