package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AiBehaviorController AI行为分析控制器
type AiBehaviorController struct {
	DB *gorm.DB
}

// ============ 请求结构 ============

type BehaviorListRequest struct {
	Page         int    `form:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" binding:"min=1,max=100"`
	DeviceID     string `form:"device_id"`
	PetUUID      string `form:"pet_uuid"`
	SessionID    string `form:"session_id"`
	BehaviorType string `form:"behavior_type"`
	InferenceMode string `form:"inference_mode"`
	IsAnomaly   *bool  `form:"is_anomaly"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
}

// ============ 行为日志 ============

// List 获取AI行为日志列表
func (c *AiBehaviorController) List(ctx *gin.Context) {
	var req BehaviorListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var list []models.AiBehavior
	var total int64

	query := c.DB.Model(&models.AiBehavior{})

	if req.DeviceID != "" {
		query = query.Where("device_id = ?", req.DeviceID)
	}
	if req.PetUUID != "" {
		query = query.Where("pet_uuid = ?", req.PetUUID)
	}
	if req.SessionID != "" {
		query = query.Where("session_id = ?", req.SessionID)
	}
	if req.BehaviorType != "" {
		query = query.Where("behavior_type = ?", req.BehaviorType)
	}
	if req.InferenceMode != "" {
		query = query.Where("inference_mode = ?", req.InferenceMode)
	}
	if req.IsAnomaly != nil {
		query = query.Where("is_anomaly = ?", *req.IsAnomaly)
	}
	if req.StartDate != "" {
		query = query.Where("occurred_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("occurred_at <= ?", req.EndDate)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("occurred_at DESC").Find(&list)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": list,
			"pagination": gin.H{"total": total, "page": req.Page, "page_size": req.PageSize},
		},
	})
}

// Stats 获取AI行为统计
func (c *AiBehaviorController) Stats(ctx *gin.Context) {
	deviceID := ctx.Query("device_id")
	petUUID := ctx.Query("pet_uuid")
	period := ctx.DefaultQuery("period", "7d") // 7d, 30d

	var startDate time.Time
	now := time.Now()
	switch period {
	case "24h":
		startDate = now.AddDate(0, 0, -1)
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	default:
		startDate = now.AddDate(0, 0, -7)
	}

	// 按行为类型统计
	type BehaviorTypeStat struct {
		BehaviorType string `json:"behavior_type"`
		TotalCount   int64  `json:"total_count"`
		AnomalyCount int64  `json:"anomaly_count"`
		AvgConfidence float64 `json:"avg_confidence"`
		AvgLatencyMs float64 `json:"avg_latency_ms"`
	}
	var behaviorStats []BehaviorTypeStat
	c.DB.Model(&models.AiBehavior{}).
		Select("behavior_type, COUNT(*) as total_count, SUM(CASE WHEN is_anomaly THEN 1 ELSE 0 END) as anomaly_count, AVG(confidence) as avg_confidence, AVG(latency_ms) as avg_latency_ms").
		Where("occurred_at >= ?", startDate).
		Group("behavior_type").
		Scan(&behaviorStats)

	// 推理模式分布
	type InferenceModeStat struct {
		InferenceMode string `json:"inference_mode"`
		Count         int64  `json:"count"`
	}
	var inferenceStats []InferenceModeStat
	c.DB.Model(&models.AiBehavior{}).
		Select("inference_mode, COUNT(*) as count").
		Where("occurred_at >= ?", startDate).
		Group("inference_mode").
		Scan(&inferenceStats)

	// 异常趋势（按天）
	type AnomalyTrend struct {
		Date         string `json:"date"`
		TotalCount   int64  `json:"total_count"`
		AnomalyCount int64  `json:"anomaly_count"`
	}
	var anomalyTrend []AnomalyTrend
	c.DB.Model(&models.AiBehavior{}).
		Select("DATE(occurred_at) as date, COUNT(*) as total_count, SUM(CASE WHEN is_anomaly THEN 1 ELSE 0 END) as anomaly_count").
		Where("occurred_at >= ?", startDate).
		Group("DATE(occurred_at)").
		Order("date ASC").
		Scan(&anomalyTrend)

	// 总体统计
	var totalStats struct {
		TotalCount    int64
		AnomalyCount  int64
		AvgConfidence float64
		AvgLatencyMs  float64
	}
	c.DB.Model(&models.AiBehavior{}).
		Select("COUNT(*) as total_count, SUM(CASE WHEN is_anomaly THEN 1 ELSE 0 END) as anomaly_count, AVG(confidence) as avg_confidence, AVG(latency_ms) as avg_latency_ms").
		Where("occurred_at >= ?", startDate).
		Scan(&totalStats)

	// Top Actions
	var topActions []struct {
		ActionCode string `json:"action_code"`
		ActionName string `json:"action_name"`
		Count      int64  `json:"count"`
	}
	c.DB.Model(&models.AiBehavior{}).
		Select("action_code, action_name, COUNT(*) as count").
		Where("occurred_at >= ? AND action_code != ''", startDate).
		Group("action_code, action_name").
		Order("count DESC").
		Limit(10).
		Scan(&topActions)

	// Device/Pet specific filters
	filterQuery := c.DB.Model(&models.AiBehavior{}).Where("occurred_at >= ?", startDate)
	if deviceID != "" {
		filterQuery = filterQuery.Where("device_id = ?", deviceID)
	}
	if petUUID != "" {
		filterQuery = filterQuery.Where("pet_uuid = ?", petUUID)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period": period,
			"start_date": startDate.Format("2006-01-02"),
			"end_date": now.Format("2006-01-02"),
			"total": gin.H{
				"total_count":    totalStats.TotalCount,
				"anomaly_count":  totalStats.AnomalyCount,
				"anomaly_rate":   float64(totalStats.AnomalyCount) / float64(totalStats.TotalCount) * 100,
				"avg_confidence": totalStats.AvgConfidence,
				"avg_latency_ms": totalStats.AvgLatencyMs,
			},
			"by_behavior_type":    behaviorStats,
			"by_inference_mode":   inferenceStats,
			"anomaly_trend":       anomalyTrend,
			"top_actions":         topActions,
		},
	})
}

// RecordBehavior 记录AI行为（供内部服务调用）
func (c *AiBehaviorController) RecordBehavior(ctx *gin.Context) {
	var behavior models.AiBehavior
	if err := ctx.ShouldBindJSON(&behavior); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if behavior.OccurredAt.IsZero() {
		behavior.OccurredAt = time.Now()
	}

	if err := c.DB.Create(&behavior).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "记录失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": behavior})
}

// GetBehavior 获取单条行为记录
func (c *AiBehaviorController) GetBehavior(ctx *gin.Context) {
	id := ctx.Param("id")
	var behavior models.AiBehavior
	if err := c.DB.First(&behavior, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": behavior})
}

// ============ 异常告警 ============

// ListAlerts 获取异常告警列表
func (c *AiBehaviorController) ListAlerts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var alerts []models.AiAnomalyAlert
	var total int64

	query := c.DB.Model(&models.AiAnomalyAlert{})

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if level := ctx.Query("anomaly_level"); level != "" {
		query = query.Where("anomaly_level = ?", level)
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("occurred_at DESC").Find(&alerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": alerts,
			"pagination": gin.H{"total": total, "page": page, "page_size": pageSize},
		},
	})
}

// AcknowledgeAlert 确认告警
func (c *AiBehaviorController) AcknowledgeAlert(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.AiAnomalyAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	now := time.Now()
	acknowledgedBy := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		acknowledgedBy = uid.(string)
	}

	if err := c.DB.Model(&alert).Updates(map[string]interface{}{
		"status":          "acknowledged",
		"acknowledged_at": now,
		"acknowledged_by": acknowledgedBy,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "确认成功"})
}

// ResolveAlert 解决告警
func (c *AiBehaviorController) ResolveAlert(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.AiAnomalyAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Note string `json:"note"`
	}
	ctx.ShouldBindJSON(&req)

	now := time.Now()
	resolvedBy := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		resolvedBy = uid.(string)
	}

	updates := map[string]interface{}{
		"status":     "resolved",
		"resolved_at": now,
		"resolved_by": resolvedBy,
	}
	if req.Note != "" {
		updates["resolve_note"] = req.Note
	}

	if err := c.DB.Model(&alert).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "解决成功"})
}
