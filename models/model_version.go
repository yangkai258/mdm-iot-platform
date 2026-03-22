package models

import (
	"time"

	"gorm.io/gorm"
)

// ModelVersionStatus 模型版本状态
type ModelVersionStatus string

const (
	VersionStatusDraft      ModelVersionStatus = "draft"
	VersionStatusTesting    ModelVersionStatus = "testing"
	VersionStatusStaging    ModelVersionStatus = "staging"
	VersionStatusProduction ModelVersionStatus = "production"
	VersionStatusDeprecated ModelVersionStatus = "deprecated"
)

// ModelVersion 模型版本（增强版，支持分片模型）
type ModelVersion struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	ModelID      uint             `json:"model_id" gorm:"index;not null"` // 关联 AIModelConfig.ID
	Version      string           `json:"version" gorm:"size:32;not null"`
	ModelName    string           `json:"model_name" gorm:"size:128"`
	Description  string           `json:"description" gorm:"type:text"`
	Status       ModelVersionStatus `json:"status" gorm:"size:20;default:'draft'"`
	ModelPath    string           `json:"model_path" gorm:"size:512"`
	Config       string           `json:"config" gorm:"type:jsonb"`
	Metrics      string           `json:"metrics" gorm:"type:jsonb"`
	TotalShards  int              `json:"total_shards" gorm:"default:0"`           // 总分片数
	VerifiedShards int            `json:"verified_shards" gorm:"default:0"`        // 已验证分片数
	IsSharded    bool             `json:"is_sarded" gorm:"default:false"`          // 是否分片模型
	RollbackFrom string           `json:"rollback_from" gorm:"size:32"`
	RollbackTo   string           `json:"rollback_to" gorm:"size:32"`
	PublishedBy  uint             `json:"published_by"`
	PublishedAt  *time.Time       `json:"published_at"`
	DeprecatedAt *time.Time       `json:"deprecated_at"`
	OrgID        uint             `gorm:"index" json:"org_id"`
	CreateUserID uint             `gorm:"index" json:"create_user_id"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName 表名
func (ModelVersion) TableName() string {
	return "ai_model_versions_v2"
}

// ModelVersionCreate 创建版本请求
type ModelVersionCreate struct {
	Version     string `json:"version" binding:"required" gorm:"size:32"`
	ModelName   string `json:"model_name" gorm:"size:128"`
	Description string `json:"description" gorm:"type:text"`
	ModelPath   string `json:"model_path" gorm:"size:512"`
	Config      string `json:"config" gorm:"type:jsonb"`
	Metrics     string `json:"metrics" gorm:"type:jsonb"`
	IsSharded  bool   `json:"is_sharded" gorm:"default:false"`
	TotalShards int   `json:"total_shards" gorm:"default:0"`
}

// ModelVersionUpdate 更新版本请求
type ModelVersionUpdate struct {
	ModelName   string `json:"model_name" gorm:"size:128"`
	Description string `json:"description" gorm:"type:text"`
	ModelPath   string `json:"model_path" gorm:"size:512"`
	Config      string `json:"config" gorm:"type:jsonb"`
	Metrics     string `json:"metrics" gorm:"type:jsonb"`
	TotalShards int    `json:"total_shards"`
	IsSharded  bool   `json:"is_sharded"`
}

// ModelVersionPublish 发布版本请求
type ModelVersionPublish struct {
	Version string `json:"version" binding:"required" gorm:"size:32"`
}

// ModelRollback 回滚请求
type ModelRollback struct {
	TargetVersion string `json:"target_version" binding:"required" gorm:"size:32"`
	Reason        string `json:"reason" gorm:"type:text"`
}

// DeployShardedModel 模型分片部署历史
type DeployShardedModel struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ModelID        uint           `json:"model_id" gorm:"index;not null"`
	VersionID      uint           `json:"version_id" gorm:"index;not null"`
	Version        string         `json:"version" gorm:"size:32;not null"`
	TargetEnv      string         `json:"target_env" gorm:"size:20;not null"` // staging / production
	ReplicaCount   int            `json:"replica_count" gorm:"default:1"`
	ResourceConfig string         `json:"resource_config" gorm:"type:jsonb"`
	Status         string         `json:"status" gorm:"size:20;default:'pending'"` // pending / deploying / online / failed / rolled_back
	ErrorMessage   string         `json:"error_message" gorm:"type:text"`
	StartedAt      time.Time      `json:"started_at"`
	CompletedAt    *time.Time     `json:"completed_at"`
	DeployedBy     uint           `json:"deployed_by"`
	OrgID          uint           `gorm:"index" json:"org_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// TableName 表名
func (DeployShardedModel) TableName() string {
	return "ai_deploy_sharded"
}
