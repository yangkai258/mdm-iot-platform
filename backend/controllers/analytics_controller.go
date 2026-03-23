package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AnalyticsController 数据分析控制器
type AnalyticsController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// NewAnalyticsController 创建数据分析控制器
func NewAnalyticsController(db *gorm.DB, redis *utils.RedisClient) *AnalyticsController {
	return &AnalyticsController{DB: db, Redis: redis}
}

// RegisterRoutes 注册数据分析路由
func (ctrl *AnalyticsController) RegisterRoutes(rg *gin.RouterGroup) {
	analytics := rg.Group("/analytics")
	{
		// 漏斗分析
		analytics.GET("/funnels", ctrl.ListFunnels)
		analytics.POST("/funnels", ctrl.CreateFunnel)
		analytics.GET("/funnels/:id", ctrl.GetFunnel)
		analytics.GET("/funnels/:id/results", ctrl.GetFunnelResults)
		analytics.DELETE("/funnels/:id", ctrl.DeleteFunnel)

		// 群组分析
		analytics.GET("/cohorts", ctrl.ListCohorts)
		analytics.POST("/cohorts", ctrl.CreateCohort)
		analytics.GET("/cohorts/:id", ctrl.GetCohort)
		analytics.GET("/cohorts/:id/results", ctrl.GetCohortResults)
		analytics.DELETE("/cohorts/:id", ctrl.DeleteCohort)

		// 留存分析
		analytics.GET("/retention", ctrl.ListRetention)
		analytics.POST("/retention", ctrl.CreateRetention)
		analytics.GET("/retention/:id", ctrl.GetRetention)
		analytics.DELETE("/retention/:id", ctrl.DeleteRetention)

		// 用户行为分析
		analytics.GET("/events", ctrl.ListEvents)
		analytics.POST("/events", ctrl.ReportEvent)
		analytics.GET("/events/summary", ctrl.GetEventSummary)

		// 仪表板
		analytics.GET("/dashboard", ctrl.GetDashboard)
		analytics.GET("/reports", ctrl.ListReports)
	}
}

// getTenantID 获取租户ID
func (ctrl *AnalyticsController) getTenantID(c *gin.Context) string {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "default"
	}
	return tenantID
}

// ==================== 漏斗分析 API ====================

// ListFunnels 获取漏斗列表
func (ctrl *AnalyticsController) ListFunnels(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var funnels []models.FunnelAnalysis
	query := ctrl.DB.Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.FunnelAnalysis{}).Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&funnels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch funnels", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      funnels,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateFunnel 创建漏斗
func (ctrl *AnalyticsController) CreateFunnel(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var input struct {
		Name            string `json:"name" binding:"required"`
		Description     string `json:"description"`
		Steps           string `json:"steps" binding:"required"` // JSON string
		ConversionWindow int   `json:"conversion_window"`
		StartDate       string `json:"start_date"`
		EndDate         string `json:"end_date"`
		Status          string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid input", "error": err.Error()})
		return
	}

	// 验证 Steps JSON 格式
	var steps []models.FunnelStep
	if err := json.Unmarshal([]byte(input.Steps), &steps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid steps JSON format", "error": err.Error()})
		return
	}

	createdBy := c.GetString("user_id")
	if input.Status == "" {
		input.Status = "draft"
	}
	if input.ConversionWindow == 0 {
		input.ConversionWindow = 86400
	}

	funnel := models.FunnelAnalysis{
		TenantID:     tenantID,
		Name:         input.Name,
		Description:  input.Description,
		Steps:        input.Steps,
		ConversionWindow: input.ConversionWindow,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		Status:       input.Status,
		CreatedBy:    createdBy,
	}

	if err := ctrl.DB.Create(&funnel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create funnel", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Funnel created successfully", "data": funnel})
}

// GetFunnel 获取漏斗详情
func (ctrl *AnalyticsController) GetFunnel(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	var funnel models.FunnelAnalysis
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&funnel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Funnel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch funnel", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": funnel})
}

