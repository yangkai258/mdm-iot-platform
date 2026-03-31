package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UnbindDevice 解绑设备
func UnbindDevice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		snCode := ctx.Param("sn_code")

		var device models.Device
		result := db.Where("sn_code = ?", snCode).First(&device)
		
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":      4002,
				"message":   "设备不存在",
				"error_code": "ERR_DEVICE_002",
			})
			return
		}

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "服务器内部错误",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		// 解绑用户
		device.BindUserID = nil
		device.LifecycleStatus = 1 // 改为未绑定状态
		
		if err := db.Save(&device).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "解绑失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id": device.DeviceID,
				"sn_code":   device.SnCode,
				"status":    "unbound",
			},
		})
	}
}

// UpdateDeviceStatus 更新设备状态
func UpdateDeviceStatus(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceID := ctx.Param("device_id")

		var device models.Device
		result := db.Where("device_id = ?", deviceID).First(&device)

		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":      4002,
				"message":   "设备不存在",
				"error_code": "ERR_DEVICE_002",
			})
			return
		}

		// 解析请求
		var req struct {
			Status     int     `json:"status"`
			Battery    float64 `json:"battery"`     // 电池电量（可选）
			IsOnline   bool    `json:"is_online"`   // 在线状态（可选）
			IsJailbreak bool   `json:"is_jailbreak"` // 越狱状态（可选）
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "参数错误",
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		// 验证状态值 (1:待激活 2:服役中 3:维修中 4:已挂失 5:已报废)
		if req.Status < 1 || req.Status > 5 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "无效的状态值",
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		// 构建告警检查数据（如果有电池或在线状态数据）
		if req.Battery > 0 || req.IsOnline || req.IsJailbreak {
			alertData := map[string]interface{}{
				"battery":          req.Battery,
				"is_online":        req.IsOnline,
				"battery_low":      req.Battery > 0 && req.Battery < 15,
				"battery_critical": req.Battery > 0 && req.Battery < 5,
				"is_jailbroken":    req.IsJailbreak,
			}
			// 调用告警检查（异步避免阻塞）
			go CheckAlerts(db, deviceID, alertData)
		}

		device.LifecycleStatus = req.Status
		if err := db.Save(&device).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "更新失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":         device.DeviceID,
				"lifecycle_status": device.LifecycleStatus,
			},
		})
	}
}

// SetDesiredStateRequest 设置设备影子期望状态请求
type SetDesiredStateRequest struct {
	DesiredNRDEnabled bool   `json:"desired_nrd_enabled"`
	DesiredNRDStart   string `json:"desired_nrd_start"`
	DesiredNRDEnd     string `json:"desired_nrd_end"`
	DesiredDNDEnabled bool   `json:"desired_dnd_enabled"`
	DesiredDNDStart   string `json:"desired_dnd_start"`
	DesiredDNDEnd     string `json:"desired_dnd_end"`
	DesiredVolume     *int   `json:"desired_volume"`
	DesiredBrightness *int   `json:"desired_brightness"`
	DesiredPowerSave  bool   `json:"desired_power_save"`
	DesiredVersion    string `json:"desired_version"`
}

