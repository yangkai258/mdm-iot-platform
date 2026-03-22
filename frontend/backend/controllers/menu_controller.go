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

// TenantMenuController 菜单管理控制器（多租户版本）
type TenantMenuController struct {
	DB *gorm.DB
}

// List 菜单列表（树形）
// GET /api/v1/menus
func (c *TenantMenuController) List(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var menus []models.Menu

	query := c.DB.Model(&models.Menu{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	tree := buildTenantMenuTree(menus, 0)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tree})
}

// buildTenantMenuTree 构建菜单树
func buildTenantMenuTree(menus []models.Menu, parentID uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, m := range menus {
		var pid uint
		if m.ParentID != nil {
			pid = *m.ParentID
		}
		if pid == parentID {
			node := map[string]interface{}{
				"id":         m.ID,
				"parent_id":  m.ParentID,
				"menu_name":  m.MenuName,
				"menu_code":  m.MenuCode,
				"icon":       m.Icon,
				"route_path": m.RoutePath,
				"component":  m.Component,
				"permission": m.Permission,
				"menu_type":  m.MenuType,
				"sort":       m.Sort,
				"status":     m.Status,
			}
			children := buildTenantMenuTree(menus, m.ID)
			if len(children) > 0 {
				node["children"] = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// Create 创建菜单
// POST /api/v1/menus
func (c *TenantMenuController) Create(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "租户ID不能为空"})
		return
	}

	var req struct {
		ParentID   *uint    `json:"parent_id"`
		MenuName   string   `json:"menu_name" binding:"required"`
		MenuCode   string   `json:"menu_code"`
		Icon       string   `json:"icon"`
		RoutePath  string   `json:"route_path"`
		Component  string   `json:"component"`
		Permission string   `json:"permission"`
		MenuType   int      `json:"menu_type"`
		Sort       int      `json:"sort"`
		Status     int      `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	menu := models.Menu{
		ParentID:   req.ParentID,
		MenuName:   req.MenuName,
		MenuCode:   req.MenuCode,
		Icon:       req.Icon,
		RoutePath:  req.RoutePath,
		Component:  req.Component,
		Permission: req.Permission,
		MenuType:   req.MenuType,
		Sort:       req.Sort,
		Status:     req.Status,
		TenantID:   tenantID,
	}
	if menu.MenuType == 0 {
		menu.MenuType = 1
	}
	if menu.Status == 0 {
		menu.Status = 1
	}

	db := plugins.WithTenantID(c.DB, tenantID)
	if err := db.Create(&menu).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": menu})
}

// Update 更新菜单
// PUT /api/v1/menus/:id
func (c *TenantMenuController) Update(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var menu models.Menu
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&menu).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		ParentID   *uint   `json:"parent_id"`
		MenuName   string  `json:"menu_name"`
		MenuCode   string  `json:"menu_code"`
		Icon       string  `json:"icon"`
		RoutePath  string  `json:"route_path"`
		Component  string  `json:"component"`
		Permission string  `json:"permission"`
		MenuType   int     `json:"menu_type"`
		Sort       int     `json:"sort"`
		Status     int     `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.ParentID != nil {
		updates["parent_id"] = req.ParentID
	}
	if req.MenuName != "" {
		updates["menu_name"] = req.MenuName
	}
	if req.MenuCode != "" {
		updates["menu_code"] = req.MenuCode
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.RoutePath != "" {
		updates["route_path"] = req.RoutePath
	}
	if req.Component != "" {
		updates["component"] = req.Component
	}
	if req.Permission != "" {
		updates["permission"] = req.Permission
	}
	if req.MenuType > 0 {
		updates["menu_type"] = req.MenuType
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := c.DB.Model(&menu).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.DB.First(&menu, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": menu})
}

// Delete 删除菜单
// DELETE /api/v1/menus/:id
func (c *TenantMenuController) Delete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	// 检查是否有子菜单
	var childCount int64
	c.DB.Model(&models.Menu{}).Where("parent_id = ? AND tenant_id = ?", id, tenantID).Count(&childCount)
	if childCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该菜单下有子菜单，无法删除"})
		return
	}

	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Menu{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Get 获取单个菜单
// GET /api/v1/menus/:id
func (c *TenantMenuController) Get(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var menu models.Menu
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&menu).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": menu})
}

// parseUint 解析uint
func parseUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint(v)
}
