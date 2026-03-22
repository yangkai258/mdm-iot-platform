package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SleepStage 睡眠阶段
const (
	SleepStageAwake    = "awake"    // 清醒
	SleepStageREM      = "rem"      // REM睡眠
	SleepStageLight    = "light"    // 浅睡眠
	SleepStageDeep     = "deep"     // 深睡眠
	SleepStageCore     = "core"     // 核心睡眠
)

// SleepQuality 睡眠质量等级
const (
	SleepQualityExcellent = "excellent" // 优秀
	SleepQualityGood      = "good"      // 良好
	SleepQualityFair      = "fair"      // 一般
	SleepQualityPoor      = "poor"      // 较差
	SleepQualityVeryPoor  = "very_poor" // 很差
)

// SleepRecord 睡眠记录
type SleepRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RecordUUID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"record_uuid"`
	PetUUID         string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID        string         `gorm:"type:varchar(64);index" json:"device_id"`
	SleepDate       time.Time      `gorm:"type:date;not null;index" json:"sleep_date"`
	BedTime         *time.Time     `gorm:"type:timestamp" json:"bed_time"`
	WakeTime        *time.Time     `gorm:"type:timestamp" json:"wake_time"`
	TotalDuration   int            `gorm:"type:int" json:"total_duration"`
	TotalDurationMinutes int       `gorm:"type:int" json:"total_duration_minutes"`
	REMCount        int            `gorm:"type:int" json:"rem_count"`
	REMTime         int            `gorm:"type:int" json:"rem_time"`
	LightSleepCount int            `gorm:"type:int" json:"light_sleep_count"`
	LightSleepTime  int            `gorm:"type:int" json:"light_sleep_time"`
	DeepSleepCount  int            `gorm:"type:int" json:"deep_sleep_count"`
	DeepSleepTime   int            `gorm:"type:int" json:"deep_sleep_time"`
	CoreSleepCount  int            `gorm:"type:int" json:"core_sleep_count"`
	CoreSleepTime   int            `gorm:"type:int" json:"core_sleep_time"`
	AwakeCount      int            `gorm:"type:int" json:"awake_count"`
	AwakeTime       int            `gorm:"type:int" json:"awake_time"`
	AvgHeartRate    float64        `gorm:"type:decimal(6,2)" json:"avg_heart_rate"`
	MinHeartRate    float64        `gorm:"type:decimal(6,2)" json:"min_heart_rate"`
	AvgRespiratoryRate float64    `gorm:"type:decimal(6,2)" json:"avg_respiratory_rate"`
	QualityScore    int            `gorm:"type:int" json:"quality_score"`
	QualityLevel    string         `gorm:"type:varchar(16)" json:"quality_level"`
	SleepEfficiency float64       `gorm:"type:decimal(5,2)" json:"sleep_efficiency"`
	Latency         int            `gorm:"type:int" json:"latency"`
	Midpoint        int            `gorm:"type:int" json:"midpoint"`
	Restlessness    float64        `gorm:"type:decimal(5,2)" json:"restlessness"`
	StageDetails    JSON           `gorm:"type:jsonb" json:"stage_details"`
	Events          JSON           `gorm:"type:jsonb" json:"events"`
	Environment     JSON           `gorm:"type:jsonb" json:"environment"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Tags            JSON           `gorm:"type:jsonb" json:"tags"`
	DataSource      string         `gorm:"type:varchar(32);default:'device'" json:"data_source"`
	IsNap           bool           `gorm:"type:boolean;default:false" json:"is_nap"`
	IsGoalAchieved  bool           `gorm:"type:boolean;default:false" json:"is_goal_achieved"`
	TenantID        string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SleepRecord) TableName() string {
	return "sleep_records"
}

// BeforeCreate 创建前自动生成 UUID
func (s *SleepRecord) BeforeCreate(tx *gorm.DB) error {
	if s.RecordUUID == "" {
		s.RecordUUID = uuid.New().String()
	}
	return nil
}

