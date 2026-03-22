package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"
	plugins "mdm-backend/plugins"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRoleController 角色管理控制器（多租户版本）
type NewRoleController struct {
	DB *gorm.DB
}

// List 角色列表
// GET /api/v1/roles
func (c *NewRoleController) List(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var roles []models.Role
	var total int64

	query := c.DB.Model(&models.Role{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("role_name LIKE ? OR role_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 时间范围筛选
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&roles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      roles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建角色
// POST /api/v1/roles
func (c *NewRoleController) Create(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "租户ID不能为空"})
		return
	}

	var req struct {
		RoleName    string `json:"role_name" binding:"required"`
		RoleCode    string `json:"role_code" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 检查code唯一性
	var count int64
	c.DB.Model(&models.Role{}).Where("role_code = ? AND tenant_id = ?", req.RoleCode, tenantID).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "角色编码已存在"})
		return
	}

	role := models.Role{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Description: req.Description,
		Status:      req.Status,
		TenantID:    tenantID,
	}
	if role.Status == 0 {
		role.Status = 1
	}

	db := plugins.WithTenantID(c.DB, tenantID)
	if err := db.Create(&role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": role})
}

// Update 更新角色
// PUT /api/v1/roles/:id
func (c *NewRoleController) Update(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var role models.Role
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		RoleName    string `json:"role_name"`
		RoleCode    string `json:"role_code"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 如果改code，检查唯一性
	if req.RoleCode != "" && req.RoleCode != role.RoleCode {
		var cnt int64
		c.DB.Model(&models.Role{}).Where("role_code = ? AND tenant_id = ? AND id != ?",
			req.RoleCode, tenantID, id).Count(&cnt)
		if cnt > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "角色编码已存在"})
			return
		}
	}

	updates := map[string]interface{}{}
	if req.RoleName != "" {
		updates["role_name"] = req.RoleName
	}
	if req.RoleCode != "" {
		updates["role_code"] = req.RoleCode
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := c.DB.Model(&role).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.DB.First(&role, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": role})
}

// Delete 删除角色
// DELETE /api/v1/roles/:id
func (c *NewRoleController) Delete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	// 检查是否有用户使用此角色
	var userCount int64
	c.DB.Model(&models.SysUserExt{}).Where("role_ids LIKE ?", "%"+id+"%").Count(&userCount)
	if userCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该角色下有用户，无法删除"})
		return
	}

	// 删除角色菜单关联
	c.DB.Where("role_id = ?", id).Delete(&models.RoleMenu{})
	// 删除角色API权限关联
	c.DB.Where("role_id = ?", id).Delete(&models.RoleApiPermission{})
	// 删除角色权限组关联
	c.DB.Where("role_id = ?", id).Delete(&models.RolePermissionGroup{})

	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Role{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetPermissions 获取角色权限
// GET /api/v1/roles/:id/permissions
func (c *NewRoleController) GetPermissions(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	roleIDStr := ctx.Param("id")

	// 获取角色关联的菜单ID
	var roleMenus []models.RoleMenu
	c.DB.Where("role_id = ?", roleIDStr).Find(&roleMenus)
	var menuIDs []uint
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	// 获取角色关联的API权限ID
	var roleApis []models.RoleApiPermission
	c.DB.Where("role_id = ?", roleIDStr).Find(&roleApis)
	var apiIDs []uint
	for _, ra := range roleApis {
		apiIDs = append(apiIDs, ra.ApiPermissionID)
	}

	// 获取角色关联的权限组ID
	var roleGroups []models.RolePermissionGroup
	c.DB.Where("role_id = ?", roleIDStr).Find(&roleGroups)
	var groupIDs []uint
	for _, rg := range roleGroups {
		groupIDs = append(groupIDs, rg.PermissionGroupID)
	}

	// 获取菜单详情
	var menus []models.Menu
	if len(menuIDs) > 0 && tenantID != "" {
		c.DB.Where("id IN ? AND tenant_id = ?", menuIDs, tenantID).Find(&menus)
	}

	// 获取API权限详情
	var apis []models.ApiPermission
	if len(apiIDs) > 0 && tenantID != "" {
		c.DB.Where("id IN ? AND tenant_id = ?", apiIDs, tenantID).Find(&apis)
	}

	// 获取权限组详情
	var groups []models.PermGroup
	if len(groupIDs) > 0 && tenantID != "" {
		c.DB.Where("id IN ? AND tenant_id = ?", groupIDs, tenantID).Find(&groups)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"menu_ids":  menuIDs,
			"api_ids":   apiIDs,
			"group_ids": groupIDs,
			"menus":     menus,
			"apis":      apis,
			"groups":    groups,
		},
	})
}

