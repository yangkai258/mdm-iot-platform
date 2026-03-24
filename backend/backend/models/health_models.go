package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HealthWarning 健康预警
type HealthWarning struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	WarningUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"warning_uuid"`
	PetUUID          string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	DeviceID         string         `gorm:"size:64;index" json:"device_id"`
	Category         string         `gorm:"size:32;not null;index" json:"category"` // vital, behavior, environment
	Level            string         `gorm:"size:16;not null;index" json:"level"`    // info, warning, critical
	Title            string         `gorm:"size:128;not null" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	Symptoms         json.RawMessage `gorm:"type:jsonb" json:"symptoms"`
	Suggestions      json.RawMessage `gorm:"type:jsonb" json:"suggestions"`
	RelatedRecords   json.RawMessage `gorm:"type:jsonb" json:"related_records"`
	TriggerData      json.RawMessage `gorm:"type:jsonb" json:"trigger_data"`
	SourceType       string         `gorm:"size:32" json:"source_type"`
	SourceID         string         `gorm:"size:64" json:"source_id"`
	Status           string         `gorm:"size:20;default:'active';index" json:"status"`
	Priority         int64          `gorm:"default:0" json:"priority"`
	Severity         int64          `gorm:"default:0" json:"severity"`
	StartTime        time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime          *time.Time     `json:"end_time"`
	Duration         int64          `json:"duration"`
	ResolvedAt       *time.Time     `json:"resolved_at"`
	ResolvedBy       int64          `json:"resolved_by"`
	DismissedAt      *time.Time     `json:"dismissed_at"`
	DismissedBy      int64          `json:"dismissed_by"`
	DismissReason    string         `gorm:"type:text" json:"dismiss_reason"`
	AckedAt          *time.Time     `json:"acked_at"`
	AckedBy          int64          `json:"acked_by"`
	NotifiedAt       *time.Time     `json:"notified_at"`
	NotificationCount int64         `gorm:"default:0" json:"notification_count"`
	TenantID         string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (h *HealthWarning) BeforeCreate(tx *gorm.DB) error {
	if h.WarningUUID == "" {
		h.WarningUUID = uuid.New().String()
	}
	return nil
}

// VitalRecord 体征记录
type VitalRecord struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	RecordUUID     string         `gorm:"size:64;uniqueIndex;not null" json:"record_uuid"`
	PetUUID        string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	DeviceID       string         `gorm:"size:64;index" json:"device_id"`
	VitalType      string         `gorm:"size:32;not null;index" json:"vital_type"` // heart_rate, temperature, blood_pressure
	Value          float64        `gorm:"type:numeric(10,2);not null" json:"value"`
	Unit           string         `gorm:"size:16" json:"unit"`
	MinValue       float64        `gorm:"type:numeric(10,2)" json:"min_value"`
	MaxValue       float64        `gorm:"type:numeric(10,2)" json:"max_value"`
	AvgValue       float64        `gorm:"type:numeric(10,2)" json:"avg_value"`
	NormalRangeMin float64        `gorm:"type:numeric(10,2)" json:"normal_range_min"`
	NormalRangeMax float64        `gorm:"type:numeric(10,2)" json:"normal_range_max"`
	IsAbnormal     bool           `gorm:"default:false;index" json:"is_abnormal"`
	AbnormalLevel  string         `gorm:"size:16;index" json:"abnormal_level"`
	Severity       int64          `gorm:"default:0" json:"severity"`
	RecordedAt     time.Time      `gorm:"not null;index" json:"recorded_at"`
	DataSource     string         `gorm:"size:32;default:'device'" json:"data_source"`
	Metadata       json.RawMessage `gorm:"type:jsonb" json:"metadata"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (v *VitalRecord) BeforeCreate(tx *gorm.DB) error {
	if v.RecordUUID == "" {
		v.RecordUUID = uuid.New().String()
	}
	return nil
}

// VitalTrend 体征趋势
type VitalTrend struct {
	PetUUID       string         `gorm:"primaryKey;size:64" json:"pet_uuid"`
	VitalType     string         `gorm:"primaryKey;size:32" json:"vital_type"`
	CurrentValue  float64        `gorm:"type:numeric" json:"current_value"`
	AvgValue      float64        `gorm:"type:numeric" json:"avg_value"`
	MinValue      float64        `gorm:"type:numeric" json:"min_value"`
	MaxValue      float64        `gorm:"type:numeric" json:"max_value"`
	Trend         string         `gorm:"type:text" json:"trend"` // rising, falling, stable
	TrendRate     float64        `gorm:"type:numeric" json:"trend_rate"`
	IsAbnormal    bool           `gorm:"default:false" json:"is_abnormal"`
	AbnormalLevel string         `gorm:"type:text" json:"abnormal_level"`
	LastUpdated   time.Time      `json:"last_updated"`
}

// ExerciseGoal 运动目标
type ExerciseGoal struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	GoalUUID     string         `gorm:"size:64;uniqueIndex;not null" json:"goal_uuid"`
	PetUUID      string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	GoalType     string         `gorm:"size:32;not null" json:"goal_type"` // steps, distance, calories, duration
	TargetValue  float64        `gorm:"type:numeric(10,2);not null" json:"target_value"`
	Unit         string         `gorm:"size:16" json:"unit"`
	StartDate    time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate      *time.Time     `gorm:"type:date" json:"end_date"`
	CurrentValue float64        `gorm:"type:numeric(10,2);default:0" json:"current_value"`
	Progress     float64        `gorm:"type:numeric(5,2)" json:"progress"`
	Status       string         `gorm:"size:20;default:'active'" json:"status"`
	Priority     int64          `gorm:"default:0" json:"priority"`
	Notes        string         `gorm:"type:text" json:"notes"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (e *ExerciseGoal) BeforeCreate(tx *gorm.DB) error {
	if e.GoalUUID == "" {
		e.GoalUUID = uuid.New().String()
	}
	return nil
}

