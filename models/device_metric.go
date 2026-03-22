package models

import (
	"time"
)

// DeviceMetric 设备指标数据
type DeviceMetric struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DeviceID    string    `gorm:"type:varchar(36);index" json:"device_id"`
	MetricType  string    `gorm:"type:varchar(50);index" json:"metric_type"` // cpu, memory, battery, temperature, etc.
	MetricName  string    `gorm:"type:varchar(100)" json:"metric_name"`
	MetricValue float64   `gorm:"type:decimal(10,4)" json:"metric_value"`
	Unit        string    `gorm:"type:varchar(20)" json:"unit"` // %, °C, MB, etc.
	Timestamp   time.Time `gorm:"index" json:"timestamp"`
	Tags        string    `gorm:"type:text" json:"tags"`          // JSON格式的标签
	ExtraData   JSON      `gorm:"type:jsonb" json:"extra_data"`  // 额外的指标数据
	TenantID    string    `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (DeviceMetric) TableName() string {
	return "device_metrics"
}

// DeviceMetricResponse 设备指标响应
type DeviceMetricResponse struct {
	ID          uint      `json:"id"`
	DeviceID    string    `json:"device_id"`
	MetricType  string    `json:"metric_type"`
	MetricName  string    `json:"metric_name"`
	MetricValue float64   `json:"metric_value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Tags        string    `json:"tags,omitempty"`
	ExtraData   JSON      `json:"extra_data,omitempty"`
}

// ToResponse 转换为响应结构
func (m *DeviceMetric) ToResponse() *DeviceMetricResponse {
	return &DeviceMetricResponse{
		ID:          m.ID,
		DeviceID:    m.DeviceID,
		MetricType:  m.MetricType,
		MetricName:  m.MetricName,
		MetricValue: m.MetricValue,
		Unit:        m.Unit,
		Timestamp:   m.Timestamp,
		Tags:        m.Tags,
		ExtraData:   m.ExtraData,
	}
}
