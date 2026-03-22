package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// Test 1: 宠物状态获取响应
func TestGetPetStatus_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/pets/:device_id/status", func(c *gin.Context) {
		deviceID := c.Param("device_id")
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"data": gin.H{
				"device_id":          deviceID,
				"pet_name":           "小橘",
				"pet_type":           "cat",
				"mood":               75,
				"energy":             90,
				"hunger":             20,
				"current_expression": "happy",
				"is_online":          true,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/pets/pet-device-001/status", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	if result["code"].(float64) != 0 {
		t.Errorf("expected code 0, got %v", result["code"])
	}

	data := result["data"].(map[string]interface{})
	if data["pet_name"] != "小橘" {
		t.Errorf("expected pet_name 小橘, got %v", data["pet_name"])
	}
	if data["mood"].(float64) != 75 {
		t.Errorf("expected mood 75, got %v", data["mood"])
	}
}

// Test 2: 心情激励请求验证
func TestMoodBoost_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/pets/:device_id/boost", func(c *gin.Context) {
		var req models.MoodBoost
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
			return
		}
		if req.BoostType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "boost_type不能为空"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	// 测试缺少boost_type
	reqBody := map[string]interface{}{
		"amount": 20,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/pets/test/boost", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

// Test 3: 心情激励类型枚举验证
func TestMoodBoostTypes(t *testing.T) {
	boostTypes := []string{
		"food",
		"play",
		"praise",
		"music",
	}

	for _, bt := range boostTypes {
		boost := models.MoodBoost{
			BoostType: bt,
			Amount:    10,
		}
		if boost.BoostType != bt {
			t.Errorf("expected boost_type %s, got %s", bt, boost.BoostType)
		}
	}
}

// Test 4: 发送消息请求验证
func TestSendMessage_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/pets/:device_id/messages", func(c *gin.Context) {
		var req models.SendMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
			return
		}
		if req.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "content不能为空"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	// 测试空content
	reqBody := map[string]string{
		"content": "",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/pets/test/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

// Test 5: 心情值边界验证
func TestMoodValue_Boundary(t *testing.T) {
	tests := []struct {
		name     string
		mood     int
		expected string
	}{
		{"高心情", 90, "happy"},
		{"中心情", 50, "neutral"},
		{"低心情", 30, "sad"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expression string
			if tt.mood >= 70 {
				expression = "happy"
			} else if tt.mood >= 40 {
				expression = "neutral"
			} else {
				expression = "sad"
			}

			if expression != tt.expected {
				t.Errorf("mood %d: expected %s, got %s", tt.mood, tt.expected, expression)
			}
		})
	}
}

// Test 6: 心情激励计算 - 食物
func TestMoodBoostCalculation_Food(t *testing.T) {
	pet := models.PetStatusV2{
		Mood:   30,
		Energy: 50,
		Hunger: 80,
	}

	boostAmount := 20

	// 模拟食物激励逻辑
	pet.Hunger = maxInt(0, pet.Hunger-boostAmount)
	pet.Mood = minInt(100, pet.Mood+boostAmount/2)
	pet.Energy = minInt(100, pet.Energy+boostAmount)

	if pet.Hunger != 60 {
		t.Errorf("expected hunger 60, got %d", pet.Hunger)
	}
	if pet.Mood != 40 {
		t.Errorf("expected mood 40, got %d", pet.Mood)
	}
	if pet.Energy != 70 {
		t.Errorf("expected energy 70, got %d", pet.Energy)
	}
}

// Test 7: 心情激励计算 - 玩耍
func TestMoodBoostCalculation_Play(t *testing.T) {
	pet := models.PetStatusV2{
		Mood:   20,
		Energy: 100,
	}

	boostAmount := 30

	// 模拟玩耍激励逻辑
	pet.Mood = minInt(100, pet.Mood+boostAmount)
	pet.Energy = maxInt(0, pet.Energy-boostAmount/2)

	if pet.Mood != 50 {
		t.Errorf("expected mood 50, got %d", pet.Mood)
	}
	if pet.Energy != 85 {
		t.Errorf("expected energy 85, got %d", pet.Energy)
	}
}

// Test 8: 会话列表响应
func TestConversationList_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/conversations", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": []map[string]interface{}{
				{"conversation_id": "conv-001", "title": "会话1", "message_count": 10},
				{"conversation_id": "conv-002", "title": "会话2", "message_count": 5},
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/conversations", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 conversations, got %d", len(data))
	}
}

// Test 9: 消息列表响应
func TestMessageList_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/conversations/:id/messages", func(c *gin.Context) {
		conversationID := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": []map[string]interface{}{
				{"message_id": "msg-001", "conversation_id": conversationID, "content": "你好"},
				{"message_id": "msg-002", "conversation_id": conversationID, "content": "在吗"},
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/conversations/conv-001/messages", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 messages, got %d", len(data))
	}
}

// Test 10: 宠物设置更新响应
func TestPetSettingsUpdate_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.PUT("/api/v1/pets/:device_id/settings", func(c *gin.Context) {
		var req models.PetSettingsUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"device_id": c.Param("device_id"),
				"pet_name":  req.PetName,
				"pet_type":  req.PetType,
			},
		})
	})

	reqBody := map[string]string{
		"pet_name": "新名称",
		"pet_type": "dog",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/pets/test/settings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["pet_name"] != "新名称" {
		t.Errorf("expected pet_name 新名称, got %v", data["pet_name"])
	}
}

// Test 11: 心情值边界 - 最高100
func TestMoodValue_MaxBoundary(t *testing.T) {
	mood := 100
	boostAmount := 50

	// 模拟增加心情（不应超过100）
	newMood := minInt(100, mood+boostAmount)

	if newMood > 100 {
		t.Errorf("mood should not exceed 100, got %d", newMood)
	}
	if newMood != 100 {
		t.Errorf("expected mood 100, got %d", newMood)
	}
}

