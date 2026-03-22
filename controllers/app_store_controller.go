package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppStoreController 企业应用商店控制器
type AppStoreController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// ==================== 请求结构体 ====================

// CreateStoreAppRequest 创建应用请求
type CreateStoreAppRequest struct {
	Name         string  `json:"name" binding:"required,min=1,max=128"`
	BundleID     string  `json:"bundle_id" binding:"required,max=256"`
	Description  string  `json:"description"`
	IconURL      string  `json:"icon_url"`
	Screenshots  string  `json:"screenshots"` // JSON数组
	Category     string  `json:"category"`
	Developer    string  `json:"developer"`
	Platform     string  `json:"platform"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	MinOSVersion string  `json:"min_os_version"`
}

// UpdateStoreAppRequest 更新应用请求
type UpdateStoreAppRequest struct {
	Name         string  `json:"name" binding:"omitempty,min=1,max=128"`
	Description  string  `json:"description"`
	IconURL      string  `json:"icon_url"`
	Screenshots  string  `json:"screenshots"`
	Category     string  `json:"category"`
	Developer    string  `json:"developer"`
	Platform     string  `json:"platform"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	MinOSVersion string  `json:"min_os_version"`
}

// CreateAppVersionRequest 创建版本请求
type CreateAppVersionRequest struct {
	Version      string `json:"version" binding:"required,max=32"`
	BuildNumber  string `json:"build_number"`
	FileSize     int64  `json:"file_size"`
	FileURL      string `json:"file_url" binding:"required"`
	FileMD5      string `json:"file_md5"`
	MinOSVersion string `json:"min_os_version"`
	ReleaseNotes string `json:"release_notes"`
	IsMandatory  bool   `json:"is_mandatory"`
}

// InstallAppRequest 安装应用请求
type InstallAppRequest struct {
	AppID        uint   `json:"app_id" binding:"required"`
	AppVersionID uint   `json:"app_version_id"`
	DeviceID     uint   `json:"device_id"`
	DeviceUUID   string `json:"device_uuid"`
}

// ReviewActionRequest 审核操作请求
type ReviewActionRequest struct {
	Reason string `json:"reason"`
}

// ==================== 应用商店 API ====================

// ListStoreApps 获取应用列表
// GET /api/v1/store/apps
func (c *AppStoreController) ListStoreApps(ctx *gin.Context) {
	var apps []models.StoreApp
	query := c.DB.Model(&models.StoreApp{})

	// 过滤条件
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if platform := ctx.Query("platform"); platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR bundle_id LIKE ? OR developer LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 只显示已发布的应用（前端浏览）
	if ctx.Query("published_only") == "true" {
		query = query.Where("is_published = ?", true)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  apps,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreateStoreApp 上架应用
// POST /api/v1/store/apps
func (c *AppStoreController) CreateStoreApp(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req CreateStoreAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 检查 BundleID 是否已存在
	var exist models.StoreApp
	if err := c.DB.Where("bundle_id = ?", req.BundleID).First(&exist).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "BundleID 已存在"})
		return
	}

	app := models.StoreApp{
		Name:         req.Name,
		BundleID:     req.BundleID,
		Description:  req.Description,
		IconURL:      req.IconURL,
		Screenshots:  req.Screenshots,
		Category:     req.Category,
		Developer:    req.Developer,
		DeveloperID:  userID,
		Platform:     req.Platform,
		Price:        req.Price,
		Currency:     req.Currency,
		MinOSVersion: req.MinOSVersion,
		Status:       models.StoreAppStatusPending,
		ReviewStatus: models.StoreAppReviewPending,
	}

	if err := c.DB.Create(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
		return
	}

	// 创建审核记录
	review := models.StoreReview{
		AppID:      app.ID,
		ReviewerID: 0,
		Status:     models.StoreReviewPending,
		Action:     "create",
	}
	c.DB.Create(&review)

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "创建成功，应用待审核", "data": app})
}

// GetStoreApp 获取应用详情
// GET /api/v1/store/apps/:id
func (c *AppStoreController) GetStoreApp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var app models.StoreApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 获取最新版本信息
	var latestVersion models.StoreAppVersion
	c.DB.Where("app_id = ? AND is_latest = ?", app.ID, true).First(&latestVersion)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"app":     app,
		"version": latestVersion,
	}})
}

