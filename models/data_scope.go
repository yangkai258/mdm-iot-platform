package models

import (
	"strconv"
	"strings"
)

// DataScopeType 数据权限范围类型
type DataScopeType int

const (
	DataScopeAll        DataScopeType = 1 // 全部数据
	DataScopeOrg       DataScopeType = 2 // 本部门数据
	DataScopeOrgAndChildren DataScopeType = 3 // 本部门及下属数据
	DataScopeSelf      DataScopeType = 4 // 仅本人数据
	DataScopeCustom   DataScopeType = 5 // 自定义数据范围
)

// DataScope 数据权限配置
type DataScope struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	RoleID    uint          `gorm:"index" json:"role_id"`        // 角色ID
	ScopeType DataScopeType `gorm:"default:1" json:"scope_type"` // 权限范围类型
	DeptIDs   string        `gorm:"size:500" json:"dept_ids"`   // 自定义部门ID列表，逗号分隔
	StoreIDs  string        `gorm:"size:500" json:"store_ids"`  // 自定义店铺ID列表
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

// DataScopeRole 数据权限角色（查询结果）
type DataScopeRole struct {
	RoleID     uint   `gorm:"column:role_id" json:"role_id"`
	ScopeType  DataScopeType `gorm:"column:scope_type" json:"scope_type"`
	DeptIDs    string `gorm:"column:dept_ids" json:"dept_ids"`
	StoreIDs   string `gorm:"column:store_ids" json:"store_ids"`
}

// ParseScopeValue 解析自定义数据范围的值
// scopeType: 权限类型 (DataScopeCustom = 5 时使用)
// value: ID列表，逗号分隔
// 返回: 字段名, ID列表
func ParseScopeValue(scopeType DataScopeType, value string) (string, []uint) {
	var ids []uint
	for _, v := range strings.Split(value, ",") {
		if id, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
			ids = append(ids, uint(id))
		}
	}
	// 根据 scopeType 判断字段名
	// DataScopeCustom: scope_value 格式为 "dept_ids:1,2,3" 或 "store_ids:1,2,3"
	field := "dept_ids" // 默认
	return field, ids
}
