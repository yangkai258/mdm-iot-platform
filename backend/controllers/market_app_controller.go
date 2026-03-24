package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MarketAppController 应用市场控制器
type MarketAppController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册应用市场路由
func (c *MarketAppController) RegisterRoutes(api *gin.RouterGroup) {
	apps := api.Group("/market/apps")
	{
		apps.GET("", c.ListApps)
		apps.GET("/:id", c.GetApp)
		apps.POST("", c.CreateApp)
		apps.PUT("/:id", c.UpdateApp)
		apps.DELETE("/:id", c.DeleteApp)
		apps.POST("/:id/publish", c.PublishApp)
		apps.POST("/:id/unpublish", c.UnpublishApp)
		apps.GET("/categories", c.ListCategories)
		apps.GET("/featured", c.GetFeaturedApps)
		apps.POST("/:id/rate", c.RateApp)
	}
}

// ListApps 获取应用列表
func (c *MarketAppController) ListApps(ctx *gin.Context) {
	var apps []models.App
	var total int64

	query := c.DB.Model(&models.App{})

	// 搜索
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 类型筛选
	if appType := ctx.Query("type"); appType != "" {
		query = query.Where("app_type = ?", appType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      apps,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetApp 获取应用详情
func (c *MarketAppController) GetApp(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.App
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 获取版本信息
	var versions []models.AppVersion
	c.DB.Where("app_id = ?", id).Order("version_code DESC").Limit(5).Find(&versions)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"app":      app,
		"versions": versions,
	}})
}

// CreateApp 创建应用
func (c *MarketAppController) CreateApp(ctx *gin.Context) {
	var app models.App
	if err := ctx.ShouldBindJSON(&app); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	app.Status = 1 // 1: enabled (draft)
	if err := c.DB.Create(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": app})
}

// UpdateApp 更新应用
func (c *MarketAppController) UpdateApp(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.App
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		IconURL     string `json:"icon_url"`
		ScreenshotURLs string `json:"screenshot_urls"`
		Version     string `json:"version"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.Name != "" {
		updates["name"] = updateData.Name
	}
	if updateData.Description != "" {
		updates["description"] = updateData.Description
	}
	if updateData.Category != "" {
		updates["category"] = updateData.Category
	}
	if updateData.IconURL != "" {
		updates["icon_url"] = updateData.IconURL
	}
	if updateData.ScreenshotURLs != "" {
		updates["screenshot_urls"] = updateData.ScreenshotURLs
	}
	if updateData.Version != "" {
		updates["version"] = updateData.Version
	}

	if err := c.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&app, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": app})
}

// DeleteApp 删除应用
func (c *MarketAppController) DeleteApp(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.App
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	if err := c.DB.Delete(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PublishApp 发布应用
func (c *MarketAppController) PublishApp(ctx *gin.Context) {
	id := ctx.Param("id")
	result := c.DB.Model(&models.App{}).Where("id = ? AND status = ?", id, 1).
		Update("status", 1)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布失败"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "应用不存在或已发布"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UnpublishApp 下架应用
func (c *MarketAppController) UnpublishApp(ctx *gin.Context) {
	id := ctx.Param("id")
	result := c.DB.Model(&models.App{}).Where("id = ? AND status = ?", id, 1).
		Update("status", 0)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "应用不存在或未发布"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ListCategories 获取应用分类
func (c *MarketAppController) ListCategories(ctx *gin.Context) {
	categories := []gin.H{
		{"id": "education", "name": "教育", "icon": "book"},
		{"id": "entertainment", "name": "娱乐", "icon": "gamepad"},
		{"id": "health", "name": "健康", "icon": "heart"},
		{"id": "social", "name": "社交", "icon": "users"},
		{"id": "productivity", "name": "效率", "icon": "tools"},
		{"id": "lifestyle", "name": "生活", "icon": "home"},
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": categories})
}

// GetFeaturedApps 获取精选应用
func (c *MarketAppController) GetFeaturedApps(ctx *gin.Context) {
	var apps []models.App
	c.DB.Where("status = ?", 1).
		Order("created_at DESC").Limit(10).Find(&apps)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": apps})
}

// RateApp 评价应用
func (c *MarketAppController) RateApp(ctx *gin.Context) {
	id := ctx.Param("id")
	appID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的应用ID"})
		return
	}

	var ratingData struct {
		UserID  uint    `json:"user_id" binding:"required"`
		Score   float64 `json:"score" binding:"required,min=1,max=5"`
		Comment string  `json:"comment"`
	}
	if err := ctx.ShouldBindJSON(&ratingData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 创建评分记录
	rating := models.Rating{
		ItemType: "app",
		ItemID:   uint(appID),
		UserID:   ratingData.UserID,
		Rating:   int(ratingData.Score),
		Review:   ratingData.Comment,
		CreatedAt: time.Now(),
	}

	if err := c.DB.Create(&rating).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "评价失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rating})
}