// SleepAnalysis 睡眠分析报告
type SleepAnalysis struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	AnalysisUUID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"analysis_uuid"`
	PetUUID           string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	AnalysisType      string         `gorm:"type:varchar(32);not null" json:"analysis_type"`
	StartDate         time.Time      `gorm:"type:date;not null;index" json:"start_date"`
	EndDate           time.Time      `gorm:"type:date;not null;index" json:"end_date"`
	AvgTotalDuration  int            `gorm:"type:int" json:"avg_total_duration"`
	AvgREMTime        int            `gorm:"type:int" json:"avg_rem_time"`
	AvgLightSleepTime int            `gorm:"type:int" json:"avg_light_sleep_time"`
	AvgDeepSleepTime  int            `gorm:"type:int" json:"avg_deep_sleep_time"`
	AvgCoreSleepTime  int            `gorm:"type:int" json:"avg_core_sleep_time"`
	AvgAwakeTime      int            `gorm:"type:int" json:"avg_awake_time"`
	AvgQualityScore   float64        `gorm:"type:decimal(5,2)" json:"avg_quality_score"`
	AvgSleepEfficiency float64       `gorm:"type:decimal(5,2)" json:"avg_sleep_efficiency"`
	AvgLatency        int            `gorm:"type:int" json:"avg_latency"`
	AvgRestlessness   float64        `gorm:"type:decimal(5,2)" json:"avg_restlessness"`
	TotalNaps         int            `gorm:"type:int" json:"total_naps"`
	TotalNapTime      int            `gorm:"type:int" json:"total_nap_time"`
	NightSleepCount   int            `gorm:"type:int" json:"night_sleep_count"`
	BestSleepDate     *time.Time     `gorm:"type:date" json:"best_sleep_date"`
	BestSleepQuality  int            `gorm:"type:int" json:"best_sleep_quality"`
	WorstSleepDate    *time.Time     `gorm:"type:date" json:"worst_sleep_date"`
	WorstSleepQuality int            `gorm:"type:int" json:"worst_sleep_quality"`
	TotalRecordDays   int            `gorm:"type:int" json:"total_record_days"`
	GoalAchievedDays  int            `gorm:"type:int" json:"goal_achieved_days"`
	TrendAnalysis     JSON           `gorm:"type:jsonb" json:"trend_analysis"`
	IssuesDetected    JSON           `gorm:"type:jsonb" json:"issues_detected"`
	Recommendations   JSON           `gorm:"type:jsonb" json:"recommendations"`
	TenantID          string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SleepAnalysis) TableName() string {
	return "sleep_analyses"
}

// BeforeCreate 创建前自动生成 UUID
func (s *SleepAnalysis) BeforeCreate(tx *gorm.DB) error {
	if s.AnalysisUUID == "" {
		s.AnalysisUUID = uuid.New().String()
	}
	return nil
}

// ReqSleepReport 设备上报睡眠数据请求
type ReqSleepReport struct {
	BedTime            string  `json:"bed_time"`
	WakeTime           string  `json:"wake_time"`
	TotalDuration      int     `json:"total_duration"`
	REMCount           int     `json:"rem_count"`
	REMTime            int     `json:"rem_time"`
	LightSleepCount    int     `json:"light_sleep_count"`
	LightSleepTime     int     `json:"light_sleep_time"`
	DeepSleepCount     int     `json:"deep_sleep_count"`
	DeepSleepTime      int     `json:"deep_sleep_time"`
	CoreSleepCount     int     `json:"core_sleep_count"`
	CoreSleepTime      int     `json:"core_sleep_time"`
	AwakeCount         int     `json:"awake_count"`
	AwakeTime          int     `json:"awake_time"`
	AvgHeartRate       float64 `json:"avg_heart_rate"`
	MinHeartRate       float64 `json:"min_heart_rate"`
	AvgRespiratoryRate float64 `json:"avg_respiratory_rate"`
	QualityScore       int     `json:"quality_score"`
	SleepEfficiency    float64 `json:"sleep_efficiency"`
	Latency            int     `json:"latency"`
	Restlessness       float64 `json:"restlessness"`
	StageDetails       []struct {
		Stage string `json:"stage"`
		Time  int    `json:"time"`
	} `json:"stage_details"`
	Environment struct {
		Temperature float64 `json:"temperature"`
		Humidity    float64 `json:"humidity"`
		NoiseLevel  string  `json:"noise_level"`
	} `json:"environment"`
	Notes  string   `json:"notes"`
	Tags   []string `json:"tags"`
	IsNap  bool     `json:"is_nap"`
}

// RespSleepList 睡眠记录列表响应
type RespSleepList struct {
	Records  []SleepRecord `json:"records"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// RespSleepAnalysis 睡眠分析响应
type RespSleepAnalysis struct {
	Analysis    *SleepAnalysis `json:"analysis"`
	DailyRecords []SleepRecord `json:"daily_records"`
	Trends      []SleepTrend   `json:"trends"`
}

// SleepTrend 睡眠趋势
type SleepTrend struct {
	Date           time.Time `json:"date"`
	TotalDuration  int       `json:"total_duration"`
	QualityScore   int       `json:"quality_score"`
	QualityLevel   string    `json:"quality_level"`
	DeepSleepTime  int       `json:"deep_sleep_time"`
	REMTime        int       `json:"rem_time"`
	AwakeTime      int       `json:"awake_time"`
}
