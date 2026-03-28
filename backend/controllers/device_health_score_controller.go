package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceHealthScoreController 设备健康评分控制器
type DeviceHealthScoreController struct {
	DB *gorm.DB
}

// NewDeviceHealthScoreController 创建控制器
func NewDeviceHealthScoreController(db *gorm.DB) *DeviceHealthScoreController {
	return &DeviceHealthScoreController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *DeviceHealthScoreController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/devices/:device_id/health-score", ctrl.GetHealthScore)
	rg.POST("/devices/:device_id/health-score/calculate", ctrl.CalculateHealthScore)
	rg.GET("/devices/:device_id/health-history", ctrl.GetHealthHistory)
	rg.GET("/devices/health-report", ctrl.GetOverallHealthReport)
}

// PerformanceRecord 性能记录
type PerformanceRecord struct {
	CPUUsage      int   `json:"cpu_usage"`
	MemoryUsage   int   `json:"memory_usage"`
	NetworkLatency int  `json:"network_latency"`
}

// HealthScore 健康评分结构
type HealthScore struct {
	DeviceID       string        `json:"device_id"`
	TotalScore    float64       `json:"total_score"`
	Grade         string        `json:"grade"`
	UptimeScore   float64       `json:"uptime_score"`
	PerfScore     float64       `json:"perf_score"`
	SecurityScore float64       `json:"security_score"`
	BehaviorScore float64       `json:"behavior_score"`
	Issues        []HealthIssue `json:"issues"`
	CalculatedAt  time.Time     `json:"calculated_at"`
}

