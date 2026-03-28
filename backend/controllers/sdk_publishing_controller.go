package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SDKPublishingController SDK发布控制器
type SDKPublishingController struct {
	DB *gorm.DB
}

// NewSDKPublishingController 创建控制器
func NewSDKPublishingController(db *gorm.DB) *SDKPublishingController {
	return &SDKPublishingController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *SDKPublishingController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/sdks", ctrl.ListSDKs)
	rg.GET("/sdks/:id", ctrl.GetSDK)
	rg.POST("/sdks", ctrl.CreateSDK)
	rg.PUT("/sdks/:id", ctrl.UpdateSDK)
	rg.DELETE("/sdks/:id", ctrl.DeleteSDK)
	rg.GET("/sdks/:id/versions", ctrl.ListVersions)
	rg.POST("/sdks/:id/versions", ctrl.CreateVersion)
	rg.GET("/sdks/:id/versions/:version_id", ctrl.GetVersion)
	rg.POST("/sdks/:id/versions/:version_id/publish", ctrl.PublishVersion)
	rg.GET("/sdks/:id/download", ctrl.DownloadSDK)
	rg.POST("/sdks/search", ctrl.SearchSDKs)
}

// ListSDKs 获取SDK列表
func (ctrl *SDKPublishingController) ListSDKs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	platform := c.Query("platform")
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := ctrl.DB.Model(&models.SDKPackage{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", "published")
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	var list []models.SDKPackage
	query.Count(&total)

	query.Order("download_count DESC, star_count DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSDK 获取单个SDK
func (ctrl *SDKPublishingController) GetSDK(c *gin.Context) {
	id := c.Param("id")

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": sdk})
}

// CreateSDK 创建SDK
func (ctrl *SDKPublishingController) CreateSDK(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		DisplayName string   `json:"display_name"`
		Description string   `json:"description"`
		Category    string   `json:"category" binding:"required"`
		Platform   string   `json:"platform"`
		Language   string   `json:"language"`
		Tags       []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sdkID := fmt.Sprintf("SDK-%s-%d", req.Name, time.Now().Unix())

	tags := ""
	if len(req.Tags) > 0 {
		data, _ := json.Marshal(req.Tags)
		tags = string(data)
	}

	sdk := models.SDKPackage{
		SDKID:        sdkID,
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Category:     req.Category,
		Platform:    req.Platform,
		Language:    req.Language,
		Tags:        tags,
		Status:      "draft",
		DeveloperID: c.GetString("user_id"),
		CreatedBy:   c.GetString("username"),
	}

	if err := ctrl.DB.Create(&sdk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": sdk})
}

// UpdateSDK 更新SDK
func (ctrl *SDKPublishingController) UpdateSDK(c *gin.Context) {
	id := c.Param("id")

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	var req struct {
		DisplayName string   `json:"display_name"`
		Description string   `json:"description"`
		Tags       []string `json:"tags"`
		DocURL     string   `json:"doc_url"`
		GitHubURL  string   `json:"github_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DocURL != "" {
		updates["doc_url"] = req.DocURL
	}
	if req.GitHubURL != "" {
		updates["github_url"] = req.GitHubURL
	}
	if len(req.Tags) > 0 {
		data, _ := json.Marshal(req.Tags)
		updates["tags"] = string(data)
	}

	ctrl.DB.Model(&sdk).Updates(updates)
	ctrl.DB.Where("id = ?", sdk.ID).First(&sdk)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": sdk})
}

// DeleteSDK 删除SDK
func (ctrl *SDKPublishingController) DeleteSDK(c *gin.Context) {
	id := c.Param("id")

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	ctrl.DB.Model(&sdk).Updates(map[string]interface{}{"status": "archived"})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已归档"})
}

// ListVersions 获取SDK版本列表
func (ctrl *SDKPublishingController) ListVersions(c *gin.Context) {
	id := c.Param("id")

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	var versions []models.SDKVersion
	ctrl.DB.Where("sdk_id = ?", sdk.SDKID).Order("created_at DESC").Find(&versions)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": versions})
}

// CreateVersion 创建SDK版本
func (ctrl *SDKPublishingController) CreateVersion(c *gin.Context) {
	id := c.Param("id")

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	var req struct {
		Version     string `json:"version" binding:"required"`
		FileURL    string `json:"file_url"`
		FileSize   int64  `json:"file_size"`
		FileHash   string `json:"file_hash"`
		ReleaseNotes string `json:"release_notes"`
		MinPlatformVersion string `json:"min_platform_version"`
		Dependencies string `json:"dependencies"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 计算文件哈希
	if req.FileHash == "" && req.FileURL != "" {
		hash := sha256.Sum256([]byte(req.FileURL))
		req.FileHash = hex.EncodeToString(hash[:])
	}

	status := "stable"
	if req.Status != "" {
		status = req.Status
	}

	version := models.SDKVersion{
		SDKID:               sdk.SDKID,
		Version:            req.Version,
		FileURL:            req.FileURL,
		FileSize:           req.FileSize,
		FileHash:           req.FileHash,
		ReleaseNotes:       req.ReleaseNotes,
		MinPlatformVersion: req.MinPlatformVersion,
		Dependencies:       req.Dependencies,
		Status:             status,
		CreatedBy:          c.GetString("username"),
	}

	if err := ctrl.DB.Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	// 更新SDK的当前版本和版本数
	updates := map[string]interface{}{
		"current_version": req.Version,
		"version_count":   sdk.VersionCount + 1,
	}
	ctrl.DB.Model(&sdk).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "版本创建成功", "data": version})
}

// GetVersion 获取SDK版本详情
func (ctrl *SDKPublishingController) GetVersion(c *gin.Context) {
	id := c.Param("id")
	versionID := c.Param("version_id")

	var version models.SDKVersion
	query := ctrl.DB.Model(&models.SDKVersion{})

	if versionID != "" {
		if v, err := strconv.ParseUint(versionID, 10, 32); err == nil {
			query = query.Where("id = ?", v)
		} else {
			query = query.Where("version = ?", versionID)
		}
	} else {
		query = query.Where("sdk_id = ?", id)
	}

	if err := query.First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": version})
}

// PublishVersion 发布版本
func (ctrl *SDKPublishingController) PublishVersion(c *gin.Context) {
	id := c.Param("id")
	versionID := c.Param("version_id")

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	var version models.SDKVersion
	if err := ctrl.DB.Where("id = ? OR version = ?", versionID, versionID).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	// 更新版本状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":        "stable",
		"release_date":  now,
	}

	// 如果是推荐的，之前的取消推荐
	if version.IsRecommended {
		ctrl.DB.Model(&models.SDKVersion{}).
			Where("sdk_id = ? AND id != ?", sdk.SDKID, version.ID).
			Update("is_recommended", false)
	}

	ctrl.DB.Model(&version).Updates(updates)

	// 更新SDK状态为已发布
	if sdk.Status == "draft" {
		ctrl.DB.Model(&sdk).Updates(map[string]interface{}{
			"status":          "published",
			"current_version": version.Version,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发布成功", "data": version})
}

// DownloadSDK 下载SDK
func (ctrl *SDKPublishingController) DownloadSDK(c *gin.Context) {
	id := c.Param("id")
	versionID := c.Query("version") // 可选，指定版本

	var sdk models.SDKPackage
	if err := ctrl.DB.Where("id = ? OR sdk_id = ?", id, id).First(&sdk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "SDK不存在"})
		return
	}

	var version models.SDKVersion
	query := ctrl.DB.Where("sdk_id = ?", sdk.SDKID)
	if versionID != "" {
		query = query.Where("version = ?", versionID)
	} else {
		query = query.Where("is_recommended = ?", true)
	}

	if err := query.First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "无可用版本"})
		return
	}

	// 记录下载
	download := models.SDKDownload{
		SDKID:    sdk.SDKID,
		VersionID: version.ID,
		Version:  version.Version,
		UserID:   c.GetString("user_id"),
		UserName: c.GetString("username"),
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	ctrl.DB.Create(&download)

	// 更新下载计数
	ctrl.DB.Model(&sdk).Update("download_count", sdk.DownloadCount+1)
	ctrl.DB.Model(&version).Update("download_count", version.DownloadCount+1)

	// 重定向到文件URL
	if version.FileURL != "" {
		c.Redirect(http.StatusFound, version.FileURL)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"version": version.Version,
				"file_url": version.FileURL,
				"file_size": version.FileSize,
				"checksum": version.FileHash,
			},
		})
	}
}

// SearchSDKs 搜索SDK
func (ctrl *SDKPublishingController) SearchSDKs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	category := c.Query("category")
	platform := c.Query("platform")
	sort := c.DefaultQuery("sort", "relevance")

	query := ctrl.DB.Model(&models.SDKPackage{}).Where("status = ?", "published")

	if keyword != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	switch sort {
	case "downloads":
		query = query.Order("download_count DESC")
	case "stars":
		query = query.Order("star_count DESC")
	case "rating":
		query = query.Order("rating DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("download_count DESC")
	}

	var total int64
	var list []models.SDKPackage
	query.Count(&total)

	query.Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
