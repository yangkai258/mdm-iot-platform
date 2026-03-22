package models

import (
	"time"
)

// EmotionType 情绪类型枚举
const (
	EmotionTypeHappy   = "happy"
	EmotionTypeSad     = "sad"
	EmotionTypeAngry   = "angry"
	EmotionTypeCalm    = "calm"
	EmotionTypeAnxious = "anxious"
	EmotionTypeLonely  = "lonely"
	EmotionTypeTired   = "tired"
)

// EmotionSource 情绪来源枚举
const (
	EmotionSourceVoice    = "voice"
	EmotionSourceText     = "text"
	EmotionSourceFace     = "face"
	EmotionSourceBehavior = "behavior"
)

// SubjectType 被试主体类型
const (
	SubjectTypeUser = "user"
	SubjectTypePet  = "pet"
)

// EmotionRecord 情绪记录
type EmotionRecord struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SubjectType  string     `gorm:"type:varchar(20);not null;index" json:"subject_type"`
	SubjectID    uint       `gorm:"not null;index" json:"subject_id"`
	EmotionType  string     `gorm:"type:varchar(30);not null;index" json:"emotion_type"`
	Intensity    float64    `gorm:"not null" json:"intensity"`
	Source       string     `gorm:"type:varchar(20)" json:"source"`
	Confidence   float64    `gorm:"default:0" json:"confidence"`
	Context      JSON       `gorm:"type:jsonb" json:"context"`
	TriggerEvent string     `gorm:"type:text" json:"trigger_event"`
	Tags         StringArray `gorm:"type:text[]" json:"tags"`
	Note         string     `gorm:"type:text" json:"note"`
	RecordedAt   time.Time  `gorm:"not null;index" json:"recorded_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (EmotionRecord) TableName() string {
	return "emotion_records"
}

// PetEmotionAction 宠物情绪动作
type PetEmotionAction struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmotionType  string     `gorm:"type:varchar(30);not null;index" json:"emotion_type"`
	ActionName   string     `gorm:"type:varchar(128);not null" json:"action_name"`
	ActionCode   string     `gorm:"type:varchar(64);not null;uniqueIndex" json:"action_code"`
	Description  string     `gorm:"type:text" json:"description"`
	Parameters   JSON       `gorm:"type:jsonb" json:"parameters"`
	Priority     int        `gorm:"type:int;default:0" json:"priority"`
	MinIntensity float64    `gorm:"default:0" json:"min_intensity"`
	MaxIntensity float64    `gorm:"default:100" json:"max_intensity"`
	Duration     int        `gorm:"type:int;default:0" json:"duration"`
	Enabled      bool       `gorm:"default:true" json:"enabled"`
	Icon         string     `gorm:"type:varchar(64)" json:"icon"`
	AnimationURL string     `gorm:"type:varchar(255)" json:"animation_url"`
	SoundURL     string     `gorm:"type:varchar(255)" json:"sound_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (PetEmotionAction) TableName() string {
	return "pet_emotion_actions"
}

// EmotionResponseConfig 情绪响应配置
type EmotionResponseConfig struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PetID         uint       `gorm:"not null;index" json:"pet_id"`
	EmotionType   string     `gorm:"type:varchar(30);not null" json:"emotion_type"`
	Strategy      string     `gorm:"type:varchar(20);not null" json:"strategy"`
	ActionCode    string     `gorm:"type:varchar(64)" json:"action_code"`
	ActionParam   JSON       `gorm:"type:jsonb" json:"action_param"`
	ResponseDelay int        `gorm:"type:int;default:0" json:"response_delay"`
	Enabled       bool       `gorm:"default:true" json:"enabled"`
	Threshold     float64    `gorm:"default:30" json:"threshold"`
	Cooldown      int        `gorm:"type:int;default:60000" json:"cooldown"`
	LastTriggered *time.Time `json:"last_triggered"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (EmotionResponseConfig) TableName() string {
	return "emotion_response_configs"
}

// EmotionReport 情绪报告
type EmotionReport struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	PetID            uint      `gorm:"not null;index" json:"pet_id"`
	ReportType       string    `gorm:"type:varchar(20);not null" json:"report_type"`
	StartDate        time.Time `gorm:"not null" json:"start_date"`
	EndDate          time.Time `gorm:"not null" json:"end_date"`
	Summary          JSON      `gorm:"type:jsonb" json:"summary"`
	EmotionStats     JSON      `gorm:"type:jsonb" json:"emotion_stats"`
	InteractionStats JSON      `gorm:"type:jsonb" json:"interaction_stats"`
	TrendAnalysis    JSON      `gorm:"type:jsonb" json:"trend_analysis"`
	Recommendations  JSON      `gorm:"type:jsonb" json:"recommendations"`
	GeneratedAt      time.Time `gorm:"not null" json:"generated_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (EmotionReport) TableName() string {
	return "emotion_reports"
}
