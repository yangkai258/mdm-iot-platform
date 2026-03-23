package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DigitalTwinController 数字孪生控制器
type DigitalTwinController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册数字孪生路由
func (d *DigitalTwinController) RegisterRoutes(r *gin.RouterGroup) {
	// 生命体征路由
	r.GET("/vitals/dashboard", d.GetVitalsDashboard)
	r.GET("/vitals/realtime", d.GetRealtimeVitals)
	r.GET("/vitals/history", d.GetVitalsHistory)
	r.GET("/vitals/alerts", d.GetHealthAlerts)
	r.POST("/vitals/alerts/:id/confirm", d.ConfirmAlert)
	r.POST("/vitals/alerts/:id/ignore", d.IgnoreAlert)

	// 行为路由
	r.GET("/behavior/prediction", d.GetBehaviorPrediction)
	r.GET("/behavior/history", d.GetBehaviorHistory)

	// 历史回放路由
	r.GET("/replay/:pet_id", d.GetReplay)

	// 精彩瞬间路由
	r.GET("/highlights", d.GetHighlights)
}

// GetVitalsDashboard 获取生命体征仪表盘数据
// @Summary 生命体征仪表盘
// @Description 获取心率/呼吸率/体温等实时数据汇总
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint false "宠物ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/vitals/dashboard [GET]
func (d *DigitalTwinController) GetVitalsDashboard(c *gin.Context) {
	petIDStr := c.Query("pet_id")

	var records []models.VitalRecord
	query := d.DB.Model(&models.VitalRecord{}).Order("recorded_at DESC")

	if petIDStr != "" {
		petID, _ := strconv.ParseUint(petIDStr, 10, 32)
		query = query.Where("pet_id = ?", petID)
	}

	// 获取最近24小时的数据
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	query = query.Where("recorded_at >= ?", oneDayAgo)

	if err := query.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 按类型分组汇总
	dashboard := make(map[string]interface{})
	typeSummary := make(map[string]map[string]interface{})

	vitalTypes := []string{"heart_rate", "breathing", "temperature"}
	for _, vt := range vitalTypes {
		typeSummary[vt] = map[string]interface{}{
			"latest":  nil,
			"avg":     0.0,
			"min":     nil,
			"max":     nil,
			"count":   0,
			"unit":    "",
			"records": []models.VitalRecord{},
		}
	}

	var heartRateSum, breathingSum, tempSum float64
	var heartRateCount, breathingCount, tempCount int
	var latestHeartRate, latestBreathing, latestTemp *float64

	for _, r := range records {
		switch r.Type {
		case "heart_rate":
			if latestHeartRate == nil {
				lp := r.Value
				latestHeartRate = &lp
			}
			heartRateSum += r.Value
			heartRateCount++
			if typeSummary["heart_rate"]["max"] == nil || r.Value > typeSummary["heart_rate"]["max"].(float64) {
				typeSummary["heart_rate"]["max"] = r.Value
			}
			if typeSummary["heart_rate"]["min"] == nil || r.Value < typeSummary["heart_rate"]["min"].(float64) {
				typeSummary["heart_rate"]["min"] = r.Value
			}
			typeSummary["heart_rate"]["unit"] = r.Unit
		case "breathing":
			if latestBreathing == nil {
				lp := r.Value
				latestBreathing = &lp
			}
			breathingSum += r.Value
			breathingCount++
			if typeSummary["breathing"]["max"] == nil || r.Value > typeSummary["breathing"]["max"].(float64) {
				typeSummary["breathing"]["max"] = r.Value
			}
			if typeSummary["breathing"]["min"] == nil || r.Value < typeSummary["breathing"]["min"].(float64) {
				typeSummary["breathing"]["min"] = r.Value
			}
			typeSummary["breathing"]["unit"] = r.Unit
		case "temperature":
			if latestTemp == nil {
				lp := r.Value
				latestTemp = &lp
			}
			tempSum += r.Value
			tempCount++
			if typeSummary["temperature"]["max"] == nil || r.Value > typeSummary["temperature"]["max"].(float64) {
				typeSummary["temperature"]["max"] = r.Value
			}
			if typeSummary["temperature"]["min"] == nil || r.Value < typeSummary["temperature"]["min"].(float64) {
				typeSummary["temperature"]["min"] = r.Value
			}
			typeSummary["temperature"]["unit"] = r.Unit
		}
	}

	// 直接设置各类型汇总数据
	heartRateSummary := map[string]interface{}{
		"latest": latestHeartRate,
		"avg":    0.0,
		"min":    nil,
		"max":    nil,
		"count":  heartRateCount,
		"unit":   "",
	}
	if heartRateCount > 0 {
		heartRateSummary["avg"] = heartRateSum / float64(heartRateCount)
	}

	breathingSummary := map[string]interface{}{
		"latest": latestBreathing,
		"avg":    0.0,
		"min":    nil,
		"max":    nil,
		"count":  breathingCount,
		"unit":   "",
	}
	if breathingCount > 0 {
		breathingSummary["avg"] = breathingSum / float64(breathingCount)
	}

	tempSummary := map[string]interface{}{
		"latest": latestTemp,
		"avg":    0.0,
		"min":    nil,
		"max":    nil,
		"count":  tempCount,
		"unit":   "",
	}
	if tempCount > 0 {
		tempSummary["avg"] = tempSum / float64(tempCount)
	}

	dashboard["heart_rate"] = heartRateSummary
	dashboard["breathing"] = breathingSummary
	dashboard["temperature"] = tempSummary

	dashboard["summary"] = map[string]interface{}{
		"total_records": len(records),
		"time_range":    map[string]string{"start": oneDayAgo.Format(time.RFC3339), "end": time.Now().Format(time.RFC3339)},
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": dashboard})
}

// GetRealtimeVitals 获取实时生命体征数据
// @Summary 实时生命体征
// @Description 获取当前心率、呼吸、体温等
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint false "宠物ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/vitals/realtime [GET]
func (d *DigitalTwinController) GetRealtimeVitals(c *gin.Context) {
	petIDStr := c.Query("pet_id")

	vitals := make(map[string]interface{})
	vitalTypes := []string{"heart_rate", "breathing", "temperature"}

	for _, vt := range vitalTypes {
		query := d.DB.Model(&models.VitalRecord{}).
			Where("type = ?", vt).
			Order("recorded_at DESC")

		if petIDStr != "" {
			petID, _ := strconv.ParseUint(petIDStr, 10, 32)
			query = query.Where("pet_id = ?", petID)
		}

		var record models.VitalRecord
		if err := query.First(&record).Error; err == nil {
			vitals[vt] = map[string]interface{}{
				"value":       record.Value,
				"unit":        record.Unit,
				"recorded_at": record.RecordedAt,
				"source":      record.Source,
			}
		} else {
			vitals[vt] = nil
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": vitals})
}

// GetVitalsHistory 获取历史生命体征查询
// @Summary 历史生命体征
// @Description 查询历史生命体征记录
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint true "宠物ID"
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Param type query string false "类型 heart_rate/breathing/temperature"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/vitals/history [GET]
func (d *DigitalTwinController) GetVitalsHistory(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	vitalType := c.Query("type")

	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	// 解析日期，默认最近7天
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(24*time.Hour - time.Second)
		}
	}

	query := d.DB.Model(&models.VitalRecord{}).
		Where("pet_id = ?", petID).
		Where("recorded_at >= ? AND recorded_at <= ?", startDate, endDate).
		Order("recorded_at DESC")

	if vitalType != "" {
		query = query.Where("type = ?", vitalType)
	}

	var records []models.VitalRecord
	if err := query.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 按日期分组
	dailyData := make(map[string][]models.VitalRecord)
	for _, r := range records {
		dateKey := r.RecordedAt.Format("2006-01-02")
		dailyData[dateKey] = append(dailyData[dateKey], r)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"records":    records,
			"daily_data": dailyData,
			"total":      len(records),
			"time_range": map[string]string{
				"start": startDate.Format("2006-01-02"),
				"end":   endDate.Format("2006-01-02"),
			},
		},
	})
}

