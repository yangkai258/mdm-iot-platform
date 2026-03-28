package models

import (
	"time"
)

// MapIntegrationConfig 地图服务集成配置
type MapIntegrationConfig struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Provider     string         `gorm:"type:varchar(50);uniqueIndex" json:"provider"` // amap/gaode/google
	APIKey       string         `gorm:"type:varchar(500)" json:"api_key"`
	APISecret    string         `gorm:"type:varchar(500)" json:"api_secret"`
	IsActive     bool           `gorm:"default:false" json:"is_active"`

	// 服务配置
	Services     string         `gorm:"type:varchar(200)" json:"services"` // geocoding/routing/location
	QuotaLimit  int            `gorm:"default:0" json:"quota_limit"`     // 日配额
	QuotaUsed   int            `gorm:"default:0" json:"quota_used"`      // 已使用
	QuotaResetAt *time.Time    `json:"quota_reset_at"`                   // 配额重置时间

	// 状态
	Status       string         `gorm:"type:varchar(20)" json:"status"`    // active/inactive/error
	ErrorMessage string         `gorm:"type:text" json:"error_message"`

	CreatedBy    string         `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedAt    time.Time     `json:"updated_at"`
	CreatedAt    time.Time     `json:"created_at"`
}

// TableName 表名
func (MapIntegrationConfig) TableName() string {
	return "map_integration_configs"
}

// MapServiceLog 地图服务调用日志
type MapServiceLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Provider     string         `gorm:"type:varchar(50)" json:"provider"`
	ServiceType string         `gorm:"type:varchar(50)" json:"service_type"` // geocoding/routing/geofencing
	Endpoint    string         `gorm:"type:varchar(200)" json:"endpoint"`

	// 请求信息
	RequestData string         `gorm:"type:text" json:"request_data"`
	ResponseData string        `gorm:"type:text" json:"response_data"`

	// 统计
	Latency     int            `gorm:"default:0" json:"latency"`      // 毫秒
	StatusCode int            `gorm:"default:0" json:"status_code"`

	// 费用
	Cost        float64        `gorm:"type:decimal(10,4);default:0" json:"cost"`

	DeviceID    string         `gorm:"type:varchar(64);index" json:"device_id"`
	UserID     string         `gorm:"type:varchar(64);index" json:"user_id"`

	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 表名
func (MapServiceLog) TableName() string {
	return "map_service_logs"
}

// PetLocation 宠物位置记录
type PetLocation struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PetID       string         `gorm:"type:varchar(64);index" json:"pet_id"`
	DeviceID    string         `gorm:"type:varchar(64);index" json:"device_id"`

	// 位置信息
	Latitude    float64        `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude   float64        `gorm:"type:decimal(10,7)" json:"longitude"`
	Altitude    float64        `gorm:"type:decimal(10,2)" json:"altitude"`
	Address     string         `gorm:"type:varchar(500)" json:"address"`      // 反 geocoding 地址
	POIName     string         `gorm:"type:varchar(200)" json:"poi_name"`    // 附近地标

	// 精度
	Accuracy    float64        `gorm:"type:decimal(8,2)" json:"accuracy"`     // 米

	// 状态
	BatteryLevel int           `gorm:"default:100" json:"battery_level"`

	LocationType string         `gorm:"type:varchar(20)" json:"location_type"` // gps/wifi/cell

	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 表名
func (PetLocation) TableName() string {
	return "pet_locations"
}
