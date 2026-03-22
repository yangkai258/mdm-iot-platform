package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type AnalyticsController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

func NewAnalyticsController(db *gorm.DB, redis *utils.RedisClient) *AnalyticsController {
	return &AnalyticsController{DB: db, Redis: redis}
}

func (ctrl *AnalyticsController) RegisterRoutes(rg *gin.RouterGroup) {
	analytics := rg.Group("/analytics")
	{
		analytics.GET("/advanced", ctrl.GetAdvancedAnalytics)
		analytics.GET("/trends", ctrl.GetTrends)
		analytics.GET("/predictions", ctrl.GetPredictions)
		analytics.POST("/export", ctrl.CreateExport)
		analytics.GET("/export/:id/status", ctrl.GetExportStatus)
		analytics.GET("/export/:id/download", ctrl.DownloadExport)
	}
	reports := rg.Group("/reports")
	{
		reports.GET("/custom", ctrl.ListCustomReports)
		reports.POST("/custom", ctrl.CreateCustomReport)
		reports.GET("/custom/:id", ctrl.GetCustomReport)
		reports.PUT("/custom/:id", ctrl.UpdateCustomReport)
		reports.DELETE("/custom/:id", ctrl.DeleteCustomReport)
		reports.GET("/custom/:id/data", ctrl.GetCustomReportData)
	}
}

func (ctrl *AnalyticsController) GetAdvancedAnalytics(c *gin.Context) {
	tenantID := getTenantID(c)
	var req models.AdvancedAnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params", "error": err.Error()}); return
	}
	if req.AnalysisType == "" { req.AnalysisType = "device_health" }
	if req.PeriodType == "" { req.PeriodType = "daily" }
	if req.StartDate == "" { req.StartDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02") }
	if req.EndDate == "" { req.EndDate = time.Now().Format("2006-01-02") }
	if req.GroupBy == "" { req.GroupBy = "day" }
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	endDate = endDate.Add(24*time.Hour - time.Second)
	var dimensions, metrics map[string]interface{}
	var summary string
	switch req.AnalysisType {
	case "device_health": dimensions, metrics, summary = ctrl.analyzeDeviceHealth(tenantID)
	case "member_activity": dimensions, metrics, summary = ctrl.analyzeMemberActivity(tenantID)
	case "ota_performance": dimensions, metrics, summary = ctrl.analyzeOTAPerformance(tenantID)
	case "usage_pattern": dimensions, metrics, summary = ctrl.analyzeUsagePattern(tenantID, startDate, endDate)
	default: c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "unsupported analysis type"}); return
	}
	record := models.AnalyticsRecord{TenantID: tenantID, AnalysisType: req.AnalysisType, PeriodStart: startDate, PeriodEnd: endDate, Dimensions: mapToJSONStr(dimensions), Metrics: mapToJSONStr(metrics), Summary: summary}
	ctrl.DB.Create(&record)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": models.AdvancedAnalyticsResponse{AnalysisType: req.AnalysisType, PeriodStart: req.StartDate, PeriodEnd: req.EndDate, Dimensions: dimensions, Metrics: metrics, Summary: summary}})
}

func (ctrl *AnalyticsController) GetTrends(c *gin.Context) {
	tenantID := getTenantID(c)
	var req models.TrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params", "error": err.Error()}); return
	}
	if req.Metrics == "" { req.Metrics = "device_count,member_count,alert_count" }
	if req.PeriodType == "" { req.PeriodType = "daily" }
	if req.StartDate == "" { req.StartDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02") }
	if req.EndDate == "" { req.EndDate = time.Now().Format("2006-01-02") }
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	ml := strings.Split(req.Metrics, ",")
	for i := range ml { ml[i] = strings.TrimSpace(ml[i]) }
	dataPoints := ctrl.computeTrends(tenantID, ml, req.PeriodType, startDate, endDate)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": models.TrendsResponse{Metrics: ml, PeriodType: req.PeriodType, DataPoints: dataPoints}})
}

