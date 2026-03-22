package models

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 企业应用商店模型 ====================

// StoreApp 企业应用商店应用
type StoreApp struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	BundleID     string         `gorm:"type:varchar(256);index;not null" json:"bundle_id"`
	Description  string         `gorm:"type:text" json:"description"`
	IconURL      string         `gorm:"type:varchar(512)" json:"icon_url"`
	Screenshots  string         `gorm:"type:text" json:"screenshots"`              // JSON数组，截图URL列表
	Category     string         `gorm:"type:varchar(64)" json:"category"`          // 应用分类
	Developer    string         `gorm:"type:varchar(128)" json:"developer"`        // 开发者名称
	DeveloperID  uint           `gorm:"index" json:"developer_id"`                 // 开发者用户ID
	Platform     string         `gorm:"type:varchar(32)" json:"platform"`          // ios/android/windows/mac/multi
	Price        float64        `gorm:"type:decimal(10,2);default:0" json:"price"` // 价格，0表示免费
	Currency     string         `gorm:"type:varchar(8);default:CNY" json:"currency"`
	Rating       float32        `gorm:"type:float;default:0" json:"rating"`           // 平均评分
	InstallCount int            `gorm:"default:0" json:"install_count"`               // 安装次数
	Status       int            `gorm:"type:smallint;default:0" json:"status"`        // 0:待审核 1:已发布 2:已下架 3:审核拒绝
	ReviewStatus int            `gorm:"type:smallint;default:0" json:"review_status"` // 审核状态：0:待审核 1:通过 2:拒绝
	ReviewNote   string         `gorm:"type:text" json:"review_note"`                 // 审核备注
	ReviewerID   uint           `gorm:"index" json:"reviewer_id"`                     // 审核人ID
	ReviewedAt   *time.Time     `json:"reviewed_at"`                                  // 审核时间
	PublishedAt  *time.Time     `json:"published_at"`                                 // 发布时间
	IsPublished  bool           `gorm:"default:false" json:"is_published"`            // 是否已发布
	Version      string         `gorm:"type:varchar(32)" json:"version"`              // 当前最新版本
	MinOSVersion string         `gorm:"type:varchar(32)" json:"min_os_version"`       // 最低系统版本
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (StoreApp) TableName() string {
	return "store_apps"
}

// StoreAppStatus 应用状态常量
const (
	StoreAppStatusPending   = 0 // 待审核
	StoreAppStatusPublished = 1 // 已发布
	StoreAppStatusOffline   = 2 // 已下架
	StoreAppStatusRejected  = 3 // 审核拒绝
)

// StoreAppReviewStatus 审核状态常量
const (
	StoreAppReviewPending  = 0 // 待审核
	StoreAppReviewApproved = 1 // 审核通过
	StoreAppReviewRejected = 2 // 审核拒绝
)

// StoreAppVersion 企业应用商店应用版本
type StoreAppVersion struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppID         uint           `gorm:"not null;index" json:"app_id"`
	Version       string         `gorm:"type:varchar(32);not null" json:"version"`
	BuildNumber   string         `gorm:"type:varchar(64)" json:"build_number"`
	FileSize      int64          `gorm:"default:0" json:"file_size"` // 文件大小(字节)
	FileURL       string         `gorm:"type:varchar(512)" json:"file_url"`
	FileMD5       string         `gorm:"type:varchar(32)" json:"file_md5"`
	MinOSVersion  string         `gorm:"type:varchar(32)" json:"min_os_version"`
	ReleaseNotes  string         `gorm:"type:text" json:"release_notes"`
	IsMandatory   bool           `gorm:"default:false" json:"is_mandatory"`     // 是否强制更新
	IsActive      bool           `gorm:"default:true" json:"is_active"`         // 是否激活
	IsLatest      bool           `gorm:"default:false" json:"is_latest"`        // 是否为最新版本
	DownloadCount int            `gorm:"default:0" json:"download_count"`       // 下载次数
	Status        int            `gorm:"type:smallint;default:0" json:"status"` // 0:待审核 1:已发布 2:已下架
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (StoreAppVersion) TableName() string {
	return "store_app_versions"
}

// StoreAppVersionStatus 版本状态常量
const (
	StoreAppVersionPending   = 0 // 待审核
	StoreAppVersionPublished = 1 // 已发布
	StoreAppVersionOffline   = 2 // 已下架
)

// StoreInstallation 应用安装记录
type StoreInstallation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AppID         uint           `gorm:"not null;index" json:"app_id"`
	AppVersionID  uint           `gorm:"index" json:"app_version_id"`                // 安装的应用版本ID
	DeviceID      uint           `gorm:"index" json:"device_id"`                     // 设备ID
	DeviceUUID    string         `gorm:"type:varchar(128);index" json:"device_uuid"` // 设备UUID
	UserID        uint           `gorm:"index" json:"user_id"`                       // 安装用户ID
	TenantID      uint           `gorm:"index" json:"tenant_id"`                     // 租户ID
	Status        int            `gorm:"type:smallint;default:0" json:"status"`      // 0:安装中 1:安装成功 2:安装失败 3:更新中 4:卸载中 5:已卸载
	Progress      int            `gorm:"type:smallint;default:0" json:"progress"`    // 安装进度 0-100
	ErrorMsg      string         `gorm:"type:text" json:"error_msg"`                 // 错误信息
	InstalledAt   *time.Time     `json:"installed_at"`                               // 安装成功时间
	UninstalledAt *time.Time     `json:"uninstalled_at"`                             // 卸载时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (StoreInstallation) TableName() string {
	return "store_installations"
}

// StoreInstallationStatus 安装状态常量
const (
	StoreInstallPending   = 0 // 安装中
	StoreInstallSuccess   = 1 // 安装成功
	StoreInstallFailed    = 2 // 安装失败
	StoreUpdatePending    = 3 // 更新中
	StoreUninstallPending = 4 // 卸载中
	StoreUninstalled      = 5 // 已卸载
)

// StoreReview 应用审核记录
type StoreReview struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AppID      uint           `gorm:"not null;index" json:"app_id"`
	VersionID  uint           `gorm:"index" json:"version_id"`               // 版本ID（可选，审核版本时填写）
	ReviewerID uint           `gorm:"index" json:"reviewer_id"`              // 审核人ID
	Status     int            `gorm:"type:smallint;default:0" json:"status"` // 0:待审核 1:通过 2:拒绝
	Action     string         `gorm:"type:varchar(32)" json:"action"`        // approve/reject/publish
	Reason     string         `gorm:"type:text" json:"reason"`               // 审核意见/拒绝原因
	ReviewedAt *time.Time     `json:"reviewed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (StoreReview) TableName() string {
	return "store_reviews"
}

// StoreReviewStatus 审核状态常量
const (
	StoreReviewPending  = 0 // 待审核
	StoreReviewApproved = 1 // 审核通过
	StoreReviewRejected = 2 // 审核拒绝
)
