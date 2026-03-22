package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SmartHomeDevice 米家智能家居设备
type SmartHomeDevice struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DeviceUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"device_uuid"`
	DeviceName    string         `gorm:"type:varchar(128);not null" json:"device_name"`
	DeviceType    string         `gorm:"type:varchar(32);not null" json:"device_type"` // light/switch/sensor/thermostat/camera/lock/air_conditioner/fan/humidifier
	Brand         string         `gorm:"type:varchar(32);default:'xiaomi'" json:"brand"` // xiaomi/philips/mijia/aqara/...
	XiaomiDeviceID string       `gorm:"type:varchar(128);index" json:"xiaomi_device_id"` // 米家设备 DID
	XiaomiToken   string         `gorm:"type:varchar(256)" json:"xiaomi_token"`          // 米家 API Token
	RoomName      string         `gorm:"type:varchar(64)" json:"room_name"`             // 房间名称
	OnlineStatus  string         `gorm:"type:varchar(16);default:'offline'" json:"online_status"` // online/offline/unknown
	DeviceStatus  string         `gorm:"type:text" json:"device_status"`                 // JSON 格式的设备状态
	IsEnabled     bool           `gorm:"type:boolean;default:true" json:"is_enabled"`
	HouseholdID   *uint          `gorm:"index" json:"household_id"`
	OwnerID       uint           `gorm:"index;not null" json:"owner_id"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SmartHomeDevice) TableName() string {
	return "smart_home_devices"
}

// BeforeCreate 创建前自动生成 UUID
func (s *SmartHomeDevice) BeforeCreate(tx *gorm.DB) error {
	if s.DeviceUUID == "" {
		s.DeviceUUID = uuid.New().String()
	}
	return nil
}

// SmartHomeTrigger 设备联动触发器
type SmartHomeTrigger struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TriggerUUID   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"trigger_uuid"`
	TriggerName   string         `gorm:"type:varchar(128);not null" json:"trigger_name"`
	TriggerType   string         `gorm:"type:varchar(32);not null" json:"trigger_type"` // device_status/time/scene/pet_behavior
	SourceDevice  string         `gorm:"type:varchar(64)" json:"source_device"`           // 触发源设备UUID
	ConditionExpr string         `gorm:"type:text" json:"condition_expr"`               // 条件表达式 JSON
	ActionExpr    string         `gorm:"type:text;not null" json:"action_expr"`         // 执行动作 JSON
	IsEnabled     bool           `gorm:"type:boolean;default:true" json:"is_enabled"`
	LastTriggered *time.Time     `json:"last_triggered"`
	TriggerCount  int            `gorm:"type:int;default:0" json:"trigger_count"`
	HouseholdID   *uint          `gorm:"index" json:"household_id"`
	OwnerID       uint           `gorm:"index;not null" json:"owner_id"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SmartHomeTrigger) TableName() string {
	return "smart_home_triggers"
}

// BeforeCreate 创建前自动生成 UUID
func (s *SmartHomeTrigger) BeforeCreate(tx *gorm.DB) error {
	if s.TriggerUUID == "" {
		s.TriggerUUID = uuid.New().String()
	}
	return nil
}
