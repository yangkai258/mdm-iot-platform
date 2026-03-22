package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeveloperController developer platform API controller
type DeveloperController struct {
	DB *gorm.DB
}

// CreateAppRequest create app request
type CreateAppRequest struct {
	AppName     string `json:"app_name" binding:"required,min=2,max=128"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
	LogoURL     string `json:"logo_url"`
	Category    string `json:"category"`
	Platform    string `json:"platform"`
}

// UpdateAppRequest update app request
type UpdateAppRequest struct {
	AppName     string `json:"app_name" binding:"required,min=2,max=128"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
	LogoURL     string `json:"logo_url"`
	Category    string `json:"category"`
	Platform    string `json:"platform"`
	Status      *int   `json:"status"`
}

// CreateKeyRequest create API key request
type CreateKeyRequest struct {
	Name      string   `json:"name" binding:"required,max=128"`
	Scopes    []string `json:"scopes"`
	ExpiresIn int      `json:"expires_in"` // 0 = never expires, unit: days
}

// CreateApp creates an app
// POST /api/v1/developer/apps
func (c *DeveloperController) CreateApp(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	var req CreateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request", "error": err.Error()})
		return
	}

	appKey := fmt.Sprintf("mdm_app_%s", uuid.New().String()[:8])

	app := models.DeveloperApp{
		UserID:      userID,
		AppName:     req.AppName,
		AppKey:      appKey,
		Description: req.Description,
		WebsiteURL:  req.WebsiteURL,
		LogoURL:     req.LogoURL,
		Category:    req.Category,
		Platform:    req.Platform,
		Status:      models.DeveloperAppStatusEnabled,
	}

	if err := c.DB.Create(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "created",
		"data":    app,
	})
}

// ListApps lists apps
// GET /api/v1/developer/apps
func (c *DeveloperController) ListApps(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	offset := (page - 1) * pageSize

	var total int64
	var apps []models.DeveloperApp

	query := c.DB.Model(&models.DeveloperApp{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": apps,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetApp gets an app
// GET /api/v1/developer/apps/:id
func (c *DeveloperController) GetApp(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "app not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    app,
	})
}

// UpdateApp updates an app
// PUT /api/v1/developer/apps/:id
func (c *DeveloperController) UpdateApp(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", id, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "app not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	var req UpdateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"app_name":    req.AppName,
		"description": req.Description,
		"website_url": req.WebsiteURL,
		"logo_url":    req.LogoURL,
		"category":    req.Category,
		"platform":    req.Platform,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := c.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed", "error": err.Error()})
		return
	}

	c.DB.First(&app, id)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "updated",
		"data":    app,
	})
}

// DeleteApp deletes an app
// DELETE /api/v1/developer/apps/:id
func (c *DeveloperController) DeleteApp(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid app id"})
		return
	}

	result := c.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.DeveloperApp{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed", "error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "app not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "deleted",
	})
}

// CreateKey creates an API key
// POST /api/v1/developer/apps/:id/keys
func (c *DeveloperController) CreateKey(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	appIDStr := ctx.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "app not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	var req CreateKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request", "error": err.Error()})
		return
	}

	secret, hash, err := models.GenerateKeyPair()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "key generation failed", "error": err.Error()})
		return
	}

	keyID := uuid.New().String()[:8]
	keyPrefix := secret[:16]

	apiKey := models.APIKey{
		UserID:    userID,
		AppID:    uint(appID),
		KeyID:    keyID,
		KeyPrefix: keyPrefix,
		KeyHash:  hash,
		Name:     req.Name,
		Scopes:   models.JSON{},
		Status:   models.APIKeyStatusEnabled,
	}

	if len(req.Scopes) > 0 {
		for _, s := range req.Scopes {
			apiKey.Scopes[s] = true
		}
	}

	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * 24 * time.Hour)
		apiKey.ExpiresAt = &expiresAt
	}

	if err := c.DB.Create(&apiKey).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create key failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "created, please save the key securely",
		"data": gin.H{
			"key_id":     apiKey.KeyID,
			"api_key":    secret,
			"key_prefix": keyPrefix,
			"name":       apiKey.Name,
			"scopes":     apiKey.Scopes,
			"expires_at":  apiKey.ExpiresAt,
			"created_at":  apiKey.CreatedAt,
		},
	})
}

