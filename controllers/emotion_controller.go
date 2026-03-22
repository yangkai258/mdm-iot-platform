package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmotionController 情绪控制器
type EmotionController struct {
	DB *gorm.DB
}

// RegisterEmotionRoutes 注册情绪相关路由
func (c *EmotionController) RegisterEmotionRoutes(r *gin.RouterGroup) {
	// 情绪识别
	r.POST("/emotion/recognize/voice", c.RecognizeVoice)
	r.POST("/emotion/recognize/text", c.RecognizeText)
	r.POST("/emotion/recognize/pet", c.RecognizePet)
	r.POST("/emotion/recognize/batch", c.RecognizeBatch)

	// 情绪响应
	r.GET("/emotion/:pet_id/response-strategies", c.GetResponseStrategies)
	r.POST("/emotion/:pet_id/response", c.TriggerEmotionResponse)
	r.PUT("/emotion/:pet_id/response-config", c.UpdateResponseConfig)

	// 情绪日志
	r.GET("/emotion/logs", c.ListEmotionLogs)
	r.GET("/emotion/logs/:id", c.GetEmotionLog)
	r.PUT("/emotion/logs/:id", c.UpdateEmotionLog)
	r.DELETE("/emotion/logs/:id", c.DeleteEmotionLog)

	// 情绪动作库 CRUD
	r.GET("/emotion/actions", c.ListEmotionActions)
	r.POST("/emotion/actions", c.CreateEmotionAction)
	r.GET("/emotion/actions/:id", c.GetEmotionAction)
	r.PUT("/emotion/actions/:id", c.UpdateEmotionAction)
	r.DELETE("/emotion/actions/:id", c.DeleteEmotionAction)

	// 情绪报告
	r.GET("/emotion/reports", c.ListEmotionReports)
	r.GET("/emotion/reports/:id", c.GetEmotionReport)
	r.POST("/emotion/reports/generate", c.GenerateEmotionReport)

	// 家庭情绪地图
	r.GET("/emotion/family-map", c.GetFamilyEmotionMap)
}

// ==================== 情绪识别 API ====================

// RecognizeVoiceRequest 语音情绪识别请求
type RecognizeVoiceRequest struct {
	SubjectType string `json:"subject_type" binding:"required,oneof=user pet"`
	SubjectID   uint   `json:"subject_id" binding:"required"`
	AudioURL    string `json:"audio_url"`
	AudioData   string `json:"audio_data"`
	Language    string `json:"language"`
}

// RecognizeTextRequest 文字情绪识别请求
type RecognizeTextRequest struct {
	SubjectType string `json:"subject_type" binding:"required,oneof=user pet"`
	SubjectID   uint   `json:"subject_id" binding:"required"`
	Text        string `json:"text" binding:"required"`
}

// RecognizePetRequest 宠物情绪识别请求
type RecognizePetRequest struct {
	PetID        uint                   `json:"pet_id" binding:"required"`
	Behaviors    []map[string]interface{} `json:"behaviors"`
	VoiceData    string                 `json:"voice_data"`
	Context      models.JSON            `json:"context"`
	TriggerEvent string                 `json:"trigger_event"`
}

// BatchRecognizeRequest 批量识别请求
type BatchRecognizeRequest struct {
	Items []BatchRecognizeItem `json:"items" binding:"required,dive"`
}

// BatchRecognizeItem 批量识别项
type BatchRecognizeItem struct {
	SubjectType string `json:"subject_type" binding:"required,oneof=user pet"`
	SubjectID   uint   `json:"subject_id" binding:"required"`
	Source      string `json:"source" binding:"required,oneof=voice text face behavior"`
	Data        string `json:"data"`
	RecordedAt  string `json:"recorded_at"`
}

