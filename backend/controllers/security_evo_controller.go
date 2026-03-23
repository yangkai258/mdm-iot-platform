package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SecurityEvoController Sprint 32 高级安全功能控制器
type SecurityEvoController struct {
	DB *gorm.DB
}

// RegisterSecurityEvoRoutes 注册安全演进路由
func (ctrl *SecurityEvoController) RegisterSecurityEvoRoutes(rg *gin.RouterGroup) {
	sec := &SecurityEvoController{DB: ctrl.DB}

	// 安全审计
	audit := rg.Group("/security/audit")
	{
		audit.GET("/logs", sec.GetAuditLogs)
		audit.GET("/reports", sec.ListAuditReports)
		audit.POST("/reports", sec.CreateAuditReport)
		audit.GET("/reports/:id", sec.GetAuditReport)
	}

	// 合规管理
	compliance := rg.Group("/security/compliance")
	{
		compliance.GET("/reports", sec.ListComplianceReports)
		compliance.POST("/reports", sec.CreateComplianceReport)
		compliance.GET("/status", sec.GetComplianceStatus)
	}

	// 数据导出
	dataExport := rg.Group("/security/data-export")
	{
		dataExport.POST("", sec.CreateDataExport)
		dataExport.GET("/:id", sec.GetDataExport)
		dataExport.GET("/:id/download", sec.DownloadDataExport)
	}

	// GDPR请求
	gdpr := rg.Group("/security/gdpr")
	{
		gdpr.GET("/requests", sec.ListGDPRRequests)
		gdpr.POST("/requests", sec.CreateGDPRRequest)
		gdpr.GET("/requests/:id", sec.GetGDPRRequest)
		gdpr.PUT("/requests/:id/process", sec.ProcessGDPRRequest)
	}

	// 同意记录
	consent := rg.Group("/security/consent")
	{
		consent.GET("/records", sec.ListConsentRecords)
		consent.POST("/records", sec.RecordConsent)
		consent.GET("/records/:id", sec.GetConsentRecord)
		consent.PUT("/records/:id/withdraw", sec.WithdrawConsent)
	}
}

// ============ 安全审计 API ============

// GetAuditLogs 获取审计日志
func (ctrl *SecurityEvoController) GetAuditLogs(ctx *gin.Context) {
	var logs []models.AuditLog
	query := ctrl.DB.Model(&models.AuditLog{})

	// 租户过滤
	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 过滤条件
	if action := ctx.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if module := ctx.Query("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	if resourceType := ctx.Query("resource_type"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)
	if err := query.Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询审计日志失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  logs,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// ListAuditReports 获取审计报告列表
func (ctrl *SecurityEvoController) ListAuditReports(ctx *gin.Context) {
	var reports []models.AuditReport
	query := ctrl.DB.Model(&models.AuditReport{})

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if reportType := ctx.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("period_start >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("period_end <= ?", endDate+" 23:59:59")
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)
	if err := query.Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询审计报告失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  reports,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreateAuditReportRequest 生成审计报告请求
type CreateAuditReportRequest struct {
	ReportName string `json:"report_name" binding:"required"`
	ReportType string `json:"report_type" binding:"required"` // security, compliance, access, operations
	PeriodStart string `json:"period_start" binding:"required"`
	PeriodEnd   string `json:"period_end" binding:"required"`
	Format     string `json:"format"` // pdf, xlsx, csv, json
}

// CreateAuditReport 生成审计报告
func (ctrl *SecurityEvoController) CreateAuditReport(ctx *gin.Context) {
	var req CreateAuditReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "period_start 格式错误，应为 YYYY-MM-DD"})
		return
	}
	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "period_end 格式错误，应为 YYYY-MM-DD"})
		return
	}

	format := req.Format
	if format == "" {
		format = "pdf"
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	report := models.AuditReport{
		ReportName:  req.ReportName,
		ReportType:  req.ReportType,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Format:      format,
		Status:      1, // 生成中
		GeneratedBy: uid,
		TenantID:    tid,
	}

	if err := ctrl.DB.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建审计报告失败"})
		return
	}

	// 生成报告摘要
	summary := fmt.Sprintf("审计报告 %s，周期 %s 至 %s，类型: %s",
		req.ReportName, req.PeriodStart, req.PeriodEnd, req.ReportType)

	var logCount int64
	ctrl.DB.Model(&models.AuditLog{}).Where("created_at BETWEEN ? AND ?", periodStart, periodEnd.Add(24*time.Hour)).Count(&logCount)

	var violationCount int64
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("created_at BETWEEN ? AND ?", periodStart, periodEnd.Add(24*time.Hour)).Count(&violationCount)

	detail := fmt.Sprintf(`{"log_count": %d, "violation_count": %d, "report_type": "%s"}`, logCount, violationCount, req.ReportType)

	ctrl.DB.Model(&report).Updates(map[string]interface{}{
		"status":  2, // 已完成
		"summary": summary,
		"detail":  detail,
	})

	report.Status = 2
	report.Summary = summary
	report.Detail = detail

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": report,
	})
}

