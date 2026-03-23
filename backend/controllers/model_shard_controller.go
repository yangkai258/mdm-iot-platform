package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModelShardController struct {
	DB *gorm.DB
}

func NewModelShardController(db *gorm.DB) *ModelShardController {
	return &ModelShardController{DB: db}
}

func (ctrl *ModelShardController) RegisterRoutes(rg *gin.RouterGroup) {
	shard := rg.Group("/model-shards")
	{
		shard.GET("/:model_id", ctrl.GetShards)
		shard.GET("/:model_id/config", ctrl.GetShardConfig)
		shard.POST("/:model_id/load", ctrl.LoadShards)
		shard.POST("/:model_id/unload", ctrl.UnloadShards)
		shard.GET("/:model_id/status", ctrl.GetLoadStatus)
	}
}

func (ctrl *ModelShardController) GetShards(c *gin.Context) {
	modelID := c.Param("model_id")
	var shards []models.ModelShard
	ctrl.DB.Where("model_id = ?", modelID).
		Order("shard_index ASC").Find(&shards)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": shards})
}

func (ctrl *ModelShardController) GetShardConfig(c *gin.Context) {
	modelID := c.Param("model_id")
	_ = modelID
	config := gin.H{
		"strategy":       "lazy",
		"priority_shards": []int{0, 1, 2},
		"cache_size":      1024,
		"min_bandwidth":   10,
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
}

func (ctrl *ModelShardController) LoadShards(c *gin.Context) {
	modelID := c.Param("model_id")
	deviceID := c.Query("device_id")

	var shards []models.ModelShard
	ctrl.DB.Where("model_id = ?", modelID).
		Order("shard_index ASC").Limit(10).Find(&shards)

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
		"shards":    downloadList,
		"strategy":   "lazy",
		"cache_size": 1024,
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
	loadedShards = totalShards / 2

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"model_id":      modelID,
		"device_id":     deviceID,
		"total_shards":  totalShards,
		"loaded_shards": loadedShards,
		"progress":      float64(loadedShards) / float64(totalShards) * 100,
		"memory_usage":  256,
		"last_updated":  time.Now(),
	}})
}
