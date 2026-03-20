package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppController 应用管理控制器
type AppController struct {
	DB *gorm.DB
}

// ===== 请求结构 =====

// CreateAppRequest 创建应用请求
type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	BundleID    string `json:"bundle_id" binding:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Category    string `json:"category"`
	Developer   string `json:"developer"`
	Platform    string `json:"platform"`
}

// UpdateAppRequest 更新应用请求
type UpdateAppRequest struct {
	Name        string `json:"name"`
	BundleID    string `json:"bundle_id"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Category    string `json:"category"`
	Developer   string `json:"developer"`
	Platform    string `json:"platform"`
	Status      *int   `json:"status"`
}

// CreateAppVersionRequest 创建版本请求
type CreateAppVersionRequest struct {
	Version      string `json:"version" binding:"required"`
	BuildNumber  string `json:"build_number"`
	FileSize     int64  `json:"file_size"`
	FileURL      string `json:"file_url" binding:"required"`
	FileMD5      string `json:"file_md5"`
	MinOSVersion string `json:"min_os_version"`
	ReleaseNotes string `json:"release_notes"`
	IsMandatory  bool   `json:"is_mandatory"`
}

// CreateDistributionRequest 创建分发任务请求
type CreateDistributionRequest struct {
	Name             string   `json:"name" binding:"required"`
	AppID            uint     `json:"app_id" binding:"required"`
	VersionID        uint     `json:"version_id" binding:"required"`
	DistributionType string   `json:"distribution_type" binding:"required"` // device / user / group
	TargetIDs        []string `json:"target_ids" binding:"required"`
}

// ===== 应用 CRUD =====

// List 获取应用列表
func (c *AppController) List(ctx *gin.Context) {
	var apps []models.App
	query := c.DB.Model(&models.App{})

	// 关键词过滤
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name ILIKE ? OR bundle_id ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	// 平台过滤
	if platform := ctx.Query("platform"); platform != "" {
		query = query.Where("platform = ?", platform)
	}
	// 状态过滤
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

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": apps,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// Get 获取应用详情
func (c *AppController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.App
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "应用不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取最新版本
	var latestVersion models.AppVersion
	c.DB.Where("app_id = ? AND is_active = ?", id, true).Order("created_at DESC").First(&latestVersion)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"app":            app,
			"latest_version": latestVersion,
		},
	})
}

// Create 创建应用
func (c *AppController) Create(ctx *gin.Context) {
	var req CreateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error(), "error_code": "ERR_VALIDATION"})
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

	if err := c.DB.Create(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建应用失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    app,
	})
}

// Update 更新应用
func (c *AppController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.App
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "应用不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var req UpdateAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error(), "error_code": "ERR_VALIDATION"})
		return
	}

	// 选择性更新
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.BundleID != "" {
		updates["bundle_id"] = req.BundleID
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

	if err := c.DB.Model(&app).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新应用失败", "error_code": "ERR_INTERNAL"})
		return
	}

	c.DB.First(&app, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    app,
	})
}

// Delete 删除应用（软删除）
func (c *AppController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var app models.App
	if err := c.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "应用不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if err := c.DB.Delete(&app).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除应用失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ===== 版本管理 =====

// ListVersions 获取版本列表
func (c *AppController) ListVersions(ctx *gin.Context) {
	appID := ctx.Param("id")
	var versions []models.AppVersion
	query := c.DB.Model(&models.AppVersion{}).Where("app_id = ?", appID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&versions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": versions,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// CreateVersion 添加版本
func (c *AppController) CreateVersion(ctx *gin.Context) {
	appID := ctx.Param("id")

	// 检查应用是否存在
	var app models.App
	if err := c.DB.First(&app, appID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "应用不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var req CreateAppVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error(), "error_code": "ERR_VALIDATION"})
		return
	}

	version := models.AppVersion{
		AppID:        app.ID,
		Version:      req.Version,
		BuildNumber:  req.BuildNumber,
		FileSize:     req.FileSize,
		FileURL:      req.FileURL,
		FileMD5:      req.FileMD5,
		MinOSVersion: req.MinOSVersion,
		ReleaseNotes: req.ReleaseNotes,
		IsMandatory:  req.IsMandatory,
		IsActive:     true,
	}

	if err := c.DB.Create(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建版本失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    version,
	})
}

// DeleteVersion 删除版本
func (c *AppController) DeleteVersion(ctx *gin.Context) {
	versionID := ctx.Param("version_id")
	var version models.AppVersion
	if err := c.DB.First(&version, versionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "版本不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if err := c.DB.Delete(&version).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除版本失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ===== 分发任务 =====

// CreateDistribution 创建分发任务
func (c *AppController) CreateDistribution(ctx *gin.Context) {
	var req CreateDistributionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error(), "error_code": "ERR_VALIDATION"})
		return
	}

	// 校验应用和版本
	var app models.App
	if err := c.DB.First(&app, req.AppID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "应用不存在", "error_code": "ERR_NOT_FOUND"})
		return
	}
	var version models.AppVersion
	if err := c.DB.First(&version, req.VersionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "版本不存在", "error_code": "ERR_NOT_FOUND"})
		return
	}
	if version.AppID != req.AppID {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "版本不属于该应用", "error_code": "ERR_VALIDATION"})
		return
	}

	// 校验 distribution_type
	if req.DistributionType != "device" && req.DistributionType != "user" && req.DistributionType != "group" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "distribution_type 必须为 device/user/group 之一", "error_code": "ERR_VALIDATION"})
		return
	}

	// 序列化 target_ids
	targetIDsJSON, _ := json.Marshal(req.TargetIDs)

	// 获取当前用户（如果有）
	createdBy := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		createdBy = uid.(string)
	}

	dist := models.AppDistribution{
		Name:             req.Name,
		AppID:            req.AppID,
		VersionID:        req.VersionID,
		DistributionType: req.DistributionType,
		TargetIDs:        string(targetIDsJSON),
		TargetCount:      len(req.TargetIDs),
		PendingCount:     len(req.TargetIDs),
		Status:           "pending",
		CreatedBy:        createdBy,
	}

	if err := c.DB.Create(&dist).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建分发任务失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 模拟初始化安装记录（实际场景由设备回调更新状态）
	for _, targetID := range req.TargetIDs {
		record := models.AppInstallRecord{
			DistributionID: &dist.ID,
			DeviceID:       targetID,
			AppID:          req.AppID,
			VersionID:      req.VersionID,
			Status:         "pending",
		}
		c.DB.Create(&record)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    dist,
	})
}

