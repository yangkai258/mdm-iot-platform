package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExperimentController AI研究实验控制器
type ExperimentController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册实验路由
func (ctrl *ExperimentController) RegisterRoutes(rg *gin.RouterGroup) {
	experiments := rg.Group("/research/experiments")
	{
		experiments.GET("", ctrl.ListExperiments)
		experiments.POST("", ctrl.CreateExperiment)
		experiments.GET("/:id", ctrl.GetExperiment)
		experiments.PUT("/:id", ctrl.UpdateExperiment)
		experiments.DELETE("/:id", ctrl.DeleteExperiment)
		experiments.POST("/:id/start", ctrl.StartExperiment)
		experiments.POST("/:id/stop", ctrl.StopExperiment)
		experiments.GET("/:id/results", ctrl.GetExperimentResults)
		experiments.POST("/:id/collaborators", ctrl.AddCollaborator)
		experiments.GET("/:id/analysis", ctrl.GetAnalysis)
	}
}

// ListExperiments 获取实验列表
// GET /api/v1/research/experiments
func (ctrl *ExperimentController) ListExperiments(c *gin.Context) {
	platformID := c.Query("platform_id")
	status := c.Query("status")
	createdBy := c.Query("created_by")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var experiments []models.ResearchExperiment
	query := ctrl.DB.Model(&models.ResearchExperiment{})

	if platformID != "" {
		query = query.Where("platform_id = ?", platformID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&experiments)

	c.JSON(http.StatusOK, gin.H{
		"experiments": experiments,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

// CreateExperiment 创建实验
// POST /api/v1/research/experiments
func (ctrl *ExperimentController) CreateExperiment(c *gin.Context) {
	var req struct {
		PlatformID  uint                  `json:"platform_id" binding:"required"`
		Name        string                `json:"name" binding:"required"`
		Description string                `json:"description"`
		Config      map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configJSON, _ := json.Marshal(req.Config)

	experiment := models.ResearchExperiment{
		PlatformID:  req.PlatformID,
		Name:        req.Name,
		Description: req.Description,
		Config:      string(configJSON),
		Status:      "draft",
	}

	ctrl.DB.Create(&experiment)
	c.JSON(http.StatusOK, experiment)
}

// GetExperiment 获取实验详情
// GET /api/v1/research/experiments/:id
func (ctrl *ExperimentController) GetExperiment(c *gin.Context) {
	id := c.Param("id")
	var experiment models.ResearchExperiment

	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	c.JSON(http.StatusOK, experiment)
}

// UpdateExperiment 更新实验
// PUT /api/v1/research/experiments/:id
func (ctrl *ExperimentController) UpdateExperiment(c *gin.Context) {
	id := c.Param("id")
	var experiment models.ResearchExperiment

	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 禁止通过更新接口修改状态和时间字段
	delete(req, "status")
	delete(req, "started_at")
	delete(req, "completed_at")
	delete(req, "created_at")

	ctrl.DB.Model(&experiment).Updates(req)
	c.JSON(http.StatusOK, experiment)
}

// DeleteExperiment 删除实验
// DELETE /api/v1/research/experiments/:id
func (ctrl *ExperimentController) DeleteExperiment(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.ResearchExperiment{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// StartExperiment 启动实验
// POST /api/v1/research/experiments/:id/start
func (ctrl *ExperimentController) StartExperiment(c *gin.Context) {
	id := c.Param("id")
	var experiment models.ResearchExperiment

	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	if experiment.Status != "draft" && experiment.Status != "stopped" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Experiment can only be started from draft or stopped state"})
		return
	}

	now := time.Now()
	ctrl.DB.Model(&experiment).Updates(map[string]interface{}{
		"status":     "running",
		"started_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "started"})
}

// StopExperiment 停止实验
// POST /api/v1/research/experiments/:id/stop
func (ctrl *ExperimentController) StopExperiment(c *gin.Context) {
	id := c.Param("id")
	var experiment models.ResearchExperiment

	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	if experiment.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only running experiments can be stopped"})
		return
	}

	now := time.Now()
	ctrl.DB.Model(&experiment).Updates(map[string]interface{}{
		"status":       "stopped",
		"completed_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "stopped"})
}

// GetExperimentResults 获取实验结果
// GET /api/v1/research/experiments/:id/results
func (ctrl *ExperimentController) GetExperimentResults(c *gin.Context) {
	id := c.Param("id")
	var experiment models.ResearchExperiment

	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	// 查询实验结果记录
	var results []models.ResearchExperimentResult
	ctrl.DB.Where("experiment_id = ?", id).Order("created_at DESC").Find(&results)

	c.JSON(http.StatusOK, gin.H{
		"experiment_id": id,
		"status":       experiment.Status,
		"config":       experiment.Config,
		"results":      results,
	})
}

// AddCollaborator 添加实验协作者
// POST /api/v1/research/experiments/:id/collaborators
func (ctrl *ExperimentController) AddCollaborator(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查实验是否存在
	var experiment models.ResearchExperiment
	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	// 检查是否已是协作者
	var existing models.ResearchCollaborator
	if err := ctrl.DB.Where("experiment_id = ? AND user_id = ?", id, req.UserID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a collaborator"})
		return
	}

	collab := models.ResearchCollaborator{
		ExperimentID: toUintResearch(id),
		UserID:       req.UserID,
		Role:         req.Role,
		Status:       "pending",
		InvitedAt:    time.Now(),
	}

	ctrl.DB.Create(&collab)
	c.JSON(http.StatusOK, collab)
}

// GetAnalysis 获取实验分析
// GET /api/v1/research/experiments/:id/analysis
func (ctrl *ExperimentController) GetAnalysis(c *gin.Context) {
	id := c.Param("id")
	var experiment models.ResearchExperiment

	if err := ctrl.DB.First(&experiment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	// 获取所有结果
	var results []models.ResearchExperimentResult
	ctrl.DB.Where("experiment_id = ?", id).Find(&results)

	// 获取协作者列表
	var collaborators []models.ResearchCollaborator
	ctrl.DB.Where("experiment_id = ?", id).Find(&collaborators)

	// 简单的分析摘要
	var totalResults int64
	ctrl.DB.Model(&models.ResearchExperimentResult{}).Where("experiment_id = ?", id).Count(&totalResults)

	analysis := map[string]interface{}{
		"experiment_id":  id,
		"status":        experiment.Status,
		"name":          experiment.Name,
		"total_runs":    totalResults,
		"collaborators": collaborators,
		"duration":      getDuration(experiment.StartedAt, experiment.CompletedAt),
		"summary":       generateAnalysisSummary(experiment, totalResults),
	}

	c.JSON(http.StatusOK, analysis)
}

func toUintResearch(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}

func getDuration(started, completed *time.Time) interface{} {
	if started == nil {
		return nil
	}
	end := time.Now()
	if completed != nil {
		end = *completed
	}
	duration := end.Sub(*started)
	return duration.Seconds()
}

func generateAnalysisSummary(exp models.ResearchExperiment, totalResults int64) string {
	switch exp.Status {
	case "running":
		return "Experiment is currently running"
	case "completed":
		return "Experiment completed successfully"
	case "failed":
		return "Experiment failed"
	case "stopped":
		return "Experiment was stopped by user"
	default:
		return "Experiment is in draft state"
	}
}
