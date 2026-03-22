package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NotificationController 通知控制器
type NotificationController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// ==================== 通知相关 ====================

// ListNotifications 获取通知列表
func (c *NotificationController) ListNotifications(ctx *gin.Context) {
	var notifications []models.Notification
	query := c.DB.Model(&models.Notification{})

	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := ctx.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      notifications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// SendNotification 发送通知（MQTT下发）
func (c *NotificationController) SendNotification(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	var req SendNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 如果指定了模板ID，先替换变量
	title := req.Title
	content := req.Content
	if req.TemplateID > 0 {
		var tmpl models.NotificationTemplate
		if err := c.DB.Where("id = ? AND enabled = ?", req.TemplateID, true).First(&tmpl).Error; err == nil {
			title = ReplaceTemplateVariables(tmpl.TitleTpl, req.Variables)
			content = ReplaceTemplateVariables(tmpl.ContentTpl, req.Variables)
		}
	}

	// 构建MQTT消息
	notifPayload := map[string]interface{}{
		"type":      "notification",
		"title":     title,
		"content":   content,
		"priority":  req.Priority,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// 通过 MQTT 下发到设备
	if mqtt.GlobalMQTTClient != nil {
		// 使用 /device/{id}/down/notification 作为通知下发Topic
		topic := fmt.Sprintf("/device/%s/down/notification", deviceID)
		payload, _ := json.Marshal(notifPayload)
		token := mqtt.GlobalMQTTClient.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":       5001,
				"message":    "通知下发失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}
	}

	now := time.Now()
	// 保存通知记录
	notif := models.Notification{
		DeviceID:   deviceID,
		Title:      title,
		Content:    content,
		Priority:   req.Priority,
		Channel:    req.Channel,
		Status:     "sent",
		SentAt:     &now,
		CreatedBy:  req.CreatedBy,
	}
	c.DB.Create(&notif)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"id":        notif.ID,
			"device_id": deviceID,
			"status":    "sent",
			"sent_at":   now,
		},
	})
}

// GetNotification 获取通知详情
func (c *NotificationController) GetNotification(ctx *gin.Context) {
	id := ctx.Param("id")

	var notif models.Notification
	if err := c.DB.First(&notif, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "通知不存在",
				"error_code": "ERR_NOT_FOUND",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    notif,
	})
}

// SendNotificationFromTemplate 通过模板发送通知
func (c *NotificationController) SendNotificationFromTemplate(ctx *gin.Context) {
	var req SendNotificationFromTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 获取模板
	var tmpl models.NotificationTemplate
	if err := c.DB.Where("id = ?", req.TemplateID).First(&tmpl).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "模板不存在",
				"error_code": "ERR_NOT_FOUND",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 替换变量
	title := ReplaceTemplateVariables(tmpl.TitleTpl, req.Variables)
	content := ReplaceTemplateVariables(tmpl.ContentTpl, req.Variables)

	now := time.Now()

	// 收集所有目标设备ID
	var deviceIDs []string
	switch req.TargetType {
	case "all":
		var devices []models.Device
		c.DB.Select("device_id").Find(&devices)
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.DeviceID)
		}
	case "device":
		deviceIDs = req.TargetIDs
	case "user":
		var devices []models.Device
		c.DB.Select("device_id").Where("owner_id IN ?", req.TargetIDs).Find(&devices)
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.DeviceID)
		}
	}

	// 批量创建通知记录（避免 N+1 查询）
	var notifications []models.Notification
	if len(deviceIDs) > 0 {
		for _, deviceID := range deviceIDs {
			notifications = append(notifications, models.Notification{
				DeviceID:  deviceID,
				Title:     title,
				Content:   content,
				Priority:  1,
				Channel:   "push",
				Status:    "sent",
				SentAt:    &now,
				CreatedBy: req.CreatedBy,
			})
		}
		// 批量插入
		if err := c.DB.Create(&notifications).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":       5001,
				"message":    "通知发送失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"sent_count":   len(notifications),
			"notifications": notifications,
		},
	})
}

func (c *NotificationController) sendToDevice(deviceID, title, content, createdBy string, now *time.Time) models.Notification {
	notif := models.Notification{
		DeviceID:  deviceID,
		Title:     title,
		Content:   content,
		Priority:  1,
		Channel:   "push",
		Status:    "sent",
		SentAt:    now,
		CreatedBy: createdBy,
	}
	c.DB.Create(&notif)
	return notif
}

