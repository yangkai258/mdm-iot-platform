package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InferenceStatus 推理状态
type InferenceStatus string

const (
	InferenceStatusPending   InferenceStatus = "pending"
	InferenceStatusQueued    InferenceStatus = "queued"
	InferenceStatusRunning   InferenceStatus = "running"
	InferenceStatusCompleted InferenceStatus = "completed"
	InferenceStatusFailed    InferenceStatus = "failed"
)

// AIInference 推理任务主表
type AIInference struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	InferenceKey  string           `json:"inference_key" gorm:"uniqueIndex;size:64;not null"` // 推理唯一标识
	ModelID       uint             `json:"model_id" gorm:"index;not null"`                  // 使用的模型ID
	ModelKey      string           `json:"model_key" gorm:"size:128"`                        // 模型key（如 gpt-4）
	Mode          string           `json:"mode" gorm:"size:32;not null"`                     // 推理模式：chat, completion, embedding, image_gen
	InputTokens   int              `json:"input_tokens" gorm:"type:integer"`                // 输入token数
	OutputTokens  int              `json:"output_tokens" gorm:"type:integer"`                // 输出token数
	TotalTokens   int              `json:"total_tokens" gorm:"type:integer"`                // 总token数
	InputData     string           `json:"input_data" gorm:"type:jsonb"`                     // 输入数据
	OutputData    string           `json:"output_data" gorm:"type:jsonb"`                    // 输出数据
	Prompt        string           `json:"prompt" gorm:"type:text"`                         // 提示词
	Response      string           `json:"response" gorm:"type:text"`                       // 响应内容
	ErrorMessage  string           `json:"error_message" gorm:"type:text"`                  // 错误信息
	Status        InferenceStatus  `json:"status" gorm:"size:20;default:'pending'"`
	LatencyMs     int64            `json:"latency_ms"`                                      // 延迟（毫秒）
	Cost          float64          `json:"cost" gorm:"type:decimal(10,6)"`                // 推理成本
	StreamEnabled bool             `json:"stream_enabled" gorm:"default:false"`             // 是否启用流式
	Extra         string           `json:"extra" gorm:"type:jsonb"`                         // 扩展字段
	StartedAt     *time.Time       `json:"started_at"`
	CompletedAt   *time.Time       `json:"completed_at"`
	DeviceID      string          `json:"device_id" gorm:"size:36;index"`                  // 关联设备ID（可选）
	UserID        uint             `json:"user_id" gorm:"index"`                           // 用户ID
	OrgID         uint             `gorm:"index" json:"org_id"`
	CreateUserID  uint             `gorm:"index" json:"create_user_id"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (i *AIInference) BeforeCreate(tx *gorm.DB) error {
	if i.InferenceKey == "" {
		i.InferenceKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (AIInference) TableName() string {
	return "ai_inference"
}

// AIInferenceRequest 推理请求
type AIInferenceRequest struct {
	ModelID       uint   `json:"model_id" binding:"required"`                              // 模型ID
	Mode          string `json:"mode" binding:"required" gorm:"size:32"`                  // 推理模式
	Prompt        string `json:"prompt" gorm:"type:text"`                                 // 提示词
	InputData     string `json:"input_data" gorm:"type:jsonb"`                            // 结构化输入
	MaxTokens     int    `json:"max_tokens" gorm:"default:2048"`                          // 最大token数
	Temperature   float64 `json:"temperature" gorm:"default:0.7"`                         // 温度参数
	TopP          float64 `json:"top_p" gorm:"default:1.0"`                               // top_p参数
	StreamEnabled bool    `json:"stream_enabled" gorm:"default:false"`                    // 流式推理
	DeviceID      string  `json:"device_id" gorm:"size:36"`                                // 设备ID
	Extra         string  `json:"extra" gorm:"type:jsonb"`                                 // 扩展参数
}

// AIInferenceStreamRequest 流式推理请求
type AIInferenceStreamRequest struct {
	ModelID     uint    `json:"model_id" binding:"required"`
	Prompt      string  `json:"prompt" binding:"required" gorm:"type:text"`
	MaxTokens   int     `json:"max_tokens" gorm:"default:2048"`
	Temperature float64 `json:"temperature" gorm:"default:0.7"`
	DeviceID    string  `json:"device_id" gorm:"size:36"`
}

// AIInferenceResult 推理结果响应
type AIInferenceResult struct {
	InferenceKey string `json:"inference_key"`
	ModelID      uint   `json:"model_id"`
	ModelKey     string `json:"model_key"`
	Status       string `json:"status"`
	Response     string `json:"response"`
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
	TotalTokens  int    `json:"total_tokens"`
	LatencyMs    int64  `json:"latency_ms"`
	Cost         float64 `json:"cost"`
	CreatedAt    time.Time `json:"created_at"`
}

// AIInferenceListQuery 推理列表查询
type AIInferenceListQuery struct {
	ModelID   uint   `form:"model_id"`
	ModelKey  string `form:"model_key"`
	Status    string `form:"status"`
	DeviceID  string `form:"device_id"`
	UserID    uint   `form:"user_id"`
}
