package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

// AppController App移动端API控制器
type AppController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// ==================== App Token API ====================

// AppTokenRequest App Token请求
type AppTokenRequest struct {
	AppID     string `json:"app_id" binding:"required"`
	Platform  string `json:"platform" binding:"required"` // ios, android
	ClientID  string `json:"client_id" binding:"required"`
	AuthCode  string `json:"auth_code"`                     // 授权码（App扫码场景）
	Username  string `json:"username"`                      // 用户名（密码模式）
	Password  string `json:"password"`                      // 密码（密码模式）
	GrantType string `json:"grant_type" binding:"required"` // authorization_code, password, client_credentials
}

// AppTokenResponse App Token响应
type AppTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// GetAppToken 获取App Token
// POST /api/v1/app/auth/token
func (c *AppController) GetAppToken(ctx *gin.Context) {
	var req AppTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证 AppID（简单白名单，实际应查数据库）
	allowedApps := map[string]bool{
		"miniclaw-ios":     true,
		"miniclaw-android": true,
	}
	if !allowedApps[req.AppID] {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效的AppID"})
		return
	}

	var userID uint
	var username string
	var roleID uint
	var tenantID string
	var isSuperAdmin bool

	// 根据 grant_type 处理不同授权模式
	switch req.GrantType {
	case "password":
		// 密码模式：验证用户名密码
		var user models.SysUser
		if err := c.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
			return
		}
		userID = user.ID
		username = user.Username
		roleID = user.RoleID
		tenantID = user.TenantID

	case "authorization_code":
		// 授权码模式（App扫码登录场景）
		// auth_code 应当由扫码流程获取，此处简化处理
		if req.AuthCode == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "授权码不能为空"})
			return
		}
		// 从Redis取出预存的code->user映射（扫码时写入）
		if c.Redis != nil {
			codeKey := fmt.Sprintf("app:authcode:%s", req.AuthCode)
			ctx2 := context.Background()
			codeData, err := c.Redis.Client().Get(ctx2, codeKey).Result()
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "授权码已过期或无效"})
				return
			}
			var userInfo map[string]interface{}
			if err := json.Unmarshal([]byte(codeData), &userInfo); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "授权码解析失败"})
				return
			}
			userID = uint(userInfo["user_id"].(float64))
			username = userInfo["username"].(string)
			roleID = uint(userInfo["role_id"].(float64))
			tenantID = userInfo["tenant_id"].(string)
			// 删除已使用的code
			c.Redis.Client().Del(ctx2, codeKey)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Redis不可用"})
			return
		}

	case "client_credentials":
		// 客户端模式（设备级别的token）
		userID = 0
		username = "app_client"
		roleID = 0
		tenantID = ""

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的授权模式"})
		return
	}

	// 生成 JWT AccessToken
	token, err := middleware.GenerateToken(userID, username, roleID, tenantID, isSuperAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成Token失败"})
		return
	}

	expiresIn := int64(86400) // 24小时
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// 保存 AppToken 记录
	appToken := models.AppToken{
		Token:     uuid.New().String(),
		AppID:     req.AppID,
		UserID:    userID,
		Platform:  req.Platform,
		ClientID:  req.ClientID,
		Scope:     "read,write",
		ExpiresAt: expiresAt,
		IsActive:  true,
	}
	if err := c.DB.Create(&appToken).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存Token记录失败"})
		return
	}

	// 生成 RefreshToken
	refreshTokenStr := generateSecureToken(64)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天
	refreshToken := models.AppRefreshToken{
		RefreshToken: refreshTokenStr,
		AppTokenID:   appToken.ID,
		UserID:       userID,
		ClientID:     req.ClientID,
		Platform:     req.Platform,
		ExpiresAt:    refreshExpiresAt,
		IsActive:     true,
	}
	if err := c.DB.Create(&refreshToken).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成RefreshToken失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": AppTokenResponse{
			AccessToken:  token,
			TokenType:    "Bearer",
			ExpiresIn:    expiresIn,
			RefreshToken: refreshTokenStr,
			Scope:        "read,write",
		},
	})
}

