package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DigitalTwinController 宠物数字孪生控制器 (Sprint 18)
type DigitalTwinController struct {
	DB *gorm.DB
}

// RegisterDigitalTwinRoutes 注册数字孪生路由
func (d *DigitalTwinController) RegisterDigitalTwinRoutes(r *gin.RouterGroup) {
	// 生命体征 API
	r.GET("/digital-twin/:pet_id/vitals", d.GetVitals)
	r.GET("/digital-twin/:pet_id/vitals/history", d.GetVitalsHistory)
	r.POST("/digital-twin/:pet_id/vitals/report", d.ReportVitals)

	// 行为预测 API
	r.GET("/digital-twin/:pet_id/predictions", d.GetPredictions)
	r.GET("/digital-twin/:pet_id/predictions/short-term", d.GetShortTermPredictions)
	r.GET("/digital-twin/:pet_id/predictions/intent", d.GetIntentPrediction)

	// 历史回放 API
	r.GET("/digital-twin/:pet_id/timeline", d.GetTimeline)
	r.GET("/digital-twin/:pet_id/events/:event_id", d.GetEventDetail)

	// 健康预警 API
	r.GET("/digital-twin/:pet_id/alerts", d.GetAlerts)
	r.POST("/digital-twin/:pet_id/alerts/:id/ack", d.AckAlert)
	r.POST("/digital-twin/:pet_id/alerts/:id/ignore", d.IgnoreAlert)

	// 高光时刻 API
	r.GET("/digital-twin/:pet_id/highlights", d.GetHighlights)
	r.POST("/digital-twin/:pet_id/highlights", d.CreateHighlight)
}

// GetVitals 获取实时生命体征
func (d *DigitalTwinController) GetVitals(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)
	vitalType := c.Query("vital_type")

	var records []models.VitalRecord
	var query *gorm.DB
	if tenantID != "" {
		query = d.DB.Where("pet_uuid = ? AND tenant_id = ?::uuid", petID, tenantID).Order("recorded_at DESC")
	} else {
		query = d.DB.Where("pet_uuid = ?", petID).Order("recorded_at DESC")
	}
	if vitalType != "" {
		query = query.Where("vital_type = ?", vitalType)
	}
	if err := query.Limit(20).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var trends []models.VitalTrend
	var trendsQuery *gorm.DB
	if tenantID != "" {
		trendsQuery = d.DB.Where("pet_uuid = ? AND tenant_id = ?::uuid", petID, tenantID)
	} else {
		trendsQuery = d.DB.Where("pet_uuid = ?", petID)
	}
	trendsQuery.Order("last_updated DESC").Limit(10).Find(&trends)

	var alerts []models.HealthAlert
	var alertsQuery *gorm.DB
	if tenantID != "" {
		alertsQuery = d.DB.Where("pet_uuid = ? AND tenant_id = ?::uuid AND status = ?", petID, tenantID, "active")
	} else {
		alertsQuery = d.DB.Where("pet_uuid = ? AND status = ?", petID, "active")
	}
	alertsQuery.Order("occurred_at DESC").Limit(5).Find(&alerts)

	var lastUpdatedAt time.Time
	if len(records) > 0 {
		lastUpdatedAt = records[0].RecordedAt
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespVitals{
			PetUUID:       petID,
			LatestRecords: records,
			Trend:         trends,
			Alerts:        alerts,
			LastUpdatedAt: lastUpdatedAt,
		},
	})
}

// GetVitalsHistory 获取生命体征历史
func (d *DigitalTwinController) GetVitalsHistory(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	vitalType := c.Query("vital_type")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := d.DB.Model(&models.VitalRecord{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)
	if vitalType != "" {
		query = query.Where("vital_type = ?", vitalType)
	}
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("recorded_at >= ?", t)
		}
	}
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("recorded_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var records []models.VitalRecord
	if err := query.Order("recorded_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	dateStr := time.Now().Format("2006-01-02")
	var stats models.VitalStats
	d.DB.Where("pet_uuid = ? AND tenant_id = ? AND stats_date = ?", petID, tenantID, dateStr).First(&stats)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespVitalHistory{
			Records:  records,
			Stats:   stats,
			Total:   total,
			Page:    page,
			PageSize: pageSize,
		},
	})
}