// HealthIssue 健康问题
type HealthIssue struct {
	Category    string  `json:"category"`
	Level       string  `json:"level"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ScoreImpact float64 `json:"score_impact"`
}

// GetHealthScore 获取设备健康评分
func (ctrl *DeviceHealthScoreController) GetHealthScore(c *gin.Context) {
	deviceID := c.Param("device_id")

	var record models.DeviceHealthScore
	if err := ctrl.DB.Where("device_id = ?", deviceID).
		Order("calculated_at DESC").First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "暂无健康评分，请先计算"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var issues []HealthIssue
	if record.IssuesJSON != "" {
		json.Unmarshal([]byte(record.IssuesJSON), &issues)
	}

	score := HealthScore{
		DeviceID:       record.DeviceID,
		TotalScore:    record.TotalScore,
		Grade:         record.Grade,
		UptimeScore:   record.UptimeScore,
		PerfScore:     record.PerfScore,
		SecurityScore: record.SecurityScore,
		BehaviorScore: record.BehaviorScore,
		Issues:        issues,
		CalculatedAt:  record.CalculatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": score})
}

// CalculateHealthScore 计算设备健康评分
func (ctrl *DeviceHealthScoreController) CalculateHealthScore(c *gin.Context) {
	deviceID := c.Param("device_id")

	// 获取设备信息
	var device models.Device
	if err := ctrl.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 获取设备影子
	var shadow struct {
		BatteryLevel  int       `gorm:"column:battery_level"`
		IsOnline      bool      `gorm:"column:is_online"`
		LastHeartbeat time.Time `gorm:"column:last_heartbeat"`
		IsJailbroken bool      `gorm:"column:is_jailbroken"`
		RootStatus   string    `gorm:"column:root_status"`
	}
	ctrl.DB.Table("device_shadows").Where("device_id = ?", deviceID).First(&shadow)

	// 获取性能历史 - 使用raw query
	var perfRecords []PerformanceRecord
	rows, err := ctrl.DB.Table("device_performance_history").
		Select("cpu_usage, memory_usage, network_latency").
		Where("device_id = ?", deviceID).
		Order("recorded_at DESC").Limit(100).Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r PerformanceRecord
			rows.Scan(&r.CPUUsage, &r.MemoryUsage, &r.NetworkLatency)
			perfRecords = append(perfRecords, r)
		}
	}

	// 获取告警统计
	var alertCount int64
	ctrl.DB.Model(&models.DeviceAlert{}).
		Where("device_id = ? AND created_at > ?", deviceID, time.Now().AddDate(0,0,-7)).
		Count(&alertCount)

	var unresolvedCount int64
	ctrl.DB.Model(&models.DeviceAlert{}).
		Where("device_id = ? AND resolved = ?", deviceID, false).
		Count(&unresolvedCount)

	// 计算各项得分
	issues := []HealthIssue{}
	var uptimeScore, perfScore, securityScore, behaviorScore float64 = 100, 100, 100, 100

	// 1. 在线时长得分
	if shadow.LastHeartbeat.IsZero() {
		uptimeScore = 0
		issues = append(issues, HealthIssue{Category: "uptime", Level: "critical", Title: "设备离线", Description: "设备已离线", ScoreImpact: -30})
	} else {
		minutesSince := time.Since(shadow.LastHeartbeat).Minutes()
		if minutesSince > 60 {
			uptimeScore = 100 - (minutesSince/60)*5
			if uptimeScore < 0 {
				uptimeScore = 0
			}
			if minutesSince > 1440 {
				uptimeScore = 0
				issues = append(issues, HealthIssue{Category: "uptime", Level: "critical", Title: "设备超过24小时无响应", Description: "设备可能已关机或网络异常", ScoreImpact: -40})
			} else if minutesSince > 60 {
				issues = append(issues, HealthIssue{Category: "uptime", Level: "warning", Title: "设备心跳延迟", Description: "设备超过1小时未上报心跳", ScoreImpact: -10})
			}
		}
	}

	// 2. 性能得分
	if len(perfRecords) > 0 {
		var avgCPU, avgMemory, avgNetwork float64
		for _, p := range perfRecords {
			avgCPU += float64(p.CPUUsage)
			avgMemory += float64(p.MemoryUsage)
			if p.NetworkLatency > 0 {
				avgNetwork += float64(p.NetworkLatency)
			}
		}
		avgCPU /= float64(len(perfRecords))
		avgMemory /= float64(len(perfRecords))
		avgNetwork /= float64(len(perfRecords))

		if avgCPU > 80 {
			perfScore -= 20
			issues = append(issues, HealthIssue{Category: "performance", Level: "warning", Title: "CPU使用率过高", Description: "平均CPU使用率超过80%", ScoreImpact: -20})
		} else if avgCPU > 60 {
			perfScore -= 10
		}

		if avgMemory > 85 {
			perfScore -= 15
			issues = append(issues, HealthIssue{Category: "performance", Level: "warning", Title: "内存使用率过高", Description: "平均内存使用率超过85%", ScoreImpact: -15})
		}

		if avgNetwork > 500 {
			perfScore -= 10
			issues = append(issues, HealthIssue{Category: "performance", Level: "info", Title: "网络延迟较高", Description: "平均网络延迟超过500ms", ScoreImpact: -10})
		}
	}

	if perfScore < 0 {
		perfScore = 0
	}

	// 3. 安全得分
	if shadow.IsJailbroken {
		securityScore = 0
		issues = append(issues, HealthIssue{Category: "security", Level: "critical", Title: "设备已越狱", Description: "设备安全状态严重风险", ScoreImpact: -100})
	} else if shadow.RootStatus == "rooted" {
		securityScore = 20
		issues = append(issues, HealthIssue{Category: "security", Level: "critical", Title: "设备已Root", Description: "设备已获取Root权限", ScoreImpact: -80})
	}

	// 4. 行为得分
	if alertCount > 50 {
		behaviorScore -= 30
		issues = append(issues, HealthIssue{Category: "behavior", Level: "warning", Title: "告警过于频繁", Description: "7天内产生超过50条告警", ScoreImpact: -30})
	} else if alertCount > 20 {
		behaviorScore -= 15
		issues = append(issues, HealthIssue{Category: "behavior", Level: "info", Title: "告警较多", Description: "7天内产生较多告警", ScoreImpact: -15})
	}

	if unresolvedCount > 5 {
		behaviorScore -= 10
		issues = append(issues, HealthIssue{Category: "behavior", Level: "warning", Title: "未处理告警较多", Description: "有较多未处理的告警", ScoreImpact: -10})
	}

	if behaviorScore < 0 {
		behaviorScore = 0
	}

	// 计算总分
	totalScore := (uptimeScore*0.25 + perfScore*0.25 + securityScore*0.3 + behaviorScore*0.2)

	// 确定等级
	grade := "E"
	if totalScore >= 90 {
		grade = "A"
	} else if totalScore >= 75 {
		grade = "B"
	} else if totalScore >= 60 {
		grade = "C"
	} else if totalScore >= 40 {
		grade = "D"
	}

	// 保存评分记录
	issuesJSON, _ := json.Marshal(issues)
	record := models.DeviceHealthScore{
		DeviceID:       deviceID,
		TotalScore:    totalScore,
		Grade:         grade,
		UptimeScore:   uptimeScore,
		PerfScore:     perfScore,
		SecurityScore: securityScore,
		BehaviorScore: behaviorScore,
		IssuesJSON:    string(issuesJSON),
		CalculatedAt:  time.Now(),
	}

	existing := models.DeviceHealthScore{}
	if err := ctrl.DB.Where("device_id = ?", deviceID).First(&existing).Error; err == nil {
		record.ID = existing.ID
		ctrl.DB.Save(&record)
	} else {
		ctrl.DB.Create(&record)
	}

	score := HealthScore{
		DeviceID:       deviceID,
		TotalScore:    totalScore,
		Grade:         grade,
		UptimeScore:   uptimeScore,
		PerfScore:     perfScore,
		SecurityScore: securityScore,
		BehaviorScore: behaviorScore,
		Issues:        issues,
		CalculatedAt:  record.CalculatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "健康评分计算完成",
		"data":    score,
	})
}

// GetHealthHistory 获取健康评分历史
func (ctrl *DeviceHealthScoreController) GetHealthHistory(c *gin.Context) {
	deviceID := c.Param("device_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	var list []models.DeviceHealthScore
	ctrl.DB.Model(&models.DeviceHealthScore{}).Where("device_id = ?", deviceID).Count(&total)

	ctrl.DB.Where("device_id = ?", deviceID).
		Order("calculated_at DESC").
		Offset((page-1)*pageSize).Limit(pageSize).
		Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetOverallHealthReport 获取整体健康报告
func (ctrl *DeviceHealthScoreController) GetOverallHealthReport(c *gin.Context) {
	var gradeA, gradeB, gradeC, gradeD, gradeE int64
	ctrl.DB.Model(&models.DeviceHealthScore{}).Where("grade = ?", "A").Count(&gradeA)
	ctrl.DB.Model(&models.DeviceHealthScore{}).Where("grade = ?", "B").Count(&gradeB)
	ctrl.DB.Model(&models.DeviceHealthScore{}).Where("grade = ?", "C").Count(&gradeC)
	ctrl.DB.Model(&models.DeviceHealthScore{}).Where("grade = ?", "D").Count(&gradeD)
	ctrl.DB.Model(&models.DeviceHealthScore{}).Where("grade = ?", "E").Count(&gradeE)

	var avgScore float64
	ctrl.DB.Model(&models.DeviceHealthScore{}).Select("AVG(total_score)").Row().Scan(&avgScore)

	var totalDevices int64
	ctrl.DB.Model(&models.Device{}).Count(&totalDevices)

	var onlineDevices int64
	ctrl.DB.Model(&models.DeviceShadow{}).Where("is_online = ?", true).Count(&onlineDevices)

	var scoredDevices int64
	ctrl.DB.Model(&models.DeviceHealthScore{}).Count(&scoredDevices)

	gradeAToFloat := float64(gradeA)
	scoredDevicesFloat := float64(scoredDevices)
	healthRate := 0.0
	if scoredDevicesFloat > 0 {
		healthRate = (gradeAToFloat + float64(gradeB)) / scoredDevicesFloat * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total_devices":     totalDevices,
			"online_devices":    onlineDevices,
			"scored_devices":    scoredDevices,
			"average_score":     avgScore,
			"grade_distribution": gin.H{
				"A": gradeA,
				"B": gradeB,
				"C": gradeC,
				"D": gradeD,
				"E": gradeE,
			},
			"health_rate": healthRate,
		},
	})
}