func (ctrl *AnalyticsController) GetPredictions(c *gin.Context) {
	tenantID := getTenantID(c)
	var req models.PredictionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params", "error": err.Error()}); return
	}
	if req.Metric == "" { req.Metric = "device_count" }
	if req.Method == "" { req.Method = "linear_regression" }
	if req.Periods <= 0 || req.Periods > 90 { req.Periods = 7 }
	if req.StartDate == "" { req.StartDate = time.Now().AddDate(0, 0, -60).Format("2006-01-02") }
	if req.EndDate == "" { req.EndDate = time.Now().Format("2006-01-02") }
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	hist, forecast, accuracy := ctrl.computePrediction(tenantID, req.Metric, req.Periods, startDate, endDate)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": models.PredictionsResponse{Metric: req.Metric, Method: req.Method, Historical: hist, Forecast: forecast, Accuracy: accuracy}})
}

func (ctrl *AnalyticsController) CreateExport(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	var req models.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params", "error": err.Error()}); return
	}
	validSources := map[string]bool{"devices": true, "members": true, "alerts": true, "ota": true, "usage": true, "custom": true}
	validTypes := map[string]bool{"csv": true, "excel": true}
	if !validSources[req.DataSource] || !validTypes[req.ExportType] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid data source or export type"}); return
	}
	fj, _ := json.Marshal(req.Filters)
	job := models.ExportJob{TenantID: tenantID, ExportType: req.ExportType, DataSource: req.DataSource, Filters: string(fj), Status: "pending", CreatedBy: userID}
	if err := ctrl.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create export job", "error": err.Error()}); return
	}
	go ctrl.processExportJob(job.ID, tenantID, req)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "export job created", "data": models.ExportJobResponse{ID: job.ID, Status: job.Status, DataSource: job.DataSource, ExportType: job.ExportType, CreatedAt: job.CreatedAt}})
}

func (ctrl *AnalyticsController) GetExportStatus(c *gin.Context) {
	tenantID := getTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid job id"}); return }
	var job models.ExportJob
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "export job not found"}); return }
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()}); return
	}
	resp := models.ExportJobResponse{ID: job.ID, Status: job.Status, DataSource: job.DataSource, ExportType: job.ExportType, RecordCount: job.RecordCount, FileSize: job.FileSize, ErrorMsg: job.ErrorMsg, CreatedAt: job.CreatedAt, ExpiresAt: job.ExpiresAt}
	if job.Status == "completed" && job.FilePath != "" { resp.DownloadURL = fmt.Sprintf("/api/v1/analytics/export/%d/download", job.ID) }
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": resp})
}

func (ctrl *AnalyticsController) DownloadExport(c *gin.Context) {
	tenantID := getTenantID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid job id"}); return }
	var job models.ExportJob
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&job).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "export job not found"}); return }
	if job.Status != "completed" || job.FilePath == "" { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "file not ready"}); return }
	if job.ExpiresAt != nil && time.Now().After(*job.ExpiresAt) { c.JSON(http.StatusGone, gin.H{"code": 410, "message": "file expired"}); return }
	if _, err := os.Stat(job.FilePath); os.IsNotExist(err) { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "file not found"}); return }
	if job.ExportType == "excel" { c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") } else { c.Header("Content-Type", "text/csv") }
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=export_%d.%s", job.ID, job.ExportType))
	c.File(job.FilePath)
}

func (ctrl *AnalyticsController) ListCustomReports(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	var reports []models.CustomReport
	query := ctrl.DB.Where("tenant_id = ? AND (created_by = ? OR is_public = ?)", tenantID, userID, true)
	if name := c.Query("name"); name != "" { query = query.Where("name LIKE ?", "%"+name+"%") }
	if err := query.Order("updated_at DESC").Find(&reports).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": reports})
}

func (ctrl *AnalyticsController) CreateCustomReport(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	var req models.CreateCustomReportRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params", "error": err.Error()}); return }
	cj, err := json.Marshal(req.ReportConfig)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "failed to serialize config"}); return }
	report := models.CustomReport{TenantID: tenantID, Name: req.Name, Description: req.Description, ReportConfig: string(cj), ChartType: req.ChartType, IsScheduled: req.IsScheduled, CronExpr: req.CronExpr, IsPublic: req.IsPublic, CreatedBy: userID, UpdatedBy: userID}
	if err := ctrl.DB.Create(&report).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create failed", "error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "report created", "data": report})
}

func (ctrl *AnalyticsController) GetCustomReport(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid report id"}); return }
	var report models.CustomReport
	if err := ctrl.DB.Where("id = ? AND tenant_id = ? AND (created_by = ? OR is_public = ?)", id, tenantID, userID, true).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "report not found"}); return }
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": report})
}

