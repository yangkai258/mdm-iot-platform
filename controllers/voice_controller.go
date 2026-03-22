package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VoiceController 声音配置控制器
type VoiceController struct {
	DB *gorm.DB
}

// RegisterVoiceRoutes 注册声音配置路由
func (ctrl *VoiceController) RegisterVoiceRoutes(api *gin.RouterGroup) {
	api.GET("/voices", ctrl.List)
	api.POST("/voices", ctrl.Create)
	api.GET("/voices/:id", ctrl.Get)
	api.PUT("/voices/:id", ctrl.Update)
	api.DELETE("/voices/:id", ctrl.Delete)
	api.POST("/voices/preview", ctrl.Preview)
}

// List 获取声音配置列表
// GET /api/v1/voices
func (ctrl *VoiceController) List(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	userID := getUserIDFromContext(c)
	isPublic := c.Query("is_public")
	provider := c.Query("provider")
	language := c.Query("language")
	gender := c.Query("gender")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	// 如果没有指定用户且没有请求公开列表，默认返回自己的配置
	query := ctrl.DB.Model(&models.VoiceConfig{})

	if userID > 0 && isPublic != "true" {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	} else if isPublic == "true" {
		query = query.Where("is_public = ?", true)
	}

	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if language != "" {
		query = query.Where("language = ?", language)
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	// 排序
	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	var total int64
	query.Count(&total)

	var configs []models.VoiceConfig
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query voice configs: " + err.Error(),
		})
		return
	}

	responses := make([]*models.VoiceConfigResponse, len(configs))
	for i := range configs {
		responses[i] = configs[i].ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      responses,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Get 获取声音配置详情
// GET /api/v1/voices/:id
func (ctrl *VoiceController) Get(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var config models.VoiceConfig
	if err := ctrl.DB.Where("id = ? OR config_uuid = ?", id, id).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Voice config not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限（自己的或公开的）
	if config.UserID != userID && !config.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config.ToResponse(),
	})
}

