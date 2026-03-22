package models

import (
	"time"

	"gorm.io/gorm"
)

// MeshDevice BLE Mesh 设备
type MeshDevice struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	DeviceID        string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"device_id"`
	MeshUUID        string         `gorm:"type:varchar(64);index" json:"mesh_uuid"` // Mesh 网络内唯一标识
	ParentDeviceID  string         `gorm:"type:varchar(64)" json:"parent_device_id"` // 上级设备（空=根节点）
	MeshAddress     string         `gorm:"type:varchar(16)" json:"mesh_address"`    // Mesh 单播地址
	Role            string         `gorm:"type:varchar(16);default:'node'" json:"role"` // relay, proxy, friend, node
	HopCount        int            `gorm:"default:0" json:"hop_count"`               // 到根节点的跳数
	SignalStrength  int            `gorm:"type:smallint" json:"signal_strength"`    // 信号强度 dBm
	ConnectionStatus string         `gorm:"type:varchar(16);default:'disconnected'" json:"connection_status"` // disconnected, connecting, connected
	ConnectedAt     *time.Time     `json:"connected_at"`
	LastSeenAt      *time.Time     `json:"last_seen_at"`
	HardwareModel   string         `gorm:"type:varchar(64)" json:"hardware_model"`
	FirmwareVersion string         `gorm:"type:varchar(32)" json:"firmware_version"`
	OrgID           uint           `gorm:"index" json:"org_id"`
	CreateUserID    uint           `gorm:"index" json:"create_user_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MeshDevice) TableName() string {
	return "mesh_devices"
}

// MeshNetwork BLE Mesh 网络
type MeshNetwork struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(128);not null" json:"name"`
	NetworkID      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"network_id"` // Mesh Network ID
	SSID           string         `gorm:"type:varchar(64)" json:"ssid"`    // Mesh 网络名称
	SecurityKey    string         `gorm:"type:varchar(128)" json:"-"`       // 安全密钥（不返回给前端）
	EncryptionType string         `gorm:"type:varchar(16);default:'aes256'" json:"encryption_type"` // aes128, aes256
	Channel        int            `gorm:"default:6" json:"channel"`        // BLE 信道
	GatewayDeviceID string        `gorm:"type:varchar(64)" json:"gateway_device_id"` // Mesh 网关设备
	Status         string         `gorm:"type:varchar(16);default:'inactive'" json:"status"` // inactive, active, configuring
	DeviceCount    int            `gorm:"default:0" json:"device_count"`   // 成员设备数量
	OrgID          uint           `gorm:"index" json:"org_id"`
	CreateUserID   uint           `gorm:"index" json:"create_user_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (MeshNetwork) TableName() string {
	return "mesh_networks"
}

// MeshNetworkMember Mesh 网络成员
type MeshNetworkMember struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	NetworkID    uint           `gorm:"not null;index" json:"network_id"`
	DeviceID     string         `gorm:"type:varchar(64);not null" json:"device_id"`
	MeshAddress  string         `gorm:"type:varchar(16)" json:"mesh_address"`
	Role         string         `gorm:"type:varchar(16);default:'node'" json:"role"`
	JoinedAt     *time.Time     `json:"joined_at"`
	LatencyMs    float64        `gorm:"default:0" json:"latency_ms"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	// GORM 关联
	MeshNetwork  MeshNetwork     `gorm:"foreignKey:NetworkID" json:"mesh_network,omitempty"`
}

// TableName 指定表名
func (MeshNetworkMember) TableName() string {
	return "mesh_network_members"
}