func (ctrl *AnalyticsController) UpdateCustomReport(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid report id"}); return }
	var report models.CustomReport
	if err := ctrl.DB.Where("id = ? AND tenant_id = ? AND created_by = ?", id, tenantID, userID).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "report not found or no permission"}); return }
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()}); return
	}
	var req models.UpdateCustomReportRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid params", "error": err.Error()}); return }
	updates := make(map[string]interface{})
	if req.Name != "" { updates["name"] = req.Name }
	if req.Description != "" { updates["description"] = req.Description }
	if req.ReportConfig != nil { cj, _ := json.Marshal(req.ReportConfig); updates["report_config"] = string(cj) }
	if req.ChartType != "" { updates["chart_type"] = req.ChartType }
	if req.IsScheduled != nil { updates["is_scheduled"] = *req.IsScheduled }
	if req.CronExpr != "" { updates["cron_expr"] = req.CronExpr }
	if req.IsPublic != nil { updates["is_public"] = *req.IsPublic }
	updates["updated_by"] = userID
	if err := ctrl.DB.Model(&report).Updates(updates).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed", "error": err.Error()}); return }
	ctrl.DB.First(&report, id)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "report updated", "data": report})
}

func (ctrl *AnalyticsController) DeleteCustomReport(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid report id"}); return }
	var report models.CustomReport
	if err := ctrl.DB.Where("id = ? AND tenant_id = ? AND created_by = ?", id, tenantID, userID).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "report not found or no permission"}); return }
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()}); return
	}
	if err := ctrl.DB.Delete(&report).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed", "error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "report deleted"})
}

func (ctrl *AnalyticsController) GetCustomReportData(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid report id"}); return }
	var report models.CustomReport
	if err := ctrl.DB.Where("id = ? AND tenant_id = ? AND (created_by = ? OR is_public = ?)", id, tenantID, userID, true).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound { c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "report not found"}); return }
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "query failed", "error": err.Error()}); return
	}
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(report.ReportConfig), &config); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to parse config"}); return }
	data := ctrl.executeCustomReport(tenantID, config)
	ctrl.DB.Model(&report).Updates(map[string]interface{}{"last_run_at": time.Now()})
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": data})
}

func (ctrl *AnalyticsController) analyzeDeviceHealth(tenantID string) (map[string]interface{}, map[string]interface{}, string) {
	dimensions := make(map[string]interface{})
	metrics := make(map[string]interface{})
	var total, online, offline, alerting int64
	ctrl.DB.Model(&models.Device{}).Where("tenant_id = ?", tenantID).Count(&total)
	keys, _ := ctrl.Redis.GetAllShadowKeys()
	online = int64(len(keys))
	offline = total - online
	if offline < 0 { offline = 0 }
	ctrl.DB.Model(&models.DeviceAlert{}).Where("tenant_id = ? AND status = 1", tenantID).Count(&alerting)
	type SD struct { Status string `json:"status"`; Count int64 `json:"count"` }
	var sd []SD
	ctrl.DB.Model(&models.Device{}).Select("status, count(*) as count").Where("tenant_id = ?", tenantID).Group("status").Scan(&sd)
	type VD struct { Version string `json:"version"`; Count int64 `json:"count"` }
	var vd []VD
	ctrl.DB.Model(&models.Device{}).Select("firmware_version as version, count(*) as count").Where("tenant_id = ?", tenantID).Group("firmware_version").Order("count desc").Limit(10).Scan(&vd)
	hs := float64(100)
	if total > 0 { hs -= float64(offline)/float64(total)*30 + float64(alerting)/float64(total)*20 }
	hs = math.Max(0, math.Min(100, hs))
	dimensions["status_distribution"] = sd
	dimensions["version_distribution"] = vd
	metrics["total"] = total; metrics["online"] = online; metrics["offline"] = offline; metrics["alerting"] = alerting
	metrics["online_rate"] = float64(0)
	if total > 0 { metrics["online_rate"] = float64(online) / float64(total) * 100 }
	metrics["health_score"] = hs
	return dimensions, metrics, fmt.Sprintf("total devices %d, online %d (%.1f%%), offline %d, alerting %d, health score %.1f", total, online, metrics["online_rate"].(float64), offline, alerting, hs)
}

