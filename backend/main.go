package main

import (
	"log"
	"os"
	"strings"

	"mdm-backend/controllers"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/mqtt"
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
		&models.PetProfile{},
		&models.OTAPackage{},
		&models.OTADeployment{},
		&models.OTAProgress{},
		&models.CommandHistory{},
		// 系统表
		&models.SysUser{},
		&models.SysRole{},
		&models.SysMenu{},
		&models.SysDictionary{},
		&models.SysOperationLog{},
		&models.SysLoginLog{},
		// 告警表
		&models.DeviceAlertRule{},
		&models.DeviceAlert{},
		// 合规表
		&models.CompliancePolicy{},
		&models.ComplianceViolation{},
		// 会员表
		&models.MemberOrder{},
		&models.MemberUpgradeRecord{},
		// 应用管理表
		&models.App{},
		&models.AppVersion{},
		&models.AppDistribution{},
		&models.AppInstallRecord{},
		&models.AppLicense{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
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
	mqttHandler, err := mqtt.InitMQTT(db, redisClient, alertCallback, complianceCallback)
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

	// 操作日志中间件
	r.Use(middleware.OperationLog(db))

	// 注册业务路由
	controllers.RegisterRoutes(r, db, redisClient)

	// 注册系统管理路由
	menuCtrl := &controllers.MenuController{DB: db}
	dictCtrl := &controllers.DictController{DB: db}
	logCtrl := &controllers.LogController{DB: db}

	sys := r.Group("/api/v1")
	{
		sys.GET("/menus/tree", menuCtrl.GetMenuTree)
		sys.GET("/dicts/:type", dictCtrl.GetDictByType)
		sys.GET("/logs/operations", logCtrl.GetOperationLogs)
		sys.GET("/logs/login", logCtrl.GetLoginLogs)

		// 告警管理
		alertCtrl := &controllers.AlertController{DB: db}
		sys.GET("/alerts/rules", alertCtrl.GetRules)
		sys.POST("/alerts/rules", alertCtrl.CreateRule)
		sys.GET("/alerts", alertCtrl.GetAlerts)
		sys.GET("/dashboard/stats", alertCtrl.GetDashboardStats)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

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