// DeleteNotification 删除通知
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Delete(&models.Notification{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "删除失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4004,
			"message":    "通知不存在",
			"error_code": "ERR_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ==================== 通知模板相关 ====================

// ListTemplates 获取模板列表
func (c *NotificationController) ListTemplates(ctx *gin.Context) {
	var templates []models.NotificationTemplate
	query := c.DB.Model(&models.NotificationTemplate{})

	if code := ctx.Query("code"); code != "" {
		query = query.Where("code = ?", code)
	}
	if channel := ctx.Query("channel"); channel != "" {
		query = query.Where("channel = ?", channel)
	}
	if enabledStr := ctx.Query("enabled"); enabledStr != "" {
		query = query.Where("enabled = ?", enabledStr)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&templates)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      templates,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateTemplate 创建模板
func (c *NotificationController) CreateTemplate(ctx *gin.Context) {
	var tmpl models.NotificationTemplate
	if err := ctx.ShouldBindJSON(&tmpl); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查编码唯一性
	var exists models.NotificationTemplate
	if err := c.DB.Where("code = ?", tmpl.Code).First(&exists).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4006,
			"message":    "模板编码已存在",
			"error_code": "ERR_DUPLICATE",
		})
		return
	}

	if err := c.DB.Create(&tmpl).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tmpl,
	})
}

// UpdateTemplate 更新模板
func (c *NotificationController) UpdateTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	var tmpl models.NotificationTemplate
	if err := c.DB.First(&tmpl, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "模板不存在",
				"error_code": "ERR_NOT_FOUND",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 不允许修改 code
	delete(updates, "code")

	if err := c.DB.Model(&tmpl).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "更新失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	c.DB.First(&tmpl, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tmpl,
	})
}