// GetHealthAlerts 获取健康预警列表
// @Summary 健康预警列表
// @Description 获取异常预警记录
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint false "宠物ID"
// @Param status query string false "状态 pending/confirmed/ignored/resolved"
// @Param level query string false "级别 warning/critical"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/vitals/alerts [GET]
func (d *DigitalTwinController) GetHealthAlerts(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	status := c.Query("status")
	level := c.Query("level")

	query := d.DB.Model(&models.HealthAlert{}).Order("detected_at DESC")

	if petIDStr != "" {
		petID, _ := strconv.ParseUint(petIDStr, 10, 32)
		query = query.Where("pet_id = ?", petID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if level != "" {
		query = query.Where("level = ?", level)
	}

	// 默认显示最近30天的未处理预警
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	query = query.Where("detected_at >= ?", thirtyDaysAgo)

	var alerts []models.HealthAlert
	if err := query.Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 统计
	stats := map[string]int{
		"total":     len(alerts),
		"pending":   0,
		"critical": 0,
		"warning":  0,
	}
	for _, a := range alerts {
		if a.Status == "pending" {
			stats["pending"]++
		}
		if a.Level == "critical" {
			stats["critical"]++
		}
		if a.Level == "warning" {
			stats["warning"]++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"alerts": alerts,
			"stats":  stats,
		},
	})
}

// ConfirmAlert 确认预警
// @Summary 确认预警
// @Description 确认指定的健康预警
// @Tags digital-twin
// @Produce json
// @Param id path uint true "预警ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/vitals/alerts/{id}/confirm [POST]
func (d *DigitalTwinController) ConfirmAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	now := time.Now()
	result := d.DB.Model(&models.HealthAlert{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "confirmed",
			"confirmed_at": &now,
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "预警已确认"})
}

