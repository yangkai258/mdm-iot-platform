package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberTagController 会员标签控制器
type MemberTagController struct {
	DB *gorm.DB
}

func NewMemberTagController(db *gorm.DB) *MemberTagController {
	return &MemberTagController{DB: db}
}

// RegisterRoutes 注册会员标签路由
func (ctrl *MemberTagController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tags", ctrl.ListTags)
	rg.POST("/tags", ctrl.CreateTag)
	rg.GET("/members/:id/tags", ctrl.GetMemberTags)
	rg.POST("/members/:id/tags", ctrl.AddMemberTag)
	rg.DELETE("/members/:id/tags/:tag_id", ctrl.RemoveMemberTag)
}

// ListTags 获取标签列表
// GET /api/v1/tags
func (ctrl *MemberTagController) ListTags(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.MemberTagDef{}).Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("tag_name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var tags []models.MemberTagDef
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&tags)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      tags,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateTag 创建标签定义
// POST /api/v1/tags
func (ctrl *MemberTagController) CreateTag(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		TagName     string `json:"tag_name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Sort        int    `json:"sort"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tag := models.MemberTagDef{
		TagName:     req.TagName,
		TagColor:    req.Color,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
		TenantID:    tenantID,
	}

	if err := ctrl.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tag})
}

// GetMemberTags 获取会员的标签
// GET /api/v1/members/:id/tags
func (ctrl *MemberTagController) GetMemberTags(c *gin.Context) {
	memberID := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var relations []models.MemberTagRelation
	ctrl.DB.Where("member_id = ? AND tenant_id = ?", memberID, tenantID).Find(&relations)

	var tagIDs []uint
	for _, r := range relations {
		tagIDs = append(tagIDs, r.TagID)
	}

	var tags []models.MemberTagDef
	if len(tagIDs) > 0 {
		ctrl.DB.Where("id IN ?", tagIDs).Find(&tags)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tags})
}

// AddMemberTag 给会员添加标签
// POST /api/v1/members/:id/tags
func (ctrl *MemberTagController) AddMemberTag(c *gin.Context) {
	memberID := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var req struct {
		TagID uint `json:"tag_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查是否已存在
	var existing models.MemberTagRelation
	if err := ctrl.DB.Where("member_id = ? AND tag_id = ? AND tenant_id = ?", memberID, req.TagID, tenantID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "标签已存在", "data": existing})
		return
	}

	relation := models.MemberTagRelation{
		MemberID: parseUintTag(memberID),
		TagID:    req.TagID,
		TenantID: tenantID,
	}

	if err := ctrl.DB.Create(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": relation})
}

// RemoveMemberTag 删除会员标签
// DELETE /api/v1/members/:id/tags/:tag_id
func (ctrl *MemberTagController) RemoveMemberTag(c *gin.Context) {
	memberID := c.Param("id")
	tagID := c.Param("tag_id")
	tenantID := middleware.GetTenantID(c)

	if err := ctrl.DB.Where("member_id = ? AND tag_id = ? AND tenant_id = ?", memberID, tagID, tenantID).Delete(&models.MemberTagRelation{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func parseUintTag(s string) uint {
	id, _ := strconv.ParseUint(s, 10, 32)
	return uint(id)
}
