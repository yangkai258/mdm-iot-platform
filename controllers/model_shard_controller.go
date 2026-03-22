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

// ModelShardController 模型分片管理控制器
type ModelShardController struct {
	DB *gorm.DB
}

// ListShards 获取模型分片列表
// GET /api/v1/ai/models/:id/shards
func (c *ModelShardController) ListShards(ctx *gin.Context) {
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

	version := ctx.Query("version")
	query := c.DB.Model(&models.ModelShard{}).Where("model_id = ?", modelID)
	if version != "" {
		query = query.Where("version = ?", version)
	}

	var shards []models.ModelShard
	if err := query.Order("shard_index ASC").Find(&shards).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 统计
	var totalSize int64
	var verifiedCount int64
	c.DB.Model(&models.ModelShard{}).Where("model_id = ?", modelID).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize)
	c.DB.Model(&models.ModelShard{}).Where("model_id = ? AND status = ?", modelID, models.ShardStatusVerified).Count(&verifiedCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"model_id":       modelID,
			"total_shards":   len(shards),
			"verified_shards": verifiedCount,
			"total_size":     totalSize,
			"list":           shards,
		},
	})
}

// CreateShard 创建分片记录（上传模型分片）
// POST /api/v1/ai/models/:id/shards
func (c *ModelShardController) CreateShard(ctx *gin.Context) {
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

	var req models.ModelShardCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查分片序号是否已存在
	var existing models.ModelShard
	if err := c.DB.Where("model_id = ? AND version = ? AND shard_index = ?", modelID, req.Version, req.ShardIndex).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":       4009,
			"message":    fmt.Sprintf("分片 %d 已存在", req.ShardIndex),
			"error_code": "ERR_SHARD_EXISTS",
		})
		return
	}

	shard := models.ModelShard{
		ModelID:    uint(modelID),
		Version:    req.Version,
		ShardIndex: req.ShardIndex,
		ShardName:  req.ShardName,
		FileSize:   req.FileSize,
		FilePath:   req.FilePath,
		FileMD5:    req.FileMD5,
		FileSHA256: req.FileSHA256,
		Checksum:   req.Checksum,
		Status:     models.ShardStatusPending,
		CreatedBy:  getUserID(ctx),
		OrgID:      getOrgID(ctx),
	}

	if err := c.DB.Create(&shard).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建分片记录失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "success",
		"data":    shard,
	})
}

// GetShard 获取分片详情
// GET /api/v1/ai/models/:id/shards/:shard_id
func (c *ModelShardController) GetShard(ctx *gin.Context) {
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

	shardIDStr := ctx.Param("shard_id")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的分片ID",
			"error_code": "ERR_INVALID_SHARD_ID",
		})
		return
	}

	var shard models.ModelShard
	if err := c.DB.Where("id = ? AND model_id = ?", shardID, modelID).First(&shard).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "分片不存在",
			"error_code": "ERR_SHARD_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    shard,
	})
}

// DeleteShard 删除分片
// DELETE /api/v1/ai/models/:id/shards/:shard_id
func (c *ModelShardController) DeleteShard(ctx *gin.Context) {
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

	shardIDStr := ctx.Param("shard_id")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的分片ID",
			"error_code": "ERR_INVALID_SHARD_ID",
		})
		return
	}

	result := c.DB.Where("id = ? AND model_id = ?", shardID, modelID).Delete(&models.ModelShard{})
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "分片不存在",
			"error_code": "ERR_SHARD_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// VerifyShard 验证分片完整性
// POST /api/v1/ai/models/:id/shards/:shard_id/verify
func (c *ModelShardController) VerifyShard(ctx *gin.Context) {
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

	shardIDStr := ctx.Param("shard_id")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的分片ID",
			"error_code": "ERR_INVALID_SHARD_ID",
		})
		return
	}

	var req models.ModelShardVerify
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var shard models.ModelShard
	if err := c.DB.Where("id = ? AND model_id = ?", shardID, modelID).First(&shard).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "分片不存在",
			"error_code": "ERR_SHARD_NOT_FOUND",
		})
		return
	}

	// 根据算法验证
	var verified bool
	var actualValue string
	now := time.Now()

	switch req.Algorithm {
	case "md5":
		actualValue = shard.FileMD5
	case "sha256":
		actualValue = shard.FileSHA256
	case "checksum":
		actualValue = shard.Checksum
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "不支持的校验算法，仅支持 md5/sha256/checksum",
			"error_code": "ERR_UNSUPPORTED_ALGORITHM",
		})
		return
	}

	if req.Expected != "" {
		verified = actualValue == req.Expected
	} else {
		// 无 expected 时，只返回当前记录的校验值
		verified = actualValue != ""
	}

	status := models.ShardStatusVerified
	verifyMessage := fmt.Sprintf("验证通过 (%s)", req.Algorithm)
	if !verified {
		status = models.ShardStatusFailed
		verifyMessage = fmt.Sprintf("验证失败：期望值 %s，实际值 %s", req.Expected, actualValue)
	}

	c.DB.Model(&shard).Updates(map[string]interface{}{
		"status":         status,
		"verify_message": verifyMessage,
		"verified_at":    &now,
		"verified_by":    getUserID(ctx),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"shard_id":      shard.ID,
			"verified":       verified,
			"algorithm":      req.Algorithm,
			"expected":       req.Expected,
			"actual":         actualValue,
			"verify_message": verifyMessage,
		},
	})
}

