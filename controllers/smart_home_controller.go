package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SmartHomeController 智能家居控制器
type SmartHomeController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterRoutes 注册智能家居相关路由
func (s *SmartHomeController) RegisterRoutes(r *gin.RouterGroup) {
	integrations := r.Group("/integrations/smart-home")
	{
		integrations.POST("/devices", s.RegisterDevice)
		integrations.GET("/devices", s.ListDevices)
		integrations.GET("/devices/:device_id", s.GetDevice)
		integrations.PUT("/devices/:device_id", s.UpdateDevice)
		integrations.DELETE("/devices/:device_id", s.DeleteDevice)
		integrations.POST("/trigger", s.TriggerDevice)
		integrations.GET("/status/:device_id", s.GetDeviceStatus)
		integrations.GET("/triggers", s.ListTriggers)
		integrations.POST("/triggers", s.CreateTrigger)
		integrations.DELETE("/triggers/:id", s.DeleteTrigger)
	}
}

// RegisterDeviceRequest 注册设备请求
type RegisterDeviceRequest struct {
	DeviceName      string `json:"device_name" binding:"required"`
	DeviceType      string `json:"device_type" binding:"required"`
	Brand           string `json:"brand"`
	XiaomiDeviceID  string `json:"xiaomi_device_id"`
	XiaomiToken     string `json:"xiaomi_token"`
	RoomName        string `json:"room_name"`
	DeviceStatus    string `json:"device_status"`
	HouseholdID     *uint  `json:"household_id"`
}

// RegisterDevice 注册米家设备
func (s *SmartHomeController) RegisterDevice(c *gin.Context) {
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 从 context 获取用户信息
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	device := models.SmartHomeDevice{
		DeviceName:      req.DeviceName,
		DeviceType:     req.DeviceType,
		Brand:          req.Brand,
		XiaomiDeviceID: req.XiaomiDeviceID,
		XiaomiToken:    req.XiaomiToken,
		RoomName:       req.RoomName,
		DeviceStatus:   req.DeviceStatus,
		OnlineStatus:   "offline",
		IsEnabled:      true,
		HouseholdID:    req.HouseholdID,
		OwnerID:        userID,
		TenantID:       tenantID,
	}

	if device.Brand == "" {
		device.Brand = "xiaomi"
	}
	if device.DeviceStatus == "" {
		device.DeviceStatus = "{}"
	}

	if err := s.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建设备失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "设备注册成功",
		"data": device,
	})
}

