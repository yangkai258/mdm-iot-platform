package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionController 权限管理控制器
type PermissionController struct {
	DB *gorm.DB
}

// List 权限列表
func (c *PermissionController) List(ctx *gin.Context) {
	var permissions []models.SysPermission
	
	query := c.DB.Model(&models.SysPermission{})
	
	if parentID := ctx.Query("parent_id"); parentID != "" {
		query = query.Where("parent_id = ?", parentID)
	}
	
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if err := query.Order("sort ASC, id ASC").Find(&permissions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	
	// 构建树形结构
	tree := buildPermTree(permissions, 0)
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tree})
}

func buildPermTree(perms []models.SysPermission, parentID uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, p := range perms {
		if p.ParentID == parentID {
			node := map[string]interface{}{
				"id":         p.ID,
				"name":       p.Name,
				"code":       p.Code,
				"type":       p.Type,
				"path":       p.Path,
				"icon":       p.Icon,
				"sort":       p.Sort,
				"visible":    p.Visible,
				"permission": p.Permission,
				"status":     p.Status,
			}
			children := buildPermTree(perms, p.ID)
			if len(children) > 0 {
				node["children"] = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// Create 创建权限
func (c *PermissionController) Create(ctx *gin.Context) {
	var perm models.SysPermission
	if err := ctx.ShouldBindJSON(&perm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	if err := c.DB.Create(&perm).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": perm})
}

// Update 更新权限
func (c *PermissionController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var perm models.SysPermission
	if err := c.DB.First(&perm, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "权限不存在"})
		return
	}
	
	if err := ctx.ShouldBindJSON(&perm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	if err := c.DB.Save(&perm).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": perm})
}

// Delete 删除权限
func (c *PermissionController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	
	// 检查是否有子权限
	var count int64
	c.DB.Model(&models.SysPermission{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先删除子权限"})
		return
	}
	
	if err := c.DB.Delete(&models.SysPermission{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// RoleController 角色管理控制器
type RoleController struct {
	DB *gorm.DB
}

// List 角色列表
func (c *RoleController) List(ctx *gin.Context) {
	var roles []models.SysRole
	var total int64
	
	query := c.DB.Model(&models.SysRole{})
	
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	query.Count(&total)
	
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize
	
	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&roles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":  roles,
			"total": total,
			"page":  page,
			"page_size": pageSize,
		},
	})
}

// Create 创建角色
func (c *RoleController) Create(ctx *gin.Context) {
	var role models.SysRole
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	if err := c.DB.Create(&role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": role})
}

// Update 更新角色
func (c *RoleController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var role models.SysRole
	if err := c.DB.First(&role, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}
	
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	if err := c.DB.Save(&role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": role})
}

// Delete 删除角色
func (c *RoleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.SysRole{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetPerms 获取角色权限
func (c *RoleController) GetPerms(ctx *gin.Context) {
	roleID := ctx.Param("id")
	
	var rolePerms []models.SysRolePermission
	c.DB.Where("role_id = ?", roleID).Find(&rolePerms)
	
	var permIDs []uint
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": permIDs})
}

// SetPerms 设置角色权限
func (c *RoleController) SetPerms(ctx *gin.Context) {
	roleID := ctx.Param("id")
	
	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	// 删除旧权限
	c.DB.Where("role_id = ?", roleID).Delete(&models.SysRolePermission{})
	
	// 添加新权限
	for _, permID := range req.PermissionIDs {
		rolePerm := models.SysRolePermission{
			RoleID:       uint(parseUint(roleID)),
			PermissionID: permID,
		}
		c.DB.Create(&rolePerm)
	}
	
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}