// IgnoreAlert 忽略预警
// @Summary 忽略预警
// @Description 忽略指定的健康预警
// @Tags digital-twin
// @Produce json
// @Param id path uint true "预警ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/vitals/alerts/{id}/ignore [POST]
func (d *DigitalTwinController) IgnoreAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	result := d.DB.Model(&models.HealthAlert{}).
		Where("id = ?", id).
		Update("status", "ignored")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "预警已忽略"})
}

// GetBehaviorPrediction 获取行为预测
// @Summary 行为预测
// @Description 预测下一步行为
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint true "宠物ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/behavior/prediction [GET]
func (d *DigitalTwinController) GetBehaviorPrediction(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	// 基于最近的行为模式预测
	var recentBehaviors []models.BehaviorEvent
	if err := d.DB.Model(&models.BehaviorEvent{}).
		Where("pet_id = ?", petID).
		Order("start_time DESC").
		Limit(10).
		Find(&recentBehaviors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 简单的基于频率的预测
	behaviorFreq := make(map[string]int)
	for _, b := range recentBehaviors {
		behaviorFreq[b.Type]++
	}

	var mostLikely string
	maxFreq := 0
	for bt, freq := range behaviorFreq {
		if freq > maxFreq {
			maxFreq = freq
			mostLikely = bt
		}
	}

	// 预测下一次行为的时间（基于历史平均间隔）
	prediction := map[string]interface{}{
		"predicted_behavior": mostLikely,
		"confidence":         0.0,
		"reason":             "基于最近行为频率分析",
	}

	if mostLikely != "" && len(recentBehaviors) > 0 {
		// 简单计算置信度
		prediction["confidence"] = float64(maxFreq) / float64(len(recentBehaviors))
	}

	// 计算预测下一次行为的时间
	var avgInterval int64
	if len(recentBehaviors) >= 2 {
		var totalInterval int64
		for i := 1; i < len(recentBehaviors); i++ {
			totalInterval += int64(recentBehaviors[i-1].StartTime.Sub(recentBehaviors[i].StartTime).Seconds())
		}
		avgInterval = totalInterval / int64(len(recentBehaviors)-1)
	}
	prediction["estimated_next_time"] = time.Now().Add(time.Duration(avgInterval) * time.Second).Format(time.RFC3339)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": prediction})
}

// GetBehaviorHistory 获取行为历史
// @Summary 行为历史
// @Description 查询宠物行为历史记录
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint true "宠物ID"
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/behavior/history [GET]
func (d *DigitalTwinController) GetBehaviorHistory(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if petIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pet_id不能为空"})
		return
	}

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	// 解析日期，默认最近7天
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(24*time.Hour - time.Second)
		}
	}

	var events []models.BehaviorEvent
	if err := d.DB.Model(&models.BehaviorEvent{}).
		Where("pet_id = ?", petID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Order("start_time DESC").
		Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 统计各类型行为时长
	typeStats := make(map[string]map[string]interface{})
	for _, e := range events {
		if _, ok := typeStats[e.Type]; !ok {
			typeStats[e.Type] = map[string]interface{}{
				"count":   0,
				"total_duration": 0,
			}
		}
		typeStats[e.Type]["count"] = typeStats[e.Type]["count"].(int) + 1
		typeStats[e.Type]["total_duration"] = typeStats[e.Type]["total_duration"].(int) + e.Duration
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"events":    events,
			"stats":     typeStats,
			"total":     len(events),
			"time_range": map[string]string{
				"start": startDate.Format("2006-01-02"),
				"end":   endDate.Format("2006-01-02"),
			},
		},
	})
}