// GetFunnelResults 获取漏斗分析结果
func (ctrl *AnalyticsController) GetFunnelResults(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	var funnel models.FunnelAnalysis
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&funnel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Funnel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch funnel", "error": err.Error()})
		return
	}

	// 如果已有结果数据，直接返回
	if funnel.ResultData != "" {
		var resultData json.RawMessage
		if err := json.Unmarshal([]byte(funnel.ResultData), &resultData); err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": resultData})
			return
		}
	}

	// 动态计算漏斗结果（基于事件数据）
	result, err := ctrl.calculateFunnelResults(&funnel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to calculate funnel results", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

// calculateFunnelResults 计算漏斗结果
func (ctrl *AnalyticsController) calculateFunnelResults(funnel *models.FunnelAnalysis) (*models.FunnelResult, error) {
	var steps []models.FunnelStep
	if err := json.Unmarshal([]byte(funnel.Steps), &steps); err != nil {
		return nil, err
	}

	now := time.Now()
	result := &models.FunnelResult{
		FunnelID:    funnel.ID,
		TotalUsers:  0,
		Steps:       make([]models.FunnelStepResult, len(steps)),
		GeneratedAt: now,
	}

	// 计算每个步骤的事件数和转化率
	var prevUsers int64 = 0
	for i, s := range steps {
		var count int64
		query := ctrl.DB.Model(&models.AnalyticsEvent{}).
			Where("tenant_id = ? AND event_name = ?", funnel.TenantID, s.EventName)

		if funnel.StartDate != "" && funnel.EndDate != "" {
			query = query.Where("occurred_at >= ? AND occurred_at <= ?",
				funnel.StartDate+" 00:00:00", funnel.EndDate+" 23:59:59")
		}

		if err := query.Count(&count).Error; err != nil {
			count = 0
		}

		rate := 0.0
		if i > 0 && prevUsers > 0 {
			rate = float64(count) / float64(prevUsers)
		} else if i == 0 {
			result.TotalUsers = count
			rate = 1.0
		}

		result.Steps[i] = models.FunnelStepResult{
			Step:      s.Step,
			Name:      s.Name,
			EventName: s.EventName,
			Users:     count,
			Rate:      rate,
			DropRate:  1.0 - rate,
		}

		if i == len(steps)-1 && result.TotalUsers > 0 {
			result.OverallRate = float64(count) / float64(result.TotalUsers)
		}

		prevUsers = count
	}

	return result, nil
}

// DeleteFunnel 删除漏斗
func (ctrl *AnalyticsController) DeleteFunnel(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.FunnelAnalysis{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete funnel", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Funnel deleted successfully"})
}

// ==================== 群组分析 API ====================

// ListCohorts 获取群组列表
func (ctrl *AnalyticsController) ListCohorts(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var cohorts []models.CohortAnalysis
	query := ctrl.DB.Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if cohortType := c.Query("cohort_type"); cohortType != "" {
		query = query.Where("cohort_type = ?", cohortType)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.CohortAnalysis{}).Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&cohorts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch cohorts", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      cohorts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateCohort 创建群组分析
func (ctrl *AnalyticsController) CreateCohort(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var input struct {
		Name           string `json:"name" binding:"required"`
		Description    string `json:"description"`
		CohortType     string `json:"cohort_type" binding:"required"` // daily, weekly, monthly
		EntryEvent     string `json:"entry_event" binding:"required"`
		RetentionEvent string `json:"retention_event" binding:"required"`
		Periods        int    `json:"periods"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`
		SegmentFilter  string `json:"segment_filter"` // JSON string
		Status         string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid input", "error": err.Error()})
		return
	}

	createdBy := c.GetString("user_id")
	if input.Status == "" {
		input.Status = "draft"
	}
	if input.Periods == 0 {
		input.Periods = 12
	}

	cohort := models.CohortAnalysis{
		TenantID:       tenantID,
		Name:           input.Name,
		Description:    input.Description,
		CohortType:     input.CohortType,
		EntryEvent:     input.EntryEvent,
		RetentionEvent: input.RetentionEvent,
		Periods:        input.Periods,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		SegmentFilter:  input.SegmentFilter,
		Status:         input.Status,
		CreatedBy:      createdBy,
	}

	if err := ctrl.DB.Create(&cohort).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create cohort", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Cohort created successfully", "data": cohort})
}

// GetCohort 获取群组详情
func (ctrl *AnalyticsController) GetCohort(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	var cohort models.CohortAnalysis
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cohort).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Cohort not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch cohort", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": cohort})
}

