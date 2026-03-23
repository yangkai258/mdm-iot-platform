package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModelShardController 模型分片管理控制器
type ModelShardController struct {
	DB *gorm.DB
}

// NewModelShardController 创建控制器实例
func NewModelShardController(db *gorm.DB) *ModelShardController {
	return &ModelShardController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *ModelShardController) RegisterRoutes(rg *gin.RouterGroup) {
	shard := rg.Group("/model-shards")
	{
		// 分片管理
		shard.GET("", ctrl.ListShards)                      // 列出所有分片(分页)
		shard.GET("/:model_id", ctrl.GetShards)              // 获取模型的所有分片
		shard.GET("/:model_id/config", ctrl.GetShardConfig)  // 获取分片配置
		shard.GET("/:model_id/status", ctrl.GetLoadStatus)   // 获取加载状态
		shard.POST("/:model_id/load", ctrl.LoadShards)       // 加载分片
		shard.POST("/:model_id/unload", ctrl.UnloadShards)    // 卸载分片
		shard.POST("/:model_id/prefetch", ctrl.PrefetchShards) // 预取分片
		shard.PUT("/:model_id/shards/:shard_id/status", ctrl.UpdateShardStatus) // 更新分片状态
		
		// Edge Model 管理
		edgeModel := rg.Group("/edge-models")
		{
			edgeModel.GET("", ctrl.ListEdgeModels)
			edgeModel.GET("/:id", ctrl.GetEdgeModel)
			edgeModel.POST("", ctrl.CreateEdgeModel)
			edgeModel.PUT("/:id", ctrl.UpdateEdgeModel)
			edgeModel.DELETE("/:id", ctrl.DeleteEdgeModel)
			edgeModel.POST("/:id/publish", ctrl.PublishEdgeModel)
			edgeModel.POST("/:id/deprecate", ctrl.DeprecateEdgeModel)
		}
	}
}

// ListShards 列出所有分片(分页)
func (ctrl *ModelShardController) ListShards(c *gin.Context) {
	var shards []models.ModelShard
	var total int64

	query := ctrl.DB.Model(&models.ModelShard{})

	// 按模型ID筛选
	if modelID := c.Query("model_id"); modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}

	// 按状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 按是否已加载筛选
	if isLoaded := c.Query("is_loaded"); isLoaded != "" {
		query = query.Where("is_loaded = ?", isLoaded)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("model_id ASC, shard_index ASC").Find(&shards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      shards,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetShards 获取模型的所有分片
func (ctrl *ModelShardController) GetShards(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	var shards []models.ModelShard
	ctrl.DB.Where("model_id = ?", uint(modelID)).
		Order("shard_index ASC").Find(&shards)

	if len(shards) == 0 {
		// 如果没有分片，返回空数组而不是错误
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []models.ModelShard{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": shards})
}

// GetShardConfig 获取分片配置
func (ctrl *ModelShardController) GetShardConfig(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	// 获取模型信息
	var edgeModel models.EdgeModel
	if err := ctrl.DB.First(&edgeModel, uint(modelID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模型不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 获取分片数量
	var shardCount int64
	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ?", uint(modelID)).Count(&shardCount)

	config := gin.H{
		"model_id":         modelID,
		"model_name":       edgeModel.Name,
		"version":          edgeModel.Version,
		"strategy":         "lazy", // 懒加载策略
		"priority_shards":  []int{0, 1, 2}, // 优先加载的分片索引
		"cache_size":       edgeModel.MinMemoryMB,
		"min_bandwidth":    10, // Mbps
		"total_shards":     shardCount,
		"quantization":     edgeModel.Quantization,
		"input_shape":      edgeModel.InputShape,
		"output_shape":     edgeModel.OutputShape,
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
}

// GetLoadStatus 获取加载状态
func (ctrl *ModelShardController) GetLoadStatus(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}
	deviceID := c.Query("device_id")

	var totalShards int64
	var loadedShards int64

	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ?", uint(modelID)).Count(&totalShards)
	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ? AND is_loaded = ?", uint(modelID), true).Count(&loadedShards)

	progress := 0.0
	if totalShards > 0 {
		progress = float64(loadedShards) / float64(totalShards) * 100
	}

	// 计算内存使用
	var memoryUsage int64
	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ? AND is_loaded = ?", uint(modelID), true).
		Select("COALESCE(SUM(memory_usage_mb), 0)").Scan(&memoryUsage)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"model_id":      uint(modelID),
		"device_id":     deviceID,
		"total_shards":  totalShards,
		"loaded_shards": loadedShards,
		"progress":      progress,
		"memory_usage":  memoryUsage,
		"last_updated":  time.Now(),
	}})
}

// LoadShards 加载分片
func (ctrl *ModelShardController) LoadShards(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}
	deviceID := c.Query("device_id")

	// 获取优先加载的分片
	var shards []models.ModelShard
	query := ctrl.DB.Where("model_id = ?", uint(modelID)).
		Order("priority DESC, shard_index ASC").Limit(10)

	// 如果指定了优先分片，只加载优先分片
	if priorityOnly := c.Query("priority_only"); priorityOnly == "true" {
		query = query.Where("priority > 0")
	}

	query.Find(&shards)

	if len(shards) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "没有需要加载的分片", "data": gin.H{
			"device_id": deviceID,
			"shards":    []interface{}{},
		}})
		return
	}

	downloadList := make([]map[string]interface{}, 0)
	for _, s := range shards {
		downloadList = append(downloadList, map[string]interface{}{
			"shard_id":     s.ID,
			"shard_index":  s.ShardIndex,
			"url":          s.FileURL,
			"size":         s.FileSize,
			"hash":         s.ShardHash,
			"layer_range":  s.LayerRange,
			"dependencies": s.Dependencies,
			"is_base":      s.IsBaseShard,
		})
	}

	// 更新分片状态为加载中
	now := time.Now()
	for _, s := range shards {
		ctrl.DB.Model(&s).Updates(map[string]interface{}{
			"status":           "loading",
			"load_duration_ms": 0,
		})
		_ = now // suppress unused warning
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"device_id":  deviceID,
		"model_id":   uint(modelID),
		"shards":     downloadList,
		"strategy":   "lazy",
		"cache_size": 1024,
	}})
}