// DeploySharded 部署分片模型
// POST /api/v1/ai/models/:id/deploy/sharded
func (c *ModelShardController) DeploySharded(ctx *gin.Context) {
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

	var req models.DeployShardedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查目标版本的所有分片是否都已验证
	var shards []models.ModelShard
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, req.Version).Find(&shards).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询分片失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	if len(shards) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "该版本没有分片记录",
			"error_code": "ERR_NO_SHARDS",
		})
		return
	}

	var unverified []int
	for _, s := range shards {
		if s.Status != models.ShardStatusVerified {
			unverified = append(unverified, s.ShardIndex)
		}
	}
	if len(unverified) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    fmt.Sprintf("以下分片未验证通过: %v", unverified),
			"error_code": "ERR_SHARDS_NOT_VERIFIED",
		})
		return
	}

	// 查询版本记录
	var version models.ModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelID, req.Version).First(&version).Error; err != nil {
		// 版本不存在则自动创建
		version = models.ModelVersion{
			ModelID:       uint(modelID),
			Version:       req.Version,
			Status:        models.VersionStatusDraft,
			IsSharded:     true,
			TotalShards:   len(shards),
			VerifiedShards: len(shards),
			CreateUserID:  getUserID(ctx),
			OrgID:         getOrgID(ctx),
		}
		c.DB.Create(&version)
	}

	// 创建部署记录
	deploy := models.DeployShardedModel{
		ModelID:        uint(modelID),
		VersionID:      version.ID,
		Version:        req.Version,
		TargetEnv:      req.TargetEnv,
		ReplicaCount:   req.ReplicaCount,
		ResourceConfig: req.ResourceConfig,
		Status:         "deploying",
		StartedAt:      time.Now(),
		DeployedBy:     getUserID(ctx),
		OrgID:          getOrgID(ctx),
	}

	if err := c.DB.Create(&deploy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建部署记录失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 更新模型状态
	c.DB.Model(&model).Updates(map[string]interface{}{
		"status":     models.ModelStatusDeploying,
		"deployed_at": time.Now(),
		"deployed_by": getUserID(ctx),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"deploy_id":     deploy.ID,
			"model_id":      modelID,
			"version":       req.Version,
			"target_env":    req.TargetEnv,
			"total_shards":  len(shards),
			"status":        "deploying",
		},
	})
}

// UpdateShard 更新分片信息（仅更新文件路径、校验值等）
// PUT /api/v1/ai/models/:id/shards/:shard_id
func (c *ModelShardController) UpdateShard(ctx *gin.Context) {
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

	shardIDStr := ctx.Param("shard_id")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4001,
			"message":    "无效的分片ID",
			"error_code": "ERR_INVALID_SHARD_ID",
		})
		return
	}

	var req models.ModelShardUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var shard models.ModelShard
	if err := c.DB.Where("id = ? AND model_id = ?", shardID, modelID).First(&shard).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "分片不存在",
			"error_code": "ERR_SHARD_NOT_FOUND",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.FilePath != "" {
		updates["file_path"] = req.FilePath
	}
	if req.FileMD5 != "" {
		updates["file_md5"] = req.FileMD5
	}
	if req.FileSHA256 != "" {
		updates["file_sha256"] = req.FileSHA256
	}
	if req.Checksum != "" {
		updates["checksum"] = req.Checksum
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "没有需要更新的字段",
			"error_code": "ERR_NO_UPDATES",
		})
		return
	}

	if err := c.DB.Model(&shard).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "更新失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	c.DB.First(&shard, shard.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    shard,
	})
}
