package plugins

import (
	"strconv"
	"strings"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataScopePlugin 数据权限范围过滤插件
// 根据用户的数据权限范围自动过滤查询结果
type DataScopePlugin struct{}

// Name 实现 gorm.Plugin 接口
func (p *DataScopePlugin) Name() string {
	return "data_scope"
}

// Initialize 注册 callbacks
func (p *DataScopePlugin) Initialize(db *gorm.DB) error {
	err := db.Callback().Query().Before("gorm:query").Register("data_scope:before_query", dataScopeBeforeQuery)
	if err != nil {
		return err
	}
	return nil
}

// DataScopeContextKey 在 db Session 中传递数据权限信息的 key
const dataScopeCtxKey = "data_scope"

// DataScopeInfo 数据权限范围信息
type DataScopeInfo struct {
	UserID      uint                   // 用户ID
	RoleID      uint                   // 角色ID
	DeptID      uint                   // 部门ID
	TenantID    string                 // 租户ID
	ScopeType   models.DataScopeType   // 权限范围类型
	DeptIDs     []uint                 // 自定义部门ID列表
	StoreIDs    []uint                 // 自定义店铺ID列表
	IsSuperAdmin bool                  // 是否超管
}

// WithDataScope 返回附加了数据权限上下文的 *gorm.DB
func WithDataScope(db *gorm.DB, info *DataScopeInfo) *gorm.DB {
	return db.Session(&gorm.Session{}).Set(dataScopeCtxKey, info)
}

// GetDataScope 从 db Instance 获取数据权限信息
func GetDataScope(db *gorm.DB) *DataScopeInfo {
	if v, ok := db.Get(dataScopeCtxKey); ok {
		if info, ok := v.(*DataScopeInfo); ok {
			return info
		}
	}
	return nil
}

// dataScopeBeforeQuery 在查询前自动注入数据权限过滤条件
func dataScopeBeforeQuery(db *gorm.DB) {
	scopeInfo := GetDataScope(db)
	if scopeInfo == nil {
		return
	}

	// 超管或未设置角色不过滤
	if scopeInfo.IsSuperAdmin || scopeInfo.RoleID == 0 {
		return
	}

	// 判断是否需要数据权限过滤
	tableName := getTableName(db)
	if !needsDataScope(tableName) {
		return
	}

	// 避免重复注入
	if hasDataScopeCondition(db) {
		return
	}

	// 根据权限范围类型注入过滤条件
	switch scopeInfo.ScopeType {
	case models.DataScopeAll:
		// 全部数据，不过滤
		return
	case models.DataScopeOwn:
		// 仅本人数据
		db.Where("create_user_id = ?", scopeInfo.UserID)
	case models.DataScopeDept:
		// 本部门数据
		db.Where("dept_id = ?", scopeInfo.DeptID)
	case models.DataScopeDeptAndSub:
		// 本部门及下属数据
		db.Where("dept_id = ?", scopeInfo.DeptID)
	case models.DataScopeCustom:
		// 自定义数据范围
		applyCustomDataScope(db, scopeInfo)
	}
}

// needsDataScope 判断表是否需要数据权限过滤
func needsDataScope(tableName string) bool {
	// 需要数据权限过滤的业务表
	scopeTables := map[string]bool{
		"devices":            true,
		"members":            true,
		"device_alerts":      true,
		"ota_deployments":    true,
		"notifications":      true,
		"policies":           true,
		"compliance_violations": true,
	}
	name := strings.TrimPrefix(tableName, "public.")
	return scopeTables[name]
}

// hasDataScopeCondition 检查是否已有数据权限过滤条件
func hasDataScopeCondition(db *gorm.DB) bool {
	if db.Statement == nil {
		return false
	}
	conditions := db.Statement.SQL.String()
	keywords := []string{"create_user_id", "dept_id", "data_scope"}
	lowerCond := strings.ToLower(conditions)
	for _, kw := range keywords {
		if strings.Contains(lowerCond, kw) {
			return true
		}
	}
	return false
}

// applyCustomDataScope 应用自定义数据范围过滤
func applyCustomDataScope(db *gorm.DB, info *DataScopeInfo) {
	if len(info.DeptIDs) > 0 {
		db.Where("dept_id IN ?", info.DeptIDs)
		return
	}
	if len(info.StoreIDs) > 0 {
		db.Where("store_id IN ?", info.StoreIDs)
		return
	}
	// 没有配置自定义范围时默认为仅本人
	db.Where("create_user_id = ?", info.UserID)
}

// LoadDataScopeFromDB 从数据库加载用户的数据权限配置
func LoadDataScopeFromDB(db *gorm.DB, userID uint) *DataScopeInfo {
	info := &DataScopeInfo{
		UserID:    userID,
		ScopeType: models.DataScopeAll,
	}

	// 获取用户信息和角色
	var user models.SysUser
	if err := db.First(&user, userID).Error; err != nil {
		return info
	}
	info.RoleID = user.RoleID

	// 获取 SysUserExt 扩展信息
	var ext models.SysUserExt
	if err := db.Where("user_id = ?", userID).First(&ext).Error; err == nil {
		info.DeptID = ext.DeptID
		info.ScopeType = models.DataScopeType(ext.DataScope)
	}

	// 获取角色的数据权限配置
	if info.ScopeType == 0 {
		info.ScopeType = models.DataScopeAll
	}

	return info
}

// parseUintSlice 解析逗号分隔的uint字符串
func parseUintSlice(s string) []uint {
	if s == "" {
		return nil
	}
	var result []uint
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseUint(p, 10, 32); err == nil {
			result = append(result, uint(v))
		}
	}
	return result
}

// DataScopeMiddleware Gin 中间件：为请求自动加载数据权限
func DataScopeMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除不需要数据权限控制的路径
		path := c.Request.URL.Path
		if isExcludedPath(path) {
			c.Next()
			return
		}

		// 获取用户ID
		var userID uint
		if uid, exists := c.Get("user_id"); exists {
			if u, ok := uid.(uint); ok {
				userID = u
			} else if u, ok := uid.(float64); ok {
				userID = uint(u)
			}
		}

		if userID > 0 {
			info := LoadDataScopeFromDB(db, userID)
			// 判断是否超管
			if claims, ok := c.Get("claims"); ok {
				if m, ok := claims.(map[string]interface{}); ok {
					if isSuper, ok := m["is_super_admin"].(bool); ok {
						info.IsSuperAdmin = isSuper
					}
				}
			}
			c.Set("data_scope", info)
		}

		c.Next()
	}
}

// isExcludedPath 判断路径是否不需要数据权限过滤
func isExcludedPath(path string) bool {
	excludedPrefixes := []string{
		"/api/v1/auth",
		"/api/v1/admin",
		"/health",
	}
	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
