package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResearchPlatformController 研究平台控制器
type ResearchPlatformController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册研究平台路由
func (ctrl *ResearchPlatformController) RegisterRoutes(rg *gin.RouterGroup) {
	research := rg.Group("/research")
	{
		// 研究平台管理
		research.GET("/platforms", ctrl.ListPlatforms)
		research.POST("/platforms", ctrl.CreatePlatform)
		research.GET("/platforms/:id", ctrl.GetPlatform)
		research.PUT("/platforms/:id", ctrl.UpdatePlatform)
		research.DELETE("/platforms/:id", ctrl.DeletePlatform)
	}
}

// ListPlatforms 获取研究平台列表
// GET /api/v1/research/platforms
func (ctrl *ResearchPlatformController) ListPlatforms(c *gin.Context) {
	platformType := c.Query("type")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var platforms []models.ResearchPlatform
	query := ctrl.DB.Model(&models.ResearchPlatform{})

	if platformType != "" {
		query = query.Where("platform_type = ?", platformType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&platforms)

	c.JSON(http.StatusOK, gin.H{
		"platforms": platforms,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreatePlatform 创建研究平台
// POST /api/v1/research/platforms
func (ctrl *ResearchPlatformController) CreatePlatform(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		PlatformType  string `json:"platform_type" binding:"required"`
		WebsiteURL    string `json:"website_url"`
		Documentation string `json:"documentation"`
		PricingInfo   string `json:"pricing_info"`
		OwnerID       uint   `json:"owner_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	platform := models.ResearchPlatform{
		Name:          req.Name,
		Description:   req.Description,
		PlatformType:  req.PlatformType,
		WebsiteURL:    req.WebsiteURL,
		Documentation: req.Documentation,
		PricingInfo:   req.PricingInfo,
		OwnerID:       req.OwnerID,
		Status:        "active",
	}

	ctrl.DB.Create(&platform)
	c.JSON(http.StatusOK, platform)
}

// GetPlatform 获取研究平台详情
// GET /api/v1/research/platforms/:id
func (ctrl *ResearchPlatformController) GetPlatform(c *gin.Context) {
	id := c.Param("id")
	var platform models.ResearchPlatform

	if err := ctrl.DB.First(&platform, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Platform not found"})
		return
	}

	c.JSON(http.StatusOK, platform)
}

// UpdatePlatform 更新研究平台
// PUT /api/v1/research/platforms/:id
func (ctrl *ResearchPlatformController) UpdatePlatform(c *gin.Context) {
	id := c.Param("id")
	var platform models.ResearchPlatform

	if err := ctrl.DB.First(&platform, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Platform not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&platform).Updates(req)
	c.JSON(http.StatusOK, platform)
}

// DeletePlatform 删除研究平台
// DELETE /api/v1/research/platforms/:id
func (ctrl *ResearchPlatformController) DeletePlatform(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.ResearchPlatform{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