// ExerciseRecord 运动记录
type ExerciseRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RecordUUID      string         `gorm:"size:64;uniqueIndex;not null" json:"record_uuid"`
	PetUUID         string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	DeviceID        string         `gorm:"size:64;index" json:"device_id"`
	ExerciseType    string         `gorm:"size:32;not null;index" json:"exercise_type"` // walk, run, swim, play, training
	StartTime       time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime         *time.Time     `json:"end_time"`
	Duration        int64          `json:"duration"`
	DurationMinutes int64          `json:"duration_minutes"`
	Steps           int64          `json:"steps"`
	Distance        float64        `gorm:"type:numeric(10,2)" json:"distance"`
	DistanceUnit    string         `gorm:"size:16" json:"distance_unit"`
	CaloriesBurned  float64        `gorm:"type:numeric(10,2)" json:"calories_burned"`
	AvgHeartRate    float64        `gorm:"type:numeric(6,2)" json:"avg_heart_rate"`
	MaxHeartRate    float64        `gorm:"type:numeric(6,2)" json:"max_heart_rate"`
	MinHeartRate    float64        `gorm:"type:numeric(6,2)" json:"min_heart_rate"`
	AvgSpeed        float64        `gorm:"type:numeric(6,2)" json:"avg_speed"`
	MaxSpeed        float64        `gorm:"type:numeric(6,2)" json:"max_speed"`
	ElevationGain   float64        `gorm:"type:numeric(8,2)" json:"elevation_gain"`
	RouteData       json.RawMessage `gorm:"type:jsonb" json:"route_data"`
	Intensity       string         `gorm:"size:16" json:"intensity"`
	IntensityScore  int64          `json:"intensity_score"`
	QualityScore    int64          `json:"quality_score"`
	Weather         string         `gorm:"size:32" json:"weather"`
	Temperature     float64        `gorm:"type:numeric(5,2)" json:"temperature"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Tags            string         `gorm:"type:text" json:"tags"`
	DataSource      string         `gorm:"size:32;default:'device'" json:"data_source"`
	IsGoalAchieved  bool           `gorm:"default:false" json:"is_goal_achieved"`
	GoalID          int64          `gorm:"index" json:"goal_id"`
	TenantID        string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (e *ExerciseRecord) BeforeCreate(tx *gorm.DB) error {
	if e.RecordUUID == "" {
		e.RecordUUID = uuid.New().String()
	}
	return nil
}

// ExerciseSummary 运动摘要
type ExerciseSummary struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	PetUUID              string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	SummaryDate          time.Time      `gorm:"type:date;not null;index" json:"summary_date"`
	TotalDuration        int64          `json:"total_duration"`
	TotalSteps           int64          `json:"total_steps"`
	TotalDistance        float64        `gorm:"type:numeric(10,2)" json:"total_distance"`
	TotalCalories        float64        `gorm:"type:numeric(10,2)" json:"total_calories"`
	AvgHeartRate         float64        `gorm:"type:numeric(6,2)" json:"avg_heart_rate"`
	ExerciseCount        int64          `json:"exercise_count"`
	WalkCount            int64          `json:"walk_count"`
	RunCount             int64          `json:"run_count"`
	SwimCount            int64          `json:"swim_count"`
	PlayCount            int64          `json:"play_count"`
	TrainingCount        int64          `json:"training_count"`
	OtherCount           int64          `json:"other_count"`
	ActiveMinutes        int64          `json:"active_minutes"`
	GoalAchievedCount    int64          `json:"goal_achieved_count"`
	GoalTotalCount       int64          `json:"goal_total_count"`
	GoalAchievementRate  float64        `gorm:"type:numeric(5,2)" json:"goal_achievement_rate"`
	IntensityDistribution json.RawMessage `gorm:"type:jsonb" json:"intensity_distribution"`
	TenantID             string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

// ExerciseTrend 运动趋势
type ExerciseTrend struct {
	Date           time.Time `gorm:"primaryKey" json:"date"`
	TotalDuration  int64     `json:"total_duration"`
	TotalSteps     int64     `json:"total_steps"`
	TotalDistance  float64   `gorm:"type:numeric" json:"total_distance"`
	TotalCalories  float64   `gorm:"type:numeric" json:"total_calories"`
}

// HealthAlertRule 健康告警规则
type HealthAlertRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	RuleUUID       string         `gorm:"size:64;uniqueIndex;not null" json:"rule_uuid"`
	PetUUID        string         `gorm:"size:64;index" json:"pet_uuid"`
	RuleName       string         `gorm:"size:64;not null" json:"rule_name"`
	AlertType      string         `gorm:"size:32;not null" json:"alert_type"`
	Condition      json.RawMessage `gorm:"type:jsonb;not null" json:"condition"`
	AlertLevel     string         `gorm:"size:16;not null" json:"alert_level"`
	TitleTemplate  string         `gorm:"size:128" json:"title_template"`
	Suggestion     string         `gorm:"type:text" json:"suggestion"`
	CooldownPeriod int64          `gorm:"default:3600" json:"cooldown_period"`
	IsEnabled      bool           `gorm:"default:true" json:"is_enabled"`
	Priority       int64          `gorm:"default:5" json:"priority"`
	NotifyChannels []string       `gorm:"type:text[]" json:"notify_channels"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (h *HealthAlertRule) BeforeCreate(tx *gorm.DB) error {
	if h.RuleUUID == "" {
		h.RuleUUID = uuid.New().String()
	}
	return nil
}

