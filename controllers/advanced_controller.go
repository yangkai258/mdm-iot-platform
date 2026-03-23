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

// AdvancedController Advanced features controller (child mode / elderly care / family album / pet finder)
type AdvancedController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterAdvancedRoutes registers advanced feature routes under /api/v1/advanced
func (ctrl *AdvancedController) RegisterAdvancedRoutes(api *gin.RouterGroup) {
	ac := &AdvancedController{DB: ctrl.DB, Redis: ctrl.Redis}

	// Child mode
	childMode := api.Group("/advanced/child-mode")
	childMode.GET("/config", ac.GetChildModeConfig)
	childMode.PUT("/config", ac.UpdateChildModeConfig)
	childMode.GET("/stats", ac.GetChildModeStats)

	// Elderly care
	elderly := api.Group("/advanced/elderly")
	elderly.GET("/config", ac.GetElderlyConfig)
	elderly.PUT("/config", ac.UpdateElderlyConfig)
	elderly.GET("/health", ac.GetElderlyHealth)

	// Family album
	album := api.Group("/advanced/album")
	album.GET("/photos", ac.ListPhotos)
	album.POST("/photos", ac.UploadPhoto)
	album.DELETE("/photos/:id", ac.DeletePhoto)
	album.GET("/share/:id", ac.SharePhoto)
	album.GET("/photos/:id/comments", ac.ListPhotoComments)
	album.POST("/photos/:id/comments", ac.AddPhotoComment)
	album.POST("/photos/:id/like", ac.LikePhoto)
	album.DELETE("/photos/:id/like", ac.UnlikePhoto)

	// Pet finder
	petFinder := api.Group("/advanced/pet-finder")
	petFinder.GET("/reports", ac.ListPetFinderReports)
	petFinder.POST("/reports", ac.CreatePetFinderReport)
	petFinder.GET("/reports/:id", ac.GetPetFinderReport)
	petFinder.PUT("/reports/:id", ac.UpdatePetFinderReport)
	petFinder.DELETE("/reports/:id", ac.ClosePetFinderReport)
	petFinder.POST("/reports/:id/sighting", ac.AddSighting)
	petFinder.GET("/reports/:id/sightings", ac.ListSightings)
	petFinder.GET("/nearby", ac.GetNearbyPets)

	// Vaccination
	vaccination := api.Group("/advanced/vaccination")
	vaccination.GET("/records", ac.ListVaccinations)
	vaccination.POST("/records", ac.CreateVaccination)
	vaccination.PUT("/records/:id", ac.UpdateVaccination)
	vaccination.DELETE("/records/:id", ac.DeleteVaccination)
	vaccination.GET("/records/:id", ac.GetVaccination)

	// Diet
	diet := api.Group("/advanced/diet")
	diet.GET("/records", ac.ListDietRecords)
	diet.POST("/records", ac.CreateDietRecord)
	diet.PUT("/records/:id", ac.UpdateDietRecord)
	diet.DELETE("/records/:id", ac.DeleteDietRecord)
	diet.GET("/summary", ac.GetDietSummary)
}

// ==================== Child Mode APIs ====================

// GetChildModeConfig returns the child mode configuration for the user
func (ac *AdvancedController) GetChildModeConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ChildModeConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err == gorm.ErrRecordNotFound {
		config = models.ChildModeConfig{
			UserID:              userID,
			DeviceID:            deviceID,
			ContentFilterLevel: "moderate",
			SessionDuration:     30,
			BreakDuration:       10,
			AllowedStartTime:    "08:00",
			AllowedEndTime:      "20:00",
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// UpdateChildModeConfig updates the child mode configuration
func (ac *AdvancedController) UpdateChildModeConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.ChildModeConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	var config models.ChildModeConfig
	query := ac.DB.Where("user_id = ?", userID)
	if req.DeviceID != "" {
		query = query.Where("device_id = ?", req.DeviceID)
	}

	if err := query.First(&config).Error; err == gorm.ErrRecordNotFound {
		req.UserID = userID
		if err := ac.DB.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create config"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query config"})
		return
	}

	req.ID = config.ID
	req.UserID = userID
	if err := ac.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// GetChildModeStats returns usage statistics for child mode