// ReportVitals 设备上报体征数据
func (d *DigitalTwinController) ReportVitals(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	var req models.ReqVitalReport
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	deviceID := c.Query("device_id")
	if deviceID == "" {
		var binding models.PetDeviceBinding
		if err := d.DB.Where("pet_uuid = ? AND is_active = ?", petID, true).First(&binding).Error; err == nil {
			deviceID = binding.DeviceID
		}
	}

	unit := req.Unit
	if unit == "" {
		unit = d.getDefaultUnit(req.VitalType)
	}

	record := models.VitalRecord{
		PetUUID:    petID,
		DeviceID:   deviceID,
		VitalType:  req.VitalType,
		Value:      req.Value,
		Unit:       unit,
		DataSource: "device",
		RecordedAt: time.Now(),
		Metadata:   req.Metadata,
		TenantID:   tenantID,
	}

	d.checkVitalAbnormal(&record)

	if err := d.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败"})
		return
	}

	if record.IsAbnormal {
		d.createVitalAlert(petID, deviceID, &record, tenantID)
	}
	d.updateVitalStats(petID, deviceID, &record, tenantID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"record_uuid":    record.RecordUUID,
			"is_abnormal":    record.IsAbnormal,
			"abnormal_level": record.AbnormalLevel,
		},
	})
}

func (d *DigitalTwinController) getDefaultUnit(vitalType string) string {
	units := map[string]string{
		models.VitalTypeHeartRate:    "bpm",
		models.VitalTypeBodyTemp:     "°C",
		models.VitalTypeActivity:     "steps/min",
		models.VitalTypeSleepQuality: "score",
		models.VitalTypeCalories:     "kcal",
		models.VitalTypeWaterIntake:  "ml",
		models.VitalTypeWeight:       "kg",
		models.VitalTypeBloodOxygen:  "%",
		models.VitalTypeRespiratory:  "breaths/min",
	}
	if u, ok := units[vitalType]; ok {
		return u
	}
	return ""
}

func (d *DigitalTwinController) checkVitalAbnormal(record *models.VitalRecord) {
	normalRanges := map[string]struct{ Min, Max float64 }{
		models.VitalTypeHeartRate:    {Min: 60, Max: 140},
		models.VitalTypeBodyTemp:     {Min: 38.0, Max: 39.2},
		models.VitalTypeBloodOxygen:  {Min: 95, Max: 100},
		models.VitalTypeRespiratory: {Min: 10, Max: 30},
	}

	if rangeData, ok := normalRanges[record.VitalType]; ok {
		record.NormalRangeMin = rangeData.Min
		record.NormalRangeMax = rangeData.Max
		if record.Value < rangeData.Min || record.Value > rangeData.Max {
			record.IsAbnormal = true
			record.AbnormalLevel = "warning"
			if record.Value < rangeData.Min*0.8 || record.Value > rangeData.Max*1.2 {
				record.AbnormalLevel = "critical"
				record.Severity = 7
			} else {
				record.Severity = 4
			}
		}
	}
}

func (d *DigitalTwinController) createVitalAlert(petUUID, deviceID string, record *models.VitalRecord, tenantID string) {
	alert := models.HealthAlert{
		PetUUID:       petUUID,
		DeviceID:      deviceID,
		AlertType:     models.AlertTypeVital,
		AlertLevel:    record.AbnormalLevel,
		Title:         d.getVitalAlertTitle(record),
		Description:   d.getVitalAlertDesc(record),
		TriggerValue:  record.Value,
		ThresholdValue: record.NormalRangeMax,
		Unit:          record.Unit,
		NormalRange:   strconv.FormatFloat(record.NormalRangeMin, 'f', 1, 64) + "-" + strconv.FormatFloat(record.NormalRangeMax, 'f', 1, 64),
		Suggestion:    d.getVitalSuggestion(record),
		Urgency:       record.Severity,
		RelatedVitals: models.JSON{"record_uuid": record.RecordUUID, "vital_type": record.VitalType},
		Status:        models.AlertStatusActive,
		OccurredAt:    record.RecordedAt,
		TenantID:      tenantID,
	}
	d.DB.Create(&alert)
}