// HealthMonitorSetting 健康监控设置
type HealthMonitorSetting struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	ProfileID             string         `gorm:"size:36;uniqueIndex;not null" json:"profile_id"`
	SettingID             string         `gorm:"size:36;uniqueIndex;not null" json:"setting_id"`
	HeartRateMonitor      bool           `gorm:"default:true" json:"heart_rate_monitor"`
	HeartRateMin          int16          `gorm:"default:50" json:"heart_rate_min"`
	HeartRateMax          int16          `gorm:"default:120" json:"heart_rate_max"`
	BpMonitor             bool           `gorm:"default:true" json:"bp_monitor"`
	BpSystolicMin         int16          `gorm:"default:90" json:"bp_systolic_min"`
	BpSystolicMax         int16          `gorm:"default:140" json:"bp_systolic_max"`
	BpDiastolicMin        int16          `gorm:"default:60" json:"bp_diastolic_min"`
	BpDiastolicMax        int16          `gorm:"default:90" json:"bp_diastolic_max"`
	SleepMonitor          bool           `gorm:"default:true" json:"sleep_monitor"`
	ActivityMonitor       bool           `gorm:"default:true" json:"activity_monitor"`
	StepGoal              int64          `gorm:"default:6000" json:"step_goal"`
	FallDetection         bool           `gorm:"default:true" json:"fall_detection"`
	FallAlertEnabled      bool           `gorm:"default:true" json:"fall_alert_enabled"`
	EmergencyAlertEnabled bool           `gorm:"default:true" json:"emergency_alert_enabled"`
	AlertThreshold        int64          `gorm:"default:3" json:"alert_threshold"`
	CheckInterval         int64          `gorm:"default:60" json:"check_interval"`
	ReportFrequency       string         `gorm:"size:20;default:'daily'" json:"report_frequency"`
	Enabled               bool           `gorm:"default:true" json:"enabled"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (h *HealthMonitorSetting) BeforeCreate(tx *gorm.DB) error {
	if h.ProfileID == "" {
		h.ProfileID = uuid.New().String()
	}
	if h.SettingID == "" {
		h.SettingID = uuid.New().String()
	}
	return nil
}

// HealthAlert 健康告警
type HealthAlert struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	AlertUUID        string         `gorm:"size:64;uniqueIndex;not null" json:"alert_uuid"`
	PetUUID          string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	DeviceID         string         `gorm:"size:64;index" json:"device_id"`
	AlertType        string         `gorm:"size:32;not null;index" json:"alert_type"`
	AlertLevel       string         `gorm:"size:16;not null;index" json:"alert_level"`
	Title            string         `gorm:"size:128;not null" json:"title"`
	Description      string         `gorm:"type:text" json:"description"`
	TriggerValue     float64        `gorm:"type:numeric(10,2)" json:"trigger_value"`
	ThresholdValue   float64        `gorm:"type:numeric(10,2)" json:"threshold_value"`
	Unit             string         `gorm:"size:16" json:"unit"`
	NormalRange      string         `gorm:"size:64" json:"normal_range"`
	Suggestion       string         `gorm:"type:text" json:"suggestion"`
	Urgency          int64          `gorm:"default:5" json:"urgency"`
	RelatedVitals    json.RawMessage `gorm:"type:jsonb" json:"related_vitals"`
	RelatedBehaviors json.RawMessage `gorm:"type:jsonb" json:"related_behaviors"`
	Status           string         `gorm:"size:16;default:'active';index" json:"status"`
	OccurredAt       time.Time      `gorm:"not null;index" json:"occurred_at"`
	AcknowledgedAt   *time.Time     `json:"acknowledged_at"`
	AcknowledgedBy   int64          `json:"acknowledged_by"`
	ResolvedAt       *time.Time     `json:"resolved_at"`
	ResolvedBy       int64          `json:"resolved_by"`
	IgnoredAt        *time.Time     `json:"ignored_at"`
	IgnoredBy        int64          `json:"ignored_by"`
	IgnoreReason     string         `gorm:"type:text" json:"ignore_reason"`
	NotifyChannels   []string       `gorm:"type:text[]" json:"notify_channels"`
	IsNotified       bool           `gorm:"default:false" json:"is_notified"`
	NotifiedAt       *time.Time     `json:"notified_at"`
	Tags             []string       `gorm:"type:text[]" json:"tags"`
	TenantID         string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (h *HealthAlert) BeforeCreate(tx *gorm.DB) error {
	if h.AlertUUID == "" {
		h.AlertUUID = uuid.New().String()
	}
	return nil
}

// PetVaccination 宠物疫苗接种记录
type PetVaccination struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	PetID            int64          `gorm:"not null;index" json:"pet_id"`
	UserID           int64          `gorm:"not null;index" json:"user_id"`
	VaccineName      string         `gorm:"size:100;not null" json:"vaccine_name"`
	VaccineType      string         `gorm:"size:50" json:"vaccine_type"`
	LotNumber        string         `gorm:"size:100" json:"lot_number"`
	Manufacturer     string         `gorm:"size:200" json:"manufacturer"`
	InoculationDate  time.Time      `gorm:"not null" json:"inoculation_date"`
	InoculationAge   string         `gorm:"size:50" json:"inoculation_age"`
	InoculationSite  string         `gorm:"size:100" json:"inoculation_site"`
	Inoculator       string         `gorm:"size:100" json:"inoculator"`
	VetClinic        string         `gorm:"size:200" json:"vet_clinic"`
	NextDoseDate     *time.Time     `json:"next_dose_date"`
	NextDoseMemo     string         `gorm:"size:255" json:"next_dose_memo"`
	AdverseReactions string         `gorm:"type:text" json:"adverse_reactions"`
	AdverseDetail    string         `gorm:"type:text" json:"adverse_detail"`
	CertificateURL   string         `gorm:"size:500" json:"certificate_url"`
	Remark           string         `gorm:"type:text" json:"remark"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// PetDietRecord 宠物饮食记录
