package models

// 权限点定义
// 格式: resource:action
const (
	// 租户权限
	PermissionTenantView   = "tenant:view"   // 查看租户
	PermissionTenantManage = "tenant:manage" // 管理租户

	// 用户权限
	PermissionUserView   = "user:view"   // 查看用户
	PermissionUserManage = "user:manage" // 管理用户

	// 设备权限
	PermissionDeviceView   = "device:view"   // 查看设备
	PermissionDeviceManage = "device:manage" // 管理设备
	PermissionDeviceControl = "device:control" // 控制设备

	// OTA权限
	PermissionOTAView   = "ota:view"   // 查看OTA
	PermissionOTADeploy = "ota:deploy" // 部署OTA

	// 告警权限
	PermissionAlertView    = "alert:view"    // 查看告警
	PermissionAlertManage = "alert:manage"  // 管理告警

	// 会员权限
	PermissionMemberView   = "member:view"   // 查看会员
	PermissionMemberManage = "member:manage" // 管理会员

	// 策略权限
	PermissionPolicyView   = "policy:view"   // 查看策略
	PermissionPolicyManage = "policy:manage" // 管理策略

	// 通知权限
	PermissionNotificationView   = "notification:view"   // 查看通知
	PermissionNotificationManage = "notification:manage" // 管理通知

	// 应用权限
	PermissionAppView   = "app:view"   // 查看应用
	PermissionAppManage = "app:manage" // 管理应用

	// 系统权限
	PermissionSystemView   = "system:view"   // 查看系统
	PermissionSystemManage = "system:manage" // 管理系统
	PermissionRoleManage    = "role:manage"  // 角色管理
	PermissionUserRole      = "user:role"   // 用户角色

	// 数据权限
	PermissionDataExport = "data:export" // 数据导出
	PermissionDataImport = "data:import" // 数据导入
)

// PermissionGroup 权限分组
type PermissionGroup struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" gorm:"type:jsonb"`
	Status      int      `json:"status"`
}

// AllPermissions 所有权限点列表
var AllPermissions = []string{
	PermissionTenantView,
	PermissionTenantManage,
	PermissionUserView,
	PermissionUserManage,
	PermissionDeviceView,
	PermissionDeviceManage,
	PermissionDeviceControl,
	PermissionOTAView,
	PermissionOTADeploy,
	PermissionAlertView,
	PermissionAlertManage,
	PermissionMemberView,
	PermissionMemberManage,
	PermissionPolicyView,
	PermissionPolicyManage,
	PermissionNotificationView,
	PermissionNotificationManage,
	PermissionAppView,
	PermissionAppManage,
	PermissionSystemView,
	PermissionSystemManage,
	PermissionRoleManage,
	PermissionUserRole,
	PermissionDataExport,
	PermissionDataImport,
}

// PermissionGroups 权限分组定义
var PermissionGroups = []PermissionGroup{
	{
		ID:   1,
		Name: "租户管理",
		Code: "tenant",
		Description: "租户相关权限",
		Permissions: []string{PermissionTenantView, PermissionTenantManage},
		Status: 1,
	},
	{
		ID:   2,
		Name: "用户管理",
		Code: "user",
		Description: "用户相关权限",
		Permissions: []string{PermissionUserView, PermissionUserManage, PermissionUserRole},
		Status: 1,
	},
	{
		ID:   3,
		Name: "设备管理",
		Code: "device",
		Description: "设备相关权限",
		Permissions: []string{PermissionDeviceView, PermissionDeviceManage, PermissionDeviceControl},
		Status: 1,
	},
	{
		ID:   4,
		Name: "OTA管理",
		Code: "ota",
		Description: "OTA升级相关权限",
		Permissions: []string{PermissionOTAView, PermissionOTADeploy},
		Status: 1,
	},
	{
		ID:   5,
		Name: "告警管理",
		Code: "alert",
		Description: "告警相关权限",
		Permissions: []string{PermissionAlertView, PermissionAlertManage},
		Status: 1,
	},
	{
		ID:   6,
		Name: "会员管理",
		Code: "member",
		Description: "会员相关权限",
		Permissions: []string{PermissionMemberView, PermissionMemberManage},
		Status: 1,
	},
	{
		ID:   7,
		Name: "策略管理",
		Code: "policy",
		Description: "策略相关权限",
		Permissions: []string{PermissionPolicyView, PermissionPolicyManage},
		Status: 1,
	},
	{
		ID:   8,
		Name: "通知管理",
		Code: "notification",
		Description: "通知相关权限",
		Permissions: []string{PermissionNotificationView, PermissionNotificationManage},
		Status: 1,
	},
	{
		ID:   9,
		Name: "应用管理",
		Code: "app",
		Description: "应用相关权限",
		Permissions: []string{PermissionAppView, PermissionAppManage},
		Status: 1,
	},
	{
		ID:   10,
		Name: "系统管理",
		Code: "system",
		Description: "系统设置和角色权限",
		Permissions: []string{PermissionSystemView, PermissionSystemManage, PermissionRoleManage},
		Status: 1,
	},
	{
		ID:   11,
		Name: "数据操作",
		Code: "data",
		Description: "数据导入导出",
		Permissions: []string{PermissionDataExport, PermissionDataImport},
		Status: 1,
	},
}