func (d *DigitalTwinController) getVitalAlertTitle(record *models.VitalRecord) string {
	titles := map[string]string{
		models.VitalTypeHeartRate:   "心率异常",
		models.VitalTypeBodyTemp:    "体温异常",
		models.VitalTypeBloodOxygen: "血氧异常",
		models.VitalTypeRespiratory: "呼吸异常",
	}
	if t, ok := titles[record.VitalType]; ok {
		return t
	}
	return "体征数据异常"
}

func (d *DigitalTwinController) getVitalAlertDesc(record *models.VitalRecord) string {
	return "检测到宠物" + record.VitalType + "为" + strconv.FormatFloat(record.Value, 'f', 2, 64) + record.Unit +
		"，超出正常范围(" + strconv.FormatFloat(record.NormalRangeMin, 'f', 1, 64) + "-" + strconv.FormatFloat(record.NormalRangeMax, 'f', 1, 64) + ")"
}

func (d *DigitalTwinController) getVitalSuggestion(record *models.VitalRecord) string {
	if record.VitalType == models.VitalTypeHeartRate {
		if record.Value > record.NormalRangeMax {
			return "建议让宠物休息，保持安静，避免剧烈运动。如果持续异常，请咨询兽医。"
		}
		return "心率偏低，建议观察宠物状态，如有不适请及时就医。"
	}
	if record.VitalType == models.VitalTypeBodyTemp {
		if record.Value > record.NormalRangeMax {
			return "体温偏高，可能存在感染或发热。建议物理降温并观察，如持续高温请就医。"
		}
		return "体温偏低，请注意保暖，观察宠物状态。"
	}
	return "建议持续观察体征变化，如有异常持续或加重，请咨询兽医。"
}

func (d *DigitalTwinController) updateVitalStats(petUUID, deviceID string, record *models.VitalRecord, tenantID string) {
	dateStr := record.RecordedAt.Format("2006-01-02")
	statsDate, _ := time.Parse("2006-01-02", dateStr)

	var stats models.VitalStats
	err := d.DB.Where("pet_uuid = ? AND tenant_id = ? AND stats_date = ?", petUUID, tenantID, statsDate).First(&stats).Error

	if err == gorm.ErrRecordNotFound {
		stats = models.VitalStats{PetUUID: petUUID, DeviceID: deviceID, StatsDate: statsDate, TenantID: tenantID}
	}

	switch record.VitalType {
	case models.VitalTypeHeartRate:
		stats.HeartRateAvg = record.Value
		stats.HeartRateMin = record.Value
		stats.HeartRateMax = record.Value
	case models.VitalTypeBodyTemp:
		stats.BodyTempAvg = record.Value
	case models.VitalTypeActivity:
		stats.ActivityTotal += record.Value
	case models.VitalTypeCalories:
		stats.CaloriesBurned += record.Value
	case models.VitalTypeWaterIntake:
		stats.WaterIntake += record.Value
	case models.VitalTypeWeight:
		stats.Weight = record.Value
	}

	if err == gorm.ErrRecordNotFound {
		d.DB.Create(&stats)
	} else {
		d.DB.Save(&stats)
	}
}

// GetPredictions 获取行为预测
func (d *DigitalTwinController) GetPredictions(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	var predictions []models.BehaviorPrediction
	query := d.DB.Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID).
		Where("time_window_end >= ?", time.Now())
	if predType := c.Query("type"); predType != "" {
		query = query.Where("prediction_type = ?", predType)
	}
	if err := query.Order("probability DESC, created_at DESC").Limit(10).Find(&predictions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var lastUpdatedAt time.Time
	if len(predictions) > 0 {
		lastUpdatedAt = predictions[0].CreatedAt
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespPredictions{
			PetUUID:       petID,
			Predictions:   predictions,
			LastUpdatedAt: lastUpdatedAt,
		},
	})
}

