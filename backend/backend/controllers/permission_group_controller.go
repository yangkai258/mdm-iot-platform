package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"
	plugins "mdm-backend/plugins"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionGroupController 权限组控制器（多租户版本）
type PermissionGroupController struct {
	DB *gorm.DB
}

// List 权限组列表
// GET /api/v1/permission-groups
func (c *PermissionGroupController) List(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var groups []models.PermGroup
	var total int64

	query := c.DB.Model(&models.PermGroup{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("group_name LIKE ? OR group_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&groups).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      groups,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建权限组
// POST /api/v1/permission-groups
func (c *PermissionGroupController) Create(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "租户ID不能为空"})
		return
	}

	var req struct {
		GroupName   string `json:"group_name" binding:"required"`
		GroupCode   string `json:"group_code" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 检查code唯一性
	var count int64
	c.DB.Model(&models.PermGroup{}).Where("group_code = ? AND tenant_id = ?", req.GroupCode, tenantID).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "权限组编码已存在"})
		return
	}

	group := models.PermGroup{
		GroupName:   req.GroupName,
		GroupCode:   req.GroupCode,
		Description: req.Description,
		Status:      req.Status,
		TenantID:    tenantID,
	}
	if group.Status == 0 {
		group.Status = 1
	}

	db := plugins.WithTenantID(c.DB, tenantID)
	if err := db.Create(&group).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}

// Update 更新权限组
// PUT /api/v1/permission-groups/:id
func (c *PermissionGroupController) Update(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var group models.PermGroup
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "权限组不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		GroupName   string `json:"group_name"`
		GroupCode   string `json:"group_code"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 如果改code，检查唯一性
	if req.GroupCode != "" && req.GroupCode != group.GroupCode {
		var count int64
		c.DB.Model(&models.PermGroup{}).Where("group_code = ? AND tenant_id = ? AND id != ?",
			req.GroupCode, tenantID, id).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "权限组编码已存在"})
			return
		}
	}

	updates := map[string]interface{}{}
	if req.GroupName != "" {
		updates["group_name"] = req.GroupName
	}
	if req.GroupCode != "" {
		updates["group_code"] = req.GroupCode
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := c.DB.Model(&group).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.DB.First(&group, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}

// Delete 删除权限组
// DELETE /api/v1/permission-groups/:id
func (c *PermissionGroupController) Delete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	// 检查是否有角色关联此权限组
	var roleCount int64
	c.DB.Model(&models.RolePermissionGroup{}).Where("permission_group_id = ?", id).Count(&roleCount)
	if roleCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该权限组已被角色使用，无法删除"})
		return
	}

	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.PermGroup{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Get 获取单个权限组
// GET /api/v1/permission-groups/:id
func (c *PermissionGroupController) Get(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var group models.PermGroup
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "权限组不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}
