package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIInferenceController AI推理控制器
type AIInferenceController struct {
	DB *gorm.DB
}

// NewAIInferenceController 创建控制器
func NewAIInferenceController(db *gorm.DB) *AIInferenceController {
	return &AIInferenceController{DB: db}
}

// RegisterRoutes 注册AI推理路由
func (c *AIInferenceController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/ai/inference", c.CreateInference)
	rg.GET("/ai/inference/:id", c.GetInference)
	rg.GET("/ai/inference", c.ListInferences)
	rg.POST("/ai/inference/:id/cancel", c.CancelInference)
}

// CreateInference 创建推理请求
// POST /api/v1/ai/inference
func (c *AIInferenceController) CreateInference(ctx *gin.Context) {
	var req models.AIInferenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 检查模型是否存在
	var model models.AIModelConfig
	if err := c.DB.First(&model, req.ModelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 检查模型是否在线
	if model.Status != models.ModelStatusOnline {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "模型未上线，无法推理"})
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

	now := time.Now()
	inference := models.AIInference{
		ModelID:       req.ModelID,
		ModelKey:      model.ModelKey,
		Mode:          req.Mode,
		InputData:     req.InputData,
		Prompt:        req.Prompt,
		Status:        models.InferenceStatusPending,
		StreamEnabled: req.StreamEnabled,
		DeviceID:      req.DeviceID,
		UserID:        userID,
		OrgID:         orgID,
		CreateUserID:  userID,
	}

	if err := c.DB.Create(&inference).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "创建推理任务失败: " + err.Error(),
		})
		return
	}

	// 模拟异步推理处理
	go c.processInference(inference.ID, req, model)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "推理任务已创建",
		"data": gin.H{
			"inference_key": inference.InferenceKey,
			"model_id":      inference.ModelID,
			"model_key":     inference.ModelKey,
			"status":        inference.Status,
			"created_at":    now,
		},
	})
}

// processInference 异步处理推理
func (c *AIInferenceController) processInference(inferenceID uint, req models.AIInferenceRequest, model models.AIModelConfig) {
	now := time.Now()

	// 更新状态为运行中
	c.DB.Model(&models.AIInference{}).Where("id = ?", inferenceID).Updates(map[string]interface{}{
		"status":     models.InferenceStatusRunning,
		"started_at": now,
	})

	// 模拟推理延迟（根据请求复杂度）
	latencyMs := int64(100 + req.MaxTokens*10)
	time.Sleep(time.Duration(latencyMs) * time.Millisecond)

	// 模拟生成响应
	response := generateMockResponse(req.Prompt, req.Mode, model.ModelType)
	inputTokens := len(req.Prompt) / 4 // 粗略估算
	outputTokens := len(response) / 4

	completedAt := time.Now()

	// 计算成本
	cost := float64(inputTokens+outputTokens) / 1000.0 * model.PricePer1K

	// 更新推理结果
	c.DB.Model(&models.AIInference{}).Where("id = ?", inferenceID).Updates(map[string]interface{}{
		"status":        models.InferenceStatusCompleted,
		"response":      response,
		"input_tokens":  inputTokens,
		"output_tokens": outputTokens,
		"total_tokens":  inputTokens + outputTokens,
		"latency_ms":    completedAt.Sub(now).Milliseconds(),
		"cost":          cost,
		"completed_at":  completedAt,
	})
}

// generateMockResponse 生成模拟响应
func generateMockResponse(prompt, mode, modelType string) string {
	switch mode {
	case "chat":
		return `{"role":"assistant","content":"这是一条模拟的AI回复。您发送的消息是: ` + prompt + `"}`
	case "completion":
		return "这是模拟的文本补全结果: " + prompt + " [补全内容]"
	case "embedding":
		return `{"embedding":[0.1,0.2,0.3,0.4,0.5],"dimensions":5}`
	case "image_gen":
		return `{"image_url":"https://example.com/generated_image.png","prompt":"` + prompt + `"}`
	default:
		return "模拟响应: " + prompt
	}
}

// GetInference 获取推理结果
// GET /api/v1/ai/inference/:id
func (c *AIInferenceController) GetInference(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// 尝试作为uint ID查询
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		var inference models.AIInference
		if err := c.DB.First(&inference, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "推理任务不存在"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    inference,
		})
		return
	}

	// 尝试作为inference_key查询
	var inference models.AIInference
	if err := c.DB.Where("inference_key = ?", idStr).First(&inference).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "推理任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    inference,
	})
}

// ListInferences 推理列表
// GET /api/v1/ai/inference
func (c *AIInferenceController) ListInferences(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIInference{})

	// 过滤条件
	if modelIDStr := ctx.Query("model_id"); modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if modelKey := ctx.Query("model_key"); modelKey != "" {
		query = query.Where("model_key = ?", modelKey)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	var total int64
	query.Count(&total)

	var inferences []models.AIInference
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&inferences).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": inferences,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CancelInference 取消推理
// POST /api/v1/ai/inference/:id/cancel
func (c *AIInferenceController) CancelInference(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var inference models.AIInference
	var err error

	// 尝试作为uint ID查询
	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&inference, id).Error
	} else {
		// 尝试作为inference_key查询
		err = c.DB.Where("inference_key = ?", idStr).First(&inference).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "推理任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 只允许取消pending或running状态的任务
	if inference.Status != models.InferenceStatusPending && inference.Status != models.InferenceStatusRunning {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "当前状态不允许取消"})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&inference).Updates(map[string]interface{}{
		"status":       models.InferenceStatusFailed,
		"error_message": "用户取消",
		"completed_at":  now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "推理任务已取消",
		"data": gin.H{
			"inference_key": inference.InferenceKey,
			"status":       models.InferenceStatusFailed,
		},
	})
}
