package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MarketController content marketplace controller
type MarketController struct {
	DB *gorm.DB
}

// RegisterMarketRoutes registers content market routes
func (ctrl *MarketController) RegisterMarketRoutes(api *gin.RouterGroup) {
	// Emoticon market
	emoticons := api.Group("/market/emoticons")
	{
		emoticons.GET("", ctrl.ListEmoticonPacks)
		emoticons.POST("", ctrl.CreateEmoticonPack)
		emoticons.GET("/:id", ctrl.GetEmoticonPack)
		emoticons.DELETE("/:id", ctrl.DeleteEmoticonPack)
		emoticons.POST("/:id/publish", ctrl.PublishEmoticonPack)
	}

	// Action resource library
	actions := api.Group("/market/actions")
	{
		actions.GET("", ctrl.ListActionResources)
		actions.POST("", ctrl.CreateActionResource)
		actions.GET("/:id", ctrl.GetActionResource)
		actions.DELETE("/:id", ctrl.DeleteActionResource)
		actions.POST("/:id/publish", ctrl.PublishActionResource)
	}

	// Voice customization
	voices := api.Group("/market/voices")
	{
		voices.GET("", ctrl.ListVoiceConfigs)
		voices.POST("", ctrl.CreateVoiceConfig)
		voices.GET("/:id", ctrl.GetVoiceConfig)
		voices.DELETE("/:id", ctrl.DeleteVoiceConfig)
		voices.POST("/:id/apply", ctrl.ApplyVoiceConfig)
	}
}

// ============ Emoticon Pack APIs ============

