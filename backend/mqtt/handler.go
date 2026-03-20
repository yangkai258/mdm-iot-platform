package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

// InitMQTT 初始化 MQTT 客户端
func InitMQTT(db *gorm.DB, redisClient *utils.RedisClient, alertCB AlertCallback, complianceCB ComplianceCallback, geofenceCB GeofenceCallback) (mqtt.Client, error) {
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("mdm-backend")
	opts.SetUsername("admin")
	opts.SetPassword("public")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	client := mqtt.NewClient(opts)
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("MQTT 连接失败: %w", token.Error())
	}

	log.Printf("[MQTT] 已连接到: %s", broker)

	// 设置订阅处理器（传入告警回调、合规回调和地理围栏回调）
	handler := NewHandler(db, redisClient, alertCB, complianceCB, geofenceCB)
	handler.SetupSubscriber(client)
	handler.StartHeartbeatChecker()

	return client, nil
}

// MQTTConfig MQTT 配置
type MQTTConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

// AlertCallback 告警回调接口，避免循环导入
type AlertCallback func(deviceID string, data map[string]interface{})

// ComplianceCallback 合规检查回调接口
type ComplianceCallback func(db *gorm.DB, deviceID string, data map[string]interface{})

// GeofenceCallback 地理围栏检查回调接口
type GeofenceCallback func(db *gorm.DB, deviceID string, lat, lng float64, alertType string)

// Handler MQTT 消息处理
type Handler struct {
	DB           *gorm.DB
	Redis        *utils.RedisClient
	AlertCB      AlertCallback
	ComplianceCB ComplianceCallback
	GeofenceCB   GeofenceCallback
}

// GlobalMQTTClient 全局 MQTT 客户端，供其他包注入使用
var GlobalMQTTClient mqtt.Client

// SetGlobalMQTTClient 设置全局 MQTT 客户端
func SetGlobalMQTTClient(client mqtt.Client) {
	GlobalMQTTClient = client
}

