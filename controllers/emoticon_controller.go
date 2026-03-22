package controllers

import (
	"net/http"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmoticonController 表情包控制器
type EmoticonController struct {
	DB *gorm.DB
}

// RegisterEmoticonRoutes 注册表情包路由
func (ctrl *EmoticonController) RegisterEmoticonRoutes(api *gin.RouterGroup) {
	api.GET("/emoticons", ctrl.List)
	api.POST("/emoticons", ctrl.Create)
	api.GET("/emoticons/categories", ctrl.GetCategories)
	api.GET("/emoticons/:id", ctrl.Get)
	api.DELETE("/emoticons/:id", ctrl.Delete)
	api.POST("/emoticons/:id/purchase", ctrl.Purchase)
}

// List 获取表情包列表（分页+筛选）
// GET /api/v1/emoticons
func (ctrl *EmoticonController) List(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	categoryID := c.Query("category_id")
	keyword := c.Query("keyword")
	isPremium := c.Query("is_premium")
	isFeatured := c.Query("is_featured")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	creatorType := c.Query("creator_type")

	query := ctrl.DB.Model(&models.Emoticon{}).Where("status = ?", "active")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if isPremium != "" {
		if isPremium == "true" {
			query = query.Where("is_premium = ?", true)
		} else if isPremium == "false" {
			query = query.Where("is_premium = ?", false)
		}
	}
	if isFeatured != "" {
		if isFeatured == "true" {
			query = query.Where("is_featured = ?", true)
		}
	}
	if creatorType != "" {
		query = query.Where("creator_type = ?", creatorType)
	}

	// 排序
	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	var total int64
	query.Count(&total)

	var emoticons []models.Emoticon
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&emoticons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query emoticons: " + err.Error(),
		})
		return
	}

	responses := make([]*models.EmoticonResponse, len(emoticons))
	for i := range emoticons {
		responses[i] = emoticons[i].ToResponse()
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

// GetCategories 获取表情包分类列表
// GET /api/v1/emoticons/categories
func (ctrl *EmoticonController) GetCategories(c *gin.Context) {
	isActive := c.DefaultQuery("is_active", "true")

	query := ctrl.DB.Model(&models.EmoticonCategory{})
	if isActive == "true" {
		query = query.Where("is_active = ?", true)
	}

	var categories []models.EmoticonCategory
	if err := query.Order("sort_order ASC, created_at DESC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 统计每个分类的表情包数量
	for i := range categories {
		var count int64
		ctrl.DB.Model(&models.Emoticon{}).Where("category_id = ? AND status = ?", categories[i].CategoryID, "active").Count(&count)
		categories[i].EmoticonCount = int(count)
	}

	responses := make([]*models.CategoryResponse, len(categories))
	for i := range categories {
		responses[i] = categories[i].ToCategoryResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// Get 获取表情包详情
// GET /api/v1/emoticons/:id
func (ctrl *EmoticonController) Get(c *gin.Context) {
	id := c.Param("id")

	var emoticon models.Emoticon
	if err := ctrl.DB.Where("id = ? OR emoticon_uuid = ?", id, id).First(&emoticon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Emoticon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	resp := emoticon.ToResponse()

	// 获取分类名称
	var category models.EmoticonCategory
	if err := ctrl.DB.Where("category_id = ?", emoticon.CategoryID).First(&category).Error; err == nil {
		resp.CategoryName = category.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": resp,
	})
}

// Create 创建表情包
// POST /api/v1/emoticons
func (ctrl *EmoticonController) Create(c *gin.Context) {
	var req struct {
		Name         string   `json:"name" binding:"required"`
		CategoryID   string   `json:"category_id" binding:"required"`
		ThumbnailURL string   `json:"thumbnail_url"`
		ImageURL     string   `json:"image_url"`
		GifURL       string   `json:"gif_url"`
		Tags         []string `json:"tags"`
		Description  string   `json:"description"`
		Price        float64  `json:"price"`
		Currency     string   `json:"currency"`
		IsPremium    bool     `json:"is_premium"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	// 获取用户ID（从上下文中，假设已认证）
	userID := getUserIDFromContext(c)
	creatorType := "user"
	if userID == 0 {
		creatorType = "system"
	}

	if req.Currency == "" {
		req.Currency = "CNY"
	}

	emoticon := models.Emoticon{
		Name:         req.Name,
		CategoryID:   req.CategoryID,
		CreatorID:    userID,
		CreatorType:  creatorType,
		ThumbnailURL: req.ThumbnailURL,
		ImageURL:     req.ImageURL,
		GifURL:       req.GifURL,
		Tags:         req.Tags,
		Description:  req.Description,
		Price:        req.Price,
		Currency:     req.Currency,
		IsPremium:    req.IsPremium,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := ctrl.DB.Create(&emoticon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create emoticon: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": emoticon.ToResponse(),
	})
}

// Delete 删除表情包
// DELETE /api/v1/emoticons/:id
func (ctrl *EmoticonController) Delete(c *gin.Context) {
	id := c.Param("id")

	var emoticon models.Emoticon
	if err := ctrl.DB.Where("id = ? OR emoticon_uuid = ?", id, id).First(&emoticon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Emoticon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限（只能删除自己的或系统表情包）
	userID := getUserIDFromContext(c)
	if emoticon.CreatorType != "system" && emoticon.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	if err := ctrl.DB.Delete(&emoticon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete emoticon: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Emoticon deleted successfully",
	})
}

// Purchase 购买表情包
// POST /api/v1/emoticons/:id/purchase
func (ctrl *EmoticonController) Purchase(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		PaymentMethod string  `json:"payment_method"` // coin/cash/free
		TransactionID string  `json:"transaction_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// 使用默认值
		req.PaymentMethod = "free"
	}

	var emoticon models.Emoticon
	if err := ctrl.DB.Where("id = ? OR emoticon_uuid = ?", id, id).First(&emoticon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Emoticon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查是否已购买
	userID := getUserIDFromContext(c)
	var existingPurchase models.EmoticonPurchase
	if err := ctrl.DB.Where("user_id = ? AND emoticon_id = ?", userID, emoticon.ID).First(&existingPurchase).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Already purchased",
			"data": gin.H{
				"purchase_uuid": existingPurchase.PurchaseUUID,
				"purchased_at":   existingPurchase.CreatedAt.Format(time.RFC3339),
			},
		})
		return
	}

	// 免费表情包直接记录购买
	price := emoticon.Price
	if price == 0 {
		req.PaymentMethod = "free"
	}

	purchase := models.EmoticonPurchase{
		UserID:        userID,
		EmoticonID:    emoticon.ID,
		EmoticonUUID:  emoticon.EmoticonUUID,
		Price:         price,
		Currency:      emoticon.Currency,
		PaymentMethod: req.PaymentMethod,
		TransactionID: req.TransactionID,
		CreatedAt:     time.Now(),
	}

	if err := ctrl.DB.Create(&purchase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to record purchase: " + err.Error()})
		return
	}

	// 增加下载计数
	ctrl.DB.Model(&emoticon).Update("download_count", emoticon.DownloadCount+1)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Purchase successful",
		"data": gin.H{
			"purchase_uuid": purchase.PurchaseUUID,
			"emoticon_uuid": emoticon.EmoticonUUID,
			"price":         price,
			"currency":      emoticon.Currency,
			"purchased_at":  purchase.CreatedAt.Format(time.RFC3339),
		},
	})
}

// CreateCategory 创建表情包分类（管理员）
// POST /api/v1/emoticons/categories
func (ctrl *EmoticonController) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		NameEn      string `json:"name_en"`
		IconURL     string `json:"icon_url"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	category := models.EmoticonCategory{
		Name:        req.Name,
		NameEn:      req.NameEn,
		IconURL:     req.IconURL,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := ctrl.DB.Create(&category).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "Category name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create category: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": category.ToCategoryResponse(),
	})
}
