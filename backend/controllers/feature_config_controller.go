package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FeatureConfigController 功能配置管理
type FeatureConfigController struct {
	DB *gorm.DB
}

// NewFeatureConfigController 创建控制器
func NewFeatureConfigController(db *gorm.DB) *FeatureConfigController {
	return &FeatureConfigController{DB: db}
}

// ListGroups 获取所有分组(含功能项)
func (ctrl *FeatureConfigController) ListGroups(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	var groups []models.FeatureGroup
	err := ctrl.DB.Where("tenant_id = ?", tenantID).
		Order("sort ASC, created_at DESC").
		Preload("Features", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC, created_at DESC")
		}).
		Find(&groups).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取分组列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
}

// GetGroup 获取单个分组
func (ctrl *FeatureConfigController) GetGroup(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	var group models.FeatureGroup
	err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("Features", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC, created_at DESC")
		}).
		First(&group).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "分组不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取分组失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": group})
}

// CreateGroup 创建分组
func (ctrl *FeatureConfigController) CreateGroup(c *gin.Context) {
	var req models.FeatureGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误", "error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	group := models.FeatureGroup{
		GroupName:   req.GroupName,
		GroupCode:   req.GroupCode,
		Icon:        req.Icon,
		Color:       req.Color,
		Sort:        req.Sort,
		Description: req.Description,
		Status:      req.Status,
		TenantID:    tenantID,
	}

	// 如果没传GroupCode，自动生成
	if group.GroupCode == "" {
		group.GroupCode = "group_" + strconv.FormatUint(uint64(group.ID), 10)
	}

	err := ctrl.DB.Create(&group).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建分组失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": group, "message": "创建成功"})
}

// UpdateGroup 更新分组
func (ctrl *FeatureConfigController) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var req models.FeatureGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误", "error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	updates := map[string]interface{}{
		"group_name":  req.GroupName,
		"icon":        req.Icon,
		"color":       req.Color,
		"sort":        req.Sort,
		"description": req.Description,
		"status":      req.Status,
	}

	if req.GroupCode != "" {
		updates["group_code"] = req.GroupCode
	}

	err := ctrl.DB.Model(&models.FeatureGroup{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(updates).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新分组失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// DeleteGroup 删除分组(同时删除组内所有功能项)
func (ctrl *FeatureConfigController) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	tx := ctrl.DB.Begin()

	// 删除组内功能项
	if err := tx.Where("group_id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.FeatureItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "删除组内功能项失败", "error": err.Error()})
		return
	}

	// 删除分组
	if err := tx.Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.FeatureGroup{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "删除分组失败", "error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListFeatures 获取所有功能项
func (ctrl *FeatureConfigController) ListFeatures(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	groupID := c.Query("group_id")

	var items []models.FeatureItem
	query := ctrl.DB.Where("tenant_id = ?", tenantID)

	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	err := query.Order("sort ASC, created_at DESC").Find(&items).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取功能列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": items})
}

// GetFeature 获取单个功能项
func (ctrl *FeatureConfigController) GetFeature(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	var item models.FeatureItem
	err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "功能不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取功能失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": item})
}

// CreateFeature 创建功能项
func (ctrl *FeatureConfigController) CreateFeature(c *gin.Context) {
	var req models.FeatureItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误", "error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	item := models.FeatureItem{
		GroupID:     req.GroupID,
		FeatureName: req.FeatureName,
		FeatureCode: req.FeatureCode,
		Icon:        req.Icon,
		RoutePath:   req.RoutePath,
		Component:   req.Component,
		ApiPaths:    req.ApiPaths,
		Permission:  req.Permission,
		Sort:        req.Sort,
		Status:      req.Status,
		IsDefault:   req.IsDefault,
		Badge:       req.Badge,
		Description: req.Description,
		TenantID:    tenantID,
	}

	if item.FeatureCode == "" {
		item.FeatureCode = "feat_" + strconv.FormatUint(uint64(item.ID), 10)
	}

	err := ctrl.DB.Create(&item).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建功能失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "创建成功"})
}

// UpdateFeature 更新功能项
func (ctrl *FeatureConfigController) UpdateFeature(c *gin.Context) {
	id := c.Param("id")
	var req models.FeatureItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误", "error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	updates := map[string]interface{}{
		"feature_name": req.FeatureName,
		"icon":         req.Icon,
		"route_path":   req.RoutePath,
		"component":    req.Component,
		"api_paths":    req.ApiPaths,
		"permission":   req.Permission,
		"sort":         req.Sort,
		"status":       req.Status,
		"is_default":   req.IsDefault,
		"badge":        req.Badge,
		"description":  req.Description,
	}

	if req.GroupID != nil {
		updates["group_id"] = req.GroupID
	}
	if req.FeatureCode != "" {
		updates["feature_code"] = req.FeatureCode
	}

	err := ctrl.DB.Model(&models.FeatureItem{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(updates).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新功能失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// DeleteFeature 删除功能项
func (ctrl *FeatureConfigController) DeleteFeature(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.FeatureItem{}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "删除功能失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// BatchSort 批量排序(支持分组和功能的混合排序)
func (ctrl *FeatureConfigController) BatchSort(c *gin.Context) {
	var req models.FeatureSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "参数错误", "error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	tx := ctrl.DB.Begin()

	for _, item := range req.Items {
		// 更新分组排序
		if item.GroupID == nil {
			err := tx.Model(&models.FeatureGroup{}).
				Where("id = ? AND tenant_id = ?", item.ID, tenantID).
				Update("sort", item.Sort).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新排序失败", "error": err.Error()})
				return
			}
		} else {
			// 更新功能项排序和分组
			updates := map[string]interface{}{
				"sort":     item.Sort,
				"group_id": item.GroupID,
			}
			err := tx.Model(&models.FeatureItem{}).
				Where("id = ? AND tenant_id = ?", item.ID, tenantID).
				Updates(updates).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新排序失败", "error": err.Error()})
				return
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "排序更新成功"})
}

// RegisterRoutes 注册路由
func (ctrl *FeatureConfigController) RegisterRoutes(rg *gin.RouterGroup) {
	groups := rg.Group("/feature-config")
	{
		// 分组管理
		groups.GET("/groups", ctrl.ListGroups)
		groups.GET("/groups/:id", ctrl.GetGroup)
		groups.POST("/groups", ctrl.CreateGroup)
		groups.PUT("/groups/:id", ctrl.UpdateGroup)
		groups.DELETE("/groups/:id", ctrl.DeleteGroup)

		// 功能项管理
		groups.GET("/features", ctrl.ListFeatures)
		groups.GET("/features/:id", ctrl.GetFeature)
		groups.POST("/features", ctrl.CreateFeature)
		groups.PUT("/features/:id", ctrl.UpdateFeature)
		groups.DELETE("/features/:id", ctrl.DeleteFeature)

		// 批量排序
		groups.POST("/sort", ctrl.BatchSort)
	}
}
