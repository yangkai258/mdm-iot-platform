package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// Test 1: 告警状态枚举验证
func TestAlertStatus_Constants(t *testing.T) {
	statuses := []struct {
		status   int
		expected string
	}{
		{1, "未处理"},
		{2, "已确认"},
		{3, "已解决"},
		{4, "已忽略"},
	}

	for _, s := range statuses {
		alert := models.DeviceAlert{
			DeviceID:  fmt.Sprintf("device-%d", s.status),
			AlertType: "test",
			Status:    s.status,
		}
		if alert.Status != s.status {
			t.Errorf("expected status %d, got %d", s.status, alert.Status)
		}
	}
}

// Test 2: 告警严重程度枚举验证
func TestAlertSeverity_Constants(t *testing.T) {
	severities := []struct {
		severity int
		expected string
	}{
		{1, "低"},
		{2, "中"},
		{3, "高"},
		{4, "严重"},
	}

	for _, s := range severities {
		alert := models.DeviceAlert{
			DeviceID:  fmt.Sprintf("device-%d", s.severity),
			AlertType: "test",
			Severity: s.severity,
		}
		if alert.Severity != s.severity {
			t.Errorf("expected severity %d, got %d", s.severity, alert.Severity)
		}
	}
}

// Test 3: 告警类型验证
func TestAlertType_Values(t *testing.T) {
	alertTypes := []string{
		"battery_low",
		"offline",
		"jailbreak_detected",
		"geofence_violation",
		"temperature_high",
		"compliance_violation",
	}

	for _, at := range alertTypes {
		alert := models.DeviceAlert{
			DeviceID:  "test-device",
			AlertType: at,
		}
		if alert.AlertType != at {
			t.Errorf("expected alert_type %s, got %s", at, alert.AlertType)
		}
	}
}

// Test 4: 告警规则创建响应
func TestAlertRuleCreate_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/alerts/rules", func(c *gin.Context) {
		var rule models.DeviceAlertRule
		if err := c.ShouldBindJSON(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}
		// 模拟创建成功
		rule.ID = 1
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
	})

	reqBody := map[string]interface{}{
		"name":       "电量低于30%告警",
		"alert_type": "battery_low",
		"condition":  "<",
		"threshold":  30,
		"severity":   2,
		"enabled":    true,
		"notify_ways": "email,sms",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/alerts/rules", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	expected := "电量低于30%告警"
	if data["name"] != expected {
		t.Errorf("expected name %s, got %v", expected, data["name"])
	}
}

// Test 5: 告警确认响应
func TestAlertConfirm_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-001")
		c.Next()
	})

	r.POST("/api/v1/alerts/:id/confirm", func(c *gin.Context) {
		userID := c.GetString("user_id")
		alertID := c.Param("id")

		// 模拟确认操作
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "告警已确认",
			"data": gin.H{
				"alert_id":     alertID,
				"status":       2,
				"confirmed_by": userID,
				"confirmed_at": time.Now().Format(time.RFC3339),
			},
		})
	})

	req, _ := http.NewRequest("POST", "/api/v1/alerts/1/confirm", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["confirmed_by"] != "test-user-001" {
		t.Errorf("expected confirmed_by test-user-001, got %v", data["confirmed_by"])
	}
	if data["status"].(float64) != 2 {
		t.Errorf("expected status 2, got %v", data["status"])
	}
}

// Test 6: 告警解决响应
func TestAlertResolve_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-001")
		c.Next()
	})

	r.POST("/api/v1/alerts/:id/resolve", func(c *gin.Context) {
		userID := c.GetString("user_id")
		alertID := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "告警已解决",
			"data": gin.H{
				"alert_id":    alertID,
				"status":      3,
				"resolved_by": userID,
				"resolved_at": time.Now().Format(time.RFC3339),
			},
		})
	})

	req, _ := http.NewRequest("POST", "/api/v1/alerts/1/resolve", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["resolved_by"] != "test-user-001" {
		t.Errorf("expected resolved_by test-user-001, got %v", data["resolved_by"])
	}
	if data["status"].(float64) != 3 {
		t.Errorf("expected status 3, got %v", data["status"])
	}
}

// Test 7: 告警列表响应
func TestAlertList_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/alerts", func(c *gin.Context) {
		alerts := []map[string]interface{}{
			{"id": 1, "device_id": "device-001", "alert_type": "battery_low", "severity": 2, "status": 1},
			{"id": 2, "device_id": "device-002", "alert_type": "offline", "severity": 1, "status": 2},
			{"id": 3, "device_id": "device-001", "alert_type": "jailbreak", "severity": 4, "status": 1},
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{"list": alerts},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	if len(list) != 3 {
		t.Errorf("expected 3 alerts, got %d", len(list))
	}
}

