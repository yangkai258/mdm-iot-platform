package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Experiment AI行为实验
type Experiment struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ExpID        string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"exp_id"`
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Hypothesis   string         `gorm:"type:text" json:"hypothesis"`           // 实验假设
	TargetModel  string         `gorm:"type:varchar(64)" json:"target_model"`  // 目标AI模型版本
	Variables    string         `gorm:"type:jsonb" json:"variables"`           // 实验变量 JSON: [{"name":"temperature","type":"float","min":0,"max":40}]
	ControlGroup string         `gorm:"type:text" json:"control_group"`        // 对照组配置
	TestGroup    string         `gorm:"type:text" json:"test_group"`            // 实验组配置
	SampleSize   int            `gorm:"default:0" json:"sample_size"`           // 样本数量
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, running, paused, completed, cancelled
	StartTime    *time.Time     `gorm:"index" json:"start_time"`
	EndTime      *time.Time     `json:"end_time"`
	Duration     int            `gorm:"default:0" json:"duration"`             // 持续时间(秒)
	Metrics      string         `gorm:"type:jsonb" json:"metrics"`             // 实验指标 JSON: [{"name":"engagement","type":"score"}]
	Results      string         `gorm:"type:text" json:"results"`              // 实验结果摘要
	Conclusion   string         `gorm:"type:text" json:"conclusion"`          // 结论
	Tags         string         `gorm:"type:varchar(256)" json:"tags"`        // 标签，逗号分隔
	OrgID        uint           `gorm:"index" json:"org_id"`
	CreateUserID uint           `gorm:"index" json:"create_user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (e *Experiment) BeforeCreate(tx *gorm.DB) error {
	if e.ExpID == "" {
		e.ExpID = uuid.New().String()
	}
	return nil
}

// ExperimentRun 实验运行记录
type ExperimentRun struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RunID         string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"run_id"`
	ExpID         string         `gorm:"type:varchar(36);index;not null" json:"exp_id"`
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, running, completed, failed
	StartTime     *time.Time     `gorm:"index" json:"start_time"`
	EndTime       *time.Time     `json:"end_time"`
	Duration      int            `gorm:"default:0" json:"duration"`          // 持续时间(秒)
	ParticipantCount int         `gorm:"default:0" json:"participant_count"` // 参与数量
	ControlCount    int          `gorm:"default:0" json:"control_count"`      // 对照组数量
	TestCount       int          `gorm:"default:0" json:"test_count"`         // 实验组数量
	ControlMetrics string       `gorm:"type:jsonb" json:"control_metrics"`   // 对照组指标
	TestMetrics     string       `gorm:"type:jsonb" json:"test_metrics"`      // 实验组指标
	StatisticalSig string        `gorm:"type:varchar(32)" json:"statistical_sig"` // 统计显著性: significant, not_significant, unclear
	PValue         float64       `gorm:"type:decimal(6,4)" json:"p_value"`    // P值
	ErrorMsg       string        `gorm:"type:text" json:"error_msg"`
	OrgID          uint          `gorm:"index" json:"org_id"`
	CreateUserID   uint          `gorm:"index" json:"create_user_id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (e *ExperimentRun) BeforeCreate(tx *gorm.DB) error {
	if e.RunID == "" {
		e.RunID = uuid.New().String()
	}
	return nil
}

// ExperimentParticipant 实验参与者
type ExperimentParticipant struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ParticipantID string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"participant_id"`
	RunID         string         `gorm:"type:varchar(36);index;not null" json:"run_id"`
	ExpID         string         `gorm:"type:varchar(36);index;not null" json:"exp_id"`
	AnonymizedID  string         `gorm:"type:varchar(36);index" json:"anonymized_id"`  // 匿名化ID
	Group         string         `gorm:"type:varchar(16);not null" json:"group"`       // control, test
	Variables     string         `gorm:"type:jsonb" json:"variables"`                  // 各变量值
	Outcome       string         `gorm:"type:jsonb" json:"outcome"`                    // 结果指标
	EnrolledAt    time.Time      `gorm:"index" json:"enrolled_at"`
	CompletedAt   *time.Time      `json:"completed_at"`
	OrgID         uint           `gorm:"index" json:"org_id"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (e *ExperimentParticipant) BeforeCreate(tx *gorm.DB) error {
	if e.ParticipantID == "" {
		e.ParticipantID = uuid.New().String()
	}
	return nil
}