// SetDesiredState 设置设备影子期望状态（NRD/免打扰）
// PUT /api/v1/devices/:device_id/desired-state
func SetDesiredState(db *gorm.DB, mqttClient interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceID := ctx.Param("device_id")

		// 验证设备是否存在
		var device models.Device
		if err := db.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":      4002,
				"message":   "设备不存在",
				"error_code": "ERR_DEVICE_002",
			})
			return
		}

		var req SetDesiredStateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "参数错误: " + err.Error(),
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		// 验证时间格式
		if req.DesiredNRDStart != "" && !isValidTimeFormat(req.DesiredNRDStart) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "NRD开始时间格式无效，请使用 HH:MM"})
			return
		}
		if req.DesiredNRDEnd != "" && !isValidTimeFormat(req.DesiredNRDEnd) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "NRD结束时间格式无效，请使用 HH:MM"})
			return
		}
		if req.DesiredDNDStart != "" && !isValidTimeFormat(req.DesiredDNDStart) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "DND开始时间格式无效，请使用 HH:MM"})
			return
		}
		if req.DesiredDNDEnd != "" && !isValidTimeFormat(req.DesiredDNDEnd) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "DND结束时间格式无效，请使用 HH:MM"})
			return
		}
		if req.DesiredVolume != nil && (*req.DesiredVolume < 0 || *req.DesiredVolume > 100) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "音量必须在 0-100 之间"})
			return
		}
		if req.DesiredBrightness != nil && (*req.DesiredBrightness < 0 || *req.DesiredBrightness > 100) {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "亮度必须在 0-100 之间"})
			return
		}

		// 更新或创建设备影子记录
		updates := map[string]interface{}{
			"desired_nrd_enabled":  req.DesiredNRDEnabled,
			"desired_nrd_start":    req.DesiredNRDStart,
			"desired_nrd_end":      req.DesiredNRDEnd,
			"desired_dnd_enabled":  req.DesiredDNDEnabled,
			"desired_dnd_start":    req.DesiredDNDStart,
			"desired_dnd_end":      req.DesiredDNDEnd,
			"desired_power_save":   req.DesiredPowerSave,
			"desired_version":      req.DesiredVersion,
		}
		if req.DesiredVolume != nil {
			updates["desired_volume"] = *req.DesiredVolume
		}
		if req.DesiredBrightness != nil {
			updates["desired_brightness"] = *req.DesiredBrightness
		}

		var shadow models.DeviceShadow
		result := db.Where("device_id = ?", deviceID).First(&shadow)
		if result.Error == gorm.ErrRecordNotFound {
			shadow = models.DeviceShadow{DeviceID: deviceID}
			for k, v := range updates {
				switch k {
				case "desired_nrd_enabled": shadow.DesiredNRDEnabled = v.(bool)
				case "desired_nrd_start": shadow.DesiredNRDStart = v.(string)
				case "desired_nrd_end": shadow.DesiredNRDEnd = v.(string)
				case "desired_dnd_enabled": shadow.DesiredDNDEnabled = v.(bool)
				case "desired_dnd_start": shadow.DesiredDNDStart = v.(string)
				case "desired_dnd_end": shadow.DesiredDNDEnd = v.(string)
				case "desired_power_save": shadow.DesiredPowerSave = v.(bool)
				case "desired_version": shadow.DesiredVersion = v.(string)
				case "desired_volume": shadow.DesiredVolume = v.(*int)
				case "desired_brightness": shadow.DesiredBrightness = v.(*int)
				}
			}
			db.Create(&shadow)
		} else {
			db.Model(&models.DeviceShadow{}).Where("device_id = ?", deviceID).Updates(updates)
		}

		// 如果设备在线，立即通过 MQTT 下发期望状态
		if shadow.IsOnline && mqttClient != nil {
			syncDesiredStateToDeviceNow(db, deviceID, mqttClient)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id": deviceID,
				"synced":    shadow.IsOnline,
			},
		})
	}
}

// GetDesiredState 获取设备影子期望状态
// GET /api/v1/devices/:device_id/desired-state
func GetDesiredState(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceID := ctx.Param("device_id")

		var shadow models.DeviceShadow
		result := db.Where("device_id = ?", deviceID).First(&shadow)
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"message": "success",
				"data": gin.H{
					"device_id":           deviceID,
					"desired_nrd_enabled":  false,
					"desired_dnd_enabled": false,
					"desired_power_save":  false,
				},
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":           shadow.DeviceID,
				"desired_nrd_enabled":  shadow.DesiredNRDEnabled,
				"desired_nrd_start":    shadow.DesiredNRDStart,
				"desired_nrd_end":      shadow.DesiredNRDEnd,
				"desired_dnd_enabled":  shadow.DesiredDNDEnabled,
				"desired_dnd_start":   shadow.DesiredDNDStart,
				"desired_dnd_end":     shadow.DesiredDNDEnd,
				"desired_volume":      shadow.DesiredVolume,
				"desired_brightness":  shadow.DesiredBrightness,
				"desired_power_save":  shadow.DesiredPowerSave,
				"desired_version":     shadow.DesiredVersion,
			},
		})
	}
}

// GetReportedState 获取设备影子报告状态
// GET /api/v1/devices/:device_id/reported-state
func GetReportedState(db *gorm.DB, redisClient *utils.RedisClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceID := ctx.Param("device_id")

		// 从 Redis 获取设备影子
		shadow, err := redisClient.GetDeviceShadow(deviceID)
		if err != nil || shadow == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 0,
				"message": "success",
				"data": gin.H{
					"device_id":      deviceID,
					"is_online":      false,
					"battery_level":  0,
					"current_mode":   "",
					"last_heartbeat": nil,
					"last_ip":        "",
					"is_jailbroken": false,
					"root_status":   "",
					"latitude":       0,
					"longitude":      0,
				},
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":      shadow.DeviceID,
				"is_online":      shadow.IsOnline,
				"battery_level":  shadow.BatteryLevel,
				"current_mode":   shadow.CurrentMode,
				"last_heartbeat": shadow.LastHeartbeat,
				"last_ip":        shadow.LastIP,
				"is_jailbroken":  shadow.IsJailbroken,
				"root_status":    shadow.RootStatus,
				"latitude":       shadow.Latitude,
				"longitude":      shadow.Longitude,
			},
		})
	}
}

