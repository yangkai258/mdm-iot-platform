package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
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

// ReportController 报表统计控制器
type ReportController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// NewReportController 创建报表控制器
func NewReportController(db *gorm.DB, redis *utils.RedisClient) *ReportController {
	return &ReportController{DB: db, Redis: redis}
}

// RegisterRoutes 注册报表路由
func (ctrl *ReportController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/reports/stats", ctrl.GetStats)
	rg.GET("/reports/devices", ctrl.GetDeviceStats)
	rg.GET("/reports/members", ctrl.GetMemberStats)
	rg.GET("/reports/ota", ctrl.GetOTAStats)
	rg.GET("/reports/trend", ctrl.GetTrend)
	rg.GET("/reports/export", ctrl.ExportReport)
	rg.POST("/reports/generate", ctrl.GenerateReport)
}

// GetStats 获取综合统计数据
// @Summary 综合统计
// @Description 获取设备、会员、OTA的综合统计数据
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} models.ReportResponse
// @Router /api/v1/reports/stats [get]
func (ctrl *ReportController) GetStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	// 设备统计
	deviceStats, err := ctrl.getDeviceStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取设备统计失败", "error": err.Error()})
		return
	}

	// 会员统计
	memberStats, err := ctrl.getMemberStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取会员统计失败", "error": err.Error()})
		return
	}

	// OTA统计
	otaStats, err := ctrl.getOTAStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取OTA统计失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": models.ReportResponse{
			ReportType: "current",
			ReportDate: time.Now().Format("2006-01-02"),
			Device:     deviceStats,
			Member:     memberStats,
			OTA:        otaStats,
		},
	})
}

// GetDeviceStats 获取设备统计
// @Summary 设备统计
// @Description 获取设备总数、在线数、告警数等
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} models.DeviceStats
// @Router /api/v1/reports/devices [get]
func (ctrl *ReportController) GetDeviceStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	stats, err := ctrl.getDeviceStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取设备统计失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": stats})
}

// GetMemberStats 获取会员统计
// @Summary 会员统计
// @Description 获取会员总数、活跃数、新增数
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} models.MemberStats
// @Router /api/v1/reports/members [get]
func (ctrl *ReportController) GetMemberStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	stats, err := ctrl.getMemberStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取会员统计失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": stats})
}

// GetOTAStats 获取OTA统计
// @Summary OTA统计
// @Description 获取OTA升级成功率、版本分布
// @Tags reports
// @Accept json
// @Produce json
// @Success 200 {object} models.OTAStats
// @Router /api/v1/reports/ota [get]
func (ctrl *ReportController) GetOTAStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	stats, err := ctrl.getOTAStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取OTA统计失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": stats})
}

// GetTrend 获取每日趋势
// @Summary 每日趋势
// @Description 获取最近N天的趋势数据
// @Tags reports
// @Param days query int false "天数，默认7天"
// @Accept json
// @Produce json
// @Success 200 {array} models.DailyTrend
// @Router /api/v1/reports/trend [get]
func (ctrl *ReportController) GetTrend(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	trends, err := ctrl.getDailyTrends(tenantID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取趋势数据失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": trends})
}

// ExportReport 导出报表
// @Summary 导出报表
// @Description 导出CSV或Excel格式的报表数据
// @Tags reports
// @Param report_type query string false "报表类型: daily, weekly, monthly"
// @Param export_type query string false "导出格式: csv, excel"
// @Param start_date query string false "开始日期 YYYY-MM-DD"
// @Param end_date query string false "结束日期 YYYY-MM-DD"
// @Accept json
// @Produce json
// @Success 200 {file} file
// @Router /api/v1/reports/export [get]
func (ctrl *ReportController) ExportReport(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	exportType := c.DefaultQuery("export_type", "csv")
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	reportType := c.DefaultQuery("report_type", "daily")

	// 获取统计数据
	deviceStats, err := ctrl.getDeviceStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取设备统计失败"})
		return
	}
	memberStats, err := ctrl.getMemberStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取会员统计失败"})
		return
	}
	otaStats, err := ctrl.getOTAStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取OTA统计失败"})
		return
	}

	// 生成临时目录
	tmpDir := filepath.Join(os.TempDir(), "mdm-reports")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建临时目录失败"})
		return
	}

	filename := fmt.Sprintf("report_%s_%s_%s", reportType, startDate, endDate)

	switch strings.ToLower(exportType) {
	case "excel", "xlsx":
		filePath := filepath.Join(tmpDir, filename+".xlsx")
		if err := ctrl.exportToExcel(filePath, reportType, startDate, endDate, deviceStats, memberStats, otaStats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出Excel失败", "error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
		c.File(filePath)
	default:
		filePath := filepath.Join(tmpDir, filename+".csv")
		if err := ctrl.exportToCSV(filePath, reportType, startDate, endDate, deviceStats, memberStats, otaStats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出CSV失败", "error": err.Error()})
			return
		}
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
		c.File(filePath)
	}
}

