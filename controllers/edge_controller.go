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

// EdgeController 端侧推理控制器
type EdgeController struct {
	DB *gorm.DB
}

// CreateModelRequest 上传模型请求
type CreateModelRequest struct {
	Name           string  `json:"name" binding:"required"`
	Version        string  `json:"version" binding:"required"`
	Description    string  `json:"description"`
	ModelType      string  `json:"model_type" binding:"required"` // yolo, mobilenet, bert, whisper, tflite, onnx
	ModelFormat    string  `json:"model_format" binding:"required"` // tflite, onnx, pt, h5
	FileSize       int64   `json:"file_size"`
	FileURL        string  `json:"file_url"`
	FileMD5        string  `json:"file_md5"`
	InputShape     string  `json:"input_shape"`
	OutputShape    string  `json:"output_shape"`
	InputDtype     string  `json:"input_dtype"`
	OutputDtype    string  `json:"output_dtype"`
	Framework      string  `json:"framework"`
	HardwareTarget string  `json:"hardware_target"`
	MaxBatchSize   int     `json:"max_batch_size"`
	AvgLatencyMs   float64 `json:"avg_latency_ms"`
	Accuracy       float64 `json:"accuracy"`
	MemoryUsageKB  int     `json:"memory_usage_kb"`
}

// DeployModelRequest 部署模型请求
type DeployModelRequest struct {
	DeviceIDs        []string `json:"device_ids" binding:"required"`
	RuntimeVersion   string   `json:"runtime_version"`
}

// ListModels 获取模型列表
// GET /api/v1/edge/models
func (c *EdgeController) ListModels(ctx *gin.Context) {
	var modelList []models.EdgeModel
	query := c.DB.Model(&models.EdgeModel{})

	// 按名称/类型过滤
	if name := ctx.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if modelType := ctx.Query("model_type"); modelType != "" {
		query = query.Where("model_type = ?", modelType)
	}
	if hardwareTarget := ctx.Query("hardware_target"); hardwareTarget != "" {
		query = query.Where("hardware_target = ?", hardwareTarget)
	}

	// 分页
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize)

	if err := query.Find(&modelList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询模型列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      modelList,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateModel 上传模型
// POST /api/v1/edge/models
func (c *EdgeController) CreateModel(ctx *gin.Context) {
	var req CreateModelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 校验 org_id（从中间件获取）
	orgID, _ := ctx.Get("org_id")
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	model := models.EdgeModel{
		Name:           req.Name,
		Version:        req.Version,
		Description:    req.Description,
		ModelType:      req.ModelType,
		ModelFormat:    req.ModelFormat,
		FileSize:       req.FileSize,
		FileURL:        req.FileURL,
		FileMD5:        req.FileMD5,
		InputShape:     req.InputShape,
		OutputShape:    req.OutputShape,
		Input_dtype:    req.InputDtype,
		Output_dtype:   req.OutputDtype,
		Framework:      req.Framework,
		HardwareTarget: req.HardwareTarget,
		MaxBatchSize:   req.MaxBatchSize,
		AvgLatencyMs:   req.AvgLatencyMs,
		Accuracy:       req.Accuracy,
		MemoryUsageKB:  req.MemoryUsageKB,
		IsActive:       true,
		OrgID:          orgID.(uint),
		CreateUserID:   uid,
	}

	if model.Input_dtype == "" {
		model.Input_dtype = "float32"
	}
	if model.Output_dtype == "" {
		model.Output_dtype = "float32"
	}

	if err := c.DB.Create(&model).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "创建模型失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": model,
	})
}

// GetModel 获取模型详情
// GET /api/v1/edge/models/:id
func (c *EdgeController) GetModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的模型ID",
		})
		return
	}

	var model models.EdgeModel
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "模型不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询模型失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": model,
	})
}

