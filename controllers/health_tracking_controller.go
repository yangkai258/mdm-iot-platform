package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthTrackingCtrl 健康追踪控制器
type HealthTrackingCtrl struct {
	DB *gorm.DB
}

// RegisterHealthRoutes 注册健康追踪路由
func (h *HealthTrackingCtrl) RegisterHealthRoutes(r *gin.RouterGroup) {
	// 早期疾病预警
	r.GET("/health/:pet_id/early-warning", h.ListEarlyWarnings)
	r.GET("/health/:pet_id/early-warning/:id", h.GetEarlyWarning)
	r.POST("/health/:pet_id/early-warning/:id/ack", h.AckEarlyWarning)
	r.POST("/health/:pet_id/early-warning/:id/dismiss", h.DismissEarlyWarning)
	r.PUT("/health/:pet_id/early-warning/:id", h.UpdateEarlyWarning)
	r.DELETE("/health/:pet_id/early-warning/:id", h.DeleteEarlyWarning)

	// 运动追踪
	r.GET("/health/:pet_id/exercise", h.ListExerciseRecords)
	r.POST("/health/:pet_id/exercise", h.ReportExercise)
	r.GET("/health/:pet_id/exercise/summary", h.GetExerciseSummary)
	r.GET("/health/:pet_id/exercise/:id", h.GetExerciseRecord)
	r.PUT("/health/:pet_id/exercise/:id", h.UpdateExerciseRecord)
	r.DELETE("/health/:pet_id/exercise/:id", h.DeleteExerciseRecord)

	// 睡眠分析
	r.GET("/health/:pet_id/sleep", h.ListSleepRecords)
	r.POST("/health/:pet_id/sleep", h.ReportSleep)
	r.GET("/health/:pet_id/sleep/analysis", h.GetSleepAnalysis)
	r.GET("/health/:pet_id/sleep/:id", h.GetSleepRecord)
	r.PUT("/health/:pet_id/sleep/:id", h.UpdateSleepRecord)
	r.DELETE("/health/:pet_id/sleep/:id", h.DeleteSleepRecord)

	// 健康报告
	r.GET("/health/:pet_id/report", h.GetHealthReport)
	r.GET("/health/:pet_id/report/weekly", h.GetWeeklyReport)
	r.GET("/health/:pet_id/report/monthly", h.GetMonthlyReport)
}

// verifyPetOwnership 验证宠物归属
func (h *HealthTrackingCtrl) verifyPetOwnership(petID string, userID uint, tenantID string) (*models.Pet, error) {
	var pet models.Pet
	err := h.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error
	return &pet, err
}

// parsePageParams 解析分页参数
func parsePageParams(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return
}

// ==================== 早期疾病预警 API ====================

// ListEarlyWarnings 获取预警列表
func (h *HealthTrackingCtrl) ListEarlyWarnings(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	_, err := h.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "pet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}

	page, pageSize := parsePageParams(c)
	offset := (page - 1) * pageSize

	var warnings []models.HealthWarning
	var total int64

	query := h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	if err := query.Order("severity DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&warnings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"records": warnings, "total": total, "page": page, "page_size": pageSize},
	})
}

// GetEarlyWarning 获取预警详情
func (h *HealthTrackingCtrl) GetEarlyWarning(c *gin.Context) {
	petID := c.Param("pet_id")
	warningID := c.Param("id")
	tenantID := getTenantID(c)

	var warning models.HealthWarning
	if err := h.DB.Where("warning_uuid = ? AND pet_uuid = ? AND tenant_id = ?", warningID, petID, tenantID).First(&warning).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "warning not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": warning})
}

// AckEarlyWarning 确认预警
func (h *HealthTrackingCtrl) AckEarlyWarning(c *gin.Context) {
	petID := c.Param("pet_id")
	warningID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	now := time.Now()
	result := h.DB.Model(&models.HealthWarning{}).
		Where("warning_uuid = ? AND pet_uuid = ? AND tenant_id = ?", warningID, petID, tenantID).
		Updates(map[string]interface{}{"status": models.WarningStatusAcked, "acked_at": now, "acked_by": userID})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "ack failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "warning not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "warning acked", "data": gin.H{"acked_at": now}})
}

