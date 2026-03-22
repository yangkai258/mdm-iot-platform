package notification

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WebhookService Webhook 通知服务
type WebhookService struct {
	client *http.Client
}

// NewWebhookService 创建 Webhook 服务实例
func NewWebhookService() *WebhookService {
	return &WebhookService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WebhookPayload Webhook 载荷
type WebhookPayload struct {
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp string                 `json:"timestamp"`
}

// Send 发送 Webhook 请求
func (s *WebhookService) Send(url string, payload []byte, secret string) error {
	if url == "" {
		return fmt.Errorf("webhook URL 不能为空")
	}

	reqFinal, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	reqFinal.Header.Set("Content-Type", "application/json")
	reqFinal.Header.Set("User-Agent", "MDM-Alert-Webhook/1.0")
	if secret != "" {
		signature := s.SignPayload(payload, secret)
		reqFinal.Header.Set("X-Webhook-Signature", "sha256="+signature)
		reqFinal.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	}

	resp, err := s.client.Do(reqFinal)
	if err != nil {
		return fmt.Errorf("webhook 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook 返回错误状态码: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SignPayload 使用 HMAC-SHA256 签名 payload
func (s *WebhookService) SignPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
