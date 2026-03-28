package models

import (
	"time"

	"gorm.io/gorm"
)

// KnowledgeVersion 知识库版本
type KnowledgeVersion struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	KnowledgeID  uint           `gorm:"not null;index" json:"knowledge_id"` // 关联知识ID
	Version      string         `gorm:"type:varchar(32);not null" json:"version"` // 版本号
	Content      string         `gorm:"type:text" json:"content"`         // 内容快照
	ChangeLog    string         `gorm:"type:text" json:"change_log"`      // 变更说明
	ContentHash  string         `gorm:"type:varchar(64)" json:"content_hash"` // 内容哈希
	FileURL      string         `gorm:"type:varchar(512)" json:"file_url"` // 附件URL
	FileSize     int64          `gorm:"default:0" json:"file_size"`     // 文件大小
	ChangeType   string         `gorm:"type:varchar(20)" json:"change_type"` // add/update/delete
	ChangedFields string        `gorm:"type:text" json:"changed_fields"`  // 变更字段
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/published/archived
	PublishedAt  *time.Time    `json:"published_at"`                  // 发布时间
	PublishedBy  string         `gorm:"type:varchar(64)" json:"published_by"` // 发布人
	CreatedBy    string         `gorm:"type:varchar(64)" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (KnowledgeVersion) TableName() string {
	return "knowledge_versions"
}

// KnowledgeVersionReview 版本审核记录
type KnowledgeVersionReview struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	VersionID     uint           `gorm:"not null;index" json:"version_id"`
	ReviewStatus  string         `gorm:"type:varchar(20);default:'pending'" json:"review_status"` // pending/approved/rejected
	ReviewComment string         `gorm:"type:text" json:"review_comment"` // 审核意见
	ReviewerID    string         `gorm:"type:varchar(64)" json:"reviewer_id"`
	ReviewerName  string         `gorm:"type:varchar(100)" json:"reviewer_name"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (KnowledgeVersionReview) TableName() string {
	return "knowledge_version_reviews"
}