// DismissEarlyWarning 忽略预警
func (h *HealthTrackingCtrl) DismissEarlyWarning(c *gin.Context) {
	petID := c.Param("pet_id")
	warningID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct{ Reason string `json:"reason"` }
	c.ShouldBindJSON(&input)

	now := time.Now()
	result := h.DB.Model(&models.HealthWarning{}).
		Where("warning_uuid = ? AND pet_uuid = ? AND tenant_id = ?", warningID, petID, tenantID).
		Updates(map[string]interface{}{"status": models.WarningStatusDismissed, "dismissed_at": now, "dismissed_by": userID, "dismiss_reason": input.Reason})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "dismiss failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "warning not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "warning dismissed", "data": gin.H{"dismissed_at": now}})
}

// UpdateEarlyWarning 更新预警
func (h *HealthTrackingCtrl) UpdateEarlyWarning(c *gin.Context) {
	petID := c.Param("pet_id")
	warningID := c.Param("id")
	tenantID := getTenantID(c)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params"})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.Priority != 0 {
		updates["priority"] = input.Priority
	}

	result := h.DB.Model(&models.HealthWarning{}).
		Where("warning_uuid = ? AND pet_uuid = ? AND tenant_id = ?", warningID, petID, tenantID).Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "update failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "warning not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "update success"})
}

// DeleteEarlyWarning 删除预警
func (h *HealthTrackingCtrl) DeleteEarlyWarning(c *gin.Context) {
	petID := c.Param("pet_id")
	warningID := c.Param("id")
	tenantID := getTenantID(c)

	result := h.DB.Where("warning_uuid = ? AND pet_uuid = ? AND tenant_id = ?", warningID, petID, tenantID).Delete(&models.HealthWarning{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "delete failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "warning not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "delete success"})
}

// ==================== 运动追踪 API ====================

// ListExerciseRecords 获取运动记录列表
func (h *HealthTrackingCtrl) ListExerciseRecords(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	page, pageSize := parsePageParams(c)
	offset := (page - 1) * pageSize

	var records []models.ExerciseRecord
	var total int64

	query := h.DB.Model(&models.ExerciseRecord{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)

	if exerciseType := c.Query("type"); exerciseType != "" {
		query = query.Where("exercise_type = ?", exerciseType)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("start_time >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("start_time <= ?", endDate)
	}

	query.Count(&total)
	if err := query.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"records": records, "total": total, "page": page, "page_size": pageSize},
	})
}

