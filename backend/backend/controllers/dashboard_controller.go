package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardController Dashboard控制器
type DashboardController struct {
	DB *gorm.DB
}

// EnhancedDashboardStats Dashboard统计数据（增强版）
type EnhancedDashboardStats struct {
	TotalDevices   int64 `json:"total_devices"`
	OnlineDevices  int64 `json:"online_devices"`
	OfflineDevices int64 `json:"offline_devices"`
	TotalAlerts    int64 `json:"total_alerts"`
	PendingAlerts  int64 `json:"pending_alerts"`
	ResolvedAlerts int64 `json:"resolved_alerts"`
	GeofenceAlerts int64 `json:"geofence_alerts"`
	TodayLogins    int64 `json:"today_logins"`
	TotalMembers   int64 `json:"total_members"`
	ActiveMembers  int64 `json:"active_members"`
}

// TopUserStat Top用户统计
type TopUserStat struct {
	Username string `json:"username"`
	Count    int64  `json:"count"`
}

// TopActionStat Top操作统计
type TopActionStat struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// ActivitySummary 活动摘要
type ActivitySummary struct {
	TodayActions int64            `json:"today_actions"`
	WeekActions  int64           `json:"week_actions"`
	TopUsers     []TopUserStat   `json:"top_users"`
	TopActions   []TopActionStat `json:"top_actions"`
}

// GetStats 获取Dashboard统计
// GET /api/v1/dashboard/stats
func (c *DashboardController) GetStats(ctx *gin.Context) {
	var stats EnhancedDashboardStats

	// 设备统计
	c.DB.Model(&models.Device{}).Count(&stats.TotalDevices)
	c.DB.Model(&models.DeviceShadow{}).Select("COUNT(*)").Where("is_online = ?", true).Scan(&stats.OnlineDevices)
	stats.OfflineDevices = stats.TotalDevices - stats.OnlineDevices
	if stats.OfflineDevices < 0 {
		stats.OfflineDevices = 0
	}

	// 告警统计
	c.DB.Model(&models.DeviceAlert{}).Count(&stats.TotalAlerts)
	c.DB.Model(&models.DeviceAlert{}).Where("status = ?", 1).Count(&stats.PendingAlerts)
	c.DB.Model(&models.DeviceAlert{}).Where("status = ?", 3).Count(&stats.ResolvedAlerts)
	c.DB.Model(&models.DeviceAlert{}).Where("alert_type = ?", "geofence_violation").Count(&stats.GeofenceAlerts)

	// 今日登录次数（从 activity_logs 表）
	c.DB.Model(&models.ActivityLog{}).Where("action = ? AND created_at >= CURRENT_DATE", "login").Count(&stats.TodayLogins)

	// 会员统计
	c.DB.Model(&models.Member{}).Count(&stats.TotalMembers)
	c.DB.Model(&models.Member{}).Where("status = ?", 1).Count(&stats.ActiveMembers)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}

// GetStatsSimple 简单统计（仅设备）
// GET /api/v1/dashboard/stats/simple
func (c *DashboardController) GetStatsSimple(ctx *gin.Context) {
	var stats struct {
		TotalDevices  int64 `json:"total_devices"`
		OnlineDevices int64 `json:"online_devices"`
	}

	c.DB.Model(&models.Device{}).Count(&stats.TotalDevices)
	c.DB.Model(&models.DeviceShadow{}).Select("COUNT(*)").Where("is_online = ?", true).Scan(&stats.OnlineDevices)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}

// GetActivitySummary 获取活动摘要
// GET /api/v1/dashboard/activity-summary
func (c *DashboardController) GetActivitySummary(ctx *gin.Context) {
	summary := ActivitySummary{
		TopUsers:   make([]TopUserStat, 0),
		TopActions: make([]TopActionStat, 0),
	}

	// 今日活动数
	c.DB.Model(&models.ActivityLog{}).Where("created_at >= CURRENT_DATE").Count(&summary.TodayActions)

	// 本周活动数
	c.DB.Model(&models.ActivityLog{}).Where("created_at >= NOW() - INTERVAL '7 days'").Count(&summary.WeekActions)

	// Top 5 用户
	rows, err := c.DB.Raw(`
		SELECT username, COUNT(*) as cnt
		FROM activity_logs
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY username
		ORDER BY cnt DESC
		LIMIT 5
	`).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var u TopUserStat
			rows.Scan(&u.Username, &u.Count)
			summary.TopUsers = append(summary.TopUsers, u)
		}
	}

	// Top 5 操作类型
	rows2, err := c.DB.Raw(`
		SELECT action, COUNT(*) as cnt
		FROM activity_logs
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY action
		ORDER BY cnt DESC
		LIMIT 5
	`).Rows()
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var a TopActionStat
			rows2.Scan(&a.Action, &a.Count)
			summary.TopActions = append(summary.TopActions, a)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": summary})
}
