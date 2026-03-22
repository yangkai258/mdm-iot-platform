package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FamilyModeController 家庭模式控制器（儿童模式 + 老人陪伴模式）
type FamilyModeController struct {
	DB *gorm.DB
}

// RegisterFamilyModeRoutes 注册家庭模式路由
func (c *FamilyModeController) RegisterFamilyModeRoutes(api *gin.RouterGroup) {
	// ============ 儿童模式路由 ============
	children := api.Group("/family/children")
	{
		// 儿童档案管理
		children.POST("/profile", c.CreateChildrenProfile)
		children.GET("/profiles", c.ListChildrenProfiles)
		children.PUT("/profiles/:id", c.UpdateChildrenProfile)
		children.DELETE("/profiles/:id", c.DeleteChildrenProfile)

		// 儿童内容过滤设置
		children.GET("/:id/content-filter", c.GetContentFilter)
		children.PUT("/:id/content-filter", c.UpdateContentFilter)

		// 儿童使用限制
		children.GET("/:id/usage-limits", c.GetUsageLimit)
		children.PUT("/:id/usage-limits", c.UpdateUsageLimit)
	}

	// ============ 老人陪伴模式路由 ============
	elderly := api.Group("/family/elderly")
	{
		// 老人档案管理
		elderly.POST("/profile", c.CreateElderlyProfile)
		elderly.GET("/profiles", c.ListElderlyProfiles)
		elderly.PUT("/profiles/:id", c.UpdateElderlyProfile)

		// 健康监控设置
		elderly.GET("/:id/health-monitor", c.GetHealthMonitor)
		elderly.PUT("/:id/health-monitor", c.UpdateHealthMonitor)

		// 提醒管理
		elderly.GET("/:id/reminders", c.ListElderlyReminders)
		elderly.POST("/:id/reminders", c.CreateElderlyReminder)
	}
}

// ============ 儿童模式 API ============

