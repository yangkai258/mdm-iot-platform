package models

import (
	"time"

	"gorm.io/gorm"
)

// DeviceHealthScore 设备健康评分
type DeviceHealthScore struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DeviceID    string         `gorm:"type:varchar(64);not null;index" json:"device_id"`

	// 总分和等级
	TotalScore  float64        `gorm:"type:decimal(5,2);default:0" json:"total_score"`  // 0-100
	Grade       string         `gorm:"type:varchar(1)" json:"grade"`                       // A/B/C/D/E

	// 各项得分
	UptimeScore   float64     `gorm:"type:decimal(5,2);default:100" json:"uptime_score"`
	PerfScore     float64     `gorm:"type:decimal(5,2);default:100" json:"perf_score"`
	SecurityScore float64     `gorm:"type:decimal(5,2);default:100" json:"security_score"`
	BehaviorScore float64     `gorm:"type:decimal(5,2);default:100" json:"behavior_score"`

	// 问题列表
	IssuesJSON    string       `gorm:"type:text" json:"-"`

	CalculatedAt time.Time      `json:"calculated_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 表名
func (DeviceHealthScore) TableName() string {
	return "device_health_scores"
}

// AlertDeduplicationRule 告警去重规则
type AlertDeduplicationRule struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RuleID      string         `gorm:"type:varchar(64);uniqueIndex" json:"rule_id"`

	// 规则条件
	AlertType    string        `gorm:"type:varchar(50);index" json:"alert_type"`
	DevicePattern string       `gorm:"type:varchar(200)" json:"device_pattern"` // 设备ID模式，支持通配符
	SeverityMin  int           `gorm:"default:1" json:"severity_min"`         // 最低严重程度
	SeverityMax  int           `gorm:"default:5" json:"severity_max"`         // 最高严重程度

	// 去重策略
	DedupWindowSeconds int      `gorm:"default:300" json:"dedup_window_seconds"` // 去重窗口(秒)
	DedupStrategy     string    `gorm:"type:varchar(20);default:'first'" json:"dedup_strategy"` // first/last/highest/lowest
	MaxCountPerWindow int      `gorm:"default:1" json:"max_count_per_window"`    // 窗口内最大告警数
	SuppressionType  string    `gorm:"type:varchar(20);default:'none'" json:"suppression_type"` // none/silence/merge

	// 状态
	IsActive     bool          `gorm:"default:true" json:"is_active"`
	Description string        `gorm:"type:text" json:"description"`

	CreatedBy   string        `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedAt  time.Time     `json:"updated_at"`
	CreatedAt  time.Time     `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (AlertDeduplicationRule) TableName() string {
	return "alert_deduplication_rules"
}

// AlertDeduplicationRecord 告警去重记录
type AlertDeduplicationRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RuleID      string         `gorm:"type:varchar(64);index" json:"rule_id"`
	AlertType   string         `gorm:"type:varchar(50)" json:"alert_type"`
	DeviceID    string         `gorm:"type:varchar(64);index" json:"device_id"`

	// 窗口信息
	WindowStart time.Time     `json:"window_start"`
	WindowEnd   time.Time     `json:"window_end"`
	AlertCount  int           `gorm:"default:1" json:"alert_count"`

	// 去重的告警摘要
	AlertSummary string        `gorm:"type:text" json:"alert_summary"` // JSON数组，存储去重掉的告警摘要

	// 最终处理的告警
	FinalAlertID string       `gorm:"type:varchar(64)" json:"final_alert_id"` // 最终保留的告警ID
	FinalAlertSnapshot string  `gorm:"type:text" json:"final_alert_snapshot"` // 最终告警的快照

	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (AlertDeduplicationRecord) TableName() string {
	return "alert_deduplication_records"
}
