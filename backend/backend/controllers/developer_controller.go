package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeveloperController 开发者API管理控制器
type DeveloperController struct {
	DB *gorm.DB
}

// ============ 请求结构 ============

type AppListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}

type AppCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	AppType     string `json:"app_type"`
	WebsiteURL  string `json:"website_url"`
	LogoURL     string `json:"logo_url"`
	QuotaTier   string `json:"quota_tier"`
}

type AppUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AppType     string `json:"app_type"`
	WebsiteURL  string `json:"website_url"`
	LogoURL     string `json:"logo_url"`
	QuotaTier   string `json:"quota_tier"`
	Status      string `json:"status"`
}

type CreateKeyRequest struct {
	Name      string   `json:"name" binding:"required"`
	Scopes    []string `json:"scopes"`
	ExpiresAt string   `json:"expires_at"` // RFC3339格式
}

// ============ 应用管理 ============

// ListApp 获取应用列表
func (c *DeveloperController) ListApp(ctx *gin.Context) {
	var req AppListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var list []models.DeveloperApp
	var total int64

	query := c.DB.Model(&models.DeveloperApp{})

	// 关键词过滤
	if req.Keyword != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&list)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": list,
			"pagination": gin.H{"total": total, "page": req.Page, "page_size": req.PageSize},
		},
	})
}

// GetApp 获取应用详情
func (c *DeveloperController) GetApp(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.DeveloperApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取关联的API Keys
	var keys []models.APIKey
	c.DB.Where("app_id = ?", app.ID).Order("created_at DESC").Find(&keys)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"app":  app,
			"keys": keys,
		},
	})
}

// CreateApp 创建应用
func (c *DeveloperController) CreateApp(ctx *gin.Context) {
	var req AppCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 获取当前用户
	ownerID := "system"
	ownerName := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		ownerID = uid.(string)
	}
	if uname, exists := ctx.Get("user_name"); exists {
		ownerName = uname.(string)
	}

	if req.AppType == "" {
		req.AppType = "standard"
	}
	if req.QuotaTier == "" {
		req.QuotaTier = "free"
	}

	app := models.DeveloperApp{
		Name:        req.Name,
		Description: req.Description,
		AppType:     req.AppType,
		WebsiteURL:  req.WebsiteURL,
		LogoURL:     req.LogoURL,
		QuotaTier:   req.QuotaTier,
		Status:      "active",
		OwnerID:     ownerID,
		OwnerName:   ownerName,
	}

	if err := c.DB.Create(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建应用失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": app})
}

// UpdateApp 更新应用
func (c *DeveloperController) UpdateApp(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.DeveloperApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req AppUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.AppType != "" {
		updates["app_type"] = req.AppType
	}
	if req.WebsiteURL != "" {
		updates["website_url"] = req.WebsiteURL
	}
	if req.LogoURL != "" {
		updates["logo_url"] = req.LogoURL
	}
	if req.QuotaTier != "" {
		updates["quota_tier"] = req.QuotaTier
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := c.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新应用失败"})
		return
	}

	c.DB.First(&app, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": app})
}

// DeleteApp 删除应用
func (c *DeveloperController) DeleteApp(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.DeveloperApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 删除关联的API Keys
	c.DB.Where("app_id = ?", app.ID).Delete(&models.APIKey{})

	if err := c.DB.Delete(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除应用失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ============ API Key管理 ============

// CreateKey 创建API Key
func (c *DeveloperController) CreateKey(ctx *gin.Context) {
	appID := ctx.Param("id")

	// 检查应用是否存在
	var app models.DeveloperApp
	if err := c.DB.First(&app, appID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req CreateKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 生成API Key
	rawKey, keyID := models.GenerateAPIKey()

	// Hash the key for storage
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	apiKey := models.APIKey{
		AppID:    app.ID,
		KeyID:    keyID,
		KeyHash:  keyHash,
		Name:     req.Name,
		Scopes:   req.Scopes,
		Status:   "active",
	}

	if req.ExpiresAt != "" {
		expiresAt, err := parseTime(req.ExpiresAt)
		if err == nil {
			apiKey.ExpiresAt = &expiresAt
		}
	}

	if err := c.DB.Create(&apiKey).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建Key失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"key_id":    keyID,
			"key":       rawKey, // 仅在此刻返回原始key，之后不再显示
			"name":      apiKey.Name,
			"scopes":    apiKey.Scopes,
			"expires_at": apiKey.ExpiresAt,
			"created_at": apiKey.CreatedAt,
		},
	})
}

// DeleteKey 删除API Key
func (c *DeveloperController) DeleteKey(ctx *gin.Context) {
	appID := ctx.Param("id")
	keyID := ctx.Param("keyId")

	var app models.DeveloperApp
	if err := c.DB.First(&app, appID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var key models.APIKey
	if err := c.DB.Where("app_id = ? AND key_id = ?", app.ID, keyID).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Key不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Delete(&key).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除Key失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
