package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VoiceEmotionController struct {
	DB *gorm.DB
}

func NewVoiceEmotionController(db *gorm.DB) *VoiceEmotionController {
	return &VoiceEmotionController{DB: db}
}

func (ctrl *VoiceEmotionController) RegisterRoutes(rg *gin.RouterGroup) {
	voice := rg.Group("/voice-emotion")
	{
		voice.POST("/upload", ctrl.UploadAudio)
		voice.POST("/analyze", ctrl.AnalyzeEmotion)
		voice.GET("/records", ctrl.GetRecords)
		voice.GET("/records/:id", ctrl.GetRecord)
	}
}

// UploadAudio 上传音频并分析
func (ctrl *VoiceEmotionController) UploadAudio(c *gin.Context) {
	var req struct {
		PetID    uint   `json:"pet_id" binding:"required"`
		UserID   uint   `json:"user_id" binding:"required"`
		AudioURL string `json:"audio_url" binding:"required"`
		Duration int    `json:"duration"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 模拟AI情绪分析结果
	record := models.VoiceEmotionRecord{
		PetID:       req.PetID,
		UserID:      req.UserID,
		AudioURL:    req.AudioURL,
		Duration:    req.Duration,
		EmotionType: "calm",
		Intensity:   6,
		Confidence:  0.85,
		Transcript:  "狗狗发出愉快的叫声",
	}
	ctrl.DB.Create(&record)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}

// AnalyzeEmotion 分析情绪
func (ctrl *VoiceEmotionController) AnalyzeEmotion(c *gin.Context) {
	var req struct {
		AudioURL string `json:"audio_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 模拟AI分析
	analysis := gin.H{
		"emotion":    "happy",
		"intensity":  7,
		"confidence": 0.88,
		"transcript": "宠物叫声愉悦",
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": analysis})
}

// GetRecords 获取语音情绪历史
func (ctrl *VoiceEmotionController) GetRecords(c *gin.Context) {
	petID := c.Query("pet_id")
	userID := c.Query("user_id")
	emotionType := c.Query("emotion_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := ctrl.DB.Model(&models.VoiceEmotionRecord{})
	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if emotionType != "" {
		query = query.Where("emotion_type = ?", emotionType)
	}

	var total int64
	query.Count(&total)

	var records []models.VoiceEmotionRecord
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

// GetRecord 获取记录详情
func (ctrl *VoiceEmotionController) GetRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.VoiceEmotionRecord
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}
