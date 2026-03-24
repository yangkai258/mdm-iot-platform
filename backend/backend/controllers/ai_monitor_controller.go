package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AiMonitorController AI模型监控控制器
type AiMonitorController struct {
	DB *gorm.DB
}

// ============ 监控指标 ============

// Metrics 获取监控指标
func (c *AiMonitorController) Metrics(ctx *gin.Context) {
	modelVersionID := ctx.Query("model_version_id")
	metricName := ctx.Query("metric_name")
	period := ctx.DefaultQuery("period", "1h") // 1h, 6h, 24h, 7d, 30d

	var startTime time.Time
	now := time.Now()
	switch period {
	case "1h":
		startTime = now.Add(-1 * time.Hour)
	case "6h":
		startTime = now.Add(-6 * time.Hour)
	case "24h":
		startTime = now.AddDate(0, 0, -1)
	case "7d":
		startTime = now.AddDate(0, 0, -7)
	case "30d":
		startTime = now.AddDate(0, 0, -30)
	default:
		startTime = now.Add(-1 * time.Hour)
	}

	var metrics []models.AiMonitoringMetric
	query := c.DB.Model(&models.AiMonitoringMetric{}).Where("recorded_at >= ?", startTime)

	if modelVersionID != "" {
		query = query.Where("model_version_id = ?", modelVersionID)
	}
	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}

	query.Order("recorded_at DESC").Limit(1000).Find(&metrics)

	// 聚合统计数据
	type MetricStat struct {
		MetricName  string  `json:"metric_name"`
		AvgValue    float64 `json:"avg_value"`
		MinValue    float64 `json:"min_value"`
		MaxValue    float64 `json:"max_value"`
		P50Value    float64 `json:"p50_value"`
		P95Value    float64 `json:"p95_value"`
		P99Value    float64 `json:"p99_value"`
		SampleCount int64   `json:"sample_count"`
	}
	var stats []MetricStat

	rows, err := c.DB.Model(&models.AiMonitoringMetric{}).
		Select("metric_name, AVG(metric_value) as avg_value, MIN(metric_value) as min_value, MAX(metric_value) as max_value, COUNT(*) as sample_count").
		Where("recorded_at >= ?", startTime).
		Group("metric_name").
		Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var s MetricStat
			rows.Scan(&s.MetricName, &s.AvgValue, &s.MinValue, &s.MaxValue, &s.SampleCount)
			stats = append(stats, s)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period":      period,
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   now.Format(time.RFC3339),
			"raw_metrics": metrics,
			"aggregated":  stats,
		},
	})
}

// MetricStats 获取指标聚合统计
func (c *AiMonitorController) MetricStats(ctx *gin.Context) {
	_ = ctx.Query("model_id") // reserved for future filtering
	period := ctx.DefaultQuery("period", "24h")

	var startTime time.Time
	now := time.Now()
	switch period {
	case "1h":
		startTime = now.Add(-1 * time.Hour)
	case "6h":
		startTime = now.Add(-6 * time.Hour)
	case "24h":
		startTime = now.AddDate(0, 0, -1)
	case "7d":
		startTime = now.AddDate(0, 0, -7)
	case "30d":
		startTime = now.AddDate(0, 0, -30)
	default:
		startTime = now.AddDate(0, 0, -1)
	}

	// 按指标名称和时间聚合
	type TimeSeriesPoint struct {
		Time  string  `json:"time"`
		Value float64 `json:"value"`
	}
	type MetricTimeSeries struct {
		MetricName string            `json:"metric_name"`
		DataPoints []TimeSeriesPoint `json:"data_points"`
	}
	var timeSeries []MetricTimeSeries

	rows, err := c.DB.Model(&models.AiMonitoringMetric{}).
		Select("metric_name, DATE_FORMAT(recorded_at, '%Y-%m-%d %H:00:00') as time_bucket, AVG(metric_value) as value").
		Where("recorded_at >= ?", startTime).
		Group("metric_name, time_bucket").
		Order("metric_name, time_bucket ASC").
		Rows()
	if err == nil {
		defer rows.Close()
		currentMetric := ""
		var currentSeries *MetricTimeSeries
		for rows.Next() {
			var metricName, timeStr string
			var value float64
			rows.Scan(&metricName, &timeStr, &value)
			if metricName != currentMetric {
				if currentSeries != nil {
					timeSeries = append(timeSeries, *currentSeries)
				}
				currentSeries = &MetricTimeSeries{MetricName: metricName, DataPoints: []TimeSeriesPoint{}}
				currentMetric = metricName
			}
			currentSeries.DataPoints = append(currentSeries.DataPoints, TimeSeriesPoint{Time: timeStr, Value: value})
		}
		if currentSeries != nil {
			timeSeries = append(timeSeries, *currentSeries)
		}
	}

	// 错误率统计
	var errorRateStats struct {
		TotalRequests int64
		FailedRequests int64
		ErrorRate    float64
	}
	c.DB.Model(&models.AiMonitoringMetric{}).
		Select("SUM(CASE WHEN metric_name = 'request_count' THEN metric_value ELSE 0 END) as total_requests, SUM(CASE WHEN metric_name = 'error_count' THEN metric_value ELSE 0 END) as failed_requests").
		Where("recorded_at >= ?", startTime).
		Scan(&errorRateStats)
	if errorRateStats.TotalRequests > 0 {
		errorRateStats.ErrorRate = float64(errorRateStats.FailedRequests) / float64(errorRateStats.TotalRequests) * 100
	}

	// 延迟统计
	var latencyStats struct {
		AvgLatency  float64 `json:"avg_latency"`
		MinLatency  float64 `json:"min_latency"`
		MaxLatency  float64 `json:"max_latency"`
		P50Latency  float64 `json:"p50_latency"`
		P95Latency  float64 `json:"p95_latency"`
		P99Latency  float64 `json:"p99_latency"`
	}
	c.DB.Model(&models.AiMonitoringMetric{}).
		Select("AVG(CASE WHEN percentile = '' THEN metric_value END) as avg_latency, MIN(metric_value) as min_latency, MAX(metric_value) as max_latency").
		Where("recorded_at >= ? AND metric_name = 'latency'", startTime).
		Scan(&latencyStats)

	// QPS统计
	var qpsStats struct {
		CurrentQPS float64 `json:"current_qps"`
		AvgQPS     float64 `json:"avg_qps"`
		MaxQPS     float64 `json:"max_qps"`
	}
	c.DB.Model(&models.AiMonitoringMetric{}).
		Select("AVG(CASE WHEN percentile = 'p50' THEN metric_value END) as avg_qps, MAX(metric_value) as max_qps").
		Where("recorded_at >= ? AND metric_name = 'qps'", startTime).
		Scan(&qpsStats)

	// 最近一次QPS
	var latestQPS float64
	c.DB.Model(&models.AiMonitoringMetric{}).
		Select("metric_value").
		Where("metric_name = 'qps' AND percentile = 'p50'").
		Order("recorded_at DESC").
		Limit(1).
		Scan(&latestQPS)
	qpsStats.CurrentQPS = latestQPS

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period":        period,
			"start_time":    startTime.Format(time.RFC3339),
			"end_time":      now.Format(time.RFC3339),
			"time_series":   timeSeries,
			"error_rate":    errorRateStats,
			"latency":       latencyStats,
			"qps":           qpsStats,
		},
	})
}

