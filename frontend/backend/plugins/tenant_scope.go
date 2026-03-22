package plugins

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// TenantScopePlugin GORM 插件：为所有业务表查询自动注入 tenant_id 隔离条件
type TenantScopePlugin struct{}

// Name 实现 gorm.Plugin 接口
func (p *TenantScopePlugin) Name() string {
	return "tenant_scope"
}

// Initialize 注册插件，安装 callbacks
func (p *TenantScopePlugin) Initialize(db *gorm.DB) error {
	// 创建查询前的 callback，自动追加 tenant_id 条件
	err := db.Callback().Query().Before("gorm:query").Register("tenant_scope:before_query", beforeQuery)
	if err != nil {
		return err
	}
	// 创建前的 callback，自动填充 tenant_id
	err = db.Callback().Create().Before("gorm:create").Register("tenant_scope:before_create", beforeCreate)
	if err != nil {
		return err
	}
	return nil
}

// TenantAwareTableNames 业务表名列表（必须包含 tenant_id 字段）
// 平台级表（无 tenant_id）不在此列
var TenantAwareTables = []string{
	"sys_users",
	"devices",
	"members",
	"ota_packages",
	"ota_deployments",
	"ota_progresses",
	"device_alerts",
	"device_alert_rules",
	"geofence_alerts",
	"geofence_rules",
	"alert_notifications",
	"notifications",
	"notification_templates",
	"announcements",
	"policies",
	"policy_configs",
	"policy_bindings",
	"compliance_policies",
	"compliance_violations",
	"apps",
	"app_versions",
	"app_distributions",
	"app_install_records",
	"app_licenses",
	"coupons",
	"coupon_grants",
	"promotions",
	"stores",
	"command_histories",
	// 权限系统表（多租户版本）
	"menus",
	"api_permissions",
	"roles",
	"permission_groups",
	"role_menus",
	"role_api_permissions",
	"role_permission_groups",
}

var tenantAwareMap map[string]bool

func init() {
	tenantAwareMap = make(map[string]bool)
	for _, t := range TenantAwareTables {
		tenantAwareMap[t] = true
	}
}

// tenantCtxKey 用于在 db Instance 层面传递 tenant_id 的 context key
const tenantCtxKey = "tenant_scope_tenant_id"

// WithTenantID 返回一个附加了 tenant_id 上下文的 *gorm.DB
func WithTenantID(db *gorm.DB, tenantID string) *gorm.DB {
	return db.Session(&gorm.Session{}).Set(tenantCtxKey, tenantID)
}

// GetTenantID 从 db Instance 获取 tenant_id
func GetTenantID(db *gorm.DB) string {
	if v, ok := db.Get(tenantCtxKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// beforeQuery 在查询前自动追加 tenant_id 条件
func beforeQuery(db *gorm.DB) {
	tenantID := GetTenantID(db)
	if tenantID == "" {
		return // 未设置 tenant_id，跳过（如超管查询）
	}

	// 判断是否为需要隔离的表
	tableName := getTableName(db)
	if !isTenantAware(tableName) {
		return
	}

	// 避免重复注入
	if hasTenantCondition(db) {
		return
	}

	db.Where("tenant_id = ?", tenantID)
}

// beforeCreate 在创建记录前自动填充 tenant_id
func beforeCreate(db *gorm.DB) {
	tenantID := GetTenantID(db)
	if tenantID == "" {
		return
	}

	tableName := getTableName(db)
	if !isTenantAware(tableName) {
		return
	}

	// 通过反射设置 tenant_id 字段（如果模型中有该字段且值为空）
	reflectSetTenantID(db.Statement, tenantID)
}

func getTableName(db *gorm.DB) string {
	if db.Statement != nil && db.Statement.Table != "" {
		return db.Statement.Table
	}
	// 从模型中获取表名
	if db.Statement != nil && db.Statement.Model != nil {
		return getModelTableName(db.Statement.Model)
	}
	return ""
}

func getModelTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}
	if v, ok := model.(interface{ TableName() string }); ok {
		return v.TableName()
	}
	return ""
}

func isTenantAware(tableName string) bool {
	// 支持带 schema 前缀的表名
	name := strings.TrimPrefix(tableName, "public.")
	return tenantAwareMap[name]
}

func hasTenantCondition(db *gorm.DB) bool {
	// 简单判断是否已有 tenant_id 条件
	scope := db.Statement
	if scope == nil {
		return false
	}
	// 检查 Where 条件中是否已包含 tenant_id
	conditions := scope.SQL.String()
	return strings.Contains(strings.ToLower(conditions), "tenant_id")
}

// reflectSetTenantID 通过反射设置模型的 tenant_id 字段
func reflectSetTenantID(stmt *gorm.Statement, tenantID string) {
	if stmt == nil || stmt.Model == nil {
		return
	}

	model := stmt.Model
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return
	}

	// 尝试找 tenant_id 字段（支持 snake_case 和 camelCase）
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")
		gormTag := field.Tag.Get("gorm")

		isTenantField := fieldName == "TenantID" ||
			strings.HasSuffix(fieldName, "TenantID") ||
			strings.Contains(jsonTag, "tenant_id") ||
			strings.Contains(gormTag, "tenant_id")

		if isTenantField {
			fieldValue := v.Field(i)
			// 只在字段值为空（零值）时设置
			if fieldValue.IsZero() {
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(tenantID)
				case reflect.Interface:
					if fieldValue.IsNil() {
						fieldValue.Set(reflect.ValueOf(tenantID))
					}
				}
			}
			return
		}
	}
}
