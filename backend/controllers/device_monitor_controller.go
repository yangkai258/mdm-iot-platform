package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceMonitorController 设备监控控制器
type DeviceMonitorController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterDeviceMonitorRoutes 注册设备监控路由
func RegisterDeviceMonitorRoutes(api *gin.RouterGroup, db *gorm.DB, redisClient *utils.RedisClient) {
	ctrl := &DeviceMonitorController{DB: db, Redis: redisClient}
	// 单设备监控子路由
	api.GET("/devices/:device_id/monitor/realtime", ctrl.GetRealtimeData)
	api.POST("/devices/:device_id/monitor/alert-rules", ctrl.CreateAlertRule)
	// 全局监控路由
	api.GET("/devices/monitor/dashboard", ctrl.GetDashboard)
	api.GET("/devices/monitor/metrics", ctrl.GetMetrics)
}

// ============ 实时监控数据 ============

// GetRealtimeData 获取设备实时监控数据
func (c *DeviceMonitorController) GetRealtimeData(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 1. 获取设备基本信息
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 2. 从 Redis 获取设备影子（实时数据）
	shadow, err := c.Redis.GetDeviceShadow(deviceID)
	isOnline := false
	var batteryLevel int
	var currentMode string
	var lastHeartbeat *time.Time
	var latitude, longitude float64
	var lastIP string
	var isJailbroken bool
	var rootStatus string

	if err == nil && shadow != nil {
		isOnline = shadow.IsOnline
		batteryLevel = shadow.BatteryLevel
		currentMode = shadow.CurrentMode
		lastHeartbeat = shadow.LastHeartbeat
		latitude = shadow.Latitude
		longitude = shadow.Longitude
		lastIP = shadow.LastIP
		isJailbroken = shadow.IsJailbroken
		rootStatus = shadow.RootStatus
	} else {
		// 影子不存在，说明设备离线
		isOnline = false
	}

	// 3. 获取最近的心跳记录（最近24小时）
	var recentAlerts []models.DeviceAlert
	c.DB.Where("device_id = ? AND created_at > ?", deviceID, time.Now().Add(-24*time.Hour)).
		Order("created_at DESC").
		Limit(10).
		Find(&recentAlerts)

	// 4. 计算在线时长（从最后一次心跳推算）
	var onlineDuration string
	if isOnline && lastHeartbeat != nil {
		d := time.Since(*lastHeartbeat)
		onlineDuration = formatDuration(d)
	}

	// 5. 计算设备运行时长（从创建至今，按服役状态计）
	runtime := time.Since(device.CreatedAt).Round(time.Second)
	runtimeFormatted := formatDuration(runtime)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id":    deviceID,
			"is_online":    isOnline,
			"battery_level": batteryLevel,
			"current_mode": currentMode,
			"last_heartbeat": lastHeartbeat,
			"last_ip":       lastIP,
			"is_jailbroken": isJailbroken,
			"root_status":  rootStatus,
			"location": gin.H{
				"latitude":  latitude,
				"longitude": longitude,
			},
			"online_duration": onlineDuration,
			"runtime":        runtimeFormatted,
			"lifecycle_status": device.LifecycleStatus,
			"hardware_model":  device.HardwareModel,
			"firmware_version": device.FirmwareVersion,
			"recent_alerts":   recentAlerts,
		},
	})
}

// ============ 监控仪表板数据 ============