// ============ 模型列表 ============

// ListModels 获取模型列表
func (c *AiMonitorController) ListModels(ctx *gin.Context) {
	var modelList []models.AiModel
	var total int64

	query := c.DB.Model(&models.AiModel{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("model_name ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if modelType := ctx.Query("model_type"); modelType != "" {
		query = query.Where("model_type = ?", modelType)
	}

	query.Count(&total)
	query.Order("created_at DESC").Find(&modelList)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  modelList,
			"total": total,
		},
	})
}

// GetModel 获取模型详情
func (c *AiMonitorController) GetModel(ctx *gin.Context) {
	id := ctx.Param("id")
	var model models.AiModel
	if err := c.DB.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取版本列表
	var versions []models.AiModelVersion
	c.DB.Where("model_id = ?", model.ID).Order("created_at DESC").Find(&versions)

	// 获取当前版本的最新指标
	var latestMetrics []models.AiMonitoringMetric
	if model.CurrentVersionID > 0 {
		c.DB.Where("model_version_id = ?", model.CurrentVersionID).
			Order("recorded_at DESC").Limit(100).Find(&latestMetrics)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"model":           model,
			"versions":        versions,
			"latest_metrics": latestMetrics,
		},
	})
}

// ListModelVersions 获取模型版本列表
func (c *AiMonitorController) ListModelVersions(ctx *gin.Context) {
	modelID := ctx.Param("id")
	var versions []models.AiModelVersion
	var total int64

	query := c.DB.Model(&models.AiModelVersion{}).Where("model_id = ?", modelID)

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Order("created_at DESC").Find(&versions)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  versions,
			"total": total,
		},
	})
}

// ============ 告警规则 ============

// ListAlertRules 获取告警规则列表
func (c *AiMonitorController) ListAlertRules(ctx *gin.Context) {
	var rules []models.AiAlertRule
	var total int64

	query := c.DB.Model(&models.AiAlertRule{})

	if enabled := ctx.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}
	if severity := ctx.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	query.Count(&total)
	query.Order("created_at DESC").Find(&rules)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  rules,
			"total": total,
		},
	})
}

// CreateAlertRule 创建告警规则
func (c *AiMonitorController) CreateAlertRule(ctx *gin.Context) {
	var rule models.AiAlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// UpdateAlertRule 更新告警规则
func (c *AiMonitorController) UpdateAlertRule(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.AiAlertRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var updateData models.AiAlertRule
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": updateData})
}

// DeleteAlertRule 删除告警规则
func (c *AiMonitorController) DeleteAlertRule(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.AiAlertRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Delete(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}
