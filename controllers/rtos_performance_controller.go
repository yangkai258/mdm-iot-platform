package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RTOSPerformanceController RTOS 性能监控控制器
type RTOSPerformanceController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterRTOSPerformanceRoutes 注册 RTOS 性能相关路由
func (c *RTOSPerformanceController) RegisterRTOSPerformanceRoutes(r *gin.RouterGroup) {
	// RTOS 性能监控
	r.GET("/device/:id/rtos/stats", c.GetRTOSStats)
	r.GET("/device/:id/rtos/tasks", c.GetRTOSTasks)
	r.GET("/device/:id/rtos/memory", c.GetRTOSMemory)

	// 设备性能分析
	r.GET("/device/:id/performance/history", c.GetPerformanceHistory)
	r.GET("/device/:id/performance/report", c.GetPerformanceReport)

	// 固件优化配置
	r.GET("/device/:id/optimization/config", c.GetOptimizationConfig)
	r.PUT("/device/:id/optimization/config", c.UpdateOptimizationConfig)
	r.POST("/device/:id/optimization/apply", c.ApplyOptimization)

	// 设备端上报 RTOS 数据(MQTT HTTP 回调)
	r.POST("/device/:id/rtos/report", c.DeviceRTOSReport)
}

// GetRTOSStats 获取 RTOS 性能统计
// GET /api/v1/device/:id/rtos/stats
func (c *RTOSPerformanceController) GetRTOSStats(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	// 获取最新的 RTOS 统计
	var stats models.RTOSStats
	err := c.DB.Where("device_id = ?", deviceID).Order("recorded_at DESC").First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "暂无RTOS统计数据",
				"data": gin.H{
					"stats":     nil,
					"device_id": deviceID,
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询RTOS统计失败: %v", err)})
		return
	}

	// 计算运行时间描述
	uptimeDesc := formatDuration(int64(stats.Uptime))

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"stats":       stats,
			"uptime_desc": uptimeDesc,
			"device_id":   deviceID,
			"firmware":    device.FirmwareVersion,
			"hardware":    device.HardwareModel,
		},
	})
}

// GetRTOSTasks 获取 RTOS 任务列表
// GET /api/v1/device/:id/rtos/tasks
func (c *RTOSPerformanceController) GetRTOSTasks(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	// 获取最新的任务快照
	var latestRecordTime time.Time
	c.DB.Model(&models.RTOSTask{}).Where("device_id = ?", deviceID).
		Select("MAX(recorded_at)").Scan(&latestRecordTime)

	if latestRecordTime.IsZero() {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "暂无任务数据",
			"data": gin.H{
				"tasks":     []models.RTOSTask{},
				"device_id": deviceID,
				"record_at": nil,
			},
		})
		return
	}

	// 获取该时间点的所有任务
	var tasks []models.RTOSTask
	err := c.DB.Where("device_id = ? AND recorded_at = ?", deviceID, latestRecordTime).
		Order("priority ASC, cpu_usage DESC").Find(&tasks).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询任务列表失败: %v", err)})
		return
	}

	// 计算栈使用率
	taskInfos := make([]gin.H, 0, len(tasks))
	for _, task := range tasks {
		stackPercent := float64(0)
		if task.StackSize > 0 {
			stackPercent = float64(task.StackUsed) / float64(task.StackSize) * 100
		}
		taskInfos = append(taskInfos, gin.H{
			"task_name":     task.TaskName,
			"priority":      task.Priority,
			"stack_size":    task.StackSize,
			"stack_used":    task.StackUsed,
			"stack_percent": fmt.Sprintf("%.1f%%", stackPercent),
			"state":         task.State,
			"cpu_usage":     task.CPUUsage,
			"runtime_ms":    task.Runtime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"tasks":       taskInfos,
			"device_id":   deviceID,
			"record_at":   latestRecordTime,
			"total_count": len(tasks),
		},
	})
}

