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

// AITrainingController AI训练任务控制器
type AITrainingController struct {
	DB *gorm.DB
}

// NewAITrainingController 创建控制器
func NewAITrainingController(db *gorm.DB) *AITrainingController {
	return &AITrainingController{DB: db}
}

// RegisterRoutes 注册AI训练路由
func (c *AITrainingController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ai/training", c.ListTraining)
	rg.POST("/ai/training", c.CreateTraining)
	rg.GET("/ai/training/:id", c.GetTraining)
	rg.PUT("/ai/training/:id", c.UpdateTraining)
	rg.POST("/ai/training/:id/start", c.StartTraining)
	rg.POST("/ai/training/:id/pause", c.PauseTraining)
	rg.POST("/ai/training/:id/resume", c.ResumeTraining)
	rg.POST("/ai/training/:id/cancel", c.CancelTraining)
}

// ListTraining 训练任务列表
// GET /api/v1/ai/training
func (c *AITrainingController) ListTraining(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AITraining{})

	// 过滤条件
	if modelIDStr := ctx.Query("model_id"); modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if trainingType := ctx.Query("training_type"); trainingType != "" {
		query = query.Where("training_type = ?", trainingType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if createUserIDStr := ctx.Query("create_user_id"); createUserIDStr != "" {
		if createUserID, err := strconv.ParseUint(createUserIDStr, 10, 64); err == nil {
			query = query.Where("create_user_id = ?", createUserID)
		}
	}

	var total int64
	query.Count(&total)

	var trainings []models.AITraining
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&trainings).Error; err != nil {
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
			"list": trainings,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateTraining 创建训练任务
// POST /api/v1/ai/training
func (c *AITrainingController) CreateTraining(ctx *gin.Context) {
	var req models.AITrainingCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 检查关联模型是否存在
	var model models.AIModelConfig
	if err := c.DB.First(&model, req.ModelID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "关联的AI模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
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
	training := models.AITraining{
		Name:            req.Name,
		Description:     req.Description,
		ModelID:         req.ModelID,
		BaseModelKey:    req.BaseModelKey,
		TrainingType:    req.TrainingType,
		DatasetPath:     req.DatasetPath,
		DatasetSize:     req.DatasetSize,
		ValidationPath:  req.ValidationPath,
		HyperParams:     req.HyperParams,
		ResourceConfig:  req.ResourceConfig,
		Priority:        models.TrainingPriority(req.Priority),
		Status:          models.TrainingStatusPending,
		QueuedAt:        now,
		OrgID:           orgID,
		CreateUserID:    userID,
	}

	if training.BaseModelKey == "" {
		training.BaseModelKey = model.ModelKey
	}

	if err := c.DB.Create(&training).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "创建训练任务失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"task_key":     training.TaskKey,
			"id":           training.ID,
			"name":         training.Name,
			"status":       training.Status,
			"training_type": training.TrainingType,
			"created_at":   training.CreatedAt,
		},
	})
}

// GetTraining 获取训练详情
// GET /api/v1/ai/training/:id
func (c *AITrainingController) GetTraining(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var training models.AITraining
	var err error

	// 尝试作为uint ID查询
	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&training, id).Error
	} else {
		// 尝试作为task_key查询
		err = c.DB.Where("task_key = ?", idStr).First(&training).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "训练任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    training,
	})
}

// UpdateTraining 更新训练任务
// PUT /api/v1/ai/training/:id
func (c *AITrainingController) UpdateTraining(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var training models.AITraining
	var err error

	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&training, id).Error
	} else {
		err = c.DB.Where("task_key = ?", idStr).First(&training).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "训练任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 只允许更新pending状态的任务
	if training.Status != models.TrainingStatusPending {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "只允许更新pending状态的任务"})
		return
	}

	var req models.AITrainingUpdate
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
	if req.HyperParams != "" {
		updates["hyper_params"] = req.HyperParams
	}
	if req.ResourceConfig != "" {
		updates["resource_config"] = req.ResourceConfig
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	updates["updated_at"] = time.Now()

	if err := c.DB.Model(&training).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	c.DB.First(&training, training.ID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": training})
}

// StartTraining 开始训练
// POST /api/v1/ai/training/:id/start
func (c *AITrainingController) StartTraining(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var training models.AITraining
	var err error

	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&training, id).Error
	} else {
		err = c.DB.Where("task_key = ?", idStr).First(&training).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "训练任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 检查状态
	if training.Status != models.TrainingStatusPending && training.Status != models.TrainingStatusPaused {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": fmt.Sprintf("当前状态%s不允许启动", training.Status)})
		return
	}

	now := time.Now()
	newStatus := models.TrainingStatusRunning
	if training.Status == models.TrainingStatusPaused {
		newStatus = models.TrainingStatusRunning
	} else {
		newStatus = models.TrainingStatusQueued
	}

	if err := c.DB.Model(&training).Updates(map[string]interface{}{
		"status":     newStatus,
		"started_at": now,
		"updated_at": now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 模拟异步训练处理
	if training.Status == models.TrainingStatusPending {
		go c.simulateTraining(training.ID)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "训练任务已启动",
		"data": gin.H{
			"task_key": training.TaskKey,
			"status":   newStatus,
			"started_at": now,
		},
	})
}

// PauseTraining 暂停训练
// POST /api/v1/ai/training/:id/pause
func (c *AITrainingController) PauseTraining(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var training models.AITraining
	var err error

	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&training, id).Error
	} else {
		err = c.DB.Where("task_key = ?", idStr).First(&training).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "训练任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if training.Status != models.TrainingStatusRunning {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "只允许暂停running状态的任务"})
		return
	}

	if err := c.DB.Model(&training).Updates(map[string]interface{}{
		"status":     models.TrainingStatusPaused,
		"updated_at": time.Now(),
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "训练任务已暂停",
		"data": gin.H{
			"task_key": training.TaskKey,
			"status":   models.TrainingStatusPaused,
		},
	})
}

