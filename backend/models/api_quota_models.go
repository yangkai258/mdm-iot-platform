package models

import "time"

// APIQuota API配额
type APIQuota struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"uniqueIndex" json:"user_id"`
	PlanType     string    `json:"plan_type"`
	MonthlyQuota int64     `json:"monthly_quota"`
	UsedQuota    int64     `gorm:"default:0" json:"used_quota"`
	ResetAt      time.Time `json:"reset_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// APIUsageLog API使用日志
type APIUsageLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	APIKeyID   uint      `gorm:"index" json:"api_key_id"`
	Endpoint   string    `gorm:"index" json:"endpoint"`
	Method     string    `json:"method"`
	StatusCode int       `json:"status_code"`
	Latency    int       `json:"latency"`
	QuotaUsed  int       `json:"quota_used"`
	CreatedAt  time.Time `json:"created_at"`
}
