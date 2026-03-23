package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdvancedController 高级功能控制器（儿童模式/老人陪伴/家庭相册/宠物服务）
type AdvancedController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterAdvancedRoutes 注册高级功能路由
func (ctrl *AdvancedController) RegisterAdvancedRoutes(api *gin.RouterGroup) {
	advanced := api.Group("/advanced")
	ac := &AdvancedController{DB: ctrl.DB, Redis: ctrl.Redis}

	// ============ 儿童模式 ============
	childMode := advanced.Group("/child-mode")
	childMode.Use() // 无额外中间件
	childMode.GET("/config", ac.GetChildModeConfig)
	childMode.PUT("/config", ac.UpdateChildModeConfig)
	childMode.GET("/stats", ac.GetChildModeStats)

	// ============ 老人陪伴 ============
	elderly := advanced.Group("/elderly")
	elderly.GET("/config", ac.GetElderlyConfig)
	elderly.PUT("/config", ac.UpdateElderlyConfig)
	elderly.GET("/health", ac.GetElderlyHealth)

	// ============ 家庭相册 ============
	album := advanced.Group("/album")
	album.GET("/photos", ac.ListPhotos)
	album.POST("/photos", ac.UploadPhoto)
	album.DELETE("/photos/:id", ac.DeletePhoto)
	album.GET("/share/:id", ac.SharePhoto)
	album.GET("/photos/:id/comments", ac.ListPhotoComments)
	album.POST("/photos/:id/comments", ac.AddPhotoComment)
	album.POST("/photos/:id/like", ac.LikePhoto)
	album.DELETE("/photos/:id/like", ac.UnlikePhoto)

	// ============ 寻回网络 ============
	petFinder := advanced.Group("/pet-finder")
	petFinder.GET("/reports", ac.ListPetFinderReports)
	petFinder.POST("/reports", ac.CreatePetFinderReport)
	petFinder.GET("/reports/:id", ac.GetPetFinderReport)
	petFinder.PUT("/reports/:id", ac.UpdatePetFinderReport)
	petFinder.DELETE("/reports/:id", ac.ClosePetFinderReport)
	petFinder.POST("/reports/:id/sighting", ac.AddSighting)
	petFinder.GET("/reports/:id/sightings", ac.ListSightings)
	petFinder.GET("/nearby", ac.GetNearbyPets)

	// ============ 宠物疫苗接种 ============
	vaccination := advanced.Group("/vaccination")
	vaccination.GET("/records", ac.ListVaccinations)
	vaccination.POST("/records", ac.CreateVaccination)
	vaccination.PUT("/records/:id", ac.UpdateVaccination)
	vaccination.DELETE("/records/:id", ac.DeleteVaccination)
	vaccination.GET("/records/:id", ac.GetVaccination)

	// ============ 宠物饮食记录 ============
	diet := advanced.Group("/diet")
	diet.GET("/records", ac.ListDietRecords)
	diet.POST("/records", ac.CreateDietRecord)
	diet.PUT("/records/:id", ac.UpdateDietRecord)
	diet.DELETE("/records/:id", ac.DeleteDietRecord)
	diet.GET("/summary", ac.GetDietSummary)
}

// ==================== 儿童模式 API ====================

