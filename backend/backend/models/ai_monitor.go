package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AiModel AI模型
type AiModel struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	ModelUUID          string         `gorm:"size:64;uniqueIndex;not null" json:"model_uuid"`
	ModelName          string         `gorm:"size:128;not null" json:"model_name"`
	ModelType          string         `gorm:"size:50;not null;index" json:"model_type"` // vision, audio, multimodal, rl, nlp
	Description        string         `gorm:"type:text" json:"description"`
	CurrentVersionID   uint           `json:"current_version_id"`
	Status             string         `gorm:"size:20;default:'active'" json:"status"` // active, archived
	Tags               []string       `gorm:"type:text[]" json:"tags"`
	TenantID           string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AiModel) BeforeCreate(tx *gorm.DB) error {
	if a.ModelUUID == "" {
		a.ModelUUID = uuid.New().String()
	}
	return nil
}

// AiModelVersion AI模型版本
type AiModelVersion struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	VersionUUID     string         `gorm:"size:64;uniqueIndex;not null" json:"version_uuid"`
	ModelID         uint           `gorm:"not null;index" json:"model_id"`
	Version         string         `gorm:"size:20;not null" json:"version"`
	VersionTag      string         `gorm:"size:50" json:"version_tag"` // stable, beta, candidate
	Status          string         `gorm:"size:20;not null;index" json:"status"` // online, gray, offline, draft
	ModelFileURL    string         `gorm:"type:varchar(500)" json:"model_file_url"`
	ModelSizeMB     float64        `gorm:"type:numeric(10,2)" json:"model_size_mb"`
	Config          json.RawMessage `gorm:"type:jsonb" json:"config"`
	Metrics         json.RawMessage `gorm:"type:jsonb" json:"metrics"` // 评估指标
	TrainingTaskID  uint           `json:"training_task_id"`
	DeployedAt      *time.Time     `json:"deployed_at"`
	OfflineAt       *time.Time     `json:"offline_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AiModelVersion) BeforeCreate(tx *gorm.DB) error {
	if a.VersionUUID == "" {
		a.VersionUUID = uuid.New().String()
	}
	return nil
}

// AiMonitoringMetric AI监控指标
type AiMonitoringMetric struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MetricUUID     string         `gorm:"size:64;uniqueIndex;not null" json:"metric_uuid"`
	ModelVersionID uint           `gorm:"not null;index" json:"model_version_id"`
	MetricName     string         `gorm:"size:50;not null;index" json:"metric_name"` // latency, qps, accuracy, error_rate
	MetricValue    float64        `gorm:"type:numeric(15,5);not null" json:"metric_value"`
	MetricUnit     string         `gorm:"size:20" json:"metric_unit"`
	Percentile     string         `gorm:"size:10" json:"percentile"` // p50, p95, p99
	Tags           json.RawMessage `gorm:"type:jsonb" json:"tags"`
	RecordedAt     time.Time      `gorm:"not null;index" json:"recorded_at"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (a *AiMonitoringMetric) BeforeCreate(tx *gorm.DB) error {
	if a.MetricUUID == "" {
		a.MetricUUID = uuid.New().String()
	}
	return nil
}

// AiAlertRule AI告警规则
type AiAlertRule struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RuleUUID        string         `gorm:"size:64;uniqueIndex;not null" json:"rule_uuid"`
	RuleName        string         `gorm:"size:255;not null" json:"rule_name"`
	ModelID         uint           `gorm:"index" json:"model_id"`
	MetricName      string         `gorm:"size:50;not null" json:"metric_name"`
	Condition       string         `gorm:"size:20;not null" json:"condition"` // >, <, >=, <=, ==
	Threshold       float64        `gorm:"type:numeric(15,5);not null" json:"threshold"`
	Severity        string         `gorm:"size:20;not null;index" json:"severity"` // info, warning, critical
	NotifyWays      []string       `gorm:"type:text[]" json:"notify_ways"` // email, sms, webhook
	Enabled         bool           `gorm:"default:true" json:"enabled"`
	CooldownMinutes int            `gorm:"default:10" json:"cooldown_minutes"`
	Description     string         `gorm:"type:text" json:"description"`
	TenantID        string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AiAlertRule) BeforeCreate(tx *gorm.DB) error {
	if a.RuleUUID == "" {
		a.RuleUUID = uuid.New().String()
	}
	return nil
}