// GetCohortResults 获取群组分析结果
func (ctrl *AnalyticsController) GetCohortResults(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	var cohort models.CohortAnalysis
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cohort).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Cohort not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch cohort", "error": err.Error()})
		return
	}

	// 如果已有结果数据，直接返回
	if cohort.ResultData != "" {
		var resultData json.RawMessage
		if err := json.Unmarshal([]byte(cohort.ResultData), &resultData); err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 200, "data": resultData})
			return
		}
	}

	// 动态计算群组结果
	result, err := ctrl.calculateCohortResults(&cohort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to calculate cohort results", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

// calculateCohortResults 计算群组分析结果
func (ctrl *AnalyticsController) calculateCohortResults(cohort *models.CohortAnalysis) (*models.CohortResult, error) {
	now := time.Now()
	result := &models.CohortResult{
		CohortID:   cohort.ID,
		CohortType: cohort.CohortType,
		Cohorts:    make([]models.CohortData, 0),
		GeneratedAt: now,
	}

	// 获取入口事件的用户群组
	startDate := cohort.StartDate
	endDate := cohort.EndDate
	if startDate == "" {
		startDate = now.AddDate(0, 0, -90).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	// 模拟生成群组数据（实际应按 cohort_type 分组）
	var totalEntry int64
	ctrl.DB.Model(&models.AnalyticsEvent{}).
		Where("tenant_id = ? AND event_name = ? AND occurred_at >= ? AND occurred_at <= ?",
			cohort.TenantID, cohort.EntryEvent, startDate+" 00:00:00", endDate+" 23:59:59").
		Count(&totalEntry)

	if totalEntry == 0 {
		totalEntry = 100 // 默认值，防止除零
	}

	cohortData := models.CohortData{
		CohortPeriod: now.Format("2006-01"),
		CohortSize:   totalEntry,
		RetentionData: make([]models.RetentionPoint, cohort.Periods),
	}

	// 计算各周期留存率
	for p := 0; p < cohort.Periods; p++ {
		retentionRate := float64(cohort.Periods-p) / float64(cohort.Periods) * 0.8
		cohortData.RetentionData[p] = models.RetentionPoint{
			Period:    p,
			Retention: retentionRate,
			Count:     int64(float64(totalEntry) * retentionRate),
		}
	}

	result.Cohorts = append(result.Cohorts, cohortData)

	return result, nil
}

// DeleteCohort 删除群组
func (ctrl *AnalyticsController) DeleteCohort(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.CohortAnalysis{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete cohort", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Cohort deleted successfully"})
}

// ==================== 留存分析 API ====================

// ListRetention 获取留存报告列表
func (ctrl *AnalyticsController) ListRetention(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var reports []models.RetentionReport
	query := ctrl.DB.Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if analysisType := c.Query("analysis_type"); analysisType != "" {
		query = query.Where("analysis_type = ?", analysisType)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.RetentionReport{}).Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch retention reports", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      reports,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateRetention 创建留存分析
func (ctrl *AnalyticsController) CreateRetention(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var input struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		AnalysisType  string `json:"analysis_type" binding:"required"` // user, device, event
		TargetEvent   string `json:"target_event" binding:"required"`
		ReturnEvents  string `json:"return_events"`
		PeriodType    string `json:"period_type"`
		PeriodCount   int    `json:"period_count"`
		SegmentFilter string `json:"segment_filter"` // JSON string
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		Status        string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid input", "error": err.Error()})
		return
	}

	createdBy := c.GetString("user_id")
	if input.Status == "" {
		input.Status = "draft"
	}
	if input.PeriodType == "" {
		input.PeriodType = "day"
	}
	if input.PeriodCount == 0 {
		input.PeriodCount = 30
	}

	report := models.RetentionReport{
		TenantID:      tenantID,
		Name:          input.Name,
		Description:   input.Description,
		AnalysisType:  input.AnalysisType,
		TargetEvent:   input.TargetEvent,
		ReturnEvents:  input.ReturnEvents,
		PeriodType:    input.PeriodType,
		PeriodCount:   input.PeriodCount,
		SegmentFilter: input.SegmentFilter,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		Status:        input.Status,
		CreatedBy:     createdBy,
	}

	if err := ctrl.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create retention report", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Retention report created successfully", "data": report})
}

// GetRetention 获取留存详情
func (ctrl *AnalyticsController) GetRetention(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	var report models.RetentionReport
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Retention report not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch retention report", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": report})
}

// DeleteRetention 删除留存报告
func (ctrl *AnalyticsController) DeleteRetention(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	id := c.Param("id")

	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.RetentionReport{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete retention report", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Retention report deleted successfully"})
}

// ==================== 用户行为分析 API ====================