// ReportExercise 上报运动数据
func (h *HealthTrackingCtrl) ReportExercise(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	var input models.ReqExerciseReport
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params"})
		return
	}

	startTime, _ := time.Parse(time.RFC3339, input.StartTime)
	var endTime *time.Time
	if input.EndTime != "" {
		et, _ := time.Parse(time.RFC3339, input.EndTime)
		endTime = &et
	}

	record := models.ExerciseRecord{
		PetUUID: petID, ExerciseType: input.ExerciseType, StartTime: startTime, EndTime: endTime,
		Duration: input.Duration, DurationMinutes: input.Duration, Steps: input.Steps,
		Distance: input.Distance, DistanceUnit: "km", CaloriesBurned: input.CaloriesBurned,
		AvgHeartRate: input.AvgHeartRate, MaxHeartRate: input.MaxHeartRate, MinHeartRate: input.MinHeartRate,
		AvgSpeed: input.AvgSpeed, MaxSpeed: input.MaxSpeed, Intensity: input.Intensity,
		Notes: input.Notes, DataSource: "api", TenantID: tenantID,
	}
	if input.Tags != nil {
		record.Tags = input.Tags
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "create failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// GetExerciseSummary 获取运动汇总统计
func (h *HealthTrackingCtrl) GetExerciseSummary(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var dailySummary models.ExerciseSummary
	h.DB.Where("pet_uuid = ? AND summary_date = ? AND tenant_id = ?", petID, today, tenantID).First(&dailySummary)

	var weeklySummary models.ExerciseSummary
	h.DB.Model(&models.ExerciseSummary{}).
		Where("pet_uuid = ? AND summary_date >= ? AND summary_date <= ? AND tenant_id = ?", petID, weekStart, today, tenantID).
		Select("SUM(total_duration) as total_duration, SUM(total_steps) as total_steps, SUM(total_distance) as total_distance, SUM(total_calories) as total_calories, COUNT(*) as exercise_count").
		Scan(&weeklySummary)

	var monthlySummary models.ExerciseSummary
	h.DB.Model(&models.ExerciseSummary{}).
		Where("pet_uuid = ? AND summary_date >= ? AND summary_date <= ? AND tenant_id = ?", petID, monthStart, today, tenantID).
		Select("SUM(total_duration) as total_duration, SUM(total_steps) as total_steps, SUM(total_distance) as total_distance, SUM(total_calories) as total_calories, COUNT(*) as exercise_count").
		Scan(&monthlySummary)

	var trends []models.ExerciseTrend
	h.DB.Model(&models.ExerciseSummary{}).
		Select("summary_date as date, total_duration, total_steps, total_distance, total_calories").
		Where("pet_uuid = ? AND summary_date >= ? AND summary_date <= ? AND tenant_id = ?", petID, today.AddDate(0, 0, -6), today, tenantID).
		Order("summary_date ASC").Find(&trends)

	var goals []models.ExerciseGoal
	h.DB.Where("pet_uuid = ? AND status = 'active' AND tenant_id = ?", petID, tenantID).Find(&goals)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": models.RespExerciseSummary{
		Daily: &dailySummary, Weekly: &weeklySummary, Monthly: &monthlySummary, Trends: trends, Goals: goals,
	}})
}