// DeleteTemplate 删除模板
func (c *NotificationController) DeleteTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Delete(&models.NotificationTemplate{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "删除失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4004,
			"message":    "模板不存在",
			"error_code": "ERR_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ==================== 企业公告相关 ====================

// ListAnnouncements 获取公告列表
func (c *NotificationController) ListAnnouncements(ctx *gin.Context) {
	var announcements []models.Announcement
	query := c.DB.Model(&models.Announcement{})

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if annType := ctx.Query("type"); annType != "" {
		query = query.Where("type = ?", annType)
	}
	if targetType := ctx.Query("target_type"); targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("priority DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&announcements)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      announcements,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateAnnouncement 创建公告
func (c *NotificationController) CreateAnnouncement(ctx *gin.Context) {
	var ann models.Announcement
	if err := ctx.ShouldBindJSON(&ann); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if err := c.DB.Create(&ann).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    ann,
	})
}

// UpdateAnnouncement 更新公告
func (c *NotificationController) UpdateAnnouncement(ctx *gin.Context) {
	id := ctx.Param("id")

	var ann models.Announcement
	if err := c.DB.First(&ann, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "公告不存在",
				"error_code": "ERR_NOT_FOUND",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if err := c.DB.Model(&ann).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "更新失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	c.DB.First(&ann, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    ann,
	})
}

// DeleteAnnouncement 删除公告
func (c *NotificationController) DeleteAnnouncement(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Delete(&models.Announcement{}, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "删除失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4004,
			"message":    "公告不存在",
			"error_code": "ERR_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// PublishAnnouncement 发布公告
func (c *NotificationController) PublishAnnouncement(ctx *gin.Context) {
	id := ctx.Param("id")

	var ann models.Announcement
	if err := c.DB.First(&ann, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "公告不存在",
				"error_code": "ERR_NOT_FOUND",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":      "published",
		"published_at": &now,
	}

	if err := c.DB.Model(&ann).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "发布失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	c.DB.First(&ann, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    ann,
	})
}

// ==================== 辅助函数 ====================

// SendNotificationRequest 发送通知请求
type SendNotificationRequest struct {
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Priority   int                    `json:"priority"`
	Channel    string                 `json:"channel"`
	TemplateID uint                   `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
	CreatedBy  string                 `json:"created_by"`
}

// SendNotificationFromTemplateRequest 通过模板发送通知请求
type SendNotificationFromTemplateRequest struct {
	TemplateID uint                   `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
	TargetType string                 `json:"target_type"` // all, device, user
	TargetIDs  []string               `json:"target_ids"`
	CreatedBy  string                 `json:"created_by"`
}

// RegisterRoutes 注册通知相关路由
func (c *NotificationController) RegisterRoutes(rg *gin.RouterGroup) {
	notifications := rg.Group("/notifications")
	{
		notifications.GET("", c.ListNotifications)
		notifications.POST("/push/from-template", c.SendNotificationFromTemplate)
	}

	notificationTemplates := rg.Group("/notification-templates")
	{
		notificationTemplates.GET("", c.ListTemplates)
		notificationTemplates.POST("", c.CreateTemplate)
		notificationTemplates.PUT("/:id", c.UpdateTemplate)
		notificationTemplates.DELETE("/:id", c.DeleteTemplate)
	}

	announcements := rg.Group("/announcements")
	{
		announcements.GET("", c.ListAnnouncements)
		announcements.POST("", c.CreateAnnouncement)
		announcements.PUT("/:id", c.UpdateAnnouncement)
		announcements.DELETE("/:id", c.DeleteAnnouncement)
		announcements.POST("/:id/publish", c.PublishAnnouncement)
	}
}

// PushNotification 发送推送通知
// POST /api/v1/notifications/push
func (c *NotificationController) PushNotification(ctx *gin.Context) {
	var req struct {
		Title      string   `json:"title" binding:"required"`
		Content    string   `json:"content" binding:"required"`
		TargetType string   `json:"target_type" binding:"required"` // all, device, user
		TargetIDs  []string `json:"target_ids"`
		CreatedBy  string   `json:"created_by"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now()

	// 收集所有目标设备ID
	var deviceIDs []string
	if req.TargetType == "all" {
		var devices []models.Device
		c.DB.Select("device_id").Find(&devices)
		for _, d := range devices {
			deviceIDs = append(deviceIDs, d.DeviceID)
		}
	} else {
		deviceIDs = req.TargetIDs
	}

	// 批量创建通知记录（避免 N+1 查询）
	if len(deviceIDs) > 0 {
		notifications := make([]models.Notification, 0, len(deviceIDs))
		for _, deviceID := range deviceIDs {
			notifications = append(notifications, models.Notification{
				DeviceID:  deviceID,
				Title:     req.Title,
				Content:   req.Content,
				Priority:  1,
				Channel:   "push",
				Status:    "sent",
				SentAt:    &now,
				CreatedBy: req.CreatedBy,
			})
		}
		// 批量插入
		if err := c.DB.Create(&notifications).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发送失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"sent_count": len(deviceIDs),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"sent_count": 0,
		},
	})
}

// BatchDeleteNotifications 批量删除通知
// POST /api/v1/notifications/batch-delete
func (c *NotificationController) BatchDeleteNotifications(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供要删除的通知ID列表"})
		return
	}

	if err := c.DB.Delete(&models.Notification{}, req.IDs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// WithdrawAnnouncement 撤回公告
// POST /api/v1/announcements/:id/withdraw
func (c *NotificationController) WithdrawAnnouncement(ctx *gin.Context) {
	id := ctx.Param("id")

	var ann models.Announcement
	if err := c.DB.First(&ann, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "公告不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if ann.Status != "published" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只能撤回已发布的公告"})
		return
	}

	updates := map[string]interface{}{
		"status": "withdrawn",
	}

	if err := c.DB.Model(&ann).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "撤回失败"})
		return
	}

	c.DB.First(&ann, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": ann})
}

// GetAnnouncement 获取公告详情
// GET /api/v1/announcements/:id
func (c *NotificationController) GetAnnouncement(ctx *gin.Context) {
	id := ctx.Param("id")

	var ann models.Announcement
	if err := c.DB.First(&ann, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "公告不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": ann})
}

// ReplaceTemplateVariables 替换模板变量，格式 {{variable}}
func ReplaceTemplateVariables(template string, vars map[string]interface{}) string {
	if template == "" || len(vars) == 0 {
		return template
	}

	// 匹配 {{variable}} 格式
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		// 提取变量名
		varName := match[2 : len(match)-2]
		if val, ok := vars[varName]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match // 变量不存在保持原样
	})

	return result
}
