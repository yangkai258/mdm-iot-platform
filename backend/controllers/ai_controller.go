package controllers

import (
	"fmt"
	"net/http"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIController AI聊天控制器
type AIController struct {
	DB *gorm.DB
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Message    string `json:"message" binding:"required"`
	SessionID  string `json:"session_id"`
	Model      string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *ChatData   `json:"data,omitempty"`
}

// ChatData 聊天数据
type ChatData struct {
	SessionID   string `json:"session_id"`
	Message    string `json:"message"`
	TokenCount int    `json:"token_count"`
	Model      string `json:"model"`
}

// Chat 聊天
// POST /api/v1/ai/chat
func (c *AIController) Chat(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	tenantID := middleware.GetTenantID(ctx)

	var req ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ChatResponse{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	// 生成或使用 session_id
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("sess_%d_%d", userID, time.Now().UnixMilli())
	}

	// 如果没有提供系统提示，使用默认
	systemPrompt := req.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "你是一个智能助手，可以帮助用户回答问题和完成各种任务。"
	}

	// 获取AI配置（这里使用模拟响应，生产环境应该调用真实AI服务）
	aiResponse, tokenCount := c.callAI(ctx, req.Message, systemPrompt, req.Model)

	// 保存用户消息
	userMsg := models.AIMessage{
		ConversationID: 0, // 新会话
		Role:          "user",
		Content:       req.Message,
		TokenCount:    0,
		Model:         req.Model,
		TenantID:      tenantID,
	}
	c.DB.Create(&userMsg)

	// 保存助手消息
	assistantMsg := models.AIMessage{
		ConversationID: userMsg.ID,
		Role:          "assistant",
		Content:       aiResponse,
		TokenCount:    tokenCount,
		Model:         req.Model,
		TenantID:      tenantID,
	}
	c.DB.Create(&assistantMsg)

	ctx.JSON(http.StatusOK, ChatResponse{
		Code:    0,
		Message: "success",
		Data: &ChatData{
			SessionID:   sessionID,
			Message:    aiResponse,
			TokenCount: tokenCount,
			Model:      req.Model,
		},
	})
}

// ChatHistory 处理 GET 请求，获取聊天历史
func (c *AIController) ChatHistory(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	tenantID := middleware.GetTenantID(ctx)

	var messages []models.AIMessage
	c.DB.Where("tenant_id = ? AND role = ?", tenantID, "user").
		Order("created_at DESC").
		Limit(50).
		Find(&messages)

	ctx.JSON(http.StatusOK, ChatResponse{
		Code:    0,
		Message: "success",
		Data: &ChatData{
			SessionID:   "history",
			Message:     "Chat history retrieved",
			TokenCount:  0,
			Model:       "",
		},
	})
	_ = userID
	_ = messages
}

// callAI 调用AI服务（这里使用模拟响应）
func (c *AIController) callAI(ctx *gin.Context, message, systemPrompt, model string) (string, int) {
	// TODO: 生产环境应该调用真实的AI服务（如OpenAI、Claude等）
	// 这里模拟一个简单的响应
	
	// 根据消息内容返回不同的响应
	var response string
	switch {
	case containsAny(message, []string{"你好", "hi", "hello"}):
		response = "你好！我是AI助手，有什么可以帮助你的吗？"
	case containsAny(message, []string{"设备", "device"}):
		response = "我可以帮你管理设备。你想了解设备的哪些信息？注册、绑定、还是状态监控？"
	case containsAny(message, []string{"会员", "member"}):
		response = "会员管理系统可以帮助你管理会员信息、积分、等级等。需要我帮你创建或查询会员吗？"
	case containsAny(message, []string{"帮助", "help"}):
		response = "我可以帮助你：\n1. 回答问题\n2. 管理设备和会员\n3. 提供系统使用指导\n\n请问有什么需要帮助的？"
	default:
		response = fmt.Sprintf("收到你的消息：%s\n\n我可以帮你回答问题和完成各种任务。有什么具体需要帮助的吗？", message)
	}

	return response, len(response) / 4 // 粗略估算token数
}

// containsAny 检查字符串是否包含任意关键词
func containsAny(s string, keywords []string) bool {
	for _, kw := range keywords {
		if len(s) >= len(kw) {
			for i := 0; i <= len(s)-len(kw); i++ {
				if s[i:i+len(kw)] == kw {
					return true
				}
			}
		}
	}
	return false
}

