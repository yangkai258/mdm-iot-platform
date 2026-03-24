package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WebhookController Webhook管理控制器
type WebhookController struct {
	DB *gorm.DB
}

// ============ 请求结构 ============

type WebhookListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	IsActive *bool  `form:"is_active"`
}

type WebhookCreateRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	URL         string          `json:"url" binding:"required,url"`
	Method      string          `json:"method"`
	Secret      string          `json:"secret"`
	Headers     json.RawMessage `json:"headers"`
	EventTypes  []string        `json:"event_types"`
	RetryPolicy json.RawMessage `json:"retry_policy"`
	TimeoutSec  int             `json:"timeout_sec"`
}

type WebhookUpdateRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	URL         string          `json:"url"`
	Method      string          `json:"method"`
	Secret      string          `json:"secret"`
	Headers     json.RawMessage `json:"headers"`
	EventTypes  []string        `json:"event_types"`
	RetryPolicy json.RawMessage `json:"retry_policy"`
	TimeoutSec  int             `json:"timeout_sec"`
	IsActive    *bool           `json:"is_active"`
}

type WebhookTestRequest struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

// ============ Webhook CRUD ============

// List 获取Webhook列表
func (c *WebhookController) List(ctx *gin.Context) {
	var req WebhookListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var list []models.Webhook
	var total int64

	query := c.DB.Model(&models.Webhook{})
	if req.Keyword != "" {
		query = query.Where("name ILIKE ? OR url ILIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&list)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": list,
			"pagination": gin.H{"total": total, "page": req.Page, "page_size": req.PageSize},
		},
	})
}

// Get 获取Webhook详情
func (c *WebhookController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var webhook models.Webhook
	if err := c.DB.First(&webhook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Webhook不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": webhook})
}

// Create 创建Webhook
func (c *WebhookController) Create(ctx *gin.Context) {
	var req WebhookCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.Method == "" {
		req.Method = "POST"
	}
	if req.TimeoutSec == 0 {
		req.TimeoutSec = 30
	}

	// 获取当前用户
	createdBy := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		createdBy = uid.(string)
	}

	webhook := models.Webhook{
		Name:        req.Name,
		Description: req.Description,
		URL:         req.URL,
		Method:      req.Method,
		Secret:      req.Secret,
		Headers:     req.Headers,
		EventTypes:  req.EventTypes,
		RetryPolicy: req.RetryPolicy,
		TimeoutSec:  req.TimeoutSec,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	if err := c.DB.Create(&webhook).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": webhook})
}

// Update 更新Webhook
func (c *WebhookController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var webhook models.Webhook
	if err := c.DB.First(&webhook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Webhook不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req WebhookUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Method != "" {
		updates["method"] = req.Method
	}
	if req.Secret != "" {
		updates["secret"] = req.Secret
	}
	if req.Headers != nil {
		updates["headers"] = req.Headers
	}
	if req.EventTypes != nil {
		updates["event_types"] = req.EventTypes
	}
	if req.RetryPolicy != nil {
		updates["retry_policy"] = req.RetryPolicy
	}
	if req.TimeoutSec > 0 {
		updates["timeout_sec"] = req.TimeoutSec
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := c.DB.Model(&webhook).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&webhook, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": webhook})
}

// Delete 删除Webhook
func (c *WebhookController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var webhook models.Webhook
	if err := c.DB.First(&webhook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Webhook不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Delete(&webhook).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// Test 发送测试Webhook
func (c *WebhookController) Test(ctx *gin.Context) {
	id := ctx.Param("id")
	var webhook models.Webhook
	if err := c.DB.First(&webhook, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Webhook不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req WebhookTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果没有body，使用默认测试payload
		req.Payload = json.RawMessage(`{"event": "test", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`)
	}

	startTime := time.Now()

	// 构建请求
	var reqBody io.Reader
	if req.Payload != nil {
		reqBody = bytes.NewReader(req.Payload)
	}

	httpReq, err := http.NewRequest(webhook.Method, webhook.URL, reqBody)
	if err != nil {
		c.saveWebhookLog(webhook, req.EventType, string(req.Payload), "", 0, startTime, "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求构建失败"})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if webhook.Secret != "" {
		signature := c.computeSignature(webhook.Secret, string(req.Payload))
		httpReq.Header.Set("X-Webhook-Signature", signature)
	}

	// 添加自定义头
	if webhook.Headers != nil {
		var headers map[string]string
		json.Unmarshal(webhook.Headers, &headers)
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
	}

	timeout := time.Duration(webhook.TimeoutSec) * time.Second
	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(httpReq)
	latencyMs := time.Since(startTime).Milliseconds()

	if err != nil {
		c.saveWebhookLog(webhook, req.EventType, string(req.Payload), "", 0, startTime, "timeout", err.Error())
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"code":    504,
			"message": "请求超时或失败",
			"data": gin.H{
				"latency_ms": latencyMs,
				"error":      err.Error(),
			},
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := "success"
	if resp.StatusCode >= 400 {
		result = "failed"
	}

	c.saveWebhookLog(webhook, req.EventType, string(req.Payload), string(body), resp.StatusCode, startTime, result, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"status_code":  resp.StatusCode,
			"latency_ms":   latencyMs,
			"response_body": string(body),
			"result":       result,
		},
	})
}

// Logs 获取Webhook调用日志
func (c *WebhookController) Logs(ctx *gin.Context) {
	webhookID := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var logs []models.WebhookLog
	var total int64

	query := c.DB.Model(&models.WebhookLog{}).Where("webhook_id = ?", webhookID)

	if result := ctx.Query("result"); result != "" {
		query = query.Where("result = ?", result)
	}
	if eventType := ctx.Query("event_type"); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("request_at DESC").Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": logs,
			"pagination": gin.H{"total": total, "page": page, "page_size": pageSize},
		},
	})
}

func (c *WebhookController) computeSignature(secret, body string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

func (c *WebhookController) saveWebhookLog(webhook models.Webhook, eventType, reqBody, respBody string, statusCode int, startTime time.Time, result, errMsg string) {
	now := time.Now()
	log := models.WebhookLog{
		WebhookID:   webhook.ID,
		EventType:   eventType,
		RequestBody: reqBody,
		ResponseBody: respBody,
		StatusCode:  statusCode,
		LatencyMs:   now.Sub(startTime).Milliseconds(),
		Result:      result,
		ErrorMsg:    errMsg,
		RequestAt:   startTime,
		ResponseAt:  &now,
	}
	c.DB.Create(&log)
}
