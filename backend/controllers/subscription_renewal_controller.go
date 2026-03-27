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

// SubscriptionRenewalController 订阅自动续费控制器
type SubscriptionRenewalController struct {
	DB             *gorm.DB
	RenewalService *services.SubscriptionRenewalService
}

// NewSubscriptionRenewalController 创建续费控制器
func NewSubscriptionRenewalController(db *gorm.DB) *SubscriptionRenewalController {
	return &SubscriptionRenewalController{
		DB:             db,
		RenewalService: services.NewSubscriptionRenewalService(db),
	}
}

// RegisterRoutes 注册续费路由
func (ctrl *SubscriptionRenewalController) RegisterRoutes(rg *gin.RouterGroup) {
	// 手动续费
	rg.POST("/subscriptions/:id/renew", ctrl.ManualRenew)

	// 自动续费管理
	autoRenew := rg.Group("/subscriptions/auto-renewal")
	{
		autoRenew.POST("/enable", ctrl.EnableAutoRenewal)
		autoRenew.POST("/disable", ctrl.DisableAutoRenewal)
		autoRenew.GET("/status", ctrl.GetAutoRenewalStatus)
		autoRenew.GET("/settings", ctrl.GetAutoRenewalSettings)
		autoRenew.PUT("/settings", ctrl.UpdateAutoRenewalSettings)
	}

	// 续费记录
	rg.GET("/subscriptions/:id/renewals", ctrl.GetRenewalHistory)
	rg.GET("/subscriptions/:id/renewals/:renewal_id", ctrl.GetRenewalDetail)
}

// ManualRenew 手动续费
// POST /api/v1/subscriptions/:id/renew
func (ctrl *SubscriptionRenewalController) ManualRenew(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的订阅ID"})
		return
	}

	var sub models.Subscription
	if err := ctrl.DB.First(&sub, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询订阅失败"})
		return
	}

	// 获取支付方式
	var paymentMethodID *uint
	var paymentMethod string
	if reqMethodID := c.PostForm("payment_method_id"); reqMethodID != "" {
		pid, _ := strconv.ParseUint(reqMethodID, 10, 32)
		paymentMethodID = new(uint)
		*paymentMethodID = uint(pid)

		var pm models.PaymentMethod
		if err := ctrl.DB.First(&pm, *paymentMethodID).Error; err == nil {
			paymentMethod = pm.MethodType
		}
	} else {
		// 使用用户默认支付方式
		var defaultPM models.PaymentMethod
		if err := ctrl.DB.Where("user_id = ? AND is_default = ? AND status = ?", sub.UserID, true, "active").First(&defaultPM).Error; err == nil {
			paymentMethodID = &defaultPM.ID
			paymentMethod = defaultPM.MethodType
		}
	}

	// 计算续费后的到期日期
	newExpiredDate := sub.EndDate.AddDate(0, 0, sub.Duration)

	// 创建续费记录
	renewal := models.SubscriptionRenewal{
		SubscriptionID:  sub.ID,
		UserID:          sub.UserID,
		PaymentMethodID: paymentMethodID,
		Amount:          sub.Price,
		PaymentMethod:   paymentMethod,
		PaymentStatus:   "processing",
		RenewalDate:     time.Now(),
		ExpiredDate:     newExpiredDate,
	}

	// 尝试扣费
	success, err := ctrl.RenewalService.ProcessPayment(sub, renewal)

	if success {
		// 扣费成功，更新订阅
		renewal.PaymentStatus = "success"
		ctrl.DB.Create(&renewal)

		now := time.Now()
		ctrl.DB.Model(&sub).Updates(map[string]interface{}{
			"end_date":       newExpiredDate,
			"status":         "active",
			"auto_renew":     true,
			"last_renew_at":  &now,
			"renew_count":    sub.RenewCount + 1,
			"retry_count":    0,
			"renew_fail_reason": "",
		})

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "续费成功",
			"data": gin.H{
				"subscription_id": sub.ID,
				"renewal_id":     renewal.ID,
				"new_end_date":   newExpiredDate,
				"amount":         sub.Price,
				"payment_method": paymentMethod,
			},
		})
	} else {
		// 扣费失败
		renewal.PaymentStatus = "failed"
		renewal.FailReason = err.Error()
		renewal.RetryCount = 1
		ctrl.DB.Create(&renewal)

		ctrl.DB.Model(&sub).Updates(map[string]interface{}{
			"retry_count":       1,
			"renew_fail_reason": err.Error(),
		})

		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": fmt.Sprintf("续费失败: %s", err.Error()),
			"data": gin.H{
				"subscription_id": sub.ID,
				"renewal_id":     renewal.ID,
				"fail_reason":    err.Error(),
				"retry_count":    1,
			},
		})
	}
}

