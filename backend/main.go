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

	// 自动迁移数据库表 - 已禁用
	// 数据库已完整（319表），AutoMigrate 会导致模型列名与数据库列名不匹配的错误
	// 如需迁移，请手动执行 SQL 或使用 migration 工具
	log.Printf("Database migration skipped - schema already complete")

	// 引用 models 包（被其他包间接使用）
	_ = models.Role{}

	// 仅为 system_configs 表执行自动迁移（不影响其他表）
	if err := db.AutoMigrate(&models.SystemConfig{}); err != nil {
		log.Printf("Warning: Failed to migrate SystemConfig: %v", err)
	} else {
		log.Printf("SystemConfig table ready")
	}

	// AI相关表迁移
	if err := db.AutoMigrate(&models.AIConversation{}, &models.AIMessage{}, &models.AIConfig{}); err != nil {
		log.Printf("Warning: Failed to migrate AI models: %v", err)
	} else {
		log.Printf("AI models table ready")
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
		log.Printf("Warning: Failed to initialize MQTT: %v (MQTT features will be disabled)", err)
	} else {
		defer mqttHandler.Disconnect(0)
		// 设置全局 MQTT 客户端供 CommandController 使用
		mqtt.SetGlobalMQTTClient(mqttHandler)
	}

	// OTA 后台 Worker：定期检查待下发的 OTA 任务
	otaWorker := services.NewOTAWorker(db, mqttHandler)
	controllers.SetOTAWorkerRef(otaWorker)
	go otaWorker.Start()

	// 订阅自动续费 Worker：每小时检查到期订阅并自动续费
	subscriptionRenewalWorker := services.NewSubscriptionRenewalWorker(db)
	go subscriptionRenewalWorker.Start()

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
	authCtrl := controllers.NewAuthController(db)
	r.POST("/api/v1/auth/login", authCtrl.Login)
	r.POST("/api/v1/auth/refresh", authCtrl.RefreshToken)

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

	// Sprint 27: 第三方集成路由
	integrationCtrl := &controllers.IntegrationController{DB: db}
	integrationGroup := r.Group("/api/v1")
	integrationCtrl.RegisterIntegrationRoutes(integrationGroup)

	// Sprint 28: 数据分析增强路由
	analyticsCtrl := controllers.NewAnalyticsController(db, redisClient)
	analyticsGroup := r.Group("/api/v1")
	analyticsCtrl.RegisterRoutes(analyticsGroup)

	// Sprint 23: 仿真测试路由
	simCtrl := &controllers.SimulationController{DB: db, Redis: redisClient}
	simGroup := r.Group("/api/v1")
	simCtrl.RegisterSimulationRoutes(simGroup)

	// 注册系统管理路由
	sys := r.Group("/api/v1")
	dictCtrl := &controllers.DictController{DB: db}
	logCtrl := &controllers.LogController{DB: db}

	// 租户管理路由（超管）
	tenantCtrl := &controllers.TenantController{DB: db}
	adminGroup := sys.Group("/admin")
	tenantCtrl.RegisterTenantRoutes(adminGroup)

	// Admin 套餐管理路由（/admin/packages）
	// 注意：套餐同时可通过 /admin/plans（TenantController.ListPlans）访问
	adminGroup.GET("/packages", func(c *gin.Context) {
		var packages []models.Package
		var total int64
		query := db.Model(&models.Package{})
		if err := query.Count(&total).Error; err != nil {
			c.JSON(500, gin.H{"code": "DB_ERROR", "message": "查询套餐列表失败"})
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		offset := (page - 1) * pageSize
		if err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&packages).Error; err != nil {
			c.JSON(500, gin.H{"code": "DB_ERROR", "message": "查询套餐列表失败"})
			return
		}
		// 填充兼容字段
		for i := range packages {
			packages[i].FillCompatFields()
		}
		c.JSON(200, gin.H{"code": 0, "data": gin.H{"list": packages, "total": total, "page": page, "page_size": pageSize}})
	})

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

		// 用户管理路由
		userCtrl := &controllers.UserController{DB: db}
		sys.GET("/users", userCtrl.List)
		sys.GET("/users/me", userCtrl.GetCurrentUser)
		sys.GET("/users/:id", userCtrl.Get)
		sys.POST("/users", userCtrl.Create)
		sys.PUT("/users/:id", userCtrl.Update)
		sys.DELETE("/users/:id", userCtrl.Delete)
		sys.PUT("/users/:id/status", userCtrl.UpdateStatus)
		sys.POST("/users/:id/reset-password", userCtrl.ResetPassword)
		sys.POST("/users/change-password", userCtrl.ChangePassword)
		sys.POST("/users/batch-delete", userCtrl.BatchDelete)

		sys.GET("/menus/tree", menuCtrl.List) // 复用 List 返回树形
		sys.GET("/dicts/:type", dictCtrl.GetDictByType)
		// 数据字典 CRUD（补充 /api/v1/dicts 完整路径）
		sys.GET("/dicts", dictCtrl.List)
		sys.POST("/dicts", dictCtrl.Create)
		sys.PUT("/dicts/:id", dictCtrl.Update)
		sys.DELETE("/dicts/:id", dictCtrl.Delete)
		sys.GET("/logs/operations", logCtrl.GetOperationLogs)
		sys.GET("/logs/login", logCtrl.GetLoginLogs)
		// 审计日志（基于 AuditLog 模型，区别于 SysOperationLog）
		auditCtrl := &controllers.AuditController{DB: db}
		sys.GET("/audit/logs", auditCtrl.GetAuditLogs)

		// 系统设置
		settingsCtrl := &controllers.SettingsController{DB: db}
		sys.GET("/settings", settingsCtrl.GetSettings)
		sys.GET("/settings/groups", settingsCtrl.ListSettingsByGroup)
		sys.GET("/settings/:key", settingsCtrl.GetSetting)
		sys.PUT("/settings/:key", settingsCtrl.UpdateSetting)
		sys.PUT("/settings", settingsCtrl.BatchUpdateSettings)

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

	// Sprint 17: 情感计算路由
	emotionCtrl := &controllers.EmotionController{DB: db}
	emotionCtrl.RegisterRoutes(apiV1)

	// Sprint 17 Phase 3: 情感计算 - 情绪识别/响应路由（补充 /emotion/* 别名）
	apiV1.GET("/emotion/recognition", emotionCtrl.GetRecords)
	apiV1.GET("/emotion/recognition/stats", emotionCtrl.GetRecordStats)
	apiV1.GET("/emotion/responses", emotionCtrl.GetEmotionConfig)
	apiV1.POST("/emotion/responses", emotionCtrl.UpdateEmotionConfig)

	// AI聊天路由
	aiCtrl := &controllers.AIController{DB: db}
	apiV1.POST("/ai/chat", aiCtrl.Chat)
	apiV1.GET("/ai/conversations", aiCtrl.GetConversations)
	apiV1.GET("/ai/conversations/:session_id/messages", aiCtrl.GetMessages)
	apiV1.DELETE("/ai/conversations/:session_id", aiCtrl.DeleteConversation)
	apiV1.GET("/ai/config", aiCtrl.GetAIConfig)
	apiV1.PUT("/ai/config", aiCtrl.UpdateAIConfig)

	// Sprint 26 Phase 4: App市场路由
	marketAppCtrl := &controllers.MarketAppController{DB: db}
	marketAppCtrl.RegisterRoutes(apiV1)

	// Sprint 26 Phase 4: 表情包市场路由（已由 MarketController 处理，但补充独立路由）
	// emoticon/voice/action 路由已在 marketCtrl.RegisterMarketRoutes 中注册

	// Sprint 27 Phase 4: 智能家居路由
	smartHomeCtrl := &controllers.SmartHomeController{DB: db}
	smartHomeCtrl.RegisterRoutes(apiV1)

	// Sprint 26 Phase 4: 家庭相册路由
	familyAlbumCtrl := &controllers.FamilyAlbumController{DB: db}
	familyAlbumCtrl.RegisterRoutes(apiV1)

	// Sprint 26 Phase 4: 儿童模式路由
	childModeCtrl := &controllers.ChildModeController{DB: db}
	childModeCtrl.RegisterRoutes(apiV1)

	// MiniClaw路由
	miniClawCtrl := &controllers.MiniClawController{}
	miniClawCtrl.RegisterRoutes(apiV1)

	// ============ Sprint 18: 数字孪生路由 ============
	digitalTwinCtrl := &controllers.DigitalTwinController{DB: db}
	apiV1.GET("/digital-twin/vitals/dashboard", digitalTwinCtrl.GetVitalsDashboard)
	apiV1.GET("/digital-twin/vitals/realtime", digitalTwinCtrl.GetRealtimeVitals)
	apiV1.GET("/digital-twin/vitals/history", digitalTwinCtrl.GetVitalsHistory)
	apiV1.GET("/digital-twin/vitals/alerts", digitalTwinCtrl.GetHealthAlerts)
	apiV1.POST("/digital-twin/vitals/alerts/:id/confirm", digitalTwinCtrl.ConfirmAlert)
	apiV1.POST("/digital-twin/vitals/alerts/:id/ignore", digitalTwinCtrl.IgnoreAlert)
	apiV1.GET("/digital-twin/behavior/prediction", digitalTwinCtrl.GetBehaviorPrediction)
	apiV1.GET("/digital-twin/behavior/history", digitalTwinCtrl.GetBehaviorHistory)
	apiV1.GET("/digital-twin/replay/:pet_id", digitalTwinCtrl.GetReplay)
	apiV1.GET("/digital-twin/highlights", digitalTwinCtrl.GetHighlights)

	// 通知路由
	notifCtrl := &controllers.NotificationController{DB: db}
	notifCtrl.RegisterRoutes(apiV1)

	// 补充通知路由（push 和 batch-delete）
	apiV1.POST("/notifications/push", notifCtrl.PushNotification)
	apiV1.POST("/notifications/batch-delete", notifCtrl.BatchDeleteNotifications)

	// 补充公告路由（withdraw）
	apiV1.POST("/announcements/:id/withdraw", notifCtrl.WithdrawAnnouncement)
	apiV1.GET("/announcements/:id", notifCtrl.GetAnnouncement)

	// ============ Sprint 9-10: 策略配置路由 ============
	policyCtrlExtra := &controllers.PolicyController{DB: db}
	apiV1.GET("/policy-configs", policyCtrlExtra.ListConfigs)
	apiV1.POST("/policy-configs", policyCtrlExtra.CreateConfig)
	apiV1.GET("/policy-configs/:id", policyCtrlExtra.GetConfig)
	apiV1.PUT("/policy-configs/:id", policyCtrlExtra.UpdateConfig)
	apiV1.DELETE("/policy-configs/:id", policyCtrlExtra.DeleteConfig)

	// 策略主表路由
	apiV1.GET("/policies", policyCtrlExtra.ListPolicies)
	apiV1.POST("/policies", policyCtrlExtra.CreatePolicy)
	apiV1.GET("/policies/:id", policyCtrlExtra.GetPolicy)
	apiV1.PUT("/policies/:id", policyCtrlExtra.UpdatePolicy)
	apiV1.DELETE("/policies/:id", policyCtrlExtra.DeletePolicy)
	apiV1.POST("/policies/:id/bind", policyCtrlExtra.BindPolicy)
	apiV1.POST("/policies/:id/unbind", policyCtrlExtra.UnbindPolicy)
	apiV1.GET("/policies/:id/bindings", policyCtrlExtra.GetPolicyBindings)

	// ============ Sprint 9-10: 合规规则路由 ============
	complianceCtrlExtra := &controllers.ComplianceController{DB: db}
	apiV1.GET("/compliance-rules", complianceCtrlExtra.ListRules)
	apiV1.POST("/compliance-rules", complianceCtrlExtra.CreateRule)
	apiV1.GET("/compliance-rules/:id", complianceCtrlExtra.GetRule)
	apiV1.PUT("/compliance-rules/:id", complianceCtrlExtra.UpdateRule)
	apiV1.DELETE("/compliance-rules/:id", complianceCtrlExtra.DeleteRule)
	apiV1.POST("/compliance-rules/:id/enforce", complianceCtrlExtra.EnforceRule)

	// 合规策略 API（/api/v1/compliance/policies）
	apiV1.GET("/compliance/policies", complianceCtrlExtra.ListCompliancePolicies)
	apiV1.POST("/compliance/policies", complianceCtrlExtra.CreateCompliancePolicy)
	apiV1.PUT("/compliance/policies/:id", complianceCtrlExtra.UpdateCompliancePolicy)
	apiV1.DELETE("/compliance/policies/:id", complianceCtrlExtra.DeleteCompliancePolicy)
	apiV1.GET("/compliance/policies/:id/violations", complianceCtrlExtra.GetPolicyViolations)

	// 合规违规记录
	apiV1.GET("/compliance/violations", complianceCtrlExtra.ListViolations)
	apiV1.PUT("/compliance/violations/:id/resolve", complianceCtrlExtra.ResolveViolation)

	// ============ Sprint 9-10: GDPR 路由 ============
	apiV1.POST("/gdpr/export", complianceCtrlExtra.GetGDPRDataExport)
	apiV1.POST("/gdpr/delete", complianceCtrlExtra.DeleteGDPRData)
	apiV1.GET("/gdpr/requests", complianceCtrlExtra.GetGDPRRequests)
	apiV1.GET("/gdpr/requests/:id", complianceCtrlExtra.GetGDPRRequest)
	apiV1.PUT("/gdpr/requests/:id/process", complianceCtrlExtra.ProcessGDPRRequest)

	// ============ 门店管理路由 ============
	storesCtrl := &controllers.StoresController{DB: db}
	apiV1.GET("/stores", storesCtrl.StoreList)
	apiV1.POST("/stores", storesCtrl.StoreCreate)
	apiV1.GET("/stores/:id", storesCtrl.StoreGet)
	apiV1.PUT("/stores/:id", storesCtrl.StoreUpdate)
	apiV1.DELETE("/stores/:id", storesCtrl.StoreDelete)
	apiV1.PUT("/stores/:id/status", storesCtrl.StoreUpdateStatus)
	apiV1.POST("/stores/batch-delete", storesCtrl.StoreBatchDelete)
	apiV1.GET("/stores/statistics", storesCtrl.StoreStatistics)

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

	// ============ Sprint 21: 具身智能路由 ============
	embodiedCtrl := &controllers.EmbodiedController{DB: db}
	embodiedCtrl.RegisterEmbodiedRoutes(apiV1)

	// ============ Sprint 30-31: 平台演进路由 ============
	platformEvoCtrl := &controllers.PlatformEvoController{DB: db}
	platformEvoCtrl.RegisterRoutes(apiV1)

	// Sprint 31: 数据集开放平台 + AI行为研究平台路由
	researchCtrl := &controllers.ResearchPlatformController{DB: db}
	researchCtrl.RegisterRoutes(apiV1)

	// Sprint 32: 高级安全功能路由
	securityEvoCtrl := &controllers.SecurityEvoController{DB: db}
	securityEvoCtrl.RegisterSecurityEvoRoutes(apiV1)

	// ============ Sprint 12: 数据权限路由 ============
	dataPermCtrl := &controllers.DataPermissionController{DB: db}
	apiV1.GET("/data-permissions/rules", dataPermCtrl.ListDataPermissionRules)

	// ============ 批量操作路由 ============
	controllers.RegisterBatchRoutes(apiV1, db, redisClient)

	// ============ 设备监控路由 ============
	controllers.RegisterDeviceMonitorRoutes(apiV1, db, redisClient)
	apiV1.POST("/data-permissions/rules", dataPermCtrl.CreateDataPermissionRule)
	apiV1.PUT("/data-permissions/rules/:id", dataPermCtrl.UpdateDataPermissionRule)
	apiV1.DELETE("/data-permissions/rules/:id", dataPermCtrl.DeleteDataPermissionRule)
	apiV1.GET("/data-permissions/roles/:role_id", dataPermCtrl.GetRoleDataPermissions)
	apiV1.PUT("/data-permissions/roles/:role_id", dataPermCtrl.UpdateRoleDataPermissions)
	apiV1.GET("/data-permissions/users/:user_id", dataPermCtrl.GetUserDataPermissions)
	apiV1.PUT("/data-permissions/users/:user_id", dataPermCtrl.UpdateUserDataPermissions)
	apiV1.GET("/data-permissions/columns", dataPermCtrl.GetColumnPermissions)
	apiV1.POST("/data-permissions/validate", dataPermCtrl.ValidatePermissionExpression)

	// ============ Sprint 25: 开放平台生态路由 ============
	// 初始化 Webhook 服务
	webhookSvc := services.NewWebhookService(db)
	// 初始化默认 Webhook 模板
	go webhookSvc.InitDefaultTemplates()

	// 开发者 API
	devAppCtrl := &controllers.DeveloperAppController{DB: db}
	apiKeyCtrl := &controllers.APIKeyController{DB: db}
	webhookCtrl := &controllers.WebhookController{DB: db, WebhookSvc: webhookSvc}

	// 开发者应用管理
	apiV1.GET("/developer/apps", devAppCtrl.ListApps)
	apiV1.POST("/developer/apps", devAppCtrl.CreateApp)
	apiV1.GET("/developer/apps/:id", devAppCtrl.GetApp)
	apiV1.PUT("/developer/apps/:id", devAppCtrl.UpdateApp)
	apiV1.DELETE("/developer/apps/:id", devAppCtrl.DeleteApp)
	apiV1.POST("/developer/apps/:id/regenerate-key", devAppCtrl.RegenerateKey)

	// API Key 管理
	apiV1.GET("/developer/api-keys", apiKeyCtrl.ListAPIKeys)
	apiV1.POST("/developer/api-keys", apiKeyCtrl.CreateAPIKey)
	apiV1.DELETE("/developer/api-keys/:id", apiKeyCtrl.DeleteAPIKey)

	// Webhook 市场
	apiV1.GET("/webhooks/templates", webhookCtrl.ListTemplates)
	apiV1.GET("/webhooks/templates/:id", webhookCtrl.GetTemplate)
	apiV1.POST("/webhooks/subscriptions", webhookCtrl.CreateSubscription)
	apiV1.GET("/webhooks/subscriptions", webhookCtrl.ListSubscriptions)
	apiV1.DELETE("/webhooks/subscriptions/:id", webhookCtrl.DeleteSubscription)
	apiV1.GET("/webhooks/deliveries/:id", webhookCtrl.GetDelivery)
	apiV1.POST("/webhooks/deliveries/:id/retry", webhookCtrl.RetryDelivery)

	// ============ Sprint 26: 内容市场路由 ============
	marketCtrl := &controllers.MarketController{DB: db}
	marketCtrl.RegisterMarketRoutes(apiV1)

	// ============ Sprint 29: 高级功能路由 ============
	advancedCtrl := &controllers.AdvancedController{DB: db}
	advancedCtrl.RegisterAdvancedRoutes(apiV1)

	// ============ Sprint 19: 健康医疗路由 ============
	healthCtrl := controllers.NewHealthController(db, redisClient)
	healthCtrl.RegisterRoutes(apiV1)

	// ============ 系统健康检查 API ============
	// 注意：此 API 不同于 Sprint 19 的宠物健康路由（/health/*）
	// 返回 DB/Redis/MQTT 连接状态
	systemHealthCtrl := controllers.NewSystemHealthController(db, redisClient)
	apiV1.GET("/system/health", systemHealthCtrl.GetHealth)

	// ============ 部门管理直接路径 ============
	// 注意：部门管理已在 RegisterRoutes 中注册于 /api/v1/org/departments
	// 此处补充直接路径 /api/v1/departments（与 PRD 文档一致）
	orgCtrlForDept := &controllers.OrgController{DB: db}
	apiV1.GET("/departments", orgCtrlForDept.DepartmentList)
	apiV1.GET("/departments/tree", orgCtrlForDept.DepartmentTree)
	apiV1.POST("/departments", orgCtrlForDept.DepartmentCreate)
	apiV1.GET("/departments/:id", orgCtrlForDept.DepartmentList)
	apiV1.PUT("/departments/:id", orgCtrlForDept.DepartmentUpdate)
	apiV1.DELETE("/departments/:id", orgCtrlForDept.DepartmentDelete)

	// Sprint 29: 宠物社交路由
	petSocialCtrl := controllers.NewPetSocialController(db)
	apiV1.GET("/posts", petSocialCtrl.GetPosts)
	apiV1.POST("/posts", petSocialCtrl.CreatePost)
	apiV1.GET("/posts/:id", petSocialCtrl.GetPost)
	apiV1.DELETE("/posts/:id", petSocialCtrl.DeletePost)
	apiV1.POST("/posts/:id/like", petSocialCtrl.LikePost)
	apiV1.GET("/posts/:id/comments", petSocialCtrl.GetComments)
	apiV1.POST("/posts/:id/comments", petSocialCtrl.CreateComment)
	apiV1.GET("/following", petSocialCtrl.GetFollowing)
	apiV1.POST("/follow", petSocialCtrl.Follow)
	apiV1.DELETE("/follow/:id", petSocialCtrl.Unfollow)
	apiV1.GET("/playdates", petSocialCtrl.GetPlaydates)
	apiV1.POST("/playdates", petSocialCtrl.CreatePlaydate)

	// Sprint 27: 宠物保险和医疗路由 (已由 IntegrationController 提供)
	// Sprint 15: 内容分发和应用管理路由
	contentCtrl := &controllers.ContentController{DB: db}
	contentCtrl.RegisterRoutes(apiV1)

	// Sprint 15: 发票和账单路由
	billingCtrl := &controllers.BillingController{DB: db}
	billingCtrl.RegisterRoutes(apiV1)

	// Sprint 19: 离线支持路由
	offlineCtrl := &controllers.OfflineController{DB: db}
	offlineCtrl.RegisterRoutes(apiV1)

	// Sprint 17: 语音情绪识别路由
	voiceEmotionCtrl := controllers.NewVoiceEmotionController(db)
	voiceEmotionCtrl.RegisterRoutes(apiV1)

	// Sprint 15: API配额路由
	apiQuotaCtrl := controllers.NewAPIQuotaController()
	apiQuotaCtrl.RegisterRoutes(apiV1)

	// Sprint 13: 模型分片路由
	modelShardCtrl := controllers.NewModelShardController(db)
	modelShardCtrl.RegisterRoutes(apiV1)

	// Sprint 31: 数据集开放平台路由 (research platform controller already registered above)

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