// GetRTOSMemory 获取 RTOS 内存使用
// GET /api/v1/device/:id/rtos/memory
func (c *RTOSPerformanceController) GetRTOSMemory(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	// 获取最新的内存数据
	var memory models.RTOSMemory
	err := c.DB.Where("device_id = ?", deviceID).Order("recorded_at DESC").First(&memory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "暂无内存数据",
				"data": gin.H{
					"memory":    nil,
					"device_id": deviceID,
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询内存数据失败: %v", err)})
		return
	}

	// 计算内存使用百分比
	heapUsagePercent := float64(0)
	if memory.TotalHeap > 0 {
		heapUsagePercent = float64(memory.UsedHeap) / float64(memory.TotalHeap) * 100
	}

	// 计算 PSRAM 和 Flash 使用率
	psramPercent := float64(0)
	if memory.PSRAMTotal > 0 {
		psramPercent = float64(memory.PSRAMUsed) / float64(memory.PSRAMTotal) * 100
	}
	flashPercent := float64(0)
	if memory.FlashTotal > 0 {
		flashPercent = float64(memory.FlashUsed) / float64(memory.FlashTotal) * 100
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"memory": gin.H{
				"heap": gin.H{
					"total_bytes":    memory.TotalHeap,
					"used_bytes":     memory.UsedHeap,
					"free_bytes":     memory.FreeHeap,
					"min_free_bytes": memory.MinFreeHeap,
					"usage_percent":  fmt.Sprintf("%.2f%%", heapUsagePercent),
					"frag_percent":   memory.MemoryFragPercent,
					"total_human":    formatRTOSBytes(memory.TotalHeap),
					"used_human":     formatRTOSBytes(memory.UsedHeap),
					"free_human":     formatRTOSBytes(memory.FreeHeap),
				},
				"psram": gin.H{
					"total_bytes":   memory.PSRAMTotal,
					"used_bytes":    memory.PSRAMUsed,
					"usage_percent": fmt.Sprintf("%.2f%%", psramPercent),
					"total_human":   formatRTOSBytes(memory.PSRAMTotal),
					"used_human":    formatRTOSBytes(memory.PSRAMUsed),
				},
				"flash": gin.H{
					"total_bytes":   memory.FlashTotal,
					"used_bytes":    memory.FlashUsed,
					"usage_percent": fmt.Sprintf("%.2f%%", flashPercent),
				},
			},
			"device_id": deviceID,
			"record_at": memory.RecordedAt,
		},
	})
}

// GetPerformanceHistory 获取性能历史
// GET /api/v1/device/:id/performance/history
func (c *RTOSPerformanceController) GetPerformanceHistory(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 解析查询参数
	metricType := ctx.DefaultQuery("metric_type", "cpu")
	startTimeStr := ctx.DefaultQuery("start_time", "")
	endTimeStr := ctx.DefaultQuery("end_time", "")
	limitStr := ctx.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	// 构建查询
	query := c.DB.Model(&models.DevicePerformanceHistory{}).Where("device_id = ?", deviceID)

	if metricType != "" && metricType != "all" {
		query = query.Where("metric_type = ?", metricType)
	}

	// 时间范围过滤
	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startTimeStr)
	} else {
		startTime = time.Now().AddDate(0, 0, -7)
	}
	if endTimeStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endTimeStr)
	} else {
		endTime = time.Now()
	}
	query = query.Where("recorded_at BETWEEN ? AND ?", startTime, endTime)

	// 查询数据
	var history []models.DevicePerformanceHistory
	err = query.Order("recorded_at DESC").Limit(limit).Find(&history).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询性能历史失败: %v", err)})
		return
	}

	// 按指标类型分组
	grouped := make(map[string][]gin.H)
	for _, h := range history {
		grouped[h.MetricType] = append(grouped[h.MetricType], gin.H{
			"metric_name":  h.MetricName,
			"metric_value": h.MetricValue,
			"unit":         h.Unit,
			"recorded_at":  h.RecordedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"history":     history,
			"grouped":     grouped,
			"device_id":   deviceID,
			"metric_type": metricType,
			"start_time":  startTime,
			"end_time":    endTime,
			"total":       len(history),
		},
	})
}