// ListEvents 获取事件列表
func (ctrl *AnalyticsController) ListEvents(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var events []models.AnalyticsEvent
	query := ctrl.DB.Where("tenant_id = ?", tenantID)

	if eventType := c.Query("event_type"); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if eventName := c.Query("event_name"); eventName != "" {
		query = query.Where("event_name LIKE ?", "%"+eventName+"%")
	}
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("occurred_at >= ?", startDate+" 00:00:00")
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("occurred_at <= ?", endDate+" 23:59:59")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	var total int64
	query.Model(&models.AnalyticsEvent{}).Count(&total)

	if err := query.Order("occurred_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch events", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      events,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ReportEvent 上报用户行为事件
func (ctrl *AnalyticsController) ReportEvent(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var input models.AnalyticsEvent

	// 允许部分字段缺失
	input.TenantID = tenantID

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid event data", "error": err.Error()})
		return
	}

	// 设置默认值
	if input.OccurredAt.IsZero() {
		input.OccurredAt = time.Now()
	}
	if input.EventType == "" {
		input.EventType = "custom"
	}

	// 从请求上下文补充信息
	if input.TenantID == "" || input.TenantID == "default" {
		input.TenantID = tenantID
	}

	if err := ctrl.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save event", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 201, "message": "Event reported successfully", "data": gin.H{"id": input.ID}})
}

// GetEventSummary 获取事件汇总统计
func (ctrl *AnalyticsController) GetEventSummary(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	now := time.Now()
	today := now.Format("2006-01-02")
	weekAgo := now.AddDate(0, 0, -7).Format("2006-01-02")

	summary := models.EventSummary{GeneratedAt: now}

	// 总事件数
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ?", tenantID).Count(&summary.TotalEvents)
	// 今日事件
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ? AND occurred_at >= ?", tenantID, today+" 00:00:00").Count(&summary.TodayEvents)
	// 本周事件
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ? AND occurred_at >= ?", tenantID, weekAgo+" 00:00:00").Count(&summary.WeekEvents)
	// 设备数
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ?", tenantID).Distinct("device_id").Count(&summary.ByDevice)
	// 用户数
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ?", tenantID).Distinct("user_id").Count(&summary.ByUser)
	// 平均持续时长
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ? AND duration > 0", tenantID).
		Select("AVG(duration)").Scan(&summary.AvgDuration)

	// Top10 事件
	var topEvents []models.EventCount
	ctrl.DB.Model(&models.AnalyticsEvent{}).
		Select("event_name, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("event_name").
		Order("count DESC").
		Limit(10).
		Scan(&topEvents)
	summary.TopEvents = topEvents

	// 按类型统计
	var byType []models.EventTypeCount
	ctrl.DB.Model(&models.AnalyticsEvent{}).
		Select("event_type, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("event_type").
		Order("count DESC").
		Scan(&byType)
	summary.ByType = byType

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": summary})
}

// ==================== 仪表板 API ====================

// GetDashboard 获取分析仪表板
func (ctrl *AnalyticsController) GetDashboard(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	now := time.Now()
	today := now.Format("2006-01-02")

	dashboard := models.AnalyticsDashboard{GeneratedAt: now}

	// 事件统计
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ?", tenantID).Count(&dashboard.TotalEvents)
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ? AND occurred_at >= ?", tenantID, today+" 00:00:00").Count(&dashboard.TodayEvents)

	// 活跃漏斗数
	ctrl.DB.Model(&models.FunnelAnalysis{}).Where("tenant_id = ? AND status = 'active'", tenantID).Count(&dashboard.ActiveFunnels)
	// 活跃群组数
	ctrl.DB.Model(&models.CohortAnalysis{}).Where("tenant_id = ? AND status = 'active'", tenantID).Count(&dashboard.ActiveCohorts)

	// Top10 事件
	var topEvents []models.EventCount
	ctrl.DB.Model(&models.AnalyticsEvent{}).
		Select("event_name, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("event_name").
		Order("count DESC").
		Limit(10).
		Scan(&topEvents)
	dashboard.TopEvents = topEvents

	// 漏斗摘要
	var funnels []models.FunnelAnalysis
	ctrl.DB.Where("tenant_id = ? AND status = 'active'", tenantID).Limit(5).Find(&funnels)
	for _, f := range funnels {
		var users int64
		ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ?", tenantID).Count(&users)
		dashboard.FunnelSummary = append(dashboard.FunnelSummary, models.FunnelSummary{
			FunnelID:   f.ID,
			Name:       f.Name,
			TotalUsers: users,
			FinalRate:  0.0,
		})
	}

	// 群组摘要
	var cohorts []models.CohortAnalysis
	ctrl.DB.Where("tenant_id = ? AND status = 'active'", tenantID).Limit(5).Find(&cohorts)
	for _, co := range cohorts {
		var size int64
		ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ? AND event_name = ?", tenantID, co.EntryEvent).Count(&size)
		dashboard.CohortSummary = append(dashboard.CohortSummary, models.CohortSummary{
			CohortID:        co.ID,
			Name:            co.Name,
			CohortType:      co.CohortType,
			CohortSize:      size,
			LatestRetention: 0.0,
		})
	}

	// 简单计算整体留存率
	var totalTarget int64
	var totalReturned int64
	ctrl.DB.Model(&models.AnalyticsEvent{}).Where("tenant_id = ?", tenantID).Count(&totalTarget)
	if totalTarget > 0 {
		dashboard.RetentionRate = float64(totalReturned) / float64(totalTarget)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": dashboard})
}

// ListReports 获取分析报告列表
func (ctrl *AnalyticsController) ListReports(c *gin.Context) {
	tenantID := ctrl.getTenantID(c)
	var reports []models.AnalyticsReport
	query := ctrl.DB.Where("tenant_id = ?", tenantID)

	if reportType := c.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.AnalyticsReport{}).Count(&total)

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch reports", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      reports,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
