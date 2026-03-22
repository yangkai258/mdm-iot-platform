package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetController 宠物控制器
type PetController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterRoutes 注册宠物相关路由
func (p *PetController) RegisterRoutes(r *gin.RouterGroup) {
	// P0 必须实现
	r.GET("/pets/:device_id/status", p.GetPetStatus)
	r.POST("/pets/:device_id/messages", p.SendMessage)
	r.POST("/pets/:device_id/actions", p.ExecuteAction)

	// P1 高优先级
	r.PUT("/pets/:device_id/settings", p.UpdateSettings)
	r.POST("/pets/:device_id/boost", p.MoodBoost)
	r.GET("/conversations", p.ListConversations)
	r.GET("/conversations/:id/messages", p.ListMessages)
}

// GetPetStatus 获取宠物状态
// @Summary 获取宠物状态
// @Description 获取指定设备宠物的当前状态
// @Tags pets
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/status [GET]
func (p *PetController) GetPetStatus(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var petStatus models.PetStatusV2
	result := p.DB.Where("device_id = ?", deviceID).First(&petStatus)
	if result.Error == gorm.ErrRecordNotFound {
		// 如果不存在，创建默认状态
		petStatus = models.PetStatusV2{
			DeviceID:          deviceID,
			PetName:           "小爪",
			PetType:           "cat",
			Personality:       "{}",
			Appearance:        "{}",
			Mood:              50,
			Energy:            100,
			Hunger:            0,
			CurrentExpression: "happy",
			IsOnline:          false,
		}
		if err := p.DB.Create(&petStatus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建宠物状态失败"})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询宠物状态失败"})
		return
	}

	// 检查设备是否在线（从Redis获取）
	if p.Redis != nil {
		if shadow, err := p.Redis.GetDeviceShadow(deviceID); err == nil && shadow != nil {
			petStatus.IsOnline = shadow.IsOnline
			if shadow.LastHeartbeat != nil {
				petStatus.LastSeenAt = shadow.LastHeartbeat
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": petStatus.ToResponse(),
	})
}

// SendMessage 发送消息
// @Summary 发送消息
// @Description 向宠物发送消息
// @Tags pets
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Param body body models.SendMessageRequest true "消息内容"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/messages [POST]
func (p *PetController) SendMessage(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	// 获取或创建会话
	userID := int64(1) // TODO: 从认证中间件获取
	var conversation models.Conversation
	err := p.DB.Where("device_id = ? AND user_id = ? AND status = ?", deviceID, userID, models.ConversationStatusNormal).
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		conversation = models.Conversation{
			UserID:   userID,
			DeviceID: deviceID,
			Title:    "与" + deviceID + "的对话",
			Status:   models.ConversationStatusNormal,
		}
		if err := p.DB.Create(&conversation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建会话失败"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询会话失败"})
		return
	}

	// 创建消息记录
	contentType := req.ContentType
	if contentType == 0 {
		contentType = models.ContentTypeText
	}

	message := models.Message{
		ConversationID: conversation.ConversationID,
		SenderType:     models.SenderTypeUser,
		SenderID:       &userID,
		Content:        req.Content,
		ContentType:    contentType,
		MediaURL:       req.MediaURL,
		Metadata:       req.Metadata,
	}

	if err := p.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存消息失败"})
		return
	}

	// 更新会话的最后消息
	now := time.Now()
	p.DB.Model(&conversation).Updates(map[string]interface{}{
		"last_message":    req.Content,
		"last_message_at": now,
		"message_count":   gorm.Expr("message_count + 1"),
	})

	// 通过MQTT发送消息到设备
	if mqtt.GlobalMQTTClient != nil {
		payload := map[string]interface{}{
			"message_id":   message.MessageID,
			"content":      req.Content,
			"content_type": contentType,
			"timestamp":    time.Now().Format(time.RFC3339),
		}
		if err := mqtt.PublishAction(mqtt.GlobalMQTTClient, deviceID, payload); err != nil {
			// 即使MQTT发送失败，消息已保存到数据库
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"message":         message.ToResponse(),
			"conversation_id": conversation.ConversationID,
		},
	})
}

// ExecuteAction 快捷指令下发
// @Summary 快捷指令下发
// @Description 向设备下发动作指令
// @Tags pets
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Param body body models.ActionExecuteRequest true "动作参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/actions [POST]
func (p *PetController) ExecuteAction(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req models.ActionExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	// 查找动作库中的动作
	var action models.ActionLibrary
	if err := p.DB.Where("action_id = ?", req.ActionID).First(&action).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "动作不存在"})
		return
	}

	// 通过MQTT下发动作
	if mqtt.GlobalMQTTClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "MQTT服务不可用"})
		return
	}

	// 构建动作payload
	payload := map[string]interface{}{
		"action_id":      action.ActionID,
		"action_name":    action.ActionName,
		"duration_ms":    action.DurationMs,
		"priority":       action.Priority,
		"parameters":     req.Parameters,
		"motor_commands": action.MotorCommands,
		"timestamp":      time.Now().Format(time.RFC3339),
	}

	if err := mqtt.PublishAction(mqtt.GlobalMQTTClient, deviceID, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "下发动作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": models.ActionExecuteResponse{
			ActionID: action.ActionID,
			Status:   "sent",
			Message:  "动作已下发",
		},
	})
}