// GetPerformanceReport 获取性能报告
// GET /api/v1/device/:id/performance/report
func (c *RTOSPerformanceController) GetPerformanceReport(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	// 获取最新 RTOS 统计数据
	var latestStats models.RTOSStats
	c.DB.Where("device_id = ?", deviceID).Order("recorded_at DESC").First(&latestStats)

	// 获取最新 RTOS 内存数据
	var latestMemory models.RTOSMemory
	c.DB.Where("device_id = ?", deviceID).Order("recorded_at DESC").First(&latestMemory)

	// 获取最新任务数据
	var latestTasks []models.RTOSTask
	var latestTaskTime time.Time
	c.DB.Model(&models.RTOSTask{}).Where("device_id = ?", deviceID).
		Select("MAX(recorded_at)").Scan(&latestTaskTime)
	if !latestTaskTime.IsZero() {
		c.DB.Where("device_id = ? AND recorded_at = ?", deviceID, latestTaskTime).
			Order("cpu_usage DESC").Limit(10).Find(&latestTasks)
	}

	// 计算最近24小时的统计数据
	oneDayAgo := time.Now().AddDate(0, 0, -1)

	// CPU 统计
	var cpuStats struct {
		AvgCPU float64
		MaxCPU float64
		MinCPU float64
	}
	c.DB.Model(&models.DevicePerformanceHistory{}).
		Where("device_id = ? AND metric_type = 'cpu' AND recorded_at >= ?", deviceID, oneDayAgo).
		Select("COALESCE(AVG(metric_value), 0) as avg_cpu, COALESCE(MAX(metric_value), 0) as max_cpu, COALESCE(MIN(metric_value), 0) as min_cpu").
		Scan(&cpuStats)

	// 内存统计
	var memStats struct {
		AvgUsage float64
		MaxUsage float64
	}
	if latestMemory.TotalHeap > 0 {
		memStats.MaxUsage = float64(latestMemory.UsedHeap) / float64(latestMemory.TotalHeap) * 100
	}
	if latestMemory.MinFreeHeap > 0 && latestMemory.TotalHeap > 0 {
		memStats.AvgUsage = float64(latestMemory.TotalHeap-latestMemory.MinFreeHeap) / float64(latestMemory.TotalHeap) * 100
	}

	// 性能健康评分
	healthScore := calculateHealthScore(latestStats, latestMemory)

	// 告警和建议
	alerts := generatePerformanceAlerts(latestStats, latestMemory, cpuStats, memStats)

	// Top 5 高CPU任务
	topTasks := make([]gin.H, 0, 5)
	for i, t := range latestTasks {
		if i >= 5 {
			break
		}
		topTasks = append(topTasks, gin.H{
			"task_name": t.TaskName,
			"cpu_usage": t.CPUUsage,
		})
	}

	// 获取设备影子状态
	var shadow models.DeviceShadow
	c.DB.Where("device_id = ?", deviceID).First(&shadow)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"report": gin.H{
				"generated_at":     time.Now().Format(time.RFC3339),
				"device_id":        deviceID,
				"firmware_version": device.FirmwareVersion,
				"hardware_model":   device.HardwareModel,
				"online_status":    shadow.IsOnline,
				"last_heartbeat":   shadow.LastHeartbeat,
				"health_score":     healthScore,
				"summary": gin.H{
					"uptime_seconds":   latestStats.Uptime,
					"uptime_desc":      formatDuration(latestStats.Uptime),
					"task_count":       latestStats.TaskCount,
					"context_switches": latestStats.ContextSwitchCount,
				},
				"cpu": gin.H{
					"current_usage":   latestStats.CPUUsage,
					"idle_task_usage": latestStats.IdleTaskUsage,
					"avg_24h":         fmt.Sprintf("%.2f%%", cpuStats.AvgCPU),
					"max_24h":         fmt.Sprintf("%.2f%%", cpuStats.MaxCPU),
					"min_24h":         fmt.Sprintf("%.2f%%", cpuStats.MinCPU),
				},
				"memory": gin.H{
					"heap_usage_percent": fmt.Sprintf("%.2f%%", memStats.MaxUsage),
					"heap_total":         formatRTOSBytes(latestMemory.TotalHeap),
					"heap_used":          formatRTOSBytes(latestMemory.UsedHeap),
					"heap_free":          formatRTOSBytes(latestMemory.FreeHeap),
					"min_free_heap":      formatRTOSBytes(latestMemory.MinFreeHeap),
					"frag_percent":       latestMemory.MemoryFragPercent,
				},
				"top_cpu_tasks":   topTasks,
				"alerts":          alerts,
				"recommendations": generateRecommendations(latestStats, latestMemory, alerts),
			},
		},
	})
}

