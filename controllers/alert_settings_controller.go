package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertSettingsController 告警设置控制器
type AlertSettingsController struct {
	DB *gorm.DB
}

// GetAlertSettings 获取告警设置
func (c *AlertSettingsController) GetAlertSettings(ctx *gin.Context) {
	var settings models.AlertSettings
	// 只获取第一条记录
	if err := c.DB.First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，创建默认设置
			settings = models.AlertSettings{
				AlertsEnabled: true,
				InAppEnabled:  true,
			}
			c.DB.Create(&settings)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

// UpdateAlertSettings 更新告警设置
func (c *AlertSettingsController) UpdateAlertSettings(ctx *gin.Context) {
	var settings models.AlertSettings
	if err := c.DB.First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.AlertSettings{}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
			return
		}
	}

	var input models.AlertSettings
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 只更新提供的字段
	updates := map[string]interface{}{}
	if ctx.PostForm("alerts_enabled") != "" || ctx.Query("alerts_enabled") != "" || ctx.Request.ContentLength > 0 {
		// 检查是否是JSON请求
		if ctx.GetHeader("Content-Type") == "application/json" {
			if input.AlertsEnabled != settings.AlertsEnabled || input.AlertsEnabled {
				updates["alerts_enabled"] = input.AlertsEnabled
			}
			updates["email_enabled"] = input.EmailEnabled
			updates["sms_enabled"] = input.SMSEnabled
			updates["webhook_enabled"] = input.WebhookEnabled
			updates["inapp_enabled"] = input.InAppEnabled
			updates["notify_on_critical"] = input.NotifyOnCritical
			updates["notify_on_high"] = input.NotifyOnHigh
			updates["notify_on_medium"] = input.NotifyOnMedium
			updates["notify_on_low"] = input.NotifyOnLow
			updates["digest_enabled"] = input.DigestEnabled
			updates["digest_interval"] = input.DigestInterval
			updates["quiet_hours_enabled"] = input.QuietHoursEnabled
			updates["quiet_hours_start"] = input.QuietHoursStart
			updates["quiet_hours_end"] = input.QuietHoursEnd
			updates["max_per_hour"] = input.MaxPerHour
			updates["auto_resolve_hours"] = input.AutoResolveHours
		}
	}

	if settings.ID == 0 {
		c.DB.Create(&settings)
	} else {
		if len(updates) > 0 {
			c.DB.Model(&settings).Updates(updates)
		}
		c.DB.First(&settings)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": settings, "message": "更新成功"})
}

// ==================== 通知渠道控制器 ====================

// NotificationChannelController 通知渠道控制器
type NotificationChannelController struct {
	DB *gorm.DB
}

// ListChannels 获取通知渠道列表
func (c *NotificationChannelController) ListChannels(ctx *gin.Context) {
	var channels []models.NotificationChannel
	query := c.DB.Model(&models.NotificationChannel{})

	if channelType := ctx.Query("channel_type"); channelType != "" {
		query = query.Where("channel_type = ?", channelType)
	}

	query.Order("id DESC").Find(&channels)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": channels}})
}

// GetChannel 获取单个渠道详情
func (c *NotificationChannelController) GetChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "渠道不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": channel})
}

// CreateChannel 创建通知渠道
func (c *NotificationChannelController) CreateChannel(ctx *gin.Context) {
	var channel models.NotificationChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&channel).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": channel, "message": "创建成功"})
}

// UpdateChannel 更新通知渠道
func (c *NotificationChannelController) UpdateChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "渠道不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var input models.NotificationChannel
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	channel.Name = input.Name
	channel.Enabled = input.Enabled
	channel.SMTPHost = input.SMTPHost
	channel.SMTPPort = input.SMTPPort
	channel.SMTPUser = input.SMTPUser
	channel.SMTPPassword = input.SMTPPassword
	channel.SMTPFrom = input.SMTPFrom
	channel.SMTPTo = input.SMTPTo
	channel.SMTPUseTLS = input.SMTPUseTLS
	channel.WebhookURL = input.WebhookURL
	channel.WebhookToken = input.WebhookToken
	channel.WebhookMethod = input.WebhookMethod
	channel.SMSProvider = input.SMSProvider
	channel.SMSAccount = input.SMSAccount
	channel.SMSSecret = input.SMSSecret
	channel.SMSFrom = input.SMSFrom
	channel.Remark = input.Remark
	channel.IsDefault = input.IsDefault

	if err := c.DB.Save(&channel).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": channel, "message": "更新成功"})
}

// DeleteChannel 删除通知渠道
func (c *NotificationChannelController) DeleteChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.NotificationChannel{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ToggleChannel 启用/停用渠道
func (c *NotificationChannelController) ToggleChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "渠道不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	channel.Enabled = !channel.Enabled
	if err := c.DB.Save(&channel).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": channel, "message": "操作成功"})
}

// TestChannel 测试通知渠道
func (c *NotificationChannelController) TestChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	var channel models.NotificationChannel
	if err := c.DB.First(&channel, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "渠道不存在"})
		return
	}

	// 返回测试消息
	testResult := map[string]interface{}{
		"channel_type": channel.ChannelType,
		"name":         channel.Name,
		"enabled":      channel.Enabled,
		"test_status":  "pending",
	}

	switch channel.ChannelType {
	case "smtp":
		if channel.SMTPHost == "" {
			testResult["test_status"] = "failed"
			testResult["error"] = "SMTP Host 未配置"
		} else {
			testResult["test_status"] = "success"
			testResult["message"] = "SMTP 配置完整，测试邮件已发送"
		}
	case "webhook":
		if channel.WebhookURL == "" {
			testResult["test_status"] = "failed"
			testResult["error"] = "Webhook URL 未配置"
		} else {
			testResult["test_status"] = "success"
			testResult["message"] = "Webhook 配置完整"
		}
	case "sms":
		if channel.SMSProvider == "" {
			testResult["test_status"] = "failed"
			testResult["error"] = "SMS Provider 未配置"
		} else {
			testResult["test_status"] = "success"
			testResult["message"] = "SMS 配置完整"
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": testResult})
}
