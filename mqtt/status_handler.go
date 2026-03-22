package mqtt

import (
	"encoding/json"
	"log"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

// PetStatusPayload MiniClaw宠物状态上报Payload
type PetStatusPayload struct {
	DeviceID          string  `json:"device_id"`
	Timestamp         string  `json:"timestamp"`
	Mood              int     `json:"mood"`
	Energy            int     `json:"energy"`
	Hunger            int     `json:"hunger"`
	PositionX         float64 `json:"position_x"`
	PositionY         float64 `json:"position_y"`
	CurrentExpression string  `json:"current_expression"`
	CurrentAction     string  `json:"current_action"`
	BatteryLevel      int     `json:"battery_level"`
	IsOnline          bool    `json:"is_online"`
}

// StatusHandler MiniClaw状态消息处理器
type StatusHandler struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// NewStatusHandler 创建状态处理器
func NewStatusHandler(db *gorm.DB, redisClient *utils.RedisClient) *StatusHandler {
	return &StatusHandler{
		DB:    db,
		Redis: redisClient,
	}
}

// SetupMiniClawSubscriber 设置MiniClaw状态订阅
func (h *StatusHandler) SetupMiniClawSubscriber(client mqtt.Client) {
	// 订阅宠物状态上报主题
	topic := "/miniclaw/+/up/status"
	token := client.Subscribe(topic, 0, h.PetStatusMessageHandler)
	if token.Wait() && token.Error() != nil {
		log.Printf("[MQTT] 订阅 %s 失败: %v", topic, token.Error())
	} else {
		log.Printf("[MQTT] 成功订阅: %s", topic)
	}
}

// PetStatusMessageHandler 处理宠物状态上报消息
func (h *StatusHandler) PetStatusMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] 收到宠物状态消息: %s", msg.Topic())

	var payload PetStatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("[MQTT] 解析宠物状态JSON失败: %v", err)
		return
	}

	// 解析时间戳
	timestamp, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	// 更新或创建宠物状态
	var petStatus models.PetStatusV2
	result := h.DB.Where("device_id = ?", payload.DeviceID).First(&petStatus)

	updates := map[string]interface{}{
		"mood":               payload.Mood,
		"energy":             payload.Energy,
		"hunger":             payload.Hunger,
		"position_x":         payload.PositionX,
		"position_y":         payload.PositionY,
		"current_expression": payload.CurrentExpression,
		"current_action":     payload.CurrentAction,
		"is_online":          payload.IsOnline,
		"last_seen_at":       timestamp,
		"updated_at":         time.Now(),
	}

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		petStatus = models.PetStatusV2{
			DeviceID:          payload.DeviceID,
			Mood:              payload.Mood,
			Energy:            payload.Energy,
			Hunger:            payload.Hunger,
			PositionX:         payload.PositionX,
			PositionY:         payload.PositionY,
			CurrentExpression: payload.CurrentExpression,
			CurrentAction:     payload.CurrentAction,
			IsOnline:          payload.IsOnline,
			LastSeenAt:        &timestamp,
		}
		if err := h.DB.Create(&petStatus).Error; err != nil {
			log.Printf("[MQTT] 创建宠物状态失败: %v", err)
			return
		}
		log.Printf("[MQTT] 新建宠物状态: device_id=%s, mood=%d, energy=%d", payload.DeviceID, payload.Mood, payload.Energy)
	} else if result.Error == nil {
		// 更新现有记录
		if err := h.DB.Model(&petStatus).Updates(updates).Error; err != nil {
			log.Printf("[MQTT] 更新宠物状态失败: %v", err)
			return
		}
		log.Printf("[MQTT] 更新宠物状态: device_id=%s, mood=%d, energy=%d, expression=%s",
			payload.DeviceID, payload.Mood, payload.Energy, payload.CurrentExpression)
	} else {
		log.Printf("[MQTT] 查询宠物状态失败: %v", result.Error)
		return
	}

	// 同步更新Redis设备影子
	if h.Redis != nil {
		shadow := utils.DeviceShadow{
			DeviceID:      payload.DeviceID,
			IsOnline:      payload.IsOnline,
			BatteryLevel:  payload.BatteryLevel,
			CurrentMode:   payload.CurrentAction,
			LastHeartbeat: &timestamp,
		}
		if err := h.Redis.SetDeviceShadow(payload.DeviceID, shadow, 90*time.Second); err != nil {
			log.Printf("[MQTT] 更新Redis设备影子失败: %v", err)
		}
	}
}

// HandleDeviceOnline 处理设备上线
func (h *StatusHandler) HandleDeviceOnline(deviceID string) {
	var petStatus models.PetStatusV2
	if err := h.DB.Where("device_id = ?", deviceID).First(&petStatus).Error; err == nil {
		h.DB.Model(&petStatus).Updates(map[string]interface{}{
			"is_online":    true,
			"last_seen_at": time.Now(),
		})
	}
}

// HandleDeviceOffline 处理设备离线
func (h *StatusHandler) HandleDeviceOffline(deviceID string) {
	var petStatus models.PetStatusV2
	if err := h.DB.Where("device_id = ?", deviceID).First(&petStatus).Error; err == nil {
		h.DB.Model(&petStatus).Updates(map[string]interface{}{
			"is_online":  false,
			"updated_at": time.Now(),
		})
	}
}
