package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	net_url "net/url"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WebhookController Webhook 控制器
type WebhookController struct {
	DB *gorm.DB
}

// ===== 请求结构 =====

// CreateWebhookRequest 创建 Webhook 请求
type CreateWebhookRequest struct {
	Name       string                 `json:"name" binding:"required"`
	URL        string                 `json:"url" binding:"required,url"`
	Secret     string                 `json:"secret"`
	EventTypes []string               `json:"event_types" binding:"required"`
	Headers    map[string]string      `json:"headers"`
	RetryCount *int                   `json:"retry_count"`
}

// UpdateWebhookRequest 更新 Webhook 请求
type UpdateWebhookRequest struct {
	Name       string                 `json:"name"`
	URL        string                  `json:"url"`
	Secret     string                  `json:"secret"`
	EventTypes []string                `json:"event_types"`
	Headers    map[string]string       `json:"headers"`
	Status     string                  `json:"status"`
	RetryCount *int                   `json:"retry_count"`
}

// TestWebhookRequest 测试 Webhook 请求
type TestWebhookRequest struct {
	EventType string                 `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
}

// ===== Webhook CRUD =====

// List 获取 Webhook 列表
func (c *WebhookController) List(ctx *gin.Context) {
	tenantID := ctx.GetUint("tenant_id")

	var webhooks []models.Webhook
	query := c.DB.Model(&models.Webhook{})

	// 租户过滤
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 用户过滤（如果非管理员）
	if !c.isAdmin(ctx) {
		query = query.Where("tenant_id = ?", ctx.GetUint("tenant_id"))
	}

	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&webhooks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 隐藏 Secret
	for i := range webhooks {
		if webhooks[i].Secret != "" {
			webhooks[i].Secret = "****"
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": webhooks,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// Get 获取 Webhook 详情
func (c *WebhookController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	tenantID := ctx.GetUint("tenant_id")

	var webhook models.Webhook
	query := c.DB.Where("webhook_id = ?", id)

	// 非管理员只能查看自己的
	if !c.isAdmin(ctx) {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&webhook).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "Webhook 不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 隐藏 Secret
	if webhook.Secret != "" {
		webhook.Secret = "****"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    webhook,
	})
}

// Create 创建 Webhook
func (c *WebhookController) Create(ctx *gin.Context) {
	tenantID := ctx.GetUint("tenant_id")

	var req CreateWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 验证 URL 格式
	if !isValidURL(req.URL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "URL 格式无效", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 生成 Secret
	secret := req.Secret
	if secret == "" {
		secret = generateSecret(32)
	}

	// 生成签名密钥
	signKey := generateSecret(64)

	retryCount := 3
	if req.RetryCount != nil {
		retryCount = *req.RetryCount
	}

	webhook := models.Webhook{
		WebhookID:  "wh-" + uuid.New().String(),
		Name:      req.Name,
		URL:       req.URL,
		Secret:    secret,
		EventTypes: models.WebhookStringArray(req.EventTypes),
		Status:    "active",
		TenantID:  tenantID,
		Headers:   convertMapStringToInterface(req.Headers),
		RetryCount: retryCount,
	}

	if err := c.DB.Create(&webhook).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"webhook": webhook,
			"sign_key": signKey, // 只在创建时返回
		},
	})
}

// Update 更新 Webhook
func (c *WebhookController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	tenantID := ctx.GetUint("tenant_id")

	var webhook models.Webhook
	query := c.DB.Where("webhook_id = ?", id)

	if !c.isAdmin(ctx) {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&webhook).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "Webhook 不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var req UpdateWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 更新字段
	if req.Name != "" {
		webhook.Name = req.Name
	}
	if req.URL != "" {
		if !isValidURL(req.URL) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "URL 格式无效", "error_code": "ERR_BAD_REQUEST"})
			return
		}
		webhook.URL = req.URL
	}
	if req.Secret != "" && req.Secret != "****" {
		webhook.Secret = req.Secret
	}
	if req.EventTypes != nil {
		webhook.EventTypes = models.WebhookStringArray(req.EventTypes)
	}
	if req.Headers != nil {
		webhook.Headers = convertMapStringToInterface(req.Headers)
	}
	if req.Status != "" {
		webhook.Status = req.Status
	}
	if req.RetryCount != nil {
		webhook.RetryCount = *req.RetryCount
	}

	if err := c.DB.Save(&webhook).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    webhook,
	})
}

// Delete 删除 Webhook
func (c *WebhookController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	tenantID := ctx.GetUint("tenant_id")

	query := c.DB.Where("webhook_id = ?", id)
	if !c.isAdmin(ctx) {
		query = query.Where("tenant_id = ?", tenantID)
	}

	result := query.Delete(&models.Webhook{})
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "Webhook 不存在", "error_code": "ERR_NOT_FOUND"})
		return
	}

	// 同时删除关联的事件记录
	c.DB.Where("webhook_id = ?", id).Delete(&models.WebhookEvent{})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ===== 事件相关 =====

// TestWebhook 发送测试事件
func (c *WebhookController) TestWebhook(ctx *gin.Context) {
	id := ctx.Param("id")
	tenantID := ctx.GetUint("tenant_id")

	var webhook models.Webhook
	query := c.DB.Where("webhook_id = ?", id)
	if !c.isAdmin(ctx) {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&webhook).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "Webhook 不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var req TestWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 使用默认测试事件
		req.EventType = models.EventTypeWebhookTest
		req.Payload = map[string]interface{}{
			"message": "This is a test webhook event",
			"timestamp": time.Now().Unix(),
		}
	}

	// 创建测试事件
	eventID := "evt-" + uuid.New().String()
	payload := models.JSON(req.Payload)
	if payload == nil {
		payload = models.JSON{}
	}
	payload["_test"] = true
	payload["event_id"] = eventID

	event := models.WebhookEvent{
		EventID:     eventID,
		WebhookID:   webhook.WebhookID,
		EventType:   req.EventType,
		Payload:     payload,
		Status:      models.WebhookStatusPending,
		MaxAttempts: webhook.RetryCount,
	}

	// 立即发送测试事件
	go c.deliverEvent(webhook, event)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "测试事件已发送",
		"data": gin.H{
			"event_id":   eventID,
			"event_type": req.EventType,
		},
	})
}

// GetEvents 获取事件历史
func (c *WebhookController) GetEvents(ctx *gin.Context) {
	id := ctx.Param("id")
	tenantID := ctx.GetUint("tenant_id")

	// 验证 Webhook 归属
	var webhook models.Webhook
	query := c.DB.Where("webhook_id = ?", id)
	if !c.isAdmin(ctx) {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.First(&webhook).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "Webhook 不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	var events []models.WebhookEvent
	eventQuery := c.DB.Where("webhook_id = ?", id)

	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		eventQuery = eventQuery.Where("status = ?", status)
	}

	// 事件类型过滤
	if eventType := ctx.Query("event_type"); eventType != "" {
		eventQuery = eventQuery.Where("event_type = ?", eventType)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	eventQuery.Count(&total)

	offset := (page - 1) * pageSize
	if err := eventQuery.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": events,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// ===== 事件投递 =====

// TriggerEvent 触发事件（内部使用）
func (c *WebhookController) TriggerEvent(webhookID string, eventType string, payload map[string]interface{}) error {
	// 查找 Webhook
	var webhook models.Webhook
	if err := c.DB.Where("webhook_id = ? AND status = ?", webhookID, "active").First(&webhook).Error; err != nil {
		return err
	}

	// 检查是否监听此事件类型
	hasEvent := false
	for _, et := range webhook.EventTypes {
		if et == eventType || et == "*" {
			hasEvent = true
			break
		}
	}
	if !hasEvent {
		return nil
	}

	// 创建事件记录
	eventID := "evt-" + uuid.New().String()
	event := models.WebhookEvent{
		EventID:     eventID,
		WebhookID:   webhook.WebhookID,
		EventType:   eventType,
		Payload:     models.JSON(payload),
		Status:      models.WebhookStatusPending,
		MaxAttempts: webhook.RetryCount,
	}

	// 幂等检查：相同 event_id 不重复发送
	var existing models.WebhookEvent
	if err := c.DB.Where("event_id = ?", eventID).First(&existing).Error; err == nil {
		return nil // 已存在，不重复创建
	}

	c.DB.Create(&event)

	// 异步投递
	go c.deliverEvent(webhook, event)

	return nil
}

// deliverEvent 投递事件
func (c *WebhookController) deliverEvent(webhook models.Webhook, event models.WebhookEvent) {
	// 构建请求
	reqBody, _ := json.Marshal(map[string]interface{}{
		"event_id":   event.EventID,
		"event_type": event.EventType,
		"timestamp":  time.Now().Unix(),
		"data":       event.Payload,
	})

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		c.updateEventStatus(event.EventID, models.WebhookStatusFailed, 0, err.Error(), "")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-ID", webhook.WebhookID)
	req.Header.Set("X-Event-ID", event.EventID)
	req.Header.Set("X-Event-Type", event.EventType)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	// 添加签名
	if webhook.Secret != "" {
		signature := generateSignature(webhook.Secret, reqBody)
		req.Header.Set("X-Signature", signature)
	}

	// 添加自定义请求头
	if webhook.Headers != nil {
		for key, value := range webhook.Headers {
			if strVal, ok := value.(string); ok {
				req.Header.Set(key, strVal)
			}
		}
	}

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.handleDeliveryFailure(event, webhook, err.Error())
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	responseBody := string(body)

	// 检查响应状态
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.updateEventStatus(event.EventID, models.WebhookStatusSuccess, resp.StatusCode, "", responseBody)
	} else {
		c.handleDeliveryFailure(event, webhook, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, responseBody))
	}
}

// handleDeliveryFailure 处理投递失败
func (c *WebhookController) handleDeliveryFailure(event models.WebhookEvent, webhook models.Webhook, errMsg string) {
	event.Attempts++
	event.LastError = errMsg

	if event.Attempts < event.MaxAttempts {
		// 计算下次重试时间
		nextRetry := time.Now().Add(time.Duration(webhook.RetryCount*60) * time.Second)
		event.NextRetryAt = &nextRetry
		event.Status = models.WebhookStatusPending
	} else {
		event.Status = models.WebhookStatusFailed
		now := time.Now()
		event.DeliveredAt = &now
	}

	c.DB.Save(&event)

	// 如果还没达到最大重试次数，调度下次重试
	if event.Attempts < event.MaxAttempts {
		go func() {
			time.Sleep(time.Duration(webhook.RetryCount*60) * time.Second)
			// 重新获取最新事件状态
			var latestEvent models.WebhookEvent
			if err := c.DB.Where("event_id = ?", event.EventID).First(&latestEvent).Error; err == nil {
				if latestEvent.Status == models.WebhookStatusPending {
					c.deliverEvent(webhook, latestEvent)
				}
			}
		}()
	}
}

// updateEventStatus 更新事件状态
func (c *WebhookController) updateEventStatus(eventID, status string, responseCode int, errMsg, responseBody string) {
	updates := map[string]interface{}{
		"status": status,
	}
	if responseCode > 0 {
		updates["response_code"] = responseCode
	}
	if errMsg != "" {
		updates["last_error"] = errMsg
	}
	if responseBody != "" {
		updates["response_body"] = responseBody
	}

	if status == models.WebhookStatusSuccess {
		now := time.Now()
		updates["delivered_at"] = &now
	}

	c.DB.Model(&models.WebhookEvent{}).Where("event_id = ?", eventID).Updates(updates)
}

// ===== 辅助函数 =====

func generateSecret(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = byte(32 + (i * 7) % 94) // 可打印 ASCII 字符
	}
	return string(bytes)
}

func generateSignature(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func isValidURL(str string) bool {
	_, err := net_url.Parse(str)
	return err == nil
}

func convertMapStringToInterface(m map[string]string) models.JSON {
	if m == nil {
		return models.JSON{}
	}
	result := make(models.JSON)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func (c *WebhookController) isAdmin(ctx *gin.Context) bool {
	role := ctx.GetString("role")
	return role == "admin" || role == "super_admin"
}
