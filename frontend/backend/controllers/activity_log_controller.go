package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActivityLogController 活动日志控制器
type ActivityLogController struct {
	DB *gorm.DB
}

// List 获取活动日志列表
// GET /api/v1/activity-logs
func (c *ActivityLogController) List(ctx *gin.Context) {
	var logs []models.ActivityLog
	var total int64

	query := c.DB.Model(&models.ActivityLog{})

	// 关键词搜索（用户名/资源名/IP）
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR resource_name LIKE ? OR ip LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 操作类型筛选
	if action := ctx.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}

	// 资源类型筛选
	if resourceType := ctx.Query("resource_type"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	// 用户ID筛选
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}

	// 时间范围筛选
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Select("*").
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 转换 Details []byte 为 map 以便 JSON 序列化
	logMaps := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		logMaps[i] = map[string]interface{}{
			"id":            log.ID,
			"user_id":       log.UserID,
			"username":      log.Username,
			"action":        log.Action,
			"resource_type": log.ResourceType,
			"resource_id":   log.ResourceID,
			"resource_name": log.ResourceName,
			"details":       log.GetDetails(),
			"ip":            log.IP,
			"user_agent":    log.UserAgent,
			"tenant_id":     log.TenantID,
			"created_at":    log.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      logMaps,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取单条活动日志
// GET /api/v1/activity-logs/:id
func (c *ActivityLogController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var log models.ActivityLog
	if err := c.DB.First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "日志不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": map[string]interface{}{
		"id":            log.ID,
		"user_id":       log.UserID,
		"username":      log.Username,
		"action":        log.Action,
		"resource_type": log.ResourceType,
		"resource_id":   log.ResourceID,
		"resource_name": log.ResourceName,
		"details":       log.GetDetails(),
		"ip":            log.IP,
		"user_agent":    log.UserAgent,
		"tenant_id":     log.TenantID,
		"created_at":    log.CreatedAt,
	}})
}

// GetStatistics 获取活动统计
// GET /api/v1/activity-logs/statistics
func (c *ActivityLogController) GetStatistics(ctx *gin.Context) {
	var stats struct {
		TodayCount  int64 `json:"today_count"`
		WeekCount   int64 `json:"week_count"`
		MonthCount  int64 `json:"month_count"`
		LoginCount  int64 `json:"login_count"`
		CreateCount int64 `json:"create_count"`
		UpdateCount int64 `json:"update_count"`
		DeleteCount int64 `json:"delete_count"`
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	c.DB.Model(&models.ActivityLog{}).Where("created_at >= ?", today).Count(&stats.TodayCount)
	c.DB.Model(&models.ActivityLog{}).Where("created_at >= ?", weekAgo).Count(&stats.WeekCount)
	c.DB.Model(&models.ActivityLog{}).Where("created_at >= ?", monthAgo).Count(&stats.MonthCount)
	c.DB.Model(&models.ActivityLog{}).Where("action = ?", "login").Count(&stats.LoginCount)
	c.DB.Model(&models.ActivityLog{}).Where("action = ?", "create").Count(&stats.CreateCount)
	c.DB.Model(&models.ActivityLog{}).Where("action = ?", "update").Count(&stats.UpdateCount)
	c.DB.Model(&models.ActivityLog{}).Where("action = ?", "delete").Count(&stats.DeleteCount)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": stats})
}

// RecordActivity 记录活动日志（供其他模块调用）
func RecordActivity(db *gorm.DB, userID uint, username, action, resourceType string, resourceID uint, resourceName string, details map[string]interface{}, ip string) {
	tenantID := ""
	if userID > 0 {
		var user models.SysUser
		if err := db.First(&user, userID).Error; err == nil {
			tenantID = user.TenantID
		}
	}

	log := models.ActivityLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		IP:           ip,
		TenantID:     tenantID,
	}
	log.SetDetails(details)
	db.Create(&log)
}

// ============ LoginLog 控制器 ============

// LoginLogController 登录日志控制器
type LoginLogController struct {
	DB *gorm.DB
}

// List 获取登录日志列表
// GET /api/v1/login-logs
func (c *LoginLogController) List(ctx *gin.Context) {
	var logs []models.LoginLog
	var total int64

	query := c.DB.Model(&models.LoginLog{})

	// 关键词搜索
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR ip LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 时间范围
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("login_time >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("login_time <= ?", t)
		}
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
