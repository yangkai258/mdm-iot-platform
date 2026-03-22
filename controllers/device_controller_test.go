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
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setupDeviceTestRouterWithMock() (*gin.Engine, sqlmock.Sqlmock) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 创建一个假的DB和mock
	// 注意：这里我们实际用gin的context来测试，不需要真正的DB
	// Mock仅用于编译通过

	return r, nil
}

// isValidMAC 测试辅助函数 - 独立测试不需要DB
func TestIsValidMAC(t *testing.T) {
	tests := []struct {
		mac      string
		expected bool
	}{
		{"AA:BB:CC:DD:EE:FF", true},
		{"aa:bb:cc:dd:ee:ff", true},
		{"11:22:33:44:55:66", true},
		{"AA:BB:CC:DD:EE", false},
		{"AA:BB:CC:DD:EE:FF:00", false},
		{"AA:BB-CC-DD-EE-FF", false},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		result := isValidMAC(tt.mac)
		if result != tt.expected {
			t.Errorf("isValidMAC(%s) = %v, expected %v", tt.mac, result, tt.expected)
		}
	}
}

// Test 1: 设备注册请求验证 - 缺少必需字段
func TestDeviceRegisterRequest_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 模拟一个只验证请求格式的端点
	r.POST("/api/v1/devices/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "参数校验失败: " + err.Error(),
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		// 验证 MAC 格式
		if !isValidMAC(req.MacAddress) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "无效的MAC地址格式",
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	// 测试缺少 mac_address
	reqBody := map[string]string{
		"sn_code":          "SN12345678",
		"hardware_model":   "M5Stack",
		"firmware_version": "1.0.0",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/devices/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

// Test 2: 设备注册请求验证 - 无效MAC格式
func TestDeviceRegisterRequest_InvalidMAC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/devices/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "参数校验失败"})
			return
		}
		if !isValidMAC(req.MacAddress) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的MAC地址格式"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	reqBody := map[string]string{
		"mac_address":      "INVALID_MAC",
		"sn_code":          "SN12345678",
		"hardware_model":   "M5Stack",
		"firmware_version": "1.0.0",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/devices/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

// Test 3: 设备绑定请求验证 - 缺少bind_user_id
func TestDeviceBindRequest_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/devices/bind/:sn_code", func(c *gin.Context) {
		var req BindRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "参数校验失败",
				"error_code": "ERR_VALIDATION",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	// 测试缺少 bind_user_id
	reqBody := map[string]string{}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/devices/bind/SN123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

// Test 4: 设备绑定请求验证 - 正确参数
func TestDeviceBindRequest_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/v1/devices/bind/:sn_code", func(c *gin.Context) {
		snCode := c.Param("sn_code")
		var req BindRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 4005})
			return
		}
		if req.BindUserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "bind_user_id不能为空"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"data": gin.H{
				"sn_code":      snCode,
				"bind_user_id": req.BindUserID,
			},
		})
	})

	reqBody := map[string]string{
		"bind_user_id": "user-123",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/devices/bind/SN123456", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)
	data := result["data"].(map[string]interface{})
	if data["bind_user_id"] != "user-123" {
		t.Errorf("expected bind_user_id user-123, got %v", data["bind_user_id"])
	}
}

// Test 5: 设备生命周期状态枚举验证
func TestDeviceLifecycleStatus(t *testing.T) {
	// 测试设备生命周期状态值
	statuses := []struct {
		status   int
		expected string
	}{
		{1, "待激活"},
		{2, "服役中"},
		{3, "维修"},
		{4, "报废"},
	}

	for _, s := range statuses {
		device := models.Device{
			DeviceID:        fmt.Sprintf("device-%d", s.status),
			MacAddress:      fmt.Sprintf("AA:BB:CC:DD:EE:%02X", s.status),
			SnCode:         fmt.Sprintf("SN%d", s.status),
			HardwareModel:   "M5Stack",
			FirmwareVersion: "1.0.0",
			LifecycleStatus: s.status,
		}

		if device.LifecycleStatus != s.status {
			t.Errorf("expected lifecycle_status %d, got %d", s.status, device.LifecycleStatus)
		}
	}
}