// ListKeys lists API keys
// GET /api/v1/developer/apps/:id/keys
func (c *DeveloperController) ListKeys(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	appIDStr := ctx.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid app id"})
		return
	}

	var app models.DeveloperApp
	if err := c.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "app not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	offset := (page - 1) * pageSize

	var total int64
	var keys []models.APIKey

	query := c.DB.Model(&models.APIKey{}).Where("app_id = ? AND user_id = ?", appID, userID)

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&keys).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()})
		return
	}

	result := make([]gin.H, len(keys))
	for i, k := range keys {
		result[i] = gin.H{
			"id":           k.ID,
			"key_id":       k.KeyID,
			"key_prefix":   k.KeyPrefix + "***",
			"name":         k.Name,
			"scopes":       k.Scopes,
			"rate_limit":   k.RateLimit,
			"status":       k.Status,
			"last_used_at": k.LastUsedAt,
			"expires_at":   k.ExpiresAt,
			"created_at":   k.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": result,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// DeleteKey deletes an API key
// DELETE /api/v1/developer/apps/:app_id/keys/:key_id
func (c *DeveloperController) DeleteKey(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	appIDStr := ctx.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid app id"})
		return
	}

	keyIDStr := ctx.Param("key_id")

	result := c.DB.Where("id = ? AND app_id = ? AND user_id = ?", keyIDStr, appID, userID).Delete(&models.APIKey{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed", "error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "key not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "deleted",
	})
}

// GetStats gets API usage statistics
// GET /api/v1/developer/stats
func (c *DeveloperController) GetStats(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	days := 30
	if d := ctx.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	startDate := time.Now().AddDate(0, 0, -days)

	var apps []models.DeveloperApp
	c.DB.Where("user_id = ?", userID).Find(&apps)

	appIDs := make([]uint, len(apps))
	for i, a := range apps {
		appIDs[i] = a.ID
	}

	type statItem struct {
		Date       string `json:"date"`
		CallCount  int64  `json:"call_count"`
		ErrorCount int64  `json:"error_count"`
	}

	var dailyStats []statItem

	if len(appIDs) > 0 {
		c.DB.Model(&models.UsageRecord{}).
			Select("DATE(record_date) as date, SUM(usage_value) as call_count, 0 as error_count").
			Where("user_id = ? AND usage_type = ? AND record_date >= ?", userID, models.UsageTypeAPICall, startDate).
			Group("DATE(record_date)").
			Order("date DESC").
			Scan(&dailyStats)
	}

	var totalCalls float64
	c.DB.Model(&models.UsageRecord{}).
		Select("COALESCE(SUM(usage_value), 0)").
		Where("user_id = ? AND usage_type = ? AND record_date >= ?", userID, models.UsageTypeAPICall, startDate).
		Scan(&totalCalls)

	var totalApps int64
	c.DB.Model(&models.DeveloperApp{}).Where("user_id = ?", userID).Count(&totalApps)

	var activeKeys int64
	if len(appIDs) > 0 {
		c.DB.Model(&models.APIKey{}).Where("app_id IN ? AND user_id = ? AND status = ?", appIDs, userID, models.APIKeyStatusEnabled).Count(&activeKeys)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_calls":  int64(totalCalls),
			"total_apps":   totalApps,
			"active_keys": activeKeys,
			"period_days": days,
			"daily_stats": dailyStats,
		},
	})
}

// GetQuota gets quota information
// GET /api/v1/developer/quota
func (c *DeveloperController) GetQuota(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}

	var quotas []models.UserQuota
	c.DB.Where("user_id = ?", userID).Find(&quotas)

	if len(quotas) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"is_default": true,
				"quotas": []gin.H{
					{"type": models.UsageTypeAPICall, "limit": 10000.0, "used": 0.0, "unit": "calls", "period": "monthly"},
					{"type": models.UsageTypeDevice, "limit": 10.0, "used": 0.0, "unit": "devices", "period": "monthly"},
					{"type": models.UsageTypeStorage, "limit": 1.0, "used": 0.0, "unit": "GB", "period": "monthly"},
				},
			},
		})
		return
	}

	quotaList := make([]gin.H, len(quotas))
	for i, q := range quotas {
		quotaList[i] = gin.H{
			"type":   q.QuotaType,
			"limit":  q.QuotaLimit,
			"used":   q.QuotaUsed,
			"unit":   q.Unit,
			"period": q.PeriodType,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"is_default": false,
			"quotas":     quotaList,
		},
	})
}
