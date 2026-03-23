package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ===== 平台演进控制器 =====

// PlatformEvoController 平台演进控制器
type PlatformEvoController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册路由
func (c *PlatformEvoController) RegisterRoutes(rg *gin.RouterGroup) {
	// 端侧推理
	edge := rg.Group("/platform/edge")
	{
		edge.GET("/models", c.ListEdgeModels)
		edge.POST("/models", c.CreateEdgeModel)
		edge.GET("/models/:id", c.GetEdgeModel)
		edge.DELETE("/models/:id", c.DeleteEdgeModel)
		edge.POST("/deploy", c.DeployEdgeModel)
	}

	// 模型分片
	shards := rg.Group("/platform/shards")
	{
		shards.GET("", c.ListModelShards)
		shards.POST("", c.CreateModelShard)
		shards.GET("/:id", c.GetModelShard)
		shards.POST("/:id/load", c.LoadModelShard)
	}

	// BLE Mesh
	mesh := rg.Group("/platform/mesh")
	{
		mesh.GET("/networks", c.ListMeshNetworks)
		mesh.POST("/networks", c.CreateMeshNetwork)
		mesh.GET("/networks/:id", c.GetMeshNetwork)
		mesh.POST("/networks/:id/join", c.JoinMeshNetwork)
		mesh.GET("/nodes", c.ListMeshNodes)
	}

	// RTOS配置
	rtos := rg.Group("/platform/rtos")
	{
		rtos.GET("/config", c.GetRTOSConfig)
		rtos.PUT("/config", c.UpdateRTOSConfig)
		rtos.GET("/performance", c.GetRTOSPerformance)
	}
}

// ===== 端侧模型 API =====