// UnloadShards 卸载分片
func (ctrl *ModelShardController) UnloadShards(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}
	deviceID := c.Query("device_id")

	// 获取已加载的分片
	var shards []models.ModelShard
	ctrl.DB.Where("model_id = ? AND is_loaded = ?", uint(modelID), true).Find(&shards)

	if len(shards) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "没有已加载的分片", "data": gin.H{
			"model_id":  uint(modelID),
			"device_id": deviceID,
		}})
		return
	}

	// 更新分片状态
	for _, s := range shards {
		ctrl.DB.Model(&s).Updates(map[string]interface{}{
			"is_loaded": false,
			"status":    "ready",
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "分片已卸载", "data": gin.H{
		"model_id":   uint(modelID),
		"device_id":  deviceID,
		"unloaded":   len(shards),
	}})
}

// PrefetchShards 预取分片
func (ctrl *ModelShardController) PrefetchShards(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}
	deviceID := c.Query("device_id")

	// 获取高优先级的分片
	var shards []models.ModelShard
	ctrl.DB.Where("model_id = ? AND is_loaded = ? AND priority > 0", uint(modelID), false).
		Order("priority DESC, shard_index ASC").Limit(5).Find(&shards)

	if len(shards) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "没有需要预取的分片", "data": gin.H{
			"device_id": deviceID,
			"shards":    []interface{}{},
		}})
		return
	}

	downloadList := make([]map[string]interface{}, 0)
	for _, s := range shards {
		downloadList = append(downloadList, map[string]interface{}{
			"shard_id":    s.ID,
			"shard_index": s.ShardIndex,
			"url":         s.FileURL,
			"size":        s.FileSize,
			"hash":        s.ShardHash,
			"priority":    s.Priority,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"device_id": deviceID,
		"model_id":  uint(modelID),
		"shards":    downloadList,
	}})
}

