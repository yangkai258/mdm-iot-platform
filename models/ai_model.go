package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIProvider AI模型提供商
type AIProvider string

const (
	ProviderOpenAI   AIProvider = "openai"
	ProviderAnthropic AIProvider = "anthropic"
	ProviderLocal    AIProvider = "local"
	ProviderHuggingFace AIProvider = "huggingface"
	ProviderOllama   AIProvider = "ollama"
)

// AIModelStatus AI模型状态
type AIModelStatus string

const (
	ModelStatusPending   AIModelStatus = "pending"
	ModelStatusUploading AIModelStatus = "uploading"
	ModelStatusReady     AIModelStatus = "ready"
	ModelStatusDeploying AIModelStatus = "deploying"
	ModelStatusOnline    AIModelStatus = "online"
	ModelStatusOffline   AIModelStatus = "offline"
	ModelStatusFailed    AIModelStatus = "failed"
)

// AIModelConfig AI模型主表（用于Sprint 29 AI增强功能）
type AIModelConfig struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	ModelKey     string         `json:"model_key" gorm:"uniqueIndex;size:128;not null"` // 模型唯一标识，如 "gpt-4", "llama2-7b"
	Name         string         `json:"name" gorm:"size:128;not null"`                  // 模型显示名称
	Description  string         `json:"description" gorm:"type:text"`                   // 模型描述
	Provider     AIProvider     `json:"provider" gorm:"size:32;not null"`               // 提供商
	ModelType    string         `json:"model_type" gorm:"size:32;not null"`             // 模型类型：llm, cv, tts, stt, embedding
	ModelSize    string         `json:"model_size" gorm:"size:32"`                     // 模型大小，如 "7b", "13b", "70b"
	FilePath     string         `json:"file_path" gorm:"size:512"`                      // 模型文件路径（本地模型）
	FileSize     int64          `json:"file_size" gorm:"type:bigint"`                  // 文件大小（字节）
	Checksum     string         `json:"checksum" gorm:"size:64"`                       // 文件校验和 (SHA256)
	Config       string         `json:"config" gorm:"type:jsonb"`                     // 模型配置 (JSON)
	Capabilities string         `json:"capabilities" gorm:"type:jsonb"`                 // 模型能力列表
	QuotaDaily   int64          `json:"quota_daily" gorm:"default:0"`                  // 每日配额 (0=无限制)
	QuotaMonthly int64          `json:"quota_monthly" gorm:"default:0"`                // 每月配额
	PricePer1K   float64        `json:"price_per_1k" gorm:"type:decimal(10,4);default:0"` // 每1000 token价格
	Status       AIModelStatus  `json:"status" gorm:"size:20;default:'pending'"`        // 状态
	DeployedAt   *time.Time     `json:"deployed_at"`                                    // 部署时间
	DeployedBy   uint           `json:"deployed_by"`                                    // 部署人
	OrgID        uint           `gorm:"index" json:"org_id"`                           // 组织ID
	CreateUserID uint           `gorm:"index" json:"create_user_id"`                   // 创建人
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (a *AIModelConfig) BeforeCreate(tx *gorm.DB) error {
	if a.ModelKey == "" {
		a.ModelKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (AIModelConfig) TableName() string {
	return "ai_models"
}

// AIModelUploadRequest 上传模型请求
type AIModelUploadRequest struct {
	Name         string  `json:"name" binding:"required" gorm:"size:128"`
	Description string  `json:"description" gorm:"type:text"`
	Provider    string  `json:"provider" binding:"required" gorm:"size:32"`
	ModelType   string  `json:"model_type" binding:"required" gorm:"size:32"`
	ModelSize   string  `json:"model_size" gorm:"size:32"`
	FilePath    string  `json:"file_path" gorm:"size:512"`
	FileSize    int64  `json:"file_size" gorm:"type:bigint"`
	Checksum    string  `json:"checksum" gorm:"size:64"`
	Config      string  `json:"config" gorm:"type:jsonb"`
	Capabilities string `json:"capabilities" gorm:"type:jsonb"`
}

// AIModelUpdateRequest 更新模型请求
type AIModelUpdateRequest struct {
	Name         string `json:"name" gorm:"size:128"`
	Description  string `json:"description" gorm:"type:text"`
	ModelType    string `json:"model_type" gorm:"size:32"`
	Config       string `json:"config" gorm:"type:jsonb"`
	Capabilities string `json:"capabilities" gorm:"type:jsonb"`
	QuotaDaily   int64  `json:"quota_daily"`
	QuotaMonthly int64  `json:"quota_monthly"`
	PricePer1K   float64 `json:"price_per_1k" gorm:"type:decimal(10,4)"`
}

// AIModelDeployRequest 部署模型请求
type AIModelDeployRequest struct {
	TargetEnv string `json:"target_env" binding:"required" gorm:"size:20"` // staging, production
	ReplicaCount int `json:"replica_count" gorm:"default:1"`                // 副本数
	ResourceConfig string `json:"resource_config" gorm:"type:jsonb"`      // 资源配置
}

// AIModelDeployHistory 模型部署历史
type AIModelDeployHistory struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ModelID      uint      `json:"model_id" gorm:"index;not null"`
	TargetEnv    string    `json:"target_env" gorm:"size:20;not null"`    // staging, production
	ReplicaCount int       `json:"replica_count" gorm:"default:1"`
	ResourceConfig string  `json:"resource_config" gorm:"type:jsonb"`
	Status       string    `json:"status" gorm:"size:20;default:'pending'"` // pending, deploying, success, failed, rolled_back
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	ErrorMessage string    `json:"error_message" gorm:"type:text"`
	DeployedBy   uint      `json:"deployed_by"`
	RollbackFrom uint      `json:"rollback_from"`                           // 回滚到哪个历史版本
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 表名
func (AIModelDeployHistory) TableName() string {
	return "ai_model_deploy_history"
}
