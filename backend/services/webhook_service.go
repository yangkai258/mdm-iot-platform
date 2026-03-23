package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// WebhookService Webhook 服务
type WebhookService struct {
	db *gorm.DB
}

// NewWebhookService 创建 Webhook 服务
func NewWebhookService(db *gorm.DB) *WebhookService {
	return &WebhookService{db: db}
}

// DeliverWebhook 投递 Webhook 事件
func (s *WebhookService) DeliverWebhook(eventID string) error {
	var event models.WebhookEventRecord
	if err := s.db.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	var sub models.WebhookSubscription
	if err := s.db.Where("id = ?", event.SubscriptionID).First(&sub).Error; err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	return s.deliverToURL(&event, &sub)
}

// deliverToURL 投递事件到目标 URL
func (s *WebhookService) deliverToURL(event *models.WebhookEventRecord, sub *models.WebhookSubscription) error {
	// 构建请求
	payloadBytes, _ := json.Marshal(event.Payload)
	reqBody, _ := io.ReadAll(bytes.NewReader(payloadBytes))

	req, err := http.NewRequest("POST", sub.URL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event-ID", event.EventID)
	req.Header.Set("X-Webhook-Event-Type", event.EventType)
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	// 添加签名
	signature := s.signPayload(reqBody, sub.Secret)
	req.Header.Set("X-Webhook-Signature", signature)

	// 添加自定义头
	for k, v := range sub.Headers {
		if vs, ok := v.(string); ok {
			req.Header.Set(k, vs)
		}
	}

	// 发送请求
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		s.saveDelivery(event, sub, resp, duration, err.Error())
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, _ := io.ReadAll(resp.Body)

	// 保存投递记录
	s.saveDelivery(event, sub, resp, duration, "")

	// 更新事件状态
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.db.Model(event).Updates(map[string]interface{}{
			"status":       models.WebhookStatusSuccess,
			"attempts":     event.Attempts + 1,
			"response_code": resp.StatusCode,
			"response_body": string(body),
			"delivered_at": time.Now(),
		})
	} else {
		s.db.Model(event).Updates(map[string]interface{}{
			"status":        models.WebhookStatusFailed,
			"attempts":      event.Attempts + 1,
			"response_code": resp.StatusCode,
			"response_body": string(body),
			"last_error":    fmt.Sprintf("HTTP %d", resp.StatusCode),
		})
	}

	return nil
}