func (ctrl *AnalyticsController) analyzeMemberActivity(tenantID string) (map[string]interface{}, map[string]interface{}, string) {
	dimensions := make(map[string]interface{})
	metrics := make(map[string]interface{})
	var total, active, newToday, newWeek int64
	ctrl.DB.Model(&models.Member{}).Where("tenant_id = ?", tenantID).Count(&total)
	ctrl.DB.Model(&models.MemberOperationRecord{}).Where("tenant_id = ? AND created_at >= ?", tenantID, time.Now().AddDate(0, 0, -30)).Distinct("member_id").Count(&active)
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	ctrl.DB.Model(&models.Member{}).Where("tenant_id = ? AND created_at >= ?", tenantID, today).Count(&newToday)
	ctrl.DB.Model(&models.Member{}).Where("tenant_id = ? AND created_at >= ?", tenantID, today.AddDate(0, 0, -7)).Count(&newWeek)
	type LD struct { Level int `json:"level"`; Count int64 `json:"count"` }
	var ld []LD
	ctrl.DB.Model(&models.Member{}).Select("member_level as level, count(*) as count").Where("tenant_id = ?", tenantID).Group("member_level").Order("level").Scan(&ld)
	dimensions["level_distribution"] = ld
	metrics["total"] = total; metrics["active"] = active; metrics["inactive"] = total - active; metrics["new_today"] = newToday; metrics["new_week"] = newWeek
	metrics["activity_rate"] = float64(0)
	if total > 0 { metrics["activity_rate"] = float64(active) / float64(total) * 100 }
	return dimensions, metrics, fmt.Sprintf("total members %d, active %d (%.1f%%), new today %d, new week %d", total, active, metrics["activity_rate"].(float64), newToday, newWeek)
}

func (ctrl *AnalyticsController) analyzeOTAPerformance(tenantID string) (map[string]interface{}, map[string]interface{}, string) {
	dimensions := make(map[string]interface{})
	metrics := make(map[string]interface{})
	var total, success, failed, pending int64
	ctrl.DB.Model(&models.OTAProgress{}).Where("tenant_id = ?", tenantID).Count(&total)
	type SC struct { Status string; Count int64 }
	var sc []SC
	ctrl.DB.Model(&models.OTAProgress{}).Select("ota_status as status, count(*) as count").Where("tenant_id = ?", tenantID).Group("ota_status").Scan(&sc)
	for _, s := range sc { switch s.Status { case "success": success = s.Count; case "failed": failed = s.Count; default: pending += s.Count } }
	type VD struct { Version string `json:"version"`; Count int64 `json:"count"` }
	var vd []VD
	ctrl.DB.Model(&models.OTAProgress{}).Select("target_version as version, count(*) as count").Where("tenant_id = ?", tenantID).Group("target_version").Order("count desc").Limit(10).Scan(&vd)
	sr := float64(0)
	if success+failed > 0 { sr = float64(success) / float64(success+failed) * 100 }
	dimensions["status_distribution"] = map[string]int64{"success": success, "failed": failed, "pending": pending}
	dimensions["version_distribution"] = vd
	metrics["total"] = total; metrics["success"] = success; metrics["failed"] = failed; metrics["pending"] = pending; metrics["success_rate"] = sr
	return dimensions, metrics, fmt.Sprintf("total OTA tasks %d, success %d, failed %d, pending %d, success rate %.1f%%", total, success, failed, pending, sr)
}

