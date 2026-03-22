package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ RTOS 性能监控模型 ============

// RTOSStats RTOS 性能统计数据
type RTOSStats struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	DeviceID           string         `gorm:"type:varchar(36);index;not null" json:"device_id"`
	CPUUsage           float64        `gorm:"type:decimal(5,2);default:0" json:"cpu_usage"`       // CPU 使用率 0-100%
	TaskCount          int            `gorm:"default:0" json:"task_count"`                        // 任务数量
	Uptime             int64          `gorm:"default:0" json:"uptime"`                            // 运行时间(秒)
	IRQCount           int64          `gorm:"default:0" json:"irq_count"`                         // 中断次数
	ContextSwitchCount int64          `gorm:"default:0" json:"context_switch_count"`              // 上下文切换次数
	Frequency          int            `gorm:"default:0" json:"frequency"`                         // CPU 频率(MHz)
	InterruptLevel     int            `gorm:"default:0" json:"interrupt_level"`                   // 当前中断优先级
	TickRate           int            `gorm:"default:0" json:"tick_rate"`                         // 系统 Tick 率(Hz)
	IdleTaskUsage      float64        `gorm:"type:decimal(5,2);default:0" json:"idle_task_usage"` // 空闲任务 CPU 使用率
	RecordedAt         time.Time      `gorm:"index" json:"recorded_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RTOSStats) TableName() string {
	return "rtos_stats"
}

// RTOSMemory RTOS 内存使用数据
type RTOSMemory struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	DeviceID          string         `gorm:"type:varchar(36);index;not null" json:"device_id"`
	TotalHeap         int64          `gorm:"default:0" json:"total_heap"`                            // 总堆内存(bytes)
	UsedHeap          int64          `gorm:"default:0" json:"used_heap"`                             // 已用堆内存(bytes)
	FreeHeap          int64          `gorm:"default:0" json:"free_heap"`                             // 空闲堆内存(bytes)
	MinFreeHeap       int64          `gorm:"default:0" json:"min_free_heap"`                         // 历史最小空闲堆内存(bytes)
	MemoryFragPercent float64        `gorm:"type:decimal(5,2);default:0" json:"memory_frag_percent"` // 内存碎片率(%)
	PSRAMTotal        int64          `gorm:"default:0" json:"psram_total"`                           // PSRAM 总大小
	PSRAMUsed         int64          `gorm:"default:0" json:"psram_used"`                            // PSRAM 已用
	FlashTotal        int64          `gorm:"default:0" json:"flash_total"`                           // Flash 总大小
	FlashUsed         int64          `gorm:"default:0" json:"flash_used"`                            // Flash 已用
	RecordedAt        time.Time      `gorm:"index" json:"recorded_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RTOSMemory) TableName() string {
	return "rtos_memory"
}

// RTOSTask RTOS 任务信息
type RTOSTask struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	DeviceID      string         `gorm:"type:varchar(36);index;not null" json:"device_id"`
	TaskName      string         `gorm:"type:varchar(64);not null" json:"task_name"`      // 任务名称
	Priority      int            `gorm:"default:0" json:"priority"`                       // 优先级(数值越小优先级越高)
	StackSize     int            `gorm:"default:0" json:"stack_size"`                     // 栈大小(bytes)
	StackUsed     int            `gorm:"default:0" json:"stack_used"`                     // 栈已用(bytes)
	State         string         `gorm:"type:varchar(20);default:'running'" json:"state"` // 状态: running/ready/blocked/suspended/deleted
	CPUUsage      float64        `gorm:"type:decimal(5,2);default:0" json:"cpu_usage"`    // 该任务 CPU 使用率
	LastWakeTime  int64          `gorm:"default:0" json:"last_wake_time"`                 // 上次唤醒时间(系统 ticks)
	Runtime       int64          `gorm:"default:0" json:"runtime"`                        // 累计运行时间(ms)
	RecordVersion int            `gorm:"default:1" json:"record_version"`                 // 记录版本号
	RecordedAt    time.Time      `gorm:"index" json:"recorded_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RTOSTask) TableName() string {
	return "rtos_tasks"
}

// ============ 设备性能历史模型 ============

// DevicePerformanceHistory 设备性能历史数据
type DevicePerformanceHistory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeviceID    string         `gorm:"type:varchar(36);index;not null" json:"device_id"`
	MetricType  string         `gorm:"type:varchar(20);not null;index" json:"metric_type"` // cpu/memory/network/disk/battery/temperature
	MetricName  string         `gorm:"type:varchar(64);not null" json:"metric_name"`       // 具体指标名
	MetricValue float64        `gorm:"type:decimal(15,4);not null" json:"metric_value"`    // 指标值
	Unit        string         `gorm:"type:varchar(16)" json:"unit"`                       // 单位: %, MB, GB, ms, mV, °C
	Tags        string         `gorm:"type:text" json:"tags"`                              // 标签 JSON，如 {"core":0}
	SnapshotID  string         `gorm:"type:varchar(36);index" json:"snapshot_id"`          // 快照ID，用于批量查询
	RecordedAt  time.Time      `gorm:"index" json:"recorded_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DevicePerformanceHistory) TableName() string {
	return "device_performance_history"
}