// signPayload 计算 Webhook 签名
func (s *WebhookService) signPayload(body []byte, secret string) string {
	if secret == "" {
		return ""
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

// saveDelivery 保存投递记录
func (s *WebhookService) saveDelivery(event *models.WebhookEventRecord, sub *models.WebhookSubscription, resp *http.Response, durationMs int64, errMsg string) {
	payloadBytes, _ := json.Marshal(event.Payload)
	delivery := models.WebhookDelivery{
		EventID:        event.EventID,
		SubscriptionID: sub.ID,
		URL:            sub.URL,
		RequestHeaders: sub.Headers,
		RequestBody:    string(payloadBytes),
		Attempts:       1,
		DurationMs:     int(durationMs),
	}

	if resp != nil {
		delivery.ResponseCode = resp.StatusCode
		body, _ := io.ReadAll(resp.Body)
		delivery.ResponseBody = string(body)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			delivery.Status = "success"
		} else {
			delivery.Status = "failed"
		}
	} else {
		delivery.Status = "failed"
		delivery.ErrorMessage = errMsg
	}

	s.db.Create(&delivery)
}

// RetryDelivery 重试投递
func (s *WebhookService) RetryDelivery(deliveryID uint) error {
	var delivery models.WebhookDelivery
	if err := s.db.First(&delivery, deliveryID).Error; err != nil {
		return fmt.Errorf("delivery not found: %w", err)
	}

	var event models.WebhookEventRecord
	if err := s.db.Where("event_id = ?", delivery.EventID).First(&event).Error; err != nil {
		return fmt.Errorf("event not found: %w", err)
	}

	var sub models.WebhookSubscription
	if err := s.db.First(&sub, delivery.SubscriptionID).Error; err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	delivery.Attempts++
	s.db.Save(&delivery)

	return s.deliverToURL(&event, &sub)
}

// PublishEvent 发布事件到所有订阅者
func (s *WebhookService) PublishEvent(eventType string, payload models.JSON) error {
	// 生成事件 ID
	eventID := fmt.Sprintf("evt_%d", time.Now().UnixNano())

	// 查找所有订阅了该事件类型的活跃订阅
	var subs []models.WebhookSubscription
	if err := s.db.Where("status = 'active'").Find(&subs).Error; err != nil {
		return err
	}

	for _, sub := range subs {
		// 检查是否订阅了该事件类型
		hasEvent := false
		for _, et := range sub.EventTypes {
			if et == eventType || et == "*" {
				hasEvent = true
				break
			}
		}
		if !hasEvent {
			continue
		}

		// 创建事件记录
		event := models.WebhookEventRecord{
			SubscriptionID: sub.ID,
			EventID:       eventID,
			EventType:     eventType,
			Payload:       payload,
			Status:        models.WebhookStatusPending,
			Attempts:      0,
			MaxAttempts:   sub.RetryCount,
		}
		if err := s.db.Create(&event).Error; err != nil {
			continue
		}

		// 异步投递
		go func(ev models.WebhookEventRecord, sb models.WebhookSubscription) {
			s.DeliverWebhook(ev.EventID)
		}(event, sub)
	}

	return nil
}

// GetWebhookTemplates 获取所有 Webhook 模板
func (s *WebhookService) GetWebhookTemplates() ([]models.WebhookTemplate, error) {
	var templates []models.WebhookTemplate
	err := s.db.Where("is_active = ?", true).Find(&templates).Error
	return templates, err
}

// GetTemplateByID 获取模板详情
func (s *WebhookService) GetTemplateByID(id uint) (*models.WebhookTemplate, error) {
	var tpl models.WebhookTemplate
	if err := s.db.First(&tpl, id).Error; err != nil {
		return nil, err
	}
	return &tpl, nil
}

// InitDefaultTemplates 初始化默认的 Webhook 模板
func (s *WebhookService) InitDefaultTemplates() error {
	templates := []models.WebhookTemplate{
		{
			Name:        "设备告警",
			EventType:   "device.alert",
			Description: "当设备产生告警时触发",
			Category:    "device",
			PayloadExample: models.JSON{
				"device_id":   "device_001",
				"alert_type":  "temperature_high",
				"alert_value": 85.5,
				"threshold":   80.0,
				"timestamp":   time.Now().Unix(),
			},
		},
		{
			Name:        "订阅创建",
			EventType:   "subscription.created",
			Description: "当用户创建新订阅时触发",
			Category:    "subscription",
			PayloadExample: models.JSON{
				"subscription_id": "sub_001",
				"user_id":        123,
				"plan":           "pro",
				"start_time":     time.Now().Unix(),
			},
		},
		{
			Name:        "订阅续费",
			EventType:   "subscription.renewed",
			Description: "当订阅续费成功时触发",
			Category:    "subscription",
			PayloadExample: models.JSON{
				"subscription_id": "sub_001",
				"renewed_until":   time.Now().AddDate(1, 0, 0).Unix(),
			},
		},
		{
			Name:        "订阅过期",
			EventType:   "subscription.expired",
			Description: "当订阅过期时触发",
			Category:    "subscription",
			PayloadExample: models.JSON{
				"subscription_id": "sub_001",
				"expired_at":      time.Now().Unix(),
			},
		},
		{
			Name:        "支付成功",
			EventType:   "payment.success",
			Description: "当支付成功时触发",
			Category:    "payment",
			PayloadExample: models.JSON{
				"order_id":   "order_001",
				"amount":     99.00,
				"currency":   "CNY",
				"paid_at":    time.Now().Unix(),
			},
		},
		{
			Name:        "配额超限",
			EventType:   "quota.exceeded",
			Description: "当 API 调用配额超限时触发",
			Category:    "alert",
			PayloadExample: models.JSON{
				"app_id":       123,
				"quota_type":   "daily_calls",
				"used":         10000,
				"limit":        10000,
				"reset_at":     time.Now().AddDate(0, 0, 1).Unix(),
			},
		},
	}

	for _, tpl := range templates {
		var existing models.WebhookTemplate
		if s.db.Where("event_type = ?", tpl.EventType).First(&existing).Error != nil {
			s.db.Create(&tpl)
		}
	}

	return nil
}
