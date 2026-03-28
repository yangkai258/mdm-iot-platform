package models

import (
	"time"

	"gorm.io/gorm"
)

// DeviceShadowSnapshot 设备影子快照
type DeviceShadowSnapshot struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DeviceID     string         `gorm:"type:varchar(64);not null;index" json:"device_id"`
	SnapshotID   string         `gorm:"type:varchar(64);uniqueIndex" json:"snapshot_id"` // 快照唯一标识

	// 快照版本信息
	Version      int            `gorm:"default:1" json:"version"`
	SnapshotType string         `gorm:"type:varchar(20);default:'manual'" json:"snapshot_type"` // manual/auto/scheduled
	Reason      string         `gorm:"type:varchar(200)" json:"reason"` // 创建原因

	// Desired State (云端期望)
	DesiredState string         `gorm:"type:text" json:"desired_state"` // JSON
	DesiredVersion int          `gorm:"default:0" json:"desired_version"`

	// Reported State (设备上报)
	ReportedState string        `gorm:"type:text" json:"reported_state"` // JSON
	ReportedVersion int          `gorm:"default:0" json:"reported_version"`

	// 元数据
	Metadata    string          `gorm:"type:text" json:"metadata"`     // JSON: 设备信息摘要
	Delta       string          `gorm:"type:text" json:"delta"`       // 与上一版本的差异
	Tags        string          `gorm:"type:varchar(500)" json:"tags"` // 快照标签

	// 状态对比
	StateDiff  int             `gorm:"default:0" json:"state_diff"`   // 差异字段数量
	IsHealthy   bool            `gorm:"default:true" json:"is_healthy"` // 状态是否健康

	// 存储信息
	FileURL     string          `gorm:"type:varchar(512)" json:"file_url"` // 快照文件URL
	FileSize    int64           `gorm:"default:0" json:"file_size"`    // 文件大小
	Checksum    string          `gorm:"type:varchar(64)" json:"checksum"` // 文件校验和

	// 时间信息
	ExpiresAt   *time.Time     `json:"expires_at"`                   // 过期时间
	CreatedBy   string          `gorm:"type:varchar(64)" json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (DeviceShadowSnapshot) TableName() string {
	return "device_shadow_snapshots"
}

// DeviceShadowSnapshotExport 快照导出记录
type DeviceShadowSnapshotExport struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SnapshotID   string         `gorm:"type:varchar(64);index" json:"snapshot_id"`
	DeviceID    string         `gorm:"type:varchar(64)" json:"device_id"`
	Format      string         `gorm:"type:varchar(20)" json:"format"` // json/csv/xml
	FileURL     string         `gorm:"type:varchar(512)" json:"file_url"`
	FileSize    int64          `gorm:"default:0" json:"file_size"`
	DownloadCount int           `gorm:"default:0" json:"download_count"`
	CreatedBy   string          `gorm:"type:varchar(64)" json:"created_by"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 表名
func (DeviceShadowSnapshotExport) TableName() string {
	return "device_shadow_snapshot_exports"
}