// ============ 固件优化配置模型 ============

// FirmwareOptimizationConfig 固件优化配置
type FirmwareOptimizationConfig struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	DeviceID         string         `gorm:"type:varchar(36);index;not null" json:"device_id"`
	OptimizationType string         `gorm:"type:varchar(32);not null" json:"optimization_type"` // power_save/performance/balance/network_low_latency
	ConfigJSON       string         `gorm:"type:jsonb;default:'{}'" json:"config_json"`         // 详细配置 JSON
	IsEnabled        bool           `gorm:"default:true" json:"is_enabled"`                     // 是否启用
	IsApplied        bool           `gorm:"default:false" json:"is_applied"`                    // 是否已应用(设备确认)
	AppliedAt        *time.Time     `json:"applied_at"`                                         // 应用时间
	AppliedVersion   string         `gorm:"type:varchar(32)" json:"applied_version"`            // 固件版本(应用时的版本)
	CreatedBy        string         `gorm:"type:varchar(64)" json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (FirmwareOptimizationConfig) TableName() string {
	return "firmware_optimization_configs"
}

// OptimizationConfigItem 单个优化项配置
type OptimizationConfigItem struct {
	Enabled       bool        `json:"enabled"`
	Value         interface{} `json:"value,omitempty"`
	ThresholdHigh float64     `json:"threshold_high,omitempty"`
	ThresholdLow  float64     `json:"threshold_low,omitempty"`
	IntervalMs    int         `json:"interval_ms,omitempty"`
}

// PowerSaveConfig 节能模式配置
type PowerSaveConfig struct {
	CPUFrequencyMax     int  `json:"cpu_frequency_max"`     // 最大 CPU 频率(MHz), 0=不限
	CPUFrequencyMin     int  `json:"cpu_frequency_min"`     // 最小 CPU 频率(MHz)
	WiFiPowerSave       bool `json:"wifi_power_save"`       // WiFi 节能模式
	BluetoothPowerSave  bool `json:"bluetooth_power_save"`  // 蓝牙节能模式
	ScreenBrightness    int  `json:"screen_brightness"`     // 屏幕亮度 0-100
	DeepSleepTimeout    int  `json:"deep_sleep_timeout"`    // 深度休眠超时(秒), 0=禁用
	TaskWatchdogTimeout int  `json:"task_watchdog_timeout"` // 任务看门狗超时(秒)
}

// PerformanceConfig 性能模式配置
type PerformanceConfig struct {
	CPUFrequencyMax  int  `json:"cpu_frequency_max"`     // 最大 CPU 频率(MHz)
	WiFiPerformance  bool `json:"wifi_performance"`      // WiFi 性能模式
	BluetoothPerf    bool `json:"bluetooth_performance"` // 蓝牙高性能模式
	CacheEnabled     bool `json:"cache_enabled"`         // 启用缓存
	PreFetchEnabled  bool `json:"prefetch_enabled"`      // 启用预取
	MultiCoreEnabled bool `json:"multi_core_enabled"`    // 多核模式
}

// BalanceConfig 均衡模式配置
type BalanceConfig struct {
	CPUFrequency      int  `json:"cpu_frequency"`       // CPU 频率(MHz)
	WiFiBalance       bool `json:"wifi_balance"`        // WiFi 均衡模式
	AutoTickRate      bool `json:"auto_tick_rate"`      // 自动 Tick 率调节
	TickRateHz        int  `json:"tick_rate_hz"`        // Tick 率(Hz)
	TaskPriorityBoost bool `json:"task_priority_boost"` // 任务优先级提升
}
