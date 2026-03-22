package models

import (
	"time"

	"gorm.io/gorm"
)

// OTAPackage OTA固件包
// 注意：兼容 device.go 中已定义的 OTAPackage，扩展字段以匹配 PRD
type OTAPackage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	Version        string         `gorm:"type:varchar(32);not null" json:"version"`
	HardwareModel  string         `gorm:"type:varchar(64);not null;index" json:"hardware_model"`
	FileSize       int64          `gorm:"default:0" json:"file_size"`
	FileURL        string         `gorm:"type:varchar(512);not null" json:"file_url"`
	FileMD5        string         `gorm:"type:varchar(32)" json:"file_md5"`
	UploadSource   string         `gorm:"type:varchar(16);default:'local'" json:"upload_source"` // local / remote
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	IsMandatory    bool           `gorm:"default:false" json:"is_mandatory"`
	AllowDowngrade bool           `gorm:"default:false" json:"allow_downgrade"`
	ReleaseNotes   string         `gorm:"type:text" json:"release_notes"`
	CreatedBy      string         `gorm:"type:varchar(64);not null" json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OTAPackage) TableName() string {
	return "ota_packages"
}

// OTADeployment OTA部署任务
type OTADeployment struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	Name                     string         `gorm:"type:varchar(64);not null" json:"name"`
	PackageID                uint           `gorm:"not null" json:"package_id"`
	HardwareModel            string         `gorm:"type:varchar(64);not null" json:"hardware_model"`
	StrategyType             string         `gorm:"type:varchar(16);not null" json:"strategy_type"` // full / percentage / whitelist
	StrategyConfig           string         `gorm:"type:jsonb;default:'{}'" json:"strategy_config"` // JSON配置
	TargetDeviceCount        int            `gorm:"default:0" json:"target_device_count"`
	PendingCount             int            `gorm:"default:0" json:"pending_count"`
	RunningCount             int            `gorm:"default:0" json:"running_count"`
	SuccessCount             int            `gorm:"default:0" json:"success_count"`
	FailedCount              int            `gorm:"default:0" json:"failed_count"`
	Status                   string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/running/paused/completed/failed/cancelled
	PauseReason              string         `gorm:"type:varchar(256)" json:"pause_reason"`
	AutoPaused               bool           `gorm:"default:false" json:"auto_paused"`
	PauseOnFailureThreshold  float64        `gorm:"default:20.0" json:"pause_on_failure_threshold"`
	ScheduledAt              *time.Time     `json:"scheduled_at"`
	CancelledBy               string         `gorm:"type:varchar(64)" json:"cancelled_by"`
	CancelledAt              *time.Time     `json:"cancelled_at"`
	CompletedAt              *time.Time     `json:"completed_at"`
	CreatedBy                string         `gorm:"type:varchar(64);not null" json:"created_by"`
	CreateUserID             uint           `gorm:"index" json:"create_user_id"` // 创建人ID（用于数据权限）
	OrgID                    uint           `gorm:"index" json:"org_id"`          // 组织ID（用于数据权限）
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OTADeployment) TableName() string {
	return "ota_deployments"
}

// OTAProgress 设备OTA升级进度
type OTAProgress struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DeploymentID  uint           `gorm:"not null;index" json:"deployment_id"`
	DeviceID      string         `gorm:"type:varchar(64);not null;index" json:"device_id"`
	PackageID     uint           `gorm:"not null" json:"package_id"`
	FromVersion   string         `gorm:"type:varchar(32)" json:"from_version"`
	ToVersion     string         `gorm:"type:varchar(32);not null" json:"to_version"`
	OTAStatus     string         `gorm:"type:varchar(16);default:'pending'" json:"ota_status"` // pending/downloading/verifying/flashing/success/failed
	OTAMessage    string         `gorm:"type:varchar(256)" json:"ota_message"`
	Progress      int            `gorm:"default:0" json:"progress_percent"` // 0-100
	RetryCount    int            `gorm:"default:0" json:"retry_count"`
	StartedAt     *time.Time     `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OTAProgress) TableName() string {
	return "ota_progress"
}

// DeploymentPackageJoin 用于 JOIN 查询
type DeploymentPackageJoin struct {
	OTADeployment
	PackageVersion string `gorm:"type:varchar(32)" json:"package_version"`
}
