package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"mdm-backend/controllers"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/multi_region"
	plugins "mdm-backend/plugins"
	"mdm-backend/services"
	"mdm-backend/timezone"
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

		// 具身智能表
		&models.EmbodiedAIState{},
		&models.ExplorationSession{},
		&models.EnvironmentMap{},

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
		// Sprint 25: 安全与合规 - 审计日志
		&models.AuditLog{},
		&models.EncryptionKey{},
		&models.DataAnonymizationRecord{},
		&models.GDPRRequest{},
		// Sprint 27: 开发者平台
		&models.DeveloperApp{},
		&models.APIKey{},
		// Sprint 28: 数据分析增强
		&models.AnalyticsRecord{},
		&models.ExportJob{},
		&models.CustomReport{},
		// Sprint 30: 性能优化
		&models.PerformanceMetric{},
		// Sprint 31: 国际化扩展
		&models.Translation{},
		// Sprint 29: AI 增强功能
		&models.AIModelConfig{},
		&models.AIModelDeployHistory{},
		&models.AIInference{},
		&models.AITraining{},
		// Sprint 35: AI 伦理与公平性
		&models.FairnessTest{},
		&models.BiasDetection{},
		&models.FairnessMetrics{},
		&models.AIAuditLog{},
		&models.AIAuditReport{},
		// Sprint 33: 模型分片加载
		&models.ModelShard{},
		&models.ModelVersion{},
		&models.DeployShardedModel{},
		// Sprint 34+: 企业应用商店
		&models.StoreApp{},
		&models.StoreAppVersion{},
		&models.StoreInstallation{},
		&models.StoreReview{},
		// Sprint 34: 保险理赔对接
		&models.InsuranceProduct{},
		&models.InsuranceClaim{},
		&models.InsuranceClaimDocument{},
		&models.PetHealthRecord{},
		// Sprint 32: 高级安全功能 (暂未实现)
		// &models.SecuritySession{},
		// &models.TwoFactorAuth{},
		// &models.LoginAttempt{},
		// &models.SecurityAudit{},
		// &models.SecurityReport{},
		// DaaS 设备即服务
		&models.DaaSContract{},
		&models.DaaSDeviceRental{},
		&models.DaaSBilling{},
		// Sprint 35: RTOS 优化和设备端性能监控
		&models.RTOSStats{},
		&models.RTOSMemory{},
		&models.RTOSTask{},
		&models.DevicePerformanceHistory{},
		&models.FirmwareOptimizationConfig{},
		// 宠物寻回网络
		&models.PetLostReport{},
		&models.SightingReport{},
		&models.FinderAlert{},
		// 家庭模式：儿童模式和老人陪伴模式
		&models.ChildrenProfile{},
		&models.ContentFilterRule{},
		&models.UsageLimit{},
		&models.ElderlyProfile{},
		&models.HealthMonitorSetting{},
		&models.ElderlyReminder{},
		// 仿真测试平台
		&models.VirtualPet{},
		&models.SimulationEnvironment{},
		&models.SimulationRun{},
		&models.SimulationMetrics{},
		// Sprint 29: 高级功能模型
		&models.ChildModeConfig{},
		&models.ElderlyCareConfig{},
		&models.FamilyAlbum{},
		&models.FamilyAlbumComment{},
		&models.FamilyAlbumLike{},
		&models.PetVaccination{},
		&models.VaccinationReminder{},
		&models.PetDietRecord{},
		&models.PetFinderReport{},
		&models.PetFinderSighting{},
		&models.PetFinderAlert{},
		// 宠物社交
		&models.PetSocialPost{},
		&models.PetSocialComment{},
		&models.PetSocialFollow{},
		&models.PetSocialLike{},
		// Sprint 17-18: 健康医疗和情感计算
		&models.EmotionRecord{},
		&models.EmotionResponseConfig{},
		&models.EmotionReport{},
		&models.ExerciseRecord{},
		&models.ExerciseSummary{},
		&models.ExerciseTrend{},
		&models.HealthAlert{},
		&models.HealthAlertRule{},
		&models.HealthWarning{},
		&models.ExerciseGoal{},
		&models.VitalRecord{},
		&models.VitalTrend{},
		// Sprint 15-16: 订阅管理
		&models.SubscriptionPlan{},
		&models.UserSubscription{},
		&models.SubscriptionChange{},
		// BLE Mesh 网络
		&models.MeshNetwork{},
		&models.MeshDevice{},
		&models.MeshNetworkMember{},
		// 宠物商店
		&models.Product{},
		&models.ProductCategory{},
		&models.Order{},
		&models.OrderItem{},
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

	// Sprint 28: 数据分析增强路由
	analyticsCtrl := controllers.NewAnalyticsController(db, redisClient)
	analyticsCtrl.RegisterRoutes(apiV1)

	// ============ Sprint 30: 性能优化路由 ============
	perfCtrl := &controllers.PerformanceController{DB: db, Redis: redisClient}
	perfGroup := apiV1.Group("/performance")
	perfCtrl.RegisterPerformanceRoutes(perfGroup)

	// ============ Sprint 35: RTOS 优化和设备端性能监控 ============
	rtosCtrl := &controllers.RTOSPerformanceController{DB: db, Redis: redisClient}
	rtosCtrl.RegisterRTOSPerformanceRoutes(apiV1)

	// 宠物控制台路由
	petConsoleCtrl := &controllers.PetConsoleController{}
	petConsoleCtrl.RegisterRoutes(apiV1)

	// 宠物寻回网络路由
	petFinderCtrl := &controllers.PetFinderCtrl{DB: db}
	petFinderCtrl.RegisterPetFinderRoutes(apiV1)

	// ============ Sprint 34: 保险理赔对接 ============
	insuranceCtrl := &controllers.InsuranceController{DB: db}
	insuranceCtrl.RegisterInsuranceRoutes(apiV1)

	// ============ 订单路由 (PetShopController) ============
	petShopCtrl := &controllers.PetShopController{DB: db}
	petShopCtrl.RegisterRoutes(apiV1)

	// ============ 宠物社交路由 (PetSocialController) ============
	petSocialCtrl := &controllers.PetSocialController{DB: db}
	petSocialCtrl.RegisterRoutes(apiV1)

	// ============ 健康报告路由 (HealthTrackingCtrl) ============
	healthCtrl := &controllers.HealthTrackingCtrl{DB: db}
	healthCtrl.RegisterHealthRoutes(apiV1)

	// ============ 情绪记录路由 (EmotionController) ============
	emotionCtrl := &controllers.EmotionController{DB: db}
	emotionCtrl.RegisterEmotionRoutes(apiV1)

	// ============ 数字孪生路由 (DigitalTwinController) ============
	digitalTwinCtrl := &controllers.DigitalTwinController{DB: db}
	digitalTwinCtrl.RegisterDigitalTwinRoutes(apiV1)

	// ============ 订阅路由 (SubscriptionController) ============
	subscriptionCtrl := &controllers.SubscriptionController{DB: db}
	subsGroup := apiV1.Group("/subscriptions")
	{
		subsGroup.GET("/plans", subscriptionCtrl.ListPlans)
		subsGroup.GET("/plans/:id", subscriptionCtrl.GetPlan)
		subsGroup.POST("/plans", subscriptionCtrl.CreatePlan)
		subsGroup.PUT("/plans/:id", subscriptionCtrl.UpdatePlan)
		subsGroup.DELETE("/plans/:id", subscriptionCtrl.DeletePlan)
		subsGroup.GET("/current", subscriptionCtrl.GetCurrentSubscription)
		subsGroup.POST("/subscribe", subscriptionCtrl.Subscribe)
		subsGroup.GET("/:id", subscriptionCtrl.GetSubscription)
		subsGroup.POST("/:id/cancel", subscriptionCtrl.CancelSubscription)
		subsGroup.POST("/:id/renew", subscriptionCtrl.RenewSubscription)
	}

	// ============ 宠物社交动态路由 ============
	// 注意：pet-social/feed 控制器暂未实现，跳过

	// ============ 家庭模式：儿童模式 + 老人陪伴模式 ============
	familyModeCtrl := &controllers.FamilyModeController{DB: db}
	familyModeCtrl.RegisterFamilyModeRoutes(apiV1)

	// ============ Sprint 29: 高级功能 API ============
	advancedCtrl := &controllers.AdvancedController{DB: db, Redis: redisClient}
	advancedCtrl.RegisterAdvancedRoutes(apiV1)

	// MiniClaw路由
	miniClawCtrl := &controllers.MiniClawController{}
	miniClawCtrl.RegisterRoutes(apiV1)

	// ============ 仿真测试平台 ============
	simulationCtrl := controllers.NewSimulationController(db)
	simulationCtrl.RegisterRoutes(apiV1)

	// ============ Sprint 29: AI 增强功能 ============
	// AI模型管理路由
	aiModelCtrl := controllers.NewAIModelController(db)
	aiModelCtrl.RegisterRoutes(apiV1)

	// AI推理路由
	aiInferenceCtrl := controllers.NewAIInferenceController(db)
	aiInferenceCtrl.RegisterRoutes(apiV1)

	// AI训练任务路由
	aiTrainingCtrl := controllers.NewAITrainingController(db)
	aiTrainingCtrl.RegisterRoutes(apiV1)

	// ============ Sprint 36: 具身智能 ============
	// 具身智能路由
	embodiedAIController := controllers.NewEmbodiedAIController(db)
	embodiedAIGroup := apiV1.Group("/ai/embodied")
	{
		embodiedAIGroup.POST("/perceive", embodiedAIController.Perceive)
		embodiedAIGroup.POST("/spatial", embodiedAIController.Spatial)
		embodiedAIGroup.POST("/explore", embodiedAIController.Explore)
		embodiedAIGroup.POST("/interact", embodiedAIController.Interact)
		embodiedAIGroup.GET("/state", embodiedAIController.GetState)
		embodiedAIGroup.GET("/capabilities", embodiedAIController.GetCapabilities)
		// 探索会话管理
		embodiedAIGroup.GET("/explore/sessions", embodiedAIController.ListExplorationSessions)
		embodiedAIGroup.GET("/explore/sessions/:id", embodiedAIController.GetExplorationSession)
		embodiedAIGroup.POST("/explore/sessions/:id/stop", embodiedAIController.StopExplorationSession)
	}

	// ============ Sprint 33: 模型分片加载 ============
	// 模型分片管理路由
	modelShardCtrl := &controllers.ModelShardController{DB: db}
	modelVersionCtrl := &controllers.ModelVersionController{DB: db}

	// 模型分片路由 /api/v1/ai/models/:id/shards
	aiModelShards := apiV1.Group("/ai/models/:id/shards")
	{
		aiModelShards.GET("", modelShardCtrl.ListShards)
		aiModelShards.POST("", modelShardCtrl.CreateShard)
		aiModelShards.PUT("/:shard_id", modelShardCtrl.UpdateShard)
		aiModelShards.GET("/:shard_id", modelShardCtrl.GetShard)
		aiModelShards.DELETE("/:shard_id", modelShardCtrl.DeleteShard)
		aiModelShards.POST("/:shard_id/verify", modelShardCtrl.VerifyShard)
	}

	// 分片模型部署 /api/v1/ai/models/:id/deploy/sharded
	apiV1.POST("/ai/models/:id/deploy/sharded", modelShardCtrl.DeploySharded)

	// 模型版本路由 /api/v1/ai/models/:id/versions
	aiModelVersions := apiV1.Group("/ai/models/:id/versions")
	{
		aiModelVersions.GET("", modelVersionCtrl.ListVersions)
		aiModelVersions.POST("", modelVersionCtrl.CreateVersion)
		aiModelVersions.GET("/simple", modelVersionCtrl.GetVersionsSimple)
		aiModelVersions.GET("/:version", modelVersionCtrl.GetVersion)
		aiModelVersions.PUT("/:version", modelVersionCtrl.UpdateVersion)
		aiModelVersions.POST("/:version/publish", modelVersionCtrl.PublishVersion)
	}

	// 模型回滚 /api/v1/ai/models/:id/rollback
	apiV1.POST("/ai/models/:id/rollback", modelVersionCtrl.Rollback)

	// ============ Sprint 35: AI 伦理与公平性 ============
	// AI公平性测试路由
	aiFairnessCtrl := controllers.NewAIFairnessController(db)
	aiFairnessCtrl.RegisterRoutes(apiV1)

	// ============ 研究平台路由 ============
	researchCtrl := &controllers.ResearchController{DB: db}
	researchCtrl.RegisterRoutes(apiV1)

	// 通知路由
	notifCtrl := &controllers.NotificationController{DB: db}
	notifCtrl.RegisterRoutes(apiV1)

	// 补充通知路由（push 和 batch-delete）
	apiV1.POST("/notifications/push", notifCtrl.PushNotification)
	apiV1.POST("/notifications/batch-delete", notifCtrl.BatchDeleteNotifications)

	// 补充公告路由（withdraw）
	apiV1.POST("/announcements/:id/withdraw", notifCtrl.WithdrawAnnouncement)
	apiV1.GET("/announcements/:id", notifCtrl.GetAnnouncement)

	// ============ BLE Mesh 网络路由 ============
	meshCtrl := &controllers.MeshController{DB: db}
	meshCtrl.RegisterMeshRoutes(apiV1)

	// ============ 合规策略和规则 ============
	complianceCtrl := &controllers.ComplianceController{DB: db}
	complianceCtrl.RegisterRoutes(apiV1)

	// ============ 设备影子 ============
	shadowCtrl := &controllers.DeviceShadowController{DB: db}
	shadowCtrl.RegisterRoutes(apiV1)
	// apiV1.DELETE("/compliance/policies/:id", complianceCtrlExtra.DeleteRule)

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

	// Sprint 25: 安全与合规路由 (security_controller.go, gdpr_controller.go, audit_controller.go)
	// 已注释 - 相关控制器文件存在编码问题，待后续 Sprint 修复

	// ============ Sprint 27: 开发者平台 API ============
	devCtrl := &controllers.DeveloperController{DB: db}
	developerGroup := apiV1.Group("/developer")
	{
		// 应用管理
		developerGroup.POST("/apps", devCtrl.CreateApp)
		developerGroup.GET("/apps", devCtrl.ListApps)
		developerGroup.GET("/apps/:id", devCtrl.GetApp)
		developerGroup.PUT("/apps/:id", devCtrl.UpdateApp)
		developerGroup.DELETE("/apps/:id", devCtrl.DeleteApp)
		// API Key 管理
		developerGroup.POST("/apps/:id/keys", devCtrl.CreateKey)
		developerGroup.GET("/apps/:id/keys", devCtrl.ListKeys)
		developerGroup.DELETE("/apps/:id/keys/:key_id", devCtrl.DeleteKey)
		// 统计与配额
		developerGroup.GET("/stats", devCtrl.GetStats)
		developerGroup.GET("/quota", devCtrl.GetQuota)
	}

	// ============ 企业应用商店 API ============
	appStoreCtrl := &controllers.AppStoreController{DB: db, Redis: redisClient}
	storeGroup := apiV1.Group("/store")
	{
		// 应用商店 - 应用管理
		storeGroup.GET("/apps", appStoreCtrl.ListStoreApps)
		storeGroup.POST("/apps", appStoreCtrl.CreateStoreApp)
		storeGroup.GET("/apps/:id", appStoreCtrl.GetStoreApp)
		storeGroup.PUT("/apps/:id", appStoreCtrl.UpdateStoreApp)
		storeGroup.DELETE("/apps/:id", appStoreCtrl.DeleteStoreApp)
		storeGroup.POST("/apps/:id/publish", appStoreCtrl.PublishStoreApp)
		// 应用版本管理
		storeGroup.POST("/apps/:id/versions", appStoreCtrl.CreateAppVersion)
		storeGroup.GET("/apps/:id/versions", appStoreCtrl.ListAppVersions)
		storeGroup.POST("/apps/:id/versions/:version_id/set-latest", appStoreCtrl.SetLatestVersion)
		// 应用审核历史
		storeGroup.GET("/apps/:id/reviews", appStoreCtrl.GetAppReviews)
	}

	// ============ 企业应用商店 - 安装管理 API ============
	storeInstallGroup := apiV1.Group("/store/installations")
	{
		storeInstallGroup.GET("", appStoreCtrl.ListInstallations)
		storeInstallGroup.POST("", appStoreCtrl.InstallApp)
		storeInstallGroup.DELETE("/:id", appStoreCtrl.UninstallApp)
		storeInstallGroup.GET("/:id/status", appStoreCtrl.GetInstallationStatus)
		storeInstallGroup.PUT("/:id/status", appStoreCtrl.UpdateInstallationStatus)
	}

	// ============ 企业应用商店 - 审核管理 API ============
	storeReviewGroup := apiV1.Group("/store/reviews")
	{
		storeReviewGroup.GET("", appStoreCtrl.ListReviews)
		storeReviewGroup.POST("/:id/approve", appStoreCtrl.ApproveReview)
		storeReviewGroup.POST("/:id/reject", appStoreCtrl.RejectReview)
	}

	// ============ Sprint 31: 国际化扩展 API ============
	// 多语言翻译管理
	i18nCtrl := controllers.NewI18nController(db)
	i18nGroup := apiV1.Group("/i18n")
	{
		i18nGroup.GET("/translations", i18nCtrl.ListTranslations)
		i18nGroup.POST("/translations", i18nCtrl.CreateTranslation)
		i18nGroup.GET("/translations/:id", i18nCtrl.GetTranslation)
		i18nGroup.PUT("/translations/:id", i18nCtrl.UpdateTranslation)
		i18nGroup.DELETE("/translations/:id", i18nCtrl.DeleteTranslation)
		i18nGroup.GET("/locales", i18nCtrl.GetSupportedLocales)
		i18nGroup.GET("/namespaces", i18nCtrl.GetNamespaces)
	}

	// ============ Sprint 32: 高级安全功能 API ============
	// securityCtrl := &controllers.SecurityController{DB: db, Redis: redisClient}
	// // 2FA APIs
	// apiV1.POST("/security/2fa/enable", securityCtrl.Enable2FA)
	// apiV1.POST("/security/2fa/verify", securityCtrl.Verify2FA)
	// apiV1.POST("/security/2fa/disable", securityCtrl.Disable2FA)
	// // Session management APIs
	// apiV1.GET("/security/sessions", securityCtrl.GetSessions)
	// apiV1.DELETE("/security/sessions/:id", securityCtrl.DeleteSession)
	// apiV1.DELETE("/security/sessions/all", securityCtrl.DeleteAllSessions)
	// // Security audit APIs
	// apiV1.GET("/security/audit", securityCtrl.GetSecurityAudits)
	// apiV1.POST("/security/audit/report", securityCtrl.GenerateSecurityReport)

	// 区域管理 API
	_ = &controllers.RegionController{DB: db} // regionCtrl
	regionSvc := multi_region.NewRegionService(db)
	apiV1.GET("/regions", func(c *gin.Context) {
		regions, err := regionSvc.ListRegions()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 0, "data": regions})
	})
	apiV1.GET("/regions/:id", func(c *gin.Context) {
		id := c.Param("id")
		regionID, _ := strconv.ParseUint(id, 10, 32)
		region, err := regionSvc.GetRegionByID(uint(regionID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "region not found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 0, "data": region})
	})
	apiV1.PUT("/regions/:id", func(c *gin.Context) {
		id := c.Param("id")
		regionID, _ := strconv.ParseUint(id, 10, 32)
		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := regionSvc.UpdateRegion(uint(regionID), updates); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		region, _ := regionSvc.GetRegionByID(uint(regionID))
		c.JSON(200, gin.H{"code": 0, "data": region})
	})

	// 时区管理 API
	timezoneCtrl := &controllers.TimezoneController{DB: db}
	apiV1.GET("/timezones", timezoneCtrl.GetSupportedTimezones)
	apiV1.GET("/timezones/:id", func(c *gin.Context) {
		// 获取时区详情 - 根据ID返回时区配置
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		tzSvc := timezone.NewTimezoneService(db)
		config, err := tzSvc.GetTimezoneConfigByID(uint(id))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "timezone not found"})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 0, "data": config})
	})

	// ============ 会员管理路由 (9个子模块) ============
	memberCtrl := controllers.NewMemberController(db)
	memberGroup := apiV1.Group("/members")
	{
		// 会员基础 CRUD
		memberGroup.GET("/list", memberCtrl.MemberList)
		memberGroup.POST("/create", memberCtrl.MemberCreate)
		memberGroup.PUT("/update", memberCtrl.MemberUpdate)
		memberGroup.DELETE("/delete", memberCtrl.MemberDelete)
		memberGroup.GET("/detail/:id", memberCtrl.MemberDetail)
		// 会员卡
		memberGroup.GET("/cards/list", memberCtrl.CardList)
		memberGroup.POST("/cards/create", memberCtrl.CardCreate)
		memberGroup.PUT("/cards/update", memberCtrl.CardUpdate)
		memberGroup.DELETE("/cards/delete", memberCtrl.CardDelete)
		// 优惠券
		memberGroup.GET("/coupons/list", memberCtrl.CouponList)
		memberGroup.POST("/coupons/create", memberCtrl.CouponCreate)
		memberGroup.PUT("/coupons/update", memberCtrl.CouponUpdate)
		memberGroup.DELETE("/coupons/delete", memberCtrl.CouponDelete)
		// 门店
		memberGroup.GET("/stores/list", memberCtrl.StoreList)
		memberGroup.POST("/stores/create", memberCtrl.StoreCreate)
		memberGroup.PUT("/stores/update", memberCtrl.StoreUpdate)
		memberGroup.DELETE("/stores/delete", memberCtrl.StoreDelete)
		// 积分
		memberGroup.GET("/points/records", memberCtrl.PointsRecordList)
		// 标签
		memberGroup.GET("/tags/list", memberCtrl.TagList)
		memberGroup.POST("/tags/create", memberCtrl.TagCreate)
		memberGroup.PUT("/tags/update", memberCtrl.TagUpdate)
		memberGroup.DELETE("/tags/delete", memberCtrl.TagDelete)
		// 等级
		memberGroup.GET("/levels/list", memberCtrl.LevelList)
		// 操作记录
		memberGroup.GET("/operations/list", memberCtrl.OperationRecordList)
	}

	// ============ DaaS 设备即服务 API ============
	daasCtrl := &controllers.DaaSController{DB: db, Redis: redisClient}
	daasGroup := apiV1.Group("/daas")
	{
		// 租赁合同管理
		daasGroup.GET("/contracts", daasCtrl.ListContracts)
		daasGroup.POST("/contracts", daasCtrl.CreateContract)
		daasGroup.GET("/contracts/:id", daasCtrl.GetContract)
		daasGroup.PUT("/contracts/:id", daasCtrl.UpdateContract)
		daasGroup.POST("/contracts/:id/terminate", daasCtrl.TerminateContract)

		// 设备租赁管理
		daasGroup.GET("/devices", daasCtrl.ListDaasDevices)
		daasGroup.POST("/devices/:device_id/rent", daasCtrl.RentDevice)
		daasGroup.POST("/devices/:device_id/return", daasCtrl.ReturnDevice)

		// 租赁账单
		daasGroup.GET("/billing", daasCtrl.ListBillings)
		daasGroup.POST("/billing/calculate", daasCtrl.CalculateBilling)
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
