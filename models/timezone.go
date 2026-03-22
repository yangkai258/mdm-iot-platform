package models

import (
	"time"
)

// TimezoneConfig 时区配置
type TimezoneConfig struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	EntityType     string    `json:"entity_type" gorm:"size:32"` // user/tenant/system
	EntityID       uint      `json:"entity_id" gorm:"index"`
	Timezone       string    `json:"timezone" gorm:"size:64"` // Asia/Shanghai, America/New_York
	DatetimeFormat string    `json:"datetime_format" gorm:"size:50;default:'%Y-%m-%d %H:%M:%S'"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
