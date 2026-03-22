package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceMonitorController 设备监控控制器
type DeviceMonitorController struct {
	DB *gorm.DB
}

// RegisterDeviceMonitorRoutes 注册设备监控路由
func (ctrl *DeviceMonitorController) RegisterDeviceMonitorRoutes(api *gin.RouterGroup) {
	api.GET("/monitor/devices/:device_id/metrics", ctrl.GetMetrics)
	api.GET("/monitor/devices/:device_id/realtime", ctrl.GetRealtime)
	api.GET("/monitor/devices/:device_id/history", ctrl.GetHistory)
	api.GET("/monitor/devices/:device_id/alerts", ctrl.GetAlerts)
}

// GetMetrics 获取设备指标数据
// GET /api/v1/monitor/devices/:device_id/metrics
func (ctrl *DeviceMonitorController) GetMetrics(c *gin.Context) {
	deviceID := c.Param("device_id")

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	metricType := c.Query("metric_type")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	sortBy := c.DefaultQuery("sort_by", "timestamp")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.DeviceMetric{}).Where("device_id = ?", deviceID)

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}
	if startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			query = query.Where("timestamp >= ?", t)
		}
	}
	if endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			query = query.Where("timestamp <= ?", t)
		}
	}

	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	var total int64
	query.Count(&total)

	var metrics []models.DeviceMetric
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query device metrics: " + err.Error(),
		})
		return
	}

	responses := make([]*models.DeviceMetricResponse, len(metrics))
	for i := range metrics {
		responses[i] = metrics[i].ToResponse()
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

// GetRealtime 获取实时监控数据
// GET /api/v1/monitor/devices/:device_id/realtime
func (ctrl *DeviceMonitorController) GetRealtime(c *gin.Context) {
	deviceID := c.Param("device_id")

	// 获取最新一条各类型指标
	var metrics []models.DeviceMetric
	if err := ctrl.DB.Model(&models.DeviceMetric{}).
		Where("device_id = ?", deviceID).
		Order("timestamp desc").
		Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query realtime metrics: " + err.Error(),
		})
		return
	}

	// 按 metric_type 分组，每类取最新一条
	latestMap := make(map[string]*models.DeviceMetric)
	for i := range metrics {
		mt := metrics[i].MetricType
		if _, exists := latestMap[mt]; !exists {
			latestMap[mt] = &metrics[i]
		}
	}

	var results []*models.DeviceMetricResponse
	for _, m := range latestMap {
		results = append(results, m.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":  deviceID,
			"metrics":    results,
			"fetched_at": time.Now(),
		},
	})
}

// GetHistory 获取历史数据（按时间范围）
// GET /api/v1/monitor/devices/:device_id/history
func (ctrl *DeviceMonitorController) GetHistory(c *gin.Context) {
	deviceID := c.Param("device_id")

	metricType := c.Query("metric_type")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	interval := c.DefaultQuery("interval", "hour") // hour/day/minute

	if startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "start_time and end_time are required"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid start_time format, use RFC3339"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid end_time format, use RFC3339"})
		return
	}

	query := ctrl.DB.Model(&models.DeviceMetric{}).
		Where("device_id = ? AND timestamp >= ? AND timestamp <= ?", deviceID, startTime, endTime)

	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}

	var metrics []models.DeviceMetric
	if err := query.Order("timestamp asc").Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query history metrics: " + err.Error(),
		})
		return
	}

	// 按时间窗口聚合
	type HistoryPoint struct {
		Timestamp time.Time `json:"timestamp"`
		AvgValue  float64   `json:"avg_value"`
		MinValue  float64   `json:"min_value"`
		MaxValue  float64   `json:"max_value"`
		Count     int       `json:"count"`
	}

	// 简化处理：按分钟聚合
	windowMinutes := 1
	switch interval {
	case "day":
		windowMinutes = 1440
	case "hour":
		windowMinutes = 60
	case "minute":
		windowMinutes = 1
	}

	type MetricSeries struct {
		MetricType string         `json:"metric_type"`
		Points     []HistoryPoint `json:"points"`
	}

	seriesMap := make(map[string][]HistoryPoint)
	for _, m := range metrics {
		// 按窗口对齐时间
		t := m.Timestamp.Truncate(time.Duration(windowMinutes) * time.Minute)
		st := m.MetricType
		points := seriesMap[st]
		found := false
		for i := range points {
			if points[i].Timestamp.Equal(t) {
				points[i].Count++
				points[i].AvgValue = (points[i].AvgValue*float64(points[i].Count-1) + m.MetricValue) / float64(points[i].Count)
				if m.MetricValue < points[i].MinValue {
					points[i].MinValue = m.MetricValue
				}
				if m.MetricValue > points[i].MaxValue {
					points[i].MaxValue = m.MetricValue
				}
				found = true
				break
			}
		}
		if !found {
			points = append(points, HistoryPoint{
				Timestamp: t,
				AvgValue:  m.MetricValue,
				MinValue:  m.MetricValue,
				MaxValue:  m.MetricValue,
				Count:     1,
			})
		}
		seriesMap[st] = points
	}

	var series []MetricSeries
	for mt, pts := range seriesMap {
		series = append(series, MetricSeries{
			MetricType: mt,
			Points:     pts,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":  deviceID,
			"start_time": startTime,
			"end_time":   endTime,
			"interval":   interval,
			"series":     series,
		},
	})
}

// GetAlerts 获取设备告警历史
// GET /api/v1/monitor/devices/:device_id/alerts
func (ctrl *DeviceMonitorController) GetAlerts(c *gin.Context) {
	deviceID := c.Param("device_id")

	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	alertType := c.Query("alert_type")
	status := c.Query("status")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	query := ctrl.DB.Model(&models.DeviceAlert{}).Where("device_id = ?", deviceID)

	if alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var alerts []models.DeviceAlert
	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query device alerts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      alerts,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ReportMetric 上报设备指标（由设备或MQTT调用）
// POST /api/v1/monitor/metrics
func (ctrl *DeviceMonitorController) ReportMetric(c *gin.Context) {
	var req struct {
		DeviceID    string                 `json:"device_id" binding:"required"`
		MetricType  string                 `json:"metric_type" binding:"required"`
		MetricName  string                 `json:"metric_name"`
		MetricValue float64                `json:"metric_value" binding:"required"`
		Unit        string                 `json:"unit"`
		Tags        map[string]interface{} `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	if req.MetricName == "" {
		req.MetricName = req.MetricType
	}

	tagsJSON, _ := json.Marshal(req.Tags)
	metric := models.DeviceMetric{
		DeviceID:    req.DeviceID,
		MetricType:  req.MetricType,
		MetricName:  req.MetricName,
		MetricValue: req.MetricValue,
		Unit:        req.Unit,
		Timestamp:   time.Now(),
		Tags:        string(tagsJSON),
	}

	if err := ctrl.DB.Create(&metric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create metric: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Metric recorded successfully",
		"data":    metric.ToResponse(),
	})
}