// GetDashboard 获取监控仪表板汇总数据
func (c *DeviceMonitorController) GetDashboard(ctx *gin.Context) {
	// 1. 设备总数和在线/离线统计
	var totalDevices int64
	c.DB.Model(&models.Device{}).Count(&totalDevices)

	var onlineDevices int64
	// 在线设备：从 Redis 统计
	keys, _ := c.Redis.GetAllShadowKeys()
	onlineDevices = int64(len(keys))

	offlineDevices := totalDevices - onlineDevices
	if offlineDevices < 0 {
		offlineDevices = 0
	}

	// 2. 各生命周期状态设备数量
	type StatusCount struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}
	var statusCounts []StatusCount
	c.DB.Model(&models.Device{}).
		Select("lifecycle_status as status, count(*) as count").
		Group("lifecycle_status").
		Scan(&statusCounts)

	// 3. 各型号设备分布
	type ModelCount struct {
		Model string `json:"model"`
		Count int64  `json:"count"`
	}
	var modelCounts []ModelCount
	c.DB.Model(&models.Device{}).
		Select("hardware_model as model, count(*) as count").
		Group("hardware_model").
		Order("count DESC").
		Limit(10).
		Scan(&modelCounts)

	// 4. 今日告警统计
	var todayAlerts int64
	c.DB.Model(&models.DeviceAlert{}).
		Where("created_at > ?", time.Now().Truncate(24*time.Hour)).
		Count(&todayAlerts)

	var todayResolvedAlerts int64
	c.DB.Model(&models.DeviceAlert{}).
		Where("status = 3 AND resolved_at > ?", time.Now().Truncate(24*time.Hour)).
		Count(&todayResolvedAlerts)

	// 5. 活跃告警（未处理）
	var activeAlerts int64
	c.DB.Model(&models.DeviceAlert{}).
		Where("status IN (1, 2)").
		Count(&activeAlerts)

	// 6. 近7天告警趋势
	var alertTrend []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	c.DB.Model(&models.DeviceAlert{}).
		Select("DATE(created_at) as date, count(*) as count").
		Where("created_at > ?", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&alertTrend)

	// 7. 在线设备列表（取前20个）
	type DeviceOnlineInfo struct {
		DeviceID   string     `json:"device_id"`
		IsOnline   bool       `json:"is_online"`
		LastSeen   *time.Time `json:"last_seen"`
		Battery    int        `json:"battery_level"`
		CurrentMode string    `json:"current_mode"`
	}
	onlineDevicesList := make([]DeviceOnlineInfo, 0)
	for _, key := range keys {
		if len(onlineDevicesList) >= 20 {
			break
		}
		deviceID := key[7:] // 去掉 "shadow:" 前缀
		shadow, err := c.Redis.GetDeviceShadow(deviceID)
		if err == nil && shadow != nil {
			onlineDevicesList = append(onlineDevicesList, DeviceOnlineInfo{
				DeviceID:    deviceID,
				IsOnline:    shadow.IsOnline,
				LastSeen:    shadow.LastHeartbeat,
				Battery:     shadow.BatteryLevel,
				CurrentMode: shadow.CurrentMode,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"summary": gin.H{
				"total_devices":     totalDevices,
				"online_devices":     onlineDevices,
				"offline_devices":   offlineDevices,
				"online_rate":        safeDivide(float64(onlineDevices), float64(totalDevices)),
				"active_alerts":     activeAlerts,
				"today_alerts":      todayAlerts,
				"today_resolved":    todayResolvedAlerts,
			},
			"status_distribution": statusCounts,
			"model_distribution":  modelCounts,
			"alert_trend":          alertTrend,
			"online_devices":        onlineDevicesList,
		},
	})
}

// ============ 设备指标统计 ============

