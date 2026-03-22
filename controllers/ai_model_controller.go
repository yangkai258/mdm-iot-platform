package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIModelController AI模型管理控制器
type AIModelController struct {
	DB *gorm.DB
}

// NewAIModelController 创建控制器
func NewAIModelController(db *gorm.DB) *AIModelController {
	return &AIModelController{DB: db}
}

// RegisterRoutes 注册AI模型路由
func (c *AIModelController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ai/models", c.ListModels)
	rg.POST("/ai/models", c.CreateModel)
	rg.GET("/ai/models/:id", c.GetModel)
	rg.PUT("/ai/models/:id", c.UpdateModel)
	rg.DELETE("/ai/models/:id", c.DeleteModel)
	rg.POST("/ai/models/:id/deploy", c.DeployModel)
	rg.GET("/ai/models/:id/deploy-history", c.GetDeployHistory)
}

// ListModels 模型列表
// GET /api/v1/ai/models
func (c *AIModelController) ListModels(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIModelConfig{})

	// 过滤条件
	if provider := ctx.Query("provider"); provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if modelType := ctx.Query("model_type"); modelType != "" {
		query = query.Where("model_type = ?", modelType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR model_key LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var modelsList []models.AIModelConfig
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&modelsList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": modelsList,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateModel 上传/注册模型
// POST /api/v1/ai/models
func (c *AIModelController) CreateModel(ctx *gin.Context) {
	var req models.AIModelUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	// 获取组织ID
	var orgID uint
	if oid, exists := ctx.Get("org_id"); exists {
		if id, ok := oid.(uint); ok {
			orgID = id
		}
	}

	// 检查模型key是否已存在
	var existing models.AIModelConfig
	if err := c.DB.Where("model_key = ?", req.Name).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":    4009,
			"message": "模型已存在",
		})
		return
	}

	now := time.Now()
	model := models.AIModelConfig{
		Name:         req.Name,
		ModelKey:     req.Name,
		Description:  req.Description,
		Provider:     models.AIProvider(req.Provider),
		ModelType:    req.ModelType,
		ModelSize:    req.ModelSize,
		FilePath:     req.FilePath,
		FileSize:     req.FileSize,
		Checksum:     req.Checksum,
		Config:       req.Config,
		Capabilities: req.Capabilities,
		Status:       models.ModelStatusPending,
		CreateUserID: userID,
		OrgID:        orgID,
	}

	if req.FilePath != "" {
		model.Status = models.ModelStatusReady
	}

	if err := c.DB.Create(&model).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "创建失败: " + err.Error(),
		})
		return
	}

	// 如果提供了文件路径，直接更新状态为ready
	if req.FilePath != "" {
		c.DB.Model(&model).Updates(map[string]interface{}{
			"status":     models.ModelStatusReady,
			"updated_at": now,
		})
		model.Status = models.ModelStatusReady
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    model,
	})
}

// GetModel 获取模型详情
// GET /api/v1/ai/models/:id
func (c *AIModelController) GetModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的模型ID"})
		return
	}

	var model models.AIModelConfig
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": model})
}

// UpdateModel 更新模型信息
// PUT /api/v1/ai/models/:id
func (c *AIModelController) UpdateModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的模型ID"})
		return
	}

	var model models.AIModelConfig
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	var req models.AIModelUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ModelType != "" {
		updates["model_type"] = req.ModelType
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	if req.Capabilities != "" {
		updates["capabilities"] = req.Capabilities
	}
	if req.QuotaDaily > 0 {
		updates["quota_daily"] = req.QuotaDaily
	}
	if req.QuotaMonthly > 0 {
		updates["quota_monthly"] = req.QuotaMonthly
	}
	if req.PricePer1K > 0 {
		updates["price_per_1k"] = req.PricePer1K
	}
	updates["updated_at"] = time.Now()

	if err := c.DB.Model(&model).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	c.DB.First(&model, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": model})
}

// DeleteModel 删除模型
// DELETE /api/v1/ai/models/:id
func (c *AIModelController) DeleteModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的模型ID"})
		return
	}

	var model models.AIModelConfig
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 不允许删除在线模型
	if model.Status == models.ModelStatusOnline {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "不允许删除在线模型，请先下线"})
		return
	}

	if err := c.DB.Delete(&model).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// DeployModel 部署模型
// POST /api/v1/ai/models/:id/deploy
func (c *AIModelController) DeployModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的模型ID"})
		return
	}

	var model models.AIModelConfig
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 检查模型状态是否可以部署
	if model.Status == models.ModelStatusOnline {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "模型已在生产环境上线"})
		return
	}
	if model.Status != models.ModelStatusReady && model.Status != models.ModelStatusOffline {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": fmt.Sprintf("模型状态为%s，无法部署", model.Status)})
		return
	}

	var req models.AIModelDeployRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 使用默认值
		req = models.AIModelDeployRequest{
			TargetEnv:    "staging",
			ReplicaCount: 1,
		}
	}

	// 获取当前用户ID
	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	now := time.Now()

	// 创建部署历史记录
	history := models.AIModelDeployHistory{
		ModelID:        uint(id),
		TargetEnv:      req.TargetEnv,
		ReplicaCount:   req.ReplicaCount,
		ResourceConfig: req.ResourceConfig,
		Status:         "deploying",
		StartedAt:      now,
		DeployedBy:     userID,
	}

	if err := c.DB.Create(&history).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建部署记录失败: " + err.Error()})
		return
	}

	// 更新模型状态为部署中
	c.DB.Model(&model).Updates(map[string]interface{}{
		"status":     models.ModelStatusDeploying,
		"updated_at": now,
	})

	// 模拟部署完成（实际应该调用K8s/Docker等编排服务）
	go func() {
		// 模拟部署延迟
		time.Sleep(2 * time.Second)

		// 更新部署历史
		completedAt := time.Now()
		c.DB.Model(&models.AIModelDeployHistory{}).Where("id = ?", history.ID).Updates(map[string]interface{}{
			"status":       "success",
			"completed_at": completedAt,
		})

		// 更新模型状态为在线
		c.DB.Model(&models.AIModelConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":      models.ModelStatusOnline,
			"deployed_at": completedAt,
			"deployed_by": userID,
			"updated_at":  completedAt,
		})
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "部署任务已创建",
		"data": gin.H{
			"deploy_id":    history.ID,
			"model_id":     id,
			"status":       "deploying",
			"target_env":   req.TargetEnv,
			"started_at":   now,
		},
	})
}

// GetDeployHistory 获取部署历史
// GET /api/v1/ai/models/:id/deploy-history
func (c *AIModelController) GetDeployHistory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的模型ID"})
		return
	}

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	var total int64
	query := c.DB.Model(&models.AIModelDeployHistory{}).Where("model_id = ?", id)
	query.Count(&total)

	var history []models.AIModelDeployHistory
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&history).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": history,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}