type PetDietRecord struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PetID          int64          `gorm:"not null;index" json:"pet_id"`
	UserID         int64          `gorm:"not null;index" json:"user_id"`
	DeviceID       string         `gorm:"size:36;index" json:"device_id"`
	MealType       string         `gorm:"size:20;not null" json:"meal_type"` // breakfast, lunch, dinner, snack
	FoodName       string         `gorm:"size:200" json:"food_name"`
	FoodBrand      string         `gorm:"size:100" json:"food_brand"`
	FoodType       string         `gorm:"size:50" json:"food_type"`
	Amount         float64        `gorm:"type:numeric(10,2)" json:"amount"`
	AmountUnit     string         `gorm:"size:20;default:'g'" json:"amount_unit"`
	Calories       int64          `gorm:"default:0" json:"calories"`
	FeedingMethod  string         `gorm:"size:50" json:"feeding_method"`
	AutoFeederID   string         `gorm:"size:36" json:"auto_feeder_id"`
	EatTime        *time.Time     `json:"eat_time"`
	PlannedTime    *time.Time     `json:"planned_time"`
	Duration       int64          `gorm:"default:0" json:"duration"`
	Status         string         `gorm:"size:20;default:'completed'" json:"status"`
	LeftOver       float64        `gorm:"type:numeric(10,2);default:0" json:"left_over"`
	AppetiteScore  int64          `gorm:"default:5" json:"appetite_score"`
	HealthNote     string         `gorm:"type:text" json:"health_note"`
	PhotoURL       string         `gorm:"size:500" json:"photo_url"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// PetHealthRecord 宠物健康记录
type PetHealthRecord struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	RecordUUID         string         `gorm:"size:64;uniqueIndex;not null" json:"record_uuid"`
	PetUUID            string         `gorm:"size:64;not null;index" json:"pet_uuid"`
	RecordType         string         `gorm:"size:32;not null" json:"record_type"` // checkup, surgery, dental, vaccination
	RecordDate         time.Time      `gorm:"type:date;not null" json:"record_date"`
	Title              string         `gorm:"size:128;not null" json:"title"`
	Hospital           string         `gorm:"size:128" json:"hospital"`
	Doctor             string         `gorm:"size:64" json:"doctor"`
	VetName            string         `gorm:"size:64" json:"vet_name"`
	Diagnosis          string         `gorm:"type:text" json:"diagnosis"`
	Treatment          string         `gorm:"type:text" json:"treatment"`
	Prescription       string         `gorm:"type:text" json:"prescription"`
	Medications        []string       `gorm:"type:text[]" json:"medications"`
	Cost               float64        `gorm:"type:numeric(10,2)" json:"cost"`
	VaccineName        string         `gorm:"size:64" json:"vaccine_name"`
	NextDueDate        *time.Time     `json:"next_due_date"`
	Weight             float64        `gorm:"type:numeric(5,2)" json:"weight"`
	Notes              string         `gorm:"type:text" json:"notes"`
	Attachments        []string       `gorm:"type:text[]" json:"attachments"`
	IsInsured          bool           `gorm:"default:false" json:"is_insured"`
	InsuranceClaimUUID string         `gorm:"size:64" json:"insurance_claim_uuid"`
	TenantID           string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (p *PetHealthRecord) BeforeCreate(tx *gorm.DB) error {
	if p.RecordUUID == "" {
		p.RecordUUID = uuid.New().String()
	}
	return nil
}