// ListDevices 获取设备列表
func (s *SmartHomeController) ListDevices(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var devices []models.SmartHomeDevice
	query := s.DB.Where("owner_id = ? AND tenant_id = ?", userID, tenantID)

	// 过滤参数
	if deviceType := c.Query("device_type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if roomName := c.Query("room_name"); roomName != "" {
		query = query.Where("room_name = ?", roomName)
	}
	if status := c.Query("online_status"); status != "" {
		query = query.Where("online_status = ?", status)
	}
	if householdID := c.Query("household_id"); householdID != "" {
		query = query.Where("household_id = ?", householdID)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.SmartHomeDevice{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": devices,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetDevice 获取设备详情
func (s *SmartHomeController) GetDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var device models.SmartHomeDevice
	if err := s.DB.Where("device_uuid = ? AND owner_id = ? AND tenant_id = ?", deviceID, userID, tenantID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": device})
}

// UpdateDevice 更新设备
func (s *SmartHomeController) UpdateDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var device models.SmartHomeDevice
	if err := s.DB.Where("device_uuid = ? AND owner_id = ? AND tenant_id = ?", deviceID, userID, tenantID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备失败"})
		return
	}

	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.DeviceName != "" {
		updates["device_name"] = req.DeviceName
	}
	if req.DeviceType != "" {
		updates["device_type"] = req.DeviceType
	}
	if req.Brand != "" {
		updates["brand"] = req.Brand
	}
	if req.RoomName != "" {
		updates["room_name"] = req.RoomName
	}
	if req.DeviceStatus != "" {
		updates["device_status"] = req.DeviceStatus
	}
	if req.XiaomiDeviceID != "" {
		updates["xiaomi_device_id"] = req.XiaomiDeviceID
	}
	if req.HouseholdID != nil {
		updates["household_id"] = req.HouseholdID
	}

	if err := s.DB.Model(&device).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新设备失败"})
		return
	}

	s.DB.First(&device, device.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": device})
}

// DeleteDevice 删除设备
func (s *SmartHomeController) DeleteDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	result := s.DB.Where("device_uuid = ? AND owner_id = ? AND tenant_id = ?", deviceID, userID, tenantID).Delete(&models.SmartHomeDevice{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// TriggerDeviceRequest 触发设备联动请求
type TriggerDeviceRequest struct {
	DeviceID     string `json:"device_id" binding:"required"`
	Action       string `json:"action" binding:"required"` // on/off/toggle/set
	Params       string `json:"params"`                   // JSON 格式参数
	TriggerType  string `json:"trigger_type"`
}

// TriggerDevice 触发设备联动
func (s *SmartHomeController) TriggerDevice(c *gin.Context) {
	var req TriggerDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	// 查询设备
	var device models.SmartHomeDevice
	if err := s.DB.Where("device_uuid = ? AND owner_id = ? AND tenant_id = ?", req.DeviceID, userID, tenantID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备失败"})
		return
	}

	if !device.IsEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备已禁用"})
		return
	}

	// 更新设备状态
	newStatus := make(map[string]interface{})
	if req.Action == "on" {
		newStatus["power"] = "on"
	} else if req.Action == "off" {
		newStatus["power"] = "off"
	} else if req.Action == "toggle" {
		var currentStatus map[string]interface{}
		json.Unmarshal([]byte(device.DeviceStatus), &currentStatus)
		if currentStatus["power"] == "on" {
			newStatus["power"] = "off"
		} else {
			newStatus["power"] = "on"
		}
	} else if req.Action == "set" && req.Params != "" {
		json.Unmarshal([]byte(req.Params), &newStatus)
	}

	newStatusJSON, _ := json.Marshal(newStatus)

	// 模拟设备在线并更新状态
	updates := map[string]interface{}{
		"device_status":  string(newStatusJSON),
		"online_status":  "online",
	}
	s.DB.Model(&device).Updates(updates)

	// 如果有触发器类型，创建触发记录
	if req.TriggerType != "" {
		trigger := models.SmartHomeTrigger{
			TriggerName:  "手动触发",
			TriggerType:  req.TriggerType,
			SourceDevice: req.DeviceID,
			ActionExpr:   string(newStatusJSON),
			TriggerCount: 1,
			OwnerID:      userID,
			TenantID:     tenantID,
		}
		now := time.Now()
		trigger.LastTriggered = &now
		s.DB.Create(&trigger)
	}

	// 模拟同步获取最新状态
	s.DB.First(&device, device.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "触发成功",
		"data": gin.H{
			"device_uuid":  device.DeviceUUID,
			"action":      req.Action,
			"new_status":  newStatus,
			"triggered_at": time.Now(),
		},
	})
}

// GetDeviceStatus 获取设备状态
func (s *SmartHomeController) GetDeviceStatus(c *gin.Context) {
	deviceID := c.Param("device_id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var device models.SmartHomeDevice
	if err := s.DB.Where("device_uuid = ? AND owner_id = ? AND tenant_id = ?", deviceID, userID, tenantID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备失败"})
		return
	}

	// 解析设备状态 JSON
	var status map[string]interface{}
	json.Unmarshal([]byte(device.DeviceStatus), &status)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"device_uuid":   device.DeviceUUID,
			"device_name":   device.DeviceName,
			"device_type":   device.DeviceType,
			"online_status": device.OnlineStatus,
			"device_status": status,
			"is_enabled":    device.IsEnabled,
			"room_name":     device.RoomName,
			"updated_at":    device.UpdatedAt,
		},
	})
}

// ListTriggers 获取触发器列表
func (s *SmartHomeController) ListTriggers(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var triggers []models.SmartHomeTrigger
	query := s.DB.Where("owner_id = ? AND tenant_id = ?", userID, tenantID)

	if triggerType := c.Query("trigger_type"); triggerType != "" {
		query = query.Where("trigger_type = ?", triggerType)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.SmartHomeTrigger{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&triggers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": triggers,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateTriggerRequest 创建触发器请求
type CreateTriggerRequest struct {
	TriggerName   string `json:"trigger_name" binding:"required"`
	TriggerType   string `json:"trigger_type" binding:"required"`
	SourceDevice  string `json:"source_device"`
	ConditionExpr string `json:"condition_expr"`
	ActionExpr    string `json:"action_expr" binding:"required"`
	HouseholdID   *uint  `json:"household_id"`
}

// CreateTrigger 创建触发器
func (s *SmartHomeController) CreateTrigger(c *gin.Context) {
	var req CreateTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	trigger := models.SmartHomeTrigger{
		TriggerName:   req.TriggerName,
		TriggerType:   req.TriggerType,
		SourceDevice:  req.SourceDevice,
		ConditionExpr: req.ConditionExpr,
		ActionExpr:    req.ActionExpr,
		IsEnabled:     true,
		HouseholdID:   req.HouseholdID,
		OwnerID:      userID,
		TenantID:     tenantID,
	}

	if err := s.DB.Create(&trigger).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建触发器失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "创建成功", "data": trigger})
}

// DeleteTrigger 删除触发器
func (s *SmartHomeController) DeleteTrigger(c *gin.Context) {
	triggerID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	result := s.DB.Where("id = ? AND owner_id = ? AND tenant_id = ?", triggerID, userID, tenantID).Delete(&models.SmartHomeTrigger{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "触发器不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}