// GetAuditReport 获取审计报告详情
func (ctrl *SecurityEvoController) GetAuditReport(ctx *gin.Context) {
	id := ctx.Param("id")

	var report models.AuditReport
	query := ctrl.DB.Where("id = ?", id)

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询报告失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": report,
	})
}

// ============ 合规管理 API ============

// ListComplianceReports 获取合规报告列表
func (ctrl *SecurityEvoController) ListComplianceReports(ctx *gin.Context) {
	var reports []models.ComplianceReport
	query := ctrl.DB.Model(&models.ComplianceReport{})

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if regulationType := ctx.Query("regulation_type"); regulationType != "" {
		query = query.Where("regulation_type = ?", regulationType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)
	if err := query.Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询合规报告失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  reports,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreateComplianceReportRequest 生成合规报告请求
type CreateComplianceReportRequest struct {
	ReportName     string  `json:"report_name" binding:"required"`
	RegulationType string  `json:"regulation_type" binding:"required"` // gdpr, ccpa, hipaa, sox, iso27001
	Scope          string  `json:"scope"`
	PeriodStart    string  `json:"period_start" binding:"required"`
	PeriodEnd      string  `json:"period_end" binding:"required"`
	Findings       string  `json:"findings"`
	Recommendations string `json:"recommendations"`
}

// CreateComplianceReport 生成合规报告
func (ctrl *SecurityEvoController) CreateComplianceReport(ctx *gin.Context) {
	var req CreateComplianceReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "period_start 格式错误，应为 YYYY-MM-DD"})
		return
	}
	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "period_end 格式错误，应为 YYYY-MM-DD"})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 计算合规评分
	var totalViolations int64
	ctrl.DB.Model(&models.ComplianceViolation{}).
		Where("created_at BETWEEN ? AND ?", periodStart, periodEnd.Add(24*time.Hour)).
		Count(&totalViolations)

	var totalPolicies int64
	ctrl.DB.Model(&models.CompliancePolicy{}).Where("enabled = ?", true).Count(&totalPolicies)

	score := 100.0
	if totalPolicies > 0 {
		score = float64(totalPolicies-totalViolations) / float64(totalPolicies) * 100
		if score < 0 {
			score = 0
		}
	}

	report := models.ComplianceReport{
		ReportName:      req.ReportName,
		RegulationType:  req.RegulationType,
		Scope:           req.Scope,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		Status:          1, // 草稿
		Score:           score,
		ViolationsCount: int(totalViolations),
		Findings:        req.Findings,
		Recommendations: req.Recommendations,
		TenantID:        tid,
	}

	if err := ctrl.DB.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建合规报告失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": report,
	})
}

