package models

import (
	"time"
)

// SensorEvent 传感器事件
type SensorEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EventID     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"event_id"`
	DeviceID    string    `gorm:"type:varchar(64);index;not null" json:"device_id"`
	SensorType  string    `gorm:"type:varchar(50);index" json:"sensor_type"` // temperature, humidity, battery, etc.
	SensorValue float64   `gorm:"type:decimal(10,2)" json:"sensor_value"`
	Unit        string    `gorm:"type:varchar(20)" json:"unit"` // °C, %, V, etc.
	Threshold   float64   `gorm:"type:decimal(10,2)" json:"threshold"`
	IsAbnormal  bool      `gorm:"type:boolean;default:false" json:"is_abnormal"`
	EventType   string    `gorm:"type:varchar(20);default:'normal'" json:"event_type"` // normal, warning, alert
	Description string    `gorm:"type:text" json:"description"`
	TenantID    string    `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (SensorEvent) TableName() string {
	return "sensor_events"
}