// StatusPayload 心跳上报 JSON 结构
type StatusPayload struct {
	DeviceID         string  `json:"device_id"`
	Timestamp        string  `json:"timestamp"`
	ConnectionStatus string  `json:"connection_status"` // online, offline, poor_network
	BatteryLevel     int     `json:"battery_level"`     // 0-100
	ChargingStatus   bool    `json:"charging_status"`
	CurrentMode      string  `json:"current_mode"` // sleeping, roaming, listening, talking, idle
	RSSI             *int    `json:"rssi,omitempty"`
	// 越狱/ROOT检测
	IsJailbroken bool   `json:"is_jailbroken"`
	RootStatus   string `json:"root_status"` // normal, rooted, jailbroken
	// 地理位置
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// PropertyPayload 属性上报 JSON 结构
type PropertyPayload struct {
	DeviceID        string `json:"device_id"`
	FirmwareVersion string `json:"firmware_version"`
	HardwareModel   string `json:"hardware_model"`
	LastIPAddress   string `json:"last_ip_address,omitempty"`
}

// NewHandler 创建 MQTT 处理器
func NewHandler(db *gorm.DB, redisClient *utils.RedisClient, alertCB AlertCallback, complianceCB ComplianceCallback, geofenceCB GeofenceCallback) *Handler {
	return &Handler{
		DB:           db,
		Redis:        redisClient,
		AlertCB:      alertCB,
		ComplianceCB: complianceCB,
		GeofenceCB:   geofenceCB,
	}
}

// SetupSubscriber 设置订阅回调
func (h *Handler) SetupSubscriber(client mqtt.Client) {
	// 订阅设备心跳上报主题
	topic := "/mdm/device/+/up/status"
	token := client.Subscribe(topic, 0, h.StatusMessageHandler)
	if token.Wait() && token.Error() != nil {
		log.Printf("[MQTT] 订阅 %s 失败: %v", topic, token.Error())
	} else {
		log.Printf("[MQTT] 成功订阅: %s", topic)
	}

	// 订阅设备属性上报主题
	topicProp := "/mdm/device/+/up/property"
	tokenProp := client.Subscribe(topicProp, 0, h.PropertyMessageHandler)
	if tokenProp.Wait() && tokenProp.Error() != nil {
		log.Printf("[MQTT] 订阅 %s 失败: %v", topicProp, tokenProp.Error())
	} else {
		log.Printf("[MQTT] 成功订阅: %s", topicProp)
	}

	// 订阅越狱/ROOT告警主题
	topicJail := "/mdm/device/+/up/jailbreak_alert"
	tokenJail := client.Subscribe(topicJail, 0, h.JailbreakAlertHandler)
	if tokenJail.Wait() && tokenJail.Error() != nil {
		log.Printf("[MQTT] 订阅 %s 失败: %v", topicJail, tokenJail.Error())
	} else {
		log.Printf("[MQTT] 成功订阅: %s", topicJail)
	}
}

// StatusMessageHandler 心跳消息处理
func (h *Handler) StatusMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] 收到心跳消息: %s", msg.Topic())

	var payload StatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("[MQTT] 解析心跳JSON失败: %v", err)
		return
	}

	// 解析时间戳
	timestamp, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	// 在线状态判断
	isOnline := payload.ConnectionStatus == "online"
	if payload.ConnectionStatus == "poor_network" {
		// 弱网也算在线
		isOnline = true
	}

	// 更新 Redis 设备影子 (TTL 90秒)
	shadow := utils.DeviceShadow{
		DeviceID:      payload.DeviceID,
		IsOnline:      isOnline,
		BatteryLevel:  payload.BatteryLevel,
		CurrentMode:   payload.CurrentMode,
		LastHeartbeat: &timestamp,
		IsJailbroken: payload.IsJailbroken,
		RootStatus:   payload.RootStatus,
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
	}

	if err := h.Redis.SetDeviceShadow(payload.DeviceID, shadow, 90*time.Second); err != nil {
		log.Printf("[MQTT] 更新设备影子失败: %v", err)
		return
	}

	// 同步更新 PostgreSQL 设备影子表
	h.syncShadowToDB(payload.DeviceID, isOnline, payload.BatteryLevel, payload.CurrentMode, &timestamp,
		payload.IsJailbroken, payload.RootStatus, payload.Latitude, payload.Longitude)

	log.Printf("[MQTT] 设备 %s 心跳更新: online=%v, battery=%d%%, mode=%s, jailbreak=%v, root=%s, lat=%.4f, lng=%.4f",
		payload.DeviceID, isOnline, payload.BatteryLevel, payload.CurrentMode,
		payload.IsJailbroken, payload.RootStatus, payload.Latitude, payload.Longitude)

	// 构建告警检查数据
	alertData := map[string]interface{}{
		"battery":          float64(payload.BatteryLevel),
		"is_online":        isOnline,
		"mode":             payload.CurrentMode,
		"battery_low":      payload.BatteryLevel < 15,
		"battery_critical": payload.BatteryLevel < 5,
		"is_jailbroken":    payload.IsJailbroken,
		"root_status":      payload.RootStatus,
	}

	// 触发告警检查（包含越狱检测）
	if h.AlertCB != nil {
		h.AlertCB(payload.DeviceID, alertData)
	}

	// 触发合规检查
	if h.ComplianceCB != nil {
		h.ComplianceCB(h.DB, payload.DeviceID, alertData)
	}

	// 触发地理围栏检查（当有有效位置时）
	if h.GeofenceCB != nil && payload.Latitude != 0 && payload.Longitude != 0 {
		h.GeofenceCB(h.DB, payload.DeviceID, payload.Latitude, payload.Longitude, "")
	}
}

// JailbreakAlertPayload 越狱告警消息结构
type JailbreakAlertPayload struct {
	DeviceID    string `json:"device_id"`
	Timestamp   string `json:"timestamp"`
	IsJailbreak bool   `json:"is_jailbreak"`
	RootStatus  string `json:"root_status"` // normal, rooted, jailbroken
	Details     string `json:"details"`     // 越狱详情描述
}

