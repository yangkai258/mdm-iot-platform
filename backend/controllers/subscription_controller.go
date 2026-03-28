package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionController struct {
	DB                *gorm.DB
	RenewalService    *services.SubscriptionRenewalService
}

func NewSubscriptionController(db *gorm.DB) *SubscriptionController {
	return &SubscriptionController{
		DB:             db,
		RenewalService: services.NewSubscriptionRenewalService(db),
	}
}

func (ctrl *SubscriptionController) RegisterRoutes(rg *gin.RouterGroup) {
	sub := rg.Group("/subscriptions")
	{
		sub.GET("", ctrl.ListSubscriptions)
		sub.POST("", ctrl.CreateSubscription)
		sub.GET("/:id", ctrl.GetSubscription)
		sub.PUT("/:id", ctrl.UpdateSubscription)
		sub.DELETE("/:id", ctrl.CancelSubscription)
		sub.POST("/webhook/payment", ctrl.PaymentWebhook)
	}
}

// ListSubscriptions 订阅列表
func (ctrl *SubscriptionController) ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	var subs []models.Subscription
	query := ctrl.DB.Model(&models.Subscription{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	query.Order("created_at DESC").Find(&subs)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": subs})
}

// CreateSubscription 创建订阅
func (ctrl *SubscriptionController) CreateSubscription(c *gin.Context) {
	var req struct {
		UserID    uint    `json:"user_id" binding:"required"`
		PlanName  string  `json:"plan_name" binding:"required"`
		PlanType  string  `json:"plan_type" binding:"required"`
		Price     float64 `json:"price" binding:"required"`
		Duration  int     `json:"duration"` // 天数
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Duration == 0 {
		req.Duration = 30
	}

	now := time.Now()
	sub := models.Subscription{
		UserID:    req.UserID,
		PlanName:  req.PlanName,
		PlanType:  req.PlanType,
		Price:     req.Price,
		Duration:  req.Duration,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, req.Duration),
		Status:    "active",
		AutoRenew: true,
	}
	ctrl.DB.Create(&sub)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": sub})
}

// GetSubscription 获取订阅详情
func (ctrl *SubscriptionController) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": sub})
}

// UpdateSubscription 更新订阅
func (ctrl *SubscriptionController) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	ctrl.DB.Model(&sub).Updates(req)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": sub})
}

// CancelSubscription 取消订阅
func (ctrl *SubscriptionController) CancelSubscription(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.Subscription{}).Where("id = ?", id).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订阅已取消"})
}

// AutoRenew 续费
func (ctrl *SubscriptionController) AutoRenew(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}

	newEndDate := sub.EndDate.AddDate(0, 0, sub.Duration)
	now := time.Now()
	ctrl.DB.Model(&sub).Updates(map[string]interface{}{
		"end_date":       newEndDate,
		"status":        "active",
		"auto_renew":     true,
		"last_renew_at": &now,
		"renew_count":   sub.RenewCount + 1,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "续费成功", "data": gin.H{
		"subscription_id": sub.ID,
		"new_end_date":   newEndDate,
	}})
}

// GetRenewalStatus 续费状态
func (ctrl *SubscriptionController) GetRenewalStatus(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}

	daysUntilExpiry := int(time.Until(sub.EndDate).Hours() / 24)
	renewalStatus := "active"
	if daysUntilExpiry <= 3 {
		renewalStatus = "expiring_soon"
	}
	if daysUntilExpiry <= 0 {
		renewalStatus = "expired"
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"subscription_id":    sub.ID,
		"auto_renew":        sub.AutoRenew,
		"end_date":          sub.EndDate,
		"days_until_expiry": daysUntilExpiry,
		"renewal_status":    renewalStatus,
		"next_renewal_amount": sub.Price,
	}})
}

// CancelAutoRenewal 取消自动续费
func (ctrl *SubscriptionController) CancelAutoRenewal(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.Subscription{}).Where("id = ?", id).Update("auto_renew", false)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已取消自动续费"})
}

// PaymentWebhook 支付回调
func (ctrl *SubscriptionController) PaymentWebhook(c *gin.Context) {
	var notify struct {
		OrderNo string `json:"order_no"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&notify); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if notify.Status == "success" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "支付确认成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "支付失败"})
	}
}

// GetRenewalLogs 获取续费日志
func (ctrl *SubscriptionController) GetRenewalLogs(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var logs []models.SubscriptionRenewalLog
	var total int64

	ctrl.DB.Model(&models.SubscriptionRenewalLog{}).Where("subscription_id = ?", id).Count(&total)
	ctrl.DB.Where("subscription_id = ?", id).Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"list": logs, "total": total, "page": page, "page_size": pageSize,
	}})
}

// ResumeSubscription 恢复已暂停的订阅
func (ctrl *SubscriptionController) ResumeSubscription(c *gin.Context) {
	id := c.Param("id")

	var sub models.Subscription
	if err := ctrl.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}

	if sub.Status != "suspended" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有暂停状态的订阅才能恢复"})
		return
	}

	sub.Status = "active"
	sub.RetryCount = 0
	sub.RenewFailReason = ""
	ctrl.DB.Save(&sub)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订阅已恢复", "data": sub})
}

func uintVal(s string) uint {
	var v uint
	fmt.Sscanf(s, "%d", &v)
	return v
}