// UpdateStoreApp 更新应用
// PUT /api/v1/store/apps/:id
func (c *AppStoreController) UpdateStoreApp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var app models.StoreApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	var req UpdateStoreAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IconURL != "" {
		updates["icon_url"] = req.IconURL
	}
	if req.Screenshots != "" {
		updates["screenshots"] = req.Screenshots
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
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.MinOSVersion != "" {
		updates["min_os_version"] = req.MinOSVersion
	}

	// 如果应用已发布，更新需要重新审核
	if app.IsPublished {
		updates["is_published"] = false
		updates["status"] = models.StoreAppStatusPending
		updates["review_status"] = models.StoreAppReviewPending
	}

	if err := c.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}

	c.DB.First(&app, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": app})
}

// DeleteStoreApp 下架应用
// DELETE /api/v1/store/apps/:id
func (c *AppStoreController) DeleteStoreApp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var app models.StoreApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 软删除（下架）
	if err := c.DB.Model(&app).Updates(map[string]interface{}{
		"status":       models.StoreAppStatusOffline,
		"is_published": false,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "下架失败", "error": err.Error()})
		return
	}

	// 下架所有版本
	c.DB.Model(&models.StoreAppVersion{}).Where("app_id = ?", id).Updates(map[string]interface{}{
		"status":    models.StoreAppVersionOffline,
		"is_latest": false,
	})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "下架成功"})
}

// PublishStoreApp 发布应用
// POST /api/v1/store/apps/:id/publish
func (c *AppStoreController) PublishStoreApp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var app models.StoreApp
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&app).Updates(map[string]interface{}{
		"is_published": true,
		"status":       models.StoreAppStatusPublished,
		"published_at": &now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布失败", "error": err.Error()})
		return
	}

	// 发布最新版本
	c.DB.Model(&models.StoreAppVersion{}).Where("app_id = ? AND is_latest = ?", id, true).
		Update("status", models.StoreAppVersionPublished)

	// 创建发布审核记录
	review := models.StoreReview{
		AppID:      app.ID,
		ReviewerID: getUserIDFromContext(ctx),
		Status:     models.StoreReviewApproved,
		Action:     "publish",
		Reason:     "应用发布",
		ReviewedAt: &now,
	}
	c.DB.Create(&review)

	c.DB.First(&app, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "发布成功", "data": app})
}

// ==================== 应用版本 API ====================

// CreateAppVersion 上传应用版本
// POST /api/v1/store/apps/:id/versions
func (c *AppStoreController) CreateAppVersion(ctx *gin.Context) {
	idStr := ctx.Param("id")
	appID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var app models.StoreApp
	if err := c.DB.First(&app, appID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	var req CreateAppVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 检查版本号是否已存在
	var exist models.StoreAppVersion
	if err := c.DB.Where("app_id = ? AND version = ?", appID, req.Version).First(&exist).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "版本号已存在"})
		return
	}

	version := models.StoreAppVersion{
		AppID:        uint(appID),
		Version:      req.Version,
		BuildNumber:  req.BuildNumber,
		FileSize:     req.FileSize,
		FileURL:      req.FileURL,
		FileMD5:      req.FileMD5,
		MinOSVersion: req.MinOSVersion,
		ReleaseNotes: req.ReleaseNotes,
		IsMandatory:  req.IsMandatory,
		IsActive:     true,
		Status:       models.StoreAppVersionPending,
	}

	// 如果是第一个版本，自动设为最新
	var count int64
	c.DB.Model(&models.StoreAppVersion{}).Where("app_id = ?", appID).Count(&count)
	if count == 0 {
		version.IsLatest = true
	}

	if err := c.DB.Create(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建版本失败", "error": err.Error()})
		return
	}

	// 更新应用版本号
	c.DB.Model(&app).Update("version", req.Version)

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "版本创建成功", "data": version})
}

