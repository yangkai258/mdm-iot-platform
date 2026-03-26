package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuController 菜单控制器
type MenuController struct {
	DB *gorm.DB
}

// GetMenuTree 获取菜单树
func (c *MenuController) GetMenuTree(ctx *gin.Context) {
	var menus []models.SysMenu
	if err := c.DB.Order("sort ASC").Find(&menus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	// 构建树形结构
	tree := buildMenuTree(menus, 0)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": tree,
	})
}

func buildMenuTree(menus []models.SysMenu, parentID uint) []interface{} {
	var result []interface{}
	for _, m := range menus {
		if m.ParentID == parentID {
			item := map[string]interface{}{
				"id":         m.ID,
				"name":       m.Name,
				"path":       m.Path,
				"component":  m.Component,
				"icon":       m.Icon,
				"sort":       m.Sort,
				"visible":    m.Visible,
				"permission": m.Permission,
				"type":       m.Type,
			}
			children := buildMenuTree(menus, m.ID)
			if len(children) > 0 {
				item["children"] = children
			}
			result = append(result, item)
		}
	}
	return result
}

// DictController 字典控制器
type DictController struct {
	DB *gorm.DB
}

// GetDictByType 根据类型获取字典
func (c *DictController) GetDictByType(ctx *gin.Context) {
	dictType := ctx.Param("type")

	var dicts []models.SysDictionary
	if err := c.DB.Where("type = ? AND status = 1", dictType).Order("sort ASC").Find(&dicts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dicts,
	})
}

// List 获取字典列表
// GET /api/v1/dicts
func (c *DictController) List(ctx *gin.Context) {
	var dicts []models.SysDictionary
	var total int64

	query := c.DB.Model(&models.SysDictionary{})

	// 关键字筛选
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR type LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "20")
	offset := (parseInt(page) - 1) * parseInt(pageSize)

	if err := query.Offset(offset).Limit(parseInt(pageSize)).Order("id DESC").Find(&dicts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  dicts,
			"total": total,
			"page":  parseInt(page),
			"page_size": parseInt(pageSize),
		},
	})
}

// Create 创建字典
// POST /api/v1/dicts
func (c *DictController) Create(ctx *gin.Context) {
	var req struct {
		Type   string `json:"type" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Label  string `json:"label"`
		Value  string `json:"value" binding:"required"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
		Remark string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 检查类型+值唯一性
	var exist models.SysDictionary
	if err := c.DB.Where("type = ? AND value = ?", req.Type, req.Value).First(&exist).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "字典类型和值组合已存在"})
		return
	}

	dict := models.SysDictionary{
		Type:   req.Type,
		Name:   req.Name,
		Label:  req.Label,
		Value:  req.Value,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}
	if dict.Status == 0 {
		dict.Status = 1
	}

	if err := c.DB.Create(&dict).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": dict, "message": "success"})
}

// Update 更新字典
// PUT /api/v1/dicts/:id
func (c *DictController) Update(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var dict models.SysDictionary
	if err := c.DB.First(&dict, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Type   string `json:"type"`
		Name   string `json:"name"`
		Label  string `json:"label"`
		Value  string `json:"value"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
		Remark string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查类型+值唯一性（排除自身）
	if req.Type != "" && req.Value != "" {
		var exist models.SysDictionary
		if err := c.DB.Where("type = ? AND value = ? AND id != ?", req.Type, req.Value, id).First(&exist).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "字典类型和值组合已存在"})
			return
		}
	}

	updates := map[string]interface{}{}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Label != "" {
		updates["label"] = req.Label
	}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := c.DB.Model(&dict).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Delete 删除字典
// DELETE /api/v1/dicts/:id
func (c *DictController) Delete(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var dict models.SysDictionary
	if err := c.DB.First(&dict, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Delete(&dict).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// LogController 日志控制器
type LogController struct {
	DB *gorm.DB
}

// GetOperationLogs 获取操作日志
func (c *LogController) GetOperationLogs(ctx *gin.Context) {
	query := c.DB.Model(&models.SysOperationLog{})

	// 过滤条件
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if module := ctx.Query("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	if startTime := ctx.Query("start_time"); startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var logs []models.SysOperationLog
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "20")
	
	query.Order("id DESC").Offset((parseInt(page) - 1) * parseInt(pageSize)).Limit(parseInt(pageSize)).Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":       logs,
			"pagination": gin.H{
				"total":    total,
				"current":  parseInt(page),
				"pageSize": parseInt(pageSize),
			},
		},
	})
}

// GetLoginLogs 获取登录日志
func (c *LogController) GetLoginLogs(ctx *gin.Context) {
	var logs []models.SysLoginLog
	c.DB.Order("id DESC").Limit(100).Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": logs,
		},
	})
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