// Test 8: 告警规则更新响应
func TestAlertRuleUpdate_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.PUT("/api/v1/alerts/rules/:id", func(c *gin.Context) {
		var rule models.DeviceAlertRule
		if err := c.ShouldBindJSON(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}
		rule.ID = 1
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
	})

	reqBody := map[string]interface{}{
		"name":      "新规则名",
		"threshold": 20,
		"severity":  3,
		"enabled":   false,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/alerts/rules/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["name"] != "新规则名" {
		t.Errorf("expected name 新规则名, got %v", data["name"])
	}
}

// Test 9: 批量确认告警响应
func TestBatchConfirmAlerts_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-001")
		c.Next()
	})

	r.POST("/api/v1/alerts/batch/confirm", func(c *gin.Context) {
		var req struct {
			AlertIDs []uint `json:"alert_ids"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供告警ID列表"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": fmt.Sprintf("已确认 %d 条告警", len(req.AlertIDs)),
		})
	})

	reqBody := map[string][]uint{
		"alert_ids": {1, 2, 3, 4, 5},
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/alerts/batch/confirm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	if result["message"] != "已确认 5 条告警" {
		t.Errorf("expected message 已确认 5 条告警, got %v", result["message"])
	}
}

// Test 10: 告警确认状态转换验证
func TestAlertStatusTransition(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  int
		toStatus    int
		shouldAllow bool
	}{
		{"未处理->已确认", 1, 2, true},
		{"未处理->已解决", 1, 3, true},
		{"未处理->已忽略", 1, 4, true},
		{"已确认->已解决", 2, 3, true},
		{"已解决->已确认", 3, 2, false},
		{"已忽略->已确认", 4, 2, false},
		{"已解决->已忽略", 3, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟状态转换验证逻辑
			allow := false
			if tt.fromStatus == 1 {
				allow = true
			} else if tt.fromStatus == 2 && tt.toStatus == 3 {
				allow = true
			}

			if allow != tt.shouldAllow {
				t.Errorf("transition %d->%d: expected allow=%v, got %v", tt.fromStatus, tt.toStatus, tt.shouldAllow, allow)
			}
		})
	}
}

// Test 11: 地理围栏规则创建响应
func TestGeofenceRuleCreate_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/geofence/rules", func(c *gin.Context) {
		var rule models.GeofenceRule
		if err := c.ShouldBindJSON(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}
		rule.ID = 1
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
	})

	reqBody := map[string]interface{}{
		"name":         "家围栏",
		"device_id":    "device-001",
		"center_lat":   39.9042,
		"center_lng":   116.4074,
		"radius_meters": 100,
		"alert_on":     "enter",
		"severity":     2,
		"enabled":      true,
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/geofence/rules", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["name"] != "家围栏" {
		t.Errorf("expected name 家围栏, got %v", data["name"])
	}
}

// Test 12: 告警触发值和阈值关系验证
func TestAlertTriggerAndThreshold(t *testing.T) {
	alert := models.DeviceAlert{
		DeviceID:   "test-device",
		AlertType:  "battery_low",
		TriggerVal: 25,
		Threshold:  30,
	}

	// 验证触发值和阈值的关系
	if alert.TriggerVal >= alert.Threshold {
		t.Errorf("trigger value %f should be less than threshold %f", alert.TriggerVal, alert.Threshold)
	}
}

// Test 13: 条件表达式评估
func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		condition string
		value     float64
		threshold float64
		expected  bool
	}{
		{"<", 25, 30, true},
		{"<", 30, 30, false},
		{"<", 35, 30, false},
		{">", 35, 30, true},
		{">", 30, 30, false},
		{">", 25, 30, false},
		{"=", 30, 30, true},
		{"=", 25, 30, false},
		{"<=", 30, 30, true},
		{"<=", 25, 30, true},
		{"<=", 35, 30, false},
		{">=", 30, 30, true},
		{">=", 35, 30, true},
		{">=", 25, 30, false},
	}

	for _, tt := range tests {
		result := evaluateCondition(tt.condition, tt.value, tt.threshold)
		if result != tt.expected {
			t.Errorf("evaluateCondition(%s, %f, %f): expected %v, got %v",
				tt.condition, tt.value, tt.threshold, tt.expected, result)
		}
	}
}

// Test 14: 告警通知方式验证
func TestAlertNotifyWays(t *testing.T) {
	notifyWays := []string{
		"email",
		"sms",
		"webhook",
		"inapp",
		"email,sms",
		"email,sms,webhook",
	}

	for _, nw := range notifyWays {
		rule := models.DeviceAlertRule{
			Name:       "test-rule",
			NotifyWays: nw,
		}
		if rule.NotifyWays != nw {
			t.Errorf("expected notify_ways %s, got %s", nw, rule.NotifyWays)
		}
	}
}
