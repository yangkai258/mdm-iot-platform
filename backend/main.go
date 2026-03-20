package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"mdm-backend/controllers"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/mqtt"
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

	// 初始化 MQTT（传入告警回调）
	alertCallback := func(deviceID string, data map[string]interface{}) {
		controllers.CheckAlerts(db, deviceID, data)
	}
	mqttHandler, err := mqtt.InitMQTT(redisClient, alertCallback)
	if err != nil {
		log.Fatalf("Failed to initialize MQTT: %v", err)
	}
	defer mqttHandler.Disconnect(0)
	// 设置全局 MQTT 客户端供 CommandController 使用
	mqtt.SetGlobalMQTTClient(mqttHandler)

	// OTA 后台 Worker：定期检查待下发的 OTA 任务
	go startOTAWorker(db)

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

	fmt.Println("MDM Backend started successfully")
}

// startOTAWorker 启动 OTA 后台检查 goroutine
func startOTAWorker(db *gorm.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		checkOTADeployments(db)
	}
}

// checkOTADeployments 检查并下发 OTA 任务
func checkOTADeployments(db *gorm.DB) {
	// 查询待下发的部署任务
	var deployments []models.OTADeployment
	if err := db.Model(&models.OTADeployment{}).Where("status = ?", "pending").Find(&deployments).Error; err != nil {
		log.Printf("[OTA-Worker] 查询部署任务失败: %v", err)
		return
	}

	for _, dep := range deployments {
		log.Printf("[OTA-Worker] 处理部署任务 #%d", dep.ID)

		// 更新状态为 rolling
		db.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Update("status", "rolling")

		// 查询目标固件包
		var pkg models.OTAPackage
		if err := db.First(&pkg, dep.PackageID).Error; err != nil {
			log.Printf("[OTA-Worker] 固件包 #%d 不存在: %v", dep.PackageID, err)
			continue
		}

		// 查询目标设备（根据 rollout strategy）
		var devices []models.Device
		query := db.Model(&models.Device{}).Where("hardware_model = ? AND lifecycle_status = ?", dep.TargetHardware, 2)
		if dep.RolloutStrategy == "whitelist" {
			// TODO: 支持白名单模式
		} else if dep.RolloutStrategy == "percentage" {
			// 按百分比随机选取
			query = query.Where("id % 100 < ?", dep.Percentage)
		}
		if err := query.Find(&devices).Error; err != nil {
			log.Printf("[OTA-Worker] 查询设备失败: %v", err)
			continue
		}

		for _, device := range devices {
			// 跳过已是最新版本的设备
			if device.FirmwareVersion == pkg.VersionCode {
				continue
			}

			// 构建 OTA 指令
			otaCmd := map[string]interface{}{
				"cmd_id":   fmt.Sprintf("ota-%d-%s", dep.ID, device.DeviceID),
				"cmd_type": "ota",
				"ota": map[string]interface{}{
					"version": pkg.VersionCode,
					"url":     pkg.BinURL,
					"md5":     pkg.Md5Hash,
				},
				"timestamp": time.Now().Format(time.RFC3339),
			}

			// 通过 MQTT 下发
			if mqtt.GlobalMQTTClient != nil {
				topic := fmt.Sprintf("/mdm/device/%s/down/cmd", device.DeviceID)
				payload, _ := json.Marshal(otaCmd)
				token := mqtt.GlobalMQTTClient.Publish(topic, 0, false, payload)
				token.Wait()
				if token.Error() != nil {
					log.Printf("[OTA-Worker] 设备 %s 下发失败: %v", device.DeviceID, token.Error())
				} else {
					log.Printf("[OTA-Worker] 已向设备 %s 下发 OTA 指令: %s", device.DeviceID, pkg.VersionCode)
				}
			}

			// 记录升级进度
			progress := models.OTAProgress{
				DeviceID:      device.DeviceID,
				DeploymentID:  dep.ID,
				PackageID:     dep.PackageID,
				TargetVersion: pkg.VersionCode,
				Status:        "pending",
			}
			db.Create(&progress)
		}

		log.Printf("[OTA-Worker] 部署任务 #%d 已处理，共 %d 台设备", dep.ID, len(devices))
	}
}
