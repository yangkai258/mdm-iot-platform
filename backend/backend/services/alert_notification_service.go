package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// ===== SMTP 发送告警邮件 =====

// SendAlertEmail 通过 SMTP 发送告警邮件
// 环境变量：SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASSWORD, SMTP_FROM, SMTP_USE_TLS
func SendAlertEmail(db *gorm.DB, alert *models.DeviceAlert, toEmail, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")
	useTLS := os.Getenv("SMTP_USE_TLS")

	if smtpHost == "" {
		return fmt.Errorf("SMTP_HOST not configured")
	}
	if toEmail == "" {
		return fmt.Errorf("recipient email is empty")
	}

	// 默认端口
	if smtpPort == "" {
		smtpPort = "587"
	}
	// 默认发件人
	if smtpFrom == "" {
		smtpFrom = smtpUser
	}

	// 构建邮件内容
	mime := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"From: " + smtpFrom + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// SMTP 地址
	smtpAddr := smtpHost + ":" + smtpPort

	// 记录通知
	record := &models.AlertNotification{
		AlertID:   alert.ID,
		AlertType: "email",
		Status:    "pending",
		Recipient: toEmail,
		Subject:   subject,
		Content:   body,
	}
	db.Create(record)

	var errMsg string
	var sentAt *time.Time

	if useTLS == "true" || smtpPort == "465" {
		// TLS/SSL 模式
		err := sendMailTLS(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, toEmail, subject, body)
		if err != nil {
			errMsg = err.Error()
			record.Status = "failed"
			record.ErrorMsg = errMsg
			db.Save(record)
			log.Printf("[AlertNotification] SMTP TLS send failed: %v", err)
			return err
		}
	} else {
		// STARTTLS / 普通模式
		auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
		err := smtp.SendMail(smtpAddr, auth, smtpFrom, []string{toEmail}, []byte(mime))
		if err != nil {
			errMsg = err.Error()
			record.Status = "failed"
			record.ErrorMsg = errMsg
			db.Save(record)
			log.Printf("[AlertNotification] SMTP send failed: %v", err)
			return err
		}
	}

	now := time.Now()
	sentAt = &now
	record.Status = "sent"
	record.SentAt = sentAt
	db.Save(record)

	log.Printf("[AlertNotification] SMTP email sent: AlertID=%d, To=%s, Subject=%s", alert.ID, toEmail, subject)
	return nil
}

// sendMailTLS 使用 TLS 发送邮件
func sendMailTLS(host, port, user, pass, from, to, subject, body string) error {
	addr := host + ":" + port

	auth := smtp.PlainAuth("", user, pass, host)

	msg := buildMimeEmail(from, to, subject, body)

	err := smtp.SendMail(addr, auth, user, []string{to}, []byte(msg))
	return err
}

// buildMimeEmail 构建标准 MIME 邮件
func buildMimeEmail(from, to, subject, body string) string {
	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	sb.WriteString("\r\n")
	return sb.String()
}

// ===== Webhook POST 发送告警 =====

// SendAlertWebhook 通过 Webhook POST 发送告警到外部系统
// 环境变量：WEBHOOK_URL, WEBHOOK_TOKEN
func SendAlertWebhook(db *gorm.DB, alert *models.DeviceAlert, subject, body string) error {
	webhookURL := os.Getenv("WEBHOOK_URL")
	webhookToken := os.Getenv("WEBHOOK_TOKEN")

	if webhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL not configured")
	}

	// 构建 payload
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
		"created_at":  alert.CreatedAt.Format(time.RFC3339),
		"status":      alert.Status,
	}

	payloadBytes, _ := json.Marshal(payload)

	// 创建请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MDM-AlertService/1.0")

	// 如果配置了 Token，加入 Authorization 头
	if webhookToken != "" {
		req.Header.Set("Authorization", "Bearer "+webhookToken)
	}

	// 发送请求，10 秒超时
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// 记录失败
		record := &models.AlertNotification{
			AlertID:   alert.ID,
			AlertType: "webhook",
			Status:    "failed",
			Recipient: webhookURL,
			Subject:   subject,
			Content:   body,
			ErrorMsg:  err.Error(),
		}
		db.Create(record)
		log.Printf("[AlertNotification] Webhook send failed: AlertID=%d, Error=%v", alert.ID, err)
		return fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体（最多 1KB）
	buf := make([]byte, 1024)
	respLen := 0
	if n, _ := resp.Body.Read(buf); n > 0 {
		respLen = n
	}
	respBody := string(buf[:respLen])

	// 记录结果
	now := time.Now()
	status := "sent"
	errMsg := ""
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = "failed"
		errMsg = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, respBody)
	}

	record := &models.AlertNotification{
		AlertID:   alert.ID,
		AlertType: "webhook",
		Status:    status,
		Recipient: webhookURL,
		Subject:   subject,
		Content:   body,
		ErrorMsg:  errMsg,
		SentAt:    &now,
	}
	db.Create(record)

	if status == "failed" {
		log.Printf("[AlertNotification] Webhook failed: AlertID=%d, HTTP=%d, Body=%s", alert.ID, resp.StatusCode, respBody)
		return fmt.Errorf("webhook returned HTTP %d: %s", resp.StatusCode, respBody)
	}

	log.Printf("[AlertNotification] Webhook sent: AlertID=%d, HTTP=%d", alert.ID, resp.StatusCode)
	return nil
}

// ===== 综合告警通知入口 =====

// SendAlertViaSMTPOrWebhook 发送告警通知（SMTP 或 Webhook）
// channel: "email" 或 "webhook"
func SendAlertViaSMTPOrWebhook(db *gorm.DB, alert *models.DeviceAlert, channel, to, subject, body string) {
	var err error
	switch channel {
	case "email":
		err = SendAlertEmail(db, alert, to, subject, body)
	case "webhook":
		err = SendAlertWebhook(db, alert, subject, body)
	default:
		log.Printf("[AlertNotification] Unknown channel: %s", channel)
		return
	}

	if err != nil {
		log.Printf("[AlertNotification] Failed to send via %s: AlertID=%d, Error=%v", channel, alert.ID, err)
	}
}
