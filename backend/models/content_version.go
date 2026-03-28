package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentVersion 内容版本
type ContentVersion struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ContentID   uint           `gorm:"not null;index" json:"content_id"`   // content_files.id
	Version     string         `gorm:"type:varchar(32);not null" json:"version"` // 语义版本 如 "1.0.0"
	VersionNum   int            `gorm:"default:1" json:"version_num"`      // 数字版本号

	// 版本内容
	Title        string         `gorm:"type:varchar(200)" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	FileURL      string         `gorm:"type:varchar(512)" json:"file_url"`
	FileSize     int64          `gorm:"default:0" json:"file_size"`
	ContentHash  string         `gorm:"type:varchar(64)" json:"content_hash"` // 内容哈希
	ThumbnailURL string         `gorm:"type:varchar(512)" json:"thumbnail_url"`

	// 版本信息
	ChangeLog    string         `gorm:"type:text" json:"change_log"`       // 变更说明
	ChangeType   string         `gorm:"type:varchar(20)" json:"change_type"` // add/update/delete/fix
	ChangeSize   int64          `gorm:"default:0" json:"change_size"`      // 相对于上一版本的变化大小

	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/published/archived
	IsLatest     bool           `gorm:"default:false" json:"is_latest"`

	PublishedAt  *time.Time     `json:"published_at"`
	PublishedBy   string         `gorm:"type:varchar(64)" json:"published_by"`

	CreatedBy    string         `gorm:"type:varchar(64)" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ContentVersion) TableName() string {
	return "content_versions"
}

// ContentVersionReview 内容版本审核
type ContentVersionReview struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	VersionID    uint           `gorm:"not null;index" json:"version_id"`
	ReviewType   string         `gorm:"type:varchar(20)" json:"review_type"` // quality/compliance/security
	ReviewStatus string         `gorm:"type:varchar(20);default:'pending'" json:"review_status"` // pending/approved/rejected
	ReviewScore  float64        `gorm:"type:decimal(5,2)" json:"review_score"`  // 审核评分
	ReviewReport string         `gorm:"type:text" json:"review_report"`    // 审核报告
	Issues       string         `gorm:"type:text" json:"issues"`        // 发现的问题
	ReviewerID   string         `gorm:"type:varchar(64)" json:"reviewer_id"`
	ReviewerName string         `gorm:"type:varchar(100)" json:"reviewer_name"`
	ReviewedAt   *time.Time     `json:"reviewed_at"`
	CreatedAt    time.Time      `json:"created_at"`
}

// TableName 表名
func (ContentVersionReview) TableName() string {
	return "content_version_reviews"
}
