package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIMonitorController AI行为监控控制器
type AIMonitorController struct {
	DB *gorm.DB
}

// NewAIMonitorController 创建控制器
func NewAIMonitorController(db *gorm.DB) *AIMonitorController {
	return &AIMonitorController{DB: db}
}

// ReportEvent 上报 AI 行为事件
// POST /api/v1/ai/monitor/events
func (c *AIMonitorController) ReportEvent(ctx *gin.Context) {
	var event models.AIBehaviorEvent
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 支持批量上报
	var events []models.AIBehaviorEvent
	if list, ok := ctx.Get("events"); ok {
		if e, ok := list.([]models.AIBehaviorEvent); ok {
			events = e
		}
	}
	if len(events) == 0 {
		events = []models.AIBehaviorEvent{event}
	}

	// 转换为日志并批量写入
	logs := make([]models.AIBehaviorLog, 0, len(events))
	now := time.Now()
	for _, e := range events {
		log := e.ToAIBehaviorLog()
		log.CreatedAt = now
		logs = append(logs, *log)
	}

	if err := c.DB.CreateInBatches(logs, 100).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "保存失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"count": len(logs),
		},
	})
}

// ListEvents 行为事件列表（分页+筛选）
// GET /api/v1/ai/monitor/events
func (c *AIMonitorController) ListEvents(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.AIBehaviorLog{})

	// 筛选条件
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if modelName := ctx.Query("model_name"); modelName != "" {
		query = query.Where("model_name = ?", modelName)
	}
	if eventType := ctx.Query("event_type"); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "查询失败: " + err.Error(),
		})
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

// GetEvent 行为事件详情
// GET /api/v1/ai/monitor/events/:id
func (c *AIMonitorController) GetEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "无效的ID",
		})
		return
	}

	var log models.AIBehaviorLog
	if err := c.DB.First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4004,
				"message": "记录不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    log,
	})
}

// GetStats 行为统计（按类型/设备/时间聚合）
// GET /api/v1/ai/monitor/stats
func (c *AIMonitorController) GetStats(ctx *gin.Context) {
	// 默认统计最近24小时
	hours := 24
	if h := ctx.Query("hours"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil && parsed > 0 {
			hours = parsed
		}
	}

	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	// 按事件类型统计
	type EventStat struct {
		EventType string `json:"event_type"`
		Count     int64  `json:"count"`
	}
	var eventStats []EventStat
	c.DB.Model(&models.AIBehaviorLog{}).
		Select("event_type, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("event_type").
		Scan(&eventStats)

	// 按状态统计
	type StatusStat struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var statusStats []StatusStat
	c.DB.Model(&models.AIBehaviorLog{}).
		Select("status, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("status").
		Scan(&statusStats)

	// 按设备统计
	type DeviceStat struct {
		DeviceID string `json:"device_id"`
		Count    int64  `json:"count"`
	}
	var deviceStats []DeviceStat
	c.DB.Model(&models.AIBehaviorLog{}).
		Select("device_id, COUNT(*) as count").
		Where("created_at >= ? AND device_id != ''", since).
		Group("device_id").
		Order("count DESC").
		Limit(10).
		Scan(&deviceStats)

	// 总体指标
	var totalLogs int64
	var avgLatency float64
	var anomalyCount int64
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ?", since).Count(&totalLogs)
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ?", since).Select("COALESCE(AVG(latency_ms), 0)").Scan(&avgLatency)
	c.DB.Model(&models.AIBehaviorLog{}).Where("created_at >= ? AND status = ?", since, "anomaly").Count(&anomalyCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period_hours": hours,
			"summary": gin.H{
				"total_events":   totalLogs,
				"avg_latency_ms": avgLatency,
				"anomaly_count":  anomalyCount,
			},
			"by_event_type": eventStats,
			"by_status":     statusStats,
			"by_device":     deviceStats,
		},
	})
}
