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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

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
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	// 支持从环境变量扩展白名单，逗号分隔
	if extra := os.Getenv("CORS_ALLOWED_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
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
	authCtrl := &controllers.AuthController{DB: db}
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

		sys.GET("/dashboard/stats", alertCtrl.GetDashboardStats)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

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
