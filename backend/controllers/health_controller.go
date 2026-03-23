package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthController 健康医疗控制器
type HealthController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// NewHealthController 创建健康控制器
func NewHealthController(db *gorm.DB, redis *utils.RedisClient) *HealthController {
	return &HealthController{DB: db, Redis: redis}
}

// RegisterRoutes 注册健康相关路由
func (ctrl *HealthController) RegisterRoutes(rg *gin.RouterGroup) {
	health := rg.Group("/health")
	{
		// 运动统计
		health.GET("/exercise/stats", ctrl.GetExerciseStats)
		health.GET("/exercise/history", ctrl.GetExerciseHistory)

		// 睡眠分析
		health.GET("/sleep/stats", ctrl.GetSleepStats)
		health.GET("/sleep/history", ctrl.GetSleepHistory)

		// 健康预警
		health.GET("/warnings", ctrl.ListWarnings)
		health.POST("/warnings/:id/confirm", ctrl.ConfirmWarning)
		health.POST("/warnings/:id/ignore", ctrl.IgnoreWarning)

		// 健康报告
		health.GET("/report", ctrl.GetHealthReport)

		// 体重追踪
		health.GET("/body-weight/history", ctrl.GetBodyWeightHistory)
		health.POST("/body-weight/record", ctrl.RecordBodyWeight)

		// 饮食记录
		health.GET("/diet/history", ctrl.GetDietHistory)
		health.POST("/diet/record", ctrl.RecordDiet)
	}
}

// getDateRange 根据 range 参数返回开始和结束时间
func getDateRange(rangeType string) (startDate, endDate time.Time) {
	now := time.Now()
	endDate = now
	switch rangeType {
	case "day":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		startDate = now.AddDate(0, 0, -7)
	}
	return
}

// parseTimeParam 解析时间参数
func parseTimeParam(c *gin.Context, param string) time.Time {
	if t := c.Query(param); t != "" {
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

// ==================== 运动统计 API ====================

// GetExerciseStats 获取运动统计
func (ctrl *HealthController) GetExerciseStats(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	memberIDStr := c.Query("member_id")
	rangeType := c.DefaultQuery("range", "week")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	_, _ = strconv.ParseUint(memberIDStr, 10, 32) // member_id for filtering

	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	startDate, endDate := getDateRange(rangeType)

	var stats struct {
		TotalSteps    int     `json:"total_steps"`
		TotalDistance float64 `json:"total_distance"`
		TotalCalories int     `json:"total_calories"`
		TotalDuration int     `json:"total_duration"`
		AvgSteps      int     `json:"avg_steps"`
		AvgDistance   float64 `json:"avg_distance"`
		GoalAchieved  int     `json:"goal_achieved"` // days goal was met
		RecordCount   int     `json:"record_count"`
	}

	rows := ctrl.DB.Model(&models.ExerciseRecord{}).
		Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Select("COALESCE(SUM(steps),0) as total_steps, COALESCE(SUM(distance),0) as total_distance, COALESCE(SUM(calories),0) as total_calories, COALESCE(SUM(duration),0) as total_duration, COUNT(*) as record_count").
		Row()

	rows.Scan(&stats.TotalSteps, &stats.TotalDistance, &stats.TotalCalories, &stats.TotalDuration, &stats.RecordCount)

	if stats.RecordCount > 0 {
		stats.AvgSteps = stats.TotalSteps / stats.RecordCount
		stats.AvgDistance = stats.TotalDistance / float64(stats.RecordCount)
	}

	// Count goal achievements
	var goalAchieved int64
	ctrl.DB.Model(&models.ExerciseRecord{}).
		Where("pet_id = ? AND date >= ? AND date <= ? AND goal > 0 AND steps >= goal", petID, startDate, endDate).
		Count(&goalAchieved)
	stats.GoalAchieved = int(goalAchieved)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"range":          rangeType,
			"start_date":     startDate.Format("2006-01-02"),
			"end_date":       endDate.Format("2006-01-02"),
			"total_steps":    stats.TotalSteps,
			"total_distance": stats.TotalDistance,
			"total_calories": stats.TotalCalories,
			"total_duration": stats.TotalDuration,
			"avg_steps":      stats.AvgSteps,
			"avg_distance":   stats.AvgDistance,
			"goal_achieved":  stats.GoalAchieved,
			"record_count":   stats.RecordCount,
		},
	})
}

// GetExerciseHistory 获取运动历史
func (ctrl *HealthController) GetExerciseHistory(c *gin.Context) {
	petIDStr := c.Query("pet_id")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	startDate := parseTimeParam(c, "start_date")
	endDate := parseTimeParam(c, "end_date")
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	var records []models.ExerciseRecord
	query := ctrl.DB.Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Order("date DESC")

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	var total int64
	query.Model(&models.ExerciseRecord{}).Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch exercise history", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"records": records,
			"total":   total,
			"page":    page,
			"page_size": pageSize,
		},
	})
}

// ==================== 睡眠分析 API ====================

