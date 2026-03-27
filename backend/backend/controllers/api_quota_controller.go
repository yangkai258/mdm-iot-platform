package controllers

import (
	"net/http"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIQuotaController API配额控制器
type APIQuotaController struct {
	DB *gorm.DB
}

// GetQuota 获取当前租户的API配额
// GET /api/v1/quotas
func (c *APIQuotaController) GetQuota(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	quota, err := models.GetOrCreateAPIQuota(c.DB, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 2, "message": "获取配额失败: " + err.Error()})
		return
	}

	// 判断是否已超限
	exceeded := quota.MonthlyCalls != -1 && quota.UsedCalls >= quota.MonthlyCalls

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"tenant_id":     quota.TenantID,
			"plan":          quota.Plan,
			"monthly_calls": quota.MonthlyCalls,
			"used_calls":    quota.UsedCalls,
			"remaining":     quota.MonthlyCalls - quota.UsedCalls,
			"reset_date":    quota.ResetDate.Format(time.RFC3339),
			"exceeded":      exceeded,
			"unlimited":     quota.MonthlyCalls == -1,
		},
	})
}

// GetUsage 获取API使用量（最近30天调用统计）
// GET /api/v1/quotas/usage
func (c *APIQuotaController) GetUsage(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	// 最近30天调用统计
	since := time.Now().AddDate(0, 0, -30)

	var totalCalls int64
	c.DB.Model(&models.APIUsageLog{}).
		Where("tenant_id = ? AND call_at >= ?", tenantID, since).
		Count(&totalCalls)

	// 按天统计
	type DailyUsage struct {
		Date  string `json:"date"`
		Calls int64  `json:"calls"`
	}
	var dailyStats []DailyUsage
	c.DB.Model(&models.APIUsageLog{}).
		Select("DATE(call_at) as date, COUNT(*) as calls").
		Where("tenant_id = ? AND call_at >= ?", tenantID, since).
		Group("DATE(call_at)").
		Order("date DESC").
		Scan(&dailyStats)

	// 按接口路径统计Top 10
	type TopEndpoint struct {
		Path  string `json:"path"`
		Calls int64  `json:"calls"`
	}
	var topEndpoints []TopEndpoint
	c.DB.Model(&models.APIUsageLog{}).
		Select("path, COUNT(*) as calls").
		Where("tenant_id = ? AND call_at >= ?", tenantID, since).
		Group("path").
		Order("calls DESC").
		Limit(10).
		Scan(&topEndpoints)

	// 错误率统计
	var errorCalls int64
	c.DB.Model(&models.APIUsageLog{}).
		Where("tenant_id = ? AND call_at >= ? AND status_code >= 400", tenantID, since).
		Count(&errorCalls)

	errorRate := float64(0)
	if totalCalls > 0 {
		errorRate = float64(errorCalls) / float64(totalCalls) * 100
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period":        "30d",
			"total_calls":   totalCalls,
			"daily_stats":   dailyStats,
			"top_endpoints": topEndpoints,
			"error_rate":    errorRate,
			"error_calls":   errorCalls,
		},
	})
}

// UpgradeQuotaRequest 升级配额请求
type UpgradeQuotaRequest struct {
	Plan string `json:"plan" binding:"required"` // basic/pro/enterprise
}

// UpgradeQuota 升级API配额（实际业务中会调用支付）
// POST /api/v1/quotas/upgrade
func (c *APIQuotaController) UpgradeQuota(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	var req UpgradeQuotaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 2, "message": "参数错误: " + err.Error()})
		return
	}

	validPlans := map[string]bool{"free": true, "basic": true, "pro": true, "enterprise": true}
	if !validPlans[req.Plan] {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 3, "message": "无效的套餐类型"})
		return
	}

	// 更新租户套餐
	if err := c.DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Update("plan", req.Plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": "更新套餐失败: " + err.Error()})
		return
	}

	// 重新获取配额（会使用新套餐的配额）
	quota, err := models.GetOrCreateAPIQuota(c.DB, tenantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5, "message": "获取配额失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "套餐升级成功",
		"data": gin.H{
			"tenant_id":     quota.TenantID,
			"plan":          quota.Plan,
			"monthly_calls": quota.MonthlyCalls,
			"used_calls":    quota.UsedCalls,
			"remaining":     quota.MonthlyCalls - quota.UsedCalls,
			"reset_date":    quota.ResetDate.Format(time.RFC3339),
			"unlimited":     quota.MonthlyCalls == -1,
		},
	})
}
