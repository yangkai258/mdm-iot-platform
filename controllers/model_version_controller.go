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

// ModelVersionController 模型版本管理控制器
type ModelVersionController struct {
	DB *gorm.DB
}

// ListVersions 获取模型版本列表
// GET /api/v1/ai/models/:id/versions
func (c *ModelVersionController) ListVersions(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	// 验证模型是否存在
	var model models.AIModelConfig
	if err := c.DB.First(&model, modelID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "模型不存在",
			"error_code": "ERR_MODEL_NOT_FOUND",
		})
		return
	}

	status := ctx.Query("status")
	query := c.DB.Model(&models.ModelVersion{}).Where("model_id = ?", modelID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var versions []models.ModelVersion
	if err := query.Order("created_at DESC").Find(&versions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"model_id":    modelID,
			"total_count": len(versions),
			"list":        versions,
		},
	})
}

// CreateVersion 创建新版本
// POST /api/v1/ai/models/:id/versions
func (c *ModelVersionController) CreateVersion(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	// 验证模型是否存在
	var model models.AIModelConfig
	if err := c.DB.First(&model, modelID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "模型不存在",
			"error_code": "ERR_MODEL_NOT_FOUND",
		})
		return
	}

	var req models.ModelVersionCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查版本号是否已存在
	var existing models.ModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, req.Version).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":       4009,
			"message":    fmt.Sprintf("版本 %s 已存在", req.Version),
			"error_code": "ERR_VERSION_EXISTS",
		})
		return
	}

	version := models.ModelVersion{
		ModelID:      uint(modelID),
		Version:      req.Version,
		ModelName:    req.ModelName,
		Description:  req.Description,
		Status:       models.VersionStatusDraft,
		ModelPath:    req.ModelPath,
		Config:       req.Config,
		Metrics:      req.Metrics,
		IsSharded:   req.IsSharded,
		TotalShards: req.TotalShards,
		CreateUserID: getUserID(ctx),
		OrgID:        getOrgID(ctx),
	}

	if err := c.DB.Create(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建版本失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "success",
		"data":    version,
	})
}

// GetVersion 获取版本详情
// GET /api/v1/ai/models/:id/versions/:version
func (c *ModelVersionController) GetVersion(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	versionStr := ctx.Param("version")

	var version models.ModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, versionStr).First(&version).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "版本不存在",
			"error_code": "ERR_VERSION_NOT_FOUND",
		})
		return
	}

	// 获取分片统计
	var shards []models.ModelShard
	c.DB.Where("model_id = ? AND version = ?", modelID, versionStr).Find(&shards)
	var totalSize int64
	var verifiedCount int64
	for _, s := range shards {
		totalSize += s.FileSize
		if s.Status == models.ShardStatusVerified {
			verifiedCount++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"version":          version,
			"shards":           shards,
			"total_shards":     len(shards),
			"verified_shards":   verifiedCount,
			"total_size":        totalSize,
		},
	})
}

// UpdateVersion 更新版本信息
// PUT /api/v1/ai/models/:id/versions/:version
func (c *ModelVersionController) UpdateVersion(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	versionStr := ctx.Param("version")

	var version models.ModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, versionStr).First(&version).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "版本不存在",
			"error_code": "ERR_VERSION_NOT_FOUND",
		})
		return
	}

	// 只有 draft 状态可以更新
	if version.Status != models.VersionStatusDraft {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "只有草稿状态可以更新",
			"error_code": "ERR_VERSION_NOT_EDITABLE",
		})
		return
	}

	var req models.ModelVersionUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.ModelName != "" {
		updates["model_name"] = req.ModelName
	}
	if req.Description != "" {
		updates["description"] = req.Description
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
	if req.TotalShards > 0 {
		updates["total_shards"] = req.TotalShards
	}
	updates["is_sharded"] = req.IsSharded

	if err := c.DB.Model(&version).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "更新失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	c.DB.First(&version, version.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    version,
	})
}

