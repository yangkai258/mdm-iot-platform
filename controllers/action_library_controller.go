package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActionLibraryController 动作库控制器
type ActionLibraryController struct {
	DB *gorm.DB
}

// RegisterActionLibraryRoutes 注册动作库路由
func (ctrl *ActionLibraryController) RegisterActionLibraryRoutes(api *gin.RouterGroup) {
	api.GET("/action-library", ctrl.List)
	api.GET("/action-library/categories", ctrl.GetCategories)
	api.GET("/action-library/:id", ctrl.Get)
	api.POST("/action-library", ctrl.Create)
	api.PUT("/action-library/:id", ctrl.Update)
	api.DELETE("/action-library/:id", ctrl.Delete)
}

// List 获取动作列表（分页+筛选）
// GET /api/v1/action-library
func (ctrl *ActionLibraryController) List(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	category := c.Query("category")
	keyword := c.Query("keyword")
	isEmergency := c.Query("is_emergency")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.ActionLibrary{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("action_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if isEmergency != "" {
		if isEmergency == "true" {
			query = query.Where("is_emergency = ?", true)
		} else if isEmergency == "false" {
			query = query.Where("is_emergency = ?", false)
		}
	}

	// 排序
	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	var total int64
	query.Count(&total)

	var actions []models.ActionLibrary
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query action library: " + err.Error(),
		})
		return
	}

	responses := make([]*models.ActionLibraryResponse, len(actions))
	for i := range actions {
		responses[i] = actions[i].ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      responses,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetCategories 获取动作分类列表
// GET /api/v1/action-library/categories
func (ctrl *ActionLibraryController) GetCategories(c *gin.Context) {
	categories := []map[string]interface{}{
		{"key": models.ActionCategoryEmotion, "label": "情绪动作", "count": 0},
		{"key": models.ActionCategoryGreeting, "label": "问候动作", "count": 0},
		{"key": models.ActionCategoryPlay, "label": "玩耍动作", "count": 0},
		{"key": models.ActionCategoryUtility, "label": "功能动作", "count": 0},
	}

	// 统计每个分类的数量
	for i := range categories {
		var count int64
		ctrl.DB.Model(&models.ActionLibrary{}).Where("category = ?", categories[i]["key"]).Count(&count)
		categories[i]["count"] = count
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": categories,
	})
}

// Get 获取动作详情
// GET /api/v1/action-library/:id
func (ctrl *ActionLibraryController) Get(c *gin.Context) {
	id := c.Param("id")

	var action models.ActionLibrary
	if err := ctrl.DB.Where("id = ? OR action_id = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// Create 创建动作
// POST /api/v1/action-library
func (ctrl *ActionLibraryController) Create(c *gin.Context) {
	var req struct {
		ActionName       string `json:"action_name" binding:"required"`
		ActionNameEn     string `json:"action_name_en"`
		Category         string `json:"category" binding:"required"`
		Description      string `json:"description"`
		DurationMs       int    `json:"duration_ms" binding:"required"`
		Priority         int    `json:"priority"`
		IsEmergency      bool   `json:"is_emergency"`
		CompatibleModels string `json:"compatible_models"`
		Parameters       string `json:"parameters"`
		AnimationData    string `json:"animation_data"`
		MotorCommands    string `json:"motor_commands"`
		AudioFile        string `json:"audio_file"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	if req.Priority == 0 {
		req.Priority = 5
	}
	if req.CompatibleModels == "" {
		req.CompatibleModels = "[]"
	}
	if req.Parameters == "" {
		req.Parameters = "{}"
	}
	if req.AnimationData == "" {
		req.AnimationData = "{}"
	}
	if req.MotorCommands == "" {
		req.MotorCommands = "{}"
	}

	action := models.ActionLibrary{
		ActionName:       req.ActionName,
		ActionNameEn:     req.ActionNameEn,
		Category:         req.Category,
		Description:      req.Description,
		DurationMs:       req.DurationMs,
		Priority:         req.Priority,
		IsEmergency:      req.IsEmergency,
		CompatibleModels: req.CompatibleModels,
		Parameters:       req.Parameters,
		AnimationData:    req.AnimationData,
		MotorCommands:    req.MotorCommands,
		AudioFile:        req.AudioFile,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := ctrl.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create action: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// Update 更新动作
// PUT /api/v1/action-library/:id
func (ctrl *ActionLibraryController) Update(c *gin.Context) {
	id := c.Param("id")

	var action models.ActionLibrary
	if err := ctrl.DB.Where("id = ? OR action_id = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	var req struct {
		ActionName       string `json:"action_name"`
		ActionNameEn     string `json:"action_name_en"`
		Category         string `json:"category"`
		Description      string `json:"description"`
		DurationMs       int    `json:"duration_ms"`
		Priority         int    `json:"priority"`
		IsEmergency      *bool  `json:"is_emergency"`
		CompatibleModels string `json:"compatible_models"`
		Parameters       string `json:"parameters"`
		AnimationData    string `json:"animation_data"`
		MotorCommands    string `json:"motor_commands"`
		AudioFile        string `json:"audio_file"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.ActionName != "" {
		updates["action_name"] = req.ActionName
	}
	if req.ActionNameEn != "" {
		updates["action_name_en"] = req.ActionNameEn
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DurationMs > 0 {
		updates["duration_ms"] = req.DurationMs
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	if req.IsEmergency != nil {
		updates["is_emergency"] = *req.IsEmergency
	}
	if req.CompatibleModels != "" {
		updates["compatible_models"] = req.CompatibleModels
	}
	if req.Parameters != "" {
		updates["parameters"] = req.Parameters
	}
	if req.AnimationData != "" {
		updates["animation_data"] = req.AnimationData
	}
	if req.MotorCommands != "" {
		updates["motor_commands"] = req.MotorCommands
	}
	if req.AudioFile != "" {
		updates["audio_file"] = req.AudioFile
	}

	if err := ctrl.DB.Model(&action).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update action: " + err.Error()})
		return
	}

	ctrl.DB.First(&action, action.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// Delete 删除动作
// DELETE /api/v1/action-library/:id
func (ctrl *ActionLibraryController) Delete(c *gin.Context) {
	id := c.Param("id")

	var action models.ActionLibrary
	if err := ctrl.DB.Where("id = ? OR action_id = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if err := ctrl.DB.Delete(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete action: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Action deleted successfully",
	})
}
