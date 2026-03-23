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

type ContentController struct {
	DB *gorm.DB
}

func (ctrl *ContentController) RegisterRoutes(rg *gin.RouterGroup) {
	content := rg.Group("/content")
	{
		content.GET("/files", ctrl.GetFiles)
		content.POST("/files", ctrl.CreateFile)
		content.GET("/files/:id", ctrl.GetFile)
		content.PUT("/files/:id", ctrl.UpdateFile)
		content.DELETE("/files/:id", ctrl.DeleteFile)
		content.GET("/files/:id/download", ctrl.DownloadFile)
		content.POST("/files/:id/distribute", ctrl.DistributeFile)
		content.GET("/distributions", ctrl.GetDistributions)
	}

	apps := rg.Group("/apps")
	{
		apps.GET("/packages", ctrl.GetPackages)
		apps.POST("/packages", ctrl.CreatePackage)
		apps.GET("/packages/:id", ctrl.GetPackage)
		apps.PUT("/packages/:id", ctrl.UpdatePackage)
		apps.DELETE("/packages/:id", ctrl.DeletePackage)
		apps.POST("/packages/:id/submit-review", ctrl.SubmitReview)
		apps.POST("/packages/:id/approve", ctrl.ApprovePackage)
		apps.POST("/packages/:id/reject", ctrl.RejectPackage)

		apps.GET("/installs", ctrl.GetInstalls)
		apps.POST("/installs", ctrl.CreateInstall)
	}
}

func (ctrl *ContentController) GetFiles(c *gin.Context) {
	category := c.Query("category")
	fileType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var files []models.ContentFile
	query := ctrl.DB.Model(&models.ContentFile{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&files)

	c.JSON(http.StatusOK, gin.H{
		"files":     files,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *ContentController) CreateFile(c *gin.Context) {
	var req struct {
		TenantID    uint     `json:"tenant_id"`
		UploaderID  uint     `json:"uploader_id" binding:"required"`
		FileName    string   `json:"file_name" binding:"required"`
		FileType    string   `json:"file_type" binding:"required"`
		FileSize    int64    `json:"file_size"`
		FileURL     string   `json:"file_url" binding:"required"`
		ThumbnailURL string  `json:"thumbnail_url"`
		Category    string   `json:"category" binding:"required"`
		Tags        []string `json:"tags"`
		IsPublic    bool     `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagsJSON, _ := json.Marshal(req.Tags)

	file := models.ContentFile{
		TenantID:    req.TenantID,
		UploaderID:  req.UploaderID,
		FileName:    req.FileName,
		FileType:    req.FileType,
		FileSize:    req.FileSize,
		FileURL:     req.FileURL,
		ThumbnailURL: req.ThumbnailURL,
		Category:    req.Category,
		Tags:        string(tagsJSON),
		IsPublic:    req.IsPublic,
	}

	ctrl.DB.Create(&file)
	c.JSON(http.StatusOK, file)
}

func (ctrl *ContentController) GetFile(c *gin.Context) {
	id := c.Param("id")
	var file models.ContentFile

	if err := ctrl.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (ctrl *ContentController) UpdateFile(c *gin.Context) {
	id := c.Param("id")
	var file models.ContentFile

	if err := ctrl.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&file).Updates(req)
	c.JSON(http.StatusOK, file)
}

func (ctrl *ContentController) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.ContentFile{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *ContentController) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	var file models.ContentFile

	if err := ctrl.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctrl.DB.Model(&file).UpdateColumn("download_count", gorm.Expr("download_count + 1"))

	c.JSON(http.StatusOK, gin.H{
		"download_url": file.FileURL,
		"file_name":    file.FileName,
	})
}

func (ctrl *ContentController) DistributeFile(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   string `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dist := models.ContentDistribution{
		ContentID:  toUintContent(id),
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Status:     "sent",
	}

	now := time.Now()
	dist.SentAt = &now

	ctrl.DB.Create(&dist)
	c.JSON(http.StatusOK, dist)
}

func (ctrl *ContentController) GetDistributions(c *gin.Context) {
	contentID := c.Query("content_id")
	var dists []models.ContentDistribution
	query := ctrl.DB.Model(&models.ContentDistribution{})

	if contentID != "" {
		query = query.Where("content_id = ?", contentID)
	}

	query.Find(&dists)
	c.JSON(http.StatusOK, dists)
}

func (ctrl *ContentController) GetPackages(c *gin.Context) {
	status := c.Query("status")
	platform := c.Query("platform")
	var packages []models.AppPackage
	query := ctrl.DB.Model(&models.AppPackage{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	query.Order("created_at DESC").Find(&packages)
	c.JSON(http.StatusOK, packages)
}

func (ctrl *ContentController) CreatePackage(c *gin.Context) {
	var pkg models.AppPackage
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.Status = "draft"
	ctrl.DB.Create(&pkg)
	c.JSON(http.StatusOK, pkg)
}

func (ctrl *ContentController) GetPackage(c *gin.Context) {
	id := c.Param("id")
	var pkg models.AppPackage

	if err := ctrl.DB.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	c.JSON(http.StatusOK, pkg)
}

func (ctrl *ContentController) UpdatePackage(c *gin.Context) {
	id := c.Param("id")
	var pkg models.AppPackage

	if err := ctrl.DB.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&pkg).Updates(req)
	c.JSON(http.StatusOK, pkg)
}

func (ctrl *ContentController) DeletePackage(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.AppPackage{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *ContentController) SubmitReview(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.AppPackage{}).Where("id = ?", id).Update("status", "pending_review")
	c.JSON(http.StatusOK, gin.H{"message": "submitted for review"})
}

func (ctrl *ContentController) ApprovePackage(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	ctrl.DB.Model(&models.AppPackage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "approved",
		"published_at": &now,
	})
	c.JSON(http.StatusOK, gin.H{"message": "approved"})
}

func (ctrl *ContentController) RejectPackage(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Comment string `json:"comment" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&models.AppPackage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         "rejected",
		"review_comment": req.Comment,
	})
	c.JSON(http.StatusOK, gin.H{"message": "rejected"})
}

func (ctrl *ContentController) GetInstalls(c *gin.Context) {
	appID := c.Query("app_id")
	deviceID := c.Query("device_id")
	var installs []models.AppInstall
	query := ctrl.DB.Model(&models.AppInstall{})

	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	query.Find(&installs)
	c.JSON(http.StatusOK, installs)
}

func (ctrl *ContentController) CreateInstall(c *gin.Context) {
	var install models.AppInstall
	if err := c.ShouldBindJSON(&install); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	install.Status = "installed"
	install.InstallAt = time.Now()
	ctrl.DB.Create(&install)
	c.JSON(http.StatusOK, install)
}

func toUintContent(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}