// GetComplianceStatus 获取合规状态
func (ctrl *SecurityEvoController) GetComplianceStatus(ctx *gin.Context) {
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	var totalPolicies int64
	ctrl.DB.Model(&models.CompliancePolicy{}).Where("enabled = ? AND tenant_id = ?", true, tid).Count(&totalPolicies)

	var openViolations int64
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("status IN ? AND tenant_id = ?", []int{1, 2}, tid).Count(&openViolations)

	var resolvedViolations int64
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("status = ? AND tenant_id = ?", 3, tid).Count(&resolvedViolations)

	var totalViolations int64
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("tenant_id = ?", tid).Count(&totalViolations)

	score := 100.0
	if totalPolicies > 0 {
		score = float64(totalPolicies-openViolations) / float64(totalPolicies) * 100
		if score < 0 {
			score = 0
		}
	}

	// 按严重程度统计
	var criticalCount int64
	var highCount int64
	var mediumCount int64
	var lowCount int64

	ctrl.DB.Model(&models.ComplianceViolation{}).
		Select("severity, count(*) as count").
		Where("tenant_id = ? AND status IN ?", tid, []int{1, 2}).
		Group("severity").
		Find(&[]struct {
			Severity int
			Count    int64
		}{})

	// 单独查询各级别
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("tenant_id = ? AND status IN ? AND severity = ?", tid, []int{1, 2}, 4).Count(&criticalCount)
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("tenant_id = ? AND status IN ? AND severity = ?", tid, []int{1, 2}, 3).Count(&highCount)
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("tenant_id = ? AND status IN ? AND severity = ?", tid, []int{1, 2}, 2).Count(&mediumCount)
	ctrl.DB.Model(&models.ComplianceViolation{}).Where("tenant_id = ? AND status IN ? AND severity = ?", tid, []int{1, 2}, 1).Count(&lowCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"overall_score":      score,
			"total_policies":     totalPolicies,
			"open_violations":    openViolations,
			"resolved_violations": resolvedViolations,
			"violations_by_severity": gin.H{
				"critical": criticalCount,
				"high":     highCount,
				"medium":   mediumCount,
				"low":      lowCount,
			},
		},
	})
}

// ============ 数据导出 API ============

// CreateDataExportRequest 创建导出任务请求
type CreateDataExportRequest struct {
	ExportType string   `json:"export_type" binding:"required"` // full, partial, audit, compliance
	DataTypes  string   `json:"data_types"`                      // devices, members, activity, all
	Format     string   `json:"format"`                          // json, csv, xlsx, zip
	Filters    []string `json:"filters"`                        // 过滤条件
}

// CreateDataExport 创建数据导出任务
func (ctrl *SecurityEvoController) CreateDataExport(ctx *gin.Context) {
	var req CreateDataExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	format := req.Format
	if format == "" {
		format = "json"
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	exportID := fmt.Sprintf("exp_%d_%d", time.Now().UnixNano(), uid)

	dataExport := models.DataExport{
		ExportID:    exportID,
		ExportType: req.ExportType,
		DataTypes:  req.DataTypes,
		Format:     format,
		Status:     2, // 处理中
		RequesterID: uid,
		TenantID:   tid,
	}

	if req.Filters != nil {
		filterJSON, _ := strings.Join(req.Filters, ","), ""
		dataExport.Filters = fmt.Sprintf(`{"filters": ["%s"]}`, strings.Join(req.Filters, `","`))
		_ = filterJSON
	}

	if err := ctrl.DB.Create(&dataExport).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建导出任务失败"})
		return
	}

	// 模拟导出完成（实际应使用异步任务队列）
	exportDir := "./exports"
	os.MkdirAll(exportDir, 0755)
	fileName := fmt.Sprintf("%s_%s.%s", exportID, req.DataTypes, format)
	filePath := filepath.Join(exportDir, fileName)

	// 创建空文件作为占位
	expiredAt := time.Now().Add(7 * 24 * time.Hour)
	ctrl.DB.Model(&dataExport).Updates(map[string]interface{}{
		"status":    3, // 完成
		"file_path": filePath,
		"file_size": 0,
		"expires_at": expiredAt,
	})

	dataExport.Status = 3
	dataExport.FilePath = filePath
	dataExport.ExpiresAt = &expiredAt

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": dataExport,
	})
}

// GetDataExport 获取导出任务详情
func (ctrl *SecurityEvoController) GetDataExport(ctx *gin.Context) {
	id := ctx.Param("id")

	var exp models.DataExport
	query := ctrl.DB.Where("export_id = ?", id)

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&exp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "导出任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询导出任务失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": exp,
	})
}