// GenerateReport 生成并存储报表
// @Summary 生成报表
// @Description 生成日报/周报/月报并存入数据库
// @Tags reports
// @Param report_type body string true "报表类型: daily, weekly, monthly"
// @Accept json
// @Produce json
// @Success 200 {object} models.ReportRecord
// @Router /api/v1/reports/generate [post]
func (ctrl *ReportController) GenerateReport(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}

	var req struct {
		ReportType string `json:"report_type" binding:"required"`
		ReportDate string `json:"report_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 兼容不传body的情况，自动使用当天
		req.ReportType = c.DefaultQuery("report_type", "daily")
		req.ReportDate = time.Now().Format("2006-01-02")
	}
	if req.ReportDate == "" {
		req.ReportDate = time.Now().Format("2006-01-02")
	}

	// 收集统计数据
	deviceStats, err := ctrl.getDeviceStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取设备统计失败"})
		return
	}
	memberStats, err := ctrl.getMemberStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取会员统计失败"})
		return
	}
	otaStats, err := ctrl.getOTAStats(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取OTA统计失败"})
		return
	}

	deviceJSON, _ := json.Marshal(deviceStats)
	memberJSON, _ := json.Marshal(memberStats)
	otaJSON, _ := json.Marshal(otaStats)

	record := models.ReportRecord{
		TenantID:    tenantID,
		ReportType:  req.ReportType,
		ReportDate:  req.ReportDate,
		DeviceStats: string(deviceJSON),
		MemberStats: string(memberJSON),
		OTAStats:    string(otaJSON),
	}

	if err := ctrl.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存报表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "报表生成成功", "data": record})
}

// ===================== 内部方法 =====================

// getDeviceStats 获取设备统计数据
func (ctrl *ReportController) getDeviceStats(tenantID string) (models.DeviceStats, error) {
	stats := models.DeviceStats{}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	// 设备总数（按租户）
	if err := ctrl.DB.Model(&models.Device{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.Total).Error; err != nil {
		return stats, err
	}

	// 今日新增
	if err := ctrl.DB.Model(&models.Device{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, today).
		Count(&stats.NewToday).Error; err != nil {
		return stats, err
	}

	// 本周新增
	if err := ctrl.DB.Model(&models.Device{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, weekAgo).
		Count(&stats.NewWeek).Error; err != nil {
		return stats, err
	}

	// 本月新增
	if err := ctrl.DB.Model(&models.Device{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, monthAgo).
		Count(&stats.NewMonth).Error; err != nil {
		return stats, err
	}

	// 活跃设备（90天内有心跳）
	activeThreshold := now.AddDate(0, 0, -90)
	if err := ctrl.DB.Model(&models.DeviceShadow{}).
		Joins("JOIN mdm_devices ON mdm_devices.device_id = device_shadows.device_id").
		Where("mdm_devices.tenant_id = ? AND device_shadows.last_heartbeat >= ?", tenantID, activeThreshold).
		Count(&stats.Active).Error; err != nil {
		// 影子表可能为空，降级处理
		stats.Active = 0
	}

	// 不活跃设备
	stats.Inactive = stats.Total - stats.Active

	// 在线设备：从 Redis 获取
	keys, err := ctrl.Redis.GetAllShadowKeys()
	if err != nil {
		// Redis 不可用时降级
		stats.Online = 0
		stats.Offline = stats.Total
	} else {
		stats.Online = int64(len(keys))
		stats.Offline = stats.Total - stats.Online
	}

	// 告警设备数（未处理的告警）
	if err := ctrl.DB.Model(&models.DeviceAlert{}).
		Where("tenant_id = ? AND status = 1", tenantID).
		Count(&stats.Alerting).Error; err != nil {
		stats.Alerting = 0
	}

	return stats, nil
}

// getMemberStats 获取会员统计数据
func (ctrl *ReportController) getMemberStats(tenantID string) (models.MemberStats, error) {
	stats := models.MemberStats{}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	// 会员总数
	if err := ctrl.DB.Model(&models.Member{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.Total).Error; err != nil {
		return stats, err
	}

	// 活跃会员（30天内有操作记录）
	activeThreshold := now.AddDate(0, 0, -30)
	var activeCount int64
	if err := ctrl.DB.Model(&models.MemberOperationRecord{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, activeThreshold).
		Distinct("member_id").
		Count(&activeCount).Error; err != nil {
		stats.Active = 0
	} else {
		stats.Active = activeCount
	}

	// 今日新增
	if err := ctrl.DB.Model(&models.Member{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, today).
		Count(&stats.NewToday).Error; err != nil {
		return stats, err
	}

	// 本周新增
	if err := ctrl.DB.Model(&models.Member{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, weekAgo).
		Count(&stats.NewWeek).Error; err != nil {
		return stats, err
	}

	// 本月新增
	if err := ctrl.DB.Model(&models.Member{}).
		Where("tenant_id = ? AND created_at >= ?", tenantID, monthAgo).
		Count(&stats.NewMonth).Error; err != nil {
		return stats, err
	}

	// 等级分布
	var levelCounts []models.LevelCount
	if err := ctrl.DB.Model(&models.Member{}).
		Select("member_level as level, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("member_level").
		Order("member_level").
		Scan(&levelCounts).Error; err != nil {
		stats.ByLevel = []models.LevelCount{}
	} else {
		stats.ByLevel = levelCounts
	}

	return stats, nil
}

// getOTAStats 获取OTA统计数据
func (ctrl *ReportController) getOTAStats(tenantID string) (models.OTAStats, error) {
	stats := models.OTAStats{}

	// 部署任务总数
	if err := ctrl.DB.Model(&models.OTADeployment{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.TotalDeployments).Error; err != nil {
		return stats, err
	}

	// 成功/失败/待处理
	if err := ctrl.DB.Model(&models.OTAProgress{}).
		Where("tenant_id = ?", tenantID).
		Group("ota_status").
		Select("ota_status, count(*) as count").
		Scan(&[]struct {
		OTAStatus string
		Count     int64
	}{}).Error; err != nil {
		// ignore
	}

	type StatusCount struct {
		OTAStatus string
		Count     int64
	}
	var statusCounts []StatusCount
	ctrl.DB.Model(&models.OTAProgress{}).
		Where("tenant_id = ?", tenantID).
		Select("ota_status, count(*) as count").
		Group("ota_status").
		Scan(&statusCounts)

	for _, sc := range statusCounts {
		switch sc.OTAStatus {
		case "success":
			stats.SuccessCount = sc.Count
		case "failed":
			stats.FailedCount = sc.Count
		default:
			stats.PendingCount += sc.Count
		}
	}

	// 计算成功率
	total := stats.SuccessCount + stats.FailedCount
	if total > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(total) * 100
	}

	// 固件版本分布
	var versionCounts []models.VersionCount
	if err := ctrl.DB.Model(&models.Device{}).
		Select("firmware_version as version, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("firmware_version").
		Order("count(*) desc").
		Limit(20).
		Scan(&versionCounts).Error; err != nil {
		stats.VersionDistribution = []models.VersionCount{}
	} else {
		stats.VersionDistribution = versionCounts
	}

	return stats, nil
}

// getDailyTrends 获取每日趋势数据
func (ctrl *ReportController) getDailyTrends(tenantID string, days int) ([]models.DailyTrend, error) {
	trends := make([]models.DailyTrend, 0, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		dayEnd := dayStart.AddDate(0, 0, 1)

		trend := models.DailyTrend{Date: dateStr}

		// 当日设备总数
		ctrl.DB.Model(&models.Device{}).
			Where("tenant_id = ?", tenantID).
			Count(&trend.DeviceTotal)

		// 当日新增设备
		ctrl.DB.Model(&models.Device{}).
			Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, dayStart, dayEnd).
			Count(&trend.DeviceNew)

		// 当日会员总数
		ctrl.DB.Model(&models.Member{}).
			Where("tenant_id = ?", tenantID).
			Count(&trend.MemberTotal)

		// 当日新增会员
		ctrl.DB.Model(&models.Member{}).
			Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, dayStart, dayEnd).
			Count(&trend.MemberNew)

		// 当日OTA完成数
		ctrl.DB.Model(&models.OTAProgress{}).
			Where("tenant_id = ? AND completed_at >= ? AND completed_at < ?", tenantID, dayStart, dayEnd).
			Count(&trend.OTAComplete)

		// 当日OTA成功率
		var successCount, totalCount int64
		ctrl.DB.Model(&models.OTAProgress{}).
			Where("tenant_id = ? AND completed_at >= ? AND completed_at < ?", tenantID, dayStart, dayEnd).
			Count(&totalCount)
		ctrl.DB.Model(&models.OTAProgress{}).
			Where("tenant_id = ? AND ota_status = 'success' AND completed_at >= ? AND completed_at < ?", tenantID, dayStart, dayEnd).
			Count(&successCount)
		if totalCount > 0 {
			trend.OTASuccessRate = float64(successCount) / float64(totalCount) * 100
		}

		trends = append(trends, trend)
	}

	return trends, nil
}

// exportToCSV 导出为CSV
func (ctrl *ReportController) exportToCSV(
	filePath, reportType, startDate, endDate string,
	deviceStats models.DeviceStats,
	memberStats models.MemberStats,
	otaStats models.OTAStats,
) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 标题行
	writer.Write([]string{"报表类型", "统计维度", "指标名称", "数值"})

	// 设备统计
	deviceRows := [][]string{
		{reportType, "设备统计", "设备总数", strconv.FormatInt(deviceStats.Total, 10)},
		{reportType, "设备统计", "在线设备", strconv.FormatInt(deviceStats.Online, 10)},
		{reportType, "设备统计", "离线设备", strconv.FormatInt(deviceStats.Offline, 10)},
		{reportType, "设备统计", "告警设备", strconv.FormatInt(deviceStats.Alerting, 10)},
		{reportType, "设备统计", "活跃设备", strconv.FormatInt(deviceStats.Active, 10)},
		{reportType, "设备统计", "今日新增", strconv.FormatInt(deviceStats.NewToday, 10)},
		{reportType, "设备统计", "本周新增", strconv.FormatInt(deviceStats.NewWeek, 10)},
		{reportType, "设备统计", "本月新增", strconv.FormatInt(deviceStats.NewMonth, 10)},
	}
	for _, row := range deviceRows {
		writer.Write(row)
	}

	// 会员统计
	memberRows := [][]string{
		{reportType, "会员统计", "会员总数", strconv.FormatInt(memberStats.Total, 10)},
		{reportType, "会员统计", "活跃会员", strconv.FormatInt(memberStats.Active, 10)},
		{reportType, "会员统计", "今日新增", strconv.FormatInt(memberStats.NewToday, 10)},
		{reportType, "会员统计", "本周新增", strconv.FormatInt(memberStats.NewWeek, 10)},
		{reportType, "会员统计", "本月新增", strconv.FormatInt(memberStats.NewMonth, 10)},
	}
	for _, row := range memberRows {
		writer.Write(row)
	}

	// OTA统计
	otaRows := [][]string{
		{reportType, "OTA统计", "部署任务总数", strconv.FormatInt(otaStats.TotalDeployments, 10)},
		{reportType, "OTA统计", "成功数", strconv.FormatInt(otaStats.SuccessCount, 10)},
		{reportType, "OTA统计", "失败数", strconv.FormatInt(otaStats.FailedCount, 10)},
		{reportType, "OTA统计", "待处理数", strconv.FormatInt(otaStats.PendingCount, 10)},
		{reportType, "OTA统计", "成功率(%)", fmt.Sprintf("%.2f", otaStats.SuccessRate)},
	}
	for _, row := range otaRows {
		writer.Write(row)
	}

	// 固件版本分布
	for _, vc := range otaStats.VersionDistribution {
		writer.Write([]string{
			reportType, "固件版本分布",
			fmt.Sprintf("v%s", vc.Version),
			strconv.FormatInt(vc.Count, 10),
		})
	}

	// 会员等级分布
	for _, lc := range memberStats.ByLevel {
		writer.Write([]string{
			reportType, "会员等级分布",
			fmt.Sprintf("等级%d", lc.Level),
			strconv.FormatInt(lc.Count, 10),
		})
	}

	return nil
}

// exportToExcel 导出为Excel
func (ctrl *ReportController) exportToExcel(
	filePath, reportType, startDate, endDate string,
	deviceStats models.DeviceStats,
	memberStats models.MemberStats,
	otaStats models.OTAStats,
) error {
	f := excelize.NewFile()
	defer f.Close()

	// 删除默认 Sheet
	f.DeleteSheet("Sheet1")
	// 创建各 Sheet（返回 sheet 索引）
	f.NewSheet("设备统计")
	f.NewSheet("会员统计")
	f.NewSheet("OTA统计")
	overviewSheetIdx, _ := f.NewSheet("报表概览")

	// --- 报表概览 Sheet ---
	f.SetCellValue("报表概览", "A1", "报表类型")
	f.SetCellValue("报表概览", "B1", reportType)
	f.SetCellValue("报表概览", "A2", "开始日期")
	f.SetCellValue("报表概览", "B2", startDate)
	f.SetCellValue("报表概览", "A3", "结束日期")
	f.SetCellValue("报表概览", "B3", endDate)
	f.SetCellValue("报表概览", "A4", "生成时间")
	f.SetCellValue("报表概览", "B4", time.Now().Format("2006-01-02 15:04:05"))

	// --- 设备统计 Sheet ---
	f.SetCellValue("设备统计", "A1", "指标")
	f.SetCellValue("设备统计", "B1", "数值")
	deviceData := [][]interface{}{
		{"设备总数", deviceStats.Total},
		{"在线设备", deviceStats.Online},
		{"离线设备", deviceStats.Offline},
		{"告警设备", deviceStats.Alerting},
		{"活跃设备", deviceStats.Active},
		{"不活跃设备", deviceStats.Inactive},
		{"今日新增", deviceStats.NewToday},
		{"本周新增", deviceStats.NewWeek},
		{"本月新增", deviceStats.NewMonth},
	}
	for i, row := range deviceData {
		f.SetCellValue("设备统计", fmt.Sprintf("A%d", i+2), row[0])
		f.SetCellValue("设备统计", fmt.Sprintf("B%d", i+2), row[1])
	}

	// --- 会员统计 Sheet ---
	f.SetCellValue("会员统计", "A1", "指标")
	f.SetCellValue("会员统计", "B1", "数值")
	memberData := [][]interface{}{
		{"会员总数", memberStats.Total},
		{"活跃会员", memberStats.Active},
		{"今日新增", memberStats.NewToday},
		{"本周新增", memberStats.NewWeek},
		{"本月新增", memberStats.NewMonth},
	}
	for i, row := range memberData {
		f.SetCellValue("会员统计", fmt.Sprintf("A%d", i+2), row[0])
		f.SetCellValue("会员统计", fmt.Sprintf("B%d", i+2), row[1])
	}

	// 会员等级分布
	f.SetCellValue("会员统计", "D1", "等级")
	f.SetCellValue("会员统计", "E1", "人数")
	for i, lc := range memberStats.ByLevel {
		f.SetCellValue("会员统计", fmt.Sprintf("D%d", i+2), fmt.Sprintf("等级%d", lc.Level))
		f.SetCellValue("会员统计", fmt.Sprintf("E%d", i+2), lc.Count)
	}

	// --- OTA统计 Sheet ---
	f.SetCellValue("OTA统计", "A1", "指标")
	f.SetCellValue("OTA统计", "B1", "数值")
	otaData := [][]interface{}{
		{"部署任务总数", otaStats.TotalDeployments},
		{"成功数", otaStats.SuccessCount},
		{"失败数", otaStats.FailedCount},
		{"待处理数", otaStats.PendingCount},
		{"成功率(%)", fmt.Sprintf("%.2f%%", otaStats.SuccessRate)},
	}
	for i, row := range otaData {
		f.SetCellValue("OTA统计", fmt.Sprintf("A%d", i+2), row[0])
		f.SetCellValue("OTA统计", fmt.Sprintf("B%d", i+2), row[1])
	}

	// 固件版本分布
	f.SetCellValue("OTA统计", "D1", "固件版本")
	f.SetCellValue("OTA统计", "E1", "设备数量")
	for i, vc := range otaStats.VersionDistribution {
		f.SetCellValue("OTA统计", fmt.Sprintf("D%d", i+2), fmt.Sprintf("v%s", vc.Version))
		f.SetCellValue("OTA统计", fmt.Sprintf("E%d", i+2), vc.Count)
	}

	// 设置活动 Sheet
	f.SetActiveSheet(overviewSheetIdx)

	return f.SaveAs(filePath)
}
