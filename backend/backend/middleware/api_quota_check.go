package middleware

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIQuotaCheck API配额检查中间件
// 每次 API 调用检查配额，配额耗尽时返回 429 Too Many Requests
func APIQuotaCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID == "" {
			// 无租户ID，跳过检查（内部接口）
			c.Next()
			return
		}

		quota, err := models.GetOrCreateAPIQuota(db, tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "DB_ERROR",
				"message": "配额查询失败",
			})
			c.Abort()
			return
		}

		// -1 表示不限
		if quota.MonthlyCalls == -1 {
			// 记录使用（异步更好，这里同步简化）
			models.IncrementAPICalls(db, tenantID)
			c.Next()
			return
		}

		// 检查是否已超限
		if quota.UsedCalls >= quota.MonthlyCalls {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "QUOTA_EXCEEDED",
				"message": "本月API调用配额已用尽，请升级套餐",
				"data": gin.H{
					"plan":          quota.Plan,
					"monthly_calls": quota.MonthlyCalls,
					"used_calls":    quota.UsedCalls,
					"reset_date":    quota.ResetDate.Format(time.RFC3339),
					"upgrade_url":   "/api/v1/quotas/upgrade",
				},
			})
			c.Abort()
			return
		}

		// 配额信息存入 context，供后续 handler 使用
		c.Set("api_quota", map[string]interface{}{
			"tenant_id":    tenantID,
			"plan":         quota.Plan,
			"monthly_calls": quota.MonthlyCalls,
			"used_calls":   quota.UsedCalls,
			"remaining":    quota.MonthlyCalls - quota.UsedCalls,
		})

		c.Next()
	}
}

// RecordAPIUsage 记录API调用（中间件栈末尾调用）
func RecordAPIUsage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// 请求结束后记录使用量
		tenantID := GetTenantID(c)
		if tenantID == "" {
			return
		}

		// 只记录 2xx/4xx 响应
		statusCode := c.Writer.Status()
		if statusCode >= 500 && statusCode < 600 {
			// 5xx 不计入配额（服务端问题）
			return
		}

		latencyMs := time.Since(start).Milliseconds()

		// 异步记录（不阻塞响应）
		go func() {
			models.IncrementAPICalls(db, tenantID)

			// 记录调用明细（采样，避免存储爆炸）
			log := models.APIUsageLog{
				TenantID:   tenantID,
				Path:       c.Request.URL.Path,
				Method:     c.Request.Method,
				StatusCode: statusCode,
				LatencyMs:  latencyMs,
				IP:         c.ClientIP(),
				UserAgent:  c.Request.UserAgent(),
				CallAt:     time.Now(),
			}
			db.Create(&log)
		}()
	}
}
