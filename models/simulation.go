package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VirtualPet 虚拟宠物
type VirtualPet struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PetID        string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"pet_id"`
	Name         string         `gorm:"type:varchar(64);not null" json:"name"`
	Species      string         `gorm:"type:varchar(32);not null" json:"species"` // cat, dog, rabbit, etc.
	Personality  string         `gorm:"type:varchar(32);default:'lively'" json:"personality"`
	Mood         string         `gorm:"type:varchar(20);default:'happy'" json:"mood"` // happy, neutral, sad, excited
	Health       int            `gorm:"type:smallint;default:100" json:"health"`       // 0-100
	Hunger       int            `gorm:"type:smallint;default:80" json:"hunger"`        // 0-100
	Energy       int            `gorm:"type:smallint;default:100" json:"energy"`       // 0-100
	Happiness    int            `gorm:"type:smallint;default:80" json:"happiness"`     // 0-100
	Age          int            `gorm:"type:int;default:0" json:"age"`                 // 月龄
	Weight       float64        `gorm:"type:decimal(6,2);default:0" json:"weight"`       // 体重 kg
	AvatarURL    string         `gorm:"type:varchar(256)" json:"avatar_url"`
	CustomAttrs  string         `gorm:"type:text" json:"custom_attrs"` // JSON 扩展属性
	CreateUserID uint           `gorm:"index" json:"create_user_id"`
	OrgID        uint           `gorm:"index" json:"org_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (v *VirtualPet) BeforeCreate(tx *gorm.DB) error {
	if v.PetID == "" {
		v.PetID = uuid.New().String()
	}
	return nil
}

// SimulationEnvironment 测试环境
type SimulationEnvironment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	EnvID       string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"env_id"`
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	Description string         `gorm:"type:varchar(512)" json:"description"`
	SceneType   string         `gorm:"type:varchar(32);not null" json:"scene_type"` // indoor, outdoor, lab, field
	SceneConfig string         `gorm:"type:text" json:"scene_config"`               // JSON 场景配置
	Parameters  string         `gorm:"type:text" json:"parameters"`                 // JSON 环境参数
	Status      string         `gorm:"type:varchar(20);default:'idle'" json:"status"` // idle, running, paused, stopped
	Tags        string         `gorm:"type:varchar(256)" json:"tags"`              // 逗号分隔标签
	CreateUserID uint          `gorm:"index" json:"create_user_id"`
	OrgID       uint           `gorm:"index" json:"org_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (s *SimulationEnvironment) BeforeCreate(tx *gorm.DB) error {
	if s.EnvID == "" {
		s.EnvID = uuid.New().String()
	}
	return nil
}

// SimulationRun 测试运行记录
type SimulationRun struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RunID         string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"run_id"`
	Name          string         `gorm:"type:varchar(128);not null" json:"name"`
	EnvID         string         `gorm:"type:varchar(36);index" json:"env_id"`
	PetID         string         `gorm:"type:varchar(36);index" json:"pet_id"`
	ScenarioConfig string       `gorm:"type:text" json:"scenario_config"` // JSON 测试场景配置
	ResultData    string         `gorm:"type:text" json:"result_data"`    // JSON 测试结果
	Logs          string         `gorm:"type:longtext" json:"logs"`       // 测试日志
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, running, success, failed, cancelled
	Duration      int            `gorm:"type:int;default:0" json:"duration"` // 运行时长（秒）
	ErrorMsg      string         `gorm:"type:text" json:"error_msg"`
	MetricsSummary string        `gorm:"type:text" json:"metrics_summary"` // JSON 关键指标摘要
	CreateUserID  uint           `gorm:"index" json:"create_user_id"`
	OrgID         uint           `gorm:"index" json:"org_id"`
	StartedAt     *time.Time      `json:"started_at"`
	CompletedAt  *time.Time      `json:"completed_at"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (s *SimulationRun) BeforeCreate(tx *gorm.DB) error {
	if s.RunID == "" {
		s.RunID = uuid.New().String()
	}
	return nil
}

// SimulationMetrics 性能指标记录
type SimulationMetrics struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MetricID       string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"metric_id"`
	RunID          string         `gorm:"type:varchar(36);index" json:"run_id"`
	PetID          string         `gorm:"type:varchar(36);index" json:"pet_id"`
	EnvID          string         `gorm:"type:varchar(36);index" json:"env_id"`
	MetricType     string         `gorm:"type:varchar(32);not null" json:"metric_type"` // behavior, performance, stress, endurance
	MetricName     string         `gorm:"type:varchar(64);not null" json:"metric_name"`
	MetricValue    float64        `gorm:"type:decimal(12,4)" json:"metric_value"`
	Unit           string         `gorm:"type:varchar(16)" json:"unit"`
	Tags           string         `gorm:"type:varchar(256)" json:"tags"` // JSON 标签
	Snapshots      string         `gorm:"type:longtext" json:"snapshots"` // JSON 时序快照
	OrgID          uint           `gorm:"index" json:"org_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (s *SimulationMetrics) BeforeCreate(tx *gorm.DB) error {
	if s.MetricID == "" {
		s.MetricID = uuid.New().String()
	}
	return nil
}