// ListEmoticonPacks GET /api/v1/market/emoticons
func (ctrl *MarketController) ListEmoticonPacks(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	packType := c.Query("pack_type")
	keyword := c.Query("keyword")
	status := c.DefaultQuery("status", "published")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.EmoticonPack{}).Where("status = ?", status)

	if packType != "" {
		query = query.Where("pack_type = ?", packType)
	}
	if keyword != "" {
		query = query.Where("pack_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	validSortFields := map[string]bool{"created_at": true, "downloads": true, "rating_avg": true, "price": true}
	if validSortFields[sortBy] {
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Count(&total)

	var packs []models.EmoticonPack
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&packs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emoticon packs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       packs,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateEmoticonPack POST /api/v1/market/emoticons
func (ctrl *MarketController) CreateEmoticonPack(c *gin.Context) {
	var req struct {
		PackName     string         `json:"pack_name" binding:"required"`
		PackType     string         `json:"pack_type" binding:"required"`
		Description  string         `json:"description"`
		ThumbnailURL string         `json:"thumbnail_url"`
		PreviewURL   string         `json:"preview_url"`
		Emoticons    models.JSON    `json:"emoticons"`
		Price        float64        `json:"price"`
		IsFree       bool           `json:"is_free"`
		Tags         []string       `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	isFree := req.Price == 0

	pack := models.EmoticonPack{
		UserID:       userID,
		PackName:     req.PackName,
		PackType:     req.PackType,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		PreviewURL:   req.PreviewURL,
		Emoticons:    req.Emoticons,
		Price:        req.Price,
		IsFree:       isFree,
		Status:       "draft",
		Tags:         req.Tags,
	}

	if err := ctrl.DB.Create(&pack).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create emoticon pack"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": pack})
}

// GetEmoticonPack GET /api/v1/market/emoticons/:id
func (ctrl *MarketController) GetEmoticonPack(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var pack models.EmoticonPack
	if err := ctrl.DB.First(&pack, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Emoticon pack not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emoticon pack"})
		return
	}

	ctrl.DB.Model(&pack).UpdateColumn("downloads", gorm.Expr("downloads + ?", 1))

	c.JSON(http.StatusOK, gin.H{"item": pack})
}

// DeleteEmoticonPack DELETE /api/v1/market/emoticons/:id
func (ctrl *MarketController) DeleteEmoticonPack(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := getUserIDFromContext(c)

	var pack models.EmoticonPack
	if err := ctrl.DB.First(&pack, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emoticon pack not found"})
		return
	}

	if pack.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this pack"})
		return
	}

	if err := ctrl.DB.Delete(&pack).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete emoticon pack"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emoticon pack deleted"})
}

// PublishEmoticonPack POST /api/v1/market/emoticons/:id/publish
func (ctrl *MarketController) PublishEmoticonPack(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := getUserIDFromContext(c)

	var pack models.EmoticonPack
	if err := ctrl.DB.First(&pack, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emoticon pack not found"})
		return
	}

	if pack.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if pack.Status != "draft" && pack.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft or rejected packs can be published"})
		return
	}

	pack.Status = "pending"
	now := time.Now()
	pack.ReviewedAt = nil

	if err := ctrl.DB.Save(&pack).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish emoticon pack"})
		return
	}

	review := models.ContentReview{
		ItemType:    "emoticon",
		ItemID:      uint(id),
		UserID:      userID,
		Status:      "pending",
		SubmittedAt: now,
	}
	ctrl.DB.Create(&review)

	c.JSON(http.StatusOK, gin.H{"item": pack, "message": "Emoticon pack submitted for review"})
}

// ============ Action Resource APIs ============

// ListActionResources GET /api/v1/market/actions
func (ctrl *MarketController) ListActionResources(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	actionType := c.Query("action_type")
	difficulty := c.Query("difficulty")
	keyword := c.Query("keyword")
	status := c.DefaultQuery("status", "published")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.ActionResource{}).Where("status = ?", status)

	if actionType != "" {
		query = query.Where("action_type = ?", actionType)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if keyword != "" {
		query = query.Where("action_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	validSortFields := map[string]bool{"created_at": true, "downloads": true, "rating_avg": true, "price": true, "duration_sec": true}
	if validSortFields[sortBy] {
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Count(&total)

	var actions []models.ActionResource
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch action resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       actions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateActionResource POST /api/v1/market/actions
func (ctrl *MarketController) CreateActionResource(c *gin.Context) {
	var req struct {
		ActionName   string         `json:"action_name" binding:"required"`
		ActionType   string         `json:"action_type" binding:"required"`
		Description  string         `json:"description"`
		Difficulty   string         `json:"difficulty"`
		ThumbnailURL string         `json:"thumbnail_url"`
		VideoURL     string         `json:"video_url"`
		MotionData   models.JSON   `json:"motion_data"`
		Price        float64        `json:"price"`
		IsFree       bool           `json:"is_free"`
		Tags         []string       `json:"tags"`
		DurationSec  int            `json:"duration_sec"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	isFree := req.Price == 0

	action := models.ActionResource{
		UserID:       userID,
		ActionName:   req.ActionName,
		ActionType:   req.ActionType,
		Description:  req.Description,
		Difficulty:   req.Difficulty,
		ThumbnailURL: req.ThumbnailURL,
		VideoURL:     req.VideoURL,
		MotionData:   req.MotionData,
		Price:        req.Price,
		IsFree:       isFree,
		Status:       "draft",
		Tags:         req.Tags,
		DurationSec:  req.DurationSec,
	}

	if err := ctrl.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create action resource"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": action})
}

// GetActionResource GET /api/v1/market/actions/:id
func (ctrl *MarketController) GetActionResource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var action models.ActionResource
	if err := ctrl.DB.First(&action, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Action resource not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch action resource"})
		return
	}

	ctrl.DB.Model(&action).UpdateColumn("downloads", gorm.Expr("downloads + ?", 1))

	c.JSON(http.StatusOK, gin.H{"item": action})
}

// DeleteActionResource DELETE /api/v1/market/actions/:id
func (ctrl *MarketController) DeleteActionResource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := getUserIDFromContext(c)

	var action models.ActionResource
	if err := ctrl.DB.First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action resource not found"})
		return
	}

	if action.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this action"})
		return
	}

	if err := ctrl.DB.Delete(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete action resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action resource deleted"})
}

// PublishActionResource POST /api/v1/market/actions/:id/publish
func (ctrl *MarketController) PublishActionResource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := getUserIDFromContext(c)

	var action models.ActionResource
	if err := ctrl.DB.First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action resource not found"})
		return
	}

	if action.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if action.Status != "draft" && action.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft or rejected actions can be published"})
		return
	}

	action.Status = "pending"
	now := time.Now()
	action.ReviewedAt = nil

	if err := ctrl.DB.Save(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish action resource"})
		return
	}

	review := models.ContentReview{
		ItemType:    "action",
		ItemID:      uint(id),
		UserID:      userID,
		Status:      "pending",
		SubmittedAt: now,
	}
	ctrl.DB.Create(&review)

	c.JSON(http.StatusOK, gin.H{"item": action, "message": "Action resource submitted for review"})
}

