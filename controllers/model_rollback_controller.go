package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelRollbackController 模型热回滚控制器
type ModelRollbackController struct {
	DB *gorm.DB
}

// NewModelRollbackController 创建控制器
func NewModelRollbackController(db *gorm.DB) *ModelRollbackController {
	return &ModelRollbackController{DB: db}
}

// Rollback 热回滚
// POST /api/v1/ai/models/:id/rollback
// 请求体: {"target_version": "v1.2.3", "reason": "..."}
func (c *ModelRollbackController) Rollback(ctx *gin.Context) {
	modelIDStr := ctx.Param("id")

	var req models.AIModelRollback
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": err.Error()})
		return
	}

	// 查找目标版本
	var targetVersion models.AIModelVersion
	if err := c.DB.Where("model_id = ? AND version = ?", modelIDStr, req.TargetVersion).First(&targetVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "目标版本不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 查找当前生产版本
	var currentProd models.AIModelVersion
	c.DB.Where("model_id = ? AND status = ?", modelIDStr, "production").First(&currentProd)

	// 获取用户ID
	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	// 创建回滚任务
	taskID := "rollback-" + uuid.New().String()
	task := models.AIRollbackTask{
		TaskID:      taskID,
		ModelID:     modelIDStr,
		FromVersion: currentProd.Version,
		ToVersion:   req.TargetVersion,
		Reason:      req.Reason,
		Status:      "in_progress",
		TriggeredBy: userID,
	}

	if err := c.DB.Create(&task).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 异步执行回滚（避免阻塞）
	go func() {
		// 将原生产版本标记为 staging
		if currentProd.ID != 0 {
			c.DB.Model(&currentProd).Updates(map[string]interface{}{
				"status":     "staging",
				"updated_at": time.Now(),
			})
		}

		// 将目标版本升级为 production
		c.DB.Model(&targetVersion).Updates(map[string]interface{}{
			"status":        "production",
			"rollback_from": currentProd.Version,
			"updated_at":    time.Now(),
		})

		// 记录行为日志
		behaviorLog := models.AIBehaviorLog{
			LogID:        uuid.New().String(),
			ModelID:      modelIDStr,
			ModelName:    targetVersion.ModelName,
			ModelVersion: req.TargetVersion,
			EventType:    "rollback",
			Status:       "success",
			Metadata:     `{"task_id":"` + taskID + `","reason":"` + req.Reason + `"}`,
			CreatedAt:    time.Now(),
		}
		c.DB.Create(&behaviorLog)

		// 更新任务状态
		now := time.Now()
		c.DB.Model(&task).Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": &now,
			"updated_at":   now,
		})
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"task_id": taskID,
			"status":  "started",
		},
	})
}

// GetRollbackTask 获取回滚任务状态
// GET /api/v1/ai/rollback/tasks/:task_id
func (c *ModelRollbackController) GetRollbackTask(ctx *gin.Context) {
	taskID := ctx.Param("task_id")

	var task models.AIRollbackTask
	if err := c.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": task})
}

// GetRollbackHistory 获取回滚历史
// GET /api/v1/ai/rollback/history
func (c *ModelRollbackController) GetRollbackHistory(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIRollbackTask{})

	if modelID := ctx.Query("model_id"); modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var tasks []models.AIRollbackTask
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": tasks,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// RollbackModelVersion 旧版回滚兼容（通过版本ID）
// POST /api/v1/ai/models/:id/rollback
func (c *ModelRollbackController) RollbackByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的ID"})
		return
	}

	var version models.AIModelVersion
	if err := c.DB.First(&version, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "版本不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 获取用户ID
	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	var req struct {
		Reason string `json:"reason"`
	}
	ctx.ShouldBindJSON(&req)

	// 查找当前生产版本
	var currentProd models.AIModelVersion
	c.DB.Where("model_id = ? AND status = ?", version.ModelID, "production").First(&currentProd)

	taskID := "rollback-" + uuid.New().String()
	task := models.AIRollbackTask{
		TaskID:      taskID,
		ModelID:     version.ModelID,
		FromVersion: currentProd.Version,
		ToVersion:   version.Version,
		Reason:      req.Reason,
		Status:      "in_progress",
		TriggeredBy: userID,
	}

	if err := c.DB.Create(&task).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 异步执行
	go func() {
		if currentProd.ID != 0 {
			c.DB.Model(&currentProd).Updates(map[string]interface{}{
				"status":     "staging",
				"updated_at": time.Now(),
			})
		}
		c.DB.Model(&version).Updates(map[string]interface{}{
			"status":        "production",
			"rollback_from": currentProd.Version,
			"updated_at":    time.Now(),
		})
		now := time.Now()
		c.DB.Model(&task).Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": &now,
			"updated_at":   now,
		})
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"task_id": taskID,
			"status":  "started",
		},
	})
}
