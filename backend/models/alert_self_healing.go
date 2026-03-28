package models

import (
	"time"

	"gorm.io/gorm"
)

// AlertSelfHealing 告警自愈建议
type AlertSelfHealing struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AlertType       string         `gorm:"type:varchar(50);not null;index" json:"alert_type"` // 告警类型
	AlertSubType    string         `gorm:"type:varchar(50)" json:"alert_sub_type"`          // 告警子类型
	Severity        int            `gorm:"default:2" json:"severity"`                       // 严重程度 1-5
	Title           string         `gorm:"type:varchar(200);not null" json:"title"`          // 问题标题
	RootCause       string         `gorm:"type:text" json:"root_cause"`                     // 根本原因
	Recommendation  string         `gorm:"type:text;not null" json:"recommendation"`        // 修复建议
	SelfHealingSteps []SelfHealingStep `gorm:"-" json:"self_healing_steps"`                // 自愈步骤（不存库）
	StepsJSON       string         `gorm:"type:text" json:"-"`                             // 自愈步骤JSON存储
	SuccessRate     float64        `gorm:"type:decimal(5,2);default:0" json:"success_rate"` // 历史成功率
	UsedCount       int            `gorm:"default:0" json:"used_count"`                    // 使用次数
	SuccessCount    int            `gorm:"default:0" json:"success_count"`                // 成功次数
	IsActive        bool           `gorm:"default:true" json:"is_active"`                  // 是否启用
	Tags            string         `gorm:"type:varchar(500)" json:"tags"`                  // 相关标签
	CreatedBy       string         `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (AlertSelfHealing) TableName() string {
	return "alert_self_healing"
}

// SelfHealingStep 自愈步骤
type SelfHealingStep struct {
	StepOrder int    `json:"step_order"` // 步骤序号
	Action    string `json:"action"`     // 操作描述
	Command   string `json:"command"`    // 执行的命令/API
	Timeout   int    `json:"timeout"`    // 超时时间(秒)
	Expected  string `json:"expected"`   // 预期结果
}

// AfterFind 加载后解析StepsJSON
func (a *AlertSelfHealing) AfterFind(tx *gorm.DB) error {
	if a.StepsJSON != "" {
		// 解析JSON到SelfHealingSteps
		// 这里简单处理，实际使用json.Unmarshal
	}
	return nil
}

// AlertSelfHealingRecord 自愈执行记录
type AlertSelfHealingRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AlertID         string         `gorm:"type:varchar(64);not null;index" json:"alert_id"`    // 关联告警ID
	SelfHealingID   uint           `gorm:"not null;index" json:"self_healing_id"`              // 使用的自愈方案ID
	AlertType       string         `gorm:"type:varchar(50)" json:"alert_type"`
	TriggerCondition string        `gorm:"type:text" json:"trigger_condition"`                // 触发条件
	StepsExecuted   int            `gorm:"default:0" json:"steps_executed"`                   // 已执行步骤数
	StepsTotal      int            `gorm:"default:0" json:"steps_total"`                     // 总步骤数
	Status          string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/running/success/failed/partial
	ExecutedBy      string         `gorm:"type:varchar(64)" json:"executed_by"`               // 执行者 (system/user)
	Result          string         `gorm:"type:text" json:"result"`                         // 执行结果
	ErrorMessage    string         `gorm:"type:text" json:"error_message"`                  // 错误信息
	StartedAt       *time.Time     `json:"started_at"`
	CompletedAt     *time.Time     `json:"completed_at"`
	Duration        int            `gorm:"default:0" json:"duration"`                        // 耗时(秒)
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (AlertSelfHealingRecord) TableName() string {
	return "alert_self_healing_records"
}
