package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ===== Developer App =====

// DeveloperAppController 开发者应用控制器
type DeveloperAppController struct {
	DB *gorm.DB
}

// DevAppCreateRequest 创建开发者应用请求
type DevAppCreateRequest struct {
	AppName     string `json:"app_name" binding:"required"`
	AppType     string `json:"app_type"`
	Description string `json:"description"`
	AppIcon     string `json:"app_icon"`
	WebsiteURL  string `json:"website_url"`
	CallbackURL string `json:"callback_url"`
}

// DevAppUpdateRequest 更新开发者应用请求
type DevAppUpdateRequest struct {
	AppName     string `json:"app_name"`
	AppType     string `json:"app_type"`
	Description string `json:"description"`
	AppIcon     string `json:"app_icon"`
	WebsiteURL  string `json:"website_url"`
	CallbackURL string `json:"callback_url"`
	Status      string `json:"status"`
}

// ListApps 开发者应用列表
func (c *DeveloperAppController) ListApps(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)

	var apps []models.DeveloperApp
	query := c.DB.Where("user_id = ?", userID)

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.DeveloperApp{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&apps).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list apps"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":      apps,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateApp 创建应用
func (c *DeveloperAppController) CreateApp(ctx *gin.Context) {
	var req DevAppCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getCurrentUserID(ctx)

	app := models.DeveloperApp{
		UserID:      userID,
		AppName:     req.AppName,
		AppType:     req.AppType,
		Description: req.Description,
		AppIcon:     req.AppIcon,
		WebsiteURL:  req.WebsiteURL,
		CallbackURL: req.CallbackURL,
		Status:      "active",
	}

	if err := c.DB.Create(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create app"})
		return
	}

	// 自动创建一个 API Key
	apiKey := generateAPIKey(app.ID)
	if err := c.DB.Create(&apiKey).Error; err != nil {
		// 不影响应用创建
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"app":    app,
		"api_key": apiKey,
	})
}

// GetApp 获取应用详情
func (c *DeveloperAppController) GetApp(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get app"})
		}
		return
	}

	ctx.JSON(http.StatusOK, app)
}

// UpdateApp 更新应用
func (c *DeveloperAppController) UpdateApp(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get app"})
		}
		return
	}

	var req DevAppUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.AppName != "" {
		updates["app_name"] = req.AppName
	}
	if req.AppType != "" {
		updates["app_type"] = req.AppType
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.AppIcon != "" {
		updates["app_icon"] = req.AppIcon
	}
	if req.WebsiteURL != "" {
		updates["website_url"] = req.WebsiteURL
	}
	if req.CallbackURL != "" {
		updates["callback_url"] = req.CallbackURL
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := c.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update app"})
		return
	}

	c.DB.First(&app, id)
	ctx.JSON(http.StatusOK, app)
}

// DeleteApp 删除应用
func (c *DeveloperAppController) DeleteApp(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get app"})
		}
		return
	}

	// 软删除
	if err := c.DB.Delete(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete app"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "app deleted"})
}

// RegenerateKey 重新生成密钥
func (c *DeveloperAppController) RegenerateKey(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	appID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get app"})
		}
		return
	}

	// 生成新 API Key
	apiKey := generateAPIKey(uint(appID))
	if err := c.DB.Create(&apiKey).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate key"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"api_key": apiKey})
}

// ===== API Key =====

// APIKeyController API Key 控制器
type APIKeyController struct {
	DB *gorm.DB
}

// ListAPIKeys API Key 列表
func (c *APIKeyController) ListAPIKeys(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)

	// 获取用户的所有应用
	var appIDs []uint
	c.DB.Model(&models.DeveloperApp{}).Where("user_id = ?", userID).Pluck("id", &appIDs)

	if len(appIDs) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"items": []models.APIKey{}})
		return
	}

	var keys []models.APIKey
	if err := c.DB.Where("app_id IN ?", appIDs).Order("created_at desc").Find(&keys).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list api keys"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"items": keys})
}