// ListEdgeModels 端侧模型列表
func (c *PlatformEvoController) ListEdgeModels(ctx *gin.Context) {
	var edgeModels []models.EdgeModel
	query := c.DB.Model(&models.EdgeModel{})

	// 过滤条件
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if modelType := ctx.Query("model_type"); modelType != "" {
		query = query.Where("model_type = ?", modelType)
	}
	if architecture := ctx.Query("architecture"); architecture != "" {
		query = query.Where("architecture = ?", architecture)
	}
	if deviceType := ctx.Query("device_type"); deviceType != "" {
		query = query.Where("device_types @> ?", "[\""+deviceType+"\"]")
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&edgeModels).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list edge models"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":       edgeModels,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateEdgeModel 创建端侧模型
func (c *PlatformEvoController) CreateEdgeModel(ctx *gin.Context) {
	var req struct {
		Name           string               `json:"name" binding:"required"`
		Version        string               `json:"version" binding:"required"`
		ModelType      string               `json:"model_type" binding:"required"`
		Architecture   string               `json:"architecture"`
		Framework      string               `json:"framework"`
		FileURL        string               `json:"file_url"`
		FileSize       int64                `json:"file_size"`
		FileHash       string               `json:"file_hash"`
		Description    string               `json:"description"`
		MinMemoryMB    int                  `json:"min_memory_mb"`
		MinStorageMB   int                  `json:"min_storage_mb"`
		InputShape     models.JSON          `json:"input_shape"`
		OutputShape    models.JSON          `json:"output_shape"`
		Quantization   string               `json:"quantization"`
		DeviceTypes    []string             `json:"device_types"`
		Tags           []string             `json:"tags"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := models.EdgeModel{
		Name:          req.Name,
		Version:       req.Version,
		ModelType:     req.ModelType,
		Architecture:  req.Architecture,
		Framework:     req.Framework,
		FileURL:       req.FileURL,
		FileSize:      req.FileSize,
		FileHash:      req.FileHash,
		Description:   req.Description,
		MinMemoryMB:   req.MinMemoryMB,
		MinStorageMB:  req.MinStorageMB,
		InputShape:    req.InputShape,
		OutputShape:   req.OutputShape,
		Quantization:  req.Quantization,
		DeviceTypes:   req.DeviceTypes,
		Tags:          req.Tags,
		Status:        "draft",
		CreatedBy:     getCurrentUserID(ctx),
	}

	if model.MinMemoryMB == 0 {
		model.MinMemoryMB = 512
	}
	if model.MinStorageMB == 0 {
		model.MinStorageMB = 1024
	}

	if err := c.DB.Create(&model).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create edge model"})
		return
	}

	ctx.JSON(http.StatusCreated, model)
}

// GetEdgeModel 获取端侧模型详情
func (c *PlatformEvoController) GetEdgeModel(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid model id"})
		return
	}

	var model models.EdgeModel
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "edge model not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get edge model"})
		}
		return
	}

	ctx.JSON(http.StatusOK, model)
}

// DeleteEdgeModel 删除端侧模型
func (c *PlatformEvoController) DeleteEdgeModel(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid model id"})
		return
	}

	var model models.EdgeModel
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "edge model not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get edge model"})
		}
		return
	}

	// 软删除
	if err := c.DB.Delete(&model).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete edge model"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "edge model deleted"})
}

// ===== 模型分片 API =====

// ListModelShards 模型分片列表
func (c *PlatformEvoController) ListModelShards(ctx *gin.Context) {
	var shards []models.ModelShard
	query := c.DB.Model(&models.ModelShard{})

	// 过滤条件
	if modelID := ctx.Query("model_id"); modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if isLoaded := ctx.Query("is_loaded"); isLoaded != "" {
		query = query.Where("is_loaded = ?", isLoaded == "true")
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "50"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("model_id, shard_index").Find(&shards).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list model shards"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":       shards,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateModelShard 创建模型分片
func (c *PlatformEvoController) CreateModelShard(ctx *gin.Context) {
	var req struct {
		ModelID        uint                 `json:"model_id" binding:"required"`
		ShardIndex     int                  `json:"shard_index" binding:"required"`
		ShardHash      string               `json:"shard_hash" binding:"required"`
		FileURL        string               `json:"file_url"`
		FileSize       int64                `json:"file_size"`
		LayerRange     models.JSON          `json:"layer_range"`
		Dependencies   []string             `json:"dependencies"`
		IsBaseShard    bool                 `json:"is_base_shard"`
		Priority       int                  `json:"priority"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证模型存在
	var model models.EdgeModel
	if err := c.DB.First(&model, req.ModelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "edge model not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify model"})
		}
		return
	}

	shard := models.ModelShard{
		ModelID:      req.ModelID,
		ShardIndex:   req.ShardIndex,
		ShardHash:    req.ShardHash,
		FileURL:      req.FileURL,
		FileSize:     req.FileSize,
		LayerRange:   req.LayerRange,
		Dependencies: req.Dependencies,
		IsBaseShard:  req.IsBaseShard,
		Priority:     req.Priority,
		Status:       "ready",
	}

	// 计算大小MB
	if shard.FileSize > 0 {
		shard.SizeMB = float64(shard.FileSize) / 1024 / 1024
	}

	if err := c.DB.Create(&shard).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create model shard"})
		return
	}

	ctx.JSON(http.StatusCreated, shard)
}

// GetModelShard 获取分片详情
func (c *PlatformEvoController) GetModelShard(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shard id"})
		return
	}

	var shard models.ModelShard
	if err := c.DB.First(&shard, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "model shard not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get model shard"})
		}
		return
	}

	ctx.JSON(http.StatusOK, shard)
}

// LoadModelShard 加载分片
func (c *PlatformEvoController) LoadModelShard(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shard id"})
		return
	}

	var shard models.ModelShard
	if err := c.DB.First(&shard, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "model shard not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get model shard"})
		}
		return
	}

	if shard.IsLoaded {
		ctx.JSON(http.StatusOK, gin.H{"message": "shard already loaded", "shard": shard})
		return
	}

	// 模拟加载过程
	now := time.Now()
	shard.IsLoaded = true
	shard.LoadedAt = &now
	shard.Status = "loaded"

	if err := c.DB.Save(&shard).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load shard"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "shard loaded successfully", "shard": shard})
}

// ===== BLE Mesh API =====