// EmotionRecognitionResponse 情绪识别响应
type EmotionRecognitionResponse struct {
	EmotionType string       `json:"emotion_type"`
	Intensity   float64      `json:"intensity"`
	Confidence  float64      `json:"confidence"`
	Source      string       `json:"source"`
	Context     models.JSON  `json:"context"`
	RecordID    uint         `json:"record_id"`
}

// RecognizeVoice 语音情绪识别
func (c *EmotionController) RecognizeVoice(ctx *gin.Context) {
	var req RecognizeVoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	emotionType, intensity, confidence := c.simulateEmotionRecognition("voice", req.SubjectType)

	record := models.EmotionRecord{
		SubjectType: req.SubjectType,
		SubjectID:   req.SubjectID,
		EmotionType: emotionType,
		Intensity:   intensity,
		Source:      "voice",
		Confidence:  confidence,
		Context: models.JSON{
			"audio_url": req.AudioURL,
			"language":  req.Language,
		},
		RecordedAt: time.Now(),
	}

	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建情绪记录失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": EmotionRecognitionResponse{
			EmotionType: emotionType,
			Intensity:   intensity,
			Confidence:  confidence,
			Source:      "voice",
			Context:     record.Context,
			RecordID:    record.ID,
		},
	})
}

// RecognizeText 文字情绪识别
func (c *EmotionController) RecognizeText(ctx *gin.Context) {
	var req RecognizeTextRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	emotionType, intensity, confidence := c.analyzeTextEmotion(req.Text)

	record := models.EmotionRecord{
		SubjectType: req.SubjectType,
		SubjectID:   req.SubjectID,
		EmotionType: emotionType,
		Intensity:   intensity,
		Source:      "text",
		Confidence:  confidence,
		Context: models.JSON{
			"text_length": len(req.Text),
		},
		TriggerEvent: req.Text,
		RecordedAt:   time.Now(),
	}

	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建情绪记录失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": EmotionRecognitionResponse{
			EmotionType: emotionType,
			Intensity:   intensity,
			Confidence:  confidence,
			Source:      "text",
			Context:     record.Context,
			RecordID:    record.ID,
		},
	})
}

// RecognizePet 宠物情绪识别
func (c *EmotionController) RecognizePet(ctx *gin.Context) {
	var req RecognizePetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	emotionType, intensity, confidence := c.analyzePetEmotion(req.Behaviors, req.VoiceData)

	record := models.EmotionRecord{
		SubjectType: "pet",
		SubjectID:   req.PetID,
		EmotionType: emotionType,
		Intensity:   intensity,
		Source:      "behavior",
		Confidence:  confidence,
		Context:     req.Context,
		TriggerEvent: req.TriggerEvent,
		RecordedAt:  time.Now(),
	}

	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建情绪记录失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": EmotionRecognitionResponse{
			EmotionType: emotionType,
			Intensity:   intensity,
			Confidence:  confidence,
			Source:      "behavior",
			Context:     record.Context,
			RecordID:    record.ID,
		},
	})
}

