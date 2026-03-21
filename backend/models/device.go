package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Device 设备资产主表
type Device struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	DeviceID        string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"device_id"`
	MacAddress      string         `gorm:"type:varchar(17);uniqueIndex;not null" json:"mac_address"`
	SnCode          string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"sn_code"`
	HardwareModel   string         `gorm:"type:varchar(32);not null" json:"hardware_model"`
	FirmwareVersion string         `gorm:"type:varchar(32);not null" json:"firmware_version"`
	BindUserID      *string        `gorm:"type:varchar(36);index" json:"bind_user_id"`
	LifecycleStatus int            `gorm:"type:smallint;default:1" json:"lifecycle_status"` // 1:待激活 2:服役中 3:维修 4:报废
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (d *Device) BeforeCreate(tx *gorm.DB) error {
	if d.DeviceID == "" {
		d.DeviceID = uuid.New().String()
	}
	return nil
}

// DeviceShadow 设备影子表
type DeviceShadow struct {
	DeviceID       string                 `gorm:"primaryKey;type:varchar(36)" json:"device_id"`
	IsOnline       bool                   `gorm:"default:false" json:"is_online"`
	BatteryLevel   int                    `gorm:"type:smallint" json:"battery_level"`
	CurrentMode    string                 `gorm:"type:varchar(20);default:'idle'" json:"current_mode"`
	LastIP         string                 `gorm:"type:varchar(45)" json:"last_ip"`
	LastHeartbeat  *time.Time             `gorm:"index" json:"last_heartbeat"`
	DesiredConfig  map[string]interface{} `gorm:"type:jsonb" json:"desired_config"`
	// 越狱/ROOT检测
	IsJailbroken bool   `gorm:"default:false" json:"is_jailbroken"`
	RootStatus   string `gorm:"type:varchar(20);default:'normal'" json:"root_status"` // normal, rooted, jailbroken
	// 地理位置
	Latitude  float64 `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(10,7)" json:"longitude"`

	// ============ 设备影子期望状态（NRD/免打扰） ============
	// 期望状态由管理平台下发，设备上线时同步
	DesiredNRDEnabled bool   `gorm:"default:false" json:"desired_nrd_enabled"` // NRD（夜间休息模式）是否启用
	DesiredNRDStart   string `gorm:"type:varchar(8);default:'23:00'" json:"desired_nrd_start"` // NRD 开始时间 (HH:MM)
	DesiredNRDEnd     string `gorm:"type:varchar(8);default:'07:00'" json:"desired_nrd_end"`   // NRD 结束时间 (HH:MM)
	DesiredDNDEnabled bool   `gorm:"default:false" json:"desired_dnd_enabled"` // 免打扰模式是否启用
	DesiredDNDStart   string `gorm:"type:varchar(8)" json:"desired_dnd_start"` // 免打扰开始时间 (HH:MM)
	DesiredDNDEnd     string `gorm:"type:varchar(8)" json:"desired_dnd_end"`   // 免打扰结束时间 (HH:MM)
	DesiredVolume     *int   `gorm:"type:smallint" json:"desired_volume"`     // 期望音量 (0-100)
	DesiredBrightness *int   `gorm:"type:smallint" json:"desired_brightness"` // 期望亮度 (0-100)
	DesiredPowerSave  bool   `gorm:"default:false" json:"desired_power_save"` // 节能模式
	DesiredVersion    string `gorm:"type:varchar(32)" json:"desired_version"` // 期望固件版本（用于OTA指令）
}

// CommandHistory 指令历史
type CommandHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeviceID  string    `gorm:"type:varchar(36);index" json:"device_id"`
	CmdID     string    `gorm:"type:varchar(36)" json:"cmd_id"`
	CmdType   string    `gorm:"type:varchar(20)" json:"cmd_type"`
	Action    string    `gorm:"type:varchar(50)" json:"action"`
	Status    string    `gorm:"type:varchar(20)" json:"status"` // sent, delivered, executed, failed
	SentAt    time.Time `json:"sent_at"`
	CreatedAt time.Time `json:"created_at"`
}

// PetProfile 宠物配置
type PetProfile struct {
	DeviceID         string                 `gorm:"primaryKey;type:varchar(36)" json:"device_id"`
	PetName          string                 `gorm:"type:varchar(64);default:'Mimi'" json:"pet_name"`
	Personality      string                 `gorm:"type:varchar(32);default:'lively'" json:"personality"`
	InteractionFreq  string                 `gorm:"type:varchar(16);default:'medium'" json:"interaction_freq"`
	DNDStartTime    string                 `gorm:"type:varchar(8);default:'23:00'" json:"dnd_start_time"`
	DNDEndTime      string                 `gorm:"type:varchar(8);default:'08:00'" json:"dnd_end_time"`
	CustomRules     map[string]interface{} `gorm:"type:jsonb" json:"custom_rules"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// DeviceData MQTT设备数据消息
type DeviceData struct {
	DeviceID     string                 `json:"device_id"`
	Battery      int                    `json:"battery"`
	Mode         string                 `json:"mode"`
	IP           string                 `json:"ip"`
	Timestamp   int64                  `json:"timestamp"`
	ExtraData   map[string]interface{} `json:"extra_data"`
}
