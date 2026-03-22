package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AlertType 预警类型
const (
	AlertTypeVital       = "vital"
	AlertTypeBehavior    = "behavior"
	AlertTypeEnvironment = "environment"
	AlertTypeDisease     = "disease"
	AlertTypeReminder    = "reminder"
)

// AlertLevel 预警级别
const (
	AlertLevelInfo      = "info"
	AlertLevelWarning   = "warning"
	AlertLevelCritical  = "critical"
	AlertLevelEmergency = "emergency"
)

// AlertStatus 预警状态
const (
	AlertStatusActive   = "active"
	AlertStatusAcked    = "acked"
	AlertStatusResolved = "resolved"
	AlertStatusIgnored  = "ignored"
)

// HealthAlert 健康预警 (Sprint 18)
type HealthAlert struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AlertUUID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"alert_uuid"`
	PetUUID         string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID        string         `gorm:"type:varchar(64);index" json:"device_id"`
	AlertType       string         `gorm:"type:varchar(32);not null;index" json:"alert_type"`
	AlertLevel      string         `gorm:"type:varchar(16);not null;index" json:"alert_level"`
	Title           string         `gorm:"type:varchar(128);not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	TriggerValue    float64        `gorm:"type:decimal(10,2)" json:"trigger_value"`
	ThresholdValue  float64        `gorm:"type:decimal(10,2)" json:"threshold_value"`
	Unit            string         `gorm:"type:varchar(16)" json:"unit"`
	NormalRange     string         `gorm:"type:varchar(64)" json:"normal_range"`
	Suggestion      string         `gorm:"type:text" json:"suggestion"`
	Urgency         int            `gorm:"type:int;default:5" json:"urgency"`
	RelatedVitals   JSON           `gorm:"type:jsonb" json:"related_vitals"`
	RelatedBehaviors JSON          `gorm:"type:jsonb" json:"related_behaviors"`
	Status          string         `gorm:"type:varchar(16);default:'active';index" json:"status"`
	OccurredAt      time.Time      `gorm:"type:timestamp;not null;index" json:"occurred_at"`
	AcknowledgedAt  *time.Time     `json:"acknowledged_at"`
	AcknowledgedBy  *uint          `json:"acknowledged_by"`
	ResolvedAt      *time.Time     `json:"resolved_at"`
	ResolvedBy      *uint          `json:"resolved_by"`
	IgnoredAt       *time.Time     `json:"ignored_at"`
	IgnoredBy       *uint          `json:"ignored_by"`
	IgnoreReason    string         `gorm:"type:text" json:"ignore_reason"`
	NotifyChannels  StringArray    `gorm:"type:text[]" json:"notify_channels"`
	IsNotified      bool           `gorm:"type:boolean;default:false" json:"is_notified"`
	NotifiedAt      *time.Time     `json:"notified_at"`
	Tags            StringArray    `gorm:"type:text[]" json:"tags"`
	TenantID        string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HealthAlert) TableName() string {
	return "health_alerts"
}

func (h *HealthAlert) BeforeCreate(tx *gorm.DB) error {
	if h.AlertUUID == "" {
		h.AlertUUID = uuid.New().String()
	}
	return nil
}

// HealthAlertRule 健康预警规则 (可配置规则模板)
type HealthAlertRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	RuleUUID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"rule_uuid"`
	PetUUID        string         `gorm:"type:varchar(64);index" json:"pet_uuid"`
	RuleName       string         `gorm:"type:varchar(64);not null" json:"rule_name"`
	AlertType      string         `gorm:"type:varchar(32);not null" json:"alert_type"`
	Condition      JSON           `gorm:"type:jsonb;not null" json:"condition"`
	AlertLevel     string         `gorm:"type:varchar(16);not null" json:"alert_level"`
	TitleTemplate  string         `gorm:"type:varchar(128)" json:"title_template"`
	Suggestion     string         `gorm:"type:text" json:"suggestion"`
	CooldownPeriod int            `gorm:"type:int;default:3600" json:"cooldown_period"`
	IsEnabled      bool           `gorm:"type:boolean;default:true" json:"is_enabled"`
	Priority       int            `gorm:"type:int;default:5" json:"priority"`
	NotifyChannels StringArray    `gorm:"type:text[]" json:"notify_channels"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HealthAlertRule) TableName() string {
	return "health_alert_rules"
}

func (h *HealthAlertRule) BeforeCreate(tx *gorm.DB) error {
	if h.RuleUUID == "" {
		h.RuleUUID = uuid.New().String()
	}
	return nil
}

// ReqAlertQuery 预警查询请求
type ReqAlertQuery struct {
	AlertType  string `form:"alert_type" json:"alert_type"`
	AlertLevel string `form:"alert_level" json:"alert_level"`
	Status     string `form:"status" json:"status"`
	StartTime  string `form:"start_time" json:"start_time"`
	EndTime    string `form:"end_time" json:"end_time"`
	IsNotified *bool  `form:"is_notified" json:"is_notified"`
	Page       int    `form:"page" json:"page"`
	PageSize   int    `form:"page_size" json:"page_size"`
}

// RespAlertList 预警列表响应
type RespAlertList struct {
	Alerts     []HealthAlert `json:"alerts"`
	Total      int64         `json:"total"`
	UnackCount int64         `json:"unack_count"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

// ReqAckAlert 确认预警请求
type ReqAckAlert struct {
	Notes string `json:"notes"`
}

// ReqIgnoreAlert 忽略预警请求
type ReqIgnoreAlert struct {
	Reason string `json:"reason" binding:"required"`
}

// ReqCreateAlert 创建预警请求
type ReqCreateAlert struct {
	AlertType    string  `json:"alert_type" binding:"required"`
	AlertLevel   string  `json:"alert_level" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	TriggerValue float64 `json:"trigger_value"`
	Suggestion   string  `json:"suggestion"`
}
