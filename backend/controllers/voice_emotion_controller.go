package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// VoiceEmotionController 语音情绪控制器
type VoiceEmotionController struct {
	DB *gorm.DB
}

// NewVoiceEmotionController 创建控制器实例
func NewVoiceEmotionController(db *gorm.DB) *VoiceEmotionController {
	return &VoiceEmotionController{DB: db}
}

// RegisterRoutes 注册语音情绪路由
func (c *VoiceEmotionController) RegisterRoutes(api *gin.RouterGroup) {
	voice := api.Group("/voice-emotion")
	{
		voice.GET("/records", c.ListRecords)
		voice.GET("/records/:id", c.GetRecord)
		voice.POST("/records", c.CreateRecord)
		voice.PUT("/records/:id", c.UpdateRecord)
		voice.DELETE("/records/:id", c.DeleteRecord)
		voice.GET("/sessions/:session_id/records", c.GetSessionRecords)
		voice.GET("/analytics/daily", c.GetDailyAnalytics)
		voice.GET("/analytics/weekly", c.GetWeeklyAnalytics)
		voice.GET("/analytics/monthly", c.GetMonthlyAnalytics)
		voice.GET("/emotions/distribution", c.GetEmotionDistribution)
		voice.POST("/process", c.ProcessVoiceEmotion)
	}
}

