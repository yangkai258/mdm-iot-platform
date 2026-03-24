package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ChildModeController 儿童模式控制器
type ChildModeController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册儿童模式路由
func (c *ChildModeController) RegisterRoutes(api *gin.RouterGroup) {
	childMode := api.Group("/child-mode")
	{
		childMode.GET("/settings", c.ListSettings)
		childMode.GET("/settings/:id", c.GetSettings)
		childMode.POST("/settings", c.CreateSettings)
		childMode.PUT("/settings/:id", c.UpdateSettings)
		childMode.DELETE("/settings/:id", c.DeleteSettings)
		childMode.POST("/settings/:id/enable", c.EnableChildMode)
		childMode.POST("/settings/:id/disable", c.DisableChildMode)
		childMode.GET("/settings/:id/statistics", c.GetStatistics)
		childMode.POST("/settings/:id/reset-timer", c.ResetTimer)
	}
}

// ListSettings 获取儿童模式设置列表
func (c *ChildModeController) ListSettings(ctx *gin.Context) {
	var settings []models.ChildModeSettings
	var total int64

	query := c.DB.Model(&models.ChildModeSettings{})

	// 用户筛选
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 设备筛选
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 是否启用筛选
	if enabled := ctx.Query("is_enabled"); enabled != "" {
		query = query.Where("is_enabled = ?", enabled)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&settings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      settings,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSettings 获取儿童模式设置详情
func (c *ChildModeController) GetSettings(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": settings})
}

// CreateSettings 创建设置
func (c *ChildModeController) CreateSettings(ctx *gin.Context) {
	var settings models.ChildModeSettings
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	settings.Status = "active"
	if err := c.DB.Create(&settings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": settings})
}

// UpdateSettings 更新设置
func (c *ChildModeController) UpdateSettings(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		ContentFilterLevel  string `json:"content_filter_level"`
		AllowedCategories   string `json:"allowed_categories"`
		BlockedKeywords     string `json:"blocked_keywords"`
		BlockedApps         string `json:"blocked_apps"`
		DailyTimeLimit      int    `json:"daily_time_limit"`
		SessionDuration     int    `json:"session_duration"`
		BreakDuration       int    `json:"break_duration"`
		AllowedStartTime    string `json:"allowed_start_time"`
		AllowedEndTime      string `json:"allowed_end_time"`
		EmergencyContact    string `json:"emergency_contact"`
		PinCode             string `json:"pin_code"`
		ParentPhone         string `json:"parent_phone"`
		CanShareContent     *bool  `json:"can_share_content"`
		CanDownloadContent  *bool  `json:"can_download_content"`
		AllowCamera         *bool  `json:"allow_camera"`
		AllowMicrophone     *bool  `json:"allow_microphone"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.ContentFilterLevel != "" {
		updates["content_filter_level"] = updateData.ContentFilterLevel
	}
	if updateData.AllowedCategories != "" {
		updates["allowed_categories"] = updateData.AllowedCategories
	}
	if updateData.BlockedKeywords != "" {
		updates["blocked_keywords"] = updateData.BlockedKeywords
	}
	if updateData.BlockedApps != "" {
		updates["blocked_apps"] = updateData.BlockedApps
	}
	if updateData.DailyTimeLimit > 0 {
		updates["daily_time_limit"] = updateData.DailyTimeLimit
	}
	if updateData.SessionDuration > 0 {
		updates["session_duration"] = updateData.SessionDuration
	}
	if updateData.BreakDuration > 0 {
		updates["break_duration"] = updateData.BreakDuration
	}
	if updateData.AllowedStartTime != "" {
		updates["allowed_start_time"] = updateData.AllowedStartTime
	}
	if updateData.AllowedEndTime != "" {
		updates["allowed_end_time"] = updateData.AllowedEndTime
	}
	if updateData.EmergencyContact != "" {
		updates["emergency_contact"] = updateData.EmergencyContact
	}
	if updateData.PinCode != "" {
		updates["pin_code"] = updateData.PinCode
	}
	if updateData.ParentPhone != "" {
		updates["parent_phone"] = updateData.ParentPhone
	}
	if updateData.CanShareContent != nil {
		updates["can_share_content"] = *updateData.CanShareContent
	}
	if updateData.CanDownloadContent != nil {
		updates["can_download_content"] = *updateData.CanDownloadContent
	}
	if updateData.AllowCamera != nil {
		updates["allow_camera"] = *updateData.AllowCamera
	}
	if updateData.AllowMicrophone != nil {
		updates["allow_microphone"] = *updateData.AllowMicrophone
	}

	if err := c.DB.Model(&settings).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&settings, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": settings})
}

// DeleteSettings 删除设置
func (c *ChildModeController) DeleteSettings(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	if err := c.DB.Delete(&settings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// EnableChildMode 启用儿童模式
func (c *ChildModeController) EnableChildMode(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	if err := c.DB.Model(&settings).Updates(map[string]interface{}{
		"is_enabled": true,
		"status":     "active",
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "启用失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// DisableChildMode 禁用儿童模式
func (c *ChildModeController) DisableChildMode(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 验证PIN码
	var pinData struct {
		PinCode string `json:"pin_code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&pinData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供PIN码"})
		return
	}

	if settings.PinCode != "" && settings.PinCode != pinData.PinCode {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "PIN码错误"})
		return
	}

	if err := c.DB.Model(&settings).Updates(map[string]interface{}{
		"is_enabled": false,
		"status":     "paused",
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "禁用失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetStatistics 获取使用统计
func (c *ChildModeController) GetStatistics(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"today_used_minutes": settings.TodayUsedMinutes,
			"week_used_minutes":  settings.WeekUsedMinutes,
			"total_used_minutes": settings.TotalUsedMinutes,
			"total_sessions":      settings.TotalSessions,
			"last_session_at":    settings.LastSessionAt,
			"daily_time_limit":   settings.DailyTimeLimit,
			"remaining_minutes":  settings.DailyTimeLimit - settings.TodayUsedMinutes,
		},
	})
}

// ResetTimer 重置使用计时
func (c *ChildModeController) ResetTimer(ctx *gin.Context) {
	id := ctx.Param("id")
	var settings models.ChildModeSettings
	if err := c.DB.First(&settings, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 验证PIN码
	var pinData struct {
		PinCode string `json:"pin_code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&pinData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供PIN码"})
		return
	}

	if settings.PinCode != "" && settings.PinCode != pinData.PinCode {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "PIN码错误"})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&settings).Updates(map[string]interface{}{
		"today_used_minutes": 0,
		"last_session_at":    &now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "重置失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
