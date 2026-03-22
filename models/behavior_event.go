package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BehaviorType 行为类型
const (
	BehaviorTypeEating    = "eating"
	BehaviorTypeDrinking  = "drinking"
	BehaviorTypeSleeping  = "sleeping"
	BehaviorTypePlaying   = "playing"
	BehaviorTypeWalking   = "walking"
	BehaviorTypeRunning   = "running"
	BehaviorTypeResting   = "resting"
	BehaviorTypeGrooming  = "grooming"
	BehaviorTypeSocial    = "social"
	BehaviorTypeCurious   = "curious"
	BehaviorTypeAlert     = "alert"
	BehaviorTypeAffection = "affection"
	BehaviorTypeBored     = "bored"
	BehaviorTypeAnxious   = "anxious"
)

// BehaviorEvent 行为事件 (Sprint 18)
type BehaviorEvent struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	EventUUID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"event_uuid"`
	PetUUID         string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID        string         `gorm:"type:varchar(64);index" json:"device_id"`
	BehaviorType    string         `gorm:"type:varchar(32);not null;index" json:"behavior_type"`
	BehaviorName    string         `gorm:"type:varchar(64)" json:"behavior_name"`
	StartTime       time.Time      `gorm:"type:timestamp;not null;index" json:"start_time"`
	EndTime         *time.Time     `json:"end_time"`
	Duration        int            `gorm:"type:int" json:"duration"`
	Intensity       float64        `gorm:"type:decimal(5,2)" json:"intensity"`
	Confidence      float64        `gorm:"type:decimal(5,2);default:1.0" json:"confidence"`
	Trigger         string         `gorm:"type:varchar(64)" json:"trigger"`
	Location        JSON           `gorm:"type:jsonb" json:"location"`
	Context         JSON           `gorm:"type:jsonb" json:"context"`
	AssociatedVitals JSON          `gorm:"type:jsonb" json:"associated_vitals"`
	IsAnomaly       bool           `gorm:"type:boolean;default:false;index" json:"is_anomaly"`
	AnomalyReason   string         `gorm:"type:text" json:"anomaly_reason"`
	Tags            StringArray    `gorm:"type:text[]" json:"tags"`
	VideoURL        string         `gorm:"type:varchar(512)" json:"video_url"`
	ThumbnailURL    string         `gorm:"type:varchar(512)" json:"thumbnail_url"`
	Notes           string         `gorm:"type:text" json:"notes"`
	TenantID        string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BehaviorEvent) TableName() string {
	return "behavior_events"
}

func (b *BehaviorEvent) BeforeCreate(tx *gorm.DB) error {
	if b.EventUUID == "" {
		b.EventUUID = uuid.New().String()
	}
	return nil
}

// BehaviorPrediction 行为预测 (Sprint 18)
type BehaviorPrediction struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	PredictionUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"prediction_uuid"`
	PetUUID           string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID          string         `gorm:"type:varchar(64);index" json:"device_id"`
	PredictionType    string         `gorm:"type:varchar(32);not null;index" json:"prediction_type"`
	PredictedBehavior string         `gorm:"type:varchar(32);index" json:"predicted_behavior"`
	Probability       float64        `gorm:"type:decimal(5,4)" json:"probability"`
	TimeWindowStart   time.Time      `gorm:"type:timestamp;index" json:"time_window_start"`
	TimeWindowEnd     time.Time      `gorm:"type:timestamp;index" json:"time_window_end"`
	Confidence        float64        `gorm:"type:decimal(5,4)" json:"confidence"`
	ModelVersion      string         `gorm:"type:varchar(32)" json:"model_version"`
	Factors           JSON           `gorm:"type:jsonb" json:"factors"`
	Recommendation    string         `gorm:"type:text" json:"recommendation"`
	Trigger           string         `gorm:"type:varchar(100)" json:"trigger"` // 触发条件
	IsVerified        bool           `gorm:"type:boolean;default:false" json:"is_verified"`
	VerifiedAt        *time.Time     `json:"verified_at"`
	ActualBehavior    string         `gorm:"type:varchar(32)" json:"actual_behavior"`
	TenantID          string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BehaviorPrediction) TableName() string {
	return "behavior_predictions"
}

func (b *BehaviorPrediction) BeforeCreate(tx *gorm.DB) error {
	if b.PredictionUUID == "" {
		b.PredictionUUID = uuid.New().String()
	}
	return nil
}

// ShortTermPrediction 短期动作预测
type ShortTermPrediction struct {
	Behavior    string  `json:"behavior"`
	Probability float64 `json:"probability"`
	TimeIn      int     `json:"time_in"`
	Duration    int     `json:"duration"`
	Trigger     string  `json:"trigger"`
}

// IntentPrediction 意图识别预测
type IntentPrediction struct {
	Intent      string  `json:"intent"`
	Probability float64 `json:"probability"`
	Confidence  float64 `json:"confidence"`
	Suggestion  string  `json:"suggestion"`
}

// BehaviorTimeline 时间轴事件
type BehaviorTimeline struct {
	ID           uint       `json:"id"`
	EventUUID    string     `json:"event_uuid"`
	BehaviorType string    `json:"behavior_type"`
	BehaviorName string    `json:"behavior_name"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Duration     int        `json:"duration"`
	Intensity    float64    `json:"intensity"`
	IsAnomaly    bool       `json:"is_anomaly"`
	VideoURL     string     `json:"video_url,omitempty"`
	ThumbnailURL string     `json:"thumbnail_url,omitempty"`
}

// RespBehaviorHistory 行为历史响应
type RespBehaviorHistory struct {
	Events     []BehaviorEvent `json:"events"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

// RespPredictions 行为预测响应
type RespPredictions struct {
	PetUUID       string              `json:"pet_uuid"`
	Predictions   []BehaviorPrediction `json:"predictions"`
	LastUpdatedAt time.Time           `json:"last_updated_at"`
}

// RespShortTermPrediction 短期预测响应
type RespShortTermPrediction struct {
	PetUUID      string               `json:"pet_uuid"`
	Predictions  []ShortTermPrediction `json:"predictions"`
	ModelVersion string               `json:"model_version"`
}

// RespIntentPrediction 意图识别响应
type RespIntentPrediction struct {
	PetUUID       string            `json:"pet_uuid"`
	CurrentIntent IntentPrediction  `json:"current_intent"`
	Alternatives  []IntentPrediction `json:"alternatives"`
}