// GetShortTermPredictions 获取短期动作预测
func (d *DigitalTwinController) GetShortTermPredictions(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	var predictions []models.BehaviorPrediction
	d.DB.Where("pet_uuid = ? AND tenant_id = ? AND prediction_type = ? AND time_window_end >= ?",
		petID, tenantID, "short_term", time.Now()).Order("probability DESC").Limit(5).Find(&predictions)

	var shortPredictions []models.ShortTermPrediction
	for _, p := range predictions {
		shortPredictions = append(shortPredictions, models.ShortTermPrediction{
			Behavior:    p.PredictedBehavior,
			Probability: p.Probability,
			TimeIn:      int(p.TimeWindowStart.Sub(time.Now()).Seconds()),
			Duration:    int(p.TimeWindowEnd.Sub(p.TimeWindowStart).Seconds()),
			Trigger:     p.Trigger,
		})
	}

	modelVersion := "v1.0"
	if len(predictions) > 0 {
		modelVersion = predictions[0].ModelVersion
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespShortTermPrediction{
			PetUUID:      petID,
			Predictions:  shortPredictions,
			ModelVersion: modelVersion,
		},
	})
}

// GetIntentPrediction 获取意图识别
func (d *DigitalTwinController) GetIntentPrediction(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	var predictions []models.BehaviorPrediction
	d.DB.Where("pet_uuid = ? AND tenant_id = ? AND prediction_type = ? AND time_window_end >= ?",
		petID, tenantID, "intent", time.Now()).Order("probability DESC").Limit(5).Find(&predictions)

	var currentIntent models.IntentPrediction
	var alternatives []models.IntentPrediction

	if len(predictions) > 0 {
		currentIntent = models.IntentPrediction{
			Intent:      predictions[0].PredictedBehavior,
			Probability: predictions[0].Probability,
			Confidence:  predictions[0].Confidence,
			Suggestion:  predictions[0].Recommendation,
		}
		for _, p := range predictions[1:] {
			alternatives = append(alternatives, models.IntentPrediction{
				Intent:      p.PredictedBehavior,
				Probability: p.Probability,
				Confidence:  p.Confidence,
				Suggestion:  p.Recommendation,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespIntentPrediction{
			PetUUID:       petID,
			CurrentIntent: currentIntent,
			Alternatives:  alternatives,
		},
	})
}

// GetTimeline 获取时间轴事件
func (d *DigitalTwinController) GetTimeline(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	eventTypes := c.Query("event_types")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	highlightOnly := c.Query("highlight_only") == "true"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var events []models.TimelineEvent
	var highlights []models.HighlightMoment
	var total int64

	if highlightOnly {
		query := d.DB.Model(&models.HighlightMoment{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)
		if startTimeStr != "" {
			if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
				query = query.Where("occurred_at >= ?", t)
			}
		}
		if endTimeStr != "" {
			if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
				query = query.Where("occurred_at <= ?", t)
			}
		}
		query.Count(&total)
		query.Order("occurred_at DESC").Offset(offset).Limit(pageSize).Find(&highlights)
		for _, h := range highlights {
			events = append(events, models.TimelineEvent{
				ID:           h.ID,
				EventUUID:    h.MomentUUID,
				EventType:    "highlight",
				EventSubType: h.HighlightType,
				Title:        h.Title,
				Description:  h.Description,
				OccurredAt:   h.OccurredAt,
				ThumbnailURL: h.ThumbnailURL,
				VideoURL:     h.VideoURL,
				Tags:         h.Tags,
			})
		}
	} else {
		var allEvents []models.TimelineEvent

		if strings.Contains(eventTypes, "behavior") || eventTypes == "" {
			var behaviors []models.BehaviorEvent
			query := d.DB.Model(&models.BehaviorEvent{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)
			if startTimeStr != "" {
				if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
					query = query.Where("start_time >= ?", t)
				}
			}
			if endTimeStr != "" {
				if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
					query = query.Where("start_time <= ?", t)
				}
			}
			var behCount int64
			query.Count(&behCount)
			total += behCount
			query.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&behaviors)
			for _, b := range behaviors {
				allEvents = append(allEvents, models.TimelineEvent{
					ID:           b.ID,
					EventUUID:    b.EventUUID,
					EventType:    "behavior",
					EventSubType: b.BehaviorType,
					Title:        b.BehaviorName,
					OccurredAt:   b.StartTime,
					Duration:     b.Duration,
					Intensity:    b.Intensity,
					IsAnomaly:    b.IsAnomaly,
					Tags:         b.Tags,
				})
			}
		}

		if strings.Contains(eventTypes, "highlight") || eventTypes == "" {
			var hl []models.HighlightMoment
			query := d.DB.Model(&models.HighlightMoment{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)
			if startTimeStr != "" {
				if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
					query = query.Where("occurred_at >= ?", t)
				}
			}
			if endTimeStr != "" {
				if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
					query = query.Where("occurred_at <= ?", t)
				}
			}
			var hlCount int64
			query.Count(&hlCount)
			total += hlCount
			query.Order("occurred_at DESC").Offset(offset).Limit(pageSize).Find(&hl)
			for _, h := range hl {
				allEvents = append(allEvents, models.TimelineEvent{
					ID:           h.ID,
					EventUUID:    h.MomentUUID,
					EventType:    "highlight",
					EventSubType: h.HighlightType,
					Title:        h.Title,
					Description:  h.Description,
					OccurredAt:   h.OccurredAt,
					ThumbnailURL: h.ThumbnailURL,
					VideoURL:     h.VideoURL,
					Tags:         h.Tags,
				})
			}
		}
		events = allEvents
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespTimeline{
			PetUUID:    petID,
			Events:     events,
			Highlights: highlights,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
		},
	})
}

