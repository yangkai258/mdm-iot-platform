package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmotionController 情感计算控制器
type EmotionController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册情绪相关路由
func (e *EmotionController) RegisterRoutes(r *gin.RouterGroup) {
	em := r.Group("/emotions")
	{
		// 情绪记录
		em.GET("/records", e.GetRecords)
		em.POST("/records", e.CreateRecord)
		em.GET("/records/stats", e.GetRecordStats)

		// 宠物情绪响应动作
		em.GET("/pet/actions", e.GetPetActions)
		em.POST("/pet/actions", e.CreatePetAction)
		em.PUT("/pet/actions/:id", e.UpdatePetAction)
		em.DELETE("/pet/actions/:id", e.DeletePetAction)

		// 情绪响应配置
		em.GET("/pet/config", e.GetEmotionConfig)
		em.PUT("/pet/config", e.UpdateEmotionConfig)

		// 情绪报告
		em.GET("/report", e.GetEmotionReport)
	}
}

// GetRecords 获取情绪记录列表
// @Summary 获取情绪记录列表
// @Tags emotions
// @Accept json
// @Produce json
// @Param pet_id query int false "宠物ID"
// @Param user_id query int false "用户ID"
// @Param emotion_type query string false "情绪类型"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Router /api/v1/emotions/records [GET]
func (e *EmotionController) GetRecords(c *gin.Context) {
	var records []models.EmotionRecord
	query := e.DB.Model(&models.EmotionRecord{})

	// 过滤条件
	if petID := c.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if emotionType := c.Query("emotion_type"); emotionType != "" {
		query = query.Where("emotion_type = ?", emotionType)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("recorded_at >= ?", t)
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("recorded_at <= ?", t.Add(24*time.Hour))
		}
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("recorded_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  records,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreateRecord 创建情绪记录
// @Summary 创建情绪记录
// @Tags emotions
// @Accept json
// @Produce json
// @Param body body CreateEmotionRecordRequest true "请求体"
// @Router /api/v1/emotions/records [POST]
func (e *EmotionController) CreateRecord(c *gin.Context) {
	var req CreateEmotionRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	record := models.EmotionRecord{
		PetID:       req.PetID,
		UserID:      req.UserID,
		Source:      req.Source,
		EmotionType: req.EmotionType,
		Intensity:   req.Intensity,
		Trigger:     req.Trigger,
		Context:     req.Context,
		AIResponse:  req.AIResponse,
		RecordedAt:  time.Now(),
	}

	if err := e.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}

// GetRecordStats 获取情绪统计
// @Summary 获取情绪统计
// @Tags emotions
// @Accept json
// @Produce json
// @Param pet_id query int true "宠物ID"
// @Param period query string false "统计周期: day, week, month"
// @Router /api/v1/emotions/records/stats [GET]
func (e *EmotionController) GetRecordStats(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}
	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	period := c.DefaultQuery("period", "week")
	var startDate time.Time
	now := time.Now()

	switch period {
	case "day":
		startDate = now.AddDate(0, 0, -1)
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, 0, -7)
	}

	var records []models.EmotionRecord
	e.DB.Where("pet_id = ? AND recorded_at >= ?", petID, startDate).
		Order("recorded_at ASC").
		Find(&records)

	// 统计各情绪类型分布
	emotionCount := make(map[string]int)
	var totalIntensity float64
	for _, r := range records {
		emotionCount[r.EmotionType]++
		totalIntensity += float64(r.Intensity)
	}

	// 找出主导情绪
	dominantEmotion := ""
	maxCount := 0
	for emotion, count := range emotionCount {
		if count > maxCount {
			maxCount = count
			dominantEmotion = emotion
		}
	}

	avgIntensity := 0.0
	if len(records) > 0 {
		avgIntensity = totalIntensity / float64(len(records))
	}

	// 计算趋势（对比上期）
	var prevCount int64
	e.DB.Model(&models.EmotionRecord{}).
		Where("pet_id = ? AND recorded_at >= ? AND recorded_at < ?", petID, startDate.Add(-startDate.Sub(now))).
		Count(&prevCount)

	trend := "stable"
	if int(prevCount) < len(records) {
		trend = "improving"
	} else if int(prevCount) > len(records) {
		trend = "declining"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period":            period,
			"total_count":       len(records),
			"emotion_distribution": emotionCount,
			"avg_intensity":     avgIntensity,
			"dominant_emotion":  dominantEmotion,
			"trend":             trend,
		},
	})
}