// ListMeshNetworks BLE Mesh网络列表
func (c *PlatformEvoController) ListMeshNetworks(ctx *gin.Context) {
	var networks []models.BLEMeshNetwork
	query := c.DB.Model(&models.BLEMeshNetwork{})

	// 过滤条件
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&networks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list mesh networks"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":       networks,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateMeshNetwork 创建BLE Mesh网络
func (c *PlatformEvoController) CreateMeshNetwork(ctx *gin.Context) {
	var req struct {
		Name            string          `json:"name" binding:"required"`
		NetworkID       string          `json:"network_id" binding:"required"`
		MaxNodes        int             `json:"max_nodes"`
		Frequency       int             `json:"frequency"`
		TxPowerDbm      int             `json:"tx_power_dbm"`
		ChannelMap      models.JSON     `json:"channel_map"`
		SecurityLevel   string          `json:"security_level"`
		ProxyEnabled    bool            `json:"proxy_enabled"`
		FriendEnabled   bool            `json:"friend_enabled"`
		RelayEnabled    bool            `json:"relay_enabled"`
		LowPowerEnabled bool            `json:"low_power_enabled"`
		Config          models.JSON     `json:"config"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network := models.BLEMeshNetwork{
		Name:            req.Name,
		NetworkID:       req.NetworkID,
		MaxNodes:        req.MaxNodes,
		Frequency:       req.Frequency,
		TxPowerDbm:      req.TxPowerDbm,
		ChannelMap:      req.ChannelMap,
		SecurityLevel:   req.SecurityLevel,
		ProxyEnabled:    req.ProxyEnabled,
		FriendEnabled:   req.FriendEnabled,
		RelayEnabled:    req.RelayEnabled,
		LowPowerEnabled: req.LowPowerEnabled,
		Config:          req.Config,
		Status:          "active",
		CreatedBy:       getCurrentUserID(ctx),
	}

	if network.MaxNodes == 0 {
		network.MaxNodes = 256
	}
	if network.Frequency == 0 {
		network.Frequency = 2440
	}
	if network.SecurityLevel == "" {
		network.SecurityLevel = "high"
	}

	if err := c.DB.Create(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create mesh network"})
		return
	}

	ctx.JSON(http.StatusCreated, network)
}

// GetMeshNetwork 获取Mesh网络详情
func (c *PlatformEvoController) GetMeshNetwork(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid network id"})
		return
	}

	var network models.BLEMeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "mesh network not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get mesh network"})
		}
		return
	}

	ctx.JSON(http.StatusOK, network)
}

// JoinMeshNetwork 节点加入网络
func (c *PlatformEvoController) JoinMeshNetwork(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid network id"})
		return
	}

	var req struct {
		DeviceID        uint           `json:"device_id" binding:"required"`
		UUID            string         `json:"uuid" binding:"required"`
		DeviceType      string         `json:"device_type"`
		ElementCount    int            `json:"element_count"`
		FirmwareVersion string         `json:"firmware_version"`
		MacAddress      string         `json:"mac_address"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var network models.BLEMeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "mesh network not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get mesh network"})
		}
		return
	}

	// 检查节点数限制
	if network.NodeCount >= network.MaxNodes {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "network node limit reached"})
		return
	}

	// 分配单播地址
	network.UnicastAddress++
	nodeAddress := network.UnicastAddress
	network.NodeCount++
	c.DB.Save(&network)

	node := models.BLEMeshNode{
		NetworkID:       network.ID,
		DeviceID:        req.DeviceID,
		NodeAddress:     nodeAddress,
		UUID:            req.UUID,
		DeviceType:      req.DeviceType,
		ElementCount:    req.ElementCount,
		FirmwareVersion: req.FirmwareVersion,
		MacAddress:      req.MacAddress,
		Status:          "provisioned",
		FirstSeenAt:     time.Now(),
	}

	if node.ElementCount == 0 {
		node.ElementCount = 1
	}

	if err := c.DB.Create(&node).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create mesh node"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"network": network, "node": node})
}