// JailbreakAlertHandler 处理越狱/ROOT告警
func (h *Handler) JailbreakAlertHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] 收到越狱告警: %s", msg.Topic())

	var payload JailbreakAlertPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("[MQTT] 解析越狱告警JSON失败: %v", err)
		return
	}

	log.Printf("[MQTT] 设备 %s 越狱检测: is_jailbreak=%v, root_status=%s, details=%s",
		payload.DeviceID, payload.IsJailbreak, payload.RootStatus, payload.Details)

	// 更新设备影子的越狱状态
	shadow, err := h.Redis.GetDeviceShadow(payload.DeviceID)
	if err == nil && shadow != nil {
		shadow.IsJailbroken = payload.IsJailbreak
		shadow.RootStatus = payload.RootStatus
		h.Redis.SetDeviceShadow(payload.DeviceID, *shadow, 90*time.Second)
	}

	// 如果检测到越狱，创建告警
	if payload.IsJailbreak || payload.RootStatus == "rooted" || payload.RootStatus == "jailbroken" {
		extraData, _ := json.Marshal(map[string]interface{}{
			"root_status": payload.RootStatus,
			"details":     payload.Details,
			"source":      "jailbreak_alert",
		})

		alert := models.DeviceAlert{
			RuleID:     0,
			DeviceID:   payload.DeviceID,
			AlertType:  "jailbreak_detected",
			Severity:   4, // 严重
			Message:    fmt.Sprintf("设备越狱/Root检测告警: %s", payload.Details),
			TriggerVal: 1,
			Threshold:  0,
			Status:     1,
			ExtraData:  string(extraData),
		}
		h.DB.Create(&alert)
		log.Printf("[MQTT] 设备 %s 越狱告警已创建, AlertID=%d", payload.DeviceID, alert.ID)

		// 触发越狱告警通知
		if h.AlertCB != nil {
			h.AlertCB(payload.DeviceID, map[string]interface{}{
				"is_jailbroken": payload.IsJailbreak,
				"root_status":   payload.RootStatus,
				"alert_type":    "jailbreak_detected",
			})
		}
	}
}

// CheckGeofence 检查设备是否触发地理围栏
func CheckGeofence(db *gorm.DB, deviceID string, lat, lng float64, alertType string) {
	var rules []models.GeofenceRule
	db.Where("enabled = ? AND (device_id = ? OR device_id = '')", true, deviceID).Find(&rules)

	for _, rule := range rules {
		// 计算距离（米）
		distance := haversineDistance(lat, lng, rule.CenterLat, rule.CenterLng)
		inside := distance <= rule.RadiusMeters

		var triggered bool
		var eventType string
		switch rule.AlertOn {
		case "enter":
			triggered = inside
			eventType = "enter"
		case "exit":
			triggered = !inside
			eventType = "exit"
		case "both":
			triggered = true
			if alertType != "" {
				eventType = alertType
			} else {
				eventType = "enter"
			}
		}

		if triggered {
			extraData, _ := json.Marshal(map[string]interface{}{
				"rule_id":        rule.ID,
				"rule_name":      rule.Name,
				"distance_m":     distance,
				"radius_m":       rule.RadiusMeters,
				"center_lat":    rule.CenterLat,
				"center_lng":     rule.CenterLng,
				"current_lat":   lat,
				"current_lng":    lng,
				"event_type":     eventType,
			})

			alert := models.DeviceAlert{
				RuleID:     rule.ID,
				DeviceID:   deviceID,
				AlertType:  "geofence_violation",
				Severity:   rule.Severity,
				Message:    fmt.Sprintf("地理围栏告警[%s]: %s", eventType, rule.Name),
				TriggerVal: distance,
				Threshold:  rule.RadiusMeters,
				Status:     1,
				ExtraData:  string(extraData),
			}
			db.Create(&alert)

			// 创建地理围栏告警记录
			geoAlert := models.GeofenceAlert{
				RuleID:    rule.ID,
				DeviceID:  deviceID,
				AlertType: eventType,
				Latitude:  lat,
				Longitude: lng,
				Severity:  rule.Severity,
				Message:   fmt.Sprintf("设备%s地理围栏[%s]", eventType, rule.Name),
				Status:    1,
				AlertID:  alert.ID,
			}
			db.Create(&geoAlert)

			log.Printf("[Geofence] 设备 %s 触发地理围栏[%s], 距离=%.2fm, 半径=%.2fm",
				deviceID, rule.Name, distance, rule.RadiusMeters)
		}
	}
}

// haversineDistance 计算两点之间的球面距离（米）
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // 地球半径（米）
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// PropertyMessageHandler 属性消息处理
func (h *Handler) PropertyMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[MQTT] 收到属性消息: %s", msg.Topic())

	var payload PropertyPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("[MQTT] 解析属性JSON失败: %v", err)
		return
	}

	// 处理属性更新，更新设备固件版本等信息
	log.Printf("[MQTT] 设备 %s 属性更新: firmware=%s, model=%s",
		payload.DeviceID, payload.FirmwareVersion, payload.HardwareModel)

	// 更新设备影子中的 IP 信息
	if payload.LastIPAddress != "" {
		shadow, err := h.Redis.GetDeviceShadow(payload.DeviceID)
		if err == nil && shadow != nil {
			shadow.LastIP = payload.LastIPAddress
			h.Redis.SetDeviceShadow(payload.DeviceID, *shadow, 90*time.Second)
		}
	}
}