// RecognizeBatch 批量识别
func (c *EmotionController) RecognizeBatch(ctx *gin.Context) {
	var req BatchRecognizeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	var results []EmotionRecognitionResponse

	for _, item := range req.Items {
		var emotionType string
		var intensity, confidence float64
		var recordedAt time.Time

		if item.RecordedAt != "" {
			parsed, err := time.Parse(time.RFC3339, item.RecordedAt)
			if err == nil {
				recordedAt = parsed
			} else {
				recordedAt = time.Now()
			}
		} else {
			recordedAt = time.Now()
		}

		switch item.Source {
		case "voice":
			emotionType, intensity, confidence = c.simulateEmotionRecognition("voice", item.SubjectType)
		case "text":
			emotionType, intensity, confidence = c.analyzeTextEmotion(item.Data)
		case "face":
			emotionType, intensity, confidence = c.simulateEmotionRecognition("face", item.SubjectType)
		case "behavior":
			emotionType, intensity, confidence = c.simulateEmotionRecognition("behavior", item.SubjectType)
		default:
			emotionType, intensity, confidence = "calm", 50.0, 0.5
		}

		record := models.EmotionRecord{
			SubjectType: item.SubjectType,
			SubjectID:   item.SubjectID,
			EmotionType: emotionType,
			Intensity:   intensity,
			Source:      item.Source,
			Confidence:  confidence,
			Context: models.JSON{
				"batch_id": time.Now().UnixNano(),
			},
			RecordedAt: recordedAt,
		}

		if err := c.DB.Create(&record).Error; err == nil {
			results = append(results, EmotionRecognitionResponse{
				EmotionType: emotionType,
				Intensity:   intensity,
				Confidence:  confidence,
				Source:      item.Source,
				RecordID:    record.ID,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":   len(req.Items),
			"success": len(results),
			"results": results,
		},
	})
}

// ==================== 情绪响应 API ====================

// GetResponseStrategies 获取响应策略
func (c *EmotionController) GetResponseStrategies(ctx *gin.Context) {
	petID, err := strconv.ParseUint(ctx.Param("pet_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的宠物ID"})
		return
	}

	var configs []models.EmotionResponseConfig
	if err := c.DB.Where("pet_id = ?", petID).Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取响应策略失败"})
		return
	}

	if len(configs) == 0 {
		configs = c.getDefaultResponseConfigs(uint(petID))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    configs,
	})
}

