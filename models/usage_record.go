package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== 用量记录 =====

// UsageRecord 用量记录
type UsageRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	RecordID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"record_id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	UsageType   string         `gorm:"type:varchar(32);not null" json:"usage_type"` // api_call/device/storage/bandwidth
	UsageValue  float64       `gorm:"type:decimal(12,2);default:0" json:"usage_value"`
	Unit        string         `gorm:"type:varchar(16)" json:"unit"` // 次/GB/MB/个
	QuotaLimit  float64        `gorm:"type:decimal(12,2);default:0" json:"quota_limit"`
	QuotaUsed   float64        `gorm:"type:decimal(12,2);default:0" json:"quota_used"`
	RecordDate  time.Time      `gorm:"index" json:"record_date"`
	PeriodStart time.Time      `json:"period_start"`
	PeriodEnd   time.Time      `json:"period_end"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UsageRecord) TableName() string {
	return "usage_records"
}

// UsageType 常量
const (
	UsageTypeAPICall   = "api_call"   // API 调用
	UsageTypeDevice    = "device"     // 设备数
	UsageTypeStorage   = "storage"    // 存储空间
	UsageTypeBandwidth = "bandwidth"  // 带宽
)

// ===== 用户配额 =====

// UserQuota 用户配额
type UserQuota struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	QuotaID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"quota_id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	QuotaType  string         `gorm:"type:varchar(32);not null" json:"quota_type"` // api_call/device/storage/bandwidth
	QuotaLimit float64        `gorm:"type:decimal(12,2);default:0" json:"quota_limit"`
	QuotaUsed  float64        `gorm:"type:decimal(12,2);default:0" json:"quota_used"`
	Unit       string         `gorm:"type:varchar(16)" json:"unit"`
	PeriodType string         `gorm:"type:varchar(20);default:'monthly'" json:"period_type"` // monthly/yearly/daily
	PeriodStart time.Time     `json:"period_start"`
	PeriodEnd  time.Time      `json:"period_end"`
	UpdatedAt  time.Time      `json:"updated_at"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserQuota) TableName() string {
	return "user_quotas"
}