// GetSleepStats 获取睡眠统计
func (ctrl *HealthController) GetSleepStats(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	rangeType := c.DefaultQuery("range", "week")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	startDate, endDate := getDateRange(rangeType)

	var stats struct {
		TotalDeepSleep  int `json:"total_deep_sleep"`
		TotalLightSleep int `json:"total_light_sleep"`
		TotalREMSleep   int `json:"total_rem_sleep"`
		TotalSleep      int `json:"total_sleep"`
		AvgDeepSleep    int `json:"avg_deep_sleep"`
		AvgLightSleep   int `json:"avg_light_sleep"`
		AvgREMSleep     int `json:"avg_rem_sleep"`
		AvgTotalSleep   int `json:"avg_total_sleep"`
		RecordCount     int `json:"record_count"`
	}

	rows := ctrl.DB.Model(&models.SleepRecord{}).
		Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Select("COALESCE(SUM(deep_sleep),0) as total_deep_sleep, COALESCE(SUM(light_sleep),0) as total_light_sleep, COALESCE(SUM(rem_sleep),0) as total_rem_sleep, COALESCE(SUM(total_sleep),0) as total_sleep, COUNT(*) as record_count").
		Row()

	rows.Scan(&stats.TotalDeepSleep, &stats.TotalLightSleep, &stats.TotalREMSleep, &stats.TotalSleep, &stats.RecordCount)

	if stats.RecordCount > 0 {
		stats.AvgDeepSleep = stats.TotalDeepSleep / stats.RecordCount
		stats.AvgLightSleep = stats.TotalLightSleep / stats.RecordCount
		stats.AvgREMSleep = stats.TotalREMSleep / stats.RecordCount
		stats.AvgTotalSleep = stats.TotalSleep / stats.RecordCount
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"range":            rangeType,
			"start_date":       startDate.Format("2006-01-02"),
			"end_date":         endDate.Format("2006-01-02"),
			"total_deep_sleep": stats.TotalDeepSleep,
			"total_light_sleep": stats.TotalLightSleep,
			"total_rem_sleep":  stats.TotalREMSleep,
			"total_sleep":      stats.TotalSleep,
			"avg_deep_sleep":   stats.AvgDeepSleep,
			"avg_light_sleep":  stats.AvgLightSleep,
			"avg_rem_sleep":    stats.AvgREMSleep,
			"avg_total_sleep":  stats.AvgTotalSleep,
			"record_count":     stats.RecordCount,
		},
	})
}

// GetSleepHistory 获取睡眠历史
func (ctrl *HealthController) GetSleepHistory(c *gin.Context) {
	petIDStr := c.Query("pet_id")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	startDate := parseTimeParam(c, "start_date")
	endDate := parseTimeParam(c, "end_date")
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	var records []models.SleepRecord
	query := ctrl.DB.Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Order("date DESC")

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	var total int64
	query.Model(&models.SleepRecord{}).Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch sleep history", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"records":   records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ==================== 健康预警 API ====================

// ListWarnings 获取预警列表
func (ctrl *HealthController) ListWarnings(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	level := c.Query("level")
	status := c.DefaultQuery("status", "")

	query := ctrl.DB.Model(&models.HealthWarning{})

	if petIDStr != "" {
		petID, _ := strconv.ParseUint(petIDStr, 10, 32)
		query = query.Where("pet_id = ?", petID)
	}

	if level != "" {
		query = query.Where("level = ?", level)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	var total int64
	query.Model(&models.HealthWarning{}).Count(&total)

	var warnings []models.HealthWarning
	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&warnings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch warnings", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"warnings":  warnings,
			"total":    total,
			"page":     page,
			"page_size": pageSize,
		},
	})
}

// ConfirmWarning 确认预警
func (ctrl *HealthController) ConfirmWarning(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid warning id"})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&models.HealthWarning{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "confirmed",
		"confirmed_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to confirm warning", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Warning confirmed"})
}

// IgnoreWarning 忽略预警
func (ctrl *HealthController) IgnoreWarning(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid warning id"})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&models.HealthWarning{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "ignored",
		"ignored_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to ignore warning", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Warning ignored"})
}

// ==================== 健康报告 API ====================

