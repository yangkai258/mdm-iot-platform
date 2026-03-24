package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// EmotionRecord 情感记录
type EmotionRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SubjectType  string         `gorm:"size:20;not null;index" json:"subject_type"` // pet, user, device
	SubjectID    int64          `gorm:"not null;index" json:"subject_id"`
	EmotionType  string         `gorm:"size:30;not null;index" json:"emotion_type"` // happy, sad, angry, calm, excited
	Intensity    float64        `gorm:"not null" json:"intensity"`                 // 0-100
	Source       string         `gorm:"size:20" json:"source"`                     // ai, manual, sensor
	Confidence   float64        `gorm:"default:0" json:"confidence"`               // 0-1
	Context      json.RawMessage `gorm:"type:jsonb" json:"context"`               // 额外上下文
	TriggerEvent string         `gorm:"type:text" json:"trigger_event"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	Note         string         `gorm:"type:text" json:"note"`
	RecordedAt   time.Time      `gorm:"not null;index" json:"recorded_at"`
	CreatedAt    time.Time      `json:"created_at"`
}

// EmotionReport 情感报告
type EmotionReport struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	PetID              int64          `gorm:"not null;index" json:"pet_id"`
	ReportType         string         `gorm:"size:20;not null" json:"report_type"` // daily, weekly, monthly
	StartDate          time.Time      `gorm:"not null" json:"start_date"`
	EndDate            time.Time      `gorm:"not null" json:"end_date"`
	Summary            json.RawMessage `gorm:"type:jsonb" json:"summary"`
	EmotionStats       json.RawMessage `gorm:"type:jsonb" json:"emotion_stats"`
	InteractionStats   json.RawMessage `gorm:"type:jsonb" json:"interaction_stats"`
	TrendAnalysis      json.RawMessage `gorm:"type:jsonb" json:"trend_analysis"`
	Recommendations    json.RawMessage `gorm:"type:jsonb" json:"recommendations"`
	GeneratedAt        time.Time      `gorm:"not null" json:"generated_at"`
	CreatedAt          time.Time      `json:"created_at"`
}

// EmotionResponseConfig 情感响应配置
type EmotionResponseConfig struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PetID          int64          `gorm:"not null;index" json:"pet_id"`
	EmotionType    string         `gorm:"size:30;not null" json:"emotion_type"`
	Strategy       string         `gorm:"size:20;not null" json:"strategy"` // comfort, reward, activity, ignore
	ActionCode     string         `gorm:"size:64" json:"action_code"`
	ActionParam    json.RawMessage `gorm:"type:jsonb" json:"action_param"`
	ResponseDelay  int64          `gorm:"default:0" json:"response_delay"` // milliseconds
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	Threshold      float64        `gorm:"default:30" json:"threshold"`
	Cooldown       int64          `gorm:"default:60000" json:"cooldown"` // milliseconds
	LastTriggered  *time.Time     `json:"last_triggered"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
