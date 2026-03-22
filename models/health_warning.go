package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WarningStatus 预警状态
const (
	WarningStatusActive   = "active"   // 待处理
	WarningStatusAcked    = "acked"    // 已确认
	WarningStatusDismissed = "dismissed" // 已忽略
	WarningStatusResolved = "resolved" // 已解决
)

// WarningLevel 预警级别
const (
	WarningLevelInfo     = "info"     // 提示
	WarningLevelWarning  = "warning"  // 警告
	WarningLevelCritical = "critical" // 严重
	WarningLevelEmergency = "emergency" // 紧急
)

// WarningCategory 预警类别
const (
	WarningCategoryVitalSign  = "vital_sign"  // 生命体征异常
	WarningCategoryBehavior   = "behavior"   // 行为异常
	WarningCategoryWeight     = "weight"     // 体重异常
	WarningCategoryActivity   = "activity"   // 活动量异常
	WarningCategorySleep      = "sleep"      // 睡眠异常
	WarningCategoryEating     = "eating"     // 饮食异常
	WarningCategoryDisease    = "disease"    // 疾病预警
	WarningCategoryOther      = "other"      // 其他
)

// HealthWarning 早期疾病预警记录
type HealthWarning struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	WarningUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"warning_uuid"`
	PetUUID        string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID       string         `gorm:"type:varchar(64);index" json:"device_id"`
	Category       string         `gorm:"type:varchar(32);not null;index" json:"category"`
	Level          string         `gorm:"type:varchar(16);not null;index" json:"level"`
	Title          string         `gorm:"type:varchar(128);not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Symptoms       JSON           `gorm:"type:jsonb" json:"symptoms"`
	Suggestions    JSON           `gorm:"type:jsonb" json:"suggestions"`
	RelatedRecords JSON           `gorm:"type:jsonb" json:"related_records"`
	TriggerData    JSON           `gorm:"type:jsonb" json:"trigger_data"`
	SourceType     string         `gorm:"type:varchar(32)" json:"source_type"`
	SourceID       string         `gorm:"type:varchar(64)" json:"source_id"`
	Status         string         `gorm:"type:varchar(20);default:'active';index" json:"status"`
	Priority       int            `gorm:"type:int;default:0" json:"priority"`
	Severity       int            `gorm:"type:int;default:0" json:"severity"`
	StartTime      time.Time      `gorm:"type:timestamp;not null;index" json:"start_time"`
	EndTime        *time.Time     `gorm:"type:timestamp" json:"end_time"`
	Duration       int            `gorm:"type:int" json:"duration"`
	ResolvedAt     *time.Time     `gorm:"type:timestamp" json:"resolved_at"`
	ResolvedBy     *uint          `gorm:"type:uint" json:"resolved_by"`
	DismissedAt    *time.Time     `gorm:"type:timestamp" json:"dismissed_at"`
	DismissedBy    *uint          `gorm:"type:uint" json:"dismissed_by"`
	DismissReason  string         `gorm:"type:text" json:"dismiss_reason"`
	AckedAt        *time.Time     `gorm:"type:timestamp" json:"acked_at"`
	AckedBy        *uint          `gorm:"type:uint" json:"acked_by"`
	NotifiedAt     *time.Time     `gorm:"type:timestamp" json:"notified_at"`
	NotificationCount int         `gorm:"type:int;default:0" json:"notification_count"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HealthWarning) TableName() string {
	return "health_warnings"
}

// BeforeCreate 创建前自动生成 UUID
func (h *HealthWarning) BeforeCreate(tx *gorm.DB) error {
	if h.WarningUUID == "" {
		h.WarningUUID = uuid.New().String()
	}
	if h.Status == "" {
		h.Status = WarningStatusActive
	}
	return nil
}

// DiseasePattern 疾病模式库 (P1)
type DiseasePattern struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PatternUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"pattern_uuid"`
	Species        string         `gorm:"type:varchar(32);not null;index" json:"species"`
	Breed          string         `gorm:"type:varchar(64);index" json:"breed"`
	DiseaseName    string         `gorm:"type:varchar(128);not null" json:"disease_name"`
	DiseaseCode    string         `gorm:"type:varchar(32);index" json:"disease_code"`
	Category       string         `gorm:"type:varchar(32);index" json:"category"`
	Description    string         `gorm:"type:text" json:"description"`
	Symptoms       JSON           `gorm:"type:jsonb" json:"symptoms"`
	RiskFactors    JSON           `gorm:"type:jsonb" json:"risk_factors"`
	WarningIndicators JSON        `gorm:"type:jsonb" json:"warning_indicators"`
	MinAge         int            `gorm:"type:int" json:"min_age"`
	MaxAge         int            `gorm:"type:int" json:"max_age"`
	IncidenceRate  float64        `gorm:"type:decimal(5,4)" json:"incidence_rate"`
	Severity       int            `gorm:"type:int" json:"severity"`
	RecommendedActions JSON       `gorm:"type:jsonb" json:"recommended_actions"`
	Tags           JSON           `gorm:"type:jsonb" json:"tags"`
	IsActive       bool           `gorm:"type:boolean;default:true;index" json:"is_active"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DiseasePattern) TableName() string {
	return "disease_patterns"
}

// BeforeCreate 创建前自动生成 UUID
func (d *DiseasePattern) BeforeCreate(tx *gorm.DB) error {
	if d.PatternUUID == "" {
		d.PatternUUID = uuid.New().String()
	}
	return nil
}

// ExerciseGoal 运动目标 (P1)
type ExerciseGoal struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	GoalUUID       string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"goal_uuid"`
	PetUUID        string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	GoalType       string         `gorm:"type:varchar(32);not null" json:"goal_type"`
	TargetValue    float64        `gorm:"type:decimal(10,2);not null" json:"target_value"`
	Unit           string         `gorm:"type:varchar(16)" json:"unit"`
	StartDate      time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate        *time.Time     `gorm:"type:date" json:"end_date"`
	CurrentValue   float64        `gorm:"type:decimal(10,2);default:0" json:"current_value"`
	Progress       float64        `gorm:"type:decimal(5,2)" json:"progress"`
	Status         string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	Priority       int            `gorm:"type:int;default:0" json:"priority"`
	Notes          string         `gorm:"type:text" json:"notes"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ExerciseGoal) TableName() string {
	return "exercise_goals"
}

// BeforeCreate 创建前自动生成 UUID
func (e *ExerciseGoal) BeforeCreate(tx *gorm.DB) error {
	if e.GoalUUID == "" {
		e.GoalUUID = uuid.New().String()
	}
	return nil
}