// CreateChildrenProfile 创建儿童档案
func (c *FamilyModeController) CreateChildrenProfile(ctx *gin.Context) {
	var req struct {
		DeviceID     string  `json:"device_id"`
		Name         string  `json:"name" binding:"required"`
		Nickname     string  `json:"nickname"`
		Avatar       string  `json:"avatar"`
		BirthDate    *string `json:"birth_date"`
		Gender       string  `json:"gender"`
		ParentUserID string  `json:"parent_user_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	profile := models.ChildrenProfile{
		Name:         req.Name,
		Nickname:     req.Nickname,
		Avatar:       req.Avatar,
		Gender:       req.Gender,
		DeviceID:     req.DeviceID,
		ParentUserID: req.ParentUserID,
	}

	if req.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			profile.BirthDate = &t
			profile.Age = calculateAge(t)
		}
	}

	if err := c.DB.Create(&profile).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	// 自动创建默认的内容过滤规则
	filterRule := models.ContentFilterRule{
		ProfileID:     profile.ProfileID,
		RuleName:      "默认过滤规则",
		FilterLevel:   2,
		BlockAdult:    true,
		BlockViolence: true,
		BlockGambling: true,
		BlockAds:      true,
		Enabled:       true,
	}
	c.DB.Create(&filterRule)

	// 自动创建默认的使用限制
	usageLimit := models.UsageLimit{
		ProfileID: profile.ProfileID,
		Enabled:   true,
	}
	c.DB.Create(&usageLimit)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": profile})
}

// ListChildrenProfiles 儿童列表
func (c *FamilyModeController) ListChildrenProfiles(ctx *gin.Context) {
	var profiles []models.ChildrenProfile
	var total int64

	query := c.DB.Model(&models.ChildrenProfile{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if parentUserID := ctx.Query("parent_user_id"); parentUserID != "" {
		query = query.Where("parent_user_id = ?", parentUserID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&profiles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      profiles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// UpdateChildrenProfile 更新儿童档案
func (c *FamilyModeController) UpdateChildrenProfile(ctx *gin.Context) {
	id := ctx.Param("id")

	var profile models.ChildrenProfile
	if err := c.DB.Where("profile_id = ?", id).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "档案不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Name      string  `json:"name"`
		Nickname  string  `json:"nickname"`
		Avatar    string  `json:"avatar"`
		BirthDate *string `json:"birth_date"`
		Gender    string  `json:"gender"`
		DeviceID  string  `json:"device_id"`
		Status    *int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.DeviceID != "" {
		updates["device_id"] = req.DeviceID
	}
	if req.BirthDate != nil {
		if t, err := time.Parse("2006-01-02", *req.BirthDate); err == nil {
			updates["birth_date"] = t
			updates["age"] = calculateAge(t)
		}
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := c.DB.Model(&profile).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.Where("profile_id = ?", id).First(&profile)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": profile})
}

// DeleteChildrenProfile 删除儿童档案
func (c *FamilyModeController) DeleteChildrenProfile(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Where("profile_id = ?", id).Delete(&models.ChildrenProfile{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "档案不存在"})
		return
	}

	// 同时删除关联的内容过滤规则和使用限制
	c.DB.Where("profile_id = ?", id).Delete(&models.ContentFilterRule{})
	c.DB.Where("profile_id = ?", id).Delete(&models.UsageLimit{})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetContentFilter 获取内容过滤设置
func (c *FamilyModeController) GetContentFilter(ctx *gin.Context) {
	id := ctx.Param("id")

	var filter models.ContentFilterRule
	if err := c.DB.Where("profile_id = ?", id).First(&filter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，返回空规则
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
				"profile_id": id,
				"rule_id":    "",
				"enabled":    false,
			}})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": filter})
}

// UpdateContentFilter 更新内容过滤规则
func (c *FamilyModeController) UpdateContentFilter(ctx *gin.Context) {
	id := ctx.Param("id")

	var filter models.ContentFilterRule
	if err := c.DB.Where("profile_id = ?", id).First(&filter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在则创建
			var req models.ContentFilterRule
			if err := ctx.ShouldBindJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
				return
			}
			req.ProfileID = id
			if err := c.DB.Create(&req).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		RuleName          string   `json:"rule_name"`
		FilterLevel       *int     `json:"filter_level"`
		BlockAdult        *bool    `json:"block_adult"`
		BlockViolence     *bool    `json:"block_violence"`
		BlockGambling     *bool    `json:"block_gambling"`
		BlockAds          *bool    `json:"block_ads"`
		BlockGames        *bool    `json:"block_games"`
		AllowedCategories []string `json:"allowed_categories"`
		BlockedKeywords   []string `json:"blocked_keywords"`
		AllowedApps       []string `json:"allowed_apps"`
		BlockedApps       []string `json:"blocked_apps"`
		WhitelistMode     *bool    `json:"whitelist_mode"`
		Enabled           *bool    `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.RuleName != "" {
		updates["rule_name"] = req.RuleName
	}
	if req.FilterLevel != nil {
		updates["filter_level"] = *req.FilterLevel
	}
	if req.BlockAdult != nil {
		updates["block_adult"] = *req.BlockAdult
	}
	if req.BlockViolence != nil {
		updates["block_violence"] = *req.BlockViolence
	}
	if req.BlockGambling != nil {
		updates["block_gambling"] = *req.BlockGambling
	}
	if req.BlockAds != nil {
		updates["block_ads"] = *req.BlockAds
	}
	if req.BlockGames != nil {
		updates["block_games"] = *req.BlockGames
	}
	if req.AllowedCategories != nil {
		updates["allowed_categories"] = req.AllowedCategories
	}
	if req.BlockedKeywords != nil {
		updates["blocked_keywords"] = req.BlockedKeywords
	}
	if req.AllowedApps != nil {
		updates["allowed_apps"] = req.AllowedApps
	}
	if req.BlockedApps != nil {
		updates["blocked_apps"] = req.BlockedApps
	}
	if req.WhitelistMode != nil {
		updates["whitelist_mode"] = *req.WhitelistMode
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := c.DB.Model(&filter).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.Where("profile_id = ?", id).First(&filter)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": filter})
}

// GetUsageLimit 获取使用限制
func (c *FamilyModeController) GetUsageLimit(ctx *gin.Context) {
	id := ctx.Param("id")

	var limit models.UsageLimit
	if err := c.DB.Where("profile_id = ?", id).First(&limit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
				"profile_id": id,
				"limit_id":   "",
				"enabled":    false,
			}})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": limit})
}