// EnableAutoRenewal 启用自动续费
// POST /api/v1/subscriptions/auto-renewal/enable
func (ctrl *SubscriptionRenewalController) EnableAutoRenewal(c *gin.Context) {
	var req struct {
		SubscriptionID uint `json:"subscription_id" binding:"required"`
		PaymentMethodID uint `json:"payment_method_id"`
		ReminderDays   string `json:"reminder_days"` // "7,3,1"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证订阅存在
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, req.SubscriptionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询订阅失败"})
		return
	}

	// 验证支付方式
	if req.PaymentMethodID != 0 {
		var pm models.PaymentMethod
		if err := ctrl.DB.First(&pm, req.PaymentMethodID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "支付方式不存在"})
			return
		}
		if pm.UserID != sub.UserID {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "支付方式不属于当前用户"})
			return
		}
	}

	// 设置默认值
	reminderDays := req.ReminderDays
	if reminderDays == "" {
		reminderDays = "7,3,1"
	}

	// 创建或更新自动续费设置
	setting := models.AutoRenewalSetting{}
	err := ctrl.DB.Where("user_id = ? AND subscription_id = ?", sub.UserID, req.SubscriptionID).First(&setting).Error

	now := time.Now()
	if err == gorm.ErrRecordNotFound {
		// 新建
		setting = models.AutoRenewalSetting{
			UserID:         sub.UserID,
			SubscriptionID: req.SubscriptionID,
			Enabled:        true,
			PaymentMethodID: &req.PaymentMethodID,
			ReminderDays:   reminderDays,
			LastRemindAt:   &now,
		}
		if req.PaymentMethodID == 0 {
			setting.PaymentMethodID = nil
		}
		ctrl.DB.Create(&setting)
	} else {
		// 更新
		updates := map[string]interface{}{
			"enabled":        true,
			"reminder_days":  reminderDays,
			"last_remind_at": &now,
		}
		if req.PaymentMethodID != 0 {
			updates["payment_method_id"] = req.PaymentMethodID
		} else {
			updates["payment_method_id"] = nil
		}
		ctrl.DB.Model(&setting).Updates(updates)
	}

	// 更新订阅自动续费标志
	ctrl.DB.Model(&sub).Update("auto_renew", true)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "自动续费已启用",
		"data": gin.H{
			"subscription_id": req.SubscriptionID,
			"enabled":         true,
			"payment_method_id": req.PaymentMethodID,
			"reminder_days":   reminderDays,
		},
	})
}

