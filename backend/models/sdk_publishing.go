package models

import (
	"time"

	"gorm.io/gorm"
)

// SDKPackage SDK包
type SDKPackage struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SDKID        string         `gorm:"type:varchar(64);uniqueIndex" json:"sdk_id"`

	// 基本信息
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	DisplayName string         `gorm:"type:varchar(200)" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"type:varchar(50)" json:"category"` // ios/android/web/hardware
	Platform    string         `gorm:"type:varchar(50)" json:"platform"`
	Language    string         `gorm:"type:varchar(30)" json:"language"` // swift/kotlin/typescript/go

	// 版本信息
	CurrentVersion string         `gorm:"type:varchar(50)" json:"current_version"`
	VersionCount  int            `gorm:"default:0" json:"version_count"`

	// 统计
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	InstallCount  int            `gorm:"default:0" json:"install_count"`
	StarCount    int            `gorm:"default:0" json:"star_count"`
	Rating       float64        `gorm:"type:decimal(3,2);default:0" json:"rating"` // 0-5

	// 开发者信息
	DeveloperID  string         `gorm:"type:varchar(64)" json:"developer_id"`
	DeveloperName string        `gorm:"type:varchar(100)" json:"developer_name"`

	// 状态
	Status      string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/published/archived/deprecated
	IsOfficial   bool           `gorm:"default:false" json:"is_official"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`

	// 标签
	Tags        string         `gorm:"type:varchar(500)" json:"tags"` // JSON数组
	CategoryTags string        `gorm:"type:varchar(200)" json:"category_tags"`

	// 文档
	DocURL      string         `gorm:"type:varchar(512)" json:"doc_url"`
	GitHubURL   string         `gorm:"type:varchar(512)" json:"github_url"`
	ChangelogURL string        `gorm:"type:varchar(512)" json:"changelog_url"`

	// Icon和截图
	IconURL     string         `gorm:"type:varchar(512)" json:"icon_url"`
	Screenshots string         `gorm:"type:text" json:"screenshots"` // JSON数组

	// 许可证
	License     string         `gorm:"type:varchar(50)" json:"license"` // MIT/Apache/GPL

	// 审核
	ReviewStatus string        `gorm:"type:varchar(20);default:'pending'" json:"review_status"` // pending/approved/rejected
	ReviewComment string        `gorm:"type:text" json:"review_comment"`
	ReviewedAt   *time.Time   `json:"reviewed_at"`
	ReviewedBy   string        `gorm:"type:varchar(64)" json:"reviewed_by"`

	CreatedBy   string         `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedAt   time.Time     `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (SDKPackage) TableName() string {
	return "sdk_packages"
}

// SDKVersion SDK版本
type SDKVersion struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SDKID        string         `gorm:"type:varchar(64);index" json:"sdk_id"`
	Version     string         `gorm:"type:varchar(50);not null" json:"version"`

	// 文件信息
	FileURL     string         `gorm:"type:varchar(512)" json:"file_url"`
	FileSize    int64          `gorm:"default:0" json:"file_size"`
	FileHash    string         `gorm:"type:varchar(64)" json:"file_hash"` // SHA256

	// 版本信息
	MinPlatformVersion string   `gorm:"type:varchar(50)" json:"min_platform_version"` // 最低平台版本要求
	ReleaseDate *time.Time     `json:"release_date"`
	ReleaseNotes string        `gorm:"type:text" json:"release_notes"`

	// 兼容性
	CompatiblePlatforms string `gorm:"type:varchar(200)" json:"compatible_platforms"` // JSON数组

	// 依赖
	Dependencies string        `gorm:"type:text" json:"dependencies"` // JSON: {"dep1": ">=1.0"}

	// 状态
	Status       string        `gorm:"type:varchar(20);default:'stable'" json:"status"` // stable/beta/alpha/deprecated
	IsRecommended bool          `gorm:"default:false" json:"is_recommended"`

	// 统计
	DownloadCount int          `gorm:"default:0" json:"download_count"`

	CreatedBy    string        `gorm:"type:varchar(64)" json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (SDKVersion) TableName() string {
	return "sdk_versions"
}

// SDKDownload SDK下载记录
type SDKDownload struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SDKID        string         `gorm:"type:varchar(64);index" json:"sdk_id"`
	VersionID   uint           `json:"version_id"`
	Version     string         `gorm:"type:varchar(50)" json:"version"`
	UserID     string         `gorm:"type:varchar(64);index" json:"user_id"`
	UserName   string         `gorm:"type:varchar(100)" json:"user_name"`
	ProjectID  string         `gorm:"type:varchar(64)" json:"project_id"`
	IPAddress  string         `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent  string         `gorm:"type:varchar(500)" json:"user_agent"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 表名
func (SDKDownload) TableName() string {
	return "sdk_downloads"
}