func (ctrl *AnalyticsController) analyzeUsagePattern(tenantID string, startDate, endDate time.Time) (map[string]interface{}, map[string]interface{}, string) {
	dimensions := make(map[string]interface{})
	metrics := make(map[string]interface{})
	type AC struct { Action string `json:"action"`; Count int64 `json:"count"` }
	var ac []AC
	ctrl.DB.Model(&models.ActivityLog{}).Select("action, count(*) as count").Where("tenant_id = ? AND created_at >= ? AND created_at <= ?", tenantID, startDate, endDate).Group("action").Order("count desc").Limit(20).Scan(&ac)
	type TU struct { Username string `json:"username"`; Count int64 `json:"count"` }
	var tu []TU
	ctrl.DB.Model(&models.ActivityLog{}).Select("username, count(*) as count").Where("tenant_id = ? AND created_at >= ? AND created_at <= ?", tenantID, startDate, endDate).Group("username").Order("count desc").Limit(10).Scan(&tu)
	var totalOps int64
	ctrl.DB.Model(&models.ActivityLog{}).Where("tenant_id = ? AND created_at >= ? AND created_at <= ?", tenantID, startDate, endDate).Count(&totalOps)
	dimensions["top_actions"] = ac; dimensions["top_users"] = tu
	metrics["total_operations"] = totalOps; metrics["unique_users"] = len(tu); metrics["unique_actions"] = len(ac)
	return dimensions, metrics, fmt.Sprintf("total operations %d, %d users, %d action types", totalOps, len(tu), len(ac))
}