// UpdateSettings 更新宠物设置
// @Summary 更新宠物设置
// @Description 更新宠物的名称、类型等设置
// @Tags pets
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Param body body models.PetSettingsUpdate true "设置参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/settings [PUT]
func (p *PetController) UpdateSettings(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req models.PetSettingsUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	var petStatus models.PetStatusV2
	if err := p.DB.Where("device_id = ?", deviceID).First(&petStatus).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			petStatus = models.PetStatusV2{
				DeviceID: deviceID,
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询宠物状态失败"})
			return
		}
	}

	updates := map[string]interface{}{}
	if req.PetName != "" {
		updates["pet_name"] = req.PetName
		petStatus.PetName = req.PetName
	}
	if req.PetType != "" {
		updates["pet_type"] = req.PetType
		petStatus.PetType = req.PetType
	}
	if req.Personality != "" {
		updates["personality"] = req.Personality
		petStatus.Personality = req.Personality
	}
	if req.Appearance != "" {
		updates["appearance"] = req.Appearance
		petStatus.Appearance = req.Appearance
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := p.DB.Model(&petStatus).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新设置失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": petStatus.ToResponse(),
	})
}

// MoodBoost 心情激励
// @Summary 心情激励
// @Description 给宠物施加心情激励
// @Tags pets
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Param body body models.MoodBoost true "激励参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/boost [POST]
func (p *PetController) MoodBoost(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req models.MoodBoost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	var petStatus models.PetStatusV2
	if err := p.DB.Where("device_id = ?", deviceID).First(&petStatus).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物状态不存在"})
		return
	}

	// 根据激励类型调整心情和能量
	boostAmount := req.Amount
	if boostAmount == 0 {
		boostAmount = 10
	}

	switch req.BoostType {
	case "food":
		petStatus.Hunger = maxInt(0, petStatus.Hunger-boostAmount)
		petStatus.Mood = minInt(100, petStatus.Mood+boostAmount/2)
		petStatus.Energy = minInt(100, petStatus.Energy+boostAmount)
	case "play":
		petStatus.Mood = minInt(100, petStatus.Mood+boostAmount)
		petStatus.Energy = maxInt(0, petStatus.Energy-boostAmount/2)
	case "praise":
		petStatus.Mood = minInt(100, petStatus.Mood+boostAmount)
	case "music":
		petStatus.Mood = minInt(100, petStatus.Mood+boostAmount/2)
		petStatus.Energy = minInt(100, petStatus.Energy+boostAmount/2)
	default:
		petStatus.Mood = minInt(100, petStatus.Mood+boostAmount)
	}

	// 更新表情
	if petStatus.Mood >= 70 {
		petStatus.CurrentExpression = "happy"
	} else if petStatus.Mood >= 40 {
		petStatus.CurrentExpression = "neutral"
	} else {
		petStatus.CurrentExpression = "sad"
	}

	updates := map[string]interface{}{
		"mood":               petStatus.Mood,
		"energy":             petStatus.Energy,
		"hunger":             petStatus.Hunger,
		"current_expression": petStatus.CurrentExpression,
		"updated_at":         time.Now(),
	}

	if err := p.DB.Model(&petStatus).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新心情失败"})
		return
	}

	// 通过MQTT下发激励动作
	if mqtt.GlobalMQTTClient != nil {
		payload := map[string]interface{}{
			"type":       "boost",
			"boost_type": req.BoostType,
			"mood":       petStatus.Mood,
			"expression": petStatus.CurrentExpression,
			"timestamp":  time.Now().Format(time.RFC3339),
		}
		mqtt.PublishAction(mqtt.GlobalMQTTClient, deviceID, payload)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"mood":               petStatus.Mood,
			"energy":             petStatus.Energy,
			"hunger":             petStatus.Hunger,
			"current_expression": petStatus.CurrentExpression,
		},
	})
}

// ListConversations 获取会话列表
// @Summary 获取会话列表
// @Description 获取用户的所有会话列表
// @Tags conversations
// @Accept json
// @Produce json
// @Param device_id query string false "设备ID"
// @Param status query int false "状态"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/conversations [GET]
func (p *PetController) ListConversations(c *gin.Context) {
	userID := int64(1) // TODO: 从认证中间件获取

	var query models.ConversationListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "查询参数错误"})
		return
	}

	var conversations []models.Conversation
	db := p.DB.Where("user_id = ?", userID)
	if query.DeviceID != "" {
		db = db.Where("device_id = ?", query.DeviceID)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	if err := db.Order("last_message_at DESC").Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询会话列表失败"})
		return
	}

	var responses []*models.ConversationResponse
	for i := range conversations {
		responses = append(responses, conversations[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// ListMessages 获取消息列表
// @Summary 获取消息列表
// @Description 获取指定会话的消息列表
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path string true "会话ID"
// @Param limit query int false "数量限制"
// @Param offset query int false "偏移量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/conversations/{id}/messages [GET]
func (p *PetController) ListMessages(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "会话ID不能为空"})
		return
	}

	var query models.MessageListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "查询参数错误"})
		return
	}

	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var messages []models.Message
	if err := p.DB.Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Offset(query.Offset).
		Limit(limit).
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询消息列表失败"})
		return
	}

	var responses []*models.MessageResponse
	for i := range messages {
		responses = append(responses, messages[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// 辅助函数
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