// syncShadowToDB 同步设备影子到 PostgreSQL
func (h *Handler) syncShadowToDB(deviceID string, isOnline bool, batteryLevel int, currentMode string, lastHeartbeat *time.Time, isJailbroken bool, rootStatus string, lat, lng float64) {
	if h.DB == nil {
		return
	}

	var shadow models.DeviceShadow
	result := h.DB.Where("device_id = ?", deviceID).First(&shadow)

	if result.Error == gorm.ErrRecordNotFound {
		// 记录不存在，创建新记录
		shadow = models.DeviceShadow{
			DeviceID:      deviceID,
			IsOnline:      isOnline,
			BatteryLevel:  batteryLevel,
			CurrentMode:   currentMode,
			LastHeartbeat: lastHeartbeat,
		}
		h.DB.Create(&shadow)
	} else if result.Error == nil {
		// 记录存在，更新字段（但 lifecycle_status 不变）
		updates := map[string]interface{}{
			"is_online":      isOnline,
			"battery_level":   batteryLevel,
			"current_mode":   currentMode,
			"last_heartbeat": lastHeartbeat,
		}
		if isJailbroken {
			updates["is_jailbroken"] = isJailbroken
			updates["root_status"] = rootStatus
		}
		if lat != 0 && lng != 0 {
			updates["latitude"] = lat
			updates["longitude"] = lng
		}
		h.DB.Model(&shadow).Updates(updates)
	}
}

// StartHeartbeatChecker 启动心跳检查器（检测离线设备）
func (h *Handler) StartHeartbeatChecker() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			h.checkOfflineDevices()
		}
	}()
}

// checkOfflineDevices 检查离线设备
func (h *Handler) checkOfflineDevices() {
	// 从 Redis 获取所有设备影子
	keys, err := h.Redis.GetAllShadowKeys()
	if err != nil {
		return
	}

	now := time.Now()
	for _, key := range keys {
		shadow, err := h.Redis.GetDeviceShadow(key)
		if err != nil || shadow == nil {
			continue
		}

		// 检查心跳超时 (90秒)
		if shadow.LastHeartbeat != nil {
			elapsed := now.Sub(*shadow.LastHeartbeat)
			if elapsed > 90*time.Second && shadow.IsOnline {
				// 标记为离线
				shadow.IsOnline = false
				h.Redis.SetDeviceShadow(shadow.DeviceID, *shadow, 0)

				// 同步离线状态到 DB（只更新 last_heartbeat，不改 lifecycle_status）
				h.syncOfflineToDB(shadow.DeviceID)

				log.Printf("[MQTT] 设备 %s 心跳超时，标记为离线", shadow.DeviceID)

				// 触发离线告警检查
				if h.AlertCB != nil {
					h.AlertCB(shadow.DeviceID, map[string]interface{}{
						"is_online": false,
						"offline":   true,
						"elapsed":   elapsed.Seconds(),
					})
				}

				// 触发合规检查
				if h.ComplianceCB != nil {
					h.ComplianceCB(h.DB, shadow.DeviceID, map[string]interface{}{
						"is_online":   false,
						"offline":     true,
						"elapsed":     elapsed.Seconds(),
						"battery":     float64(shadow.BatteryLevel),
					})
				}
			}
		}
	}
}

// syncOfflineToDB 同步离线状态到 DB（只更新 last_heartbeat）
func (h *Handler) syncOfflineToDB(deviceID string) {
	if h.DB == nil {
		return
	}

	now := time.Now()
	h.DB.Model(&models.DeviceShadow{}).Where("device_id = ?", deviceID).Updates(map[string]interface{}{
		"is_online":     false,
		"last_heartbeat": now,
	})
}

// PublishCommand 下发指令到设备
func (h *Handler) PublishCommand(client mqtt.Client, deviceID string, cmd interface{}) error {
	topic := fmt.Sprintf("/mdm/device/%s/down/cmd", deviceID)
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	token := client.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Printf("[MQTT] 向设备 %s 下发指令: %s", deviceID, string(payload))
	return nil
}
