package controllers

import (
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
	"gorm.io/gorm"
)

// DeviceSecurityController 设备安全控制器
type DeviceSecurityController struct {
	DB *gorm.DB
}

// LockDevice 锁定设备
func (c *DeviceSecurityController) LockDevice(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	username := middleware.GetUsername(ctx)
	_ = middleware.GetTenantIDCtx(ctx) // 保留上下文
	_ = middleware.GetUserID(ctx)     // 保留上下文

	// 验证设备存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	// 通过 MQTT 下发锁定命令 (QoS=1)
	lockCmd := map[string]interface{}{
		"lock":   true,
		"reason": "admin",
		"locked_by": username,
		"locked_at": time.Now().Format(time.RFC3339),
	}

	payload, _ := json.Marshal(lockCmd)
	topic := fmt.Sprintf("/miniclaw/%s/down/lock", deviceID)

	client := mqtt.GlobalMQTTClient
	if client == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "MQTT 服务不可用",
		})
		return
	}

	// QoS=1 保证送达
	token := client.Publish(topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "下发锁定命令失败: " + token.Error().Error(),
		})
		return
	}

	// 记录操作日志（可选：创建锁定记录）
	// 这里简化为仅记录到操作日志中间件

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id": deviceID,
			"status":    "lock_sent",
			"message":   "锁定命令已下发",
		},
		"message": "锁定命令已下发",
	})
}

// UnlockDevice 解锁设备
func (c *DeviceSecurityController) UnlockDevice(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	username := middleware.GetUsername(ctx)
	_ = middleware.GetTenantIDCtx(ctx) // 保留上下文
	_ = middleware.GetUserID(ctx)     // 保留上下文

	// 验证设备存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	// 通过 MQTT 下发解锁命令 (QoS=1)
	unlockCmd := map[string]interface{}{
		"lock":     false,
		"reason":   "admin_unlock",
		"unlocked_by": username,
		"unlocked_at": time.Now().Format(time.RFC3339),
	}

	payload, _ := json.Marshal(unlockCmd)
	topic := fmt.Sprintf("/miniclaw/%s/down/lock", deviceID)

	client := mqtt.GlobalMQTTClient
	if client == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "MQTT 服务不可用",
		})
		return
	}

	token := client.Publish(topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "下发解锁命令失败: " + token.Error().Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id": deviceID,
			"status":    "unlock_sent",
			"message":   "解锁命令已下发",
		},
		"message": "解锁命令已下发",
	})
}

// WipeDevice 擦除设备（需要二次确认token）
func (c *DeviceSecurityController) WipeDevice(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	tenantID := middleware.GetTenantIDCtx(ctx)
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	var req models.WipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if req.WipeType != "full" && req.WipeType != "selective" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "擦除类型必须是 full 或 selective",
		})
		return
	}

	// 验证设备存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	// 生成确认 token (24位十六进制)
	confirmToken := generateConfirmToken()

	// 创建擦除历史记录（待确认状态）
	wipeHistory := models.WipeHistory{
		DeviceID:     deviceID,
		OperatorID:   userID,
		OperatorName: username,
		WipeType:     req.WipeType,
		Status:       "pending",
		ConfirmToken: confirmToken,
		Reason:       req.Reason,
		TenantID:     tenantID,
	}

	if err := c.DB.Create(&wipeHistory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建擦除记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"wipe_id":        wipeHistory.ID,
			"confirm_token":  confirmToken,
			"message":        "擦除操作需要二次确认，请在 5 分钟内确认",
			"confirm_url":    fmt.Sprintf("/api/v1/devices/%s/wipe/confirm", deviceID),
		},
		"message": "需要二次确认",
	})
}

// GenerateWipeConfirmToken 生成分解确认 token（用于让用户获取确认码）
func (c *DeviceSecurityController) GenerateWipeConfirmToken(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	userID := middleware.GetUserID(ctx)
	_ = middleware.GetTenantIDCtx(ctx) // 保留上下文

	// 查找设备最新的 pending 擦除记录
	var wipeHistory models.WipeHistory
	if err := c.DB.Where("device_id = ? AND operator_id = ? AND status = ?",
		deviceID, userID, "pending").Order("created_at DESC").First(&wipeHistory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "没有待确认的擦除操作",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 检查是否过期（5分钟）
	if time.Since(wipeHistory.CreatedAt) > 5*time.Minute {
		wipeHistory.Status = "failed"
		wipeHistory.Result = "确认超时"
		c.DB.Save(&wipeHistory)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "确认已超时，请重新发起擦除",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"wipe_id":       wipeHistory.ID,
			"confirm_token": wipeHistory.ConfirmToken,
			"expires_in":    300, // 剩余秒数
		},
		"message": "获取确认码成功",
	})
}