// GetChildModeConfig 获取儿童模式配置
func (ac *AdvancedController) GetChildModeConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ChildModeConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err == gorm.ErrRecordNotFound {
		// 返回默认配置
		config = models.ChildModeConfig{
			UserID:              userID,
			DeviceID:            deviceID,
			ContentFilterLevel:  "moderate",
			SessionDuration:     30,
			BreakDuration:       10,
			AllowedStartTime:    "08:00",
			AllowedEndTime:      "20:00",
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// UpdateChildModeConfig 更新儿童模式配置
func (ac *AdvancedController) UpdateChildModeConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.ChildModeConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var config models.ChildModeConfig
	query := ac.DB.Where("user_id = ?", userID)
	if req.DeviceID != "" {
		query = query.Where("device_id = ?", req.DeviceID)
	}

	if err := query.First(&config).Error; err == gorm.ErrRecordNotFound {
		// 创建新配置
		config = req
		config.UserID = userID
		if err := ac.DB.Create(&config).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建配置失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询配置失败"})
		return
	} else {
		// 更新现有配置（保留ID和UserID）
		req.ID = config.ID
		req.UserID = userID
		if err := ac.DB.Save(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新配置失败"})
			return
		}
		config = req
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// GetChildModeStats 获取儿童模式使用统计
func (ac *AdvancedController) GetChildModeStats(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ChildModeConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到配置"})
		return
	}

	// 补充本周每日使用数据
	var dailyStats []map[string]interface{}
	weekStart := time.Now().AddDate(0, 0, -6)
	for i := 0; i < 7; i++ {
		day := weekStart.AddDate(0, 0, i)
		dailyStats = append(dailyStats, map[string]interface{}{
			"date":       day.Format("2006-01-02"),
			"used_minutes": 0, // 实际应从使用记录表统计
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success", "data": gin.H{
			"today_used_minutes": config.TodayUsedMinutes,
			"week_used_minutes":  config.WeekUsedMinutes,
			"total_used_minutes": config.TotalUsedMinutes,
			"total_sessions":     config.TotalSessions,
			"last_session_at":    config.LastSessionAt,
			"daily_stats":        dailyStats,
		},
	})
}

// ==================== 老人陪伴 API ====================

// GetElderlyConfig 获取老人陪伴配置
func (ac *AdvancedController) GetElderlyConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ElderlyCareConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err == gorm.ErrRecordNotFound {
		config = models.ElderlyCareConfig{
			UserID:                userID,
			DeviceID:              deviceID,
			HealthMonitorEnabled:  true,
			HeartRateAlertHigh:    100,
			HeartRateAlertLow:     50,
			ActivityGoal:         6000,
			SleepMonitoring:       true,
			MedicationReminders:   true,
			CompanionEnabled:      true,
			InteractionFrequency:  "normal",
			VoiceCallEnabled:      true,
			FallDetectionEnabled:  true,
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// UpdateElderlyConfig 更新老人陪伴配置
func (ac *AdvancedController) UpdateElderlyConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.ElderlyCareConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var config models.ElderlyCareConfig
	query := ac.DB.Where("user_id = ?", userID)
	if req.DeviceID != "" {
		query = query.Where("device_id = ?", req.DeviceID)
	}

	if err := query.First(&config).Error; err == gorm.ErrRecordNotFound {
		req.UserID = userID
		if err := ac.DB.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建配置失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询配置失败"})
		return
	} else {
		req.ID = config.ID
		req.UserID = userID
		if err := ac.DB.Save(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新配置失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// GetElderlyHealth 获取老人健康数据
func (ac *AdvancedController) GetElderlyHealth(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ElderlyCareConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到配置"})
		return
	}

	// 构造健康数据（实际应从健康数据表查询）
	healthData := gin.H{
		"heart_rate":         72, // 实际从设备或健康数据表获取
		"heart_rate_high":   config.HeartRateAlertHigh,
		"heart_rate_low":    config.HeartRateAlertLow,
		"activity_steps":    4500,
		"activity_goal":     config.ActivityGoal,
		"sleep_hours":       7.5,
		"last_interaction":   config.LastInteractionAt,
		"total_interactions": config.TotalInteractions,
		"fall_detected":     false,
		"medication_taken":  true,
		"updated_at":        time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": healthData})
}

// ==================== 家庭相册 API ====================

// ListPhotos 获取照片列表
func (ac *AdvancedController) ListPhotos(c *gin.Context) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	deviceID := c.Query("device_id")
	petID := c.Query("pet_id")
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := ac.DB.Model(&models.FamilyAlbum{}).Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Count(&total)

	var photos []models.FamilyAlbum
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&photos)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success", "data": gin.H{
			"list": photos,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// UploadPhoto 上传照片
func (ac *AdvancedController) UploadPhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		DeviceID    string  `json:"device_id"`
		PetID       uint    `json:"pet_id"`
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description"`
		PhotoURL    string  `json:"photo_url" binding:"required"`
		ThumbnailURL string `json:"thumbnail_url"`
		FileSize    int64   `json:"file_size"`
		Width       int     `json:"width"`
		Height      int     `json:"height"`
		MimeType    string  `json:"mime_type"`
		Category    string  `json:"category"`
		Tags        string  `json:"tags"`
		TakenLat    float64 `json:"taken_lat"`
		TakenLng    float64 `json:"taken_lng"`
		TakenAddr   string  `json:"taken_addr"`
		TakenAt     string  `json:"taken_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var takenAt *time.Time
	if req.TakenAt != "" {
		t := parseTime(req.TakenAt)
		takenAt = &t
	}

	photo := models.FamilyAlbum{
		UserID:       userID,
		DeviceID:     req.DeviceID,
		PetID:        req.PetID,
		Title:        req.Title,
		Description:  req.Description,
		PhotoURL:     req.PhotoURL,
		ThumbnailURL: req.ThumbnailURL,
		FileSize:     req.FileSize,
		Width:        req.Width,
		Height:       req.Height,
		MimeType:     req.MimeType,
		Category:     req.Category,
		Tags:         req.Tags,
		TakenLat:     req.TakenLat,
		TakenLng:     req.TakenLng,
		TakenAddr:    req.TakenAddr,
		TakenAt:      takenAt,
	}

	if err := ac.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "上传失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": photo})
}

// DeletePhoto 删除照片
func (ac *AdvancedController) DeletePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	result := ac.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.FamilyAlbum{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "照片不存在或无权删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// SharePhoto 分享照片
func (ac *AdvancedController) SharePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var photo models.FamilyAlbum
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "照片不存在"})
		return
	}

	// 生成或更新分享令牌
	if photo.ShareToken == "" {
		photo.ShareToken = fmt.Sprintf("share_%d_%d", photo.ID, time.Now().UnixNano())
	}
	photo.IsShared = true
	// 默认7天有效期
	expiry := time.Now().Add(7 * 24 * time.Hour)
	photo.ShareExpiry = &expiry

	if err := ac.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "分享失败"})
		return
	}

	shareURL := fmt.Sprintf("/api/v1/advanced/album/shared/%s", photo.ShareToken)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"share_url":    shareURL,
		"share_token":  photo.ShareToken,
		"share_expire": photo.ShareExpiry,
	}})
}

// ListPhotoComments 获取照片评论
func (ac *AdvancedController) ListPhotoComments(c *gin.Context) {
	photoUUID := c.Param("id")

	var comments []models.FamilyAlbumComment
	ac.DB.Where("photo_uuid = ?", photoUUID).Order("created_at DESC").Find(&comments)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": comments})
}

// AddPhotoComment 添加评论
func (ac *AdvancedController) AddPhotoComment(c *gin.Context) {
	userID := getUserIDFromContext(c)
	photoUUID := c.Param("id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容不能为空"})
		return
	}

	// 获取用户名
	var userName string
	ac.DB.Table("sys_users").Where("id = ?", userID).Select("nickname").Scan(&userName)

	comment := models.FamilyAlbumComment{
		PhotoUUID: photoUUID,
		UserID:    userID,
		UserName:  userName,
		Content:   req.Content,
		ParentID:  req.ParentID,
	}

	if err := ac.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "评论失败"})
		return
	}

	// 更新照片评论数
	ac.DB.Model(&models.FamilyAlbum{}).Where("uuid = ?", photoUUID).Update("comment_count", gorm.Expr("comment_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": comment})
}

// LikePhoto 点赞照片
func (ac *AdvancedController) LikePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	photoUUID := c.Param("id")

	var existing models.FamilyAlbumLike
	if err := ac.DB.Where("photo_uuid = ? AND user_id = ?", photoUUID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已点赞"})
		return
	}

	like := models.FamilyAlbumLike{PhotoUUID: photoUUID, UserID: userID}
	if err := ac.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "点赞失败"})
		return
	}

	ac.DB.Model(&models.FamilyAlbum{}).Where("uuid = ?", photoUUID).Update("like_count", gorm.Expr("like_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UnlikePhoto 取消点赞
func (ac *AdvancedController) UnlikePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	photoUUID := c.Param("id")

	result := ac.DB.Where("photo_uuid = ? AND user_id = ?", photoUUID, userID).Delete(&models.FamilyAlbumLike{})
	if result.RowsAffected > 0 {
		ac.DB.Model(&models.FamilyAlbum{}).Where("uuid = ?", photoUUID).Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ==================== 寻回网络 API ====================

// ListPetFinderReports 获取寻宠报告列表
func (ac *AdvancedController) ListPetFinderReports(c *gin.Context) {
	userID := getUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	reportType := c.Query("report_type")
	status := c.DefaultQuery("status", "active")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := ac.DB.Model(&models.PetFinderReport{}).Where("user_id = ?", userID)
	if reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var reports []models.PetFinderReport
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success", "data": gin.H{
			"list": reports,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// CreatePetFinderReport 创建寻宠报告
func (ac *AdvancedController) CreatePetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		PetID       uint     `json:"pet_id" binding:"required"`
		ReportType  string   `json:"report_type" binding:"required"`
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description"`
		LastSeenAt  string   `json:"last_seen_at"`
		LastSeenLat float64  `json:"last_seen_lat"`
		LastSeenLng float64  `json:"last_seen_lng"`
		LastSeenAddr string  `json:"last_seen_addr"`
		ContactName string   `json:"contact_name"`
		ContactPhone string  `json:"contact_phone"`
		Reward      float64  `json:"reward"`
		RewardMemo  string   `json:"reward_memo"`
		Photos      []string `json:"photos"`
		AlertRadius float64  `json:"alert_radius"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var lastSeenAt *time.Time
	if req.LastSeenAt != "" {
		t := parseTime(req.LastSeenAt)
		lastSeenAt = &t
	}

	report := models.PetFinderReport{
		UserID:        userID,
		PetID:         req.PetID,
		ReportType:    req.ReportType,
		Title:         req.Title,
		Description:   req.Description,
		LastSeenAt:    lastSeenAt,
		LastSeenLat:   req.LastSeenLat,
		LastSeenLng:   req.LastSeenLng,
		LastSeenAddr:  req.LastSeenAddr,
		ContactName:   req.ContactName,
		ContactPhone:  req.ContactPhone,
		Reward:        req.Reward,
		RewardMemo:    req.RewardMemo,
		Status:        "active",
		AlertRadius:   req.AlertRadius,
	}
	if report.AlertRadius == 0 {
		report.AlertRadius = 5.0
	}

	// 转换 Photos slice to StringArray
	if len(req.Photos) > 0 {
		report.Photos = models.StringArray(req.Photos)
	}

	if err := ac.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// GetPetFinderReport 获取报告详情
func (ac *AdvancedController) GetPetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var report models.PetFinderReport
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}

	// 增加浏览次数
	ac.DB.Model(&report).Update("view_count", gorm.Expr("view_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// UpdatePetFinderReport 更新报告
func (ac *AdvancedController) UpdatePetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var report models.PetFinderReport
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	delete(req, "id")
	delete(req, "user_id")
	delete(req, "created_at")

	if err := ac.DB.Model(&report).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ac.DB.First(&report, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// ClosePetFinderReport 关闭报告
func (ac *AdvancedController) ClosePetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	now := time.Now()
	result := ac.DB.Model(&models.PetFinderReport{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"status": "closed", "resolved_at": now})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// AddSighting 添加目击报告
func (ac *AdvancedController) AddSighting(c *gin.Context) {
	reporterID := getUserIDFromContext(c)
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的报告ID"})
		return
	}

	var req struct {
		SightedAt   string  `json:"sighted_at"`
		SightedLat  float64 `json:"sighted_lat" binding:"required"`
		SightedLng  float64 `json:"sighted_lng" binding:"required"`
		SightedAddr string  `json:"sighted_addr"`
		Description string  `json:"description"`
		PetStatus   string  `json:"pet_status"`
		PhotoURL    string  `json:"photo_url"`
		ContactName string  `json:"contact_name"`
		ContactPhone string `json:"contact_phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var sightedAt *time.Time
	if req.SightedAt != "" {
		t := parseTime(req.SightedAt)
		sightedAt = &t
	}

	sighting := models.PetFinderSighting{
		ReportID:     uint(reportID),
		ReporterID:   reporterID,
		SightedAt:   sightedAt,
		SightedLat:  req.SightedLat,
		SightedLng:  req.SightedLng,
		SightedAddr: req.SightedAddr,
		Description: req.Description,
		PetStatus:   req.PetStatus,
		PhotoURL:    req.PhotoURL,
		ContactName: req.ContactName,
		ContactPhone: req.ContactPhone,
	}

	if err := ac.DB.Create(&sighting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "添加目击失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": sighting})
}

// ListSightings 获取目击报告列表
func (ac *AdvancedController) ListSightings(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的报告ID"})
		return
	}

	var sightings []models.PetFinderSighting
	ac.DB.Where("report_id = ?", reportID).Order("created_at DESC").Find(&sightings)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": sightings})
}

// GetNearbyPets 获取附近宠物
func (ac *AdvancedController) GetNearbyPets(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "10") // km

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的纬度"})
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的经度"})
		return
	}
	radius, _ := strconv.ParseFloat(radiusStr, 64)

	// 使用 Haversine 公式计算距离附近的寻宠报告
	// MySQL/PostgreSQL 支持：3959 * acos(cos(radians(?)) * cos(radians(last_seen_lat)) * cos(radians(last_seen_lng) - radians(?)) + sin(radians(?)) * sin(radians(last_seen_lat))))
	// 这里简化为矩形范围查询，生产环境应使用 PostGIS 或 MySQL GIS
	latDelta := radius / 111.0   // 约111km/度
	lngDelta := radius / (111.0 * 0.75)

	var reports []models.PetFinderReport
	ac.DB.Where("status = 'active' AND report_type IN ('lost','found')").
		Where("last_seen_lat BETWEEN ? AND ? AND last_seen_lng BETWEEN ? AND ?",
			lat-latDelta, lat+latDelta, lng-lngDelta, lng+lngDelta).
		Order("created_at DESC").
		Limit(50).
		Find(&reports)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"reports": reports,
		"center": gin.H{
			"lat":    lat,
			"lng":    lng,
			"radius": radius,
		},
	}})
}

