package models

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

// DataScopeRole 数据权限与角色的关联表
type DataScopeRole struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RoleID     uint   `gorm:"index" json:"role_id"`
	DataScopeID uint  `gorm:"index" json:"data_scope_id"`
	CreatedAt  string `json:"created_at"`
}

// ParseScopeValue 解析自定义数据范围的值
func ParseScopeValue(value string) []uint {
	var ids []uint
	for _, v := range strings.Split(value, ",") {
		if id, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}
