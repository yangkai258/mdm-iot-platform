package controllers

import (
	"context"
	"net/http"
	"time"

	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SystemHealthController 系统健康检查控制器
type SystemHealthController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// NewSystemHealthController 创建系统健康检查控制器
func NewSystemHealthController(db *gorm.DB, redis *utils.RedisClient) *SystemHealthController {
	return &SystemHealthController{DB: db, Redis: redis}
}

// GetHealth 获取系统健康状态
// GET /api/v1/system/health
func (c *SystemHealthController) GetHealth(ctx *gin.Context) {
	health := gin.H{
		"status":   "healthy",
		"services": gin.H{},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	allHealthy := true

	// DB 健康检查
	dbStatus := c.checkDB()
	health["services"].(gin.H)["db"] = dbStatus
	if dbStatus["status"] != "up" {
		allHealthy = false
	}

	// Redis 健康检查
	redisStatus := c.checkRedis()
	health["services"].(gin.H)["redis"] = redisStatus
	if redisStatus["status"] != "up" {
		allHealthy = false
	}

	// MQTT 健康检查
	mqttStatus := c.checkMQTT()
	health["services"].(gin.H)["mqtt"] = mqttStatus
	if mqttStatus["status"] != "up" {
		allHealthy = false
	}

	if !allHealthy {
		health["status"] = "degraded"
		ctx.JSON(http.StatusServiceUnavailable, health)
		return
	}

	ctx.JSON(http.StatusOK, health)
}

// checkDB 检查数据库连接
func (c *SystemHealthController) checkDB() gin.H {
	result := gin.H{"status": "up", "latency_ms": 0}
	sqlDB, err := c.DB.DB()
	if err != nil {
		result["status"] = "down"
		result["error"] = "failed to get underlying db"
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	if err := sqlDB.PingContext(ctx); err != nil {
		result["status"] = "down"
		result["error"] = err.Error()
		return result
	}
	result["latency_ms"] = time.Since(start).Milliseconds()
	return result
}

// checkRedis 检查 Redis 连接
func (c *SystemHealthController) checkRedis() gin.H {
	result := gin.H{"status": "up", "latency_ms": 0}
	if c.Redis == nil {
		result["status"] = "down"
		result["error"] = "redis client not initialized"
		return result
	}
	redisClient := c.Redis.Client()
	if redisClient == nil {
		result["status"] = "down"
		result["error"] = "redis client is nil"
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		result["status"] = "down"
		result["error"] = err.Error()
		return result
	}
	result["latency_ms"] = time.Since(start).Milliseconds()
	return result
}

// checkMQTT 检查 MQTT 连接
func (c *SystemHealthController) checkMQTT() gin.H {
	result := gin.H{"status": "up"}
	if mqtt.GlobalMQTTClient == nil {
		result["status"] = "down"
		result["error"] = "mqtt client not initialized"
		return result
	}
	if !mqtt.GlobalMQTTClient.IsConnected() {
		result["status"] = "down"
		result["error"] = "mqtt client not connected"
		return result
	}
	return result
}
