package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AiBehavior AI行为记录
type AiBehavior struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	BehaviorUUID  string         `gorm:"size:64;uniqueIndex;not null" json:"behavior_uuid"`
	DeviceID      string         `gorm:"size:64;index" json:"device_id"`
	PetUUID       string         `gorm:"size:64;index" json:"pet_uuid"`
	SessionID     string         `gorm:"size:64;index" json:"session_id"`
	ModelVersion  string         `gorm:"size:50" json:"model_version"`
	BehaviorType  string         `gorm:"size:50;not null;index" json:"behavior_type"` // action, emotion_response, decision, perception
	ActionCode    string         `gorm:"size:100" json:"action_code"`
	ActionName    string         `gorm:"size:128" json:"action_name"`
	TriggerEvent  string         `gorm:"size:128" json:"trigger_event"`
	InputData     json.RawMessage `gorm:"type:jsonb" json:"input_data"`
	OutputData    json.RawMessage `gorm:"type:jsonb" json:"output_data"`
	DecisionReason json.RawMessage `gorm:"type:jsonb" json:"decision_reason"`
	Confidence    float64        `gorm:"type:numeric(5,4)" json:"confidence"`
	InferenceMode string         `gorm:"size:20;index" json:"inference_mode"` // edge, cloud, hybrid
	LatencyMs     int64          `json:"latency_ms"`
	IsAnomaly     bool           `gorm:"default:false;index" json:"is_anomaly"`
	AnomalyScore  float64        `gorm:"type:numeric(5,4)" json:"anomaly_score"`
	AnomalyReason string         `gorm:"type:text" json:"anomaly_reason"`
	Location      json.RawMessage `gorm:"type:jsonb" json:"location"`
	Context       json.RawMessage `gorm:"type:jsonb" json:"context"`
	Tags          []string       `gorm:"type:text[]" json:"tags"`
	OccurredAt    time.Time      `gorm:"not null;index" json:"occurred_at"`
	TenantID      string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (a *AiBehavior) BeforeCreate(tx *gorm.DB) error {
	if a.BehaviorUUID == "" {
		a.BehaviorUUID = uuid.New().String()
	}
	return nil
}

// AiBehaviorStats AI行为统计
type AiBehaviorStats struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeviceID      string    `gorm:"size:64;index" json:"device_id"`
	PetUUID       string    `gorm:"size:64;index" json:"pet_uuid"`
	StatDate      time.Time `gorm:"type:date;not null;index" json:"stat_date"`
	BehaviorType  string    `gorm:"size:50;not null;index" json:"behavior_type"`
	TotalCount    int64     `json:"total_count"`
	AnomalyCount  int64     `json:"anomaly_count"`
	AvgConfidence float64   `gorm:"type:numeric(5,4)" json:"avg_confidence"`
	AvgLatencyMs  float64   `gorm:"type:numeric(10,2)" json:"avg_latency_ms"`
	EdgeCount     int64     `json:"edge_count"`
	CloudCount    int64     `json:"cloud_count"`
	TopActions    json.RawMessage `gorm:"type:jsonb" json:"top_actions"`
	TenantID      string    `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// AiAnomalyAlert AI异常告警
type AiAnomalyAlert struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	AlertUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"alert_uuid"`
	BehaviorID     uint           `gorm:"not null;index" json:"behavior_id"`
	DeviceID       string         `gorm:"size:64;index" json:"device_id"`
	PetUUID        string         `gorm:"size:64;index" json:"pet_uuid"`
	AnomalyType    string         `gorm:"size:50;not null;index" json:"anomaly_type"`
	AnomalyLevel   string         `gorm:"size:20;not null;index" json:"anomaly_level"` // low, medium, high, critical
	Description    string         `gorm:"type:text" json:"description"`
	TriggerData    json.RawMessage `gorm:"type:jsonb" json:"trigger_data"`
	Status         string         `gorm:"size:20;default:'active';index" json:"status"` // active, acknowledged, resolved, ignored
	OccurredAt     time.Time      `gorm:"not null;index" json:"occurred_at"`
	AcknowledgedAt *time.Time     `json:"acknowledged_at"`
	AcknowledgedBy string         `gorm:"size:64" json:"acknowledged_by"`
	ResolvedAt     *time.Time     `json:"resolved_at"`
	ResolvedBy     string         `gorm:"size:64" json:"resolved_by"`
	ResolveNote    string         `gorm:"type:text" json:"resolve_note"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (a *AiAnomalyAlert) BeforeCreate(tx *gorm.DB) error {
	if a.AlertUUID == "" {
		a.AlertUUID = uuid.New().String()
	}
	return nil
}