// DisableAutoRenewal 禁用自动续费
// POST /api/v1/subscriptions/auto-renewal/disable
func (ctrl *SubscriptionRenewalController) DisableAutoRenewal(c *gin.Context) {
	var req struct {
		SubscriptionID uint `json:"subscription_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证订阅存在
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, req.SubscriptionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询订阅失败"})
		return
	}

	// 更新自动续费设置
	setting := models.AutoRenewalSetting{}
	if err := ctrl.DB.Where("user_id = ? AND subscription_id = ?", sub.UserID, req.SubscriptionID).First(&setting).Error; err == nil {
		ctrl.DB.Model(&setting).Update("enabled", false)
	}

	// 更新订阅自动续费标志
	ctrl.DB.Model(&sub).Update("auto_renew", false)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "自动续费已禁用",
		"data": gin.H{
			"subscription_id": req.SubscriptionID,
			"enabled":         false,
		},
	})
}

// GetAutoRenewalStatus 获取自动续费状态
// GET /api/v1/subscriptions/auto-renewal/status
func (ctrl *SubscriptionRenewalController) GetAutoRenewalStatus(c *gin.Context) {
	subscriptionIDStr := c.Query("subscription_id")
	userIDStr := c.Query("user_id")

	if subscriptionIDStr == "" && userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少订阅ID或用户ID"})
		return
	}

	var results []gin.H

	if subscriptionIDStr != "" {
		// 按订阅ID查询
		subID, _ := strconv.ParseUint(subscriptionIDStr, 10, 32)
		var sub models.Subscription
		if err := ctrl.DB.First(&sub, subID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
			return
		}

		var setting models.AutoRenewalSetting
		ctrl.DB.Where("subscription_id = ?", subID).First(&setting)

		daysUntilExpiry := int(time.Until(sub.EndDate).Hours() / 24)

		results = append(results, gin.H{
			"subscription_id":    sub.ID,
			"user_id":            sub.UserID,
			"plan_name":          sub.PlanName,
			"end_date":           sub.EndDate,
			"days_until_expiry":  daysUntilExpiry,
			"auto_renew_enabled": sub.AutoRenew,
			"setting_enabled":    setting.Enabled,
			"payment_method_id":  setting.PaymentMethodID,
			"reminder_days":      setting.ReminderDays,
		})
	} else {
		// 按用户ID查询所有订阅的自动续费状态
		uid, _ := strconv.ParseUint(userIDStr, 10, 32)
		var subs []models.Subscription
		ctrl.DB.Where("user_id = ?", uid).Find(&subs)

		for _, sub := range subs {
			var setting models.AutoRenewalSetting
			ctrl.DB.Where("subscription_id = ?", sub.ID).First(&setting)

			daysUntilExpiry := int(time.Until(sub.EndDate).Hours() / 24)

			results = append(results, gin.H{
				"subscription_id":    sub.ID,
				"user_id":            sub.UserID,
				"plan_name":          sub.PlanName,
				"end_date":           sub.EndDate,
				"days_until_expiry":  daysUntilExpiry,
				"auto_renew_enabled": sub.AutoRenew,
				"setting_enabled":    setting.Enabled,
				"payment_method_id":  setting.PaymentMethodID,
				"reminder_days":      setting.ReminderDays,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": results,
	})
}

// GetAutoRenewalSettings 获取自动续费设置详情
// GET /api/v1/subscriptions/auto-renewal/settings
func (ctrl *SubscriptionRenewalController) GetAutoRenewalSettings(c *gin.Context) {
	subscriptionIDStr := c.Query("subscription_id")
	if subscriptionIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少订阅ID"})
		return
	}

	subID, _ := strconv.ParseUint(subscriptionIDStr, 10, 32)
	var sub models.Subscription
	if err := ctrl.DB.First(&sub, subID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}

	var setting models.AutoRenewalSetting
	if err := ctrl.DB.Where("subscription_id = ?", subID).First(&setting).Error; err != nil {
		// 返回默认值
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"subscription_id": sub.ID,
				"enabled":         sub.AutoRenew,
				"reminder_days":   "7,3,1",
			},
		})
		return
	}

	// 获取支付方式详情
	var paymentMethod *models.PaymentMethod
	if setting.PaymentMethodID != nil {
		var pm models.PaymentMethod
		if err := ctrl.DB.First(&pm, *setting.PaymentMethodID).Error; err == nil {
			paymentMethod = &pm
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":               setting.ID,
			"subscription_id":  setting.SubscriptionID,
			"enabled":          setting.Enabled,
			"payment_method_id": setting.PaymentMethodID,
			"payment_method":   paymentMethod,
			"reminder_days":    setting.ReminderDays,
			"last_remind_at":   setting.LastRemindAt,
		},
	})
}

// UpdateAutoRenewalSettings 更新自动续费设置
// PUT /api/v1/subscriptions/auto-renewal/settings
func (ctrl *SubscriptionRenewalController) UpdateAutoRenewalSettings(c *gin.Context) {
	var req struct {
		SubscriptionID  uint   `json:"subscription_id" binding:"required"`
		Enabled          *bool  `json:"enabled"`
		PaymentMethodID  *uint  `json:"payment_method_id"`
		ReminderDays     string `json:"reminder_days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var sub models.Subscription
	if err := ctrl.DB.First(&sub, req.SubscriptionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}

	var setting models.AutoRenewalSetting
	isNew := false
	if err := ctrl.DB.Where("subscription_id = ?", req.SubscriptionID).First(&setting).Error; err == gorm.ErrRecordNotFound {
		isNew = true
		setting = models.AutoRenewalSetting{
			UserID:         sub.UserID,
			SubscriptionID: req.SubscriptionID,
		}
	}

	// 更新字段
	if req.Enabled != nil {
		setting.Enabled = *req.Enabled
		ctrl.DB.Model(&sub).Update("auto_renew", *req.Enabled)
	}
	if req.PaymentMethodID != nil {
		setting.PaymentMethodID = req.PaymentMethodID
	}
	if req.ReminderDays != "" {
		setting.ReminderDays = req.ReminderDays
	}

	if isNew {
		ctrl.DB.Create(&setting)
	} else {
		ctrl.DB.Model(&setting).Updates(map[string]interface{}{
			"enabled":           setting.Enabled,
			"payment_method_id": setting.PaymentMethodID,
			"reminder_days":     setting.ReminderDays,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "设置已更新",
		"data": gin.H{
			"subscription_id":  setting.SubscriptionID,
			"enabled":           setting.Enabled,
			"payment_method_id": setting.PaymentMethodID,
			"reminder_days":     setting.ReminderDays,
		},
	})
}

// GetRenewalHistory 获取续费历史记录
// GET /api/v1/subscriptions/:id/renewals
func (ctrl *SubscriptionRenewalController) GetRenewalHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var renewals []models.SubscriptionRenewal
	var total int64

	query := ctrl.DB.Model(&models.SubscriptionRenewal{}).Where("subscription_id = ?", id)
	query.Count(&total)

	offset := (page - 1) * pageSize
	ctrl.DB.Where("subscription_id = ?", id).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&renewals)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      renewals,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRenewalDetail 获取续费记录详情
// GET /api/v1/subscriptions/:id/renewals/:renewal_id
func (ctrl *SubscriptionRenewalController) GetRenewalDetail(c *gin.Context) {
	renewalIDStr := c.Param("renewal_id")
	renewalID, err := strconv.ParseUint(renewalIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的续费记录ID"})
		return
	}

	var renewal models.SubscriptionRenewal
	if err := ctrl.DB.First(&renewal, renewalID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "续费记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取关联的支付方式信息
	var paymentMethod *models.PaymentMethod
	if renewal.PaymentMethodID != nil {
		var pm models.PaymentMethod
		if err := ctrl.DB.First(&pm, *renewal.PaymentMethodID).Error; err == nil {
			paymentMethod = &pm
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"renewal":        renewal,
			"payment_method": paymentMethod,
		},
	})
}
