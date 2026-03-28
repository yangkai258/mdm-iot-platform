package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KnowledgeVersionController 知识库版本控制器
type KnowledgeVersionController struct {
	DB *gorm.DB
}

// NewKnowledgeVersionController 创建控制器
func NewKnowledgeVersionController(db *gorm.DB) *KnowledgeVersionController {
	return &KnowledgeVersionController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *KnowledgeVersionController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/knowledge/versions", ctrl.ListVersions)
	rg.GET("/knowledge/versions/:id", ctrl.GetVersion)
	rg.POST("/knowledge/versions", ctrl.CreateVersion)
	rg.PUT("/knowledge/versions/:id", ctrl.UpdateVersion)
	rg.DELETE("/knowledge/versions/:id", ctrl.DeleteVersion)
	rg.GET("/knowledge/:id/versions", ctrl.ListVersionsByKnowledge)
	rg.POST("/knowledge/versions/:id/publish", ctrl.PublishVersion)
	rg.GET("/knowledge/versions/:id/reviews", ctrl.GetVersionReviews)
	rg.POST("/knowledge/versions/:id/reviews", ctrl.AddVersionReview)
}

// ListVersions 获取版本列表
func (ctrl *KnowledgeVersionController) ListVersions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	knowledgeID := c.Query("knowledge_id")
	status := c.DefaultQuery("status", "")

	if page < 1 {
		page = 1
	}

	query := ctrl.DB.Model(&models.KnowledgeVersion{})

	if knowledgeID != "" {
		query = query.Where("knowledge_id = ?", knowledgeID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	var list []models.KnowledgeVersion
	query.Count(&total)

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取版本列表失败"})
		return
	}

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
func (ctrl *KnowledgeVersionController) GetVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.KnowledgeVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": version})
}

// CreateVersion 创建新版本
func (ctrl *KnowledgeVersionController) CreateVersion(c *gin.Context) {
	var req struct {
		KnowledgeID uint   `json:"knowledge_id" binding:"required"`
		Version    string `json:"version" binding:"required"`
		Content    string `json:"content"`
		ChangeLog  string `json:"change_log"`
		FileURL    string `json:"file_url"`
		FileSize   int64  `json:"file_size"`
		ChangeType string `json:"change_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 计算内容哈希
	hash := sha256.Sum256([]byte(req.Content))
	contentHash := hex.EncodeToString(hash[:])

	version := models.KnowledgeVersion{
		KnowledgeID: req.KnowledgeID,
		Version:     req.Version,
		Content:     req.Content,
		ChangeLog:   req.ChangeLog,
		ContentHash: contentHash,
		FileURL:     req.FileURL,
		FileSize:    req.FileSize,
		ChangeType:  req.ChangeType,
		Status:      "draft",
		CreatedBy:   c.GetString("username"),
	}

	if err := ctrl.DB.Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": version})
}

// UpdateVersion 更新版本
func (ctrl *KnowledgeVersionController) UpdateVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.KnowledgeVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	if version.Status == "published" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已发布的版本不能修改"})
		return
	}

	var req struct {
		Version    string `json:"version"`
		Content    string `json:"content"`
		ChangeLog  string `json:"change_log"`
		FileURL    string `json:"file_url"`
		ChangeType string `json:"change_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Version != "" {
		updates["version"] = req.Version
	}
	if req.Content != "" {
		updates["content"] = req.Content
		hash := sha256.Sum256([]byte(req.Content))
		updates["content_hash"] = hex.EncodeToString(hash[:])
	}
	if req.ChangeLog != "" {
		updates["change_log"] = req.ChangeLog
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.ChangeType != "" {
		updates["change_type"] = req.ChangeType
	}

	if err := ctrl.DB.Model(&version).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.Where("id = ?", id).First(&version)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": version})
}

// DeleteVersion 删除版本
func (ctrl *KnowledgeVersionController) DeleteVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.KnowledgeVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	if version.Status == "published" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已发布的版本不能删除"})
		return
	}

	ctrl.DB.Delete(&version)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListVersionsByKnowledge 获取指定知识的版本列表
func (ctrl *KnowledgeVersionController) ListVersionsByKnowledge(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	var list []models.KnowledgeVersion
	ctrl.DB.Model(&models.KnowledgeVersion{}).Where("knowledge_id = ?", id).Count(&total)

	ctrl.DB.Where("knowledge_id = ?", id).Order("created_at DESC").
		Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": list, "total": total, "page": page, "page_size": pageSize},
	})
}

// PublishVersion 发布版本
func (ctrl *KnowledgeVersionController) PublishVersion(c *gin.Context) {
	id := c.Param("id")

	var version models.KnowledgeVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&version).Updates(map[string]interface{}{
		"status":       "published",
		"published_at":  now,
		"published_by":  c.GetString("username"),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发布成功", "data": version})
}

// GetVersionReviews 获取版本审核记录
func (ctrl *KnowledgeVersionController) GetVersionReviews(c *gin.Context) {
	id := c.Param("id")

	var reviews []models.KnowledgeVersionReview
	ctrl.DB.Where("version_id = ?", id).Order("created_at DESC").Find(&reviews)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": reviews})
}

// AddVersionReview 添加版本审核
func (ctrl *KnowledgeVersionController) AddVersionReview(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ReviewStatus  string `json:"review_status" binding:"required"`
		ReviewComment string `json:"review_comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var version models.KnowledgeVersion
	if err := ctrl.DB.Where("id = ?", id).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "版本不存在"})
		return
	}

	review := models.KnowledgeVersionReview{
		VersionID:     version.ID,
		ReviewStatus:  req.ReviewStatus,
		ReviewComment: req.ReviewComment,
		ReviewerID:    c.GetString("user_id"),
		ReviewerName: c.GetString("username"),
	}
	now := time.Now()
	review.ReviewedAt = &now

	ctrl.DB.Create(&review)

	// 如果审核通过，自动发布
	if req.ReviewStatus == "approved" {
		ctrl.DB.Model(&version).Updates(map[string]interface{}{
			"status":       "published",
			"published_at":  time.Now(),
			"published_by":  c.GetString("username"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "审核完成", "data": review})
}