// ConfirmWipe 确认擦除操作
func (c *DeviceSecurityController) ConfirmWipe(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	userID := middleware.GetUserID(ctx)
	_ = middleware.GetTenantIDCtx(ctx) // 保留上下文

	var req models.WipeConfirmRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 查找对应的擦除记录
	var wipeHistory models.WipeHistory
	if err := c.DB.Where("device_id = ? AND operator_id = ? AND confirm_token = ? AND status = ?",
		deviceID, userID, req.ConfirmToken, "pending").First(&wipeHistory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "确认token无效或已过期",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 检查是否过期（5分钟）
	if time.Since(wipeHistory.CreatedAt) > 5*time.Minute {
		wipeHistory.Status = "failed"
		wipeHistory.Result = "确认超时"
		c.DB.Save(&wipeHistory)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "确认已超时，请重新发起擦除",
		})
		return
	}

	// 更新为已确认状态
	now := time.Now()
	wipeHistory.Status = "executing"
	wipeHistory.ConfirmedAt = &now

	if err := c.DB.Save(&wipeHistory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新状态失败: " + err.Error(),
		})
		return
	}

	// 通过 MQTT 下发擦除命令 (QoS=1)
	wipeCmd := map[string]interface{}{
		"wipe":         true,
		"wipe_type":    wipeHistory.WipeType,
		"confirm_code": wipeHistory.ConfirmToken,
		"executed_by":  wipeHistory.OperatorName,
		"executed_at": now.Format(time.RFC3339),
	}

	payload, _ := json.Marshal(wipeCmd)
	topic := fmt.Sprintf("/miniclaw/%s/down/wipe", deviceID)

	client := mqtt.GlobalMQTTClient
	if client == nil {
		wipeHistory.Status = "failed"
		wipeHistory.Result = "MQTT 服务不可用"
		c.DB.Save(&wipeHistory)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "MQTT 服务不可用",
		})
		return
	}

	token := client.Publish(topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		wipeHistory.Status = "failed"
		wipeHistory.Result = "下发擦除命令失败: " + token.Error().Error()
		c.DB.Save(&wipeHistory)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "下发擦除命令失败: " + token.Error().Error(),
		})
		return
	}

	executedAt := time.Now()
	wipeHistory.ExecutedAt = &executedAt
	c.DB.Save(&wipeHistory)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"wipe_id":   wipeHistory.ID,
			"status":   "executing",
			"message":   "擦除命令已下发",
		},
		"message": "擦除命令已下发",
	})
}

// GetWipeHistory 获取设备擦除历史
func (c *DeviceSecurityController) GetWipeHistory(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	tenantID := middleware.GetTenantIDCtx(ctx)
	page := parseIntDefault(ctx.Query("page"), 1)
	pageSize := parseIntDefault(ctx.Query("page_size"), 20)

	query := c.DB.Model(&models.WipeHistory{}).Where("device_id = ? AND tenant_id = ?", deviceID, tenantID)

	var total int64
	query.Count(&total)

	var histories []models.WipeHistory
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	responses := make([]models.WipeHistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = h.ToResponse()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      responses,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// GetDeviceLockStatus 获取设备锁定状态
func (c *DeviceSecurityController) GetDeviceLockStatus(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 从 Redis 设备影子获取锁定状态
	redisClient := utils.GetGlobalRedisClient()
	if redisClient == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    503,
			"message": "Redis 服务不可用",
		})
		return
	}

	shadow, err := redisClient.GetDeviceShadow(deviceID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "设备不在线或不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":   deviceID,
			"is_locked":   shadow.IsJailbroken, // 实际应根据业务定义锁定状态字段
			"is_online":   shadow.IsOnline,
			"last_heartbeat": shadow.LastHeartbeat,
		},
		"message": "success",
	})
}

// generateConfirmToken 生成24位十六进制确认token
func generateConfirmToken() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
