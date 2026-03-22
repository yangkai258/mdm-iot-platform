package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIQualityController AI质量指标控制器
type AIQualityController struct {
	DB *gorm.DB
}

// NewAIQualityController 创建控制器
func NewAIQualityController(db *gorm.DB) *AIQualityController {
	return &AIQualityController{DB: db}
}

// GetMetrics 质量指标总览
// GET /api/v1/ai/quality/metrics
func (c *AIQualityController) GetMetrics(ctx *gin.Context) {
	// 默认统计最近24小时
	hours := 24
	if h := ctx.Query("hours"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	type Metric struct {
		ModelID    string  `json:"model_id"`
		ModelName  string  `json:"model_name"`
		AvgLatency float64 `json:"avg_latency_ms"`
		AvgError   float64 `json:"avg_error_rate"`
		AvgConf    float64 `json:"avg_confidence"`
		TotalCount int64   `json:"total_count"`
		AnomalyCnt int64   `json:"anomaly_count"`
	}

	var metrics []Metric
	c.DB.Model(&models.AIBehaviorLog{}).
		Select(`model_id, model_name,
			COALESCE(AVG(latency_ms), 0) as avg_latency,
			COALESCE(AVG(error_rate), 0) as avg_error,
			COALESCE(AVG(confidence), 0) as avg_conf,
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'anomaly' THEN 1 ELSE 0 END) as anomaly_cnt`).
		Where("created_at >= ?", since).
		Group("model_id, model_name").
		Scan(&metrics)

	// 总体指标
	var totalLogs int64
	var avgLatency, avgErrorRate, avgConfidence float64
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ?", since).Count(&totalLogs)
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ?", since).Select("COALESCE(AVG(latency_ms),0)").Scan(&avgLatency)
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ?", since).Select("COALESCE(AVG(error_rate),0)").Scan(&avgErrorRate)
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ?", since).Select("COALESCE(AVG(confidence),0)").Scan(&avgConfidence)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period_hours": hours,
			"overall": gin.H{
				"total_events":   totalLogs,
				"avg_latency_ms": avgLatency,
				"avg_error_rate": avgErrorRate,
				"avg_confidence": avgConfidence,
			},
			"by_model": metrics,
		},
	})
}

// GetMetricsTrend 质量趋势（时间序列）
// GET /api/v1/ai/quality/metrics/trend
func (c *AIQualityController) GetMetricsTrend(ctx *gin.Context) {
	hours := 24
	if h := ctx.Query("hours"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	modelID := ctx.Query("model_id")

	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	type TrendPoint struct {
		Time         time.Time `json:"time"`
		AvgLatency   float64   `json:"avg_latency_ms"`
		AvgError     float64   `json:"avg_error_rate"`
		AvgConf      float64   `json:"avg_confidence"`
		EventCount   int64     `json:"event_count"`
		AnomalyCount int64     `json:"anomaly_count"`
	}

	var trends []TrendPoint

	query := c.DB.Model(&models.AIBehaviorLog{}).
		Where("created_at >= ?", since)
	if modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}

	// 按小时聚合
	interval := "1 hour"
	if hours <= 6 {
		interval = "30 minutes"
	}

	c.DB.Raw(`
		SELECT
			date_trunc(?, created_at) as time,
			COALESCE(AVG(latency_ms), 0) as avg_latency,
			COALESCE(AVG(error_rate), 0) as avg_error,
			COALESCE(AVG(confidence), 0) as avg_conf,
			COUNT(*) as event_count,
			SUM(CASE WHEN status = 'anomaly' THEN 1 ELSE 0 END) as anomaly_count
		FROM ai_behavior_logs
		WHERE created_at >= ?
		`+func() string {
		if modelID != "" {
			return " AND model_id = ?"
		}
		return ""
	}()+`
		GROUP BY date_trunc(?, created_at)
		ORDER BY time ASC
	`, interval, since, func() interface{} {
		if modelID != "" {
			return modelID
		}
		return nil
	}(), interval).Scan(&trends)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period_hours": hours,
			"model_id":     modelID,
			"trend":        trends,
		},
	})
}

// GetAnomalyEvents 异常事件列表
// GET /api/v1/ai/quality/anomaly
func (c *AIQualityController) GetAnomalyEvents(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIBehaviorLog{}).Where("status = ?", "anomaly")

	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if modelName := ctx.Query("model_name"); modelName != "" {
		query = query.Where("model_name = ?", modelName)
	}
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var logs []models.AIBehaviorLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": logs,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}