// ============ Voice Config APIs ============

// ListVoiceConfigs GET /api/v1/market/voices
func (ctrl *MarketController) ListVoiceConfigs(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	voiceType := c.Query("voice_type")
	keyword := c.Query("keyword")
	status := c.DefaultQuery("status", "published")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.VoiceConfig{}).Where("status = ?", status)

	if voiceType != "" {
		query = query.Where("voice_type = ?", voiceType)
	}
	if keyword != "" {
		query = query.Where("voice_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	validSortFields := map[string]bool{"created_at": true, "downloads": true, "rating_avg": true, "price": true}
	if validSortFields[sortBy] {
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Count(&total)

	var voices []models.VoiceConfig
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&voices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch voice configs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       voices,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// CreateVoiceConfig POST /api/v1/market/voices
func (ctrl *MarketController) CreateVoiceConfig(c *gin.Context) {
	var req struct {
		VoiceName    string         `json:"voice_name" binding:"required"`
		VoiceType    string         `json:"voice_type" binding:"required"`
		Description  string         `json:"description"`
		PreviewURL   string         `json:"preview_url"`
		AudioSamples models.JSON    `json:"audio_samples"`
		VoiceParams  models.JSON    `json:"voice_params"`
		Price        float64        `json:"price"`
		IsFree       bool           `json:"is_free"`
		Tags         []string       `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	isFree := req.Price == 0

	voice := models.VoiceConfig{
		UserID:       userID,
		VoiceName:    req.VoiceName,
		VoiceType:    req.VoiceType,
		Description:  req.Description,
		PreviewURL:   req.PreviewURL,
		AudioSamples: req.AudioSamples,
		VoiceParams:  req.VoiceParams,
		Price:        req.Price,
		IsFree:       isFree,
		Status:       "draft",
		Tags:         req.Tags,
	}

	if err := ctrl.DB.Create(&voice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create voice config"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": voice})
}

// GetVoiceConfig GET /api/v1/market/voices/:id
func (ctrl *MarketController) GetVoiceConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var voice models.VoiceConfig
	if err := ctrl.DB.First(&voice, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Voice config not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch voice config"})
		return
	}

	ctrl.DB.Model(&voice).UpdateColumn("downloads", gorm.Expr("downloads + ?", 1))

	c.JSON(http.StatusOK, gin.H{"item": voice})
}

// DeleteVoiceConfig DELETE /api/v1/market/voices/:id
func (ctrl *MarketController) DeleteVoiceConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := getUserIDFromContext(c)

	var voice models.VoiceConfig
	if err := ctrl.DB.First(&voice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voice config not found"})
		return
	}

	if voice.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this voice"})
		return
	}

	if err := ctrl.DB.Delete(&voice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete voice config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voice config deleted"})
}

// ApplyVoiceConfig POST /api/v1/market/voices/:id/apply
func (ctrl *MarketController) ApplyVoiceConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)

	var voice models.VoiceConfig
	if err := ctrl.DB.First(&voice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voice config not found"})
		return
	}

	if voice.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only published voices can be applied"})
		return
	}

	// Verify device ownership
	var device models.Device
	if err := ctrl.DB.Where("id = ? AND user_id = ?", req.DeviceID, userID).First(&device).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Device not found or not authorized"})
		return
	}

	// Record purchase if paid
	if !voice.IsFree {
		purchase := models.UserPurchase{
			UserID:   userID,
			ItemType: "voice",
			ItemID:   uint(id),
			Price:    voice.Price,
			Status:   "completed",
		}
		ctrl.DB.Create(&purchase)
	}

	// TODO: publish voice config to device via MQTT
	c.JSON(http.StatusOK, gin.H{
		"message":   "Voice config applied successfully",
		"device_id": req.DeviceID,
		"voice_id":  id,
	})
}