// DownloadDataExport 下载导出文件
func (ctrl *SecurityEvoController) DownloadDataExport(ctx *gin.Context) {
	id := ctx.Param("id")

	var exp models.DataExport
	query := ctrl.DB.Where("export_id = ?", id)

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&exp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "导出任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询导出任务失败"})
		return
	}

	if exp.Status != 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "导出任务未完成，无法下载"})
		return
	}

	if exp.ExpiresAt != nil && time.Now().After(*exp.ExpiresAt) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "导出文件已过期"})
		return
	}

	// 增加下载次数
	ctrl.DB.Model(&exp).Update("download_count", exp.DownloadCount+1)

	// 检查文件是否存在
	if _, err := os.Stat(exp.FilePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "导出文件不存在"})
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(exp.FilePath)))
	ctx.File(exp.FilePath)
}

// ============ GDPR请求 API ============

// ListGDPRRequests 获取GDPR请求列表
func (ctrl *SecurityEvoController) ListGDPRRequests(ctx *gin.Context) {
	var requests []models.GDPRRequest
	query := ctrl.DB.Model(&models.GDPRRequest{})

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if requestType := ctx.Query("request_type"); requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if email := ctx.Query("email"); email != "" {
		query = query.Where("requester_email LIKE ?", "%"+email+"%")
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)
	if err := query.Find(&requests).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询GDPR请求失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  requests,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreateGDPRRequestInput 创建GDPR请求输入
type CreateGDPRRequestInput struct {
	RequestType   string `json:"request_type" binding:"required"` // data_access, data_deletion, data_portability
	RequesterEmail string `json:"requester_email" binding:"required"`
	RequesterName string `json:"requester_name"`
	RequestReason string `json:"request_reason"`
}

// CreateGDPRRequest 创建GDPR请求
func (ctrl *SecurityEvoController) CreateGDPRRequest(ctx *gin.Context) {
	var input CreateGDPRRequestInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	requestID := fmt.Sprintf("gdpr_%d_%s", time.Now().UnixNano(), strings.ToLower(input.RequestType))

	req := models.GDPRRequest{
		RequestID:      requestID,
		RequestType:    input.RequestType,
		RequesterEmail: input.RequesterEmail,
		RequesterName:  input.RequesterName,
		RequestReason:  input.RequestReason,
		Status:         1, // 待处理
		TenantID:       tid,
	}

	if err := ctrl.DB.Create(&req).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建GDPR请求失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": req,
	})
}

// GetGDPRRequest 获取GDPR请求详情
func (ctrl *SecurityEvoController) GetGDPRRequest(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.GDPRRequest
	query := ctrl.DB.Where("id = ?", id)

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&req).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "GDPR请求不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询GDPR请求失败"})
		return
	}

	// 查询扩展信息
	var extra models.GDPRRequestExtra
	ctrl.DB.Where("gdpr_request_id = ?", req.ID).First(&extra)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"request": req,
			"extra":   extra,
		},
	})
}

// ProcessGDPRRequestInput 处理GDPR请求输入
type ProcessGDPRRequestInput struct {
	Action         string `json:"action" binding:"required"` // approve, reject
	ProcessedNotes  string `json:"processed_notes"`
	RejectedReason string `json:"rejected_reason"`
}

