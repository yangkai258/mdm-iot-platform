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
	"mdm-backend/notification"
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

	// Sprint 11: /api/v1/notification/* 路由
	notificationGroup := rg.Group("/notification")
	{
		// 通知渠道 CRUD
		notificationGroup.GET("/channels", c.ListNotificationChannels)
		notificationGroup.POST("/channels", c.CreateNotificationChannel)
		notificationGroup.GET("/channels/:id", c.GetNotificationChannel)
		notificationGroup.PUT("/channels/:id", c.UpdateNotificationChannel)
		notificationGroup.DELETE("/channels/:id", c.DeleteNotificationChannel)
		notificationGroup.POST("/channels/:id/test", c.TestNotificationChannel)

		// 通知日志
		notificationGroup.GET("/logs", c.ListNotificationLogs)

		// 通知统计
		notificationGroup.GET("/stats", c.GetNotificationStats)

		// 通知模板 CRUD
		notificationGroup.GET("/templates", c.ListTemplates)
		notificationGroup.POST("/templates", c.CreateNotificationTemplate)
		notificationGroup.PUT("/templates/:id", c.UpdateNotificationTemplate)
		notificationGroup.DELETE("/templates/:id", c.DeleteNotificationTemplate)
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

// ==================== Sprint 11: 通知渠道 API (/api/v1/notification/*) ====================

// CreateNotificationChannelRequest 创建通知渠道请求
type CreateNotificationChannelRequest struct {
	ChannelType string `json:"channel_type" binding:"required"` // smtp/webhook/sms
	ChannelName string `json:"channel_name" binding:"required"`
	Config      map[string]interface{} `json:"config"`
	Enabled     bool `json:"enabled"`
	IsDefault   bool `json:"is_default"`
	Priority    int  `json:"priority"`
}

// UpdateNotificationChannelRequest 更新通知渠道请求
type UpdateNotificationChannelRequest struct {
	ChannelName string `json:"channel_name"`
	Config      map[string]interface{} `json:"config"`
	Enabled     *bool `json:"enabled"`
	IsDefault   *bool `json:"is_default"`
	Priority    *int  `json:"priority"`
}

// ListNotificationChannels 获取通知渠道列表
// GET /api/v1/notification/channels
func (c *NotificationController) ListNotificationChannels(ctx *gin.Context) {
	var channels []models.NotificationChannel
	query := c.DB.Model(&models.NotificationChannel{})

	if channelType := ctx.Query("channel_type"); channelType != "" {
		query = query.Where("channel_type = ?", channelType)
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

	query.Order("priority DESC, id DESC").Offset(offset).Limit(pageSize).Find(&channels)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      channels,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateNotificationChannel 创建通知渠道
// POST /api/v1/notification/channels
func (c *NotificationController) CreateNotificationChannel(ctx *gin.Context) {
	var req CreateNotificationChannelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 验证渠道类型
	validTypes := map[string]bool{"smtp": true, "webhook": true, "sms": true}
	if !validTypes[req.ChannelType] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "不支持的渠道类型，仅支持 smtp/webhook/sms",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 如果设置为默认，取消其他默认
	if req.IsDefault {
		c.DB.Model(&models.NotificationChannel{}).Where("is_default = ?", true).
			Updates(map[string]interface{}{"is_default": false})
	}

	// 将 config map 转为 JSON（备用，如果后续需要存储原始配置）
	_, _ = json.Marshal(req.Config)

	channel := models.NotificationChannel{
		ChannelType: req.ChannelType,
		Name:        req.ChannelName,
		Enabled:     req.Enabled,
		IsDefault:   req.IsDefault,
		Priority:    req.Priority,
	}

	// 根据渠道类型设置对应字段
	c.applyChannelConfig(&channel, req.ChannelType, req.Config)

	if err := c.DB.Create(&channel).Error; err != nil {
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
		"data":    channel,
	})
}

// GetNotificationChannel 获取单个通知渠道
// GET /api/v1/notification/channels/:id
func (c *NotificationController) GetNotificationChannel(ctx *gin.Context) {
	id := ctx.Param("id")

	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "渠道不存在",
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
		"data":    channel,
	})
}

