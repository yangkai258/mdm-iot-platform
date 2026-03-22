package notification

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status     string    `json:"status"`      // healthy / unhealthy / unknown
	Message    string    `json:"message"`
	LatencyMs  int64     `json:"latency_ms"`
	CheckedAt  time.Time `json:"checked_at"`
}

// ChannelHealth 健康检查结果
type ChannelHealth struct {
	ChannelID   uint         `json:"channel_id"`
	ChannelType string       `json:"channel_type"`
	ChannelName string       `json:"channel_name"`
	Health      HealthStatus `json:"health"`
}

// HealthChecker 健康检查器
type HealthChecker struct {
	emailService *EmailService
	smsService   *SMSService
	client       *http.Client
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(email *EmailService, sms *SMSService) *HealthChecker {
	return &HealthChecker{
		emailService: email,
		smsService:   sms,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// CheckEmailHealth 检查邮件服务健康状态
func (h *HealthChecker) CheckEmailHealth() HealthStatus {
	if h.emailService == nil {
		return HealthStatus{
			Status:    "unknown",
			Message:   "邮件服务未配置",
			CheckedAt: time.Now(),
		}
	}

	start := time.Now()
	err := h.emailService.TestConnection()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   err.Error(),
			LatencyMs: latency,
			CheckedAt: time.Now(),
		}
	}

	return HealthStatus{
		Status:    "healthy",
		Message:   "SMTP 连接正常",
		LatencyMs: latency,
		CheckedAt: time.Now(),
	}
}

// CheckSMSHealth 检查短信服务健康状态
func (h *HealthChecker) CheckSMSHealth() HealthStatus {
	if h.smsService == nil {
		return HealthStatus{
			Status:    "unknown",
			Message:   "短信服务未配置",
			CheckedAt: time.Now(),
		}
	}

	return HealthStatus{
		Status:    "healthy",
		Message:   fmt.Sprintf("短信服务商: %s", h.smsService.provider),
		LatencyMs: 0,
		CheckedAt: time.Now(),
	}
}

// CheckWebhookHealth 检查 Webhook URL 可达性
func (h *HealthChecker) CheckWebhookHealth(webhookURL string) HealthStatus {
	if webhookURL == "" {
		return HealthStatus{
			Status:    "unknown",
			Message:   "Webhook URL 未配置",
			CheckedAt: time.Now(),
		}
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", webhookURL, nil)
	if err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("请求创建失败: %v", err),
			CheckedAt: time.Now(),
		}
	}

	resp, err := h.client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return HealthStatus{
				Status:    "unhealthy",
				Message:   "Webhook 连接超时",
				LatencyMs: latency,
				CheckedAt: time.Now(),
			}
		}
		return HealthStatus{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("连接失败: %v", err),
			LatencyMs: latency,
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		return HealthStatus{
			Status:    "healthy",
			Message:   fmt.Sprintf("HTTP %d", resp.StatusCode),
			LatencyMs: latency,
			CheckedAt: time.Now(),
		}
	}

	return HealthStatus{
		Status:    "unhealthy",
		Message:   fmt.Sprintf("Webhook 返回 HTTP %d", resp.StatusCode),
		LatencyMs: latency,
		CheckedAt: time.Now(),
	}
}

// CheckChannel 检查指定渠道的健康状态
func (h *HealthChecker) CheckChannel(channelType string, cfg EmailConfig) HealthStatus {
	switch channelType {
	case "smtp", "email":
		emailSvc := NewEmailService(cfg)
		start := time.Now()
		err := emailSvc.TestConnection()
		latency := time.Since(start).Milliseconds()
		if err != nil {
			return HealthStatus{
				Status:    "unhealthy",
				Message:   err.Error(),
				LatencyMs: latency,
				CheckedAt: time.Now(),
			}
		}
		return HealthStatus{
			Status:    "healthy",
			Message:   "SMTP 连接正常",
			LatencyMs: latency,
			CheckedAt: time.Now(),
		}
	case "sms":
		return h.CheckSMSHealth()
	case "webhook":
		return h.CheckWebhookHealth(cfg.Host) // Host 字段存储 webhook URL
	}
	return HealthStatus{
		Status:    "unknown",
		Message:   "未知渠道类型",
		CheckedAt: time.Now(),
	}
}

// GlobalHealthChecker 全局健康检查器实例
var GlobalHealthChecker *HealthChecker

func init() {
	GlobalHealthChecker = NewHealthChecker(nil, nil)
}
