package models

import (
	"time"
)

// Region 多区域配置
type Region struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	RegionCode string    `json:"region_code" gorm:"uniqueIndex;size:16"` // cn-east/cn-north/us-east/eu-west
	RegionName string    `json:"region_name" gorm:"size:64"`
	RegionType string    `json:"region_type" gorm:"size:32"` // primary/replica/ai-inference
	DBHost     string    `json:"db_host" gorm:"size:256"`
	DBPort     int       `json:"db_port"`
	DBName     string    `json:"db_name" gorm:"size:64"`
	AIHost     string    `json:"ai_host" gorm:"size:256"`
	AIPort     int       `json:"ai_port"`
	Status     string    `json:"status" gorm:"size:20"` // active/inactive/failover
	IsDefault  bool      `json:"is_default" gorm:"default:false"`
	Config     string    `json:"config" gorm:"type:text"` // JSON配置
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RegionalNode 区域节点
type RegionalNode struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	NodeID     string    `json:"node_id" gorm:"uniqueIndex;size:64"`
	RegionCode string    `json:"region_code" gorm:"index;size:16"`
	NodeType   string    `json:"node_type" gorm:"size:32"` // db/redis/mqtt/ai
	NodeHost   string    `json:"node_host" gorm:"size:256"`
	NodePort   int       `json:"node_port"`
	NodeStatus string    `json:"node_status" gorm:"size:20"` // online/offline/maintenance
	Load       float64   `json:"load" gorm:"default:0"`      // 负载百分比
	Config     string    `json:"config" gorm:"type:text"`    // JSON配置
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
