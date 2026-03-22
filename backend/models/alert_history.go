package models

import (
	"time"
)

// AlertHistory 告警历史（归档已解决的告警）
type AlertHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OriginalID     uint      `gorm:"not null;index" json:"original_id"` // 原始告警ID
	RuleID         uint      `gorm:"index" json:"rule_id"`
	DeviceID       string    `gorm:"type:varchar(64);not null;index" json:"device_id"`
	AlertType      string    `gorm:"type:varchar(50);not null;index" json:"alert_type"`
	Severity       int       `gorm:"not null" json:"severity"` // 1:低 2:中 3:高 4:严重
	Message        string    `gorm:"type:text" json:"message"`
	TriggerValue   string    `gorm:"type:varchar(100)" json:"trigger_value"`
	Threshold      string    `gorm:"type:varchar(100)" json:"threshold"`
	Status         int       `gorm:"not null" json:"status"` // 1:未处理 2:已确认 3:已解决 4:已忽略
	NotifiedChannels string `gorm:"type:jsonb" json:"notified_channels"` // 通知渠道列表
	ConfirmedAt    *time.Time `json:"confirmed_at"`
	ConfirmedBy    uint       `json:"confirmed_by"`
	ResolvedAt     *time.Time `json:"resolved_at"`
	ResolvedBy     uint       `json:"resolved_by"`
	ResolveRemark  string     `gorm:"type:text" json:"resolve_remark"`
	CreatedAt      time.Time  `gorm:"not null;index" json:"created_at"`
	ResolvedAtH    *time.Time `json:"resolved_at_h"` // 解决时间（归档后）
	ArchivedAt     time.Time  `gorm:"default:NOW()" json:"archived_at"`
	TenantID       string    `gorm:"type:uuid;index" json:"tenant_id"`
}

func (AlertHistory) TableName() string {
	return "alert_history"
}

// AlertHistoryQuery 告警历史查询参数
type AlertHistoryQuery struct {
	DeviceID  string `form:"device_id"`
	AlertType string `form:"alert_type"`
	Severity  int    `form:"severity"`
	Status    int    `form:"status"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	RuleID    uint   `form:"rule_id"`
}
