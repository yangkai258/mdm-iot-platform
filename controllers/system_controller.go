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
