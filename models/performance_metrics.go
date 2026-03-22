package models

import (
	"time"
)

// PerformanceMetric 性能指标记录
type PerformanceMetric struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MetricName string   `gorm:"size:100;not null;index" json:"metric_name"`      // 指标名称
	MetricType string   `gorm:"size:50;not null" json:"metric_type"`            // 指标类型: cpu, memory, disk, network, db, cache
	MetricValue float64 `gorm:"not null" json:"metric_value"`                   // 指标值
	Unit       string   `gorm:"size:20" json:"unit"`                            // 单位: %, MB, GB, ms, count
	Tags       string   `gorm:"type:text" json:"tags"`                         // 标签 JSON
	Source     string   `gorm:"size:100" json:"source"`                         // 数据来源: system, application, database
	DeviceID   string   `gorm:"size:64;index" json:"device_id,omitempty"`       // 关联设备ID（可选）
	TenantID   string   `gorm:"size:64;index" json:"tenant_id,omitempty"`      // 租户ID
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalKeys    int64  `json:"total_keys"`     // 总键数量
	MemoryUsed   string `json:"memory_used"`    // 已用内存
	MemoryPeak   string `json:"memory_peak"`    // 峰值内存
	HitRate      string `json:"hit_rate"`       // 命中率
	MissRate     string `json:"miss_rate"`      // 丢失率
	EvictionCount int64 `json:"eviction_count"` // 驱逐数量
	Uptime       string `json:"uptime"`         // 运行时间
}

// DBStats 数据库统计信息
type DBStats struct {
	TotalConnections int    `json:"total_connections"` // 总连接数
	ActiveConnections int    `json:"active_connections"` // 活跃连接数
	IdleConnections  int    `json:"idle_connections"`  // 空闲连接数
	TotalTables      int64  `json:"total_tables"`      // 总表数
	TotalIndexes     int64  `json:"total_indexes"`     // 总索引数
	TotalSize        string `json:"total_size"`         // 总大小
	IndexSize        string `json:"index_size"`         // 索引大小
	TableSize        string `json:"table_size"`         // 表大小
	CacheHitRate     string `json:"cache_hit_rate"`    // 缓存命中率
	QueryPerSecond   float64 `json:"query_per_second"`  // 每秒查询数
	TransactionPerSec float64 `json:"transaction_per_second"` // 每秒事务数
	AvgQueryTime     string  `json:"avg_query_time"`    // 平均查询时间
}

// HealthStatus 健康检查状态
type HealthStatus struct {
	Status    string            `json:"status"`     // healthy, degraded, unhealthy
	Timestamp time.Time        `json:"timestamp"`  // 检查时间
	Duration  string           `json:"duration"`    // 检查耗时
	Components map[string]ComponentStatus `json:"components"` // 各组件状态
}

// ComponentStatus 组件健康状态
type ComponentStatus struct {
	Status  string `json:"status"`  // up, down, degraded
	Latency string `json:"latency"` // 延迟
	Message string `json:"message"` // 附加信息
}

// IndexInfo 索引信息
type IndexInfo struct {
	Schema      string `json:"schema"`      // 模式
	TableName   string `json:"table_name"` // 表名
	IndexName   string `json:"index_name"` // 索引名
	IndexSize   string `json:"index_size"` // 索引大小
	IsUnique    bool   `json:"is_unique"`  // 是否唯一
	IsValid     bool   `json:"is_valid"`   // 是否有效
	Scans       int64  `json:"scans"`      // 扫描次数
}