// GetReplay 获取历史回放数据
// @Summary 历史回放
// @Description 获取时间线数据用于回放
// @Tags digital-twin
// @Produce json
// @Param pet_id path uint true "宠物ID"
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/replay/{pet_id} [GET]
func (d *DigitalTwinController) GetReplay(c *gin.Context) {
	petIDStr := c.Param("pet_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	petID, _ := strconv.ParseUint(petIDStr, 10, 32)

	// 解析日期，默认当天
	endDate := time.Now()
	startDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.Local)
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(24*time.Hour - time.Second)
		}
	}

	// 获取生命体征记录
	var vitals []models.VitalRecord
	d.DB.Model(&models.VitalRecord{}).
		Where("pet_id = ?", petID).
		Where("recorded_at >= ? AND recorded_at <= ?", startDate, endDate).
		Order("recorded_at ASC").
		Find(&vitals)

	// 获取行为事件
	var behaviors []models.BehaviorEvent
	d.DB.Model(&models.BehaviorEvent{}).
		Where("pet_id = ?", petID).
		Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Order("start_time ASC").
		Find(&behaviors)

	// 获取预警
	var alerts []models.HealthAlert
	d.DB.Model(&models.HealthAlert{}).
		Where("pet_id = ?", petID).
		Where("detected_at >= ? AND detected_at <= ?", startDate, endDate).
		Order("detected_at ASC").
		Find(&alerts)

	// 合并为时间线
	timeline := make([]map[string]interface{}, 0)

	for _, v := range vitals {
		timeline = append(timeline, map[string]interface{}{
			"type":      "vital",
			"sub_type":  v.Type,
			"value":     v.Value,
			"unit":      v.Unit,
			"timestamp": v.RecordedAt,
			"source":    v.Source,
		})
	}

	for _, b := range behaviors {
		item := map[string]interface{}{
			"type":      "behavior",
			"sub_type":  b.Type,
			"start":     b.StartTime,
			"duration":  b.Duration,
			"timestamp": b.StartTime,
		}
		if b.EndTime != nil {
			item["end"] = b.EndTime
		}
		timeline = append(timeline, item)
	}

	for _, a := range alerts {
		timeline = append(timeline, map[string]interface{}{
			"type":      "alert",
			"sub_type":  a.Type,
			"level":     a.Level,
			"message":   a.Message,
			"status":    a.Status,
			"timestamp": a.DetectedAt,
		})
	}

	// 按时间排序
	for i := 0; i < len(timeline)-1; i++ {
		for j := i + 1; j < len(timeline); j++ {
			ti := timeline[i]["timestamp"].(time.Time)
			tj := timeline[j]["timestamp"].(time.Time)
			if ti.After(tj) {
				timeline[i], timeline[j] = timeline[j], timeline[i]
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"pet_id":     petID,
			"timeline":   timeline,
			"vitals_count": len(vitals),
			"behaviors_count": len(behaviors),
			"alerts_count":  len(alerts),
			"time_range": map[string]string{
				"start": startDate.Format("2006-01-02 15:04:05"),
				"end":   endDate.Format("2006-01-02 15:04:05"),
			},
		},
	})
}

// GetHighlights 获取精彩瞬间
// @Summary 精彩瞬间
// @Description 获取宠物的精彩瞬间记录
// @Tags digital-twin
// @Produce json
// @Param pet_id query uint false "宠物ID"
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/digital-twin/highlights [GET]
func (d *DigitalTwinController) GetHighlights(c *gin.Context) {
	petIDStr := c.Query("pet_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	query := d.DB.Model(&models.HighlightMoment{}).Order("captured_at DESC")

	if petIDStr != "" {
		petID, _ := strconv.ParseUint(petIDStr, 10, 32)
		query = query.Where("pet_id = ?", petID)
	}

	// 解析日期，默认最近30天
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(24*time.Hour - time.Second)
		}
	}
	query = query.Where("captured_at >= ? AND captured_at <= ?", startDate, endDate)

	var highlights []models.HighlightMoment
	if err := query.Find(&highlights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 按类型分组
	byType := make(map[string][]models.HighlightMoment)
	for _, h := range highlights {
		byType[h.Type] = append(byType[h.Type], h)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"highlights": highlights,
			"by_type":    byType,
			"total":      len(highlights),
			"time_range": map[string]string{
				"start": startDate.Format("2006-01-02"),
				"end":   endDate.Format("2006-01-02"),
			},
		},
	})
}