// GetEventDetail 获取事件详情
func (d *DigitalTwinController) GetEventDetail(c *gin.Context) {
	petID := c.Param("pet_id")
	eventID := c.Param("event_id")
	tenantID := middleware.GetTenantIDCtx(c)
	eventType := c.Query("type")
	if eventType == "" {
		eventType = "behavior"
	}

	var detail models.EventDetail
	detail.EventUUID = eventID

	switch eventType {
	case "behavior":
		var event models.BehaviorEvent
		if err := d.DB.Where("event_uuid = ? AND pet_uuid = ? AND tenant_id = ?", eventID, petID, tenantID).First(&event).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "事件不存在"})
			return
		}
		detail.EventType = "behavior"
		detail.Data = event

		var related []models.BehaviorEvent
		d.DB.Where("pet_uuid = ? AND tenant_id = ? AND behavior_type = ? AND id != ?",
			petID, tenantID, event.BehaviorType, event.ID).
			Order("start_time DESC").Limit(5).Find(&related)

		var relatedEvents []models.TimelineEvent
		for _, r := range related {
			relatedEvents = append(relatedEvents, models.TimelineEvent{
				ID:           r.ID,
				EventUUID:    r.EventUUID,
				EventType:    "behavior",
				EventSubType: r.BehaviorType,
				Title:        r.BehaviorName,
				OccurredAt:   r.StartTime,
				Duration:     r.Duration,
				IsAnomaly:    r.IsAnomaly,
			})
		}
		detail.RelatedEvents = relatedEvents

	case "highlight":
		var moment models.HighlightMoment
		if err := d.DB.Where("moment_uuid = ? AND pet_uuid = ? AND tenant_id = ?", eventID, petID, tenantID).First(&moment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "事件不存在"})
			return
		}
		detail.EventType = "highlight"
		detail.Data = moment

	case "alert":
		var alert models.HealthAlert
		if err := d.DB.Where("alert_uuid = ? AND pet_uuid = ? AND tenant_id = ?", eventID, petID, tenantID).First(&alert).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
			return
		}
		detail.EventType = "alert"
		detail.Data = alert

	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "未知事件类型"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": detail})
}