// ProcessGDPRRequest 处理GDPR请求
func (ctrl *SecurityEvoController) ProcessGDPRRequest(ctx *gin.Context) {
	id := ctx.Param("id")

	var req models.GDPRRequest
	if err := ctrl.DB.Where("id = ?", id).First(&req).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "GDPR请求不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询GDPR请求失败"})
		return
	}

	if req.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该请求已被处理，无法重复处理"})
		return
	}

	var input ProcessGDPRRequestInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	now := time.Now()

	if input.Action == "approve" {
		ctrl.DB.Model(&req).Updates(map[string]interface{}{
			"status":       2, // 处理中
			"processed_by": uid,
			"processed_at": now,
		})

		// 模拟处理完成
		ctrl.DB.Model(&req).Updates(map[string]interface{}{
			"status":       3, // 已完成
			"completed_at": now,
		})

		// 创建扩展记录
		extra := models.GDPRRequestExtra{
			GDPRRequestID:   req.ID,
			ProcessingNotes: input.ProcessedNotes,
			CompletionNotes: "请求已处理完成",
		}
		ctrl.DB.Create(&extra)

	} else if input.Action == "reject" {
		ctrl.DB.Model(&req).Updates(map[string]interface{}{
			"status":          4, // 已拒绝
			"processed_by":    uid,
			"processed_at":    now,
			"rejected_reason": input.RejectedReason,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "action 必须是 approve 或 reject"})
		return
	}

	ctrl.DB.Where("id = ?", id).First(&req)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": req,
	})
}

// ============ 同意记录 API ============

// ListConsentRecords 获取同意记录列表
func (ctrl *SecurityEvoController) ListConsentRecords(ctx *gin.Context) {
	var records []models.ConsentRecord
	query := ctrl.DB.Model(&models.ConsentRecord{})

	if tenantID, exists := ctx.Get("tenant_id"); exists {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if consentType := ctx.Query("consent_type"); consentType != "" {
		query = query.Where("consent_type = ?", consentType)
	}
	if consentAction := ctx.Query("consent_action"); consentAction != "" {
		query = query.Where("consent_action = ?", consentAction)
	}
	if userEmail := ctx.Query("user_email"); userEmail != "" {
		query = query.Where("user_email LIKE ?", "%"+userEmail+"%")
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)
	if err := query.Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询同意记录失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  records,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// RecordConsentInput 记录同意输入
type RecordConsentInput struct {
	ConsentType  string `json:"consent_type" binding:"required"` // privacy, terms, marketing, data_processing, third_party_sharing
	ConsentAction string `json:"consent_action" binding:"required"` // granted, withdrawn, updated
	UserID      uint   `json:"user_id"`
	UserEmail   string `json:"user_email" binding:"required"`
	DeviceID    string `json:"device_id"`
	Version     string `json:"version"`
	PolicyURL   string `json:"policy_url"`
	ConsentProof string `json:"consent_proof"`
}

// RecordConsent 记录同意
func (ctrl *SecurityEvoController) RecordConsent(ctx *gin.Context) {
	var input RecordConsentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	record := models.ConsentRecord{
		ConsentType:   input.ConsentType,
		ConsentAction: input.ConsentAction,
		UserID:        input.UserID,
		UserEmail:     input.UserEmail,
		DeviceID:      input.DeviceID,
		Version:       input.Version,
		PolicyURL:     input.PolicyURL,
		ConsentProof:  input.ConsentProof,
		IPAddress:     ctx.ClientIP(),
		UserAgent:     ctx.Request.UserAgent(),
		TenantID:      tid,
	}

	if input.ConsentAction == "withdrawn" {
		now := time.Now()
		record.WithdrawnAt = &now
	}

	if err := ctrl.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "记录同意失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": record,
	})
}

// GetConsentRecord 获取同意记录详情
func (ctrl *SecurityEvoController) GetConsentRecord(ctx *gin.Context) {
	id := ctx.Param("id")

	var record models.ConsentRecord
	if err := ctrl.DB.Where("id = ?", id).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "同意记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询同意记录失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": record,
	})
}

// WithdrawConsent 撤回同意
func (ctrl *SecurityEvoController) WithdrawConsent(ctx *gin.Context) {
	id := ctx.Param("id")

	var record models.ConsentRecord
	if err := ctrl.DB.Where("id = ?", id).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "同意记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询同意记录失败"})
		return
	}

	if record.ConsentAction == "withdrawn" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该同意记录已处于撤回状态"})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&record).Updates(map[string]interface{}{
		"consent_action": "withdrawn",
		"withdrawn_at":   now,
		"updated_at":     now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "撤回同意失败"})
		return
	}

	record.ConsentAction = "withdrawn"
	record.WithdrawnAt = &now

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": record,
	})
}
	