package main

import (
	"log"
	"os"
	"strings"

	"mdm-backend/controllers"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/mqtt"
	plugins "mdm-backend/plugins"
	"mdm-backend/services"
	"mdm-backend/utils"
	"mdm-backend/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置全局数据库实例
	models.SetDB(db)

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&models.Device{},
		&models.DeviceShadow{},
		&models.PetProfile{},
		&models.OTAPackage{},
		&models.OTADeployment{},
		&models.OTAProgress{},
		&models.CommandHistory{},
		// 系统表
		&models.SysUser{},
		&models.SysRole{},
		&models.SysMenu{},
		&models.SysPermission{},
		&models.SysRolePermission{},
		&models.SysUserRole{},
		&models.SysUserExt{},
		&models.SysDictionary{},
		&models.SysConfig{},
		&models.SysOperationLog{},
		&models.SysLoginLog{},
		// 知识库表
		&models.Knowledge{},
		// 告警表
		&models.DeviceAlertRule{},
		&models.DeviceAlert{},
		// 地理围栏表
		&models.GeofenceRule{},
		&models.GeofenceAlert{},
		// 告警通知表
		&models.AlertNotification{},
		// 告警设置表
		&models.AlertSettings{},
		&models.NotificationChannel{},
		// Sprint 11: 通知日志和告警历史
		&models.NotificationLog{},
		&models.AlertHistory{},
		// 合规表
		&models.CompliancePolicy{},
		&models.ComplianceViolation{},
		// 策略管理表
		&models.PolicyConfig{},
		&models.Policy{},
		&models.PolicyBinding{},
		// 会员表
		&models.MemberOrder{},
		&models.MemberUpgradeRecord{},
		&models.Member{},
		&models.MemberCard{},
		&models.MemberCardGroup{},
		&models.MemberLevel{},
		&models.MemberUpgradeRule{},
		&models.MemberTag{},
		&models.MemberTagRecord{},
		&models.Coupon{},
		&models.CouponGrant{},
		&models.Promotion{},
		&models.Store{},
		&models.PointsRule{},
		&models.MemberPointsRecord{},
		&models.MemberOperationRecord{},
		&models.TempMember{},

		// 应用管理表
		&models.App{},
		&models.AppVersion{},
		&models.AppDistribution{},
		&models.AppInstallRecord{},
		&models.AppLicense{},
		// 通知表
		&models.Notification{},
		&models.NotificationTemplate{},
		&models.Announcement{},
		// 租户相关表
		&models.Tenant{},
		&models.TenantQuota{},
		&models.Plan{},
		// 租户申请表
		&models.TenantApplication{},
		&models.ApprovalHistory{},
		// 基准岗位模板
		&models.PositionTemplate{},

		// 权限系统表（多租户版本）
		&models.Menu{},
		&models.ApiPermission{},
		&models.Role{},
		&models.PermissionGroup{},
		&models.RoleMenu{},
		&models.RoleApiPermission{},
		&models.RolePermissionGroup{},
		// 数据权限表
		&models.DataScope{},
		// 活动日志表
		&models.ActivityLog{},
		&models.LoginLog{},
		// Sprint 12: LDAP/AD 配置表
		&models.LDAPConfig{},
		&models.LDAPUserMapping{},
		&models.LDAPGroupRoleMapping{},
		// Sprint 12: 证书管理表
		&models.Certificate{},
		// Sprint 12: 设备安全表
		&models.WipeHistory{},
		// Sprint 12: 数据权限表
		&models.DataPermissionRule{},
		&models.UserDataPermission{},
		// Sprint 13: 多区域数据库架构
		&models.Region{},
		&models.RegionalNode{},
		// Sprint 13: 时区支持
		&models.TimezoneConfig{},
		// Sprint 13: 数据驻留
		&models.DataResidencyRule{},
		// Sprint 21: 内容生态 - 表情包市场
		&models.EmoticonCategory{},
		&models.Emoticon{},
		&models.EmoticonPurchase{},
		// Sprint 21: 内容生态 - 动作资源库
		&models.CustomAction{},
		&models.ActionMarket{},
		// Sprint 21: 内容生态 - 声音定制
		&models.VoiceConfig{},
		&models.VoiceMarketItem{},
		// Sprint 22: 移动端后端支持 - App Token & Push
		&models.AppToken{},
		&models.AppRefreshToken{},
		&models.AppPush{},
		// Sprint 22: 微信小程序后端
		&models.MiniAppDevice{},
		&models.MiniAppQRCodeBind{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 注册租户隔离 GORM 插件
	tenantPlugin := &plugins.TenantScopePlugin{}
	if err := tenantPlugin.Initialize(db); err != nil {
		log.Fatalf("Failed to initialize tenant plugin: %v", err)
	}

	// 注册数据权限范围 GORM 插件
	dataScopePlugin := &plugins.DataScopePlugin{}
	if err := dataScopePlugin.Initialize(db); err != nil {
		log.Fatalf("Failed to initialize data scope plugin: %v", err)
	}

	// 初始化 Redis
	redisClient, err := utils.InitRedis()
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	// 设置全局 Redis 客户端供其他模块使用
	utils.SetGlobalRedisClient(redisClient)

	// 初始化 MQTT（传入告警回调和合规检查回调）
	alertCallback := func(deviceID string, data map[string]interface{}) {
		controllers.CheckAlerts(db, deviceID, data)
	}
	complianceCallback := func(db *gorm.DB, deviceID string, data map[string]interface{}) {
		controllers.CheckCompliance(db, deviceID, data)
	}
	geofenceCallback := func(db *gorm.DB, deviceID string, lat, lng float64, alertType string) {
		mqtt.CheckGeofence(db, deviceID, lat, lng, alertType)
	}
	mqttHandler, err := mqtt.InitMQTT(db, redisClient, alertCallback, complianceCallback, geofenceCallback)
	if err != nil {
		log.Fatalf("Failed to initialize MQTT: %v", err)
	}
	defer mqttHandler.Disconnect(0)
	// 设置全局 MQTT 客户端供 CommandController 使用
	mqtt.SetGlobalMQTTClient(mqttHandler)

	// OTA 后台 Worker：定期检查待下发的 OTA 任务
	otaWorker := services.NewOTAWorker(db, mqttHandler)
	controllers.SetOTAWorkerRef(otaWorker)
	go otaWorker.Start()

	// 注册通知服务（支持邮件/Webhook/站内通知）
	controllers.RegisterNotificationService(services.SendAlertNotifications)

	// 初始化 Gin 路由
	r := gin.Default()

	// CORS 白名单中间件
	// 优先从环境变量 CORS_ALLOWED_ORIGINS 读取（逗号分隔），为空时使用默认值
	var allowedOrigins []string
	if extra := os.Getenv("CORS_ALLOWED_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	} else {
		// 开发环境默认值
		allowedOrigins = []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174",
		}
	}
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 注册认证路由 (不需要 JWT)
	authCtrl := &controllers.AuthController{DB: db, Redis: redisClient}
	r.POST("/api/v1/auth/login", authCtrl.Login)

	// JWT 中间件
	r.Use(middleware.JWTAuth())

	// 用户上下文中间件：从 JWT 提取用户ID/OrgID 存入 Context
	r.Use(middleware.UserContext())

	// 租户上下文中间件：从 JWT 解析 tenant_id
	r.Use(middleware.TenantContext())

	// 数据权限范围中间件：自动根据用户角色过滤数据
	r.Use(plugins.DataScopeMiddleware(db))

	// 操作日志中间件
	r.Use(middleware.OperationLog(db))

	// 注册业务路由
	controllers.RegisterRoutes(r, db, redisClient)

	// 注册系统管理路由
	sys := r.Group("/api/v1")
	dictCtrl := &controllers.DictController{DB: db}
	logCtrl := &controllers.LogController{DB: db}

	// 租户管理路由（超管）
	tenantCtrl := &controllers.TenantController{DB: db}
	adminGroup := sys.Group("/admin")
	tenantCtrl.RegisterTenantRoutes(adminGroup)

	// 租户申请审批路由（/api/v1/tenant-approvals）
	tenantApprovalCtrl := &controllers.TenantApprovalController{DB: db}
	tenantApprovalCtrl.RegisterRoutes(sys)

	{
		// 权限系统路由（多租户版本）
		menuCtrl := &controllers.TenantMenuController{DB: db}
		apiPermCtrl := &controllers.ApiPermissionController{DB: db}
		newRoleCtrl := &controllers.NewRoleController{DB: db}
		permGroupCtrl := &controllers.PermissionGroupController{DB: db}

		sys.GET("/menus", menuCtrl.List)
		sys.GET("/menus/:id", menuCtrl.Get)
		sys.POST("/menus", menuCtrl.Create)
		sys.PUT("/menus/:id", menuCtrl.Update)
		sys.DELETE("/menus/:id", menuCtrl.Delete)

		sys.GET("/api-permissions", apiPermCtrl.List)
		sys.POST("/api-permissions", apiPermCtrl.Create)
		sys.PUT("/api-permissions/:id", apiPermCtrl.Update)
		sys.DELETE("/api-permissions/:id", apiPermCtrl.Delete)
		sys.POST("/api-permissions/import", apiPermCtrl.Import)
		sys.GET("/api-permissions/export", apiPermCtrl.Export)

		sys.GET("/roles", newRoleCtrl.List)
		sys.POST("/roles", newRoleCtrl.Create)
		sys.PUT("/roles/:id", newRoleCtrl.Update)
		sys.DELETE("/roles/:id", newRoleCtrl.Delete)
		sys.GET("/roles/:id/permissions", newRoleCtrl.GetPermissions)
		sys.PUT("/roles/:id/permissions", newRoleCtrl.SetPermissions)

		sys.GET("/permission-groups", permGroupCtrl.List)
		sys.GET("/permission-groups/:id", permGroupCtrl.Get)
		sys.POST("/permission-groups", permGroupCtrl.Create)
		sys.PUT("/permission-groups/:id", permGroupCtrl.Update)
		sys.DELETE("/permission-groups/:id", permGroupCtrl.Delete)

		sys.GET("/menus/tree", menuCtrl.List) // 复用 List 返回树形
		sys.GET("/dicts/:type", dictCtrl.GetDictByType)
		sys.GET("/logs/operations", logCtrl.GetOperationLogs)
		sys.GET("/logs/login", logCtrl.GetLoginLogs)

		// 告警管理
		alertCtrl := &controllers.AlertController{DB: db}
		sys.GET("/alerts/rules", alertCtrl.GetRules)
		sys.POST("/alerts/rules", alertCtrl.CreateRule)
		sys.PUT("/alerts/rules/:id", alertCtrl.UpdateRule)
		sys.DELETE("/alerts/rules/:id", alertCtrl.DeleteRule)
		sys.GET("/alerts", alertCtrl.GetAlerts)
		sys.POST("/alerts/:id/confirm", alertCtrl.ConfirmAlert)
		sys.POST("/alerts/:id/resolve", alertCtrl.ResolveAlert)
		sys.POST("/alerts/:id/ignore", alertCtrl.IgnoreAlert)
		sys.POST("/alerts/batch/confirm", alertCtrl.BatchConfirmAlerts)
		sys.POST("/alerts/batch/resolve", alertCtrl.BatchResolveAlerts)
		sys.GET("/alerts/:id/notifications", alertCtrl.GetAlertNotifications)

		// 告警规则别名路由（/api/v1/alert-rules，供前端使用）
		sys.GET("/alert-rules", alertCtrl.GetRules)
		sys.POST("/alert-rules", alertCtrl.CreateRule)
		sys.PUT("/alert-rules/:id", alertCtrl.UpdateRule)
		sys.DELETE("/alert-rules/:id", alertCtrl.DeleteRule)
		sys.POST("/alert-rules/batch-delete", alertCtrl.BatchDeleteAlertRules)

		// 地理围栏管理
		sys.GET("/geofence/rules", alertCtrl.GetGeofenceRules)
		sys.POST("/geofence/rules", alertCtrl.CreateGeofenceRule)
		sys.PUT("/geofence/rules/:id", alertCtrl.UpdateGeofenceRule)
		sys.DELETE("/geofence/rules/:id", alertCtrl.DeleteGeofenceRule)
		sys.GET("/geofence/alerts", alertCtrl.GetGeofenceAlerts)

		// 告警设置
		alertSettingsCtrl := &controllers.AlertSettingsController{DB: db}
		sys.GET("/alerts/settings", alertSettingsCtrl.GetAlertSettings)
		sys.PUT("/alerts/settings", alertSettingsCtrl.UpdateAlertSettings)

		// 通知渠道配置
		notifChannelCtrl := &controllers.NotificationChannelController{DB: db}
		sys.GET("/notification-channels", notifChannelCtrl.ListChannels)
		sys.POST("/notification-channels", notifChannelCtrl.CreateChannel)
		sys.GET("/notification-channels/:id", notifChannelCtrl.GetChannel)
		sys.PUT("/notification-channels/:id", notifChannelCtrl.UpdateChannel)
		sys.DELETE("/notification-channels/:id", notifChannelCtrl.DeleteChannel)
		sys.POST("/notification-channels/:id/toggle", notifChannelCtrl.ToggleChannel)
		sys.POST("/notification-channels/:id/test", notifChannelCtrl.TestChannel)

		// Sprint 11: 告警历史管理
		alertHistoryCtrl := &controllers.AlertHistoryController{DB: db}
		sys.GET("/alerts/history", alertHistoryCtrl.GetAlertHistory)
		sys.GET("/alerts/history/:id", alertHistoryCtrl.GetAlertHistoryByID)
		sys.POST("/alerts/history/archive", alertHistoryCtrl.ArchiveAlert)

		// Dashboard 统计（使用独立的 DashboardController）
		dashboardCtrl := &controllers.DashboardController{DB: db}
		sys.GET("/dashboard/stats", dashboardCtrl.GetStats)
		sys.GET("/dashboard/stats/simple", dashboardCtrl.GetStatsSimple)
		sys.GET("/dashboard/activity-summary", dashboardCtrl.GetActivitySummary)

		// 活动日志（审计日志）
		activityLogCtrl := &controllers.ActivityLogController{DB: db}
		loginLogCtrl := &controllers.LoginLogController{DB: db}
		sys.GET("/activity-logs", activityLogCtrl.List)
		sys.GET("/activity-logs/statistics", activityLogCtrl.GetStatistics)
		sys.GET("/activity-logs/:id", activityLogCtrl.Get)
		sys.GET("/login-logs", loginLogCtrl.List)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// WebSocket 实时通知路由
	websocket.GetHub() // 初始化Hub
	r.GET("/ws/notifications", controllers.HandleWebSocket)

	// 知识库路由
	knowledgeCtrl := &controllers.KnowledgeController{}
	apiV1 := r.Group("/api/v1")
	knowledgeCtrl.RegisterRoutesWithDB(apiV1, db)

	// 宠物控制台路由
	petConsoleCtrl := &controllers.PetConsoleController{}
	petConsoleCtrl.RegisterRoutes(apiV1)

	// MiniClaw路由
	miniClawCtrl := &controllers.MiniClawController{}
	miniClawCtrl.RegisterRoutes(apiV1)

	// 通知路由
	notifCtrl := &controllers.NotificationController{DB: db}
	notifCtrl.RegisterRoutes(apiV1)

	// 补充通知路由（push 和 batch-delete）
	apiV1.POST("/notifications/push", notifCtrl.PushNotification)
	apiV1.POST("/notifications/batch-delete", notifCtrl.BatchDeleteNotifications)

	// 补充公告路由（withdraw）
	apiV1.POST("/announcements/:id/withdraw", notifCtrl.WithdrawAnnouncement)
	apiV1.GET("/announcements/:id", notifCtrl.GetAnnouncement)

	// 策略配置别名路由（/api/v1/policy-configs，供前端使用）
	policyCtrlExtra := &controllers.PolicyController{DB: db}
	complianceCtrlExtra := &controllers.ComplianceController{DB: db}
	apiV1.GET("/policy-configs", policyCtrlExtra.ListConfigs)
	apiV1.POST("/policy-configs", policyCtrlExtra.CreateConfig)
	apiV1.PUT("/policy-configs/:id", policyCtrlExtra.UpdateConfig)
	apiV1.DELETE("/policy-configs/:id", policyCtrlExtra.DeleteConfig)

	// 合规规则别名路由（/api/v1/compliance-rules，供前端使用）
	apiV1.GET("/compliance-rules", complianceCtrlExtra.ListRules)
	apiV1.POST("/compliance-rules", complianceCtrlExtra.CreateRule)
	apiV1.PUT("/compliance-rules/:id", complianceCtrlExtra.UpdateRule)
	apiV1.DELETE("/compliance-rules/:id", complianceCtrlExtra.DeleteRule)

	// 合规策略 API（/api/v1/compliance/policies，标准 REST 端点）
	apiV1.GET("/compliance/policies", complianceCtrlExtra.ListRules)
	apiV1.POST("/compliance/policies", complianceCtrlExtra.CreateRule)
	apiV1.PUT("/compliance/policies/:id", complianceCtrlExtra.UpdateRule)
	apiV1.DELETE("/compliance/policies/:id", complianceCtrlExtra.DeleteRule)

	// ============ 数据导入导出路由 ============
	importExportCtrl := &controllers.ImportExportController{DB: db}
	// 设备导入导出
	apiV1.GET("/devices/export", importExportCtrl.ExportDevices)
	apiV1.POST("/devices/import", importExportCtrl.ImportDevices)
	// 会员导入导出
	apiV1.GET("/members/export", importExportCtrl.ExportMembers)
	apiV1.POST("/members/import", importExportCtrl.ImportMembers)
	// 活动日志导出
	apiV1.GET("/activity-logs/export", importExportCtrl.ExportActivityLogs)

	// ============ 批量操作路由 ============
	// 设备批量操作
	apiV1.POST("/devices/batch-delete", importExportCtrl.BatchDeleteDevices)
	apiV1.POST("/devices/batch-status", importExportCtrl.BatchUpdateDeviceStatus)
	// 会员批量操作
	apiV1.POST("/members/batch-delete", importExportCtrl.BatchDeleteMembers)

	// ============ Sprint 12: LDAP/AD 路由 ============
	ldapCtrl := &controllers.LDAPController{DB: db}
	apiV1.GET("/ldap/config", ldapCtrl.GetLDAPConfig)
	apiV1.PUT("/ldap/config", ldapCtrl.UpdateLDAPConfig)
	apiV1.POST("/ldap/test", ldapCtrl.TestLDAPConnection)
	apiV1.GET("/ldap/users", ldapCtrl.GetLDAPUsers)
	apiV1.POST("/ldap/sync", ldapCtrl.SyncLDAPUsers)
	apiV1.GET("/ldap/groups", ldapCtrl.GetLDAPGroups)
	apiV1.POST("/ldap/groups/mapping", ldapCtrl.SetGroupRoleMapping)
	apiV1.GET("/ldap/group-mappings", ldapCtrl.GetGroupRoleMappings)

	// ============ Sprint 12: 证书管理路由 ============
	certCtrl := &controllers.CertificateController{DB: db}
	apiV1.GET("/certificates", certCtrl.ListCertificates)
	apiV1.POST("/certificates", certCtrl.CreateCertificate)
	apiV1.GET("/certificates/stats", certCtrl.GetCertificateStats)
	apiV1.GET("/certificates/expiring", certCtrl.GetExpiringCertificates)
	apiV1.POST("/certificates/validate", certCtrl.ValidateCertificate)
	apiV1.GET("/certificates/:id", certCtrl.GetCertificate)
	apiV1.PUT("/certificates/:id", certCtrl.UpdateCertificate)
	apiV1.DELETE("/certificates/:id", certCtrl.DeleteCertificate)
	apiV1.POST("/certificates/:id/revoke", certCtrl.RevokeCertificate)
	apiV1.POST("/certificates/upload", certCtrl.UploadCertificateFile)
	apiV1.GET("/certificates/:id/download", certCtrl.DownloadCertificate)

	// ============ Sprint 12: 设备安全路由 (锁定/擦除) ============
	deviceSecurityCtrl := &controllers.DeviceSecurityController{DB: db}
	apiV1.POST("/devices/:device_id/lock", deviceSecurityCtrl.LockDevice)
	apiV1.POST("/devices/:device_id/unlock", deviceSecurityCtrl.UnlockDevice)
	apiV1.POST("/devices/:device_id/wipe", deviceSecurityCtrl.WipeDevice)
	apiV1.POST("/devices/:device_id/wipe/confirm", deviceSecurityCtrl.ConfirmWipe)
	apiV1.POST("/devices/:device_id/wipe/token", deviceSecurityCtrl.GenerateWipeConfirmToken)
	apiV1.GET("/devices/:device_id/wipe-history", deviceSecurityCtrl.GetWipeHistory)
	apiV1.GET("/devices/:device_id/lock-status", deviceSecurityCtrl.GetDeviceLockStatus)

	// ============ Sprint 12: 数据权限路由 ============
	dataPermCtrl := &controllers.DataPermissionController{DB: db}
	apiV1.GET("/data-permissions/rules", dataPermCtrl.ListDataPermissionRules)
	apiV1.POST("/data-permissions/rules", dataPermCtrl.CreateDataPermissionRule)
	apiV1.PUT("/data-permissions/rules/:id", dataPermCtrl.UpdateDataPermissionRule)
	apiV1.DELETE("/data-permissions/rules/:id", dataPermCtrl.DeleteDataPermissionRule)
	apiV1.GET("/data-permissions/roles/:role_id", dataPermCtrl.GetRoleDataPermissions)
	apiV1.PUT("/data-permissions/roles/:role_id", dataPermCtrl.UpdateRoleDataPermissions)
	apiV1.GET("/data-permissions/users/:user_id", dataPermCtrl.GetUserDataPermissions)
	apiV1.PUT("/data-permissions/users/:user_id", dataPermCtrl.UpdateUserDataPermissions)
	apiV1.GET("/data-permissions/columns", dataPermCtrl.GetColumnPermissions)
	apiV1.POST("/data-permissions/validate", dataPermCtrl.ValidatePermissionExpression)

	// ============ Sprint 21: 内容生态 - 表情包市场 ============
	emoticonCtrl := &controllers.EmoticonController{DB: db}
	emoticonCtrl.RegisterEmoticonRoutes(apiV1)
	apiV1.POST("/emoticons/categories", emoticonCtrl.CreateCategory)

	// ============ Sprint 21: 内容生态 - 动作资源库 ============
	actionMarketCtrl := &controllers.ActionMarketController{DB: db}
	actionMarketCtrl.RegisterActionMarketRoutes(apiV1)

	// ============ Sprint 21: 内容生态 - 声音定制 ============
	voiceCtrl := &controllers.VoiceController{DB: db}
	voiceCtrl.RegisterVoiceRoutes(apiV1)

	// ============ Sprint 22: 移动端后端支持 - App API ============
	appCtrl := &controllers.AppController{DB: db, Redis: redisClient}
	appAuthGroup := r.Group("/api/v1/app/auth")
	{
		// App Token API（不需要JWT）
		appAuthGroup.POST("/token", appCtrl.GetAppToken)
		appAuthGroup.POST("/refresh", appCtrl.RefreshAppToken)
	}
	// App API（需要JWT）
	appGroup := r.Group("/api/v1/app")
	appGroup.Use(middleware.JWTAuth())
	appGroup.Use(middleware.UserContext())
	{
		// App Push API
		appGroup.POST("/push", appCtrl.SendAppPush)
		appGroup.GET("/push/history", appCtrl.GetPushHistory)
		// App 设备状态
		appGroup.GET("/device/:device_id/status", appCtrl.GetAppDeviceStatus)
		// App 设备控制
		appGroup.POST("/device/:device_id/command", appCtrl.SendAppDeviceCommand)
	}

	// ============ Sprint 22: 微信小程序后端 API ============
	miniAppCtrl := &controllers.MiniAppController{DB: db, Redis: redisClient}
	miniAppGroup := r.Group("/api/v1/miniapp")
	// 小程序API使用小程序专用认证中间件（从X-OpenID头解析openid）
	miniAppGroup.Use(controllers.MiniAppAuthMiddleware())
	{
		miniAppGroup.GET("/devices", miniAppCtrl.GetMiniAppDevices)
		miniAppGroup.POST("/bind", miniAppCtrl.BindDevice)
		miniAppGroup.POST("/unbind", miniAppCtrl.UnbindDevice)
		miniAppGroup.GET("/qrcode", miniAppCtrl.GenerateQRCode)
		miniAppGroup.GET("/device/:device_id/status", miniAppCtrl.GetDeviceStatus)
		miniAppGroup.POST("/device/:device_id/command", miniAppCtrl.SendDeviceCommand)
		miniAppGroup.POST("/push", miniAppCtrl.SendPush)
	}

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("MDM Backend starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