// GetDistribution 获取分发详情
func (c *AppController) GetDistribution(ctx *gin.Context) {
	id := ctx.Param("id")
	var dist models.AppDistribution
	if err := c.DB.First(&dist, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "分发任务不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取关联应用和版本信息
	var app models.App
	var version models.AppVersion
	c.DB.First(&app, dist.AppID)
	c.DB.First(&version, dist.VersionID)

	// 获取安装记录
	var records []models.AppInstallRecord
	c.DB.Where("distribution_id = ?", dist.ID).Find(&records)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"distribution":     dist,
			"app":              app,
			"version":          version,
			"install_records":  records,
		},
	})
}

// CancelDistribution 取消分发任务
func (c *AppController) CancelDistribution(ctx *gin.Context) {
	id := ctx.Param("id")
	var dist models.AppDistribution
	if err := c.DB.First(&dist, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "分发任务不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if dist.Status == "cancelled" || dist.Status == "completed" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "该任务无法取消", "error_code": "ERR_VALIDATION"})
		return
	}

	cancelledBy := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		cancelledBy = uid.(string)
	}
	now := time.Now()

	if err := c.DB.Model(&dist).Updates(map[string]interface{}{
		"status":       "cancelled",
		"cancelled_by": cancelledBy,
		"cancelled_at": now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "取消分发失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 更新相关安装记录状态
	c.DB.Model(&models.AppInstallRecord{}).Where("distribution_id = ? AND status = ?", dist.ID, "pending").
		Update("status", "cancelled")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ===== 统计 =====

// GetStats 获取应用安装统计
func (c *AppController) GetStats(ctx *gin.Context) {
	appID := ctx.Param("id")

	var app models.App
	if err := c.DB.First(&app, appID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "应用不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var totalInstalls int64
	var activeInstalls int64
	var failedInstalls int64

	c.DB.Model(&models.AppInstallRecord{}).Where("app_id = ?", appID).Count(&totalInstalls)
	c.DB.Model(&models.AppInstallRecord{}).Where("app_id = ? AND status = ?", appID, "installed").Count(&activeInstalls)
	c.DB.Model(&models.AppInstallRecord{}).Where("app_id = ? AND status = ?", appID, "failed").Count(&failedInstalls)

	// 按版本统计
	type VersionStat struct {
		VersionID   uint   `json:"version_id"`
		Version     string `json:"version"`
		TotalCount  int64  `json:"total_count"`
		SuccessCount int64 `json:"success_count"`
	}
	var versionStats []VersionStat
	c.DB.Table("app_install_records as r").
		Select("r.version_id, v.version, COUNT(*) as total_count, SUM(CASE WHEN r.status = 'installed' THEN 1 ELSE 0 END) as success_count").
		Joins("LEFT JOIN app_versions v ON v.id = r.version_id").
		Where("r.app_id = ?", appID).
		Group("r.version_id, v.version").
		Scan(&versionStats)

	// 按时间统计（最近30天）
	type DailyStat struct {
		Date         string `json:"date"`
		InstallCount int64  `json:"install_count"`
		FailCount    int64  `json:"fail_count"`
	}
	var dailyStats []DailyStat
	c.DB.Table("app_install_records").
		Select("DATE(created_at) as date, COUNT(*) as install_count, SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as fail_count").
		Where("app_id = ? AND created_at >= ?", appID, time.Now().AddDate(0, 0, -30)).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyStats)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"app_id":          appID,
			"total_installs":  totalInstalls,
			"active_installs": activeInstalls,
			"failed_installs": failedInstalls,
			"version_stats":   versionStats,
			"daily_stats":     dailyStats,
		},
	})
}
