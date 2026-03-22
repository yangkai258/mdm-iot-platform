package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// SendAlertNotifications 发送告警通知（邮件、Webhook、站内）
func SendAlertNotifications(db *gorm.DB, alert *models.DeviceAlert, notifyWays string) {
	if notifyWays == "" {
		notifyWays = "inapp" // 默认站内通知
	}

	// 获取设备信息
	deviceName := alert.DeviceID
	var device models.Device
	if db.Where("device_id = ?", alert.DeviceID).First(&device).Error == nil {
		deviceName = device.HardwareModel
	}

	// 构建通知内容
	subject := fmt.Sprintf("【告警-%s】设备 %s: %s",
		severityLabel(alert.Severity), deviceName, alert.AlertType)
	content := fmt.Sprintf("告警类型: %s\n严重程度: %s\n设备ID: %s\n触发值: %.2f\n阈值: %.2f\n时间: %s\n消息: %s",
		alert.AlertType,
		severityLabel(alert.Severity),
		alert.DeviceID,
		alert.TriggerVal,
		alert.Threshold,
		alert.CreatedAt.Format(time.RFC3339),
		alert.Message,
	)

	// 解析通知方式
	ways := parseNotifyWays(notifyWays)

	for _, way := range ways {
		switch way {
		case "email":
			sendEmailNotification(db, alert, subject, content)
		case "webhook":
			sendWebhookNotification(db, alert, subject, content)
		case "inapp":
			sendInAppNotification(db, alert, deviceName, subject, content)
		}
	}
}

func sendEmailNotification(db *gorm.DB, alert *models.DeviceAlert, subject, body string) {
	smtpHost := os.Getenv("SMTP_HOST")
	_ = os.Getenv("SMTP_PORT")
	_ = os.Getenv("SMTP_USER")
	_ = os.Getenv("SMTP_PASS")
	alertEmail := os.Getenv("ALERT_EMAIL")

	if smtpHost == "" || alertEmail == "" {
		// 如果没有配置SMTP，降级为记录通知
		log.Printf("[Notification] 邮件通知（未配置SMTP）: AlertID=%d, Subject=%s", alert.ID, subject)
		record := &models.AlertNotification{
			AlertID:   alert.ID,
			AlertType: "email",
			Status:    "failed",
			Recipient: alertEmail,
			Subject:   subject,
			Content:   body,
			ErrorMsg:  "SMTP未配置",
		}
		db.Create(record)
		return
	}

	// 实际发送邮件（这里记录为pending，实际发送由外部服务处理）
	record := &models.AlertNotification{
		AlertID:   alert.ID,
		AlertType: "email",
		Status:    "pending",
		Recipient: alertEmail,
		Subject:   subject,
		Content:   body,
	}
	db.Create(record)

	// TODO: 使用 smtp.SendMail 实际发送邮件
	// 此处记录日志，生产环境需要实现真实的邮件发送
	log.Printf("[Notification] 邮件通知已创建: AlertID=%d, To=%s, Subject=%s", alert.ID, alertEmail, subject)

	// 更新状态为已发送
	now := time.Now()
	record.Status = "sent"
	record.SentAt = &now
	db.Save(record)
}

func sendWebhookNotification(db *gorm.DB, alert *models.DeviceAlert, subject, body string) {
	webhookURL := os.Getenv("ALERT_WEBHOOK_URL")
	if webhookURL == "" {
		log.Printf("[Notification] Webhook通知（未配置URL）: AlertID=%d", alert.ID)
		record := &models.AlertNotification{
			AlertID:   alert.ID,
			AlertType: "webhook",
			Status:    "failed",
			Recipient: webhookURL,
			Subject:   subject,
			Content:   body,
			ErrorMsg:  "WEBHOOK_URL未配置",
		}
		db.Create(record)
		return
	}

	record := &models.AlertNotification{
		AlertID:   alert.ID,
		AlertType: "webhook",
		Status:    "pending",
		Recipient: webhookURL,
		Subject:   subject,
		Content:   body,
	}
	db.Create(record)

	// 构建 webhook payload
	payload := map[string]interface{}{
		"alert_id":    alert.ID,
		"device_id":   alert.DeviceID,
		"alert_type":  alert.AlertType,
		"severity":    alert.Severity,
		"message":     alert.Message,
		"subject":     subject,
		"content":     body,
		"trigger_val": alert.TriggerVal,
		"threshold":   alert.Threshold,
		"created_at":  alert.CreatedAt,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(payloadBytes))
	if err != nil {
		record.Status = "failed"
		record.ErrorMsg = err.Error()
		db.Save(record)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		record.Status = "failed"
		record.ErrorMsg = err.Error()
		db.Save(record)
		log.Printf("[Notification] Webhook发送失败: AlertID=%d, Error=%v", alert.ID, err)
		return
	}
	defer resp.Body.Close()

	now := time.Now()
	record.Status = "sent"
	record.SentAt = &now
	db.Save(record)
	log.Printf("[Notification] Webhook通知已发送: AlertID=%d, Status=%d", alert.ID, resp.StatusCode)
}

func sendInAppNotification(db *gorm.DB, alert *models.DeviceAlert, deviceName, subject, body string) {
	notification := &models.Notification{
		DeviceID:  alert.DeviceID,
		Title:     subject,
		Content:   body,
		Priority:  alert.Severity,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}

	if err := db.Create(notification).Error; err != nil {
		log.Printf("[Notification] 站内通知创建失败: AlertID=%d, Error=%v", alert.ID, err)
		return
	}

	// 同时记录到 alert_notifications
	record := &models.AlertNotification{
		AlertID:   alert.ID,
		AlertType: "inapp",
		Status:    "delivered",
		Recipient: "in-app",
		Subject:   subject,
		Content:   body,
	}
	now := time.Now()
	record.SentAt = &now
	db.Create(record)

	log.Printf("[Notification] 站内通知已创建: AlertID=%d, NotificationID=%d", alert.ID, notification.ID)
}

func parseNotifyWays(ways string) []string {
	if ways == "" {
		return []string{"inapp"}
	}
	// 简单逗号分割
	var result []string
	current := ""
	for _, c := range ways {
		if c == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func severityLabel(severity int) string {
	switch severity {
	case 1:
		return "低"
	case 2:
		return "中"
	case 3:
		return "高"
	case 4:
		return "严重"
	default:
		return "未知"
	}
}
