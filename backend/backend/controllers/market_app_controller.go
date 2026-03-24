package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MarketAppController App市场控制器
type MarketAppController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册App市场路由
func (ctrl *MarketAppController) RegisterRoutes(api *gin.RouterGroup) {
	apps := api.Group("/market/apps")
	{
		apps.GET("", ctrl.ListApps)
		apps.GET("/:id", ctrl.GetApp)
		apps.POST("", ctrl.CreateApp)
		apps.PUT("/:id", ctrl.UpdateApp)
		apps.DELETE("/:id", ctrl.DeleteApp)
	}
}

// ListApps GET /api/v1/market/apps - 应用列表
func (ctrl *MarketAppController) ListApps(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)
	keyword := c.Query("keyword")
	category := c.Query("category")
	platform := c.Query("platform")
	status := c.DefaultQuery("status", "1") // 默认只显示启用的
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.App{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if platform != "" {
		query = query.Where("platform = ? OR platform = 'multi'", platform)
	}

	validSortFields := map[string]bool{
		"created_at": true, "name": true, "category": true,
		"developer": true,
	}
	if validSortFields[sortBy] {
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Count(&total)

	var apps []models.App
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":       apps,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetApp GET /api/v1/market/apps/:id - 应用详情
func (ctrl *MarketAppController) GetApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的应用ID"})
		return
	}

	var app models.App
	if err := ctrl.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取最新版本信息
	var latestVersion models.AppVersion
	ctrl.DB.Where("app_id = ? AND is_active = ?", id, true).
		Order("created_at DESC").First(&latestVersion)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"app":            app,
			"latest_version": latestVersion,
		},
	})
}

// CreateApp POST /api/v1/market/apps - 发布应用
func (ctrl *MarketAppController) CreateApp(c *gin.Context) {
	var req MarketAppCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	app := models.App{
		Name:        req.Name,
		BundleID:    req.BundleID,
		Description: req.Description,
		IconURL:     req.IconURL,
		Category:    req.Category,
		Developer:   req.Developer,
		Platform:    req.Platform,
		Status:      1,
	}

	if err := ctrl.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": app})
}

// UpdateApp PUT /api/v1/market/apps/:id - 更新应用
func (ctrl *MarketAppController) UpdateApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的应用ID"})
		return
	}

	var app models.App
	if err := ctrl.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req MarketAppUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IconURL != "" {
		updates["icon_url"] = req.IconURL
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Developer != "" {
		updates["developer"] = req.Developer
	}
	if req.Platform != "" {
		updates["platform"] = req.Platform
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := ctrl.DB.Model(&app).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&app, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": app})
}

// DeleteApp DELETE /api/v1/market/apps/:id - 删除应用
func (ctrl *MarketAppController) DeleteApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的应用ID"})
		return
	}

	if err := ctrl.DB.Delete(&models.App{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ========== 请求结构体 ==========

type MarketAppCreateRequest struct {
	Name        string `json:"name"`
	BundleID    string `json:"bundle_id"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Category    string `json:"category"`
	Developer   string `json:"developer"`
	Platform    string `json:"platform"`
}

type MarketAppUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Category    string `json:"category"`
	Developer   string `json:"developer"`
	Platform    string `json:"platform"`
	Status      *int   `json:"status"`
}
