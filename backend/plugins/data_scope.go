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

// dataScopeCtxKey 在 db Session 中传递数据权限信息的 key
const dataScopeCtxKey = "data_scope"

// DataScopeInfo 数据权限范围信息（传递给 GORM Session）
type DataScopeInfo struct {
	UserID       uint                   // 用户ID
	RoleID       uint                   // 角色ID
	OrgID        uint                   // 组织/部门ID
	TenantID     string                 // 租户ID
	ScopeType    models.DataScopeType   // 权限范围类型
	OrgIDs       []uint                 // 自定义机构ID列表（org_and_children 时使用）
	CustomField  string                 // 自定义过滤字段名
	CustomValues []uint                 // 自定义过滤值
	IsSuperAdmin bool                   // 是否超管
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
	tableName := getDSTableName(db)
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
	case models.DataScopeSelf:
		// 仅本人数据（使用 create_user_id 或 BindUserID）
		applySelfScope(db, tableName, scopeInfo.UserID)
	case models.DataScopeOrg:
		// 本部门/本机构数据
		applyOrgScope(db, tableName, scopeInfo)
	case models.DataScopeOrgAndChildren:
		// 本部门及下级数据
		applyOrgAndChildrenScope(db, tableName, scopeInfo)
	case models.DataScopeCustom:
		// 自定义数据范围
		applyCustomDataScope(db, tableName, scopeInfo)
	}
}

// needsDataScope 判断表是否需要数据权限过滤
func needsDataScope(tableName string) bool {
	// 需要数据权限过滤的业务表（支持 tenant 隔离）
	scopeTables := map[string]bool{
		"devices":               true,
		"members":               true,
		"device_alerts":         true,
		"ota_deployments":      true,
		"ota_progresses":        true,
		"notifications":         true,
		"policies":              true,
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
	keywords := []string{"create_user_id", "bind_user_id", "org_id", "dept_id", "data_scope"}
	lowerCond := strings.ToLower(conditions)
	for _, kw := range keywords {
		if strings.Contains(lowerCond, kw) {
			return true
		}
	}
	return false
}

// applySelfScope 仅本人数据过滤
func applySelfScope(db *gorm.DB, tableName string, userID uint) {
	switch tableName {
	case "devices":
		// 设备表使用 BindUserID
		db.Where("bind_user_id = ?", strconv.FormatUint(uint64(userID), 10))
	default:
		db.Where("create_user_id = ?", userID)
	}
}

// applyOrgScope 本机构数据过滤
func applyOrgScope(db *gorm.DB, tableName string, info *DataScopeInfo) {
	if info.OrgID == 0 {
		// 无机构信息时降级为仅本人
		applySelfScope(db, tableName, info.UserID)
		return
	}
	db.Where("org_id = ?", info.OrgID)
}

// applyOrgAndChildrenScope 本机构及下级数据过滤
func applyOrgAndChildrenScope(db *gorm.DB, tableName string, info *DataScopeInfo) {
	if info.OrgID == 0 {
		applySelfScope(db, tableName, info.UserID)
		return
	}
	// 先获取所有下级机构ID
	orgIDs := getChildOrgIDs(db, info.OrgID, info.TenantID)
	if len(orgIDs) == 0 {
		orgIDs = []uint{info.OrgID}
	}
	db.Where("org_id IN ?", orgIDs)
}

// applyCustomDataScope 应用自定义数据范围过滤
func applyCustomDataScope(db *gorm.DB, tableName string, info *DataScopeInfo) {
	if info.CustomField != "" && len(info.CustomValues) > 0 {
		db.Where(info.CustomField+" IN ?", info.CustomValues)
		return
	}
	// 没有配置自定义范围时默认为仅本人
	applySelfScope(db, tableName, info.UserID)
}

// getChildOrgIDs 获取指定部门的所有下级部门ID（用于 org_and_children）
func getChildOrgIDs(db *gorm.DB, parentID uint, tenantID string) []uint {
	var deptIDs []uint
	query := db.Model(&models.Department{}).Where("parent_id = ?", parentID)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Pluck("id", &deptIDs).Error; err != nil {
		return nil
	}
	// 递归获取所有下级（最多3层防止无限递归）
	var allIDs []uint
	for _, deptID := range deptIDs {
		allIDs = append(allIDs, deptID)
		if childIDs := getChildOrgIDs(db, deptID, tenantID); len(childIDs) > 0 {
			allIDs = append(allIDs, childIDs...)
		}
	}
	return allIDs
}

// LoadDataScopeFromDB 从数据库加载用户的数据权限配置
// 支持多角色合并：多角色时取最大范围
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
	info.TenantID = user.TenantID

	// 获取 SysUserExt 扩展信息（OrgID）
	var ext models.SysUserExt
	if err := db.Where("user_id = ?", userID).First(&ext).Error; err == nil {
		info.OrgID = ext.DeptID // SysUserExt.DeptID 即为 OrgID
	}

	// 获取用户所有角色的数据权限配置（多角色合并）
	var roleScopes []models.DataScopeRole
	err := db.Table("data_scopes").
		Select("role_id, scope_type, scope_value").
		Where("role_id = ? AND deleted_at IS NULL", user.RoleID).
		Scan(&roleScopes).Error

	if err != nil || len(roleScopes) == 0 {
		// 无配置，默认全部
		info.ScopeType = models.DataScopeAll
		return info
	}

	// 多角色合并：取最大范围
	merged := mergeDataScopes(roleScopes)
	info.ScopeType = merged.ScopeType

	// 解析 ScopeValue
	if merged.ScopeValue != "" {
		field, values := models.ParseScopeValue(merged.ScopeType, merged.ScopeValue)
		info.CustomField = field
		info.CustomValues = values
	}

	// 如果是 org_and_children，收集所有下级机构
	if info.ScopeType == models.DataScopeOrgAndChildren && info.OrgID > 0 {
		info.OrgIDs = getChildOrgIDs(db, info.OrgID, info.TenantID)
		if len(info.OrgIDs) == 0 {
			info.OrgIDs = []uint{info.OrgID}
		}
	}

	return info
}

