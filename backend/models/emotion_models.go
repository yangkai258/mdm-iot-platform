package models

import (
	"time"
)

// EmotionRecord 情绪记录
type EmotionRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PetID       uint      `gorm:"index" json:"pet_id"`
	UserID      uint      `gorm:"index" json:"user_id"` // 如果是用户的情绪
	Source      string    `gorm:"type:varchar(32)" json:"source"`              // voice, text, pet_behavior
	EmotionType string    `gorm:"type:varchar(32)" json:"emotion_type"`        // happy, sad, angry, calm, excited, anxious
	Intensity   int       `gorm:"type:smallint" json:"intensity"`              // 1-10
	Trigger     string    `gorm:"type:varchar(256)" json:"trigger"`           // 触发原因
	Context     string    `gorm:"type:text" json:"context"`                   // 上下文
	AIResponse  string    `gorm:"type:text" json:"ai_response"`                // AI的响应动作
	RecordedAt  time.Time `gorm:"index" json:"recorded_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// PetEmotionAction 宠物情绪响应动作
type PetEmotionAction struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PetID       uint   `gorm:"index" json:"pet_id"`
	EmotionType string `gorm:"type:varchar(32)" json:"emotion_type"` // 情绪类型
	ActionType  string `gorm:"type:varchar(32)" json:"action_type"`  // gesture, sound, movement
	ActionData  string `gorm:"type:jsonb" json:"action_data"`       // JSON 动作数据
	Priority    int    `gorm:"type:int;default:0" json:"priority"`  // 优先级
	Enabled     bool   `gorm:"type:boolean;default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EmotionResponseConfig 情绪响应配置
type EmotionResponseConfig struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	PetID        uint   `gorm:"index" json:"pet_id"`
	EmotionType  string `gorm:"type:varchar(32)" json:"emotion_type"`
	ResponseMode string `gorm:"type:varchar(32)" json:"response_mode"` // comfort, encourage, play, quiet
	Actions      string `gorm:"type:jsonb" json:"actions"`             // JSON 配置的动作序列
	Enabled      bool   `gorm:"type:boolean;default:true" json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmotionReport 情绪报告
type EmotionReport struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PetID           uint      `gorm:"index" json:"pet_id"`
	UserID          uint      `gorm:"index" json:"user_id"`
	Period          string    `gorm:"type:varchar(16)" json:"period"`           // daily, weekly, monthly
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Summary         string    `gorm:"type:jsonb" json:"summary"`                // JSON 情绪摘要
	AvgIntensity    float64   `gorm:"type:float" json:"avg_intensity"`
	DominantEmotion string    `gorm:"type:varchar(32)" json:"dominant_emotion"`
	Trend           string    `gorm:"type:varchar(16)" json:"trend"`            // improving, stable, declining
	CreatedAt       time.Time `json:"created_at"`
}