func (ac *AdvancedController) GetChildModeStats(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ChildModeConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Config not found"})
		return
	}

	dailyStats := make([]map[string]interface{}, 7)
	weekStart := time.Now().AddDate(0, 0, -6)
	for i := 0; i < 7; i++ {
		day := weekStart.AddDate(0, 0, i)
		dailyStats[i] = map[string]interface{}{
			"date":         day.Format("2006-01-02"),
			"used_minutes": 0,
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"today_used_minutes": config.TodayUsedMinutes,
		"week_used_minutes":  config.WeekUsedMinutes,
		"total_used_minutes": config.TotalUsedMinutes,
		"total_sessions":     config.TotalSessions,
		"last_session_at":    config.LastSessionAt,
		"daily_stats":        dailyStats,
	}})
}

// ==================== Elderly Care APIs ====================

// GetElderlyConfig returns the elderly care configuration
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
			UserID:               userID,
			DeviceID:             deviceID,
			HealthMonitorEnabled: true,
			HeartRateAlertHigh:   100,
			HeartRateAlertLow:    50,
			ActivityGoal:         6000,
			SleepMonitoring:      true,
			MedicationReminders:  true,
			CompanionEnabled:     true,
			InteractionFrequency: "normal",
			VoiceCallEnabled:     true,
			FallDetectionEnabled: true,
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// UpdateElderlyConfig updates the elderly care configuration
func (ac *AdvancedController) UpdateElderlyConfig(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.ElderlyCareConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create config"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query config"})
		return
	}

	req.ID = config.ID
	req.UserID = userID
	if err := ac.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update config"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// GetElderlyHealth returns elderly health data
func (ac *AdvancedController) GetElderlyHealth(c *gin.Context) {
	userID := getUserIDFromContext(c)
	deviceID := c.Query("device_id")

	var config models.ElderlyCareConfig
	query := ac.DB.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Config not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"heart_rate":          72,
		"heart_rate_high":     config.HeartRateAlertHigh,
		"heart_rate_low":      config.HeartRateAlertLow,
		"activity_steps":      4500,
		"activity_goal":       config.ActivityGoal,
		"sleep_hours":         7.5,
		"last_interaction":     config.LastInteractionAt,
		"total_interactions":   config.TotalInteractions,
		"fall_detected":       false,
		"medication_taken":    true,
		"updated_at":          time.Now(),
	}})
}

// ==================== Family Album APIs ====================

// ListPhotos returns the list of family photos
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

	query := ac.DB.Model(&models.FamilyPhoto{}).Where("user_id = ?", userID)
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

	var photos []models.FamilyPhoto
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&photos)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success", "data": gin.H{
			"list": photos,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// UploadPhoto uploads a new family photo
func (ac *AdvancedController) UploadPhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		DeviceID     string `json:"device_id"`
		PetID        uint   `json:"pet_id"`
		Title        string `json:"title" binding:"required"`
		Description  string `json:"description"`
		PhotoURL     string `json:"photo_url" binding:"required"`
		ThumbnailURL string `json:"thumbnail_url"`
		FileSize     int64  `json:"file_size"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		MimeType     string `json:"mime_type"`
		Category     string `json:"category"`
		Tags         string `json:"tags"`
		TakenLat     float64 `json:"taken_lat"`
		TakenLng     float64 `json:"taken_lng"`
		TakenAddr    string  `json:"taken_addr"`
		TakenAt      string  `json:"taken_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	var takenAt *time.Time
	if req.TakenAt != "" {
		t := parseTimeAdvanced(req.TakenAt)
		takenAt = &t
	}

	photo := models.FamilyPhoto{
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": photo})
}

// DeletePhoto deletes a family photo
func (ac *AdvancedController) DeletePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	result := ac.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.FamilyPhoto{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Deleted"})
}

// SharePhoto generates a share link for a photo
func (ac *AdvancedController) SharePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	var photo models.FamilyPhoto
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Photo not found"})
		return
	}

	if photo.ShareToken == "" {
		photo.ShareToken = fmt.Sprintf("share_%d_%d", photo.ID, time.Now().UnixNano())
	}
	photo.IsShared = true
	expiry := time.Now().Add(7 * 24 * time.Hour)
	photo.ShareExpiry = &expiry

	if err := ac.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Share failed"})
		return
	}

	shareURL := fmt.Sprintf("/api/v1/advanced/album/shared/%s", photo.ShareToken)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"share_url":    shareURL,
		"share_token":  photo.ShareToken,
		"share_expire": photo.ShareExpiry,
	}})
}

// ListPhotoComments returns comments for a photo
func (ac *AdvancedController) ListPhotoComments(c *gin.Context) {
	photoUUID := c.Param("id")

	var comments []models.FamilyAlbumComment
	ac.DB.Where("photo_uuid = ?", photoUUID).Order("created_at DESC").Find(&comments)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": comments})
}