// UpdateShardStatus 更新分片状态
func (ctrl *ModelShardController) UpdateShardStatus(c *gin.Context) {
	modelIDStr := c.Param("model_id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	shardIDStr := c.Param("shard_id")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的分片ID"})
		return
	}

	var updateData struct {
		IsLoaded       *bool  `json:"is_loaded"`
		Status         string `json:"status"`
		ErrorMessage   string `json:"error_message"`
		MemoryUsageMB  *int   `json:"memory_usage_mb"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var shard models.ModelShard
	if err := ctrl.DB.First(&shard, uint(shardID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "分片不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 验证分片属于指定模型
	if shard.ModelID != uint(modelID) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "分片不属于指定模型"})
		return
	}

	// 更新字段
	updates := map[string]interface{}{}
	if updateData.IsLoaded != nil {
		updates["is_loaded"] = *updateData.IsLoaded
		if *updateData.IsLoaded {
			updates["loaded_at"] = time.Now()
		}
	}
	if updateData.Status != "" {
		updates["status"] = updateData.Status
	}
	if updateData.ErrorMessage != "" {
		updates["error_message"] = updateData.ErrorMessage
	}
	if updateData.MemoryUsageMB != nil {
		updates["memory_usage_mb"] = *updateData.MemoryUsageMB
	}

	if err := ctrl.DB.Model(&shard).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&shard, shardID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": shard})
}

// ===== Edge Model 管理 =====

// ListEdgeModels 列出边缘模型
func (ctrl *ModelShardController) ListEdgeModels(c *gin.Context) {
	var edgeModels []models.EdgeModel
	var total int64

	query := ctrl.DB.Model(&models.EdgeModel{})

	// 搜索条件
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR version LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按类型筛选
	if modelType := c.Query("model_type"); modelType != "" {
		query = query.Where("model_type = ?", modelType)
	}

	// 按状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&edgeModels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      edgeModels,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetEdgeModel 获取边缘模型详情
func (ctrl *ModelShardController) GetEdgeModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	var edgeModel models.EdgeModel
	if err := ctrl.DB.First(&edgeModel, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模型不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 获取关联的分片数量
	var shardCount int64
	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ?", uint(id)).Count(&shardCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"model":      edgeModel,
			"shard_count": shardCount,
		},
	})
}

// CreateEdgeModel 创建边缘模型
func (ctrl *ModelShardController) CreateEdgeModel(c *gin.Context) {
	var edgeModel models.EdgeModel
	if err := c.ShouldBindJSON(&edgeModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	edgeModel.Status = "draft"

	if err := ctrl.DB.Create(&edgeModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": edgeModel})
}

// UpdateEdgeModel 更新边缘模型
func (ctrl *ModelShardController) UpdateEdgeModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	var edgeModel models.EdgeModel
	if err := ctrl.DB.First(&edgeModel, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模型不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		Name          string `json:"name"`
		Version       string `json:"version"`
		ModelType     string `json:"model_type"`
		FileURL       string `json:"file_url"`
		Description   string `json:"description"`
		MinMemoryMB   *int   `json:"min_memory_mb"`
		Quantization  string `json:"quantization"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.Name != "" {
		updates["name"] = updateData.Name
	}
	if updateData.Version != "" {
		updates["version"] = updateData.Version
	}
	if updateData.ModelType != "" {
		updates["model_type"] = updateData.ModelType
	}
	if updateData.FileURL != "" {
		updates["file_url"] = updateData.FileURL
	}
	if updateData.Description != "" {
		updates["description"] = updateData.Description
	}
	if updateData.MinMemoryMB != nil {
		updates["min_memory_mb"] = *updateData.MinMemoryMB
	}
	if updateData.Quantization != "" {
		updates["quantization"] = updateData.Quantization
	}

	if err := ctrl.DB.Model(&edgeModel).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&edgeModel, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": edgeModel})
}

// DeleteEdgeModel 删除边缘模型
func (ctrl *ModelShardController) DeleteEdgeModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	// 检查是否有分片关联
	var shardCount int64
	ctrl.DB.Model(&models.ModelShard{}).Where("model_id = ?", uint(id)).Count(&shardCount)
	if shardCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该模型下还有分片，无法删除"})
		return
	}

	if err := ctrl.DB.Delete(&models.EdgeModel{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PublishEdgeModel 发布边缘模型
func (ctrl *ModelShardController) PublishEdgeModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	now := time.Now()
	result := ctrl.DB.Model(&models.EdgeModel{}).Where("id = ? AND status = ?", uint(id), "draft").
		Updates(map[string]interface{}{
			"status":       "active",
			"published_at": now,
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "模型不存在或已发布"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// DeprecateEdgeModel 弃用边缘模型
func (ctrl *ModelShardController) DeprecateEdgeModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的模型ID"})
		return
	}

	result := ctrl.DB.Model(&models.EdgeModel{}).Where("id = ?", uint(id)).
		Update("status", "deprecated")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模型不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
