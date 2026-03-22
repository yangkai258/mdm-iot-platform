package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActionMarketController 动作市场控制器
type ActionMarketController struct {
	DB *gorm.DB
}

// RegisterActionMarketRoutes 注册动作市场路由
func (ctrl *ActionMarketController) RegisterActionMarketRoutes(api *gin.RouterGroup) {
	api.GET("/actions/market", ctrl.ListMarket)
	api.GET("/actions/market/:id", ctrl.GetMarketItem)
	api.POST("/actions", ctrl.Create)
	api.GET("/actions/:id", ctrl.Get)
	api.PUT("/actions/:id", ctrl.Update)
	api.DELETE("/actions/:id", ctrl.Delete)
	api.POST("/actions/:id/publish", ctrl.Publish)
}

// ListMarket 获取动作市场列表
// GET /api/v1/actions/market
func (ctrl *ActionMarketController) ListMarket(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	category := c.Query("category")
	keyword := c.Query("keyword")
	isPremium := c.Query("is_premium")
	isFeatured := c.Query("is_featured")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.ActionMarket{}).Where("status = ?", "approved")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("action_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if isPremium != "" {
		if isPremium == "true" {
			query = query.Where("is_premium = ?", true)
		}
	}
	if isFeatured != "" {
		if isFeatured == "true" {
			query = query.Where("is_featured = ?", true)
		}
	}

	// 排序
	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}

	// 特殊排序处理
	switch sortBy {
	case "rating":
		query = query.Order("rating " + order)
	case "download_count":
		query = query.Order("download_count " + order)
	case "use_count":
		query = query.Order("use_count " + order)
	case "price":
		query = query.Order("price " + order)
	default:
		query = query.Order(sortBy + " " + order)
	}

	var total int64
	query.Count(&total)

	var items []models.ActionMarket
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query action market: " + err.Error(),
		})
		return
	}

	responses := make([]*models.ActionMarketResponse, len(items))
	for i := range items {
		responses[i] = items[i].ToResponse()
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

// GetMarketItem 获取市场动作详情
// GET /api/v1/actions/market/:id
func (ctrl *ActionMarketController) GetMarketItem(c *gin.Context) {
	id := c.Param("id")

	var item models.ActionMarket
	if err := ctrl.DB.Where("id = ? OR market_uuid = ?", id, id).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": item.ToResponse(),
	})
}

// Create 创建自定义动作
// POST /api/v1/actions
func (ctrl *ActionMarketController) Create(c *gin.Context) {
	var req struct {
		ActionName    string `json:"action_name" binding:"required"`
		ActionNameEn  string `json:"action_name_en"`
		Category      string `json:"category" binding:"required"`
		Description   string `json:"description"`
		DurationMs    int    `json:"duration_ms" binding:"required"`
		Priority      int    `json:"priority"`
		IsEmergency   bool   `json:"is_emergency"`
		Parameters    string `json:"parameters"`
		AnimationData string `json:"animation_data"`
		MotorCommands string `json:"motor_commands"`
		AudioFile     string `json:"audio_file"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(c)

	if req.Priority == 0 {
		req.Priority = 5
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

	action := models.CustomAction{
		UserID:        userID,
		ActionName:    req.ActionName,
		ActionNameEn:  req.ActionNameEn,
		Category:      req.Category,
		Description:   req.Description,
		DurationMs:    req.DurationMs,
		Priority:      req.Priority,
		IsEmergency:   req.IsEmergency,
		Parameters:    req.Parameters,
		AnimationData: req.AnimationData,
		MotorCommands: req.MotorCommands,
		AudioFile:     req.AudioFile,
		IsPublished:   false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
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

// Get 获取自定义动作详情
// GET /api/v1/actions/:id
func (ctrl *ActionMarketController) Get(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var action models.CustomAction
	if err := ctrl.DB.Where("id = ? OR action_uuid = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限
	if action.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// Update 更新自定义动作
// PUT /api/v1/actions/:id
func (ctrl *ActionMarketController) Update(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var action models.CustomAction
	if err := ctrl.DB.Where("id = ? OR action_uuid = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限
	if action.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	var req struct {
		ActionName    string `json:"action_name"`
		ActionNameEn  string `json:"action_name_en"`
		Category      string `json:"category"`
		Description   string `json:"description"`
		DurationMs    int    `json:"duration_ms"`
		Priority      int    `json:"priority"`
		IsEmergency   *bool  `json:"is_emergency"`
		Parameters    string `json:"parameters"`
		AnimationData string `json:"animation_data"`
		MotorCommands string `json:"motor_commands"`
		AudioFile     string `json:"audio_file"`
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

// Delete 删除自定义动作
// DELETE /api/v1/actions/:id
func (ctrl *ActionMarketController) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var action models.CustomAction
	if err := ctrl.DB.Where("id = ? OR action_uuid = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限
	if action.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
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

// Publish 发布动作到市场
// POST /api/v1/actions/:id/publish
func (ctrl *ActionMarketController) Publish(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var req struct {
		PreviewVideo string   `json:"preview_video"`
		ThumbnailURL string   `json:"thumbnail_url"`
		Tags         []string `json:"tags"`
		Price        float64  `json:"price"`
		Description  string   `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 使用空结构
		req = struct {
			PreviewVideo string   `json:"preview_video"`
			ThumbnailURL string   `json:"thumbnail_url"`
			Tags         []string `json:"tags"`
			Price        float64  `json:"price"`
			Description  string   `json:"description"`
		}{}
	}

	var action models.CustomAction
	if err := ctrl.DB.Where("id = ? OR action_uuid = ?", id, id).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Action not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限
	if action.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	// 获取用户名（简化处理）
	creatorName := "User"
	if userID > 0 {
		creatorName = "User-" + string(rune(userID))
	}

	now := time.Now()
	marketItem := models.ActionMarket{
		ActionID:      action.ID,
		OriginalUUID:  action.ActionUUID,
		CreatorID:     userID,
		CreatorName:   creatorName,
		ActionName:    action.ActionName,
		ActionNameEn:  action.ActionNameEn,
		Category:      action.Category,
		Description:   action.Description,
		PreviewVideo:  req.PreviewVideo,
		ThumbnailURL:  req.ThumbnailURL,
		DurationMs:    action.DurationMs,
		Parameters:    action.Parameters,
		AnimationData: action.AnimationData,
		MotorCommands: action.MotorCommands,
		Tags:          req.Tags,
		Price:         req.Price,
		Currency:      "CNY",
		Status:        "pending", // 待审核
		PublishedAt:   &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := ctrl.DB.Create(&marketItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to publish action: " + err.Error()})
		return
	}

	// 更新原始动作的发布状态
	ctrl.DB.Model(&action).Updates(map[string]interface{}{
		"is_published": true,
		"published_at": now,
		"updated_at":   now,
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Action published to market, pending review",
		"data": gin.H{
			"market_uuid": marketItem.MarketUUID,
			"status":      marketItem.Status,
			"published_at": marketItem.PublishedAt.Format(time.RFC3339),
		},
	})
}
