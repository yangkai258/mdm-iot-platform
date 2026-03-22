package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExerciseType 运动类型
const (
	ExerciseTypeWalk      = "walk"       // 散步
	ExerciseTypeRun       = "run"        // 跑步
	ExerciseTypeSwim      = "swim"       // 游泳
	ExerciseTypePlay      = "play"       // 玩耍
	ExerciseTypeTraining  = "training"   // 训练
	ExerciseTypeOther     = "other"      // 其他
)

// ExerciseRecord 运动记录
type ExerciseRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RecordUUID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"record_uuid"`
	PetUUID         string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID        string         `gorm:"type:varchar(64);index" json:"device_id"`
	ExerciseType    string         `gorm:"type:varchar(32);not null;index" json:"exercise_type"`
	StartTime       time.Time      `gorm:"type:timestamp;not null;index" json:"start_time"`
	EndTime         *time.Time     `gorm:"type:timestamp" json:"end_time"`
	Duration        int            `gorm:"type:int" json:"duration"`
	DurationMinutes int            `gorm:"type:int" json:"duration_minutes"`
	Steps           int            `gorm:"type:int" json:"steps"`
	Distance        float64        `gorm:"type:decimal(10,2)" json:"distance"`
	DistanceUnit    string         `gorm:"type:varchar(16)" json:"distance_unit"`
	CaloriesBurned  float64        `gorm:"type:decimal(10,2)" json:"calories_burned"`
	AvgHeartRate    float64        `gorm:"type:decimal(6,2)" json:"avg_heart_rate"`
	MaxHeartRate    float64        `gorm:"type:decimal(6,2)" json:"max_heart_rate"`
	MinHeartRate    float64        `gorm:"type:decimal(6,2)" json:"min_heart_rate"`
	AvgSpeed        float64        `gorm:"type:decimal(6,2)" json:"avg_speed"`
	MaxSpeed        float64        `gorm:"type:decimal(6,2)" json:"max_speed"`
	ElevationGain   float64        `gorm:"type:decimal(8,2)" json:"elevation_gain"`
	RouteData       JSON           `gorm:"type:jsonb" json:"route_data"`
	Intensity       string         `gorm:"type:varchar(16)" json:"intensity"`
	IntensityScore  int            `gorm:"type:int" json:"intensity_score"`
	QualityScore    int            `gorm:"type:int" json:"quality_score"`
	Weather         string         `gorm:"type:varchar(32)" json:"weather"`
	Temperature     float64        `gorm:"type:decimal(5,2)" json:"temperature"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Tags            JSON           `gorm:"type:jsonb" json:"tags"`
	DataSource      string         `gorm:"type:varchar(32);default:'device'" json:"data_source"`
	IsGoalAchieved  bool           `gorm:"type:boolean;default:false" json:"is_goal_achieved"`
	GoalID          *uint          `gorm:"type:uint;index" json:"goal_id"`
	TenantID        string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ExerciseRecord) TableName() string {
	return "exercise_records"
}

// BeforeCreate 创建前自动生成 UUID
func (e *ExerciseRecord) BeforeCreate(tx *gorm.DB) error {
	if e.RecordUUID == "" {
		e.RecordUUID = uuid.New().String()
	}
	return nil
}

// ExerciseSummary 运动汇总统计
type ExerciseSummary struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	PetUUID           string    `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	SummaryDate       time.Time `gorm:"type:date;not null;index" json:"summary_date"`
	TotalDuration     int       `gorm:"type:int" json:"total_duration"`
	TotalSteps        int       `gorm:"type:int" json:"total_steps"`
	TotalDistance     float64   `gorm:"type:decimal(10,2)" json:"total_distance"`
	TotalCalories     float64   `gorm:"type:decimal(10,2)" json:"total_calories"`
	AvgHeartRate      float64   `gorm:"type:decimal(6,2)" json:"avg_heart_rate"`
	ExerciseCount     int       `gorm:"type:int" json:"exercise_count"`
	WalkCount         int       `gorm:"type:int" json:"walk_count"`
	RunCount          int       `gorm:"type:int" json:"run_count"`
	SwimCount         int       `gorm:"type:int" json:"swim_count"`
	PlayCount         int       `gorm:"type:int" json:"play_count"`
	TrainingCount     int       `gorm:"type:int" json:"training_count"`
	OtherCount        int       `gorm:"type:int" json:"other_count"`
	ActiveMinutes     int       `gorm:"type:int" json:"active_minutes"`
	GoalAchievedCount int       `gorm:"type:int" json:"goal_achieved_count"`
	GoalTotalCount    int       `gorm:"type:int" json:"goal_total_count"`
	GoalAchievementRate float64 `gorm:"type:decimal(5,2)" json:"goal_achievement_rate"`
	IntensityDistribution JSON   `gorm:"type:jsonb" json:"intensity_distribution"`
	TenantID          string    `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ExerciseSummary) TableName() string {
	return "exercise_summaries"
}

// ReqExerciseReport 设备上报运动数据请求
type ReqExerciseReport struct {
	ExerciseType    string  `json:"exercise_type" binding:"required"`
	StartTime        string  `json:"start_time" binding:"required"`
	EndTime          string  `json:"end_time"`
	Duration         int    `json:"duration"`
	Steps            int    `json:"steps"`
	Distance         float64 `json:"distance"`
	CaloriesBurned   float64 `json:"calories_burned"`
	AvgHeartRate     float64 `json:"avg_heart_rate"`
	MaxHeartRate     float64 `json:"max_heart_rate"`
	MinHeartRate     float64 `json:"min_heart_rate"`
	AvgSpeed         float64 `json:"avg_speed"`
	MaxSpeed         float64 `json:"max_speed"`
	Intensity        string  `json:"intensity"`
	Notes            string  `json:"notes"`
	Tags             []string `json:"tags"`
}

// RespExerciseList 运动记录列表响应
type RespExerciseList struct {
	Records  []ExerciseRecord `json:"records"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// RespExerciseSummary 运动汇总响应
type RespExerciseSummary struct {
	Daily    *ExerciseSummary    `json:"daily"`
	Weekly   *ExerciseSummary    `json:"weekly"`
	Monthly  *ExerciseSummary    `json:"monthly"`
	Trends   []ExerciseTrend    `json:"trends"`
	Goals    []ExerciseGoal     `json:"goals"`
}

// ExerciseTrend 运动趋势
type ExerciseTrend struct {
	Date           time.Time `json:"date"`
	TotalDuration int       `json:"total_duration"`
	TotalSteps    int       `json:"total_steps"`
	TotalDistance float64   `json:"total_distance"`
	TotalCalories float64   `json:"total_calories"`
}
