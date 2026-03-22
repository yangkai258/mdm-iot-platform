package notification

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WebhookService Webhook 服务
type WebhookService struct {
	Client *http.Client
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send 发送 Webhook 请求
func (s *WebhookService) Send(targetURL string, payload []byte, secret string) error {
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		signature := s.SignPayload(payload, secret)
		req.Header.Set("X-Webhook-Signature", signature)
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook error: %d - %s", resp.StatusCode, string(body))
	}
	return nil
}

// SignPayload 签名 Webhook payload
func (s *WebhookService) SignPayload(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// WebhookServiceInterface Webhook 服务接口
type WebhookServiceInterface interface {
	Send(targetURL string, payload []byte, secret string) error
	SignPayload(payload []byte, secret string) string
}