// UpdateNotificationChannel 更新通知渠道
// PUT /api/v1/notification/channels/:id
func (c *NotificationController) UpdateNotificationChannel(ctx *gin.Context) {
	id := ctx.Param("id")

	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "渠道不存在",
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

	var req UpdateNotificationChannelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.ChannelName != "" {
		updates["name"] = req.ChannelName
	}
	if req.Config != nil {
		configJSON, _ := json.Marshal(req.Config)
		updates["config"] = configJSON
		c.applyChannelConfig(&channel, channel.ChannelType, req.Config)
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
		channel.Enabled = *req.Enabled
	}
	if req.IsDefault != nil {
		if *req.IsDefault {
			c.DB.Model(&models.NotificationChannel{}).Where("is_default = ?", true).
				Updates(map[string]interface{}{"is_default": false})
		}
		updates["is_default"] = *req.IsDefault
		channel.IsDefault = *req.IsDefault
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
		channel.Priority = *req.Priority
	}

	if err := c.DB.Model(&channel).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "更新失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	c.DB.First(&channel, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    channel,
	})
}

// DeleteNotificationChannel 删除通知渠道
// DELETE /api/v1/notification/channels/:id
func (c *NotificationController) DeleteNotificationChannel(ctx *gin.Context) {
	id := ctx.Param("id")

	result := c.DB.Delete(&models.NotificationChannel{}, id)
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
			"message":    "渠道不存在",
			"error_code": "ERR_NOT_FOUND",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// TestNotificationChannel 测试通知渠道
// POST /api/v1/notification/channels/:id/test
func (c *NotificationController) TestNotificationChannel(ctx *gin.Context) {
	id := ctx.Param("id")

	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4004,
				"message":    "渠道不存在",
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

	// 根据渠道类型执行健康检查
	var healthStatus string
	var healthMsg string

	switch channel.ChannelType {
	case "smtp":
		emailCfg := notification.EmailConfig{
			Host:     channel.SMTPHost,
			Port:     channel.SMTPPort,
			Username: channel.SMTPUser,
			Password: channel.SMTPPassword,
			From:     channel.SMTPFrom,
		}
		emailSvc := notification.NewEmailService(emailCfg, c.DB)
		err := emailSvc.TestConnection()
		if err != nil {
			healthStatus = "unhealthy"
			healthMsg = err.Error()
		} else {
			healthStatus = "healthy"
			healthMsg = "SMTP 连接成功"
		}
	case "webhook":
		webhookURL := channel.WebhookURL
		if webhookURL == "" {
			healthStatus = "unhealthy"
			healthMsg = "Webhook URL 未配置"
		} else {
			webhookSvc := notification.NewWebhookService()
			// 发送测试 payload
			testPayload := []byte(`{"test": true, "message": "MDM test notification"}`)
			err := webhookSvc.Send(webhookURL, testPayload, channel.WebhookToken)
			if err != nil {
				healthStatus = "unhealthy"
				healthMsg = err.Error()
			} else {
				healthStatus = "healthy"
				healthMsg = "Webhook 发送成功"
			}
		}
	case "sms":
		smsCfg := notification.SMSConfig{
			Provider:  channel.SMSProvider,
			AccessKey: channel.SMSAccount,
			SecretKey: channel.SMSSecret,
			SignName:  channel.SMSFrom,
		}
		smsSvc := notification.NewSMSService(smsCfg)
		// 测试手机号从查询参数获取
		testPhone := ctx.Query("test_phone")
		if testPhone != "" {
			err := smsSvc.Send([]string{testPhone}, "SMS_TEST", nil)
			if err != nil {
				healthStatus = "unhealthy"
				healthMsg = err.Error()
			} else {
				healthStatus = "healthy"
				healthMsg = "SMS 发送成功"
			}
		} else {
			healthStatus = "unknown"
			healthMsg = "请通过 ?test_phone= 参数提供测试手机号"
		}
	default:
		healthStatus = "unknown"
		healthMsg = "未知渠道类型"
	}

	now := time.Now()
	c.DB.Model(&channel).Updates(map[string]interface{}{
		"health_status":   healthStatus,
		"last_checked_at": &now,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"status":  healthStatus,
			"message": healthMsg,
		},
	})
}