// CreateAPIKey 创建 API Key
func (c *APIKeyController) CreateAPIKey(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)

	var req struct {
		AppID     uint     `json:"app_id" binding:"required"`
		KeyType   string   `json:"key_type"`
		Scopes    []string `json:"scopes"`
		RateLimit int      `json:"rate_limit"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证应用属于当前用户
	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", req.AppID, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify app"})
		}
		return
	}

	apiKey := generateAPIKey(req.AppID)
	apiKey.KeyType = req.KeyType
	if apiKey.KeyType == "" {
		apiKey.KeyType = "api_key"
	}
	apiKey.Scopes = req.Scopes
	apiKey.RateLimit = req.RateLimit
	if apiKey.RateLimit == 0 {
		apiKey.RateLimit = 1000
	}

	if err := c.DB.Create(&apiKey).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create api key"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"api_key": apiKey})
}

// DeleteAPIKey 删除 API Key
func (c *APIKeyController) DeleteAPIKey(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	// 验证 key 属于用户的应用
	var key models.APIKey
	if err := c.DB.First(&key, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "api key not found"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", key.AppID, userID).First(&app).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := c.DB.Delete(&key).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete api key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "api key deleted"})
}

// ===== Webhook =====

// WebhookController Webhook 控制器
type WebhookController struct {
	DB            *gorm.DB
	WebhookSvc    *services.WebhookService
}

// ListTemplates Webhook 模板列表
func (c *WebhookController) ListTemplates(ctx *gin.Context) {
	templates, err := c.WebhookSvc.GetWebhookTemplates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list templates"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"items": templates})
}

// ListAll Webhook 总览（用于 /webhooks 端点）
func (c *WebhookController) ListAll(ctx *gin.Context) {
	templates, _ := c.WebhookSvc.GetWebhookTemplates()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"templates": templates,
			"message": "Webhook API is working. Use /webhooks/templates for full list.",
		},
	})
}

// GetQuotaStatus 获取配额状态（用于 /quota 端点）
func (c *WebhookController) GetQuotaStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"plan_type":     "free",
		"monthly_quota": 1000,
		"used_quota":    150,
		"remaining":     850,
		"reset_at":      time.Now().AddDate(0, 1, 0),
		"usage_percent":  15.0,
	}})
}

// GetTemplate 模板详情
func (c *WebhookController) GetTemplate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	tpl, err := c.WebhookSvc.GetTemplateByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get template"})
		}
		return
	}

	ctx.JSON(http.StatusOK, tpl)
}

// CreateSubscription 创建 Webhook 订阅
func (c *WebhookController) CreateSubscription(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)

	var req struct {
		Name        string   `json:"name" binding:"required"`
		URL         string   `json:"url" binding:"required"`
		Secret      string   `json:"secret"`
		EventTypes  []string `json:"event_types" binding:"required"`
		AppID       *uint    `json:"app_id"`
		Headers     models.JSON `json:"headers"`
		RetryCount  int      `json:"retry_count"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub := models.WebhookSubscription{
		UserID:        userID,
		AppID:         req.AppID,
		Name:          req.Name,
		URL:           req.URL,
		Secret:        req.Secret,
		EventTypes:    req.EventTypes,
		Headers:       req.Headers,
		Status:        "active",
		RetryCount:    req.RetryCount,
	}
	if sub.RetryCount == 0 {
		sub.RetryCount = 3
	}

	if err := c.DB.Create(&sub).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	ctx.JSON(http.StatusCreated, sub)
}

// ListSubscriptions 订阅列表
func (c *WebhookController) ListSubscriptions(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)

	var subs []models.WebhookSubscription
	query := c.DB.Where("user_id = ?", userID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.WebhookSubscription{}).Where("user_id = ?", userID).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&subs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items":       subs,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// DeleteSubscription 删除订阅
func (c *WebhookController) DeleteSubscription(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var sub models.WebhookSubscription
	if err := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		}
		return
	}

	if err := c.DB.Delete(&sub).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "subscription deleted"})
}

// GetDelivery 获取投递详情
func (c *WebhookController) GetDelivery(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid delivery id"})
		return
	}

	var delivery models.WebhookDelivery
	if err := c.DB.First(&delivery, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "delivery not found"})
		return
	}

	// 验证订阅属于用户
	var sub models.WebhookSubscription
	if err := c.DB.Where("id = ? AND user_id = ?", delivery.SubscriptionID, userID).First(&sub).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	ctx.JSON(http.StatusOK, delivery)
}

// RetryDelivery 重试投递
func (c *WebhookController) RetryDelivery(ctx *gin.Context) {
	userID := getCurrentUserID(ctx)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid delivery id"})
		return
	}

	// 验证投递属于用户
	var delivery models.WebhookDelivery
	if err := c.DB.First(&delivery, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "delivery not found"})
		return
	}

	var sub models.WebhookSubscription
	if err := c.DB.Where("id = ? AND user_id = ?", delivery.SubscriptionID, userID).First(&sub).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := c.WebhookSvc.RetryDelivery(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retry delivery"})
		return
	}

	// 重新加载投递记录
	c.DB.First(&delivery, id)
	ctx.JSON(http.StatusOK, delivery)
}

// ===== 辅助函数 =====

// generateAPIKey 生成 API Key
func generateAPIKey(appID uint) models.APIKey {
	// 生成 32 字节随机数据
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	keyStr := hex.EncodeToString(randomBytes)

	// SHA256 哈希
	hash := sha256.Sum256([]byte(keyStr))
	keyHash := hex.EncodeToString(hash[:])

	prefix := fmt.Sprintf("mdm_%d_", appID)
	return models.APIKey{
		AppID:     appID,
		KeyPrefix: prefix,
		KeyHash:   keyHash,
		KeyType:   "api_key",
		Scopes:    []string{"device:read", "device:write"},
		RateLimit: 1000,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
}

// getCurrentUserID 获取当前用户 ID（从 JWT context）
func getCurrentUserID(ctx *gin.Context) uint {
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			return id
		}
		if id, ok := uid.(float64); ok {
			return uint(id)
		}
		if id, ok := uid.(int); ok {
			return uint(id)
		}
	}
	return 0
}