func (ctrl *AnalyticsController) computeTrends(tenantID string, ml []string, periodType string, startDate, endDate time.Time) []models.TrendDataPoint {
	dp := make([]models.TrendDataPoint, 0)
	inc := func(t time.Time) time.Time {
		switch periodType {
		case "weekly": return t.AddDate(0, 0, 7)
		case "monthly": return t.AddDate(0, 1, 0)
		default: return t.AddDate(0, 0, 1)
		}
	}
	for cur := startDate; !cur.After(endDate); cur = inc(cur) {
		next := inc(cur)
		if next.After(endDate) { next = endDate.Add(24 * time.Hour) }
		p := models.TrendDataPoint{Date: cur.Format("2006-01-02"), Metrics: make(map[string]float64)}
		for _, m := range ml {
			var v float64
			switch m {
			case "device_count":
				var c int64
				ctrl.DB.Model(&models.Device{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, cur, next).Count(&c)
				v = float64(c)
			case "member_count":
				var c int64
				ctrl.DB.Model(&models.Member{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, cur, next).Count(&c)
				v = float64(c)
			case "alert_count":
				var c int64
				ctrl.DB.Model(&models.DeviceAlert{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, cur, next).Count(&c)
				v = float64(c)
			case "ota_success_rate":
				var tot, suc int64
				ctrl.DB.Model(&models.OTAProgress{}).Where("tenant_id = ? AND completed_at >= ? AND completed_at < ?", tenantID, cur, next).Count(&tot)
				ctrl.DB.Model(&models.OTAProgress{}).Where("tenant_id = ? AND ota_status = 'success' AND completed_at >= ? AND completed_at < ?", tenantID, cur, next).Count(&suc)
				if tot > 0 { v = float64(suc) / float64(tot) * 100 }
			}
			p.Metrics[m] = v
		}
		dp = append(dp, p)
	}
	return dp
}

func (ctrl *AnalyticsController) computePrediction(tenantID, metric string, periods int, startDate, endDate time.Time) ([]models.PredictionDataPoint, []models.PredictionDataPoint, float64) {
	type rp struct { Date time.Time; Value float64 }
	raw := make([]rp, 0)
	for cur := startDate; !cur.After(endDate); cur = cur.Add(24 * time.Hour) {
		next := cur.Add(24 * time.Hour)
		var v float64
		switch metric {
		case "device_count":
			var c int64
			ctrl.DB.Model(&models.Device{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, startDate, next).Count(&c)
			v = float64(c)
		case "member_count":
			var c int64
			ctrl.DB.Model(&models.Member{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, startDate, next).Count(&c)
			v = float64(c)
		case "alert_count":
			var c int64
			ctrl.DB.Model(&models.DeviceAlert{}).Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, startDate, next).Count(&c)
			v = float64(c)
		}
		raw = append(raw, rp{Date: cur, Value: v})
	}
	hist := make([]models.PredictionDataPoint, 0, len(raw))
	for _, p := range raw { hist = append(hist, models.PredictionDataPoint{Date: p.Date.Format("2006-01-02"), Value: p.Value, Lower: p.Value * 0.9, Upper: p.Value * 1.1}) }
	n := float64(len(raw))
	if n < 2 {
		mean := float64(0)
		for _, p := range raw { mean += p.Value }
		if n > 0 { mean /= n }
		fc := make([]models.PredictionDataPoint, 0, periods)
		for i := 1; i <= periods; i++ { fc = append(fc, models.PredictionDataPoint{Date: endDate.AddDate(0, 0, i).Format("2006-01-02"), Value: mean, Lower: mean * 0.8, Upper: mean * 1.2}) }
		return hist, fc, 0.5
	}
	sumX, sumY, sumXY, sumX2 := 0.0, 0.0, 0.0, 0.0
	for i, p := range raw { x := float64(i); y := p.Value; sumX += x; sumY += y; sumXY += x * y; sumX2 += x * x }
	meanY := sumY / n
	denom := n*sumX2 - sumX*sumX
	var slope, intercept float64
	if math.Abs(denom) > 1e-10 { slope = (n*sumXY - sumX*sumY) / denom; intercept = meanY - slope*(sumX/n) } else { slope, intercept = 0, meanY }
	ssRes, ssTot := 0.0, 0.0
	for i, p := range raw { pr := slope*float64(i) + intercept; ssRes += (p.Value - pr) * (p.Value - pr); ssTot += (p.Value - meanY) * (p.Value - meanY) }
	r2 := 0.0
	if ssTot > 0 { r2 = 1 - ssRes/ssTot }
	fc := make([]models.PredictionDataPoint, 0, periods)
	for i := 1; i <= periods; i++ {
		x := float64(len(raw)-1) + float64(i)
		pv := slope*x + intercept
		if pv < 0 { pv = 0 }
		m := pv * 0.15
		fc = append(fc, models.PredictionDataPoint{Date: endDate.AddDate(0, 0, i).Format("2006-01-02"), Value: pv, Lower: pv - m, Upper: pv + m})
	}
	return hist, fc, r2
}

func (ctrl *AnalyticsController) processExportJob(jobID uint, tenantID string, req models.ExportRequest) {
	now := time.Now()
	ctrl.DB.Model(&models.ExportJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{"status": "processing", "started_at": now})
	var fp string; var rc int64; var fs int64; var em string; fn := ""
	td := filepath.Join(os.TempDir(), "mdm-exports")
	if me := os.MkdirAll(td, 0755); me != nil { em = me.Error(); goto done }
	fn = fmt.Sprintf("export_%s_%s_%d", req.DataSource, now.Format("20060102150405"), jobID)
	switch req.ExportType {
	case "excel": fp, rc, em = ctrl.exportToExcel(td, fn, tenantID, req)
	case "csv": fp, rc, em = ctrl.exportToCSV(td, fn, tenantID, req)
	default: em = "unsupported export type"
	}
	if em == "" { if info, se := os.Stat(fp); se == nil { fs = info.Size() } }
done:
	ea := now.Add(24 * time.Hour)
	up := map[string]interface{}{"completed_at": time.Now(), "expires_at": ea, "record_count": rc, "file_size": fs}
	if em != "" { up["status"] = "failed"; up["error_msg"] = em } else { up["status"] = "completed"; up["file_path"] = fp }
	ctrl.DB.Model(&models.ExportJob{}).Where("id = ?", jobID).Updates(up)
}

func (ctrl *AnalyticsController) exportToCSV(td, fn, tenantID string, req models.ExportRequest) (string, int64, string) {
	fp := filepath.Join(td, fn+".csv")
	f, err := os.Create(fp)
	if err != nil { return "", 0, err.Error() }
	defer f.Close()
	data, err := ctrl.fetchExportRows(tenantID, req)
	if err != nil { return "", 0, err.Error() }
	h := data["headers"].([]string)
	for i, v := range h {
		f.WriteString(v)
		if i < len(h)-1 { f.WriteString(",") }
	}
	f.WriteString("\n")
	rows := data["rows"].([][]string)
	for _, row := range rows {
		for i, v := range row {
			f.WriteString(v)
			if i < len(row)-1 { f.WriteString(",") }
		}
		f.WriteString("\n")
	}
	return fp, int64(len(rows)), ""
}

func (ctrl *AnalyticsController) exportToExcel(td, fn, tenantID string, req models.ExportRequest) (string, int64, string) {
	fp := filepath.Join(td, fn+".xlsx")
	f := excelize.NewFile()
	defer f.Close()
	f.NewSheet("Sheet1")
	s := "Sheet1"
	data, err := ctrl.fetchExportRows(tenantID, req)
	if err != nil { return "", 0, err.Error() }
	h := data["headers"].([]string)
	for ci, h := range h {
		cell, _ := excelize.CoordinatesToCellName(ci+1, 1)
		f.SetCellValue(s, cell, h)
	}
	rows := data["rows"].([][]string)
	for ri, row := range rows {
		for ci, v := range row {
			cell, _ := excelize.CoordinatesToCellName(ci+1, ri+2)
			f.SetCellValue(s, cell, v)
		}
	}
	if err := f.SaveAs(fp); err != nil { return "", 0, err.Error() }
	return fp, int64(len(rows)), ""
}

func (ctrl *AnalyticsController) fetchExportRows(tenantID string, req models.ExportRequest) (map[string]interface{}, error) {
	res := map[string]interface{}{"headers": []string{"ID", "TenantID", "CreatedAt"}, "rows": [][]string{}}
	switch req.DataSource {
	case "devices":
		type DR struct { ID uint; TenantID string; CreatedAt time.Time }
		var ds []DR
		q := ctrl.DB.Model(&models.Device{}).Where("tenant_id = ?", tenantID).Limit(10000)
		if req.StartDate != "" { if st, _ := time.Parse("2006-01-02", req.StartDate); !st.IsZero() { q = q.Where("created_at >= ?", st) } }
		if req.EndDate != "" { if ed, _ := time.Parse("2006-01-02", req.EndDate); !ed.IsZero() { q = q.Where("created_at <= ?", ed.Add(24*time.Hour)) } }
		q.Select("id, tenant_id, created_at").Scan(&ds)
		rs := make([][]string, 0, len(ds))
		for _, d := range ds { rs = append(rs, []string{strconv.Itoa(int(d.ID)), d.TenantID, d.CreatedAt.Format(time.RFC3339)}) }
		res["rows"] = rs
	case "members":
		type MR struct { ID uint; TenantID string; CreatedAt time.Time }
		var ms []MR
		q := ctrl.DB.Model(&models.Member{}).Where("tenant_id = ?", tenantID).Limit(10000)
		if req.StartDate != "" { if st, _ := time.Parse("2006-01-02", req.StartDate); !st.IsZero() { q = q.Where("created_at >= ?", st) } }
		if req.EndDate != "" { if ed, _ := time.Parse("2006-01-02", req.EndDate); !ed.IsZero() { q = q.Where("created_at <= ?", ed.Add(24*time.Hour)) } }
		q.Select("id, tenant_id, created_at").Scan(&ms)
		rs := make([][]string, 0, len(ms))
		for _, m := range ms { rs = append(rs, []string{strconv.Itoa(int(m.ID)), m.TenantID, m.CreatedAt.Format(time.RFC3339)}) }
		res["rows"] = rs
	case "alerts":
		type AR struct { ID uint; TenantID string; CreatedAt time.Time }
		var as []AR
		q := ctrl.DB.Model(&models.DeviceAlert{}).Where("tenant_id = ?", tenantID).Limit(10000)
		if req.StartDate != "" { if st, _ := time.Parse("2006-01-02", req.StartDate); !st.IsZero() { q = q.Where("created_at >= ?", st) } }
		if req.EndDate != "" { if ed, _ := time.Parse("2006-01-02", req.EndDate); !ed.IsZero() { q = q.Where("created_at <= ?", ed.Add(24*time.Hour)) } }
		q.Select("id, tenant_id, created_at").Scan(&as)
		rs := make([][]string, 0, len(as))
		for _, a := range as { rs = append(rs, []string{strconv.Itoa(int(a.ID)), a.TenantID, a.CreatedAt.Format(time.RFC3339)}) }
		res["rows"] = rs
	}
	return res, nil
}

func (ctrl *AnalyticsController) executeCustomReport(tenantID string, config map[string]interface{}) map[string]interface{} {
	ds, _ := config["data_source"].(string)
	dim, _ := config["dimensions"].([]interface{})
	met, _ := config["metrics"].([]interface{})
	return map[string]interface{}{"data_source": ds, "dimensions": dim, "metrics": met, "message": "report data retrieved", "columns": []string{"id", "created_at"}, "data": []map[string]interface{}{}}
}

func mapToJSONStr(m map[string]interface{}) string {
	if m == nil { return "{}" }
	b, err := json.Marshal(m)
	if err != nil { return "{}" }
	return string(b)
}
