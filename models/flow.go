package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// FlowNodeType 流程节点类型
type FlowNodeType string

const (
	FlowNodeStart    FlowNodeType = "start"
	FlowNodeApproval FlowNodeType = "approval"
	FlowNodeEnd      FlowNodeType = "end"
)

// FlowInstanceStatus 流程实例状态
type FlowInstanceStatus string

const (
	FlowStatusRunning FlowInstanceStatus = "running"
	FlowStatusApproved FlowInstanceStatus = "approved"
	FlowStatusRejected FlowInstanceStatus = "rejected"
	FlowStatusCanceled FlowInstanceStatus = "canceled"
)

// FlowTaskStatus 流程任务状态
type FlowTaskStatus string

const (
	TaskStatusPending   FlowTaskStatus = "pending"
	TaskStatusApproved  FlowTaskStatus = "approved"
	TaskStatusRejected  FlowTaskStatus = "rejected"
)

// FlowNode 流程节点定义（JSONB存储在flow_definition中）
type FlowNode struct {
	NodeID     string        `json:"node_id"`
	NodeName   string        `json:"node_name"`
	NodeType   FlowNodeType `json:"node_type"`
	ApproverID *uint         `json:"approver_id,omitempty"` // 审批人用户ID
	NextNodeID string        `json:"next_node_id"`         // 下一个节点ID，end时为空
}

// FlowDefinition 流程定义表
type FlowDefinition struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	Version     int            `gorm:"type:smallint;default:1" json:"version"`
	Nodes       json.RawMessage `gorm:"type:jsonb;not null" json:"nodes"` // []FlowNode，顺序存储
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// FlowInstance 流程实例表
type FlowInstance struct {
	ID               uint               `gorm:"primaryKey" json:"id"`
	InstanceUUID     string              `gorm:"type:varchar(36);uniqueIndex;not null" json:"instance_uuid"`
	FlowDefinitionID uint                `gorm:"index;not null" json:"flow_definition_id"`
	FlowName         string              `gorm:"type:varchar(100)" json:"flow_name"`
	BusinessKey      string              `gorm:"type:varchar(100);index" json:"business_key"` // 关联业务ID，如设备ID
	BusinessType     string              `gorm:"type:varchar(50)" json:"business_type"`      // 业务类型，如device_enrollment
	Status           FlowInstanceStatus `gorm:"type:varchar(20);default:'running'" json:"status"`
	CurrentNodeID    string              `gorm:"type:varchar(50)" json:"current_node_id"`
	InitiatorID      uint                `gorm:"index" json:"initiator_id"`
	TenantID         string              `gorm:"type:uuid;index" json:"tenant_id"`
	FormData         json.RawMessage     `gorm:"type:jsonb" json:"form_data"` // 申请表单数据
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	DeletedAt        gorm.DeletedAt      `gorm:"index" json:"-"`
}

// FlowTask 流程任务表
type FlowTask struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	InstanceID    uint            `gorm:"index;not null" json:"instance_id"`
	InstanceUUID  string          `gorm:"type:varchar(36);index" json:"instance_uuid"`
	NodeID        string          `gorm:"type:varchar(50);index" json:"node_id"`
	NodeName      string          `gorm:"type:varchar(100)" json:"node_name"`
	ApproverID    uint            `gorm:"index" json:"approver_id"` // 当前待审批人
	Status        FlowTaskStatus  `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Remark        string          `gorm:"type:varchar(500)" json:"remark"` // 审批意见
	CompletedAt   *time.Time      `json:"completed_at,omitempty"`
	TenantID      string          `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}
