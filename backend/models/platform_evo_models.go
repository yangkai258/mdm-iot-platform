package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== 端侧模型 =====

// EdgeModel 端侧AI模型
type EdgeModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(128);not null" json:"name"`
	Version         string         `gorm:"type:varchar(32);not null" json:"version"`
	ModelType       string         `gorm:"type:varchar(64);not null" json:"model_type"` // vision/nlp/speech/mixed
	Architecture    string         `gorm:"type:varchar(64)" json:"architecture"`        // tensorflow/pytorch/onnx/custom
	Framework       string         `gorm:"type:varchar(32)" json:"framework"`            // tflite/onnxruntime/tensorrt
	FileURL         string         `gorm:"type:varchar(512)" json:"file_url"`
	FileSize        int64          `gorm:"type:bigint" json:"file_size"` // bytes
	FileHash        string         `gorm:"type:varchar(128)" json:"file_hash"`
	Description     string         `gorm:"type:text" json:"description"`
	MinMemoryMB     int            `gorm:"default:512" json:"min_memory_mb"`
	MinStorageMB    int            `gorm:"default:1024" json:"min_storage_mb"`
	InputShape      JSON           `gorm:"type:jsonb" json:"input_shape"`      // {"width":224,"height":224,"channels":3}
	OutputShape     JSON           `gorm:"type:jsonb" json:"output_shape"`     // {"classes":1000}
	Quantization    string         `gorm:"type:varchar(32)" json:"quantization"` // none/int8/fp16
	Accuracy        float64        `gorm:"type:decimal(5,2)" json:"accuracy"`  // 推理精度指标
	LatencyMs       int            `gorm:"default:0" json:"latency_ms"`        // 推理延迟
	ThroughputFPS   float64        `gorm:"type:decimal(10,2)" json:"throughput_fps"`
	Status          string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/active/deprecated
	DeviceTypes     StringArray    `gorm:"type:text" json:"device_types"`      // 兼容设备类型列表
	Tags            StringArray    `gorm:"type:text" json:"tags"`
	DownloadCount   int            `gorm:"default:0" json:"download_count"`
	Rating          float64        `gorm:"type:decimal(3,2)" json:"rating"`
	CreatedBy       uint           `gorm:"index" json:"created_by"`
	ApprovedBy      *uint          `gorm:"index" json:"approved_by"`
	ApprovedAt     *time.Time     `json:"approved_at"`
	PublishedAt    *time.Time     `json:"published_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (EdgeModel) TableName() string {
	return "edge_models"
}

// ===== 模型分片 =====

// ModelShard 模型分片
type ModelShard struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	ModelID         uint           `gorm:"index;not null" json:"model_id"`
	ShardIndex      int            `gorm:"not null" json:"shard_index"` // 分片序号
	ShardHash       string         `gorm:"type:varchar(128);not null" json:"shard_hash"`
	FileURL         string         `gorm:"type:varchar(512)" json:"file_url"`
	FileSize        int64          `gorm:"type:bigint" json:"file_size"`
	LayerRange      JSON           `gorm:"type:jsonb" json:"layer_range"` // {"start_layer":0,"end_layer":50}
	Dependencies    StringArray    `gorm:"type:text" json:"dependencies"` // 依赖的其他分片ID列表
	SizeMB          float64        `gorm:"type:decimal(10,2)" json:"size_mb"`
	IsBaseShard     bool           `gorm:"default:false" json:"is_base_shard"` // 是否基础分片
	IsLoaded        bool           `gorm:"default:false" json:"is_loaded"`
	LoadedAt       *time.Time     `json:"loaded_at"`
	LoadDurationMs  int            `gorm:"default:0" json:"load_duration_ms"`
	MemoryUsageMB   int            `gorm:"default:0" json:"memory_usage_mb"`
	Priority        int            `gorm:"default:0" json:"priority"`        // 加载优先级
	Status          string         `gorm:"type:varchar(20);default:'ready'" json:"status"` // ready/loading/loaded/error
	ErrorMessage    string         `gorm:"type:text" json:"error_message"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ModelShard) TableName() string {
	return "model_shards"
}