// ListAppVersions 获取应用版本列表
// GET /api/v1/store/apps/:id/versions
func (c *AppStoreController) ListAppVersions(ctx *gin.Context) {
	idStr := ctx.Param("id")
	appID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var versions []models.StoreAppVersion
	query := c.DB.Where("app_id = ?", appID)

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("id DESC").Find(&versions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": versions})
}

// SetLatestVersion 设置最新版本
// POST /api/v1/store/apps/:id/versions/:version_id/set-latest
func (c *AppStoreController) SetLatestVersion(ctx *gin.Context) {
	appIDStr := ctx.Param("id")
	appID, _ := strconv.ParseUint(appIDStr, 10, 32)
	versionIDStr := ctx.Param("version_id")
	versionID, _ := strconv.ParseUint(versionIDStr, 10, 32)

	// 取消所有版本的最新标记
	c.DB.Model(&models.StoreAppVersion{}).Where("app_id = ?", appID).Update("is_latest", false)

	// 设置指定版本为最新
	if err := c.DB.Model(&models.StoreAppVersion{}).Where("id = ?", versionID).Update("is_latest", true).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "设置失败", "error": err.Error()})
		return
	}

	// 更新应用当前版本号
	var version models.StoreAppVersion
	if err := c.DB.First(&version, versionID).Error; err == nil {
		c.DB.Model(&models.StoreApp{}).Where("id = ?", appID).Update("version", version.Version)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置成功"})
}

// ==================== 应用安装管理 API ====================

// ListInstallations 获取安装列表
// GET /api/v1/store/installations
func (c *AppStoreController) ListInstallations(ctx *gin.Context) {
	var installations []models.StoreInstallation
	query := c.DB.Model(&models.StoreInstallation{})

	// 过滤条件
	if appID := ctx.Query("app_id"); appID != "" {
		query = query.Where("app_id = ?", appID)
	}
	if deviceUUID := ctx.Query("device_uuid"); deviceUUID != "" {
		query = query.Where("device_uuid = ?", deviceUUID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Preload("StoreApp").Order("id DESC").Offset(offset).Limit(pageSize).Find(&installations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  installations,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// InstallApp 安装应用
// POST /api/v1/store/installations
func (c *AppStoreController) InstallApp(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	tenantID := getTenantIDFromContext(ctx)

	var req InstallAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 检查应用是否存在
	var app models.StoreApp
	if err := c.DB.First(&app, req.AppID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "应用不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 检查是否已安装
	var existing models.StoreInstallation
	if err := c.DB.Where("app_id = ? AND device_uuid = ? AND status NOT IN (?, ?)",
		req.AppID, req.DeviceUUID, models.StoreUninstallPending, models.StoreUninstalled).
		First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "应用已安装"})
		return
	}

	// 获取最新版本
	var latestVersion models.StoreAppVersion
	if req.AppVersionID > 0 {
		c.DB.First(&latestVersion, req.AppVersionID)
	} else {
		c.DB.Where("app_id = ? AND is_latest = ?", req.AppID, true).First(&latestVersion)
	}

	installation := models.StoreInstallation{
		AppID:        req.AppID,
		AppVersionID: latestVersion.ID,
		DeviceID:     req.DeviceID,
		DeviceUUID:   req.DeviceUUID,
		UserID:       userID,
		TenantID:     tenantID,
		Status:       models.StoreInstallPending,
		Progress:     0,
	}

	if err := c.DB.Create(&installation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建安装记录失败", "error": err.Error()})
		return
	}

	// 更新应用安装次数
	c.DB.Model(&app).Update("install_count", gorm.Expr("install_count + 1"))

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "安装任务已创建", "data": installation})
}

// UninstallApp 卸载应用
// DELETE /api/v1/store/installations/:id
func (c *AppStoreController) UninstallApp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var installation models.StoreInstallation
	if err := c.DB.First(&installation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "安装记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&installation).Updates(map[string]interface{}{
		"status":         models.StoreUninstallPending,
		"uninstalled_at": &now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "卸载失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "卸载任务已创建"})
}

// GetInstallationStatus 获取安装状态
// GET /api/v1/store/installations/:id/status
func (c *AppStoreController) GetInstallationStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var installation models.StoreInstallation
	if err := c.DB.First(&installation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "安装记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 获取应用信息
	var app models.StoreApp
	c.DB.First(&app, installation.AppID)

	// 获取版本信息
	var version models.StoreAppVersion
	c.DB.First(&version, installation.AppVersionID)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"id":           installation.ID,
		"app_id":       installation.AppID,
		"app_name":     app.Name,
		"device_uuid":  installation.DeviceUUID,
		"status":       installation.Status,
		"progress":     installation.Progress,
		"error_msg":    installation.ErrorMsg,
		"version":      version.Version,
		"installed_at": installation.InstalledAt,
	}})
}