// DataScopeRoleRow 用于 Scan 到结构体
type DataScopeRoleRow struct {
	RoleID     uint   `gorm:"column:role_id"`
	ScopeType  string `gorm:"column:scope_type"`
	ScopeValue string `gorm:"column:scope_value"`
}

// mergeDataScopes 合并多个角色的数据权限
func mergeDataScopes(rows []models.DataScopeRole) models.DataScopeRole {
	if len(rows) == 0 {
		return models.DataScopeRole{ScopeType: models.DataScopeAll}
	}
	if len(rows) == 1 {
		return rows[0]
	}

	// 优先级：all > org_and_children > org > self > custom
	priority := func(t models.DataScopeType) int {
		switch t {
		case models.DataScopeAll:
			return 5
		case models.DataScopeOrgAndChildren:
			return 4
		case models.DataScopeOrg:
			return 3
		case models.DataScopeSelf:
			return 2
		case models.DataScopeCustom:
			return 1
		default:
			return 0
		}
	}

	best := rows[0]
	for i := 1; i < len(rows); i++ {
		if priority(rows[i].ScopeType) > priority(best.ScopeType) {
			best = rows[i]
		}
	}
	return best
}

// DataScopeMiddleware Gin 中间件：为请求自动加载数据权限
func DataScopeMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除不需要数据权限控制的路径
		path := c.Request.URL.Path
		if isDSExcludedPath(path) {
			c.Next()
			return
		}

		// 获取用户ID
		userID := getUserIDFromContext(c)
		if userID == 0 {
			c.Next()
			return
		}

		info := LoadDataScopeFromDB(db, userID)

		// 判断是否超管
		if claims, ok := c.Get("claims"); ok {
			if m, ok := claims.(map[string]interface{}); ok {
				if isSuper, ok := m["is_super_admin"].(bool); ok {
					info.IsSuperAdmin = isSuper
				}
				if tenantID, ok := m["tenant_id"]; ok {
					info.TenantID = toStringCtx(tenantID)
				}
				if orgID, ok := m["org_id"]; ok {
					info.OrgID = toUintCtx(orgID)
				}
			}
		}

		c.Set("data_scope", info)
		c.Next()
	}
}

// getUserIDFromContext 从 Gin Context 获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	if uid, exists := c.Get("user_id"); exists {
		if u, ok := uid.(uint); ok {
			return u
		}
		if u, ok := uid.(float64); ok {
			return uint(u)
		}
		if u, ok := uid.(int); ok {
			return uint(u)
		}
	}
	return 0
}

// isDSExcludedPath 判断路径是否不需要数据权限过滤
func isDSExcludedPath(path string) bool {
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

func getDSTableName(db *gorm.DB) string {
	if db.Statement != nil && db.Statement.Table != "" {
		return db.Statement.Table
	}
	if db.Statement != nil && db.Statement.Model != nil {
		return getModelTableNameDS(db.Statement.Model)
	}
	return ""
}

func getModelTableNameDS(model interface{}) string {
	if v, ok := model.(interface{ TableName() string }); ok {
		return v.TableName()
	}
	return ""
}

func toStringCtx(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case float64:
		return strconv.FormatFloat(s, 'f', 0, 64)
	case int:
		return strconv.Itoa(s)
	default:
		return ""
	}
}

func toUintCtx(v interface{}) uint {
	switch val := v.(type) {
	case uint:
		return val
	case float64:
		return uint(val)
	case int:
		return uint(val)
	case string:
		if val == "" {
			return 0
		}
		var result uint
		for _, c := range val {
			if c >= '0' && c <= '9' {
				result = result*10 + uint(c-'0')
			} else {
				break
			}
		}
		return result
	default:
		return 0
	}
}