// GetMetrics 获取设备指标统计
func (c *DeviceMonitorController) GetMetrics(ctx *gin.Context) {
	// 支持按时间范围过滤，默认为24小时
	hoursStr := ctx.DefaultQuery("hours", "24")
	hours, _ := strconv.Atoi(hoursStr)
	if hours <= 0 || hours > 720 { // 最多30天
		hours = 24
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	// 1. 设备总数和分布
	var totalDevices int64
	c.DB.Model(&models.Device{}).Count(&totalDevices)

	var activeDevices int64 // 有过心跳的设备
	c.DB.Model(&models.DeviceShadow{}).
		Where("last_heartbeat > ?", since).
		Count(&activeDevices)

	// 2. 告警统计
	var totalAlerts int64
	c.DB.Model(&models.DeviceAlert{}).
		Where("created_at > ?", since).
		Count(&totalAlerts)

	var highSeverityAlerts int64
	c.DB.Model(&models.DeviceAlert{}).
		Where("created_at > ? AND severity >= 3", since).
		Count(&highSeverityAlerts)

	// 3. 平均在线率（从今日活跃影子 / 设备总数）
	keys, _ := c.Redis.GetAllShadowKeys()
	onlineRate := safeDivide(float64(len(keys)), float64(totalDevices))

	// 4. 固件版本分布
	type VersionCount struct {
		Version string `json:"version"`
		Count   int64  `json:"count"`
	}
	var versionCounts []VersionCount
	c.DB.Model(&models.Device{}).
		Select("firmware_version as version, count(*) as count").
		Group("firmware_version").
		Order("count DESC").
		Limit(10).
		Scan(&versionCounts)

	// 5. 设备在线时长分布
	type UptimeBucket struct {
		Bucket  string `json:"bucket"`
		Count   int64  `json:"count"`
	}
	uptimeBuckets := []UptimeBucket{
		{Bucket: "<1小时", Count: 0},
		{Bucket: "1-6小时", Count: 0},
		{Bucket: "6-24小时", Count: 0},
		{Bucket: "1-7天", Count: 0},
		{Bucket: ">7天", Count: 0},
	}

	now := time.Now()
	for _, key := range keys {
		deviceID := key[7:]
		shadow, err := c.Redis.GetDeviceShadow(deviceID)
		if err != nil || shadow.LastHeartbeat == nil {
			continue
		}
		duration := now.Sub(*shadow.LastHeartbeat)
		if duration < time.Hour {
			uptimeBuckets[0].Count++
		} else if duration < 6*time.Hour {
			uptimeBuckets[1].Count++
		} else if duration < 24*time.Hour {
			uptimeBuckets[2].Count++
		} else if duration < 7*24*time.Hour {
			uptimeBuckets[3].Count++
		} else {
			uptimeBuckets[4].Count++
		}
	}

	// 6. 告警类型分布
	type AlertTypeCount struct {
		AlertType string `json:"alert_type"`
		Count     int64  `json:"count"`
	}
	var alertTypeCounts []AlertTypeCount
	c.DB.Model(&models.DeviceAlert{}).
		Select("alert_type, count(*) as count").
		Where("created_at > ?", since).
		Group("alert_type").
		Order("count DESC").
		Scan(&alertTypeCounts)

	// 7. 设备行为模式分布（从影子获取）
	type ModeCount struct {
		Mode  string `json:"mode"`
		Count int64  `json:"count"`
	}
	var modeCounts []ModeCount
	modeMap := make(map[string]int64)
	for _, key := range keys {
		deviceID := key[7:]
		shadow, err := c.Redis.GetDeviceShadow(deviceID)
		if err != nil || shadow == nil {
			continue
		}
		if shadow.CurrentMode != "" {
			modeMap[shadow.CurrentMode]++
		}
	}
	for mode, count := range modeMap {
		modeCounts = append(modeCounts, ModeCount{Mode: mode, Count: count})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period_hours":        hours,
			"total_devices":       totalDevices,
			"active_devices":      activeDevices,
			"total_alerts":        totalAlerts,
			"high_severity_alerts": highSeverityAlerts,
			"online_rate":         onlineRate,
			"firmware_distribution": versionCounts,
			"uptime_distribution":  uptimeBuckets,
			"alert_type_distribution": alertTypeCounts,
			"mode_distribution":      modeCounts,
		},
	})
}

// ============ 设备级别告警规则 ============

// CreateAlertRule 为指定设备创建告警规则
func (c *DeviceMonitorController) CreateAlertRule(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 验证设备存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	var req struct {
		Name       string  `json:"name" binding:"required"`
		AlertType  string  `json:"alert_type" binding:"required"`  // battery_low, offline, temperature_high
		Condition  string  `json:"condition" binding:"required"`    // <, >, =, >=, <=
		Threshold  float64 `json:"threshold" binding:"required"`
		Severity   int     `json:"severity"`                        // 1-4
		Enabled    *bool   `json:"enabled"`
		NotifyWays string  `json:"notify_ways"`                      // email,sms,webhook
		Remark     string  `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	severity := 1
	if req.Severity > 0 {
		severity = req.Severity
	}

	rule := models.DeviceAlertRule{
		Name:       req.Name,
		DeviceID:   deviceID,
		AlertType:  req.AlertType,
		Condition:  req.Condition,
		Threshold:  req.Threshold,
		Severity:   severity,
		Enabled:    enabled,
		NotifyWays: req.NotifyWays,
		Remark:     req.Remark,
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "创建告警规则失败: " + err.Error(),
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "告警规则创建成功",
		"data":    rule,
	})
}

// ============ 辅助函数 ============

// safeDivide 安全除法，避免除零
func safeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return (a / b) * 100 // 返回百分比
}

// formatDuration 格式化时长
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "<1分钟"
	}
	if d < time.Hour {
		return strconv.Itoa(int(d.Minutes())) + "分钟"
	}
	if d < 24*time.Hour {
		return strconv.Itoa(int(d.Hours())) + "小时"
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return strconv.Itoa(days) + "天" + strconv.Itoa(hours) + "小时"
}