// ===== BLE Mesh 网络 =====

// BLEMeshNetwork BLE Mesh网络
type BLEMeshNetwork struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `gorm:"type:varchar(128);not null" json:"name"`
	NetworkID        string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"network_id"` // Mesh网络ID
	NetKey           string         `gorm:"type:varchar(128)" json:"-"`                              // 网络密钥（加密存储）
	AppKey           string         `gorm:"type:varchar(128)" json:"-"`                              // 应用密钥
	IVIndex          int            `gorm:"default:0" json:"iv_index"`
	UnicastAddress   int            `gorm:"default:0" json:"unicast_address"`   // 下一个单播地址
	GroupCount       int            `gorm:"default:0" json:"group_count"`       // 组数量
	NodeCount        int            `gorm:"default:0" json:"node_count"`        // 已入网节点数
	MaxNodes         int            `gorm:"default:256" json:"max_nodes"`       // 最大节点数
	Frequency        int            `gorm:"default:2440" json:"frequency"`      // 工作频率(MHz)
	TxPowerDbm       int            `gorm:"default:0" json:"tx_power_dbm"`       // 发射功率
	ChannelMap       JSON           `gorm:"type:jsonb" json:"channel_map"`      // 信道图
	SecurityLevel    string         `gorm:"type:varchar(20);default:'high'" json:"security_level"` // low/medium/high
	ProxyEnabled     bool           `gorm:"default:false" json:"proxy_enabled"` // 是否启用代理节点
	FriendEnabled    bool           `gorm:"default:false" json:"friend_enabled"` // 是否启用Friend节点
	RelayEnabled     bool           `gorm:"default:true" json:"relay_enabled"`  // 是否启用中继
	LowPowerEnabled  bool           `gorm:"default:false" json:"low_power_enabled"` // 是否启用低功耗
	Status           string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive/configuring
	Config           JSON           `gorm:"type:jsonb" json:"config"`            // 其他配置
	CreatedBy        uint           `gorm:"index" json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (BLEMeshNetwork) TableName() string {
	return "ble_mesh_networks"
}

// ===== BLE Mesh 节点 =====

// BLEMeshNode BLE Mesh节点
type BLEMeshNode struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	NetworkID        uint           `gorm:"index;not null" json:"network_id"`
	DeviceID         uint           `gorm:"index" json:"device_id"`              // 关联设备
	NodeAddress      int            `gorm:"not null" json:"node_address"`        // Mesh单播地址
	UUID             string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"uuid"`
	DeviceType       string         `gorm:"type:varchar(64)" json:"device_type"` // proxy/friend/low_power/normal
	ElementCount     int            `gorm:"default:1" json:"element_count"`    // 元素数量
	VendorID         int            `gorm:"default:0" json:"vendor_id"`
	ProductID        int            `gorm:"default:0" json:"product_id"`
	Version          string         `gorm:"type:varchar(32)" json:"version"`
	FirmwareVersion  string         `gorm:"type:varchar(64)" json:"firmware_version"`
	MacAddress       string         `gorm:"type:varchar(32)" json:"mac_address"`
	RSSI             int            `gorm:"default:-100" json:"rssi"`           // 信号强度
	BatteryLevel     int            `gorm:"default:100" json:"battery_level"`   // 电池电量(%)
	PowerState       string         `gorm:"type:varchar(20)" json:"power_state"` // on/off/sleep
	SequenceNumber   int            `gorm:"default:0" json:"sequence_number"`
	LastHeartbeat    *time.Time     `json:"last_heartbeat"`
	HeartbeatPeriod   int            `gorm:"default:0" json:"heartbeat_period"` // 心跳周期(秒)
	PubKey           JSON           `gorm:"type:jsonb" json:"pub_key"`          // 公钥信息
	Sensors          JSON           `gorm:"type:jsonb" json:"sensors"`         // 传感器数据
	Groups           StringArray    `gorm:"type:text" json:"groups"`           // 所属组
	Subscriptions    StringArray    `gorm:"type:text" json:"subscriptions"`    // 订阅地址
	Status           string         `gorm:"type:varchar(20);default:'offline'" json:"status"` // offline/provisioning/provisioned/active/error
	ErrorCode        string         `gorm:"type:varchar(32)" json:"error_code"`
	OnlineDuration   int            `gorm:"default:0" json:"online_duration"` // 在线时长(秒)
	TotalMessages    int64          `gorm:"default:0" json:"total_messages"`  // 总消息数
	FailedMessages   int64          `gorm:"default:0" json:"failed_messages"` // 失败消息数
	Location         JSON           `gorm:"type:jsonb" json:"location"`       // 位置信息
	LastSeenAt       *time.Time     `json:"last_seen_at"`
	FirstSeenAt      time.Time      `json:"first_seen_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (BLEMeshNode) TableName() string {
	return "ble_mesh_nodes"
}