// GetHealthReport 获取健康报告
func (ctrl *HealthController) GetHealthReport(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	memberIDStr := c.Query("member_id")
	reportType := c.DefaultQuery("type", "weekly")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	memberID, _ := strconv.ParseUint(memberIDStr, 10, 32)

	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	var startDate, endDate time.Time
	now := time.Now()
	if reportType == "monthly" {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		startDate = now.AddDate(0, 0, -7)
	}
	endDate = now

	// 运动统计
	var exerciseStats struct {
		TotalSteps    int
		TotalDistance float64
		TotalCalories int
	}
	ctrl.DB.Model(&models.ExerciseRecord{}).
		Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Select("COALESCE(SUM(steps),0), COALESCE(SUM(distance),0), COALESCE(SUM(calories),0)").
		Row().Scan(&exerciseStats.TotalSteps, &exerciseStats.TotalDistance, &exerciseStats.TotalCalories)

	// 睡眠统计
	var sleepStats struct {
		TotalSleep int
		AvgSleep   int
	}
	ctrl.DB.Model(&models.SleepRecord{}).
		Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Select("COALESCE(SUM(total_sleep),0), COALESCE(AVG(total_sleep),0)").
		Row().Scan(&sleepStats.TotalSleep, &sleepStats.AvgSleep)

	// 体重变化
	var latestWeight, earliestWeight float64
	ctrl.DB.Model(&models.BodyWeightRecord{}).
		Where("pet_id = ? AND recorded_at >= ? AND recorded_at <= ?", petID, startDate, endDate).
		Order("recorded_at DESC").Limit(1).Select("weight").Row().Scan(&latestWeight)
	ctrl.DB.Model(&models.BodyWeightRecord{}).
		Where("pet_id = ? AND recorded_at >= ? AND recorded_at <= ?", petID, startDate, endDate).
		Order("recorded_at ASC").Limit(1).Select("weight").Row().Scan(&earliestWeight)
	weightChange := latestWeight - earliestWeight

	// 饮食统计
	var totalCalories int
	ctrl.DB.Model(&models.DietRecord{}).
		Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate).
		Select("COALESCE(SUM(calories),0)").Row().Scan(&totalCalories)

	// 预警统计
	var warningCount int64
	ctrl.DB.Model(&models.HealthWarning{}).
		Where("pet_id = ? AND created_at >= ? AND created_at <= ?", petID, startDate, endDate).
		Count(&warningCount)

	report := gin.H{
		"report_type":    reportType,
		"pet_id":         petID,
		"member_id":      memberID,
		"start_date":     startDate.Format("2006-01-02"),
		"end_date":       endDate.Format("2006-01-02"),
		"exercise": gin.H{
			"total_steps":    exerciseStats.TotalSteps,
			"total_distance": exerciseStats.TotalDistance,
			"total_calories": exerciseStats.TotalCalories,
		},
		"sleep": gin.H{
			"total_minutes": sleepStats.TotalSleep,
			"avg_minutes":   sleepStats.AvgSleep,
		},
		"body_weight": gin.H{
			"change": weightChange,
		},
		"diet": gin.H{
			"total_calories": totalCalories,
		},
		"warnings_count": warningCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": report,
	})
}

// ==================== 体重追踪 API ====================

// GetBodyWeightHistory 获取体重历史
func (ctrl *HealthController) GetBodyWeightHistory(c *gin.Context) {
	petIDStr := c.Query("pet_id")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	startDate := parseTimeParam(c, "start_date")
	endDate := parseTimeParam(c, "end_date")
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -90)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	var records []models.BodyWeightRecord
	query := ctrl.DB.Where("pet_id = ? AND recorded_at >= ? AND recorded_at <= ?", petID, startDate, endDate).
		Order("recorded_at DESC")

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	var total int64
	query.Model(&models.BodyWeightRecord{}).Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch body weight history", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"records":   records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RecordBodyWeight 记录体重
func (ctrl *HealthController) RecordBodyWeight(c *gin.Context) {
	var req struct {
		PetID      uint      `json:"pet_id" binding:"required"`
		Weight     float64   `json:"weight" binding:"required"`
		RecordedAt time.Time `json:"recorded_at"`
		Note       string    `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request", "error": err.Error()})
		return
	}

	if req.RecordedAt.IsZero() {
		req.RecordedAt = time.Now()
	}

	record := models.BodyWeightRecord{
		PetID:      req.PetID,
		Weight:     req.Weight,
		RecordedAt: req.RecordedAt,
		Note:       req.Note,
	}

	if err := ctrl.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to record body weight", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Body weight recorded", "data": record})
}

// ==================== 饮食记录 API ====================

// GetDietHistory 获取饮食历史
func (ctrl *HealthController) GetDietHistory(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	mealType := c.Query("meal_type")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)
	if petID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id is required"})
		return
	}

	startDate := parseTimeParam(c, "start_date")
	endDate := parseTimeParam(c, "end_date")
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -30)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	var records []models.DietRecord
	query := ctrl.DB.Where("pet_id = ? AND date >= ? AND date <= ?", petID, startDate, endDate)

	if mealType != "" {
		query = query.Where("meal_type = ?", mealType)
	}

	query = query.Order("date DESC")

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	var total int64
	query.Model(&models.DietRecord{}).Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch diet history", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"records":   records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RecordDiet 记录饮食
func (ctrl *HealthController) RecordDiet(c *gin.Context) {
	var req struct {
		PetID      uint      `json:"pet_id" binding:"required"`
		MealType   string    `json:"meal_type" binding:"required"` // breakfast, lunch, dinner, snack
		FoodType   string    `json:"food_type"`
		Amount     float64   `json:"amount"`
		Calories   int       `json:"calories"`
		Note       string    `json:"note"`
		RecordedAt time.Time `json:"recorded_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request", "error": err.Error()})
		return
	}

	recordedAt := req.RecordedAt
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}

	record := models.DietRecord{
		PetID:     req.PetID,
		Date:      recordedAt,
		MealType:  req.MealType,
		FoodType:  req.FoodType,
		Amount:    req.Amount,
		Calories:  req.Calories,
		Note:      req.Note,
	}

	if err := ctrl.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to record diet", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Diet recorded", "data": record})
}
