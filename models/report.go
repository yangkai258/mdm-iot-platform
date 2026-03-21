package models

import (
	"time"

	"gorm.io/gorm"
)

// ReportRecord 报表记录（用于存储生成的报表元数据）
type ReportRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	ReportType  string         `gorm:"type:varchar(20);not null" json:"report_type"` // daily, weekly, monthly
	ReportDate  string         `gorm:"type:varchar(10);not null" json:"report_date"`  // YYYY-MM-DD
	DeviceStats string         `gorm:"type:jsonb" json:"device_stats"`               // 设备统计数据
	MemberStats string         `gorm:"type:jsonb" json:"member_stats"`               // 会员统计数据
	OTAStats    string         `gorm:"type:jsonb" json:"ota_stats"`                  // OTA 统计数据
	FilePath    string         `gorm:"type:varchar(512)" json:"file_path"`            // 导出文件路径
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ReportRecord) TableName() string {
	return "report_records"
}

// DeviceStats 设备统计数据
type DeviceStats struct {
	Total     int64 `json:"total"`
	Online    int64 `json:"online"`
	Offline   int64 `json:"offline"`
	Alerting  int64 `json:"alerting"`
	Active    int64 `json:"active"`
	Inactive  int64 `json:"inactive"`
	NewToday  int64 `json:"new_today"`
	NewWeek   int64 `json:"new_week"`
	NewMonth  int64 `json:"new_month"`
}

// MemberStats 会员统计数据
type MemberStats struct {
	Total     int64 `json:"total"`
	Active    int64 `json:"active"`    // 30天内有活动的
	NewToday  int64 `json:"new_today"`
	NewWeek   int64 `json:"new_week"`
	NewMonth  int64 `json:"new_month"`
	ByLevel   []LevelCount `json:"by_level"`
}

// LevelCount 等级分布
type LevelCount struct {
	Level int    `json:"level"`
	Count int64  `json:"count"`
}

// OTAStats OTA 统计数据
type OTAStats struct {
	TotalDeployments   int64          `json:"total_deployments"`
	SuccessCount       int64          `json:"success_count"`
	FailedCount        int64          `json:"failed_count"`
	PendingCount       int64          `json:"pending_count"`
	SuccessRate        float64        `json:"success_rate"`         // 百分比
	VersionDistribution []VersionCount `json:"version_distribution"` // 固件版本分布
}

// VersionCount 版本分布
type VersionCount struct {
	Version string `json:"version"`
	Count   int64  `json:"count"`
}

// ReportRequest 报表查询请求
type ReportRequest struct {
	ReportType string `form:"report_type" json:"report_type"` // daily, weekly, monthly
	StartDate  string `form:"start_date" json:"start_date"`   // YYYY-MM-DD
	EndDate    string `form:"end_date" json:"end_date"`       // YYYY-MM-DD
	ExportType string `form:"export_type" json:"export_type"` // csv, excel
}

// ReportResponse 报表响应
type ReportResponse struct {
	ReportType string       `json:"report_type"`
	ReportDate string       `json:"report_date"`
	Device     DeviceStats  `json:"device"`
	Member     MemberStats  `json:"member"`
	OTA        OTAStats     `json:"ota"`
}

// DailyTrend 每日趋势数据
type DailyTrend struct {
	Date          string `json:"date"`
	DeviceTotal   int64  `json:"device_total"`
	DeviceNew     int64  `json:"device_new"`
	MemberTotal   int64  `json:"member_total"`
	MemberNew     int64  `json:"member_new"`
	OTAComplete   int64  `json:"ota_complete"`
	OTASuccessRate float64 `json:"ota_success_rate"`
}
