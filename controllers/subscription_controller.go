package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionController 订阅管理控制器
type SubscriptionController struct {
	DB *gorm.DB
}

// ===== 请求结构 =====

// CreatePlanRequest 创建订阅计划请求
type CreatePlanRequest struct {
	PlanName     string                 `json:"plan_name" binding:"required"`
	PlanType     string                 `json:"plan_type" binding:"required"` // free/basic/pro/enterprise
	Price        float64                `json:"price"`
	Currency     string                 `json:"currency"`
	DurationDays int                    `json:"duration_days"`
	Features     map[string]interface{} `json:"features"`
	Quotas       map[string]interface{} `json:"quotas"`
	SortOrder    int                    `json:"sort_order"`
}

// UpdatePlanRequest 更新订阅计划请求
type UpdatePlanRequest struct {
	PlanName     string                 `json:"plan_name"`
	PlanType     string                 `json:"plan_type"`
	Price        *float64               `json:"price"`
	Currency     string                 `json:"currency"`
	DurationDays *int                   `json:"duration_days"`
	Features     map[string]interface{} `json:"features"`
	Quotas       map[string]interface{} `json:"quotas"`
	Status       string                 `json:"status"`
	SortOrder    *int                   `json:"sort_order"`
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

// PlanChangeRequest 变更计划请求
type PlanChangeRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

// CancelSubscriptionRequest 取消订阅请求
type CancelSubscriptionRequest struct {
	Reason string `json:"reason"`
}

// ===== 订阅计划管理 =====

// ListPlans 获取订阅计划列表
func (c *SubscriptionController) ListPlans(ctx *gin.Context) {
	var plans []models.SubscriptionPlan
	query := c.DB.Model(&models.SubscriptionPlan{})

	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// 类型过滤
	if planType := ctx.Query("plan_type"); planType != "" {
		query = query.Where("plan_type = ?", planType)
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

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, id ASC").Offset(offset).Limit(pageSize).Find(&plans).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": plans,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// GetPlan 获取订阅计划详情
func (c *SubscriptionController) GetPlan(ctx *gin.Context) {
	id := ctx.Param("id")
	var plan models.SubscriptionPlan

	// 尝试用 ID 查找
	if err := c.DB.First(&plan, id).Error; err != nil {
		// 尝试用 plan_id 查找
		if err := c.DB.Where("plan_id = ?", id).First(&plan).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅计划不存在", "error_code": "ERR_NOT_FOUND"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    plan,
	})
}

// CreatePlan 创建订阅计划
func (c *SubscriptionController) CreatePlan(ctx *gin.Context) {
	var req CreatePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 设置默认值
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}
	durationDays := req.DurationDays
	if durationDays == 0 {
		durationDays = 30
	}

	plan := models.SubscriptionPlan{
		PlanID:       "plan-" + uuid.New().String(),
		PlanName:     req.PlanName,
		PlanType:     req.PlanType,
		Price:        req.Price,
		Currency:     currency,
		DurationDays: durationDays,
		Features:     models.JSON(req.Features),
		Quotas:       models.JSON(req.Quotas),
		Status:       "active",
		SortOrder:    req.SortOrder,
	}

	if err := c.DB.Create(&plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    plan,
	})
}

// UpdatePlan 更新订阅计划
func (c *SubscriptionController) UpdatePlan(ctx *gin.Context) {
	id := ctx.Param("id")
	var plan models.SubscriptionPlan

	if err := c.DB.Where("plan_id = ?", id).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅计划不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var req UpdatePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 更新字段
	if req.PlanName != "" {
		plan.PlanName = req.PlanName
	}
	if req.PlanType != "" {
		plan.PlanType = req.PlanType
	}
	if req.Price != nil {
		plan.Price = *req.Price
	}
	if req.Currency != "" {
		plan.Currency = req.Currency
	}
	if req.DurationDays != nil {
		plan.DurationDays = *req.DurationDays
	}
	if req.Features != nil {
		plan.Features = models.JSON(req.Features)
	}
	if req.Quotas != nil {
		plan.Quotas = models.JSON(req.Quotas)
	}
	if req.Status != "" {
		plan.Status = req.Status
	}
	if req.SortOrder != nil {
		plan.SortOrder = *req.SortOrder
	}

	if err := c.DB.Save(&plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    plan,
	})
}

// DeletePlan 删除订阅计划（软删除）
func (c *SubscriptionController) DeletePlan(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Where("plan_id = ?", id).Delete(&models.SubscriptionPlan{})
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅计划不存在", "error_code": "ERR_NOT_FOUND"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ===== 用户订阅管理 =====

// GetCurrentSubscription 获取当前用户订阅
func (c *SubscriptionController) GetCurrentSubscription(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	var subscription models.UserSubscription
	if err := c.DB.Where("user_id = ? AND status IN ?", userID, []string{"active", "pending"}).Order("created_at DESC").First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data":    nil,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取关联的计划信息
	var plan models.SubscriptionPlan
	c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"subscription": subscription,
			"plan":         plan,
		},
	})
}

