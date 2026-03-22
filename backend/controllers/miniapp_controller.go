package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MiniAppController 微信小程序后端API控制器
type MiniAppController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// ==================== 小程序设备列表 API ====================

// GetMiniAppDevices 获取小程序设备列表
// GET /api/v1/miniapp/devices
func (c *MiniAppController) GetMiniAppDevices(ctx *gin.Context) {
	openID, _ := ctx.Get("open_id") // 从JWT解析的微信openid
	if openID == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权，缺少OpenID"})
		return
	}

	page, _ := ctx.GetQuery("page")
	pageSize, _ := ctx.GetQuery("page_size")
	pageInt := 1
	pageSizeInt := 20
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(pageSize, "%d", &pageSizeInt)
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 || pageSizeInt > 50 {
		pageSizeInt = 20
	}
	offset := (pageInt - 1) * pageSizeInt

	// 查询该用户绑定的设备
	type DeviceItem struct {
		models.MiniAppDevice
		DeviceShadow    *utils.DeviceShadow `json:"device_shadow"`
		FirmwareVersion string              `json:"firmware_version"`
		HardwareModel   string              `json:"hardware_model"`
		BindUserID      string              `json:"bind_user_id"`
	}

	var total int64
	c.DB.Model(&models.MiniAppDevice{}).Where("open_id = ? AND is_active = ?", openID, true).Count(&total)

	var bindings []models.MiniAppDevice
	if err := c.DB.Where("open_id = ? AND is_active = ?", openID, true).
		Order("bind_time DESC").Offset(offset).Limit(pageSizeInt).Find(&bindings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	deviceIDs := make([]string, len(bindings))
	for i, b := range bindings {
		deviceIDs[i] = b.DeviceID
	}

	// 批量获取设备信息
	var devices []models.Device
	if len(deviceIDs) > 0 {
		c.DB.Where("device_id IN ?", deviceIDs).Find(&devices)
	}
	deviceMap := make(map[string]models.Device)
	for _, d := range devices {
		deviceMap[d.DeviceID] = d
	}

	// 批量获取设备影子
	shadowMap := make(map[string]*utils.DeviceShadow)
	if c.Redis != nil {
		for _, did := range deviceIDs {
			shadow, _ := c.Redis.GetDeviceShadow(did)
			if shadow != nil {
				shadowMap[did] = shadow
			}
		}
	}

	// 组装结果
	items := make([]DeviceItem, len(bindings))
	for i, b := range bindings {
		items[i] = DeviceItem{
			MiniAppDevice: b,
			DeviceShadow:  shadowMap[b.DeviceID],
		}
		if dev, ok := deviceMap[b.DeviceID]; ok {
			items[i].FirmwareVersion = dev.FirmwareVersion
			items[i].HardwareModel = dev.HardwareModel
			if dev.BindUserID != nil {
				items[i].BindUserID = *dev.BindUserID
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      items,
			"total":     total,
			"page":      pageInt,
			"page_size": pageSizeInt,
		},
	})
}

// BindDeviceRequest 绑定设备请求
type BindDeviceRequest struct {
	DeviceID  string `json:"device_id" binding:"required"`
	BindToken string `json:"bind_token"` // 扫码绑定时使用
	Nickname  string `json:"nickname"`   // 设备昵称
}

