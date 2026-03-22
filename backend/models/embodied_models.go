package models

import (
	"time"

	"gorm.io/gorm"
)

// EmbodiedMap 环境地图表
type EmbodiedMap struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeviceID    string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	MapType     string         `gorm:"type:varchar(20);not null" json:"map_type"` // grid/semantic/topological
	MapData     string         `gorm:"type:jsonb;not null" json:"map_data"`        // 地图数据
	Resolution  float64        `gorm:"type:decimal(10,3)" json:"resolution"`     // 栅格分辨率
	Size        string         `gorm:"type:jsonb" json:"size"`                    // 地图尺寸
	Version     int            `gorm:"type:int;default:1" json:"version"`
	IsActive    bool           `gorm:"type:boolean;default:false" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (EmbodiedMap) TableName() string {
	return "embodied_maps"
}

// SpatialPosition 空间位置表
type SpatialPosition struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DeviceID   string    `gorm:"type:varchar(100);not null;index" json:"device_id"`
	MapID      *uint     `gorm:"index" json:"map_id"`
	PositionX  float64   `gorm:"type:decimal(10,3)" json:"position_x"`
	PositionY  float64   `gorm:"type:decimal(10,3)" json:"position_y"`
	PositionZ  float64   `gorm:"type:decimal(10,3)" json:"position_z"`
	Orientation float64  `gorm:"type:decimal(10,3)" json:"orientation"`
	Confidence float64   `gorm:"type:decimal(5,4)" json:"confidence"`
	RecordedAt time.Time `gorm:"autoUpdateTime" json:"recorded_at"`
}

// TableName 指定表名
func (SpatialPosition) TableName() string {
	return "spatial_positions"
}

// ActionExecution 动作执行记录表
type ActionExecution struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	DeviceID           string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	ActionID           uint           `gorm:"index" json:"action_id"`
	ExecutionType      string         `gorm:"type:varchar(20)" json:"execution_type"`  // triggered/scheduled/manual
	StartTime          time.Time      `gorm:"not null" json:"start_time"`
	EndTime            *time.Time     `json:"end_time"`
	Status             string         `gorm:"type:varchar(20)" json:"status"`          // running/completed/interrupted/failed
	Parameters         string         `gorm:"type:jsonb" json:"parameters"`            // 执行参数
	InterruptionReason string         `gorm:"type:text" json:"interruption_reason"`
	CreatedAt          time.Time      `json:"created_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ActionExecution) TableName() string {
	return "action_executions"
}

// SafetyZone 禁区表
type SafetyZone struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DeviceID   string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	ZoneType   string         `gorm:"type:varchar(20);not null" json:"zone_type"`   // forbidden/caution/safe
	ZoneShape  string         `gorm:"type:varchar(20);not null" json:"zone_shape"`  // rectangle/circle/polygon
	ZoneData   string         `gorm:"type:jsonb;not null" json:"zone_data"`        // 区域坐标数据
	ZoneName   string         `gorm:"type:varchar(100)" json:"zone_name"`
	IsEnabled  bool           `gorm:"type:boolean;default:true" json:"is_enabled"`
	CreatedBy  *uint          `json:"created_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SafetyZone) TableName() string {
	return "safety_zones"
}

// EmbodiedDecisionLog 具身决策日志表
type EmbodiedDecisionLog struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	DeviceID         string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	DecisionType     string         `gorm:"type:varchar(50);not null" json:"decision_type"`
	Context          string         `gorm:"type:jsonb" json:"context"`           // 决策上下文
	ChosenAction     string         `gorm:"type:varchar(100)" json:"chosen_action"`
	ActionParams     string         `gorm:"type:jsonb" json:"action_params"`
	Confidence       float64        `gorm:"type:decimal(5,4)" json:"confidence"`
	Reasoning        string         `gorm:"type:text" json:"reasoning"`
	ExecutionResult  string         `gorm:"type:varchar(20)" json:"execution_result"`
	LatencyMs        int            `gorm:"type:int" json:"latency_ms"`
	DecidedAt        time.Time      `gorm:"not null" json:"decided_at"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (EmbodiedDecisionLog) TableName() string {
	return "embodied_decision_logs"
}

// SafetyLog 安全日志表
type SafetyLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	DeviceID  string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	EventType string         `gorm:"type:varchar(50);not null" json:"event_type"`  // collision/emergency_stop/zone_violation
	Severity  string         `gorm:"type:varchar(20);not null" json:"severity"`
	Details   string         `gorm:"type:jsonb" json:"details"`
	Location  string         `gorm:"type:jsonb" json:"location"`
	Resolved  bool           `gorm:"type:boolean;default:false" json:"resolved"`
	ResolvedAt *time.Time    `json:"resolved_at"`
	CreatedAt time.Time     `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SafetyLog) TableName() string {
	return "safety_logs"
}

// PerceptionData 感知数据（用于存储设备上报的感知结果）
type PerceptionData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	DeviceID  string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"` // visual/depth/touch
	Data      string         `gorm:"type:jsonb;not null" json:"data"`
	Timestamp time.Time      `json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PerceptionData) TableName() string {
	return "perception_data"
}

// NavigationTask 导航任务
type NavigationTask struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DeviceID     string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	TargetX      float64        `gorm:"type:decimal(10,3)" json:"target_x"`
	TargetY      float64        `gorm:"type:decimal(10,3)" json:"target_y"`
	TargetZ      float64        `gorm:"type:decimal(10,3)" json:"target_z"`
	Status       string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/running/completed/failed
	ErrorMessage string         `gorm:"type:text" json:"error_message"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (NavigationTask) TableName() string {
	return "navigation_tasks"
}

// ExploreTask 探索任务
type ExploreTask struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DeviceID     string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	Status       string         `gorm:"type:varchar(20);default:'idle'" json:"status"` // idle/running/paused/completed
	Coverage     float64        `gorm:"type:decimal(5,2);default:0" json:"coverage"`     // 探索覆盖率
	Strategy     string         `gorm:"type:varchar(50)" json:"strategy"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ExploreTask) TableName() string {
	return "explore_tasks"
}

// FollowTask 跟随任务
type FollowTask struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeviceID    string         `gorm:"type:varchar(100);not null;index" json:"device_id"`
	TargetID    string         `gorm:"type:varchar(100)" json:"target_id"` // 被跟随的目标ID
	Status      string         `gorm:"type:varchar(20);default:'idle'" json:"status"` // idle/running/stopped
	Distance    float64        `gorm:"type:decimal(10,3);default:1.5" json:"distance"` // 保持距离
	StartedAt   *time.Time      `json:"started_at"`
	StoppedAt   *time.Time      `json:"stopped_at"`
	CreatedAt   time.Time       `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (FollowTask) TableName() string {
	return "follow_tasks"
}

// DecisionStrategy 决策策略
type DecisionStrategy struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DeviceID    string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"device_id"`
	Strategy    string    `gorm:"type:varchar(50);not null" json:"strategy"` // safety_first/task_oriented/interactive/exploration
	Config      string    `gorm:"type:jsonb" json:"config"`                  // 策略配置
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (DecisionStrategy) TableName() string {
	return "decision_strategies"
}