// Subscribe 订阅
func (c *SubscriptionController) Subscribe(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	var req SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 查询计划
	var plan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ? AND status = ?", req.PlanID, "active").First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅计划不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 检查是否已有活跃订阅
	var existingSub models.UserSubscription
	if err := c.DB.Where("user_id = ? AND status IN ?", userID, []string{"active", "pending"}).First(&existingSub).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "已有活跃订阅，请先取消或升级", "error_code": "ERR_SUBSCRIPTION_EXISTS"})
		return
	}

	// 创建订阅
	now := time.Now()
	subscription := models.UserSubscription{
		SubID:      "sub-" + uuid.New().String(),
		UserID:     userID,
		PlanID:     req.PlanID,
		Status:     "active",
		StartTime:  now,
		ExpireTime: now.AddDate(0, 0, plan.DurationDays),
		AutoRenew:  true,
	}

	// 开启事务
	tx := c.DB.Begin()

	if err := tx.Create(&subscription).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "订阅创建失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 创建订阅变更记录
	change := models.SubscriptionChange{
		ChangeID:    "change-" + uuid.New().String(),
		UserID:     userID,
		SubID:      subscription.SubID,
		ChangeType: "create",
		ToPlanID:   req.PlanID,
		Amount:     plan.Price,
		EffectiveAt: now,
	}
	tx.Create(&change)

	// 创建账单记录（如果是付费计划）
	if plan.Price > 0 {
		billing := models.BillingRecord{
			BillID:      "bill-" + uuid.New().String(),
			UserID:      userID,
			Type:        models.BillingTypeSubscription,
			Amount:      plan.Price,
			Currency:    plan.Currency,
			Status:      models.BillingStatusPending,
			Description: "订阅 " + plan.PlanName,
		}
		tx.Create(&billing)
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "订阅成功",
		"data": gin.H{
			"subscription": subscription,
			"plan":         plan,
		},
	})
}

// GetSubscription 获取订阅详情
func (c *SubscriptionController) GetSubscription(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	id := ctx.Param("id")

	var subscription models.UserSubscription
	query := c.DB.Where("sub_id = ?", id)

	// 非管理员只能查看自己的订阅
	if !c.isAdmin(ctx) {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取关联的计划信息
	var plan models.SubscriptionPlan
	c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"subscription": subscription,
			"plan":         plan,
		},
	})
}

// CancelSubscription 取消订阅
func (c *SubscriptionController) CancelSubscription(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	id := ctx.Param("id")
	var req CancelSubscriptionRequest
	ctx.ShouldBindJSON(&req)

	var subscription models.UserSubscription
	if err := c.DB.Where("sub_id = ? AND user_id = ?", id, userID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if subscription.Status != "active" && subscription.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "订阅已取消或过期", "error_code": "ERR_INVALID_STATUS"})
		return
	}

	// 关闭自动续费，但订阅在到期前仍然有效
	subscription.AutoRenew = false

	tx := c.DB.Begin()

	// 记录变更
	change := models.SubscriptionChange{
		ChangeID:     "change-" + uuid.New().String(),
		UserID:      userID,
		SubID:       subscription.SubID,
		ChangeType:  "cancel",
		FromPlanID:  subscription.PlanID,
		ToPlanID:    subscription.PlanID,
		ChangeReason: req.Reason,
		EffectiveAt: time.Now(),
	}
	tx.Create(&change)

	// 更新订阅状态
	subscription.Status = "cancelled"
	tx.Save(&subscription)

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功，订阅将在到期日后失效",
		"data":    subscription,
	})
}

// RenewSubscription 续费
func (c *SubscriptionController) RenewSubscription(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	id := ctx.Param("id")

	var subscription models.UserSubscription
	if err := c.DB.Where("sub_id = ? AND user_id = ?", id, userID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取计划信息计算续费价格
	var plan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "计划不存在", "error_code": "ERR_INTERNAL"})
		return
	}

	// 续费逻辑：延长到期时间
	now := time.Now()
	newExpireTime := subscription.ExpireTime

	// 如果已过期，从当前时间开始计算
	if newExpireTime.Before(now) {
		newExpireTime = now
	}

	newExpireTime = newExpireTime.AddDate(0, 0, plan.DurationDays)

	tx := c.DB.Begin()

	// 更新订阅
	subscription.ExpireTime = newExpireTime
	subscription.Status = "active"
	subscription.AutoRenew = true
	tx.Save(&subscription)

	// 记录变更
	change := models.SubscriptionChange{
		ChangeID:    "change-" + uuid.New().String(),
		UserID:     userID,
		SubID:      subscription.SubID,
		ChangeType: "renew",
		ToPlanID:   subscription.PlanID,
		Amount:     plan.Price,
		EffectiveAt: now,
	}
	tx.Create(&change)

	// 创建账单
	billing := models.BillingRecord{
		BillID:      "bill-" + uuid.New().String(),
		UserID:      userID,
		Type:        models.BillingTypeSubscription,
		Amount:      plan.Price,
		Currency:    plan.Currency,
		Status:      models.BillingStatusPending,
		Description: "续费 " + plan.PlanName,
	}
	tx.Create(&billing)

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "续费成功",
		"data": gin.H{
			"subscription": subscription,
			"plan":        plan,
		},
	})
}

