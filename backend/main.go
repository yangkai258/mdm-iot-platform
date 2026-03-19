package main

import (
	"fmt"
	"log"
	"os"

	"mdm-backend/controllers"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
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
		// 会员表
		&models.MemberOrder{},
		&models.MemberUpgradeRecord{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化 Redis
	redisClient, err := utils.InitRedis()
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 初始化 MQTT
	mqttHandler, err := mqtt.InitMQTT(redisClient)
	if err != nil {
		log.Fatalf("Failed to initialize MQTT: %v", err)
	}
	defer mqttHandler.Disconnect(0)

	// 初始化 Gin 路由
	r := gin.Default()
	
	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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

	fmt.Println("MDM Backend started successfully")
}