// DeployModel 部署模型到设备
// POST /api/v1/edge/models/:id/deploy
func (c *EdgeController) DeployModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的模型ID",
		})
		return
	}

	var req DeployModelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 检查模型是否存在
	var model models.EdgeModel
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "模型不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询模型失败: " + err.Error(),
		})
		return
	}

	orgID, _ := ctx.Get("org_id")
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	deployments := make([]models.EdgeModelDeployment, 0, len(req.DeviceIDs))

	for _, deviceID := range req.DeviceIDs {
		// 检查是否已部署
		var existing models.EdgeModelDeployment
		if err := c.DB.Where("model_id = ? AND device_id = ?", id, deviceID).First(&existing).Error; err == nil {
			// 已存在部署记录，更新状态
			existing.Status = "deploying"
			existing.RuntimeVersion = req.RuntimeVersion
			c.DB.Save(&existing)
			deployments = append(deployments, existing)
			continue
		}

		deployment := models.EdgeModelDeployment{
			ModelID:          uint(id),
			DeviceID:         deviceID,
			Status:           "deploying",
			RuntimeVersion:   req.RuntimeVersion,
			OrgID:            orgID.(uint),
			CreateUserID:     uid,
		}
		c.DB.Create(&deployment)
		deployments = append(deployments, deployment)
	}

	// TODO: 通过 MQTT 下发部署指令到设备
	// mqttHandler.PublishDeployModel(deployments)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"model_id":       id,
			"deployments":    deployments,
			"total_devices":  len(req.DeviceIDs),
			"deployed_count": len(deployments),
		},
	})
}

// GetInference 获取端侧推理结果
// GET /api/v1/edge/inference/:device_id
func (c *EdgeController) GetInference(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "设备ID不能为空",
		})
		return
	}

	// 获取当前设备上运行的模型部署
	var deployment models.EdgeModelDeployment
	if err := c.DB.Where("device_id = ? AND status = ?", deviceID, "running").
		Preload("EdgeModel").
		First(&deployment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "设备上未运行任何模型",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备模型失败: " + err.Error(),
		})
		return
	}

	// 获取最新推理日志
	var lastLog models.EdgeInferenceLog
	c.DB.Where("device_id = ?", deviceID).Order("created_at DESC").First(&lastLog)

	// 获取推理统计（最近24小时）
	var stats struct {
		TotalCount  int64   `json:"total_count"`
		SuccessRate float64 `json:"success_rate"`
		AvgLatency  float64 `json:"avg_latency_ms"`
	}
	since24h := time.Now().Add(-24 * time.Hour)
	row := c.DB.Model(&models.EdgeInferenceLog{}).
		Select("COUNT(*) as total_count, COALESCE(AVG(latency_ms), 0) as avg_latency, "+
			"COALESCE(SUM(CASE WHEN result_code = 'success' THEN 1 ELSE 0 END)::float / NULLIF(COUNT(*), 0) * 100, 0) as success_rate").
		Where("device_id = ? AND created_at >= ?", deviceID, since24h).Row()
	row.Scan(&stats.TotalCount, &stats.AvgLatency, &stats.SuccessRate)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":     deviceID,
			"model": gin.H{
				"id":            deployment.ModelID,
				"name":          deployment.EdgeModel.Name,
				"version":       deployment.EdgeModel.Version,
				"model_type":    deployment.EdgeModel.ModelType,
				"hardware_target": deployment.EdgeModel.HardwareTarget,
			},
			"deployment": gin.H{
				"id":               deployment.ID,
				"status":           deployment.Status,
				"runtime_version":  deployment.RuntimeVersion,
				"deployed_at":      deployment.DeployedAt,
			},
			"last_inference": gin.H{
				"id":           lastLog.ID,
				"latency_ms":   lastLog.LatencyMs,
				"result_code":  lastLog.ResultCode,
				"created_at":   lastLog.CreatedAt,
			},
			"stats_24h": gin.H{
				"total_count":   stats.TotalCount,
				"success_rate":  fmt.Sprintf("%.2f%%", stats.SuccessRate),
				"avg_latency_ms": fmt.Sprintf("%.2f", stats.AvgLatency),
			},
		},
	})
}

// ListDeployments 获取设备部署列表（辅助接口）
// GET /api/v1/edge/deployments
func (c *EdgeController) ListDeployments(ctx *gin.Context) {
	var deployments []models.EdgeModelDeployment
	query := c.DB.Model(&deployments)

	deviceID := ctx.Query("device_id")
	modelID := ctx.Query("model_id")
	status := ctx.Query("status")

	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if modelID != "" {
		mid, _ := strconv.ParseUint(modelID, 10, 32)
		query = query.Where("model_id = ?", mid)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("EdgeModel")

	if err := query.Find(&deployments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询部署列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      deployments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
