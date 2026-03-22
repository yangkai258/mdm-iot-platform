package models

import (
	"time"

	"gorm.io/gorm"
)

// ShardStatus 分片状态
type ShardStatus string

const (
	ShardStatusPending   ShardStatus = "pending"
	ShardStatusUploading ShardStatus = "uploading"
	ShardStatusUploaded  ShardStatus = "uploaded"
	ShardStatusVerified  ShardStatus = "verified"
	ShardStatusFailed    ShardStatus = "failed"
)

// ModelShard 模型分片
type ModelShard struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ModelID       uint           `json:"model_id" gorm:"index;not null"`
	Version       string         `json:"version" gorm:"size:32;index"`
	ShardIndex    int            `json:"shard_index" gorm:"not null"`                           // 分片序号 (0-based)
	ShardName     string         `json:"shard_name" gorm:"size:256;not null"`                    // 分片文件名
	FileSize      int64          `json:"file_size" gorm:"type:bigint;default:0"`                 // 分片大小（字节）
	FilePath      string         `json:"file_path" gorm:"size:512"`                            // 存储路径
	FileMD5       string         `json:"file_md5" gorm:"size:32"`                              // 分片 MD5
	FileSHA256    string         `json:"file_sha256" gorm:"size:64"`                            // 分片 SHA256
	Checksum      string         `json:"checksum" gorm:"size:128"`                               // 自定义校验值
	Status        ShardStatus    `json:"status" gorm:"size:20;default:'pending'"`
	VerifyMessage string         `json:"verify_message" gorm:"type:text"`                      // 验证消息
	VerifiedAt    *time.Time     `json:"verified_at"`
	VerifiedBy    uint           `json:"verified_by"`
	CreatedBy     uint           `json:"created_by"`
	OrgID         uint           `gorm:"index" json:"org_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ModelShard) TableName() string {
	return "model_shards"
}

// ModelShardCreate 创建分片请求
type ModelShardCreate struct {
	Version    string `json:"version" binding:"required" gorm:"size:32"`
	ShardIndex int    `json:"shard_index" binding:"required"`
	ShardName  string `json:"shard_name" binding:"required" gorm:"size:256"`
	FileSize   int64  `json:"file_size" gorm:"type:bigint"`
	FilePath   string `json:"file_path" gorm:"size:512"`
	FileMD5    string `json:"file_md5" gorm:"size:32"`
	FileSHA256 string `json:"file_sha256" gorm:"size:64"`
	Checksum   string `json:"checksum" gorm:"size:128"`
}

// ModelShardUpdate 更新分片请求
type ModelShardUpdate struct {
	FilePath   string `json:"file_path" gorm:"size:512"`
	FileMD5    string `json:"file_md5" gorm:"size:32"`
	FileSHA256 string `json:"file_sha256" gorm:"size:64"`
	Checksum   string `json:"checksum" gorm:"size:128"`
}

// ModelShardVerify 验证分片请求
type ModelShardVerify struct {
	Algorithm string `json:"algorithm" binding:"required"` // md5 / sha256 / checksum
	Expected  string `json:"expected"`
}

// DeployShardedRequest 部署分片模型请求
type DeployShardedRequest struct {
	Version         string `json:"version" binding:"required" gorm:"size:20"`
	TargetEnv      string `json:"target_env" binding:"required" gorm:"size:20"` // staging / production
	ReplicaCount   int    `json:"replica_count" gorm:"default:1"`
	ResourceConfig string `json:"resource_config" gorm:"type:jsonb"`
}