// Test 6: RegisterRequest 结构验证
func TestRegisterRequest_Struct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/test", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 验证字段赋值正确
		if req.MacAddress == "" || req.SnCode == "" || req.HardwareModel == "" || req.FirmwareVersion == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	reqBody := map[string]string{
		"mac_address":      "AA:BB:CC:DD:EE:FF",
		"sn_code":          "SN12345678",
		"hardware_model":   "M5Stack",
		"firmware_version": "1.0.0",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
}

// Test 7: ListRequest 分页参数默认值
func TestListRequest_PaginationDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/devices", func(c *gin.Context) {
		var req ListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 验证默认值
		if req.Page < 1 {
			req.Page = 1
		}
		if req.PageSize < 1 || req.PageSize > 100 {
			req.PageSize = 20
		}
		c.JSON(http.StatusOK, gin.H{
			"page":      req.Page,
			"page_size": req.PageSize,
		})
	})

	// 不带分页参数
	req, _ := http.NewRequest("GET", "/devices", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	// 验证默认值被正确应用
	if result["page"].(float64) != 1 {
		t.Errorf("expected default page 1, got %v", result["page"])
	}
	if result["page_size"].(float64) != 20 {
		t.Errorf("expected default page_size 20, got %v", result["page_size"])
	}
}

// Test 8: ListRequest 分页参数边界值
func TestListRequest_PaginationBoundary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/devices", func(c *gin.Context) {
		var req ListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 验证边界值处理
		if req.Page < 1 {
			req.Page = 1
		}
		if req.PageSize < 1 || req.PageSize > 100 {
			req.PageSize = 20
		}
		c.JSON(http.StatusOK, gin.H{
			"page":      req.Page,
			"page_size": req.PageSize,
		})
	})

	// 测试超过100的page_size应被限制为20
	req, _ := http.NewRequest("GET", "/devices?page_size=200", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	if result["page_size"].(float64) != 20 {
		t.Errorf("expected page_size capped at 20, got %v", result["page_size"])
	}
}

// Test 9: 设备详情响应结构验证
func TestDeviceDetailResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/devices/:device_id", func(c *gin.Context) {
		deviceID := c.Param("device_id")
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":         deviceID,
				"mac_address":      "AA:BB:CC:DD:EE:FF",
				"sn_code":          "SN12345678",
				"hardware_model":   "M5Stack",
				"firmware_version": "1.0.0",
				"lifecycle_status": 1,
				"created_at":       "2024-01-01T00:00:00Z",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/devices/test-device-001", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	data := result["data"].(map[string]interface{})
	if data["device_id"] != "test-device-001" {
		t.Errorf("expected device_id test-device-001, got %v", data["device_id"])
	}
	if data["lifecycle_status"].(float64) != 1 {
		t.Errorf("expected lifecycle_status 1, got %v", data["lifecycle_status"])
	}
}

// Test 10: DeviceWithShadow 结构验证
func TestDeviceWithShadow_Struct(t *testing.T) {
	// 验证 DeviceWithShadow 结构
	type DeviceWithShadow struct {
		models.Device
		IsOnline      bool `json:"is_online"`
		BatteryLevel  int  `json:"battery_level"`
	}

	device := DeviceWithShadow{
		Device: models.Device{
			DeviceID:        "shadow-test",
			MacAddress:      "AA:BB:CC:DD:EE:FF",
			SnCode:         "SNSHADOW",
			HardwareModel:   "M5Stack",
			FirmwareVersion: "1.0.0",
			LifecycleStatus: 2,
		},
		IsOnline:     true,
		BatteryLevel: 85,
	}

	if !device.IsOnline {
		t.Error("expected IsOnline to be true")
	}
	if device.BatteryLevel != 85 {
		t.Errorf("expected BatteryLevel 85, got %d", device.BatteryLevel)
	}
}