// BindMiniAppDevice 绑定设备
// POST /api/v1/miniapp/bind
func (c *MiniAppController) BindDevice(ctx *gin.Context) {
	openIDVal, exists := ctx.Get("open_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权，缺少OpenID"})
		return
	}
	openID := openIDVal.(string)

	userIDVal, _ := ctx.Get("user_id")
	userID := uint(0)
	if userIDVal != nil {
		userID = uint(userIDVal.(int))
	}

	var req BindDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", req.DeviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 如果有BindToken，验证二维码绑定
	if req.BindToken != "" {
		var qrBind models.MiniAppQRCodeBind
		if err := c.DB.Where("qr_code_id = ? AND device_id = ? AND is_bound = ?", req.BindToken, req.DeviceID, false).First(&qrBind).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的绑定码或已过期"})
			return
		}
		if time.Now().After(qrBind.ExpiresAt) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "绑定码已过期，请重新扫码"})
			return
		}
		// 标记二维码已使用
		qrBind.IsBound = true
		now := time.Now()
		qrBind.BoundAt = &now
		qrBind.BindUserID = userID
		qrBind.OpenID = openID
		c.DB.Save(&qrBind)
	}

	// 检查是否已绑定（不能重复绑定）
	var existing models.MiniAppDevice
	if err := c.DB.Where("device_id = ? AND is_active = ?", req.DeviceID, true).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "该设备已被其他用户绑定"})
		return
	}

	// 解绑旧绑定记录（如果有的话，软删除）
	c.DB.Model(&models.MiniAppDevice{}).Where("open_id = ? AND device_id = ?", openID, req.DeviceID).
		Update("is_active", false)

	// 创建新的绑定记录
	binding := models.MiniAppDevice{
		UserID:    userID,
		OpenID:    openID,
		DeviceID:  req.DeviceID,
		Nickname:  req.Nickname,
		BindTime:  time.Now(),
		IsActive:  true,
		BindToken: req.BindToken,
	}
	if err := c.DB.Create(&binding).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "绑定失败"})
		return
	}

	// 同时更新设备的BindUserID
	bindUserIDStr := fmt.Sprintf("%d", userID)
	c.DB.Model(&device).Updates(map[string]interface{}{
		"bind_user_id": bindUserIDStr,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id": req.DeviceID,
			"nickname":  req.Nickname,
			"bind_time": binding.BindTime,
		},
	})
}

// UnbindDeviceRequest 解绑设备请求
type UnbindDeviceRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

// UnbindMiniAppDevice 解绑设备
// POST /api/v1/miniapp/unbind
func (c *MiniAppController) UnbindDevice(ctx *gin.Context) {
	openIDVal, exists := ctx.Get("open_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权，缺少OpenID"})
		return
	}
	openID := openIDVal.(string)

	var req UnbindDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查找绑定记录
	var binding models.MiniAppDevice
	if err := c.DB.Where("device_id = ? AND open_id = ? AND is_active = ?", req.DeviceID, openID, true).First(&binding).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到绑定记录"})
		return
	}

	// 软解绑
	now := time.Now()
	binding.IsActive = false
	binding.UnbindTime = &now
	c.DB.Save(&binding)

	// 清除设备BindUserID
	c.DB.Model(&models.Device{}).Where("device_id = ?", req.DeviceID).Update("bind_user_id", nil)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id":   req.DeviceID,
			"unbind_time": now,
		},
	})
}

// GenerateQRCodeRequest 生成二维码请求
type GenerateQRCodeRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Scene    string `json:"scene"` // 场景描述
}

// GenerateQRCode 生成设备绑定二维码
// GET /api/v1/miniapp/qrcode?device_id=xxx
func (c *MiniAppController) GenerateQRCode(ctx *gin.Context) {
	deviceID := ctx.GetQuery("device_id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "device_id不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 检查设备是否已被绑定（避免恶意扫码）
	var existing models.MiniAppDevice
	if err := c.DB.Where("device_id = ? AND is_active = ?", deviceID, true).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "该设备已被绑定，请先解绑"})
		return
	}

	// 生成二维码ID和绑定Token
	qrCodeID := uuid.New().String()
	bindToken := uuid.New().String()
	expiresMinutes := 30
	expiresAt := time.Now().Add(time.Duration(expiresMinutes) * time.Minute)

	qrBind := models.MiniAppQRCodeBind{
		QRCodeID:      qrCodeID,
		DeviceID:      deviceID,
		Scene:         "device_bind",
		ExpireMinutes: expiresMinutes,
		ExpiresAt:     expiresAt,
		IsBound:       false,
	}
	if err := c.DB.Create(&qrBind).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成二维码失败"})
		return
	}

	// 构建二维码内容（实际是微信小程序码的scene参数）
	// 微信小程序码需要通过微信API生成，这里先返回二维码ID
	qrContent := map[string]interface{}{
		"qr_code_id":     qrCodeID,
		"bind_token":     bindToken,
		"device_id":      deviceID,
		"expires_at":     expiresAt,
		"expire_seconds": expiresMinutes * 60,
	}

	// TODO: 实际应调用微信API生成小程序码
	// 微信通过 /cgi-bin/wxaapp/createwxaqrcode 生成小程序码
	// 返回的二维码URL供前端生成二维码图片
	qrURL := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacode?scene=%s", qrCodeID)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"qr_code_id":     qrCodeID,
			"bind_token":     bindToken,
			"device_id":      deviceID,
			"qr_url":         qrURL,
			"expires_at":     expiresAt,
			"expire_seconds": expiresMinutes * 60,
		},
	})
}