// SetPermissions 设置角色权限
// PUT /api/v1/roles/:id/permissions
func (c *NewRoleController) SetPermissions(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	roleIDStr := ctx.Param("id")

	// 检查角色是否存在
	var role models.Role
	if err := c.DB.Where("id = ? AND tenant_id = ?", roleIDStr, tenantID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		MenuIDs  []uint `json:"menu_ids"`
		ApiIDs   []uint `json:"api_ids"`
		GroupIDs []uint `json:"group_ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	roleID, _ := strconv.ParseUint(roleIDStr, 10, 32)

	// 事务处理
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		// 删除旧关联
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleApiPermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermissionGroup{}).Error; err != nil {
			return err
		}

		// 插入新菜单关联
		for _, menuID := range req.MenuIDs {
			var cnt int64
			tx.Model(&models.Menu{}).Where("id = ? AND tenant_id = ?", menuID, tenantID).Count(&cnt)
			if cnt == 0 {
				continue
			}
			if err := tx.Create(&models.RoleMenu{RoleID: uint(roleID), MenuID: menuID}).Error; err != nil {
				return err
			}
		}

		// 插入新API权限关联
		for _, apiID := range req.ApiIDs {
			var cnt int64
			tx.Model(&models.ApiPermission{}).Where("id = ? AND tenant_id = ?", apiID, tenantID).Count(&cnt)
			if cnt == 0 {
				continue
			}
			if err := tx.Create(&models.RoleApiPermission{RoleID: uint(roleID), ApiPermissionID: apiID}).Error; err != nil {
				return err
			}
		}

		// 插入新权限组关联
		for _, groupID := range req.GroupIDs {
			var cnt int64
			tx.Model(&models.PermGroup{}).Where("id = ? AND tenant_id = ?", groupID, tenantID).Count(&cnt)
			if cnt == 0 {
				continue
			}
			if err := tx.Create(&models.RolePermissionGroup{RoleID: uint(roleID), PermissionGroupID: groupID}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "设置权限失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ 以下为旧版 RoleController（兼容 device_controller.go）============

// OldRoleController 旧版角色管理控制器（基于 SysRole）
type OldRoleController struct {
	DB *gorm.DB
}

// List 旧版角色列表
func (c *OldRoleController) List(ctx *gin.Context) {
	var roles []models.SysRole
	var total int64

	query := c.DB.Model(&models.SysRole{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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

	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&roles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      roles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 旧版获取单个角色
func (c *OldRoleController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var role models.SysRole

	if err := c.DB.First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var rolePerms []models.SysRolePermission
	c.DB.Where("role_id = ?", id).Find(&rolePerms)
	var permIDs []uint
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	var perms []models.SysPermission
	if len(permIDs) > 0 {
		c.DB.Where("id IN ?", permIDs).Find(&perms)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"role":     role,
			"perms":    perms,
			"perm_ids": permIDs,
		},
	})
}

// Create 旧版创建角色
func (c *OldRoleController) Create(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
		Sort        int    `json:"sort"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	role := models.SysRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
		Sort:        req.Sort,
	}
	if role.Status == 0 {
		role.Status = 1
	}

	if err := c.DB.Create(&role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": role})
}

// Update 旧版更新角色
func (c *OldRoleController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var role models.SysRole
	if err := c.DB.First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      int    `json:"status"`
		Sort        int    `json:"sort"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}

	if len(updates) > 0 {
		if err := c.DB.Model(&role).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.DB.First(&role, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": role})
}

// Delete 旧版删除角色
func (c *OldRoleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	var userCount int64
	c.DB.Model(&models.SysUser{}).Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该角色下有用户，无法删除"})
		return
	}

	c.DB.Where("role_id = ?", id).Delete(&models.SysRolePermission{})

	if err := c.DB.Delete(&models.SysRole{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetPermissions 旧版获取角色权限
func (c *OldRoleController) GetPermissions(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	var rolePerms []models.SysRolePermission
	c.DB.Where("role_id = ?", roleID).Find(&rolePerms)
	var permIDs []uint
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	var perms []models.SysPermission
	if len(permIDs) > 0 {
		c.DB.Where("id IN ?", permIDs).Find(&perms)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"permission_ids": permIDs,
			"permissions":    perms,
		},
	})
}

// AssignPermissions 旧版分配权限
func (c *OldRoleController) AssignPermissions(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的角色ID"})
		return
	}

	var role models.SysRole
	if err := c.DB.First(&role, roleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	err = c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&models.SysRolePermission{}).Error; err != nil {
			return err
		}

		for _, permID := range req.PermissionIDs {
			var count int64
			tx.Model(&models.SysPermission{}).Where("id = ?", permID).Count(&count)
			if count == 0 {
				continue
			}

			rp := models.SysRolePermission{
				RoleID:       uint(roleID),
				PermissionID: permID,
			}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "分配权限失败"})
		return
	}

	var perms []models.SysPermission
	if len(req.PermissionIDs) > 0 {
		c.DB.Where("id IN ? AND status = 1", req.PermissionIDs).Find(&perms)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"permissions": perms,
		},
	})
}

// ListPermissions 旧版获取所有权限点（树形）
func (c *OldRoleController) ListPermissions(ctx *gin.Context) {
	var perms []models.SysPermission

	query := c.DB.Model(&models.SysPermission{})

	if parentID := ctx.Query("parent_id"); parentID != "" {
		query = query.Where("parent_id = ?", parentID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR permission LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if permType := ctx.Query("type"); permType != "" {
		query = query.Where("type = ?", permType)
	}

	if err := query.Order("sort ASC, id ASC").Find(&perms).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	tree := BuildPermTree(perms, 0)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tree})
}

// BuildPermTree 构建权限树
func BuildPermTree(perms []models.SysPermission, parentID uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, p := range perms {
		if p.ParentID == parentID {
			node := map[string]interface{}{
				"id":         p.ID,
				"name":       p.Name,
				"code":       p.Code,
				"type":       p.Type,
				"path":       p.Path,
				"component":  p.Component,
				"icon":       p.Icon,
				"sort":       p.Sort,
				"visible":    p.Visible,
				"permission": p.Permission,
				"status":     p.Status,
			}
			children := BuildPermTree(perms, p.ID)
			if len(children) > 0 {
				node["children"] = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}