// UpdateInstallationStatus 更新安装状态（设备回调）
// PUT /api/v1/store/installations/:id/status
func (c *AppStoreController) UpdateInstallationStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var updates struct {
		Status   int    `json:"status"`
		Progress int    `json:"progress"`
		ErrorMsg string `json:"error_msg"`
	}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"status":    updates.Status,
		"progress":  updates.Progress,
		"error_msg": updates.ErrorMsg,
	}
	if updates.Status == models.StoreInstallSuccess {
		updateData["installed_at"] = &now
	}

	if err := c.DB.Model(&models.StoreInstallation{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "状态已更新"})
}

// ==================== 应用审核 API ====================

// ListReviews 获取审核列表
// GET /api/v1/store/reviews
func (c *AppStoreController) ListReviews(ctx *gin.Context) {
	var reviews []models.StoreReview
	query := c.DB.Model(&models.StoreReview{})

	// 过滤条件
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if appID := ctx.Query("app_id"); appID != "" {
		query = query.Where("app_id = ?", appID)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Preload("StoreApp").Order("id DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  reviews,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// ApproveReview 批准应用/版本
// POST /api/v1/store/reviews/:id/approve
func (c *AppStoreController) ApproveReview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var req ReviewActionRequest
	ctx.ShouldBindJSON(&req)

	userID := getUserIDFromContext(ctx)
	now := time.Now()

	// 更新审核记录
	if err := c.DB.Model(&models.StoreReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      models.StoreReviewApproved,
		"reviewer_id": userID,
		"reviewed_at": &now,
		"reason":      req.Reason,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审核失败", "error": err.Error()})
		return
	}

	// 获取审核记录
	var review models.StoreReview
	c.DB.First(&review, id)

	// 更新应用状态
	if review.VersionID > 0 {
		// 审核版本
		c.DB.Model(&models.StoreAppVersion{}).Where("id = ?", review.VersionID).
			Update("status", models.StoreAppVersionPublished)
	} else {
		// 审核应用
		c.DB.Model(&models.StoreApp{}).Where("id = ?", review.AppID).Updates(map[string]interface{}{
			"status":        models.StoreAppStatusPublished,
			"review_status": models.StoreAppReviewApproved,
			"reviewer_id":   userID,
			"reviewed_at":   &now,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "审核通过"})
}

// RejectReview 拒绝应用/版本
// POST /api/v1/store/reviews/:id/reject
func (c *AppStoreController) RejectReview(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var req ReviewActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请填写拒绝原因", "error": err.Error()})
		return
	}

	userID := getUserIDFromContext(ctx)
	now := time.Now()

	// 更新审核记录
	if err := c.DB.Model(&models.StoreReview{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      models.StoreReviewRejected,
		"reviewer_id": userID,
		"reviewed_at": &now,
		"reason":      req.Reason,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审核失败", "error": err.Error()})
		return
	}

	// 获取审核记录
	var review models.StoreReview
	c.DB.First(&review, id)

	// 更新应用状态
	if review.VersionID > 0 {
		// 拒绝版本
		c.DB.Model(&models.StoreAppVersion{}).Where("id = ?", review.VersionID).
			Update("status", models.StoreAppVersionOffline)
	} else {
		// 拒绝应用
		c.DB.Model(&models.StoreApp{}).Where("id = ?", review.AppID).Updates(map[string]interface{}{
			"status":        models.StoreAppStatusRejected,
			"review_status": models.StoreAppReviewRejected,
			"reviewer_id":   userID,
			"reviewed_at":   &now,
			"review_note":   req.Reason,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "已拒绝"})
}

// GetAppReviews 获取应用的审核历史
// GET /api/v1/store/apps/:id/reviews
func (c *AppStoreController) GetAppReviews(ctx *gin.Context) {
	idStr := ctx.Param("id")
	appID, _ := strconv.ParseUint(idStr, 10, 32)

	var reviews []models.StoreReview
	if err := c.DB.Where("app_id = ?", appID).Order("id DESC").Find(&reviews).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": reviews})
}

// getTenantIDFromContext Helper function to get tenant ID from context
func getTenantIDFromContext(ctx *gin.Context) uint {
	if v, exists := ctx.Get("tenant_id"); exists {
		if tenantID, ok := v.(uint); ok {
			return tenantID
		}
	}
	return 0
}