// GetPetActions 获取宠物情绪响应动作列表
// @Summary 获取宠物情绪响应动作列表
// @Tags emotions
// @Param pet_id query int true "宠物ID"
// @Router /api/v1/emotions/pet/actions [GET]
func (e *EmotionController) GetPetActions(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}
	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	var actions []models.PetEmotionAction
	if err := e.DB.Where("pet_id = ?", petID).Order("priority DESC, id ASC").Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": actions})
}

// CreatePetAction 创建宠物情绪响应动作
// @Summary 创建宠物情绪响应动作
// @Tags emotions
// @Accept json
// @Param body body CreatePetEmotionActionRequest true "请求体"
// @Router /api/v1/emotions/pet/actions [POST]
func (e *EmotionController) CreatePetAction(c *gin.Context) {
	var req CreatePetEmotionActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	action := models.PetEmotionAction{
		PetID:       req.PetID,
		EmotionType: req.EmotionType,
		ActionType:  req.ActionType,
		ActionData:  req.ActionData,
		Priority:    req.Priority,
		Enabled:     true,
	}

	if err := e.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": action})
}

// UpdatePetAction 更新宠物情绪响应动作
// @Summary 更新宠物情绪响应动作
// @Tags emotions
// @Accept json
// @Param id path int true "动作ID"
// @Param body body UpdatePetEmotionActionRequest true "请求体"
// @Router /api/v1/emotions/pet/actions/{id} [PUT]
func (e *EmotionController) UpdatePetAction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	var action models.PetEmotionAction
	if err := e.DB.First(&action, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req UpdatePetEmotionActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 只更新非零值字段
	updates := make(map[string]interface{})
	if req.EmotionType != "" {
		updates["emotion_type"] = req.EmotionType
	}
	if req.ActionType != "" {
		updates["action_type"] = req.ActionType
	}
	if req.ActionData != "" {
		updates["action_data"] = req.ActionData
	}
	if req.Priority != 0 {
		updates["priority"] = req.Priority
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := e.DB.Model(&action).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": action})
}

// DeletePetAction 删除宠物情绪响应动作
// @Summary 删除宠物情绪响应动作
// @Tags emotions
// @Param id path int true "动作ID"
// @Router /api/v1/emotions/pet/actions/{id} [DELETE]
func (e *EmotionController) DeletePetAction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	if err := e.DB.Delete(&models.PetEmotionAction{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetEmotionConfig 获取情绪响应配置
// @Summary 获取情绪响应配置
// @Tags emotions
// @Param pet_id query int true "宠物ID"
// @Router /api/v1/emotions/pet/config [GET]
func (e *EmotionController) GetEmotionConfig(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}
	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	var configs []models.EmotionResponseConfig
	if err := e.DB.Where("pet_id = ?", petID).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": configs})
}

// UpdateEmotionConfig 更新情绪响应配置
// @Summary 更新情绪响应配置
// @Tags emotions
// @Accept json
// @Param body body UpdateEmotionConfigRequest true "请求体"
// @Router /api/v1/emotions/pet/config [PUT]
func (e *EmotionController) UpdateEmotionConfig(c *gin.Context) {
	var req UpdateEmotionConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找或创建配置
	var config models.EmotionResponseConfig
	err := e.DB.Where("pet_id = ? AND emotion_type = ?", req.PetID, req.EmotionType).First(&config).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新配置
		enabledVal := true
		if req.Enabled != nil {
			enabledVal = *req.Enabled
		}
		config = models.EmotionResponseConfig{
			PetID:        req.PetID,
			EmotionType:  req.EmotionType,
			ResponseMode: req.ResponseMode,
			Actions:      req.Actions,
			Enabled:      enabledVal,
		}
		if err := e.DB.Create(&config).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	} else {
		// 更新现有配置
		updates := make(map[string]interface{})
		if req.ResponseMode != "" {
			updates["response_mode"] = req.ResponseMode
		}
		if req.Actions != "" {
			updates["actions"] = req.Actions
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}
		if err := e.DB.Model(&config).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
}

// GetEmotionReport 生成情绪报告
// @Summary 生成情绪报告
// @Tags emotions
// @Param pet_id query int true "宠物ID"
// @Param user_id query int false "用户ID"
// @Param period query string false "周期: daily, weekly, monthly"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Router /api/v1/emotions/report [GET]
func (e *EmotionController) GetEmotionReport(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}
	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	userIDStr := c.Query("user_id")
	var userID uint
	if userIDStr != "" {
		uid, _ := strconv.ParseUint(userIDStr, 10, 32)
		userID = uint(uid)
	}

	period := c.DefaultQuery("period", "weekly")

	var startDate, endDate time.Time
	now := time.Now()

	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = t
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = t.Add(24 * time.Hour)
		}
	}

	// 默认周期
	if startDate.IsZero() || endDate.IsZero() {
		switch period {
		case "daily":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.Add(24 * time.Hour)
		case "weekly":
			startDate = now.AddDate(0, 0, -7)
			endDate = now
		case "monthly":
			startDate = now.AddDate(0, -1, 0)
			endDate = now
		default:
			startDate = now.AddDate(0, 0, -7)
			endDate = now
		}
	}

	// 查询情绪记录
	query := e.DB.Model(&models.EmotionRecord{}).Where("pet_id = ? AND recorded_at >= ? AND recorded_at <= ?", petID, startDate, endDate)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	var records []models.EmotionRecord
	query.Find(&records)

	// 统计
	emotionCount := make(map[string]int)
	var totalIntensity float64
	for _, r := range records {
		emotionCount[r.EmotionType]++
		totalIntensity += float64(r.Intensity)
	}

	dominantEmotion := ""
	maxCount := 0
	for emotion, count := range emotionCount {
		if count > maxCount {
			maxCount = count
			dominantEmotion = emotion
		}
	}

	avgIntensity := 0.0
	if len(records) > 0 {
		avgIntensity = totalIntensity / float64(len(records))
	}

	// 计算趋势
	var prevCount int64
	prevQuery := e.DB.Model(&models.EmotionRecord{}).Where("pet_id = ? AND recorded_at < ?", petID, startDate)
	if userID > 0 {
		prevQuery = prevQuery.Where("user_id = ?", userID)
	}
	prevQuery.Count(&prevCount)

	trend := "stable"
	if int(prevCount) < len(records) {
		trend = "improving"
	} else if int(prevCount) > len(records) {
		trend = "declining"
	}

	// 序列化摘要
	summaryJSON, _ := json.Marshal(map[string]interface{}{
		"emotion_distribution": emotionCount,
		"record_count":          len(records),
	})

	report := models.EmotionReport{
		PetID:           uint(petID),
		UserID:          uint(userID),
		Period:          period,
		StartDate:       startDate,
		EndDate:         endDate,
		Summary:         string(summaryJSON),
		AvgIntensity:    avgIntensity,
		DominantEmotion: dominantEmotion,
		Trend:           trend,
		CreatedAt:       now,
	}

	// 保存或更新报告
	var existingReport models.EmotionReport
	err := e.DB.Where("pet_id = ? AND period = ? AND start_date = ? AND end_date = ?", petID, period, startDate, endDate).First(&existingReport).Error
	if err == gorm.ErrRecordNotFound {
		e.DB.Create(&report)
	} else if err == nil {
		report.ID = existingReport.ID
		e.DB.Save(&report)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": report})
}

// ========== 请求结构体 ==========

type CreateEmotionRecordRequest struct {
	PetID       uint   `json:"pet_id"`
	UserID      uint   `json:"user_id"`
	Source      string `json:"source"`
	EmotionType string `json:"emotion_type"`
	Intensity   int    `json:"intensity"`
	Trigger     string `json:"trigger"`
	Context     string `json:"context"`
	AIResponse  string `json:"ai_response"`
}

type CreatePetEmotionActionRequest struct {
	PetID       uint   `json:"pet_id"`
	EmotionType string `json:"emotion_type"`
	ActionType  string `json:"action_type"`
	ActionData  string `json:"action_data"`
	Priority    int    `json:"priority"`
}

type UpdatePetEmotionActionRequest struct {
	EmotionType string `json:"emotion_type"`
	ActionType  string `json:"action_type"`
	ActionData  string `json:"action_data"`
	Priority    int    `json:"priority"`
	Enabled     *bool  `json:"enabled"`
}

type UpdateEmotionConfigRequest struct {
	PetID        uint   `json:"pet_id"`
	EmotionType  string `json:"emotion_type"`
	ResponseMode string `json:"response_mode"`
	Actions      string `json:"actions"`
	Enabled      *bool  `json:"enabled"`
}