// TriggerEmotionResponse 触发情绪响应
func (c *EmotionController) TriggerEmotionResponse(ctx *gin.Context) {
	petID, err := strconv.ParseUint(ctx.Param("pet_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的宠物ID"})
		return
	}

	var req struct {
		EmotionType string       `json:"emotion_type" binding:"required"`
		Intensity   float64      `json:"intensity"`
		Context     models.JSON  `json:"context"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	var config models.EmotionResponseConfig
	err = c.DB.Where("pet_id = ? AND emotion_type = ? AND enabled = ?", petID, req.EmotionType, true).
		First(&config).Error

	responseData := gin.H{}
	if err == nil && config.Threshold <= req.Intensity {
		if config.LastTriggered != nil {
			elapsed := time.Since(*config.LastTriggered).Milliseconds()
			if elapsed < int64(config.Cooldown) {
				responseData["cooldown"] = true
				responseData["remaining_ms"] = config.Cooldown - int(elapsed)
			} else {
				now := time.Now()
				config.LastTriggered = &now
				c.DB.Save(&config)
				responseData["executed"] = true
				responseData["strategy"] = config.Strategy
				responseData["action_code"] = config.ActionCode
				responseData["action_param"] = config.ActionParam
			}
		} else {
			now := time.Now()
			config.LastTriggered = &now
			c.DB.Save(&config)
			responseData["executed"] = true
			responseData["strategy"] = config.Strategy
			responseData["action_code"] = config.ActionCode
			responseData["action_param"] = config.ActionParam
		}
	} else {
		responseData["executed"] = false
		responseData["reason"] = "below_threshold"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    responseData,
	})
}

// UpdateResponseConfig 更新响应配置
func (c *EmotionController) UpdateResponseConfig(ctx *gin.Context) {
	petID, err := strconv.ParseUint(ctx.Param("pet_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的宠物ID"})
		return
	}

	var req struct {
		Configs []struct {
			EmotionType   string       `json:"emotion_type" binding:"required"`
			Strategy      string       `json:"strategy" binding:"required"`
			ActionCode    string       `json:"action_code"`
			ActionParam   models.JSON  `json:"action_param"`
			ResponseDelay int          `json:"response_delay"`
			Threshold     float64      `json:"threshold"`
			Cooldown      int          `json:"cooldown"`
			Enabled       bool         `json:"enabled"`
		} `json:"configs" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	for _, cfg := range req.Configs {
		var existing models.EmotionResponseConfig
		err := c.DB.Where("pet_id = ? AND emotion_type = ?", petID, cfg.EmotionType).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			config := models.EmotionResponseConfig{
				PetID:          uint(petID),
				EmotionType:   cfg.EmotionType,
				Strategy:      cfg.Strategy,
				ActionCode:    cfg.ActionCode,
				ActionParam:   cfg.ActionParam,
				ResponseDelay: cfg.ResponseDelay,
				Threshold:     cfg.Threshold,
				Cooldown:      cfg.Cooldown,
				Enabled:       cfg.Enabled,
			}
			if err := c.DB.Create(&config).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建配置失败"})
				return
			}
		} else if err == nil {
			c.DB.Model(&existing).Updates(map[string]interface{}{
				"strategy":        cfg.Strategy,
				"action_code":     cfg.ActionCode,
				"action_param":    cfg.ActionParam,
				"response_delay":  cfg.ResponseDelay,
				"threshold":       cfg.Threshold,
				"cooldown":        cfg.Cooldown,
				"enabled":         cfg.Enabled,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ==================== 情绪日志 API ====================

// ListEmotionLogs 情绪日志列表
func (c *EmotionController) ListEmotionLogs(ctx *gin.Context) {
	var records []models.EmotionRecord
	query := c.DB.Model(&models.EmotionRecord{})

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if subjectType := ctx.Query("subject_type"); subjectType != "" {
		query = query.Where("subject_type = ?", subjectType)
	}
	if subjectID := ctx.Query("subject_id"); subjectID != "" {
		query = query.Where("subject_id = ?", subjectID)
	}
	if emotionType := ctx.Query("emotion_type"); emotionType != "" {
		query = query.Where("emotion_type = ?", emotionType)
	}
	if source := ctx.Query("source"); source != "" {
		query = query.Where("source = ?", source)
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("recorded_at >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("recorded_at <= ?", endDate)
	}

	var total int64
	query.Count(&total)

	if err := query.Order("recorded_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取情绪日志失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"items":     records,
		},
	})
}

// GetEmotionLog 获取情绪日志详情
func (c *EmotionController) GetEmotionLog(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的日志ID"})
		return
	}

	var record models.EmotionRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "日志不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取日志失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateEmotionLog 更新情绪日志
func (c *EmotionController) UpdateEmotionLog(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的日志ID"})
		return
	}

	var req struct {
		Tags      models.StringArray `json:"tags"`
		Note      string             `json:"note"`
		Intensity float64            `json:"intensity"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Note != "" {
		updates["note"] = req.Note
	}
	if req.Intensity > 0 {
		updates["intensity"] = req.Intensity
	}

	if len(updates) > 0 {
		c.DB.Model(&models.EmotionRecord{}).Where("id = ?", id).Updates(updates)
	}

	var updatedRecord models.EmotionRecord
	c.DB.First(&updatedRecord, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updatedRecord})
}

// DeleteEmotionLog 删除情绪日志
func (c *EmotionController) DeleteEmotionLog(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的日志ID"})
		return
	}

	c.DB.Delete(&models.EmotionRecord{}, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ==================== 情绪动作库 CRUD ====================

// ListEmotionActions 列出所有情绪动作
func (c *EmotionController) ListEmotionActions(ctx *gin.Context) {
	var actions []models.PetEmotionAction
	query := c.DB.Model(&models.PetEmotionAction{})

	if emotionType := ctx.Query("emotion_type"); emotionType != "" {
		query = query.Where("emotion_type = ?", emotionType)
	}
	if enabled := ctx.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	query.Order("priority DESC, emotion_type ASC").Find(&actions)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": actions})
}

// CreateEmotionAction 创建情绪动作
func (c *EmotionController) CreateEmotionAction(ctx *gin.Context) {
	var req struct {
		EmotionType   string          `json:"emotion_type" binding:"required"`
		ActionName    string          `json:"action_name" binding:"required"`
		ActionCode    string          `json:"action_code" binding:"required"`
		Description   string          `json:"description"`
		Parameters    models.JSON     `json:"parameters"`
		Priority      int             `json:"priority"`
		MinIntensity  float64         `json:"min_intensity"`
		MaxIntensity  float64         `json:"max_intensity"`
		Duration      int             `json:"duration"`
		Enabled       bool            `json:"enabled"`
		Icon          string          `json:"icon"`
		AnimationURL  string          `json:"animation_url"`
		SoundURL      string          `json:"sound_url"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	action := models.PetEmotionAction{
		EmotionType:   req.EmotionType,
		ActionName:    req.ActionName,
		ActionCode:    req.ActionCode,
		Description:   req.Description,
		Parameters:    req.Parameters,
		Priority:      req.Priority,
		MinIntensity:  req.MinIntensity,
		MaxIntensity:  req.MaxIntensity,
		Duration:      req.Duration,
		Enabled:       req.Enabled,
		Icon:          req.Icon,
		AnimationURL:  req.AnimationURL,
		SoundURL:      req.SoundURL,
	}

	if err := c.DB.Create(&action).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建情绪动作失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": action})
}

// GetEmotionAction 获取情绪动作详情
func (c *EmotionController) GetEmotionAction(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的动作ID"})
		return
	}

	var action models.PetEmotionAction
	if err := c.DB.First(&action, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "动作不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": action})
}

