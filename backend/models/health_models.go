package models

import (
	"time"
)

// ExerciseRecord 运动记录
type ExerciseRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PetID     uint      `gorm:"index;not null" json:"pet_id"`
	Date      time.Time `gorm:"index;not null" json:"date"`
	Steps     int       `gorm:"default:0" json:"steps"`
	Distance  float64   `gorm:"type:decimal(10,2);default:0" json:"distance"` // km
	Calories  int       `gorm:"default:0" json:"calories"`
	Duration  int       `gorm:"default:0" json:"duration"` // minutes
	Goal      int       `gorm:"default:0" json:"goal"`    // daily step goal
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ExerciseRecord) TableName() string {
	return "exercise_records"
}

// SleepRecord 睡眠记录
type SleepRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PetID      uint      `gorm:"index;not null" json:"pet_id"`
	Date       time.Time `gorm:"index;not null" json:"date"`
	BedTime    time.Time `json:"bed_time"`
	WakeTime   time.Time `json:"wake_time"`
	DeepSleep  int       `gorm:"default:0" json:"deep_sleep"`  // minutes
	LightSleep int       `gorm:"default:0" json:"light_sleep"` // minutes
	REMSleep   int       `gorm:"default:0" json:"rem_sleep"`  // minutes
	TotalSleep int       `gorm:"default:0" json:"total_sleep"` // minutes
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (SleepRecord) TableName() string {
	return "sleep_records"
}

// BodyWeightRecord 体重记录
type BodyWeightRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PetID      uint      `gorm:"index;not null" json:"pet_id"`
	Weight     float64   `gorm:"type:decimal(5,2);not null" json:"weight"`
	RecordedAt time.Time `gorm:"not null" json:"recorded_at"`
	Note       string    `gorm:"type:varchar(255)" json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (BodyWeightRecord) TableName() string {
	return "body_weight_records"
}

// DietRecord 饮食记录
type DietRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PetID      uint      `gorm:"index;not null" json:"pet_id"`
	Date       time.Time `gorm:"index;not null" json:"date"`
	MealType   string    `gorm:"type:varchar(20);not null" json:"meal_type"` // breakfast, lunch, dinner, snack
	FoodType   string    `gorm:"type:varchar(100)" json:"food_type"`
	Amount     float64   `gorm:"type:decimal(10,2)" json:"amount"`
	Calories   int       `gorm:"default:0" json:"calories"`
	Note       string    `gorm:"type:varchar(255)" json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (DietRecord) TableName() string {
	return "diet_records"
}

// HealthWarning 健康预警
type HealthWarning struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PetID      uint      `gorm:"index;not null" json:"pet_id"`
	MemberID   uint      `gorm:"index" json:"member_id"`
	Type       string    `gorm:"type:varchar(50);not null" json:"type"`    // weight/sleep/exercise/diet/temperature
	Level      string    `gorm:"type:varchar(20);not null" json:"level"`    // info/warning/critical
	Status     string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/confirmed/ignored
	Title      string    `gorm:"type:varchar(255)" json:"title"`
	Message    string    `gorm:"type:text" json:"message"`
	Data       JSONMap   `gorm:"type:jsonb" json:"data"` // 附加数据
	ConfirmedAt *time.Time `json:"confirmed_at"`
	IgnoredAt  *time.Time `json:"ignored_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (HealthWarning) TableName() string {
	return "health_warnings"
}

// HealthReport 健康报告
type HealthReport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PetID       uint      `gorm:"index;not null" json:"pet_id"`
	MemberID    uint      `gorm:"index" json:"member_id"`
	ReportType  string    `gorm:"type:varchar(20);not null" json:"report_type"` // weekly/monthly
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Summary     string    `gorm:"type:text" json:"summary"`
	Data        JSONMap   `gorm:"type:jsonb" json:"data"` // 报告详细数据
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (HealthReport) TableName() string {
	return "health_reports"
}
