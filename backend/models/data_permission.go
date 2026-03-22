package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// DataPermissionRule 数据权限规则表
type DataPermissionRule struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	RuleName      string          `gorm:"type:varchar(100);not null" json:"rule_name"` // 规则名称
	ResourceType  string          `gorm:"type:varchar(50);not null" json:"resource_type"` // device/pet/member/org
	RuleType      string          `gorm:"type:varchar(20);not null" json:"rule_type"` // row/column
	ResourceIDs   StringArray      `gorm:"type:text[]" json:"resource_ids"`     // 资源ID列表，空表示全部
	PermissionExpr JSONMap        `gorm:"type:jsonb" json:"permission_expr"`   // 权限表达式
	Priority      int             `gorm:"default:0" json:"priority"`             // 优先级（越大越高）
	IsActive      bool            `gorm:"default:true" json:"is_active"`        // 是否启用
	Description   string          `gorm:"type:text" json:"description"`         // 描述
	TenantID      string          `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedBy     uint            `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (DataPermissionRule) TableName() string {
	return "data_permission_rules"
}

// StringArray PostgreSQL text[] 类型
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil || len(a) == 0 {
		return "{}", nil
	}
	// GORM 会自动处理数组，这里返回原始值
	return a, nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		// 解析 PostgreSQL 数组格式 {a,b,c}
		str := string(v)
		if str == "{}" || str == "" {
			*a = StringArray{}
			return nil
		}
		// 简单解析 {xxx,yyy}
		str = str[1 : len(str)-1] // 去掉 {}
		if str == "" {
			*a = StringArray{}
			return nil
		}
		parts := splitAndTrim(str, ",")
		*a = StringArray(parts)
		return nil
	case string:
		if v == "{}" || v == "" {
			*a = StringArray{}
			return nil
		}
		parts := splitAndTrim(v[1:len(v)-1], ",")
		*a = StringArray(parts)
		return nil
	}
	return nil
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		part = trimQuotes(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimQuotes(s string) string {
	s = trimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// JSONMap JSONB 类型
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = JSONMap{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// UserDataPermission 用户数据权限配置（行级/列级权限）
type UserDataPermission struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index;not null" json:"user_id"`
	RoleID       uint       `gorm:"index" json:"role_id"` // 可选，角色级别配置
	ResourceType string     `gorm:"type:varchar(50);not null" json:"resource_type"`
	RuleType     string     `gorm:"type:varchar(20);not null" json:"rule_type"` // row/column
	ColumnFields StringArray `gorm:"type:text[]" json:"column_fields"`   // 列级权限：可访问的字段列表
	DataScope    JSONMap    `gorm:"type:jsonb" json:"data_scope"`        // 行级权限：数据范围表达式
	FilterExpr   string     `gorm:"type:text" json:"filter_expr"`       // 自定义过滤表达式
	Priority     int        `gorm:"default:0" json:"priority"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	TenantID     string     `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedBy    uint       `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (UserDataPermission) TableName() string {
	return "user_data_permissions"
}

// DataPermissionRuleRequest 创建/更新权限规则请求
type DataPermissionRuleRequest struct {
	RuleName       string          `json:"rule_name" binding:"required"`
	ResourceType   string          `json:"resource_type" binding:"required"`
	RuleType       string          `json:"rule_type" binding:"required"` // row/column
	ResourceIDs    []string        `json:"resource_ids"`                // 空表示全部
	PermissionExpr map[string]interface{} `json:"permission_expr"`
	Priority       int             `json:"priority"`
	IsActive       bool            `json:"is_active"`
	Description    string          `json:"description"`
}

// UserDataPermissionRequest 用户数据权限请求
type UserDataPermissionRequest struct {
	ResourceType string                 `json:"resource_type" binding:"required"`
	RuleType    string                 `json:"rule_type" binding:"required"`
	ColumnFields []string               `json:"column_fields"`  // 列级权限
	DataScope   map[string]interface{} `json:"data_scope"`    // 行级权限
	FilterExpr  string                 `json:"filter_expr"`
	Priority    int                    `json:"priority"`
	IsActive    bool                   `json:"is_active"`
}

// ColumnPermission 可配置列级权限字段
type ColumnPermission struct {
	Field     string `json:"field"`      // 字段名
	Label     string `json:"label"`     // 中文标签
	TableName string `json:"table_name"` // 所属表
	DataType  string `json:"data_type"`  // 数据类型
	Sensitive bool   `json:"sensitive"`  // 是否敏感字段
}

// DataPermissionValidateRequest 权限表达式验证请求
type DataPermissionValidateRequest struct {
	Expression string                 `json:"expression" binding:"required"`
	Context   map[string]interface{} `json:"context"` // 模拟上下文数据
}

// DataPermissionValidateResponse 权限表达式验证响应
type DataPermissionValidateResponse struct {
	Valid   bool     `json:"valid"`
	Error   string   `json:"error,omitempty"`
	Result  bool     `json:"result"`  // 模拟执行结果
	Details []string `json:"details"` // 解析详情
}
