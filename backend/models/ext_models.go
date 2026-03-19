package models

import (
	"time"

	"gorm.io/gorm"
)

// SysPermission 权限表
type SysPermission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ParentID    uint           `gorm:"default:0" json:"parent_id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"size:100;uniqueIndex" json:"code"`
	Type        int            `gorm:"default:1" json:"type"` // 1:菜单 2:按钮 3:接口
	Path        string         `gorm:"size:255" json:"path"`
	Component   string         `gorm:"size:255" json:"component"`
	Icon        string         `gorm:"size:50" json:"icon"`
	Sort        int            `gorm:"default:0" json:"sort"`
	Visible     int            `gorm:"default:1" json:"visible"`
	Permission  string         `gorm:"size:100" json:"permission"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// SysRolePermission 角色权限关联
type SysRolePermission struct {
	ID           uint      `gorm:"primaryKey"`
	RoleID       uint      `gorm:"index"`
	PermissionID uint      `gorm:"index"`
	CreatedAt    time.Time `json:"created_at"`
}

// SysUserRole 用户角色关联
type SysUserRole struct {
	ID      uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"index"`
	RoleID uint      `gorm:"index"`
}

// WorkflowProcess 流程定义
type WorkflowProcess struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Category    string         `gorm:"size:50" json:"category"`
	Version     string         `gorm:"size:20" json:"version"`
	ProcessKey  string         `gorm:"size:100;uniqueIndex" json:"process_key"`
	ProcessDef  string         `gorm:"type:jsonb" json:"process_def"`
	Status      int            `gorm:"default:0" json:"status"` // 0:草稿 1:已发布
	Creator     string         `gorm:"size:50" json:"creator"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// WorkflowTask 待办任务
type WorkflowTask struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProcessID   uint           `gorm:"index" json:"process_id"`
	ProcessName string         `gorm:"size:100" json:"process_name"`
	TaskKey     string         `gorm:"size:100" json:"task_key"`
	TaskName    string         `gorm:"size:100" json:"task_name"`
	Assignee    string         `gorm:"size:50" json:"assignee"`
	Owner       string         `gorm:"size:50" json:"owner"`
	Status      int            `gorm:"default:1" json:"status"` // 1:待处理 2:处理中 3:已完成 4:已拒绝
	Priority    int            `gorm:"default:2" json:"priority"` // 1:紧急 2:重要 3:一般
	BusinessKey string         `gorm:"size:100" json:"business_key"`
	FormData    string         `gorm:"type:jsonb" json:"form_data"`
	Comment     string         `gorm:"size:500" json:"comment"`
	EndTime     *time.Time    `json:"end_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// WorkflowHistory 已办任务
type WorkflowHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaskID      uint      `gorm:"index" json:"task_id"`
	ProcessID   uint      `json:"process_id"`
	ProcessName string    `gorm:"size:100" json:"process_name"`
	TaskName    string    `gorm:"size:100" json:"task_name"`
	Assignee    string    `gorm:"size:50" json:"assignee"`
	Action      string    `gorm:"size:20" json:"action"` // approve/reject/transfer
	Comment     string    `gorm:"size:500" json:"comment"`
	Result      string    `gorm:"size:20" json:"result"`
	Duration    int       `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
}

// SysConfig 系统配置
type SysConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ConfigKey string    `gorm:"size:100;uniqueIndex" json:"config_key"`
	ConfigVal string    `gorm:"type:text" json:"config_val"`
	ConfigType string   `gorm:"size:20" json:"config_type"` // string/int/json/boolean
	Remark    string    `gorm:"size:255" json:"remark"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