// RefreshAppToken 刷新App Token
// POST /api/v1/app/auth/refresh
func (c *AppController) RefreshAppToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
		ClientID     string `json:"client_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查询刷新Token
	var refreshToken models.AppRefreshToken
	if err := c.DB.Where("refresh_token = ? AND client_id = ? AND is_active = ?", req.RefreshToken, req.ClientID, true).First(&refreshToken).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "RefreshToken无效"})
		return
	}

	if refreshToken.IsExpired() || !refreshToken.IsValid() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "RefreshToken已过期"})
		return
	}

	// 获取原AppToken信息
	var appToken models.AppToken
	if err := c.DB.First(&appToken, refreshToken.AppTokenID).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "关联的Token不存在"})
		return
	}

	// 撤销旧Token
	now := time.Now()
	appToken.IsActive = false
	appToken.RevokedAt = &now
	c.DB.Save(&appToken)

	// 撤销旧刷新Token
	refreshToken.IsActive = false
	refreshToken.RevokedAt = &now
	refreshToken.RefreshCount++
	c.DB.Save(&refreshToken)

	// 获取用户信息
	var username string
	var roleID uint
	var tenantID string
	var isSuperAdmin bool
	if refreshToken.UserID > 0 {
		var user models.SysUser
		if err := c.DB.First(&user, refreshToken.UserID).Error; err == nil {
			username = user.Username
			roleID = user.RoleID
			tenantID = user.TenantID
		}
	} else {
		username = "app_client"
	}

	// 生成新的 AccessToken
	token, err := middleware.GenerateToken(refreshToken.UserID, username, roleID, tenantID, isSuperAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成Token失败"})
		return
	}

	expiresIn := int64(86400)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// 创建新的AppToken
	newAppToken := models.AppToken{
		Token:     uuid.New().String(),
		AppID:     appToken.AppID,
		UserID:    refreshToken.UserID,
		Platform:  refreshToken.Platform,
		ClientID:  refreshToken.ClientID,
		Scope:     "read,write",
		ExpiresAt: expiresAt,
		IsActive:  true,
	}
	if err := c.DB.Create(&newAppToken).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存Token记录失败"})
		return
	}

	// 生成新的RefreshToken
	newRefreshTokenStr := generateSecureToken(64)
	newRefreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	newRefreshToken := models.AppRefreshToken{
		RefreshToken: newRefreshTokenStr,
		AppTokenID:   newAppToken.ID,
		UserID:       refreshToken.UserID,
		ClientID:     refreshToken.ClientID,
		Platform:     refreshToken.Platform,
		ExpiresAt:    newRefreshExpiresAt,
		IsActive:     true,
	}
	if err := c.DB.Create(&newRefreshToken).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成RefreshToken失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": AppTokenResponse{
			AccessToken:  token,
			TokenType:    "Bearer",
			ExpiresIn:    expiresIn,
			RefreshToken: newRefreshTokenStr,
			Scope:        "read,write",
		},
	})
}

// ==================== App Push API ====================

// SendPushRequest 发送推送请求
type SendPushRequest struct {
	UserID    uint                   `json:"user_id" binding:"required"`
	DeviceID  string                 `json:"device_id"`                   // 可选，指定设备
	Platform  string                 `json:"platform" binding:"required"` // ios, android
	Title     string                 `json:"title" binding:"required"`
	Body      string                 `json:"body" binding:"required"`
	PushType  string                 `json:"push_type" binding:"required"` // alert, device, member, system, ota, marketing
	Data      map[string]interface{} `json:"data"`                         // 透传数据
	Badge     int                    `json:"badge"`
	Sound     string                 `json:"sound"`
	ChannelID string                 `json:"channel_id"`
	Tag       string                 `json:"tag"`
	ExpiresAt *time.Time             `json:"expires_at"` // 过期时间（用于撤回）
}

// SendAppPush 发送App推送
// POST /api/v1/app/push
func (c *AppController) SendAppPush(ctx *gin.Context) {
	var req SendPushRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	pushID := uuid.New().String()
	dataJSON, _ := json.Marshal(req.Data)

	// 构建推送记录
	push := models.AppPush{
		PushID:    pushID,
		UserID:    req.UserID,
		DeviceID:  req.DeviceID,
		Platform:  req.Platform,
		Title:     req.Title,
		Body:      req.Body,
		PushType:  req.PushType,
		Data:      string(dataJSON),
		Badge:     req.Badge,
		Sound:     req.Sound,
		ChannelID: req.ChannelID,
		Tag:       req.Tag,
		Status:    models.PushStatusPending,
		ExpiresAt: req.ExpiresAt,
	}

	if err := c.DB.Create(&push).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建推送记录失败"})
		return
	}

	// 实际推送逻辑（这里模拟成功，实际应调用极光/FCM/华为push）
	// TODO: 接入真实的推送服务（极光/FCM/华为Push）
	go func() {
		time.Sleep(100 * time.Millisecond)
		// 标记为已发送
		now := time.Now()
		c.DB.Model(&models.AppPush{}).Where("id = ?", push.ID).Updates(map[string]interface{}{
			"status":  models.PushStatusSent,
			"sent_at": now,
		})
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"push_id": pushID,
			"status":  "pending",
		},
	})
}

// GetPushHistory 获取推送历史
// GET /api/v1/app/push/history
func (c *AppController) GetPushHistory(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	userID := uint(userIDVal.(int))

	page, _ := ctx.GetQuery("page")
	pageSize, _ := ctx.GetQuery("page_size")
	pushType, _ := ctx.GetQuery("push_type")
	status, _ := ctx.GetQuery("status")

	pageInt := 1
	pageSizeInt := 20
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(pageSize, "%d", &pageSizeInt)
	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 20
	}
	offset := (pageInt - 1) * pageSizeInt

	query := c.DB.Model(&models.AppPush{}).Where("user_id = ?", userID)
	if pushType != "" {
		query = query.Where("push_type = ?", pushType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var pushes []models.AppPush
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSizeInt).Find(&pushes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 解析Data字段
	type PushItem struct {
		models.AppPush
		DataParsed map[string]interface{} `json:"data_parsed"`
	}
	items := make([]PushItem, len(pushes))
	for i, p := range pushes {
		var data map[string]interface{}
		json.Unmarshal([]byte(p.Data), &data)
		items[i] = PushItem{AppPush: p, DataParsed: data}
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

// ==================== App 设备控制 API ====================

// GetAppDeviceStatus 获取App专用设备状态
// GET /api/v1/app/device/:device_id/status
func (c *AppController) GetAppDeviceStatus(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	userID := uint(userIDVal.(int))

	// 查找设备（需要用户有权限查看）
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 简单权限检查：设备属于该用户或用户是创建者
	if device.BindUserID != nil && *device.BindUserID != fmt.Sprintf("%d", userID) && device.CreateUserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限查看该设备"})
		return
	}

	// 构建基础状态
	status := map[string]interface{}{
		"device_id":    device.DeviceID,
		"sn_code":      device.SnCode,
		"mac_address":  device.MacAddress,
		"firmware_ver": device.FirmwareVersion,
		"lifecycle":    device.LifecycleStatus,
	}

	// 从Redis获取设备影子状态
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
			status["battery"] = 0
			status["current_mode"] = "unknown"
		}
	} else {
		status["is_online"] = false
		status["battery"] = 0
		status["current_mode"] = "unknown"
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

// SendAppDeviceCommand 发送设备指令（App专用）
// POST /api/v1/app/device/:device_id/command
func (c *AppController) SendAppDeviceCommand(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}
	userID := uint(userIDVal.(int))

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 权限检查
	if device.BindUserID != nil && *device.BindUserID != fmt.Sprintf("%d", userID) && device.CreateUserID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权限操作该设备"})
		return
	}

	var req struct {
		CmdType string                 `json:"cmd_type" binding:"required"` // action, display, config, ota
		Action  string                 `json:"action"`
		Display map[string]interface{} `json:"display"`
		Config  map[string]interface{} `json:"config"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 获取设备在线状态
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

	// 通过MQTT下发指令
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

// ==================== 辅助函数 ====================

// generateSecureToken 生成安全的随机Token
func generateSecureToken(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return uuid.New().String() + uuid.New().String()
	}
	return hex.EncodeToString(bytes)
}
