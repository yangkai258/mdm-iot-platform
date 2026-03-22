package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TrainingStatus 训练任务状态
type TrainingStatus string

const (
	TrainingStatusPending   TrainingStatus = "pending"
	TrainingStatusQueued    TrainingStatus = "queued"
	TrainingStatusRunning   TrainingStatus = "running"
	TrainingStatusPaused    TrainingStatus = "paused"
	TrainingStatusCompleted TrainingStatus = "completed"
	TrainingStatusFailed    TrainingStatus = "failed"
	TrainingStatusCancelled TrainingStatus = "cancelled"
)

// TrainingPriority 训练优先级
type TrainingPriority int

const (
	PriorityLow    TrainingPriority = 1
	PriorityNormal TrainingPriority = 5
	PriorityHigh   TrainingPriority = 10
)

// AITraining 训练任务主表
type AITraining struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TaskKey        string         `json:"task_key" gorm:"uniqueIndex;size:64;not null"` // 任务唯一标识
	Name           string         `json:"name" gorm:"size:128;not null"`                // 任务名称
	Description    string         `json:"description" gorm:"type:text"`                // 任务描述
	ModelID        uint           `json:"model_id" gorm:"index;not null"`                // 关联的AI模型
	BaseModelKey   string         `json:"base_model_key" gorm:"size:128"`              // 基础模型key（如 llama2-7b）
	TrainingType   string         `json:"training_type" gorm:"size:32;not null"`         // 训练类型：fine_tune, rlhf, prompt_tune, domain_adapt
	DatasetPath    string         `json:"dataset_path" gorm:"size:512"`                 // 训练数据集路径
	DatasetSize    int64          `json:"dataset_size" gorm:"type:bigint"`              // 数据集大小
	ValidationPath string         `json:"validation_path" gorm:"size:512"`              // 验证集路径
	HyperParams    string         `json:"hyper_params" gorm:"type:jsonb"`              // 超参数配置
	ResourceConfig string         `json:"resource_config" gorm:"type:jsonb"`            // 资源配置：gpu, cpu, memory
	OutputPath     string         `json:"output_path" gorm:"size:512"`                  // 模型输出路径
	Status         TrainingStatus `json:"status" gorm:"size:20;default:'pending'"`
	Priority       TrainingPriority `json:"priority" gorm:"default:5"`
	Progress       int            `json:"progress" gorm:"type:smallint;default:0"`     // 进度 0-100
	Epoch          int            `json:"epoch" gorm:"default:0"`                      // 当前epoch
	TotalEpochs    int            `json:"total_epochs" gorm:"default:0"`              // 总epoch数
	CurrentStep    int64          `json:"current_step" gorm:"type:bigint;default:0"`    // 当前步数
	TotalSteps     int64          `json:"total_steps" gorm:"type:bigint;default:0"`    // 总步数
	Loss           float64        `json:"loss" gorm:"type:decimal(10,6)"`              // 当前loss
	Metrics        string         `json:"metrics" gorm:"type:jsonb"`                   // 训练指标
	LogsPath       string         `json:"logs_path" gorm:"size:512"`                   // 日志路径
	ErrorMessage   string         `json:"error_message" gorm:"type:text"`              // 错误信息
	StartedAt      *time.Time     `json:"started_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
	EstimatedTime  int            `json:"estimated_time" gorm:"type:integer"`          // 预估剩余时间（秒）
	QueuedAt       time.Time      `json:"queued_at"`
	OrgID          uint           `gorm:"index" json:"org_id"`
	CreateUserID   uint           `gorm:"index" json:"create_user_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (t *AITraining) BeforeCreate(tx *gorm.DB) error {
	if t.TaskKey == "" {
		t.TaskKey = uuid.New().String()
	}
	if t.QueuedAt.IsZero() {
		t.QueuedAt = time.Now()
	}
	return nil
}

// TableName 表名
func (AITraining) TableName() string {
	return "ai_training"
}

// AITrainingCreate 创建训练任务请求
type AITrainingCreate struct {
	Name           string `json:"name" binding:"required" gorm:"size:128"`
	Description    string `json:"description" gorm:"type:text"`
	ModelID        uint   `json:"model_id" binding:"required"`
	BaseModelKey   string `json:"base_model_key" gorm:"size:128"`
	TrainingType   string `json:"training_type" binding:"required" gorm:"size:32"`
	DatasetPath    string `json:"dataset_path" gorm:"size:512"`
	DatasetSize    int64  `json:"dataset_size" gorm:"type:bigint"`
	ValidationPath string `json:"validation_path" gorm:"size:512"`
	HyperParams    string `json:"hyper_params" gorm:"type:jsonb"`
	ResourceConfig string `json:"resource_config" gorm:"type:jsonb"`
	Priority       int    `json:"priority" gorm:"default:5"`
}

// AITrainingUpdate 更新训练任务请求
type AITrainingUpdate struct {
	Name           string `json:"name" gorm:"size:128"`
	Description    string `json:"description" gorm:"type:text"`
	HyperParams    string `json:"hyper_params" gorm:"type:jsonb"`
	ResourceConfig string `json:"resource_config" gorm:"type:jsonb"`
	Priority       int    `json:"priority"`
}

// AITrainingListQuery 训练任务列表查询
type AITrainingListQuery struct {
	ModelID      uint   `form:"model_id"`
	TrainingType string `form:"training_type"`
	Status       string `form:"status"`
	CreateUserID uint   `form:"create_user_id"`
}
