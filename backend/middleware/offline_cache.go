package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/services"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OfflineCacheMiddleware 离线缓存中间件
// 检测设备是否在线，离线时将操作存入队列
type OfflineCacheMiddleware struct {
	DB              *gorm.DB
	SyncService     *services.OfflineSyncService
	DeviceOnlineTTL time.Duration // 设备在线判定TTL，默认90秒
}

// NewOfflineCacheMiddleware 创建离线缓存中间件
func NewOfflineCacheMiddleware(db *gorm.DB) *OfflineCacheMiddleware {
	return &OfflineCacheMiddleware{
		DB:              db,
		SyncService:     services.NewOfflineSyncService(db),
		DeviceOnlineTTL: 90 * time.Second,
	}
}

// IsDeviceOnline 检测设备是否在线
// 通过 Redis 中的设备影子心跳判定
func (m *OfflineCacheMiddleware) IsDeviceOnline(deviceID string) bool {
	// 尝试从全局 Redis 客户端获取设备影子
	redisClient := getGlobalRedisClient()
	if redisClient != nil {
		shadow, err := redisClient.GetDeviceShadow(deviceID)
		if err == nil && shadow != nil && shadow.LastHeartbeat != nil {
			return time.Since(*shadow.LastHeartbeat) < m.DeviceOnlineTTL
		}
	}

	// 降级：从数据库查询
	var shadow utils.DeviceShadow
	if err := m.DB.Where("device_id = ?", deviceID).First(&shadow).Error; err != nil {
		return false
	}
	if shadow.LastHeartbeat == nil {
		return false
	}
	return time.Since(*shadow.LastHeartbeat) < m.DeviceOnlineTTL
}

// CacheResponse 缓存响应数据
func (m *OfflineCacheMiddleware) CacheResponse(deviceID, dataType string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = m.SyncService.SaveCache(deviceID, dataType, string(jsonData))
	return err
}

// EnqueueOfflineOperation 将操作加入离线队列
func (m *OfflineCacheMiddleware) EnqueueOfflineOperation(deviceID, operation, payload string) error {
	op := models.OfflineOperation{
		DeviceID:  deviceID,
		Operation: operation,
		Payload:   payload,
		CreatedAt: time.Now(),
		Status:    "pending",
	}
	return m.DB.Create(&op).Error
}

// OfflineCacheHandler 离线缓存处理器中间件
// 使用方式：作为特定路由的中间件使用
func (m *OfflineCacheMiddleware) OfflineCacheHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文或参数获取设备ID
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = c.Query("device_id")
		}

		if deviceID != "" && !m.IsDeviceOnline(deviceID) {
			// 设备离线，处理离线逻辑
			m.handleOfflineRequest(c, deviceID)
			return
		}

		// 设备在线，继续处理请求
		c.Next()
	}
}

// handleOfflineRequest 处理离线请求
func (m *OfflineCacheMiddleware) handleOfflineRequest(c *gin.Context, deviceID string) {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 读取请求体
	var payload string
	if c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		payload = string(bodyBytes)
	}

	// 根据请求方法和路径决定操作类型
	operation := m.classifyOperation(method, path)

	// 将操作加入离线队列
	if err := m.EnqueueOfflineOperation(deviceID, operation, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "离线队列保存失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回离线响应
	c.JSON(http.StatusAccepted, gin.H{
		"code":    202,
		"message": "设备离线，操作已加入离线队列",
		"data": gin.H{
			"device_id": deviceID,
			"operation": operation,
			"status":    "queued",
		},
	})
}

// classifyOperation 根据HTTP方法和路径分类操作
func (m *OfflineCacheMiddleware) classifyOperation(method, path string) string {
	// 根据路径判断数据类型
	switch {
	case contains(path, "/control"):
		return "control"
	case contains(path, "/setting"):
		return "setting"
	case contains(path, "/status"):
		return "device_status"
	case contains(path, "/sensor"):
		return "sensor_data"
	default:
		return "update"
	}
}

// contains 判断s是否包含substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// getGlobalRedisClient 获取全局 Redis 客户端
func getGlobalRedisClient() *utils.RedisClient {
	return utils.GetGlobalRedisClient()
}

// DeviceStatusCache 设备状态缓存中间件
// 用于缓存设备状态数据，减少数据库查询
func DeviceStatusCache(db *gorm.DB, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			c.Next()
			return
		}

		// 尝试从缓存获取
		cacheKey := "http:cache:device:" + deviceID

		// 在 c.Set 中存储缓存检查函数
		c.Set("cache_key", cacheKey)
		c.Next()
	}
}

// ResponseCache 响应缓存中间件
// 缓存 GET 请求的响应，减少服务器负载
func ResponseCache(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// 设置缓存头
		c.Header("Cache-Control", "public, max-age="+string(rune(int(maxAge.Seconds()))))
		c.Next()
	}
}
