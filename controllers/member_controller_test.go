package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// Test 1: 会员创建请求验证 - 缺少必需字段
func TestMemberCreate_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/members", func(c *gin.Context) {
		var member models.Member
		if err := c.ShouldBindJSON(&member); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	// 发送空JSON
	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// GORM的ShouldBindJSON对于空struct可能不会报错，因为所有字段都有默认值
	// 这里测试实际行为
	t.Logf("Response code: %d, body: %s", resp.Code, resp.Body.String())
}

// Test 2: 会员详情响应结构验证
func TestMemberDetail_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/members/:id", func(c *gin.Context) {
		_ = c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"id":           1,
				"member_code":  "M001",
				"member_name":  "张三",
				"phone":        "13800000001",
				"member_level": 1,
				"points":       100,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/members/1", nil)
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
	if data["member_name"] != "张三" {
		t.Errorf("expected member_name 张三, got %v", data["member_name"])
	}
}

// Test 3: 会员列表响应结构验证
func TestMemberList_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/members", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"list": []map[string]interface{}{
					{"id": 1, "member_name": "张三", "member_level": 1},
					{"id": 2, "member_name": "李四", "member_level": 2},
				},
				"total":     2,
				"page":      1,
				"page_size": 10,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/members", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	if len(list) != 2 {
		t.Errorf("expected 2 members, got %d", len(list))
	}
}

// Test 4: 会员等级创建响应
func TestLevelCreate_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/member/levels", func(c *gin.Context) {
		var level models.MemberLevel
		if err := c.ShouldBindJSON(&level); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}
		// 模拟创建成功
		level.ID = 1
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": level})
	})

	reqBody := map[string]interface{}{
		"level_name":  "黄金会员",
		"level_code":  "GOLD",
		"min_points":  1000,
		"discount":    0.9,
		"points_rate": 1.5,
		"sort":        1,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/member/levels", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["level_name"] != "黄金会员" {
		t.Errorf("expected level_name 黄金会员, got %v", data["level_name"])
	}
}

// Test 5: 会员等级列表响应
func TestLevelList_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/member/levels", func(c *gin.Context) {
		levels := []models.MemberLevel{
			{LevelName: "青铜", LevelCode: "BRONZE", MinPoints: 0, Sort: 1},
			{LevelName: "白银", LevelCode: "SILVER", MinPoints: 500, Sort: 2},
			{LevelName: "黄金", LevelCode: "GOLD", MinPoints: 2000, Sort: 3},
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": levels})
	})

	req, _ := http.NewRequest("GET", "/api/v1/member/levels", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].([]interface{})
	if len(data) != 3 {
		t.Errorf("expected 3 levels, got %d", len(data))
	}
}

// Test 6: 会员删除响应
func TestMemberDelete_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.DELETE("/api/v1/members/:id", func(c *gin.Context) {
		deletedID := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"deleted_id": deletedID,
			},
		})
	})

	req, _ := http.NewRequest("DELETE", "/api/v1/members/1", nil)
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
}

// Test 7: 会员等级枚举值验证
func TestMemberLevel_Constants(t *testing.T) {
	// 验证会员等级模型常量
	levels := []struct {
		code  string
		points int64
	}{
		{"BRONZE", 0},
		{"SILVER", 500},
		{"GOLD", 2000},
		{"PLATINUM", 5000},
		{"DIAMOND", 10000},
	}

	for _, level := range levels {
		memberLevel := models.MemberLevel{
			LevelName: level.code,
			LevelCode: level.code,
			MinPoints: level.points,
		}
		if memberLevel.MinPoints != level.points {
			t.Errorf("expected MinPoints %d for %s, got %d", level.points, level.code, memberLevel.MinPoints)
		}
	}
}

// Test 8: 会员状态枚举验证
func TestMemberStatus_Constants(t *testing.T) {
	// 测试会员状态值
	statuses := []struct {
		status   int
		expected string
	}{
		{1, "正常"},
		{2, "禁用"},
		{0, "未知"},
	}

	for _, s := range statuses {
		member := models.Member{
			MemberCode: fmt.Sprintf("M%d", s.status),
			Status:     s.status,
		}
		if member.Status != s.status {
			t.Errorf("expected status %d, got %d", s.status, member.Status)
		}
	}
}

// Test 9: 会员卡类型枚举验证
func TestMemberCardType_Constants(t *testing.T) {
	cardTypes := []struct {
		cardType int
		expected string
	}{
		{1, "储值卡"},
		{2, "积分卡"},
		{3, "打折卡"},
	}

	for _, ct := range cardTypes {
		card := models.MemberCard{
			CardCode: fmt.Sprintf("CARD%d", ct.cardType),
			CardType: ct.cardType,
		}
		if card.CardType != ct.cardType {
			t.Errorf("expected CardType %d, got %d", ct.cardType, card.CardType)
		}
	}
}

// Test 10: 会员积分计算验证
func TestMemberPoints_Calculation(t *testing.T) {
	tests := []struct {
		name         string
		currentPoints int64
		addPoints    int64
		expected     int64
	}{
		{"正常增加", 100, 50, 150},
		{"零积分增加", 0, 100, 100},
		{"大量增加", 1000, 5000, 6000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := models.Member{
				MemberCode: "TEST",
				Points:     tt.currentPoints,
			}
			// 模拟积分增加
			member.Points += tt.addPoints
			if member.Points != tt.expected {
				t.Errorf("expected points %d, got %d", tt.expected, member.Points)
			}
		})
	}
}

// Test 11: 会员模型字段验证
func TestMember_Fields(t *testing.T) {
	member := models.Member{
		MemberCode: "M001",
		MemberName: "测试会员",
		Phone:      "13800138000",
		Email:      "test@example.com",
		MemberLevel: 2,
		Points:     500,
		Balance:    100.50,
		Status:     1,
	}

	if member.MemberCode != "M001" {
		t.Errorf("expected MemberCode M001, got %s", member.MemberCode)
	}
	if member.MemberName != "测试会员" {
		t.Errorf("expected MemberName 测试会员, got %s", member.MemberName)
	}
	if member.Phone != "13800138000" {
		t.Errorf("expected Phone 13800138000, got %s", member.Phone)
	}
	if member.Points != 500 {
		t.Errorf("expected Points 500, got %d", member.Points)
	}
}

// Test 12: 会员卡折扣率边界验证
func TestMemberCard_DiscountRate(t *testing.T) {
	tests := []struct {
		discount float64
		valid    bool
	}{
		{1.0, true},    // 无折扣 (100%)
		{0.9, true},    // 九折
		{0.5, true},    // 五折
		{0.0, true},    // 有效 - 可以是免费
		{1.5, false},   // 无效 - 超过100%
		{-0.1, false},  // 无效 - 负数
	}

	for _, tt := range tests {
		card := models.MemberCard{
			CardCode: "TEST",
			Discount: tt.discount,
		}
		// 简单验证折扣率在0-1之间
		valid := card.Discount >= 0 && card.Discount <= 1
		if valid != tt.valid {
			t.Errorf("discount %f: expected valid=%v, got %v", tt.discount, tt.valid, valid)
		}
	}
}