// ListRecords 获取语音情绪记录列表
func (c *VoiceEmotionController) ListRecords(ctx *gin.Context) {
	var records []models.VoiceEmotionRecord
	var total int64

	query := c.DB.Model(&models.VoiceEmotionRecord{})

	// 用户筛选
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 宠物筛选
	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}

	// 设备筛选
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	// 情绪类型筛选
	if emotionType := ctx.Query("emotion_type"); emotionType != "" {
		query = query.Where("emotion_type = ?", emotionType)
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 时间范围筛选
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRecord 获取语音情绪记录详情
func (c *VoiceEmotionController) GetRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.VoiceEmotionRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// CreateRecord 创建语音情绪记录
func (c *VoiceEmotionController) CreateRecord(ctx *gin.Context) {
	var record models.VoiceEmotionRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	record.Status = "processing"
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateRecord 更新语音情绪记录
func (c *VoiceEmotionController) UpdateRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.VoiceEmotionRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		EmotionType  string  `json:"emotion_type"`
		EmotionScore float64 `json:"emotion_score"`
		Intensity    int     `json:"intensity"`
		Valence      float64 `json:"valence"`
		Arousal      float64 `json:"arousal"`
		Transcript   string  `json:"transcript"`
		Keywords     string  `json:"keywords"`
		Sentiment    string  `json:"sentiment"`
		AIResponse   string  `json:"ai_response"`
		ResponseType string  `json:"response_type"`
		Status       string  `json:"status"`
		ErrorMessage string  `json:"error_message"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.EmotionType != "" {
		updates["emotion_type"] = updateData.EmotionType
	}
	if updateData.EmotionScore > 0 {
		updates["emotion_score"] = updateData.EmotionScore
	}
	if updateData.Intensity > 0 {
		updates["intensity"] = updateData.Intensity
	}
	if updateData.Valence != 0 {
		updates["valence"] = updateData.Valence
	}
	if updateData.Arousal != 0 {
		updates["arousal"] = updateData.Arousal
	}
	if updateData.Transcript != "" {
		updates["transcript"] = updateData.Transcript
	}
	if updateData.Keywords != "" {
		updates["keywords"] = updateData.Keywords
	}
	if updateData.Sentiment != "" {
		updates["sentiment"] = updateData.Sentiment
	}
	if updateData.AIResponse != "" {
		updates["ai_response"] = updateData.AIResponse
	}
	if updateData.ResponseType != "" {
		updates["response_type"] = updateData.ResponseType
	}
	if updateData.Status != "" {
		updates["status"] = updateData.Status
	}
	if updateData.ErrorMessage != "" {
		updates["error_message"] = updateData.ErrorMessage
	}

	if err := c.DB.Model(&record).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&record, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeleteRecord 删除语音情绪记录
func (c *VoiceEmotionController) DeleteRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.VoiceEmotionRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	if err := c.DB.Delete(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetSessionRecords 获取会话的所有记录
func (c *VoiceEmotionController) GetSessionRecords(ctx *gin.Context) {
	sessionID := ctx.Param("session_id")
	var records []models.VoiceEmotionRecord
	var total int64

	query := c.DB.Model(&models.VoiceEmotionRecord{}).Where("session_id = ?", sessionID)
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at ASC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDailyAnalytics 获取每日情绪分析
func (c *VoiceEmotionController) GetDailyAnalytics(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	date := ctx.DefaultQuery("date", time.Now().Format("2006-01-02"))

	var records []models.VoiceEmotionRecord
	c.DB.Where("user_id = ? AND DATE(created_at) = ? AND status = ?", userID, date, "completed").
		Order("created_at ASC").Find(&records)

	// 计算统计数据
	emotionCount := make(map[string]int)
	var totalScore float64
	var totalIntensity float64
	var count int

	for _, r := range records {
		if r.EmotionType != "" {
			emotionCount[r.EmotionType]++
		}
		totalScore += r.EmotionScore
		totalIntensity += float64(r.Intensity)
		count++
	}

	avgScore := 0.0
	avgIntensity := 0.0
	if count > 0 {
		avgScore = totalScore / float64(count)
		avgIntensity = totalIntensity / float64(count)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"date":          date,
			"total_records": count,
			"avg_score":     avgScore,
			"avg_intensity": avgIntensity,
			"emotion_distribution": emotionCount,
		},
	})
}

// GetWeeklyAnalytics 获取每周情绪分析
func (c *VoiceEmotionController) GetWeeklyAnalytics(ctx *gin.Context) {
	userID := ctx.Query("user_id")

	// 获取最近7天的数据
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	var records []models.VoiceEmotionRecord
	c.DB.Where("user_id = ? AND created_at >= ? AND created_at <= ? AND status = ?",
		userID, startDate, endDate, "completed").Find(&records)

	// 按日期分组
	dailyStats := make(map[string]gin.H)
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		dailyStats[date] = gin.H{"count": 0, "emotions": make(map[string]int)}
	}

	for _, r := range records {
		date := r.CreatedAt.Format("2006-01-02")
		if stats, ok := dailyStats[date]; ok {
			stats["count"] = stats["count"].(int) + 1
			emotions := stats["emotions"].(map[string]int)
			emotions[r.EmotionType]++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
			"daily_stats": dailyStats,
		},
	})
}

// GetMonthlyAnalytics 获取每月情绪分析
func (c *VoiceEmotionController) GetMonthlyAnalytics(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	month := ctx.DefaultQuery("month", time.Now().Format("2006-01"))

	startDate, _ := time.Parse("2006-01", month)
	endDate := startDate.AddDate(0, 1, 0)

	var records []models.VoiceEmotionRecord
	c.DB.Where("user_id = ? AND created_at >= ? AND created_at < ? AND status = ?",
		userID, startDate, endDate, "completed").Find(&records)

	// 统计
	emotionCount := make(map[string]int)
	emotionScores := make(map[string][]float64)
	var totalScore float64
	var totalIntensity float64
	var count int

	for _, r := range records {
		if r.EmotionType != "" {
			emotionCount[r.EmotionType]++
			emotionScores[r.EmotionType] = append(emotionScores[r.EmotionType], r.EmotionScore)
		}
		totalScore += r.EmotionScore
		totalIntensity += float64(r.Intensity)
		count++
	}

	// 计算每种情绪的平均分
	emotionAvgScores := make(map[string]float64)
	for emotion, scores := range emotionScores {
		var sum float64
		for _, s := range scores {
			sum += s
		}
		emotionAvgScores[emotion] = sum / float64(len(scores))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"month":              month,
			"total_records":      count,
			"avg_score":          totalScore / float64(count),
			"avg_intensity":      totalIntensity / float64(count),
			"emotion_distribution": emotionCount,
			"emotion_avg_scores": emotionAvgScores,
		},
	})
}

// GetEmotionDistribution 获取情绪分布
func (c *VoiceEmotionController) GetEmotionDistribution(ctx *gin.Context) {
	userID := ctx.Query("user_id")

	var records []models.VoiceEmotionRecord
	query := c.DB.Model(&models.VoiceEmotionRecord{}).Where("status = ?", "completed")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	query.Find(&records)

	emotionCount := make(map[string]int)
	var total int64

	for _, r := range records {
		if r.EmotionType != "" {
			emotionCount[r.EmotionType]++
			total++
		}
	}

	// 计算百分比
	emotionPercent := make(map[string]float64)
	for emotion, count := range emotionCount {
		emotionPercent[emotion] = float64(count) / float64(total) * 100
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"total":       total,
			"distribution": emotionCount,
			"percentage":  emotionPercent,
		},
	})
}

// ProcessVoiceEmotion 处理语音情绪分析
func (c *VoiceEmotionController) ProcessVoiceEmotion(ctx *gin.Context) {
	var processData struct {
		RecordID uint   `json:"record_id" binding:"required"`
		AudioURL string `json:"audio_url"`
	}
	if err := ctx.ShouldBindJSON(&processData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var record models.VoiceEmotionRecord
	if err := c.DB.First(&record, processData.RecordID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 模拟AI处理
	// 实际项目中这里会调用AI服务进行语音情绪分析
	updates := map[string]interface{}{
		"status":        "completed",
		"emotion_type":  "calm",
		"emotion_score": 75.5,
		"intensity":     5,
		"valence":       0.2,
		"arousal":       0.1,
		"transcript":    "Sample transcript from voice",
		"sentiment":     "positive",
		"ai_response":   "听起来你今天心情不错！",
		"response_type": "encourage",
	}

	if err := c.DB.Model(&record).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "处理失败"})
		return
	}

	c.DB.First(&record, record.ID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}