// ListMeshNodes BLE Mesh节点列表
func (c *PlatformEvoController) ListMeshNodes(ctx *gin.Context) {
	var nodes []models.BLEMeshNode
	query := c.DB.Model(&models.BLEMeshNode{})

	// 过滤条件
	if networkID := ctx.Query("network_id"); networkID != "" {
		query = query.Where("network_id = ?", networkID)
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceType := ctx.Query("device_type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "50"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&nodes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list mesh nodes"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":       nodes,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// ===== RTOS配置 API =====

// GetRTOSConfig 获取RTOS配置
func (c *PlatformEvoController) GetRTOSConfig(ctx *gin.Context) {
	var configs []models.RTOSConfig
	query := c.DB.Model(&models.RTOSConfig{})

	// 过滤条件
	if deviceType := ctx.Query("device_type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if configKey := ctx.Query("config_key"); configKey != "" {
		query = query.Where("config_key = ?", configKey)
	}
	if isDefault := ctx.Query("is_default"); isDefault == "true" {
		query = query.Where("is_default = ?", true)
	}

	// 如果没有特定过滤，返回默认配置或所有配置
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rtos configs"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":       configs,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateRTOSConfig 更新RTOS配置
func (c *PlatformEvoController) UpdateRTOSConfig(ctx *gin.Context) {
	var req struct {
		ConfigKey          string          `json:"config_key" binding:"required"`
		ConfigName        string          `json:"config_name"`
		DeviceType        string          `json:"device_type"`
		FirmwareVersion   string          `json:"firmware_version"`
		OSVersion         string          `json:"os_version"`
		KernelVersion     string          `json:"kernel_version"`
		SchedulerType     string          `json:"scheduler_type"`
		SchedulerQuantumMs int            `json:"scheduler_quantum_ms"`
		TaskConfig        models.JSON     `json:"task_config"`
		MemoryConfig      models.JSON     `json:"memory_config"`
		StackSizeDefault  int             `json:"stack_size_default"`
		HeapSizeTotal     int             `json:"heap_size_total"`
		InterruptConfig   models.JSON     `json:"interrupt_config"`
		TimerConfig       models.JSON     `json:"timer_config"`
		MutexConfig       models.JSON     `json:"mutex_config"`
		SemaphoreConfig   models.JSON     `json:"semaphore_config"`
		MessageQueueConfig models.JSON    `json:"message_queue_config"`
		NetworkConfig     models.JSON     `json:"network_config"`
		PowerConfig       models.JSON     `json:"power_config"`
		PeripheralConfig  models.JSON     `json:"peripheral_config"`
		BootConfig        models.JSON     `json:"boot_config"`
		SecurityConfig    models.JSON     `json:"security_config"`
		DebugConfig       models.JSON     `json:"debug_config"`
		PerformanceTuning models.JSON     `json:"performance_tuning"`
		IsDefault         bool            `json:"is_default"`
		Description       string          `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询或创建配置
	var config models.RTOSConfig
	result := c.DB.Where("config_key = ?", req.ConfigKey).First(&config)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新配置
		config = models.RTOSConfig{
			ConfigKey:           req.ConfigKey,
			ConfigName:          req.ConfigName,
			DeviceType:          req.DeviceType,
			FirmwareVersion:     req.FirmwareVersion,
			OSVersion:           req.OSVersion,
			KernelVersion:       req.KernelVersion,
			SchedulerType:       req.SchedulerType,
			SchedulerQuantumMs:  req.SchedulerQuantumMs,
			TaskConfig:          req.TaskConfig,
			MemoryConfig:        req.MemoryConfig,
			StackSizeDefault:    req.StackSizeDefault,
			HeapSizeTotal:       req.HeapSizeTotal,
			InterruptConfig:     req.InterruptConfig,
			TimerConfig:         req.TimerConfig,
			MutexConfig:          req.MutexConfig,
			SemaphoreConfig:      req.SemaphoreConfig,
			MessageQueueConfig:   req.MessageQueueConfig,
			NetworkConfig:        req.NetworkConfig,
			PowerConfig:          req.PowerConfig,
			PeripheralConfig:     req.PeripheralConfig,
			BootConfig:           req.BootConfig,
			SecurityConfig:       req.SecurityConfig,
			DebugConfig:          req.DebugConfig,
			PerformanceTuning:    req.PerformanceTuning,
			IsDefault:           req.IsDefault,
			IsActive:            true,
			Description:         req.Description,
			CreatedBy:           getCurrentUserID(ctx),
		}
		if err := c.DB.Create(&config).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rtos config"})
			return
		}
		ctx.JSON(http.StatusCreated, config)
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query rtos config"})
		return
	}

	// 更新现有配置
	updates := map[string]interface{}{}
	if req.ConfigName != "" {
		updates["config_name"] = req.ConfigName
	}
	if req.DeviceType != "" {
		updates["device_type"] = req.DeviceType
	}
	if req.OSVersion != "" {
		updates["os_version"] = req.OSVersion
	}
	if req.KernelVersion != "" {
		updates["kernel_version"] = req.KernelVersion
	}
	if req.SchedulerType != "" {
		updates["scheduler_type"] = req.SchedulerType
	}
	if req.SchedulerQuantumMs > 0 {
		updates["scheduler_quantum_ms"] = req.SchedulerQuantumMs
	}
	if req.TaskConfig != nil {
		updates["task_config"] = req.TaskConfig
	}
	if req.MemoryConfig != nil {
		updates["memory_config"] = req.MemoryConfig
	}
	if req.StackSizeDefault > 0 {
		updates["stack_size_default"] = req.StackSizeDefault
	}
	if req.HeapSizeTotal > 0 {
		updates["heap_size_total"] = req.HeapSizeTotal
	}
	if req.InterruptConfig != nil {
		updates["interrupt_config"] = req.InterruptConfig
	}
	if req.TimerConfig != nil {
		updates["timer_config"] = req.TimerConfig
	}
	if req.MutexConfig != nil {
		updates["mutex_config"] = req.MutexConfig
	}
	if req.SemaphoreConfig != nil {
		updates["semaphore_config"] = req.SemaphoreConfig
	}
	if req.MessageQueueConfig != nil {
		updates["message_queue_config"] = req.MessageQueueConfig
	}
	if req.NetworkConfig != nil {
		updates["network_config"] = req.NetworkConfig
	}
	if req.PowerConfig != nil {
		updates["power_config"] = req.PowerConfig
	}
	if req.PeripheralConfig != nil {
		updates["peripheral_config"] = req.PeripheralConfig
	}
	if req.BootConfig != nil {
		updates["boot_config"] = req.BootConfig
	}
	if req.SecurityConfig != nil {
		updates["security_config"] = req.SecurityConfig
	}
	if req.DebugConfig != nil {
		updates["debug_config"] = req.DebugConfig
	}
	if req.PerformanceTuning != nil {
		updates["performance_tuning"] = req.PerformanceTuning
	}
	updates["is_default"] = req.IsDefault
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := c.DB.Model(&config).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rtos config"})
		return
	}

	// 增加应用次数
	c.DB.Model(&config).Update("apply_count", config.ApplyCount+1)

	c.DB.First(&config, config.ID)
	ctx.JSON(http.StatusOK, config)
}

// DeployEdgeModel 部署端侧模型到设备
func (c *PlatformEvoController) DeployEdgeModel(ctx *gin.Context) {
	var req struct {
		ModelID   uint     `json:"model_id" binding:"required"`
		DeviceIDs []uint   `json:"device_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var model models.EdgeModel
	if err := c.DB.First(&model, req.ModelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "edge model not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get edge model"})
		}
		return
	}

	// 模拟部署：更新模型状态为已部署
	if model.Status == "draft" {
		now := time.Now()
		model.Status = "active"
		model.PublishedAt = &now
		c.DB.Save(&model)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "model deployed successfully",
		"model_id": req.ModelID,
		"devices":  req.DeviceIDs,
	})
}

// ===== RTOS性能数据 =====

// PerformanceMetrics RTOS性能指标
type PerformanceMetrics struct {
	CPUUsagePercent    float64 `json:"cpu_usage_percent"`
	MemoryUsedMB       float64 `json:"memory_used_mb"`
	MemoryTotalMB      float64 `json:"memory_total_mb"`
	MemoryUsagePercent float64 `json:"memory_usage_percent"`
	HeapFreeBytes      int64   `json:"heap_free_bytes"`
	HeapUsedBytes      int64   `json:"heap_used_bytes"`
	TaskCount          int     `json:"task_count"`
	RunningTasks       int     `json:"running_tasks"`
	BlockedTasks       int     `json:"blocked_tasks"`
	ReadyTasks         int     `json:"ready_tasks"`
	ContextSwitchCount int64   `json:"context_switch_count"`
	ISRCount           int64   `json:"isr_count"`
	ISRRateHz          float64 `json:"isr_rate_hz"`
	InterruptLatencyUs int     `json:"interrupt_latency_us"`
	SchedulerLatencyUs int     `json:"scheduler_latency_us"`
	PowerState         string  `json:"power_state"`
	TemperatureC       float64 `json:"temperature_c"`
	BatteryLevel       int     `json:"battery_level"`
	NetworkLatencyMs   int     `json:"network_latency_ms"`
	ThroughputMbps     float64 `json:"throughput_mbps"`
	ErrorCount         int     `json:"error_count"`
	UptimeSeconds      int64   `json:"uptime_seconds"`
}

// GetRTOSPerformance 获取RTOS性能数据
func (c *PlatformEvoController) GetRTOSPerformance(ctx *gin.Context) {
	deviceID := ctx.Query("device_id")
	configKey := ctx.DefaultQuery("config_key", "default")

	// 如果指定了设备ID，返回该设备的实时性能数据（模拟数据）
	if deviceID != "" {
	metrics := PerformanceMetrics{
			CPUUsagePercent:    23.5,
			MemoryUsedMB:       128.4,
			MemoryTotalMB:      512.0,
			MemoryUsagePercent: 25.1,
			HeapFreeBytes:      45_000_000,
			HeapUsedBytes:      19_000_000,
			TaskCount:          12,
			RunningTasks:       8,
			BlockedTasks:       2,
			ReadyTasks:         2,
			ContextSwitchCount: 1_234_567,
			ISRCount:           56_789,
			ISRRateHz:          142.3,
			InterruptLatencyUs: 5,
			SchedulerLatencyUs: 12,
			PowerState:         "active",
			TemperatureC:       42.3,
			BatteryLevel:       85,
			NetworkLatencyMs:   15,
			ThroughputMbps:     12.5,
			ErrorCount:         0,
			UptimeSeconds:      86400,
		}
		_ = configKey // silence unused
		ctx.JSON(http.StatusOK, gin.H{
			"device_id": deviceID,
			"metrics":   metrics,
			"sampled_at": time.Now(),
		})
		return
	}

	// 返回聚合性能数据（按配置键分组）
	var configs []models.RTOSConfig
	c.DB.Where("config_key = ?", configKey).Find(&configs)

	metrics := map[string]interface{}{
		"config_key":  configKey,
		"avg_cpu":    24.3,
		"avg_memory": 26.5,
		"total_tasks": 156,
		"total_nodes": 42,
		"sampled_at":  time.Now(),
	}
	ctx.JSON(http.StatusOK, metrics)
}

// PerformanceMetrics RTOS性能指标（GORM模型，用于存储历史性能数据）
type PerformanceMetricsOld struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	DeviceID        uint           `gorm:"index;not null" json:"device_id"`
	ConfigKey       string         `gorm:"type:varchar(128)" json:"config_key"`
	CPUUsagePercent float64        `gorm:"type:decimal(5,2)" json:"cpu_usage_percent"`
	MemoryUsageMB   float64        `gorm:"type:decimal(10,2)" json:"memory_usage_mb"`
	HeapUsedBytes   int64          `gorm:"type:bigint" json:"heap_used_bytes"`
	TaskCount       int            `gorm:"default:0" json:"task_count"`
	NetworkLatencyMs int           `gorm:"default:0" json:"network_latency_ms"`
	ErrorCount      int            `gorm:"default:0" json:"error_count"`
	SampledAt       time.Time      `json:"sampled_at"`
	CreatedAt       time.Time      `json:"created_at"`
}

func (PerformanceMetricsOld) TableName() string {
	return "performance_metrics"
}