// GetOptimizationConfig 获取优化配置
// GET /api/v1/device/:id/optimization/config
func (c *RTOSPerformanceController) GetOptimizationConfig(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	// 获取配置
	var config models.FirmwareOptimizationConfig
	err := c.DB.Where("device_id = ?", deviceID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认配置
			defaultConfig := models.FirmwareOptimizationConfig{
				DeviceID:         deviceID,
				OptimizationType: "balance",
				IsEnabled:        true,
				IsApplied:        false,
			}
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "使用默认配置",
				"data": gin.H{
					"config":     defaultConfig,
					"device_id":  deviceID,
					"is_default": true,
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询优化配置失败: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"config":     config,
			"device_id":  deviceID,
			"is_default": false,
		},
	})
}

// UpdateOptimizationConfig 更新优化配置
// PUT /api/v1/device/:id/optimization/config
func (c *RTOSPerformanceController) UpdateOptimizationConfig(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	var req struct {
		OptimizationType string `json:"optimization_type" binding:"required"`
		ConfigJSON       string `json:"config_json"`
		IsEnabled        *bool  `json:"is_enabled"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": fmt.Sprintf("请求参数错误: %v", err)})
		return
	}

	// 验证优化类型
	validTypes := map[string]bool{
		"power_save":          true,
		"performance":         true,
		"balance":             true,
		"network_low_latency": true,
	}
	if !validTypes[req.OptimizationType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的优化类型，有效值: power_save, performance, balance, network_low_latency"})
		return
	}

	// 获取当前用户(从 context)
	createdBy := "system"
	if userID, exists := ctx.Get("user_id"); exists {
		createdBy = fmt.Sprintf("%v", userID)
	}

	// 查找或创建配置
	var config models.FirmwareOptimizationConfig
	err := c.DB.Where("device_id = ?", deviceID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新配置
			config = models.FirmwareOptimizationConfig{
				DeviceID:         deviceID,
				OptimizationType: req.OptimizationType,
				ConfigJSON:       req.ConfigJSON,
				IsEnabled:        true,
				CreatedBy:        createdBy,
			}
			if req.IsEnabled != nil {
				config.IsEnabled = *req.IsEnabled
			}
			if err := c.DB.Create(&config).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("创建优化配置失败: %v", err)})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询优化配置失败: %v", err)})
			return
		}
	} else {
		// 更新现有配置
		updates := map[string]interface{}{
			"optimization_type": req.OptimizationType,
			"config_json":       req.ConfigJSON,
			"is_applied":        false, // 配置已更新，需要重新应用
		}
		if req.IsEnabled != nil {
			updates["is_enabled"] = *req.IsEnabled
		}
		if err := c.DB.Model(&config).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("更新优化配置失败: %v", err)})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "优化配置已更新，请调用 apply 接口使配置生效",
		"data": gin.H{
			"device_id": deviceID,
			"applied":   false,
		},
	})
}

// ApplyOptimization 应用优化配置
// POST /api/v1/device/:id/optimization/apply
func (c *RTOSPerformanceController) ApplyOptimization(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备失败: %v", err)})
		return
	}

	// 获取设备在线状态
	var shadow models.DeviceShadow
	if err := c.DB.Where("device_id = ?", deviceID).First(&shadow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询设备状态失败: %v", err)})
		return
	}

	if !shadow.IsOnline {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "设备不在线，无法应用优化配置",
			"data": gin.H{
				"device_id": deviceID,
				"is_online": false,
			},
		})
		return
	}

	// 获取优化配置
	var config models.FirmwareOptimizationConfig
	err := c.DB.Where("device_id = ?", deviceID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "未找到优化配置，请先调用 PUT 接口设置配置"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("查询优化配置失败: %v", err)})
		return
	}

	if !config.IsEnabled {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "优化配置未启用"})
		return
	}

	// 构建 MQTT 下发命令
	cmdPayload := map[string]interface{}{
		"cmd_type":  "optimization",
		"action":    "apply_config",
		"device_id": deviceID,
		"config": gin.H{
			"optimization_type": config.OptimizationType,
			"config_json":       config.ConfigJSON,
		},
		"timestamp": time.Now().Unix(),
	}

	cmdJSON, _ := json.Marshal(cmdPayload)

	// 通过 MQTT 下发
	mqttClient := mqtt.GetGlobalMQTTClient()
	if mqttClient == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"code": 503, "message": "MQTT服务不可用"})
		return
	}

	topic := fmt.Sprintf("/device/%s/down/cmd", deviceID)
	token := mqttClient.Publish(topic, 1, false, cmdJSON)
	if token.Wait() && token.Error() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": fmt.Sprintf("下发优化命令失败: %v", token.Error())})
		return
	}

	// 记录命令历史
	cmdHistory := models.CommandHistory{
		DeviceID: deviceID,
		CmdID:    fmt.Sprintf("opt_%d", time.Now().UnixNano()),
		CmdType:  "optimization",
		Action:   "apply_config",
		Status:   "sent",
		SentAt:   time.Now(),
	}
	c.DB.Create(&cmdHistory)

	// 更新配置状态
	now := time.Now()
	c.DB.Model(&config).Updates(map[string]interface{}{
		"is_applied":      true,
		"applied_at":      &now,
		"applied_version": device.FirmwareVersion,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "优化配置已下发，请等待设备确认",
		"data": gin.H{
			"device_id":         deviceID,
			"optimization_type": config.OptimizationType,
			"applied":           true,
			"mqtt_topic":        topic,
		},
	})
}

// DeviceRTOSReport 设备端上报 RTOS 数据
// POST /api/v1/device/:id/rtos/report
func (c *RTOSPerformanceController) DeviceRTOSReport(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req struct {
		Type string      `json:"type" binding:"required"` // stats, memory, task
		Data interface{} `json:"data" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": fmt.Sprintf("请求参数错误: %v", err)})
		return
	}

	now := time.Now()
	switch req.Type {
	case "stats":
		statsData, ok := req.Data.(map[string]interface{})
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "stats 数据格式错误"})
			return
		}
		stats := models.RTOSStats{
			DeviceID:   deviceID,
			RecordedAt: now,
		}
		if v, ok := statsData["cpu_usage"].(float64); ok {
			stats.CPUUsage = v
		}
		if v, ok := statsData["task_count"].(float64); ok {
			stats.TaskCount = int(v)
		}
		if v, ok := statsData["uptime"].(float64); ok {
			stats.Uptime = int64(v)
		}
		if v, ok := statsData["irq_count"].(float64); ok {
			stats.IRQCount = int64(v)
		}
		if v, ok := statsData["context_switch_count"].(float64); ok {
			stats.ContextSwitchCount = int64(v)
		}
		if v, ok := statsData["idle_task_usage"].(float64); ok {
			stats.IdleTaskUsage = v
		}
		c.DB.Create(&stats)

		// 同步到 PerformanceHistory
		history := models.DevicePerformanceHistory{
			DeviceID:    deviceID,
			MetricType:  "cpu",
			MetricName:  "rtos_cpu_usage",
			MetricValue: stats.CPUUsage,
			Unit:        "%",
			RecordedAt:  now,
		}
		c.DB.Create(&history)

	case "memory":
		memData, ok := req.Data.(map[string]interface{})
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "memory 数据格式错误"})
			return
		}
		memory := models.RTOSMemory{
			DeviceID:   deviceID,
			RecordedAt: now,
		}
		if v, ok := memData["total_heap"].(float64); ok {
			memory.TotalHeap = int64(v)
		}
		if v, ok := memData["used_heap"].(float64); ok {
			memory.UsedHeap = int64(v)
		}
		if v, ok := memData["free_heap"].(float64); ok {
			memory.FreeHeap = int64(v)
		}
		if v, ok := memData["min_free_heap"].(float64); ok {
			memory.MinFreeHeap = int64(v)
		}
		if v, ok := memData["memory_frag_percent"].(float64); ok {
			memory.MemoryFragPercent = v
		}
		if v, ok := memData["psram_total"].(float64); ok {
			memory.PSRAMTotal = int64(v)
		}
		if v, ok := memData["psram_used"].(float64); ok {
			memory.PSRAMUsed = int64(v)
		}
		if v, ok := memData["flash_total"].(float64); ok {
			memory.FlashTotal = int64(v)
		}
		if v, ok := memData["flash_used"].(float64); ok {
			memory.FlashUsed = int64(v)
		}
		c.DB.Create(&memory)

		if memory.TotalHeap > 0 {
			memUsage := float64(memory.UsedHeap) / float64(memory.TotalHeap) * 100
			history := models.DevicePerformanceHistory{
				DeviceID:    deviceID,
				MetricType:  "memory",
				MetricName:  "rtos_heap_usage",
				MetricValue: memUsage,
				Unit:        "%",
				RecordedAt:  now,
			}
			c.DB.Create(&history)
		}

	case "tasks":
		tasksData, ok := req.Data.([]interface{})
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "tasks 数据格式错误"})
			return
		}

		var maxVersion int
		c.DB.Model(&models.RTOSTask{}).Where("device_id = ?", deviceID).
			Select("COALESCE(MAX(record_version), 0)").Scan(&maxVersion)
		newVersion := maxVersion + 1

		for _, t := range tasksData {
			taskMap, ok := t.(map[string]interface{})
			if !ok {
				continue
			}
			task := models.RTOSTask{
				DeviceID:      deviceID,
				RecordedAt:    now,
				RecordVersion: newVersion,
			}
			if v, ok := taskMap["task_name"].(string); ok {
				task.TaskName = v
			}
			if v, ok := taskMap["priority"].(float64); ok {
				task.Priority = int(v)
			}
			if v, ok := taskMap["stack_size"].(float64); ok {
				task.StackSize = int(v)
			}
			if v, ok := taskMap["stack_used"].(float64); ok {
				task.StackUsed = int(v)
			}
			if v, ok := taskMap["state"].(string); ok {
				task.State = v
			}
			if v, ok := taskMap["cpu_usage"].(float64); ok {
				task.CPUUsage = v
			}
			if v, ok := taskMap["runtime"].(float64); ok {
				task.Runtime = int64(v)
			}
			c.DB.Create(&task)
		}

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的数据类型"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "数据上报成功",
		"data": gin.H{
			"device_id":   deviceID,
			"type":        req.Type,
			"recorded_at": now,
		},
	})
}