// ==================== 辅助函数 ====================

// buildMiniAppClaims 构建小程序专用JWT Claims（用于解析openid）
// 注意：实际项目中openid应通过微信登录流程获取，此处简化处理
func buildMiniAppClaims(userID uint, openID, tenantID string) *middleware.JWTClaims {
	return &middleware.JWTClaims{
		UserID:   userID,
		Username: "miniapp",
		RoleID:   0,
		TenantID: tenantID,
	}
}

// MiniAppAuthMiddleware 小程序专用认证中间件（从请求头解析openid）
// 实际项目中openid应通过微信授权code换取
func MiniAppAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 优先从X-OpenID头获取（由网关/SSO服务填充）
		openID := ctx.GetHeader("X-OpenID")
		if openID != "" {
			ctx.Set("open_id", openID)
			ctx.Next()
			return
		}

		// 也可以从jwt token中获取openid（通过自定义claim）
		claimsVal, exists := ctx.Get("claims")
		if exists {
			claims := claimsVal.(map[string]interface{})
			if oid, ok := claims["open_id"].(string); ok && oid != "" {
				ctx.Set("open_id", oid)
				ctx.Next()
				return
			}
		}

		// 临时方案：从device_id参数获取openid（仅用于开发调试）
		openID = ctx.GetHeader("X-MiniApp-OpenID")
		if openID != "" {
			ctx.Set("open_id", openID)
			ctx.Next()
			return
		}

		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权，缺少OpenID"})
		ctx.Abort()
	}
}

// GetMiniAppDeviceStatus 获取小程序设备状态
// GET /api/v1/miniapp/device/:device_id/status
func (c *MiniAppController) GetDeviceStatus(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	openIDVal, exists := ctx.Get("open_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	openID := openIDVal.(string)

	// 验证设备是否绑定到该用户
	var binding models.MiniAppDevice
	if err := c.DB.Where("device_id = ? AND open_id = ? AND is_active = ?", deviceID, openID, true).First(&binding).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "该设备未绑定到您的小程序"})
		return
	}

	// 获取设备信息
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 构建状态
	status := map[string]interface{}{
		"device_id":    device.DeviceID,
		"sn_code":      device.SnCode,
		"nickname":     binding.Nickname,
		"firmware_ver": device.FirmwareVersion,
		"lifecycle":    device.LifecycleStatus,
	}

	// 获取设备影子
	if c.Redis != nil {
		shadow, _ := c.Redis.GetDeviceShadow(deviceID)
		if shadow != nil {
			status["is_online"] = shadow.IsOnline
			status["battery"] = shadow.BatteryLevel
			status["current_mode"] = shadow.CurrentMode
			status["last_ip"] = shadow.LastIP
			status["last_heartbeat"] = shadow.LastHeartbeat
		} else {
			status["is_online"] = false
		}
	} else {
		status["is_online"] = false
	}

	// 获取设备profile
	var profile models.PetProfile
	if err := c.DB.Where("device_id = ?", deviceID).First(&profile).Error; err == nil {
		status["pet_name"] = profile.PetName
		status["personality"] = profile.Personality
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    status,
	})
}