// ListNotificationLogs 获取通知日志
// GET /api/v1/notification/logs
func (c *NotificationController) ListNotificationLogs(ctx *gin.Context) {
	var logs []models.NotificationLog
	query := c.DB.Model(&models.NotificationLog{})

	if channelID := ctx.Query("channel_id"); channelID != "" {
		query = query.Where("channel_id = ?", channelID)
	}
	if channelType := ctx.Query("channel_type"); channelType != "" {
		query = query.Where("channel_type = ?", channelType)
	}
	if alertID := ctx.Query("alert_id"); alertID != "" {
		query = query.Where("alert_id = ?", alertID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("recipient LIKE ? OR subject LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
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

	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetNotificationStats 获取通知统计
// GET /api/v1/notification/stats
func (c *NotificationController) GetNotificationStats(ctx *gin.Context) {
	stats := models.NotificationStats{
		ByChannel: make(map[string]int64),
	}

	// 总数统计
	c.DB.Model(&models.NotificationLog{}).Count(&stats.TotalSent)
	c.DB.Model(&models.NotificationLog{}).Where("status = ?", "success").Count(&stats.TotalSuccess)
	c.DB.Model(&models.NotificationLog{}).Where("status = ?", "failed").Count(&stats.TotalFailed)

	// 今日统计
	today := time.Now().Truncate(24 * time.Hour)
	c.DB.Model(&models.NotificationLog{}).Where("created_at >= ?", today).Count(&stats.TodaySent)
	c.DB.Model(&models.NotificationLog{}).Where("status = ? AND created_at >= ?", "success", today).Count(&stats.TodaySuccess)
	c.DB.Model(&models.NotificationLog{}).Where("status = ? AND created_at >= ?", "failed", today).Count(&stats.TodayFailed)

	// 按渠道统计
	var channelStats []struct {
		ChannelType string
		Count       int64
	}
	c.DB.Model(&models.NotificationLog{}).
		Select("channel_type, COUNT(*) as count").
		Group("channel_type").
		Scan(&channelStats)

	for _, cs := range channelStats {
		stats.ByChannel[cs.ChannelType] = cs.Count
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// ==================== Sprint 11: 通知模板 API (/api/v1/notification/templates) ====================

// CreateNotificationTemplate 创建通知模板
// POST /api/v1/notification/templates
func (c *NotificationController) CreateNotificationTemplate(ctx *gin.Context) {
	var req models.NotificationTemplate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查编码唯一性
	var exists models.NotificationTemplate
	if err := c.DB.Where("code = ?", req.Code).First(&exists).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4006,
			"message":    "模板编码已存在",
			"error_code": "ERR_DUPLICATE",
		})
		return
	}

	if err := c.DB.Create(&req).Error; err != nil {
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
		"data":    req,
	})
}

// UpdateNotificationTemplate 更新通知模板
// PUT /api/v1/notification/templates/:id
func (c *NotificationController) UpdateNotificationTemplate(ctx *gin.Context) {
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

	delete(updates, "code") // 不允许修改 code

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

// DeleteNotificationTemplate 删除通知模板
// DELETE /api/v1/notification/templates/:id
func (c *NotificationController) DeleteNotificationTemplate(ctx *gin.Context) {
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

// ==================== 辅助函数 ====================

// applyChannelConfig 根据渠道类型应用配置到 channel 模型
func (c *NotificationController) applyChannelConfig(channel *models.NotificationChannel, channelType string, config map[string]interface{}) {
	switch channelType {
	case "smtp":
		if v, ok := config["smtp_host"].(string); ok {
			channel.SMTPHost = v
		}
		if v, ok := config["smtp_port"].(float64); ok {
			channel.SMTPPort = int(v)
		}
		if v, ok := config["smtp_user"].(string); ok {
			channel.SMTPUser = v
		}
		if v, ok := config["smtp_password"].(string); ok {
			channel.SMTPPassword = v
		}
		if v, ok := config["smtp_from"].(string); ok {
			channel.SMTPFrom = v
		}
		if v, ok := config["smtp_to"].(string); ok {
			channel.SMTPTo = v
		}
		if v, ok := config["smtp_use_tls"].(bool); ok {
			channel.SMTPUseTLS = v
		}
	case "webhook":
		if v, ok := config["webhook_url"].(string); ok {
			channel.WebhookURL = v
		}
		if v, ok := config["webhook_token"].(string); ok {
			channel.WebhookToken = v
		}
		if v, ok := config["webhook_method"].(string); ok {
			channel.WebhookMethod = v
		}
	case "sms":
		if v, ok := config["sms_provider"].(string); ok {
			channel.SMSProvider = v
		}
		if v, ok := config["sms_account"].(string); ok {
			channel.SMSAccount = v
		}
		if v, ok := config["sms_secret"].(string); ok {
			channel.SMSSecret = v
		}
		if v, ok := config["sms_from"].(string); ok {
			channel.SMSFrom = v
		}
	}
}

// getStringField 安全获取 map 中的字符串字段
func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// getIntField 安全获取 map 中的整数字段
func getIntField(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}
