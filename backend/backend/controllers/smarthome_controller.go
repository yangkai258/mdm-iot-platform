package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SmartHomeController 智能家居控制器
type SmartHomeController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册智能家居路由
func (ctrl *SmartHomeController) RegisterRoutes(api *gin.RouterGroup) {
	smarthome := api.Group("/smarthome")
	{
		smarthome.GET("/devices", ctrl.ListDevices)
		smarthome.GET("/devices/:id", ctrl.GetDevice)
		smarthome.POST("/devices", ctrl.CreateDevice)
		smarthome.PUT("/devices/:id", ctrl.UpdateDevice)
		smarthome.DELETE("/devices/:id", ctrl.DeleteDevice)
		smarthome.POST("/devices/:id/control", ctrl.ControlDevice)
		smarthome.GET("/platforms", ctrl.ListPlatforms)
	}
}

// ListDevices GET /api/v1/smarthome/devices - 设备列表
func (ctrl *SmartHomeController) ListDevices(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)
	platform := c.Query("platform")
	deviceType := c.Query("device_type")
	keyword := c.Query("keyword")
	isOnline := c.Query("is_online")

	query := ctrl.DB.Model(&models.SmartHomeDevice{})

	if platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if keyword != "" {
		query = query.Where("device_name LIKE ?", "%"+keyword+"%")
	}
	if isOnline != "" {
		query = query.Where("is_online = ?", isOnline)
	}

	var total int64
	query.Count(&total)

	var devices []models.SmartHomeDevice
	if err := query.Order("last_control_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":       devices,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetDevice GET /api/v1/smarthome/devices/:id - 设备详情
func (ctrl *SmartHomeController) GetDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的设备ID"})
		return
	}

	var device models.SmartHomeDevice
	if err := ctrl.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": device})
}

// CreateDevice POST /api/v1/smarthome/devices - 添加设备
func (ctrl *SmartHomeController) CreateDevice(c *gin.Context) {
	var req CreateSmartHomeDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	device := models.SmartHomeDevice{
		UserID:           req.UserID,
		IntegrationID:    req.IntegrationID,
		Platform:         req.Platform,
		PlatformDeviceID: req.PlatformDeviceID,
		DeviceName:       req.DeviceName,
		DeviceType:       req.DeviceType,
		Status:           models.JSONMap{},
		IsOnline:        true,
	}

	if err := ctrl.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": device})
}

// UpdateDevice PUT /api/v1/smarthome/devices/:id - 更新设备
func (ctrl *SmartHomeController) UpdateDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的设备ID"})
		return
	}

	var device models.SmartHomeDevice
	if err := ctrl.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req UpdateSmartHomeDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.DeviceName != "" {
		updates["device_name"] = req.DeviceName
	}
	if req.DeviceType != "" {
		updates["device_type"] = req.DeviceType
	}
	if req.Platform != "" {
		updates["platform"] = req.Platform
	}
	if req.PlatformDeviceID != "" {
		updates["platform_device_id"] = req.PlatformDeviceID
	}

	if err := ctrl.DB.Model(&device).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&device, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": device})
}

// DeleteDevice DELETE /api/v1/smarthome/devices/:id - 删除设备
func (ctrl *SmartHomeController) DeleteDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的设备ID"})
		return
	}

	if err := ctrl.DB.Delete(&models.SmartHomeDevice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ControlDevice POST /api/v1/smarthome/devices/:id/control - 控制设备
func (ctrl *SmartHomeController) ControlDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的设备ID"})
		return
	}

	var device models.SmartHomeDevice
	if err := ctrl.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req ControlDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 更新设备状态
	now := time.Now()
	updates := map[string]interface{}{
		"last_control_at": &now,
	}

	// 根据设备类型更新状态
	if req.Action == "turn_on" {
		updates["is_online"] = true
	} else if req.Action == "turn_off" {
		updates["is_online"] = false
	}

	if err := ctrl.DB.Model(&device).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "控制失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "控制成功",
		"data": gin.H{
			"device_id": id,
			"action":    req.Action,
		},
	})
}

// ListPlatforms GET /api/v1/smarthome/platforms - 支持的平台列表
func (ctrl *SmartHomeController) ListPlatforms(c *gin.Context) {
	platforms := []gin.H{
		{"id": "mi_home", "name": "小米米家", "icon": ""},
		{"id": "tmall_genie", "name": "天猫精灵", "icon": ""},
		{"id": "homekit", "name": "Apple HomeKit", "icon": ""},
		{"id": "google_home", "name": "Google Home", "icon": ""},
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": platforms})
}

// ========== 请求结构体 ==========

type CreateSmartHomeDeviceRequest struct {
	UserID           uint   `json:"user_id"`
	IntegrationID    uint   `json:"integration_id"`
	Platform         string `json:"platform"`
	PlatformDeviceID string `json:"platform_device_id"`
	DeviceName       string `json:"device_name"`
	DeviceType       string `json:"device_type"`
}

type UpdateSmartHomeDeviceRequest struct {
	DeviceName       string `json:"device_name"`
	DeviceType       string `json:"device_type"`
	Platform         string `json:"platform"`
	PlatformDeviceID string `json:"platform_device_id"`
}

type ControlDeviceRequest struct {
	Action    string                 `json:"action"` // turn_on/turn_off/set_value
	Param     string                 `json:"param"`
	Value     interface{}            `json:"value"`
	ExtraData map[string]interface{} `json:"extra_data"`
}
