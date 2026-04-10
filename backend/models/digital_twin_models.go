package models

import (
	"time"

	"gorm.io/gorm"
)

// VitalRecord 生命体征记录
type VitalRecord struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PetID      uint           `gorm:"index" json:"pet_id"`
	Type       string         `gorm:"type:varchar(32);index" json:"type"` // heart_rate, breathing, temperature
	Value      float64        `json:"value"`
	Unit       string         `gorm:"type:varchar(16)" json:"unit"`
	RecordedAt time.Time      `gorm:"index" json:"recorded_at"`
	Source     string         `gorm:"type:varchar(32)" json:"source"` // device, manual
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (VitalRecord) TableName() string {
	return "vital_records"
}

// HealthAlert 健康预警
type HealthAlert struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PetID       uint           `gorm:"index" json:"pet_id"`
	Type        string         `gorm:"type:varchar(64);index" json:"type"` // abnormal_heart_rate, fever, etc
	Level       string         `gorm:"type:varchar(16);index" json:"level"` // warning, critical
	Message     string         `gorm:"type:text" json:"message"`
	Status      string         `gorm:"type:varchar(16);index;default:pending" json:"status"` // pending, confirmed, ignored, resolved
	DetectedAt  time.Time      `gorm:"index" json:"detected_at"`
	ConfirmedAt *time.Time     `json:"confirmed_at"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HealthAlert) TableName() string {
	return "health_alerts"
}

// BehaviorEvent 行为事件
type BehaviorEvent struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	EventID   string         `gorm:"column:event_id;uniqueIndex" json:"event_id"`
	PetID     string         `gorm:"column:pet_id;index" json:"pet_id"`
	DeviceID  string         `gorm:"column:device_id" json:"device_id"`
	EventType string         `gorm:"column:event_type;index" json:"event_type"`
	EventData string         `gorm:"column:event_data;type:jsonb" json:"event_data"`
	Location  string         `gorm:"column:location;type:jsonb" json:"location"`
	Timestamp time.Time      `gorm:"column:timestamp" json:"timestamp"`
	StartTime *time.Time    `gorm:"column:start_time;index" json:"start_time"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 指定表名
func (BehaviorEvent) TableName() string {
	return "behavior_events"
}

// HighlightMoment 精彩瞬间
type HighlightMoment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PetID       uint           `gorm:"index" json:"pet_id"`
	Type        string         `gorm:"type:varchar(32);index" json:"type"` // cute, milestone, funny
	Title       string         `gorm:"type:varchar(128)" json:"title"`
	MediaURL    string         `gorm:"type:varchar(512)" json:"media_url"`
	CapturedAt  time.Time      `gorm:"index" json:"captured_at"`
	Description string         `gorm:"type:text" json:"description"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (HighlightMoment) TableName() string {
	return "highlight_moments"
}