// ResumeTraining 恢复训练
// POST /api/v1/ai/training/:id/resume
func (c *AITrainingController) ResumeTraining(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var training models.AITraining
	var err error

	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&training, id).Error
	} else {
		err = c.DB.Where("task_key = ?", idStr).First(&training).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "训练任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if training.Status != models.TrainingStatusPaused {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "只允许恢复paused状态的任务"})
		return
	}

	if err := c.DB.Model(&training).Updates(map[string]interface{}{
		"status":     models.TrainingStatusRunning,
		"updated_at": time.Now(),
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "训练任务已恢复",
		"data": gin.H{
			"task_key": training.TaskKey,
			"status":   models.TrainingStatusRunning,
		},
	})
}

// CancelTraining 取消训练
// POST /api/v1/ai/training/:id/cancel
func (c *AITrainingController) CancelTraining(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var training models.AITraining
	var err error

	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&training, id).Error
	} else {
		err = c.DB.Where("task_key = ?", idStr).First(&training).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "训练任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 不允许取消已完成或已失败的任务
	if training.Status == models.TrainingStatusCompleted || training.Status == models.TrainingStatusFailed || training.Status == models.TrainingStatusCancelled {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "当前状态不允许取消"})
		return
	}

	if err := c.DB.Model(&training).Updates(map[string]interface{}{
		"status":       models.TrainingStatusCancelled,
		"completed_at": time.Now(),
		"updated_at":   time.Now(),
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "训练任务已取消",
		"data": gin.H{
			"task_key": training.TaskKey,
			"status":   models.TrainingStatusCancelled,
		},
	})
}

// simulateTraining 模拟训练过程
func (c *AITrainingController) simulateTraining(trainingID uint) {
	totalEpochs := 10
	totalSteps := int64(1000)

	for epoch := 1; epoch <= totalEpochs; epoch++ {
		// 检查是否被暂停或取消
		var training models.AITraining
		c.DB.First(&training, trainingID)
		if training.Status == models.TrainingStatusPaused || training.Status == models.TrainingStatusCancelled {
			return
		}

		for step := int64(1); step <= totalSteps/int64(totalEpochs); step++ {
			currentStep := int64((epoch-1)*(int(totalSteps)/totalEpochs)) + step
			progress := int(float64(currentStep) / float64(totalSteps) * 100)
			loss := 2.5 - (float64(currentStep) / float64(totalSteps) * 2.0)

			c.DB.Model(&models.AITraining{}).Where("id = ?", trainingID).Updates(map[string]interface{}{
				"epoch":        epoch,
				"total_epochs": totalEpochs,
				"current_step": currentStep,
				"total_steps":  totalSteps,
				"progress":     progress,
				"loss":         loss,
				"metrics":      fmt.Sprintf(`{"accuracy":%.2f,"f1":%.2f}`, 0.8+float64(progress)/500, 0.75+float64(progress)/500),
			})

			time.Sleep(50 * time.Millisecond)
		}
	}

	// 训练完成
	completedAt := time.Now()
	// 获取训练任务的task_key用于输出路径
	var completedTraining models.AITraining
	c.DB.Select("task_key").First(&completedTraining, trainingID)
	c.DB.Model(&models.AITraining{}).Where("id = ?", trainingID).Updates(map[string]interface{}{
		"status":       models.TrainingStatusCompleted,
		"progress":     100,
		"completed_at": completedAt,
		"output_path":  "/models/output/" + completedTraining.TaskKey,
	})
}