// UpdateUsageLimit 更新使用限制
func (c *FamilyModeController) UpdateUsageLimit(ctx *gin.Context) {
	id := ctx.Param("id")

	var limit models.UsageLimit
	if err := c.DB.Where("profile_id = ?", id).First(&limit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var req models.UsageLimit
			if err := ctx.ShouldBindJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
				return
			}
			req.ProfileID = id
			if err := c.DB.Create(&req).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		DailyTimeLimit         *int    `json:"daily_time_limit"`
		WeeklyTimeLimit        *int    `json:"weekly_time_limit"`
		MaxSingleSession       *int    `json:"max_single_session"`
		AllowedStartTime       *string `json:"allowed_start_time"`
		AllowedEndTime         *string `json:"allowed_end_time"`
		AllowedDays            []int   `json:"allowed_days"`
		BreakInterval          *int    `json:"break_interval"`
		BreakDuration          *int    `json:"break_duration"`
		EyeProtectionEnabled   *bool   `json:"eye_protection_enabled"`
		EyeProtectionInterval  *int    `json:"eye_protection_interval"`
		PostureReminderEnabled *bool   `json:"posture_reminder_enabled"`
		Enabled                *bool   `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.DailyTimeLimit != nil {
		updates["daily_time_limit"] = *req.DailyTimeLimit
	}
	if req.WeeklyTimeLimit != nil {
		updates["weekly_time_limit"] = *req.WeeklyTimeLimit
	}
	if req.MaxSingleSession != nil {
		updates["max_single_session"] = *req.MaxSingleSession
	}
	if req.AllowedStartTime != nil {
		updates["allowed_start_time"] = *req.AllowedStartTime
	}
	if req.AllowedEndTime != nil {
		updates["allowed_end_time"] = *req.AllowedEndTime
	}
	if req.AllowedDays != nil {
		updates["allowed_days"] = req.AllowedDays
	}
	if req.BreakInterval != nil {
		updates["break_interval"] = *req.BreakInterval
	}
	if req.BreakDuration != nil {
		updates["break_duration"] = *req.BreakDuration
	}
	if req.EyeProtectionEnabled != nil {
		updates["eye_protection_enabled"] = *req.EyeProtectionEnabled
	}
	if req.EyeProtectionInterval != nil {
		updates["eye_protection_interval"] = *req.EyeProtectionInterval
	}
	if req.PostureReminderEnabled != nil {
		updates["posture_reminder_enabled"] = *req.PostureReminderEnabled
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := c.DB.Model(&limit).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.Where("profile_id = ?", id).First(&limit)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": limit})
}

// ============ 老人陪伴模式 API ============

// CreateElderlyProfile 创建老人档案
func (c *FamilyModeController) CreateElderlyProfile(ctx *gin.Context) {
	var req struct {
		DeviceID         string  `json:"device_id"`
		Name             string  `json:"name" binding:"required"`
		Nickname         string  `json:"nickname"`
		Avatar           string  `json:"avatar"`
		BirthDate        *string `json:"birth_date"`
		Gender           string  `json:"gender"`
		Phone            string  `json:"phone"`
		EmergencyContact string  `json:"emergency_contact"`
		EmergencyPhone   string  `json:"emergency_phone"`
		Address          string  `json:"address"`
		MedicalHistory   string  `json:"medical_history"`
		Allergies        string  `json:"allergies"`
		BloodType        string  `json:"blood_type"`
		CaregiverUserID  string  `json:"caregiver_user_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	profile := models.ElderlyProfile{
		Name:             req.Name,
		Nickname:         req.Nickname,
		Avatar:           req.Avatar,
		Gender:           req.Gender,
		Phone:            req.Phone,
		EmergencyContact: req.EmergencyContact,
		EmergencyPhone:   req.EmergencyPhone,
		Address:          req.Address,
		MedicalHistory:   req.MedicalHistory,
		Allergies:        req.Allergies,
		BloodType:        req.BloodType,
		DeviceID:         req.DeviceID,
		CaregiverUserID:  req.CaregiverUserID,
	}

	if req.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err == nil {
			profile.BirthDate = &t
			profile.Age = calculateAge(t)
		}
	}

	if err := c.DB.Create(&profile).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	// 自动创建默认的健康监控设置
	healthSetting := models.HealthMonitorSetting{
		ProfileID: profile.ProfileID,
		Enabled:   true,
	}
	c.DB.Create(&healthSetting)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": profile})
}