// PublishVersion 发布版本
// POST /api/v1/ai/models/:id/versions/:version/publish
func (c *ModelVersionController) PublishVersion(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	versionStr := ctx.Param("version")

	var version models.ModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, versionStr).First(&version).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "版本不存在",
			"error_code": "ERR_VERSION_NOT_FOUND",
		})
		return
	}

	// 检查分片模型的所有分片是否已验证
	if version.IsSharded {
		var shards []models.ModelShard
		c.DB.Where("model_id = ? AND version = ?", modelID, versionStr).Find(&shards)
		if len(shards) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":       4005,
				"message":    "分片模型必须上传所有分片后才能发布",
				"error_code": "ERR_SHARDS_INCOMPLETE",
			})
			return
		}
		for _, s := range shards {
			if s.Status != models.ShardStatusVerified {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":       4005,
					"message":    fmt.Sprintf("分片 %d 尚未验证通过", s.ShardIndex),
					"error_code": "ERR_SHARDS_NOT_VERIFIED",
				})
				return
			}
		}
	}

	// 只能是 draft 或 testing 才能发布
	if version.Status != models.VersionStatusDraft && version.Status != models.VersionStatusTesting {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "当前状态不允许发布",
			"error_code": "ERR_VERSION_NOT_PUBLISHABLE",
		})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":       models.VersionStatusProduction,
		"published_by":  getUserID(ctx),
		"published_at":  &now,
	}

	if err := c.DB.Model(&version).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "发布失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 更新模型的最新版本
	c.DB.Model(&models.AIModelConfig{}).Where("id = ?", modelID).Updates(map[string]interface{}{
		"status": models.ModelStatusOnline,
	})

	c.DB.First(&version, version.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"version":      version,
			"published_at": now,
		},
	})
}

// Rollback 回滚到上一版本
// POST /api/v1/ai/models/:id/rollback
func (c *ModelVersionController) Rollback(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	var req models.ModelRollback
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 查找当前生产版本
	var currentVersion models.ModelVersion
	if err := c.DB.Where("model_id = ? AND status = ?", modelID, models.VersionStatusProduction).First(&currentVersion).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "当前没有生产版本",
			"error_code": "ERR_NO_PRODUCTION_VERSION",
		})
		return
	}

	// 查找目标版本
	var targetVersion models.ModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, req.TargetVersion).First(&targetVersion).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "目标版本不存在",
			"error_code": "ERR_VERSION_NOT_FOUND",
		})
		return
	}

	// 目标版本不能是已废弃的
	if targetVersion.Status == models.VersionStatusDeprecated {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "不能回滚到已废弃的版本",
			"error_code": "ERR_VERSION_DEPRECATED",
		})
		return
	}

	now := time.Now()
	// 将当前版本标记为回滚来源
	c.DB.Model(&currentVersion).Updates(map[string]interface{}{
		"status":        models.VersionStatusDeprecated,
		"deprecated_at":  &now,
		"rollback_to":   req.TargetVersion,
	})

	// 激活目标版本
	c.DB.Model(&targetVersion).Updates(map[string]interface{}{
		"status":        models.VersionStatusProduction,
		"published_by":  getUserID(ctx),
		"published_at":  &now,
		"rollback_from": currentVersion.Version,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"rolled_back_from":  currentVersion.Version,
			"rolled_back_to":    targetVersion.Version,
			"reason":            req.Reason,
			"rolled_back_at":    now,
		},
	})
}

// GetVersionsSimple 获取版本列表（不带分片信息，用于下拉选择）
// GET /api/v1/ai/models/:id/versions/simple
func (c *ModelVersionController) GetVersionsSimple(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")
	modelID, err := strconv.ParseUint(modelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的模型ID",
			"error_code": "ERR_INVALID_MODEL_ID",
		})
		return
	}

	var versions []models.ModelVersion
	if err := c.DB.Model(&models.ModelVersion{}).
		Where("model_id = ?", modelID).
		Select("id, version, model_name, status, is_sharded, total_shards, verified_shards, created_at").
		Order("created_at DESC").
		Find(&versions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    versions,
	})
}
