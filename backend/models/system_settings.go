package models

import (
	"time"
)

// SystemConfig 系统配置
type SystemConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Group     string    `gorm:"type:varchar(50);index;not null" json:"group"`        // 配置分组
	ConfigKey string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`   // 配置键
	Value     string    `gorm:"type:text" json:"value"`                              // 配置值
	Type      string    `gorm:"type:varchar(20);default:'string'" json:"type"`       // 值类型: string, number, boolean, json
	Label     string    `gorm:"type:varchar(100)" json:"label"`                       // 显示标签
	Remark    string    `gorm:"type:varchar(255)" json:"remark"`                      // 备注
	Sort      int       `gorm:"default:0" json:"sort"`                                // 排序
	Status    int       `gorm:"default:1" json:"status"`                              // 状态: 0禁用 1启用
	TenantID  string    `gorm:"index" json:"tenant_id"`                               // 租户ID(空字符串表示全局配置)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemConfig) TableName() string { return "system_configs" }