// GetConversations 获取会话列表
// GET /api/v1/ai/conversations
func (c *AIController) GetConversations(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	tenantID := middleware.GetTenantID(ctx)

	var conversations []models.AIConversation
	query := c.DB.Model(&models.AIConversation{}).Where("user_id = ?", userID)

	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.Order("updated_at DESC").Limit(50).Find(&conversations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    conversations,
	})
}

// GetMessages 获取会话消息
// GET /api/v1/ai/conversations/:session_id/messages
func (c *AIController) GetMessages(ctx *gin.Context) {
	sessionID := ctx.Param("session_id")

	var messages []models.AIMessage
	if err := c.DB.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&messages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    messages,
	})
}

// DeleteConversation 删除会话
// DELETE /api/v1/ai/conversations/:session_id
func (c *AIController) DeleteConversation(ctx *gin.Context) {
	sessionID := ctx.Param("session_id")

	// 删除会话消息
	c.DB.Where("session_id = ?", sessionID).Delete(&models.AIMessage{})
	// 删除会话
	c.DB.Where("session_id = ?", sessionID).Delete(&models.AIConversation{})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// GetAIConfig 获取AI配置
// GET /api/v1/ai/config
func (c *AIController) GetAIConfig(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	var configs []models.AIConfig
	query := c.DB.Model(&models.AIConfig{}).Where("status = 1")

	if tenantID != "" {
		query = query.Where("tenant_id = ? OR tenant_id = ''", tenantID)
	} else {
		query = query.Where("tenant_id = ''")
	}

	if err := query.Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    configs,
	})
}

// UpdateAIConfig 更新AI配置
// PUT /api/v1/ai/config
func (c *AIController) UpdateAIConfig(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	var req struct {
		Provider    string  `json:"provider"`
		Model       string  `json:"model"`
		APIKey      string  `json:"api_key"`
		BaseURL     string  `json:"base_url"`
		MaxTokens   int     `json:"max_tokens"`
		Temperature float64 `json:"temperature"`
		IsDefault   int     `json:"is_default"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 如果设置为默认，先取消其他默认
	if req.IsDefault == 1 {
		c.DB.Model(&models.AIConfig{}).Where("tenant_id = ? AND is_default = 1", tenantID).Update("is_default", 0)
	}

	config := models.AIConfig{
		Provider:    req.Provider,
		Model:       req.Model,
		APIKey:      req.APIKey,
		BaseURL:     req.BaseURL,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		IsDefault:   req.IsDefault,
		TenantID:    tenantID,
		Status:      1,
	}

	if err := c.DB.Create(&config).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// GetModels 获取可用模型列表
// GET /api/v1/ai/models
func (c *AIController) GetModels(ctx *gin.Context) {
	models := []map[string]interface{}{
		{
			"model":      "simulated-v1",
			"provider":   "simulated",
			"name":       "模拟模型",
			"description": "默认模拟AI模型，用于测试和演示",
			"max_tokens": 4096,
			"supported":  true,
		},
		{
			"model":      "gpt-4o",
			"provider":   "openai",
			"name":       "GPT-4o",
			"description": "OpenAI 最新多模态大模型",
			"max_tokens": 128000,
			"supported":  true,
		},
		{
			"model":      "gpt-4o-mini",
			"provider":   "openai",
			"name":       "GPT-4o Mini",
			"description": "OpenAI 高性价比小模型",
			"max_tokens": 128000,
			"supported":  true,
		},
		{
			"model":      "claude-3-5-sonnet",
			"provider":   "anthropic",
			"name":       "Claude 3.5 Sonnet",
			"description": "Anthropic 高智能模型",
			"max_tokens": 200000,
			"supported":  true,
		},
		{
			"model":      "claude-3-5-haiku",
			"provider":   "anthropic",
			"name":       "Claude 3.5 Haiku",
			"description": "Anthropic 快速响应模型",
			"max_tokens": 200000,
			"supported":  true,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    models,
	})
}

// InitDefaultAIConfig 初始化默认AI配置
func (c *AIController) InitDefaultAIConfig(tenantID string) {
	var count int64
	c.DB.Model(&models.AIConfig{}).Where("tenant_id = ?", tenantID).Count(&count)
	if count == 0 {
		config := models.AIConfig{
			Provider:    "simulated",
			Model:       "simulated-v1",
			BaseURL:     "",
			MaxTokens:   4096,
			Temperature: 0.7,
			IsDefault:   1,
			TenantID:    tenantID,
			Status:      1,
		}
		c.DB.Create(&config)
	}
}