// UpgradeSubscription 升级订阅
func (c *SubscriptionController) UpgradeSubscription(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	id := ctx.Param("id")
	var req PlanChangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 获取当前订阅
	var subscription models.UserSubscription
	if err := c.DB.Where("sub_id = ? AND user_id = ?", id, userID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if subscription.Status != "active" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "订阅未激活，无法升级", "error_code": "ERR_INVALID_STATUS"})
		return
	}

	// 获取新旧计划
	var oldPlan, newPlan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ?", subscription.PlanID).First(&oldPlan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "当前计划不存在", "error_code": "ERR_INTERNAL"})
		return
	}
	if err := c.DB.Where("plan_id = ? AND status = ?", req.PlanID, "active").First(&newPlan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "目标计划不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 计算升级差价
	upgradeAmount := newPlan.Price - oldPlan.Price
	if upgradeAmount < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4004, "message": "目标计划价格不能低于当前计划", "error_code": "ERR_PRICE_INVALID"})
		return
	}

	tx := c.DB.Begin()

	// 更新订阅
	subscription.PlanID = newPlan.PlanID
	tx.Save(&subscription)

	// 记录变更
	change := models.SubscriptionChange{
		ChangeID:    "change-" + uuid.New().String(),
		UserID:     userID,
		SubID:      subscription.SubID,
		ChangeType: "upgrade",
		FromPlanID: oldPlan.PlanID,
		ToPlanID:   newPlan.PlanID,
		Amount:     upgradeAmount,
		EffectiveAt: time.Now(),
	}
	tx.Create(&change)

	// 创建账单（如果需要补差价）
	if upgradeAmount > 0 {
		billing := models.BillingRecord{
			BillID:      "bill-" + uuid.New().String(),
			UserID:      userID,
			Type:        models.BillingTypeUpgrade,
			Amount:      upgradeAmount,
			Currency:    newPlan.Currency,
			Status:      models.BillingStatusPending,
			Description: "升级从 " + oldPlan.PlanName + " 到 " + newPlan.PlanName,
		}
		tx.Create(&billing)
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "升级成功",
		"data": gin.H{
			"subscription": subscription,
			"old_plan":    oldPlan,
			"new_plan":    newPlan,
		},
	})
}

// DowngradeSubscription 降级订阅
func (c *SubscriptionController) DowngradeSubscription(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	id := ctx.Param("id")
	var req PlanChangeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 获取当前订阅
	var subscription models.UserSubscription
	if err := c.DB.Where("sub_id = ? AND user_id = ?", id, userID).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "订阅不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if subscription.Status != "active" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "订阅未激活，无法降级", "error_code": "ERR_INVALID_STATUS"})
		return
	}

	// 获取新旧计划
	var oldPlan, newPlan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ?", subscription.PlanID).First(&oldPlan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "当前计划不存在", "error_code": "ERR_INTERNAL"})
		return
	}
	if err := c.DB.Where("plan_id = ? AND status = ?", req.PlanID, "active").First(&newPlan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "目标计划不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	tx := c.DB.Begin()

	// 降级在当前周期结束后生效
	subscription.PlanID = newPlan.PlanID
	tx.Save(&subscription)

	// 记录变更
	change := models.SubscriptionChange{
		ChangeID:    "change-" + uuid.New().String(),
		UserID:     userID,
		SubID:      subscription.SubID,
		ChangeType: "downgrade",
		FromPlanID: oldPlan.PlanID,
		ToPlanID:   newPlan.PlanID,
		Amount:     0,
		EffectiveAt: subscription.ExpireTime, // 下个周期生效
	}
	tx.Create(&change)

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "降级成功，将在当前周期结束后生效",
		"data": gin.H{
			"subscription": subscription,
			"old_plan":    oldPlan,
			"new_plan":    newPlan,
		},
	})
}

// ListUserSubscriptions 获取用户所有订阅历史
func (c *SubscriptionController) ListUserSubscriptions(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	var subscriptions []models.UserSubscription
	query := c.DB.Where("user_id = ?", userID)

	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&subscriptions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": subscriptions,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// isAdmin 检查用户是否为管理员
func (c *SubscriptionController) isAdmin(ctx *gin.Context) bool {
	role := ctx.GetString("role")
	return role == "admin" || role == "super_admin"
}
