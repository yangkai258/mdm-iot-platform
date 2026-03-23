package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// APIQuotaController API配额
type APIQuotaController struct{}

func NewAPIQuotaController() *APIQuotaController {
	return &APIQuotaController{}
}

func (ctrl *APIQuotaController) RegisterRoutes(rg *gin.RouterGroup) {
	quota := rg.Group("/quota")
	{
		quota.GET("/status", ctrl.GetQuotaStatus)
		quota.GET("/usage", ctrl.GetUsageLog)
		quota.GET("/plans", ctrl.GetQuotaPlans)
	}
}

// GetQuotaStatus 获取配额状态
func (ctrl *APIQuotaController) GetQuotaStatus(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id required"})
		return
	}

	// 模拟返回配额数据（实际需要查询数据库）
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"plan_type":      "free",
		"monthly_quota":  1000,
		"used_quota":     150,
		"remaining":      850,
		"reset_at":       time.Now().AddDate(0, 1, 0),
		"usage_percent":   15.0,
	}})
}

// GetUsageLog 获取使用日志
func (ctrl *APIQuotaController) GetUsageLog(c *gin.Context) {
	userID := c.Query("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	// 模拟返回使用日志
	logs := []gin.H{
		{"id": 1, "endpoint": "/api/v1/devices", "method": "GET", "status_code": 200, "latency": 45, "created_at": time.Now()},
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"list": logs, "total": 1, "page": page, "page_size": pageSize,
	}})
}

// GetQuotaPlans 获取配额套餐
func (ctrl *APIQuotaController) GetQuotaPlans(c *gin.Context) {
	plans := []map[string]interface{}{
		{"plan": "free", "monthly_quota": 1000, "price": 0, "name": "免费版"},
		{"plan": "basic", "monthly_quota": 10000, "price": 99, "name": "基础版"},
		{"plan": "pro", "monthly_quota": 100000, "price": 399, "name": "专业版"},
		{"plan": "enterprise", "monthly_quota": 1000000, "price": 999, "name": "企业版"},
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": plans})
}
