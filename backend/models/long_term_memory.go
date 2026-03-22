package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetLongTermMemory 长期记忆表
type PetLongTermMemory struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	MemoryID           string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"memory_id"`
	DeviceID           string         `gorm:"type:varchar(64);index;not null" json:"device_id"`
	UserID             int64          `gorm:"not null" json:"user_id"`
	MemoryCategory     string         `gorm:"type:varchar(32);index;not null" json:"memory_category"`
	Content            string         `gorm:"type:jsonb;not null" json:"content"`
	Keywords           string         `gorm:"type:jsonb;default:'[]'" json:"keywords"`
	Embedding          string         `gorm:"type:jsonb" json:"embedding"`
	Confidence         float64        `gorm:"type:float;default:0.8" json:"confidence"`
	ReinforcementCount int            `gorm:"type:int;default:1" json:"reinforcement_count"`
	LastReinforcedAt   *time.Time     `json:"last_reinforced_at"`
	DecayScore         float64        `gorm:"type:float;default:1.0" json:"decay_score"`
	IsLocked           bool           `gorm:"type:boolean;default:false" json:"is_locked"`
	SourceMemoryID     string         `gorm:"type:varchar(64)" json:"source_memory_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetLongTermMemory) TableName() string {
	return "long_term_memory"
}

// BeforeCreate 创建前自动生成 UUID
func (l *PetLongTermMemory) BeforeCreate(tx *gorm.DB) error {
	if l.MemoryID == "" {
		l.MemoryID = uuid.New().String()
	}
	return nil
}

// MemoryCategory 常量
const (
	MemoryCategoryPreference   = "preference"   // 偏好记忆
	MemoryCategoryHabit        = "habit"        // 习惯记忆
	MemoryCategoryKnowledge    = "knowledge"    // 知识记忆
	MemoryCategoryRelationship = "relationship" // 关系记忆
)

// LongTermMemoryResponse 长期记忆响应
type LongTermMemoryResponse struct {
	MemoryID           string  `json:"memory_id"`
	DeviceID           string  `json:"device_id"`
	UserID             int64   `json:"user_id"`
	MemoryCategory     string  `json:"memory_category"`
	Content            string  `json:"content"`
	Keywords           string  `json:"keywords"`
	Embedding          string  `json:"embedding"`
	Confidence         float64 `json:"confidence"`
	ReinforcementCount int     `json:"reinforcement_count"`
	LastReinforcedAt   *string `json:"last_reinforced_at"`
	DecayScore         float64 `json:"decay_score"`
	IsLocked           bool    `json:"is_locked"`
	SourceMemoryID     string  `json:"source_memory_id"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (l *PetLongTermMemory) ToResponse() *LongTermMemoryResponse {
	resp := &LongTermMemoryResponse{
		MemoryID:           l.MemoryID,
		DeviceID:           l.DeviceID,
		UserID:             l.UserID,
		MemoryCategory:     l.MemoryCategory,
		Content:            l.Content,
		Keywords:           l.Keywords,
		Embedding:          l.Embedding,
		Confidence:         l.Confidence,
		ReinforcementCount: l.ReinforcementCount,
		DecayScore:         l.DecayScore,
		IsLocked:           l.IsLocked,
		SourceMemoryID:     l.SourceMemoryID,
		CreatedAt:          l.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          l.UpdatedAt.Format(time.RFC3339),
	}
	if l.LastReinforcedAt != nil {
		t := l.LastReinforcedAt.Format(time.RFC3339)
		resp.LastReinforcedAt = &t
	}
	return resp
}

// CreateLongTermMemoryRequest 创建长期记忆请求
type CreateLongTermMemoryRequest struct {
	MemoryCategory string                 `json:"memory_category" binding:"required"`
	Content        map[string]interface{} `json:"content" binding:"required"`
	Keywords       []string               `json:"keywords"`
	Confidence     float64                `json:"confidence"`
	IsLocked       bool                   `json:"is_locked"`
}

// LongTermMemoryListQuery 长期记忆列表查询参数
type LongTermMemoryListQuery struct {
	Category string `form:"category"`
	DeviceID string `form:"device_id"`
	Limit    int    `form:"limit,default=20"`
	Offset   int    `form:"offset,default=0"`
}
