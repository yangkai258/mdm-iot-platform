package models

import (
	"time"

	"gorm.io/gorm"
)

// FlowProcess BPMN流程定义
type FlowProcess struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ProcessName  string         `gorm:"size:200" json:"process_name"`      // 流程名称
	ProcessKey   string         `gorm:"size:50;uniqueIndex" json:"process_key"` // 流程唯一标识
	BpmnXML      string         `gorm:"type:text" json:"bpmn_xml"`         // BPMN XML内容
	Status       int            `gorm:"default:1" json:"status"`           // 状态: 1草稿 2已发布 3已下线
	Version      int            `gorm:"default:1" json:"version"`          // 版本号
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`     // 租户ID
	Description  string         `gorm:"size:500" json:"description"`       // 描述
	Category     string         `gorm:"size:50" json:"category"`           // 流程分类
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
