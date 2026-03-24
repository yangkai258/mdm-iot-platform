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
func (ctrl *ChildModeController) RegisterRoutes(api *gin.RouterGroup) {
	childMode := api.Group("/child-mode")
	{
		childMode.GET("/settings", ctrl.GetSettings)
		childMode.PUT("/settings", ctrl.UpdateSettings)
		childMode.POST("/settings/enable", ctrl.EnableChildMode)
		childMode.POST("/settings/disable", ctrl.DisableChildMode)
		childMode.GET("/usage/stats", ctrl.GetUsageStats)
	}
}

// GetSettings GET /api/v1/child-mode/settings - 获取儿童模式设置
func (ctrl *ChildModeController) GetSettings(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id不能为空"})
		return
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	var settings models.ChildModeSettings
	if err := ctrl.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认设置
			settings = models.ChildModeSettings{
				UserID:                 uint(userID),
				IsEnabled:              false,
				DailyTimeLimitMinutes:  60,
				SingleSessionMinutes:  30,
				AllowVoiceChat:        true,
				AllowVideoChat:        false,
				AllowAppInstall:       false,
				AllowAppUninstall:     false,
				AllowGameMode:         true,
				AllowSocialFeatures:   false,
				AllowInAppPurchases:   false,
				ContentFilterLevel:    "moderate",
				AllowedStartTime:      "08:00",
				AllowedEndTime:        "20:00",
				RestReminderMinutes:   30,
				WeeklyReportEnabled:   true,
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 计算今日使用时长
	today := time.Now().Format("2006-01-02")
	settings.TodayUsedMinutes = ctrl.getTodayUsageMinutes(uint(userID), today)
	settings.WeekUsedMinutes = ctrl.getWeekUsageMinutes(uint(userID))

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// UpdateSettings PUT /api/v1/child-mode/settings - 更新儿童模式设置
func (ctrl *ChildModeController) UpdateSettings(c *gin.Context) {
	var req UpdateChildModeSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id不能为空"})
		return
	}

	var settings models.ChildModeSettings
	if err := ctrl.DB.Where("user_id = ?", req.UserID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新设置
			settings = models.ChildModeSettings{
				UserID: req.UserID,
			}
			applyUpdateFields(&settings, &req)
			if err := ctrl.DB.Create(&settings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 更新现有设置
	applyUpdateFields(&settings, &req)
	if err := ctrl.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// EnableChildMode POST /api/v1/child-mode/settings/enable - 启用儿童模式
func (ctrl *ChildModeController) EnableChildMode(c *gin.Context) {
	var req EnableChildModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var settings models.ChildModeSettings
	err := ctrl.DB.Where("user_id = ?", req.UserID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		settings = models.ChildModeSettings{UserID: req.UserID, IsEnabled: true}
		if err := ctrl.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "启用失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	} else {
		settings.IsEnabled = true
		if err := ctrl.DB.Save(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "启用失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "儿童模式已启用", "data": settings})
}

// DisableChildMode POST /api/v1/child-mode/settings/disable - 禁用儿童模式
func (ctrl *ChildModeController) DisableChildMode(c *gin.Context) {
	var req DisableChildModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证家长PIN
	var settings models.ChildModeSettings
	if err := ctrl.DB.Where("user_id = ?", req.UserID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "儿童模式未启用"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if settings.ParentPin != "" && settings.ParentPin != req.ParentPin {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "家长PIN码错误"})
		return
	}

	settings.IsEnabled = false
	if err := ctrl.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "禁用失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "儿童模式已禁用"})
}

// GetUsageStats GET /api/v1/child-mode/usage/stats - 获取使用统计
func (ctrl *ChildModeController) GetUsageStats(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id不能为空"})
		return
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	today := time.Now().Format("2006-01-02")

	stats := gin.H{
		"today_used_minutes": ctrl.getTodayUsageMinutes(uint(userID), today),
		"week_used_minutes":  ctrl.getWeekUsageMinutes(uint(userID)),
		"daily_limit_minutes": 60,
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": stats})
}

// ========== 辅助方法 ==========

func applyUpdateFields(settings *models.ChildModeSettings, req *UpdateChildModeSettingsRequest) {
	if req.IsEnabled != nil {
		settings.IsEnabled = *req.IsEnabled
	}
	if req.DailyTimeLimitMinutes != nil {
		settings.DailyTimeLimitMinutes = *req.DailyTimeLimitMinutes
	}
	if req.SingleSessionMinutes != nil {
		settings.SingleSessionMinutes = *req.SingleSessionMinutes
	}
	if req.AllowVoiceChat != nil {
		settings.AllowVoiceChat = *req.AllowVoiceChat
	}
	if req.AllowVideoChat != nil {
		settings.AllowVideoChat = *req.AllowVideoChat
	}
	if req.AllowAppInstall != nil {
		settings.AllowAppInstall = *req.AllowAppInstall
	}
	if req.AllowAppUninstall != nil {
		settings.AllowAppUninstall = *req.AllowAppUninstall
	}
	if req.AllowGameMode != nil {
		settings.AllowGameMode = *req.AllowGameMode
	}
	if req.AllowSocialFeatures != nil {
		settings.AllowSocialFeatures = *req.AllowSocialFeatures
	}
	if req.AllowInAppPurchases != nil {
		settings.AllowInAppPurchases = *req.AllowInAppPurchases
	}
	if req.ContentFilterLevel != "" {
		settings.ContentFilterLevel = req.ContentFilterLevel
	}
	if req.AllowedStartTime != "" {
		settings.AllowedStartTime = req.AllowedStartTime
	}
	if req.AllowedEndTime != "" {
		settings.AllowedEndTime = req.AllowedEndTime
	}
	if req.RestReminderMinutes != nil {
		settings.RestReminderMinutes = *req.RestReminderMinutes
	}
	if req.ParentPin != "" {
		settings.ParentPin = req.ParentPin
	}
	if req.MonitorPetID != nil {
		settings.MonitorPetID = *req.MonitorPetID
	}
	if req.WeeklyReportEnabled != nil {
		settings.WeeklyReportEnabled = *req.WeeklyReportEnabled
	}
}

func (ctrl *ChildModeController) getTodayUsageMinutes(userID uint, today string) int {
	// 占位：实际应根据设备使用日志计算
	return 0
}

func (ctrl *ChildModeController) getWeekUsageMinutes(userID uint) int {
	// 占位：实际应根据设备使用日志计算
	return 0
}

// ========== 请求结构体 ==========

type UpdateChildModeSettingsRequest struct {
	UserID                 uint    `json:"user_id"`
	IsEnabled              *bool   `json:"is_enabled"`
	DailyTimeLimitMinutes  *int    `json:"daily_time_limit_minutes"`
	SingleSessionMinutes   *int    `json:"single_session_minutes"`
	AllowVoiceChat         *bool   `json:"allow_voice_chat"`
	AllowVideoChat         *bool   `json:"allow_video_chat"`
	AllowAppInstall        *bool   `json:"allow_app_install"`
	AllowAppUninstall      *bool   `json:"allow_app_uninstall"`
	AllowGameMode          *bool   `json:"allow_game_mode"`
	AllowSocialFeatures    *bool   `json:"allow_social_features"`
	AllowInAppPurchases    *bool   `json:"allow_in_app_purchases"`
	ContentFilterLevel     string  `json:"content_filter_level"`
	AllowedStartTime       string  `json:"allowed_start_time"`
	AllowedEndTime         string  `json:"allowed_end_time"`
	RestReminderMinutes    *int    `json:"rest_reminder_minutes"`
	ParentPin              string  `json:"parent_pin"`
	MonitorPetID           *uint   `json:"monitor_pet_id"`
	WeeklyReportEnabled     *bool   `json:"weekly_report_enabled"`
}

type EnableChildModeRequest struct {
	UserID uint `json:"user_id"`
}

type DisableChildModeRequest struct {
	UserID    uint   `json:"user_id"`
	ParentPin string `json:"parent_pin"`
}
