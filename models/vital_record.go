package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VitalType 生命体征类型
const (
	VitalTypeHeartRate    = "heart_rate"
	VitalTypeBodyTemp     = "body_temp"
	VitalTypeActivity     = "activity"
	VitalTypeSleepQuality = "sleep_quality"
	VitalTypeCalories     = "calories"
	VitalTypeWaterIntake  = "water_intake"
	VitalTypeWeight       = "weight"
	VitalTypeBloodOxygen  = "blood_oxygen"
	VitalTypeRespiratory  = "respiratory"
)

// VitalRecord 生命体征记录 (Sprint 18)
type VitalRecord struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	RecordUUID     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"record_uuid"`
	PetUUID        string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID       string         `gorm:"type:varchar(64);index" json:"device_id"`
	VitalType      string         `gorm:"type:varchar(32);not null;index" json:"vital_type"`
	Value          float64        `gorm:"type:decimal(10,2);not null" json:"value"`
	Unit           string         `gorm:"type:varchar(16)" json:"unit"`
	MinValue       float64        `gorm:"type:decimal(10,2)" json:"min_value"`
	MaxValue       float64        `gorm:"type:decimal(10,2)" json:"max_value"`
	AvgValue       float64        `gorm:"type:decimal(10,2)" json:"avg_value"`
	NormalRangeMin float64        `gorm:"type:decimal(10,2)" json:"normal_range_min"`
	NormalRangeMax float64        `gorm:"type:decimal(10,2)" json:"normal_range_max"`
	IsAbnormal     bool           `gorm:"type:boolean;default:false;index" json:"is_abnormal"`
	AbnormalLevel  string         `gorm:"type:varchar(16);index" json:"abnormal_level"`
	Severity       int            `gorm:"type:int;default:0" json:"severity"`
	RecordedAt     time.Time      `gorm:"type:timestamp;not null;index" json:"recorded_at"`
	DataSource     string         `gorm:"type:varchar(32);default:'device'" json:"data_source"`
	Metadata       JSON           `gorm:"type:jsonb" json:"metadata"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (VitalRecord) TableName() string {
	return "vital_records"
}

func (v *VitalRecord) BeforeCreate(tx *gorm.DB) error {
	if v.RecordUUID == "" {
		v.RecordUUID = uuid.New().String()
	}
	return nil
}

// VitalStats 生命体征日统计
type VitalStats struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PetUUID        string    `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID       string    `gorm:"type:varchar(64);index" json:"device_id"`
	StatsDate      time.Time `gorm:"type:date;not null;index" json:"stats_date"`
	HeartRateAvg   float64   `gorm:"type:decimal(6,2)" json:"heart_rate_avg"`
	HeartRateMin   float64   `gorm:"type:decimal(6,2)" json:"heart_rate_min"`
	HeartRateMax   float64   `gorm:"type:decimal(6,2)" json:"heart_rate_max"`
	BodyTempAvg    float64   `gorm:"type:decimal(5,2)" json:"body_temp_avg"`
	ActivityTotal  float64   `gorm:"type:decimal(10,2)" json:"activity_total"`
	SleepDuration  float64   `gorm:"type:decimal(6,2)" json:"sleep_duration"`
	SleepQuality   float64   `gorm:"type:decimal(5,2)" json:"sleep_quality"`
	CaloriesBurned float64   `gorm:"type:decimal(8,2)" json:"calories_burned"`
	WaterIntake    float64   `gorm:"type:decimal(8,2)" json:"water_intake"`
	Weight         float64   `gorm:"type:decimal(6,2)" json:"weight"`
	Steps          int       `gorm:"type:int" json:"steps"`
	TenantID       string    `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (VitalStats) TableName() string {
	return "vital_stats"
}

// VitalTrend 生命体征趋势
type VitalTrend struct {
	PetUUID       string    `json:"pet_uuid"`
	VitalType     string    `json:"vital_type"`
	CurrentValue  float64   `json:"current_value"`
	AvgValue      float64   `json:"avg_value"`
	MinValue      float64   `json:"min_value"`
	MaxValue      float64   `json:"max_value"`
	Trend         string    `json:"trend"`
	TrendRate     float64   `json:"trend_rate"`
	IsAbnormal    bool      `json:"is_abnormal"`
	AbnormalLevel string    `json:"abnormal_level"`
	LastUpdated   time.Time `json:"last_updated"`
}

// ReqVitalReport 设备上报体征数据请求
type ReqVitalReport struct {
	VitalType string  `json:"vital_type" binding:"required"`
	Value     float64 `json:"value" binding:"required"`
	Unit      string  `json:"unit"`
	Metadata  JSON    `json:"metadata"`
}

// RespVitals 生命体征响应
type RespVitals struct {
	PetUUID       string            `json:"pet_uuid"`
	LatestRecords []VitalRecord     `json:"latest_records"`
	Trend         []VitalTrend     `json:"trend"`
	Alerts        []HealthAlert    `json:"alerts"`
	LastUpdatedAt time.Time        `json:"last_updated_at"`
}

// RespVitalHistory 生命体征历史响应
type RespVitalHistory struct {
	Records  []VitalRecord `json:"records"`
	Stats    VitalStats    `json:"stats"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}
