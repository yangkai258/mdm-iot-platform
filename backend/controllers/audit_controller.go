package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditController 审计日志控制器
type AuditController struct {
	DB *gorm.DB
}

// GetAuditLogs 获取审计日志列表
func (c *AuditController) GetAuditLogs(ctx *gin.Context) {
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// 查询参数
	action := ctx.Query("action")
	module := ctx.Query("module")
	resourceType := ctx.Query("resource_type")
	userID := ctx.Query("user_id")
	status := ctx.Query("status")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	keyword := ctx.Query("keyword")

	// 构建查询
	query := c.DB.Model(&models.AuditLog{})

	// 租户过滤
	if tid > 0 {
		query = query.Where("tenant_id = ?", tid)
	}

	// 过滤条件
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour))
		}
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR error_msg LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	var logs []models.AuditLog
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)

	// 获取动作类型统计
	actionStats := c.getActionStats(tid)

	// 获取模块统计
	moduleStats := c.getModuleStats(tid)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":         logs,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"action_stats": actionStats,
			"module_stats": moduleStats,
		},
	})
}

// GetAuditLog 获取审计日志详情
func (c *AuditController) GetAuditLog(ctx *gin.Context) {
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	id := ctx.Param("id")

	var log models.AuditLog
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tid).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "审计日志不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	// 获取相关资源详情（如果有）
	resourceDetail := c.getResourceDetail(log.ResourceType, log.ResourceID)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"audit_log":       log,
			"resource_detail": resourceDetail,
		},
	})
}

// GetAuditStatistics 获取审计统计
func (c *AuditController) GetAuditStatistics(ctx *gin.Context) {
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 今日统计
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	c.DB.Model(&models.AuditLog{}).Where("tenant_id = ? AND DATE(created_at) = ?", tid, today).Count(&todayCount)

	// 本周统计
	weekStart := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var weekCount int64
	c.DB.Model(&models.AuditLog{}).Where("tenant_id = ? AND created_at >= ?", tid, weekStart).Count(&weekCount)

	// 失败操作统计
	var failedCount int64
	c.DB.Model(&models.AuditLog{}).Where("tenant_id = ? AND status = 2", tid).Count(&failedCount)

	// 按小时统计（今日）
	hourlyStats := c.getHourlyStats(tid)

	// 按模块统计
	moduleStats := c.getModuleStats(tid)

	// 按操作类型统计
	actionStats := c.getActionStats(tid)

	// 最近的失败操作
	var recentFailures []models.AuditLog
	c.DB.Where("tenant_id = ? AND status = 2", tid).Order("created_at DESC").Limit(10).Find(&recentFailures)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"today_count":      todayCount,
			"week_count":       weekCount,
			"failed_count":     failedCount,
			"hourly_stats":     hourlyStats,
			"module_stats":     moduleStats,
			"action_stats":     actionStats,
			"recent_failures":  recentFailures,
		},
	})
}

// ExportAuditLogs 导出审计日志
func (c *AuditController) ExportAuditLogs(ctx *gin.Context) {
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 查询参数（与 GetAuditLogs 相同）
	action := ctx.Query("action")
	module := ctx.Query("module")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	query := c.DB.Model(&models.AuditLog{}).Where("tenant_id = ?", tid)

	if action != "" {
		query = query.Where("action = ?", action)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour))
		}
	}

	var logs []models.AuditLog
	query.Order("created_at DESC").Limit(10000).Find(&logs)

	// 记录导出操作
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	username, _ := ctx.Get("username")

	exportLog := models.AuditLog{
		Action:        "export",
		Module:        "audit",
		ResourceType:  "audit_log",
		UserID:        uid,
		Username:      username.(string),
		IP:            ctx.ClientIP(),
		UserAgent:     ctx.GetHeader("User-Agent"),
		Status:        1,
		RequestMethod: ctx.Request.Method,
		RequestPath:   ctx.Request.URL.Path,
		TenantID:      tid,
	}
	c.DB.Create(&exportLog)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"export_count": len(logs),
			"logs":         logs,
			"exported_at":  time.Now(),
		},
	})
}

// ============ 辅助函数 ============

func (c *AuditController) getActionStats(tenantID uint) []map[string]interface{} {
	type Result struct {
		Action string
		Count  int64
	}
	var results []Result
	c.DB.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("action").
		Order("count DESC").
		Limit(10).
		Scan(&results)

	stats := make([]map[string]interface{}, len(results))
	for i, r := range results {
		stats[i] = map[string]interface{}{
			"action": r.Action,
			"count":  r.Count,
		}
	}
	return stats
}

func (c *AuditController) getModuleStats(tenantID uint) []map[string]interface{} {
	type Result struct {
		Module string
		Count  int64
	}
	var results []Result
	c.DB.Model(&models.AuditLog{}).
		Select("module, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("module").
		Order("count DESC").
		Scan(&results)

	stats := make([]map[string]interface{}, len(results))
	for i, r := range results {
		stats[i] = map[string]interface{}{
			"module": r.Module,
			"count":  r.Count,
		}
	}
	return stats
}

func (c *AuditController) getHourlyStats(tenantID uint) []map[string]interface{} {
	today := time.Now().Format("2006-01-02")
	type Result struct {
		Hour  int
		Count int64
	}
	var results []Result
	c.DB.Model(&models.AuditLog{}).
		Select("EXTRACT(HOUR FROM created_at) as hour, COUNT(*) as count").
		Where("tenant_id = ? AND DATE(created_at) = ?", tenantID, today).
		Group("hour").
		Order("hour").
		Scan(&results)

	stats := make([]map[string]interface{}, 24)
	for i := 0; i < 24; i++ {
		stats[i] = map[string]interface{}{
			"hour":  i,
			"count": 0,
		}
	}
	for _, r := range results {
		stats[r.Hour] = map[string]interface{}{
			"hour":  r.Hour,
			"count": r.Count,
		}
	}
	return stats
}

func (c *AuditController) getResourceDetail(resourceType, resourceID string) interface{} {
	if resourceID == "" {
		return nil
	}

	switch resourceType {
	case "device":
		var device models.Device
		if err := c.DB.Where("device_id = ?", resourceID).First(&device).Error; err == nil {
			return device
		}
	case "member":
		var member models.Member
		if err := c.DB.Where("id = ?", resourceID).First(&member).Error; err == nil {
			return member
		}
	case "sys_user":
		var user models.SysUser
		if err := c.DB.Where("id = ?", resourceID).First(&user).Error; err == nil {
			return map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.Nickname,
				"email":    user.Email,
			}
		}
	}

	return nil
}
