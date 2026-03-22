package models

import (
	"time"
)

// UserDataPermission 用户数据权限
type UserDataPermission struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index" json:"user_id"`                 // 用户ID（角色级别权限为0）
	RoleID       uint       `gorm:"index" json:"role_id"`                 // 角色ID
	ResourceType string     `gorm:"type:varchar(50);index" json:"resource_type"` // 资源类型：device, pet, member, etc.
	RuleType     string     `gorm:"type:varchar(20)" json:"rule_type"`     // 规则类型：all, own, dept, custom
	ColumnFields StringArray `gorm:"type:text" json:"column_fields"`      // 可访问的字段列表
	DataScope   JSON        `gorm:"type:jsonb" json:"data_scope"`         // 数据范围过滤条件
	FilterExpr  string      `gorm:"type:varchar(500)" json:"filter_expr"`  // 自定义过滤表达式
	Priority    int         `gorm:"default:50" json:"priority"`           // 优先级
	IsActive    bool        `gorm:"default:true" json:"is_active"`        // 是否启用
	TenantID    string      `gorm:"type:varchar(50);index" json:"tenant_id"` // 租户ID
	CreatedBy   uint        `gorm:"index" json:"created_by"`             // 创建人
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// TableName 指定表名
func (UserDataPermission) TableName() string {
	return "user_data_permissions"
}

// UserDataPermissionRequest 用户数据权限请求
type UserDataPermissionRequest struct {
	ResourceType string                 `json:"resource_type" binding:"required"`
	RuleType     string                 `json:"rule_type"`
	ColumnFields []string               `json:"column_fields"`
	DataScope    map[string]interface{} `json:"data_scope"`
	FilterExpr   string                 `json:"filter_expr"`
	Priority     int                    `json:"priority"`
	IsActive     bool                   `json:"is_active"`
}

// ColumnPermission 列级权限字段定义
type ColumnPermission struct {
	Field     string `json:"field"`      // 字段名
	Label     string `json:"label"`      // 显示名称
	TableName string `json:"table_name"` // 表名
	DataType  string `json:"data_type"` // 数据类型：string, int, datetime, decimal
	Sensitive bool   `json:"sensitive"`  // 是否敏感
}

// DataPermissionValidateRequest 权限表达式验证请求
type DataPermissionValidateRequest struct {
	Expression string `json:"expression" binding:"required"`
}

// DataPermissionValidateResponse 权限表达式验证响应
type DataPermissionValidateResponse struct {
	Valid   bool     `json:"valid"`
	Result  bool     `json:"result"`
	Details []string `json:"details"`
	Error   string   `json:"error,omitempty"`
}

// DataPermissionRule 数据权限规则
type DataPermissionRule struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	RuleName       string    `gorm:"type:varchar(100)" json:"rule_name"`
	RuleType       string    `gorm:"type:varchar(20);not null" json:"rule_type"` // row, column
	ResourceType   string    `gorm:"type:varchar(50);not null" json:"resource_type"`
	ResourceIDs    StringArray `gorm:"type:text" json:"resource_ids"`          // 资源ID列表
	PermissionExpr JSON       `gorm:"type:jsonb" json:"permission_expr"`
	Description    string     `gorm:"type:varchar(500)" json:"description"`
	Priority       int        `gorm:"default:50" json:"priority"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	TenantID       string     `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedBy      uint       `gorm:"index" json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (DataPermissionRule) TableName() string {
	return "data_permission_rules"
}

// DataPermissionRuleRequest 数据权限规则创建请求
type DataPermissionRuleRequest struct {
	Name           string                 `json:"name" binding:"required"`
	RuleName       string                 `json:"rule_name"`
	RuleType       string                 `json:"rule_type" binding:"required"`
	ResourceType   string                 `json:"resource_type" binding:"required"`
	ResourceIDs    []string               `json:"resource_ids"`
	PermissionExpr map[string]interface{} `json:"permission_expr"`
	Description    string                 `json:"description"`
	Priority       int                    `json:"priority"`
	IsActive       bool                   `json:"is_active"`
}
