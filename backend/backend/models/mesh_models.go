package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MeshNetwork BLE Mesh 网络
type MeshNetwork struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	NetworkUUID  string         `gorm:"size:64;uniqueIndex;not null" json:"network_uuid"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	NetworkKey   string         `gorm:"size:64" json:"network_key"` // encrypted
	DeviceCount  int64           `gorm:"default:0" json:"device_count"`
	Status       string         `gorm:"size:20;default:'active'" json:"status"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m *MeshNetwork) BeforeCreate(tx *gorm.DB) error {
	if m.NetworkUUID == "" {
		m.NetworkUUID = uuid.New().String()
	}
	return nil
}

// MeshNode BLE Mesh 节点
type MeshNode struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	NodeUUID       string         `gorm:"size:64;uniqueIndex;not null" json:"node_uuid"`
	NetworkUUID    string         `gorm:"size:64;not null;index" json:"network_uuid"`
	DeviceID       string         `gorm:"size:64;index" json:"device_id"` // bound device
	MacAddress     string         `gorm:"size:64" json:"mac_address"`
	NodeAddress    uint16         `gorm:"not null" json:"node_address"` // mesh unicast address
	ElementCount   int            `gorm:"default:1" json:"element_count"`
	FirmwareVersion string        `gorm:"size:32" json:"firmware_version"`
	HardwareVersion string        `gorm:"size:32" json:"hardware_version"`
	BatteryLevel   int            `gorm:"default:100" json:"battery_level"`
	RSSI           int            `json:"rssi"`
	IsProvisioned  bool           `gorm:"default:false" json:"is_provisioned"`
	IsOnline       bool           `gorm:"default:false;index" json:"is_online"`
	LastSeenAt     *time.Time     `json:"last_seen_at"`
	Status         string         `gorm:"size:20;default:'active'" json:"status"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m *MeshNode) BeforeCreate(tx *gorm.DB) error {
	if m.NodeUUID == "" {
		m.NodeUUID = uuid.New().String()
	}
	return nil
}

// MeshGroup BLE Mesh 组
type MeshGroup struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	GroupUUID    string         `gorm:"size:64;uniqueIndex;not null" json:"group_uuid"`
	NetworkUUID  string         `gorm:"size:64;not null;index" json:"network_uuid"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	GroupAddress uint16         `gorm:"not null" json:"group_address"` // mesh group address
	Description  string         `gorm:"type:text" json:"description"`
	MemberCount  int64          `gorm:"default:0" json:"member_count"`
	Status       string         `gorm:"size:20;default:'active'" json:"status"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (m *MeshGroup) BeforeCreate(tx *gorm.DB) error {
	if m.GroupUUID == "" {
		m.GroupUUID = uuid.New().String()
	}
	return nil
}

// MeshNodeGroup 节点组成员关系
type MeshNodeGroup struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	NodeUUID   string    `gorm:"size:64;not null;index" json:"node_uuid"`
	GroupUUID  string    `gorm:"size:64;not null;index" json:"group_uuid"`
	AddedAt    time.Time `json:"added_at"`
}

// MeshTelemetry BLE Mesh 遥测数据
type MeshTelemetry struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	NodeUUID     string         `gorm:"size:64;not null;index" json:"node_uuid"`
	TelemetryType string        `gorm:"size:32;not null" json:"telemetry_type"`
	Payload      string         `gorm:"type:jsonb" json:"payload"`
	RSSI         int            `json:"rssi"`
	HopCount     int            `gorm:"default:0" json:"hop_count"`
	RecordedAt   time.Time      `json:"recorded_at"`
	CreatedAt    time.Time      `json:"created_at"`
}
