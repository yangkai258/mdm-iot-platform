package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContentVersionController 内容版本控制器
type ContentVersionController struct {
	DB *gorm.DB
}

// NewContentVersionController 创建控制器
func NewContentVersionController(db *gorm.DB) *ContentVersionController {
	return &ContentVersionController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *ContentVersionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/content/versions", ctrl.ListVersions)
	rg.GET("/content/versions/:id", ctrl.GetVersion)
	rg.POST("/content/versions", ctrl.CreateVersion)
	rg.PUT("/content/versions/:id", ctrl.UpdateVersion)
	rg.DELETE("/content/versions/:id", ctrl.DeleteVersion)
	rg.GET("/content/:id/versions", ctrl.ListContentVersions)
	rg.POST("/content/versions/:id/publish", ctrl.PublishVersion)
	rg.POST("/content/versions/:id/rollback", ctrl.RollbackVersion)
	rg.GET("/content/versions/:id/reviews", ctrl.GetReviews)
	rg.POST("/content/versions/:id/reviews", ctrl.AddReview)
}

// ListVersions 获取版本列表
func (ctrl *ContentVersionController) ListVersions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	contentID := c.Query("content_id")
	status := c.Query("status")

	query := ctrl.DB.Model(&models.ContentVersion{})

	if contentID != "" {
		query = query.Where("content_id = ?", contentID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	var list []models.ContentVersion
	query.Count(&total)

	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

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

// GetVersion 获取单个版本
func (ctrl *ContentVersionController) GetVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.ContentVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": version})
}

// CreateVersion 创建新版本
func (ctrl *ContentVersionController) CreateVersion(c *gin.Context) {
	var req struct {
		ContentID   uint   `json:"content_id" binding:"required"`
		Version    string `json:"version" binding:"required"`
		Title      string `json:"title"`
		Description string `json:"description"`
		FileURL    string `json:"file_url"`
		FileSize   int64  `json:"file_size"`
		ChangeLog  string `json:"change_log"`
		ChangeType string `json:"change_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 计算版本号
	var maxVersionNum int
	ctrl.DB.Model(&models.ContentVersion{}).
		Where("content_id = ?", req.ContentID).
		Select("COALESCE(MAX(version_num), 0)").
		Row().Scan(&maxVersionNum)

	versionNum := maxVersionNum + 1

	// 计算内容哈希
	hash := sha256.Sum256([]byte(req.Description + req.FileURL))
	contentHash := hex.EncodeToString(hash[:])

	version := models.ContentVersion{
		ContentID:   req.ContentID,
		Version:    req.Version,
		VersionNum:  versionNum,
		Title:      req.Title,
		Description: req.Description,
		FileURL:     req.FileURL,
		FileSize:   req.FileSize,
		ContentHash: contentHash,
		ChangeLog:   req.ChangeLog,
		ChangeType:  req.ChangeType,
		Status:      "draft",
		IsLatest:    false,
		CreatedBy:   c.GetString("username"),
	}

	if err := ctrl.DB.Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": version})
}

// UpdateVersion 更新版本
func (ctrl *ContentVersionController) UpdateVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.ContentVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	if version.Status == "published" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已发布版本不能修改"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		FileURL     string `json:"file_url"`
		ChangeLog   string `json:"change_log"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
		hash := sha256.Sum256([]byte(req.Description + version.FileURL))
		updates["content_hash"] = hex.EncodeToString(hash[:])
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.ChangeLog != "" {
		updates["change_log"] = req.ChangeLog
	}

	ctrl.DB.Model(&version).Updates(updates)
	ctrl.DB.Where("id = ?", id).First(&version)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": version})
}

// DeleteVersion 删除版本
func (ctrl *ContentVersionController) DeleteVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.ContentVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	if version.IsLatest {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不能删除最新版本"})
		return
	}

	ctrl.DB.Delete(&version)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListContentVersions 获取指定内容的版本列表
func (ctrl *ContentVersionController) ListContentVersions(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	var list []models.ContentVersion
	ctrl.DB.Model(&models.ContentVersion{}).Where("content_id = ?", id).Count(&total)

	ctrl.DB.Where("content_id = ?", id).Order("version_num DESC").
		Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

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

// PublishVersion 发布版本
func (ctrl *ContentVersionController) PublishVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.ContentVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	// 取消之前的最新标记
	ctrl.DB.Model(&models.ContentVersion{}).
		Where("content_id = ? AND is_latest = ?", version.ContentID, true).
		Update("is_latest", false)

	// 发布新版本
	now := time.Now()
	if err := ctrl.DB.Model(&version).Updates(map[string]interface{}{
		"status":       "published",
		"is_latest":    true,
		"published_at":  now,
		"published_by":  c.GetString("username"),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发布成功", "data": version})
}

// RollbackVersion 回滚到指定版本
func (ctrl *ContentVersionController) RollbackVersion(c *gin.Context) {
	id := c.Param("id")

	var targetVersion models.ContentVersion
	if err := ctrl.DB.Where("id = ?", id).First(&targetVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	if targetVersion.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只能回滚已发布版本"})
		return
	}

	// 创建新版本，内容复制自目标版本
	var maxVersionNum int
	ctrl.DB.Model(&models.ContentVersion{}).
		Where("content_id = ?", targetVersion.ContentID).
		Select("COALESCE(MAX(version_num), 0)").Row().Scan(&maxVersionNum)

	newVersion := models.ContentVersion{
		ContentID:   targetVersion.ContentID,
		Version:    fmt.Sprintf("%d.0.0", maxVersionNum+1),
		VersionNum:  maxVersionNum + 1,
		Title:       targetVersion.Title,
		Description: targetVersion.Description,
		FileURL:     targetVersion.FileURL,
		FileSize:    targetVersion.FileSize,
		ContentHash: targetVersion.ContentHash,
		ChangeLog:   fmt.Sprintf("回滚到版本 %s", targetVersion.Version),
		ChangeType:  "fix",
		Status:      "draft",
		CreatedBy:   c.GetString("username"),
	}

	ctrl.DB.Create(&newVersion)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已创建新版本，内容与目标版本相同",
		"data": gin.H{
			"new_version": newVersion,
			"rollback_from": targetVersion.Version,
		},
	})
}

// GetReviews 获取版本审核记录
func (ctrl *ContentVersionController) GetReviews(c *gin.Context) {
	id := c.Param("id")

	var reviews []models.ContentVersionReview
	ctrl.DB.Where("version_id = ?", id).Order("created_at DESC").Find(&reviews)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": reviews})
}

// AddReview 添加版本审核
func (ctrl *ContentVersionController) AddReview(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ReviewType   string  `json:"review_type" binding:"required"`
		ReviewStatus string  `json:"review_status" binding:"required"`
		ReviewScore  float64 `json:"review_score"`
		ReviewReport string  `json:"review_report"`
		Issues       string  `json:"issues"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var version models.ContentVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	now := time.Now()
	review := models.ContentVersionReview{
		VersionID:    version.ID,
		ReviewType:   req.ReviewType,
		ReviewStatus: req.ReviewStatus,
		ReviewScore:  req.ReviewScore,
		ReviewReport: req.ReviewReport,
		Issues:       req.Issues,
		ReviewerID:   c.GetString("user_id"),
		ReviewerName: c.GetString("username"),
		ReviewedAt:   &now,
	}

	ctrl.DB.Create(&review)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "审核完成", "data": review})
}

// InitContentVersionTable 添加 version_num 和 content_hash 字段到 content_files 表（如果不存在）
// 这个函数在 migration 时调用
func AddContentVersionFields(db *gorm.DB) error {
	// 这些字段应该通过 ALTER TABLE 添加，但在 GORM AutoMigrate 中会被忽略
	// 实际使用需要手动执行 SQL
	return nil
}