// ============ 辅助函数 ============

// formatDuration 将秒数转换为可读的时间描述
func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds%60)
	}
	hours := minutes / 60
	minutes = minutes % 60
	if hours < 24 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	}
	days := hours / 24
	hours = hours % 24
	return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
}

// formatRTOSBytes 将字节数转换为可读的大小描述
func formatRTOSBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
}

// calculateHealthScore 计算性能健康评分 (0-100)
func calculateHealthScore(stats models.RTOSStats, memory models.RTOSMemory) int {
	score := 100

	// CPU 评分
	if stats.CPUUsage > 90 {
		score -= 30
	} else if stats.CPUUsage > 80 {
		score -= 20
	} else if stats.CPUUsage > 70 {
		score -= 10
	}

	// 内存评分
	if memory.TotalHeap > 0 {
		memUsage := float64(memory.UsedHeap) / float64(memory.TotalHeap) * 100
		if memUsage > 90 {
			score -= 30
		} else if memUsage > 80 {
			score -= 20
		} else if memUsage > 70 {
			score -= 10
		}
	}

	// 内存碎片评分
	if memory.MemoryFragPercent > 30 {
		score -= 10
	}

	if score < 0 {
		score = 0
	}
	return score
}

// generatePerformanceAlerts 生成性能告警
func generatePerformanceAlerts(stats models.RTOSStats, memory models.RTOSMemory, cpuStats struct {
	AvgCPU float64
	MaxCPU float64
	MinCPU float64
}, memStats struct {
	AvgUsage float64
	MaxUsage float64
}) []gin.H {
	alerts := []gin.H{}

	// CPU 高使用率告警
	if stats.CPUUsage > 90 {
		alerts = append(alerts, gin.H{
			"level":   "critical",
			"type":    "cpu_high",
			"message": fmt.Sprintf("CPU 使用率过高: %.2f%%", stats.CPUUsage),
		})
	} else if stats.CPUUsage > 80 {
		alerts = append(alerts, gin.H{
			"level":   "warning",
			"type":    "cpu_high",
			"message": fmt.Sprintf("CPU 使用率偏高: %.2f%%", stats.CPUUsage),
		})
	}

	// 内存高使用率告警
	if memory.TotalHeap > 0 {
		memUsage := float64(memory.UsedHeap) / float64(memory.TotalHeap) * 100
		if memUsage > 90 {
			alerts = append(alerts, gin.H{
				"level":   "critical",
				"type":    "memory_high",
				"message": fmt.Sprintf("堆内存使用率过高: %.2f%%", memUsage),
			})
		} else if memUsage > 80 {
			alerts = append(alerts, gin.H{
				"level":   "warning",
				"type":    "memory_high",
				"message": fmt.Sprintf("堆内存使用率偏高: %.2f%%", memUsage),
			})
		}
	}

	// 内存碎片告警
	if memory.MemoryFragPercent > 30 {
		alerts = append(alerts, gin.H{
			"level":   "warning",
			"type":    "memory_frag",
			"message": fmt.Sprintf("内存碎片率较高: %.2f%%", memory.MemoryFragPercent),
		})
	}

	// 最小空闲内存告警
	if memory.MinFreeHeap > 0 && memory.TotalHeap > 0 {
		minFreePercent := float64(memory.MinFreeHeap) / float64(memory.TotalHeap) * 100
		if minFreePercent < 10 {
			alerts = append(alerts, gin.H{
				"level":   "critical",
				"type":    "memory_low",
				"message": fmt.Sprintf("历史最低空闲内存仅剩: %.2f%%", minFreePercent),
			})
		}
	}

	return alerts
}

// generateRecommendations 生成优化建议
func generateRecommendations(stats models.RTOSStats, memory models.RTOSMemory, alerts []gin.H) []string {
	recommendations := []string{}

	if len(alerts) == 0 {
		recommendations = append(recommendations, "系统运行良好，无需特殊优化")
		return recommendations
	}

	for _, alert := range alerts {
		alertType := alert["type"].(string)
		switch alertType {
		case "cpu_high":
			recommendations = append(recommendations, "建议：检查高 CPU 占用的任务，考虑优化算法或降低任务频率")
		case "memory_high", "memory_low":
			recommendations = append(recommendations, "建议：检查内存泄漏，适当增加堆大小或优化内存分配策略")
		case "memory_frag":
			recommendations = append(recommendations, "建议：内存碎片较多，建议重启设备或优化内存分配模式")
		}
	}

	return recommendations
}