// ==================== 疫苗接种 API ====================

// ListVaccinations 获取疫苗接种列表
func (ac *AdvancedController) ListVaccinations(c *gin.Context) {
	userID := getUserIDFromContext(c)
	petIDStr := c.Query("pet_id")

	query := ac.DB.Model(&models.PetVaccination{}).Where("user_id = ?", userID)
	if petIDStr != "" {
		query = query.Where("pet_id = ?", petIDStr)
	}

	var records []models.PetVaccination
	query.Order("inoculation_date DESC").Find(&records)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": records})
}

// CreateVaccination 创建疫苗记录
func (ac *AdvancedController) CreateVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.PetVaccination
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	req.UserID = userID
	if err := ac.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// GetVaccination 获取单条疫苗记录
func (ac *AdvancedController) GetVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var record models.PetVaccination
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateVaccination 更新疫苗记录
func (ac *AdvancedController) UpdateVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var record models.PetVaccination
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	delete(req, "id")
	delete(req, "user_id")
	delete(req, "created_at")

	if err := ac.DB.Model(&record).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ac.DB.First(&record, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeleteVaccination 删除疫苗记录
func (ac *AdvancedController) DeleteVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	result := ac.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PetVaccination{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ==================== 饮食记录 API ====================

// ListDietRecords 获取饮食记录列表
func (ac *AdvancedController) ListDietRecords(c *gin.Context) {
	userID := getUserIDFromContext(c)
	petIDStr := c.Query("pet_id")
	dateStr := c.Query("date")

	query := ac.DB.Model(&models.PetDietRecord{}).Where("user_id = ?", userID)
	if petIDStr != "" {
		query = query.Where("pet_id = ?", petIDStr)
	}
	if dateStr != "" {
		t := parseTime(dateStr)
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 0, 1)
		query = query.Where("eat_time >= ? AND eat_time < ?", start, end)
	}

	var records []models.PetDietRecord
	query.Order("eat_time DESC").Find(&records)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": records})
}

// CreateDietRecord 创建饮食记录
func (ac *AdvancedController) CreateDietRecord(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.PetDietRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	req.UserID = userID
	if err := ac.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// UpdateDietRecord 更新饮食记录
func (ac *AdvancedController) UpdateDietRecord(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var record models.PetDietRecord
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	delete(req, "id")
	delete(req, "user_id")
	delete(req, "created_at")

	if err := ac.DB.Model(&record).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ac.DB.First(&record, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeleteDietRecord 删除饮食记录
func (ac *AdvancedController) DeleteDietRecord(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	result := ac.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PetDietRecord{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetDietSummary 获取饮食汇总
func (ac *AdvancedController) GetDietSummary(c *gin.Context) {
	userID := getUserIDFromContext(c)
	petIDStr := c.Query("pet_id")

	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}

	petID, _ := strconv.ParseUint(petIDStr, 10, 64)

	// 本日、本周、本月统计
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	type Summary struct {
		TotalMeals  int     `json:"total_meals"`
		TotalCalories int   `json:"total_calories"`
		AvgAppetite float64 `json:"avg_appetite"`
	}

	todaySummary := Summary{}
	weekSummary := Summary{}
	monthSummary := Summary{}

	ac.DB.Model(&models.PetDietRecord{}).
		Select("COUNT(*) as total_meals, COALESCE(SUM(calories),0) as total_calories, COALESCE(AVG(appetite_score),0) as avg_appetite").
		Where("user_id = ? AND pet_id = ? AND eat_time >= ?", userID, petID, todayStart).
		Scan(&todaySummary)

	ac.DB.Model(&models.PetDietRecord{}).
		Select("COUNT(*) as total_meals, COALESCE(SUM(calories),0) as total_calories, COALESCE(AVG(appetite_score),0) as avg_appetite").
		Where("user_id = ? AND pet_id = ? AND eat_time >= ?", userID, petID, weekStart).
		Scan(&weekSummary)

	ac.DB.Model(&models.PetDietRecord{}).
		Select("COUNT(*) as total_meals, COALESCE(SUM(calories),0) as total_calories, COALESCE(AVG(appetite_score),0) as avg_appetite").
		Where("user_id = ? AND pet_id = ? AND eat_time >= ?", userID, petID, monthStart).
		Scan(&monthSummary)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"today": todaySummary,
		"week":  weekSummary,
		"month": monthSummary,
	}})
}

// ==================== 工具函数 ====================

func parseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			t, err = time.Parse("2006-01-02", s)
			if err != nil {
				return time.Now()
			}
		}
	}
	return t
}