// ===== RTOS 配置 =====

// RTOSConfig RTOS系统配置
type RTOSConfig struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	ConfigKey         string         `gorm:"type:varchar(128);uniqueIndex;not null" json:"config_key"`
	ConfigName        string         `gorm:"type:varchar(128);not null" json:"config_name"`
	DeviceType        string         `gorm:"type:varchar(64)" json:"device_type"`    // 适用设备类型
	FirmwareVersion   string         `gorm:"type:varchar(64)" json:"firmware_version"`
	OSVersion         string         `gorm:"type:varchar(32)" json:"os_version"`     // FreeRTOS/TencentOS/AliOS
	KernelVersion     string         `gorm:"type:varchar(32)" json:"kernel_version"`
	SchedulerType     string         `gorm:"type:varchar(32)" json:"scheduler_type"` // priority/preemptive/cooperative
	SchedulerQuantumMs int           `gorm:"default:10" json:"scheduler_quantum_ms"` // 调度时间片
	TaskConfig        JSON           `gorm:"type:jsonb" json:"task_config"`    // 任务配置列表
	MemoryConfig      JSON           `gorm:"type:jsonb" json:"memory_config"` // 内存配置
	StackSizeDefault  int            `gorm:"default:4096" json:"stack_size_default"` // 默认栈大小
	HeapSizeTotal     int            `gorm:"default:65536" json:"heap_size_total"`    // 总堆大小
	HeapSizeFree      int            `gorm:"default:0" json:"heap_size_free"`
	MinFreeStackKB    int            `gorm:"default:512" json:"min_free_stack_kb"`
	InterruptConfig   JSON           `gorm:"type:jsonb" json:"interrupt_config"` // 中断配置
	TimerConfig       JSON           `gorm:"type:jsonb" json:"timer_config"`     // 定时器配置
	MutexConfig       JSON           `gorm:"type:jsonb" json:"mutex_config"`     // 互斥锁配置
	SemaphoreConfig   JSON           `gorm:"type:jsonb" json:"semaphore_config"` // 信号量配置
	MessageQueueConfig JSON          `gorm:"type:jsonb" json:"message_queue_config"`
	NetworkConfig     JSON           `gorm:"type:jsonb" json:"network_config"`   // 网络配置
	PowerConfig       JSON           `gorm:"type:jsonb" json:"power_config"`     // 电源管理配置
	PeripheralConfig  JSON           `gorm:"type:jsonb" json:"peripheral_config"` // 外设配置
	BootConfig        JSON           `gorm:"type:jsonb" json:"boot_config"`     // 启动配置
	SecurityConfig    JSON           `gorm:"type:jsonb" json:"security_config"` // 安全配置
	DebugConfig       JSON           `gorm:"type:jsonb" json:"debug_config"`    // 调试配置
	PerformanceTuning JSON           `gorm:"type:jsonb" json:"performance_tuning"` // 性能调优
	IsDefault         bool           `gorm:"default:false" json:"is_default"`  // 是否默认配置
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	ApplyCount        int            `gorm:"default:0" json:"apply_count"`      // 应用次数
	Description       string         `gorm:"type:text" json:"description"`
	CreatedBy         uint           `gorm:"index" json:"created_by"`
	ApprovedBy       *uint          `gorm:"index" json:"approved_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (RTOSConfig) TableName() string {
	return "rtos_configs"
}
