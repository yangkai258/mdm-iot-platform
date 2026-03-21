package models

import (
	"time"

	"gorm.io/gorm"
)

// DataScopeType 数据权限范围类型
type DataScopeType string

const (
	DataScopeAll          DataScopeType = "all"           // 全部数据
	DataScopeOrg          DataScopeType = "org"           // 本部门数据
	DataScopeOrgAndChildren DataScopeType = "org_and_children" // 本部门及下级数据
	DataScopeSelf         DataScopeType = "self"          // 仅本人数据
	DataScopeCustom       DataScopeType = "custom"        // 自定义数据范围
)

// DataScope 数据权限配置表（角色维度）
type DataScope struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoleID    uint           `gorm:"index;not null" json:"role_id"`                 // 角色ID
	ScopeType DataScopeType  `gorm:"type:varchar(50);default:'all'" json:"scope_type"` // 权限范围类型
	ScopeValue string        `gorm:"type:varchar(255)" json:"scope_value"`          // custom时存储具体条件，如 "dept_id IN (1,2,3)"
	TenantID  string         `gorm:"type:varchar(64);index" json:"tenant_id"`        // 租户ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (DataScope) TableName() string {
	return "data_scopes"
}

// DataScopeRole 数据权限-角色关联范围（多角色合并用）
type DataScopeRole struct {
	RoleID      uint          `gorm:"primaryKey"`
	ScopeType   DataScopeType `gorm:"type:varchar(50)"`
	ScopeValue  string        `gorm:"type:varchar(255)"`
	DeptIDs     []uint        // 自定义部门ID列表（解析 ScopeValue 而来）
	OrgIDs      []uint        // 自定义机构ID列表
}

// MergeScopes 合并多个角色的数据权限范围，返回最大范围
// 范围从大到小: all > org_and_children > org > self > custom
func MergeScopes(scopes []DataScopeRole) DataScopeRole {
	if len(scopes) == 0 {
		return DataScopeRole{ScopeType: DataScopeAll}
	}
	if len(scopes) == 1 {
		return scopes[0]
	}

	// 优先级：all > org_and_children > org > self > custom
	priority := func(t DataScopeType) int {
		switch t {
		case DataScopeAll:
			return 5
		case DataScopeOrgAndChildren:
			return 4
		case DataScopeOrg:
			return 3
		case DataScopeSelf:
			return 2
		case DataScopeCustom:
			return 1
		default:
			return 0
		}
	}

	best := scopes[0]
	for i := 1; i < len(scopes); i++ {
		if priority(scopes[i].ScopeType) > priority(best.ScopeType) {
			best = scopes[i]
		}
	}
	return best
}

// ParseScopeValue 解析 ScopeValue 字符串为具体条件
// 格式如: "dept_id:1,2,3" 或 "org_id:1,2" 或 "user_id:100"
func ParseScopeValue(scopeType DataScopeType, scopeValue string) (field string, values []uint) {
	if scopeValue == "" {
		return
	}
	parts := splitOnce(scopeValue, ":")
	if len(parts) != 2 {
		return
	}
	field = parts[0]
	for _, v := range splitComma(parts[1]) {
		if v == "" {
			continue
		}
		id := parseUint(v)
		if id > 0 {
			values = append(values, id)
		}
	}
	return
}

func splitOnce(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
	}
	return []string{s}
}

func splitComma(s string) []string {
	var result []string
	var cur string
	for _, ch := range s {
		if ch == ',' {
			result = append(result, cur)
			cur = ""
		} else {
			cur += string(ch)
		}
	}
	if cur != "" {
		result = append(result, cur)
	}
	return result
}

func parseUint(s string) uint {
	var v uint
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		v = v*10 + uint(c-'0')
	}
	return v
}