// GetAlerts 获取健康预警列表
func (d *DigitalTwinController) GetAlerts(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	alertType := c.Query("alert_type")
	alertLevel := c.Query("alert_level")
	status := c.Query("status")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := d.DB.Model(&models.HealthAlert{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)
	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if alertLevel != "" {
		query = query.Where("alert_level = ?", alertLevel)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", models.AlertStatusActive)
	}
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("occurred_at >= ?", t)
		}
	}
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("occurred_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var unackCount int64
	d.DB.Model(&models.HealthAlert{}).
		Where("pet_uuid = ? AND tenant_id = ? AND status = ?", petID, tenantID, models.AlertStatusActive).
		Count(&unackCount)

	var alerts []models.HealthAlert
	if err := query.Order("occurred_at DESC").Offset(offset).Limit(pageSize).Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": models.RespAlertList{
			Alerts:     alerts,
			Total:      total,
			UnackCount: unackCount,
			Page:       page,
			PageSize:   pageSize,
		},
	})
}

// AckAlert 确认预警
func (d *DigitalTwinController) AckAlert(c *gin.Context) {
	petID := c.Param("pet_id")
	alertID := c.Param("id")
	tenantID := middleware.GetTenantIDCtx(c)
	userID := middleware.GetUserID(c)

	var alert models.HealthAlert
	if err := d.DB.Where("alert_uuid = ? AND pet_uuid = ? AND tenant_id = ?", alertID, petID, tenantID).First(&alert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}

	if alert.Status != models.AlertStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预警已处理"})
		return
	}

	now := time.Now()
	alert.Status = models.AlertStatusAcked
	alert.AcknowledgedAt = &now
	alert.AcknowledgedBy = &userID

	if err := d.DB.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": alert})
}

// IgnoreAlert 忽略预警
func (d *DigitalTwinController) IgnoreAlert(c *gin.Context) {
	petID := c.Param("pet_id")
	alertID := c.Param("id")
	tenantID := middleware.GetTenantIDCtx(c)
	userID := middleware.GetUserID(c)

	var req models.ReqIgnoreAlert
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供忽略原因"})
		return
	}

	var alert models.HealthAlert
	if err := d.DB.Where("alert_uuid = ? AND pet_uuid = ? AND tenant_id = ?", alertID, petID, tenantID).First(&alert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}

	if alert.Status != models.AlertStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预警已处理"})
		return
	}

	now := time.Now()
	alert.Status = models.AlertStatusIgnored
	alert.IgnoredAt = &now
	alert.IgnoredBy = &userID
	alert.IgnoreReason = req.Reason

	if err := d.DB.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": alert})
}

// GetHighlights 获取高光时刻列表
func (d *DigitalTwinController) GetHighlights(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)

	highlightType := c.Query("highlight_type")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := d.DB.Model(&models.HighlightMoment{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)
	if highlightType != "" {
		query = query.Where("highlight_type = ?", highlightType)
	}
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("occurred_at >= ?", t)
		}
	}
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("occurred_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var highlights []models.HighlightMoment
	if err := query.Order("occurred_at DESC").Offset(offset).Limit(pageSize).Find(&highlights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"list": highlights, "total": total, "page": page, "page_size": pageSize},
	})
}

// CreateHighlight 创建高光时刻
func (d *DigitalTwinController) CreateHighlight(c *gin.Context) {
	petID := c.Param("pet_id")
	tenantID := middleware.GetTenantIDCtx(c)
	userID := middleware.GetUserID(c)

	var req models.ReqCreateHighlight
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	occurredAt := time.Now()
	if req.OccurredAt != "" {
		if t, err := time.Parse(time.RFC3339, req.OccurredAt); err == nil {
			occurredAt = t
		}
	}

	deviceID := ""
	var binding models.PetDeviceBinding
	if err := d.DB.Where("pet_uuid = ? AND is_active = ?", petID, true).First(&binding).Error; err == nil {
		deviceID = binding.DeviceID
	}

	moment := models.HighlightMoment{
		PetUUID:       petID,
		DeviceID:      deviceID,
		HighlightType: req.HighlightType,
		Title:         req.Title,
		Description:   req.Description,
		OccurredAt:    occurredAt,
		MediaURLs:     req.MediaURLs,
		Shareable:     req.Shareable,
		IsPublic:      req.IsPublic,
		Tags:          req.Tags,
		CapturedBy:    "manual",
		OwnerID:       userID,
		TenantID:      tenantID,
	}

	if err := d.DB.Create(&moment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": moment})
}