// UpdateEmotionAction 更新情绪动作
func (c *EmotionController) UpdateEmotionAction(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的动作ID"})
		return
	}

	var req struct {
		EmotionType   string          `json:"emotion_type"`
		ActionName    string          `json:"action_name"`
		Description   string          `json:"description"`
		Parameters    models.JSON     `json:"parameters"`
		Priority      int             `json:"priority"`
		MinIntensity  float64         `json:"min_intensity"`
		MaxIntensity  float64         `json:"max_intensity"`
		Duration      int             `json:"duration"`
		Enabled       bool            `json:"enabled"`
		Icon          string          `json:"icon"`
		AnimationURL  string          `json:"animation_url"`
		SoundURL      string          `json:"sound_url"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.EmotionType != "" { updates["emotion_type"] = req.EmotionType }
	if req.ActionName != "" { updates["action_name"] = req.ActionName }
	if req.Description != "" { updates["description"] = req.Description }
	if req.Parameters != nil { updates["parameters"] = req.Parameters }
	if req.Priority != 0 { updates["priority"] = req.Priority }
	if req.MinIntensity > 0 { updates["min_intensity"] = req.MinIntensity }
	if req.MaxIntensity > 0 { updates["max_intensity"] = req.MaxIntensity }
	if req.Duration > 0 { updates["duration"] = req.Duration }
	updates["enabled"] = req.Enabled
	if req.Icon != "" { updates["icon"] = req.Icon }
	if req.AnimationURL != "" { updates["animation_url"] = req.AnimationURL }
	if req.SoundURL != "" { updates["sound_url"] = req.SoundURL }

	c.DB.Model(&models.PetEmotionAction{}).Where("id = ?", id).Updates(updates)

	var action models.PetEmotionAction
	c.DB.First(&action, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": action})
}

// DeleteEmotionAction 删除情绪动作
func (c *EmotionController) DeleteEmotionAction(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的动作ID"})
		return
	}

	c.DB.Delete(&models.PetEmotionAction{}, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ==================== 情绪报告 API ====================

// ListEmotionReports 情绪报告列表
func (c *EmotionController) ListEmotionReports(ctx *gin.Context) {
	var reports []models.EmotionReport
	query := c.DB.Model(&models.EmotionReport{})

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if reportType := ctx.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("generated_at DESC").Offset(offset).Limit(pageSize).Find(&reports)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"items":     reports,
		},
	})
}

// GetEmotionReport 获取情绪报告详情
func (c *EmotionController) GetEmotionReport(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的报告ID"})
		return
	}

	var report models.EmotionReport
	if err := c.DB.First(&report, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "报告不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// GenerateEmotionReport 生成情绪报告
func (c *EmotionController) GenerateEmotionReport(ctx *gin.Context) {
	var req struct {
		PetID      uint   `json:"pet_id" binding:"required"`
		ReportType string `json:"report_type" binding:"required,oneof=daily weekly monthly"`
		StartDate  string `json:"start_date" binding:"required"`
		EndDate    string `json:"end_date" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error": err.Error()})
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	endDate = endDate.Add(24*time.Hour - time.Second)

	var records []models.EmotionRecord
	c.DB.Where("subject_type = 'pet' AND subject_id = ? AND recorded_at >= ? AND recorded_at <= ?",
		req.PetID, startDate, endDate).Find(&records)

	// 计算情绪统计
	emotionStatsMap := make(map[string]int)
	for _, r := range records {
		emotionStatsMap[r.EmotionType]++
	}
	emotionStats := make(map[string]interface{})
	for k, v := range emotionStatsMap {
		emotionStats[k] = v
	}

	report := models.EmotionReport{
		PetID:         req.PetID,
		ReportType:    req.ReportType,
		StartDate:     startDate,
		EndDate:       endDate,
		Summary: models.JSON{
			"total_records": len(records),
		},
		EmotionStats: models.JSON(emotionStats),
		GeneratedAt:  time.Now(),
	}

	c.DB.Create(&report)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// GetFamilyEmotionMap 家庭情绪地图
func (c *EmotionController) GetFamilyEmotionMap(ctx *gin.Context) {
	var records []models.EmotionRecord
	c.DB.Where("subject_type = 'pet'").
		Order("recorded_at DESC").
		Limit(100).
		Find(&records)

	emotionMap := make(map[uint]models.EmotionRecord)
	for _, r := range records {
		if _, ok := emotionMap[r.SubjectID]; !ok {
			emotionMap[r.SubjectID] = r
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": emotionMap})
}

// ==================== 辅助方法 ====================

func (c *EmotionController) simulateEmotionRecognition(source, subjectType string) (string, float64, float64) {
	emotions := []string{"happy", "calm", "anxious", "sad"}
	emotion := emotions[time.Now().Unix()%int64(len(emotions))]
	return emotion, 50.0 + float64(time.Now().Unix()%50), 0.7 + float64(time.Now().Unix()%30)/100
}

func (c *EmotionController) analyzeTextEmotion(text string) (string, float64, float64) {
	if len(text) < 5 {
		return "calm", 50.0, 0.6
	}
	emotions := []string{"happy", "sad", "angry", "calm"}
	emotion := emotions[int(len(text))%len(emotions)]
	return emotion, 60.0, 0.75
}

func (c *EmotionController) analyzePetEmotion(behaviors []map[string]interface{}, voiceData string) (string, float64, float64) {
	if len(behaviors) == 0 && voiceData == "" {
		return "calm", 50.0, 0.5
	}
	emotions := []string{"happy", "excited", "anxious", "tired"}
	emotion := emotions[time.Now().Unix()%int64(len(emotions))]
	return emotion, 65.0, 0.8
}

func (c *EmotionController) getDefaultResponseConfigs(petID uint) []models.EmotionResponseConfig {
	defaults := []struct {
		emotionType string
		strategy    string
		actionCode  string
	}{
		{"happy", "action", "play_music"},
		{"sad", "action", "comfort语音"},
		{"anxious", "action", "calming_music"},
		{"calm", "none", ""},
	}

	var configs []models.EmotionResponseConfig
	for _, d := range defaults {
		configs = append(configs, models.EmotionResponseConfig{
			PetID:        petID,
			EmotionType:  d.emotionType,
			Strategy:     d.strategy,
			ActionCode:   d.actionCode,
			Threshold:    30,
			Cooldown:     60000,
			Enabled:      true,
		})
	}
	return configs
}