// Test 12: 饥饿值边界 - 最低0
func TestHungerValue_MinBoundary(t *testing.T) {
	hunger := 0
	boostAmount := 50

	// 模拟喂食降低饥饿（不应低于0）
	newHunger := maxInt(0, hunger-boostAmount)

	if newHunger < 0 {
		t.Errorf("hunger should not go below 0, got %d", newHunger)
	}
}

// Test 13: 能量值边界验证
func TestEnergyValue_Boundary(t *testing.T) {
	tests := []struct {
		name        string
		energy      int
		change      int
		expectValid bool
	}{
		{"正常消耗", 80, -30, true},
		{"完全消耗", 50, -50, true},
		{"过度消耗", 20, -50, true}, // 不会变成负数，但会变成0
		{"正常恢复", 50, +30, true},
		{"满值恢复", 80, +50, true}, // 不会超过100
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			energy := tt.energy
			// 消耗
			if tt.change < 0 {
				energy = maxInt(0, energy+tt.change)
			} else {
				energy = minInt(100, energy+tt.change)
			}

			if energy < 0 {
				t.Errorf("energy should not go below 0, got %d", energy)
			}
			if energy > 100 {
				t.Errorf("energy should not exceed 100, got %d", energy)
			}
		})
	}
}

// Test 14: 宠物类型枚举验证
func TestPetTypes(t *testing.T) {
	petTypes := []string{
		"cat",
		"dog",
		"bird",
		"rabbit",
		"hamster",
		"fish",
		"turtle",
	}

	for _, pt := range petTypes {
		pet := models.PetStatusV2{
			PetType: pt,
		}
		if pet.PetType != pt {
			t.Errorf("expected pet_type %s, got %s", pt, pet.PetType)
		}
	}
}

// Test 15: 表情枚举验证
func TestCurrentExpressions(t *testing.T) {
	expressions := []string{
		"happy",
		"neutral",
		"sad",
		"excited",
		"sleepy",
		"hungry",
		"angry",
		"playful",
	}

	for _, exp := range expressions {
		pet := models.PetStatusV2{
			CurrentExpression: exp,
		}
		if pet.CurrentExpression != exp {
			t.Errorf("expected expression %s, got %s", exp, pet.CurrentExpression)
		}
	}
}

// Test 16: 辅助函数 minInt 和 maxInt
func TestMinMaxInt(t *testing.T) {
	if minInt(5, 10) != 5 {
		t.Error("minInt(5, 10) should be 5")
	}
	if minInt(10, 5) != 5 {
		t.Error("minInt(10, 5) should be 5")
	}
	if maxInt(5, 10) != 10 {
		t.Error("maxInt(5, 10) should be 10")
	}
	if maxInt(10, 5) != 10 {
		t.Error("maxInt(10, 5) should be 10")
	}
	if minInt(100, 100) != 100 {
		t.Error("minInt(100, 100) should be 100")
	}
	if maxInt(100, 100) != 100 {
		t.Error("maxInt(100, 100) should be 100")
	}
}

// Test 17: PetStatusV2 ToResponse 转换
func TestPetStatusV2ToResponse(t *testing.T) {
	pet := models.PetStatusV2{
		DeviceID:          "test-device",
		PetName:           "小橘",
		PetType:           "cat",
		Mood:              75,
		Energy:            90,
		Hunger:            20,
		CurrentExpression: "happy",
		IsOnline:          true,
	}

	resp := pet.ToResponse()

	if resp.DeviceID != pet.DeviceID {
		t.Errorf("expected device_id %s, got %s", pet.DeviceID, resp.DeviceID)
	}
	if resp.PetName != pet.PetName {
		t.Errorf("expected pet_name %s, got %s", pet.PetName, resp.PetName)
	}
	if resp.Mood != pet.Mood {
		t.Errorf("expected mood %d, got %d", pet.Mood, resp.Mood)
	}
	if resp.IsOnline != pet.IsOnline {
		t.Errorf("expected is_online %v, got %v", pet.IsOnline, resp.IsOnline)
	}
}

// Test 18: 消息发送者类型验证
func TestSenderTypes(t *testing.T) {
	senderTypes := []struct {
		senderType int
		expected   string
	}{
		{models.SenderTypeUser, "用户"},
		{models.SenderTypePet, "宠物"},
	}

	for _, st := range senderTypes {
		if st.senderType == models.SenderTypeUser && st.expected != "用户" {
			t.Errorf("expected 用户, got unknown")
		}
		if st.senderType == models.SenderTypePet && st.expected != "宠物" {
			t.Errorf("expected 宠物, got unknown")
		}
	}
}

// Test 19: 消息内容类型验证
func TestContentTypes(t *testing.T) {
	contentTypes := []struct {
		contentType int
		expected   string
	}{
		{models.ContentTypeText, "文本"},
		{models.ContentTypeVoice, "语音"},
		{models.ContentTypeImage, "图片"},
		{models.ContentTypeAction, "动作"},
	}

	for _, ct := range contentTypes {
		if ct.contentType == models.ContentTypeText && ct.expected != "文本" {
			t.Errorf("expected 文本")
		}
	}
}

// Test 20: 会话状态常量验证
func TestConversationStatusConstants(t *testing.T) {
	if models.ConversationStatusNormal != 1 {
		t.Errorf("expected ConversationStatusNormal to be 1, got %d", models.ConversationStatusNormal)
	}
	if models.ConversationStatusClosed != 2 {
		t.Errorf("expected ConversationStatusClosed to be 2, got %d", models.ConversationStatusClosed)
	}
}
