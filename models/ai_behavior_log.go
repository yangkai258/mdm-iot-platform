package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIBehaviorLog AI行为日志
type AIBehaviorLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	LogID         string    `json:"log_id" gorm:"uniqueIndex;size:64"`
	DeviceID      string    `json:"device_id" gorm:"index;size:64"`
	UserID        uint      `json:"user_id" gorm:"index"`
	ModelID       string    `json:"model_id" gorm:"size:64"`          // 模型ID
	ModelName     string    `json:"model_name" gorm:"size:64"`
	ModelVersion  string    `json:"model_version" gorm:"size:32"`
	EventType     string    `json:"event_type" gorm:"size:32"` // inference/abnormal/rollback
	InputSummary  string    `json:"input_summary" gorm:"type:text"`
	OutputSummary string    `json:"output_summary" gorm:"type:text"`
	LatencyMs     int       `json:"latency_ms"`
	ErrorRate     float64   `json:"error_rate"`
	Confidence    float64   `json:"confidence"`
	Status        string    `json:"status" gorm:"size:20"` // success/failed/anomaly
	ErrorMsg      string    `json:"error_msg" gorm:"size:512"`
	Metadata      string    `json:"metadata" gorm:"type:jsonb"`
	CreatedAt     time.Time `json:"created_at" gorm:"index"`
}

// TableName 表名
func (AIBehaviorLog) TableName() string {
	return "ai_behavior_logs"
}

// BeforeCreate 创建前自动生成 LogID
func (a *AIBehaviorLog) BeforeCreate(tx *gorm.DB) error {
	if a.LogID == "" {
		a.LogID = uuid.New().String()
	}
	return nil
}

// AIBehaviorEvent AI行为事件（上报用）
type AIBehaviorEvent struct {
	EventType  string  `json:"event_type" binding:"required"`
	DeviceID   string  `json:"device_id"`
	ModelName  string  `json:"model_name" binding:"required"`
	InputData  string  `json:"input_data"`
	OutputData string  `json:"output_data"`
	LatencyMs  int     `json:"latency_ms"`
	Confidence float64 `json:"confidence"`
	IsAnomaly  bool    `json:"is_anomaly"`
}

// ToAIBehaviorLog 将事件转为日志
func (e *AIBehaviorEvent) ToAIBehaviorLog() *AIBehaviorLog {
	status := "success"
	if e.IsAnomaly {
		status = "anomaly"
	}
	return &AIBehaviorLog{
		LogID:         uuid.New().String(),
		DeviceID:      e.DeviceID,
		ModelName:     e.ModelName,
		EventType:     e.EventType,
		InputSummary:  e.InputData,
		OutputSummary: e.OutputData,
		LatencyMs:     e.LatencyMs,
		Confidence:    e.Confidence,
		Status:        status,
	}
}

// AIRollbackTask AI回滚任务
type AIRollbackTask struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	TaskID        string     `json:"task_id" gorm:"uniqueIndex;size:64"`
	ModelID       string     `json:"model_id" gorm:"index;size:64"`
	FromVersion   string     `json:"from_version" gorm:"size:32"`
	ToVersion     string     `json:"to_version" gorm:"size:32"`
	Reason        string     `json:"reason" gorm:"type:text"`
	Status        string     `json:"status" gorm:"size:20"` // pending/in_progress/completed/failed
	TriggeredBy   uint       `json:"triggered_by"`
	AffectedCount int        `json:"affected_count"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// TableName 表名
func (AIRollbackTask) TableName() string {
	return "ai_rollback_tasks"
}