// Create 创建声音配置
// POST /api/v1/voices
func (ctrl *VoiceController) Create(c *gin.Context) {
	var req struct {
		Name       string  `json:"name" binding:"required"`
		Provider   string  `json:"provider" binding:"required"`
		VoiceID    string  `json:"voice_id"`
		Language   string  `json:"language"`
		Gender     string  `json:"gender"`
		AgeGroup   string  `json:"age_group"`
		Pitch      float64 `json:"pitch"`
		Speed      float64 `json:"speed"`
		Volume     float64 `json:"volume"`
		Emotion    string  `json:"emotion"`
		Style      string  `json:"style"`
		IsDefault  bool    `json:"is_default"`
		IsPublic   bool    `json:"is_public"`
		PreviewURL string  `json:"preview_url"`
		Settings   string  `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(c)

	// 设置默认值
	if req.Language == "" {
		req.Language = "zh-CN"
	}
	if req.Pitch == 0 {
		req.Pitch = 1.0
	}
	if req.Speed == 0 {
		req.Speed = 1.0
	}
	if req.Volume == 0 {
		req.Volume = 1.0
	}
	if req.Settings == "" {
		req.Settings = "{}"
	}

	config := models.VoiceConfig{
		UserID:     userID,
		Name:       req.Name,
		Provider:   req.Provider,
		VoiceID:    req.VoiceID,
		Language:   req.Language,
		Gender:     req.Gender,
		AgeGroup:   req.AgeGroup,
		Pitch:      req.Pitch,
		Speed:      req.Speed,
		Volume:     req.Volume,
		Emotion:    req.Emotion,
		Style:      req.Style,
		IsDefault:  req.IsDefault,
		IsPublic:   req.IsPublic,
		PreviewURL: req.PreviewURL,
		Settings:   req.Settings,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := ctrl.DB.Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create voice config: " + err.Error()})
		return
	}

	// 如果设为默认，取消其他默认
	if req.IsDefault {
		ctrl.DB.Model(&models.VoiceConfig{}).
			Where("user_id = ? AND id != ?", userID, config.ID).
			Update("is_default", false)
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": config.ToResponse(),
	})
}

// Update 更新声音配置
// PUT /api/v1/voices/:id
func (ctrl *VoiceController) Update(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var config models.VoiceConfig
	if err := ctrl.DB.Where("id = ? OR config_uuid = ?", id, id).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Voice config not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限
	if config.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	var req struct {
		Name       string  `json:"name"`
		Provider   string  `json:"provider"`
		VoiceID    string  `json:"voice_id"`
		Language   string  `json:"language"`
		Gender     string  `json:"gender"`
		AgeGroup   string  `json:"age_group"`
		Pitch      float64 `json:"pitch"`
		Speed      float64 `json:"speed"`
		Volume     float64 `json:"volume"`
		Emotion    string  `json:"emotion"`
		Style      string  `json:"style"`
		IsDefault  *bool   `json:"is_default"`
		IsPublic   *bool   `json:"is_public"`
		PreviewURL string  `json:"preview_url"`
		Settings   string  `json:"settings"`
		Status     string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Provider != "" {
		updates["provider"] = req.Provider
	}
	if req.VoiceID != "" {
		updates["voice_id"] = req.VoiceID
	}
	if req.Language != "" {
		updates["language"] = req.Language
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.AgeGroup != "" {
		updates["age_group"] = req.AgeGroup
	}
	if req.Pitch > 0 {
		updates["pitch"] = req.Pitch
	}
	if req.Speed > 0 {
		updates["speed"] = req.Speed
	}
	if req.Volume > 0 {
		updates["volume"] = req.Volume
	}
	if req.Emotion != "" {
		updates["emotion"] = req.Emotion
	}
	if req.Style != "" {
		updates["style"] = req.Style
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.PreviewURL != "" {
		updates["preview_url"] = req.PreviewURL
	}
	if req.Settings != "" {
		updates["settings"] = req.Settings
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := ctrl.DB.Model(&config).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update voice config: " + err.Error()})
		return
	}

	// 如果设为默认，取消其他默认
	if req.IsDefault != nil && *req.IsDefault {
		ctrl.DB.Model(&models.VoiceConfig{}).
			Where("user_id = ? AND id != ?", userID, config.ID).
			Update("is_default", false)
	}

	ctrl.DB.First(&config, config.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config.ToResponse(),
	})
}

// Delete 删除声音配置
// DELETE /api/v1/voices/:id
func (ctrl *VoiceController) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := getUserIDFromContext(c)

	var config models.VoiceConfig
	if err := ctrl.DB.Where("id = ? OR config_uuid = ?", id, id).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Voice config not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 检查权限
	if config.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Permission denied"})
		return
	}

	if err := ctrl.DB.Delete(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete voice config: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Voice config deleted successfully",
	})
}

// Preview 预览声音
// POST /api/v1/voices/preview
func (ctrl *VoiceController) Preview(c *gin.Context) {
	var req models.VoicePreviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	// 设置默认值
	if req.Language == "" {
		req.Language = "zh-CN"
	}
	if req.Pitch == 0 {
		req.Pitch = 1.0
	}
	if req.Speed == 0 {
		req.Speed = 1.0
	}
	if req.Volume == 0 {
		req.Volume = 1.0
	}

	// 这里应该调用实际的声音合成服务
	// 简化实现：返回模拟预览URL
	previewURL := ""
	durationMs := len(req.Text) * 100 // 估算：每秒约10个字符

	switch req.Provider {
	case "elevenlabs":
		previewURL = "https://api.elevenlabs.io/preview/mock/" + req.VoiceID
	case "azure":
		previewURL = "https://azure.tts.preview/mock"
	case "gtts":
		previewURL = "https://translate.google.com/translate_tts?tl=" + req.Language
	default:
		previewURL = "https://api.example.com/voice/preview"
	}

	// 实际项目中这里会：
	// 1. 调用第三方TTS API生成音频
	// 2. 上传到OSS/S3等存储
	// 3. 返回可访问的URL
	// 4. 记录预览日志

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": models.VoicePreviewResponse{
			PreviewURL: previewURL,
			DurationMs: durationMs,
			Format:     "mp3",
		},
	})
}
