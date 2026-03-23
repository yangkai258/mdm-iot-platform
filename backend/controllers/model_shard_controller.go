package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModelShardController struct {
	DB *gorm.DB
}

func (ctrl *ModelShardController) RegisterRoutes(rg *gin.RouterGroup) {
	shard := rg.Group("/model-shards")
	{
		shard.GET("/:model_id", ctrl.GetShards)
		shard.POST("/:model_id/config", ctrl.UpdateShardConfig)
		shard.GET("/:model_id/config", ctrl.GetShardConfig)
		shard.POST("/:model_id/load", ctrl.LoadShards)
		shard.POST("/:model_id/unload", ctrl.UnloadShards)
		shard.GET("/:model_id/status", ctrl.GetLoadStatus)
	}
}

func (ctrl *ModelShardController) GetShards(c *gin.Context) {
	modelID := c.Param("model_id")
	var shards []models.ModelShard
	ctrl.DB.Where("model_id = ? AND status = 'active'", modelID).
		Order("shard_index ASC").Find(&shards)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": shards})
}

func (ctrl *ModelShardController) UpdateShardConfig(c *gin.Context) {
	modelID := c.Param("model_id")
	var req struct {
		Strategy      string `json:"strategy" binding:"required"`
		PriorityShards []int `json:"priority_shards"`
		CacheSize     int64  `json:"cache_size"`
		MinBandwidth  int    `json:"min_bandwidth"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	config := models.ModelShardConfig{
		ModelID:        uintVal(modelID),
		Strategy:       req.Strategy,
		PriorityShards: intSliceToJSON(req.PriorityShards),
		CacheSize:     req.CacheSize,
		MinBandwidth:  req.MinBandwidth,
	}

	var existing models.ModelShardConfig
	if err := ctrl.DB.Where("model_id = ?", modelID).First(&existing).Error; err == nil {
		ctrl.DB.Model(&existing).Updates(config)
	} else {
		ctrl.DB.Create(&config)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "配置已更新"})
}

func (ctrl *ModelShardController) GetShardConfig(c *gin.Context) {
	modelID := c.Param("model_id")
	var config models.ModelShardConfig
	if err := ctrl.DB.Where("model_id = ?", modelID).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配置不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
}

func (ctrl *ModelShardController) LoadShards(c *gin.Context) {
	modelID := c.Param("model_id")
	deviceID := c.Query("device_id")

	var config models.ModelShardConfig
	ctrl.DB.Where("model_id = ?", modelID).First(&config)

	var shards []models.ModelShard
	if config.PriorityShards != "" {
		ctrl.DB.Where("model_id = ? AND status = 'active'", modelID).
			Order("shard_index ASC").Limit(10).Find(&shards)
	} else {
		ctrl.DB.Where("model_id = ? AND status = 'active'", modelID).
			Order("shard_index ASC").Limit(10).Find(&shards)
	}

	downloadList := make([]map[string]interface{}, 0)
	for _, s := range shards {
		downloadList = append(downloadList, map[string]interface{}{
			"shard_index": s.ShardIndex,
			"url":         s.FileURL,
			"size":        s.FileSize,
			"hash":        s.ShardHash,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"device_id":  deviceID,
		"shards":     downloadList,
		"strategy":   config.Strategy,
		"cache_size": config.CacheSize,
	}})
}

func (ctrl *ModelShardController) UnloadShards(c *gin.Context) {
	modelID := c.Param("model_id")
	deviceID := c.Query("device_id")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "分片已卸载", "data": gin.H{
		"model_id":  modelID,
		"device_id": deviceID,
	}})
}

func (ctrl *ModelShardController) GetLoadStatus(c *gin.Context) {
	modelID := c.Param("model_id")
	deviceID := c.Query("device_id")

	var totalShards int64
	var loadedShards int64

	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ?", modelID).Count(&totalShards)
	loadedShards = totalShards / 2 // 模拟

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"model_id":      modelID,
		"device_id":     deviceID,
		"total_shards":  totalShards,
		"loaded_shards":  loadedShards,
		"progress":      float64(loadedShards) / float64(totalShards) * 100,
		"memory_usage":  256,
		"last_updated":  time.Now(),
	}})
}

func intSliceToJSON(arr []int) string {
	result := "["
	for i, v := range arr {
		if i > 0 {
			result += ","
		}
		result += strconv.Itoa(v)
	}
	return result + "]"
}