// SendMiniAppDeviceCommand 小程序发送设备指令
// POST /api/v1/miniapp/device/:device_id/command
func (c *MiniAppController) SendDeviceCommand(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	openIDVal, exists := ctx.Get("open_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	openID := openIDVal.(string)

	// 验证绑定
	var binding models.MiniAppDevice
	if err := c.DB.Where("device_id = ? AND open_id = ? AND is_active = ?", deviceID, openID, true).First(&binding).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "该设备未绑定到您的小程序"})
		return
	}

	var req struct {
		CmdType string                 `json:"cmd_type" binding:"required"`
		Action  string                 `json:"action"`
		Display map[string]interface{} `json:"display"`
		Config  map[string]interface{} `json:"config"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查设备在线状态
	if c.Redis != nil {
		shadow, _ := c.Redis.GetDeviceShadow(deviceID)
		if shadow == nil || !shadow.IsOnline {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "设备离线，无法下发指令"})
			return
		}
	}

	cmdID := uuid.New().String()
	cmd := map[string]interface{}{
		"cmd_id":    cmdID,
		"cmd_type":  req.CmdType,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	switch req.CmdType {
	case "action":
		cmd["action"] = req.Action
	case "display":
		cmd["display"] = req.Display
	case "config":
		cmd["config"] = req.Config
	default:
		cmd["action"] = req.Action
	}

	// 下发MQTT指令
	if mqtt.GlobalMQTTClient != nil {
		if err := mqtt.PublishCommand(mqtt.GlobalMQTTClient, deviceID, req.CmdType, cmd); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "指令下发失败"})
			return
		}
	}

	// 保存指令历史
	cmdHistory := models.CommandHistory{
		DeviceID: deviceID,
		CmdID:    cmdID,
		CmdType:  req.CmdType,
		Action:   req.Action,
		Status:   "sent",
		SentAt:   time.Now(),
	}
	c.DB.Create(&cmdHistory)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"cmd_id":  cmdID,
			"status":  "sent",
			"sent_at": time.Now(),
		},
	})
}

// ==================== 小程序推送 API ====================

// SendMiniAppPushRequest 小程序推送请求
type SendMiniAppPushRequest struct {
	DeviceID string                 `json:"device_id" binding:"required"`
	Title    string                 `json:"title" binding:"required"`
	Body     string                 `json:"body" binding:"required"`
	PushType string                 `json:"push_type" binding:"required"`
	Data     map[string]interface{} `json:"data"`
}

// SendMiniAppPush 发送小程序推送（通过微信订阅消息）
// POST /api/v1/miniapp/push
func (c *MiniAppController) SendPush(ctx *gin.Context) {
	openIDVal, exists := ctx.Get("open_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	openID := openIDVal.(string)

	var req SendMiniAppPushRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证绑定
	var binding models.MiniAppDevice
	if err := c.DB.Where("device_id = ? AND open_id = ? AND is_active = ?", req.DeviceID, openID, true).First(&binding).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "该设备未绑定"})
		return
	}

	pushID := uuid.New().String()
	dataJSON, _ := json.Marshal(req.Data)

	push := models.AppPush{
		PushID:   pushID,
		Platform: "miniapp",
		ClientID: openID,
		DeviceID: req.DeviceID,
		Title:    req.Title,
		Body:     req.Body,
		PushType: req.PushType,
		Data:     string(dataJSON),
		Status:   models.PushStatusPending,
	}
	if err := c.DB.Create(&push).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建推送记录失败"})
		return
	}

	// TODO: 调用微信订阅消息API实际推送
	// POST https://api.weixin.qq.com/cgi-bin/message/subscribe/send
	// access_token=ACCESS_TOKEN&body=JSON

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"push_id": pushID,
			"status":  "pending",
		},
	})
}
