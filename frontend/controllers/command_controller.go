package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CommandController 指令控制
type CommandController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// SendCommandRequest 下发指令请求
type SendCommandRequest struct {
	CmdType string                 `json:"cmd_type" binding:"required"` // action, display, config, ota
	Action  string                 `json:"action"`
	Display map[string]interface{} `json:"display"`
	Config  map[string]interface{} `json:"config"`
	OTA     map[string]interface{} `json:"ota"`
}

// SendCommand 下发指令到设备
func (c *CommandController) SendCommand(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 获取设备在线状态
	shadow, err := c.Redis.GetDeviceShadow(deviceID)
	if err != nil || shadow == nil || !shadow.IsOnline {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4003,
			"message":   "设备离线，无法下发指令",
			"error_code": "ERR_DEVICE_003",
		})
		return
	}

	var req SendCommandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	cmdID := uuid.New().String()
	// 构建指令消息
	cmd := map[string]interface{}{
		"cmd_id":    cmdID,
		"cmd_type":  req.CmdType,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if req.Action != "" {
		cmd["action"] = req.Action
	}
	if req.Display != nil {
		cmd["display"] = req.Display
	}
	if req.Config != nil {
		cmd["config"] = req.Config
	}
	if req.OTA != nil {
		cmd["ota"] = req.OTA
	}

	// 通过 MQTT 下发
	if mqtt.GlobalMQTTClient != nil {
		topic := fmt.Sprintf("/mdm/device/%s/down/cmd", deviceID)
		payload, _ := json.Marshal(cmd)
		token := mqtt.GlobalMQTTClient.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "指令下发失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}
	}

	// 记录指令历史
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
			"cmd_id":    cmdID,
			"device_id": deviceID,
			"status":    "sent",
		},
	})
}

// GetCommandHistory 获取指令历史
func (c *CommandController) GetCommandHistory(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var history []models.CommandHistory
	if err := c.DB.Where("device_id = ?", deviceID).Order("sent_at DESC").Limit(50).Find(&history).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": history,
		},
	})
}
