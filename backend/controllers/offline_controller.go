package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OfflineController struct {
	DB *gorm.DB
}

func NewOfflineController(db *gorm.DB) *OfflineController {
	return &OfflineController{DB: db}
}

func (ctrl *OfflineController) RegisterRoutes(rg *gin.RouterGroup) {
	offline := rg.Group("/offline")
	{
		offline.GET("/cache/:device_id", ctrl.GetDeviceCache)
		offline.POST("/cache", ctrl.CreateCache)
		offline.POST("/cache/sync", ctrl.SyncCache)
		offline.POST("/cache/:id/confirm", ctrl.ConfirmCacheSync)
		offline.GET("/queue/:device_id", ctrl.GetOfflineQueue)
		offline.POST("/queue", ctrl.EnqueueOfflineAction)
		offline.POST("/queue/sync", ctrl.SyncOfflineQueue)
	}
}

// GetDeviceCache 获取设备缓存
func (ctrl *OfflineController) GetDeviceCache(c *gin.Context) {
	deviceID := c.Param("device_id")
	var caches []models.OfflineCache
	ctrl.DB.Where("device_id = ?", deviceID).Order("created_at DESC").Find(&caches)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": caches})
}

// CreateCache 创建缓存记录
func (ctrl *OfflineController) CreateCache(c *gin.Context) {
	var req struct {
		DeviceID  uint   `json:"device_id" binding:"required"`
		UserID   uint   `json:"user_id" binding:"required"`
		CacheKey string `json:"cache_key" binding:"required"`
		CacheData string `json:"cache_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	cache := models.OfflineCache{
		DeviceID:   req.DeviceID,
		UserID:    req.UserID,
		CacheKey:  req.CacheKey,
		CacheData: req.CacheData,
		SyncStatus: "pending",
	}
	ctrl.DB.Create(&cache)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cache})
}

// SyncCache 同步缓存
func (ctrl *OfflineController) SyncCache(c *gin.Context) {
	deviceID := c.Query("device_id")
	lastSync := c.Query("last_sync")

	var caches []models.OfflineCache
	query := ctrl.DB.Model(&models.OfflineCache{}).Where("device_id = ? AND sync_status = 'pending'", deviceID)
	if lastSync != "" {
		t, _ := time.Parse(time.RFC3339, lastSync)
		query = query.Where("created_at > ?", t)
	}
	query.Find(&caches)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"pending_items": len(caches),
		"items":         caches,
		"sync_time":     time.Now(),
	}})
}

// ConfirmCacheSync 确认缓存已同步
func (ctrl *OfflineController) ConfirmCacheSync(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	ctrl.DB.Model(&models.OfflineCache{}).Where("id = ?", id).Updates(map[string]interface{}{
		"sync_status": "synced",
		"synced_at":   &now,
	})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已确认"})
}

// GetOfflineQueue 获取离线队列
func (ctrl *OfflineController) GetOfflineQueue(c *gin.Context) {
	deviceID := c.Param("device_id")
	var queues []models.OfflineQueue
	ctrl.DB.Where("device_id = ?", deviceID).Order("created_at ASC").Find(&queues)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": queues})
}

// EnqueueOfflineAction 添加离线操作
func (ctrl *OfflineController) EnqueueOfflineAction(c *gin.Context) {
	var req struct {
		DeviceID   uint   `json:"device_id" binding:"required"`
		UserID   uint   `json:"user_id" binding:"required"`
		ActionType string `json:"action_type" binding:"required"`
		ActionData string `json:"action_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	queue := models.OfflineQueue{
		DeviceID:   req.DeviceID,
		UserID:   req.UserID,
		ActionType: req.ActionType,
		ActionData: req.ActionData,
		Status:     "pending",
	}
	ctrl.DB.Create(&queue)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": queue})
}

// SyncOfflineQueue 批量同步离线操作
func (ctrl *OfflineController) SyncOfflineQueue(c *gin.Context) {
	deviceID := c.Query("device_id")

	var queues []models.OfflineQueue
	ctrl.DB.Where("device_id = ? AND status = 'pending'", deviceID).
		Order("created_at ASC").Find(&queues)

	results := make([]map[string]interface{}, 0)
	for _, q := range queues {
		now := time.Now()
		ctrl.DB.Model(&q).Updates(map[string]interface{}{
			"status":  "sent",
			"sent_at": &now,
		})
		results = append(results, map[string]interface{}{
			"id":     q.ID,
			"status": "sent",
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": results})
}
