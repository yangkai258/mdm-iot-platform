package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModelVersionController 模型版本管理控制器
type ModelVersionController struct {
	DB *gorm.DB
}

// NewModelVersionController 创建控制器
func NewModelVersionController(db *gorm.DB) *ModelVersionController {
	return &ModelVersionController{DB: db}
}

// ListModels 模型列表
// GET /api/v1/ai/models
func (c *ModelVersionController) ListModels(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIModelVersion{})

	if modelID := ctx.Query("model_id"); modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}
	if modelName := ctx.Query("model_name"); modelName != "" {
		query = query.Where("model_name LIKE ?", "%"+modelName+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var versions []models.AIModelVersion
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&versions).Error; err != nil {
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
			"list": versions,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateModel 注册新模型（第一个版本）
// POST /api/v1/ai/models
func (c *ModelVersionController) CreateModel(ctx *gin.Context) {
	var req struct {
		ModelID   string `json:"model_id" binding:"required"`
		ModelName string `json:"model_name" binding:"required"`
		Version   string `json:"version" binding:"required"`
		ModelPath string `json:"model_path"`
		Config    string `json:"config"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 检查是否已存在
	var existing models.AIModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", req.ModelID, req.Version).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":    4009,
			"message": "模型版本已存在",
		})
		return
	}

	// 获取当前用户ID（从context）
	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	version := models.AIModelVersion{
		ModelID:     req.ModelID,
		ModelName:   req.ModelName,
		Version:     req.Version,
		Status:      "testing",
		ModelPath:   req.ModelPath,
		Config:      req.Config,
		PublishedBy: userID,
		PublishedAt: time.Now(),
	}

	if err := c.DB.Create(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "创建失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    version,
	})
}

// GetModel 模型详情
// GET /api/v1/ai/models/:id
func (c *ModelVersionController) GetModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的ID"})
		return
	}

	var version models.AIModelVersion
	if err := c.DB.First(&version, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": version})
}

// UpdateModel 更新模型信息
// PUT /api/v1/ai/models/:id
func (c *ModelVersionController) UpdateModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的ID"})
		return
	}

	var version models.AIModelVersion
	if err := c.DB.First(&version, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	var req models.AIModelVersionUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.ModelName != "" {
		updates["model_name"] = req.ModelName
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.ModelPath != "" {
		updates["model_path"] = req.ModelPath
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	if req.Metrics != "" {
		updates["metrics"] = req.Metrics
	}
	updates["updated_at"] = time.Now()

	if err := c.DB.Model(&version).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	c.DB.First(&version, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": version})
}

// DeleteModel 删除模型
// DELETE /api/v1/ai/models/:id
func (c *ModelVersionController) DeleteModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的ID"})
		return
	}

	var version models.AIModelVersion
	if err := c.DB.First(&version, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 不允许删除生产环境版本
	if version.Status == "production" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "不允许删除生产环境版本"})
		return
	}

	if err := c.DB.Delete(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ListVersions 版本列表
// GET /api/v1/ai/models/:id/versions
func (c *ModelVersionController) ListVersions(ctx *gin.Context) {
	modelID := ctx.Param("id")

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIModelVersion{}).Where("model_id = ?", modelID)
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var versions []models.AIModelVersion
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&versions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": versions,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateVersion 发布新版本
// POST /api/v1/ai/models/:id/versions
func (c *ModelVersionController) CreateVersion(ctx *gin.Context) {
	modelID := ctx.Param("id")

	var req models.AIModelVersionCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": err.Error()})
		return
	}

	// 检查模型是否存在
	var existingModel models.AIModelVersion
	if err := c.DB.Where("model_id = ?", modelID).First(&existingModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 检查版本是否已存在
	var existingVersion models.AIModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, req.Version).First(&existingVersion).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 4009, "message": "版本已存在"})
		return
	}

	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	version := models.AIModelVersion{
		ModelID:     modelID,
		ModelName:   existingModel.ModelName,
		Version:     req.Version,
		Status:      "testing",
		ModelPath:   req.ModelPath,
		Config:      req.Config,
		Metrics:     req.Metrics,
		PublishedBy: userID,
		PublishedAt: time.Now(),
	}

	if err := c.DB.Create(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": version})
}

// DeprecateVersion 废弃版本
// POST /api/v1/ai/models/:id/deprecate
func (c *ModelVersionController) DeprecateVersion(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的ID"})
		return
	}

	var version models.AIModelVersion
	if err := c.DB.First(&version, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if version.Status == "production" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "请先回滚生产环境版本"})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&version).Updates(map[string]interface{}{
		"status":        "deprecated",
		"deprecated_at": now,
		"updated_at":    now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	c.DB.First(&version, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": version})
}
