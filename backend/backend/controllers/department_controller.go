package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DepartmentController 部门管理控制器
type DepartmentController struct {
	DB *gorm.DB
}

// DepartmentTree 部门树形列表
func (c *DepartmentController) DepartmentTree(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var departments []models.Department

	query := c.DB.Model(&models.Department{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 公司筛选
	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if err := query.Order("sort ASC, id DESC").Find(&departments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	tree := buildDeptTree(departments, 0)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tree})
}

// DepartmentList 部门列表（扁平）
func (c *DepartmentController) DepartmentList(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var departments []models.Department
	var total int64

	query := c.DB.Model(&models.Department{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&departments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      departments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DepartmentCreate 创建部门
func (c *DepartmentController) DepartmentCreate(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var dept models.Department
	if err := ctx.ShouldBindJSON(&dept); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if tenantID != "" {
		dept.TenantID = tenantID
	}

	// 生成路径
	if dept.ParentID != nil {
		var parent models.Department
		if err := c.DB.First(&parent, *dept.ParentID).Error; err == nil {
			dept.Level = parent.Level + 1
			dept.Path = fmt.Sprintf("%s/%d", parent.Path, dept.ID)
		dept.ParentID = &parent.ID
		}
	}

	if err := c.DB.Create(&dept).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": dept})
}

// DepartmentGet 获取部门详情
func (c *DepartmentController) DepartmentGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var dept models.Department
	if err := c.DB.First(&dept, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": dept})
}

// DepartmentUpdate 更新部门
func (c *DepartmentController) DepartmentUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var dept models.Department
	if err := c.DB.First(&dept, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 允许更新的字段
	allowedFields := []string{"dept_name", "dept_code", "parent_id", "manager", "phone", "email", "status", "sort", "company_id"}
	updates := map[string]interface{}{}
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	// 如果改了父部门，重新计算层级
	if parentID, ok := updateData["parent_id"]; ok {
		if parentID == nil {
			dept.Level = 1
			dept.Path = ""
		} else {
			var parent models.Department
			if pid, err := strconv.ParseUint(fmt.Sprintf("%v", parentID), 10, 32); err == nil {
				if err := c.DB.First(&parent, pid).Error; err == nil {
					dept.ParentID = &parent.ID
					dept.Level = parent.Level + 1
					dept.Path = fmt.Sprintf("%s/%d", parent.Path, dept.ID)
				}
			}
		}
	}

	if err := c.DB.Model(&dept).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&dept, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": dept})
}

// DepartmentDelete 删除部门
func (c *DepartmentController) DepartmentDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	// 检查是否有子部门
	var count int64
	c.DB.Model(&models.Department{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先删除子部门"})
		return
	}

	// 检查是否有员工
	var empCount int64
	c.DB.Model(&models.Employee{}).Where("dept_id = ?", id).Count(&empCount)
	if empCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该部门下还有员工，无法删除"})
		return
	}

	if err := c.DB.Delete(&models.Department{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