// GetExerciseRecord 获取运动记录详情
func (h *HealthTrackingCtrl) GetExerciseRecord(c *gin.Context) {
	petID := c.Param("pet_id")
	recordID := c.Param("id")
	tenantID := getTenantID(c)

	var record models.ExerciseRecord
	if err := h.DB.Where("record_uuid = ? AND pet_uuid = ? AND tenant_id = ?", recordID, petID, tenantID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateExerciseRecord 更新运动记录
func (h *HealthTrackingCtrl) UpdateExerciseRecord(c *gin.Context) {
	petID := c.Param("pet_id")
	recordID := c.Param("id")
	tenantID := getTenantID(c)

	var input struct {
		ExerciseType   string  `json:"exercise_type"`
		Duration       int     `json:"duration"`
		Steps          int     `json:"steps"`
		Distance       float64 `json:"distance"`
		CaloriesBurned float64 `json:"calories_burned"`
		Notes          string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params"})
		return
	}

	updates := map[string]interface{}{}
	if input.ExerciseType != "" {
		updates["exercise_type"] = input.ExerciseType
	}
	if input.Duration != 0 {
		updates["duration"] = input.Duration
		updates["duration_minutes"] = input.Duration
	}
	if input.Steps != 0 {
		updates["steps"] = input.Steps
	}
	if input.Distance != 0 {
		updates["distance"] = input.Distance
	}
	if input.CaloriesBurned != 0 {
		updates["calories_burned"] = input.CaloriesBurned
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}

	result := h.DB.Model(&models.ExerciseRecord{}).
		Where("record_uuid = ? AND pet_uuid = ? AND tenant_id = ?", recordID, petID, tenantID).Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "update failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "update success"})
}

// DeleteExerciseRecord 删除运动记录
func (h *HealthTrackingCtrl) DeleteExerciseRecord(c *gin.Context) {
	petID := c.Param("pet_id")
	recordID := c.Param("id")
	tenantID := getTenantID(c)

	result := h.DB.Where("record_uuid = ? AND pet_uuid = ? AND tenant_id = ?", recordID, petID, tenantID).Delete(&models.ExerciseRecord{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "delete failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "delete success"})
}

// ==================== 睡眠分析 API ====================

// ListSleepRecords 获取睡眠记录列表
func (h *HealthTrackingCtrl) ListSleepRecords(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	page, pageSize := parsePageParams(c)
	offset := (page - 1) * pageSize

	var records []models.SleepRecord
	var total int64

	query := h.DB.Model(&models.SleepRecord{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)

	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("sleep_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("sleep_date <= ?", endDate)
	}
	if isNap := c.Query("is_nap"); isNap != "" {
		query = query.Where("is_nap = ?", isNap == "true")
	}

	query.Count(&total)
	if err := query.Order("sleep_date DESC, bed_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"records": records, "total": total, "page": page, "page_size": pageSize},
	})
}

// ReportSleep 上报睡眠数据
func (h *HealthTrackingCtrl) ReportSleep(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	var input models.ReqSleepReport
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params"})
		return
	}

	var bedTime, wakeTime *time.Time
	if input.BedTime != "" {
		bt, _ := time.Parse(time.RFC3339, input.BedTime)
		bedTime = &bt
	}
	if input.WakeTime != "" {
		wt, _ := time.Parse(time.RFC3339, input.WakeTime)
		wakeTime = &wt
	}

	qualityLevel := models.SleepQualityVeryPoor
	if input.QualityScore >= 90 {
		qualityLevel = models.SleepQualityExcellent
	} else if input.QualityScore >= 75 {
		qualityLevel = models.SleepQualityGood
	} else if input.QualityScore >= 60 {
		qualityLevel = models.SleepQualityFair
	} else if input.QualityScore >= 40 {
		qualityLevel = models.SleepQualityPoor
	}

	record := models.SleepRecord{
		PetUUID: petID, SleepDate: time.Now(), BedTime: bedTime, WakeTime: wakeTime,
		TotalDuration: input.TotalDuration, TotalDurationMinutes: input.TotalDuration,
		REMCount: input.REMCount, REMTime: input.REMTime,
		LightSleepCount: input.LightSleepCount, LightSleepTime: input.LightSleepTime,
		DeepSleepCount: input.DeepSleepCount, DeepSleepTime: input.DeepSleepTime,
		CoreSleepCount: input.CoreSleepCount, CoreSleepTime: input.CoreSleepTime,
		AwakeCount: input.AwakeCount, AwakeTime: input.AwakeTime,
		AvgHeartRate: input.AvgHeartRate, MinHeartRate: input.MinHeartRate,
		AvgRespiratoryRate: input.AvgRespiratoryRate,
		QualityScore: input.QualityScore, QualityLevel: qualityLevel,
		SleepEfficiency: input.SleepEfficiency, Latency: input.Latency,
		Restlessness: input.Restlessness, Notes: input.Notes,
		DataSource: "api", IsNap: input.IsNap, TenantID: tenantID,
	}
	if input.Tags != nil {
		record.Tags = input.Tags
	}
	if input.StageDetails != nil {
		record.StageDetails = input.StageDetails
	}
	if input.Environment.Temperature != 0 || input.Environment.Humidity != 0 {
		record.Environment = map[string]interface{}{
			"temperature": input.Environment.Temperature,
			"humidity":    input.Environment.Humidity,
			"noise_level": input.Environment.NoiseLevel,
		}
	}

	if err := h.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "create failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// GetSleepAnalysis 获取睡眠分析报告
func (h *HealthTrackingCtrl) GetSleepAnalysis(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	now := time.Now()
	startDate := now.AddDate(0, 0, -days)
	endDate := now

	var dailyRecords []models.SleepRecord
	h.DB.Where("pet_uuid = ? AND sleep_date >= ? AND sleep_date <= ? AND tenant_id = ?", petID, startDate, endDate, tenantID).
		Order("sleep_date DESC").Find(&dailyRecords)

	var trends []models.SleepTrend
	h.DB.Model(&models.SleepRecord{}).
		Select("sleep_date as date, total_duration, quality_score, quality_level, deep_sleep_time, rem_time, awake_time").
		Where("pet_uuid = ? AND sleep_date >= ? AND sleep_date <= ? AND tenant_id = ?", petID, startDate, endDate, tenantID).
		Order("sleep_date ASC").Find(&trends)

	analysis := h.generateSleepAnalysis(petID, tenantID, startDate, endDate, dailyRecords)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": models.RespSleepAnalysis{
		Analysis: &analysis, DailyRecords: dailyRecords, Trends: trends,
	}})
}

// generateSleepAnalysis 生成睡眠分析报告
func (h *HealthTrackingCtrl) generateSleepAnalysis(petID, tenantID string, startDate, endDate time.Time, records []models.SleepRecord) models.SleepAnalysis {
	var analysis models.SleepAnalysis
	analysis.PetUUID = petID
	analysis.TenantID = tenantID
	analysis.AnalysisType = "weekly"
	analysis.StartDate = startDate
	analysis.EndDate = endDate
	analysis.TotalRecordDays = len(records)

	if len(records) == 0 {
		return analysis
	}

	var totalDuration, totalREM, totalLight, totalDeep, totalCore, totalAwake int
	var totalQualityScore, totalEfficiency, totalRestlessness float64
	bestScore := 0
	worstScore := 100
	var bestDate, worstDate *time.Time

	for _, r := range records {
		totalDuration += r.TotalDuration
		totalREM += r.REMTime
		totalLight += r.LightSleepTime
		totalDeep += r.DeepSleepTime
		totalCore += r.CoreSleepTime
		totalAwake += r.AwakeTime
		totalQualityScore += float64(r.QualityScore)
		totalEfficiency += r.SleepEfficiency
		totalRestlessness += r.Restlessness

		if r.QualityScore > bestScore {
			bestScore = r.QualityScore
			bd := r.SleepDate
			bestDate = &bd
		}
		if r.QualityScore < worstScore {
			worstScore = r.QualityScore
			wd := r.SleepDate
			worstDate = &wd
		}
	}

	count := len(records)
	analysis.AvgTotalDuration = totalDuration / count
	analysis.AvgREMTime = totalREM / count
	analysis.AvgLightSleepTime = totalLight / count
	analysis.AvgDeepSleepTime = totalDeep / count
	analysis.AvgCoreSleepTime = totalCore / count
	analysis.AvgAwakeTime = totalAwake / count
	analysis.AvgQualityScore = totalQualityScore / float64(count)
	analysis.AvgSleepEfficiency = totalEfficiency / float64(count)
	analysis.AvgRestlessness = totalRestlessness / float64(count)
	analysis.BestSleepDate = bestDate
	analysis.BestSleepQuality = bestScore
	analysis.WorstSleepDate = worstDate
	analysis.WorstSleepQuality = worstScore
	analysis.TotalRecordDays = count

	var issues []map[string]interface{}
	if analysis.AvgDeepSleepTime < 60 {
		issues = append(issues, map[string]interface{}{"type": "low_deep_sleep", "description": "deep sleep insufficient"})
	}
	if analysis.AvgQualityScore < 60 {
		issues = append(issues, map[string]interface{}{"type": "low_quality", "description": "sleep quality low"})
	}
	if analysis.AvgRestlessness > 30 {
		issues = append(issues, map[string]interface{}{"type": "high_restlessness", "description": "restlessness high during sleep"})
	}
	analysis.IssuesDetected = issues

	var recommendations []map[string]interface{}
	if analysis.AvgDeepSleepTime < 60 {
		recommendations = append(recommendations, map[string]interface{}{"category": "deep_sleep", "suggestion": "increase daily exercise to improve deep sleep"})
	}
	if analysis.AvgQualityScore < 60 {
		recommendations = append(recommendations, map[string]interface{}{"category": "quality", "suggestion": "maintain regular sleep schedule, avoid intense activity before bed"})
	}
	analysis.Recommendations = recommendations

	return analysis
}

// GetSleepRecord 获取睡眠记录详情
func (h *HealthTrackingCtrl) GetSleepRecord(c *gin.Context) {
	petID := c.Param("pet_id")
	recordID := c.Param("id")
	tenantID := getTenantID(c)

	var record models.SleepRecord
	if err := h.DB.Where("record_uuid = ? AND pet_uuid = ? AND tenant_id = ?", recordID, petID, tenantID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateSleepRecord 更新睡眠记录
func (h *HealthTrackingCtrl) UpdateSleepRecord(c *gin.Context) {
	petID := c.Param("pet_id")
	recordID := c.Param("id")
	tenantID := getTenantID(c)

	var input struct {
		TotalDuration   int     `json:"total_duration"`
		QualityScore    int     `json:"quality_score"`
		QualityLevel    string  `json:"quality_level"`
		SleepEfficiency float64 `json:"sleep_efficiency"`
		Notes           string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params"})
		return
	}

	updates := map[string]interface{}{}
	if input.TotalDuration != 0 {
		updates["total_duration"] = input.TotalDuration
		updates["total_duration_minutes"] = input.TotalDuration
	}
	if input.QualityScore != 0 {
		updates["quality_score"] = input.QualityScore
	}
	if input.QualityLevel != "" {
		updates["quality_level"] = input.QualityLevel
	}
	if input.SleepEfficiency != 0 {
		updates["sleep_efficiency"] = input.SleepEfficiency
	}
	if input.Notes != "" {
		updates["notes"] = input.Notes
	}

	result := h.DB.Model(&models.SleepRecord{}).
		Where("record_uuid = ? AND pet_uuid = ? AND tenant_id = ?", recordID, petID, tenantID).Updates(updates)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "update failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "update success"})
}

// DeleteSleepRecord 删除睡眠记录
func (h *HealthTrackingCtrl) DeleteSleepRecord(c *gin.Context) {
	petID := c.Param("pet_id")
	recordID := c.Param("id")
	tenantID := getTenantID(c)

	result := h.DB.Where("record_uuid = ? AND pet_uuid = ? AND tenant_id = ?", recordID, petID, tenantID).Delete(&models.SleepRecord{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "delete failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "delete success"})
}

// ==================== 健康报告 API ====================

// RespHealthReport 健康报告响应
type RespHealthReport struct {
	PetUUID        string                     `json:"pet_uuid"`
	ReportDate     time.Time                  `json:"report_date"`
	Overview       HealthOverview             `json:"overview"`
	VitalSigns     []models.VitalRecord       `json:"vital_signs"`
	Exercise       *models.RespExerciseSummary `json:"exercise"`
	Sleep          *models.RespSleepAnalysis   `json:"sleep"`
	Warnings       []models.HealthWarning     `json:"warnings"`
	Alerts         []models.HealthAlert       `json:"alerts"`
	Recommendations []map[string]interface{}  `json:"recommendations"`
}

// HealthOverview 健康总览
type HealthOverview struct {
	HealthScore   int    `json:"health_score"`
	HealthLevel   string `json:"health_level"`
	ActivityLevel string `json:"activity_level"`
	SleepQuality  string `json:"sleep_quality"`
	RiskLevel     string `json:"risk_level"`
	Trend         string `json:"trend"`
}

// GetHealthReport 获取综合健康报告
func (h *HealthTrackingCtrl) GetHealthReport(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	_, err := h.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "pet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "query failed"})
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var vitalSigns []models.VitalRecord
	h.DB.Where("pet_uuid = ? AND recorded_at >= ? AND tenant_id = ?", petID, today.AddDate(0, 0, -7), tenantID).
		Order("recorded_at DESC").Limit(20).Find(&vitalSigns)

	var exerciseSummary models.RespExerciseSummary
	var dailyExercise models.ExerciseSummary
	h.DB.Where("pet_uuid = ? AND summary_date = ? AND tenant_id = ?", petID, today, tenantID).First(&dailyExercise)
	exerciseSummary.Daily = &dailyExercise

	var sleepAnalysis models.RespSleepAnalysis
	var sleepRecords []models.SleepRecord
	h.DB.Where("pet_uuid = ? AND sleep_date >= ? AND tenant_id = ?", petID, today.AddDate(0, 0, -7), tenantID).
		Order("sleep_date DESC").Find(&sleepRecords)
	sleepAnalysis.DailyRecords = sleepRecords

	// Get active warnings
	var warnings []models.HealthWarning
	h.DB.Where("pet_uuid = ? AND status IN ? AND tenant_id = ?", petID, []string{"active", "acked"}, tenantID).
		Order("severity DESC, created_at DESC").Limit(10).Find(&warnings)

	// Get health alerts
	var alerts []models.HealthAlert
	h.DB.Where("pet_uuid = ? AND alert_time >= ? AND tenant_id = ?", petID, today.AddDate(0, 0, -7), tenantID).
		Order("alert_time DESC").Limit(10).Find(&alerts)

	// Calculate health score
	healthScore := 85
	if len(warnings) > 0 {
		for _, w := range warnings {
			if w.Level == models.WarningLevelCritical || w.Level == models.WarningLevelEmergency {
				healthScore -= 20
			} else if w.Level == models.WarningLevelWarning {
				healthScore -= 10
			}
		}
	}
	if healthScore < 0 {
		healthScore = 0
	}

	healthLevel := "excellent"
	if healthScore < 70 {
		healthLevel = "poor"
	} else if healthScore < 85 {
		healthLevel = "good"
	}

	report := RespHealthReport{
		PetUUID:    petID,
		ReportDate: now,
		Overview: HealthOverview{
			HealthScore:   healthScore,
			HealthLevel:   healthLevel,
			ActivityLevel: "normal",
			SleepQuality:  "good",
			RiskLevel:     "low",
			Trend:         "stable",
		},
		VitalSigns: vitalSigns,
		Exercise:   &exerciseSummary,
		Sleep:      &sleepAnalysis,
		Warnings:   warnings,
		Alerts:     alerts,
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// GetWeeklyReport 获取周报
func (h *HealthTrackingCtrl) GetWeeklyReport(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	now := time.Now()
	weekStart := now.AddDate(0, 0, -7)

	// Weekly exercise summary
	var weeklyExercise models.ExerciseSummary
	h.DB.Model(&models.ExerciseSummary{}).
		Where("pet_uuid = ? AND summary_date >= ? AND summary_date <= ? AND tenant_id = ?", petID, weekStart, now, tenantID).
		Select("SUM(total_duration) as total_duration, SUM(total_steps) as total_steps, SUM(total_distance) as total_distance, SUM(total_calories) as total_calories, COUNT(*) as exercise_count").
		Scan(&weeklyExercise)

	// Weekly sleep analysis
	var sleepRecords []models.SleepRecord
	h.DB.Where("pet_uuid = ? AND sleep_date >= ? AND sleep_date <= ? AND is_nap = ? AND tenant_id = ?", petID, weekStart, now, false, tenantID).
		Find(&sleepRecords)

	var totalSleepDuration, totalDeepSleep, totalREMSleep int
	var totalQualityScore float64
	for _, r := range sleepRecords {
		totalSleepDuration += r.TotalDuration
		totalDeepSleep += r.DeepSleepTime
		totalREMSleep += r.REMTime
		totalQualityScore += float64(r.QualityScore)
	}
	avgSleepDuration := 0
	avgDeepSleep := 0
	avgREMSleep := 0
	avgQualityScore := float64(0)
	if len(sleepRecords) > 0 {
		avgSleepDuration = totalSleepDuration / len(sleepRecords)
		avgDeepSleep = totalDeepSleep / len(sleepRecords)
		avgREMSleep = totalREMSleep / len(sleepRecords)
		avgQualityScore = totalQualityScore / float64(len(sleepRecords))
	}

	// Weekly warning stats
	var warningStats struct {
		Total    int64
		Active   int64
		Critical int64
	}
	h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND created_at >= ? AND tenant_id = ?", petID, weekStart, tenantID).Count(&warningStats.Total)
	h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND created_at >= ? AND status = 'active' AND tenant_id = ?", petID, weekStart, tenantID).Count(&warningStats.Active)
	h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND created_at >= ? AND level IN ? AND tenant_id = ?", petID, weekStart, []string{"critical", "emergency"}, tenantID).Count(&warningStats.Critical)

	report := map[string]interface{}{
		"pet_uuid":          petID,
		"report_type":       "weekly",
		"start_date":        weekStart,
		"end_date":          now,
		"exercise_summary":   weeklyExercise,
		"sleep_summary": map[string]interface{}{
			"total_records":     len(sleepRecords),
			"avg_duration":      avgSleepDuration,
			"avg_deep_sleep":    avgDeepSleep,
			"avg_rem_sleep":     avgREMSleep,
			"avg_quality_score": avgQualityScore,
		},
		"warning_stats": warningStats,
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// GetMonthlyReport 获取月报
func (h *HealthTrackingCtrl) GetMonthlyReport(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := getTenantID(c)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Monthly exercise summary
	var monthlyExercise models.ExerciseSummary
	h.DB.Model(&models.ExerciseSummary{}).
		Where("pet_uuid = ? AND summary_date >= ? AND summary_date <= ? AND tenant_id = ?", petID, monthStart, now, tenantID).
		Select("SUM(total_duration) as total_duration, SUM(total_steps) as total_steps, SUM(total_distance) as total_distance, SUM(total_calories) as total_calories, COUNT(*) as exercise_count").
		Scan(&monthlyExercise)

	// Monthly sleep analysis
	var sleepRecords []models.SleepRecord
	h.DB.Where("pet_uuid = ? AND sleep_date >= ? AND sleep_date <= ? AND is_nap = ? AND tenant_id = ?", petID, monthStart, now, false, tenantID).
		Find(&sleepRecords)

	var totalSleepDuration, totalDeepSleep, totalREMSleep int
	var totalQualityScore float64
	for _, r := range sleepRecords {
		totalSleepDuration += r.TotalDuration
		totalDeepSleep += r.DeepSleepTime
		totalREMSleep += r.REMTime
		totalQualityScore += float64(r.QualityScore)
	}
	avgSleepDuration := 0
	avgDeepSleep := 0
	avgREMSleep := 0
	avgQualityScore := float64(0)
	if len(sleepRecords) > 0 {
		avgSleepDuration = totalSleepDuration / len(sleepRecords)
		avgDeepSleep = totalDeepSleep / len(sleepRecords)
		avgREMSleep = totalREMSleep / len(sleepRecords)
		avgQualityScore = totalQualityScore / float64(len(sleepRecords))
	}

	// Monthly warning stats
	var warningStats struct {
		Total    int64
		Active   int64
		Critical int64
	}
	h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND created_at >= ? AND tenant_id = ?", petID, monthStart, tenantID).Count(&warningStats.Total)
	h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND created_at >= ? AND status = 'active' AND tenant_id = ?", petID, monthStart, tenantID).Count(&warningStats.Active)
	h.DB.Model(&models.HealthWarning{}).Where("pet_uuid = ? AND created_at >= ? AND level IN ? AND tenant_id = ?", petID, monthStart, []string{"critical", "emergency"}, tenantID).Count(&warningStats.Critical)

	report := map[string]interface{}{
		"pet_uuid":          petID,
		"report_type":       "monthly",
		"start_date":        monthStart,
		"end_date":          now,
		"exercise_summary":   monthlyExercise,
		"sleep_summary": map[string]interface{}{
			"total_records":     len(sleepRecords),
			"avg_duration":      avgSleepDuration,
			"avg_deep_sleep":    avgDeepSleep,
			"avg_rem_sleep":     avgREMSleep,
			"avg_quality_score": avgQualityScore,
		},
		"warning_stats": warningStats,
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}