// AddPhotoComment adds a comment to a photo
func (ac *AdvancedController) AddPhotoComment(c *gin.Context) {
	userID := getUserIDFromContext(c)
	photoUUID := c.Param("id")

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID uint   `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Content required"})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Comment failed"})
		return
	}

	ac.DB.Model(&models.FamilyPhoto{}).Where("uuid = ?", photoUUID).Update("comment_count", gorm.Expr("comment_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": comment})
}

// LikePhoto likes a photo
func (ac *AdvancedController) LikePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	photoUUID := c.Param("id")

	var existing models.FamilyAlbumLike
	if err := ac.DB.Where("photo_uuid = ? AND user_id = ?", photoUUID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Already liked"})
		return
	}

	like := models.FamilyAlbumLike{PhotoUUID: photoUUID, UserID: userID}
	if err := ac.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Like failed"})
		return
	}

	ac.DB.Model(&models.FamilyPhoto{}).Where("uuid = ?", photoUUID).Update("like_count", gorm.Expr("like_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UnlikePhoto unlikes a photo
func (ac *AdvancedController) UnlikePhoto(c *gin.Context) {
	userID := getUserIDFromContext(c)
	photoUUID := c.Param("id")

	result := ac.DB.Where("photo_uuid = ? AND user_id = ?", photoUUID, userID).Delete(&models.FamilyAlbumLike{})
	if result.RowsAffected > 0 {
		ac.DB.Model(&models.FamilyPhoto{}).Where("uuid = ?", photoUUID).Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ==================== Pet Finder APIs ====================

// ListPetFinderReports returns the list of pet finder reports
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
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// CreatePetFinderReport creates a new pet finder report
func (ac *AdvancedController) CreatePetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		PetID        uint     `json:"pet_id" binding:"required"`
		ReportType   string   `json:"report_type" binding:"required"`
		Title        string   `json:"title" binding:"required"`
		Description  string   `json:"description"`
		LastSeenAt   string   `json:"last_seen_at"`
		LastSeenLat  float64  `json:"last_seen_lat"`
		LastSeenLng  float64  `json:"last_seen_lng"`
		LastSeenAddr string   `json:"last_seen_addr"`
		ContactName  string   `json:"contact_name"`
		ContactPhone string   `json:"contact_phone"`
		Reward       float64  `json:"reward"`
		RewardMemo   string   `json:"reward_memo"`
		Photos       []string `json:"photos"`
		AlertRadius  float64  `json:"alert_radius"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	var lastSeenAt *time.Time
	if req.LastSeenAt != "" {
		t := parseTimeAdvanced(req.LastSeenAt)
		lastSeenAt = &t
	}

	report := models.PetFinderReport{
		UserID:       userID,
		PetID:        req.PetID,
		ReportType:   req.ReportType,
		Title:        req.Title,
		Description:  req.Description,
		LastSeenAt:   lastSeenAt,
		LastSeenLat:  req.LastSeenLat,
		LastSeenLng:  req.LastSeenLng,
		LastSeenAddr: req.LastSeenAddr,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Reward:       req.Reward,
		RewardMemo:   req.RewardMemo,
		Status:       "active",
		AlertRadius:  req.AlertRadius,
	}
	if report.AlertRadius == 0 {
		report.AlertRadius = 5.0
	}
	if len(req.Photos) > 0 {
		report.Photos = models.StringArray(req.Photos)
	}

	if err := ac.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// GetPetFinderReport returns a single pet finder report
func (ac *AdvancedController) GetPetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	var report models.PetFinderReport
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Report not found"})
		return
	}

	ac.DB.Model(&report).Update("view_count", gorm.Expr("view_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// UpdatePetFinderReport updates a pet finder report
func (ac *AdvancedController) UpdatePetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	var report models.PetFinderReport
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Report not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	delete(req, "id")
	delete(req, "user_id")
	delete(req, "created_at")

	if err := ac.DB.Model(&report).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Update failed"})
		return
	}

	ac.DB.First(&report, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// ClosePetFinderReport closes a pet finder report
func (ac *AdvancedController) ClosePetFinderReport(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	now := time.Now()
	result := ac.DB.Model(&models.PetFinderReport{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"status": "closed", "resolved_at": now})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// AddSighting adds a sighting to a report
func (ac *AdvancedController) AddSighting(c *gin.Context) {
	reporterID := getUserIDFromContext(c)
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid report ID"})
		return
	}

	var req struct {
		SightedAt    string  `json:"sighted_at"`
		SightedLat   float64 `json:"sighted_lat" binding:"required"`
		SightedLng   float64 `json:"sighted_lng" binding:"required"`
		SightedAddr  string  `json:"sighted_addr"`
		Description  string  `json:"description"`
		PetStatus    string  `json:"pet_status"`
		PhotoURL     string  `json:"photo_url"`
		ContactName  string  `json:"contact_name"`
		ContactPhone string  `json:"contact_phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	var sightedAt *time.Time
	if req.SightedAt != "" {
		t := parseTimeAdvanced(req.SightedAt)
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to add sighting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": sighting})
}

// ListSightings returns sightings for a report
func (ac *AdvancedController) ListSightings(c *gin.Context) {
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid report ID"})
		return
	}

	var sightings []models.PetFinderSighting
	ac.DB.Where("report_id = ?", reportID).Order("created_at DESC").Find(&sightings)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": sightings})
}

// GetNearbyPets returns nearby lost/found pets
func (ac *AdvancedController) GetNearbyPets(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "10")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid latitude"})
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid longitude"})
		return
	}
	radius, _ := strconv.ParseFloat(radiusStr, 64)

	latDelta := radius / 111.0
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

// ==================== Vaccination APIs ====================

// ListVaccinations returns pet vaccination records
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

// CreateVaccination creates a vaccination record
func (ac *AdvancedController) CreateVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.PetVaccination
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	req.UserID = userID
	if err := ac.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Create failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// GetVaccination returns a single vaccination record
func (ac *AdvancedController) GetVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	var record models.PetVaccination
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateVaccination updates a vaccination record
func (ac *AdvancedController) UpdateVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	var record models.PetVaccination
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Record not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	delete(req, "id")
	delete(req, "user_id")
	delete(req, "created_at")

	if err := ac.DB.Model(&record).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Update failed"})
		return
	}

	ac.DB.First(&record, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeleteVaccination deletes a vaccination record
func (ac *AdvancedController) DeleteVaccination(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	result := ac.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PetVaccination{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Deleted"})
}

// ==================== Diet Record APIs ====================

// ListDietRecords returns pet diet records
func (ac *AdvancedController) ListDietRecords(c *gin.Context) {
	userID := getUserIDFromContext(c)
	petIDStr := c.Query("pet_id")
	dateStr := c.Query("date")

	query := ac.DB.Model(&models.PetDietRecord{}).Where("user_id = ?", userID)
	if petIDStr != "" {
		query = query.Where("pet_id = ?", petIDStr)
	}
	if dateStr != "" {
		t := parseTimeAdvanced(dateStr)
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 0, 1)
		query = query.Where("eat_time >= ? AND eat_time < ?", start, end)
	}

	var records []models.PetDietRecord
	query.Order("eat_time DESC").Find(&records)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": records})
}

// CreateDietRecord creates a diet record
func (ac *AdvancedController) CreateDietRecord(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req models.PetDietRecord
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	req.UserID = userID
	if err := ac.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Create failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// UpdateDietRecord updates a diet record
func (ac *AdvancedController) UpdateDietRecord(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	var record models.PetDietRecord
	if err := ac.DB.Where("id = ? AND user_id = ?", id, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Record not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request"})
		return
	}

	delete(req, "id")
	delete(req, "user_id")
	delete(req, "created_at")

	if err := ac.DB.Model(&record).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Update failed"})
		return
	}

	ac.DB.First(&record, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeleteDietRecord deletes a diet record
func (ac *AdvancedController) DeleteDietRecord(c *gin.Context) {
	userID := getUserIDFromContext(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid ID"})
		return
	}

	result := ac.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PetDietRecord{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Deleted"})
}

// GetDietSummary returns diet summary statistics
func (ac *AdvancedController) GetDietSummary(c *gin.Context) {
	userID := getUserIDFromContext(c)
	petIDStr := c.Query("pet_id")

	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id required"})
		return
	}

	petID, _ := strconv.ParseUint(petIDStr, 10, 64)

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	type Summary struct {
		TotalMeals    int     `json:"total_meals"`
		TotalCalories int     `json:"total_calories"`
		AvgAppetite   float64 `json:"avg_appetite"`
	}

	var todaySummary, weekSummary, monthSummary Summary

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

// ==================== Utility ====================

func parseTimeAdvanced(s string) time.Time {
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
