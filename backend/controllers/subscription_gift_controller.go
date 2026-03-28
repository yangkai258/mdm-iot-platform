package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SubscriptionGiftController 订阅赠送控制器
type SubscriptionGiftController struct {
	DB *gorm.DB
}

// NewSubscriptionGiftController 创建控制器
func NewSubscriptionGiftController(db *gorm.DB) *SubscriptionGiftController {
	return &SubscriptionGiftController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *SubscriptionGiftController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/subscription/gifts", ctrl.ListGifts)
	rg.GET("/subscription/gifts/:id", ctrl.GetGift)
	rg.POST("/subscription/gifts", ctrl.CreateGift)
	rg.POST("/subscription/gifts/claim", ctrl.ClaimGift)
	rg.GET("/subscription/gifts/code/:code", ctrl.GetGiftByCode)
	rg.PUT("/subscription/gifts/:id/cancel", ctrl.CancelGift)
	rg.GET("/subscription/gifts/sent", ctrl.ListSentGifts)
	rg.GET("/subscription/gifts/received", ctrl.ListReceivedGifts)
}

// generateGiftCode 生成赠送码
func generateGiftCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// ListGifts 获取赠送列表（管理员用）
func (ctrl *SubscriptionGiftController) ListGifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}

	query := ctrl.DB.Model(&models.SubscriptionGift{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	var list []models.SubscriptionGift
	query.Count(&total)

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": list, "total": total, "page": page, "page_size": pageSize},
	})
}

// GetGift 获取单个赠送记录
func (ctrl *SubscriptionGiftController) GetGift(c *gin.Context) {
	id := c.Param("id")

	var gift models.SubscriptionGift
	if err := ctrl.DB.Where("id = ?", id).First(&gift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "赠送记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gift})
}

// CreateGift 创建赠送
func (ctrl *SubscriptionGiftController) CreateGift(c *gin.Context) {
	var req struct {
		RecipientEmail string `json:"recipient_email" binding:"required"`
		RecipientName  string `json:"recipient_name"`
		PlanID        uint   `json:"plan_id" binding:"required"`
		PlanName      string `json:"plan_name"`
		Duration      int    `json:"duration"` // 默认30天
		Message       string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Duration <= 0 {
		req.Duration = 30
	}

	gift := models.SubscriptionGift{
		GiftCode:      generateGiftCode(),
		SenderID:     c.GetString("user_id"),
		SenderName:   c.GetString("username"),
		RecipientEmail: req.RecipientEmail,
		RecipientName: req.RecipientName,
		PlanID:       req.PlanID,
		PlanName:     req.PlanName,
		Duration:     req.Duration,
		Status:       "pending",
		Message:      req.Message,
	}

	// 设置过期时间（7天后）
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	gift.ExpiresAt = &expiresAt

	if err := ctrl.DB.Create(&gift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "赠送码已生成",
		"data": gin.H{
			"gift_code": gift.GiftCode,
			"expires_at": gift.ExpiresAt,
			"duration":  gift.Duration,
			"plan_name": gift.PlanName,
		},
	})
}

// GetGiftByCode 通过赠送码获取赠送信息
func (ctrl *SubscriptionGiftController) GetGiftByCode(c *gin.Context) {
	code := c.Param("code")

	var gift models.SubscriptionGift
	if err := ctrl.DB.Where("gift_code = ?", code).First(&gift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "赠送码无效"})
		return
	}

	if gift.Status == "expired" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "赠送码已过期"})
		return
	}

	if gift.Status == "claimed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "赠送码已被使用"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"gift_code":   gift.GiftCode,
			"sender_name": gift.SenderName,
			"plan_name":   gift.PlanName,
			"duration":    gift.Duration,
			"message":     gift.Message,
			"expires_at":  gift.ExpiresAt,
			"status":      gift.Status,
		},
	})
}

// ClaimGift 领取赠送
func (ctrl *SubscriptionGiftController) ClaimGift(c *gin.Context) {
	var req struct {
		GiftCode     string `json:"gift_code" binding:"required"`
		RecipientID  string `json:"recipient_id" binding:"required"`
		RecipientName string `json:"recipient_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var gift models.SubscriptionGift
	if err := ctrl.DB.Where("gift_code = ?", req.GiftCode).First(&gift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "赠送码无效"})
		return
	}

	if gift.Status == "claimed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "赠送码已被使用"})
		return
	}

	if gift.Status == "expired" || (gift.ExpiresAt != nil && gift.ExpiresAt.Before(time.Now())) {
		ctrl.DB.Model(&gift).Update("status", "expired")
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "赠送码已过期"})
		return
	}

	// 更新领取状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":       "claimed",
		"recipient_id":  req.RecipientID,
		"recipient_name": req.RecipientName,
		"claimed_at":     now,
	}

	if err := ctrl.DB.Model(&gift).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "领取失败"})
		return
	}

	// 记录使用
	usage := models.SubscriptionGiftUsage{
		GiftID:         gift.ID,
		ClaimedFeature: fmt.Sprintf("subscription:%s:%dd", gift.PlanName, gift.Duration),
	}
	ctrl.DB.Create(&usage)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "领取成功",
		"data": gin.H{
			"plan_name": gift.PlanName,
			"duration":  gift.Duration,
			"claimed_at": now,
		},
	})
}

// CancelGift 取消赠送
func (ctrl *SubscriptionGiftController) CancelGift(c *gin.Context) {
	id := c.Param("id")

	var gift models.SubscriptionGift
	if err := ctrl.DB.Where("id = ?", id).First(&gift).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "赠送记录不存在"})
		return
	}

	if gift.SenderID != c.GetString("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "只能取消自己的赠送"})
		return
	}

	if gift.Status == "claimed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已被领取无法取消"})
		return
	}

	ctrl.DB.Model(&gift).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已取消"})
}

// ListSentGifts 获取我发送的赠送
func (ctrl *SubscriptionGiftController) ListSentGifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	var list []models.SubscriptionGift
	ctrl.DB.Model(&models.SubscriptionGift{}).
		Where("sender_id = ?", c.GetString("user_id")).
		Count(&total)

	ctrl.DB.Where("sender_id = ?", c.GetString("user_id")).
		Order("created_at DESC").
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": list, "total": total, "page": page, "page_size": pageSize},
	})
}

// ListReceivedGifts 获取我收到的赠送
func (ctrl *SubscriptionGiftController) ListReceivedGifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	var list []models.SubscriptionGift
	ctrl.DB.Model(&models.SubscriptionGift{}).
		Where("recipient_id = ? OR recipient_email = ?", c.GetString("user_id"), c.GetString("email")).
		Count(&total)

	ctrl.DB.Where("recipient_id = ? OR recipient_email = ?", c.GetString("user_id"), c.GetString("email")).
		Order("created_at DESC").
		Offset((page-1)*pageSize).
		Limit(pageSize).
		Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"list": list, "total": total, "page": page, "page_size": pageSize},
	})
}