// ListElderlyProfiles 老人列表
func (c *FamilyModeController) ListElderlyProfiles(ctx *gin.Context) {
	var profiles []models.ElderlyProfile
	var total int64

	query := c.DB.Model(&models.ElderlyProfile{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR nickname LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if caregiverUserID := ctx.Query("caregiver_user_id"); caregiverUserID != "" {
		query = query.Where("caregiver_user_id = ?", caregiverUserID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&profiles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      profiles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// UpdateElderlyProfile 更新老人档案
func (c *FamilyModeController) UpdateElderlyProfile(ctx *gin.Context) {
	id := ctx.Param("id")

	var profile models.ElderlyProfile
	if err := c.DB.Where("profile_id = ?", id).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "档案不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Name             string  `json:"name"`
		Nickname         string  `json:"nickname"`
		Avatar           string  `json:"avatar"`
		BirthDate        *string `json:"birth_date"`
		Gender           string  `json:"gender"`
		Phone            string  `json:"phone"`
		EmergencyContact string  `json:"emergency_contact"`
		EmergencyPhone   string  `json:"emergency_phone"`
		Address          string  `json:"address"`
		MedicalHistory   string  `json:"medical_history"`
		Allergies        string  `json:"allergies"`
		BloodType        string  `json:"blood_type"`
		DeviceID         string  `json:"device_id"`
		Status           *int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.EmergencyContact != "" {
		updates["emergency_contact"] = req.EmergencyContact
	}
	if req.EmergencyPhone != "" {
		updates["emergency_phone"] = req.EmergencyPhone
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.MedicalHistory != "" {
		updates["medical_history"] = req.MedicalHistory
	}
	if req.Allergies != "" {
		updates["allergies"] = req.Allergies
	}
	if req.BloodType != "" {
		updates["blood_type"] = req.BloodType
	}
	if req.DeviceID != "" {
		updates["device_id"] = req.DeviceID
	}
	if req.BirthDate != nil {
		if t, err := time.Parse("2006-01-02", *req.BirthDate); err == nil {
			updates["birth_date"] = t
			updates["age"] = calculateAge(t)
		}
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := c.DB.Model(&profile).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.Where("profile_id = ?", id).First(&profile)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": profile})
}

// GetHealthMonitor 获取健康监控设置
func (c *FamilyModeController) GetHealthMonitor(ctx *gin.Context) {
	id := ctx.Param("id")

	var setting models.HealthMonitorSetting
	if err := c.DB.Where("profile_id = ?", id).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
				"profile_id": id,
				"setting_id": "",
				"enabled":    false,
			}})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": setting})
}

// UpdateHealthMonitor 更新健康监控设置
func (c *FamilyModeController) UpdateHealthMonitor(ctx *gin.Context) {
	id := ctx.Param("id")

	var setting models.HealthMonitorSetting
	if err := c.DB.Where("profile_id = ?", id).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var req models.HealthMonitorSetting
			if err := ctx.ShouldBindJSON(&req); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
				return
			}
			req.ProfileID = id
			if err := c.DB.Create(&req).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		HeartRateMonitor      *bool   `json:"heart_rate_monitor"`
		HeartRateMin          *int    `json:"heart_rate_min"`
		HeartRateMax          *int    `json:"heart_rate_max"`
		BPMonitor             *bool   `json:"bp_monitor"`
		BPSystolicMin         *int    `json:"bp_systolic_min"`
		BPSystolicMax         *int    `json:"bp_systolic_max"`
		BPDiastolicMin        *int    `json:"bp_diastolic_min"`
		BPDiastolicMax        *int    `json:"bp_diastolic_max"`
		SleepMonitor          *bool   `json:"sleep_monitor"`
		ActivityMonitor       *bool   `json:"activity_monitor"`
		StepGoal              *int    `json:"step_goal"`
		FallDetection         *bool   `json:"fall_detection"`
		FallAlertEnabled      *bool   `json:"fall_alert_enabled"`
		EmergencyAlertEnabled *bool   `json:"emergency_alert_enabled"`
		AlertThreshold        *int    `json:"alert_threshold"`
		CheckInterval         *int    `json:"check_interval"`
		ReportFrequency       *string `json:"report_frequency"`
		Enabled               *bool   `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.HeartRateMonitor != nil {
		updates["heart_rate_monitor"] = *req.HeartRateMonitor
	}
	if req.HeartRateMin != nil {
		updates["heart_rate_min"] = *req.HeartRateMin
	}
	if req.HeartRateMax != nil {
		updates["heart_rate_max"] = *req.HeartRateMax
	}
	if req.BPMonitor != nil {
		updates["bp_monitor"] = *req.BPMonitor
	}
	if req.BPSystolicMin != nil {
		updates["bp_systolic_min"] = *req.BPSystolicMin
	}
	if req.BPSystolicMax != nil {
		updates["bp_systolic_max"] = *req.BPSystolicMax
	}
	if req.BPDiastolicMin != nil {
		updates["bp_diastolic_min"] = *req.BPDiastolicMin
	}
	if req.BPDiastolicMax != nil {
		updates["bp_diastolic_max"] = *req.BPDiastolicMax
	}
	if req.SleepMonitor != nil {
		updates["sleep_monitor"] = *req.SleepMonitor
	}
	if req.ActivityMonitor != nil {
		updates["activity_monitor"] = *req.ActivityMonitor
	}
	if req.StepGoal != nil {
		updates["step_goal"] = *req.StepGoal
	}
	if req.FallDetection != nil {
		updates["fall_detection"] = *req.FallDetection
	}
	if req.FallAlertEnabled != nil {
		updates["fall_alert_enabled"] = *req.FallAlertEnabled
	}
	if req.EmergencyAlertEnabled != nil {
		updates["emergency_alert_enabled"] = *req.EmergencyAlertEnabled
	}
	if req.AlertThreshold != nil {
		updates["alert_threshold"] = *req.AlertThreshold
	}
	if req.CheckInterval != nil {
		updates["check_interval"] = *req.CheckInterval
	}
	if req.ReportFrequency != nil {
		updates["report_frequency"] = *req.ReportFrequency
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := c.DB.Model(&setting).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.Where("profile_id = ?", id).First(&setting)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": setting})
}

// ListElderlyReminders 获取提醒列表
func (c *FamilyModeController) ListElderlyReminders(ctx *gin.Context) {
	id := ctx.Param("id")
	var reminders []models.ElderlyReminder
	var total int64

	query := c.DB.Model(&models.ElderlyReminder{}).Where("profile_id = ?", id)

	if reminderType := ctx.Query("reminder_type"); reminderType != "" {
		query = query.Where("reminder_type = ?", reminderType)
	}
	if isEnabled := ctx.Query("is_enabled"); isEnabled != "" {
		query = query.Where("is_enabled = ?", isEnabled)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&reminders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      reminders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// CreateElderlyReminder 创建提醒
func (c *FamilyModeController) CreateElderlyReminder(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		Title          string `json:"title" binding:"required"`
		Content        string `json:"content"`
		ReminderType   string `json:"reminder_type" binding:"required"`
		MedicineName   string `json:"medicine_name"`
		MedicineDosage string `json:"medicine_dosage"`
		ScheduleType   string `json:"schedule_type"`
		ScheduleTime   string `json:"schedule_time" binding:"required"`
		ScheduleDays   []int  `json:"schedule_days"`
		ScheduleDates  []int  `json:"schedule_dates"`
		AdvanceNotice  int    `json:"advance_notice"`
		RepeatCount    int    `json:"repeat_count"`
		RepeatInterval int    `json:"repeat_interval"`
		IsEnabled      *bool  `json:"is_enabled"`
		DeviceID       string `json:"device_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	reminder := models.ElderlyReminder{
		ProfileID:      id,
		Title:          req.Title,
		Content:        req.Content,
		ReminderType:   req.ReminderType,
		MedicineName:   req.MedicineName,
		MedicineDosage: req.MedicineDosage,
		ScheduleType:   req.ScheduleType,
		ScheduleTime:   req.ScheduleTime,
		ScheduleDays:   req.ScheduleDays,
		ScheduleDates:  req.ScheduleDates,
		AdvanceNotice:  req.AdvanceNotice,
		RepeatCount:    req.RepeatCount,
		RepeatInterval: req.RepeatInterval,
		DeviceID:       req.DeviceID,
		IsEnabled:      true,
	}

	if req.IsEnabled != nil {
		reminder.IsEnabled = *req.IsEnabled
	}

	if err := c.DB.Create(&reminder).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": reminder})
}

// calculateAge 根据出生日期计算年龄
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