// GetStateDiff 获取设备影子状态差异
// GET /api/v1/devices/:device_id/state-diff
func GetStateDiff(db *gorm.DB, redisClient *utils.RedisClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceID := ctx.Param("device_id")

		// 获取 desired state from DB
		var dbShadow models.DeviceShadow
		dbResult := db.Where("device_id = ?", deviceID).First(&dbShadow)

		// 获取 reported state from Redis
		redisShadow, _ := redisClient.GetDeviceShadow(deviceID)

		// 构建 diff
		diffs := []map[string]interface{}{}

		// NRD 差异
		if dbShadow.DesiredNRDEnabled && redisShadow != nil {
			// NRD 已启用，设备应处于勿扰模式
			if !redisShadow.IsOnline {
				diffs = append(diffs, map[string]interface{}{
					"field":    "nrd",
					"expected": "enabled (NRD mode)",
					"actual":   "offline",
					"status":   "mismatch",
				})
			}
		}

		// 如果设备不在线，所有配置都认为不同步
		if redisShadow == nil || !redisShadow.IsOnline {
			if dbResult.Error == nil {
				if dbShadow.DesiredNRDEnabled {
					diffs = append(diffs, map[string]interface{}{
						"field":    "nrd_enabled",
						"expected": true,
						"actual":   "offline",
						"status":   "offline",
					})
				}
				if dbShadow.DesiredDNDEnabled {
					diffs = append(diffs, map[string]interface{}{
						"field":    "dnd_enabled",
						"expected": true,
						"actual":   "offline",
						"status":   "offline",
					})
				}
				if dbShadow.DesiredVolume != nil {
					diffs = append(diffs, map[string]interface{}{
						"field":    "volume",
						"expected": *dbShadow.DesiredVolume,
						"actual":   "offline",
						"status":   "offline",
					})
				}
				if dbShadow.DesiredBrightness != nil {
					diffs = append(diffs, map[string]interface{}{
						"field":    "brightness",
						"expected": *dbShadow.DesiredBrightness,
						"actual":   "offline",
						"status":   "offline",
					})
				}
				if dbShadow.DesiredPowerSave {
					diffs = append(diffs, map[string]interface{}{
						"field":    "power_save",
						"expected": true,
						"actual":   "offline",
						"status":   "offline",
					})
				}
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":    deviceID,
				"is_online":    redisShadow != nil && redisShadow.IsOnline,
				"has_diff":     len(diffs) > 0,
				"diff_count":   len(diffs),
				"diffs":        diffs,
			},
		})
	}
}

// syncDesiredStateToDeviceNow 立即同步期望状态到在线设备
func syncDesiredStateToDeviceNow(db *gorm.DB, deviceID string, mqttClient interface{}) {
	var shadow models.DeviceShadow
	if err := db.Where("device_id = ?", deviceID).First(&shadow).Error; err != nil {
		return
	}

	desiredCmd := map[string]interface{}{
		"cmd_id":    fmt.Sprintf("desired-sync-%s-%d", deviceID, time.Now().Unix()),
		"cmd_type":  "desired_sync",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	config := map[string]interface{}{}
	if shadow.DesiredNRDEnabled {
		config["nrd_enabled"] = true
		config["nrd_start"] = shadow.DesiredNRDStart
		config["nrd_end"] = shadow.DesiredNRDEnd
	}
	if shadow.DesiredDNDEnabled {
		config["dnd_enabled"] = true
		config["dnd_start"] = shadow.DesiredDNDStart
		config["dnd_end"] = shadow.DesiredDNDEnd
	}
	if shadow.DesiredVolume != nil {
		config["volume"] = *shadow.DesiredVolume
	}
	if shadow.DesiredBrightness != nil {
		config["brightness"] = *shadow.DesiredBrightness
	}
	if shadow.DesiredPowerSave {
		config["power_save"] = true
	}
	if shadow.DesiredVersion != "" {
		config["desired_version"] = shadow.DesiredVersion
	}

	desiredCmd["config"] = config

	// 通过 MQTT 下发
	topic := fmt.Sprintf("/device/%s/down/cmd", deviceID)
	payload, _ := json.Marshal(desiredCmd)

	// 使用类型断言调用 MQTT Publish
	if client, ok := mqttClient.(interface {
		Publish(topic string, qos byte, retained bool, payload interface{}) interface {
			Wait() bool
			Error() error
		}
	}); ok {
		token := client.Publish(topic, 0, false, payload)
		if token.Wait() && token.Error() == nil {
			fmt.Printf("[MQTT] 设备 %s 期望状态已实时同步\n", deviceID)
		}
	}
}

func isValidTimeFormat(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}
