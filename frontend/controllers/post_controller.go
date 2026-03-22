package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostController 岗位管理控制器
type PostController struct {
	DB *gorm.DB
}

// PostList 岗位列表
func (c *PostController) PostList(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var positions []models.Position
	var total int64

	query := c.DB.Model(&models.Position{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("pos_name LIKE ? OR pos_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if deptID := ctx.Query("dept_id"); deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}

	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&positions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      positions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// PostCreate 创建岗位
func (c *PostController) PostCreate(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var position models.Position
	if err := ctx.ShouldBindJSON(&position); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if tenantID != "" {
		position.TenantID = tenantID
	}

	if err := c.DB.Create(&position).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// PostGet 获取岗位详情
func (c *PostController) PostGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var position models.Position
	if err := c.DB.First(&position, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "岗位不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// PostUpdate 更新岗位
func (c *PostController) PostUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var position models.Position
	if err := c.DB.First(&position, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "岗位不存在"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	allowedFields := []string{"pos_name", "pos_code", "category", "level", "dept_id", "company_id", "description", "status", "sort"}
	updates := map[string]interface{}{}
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := c.DB.Model(&position).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&position, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// PostDelete 删除岗位
func (c *PostController) PostDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	// 检查是否有员工使用该岗位
	var empCount int64
	c.DB.Model(&models.Employee{}).Where("position_id = ?", id).Count(&empCount)
	if empCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该岗位下还有员工，无法删除"})
		return
	}

	if err := c.DB.Delete(&models.Position{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
