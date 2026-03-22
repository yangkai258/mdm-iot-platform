package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UsageController 用量计费控制器
type UsageController struct {
	DB *gorm.DB
}

// ===== 用量查询 =====

// GetCurrentUsage 获取当前用量
func (c *UsageController) GetCurrentUsage(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	// 获取当前周期
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	// 查询当前周期用量
	var usageRecords []models.UsageRecord
	if err := c.DB.Where("user_id = ? AND period_start = ? AND period_end = ?", userID, periodStart, periodEnd).Find(&usageRecords).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取用户订阅信息
	var subscription models.UserSubscription
	c.DB.Where("user_id = ? AND status = ?", userID, "active").Order("created_at DESC").First(&subscription)

	// 获取计划配额
	var plan models.SubscriptionPlan
	quotas := make(map[string]interface{})
	if subscription.PlanID != "" {
		c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan)
		if plan.Quotas != nil {
			quotas = plan.Quotas
		}
	}

	// 格式化返回数据
	usageList := make([]map[string]interface{}, 0)
	usageTypes := []string{models.UsageTypeAPICall, models.UsageTypeDevice, models.UsageTypeStorage, models.UsageTypeBandwidth}

	for _, usageType := range usageTypes {
		found := false
		for _, record := range usageRecords {
			if record.UsageType == usageType {
				usageList = append(usageList, map[string]interface{}{
					"type":       record.UsageType,
					"used":       record.UsageValue,
					"unit":       record.Unit,
					"quota":      quotas[usageType],
					"percentage": 0,
				})
				found = true
				break
			}
		}
		if !found {
			usageList = append(usageList, map[string]interface{}{
				"type":       usageType,
				"used":       0,
				"unit":       getUnitForType(usageType),
				"quota":      quotas[usageType],
				"percentage": 0,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period": gin.H{
				"start": periodStart,
				"end":   periodEnd,
			},
			"usage":   usageList,
			"plan":    plan,
		},
	})
}

// GetUsageHistory 获取用量历史
func (c *UsageController) GetUsageHistory(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	// 用量类型过滤
	usageType := ctx.Query("usage_type")

	var records []models.UsageRecord
	query := c.DB.Where("user_id = ?", userID)

	if usageType != "" {
		query = query.Where("usage_type = ?", usageType)
	}

	// 日期范围
	if startDate := ctx.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("record_date >= ?", t)
		}
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("record_date <= ?", t.Add(24*time.Hour-time.Second))
		}
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("record_date DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": records,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// GetQuotas 获取配额信息
func (c *UsageController) GetQuotas(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	// 获取用户订阅
	var subscription models.UserSubscription
	if err := c.DB.Where("user_id = ? AND status = ?", userID, "active").Order("created_at DESC").First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"quotas": []map[string]interface{}{},
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取计划配额
	var plan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "计划不存在", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取当前周期
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 查询已使用配额
	var usageRecords []models.UsageRecord
	c.DB.Where("user_id = ? AND period_start = ?", userID, periodStart).Find(&usageRecords)

	// 构建配额列表
	quotasList := make([]map[string]interface{}, 0)
	if plan.Quotas != nil {
		for key, limit := range plan.Quotas {
			used := float64(0)
			unit := getUnitForType(key)
			for _, record := range usageRecords {
				if record.UsageType == key {
					used = record.UsageValue
					break
				}
			}
			limitFloat, _ := parseFloat64(limit)
			percentage := float64(0)
			if limitFloat > 0 {
				percentage = (used / limitFloat) * 100
			}
			quotasList = append(quotasList, map[string]interface{}{
				"type":       key,
				"limit":      limitFloat,
				"used":       used,
				"unit":       unit,
				"percentage": percentage,
				"status":     getQuotaStatus(percentage),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"plan":       plan,
			"quotas":     quotasList,
			"period":     periodStart.Format("2006-01"),
		},
	})
}

// GetQuotaByType 获取特定配额
func (c *UsageController) GetQuotaByType(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	quotaType := ctx.Param("type")
	validTypes := map[string]bool{
		models.UsageTypeAPICall:   true,
		models.UsageTypeDevice:    true,
		models.UsageTypeStorage:   true,
		models.UsageTypeBandwidth: true,
	}

	if !validTypes[quotaType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的配额类型", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 获取用户订阅
	var subscription models.UserSubscription
	if err := c.DB.Where("user_id = ? AND status = ?", userID, "active").Order("created_at DESC").First(&subscription).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "无活跃订阅", "error_code": "ERR_NOT_FOUND"})
		return
	}

	// 获取计划配额
	var plan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "计划不存在", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取当前周期已使用量
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var usageRecord models.UsageRecord
	c.DB.Where("user_id = ? AND usage_type = ? AND period_start = ?", userID, quotaType, periodStart).First(&usageRecord)

	limit, _ := parseFloat64(plan.Quotas[quotaType])
	used := usageRecord.UsageValue
	percentage := float64(0)
	if limit > 0 {
		percentage = (used / limit) * 100
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"type":       quotaType,
			"limit":      limit,
			"used":       used,
			"unit":       usageRecord.Unit,
			"percentage": percentage,
			"status":     getQuotaStatus(percentage),
		},
	})
}

// UpdateQuota 更新配额（管理员操作）
func (c *UsageController) UpdateQuota(ctx *gin.Context) {
	// 检查管理员权限
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4031, "message": "权限不足", "error_code": "ERR_FORBIDDEN"})
		return
	}

	quotaType := ctx.Param("type")
	validTypes := map[string]bool{
		models.UsageTypeAPICall:   true,
		models.UsageTypeDevice:    true,
		models.UsageTypeStorage:   true,
		models.UsageTypeBandwidth: true,
	}

	if !validTypes[quotaType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的配额类型", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	var req struct {
		UserID uint    `json:"user_id" binding:"required"`
		Limit  float64 `json:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 获取或创建用户配额记录
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	var quota models.UserQuota
	err := c.DB.Where("user_id = ? AND quota_type = ? AND period_start = ?", req.UserID, quotaType, periodStart).First(&quota).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新配额记录
		quota = models.UserQuota{
			QuotaID:    "quota-" + uuid.New().String(),
			UserID:     req.UserID,
			QuotaType:  quotaType,
			QuotaLimit: req.Limit,
			QuotaUsed:  0,
			Unit:       getUnitForType(quotaType),
			PeriodType: "monthly",
			PeriodStart: periodStart,
			PeriodEnd:  periodEnd,
		}
		c.DB.Create(&quota)
	} else {
		// 更新现有配额
		quota.QuotaLimit = req.Limit
		c.DB.Save(&quota)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    quota,
	})
}

// GetUsageStats 获取用量统计
func (c *UsageController) GetUsageStats(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	// 获取最近6个月的统计
	now := time.Now()
	stats := make([]map[string]interface{}, 0)

	for i := 5; i >= 0; i-- {
		month := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, now.Location())
		monthEnd := month.AddDate(0, 1, 0).Add(-time.Second)

		var totalAPI, totalStorage, totalBandwidth float64
		c.DB.Model(&models.UsageRecord{}).
			Select("usage_type, SUM(usage_value) as total").
			Where("user_id = ? AND period_start = ? AND period_end = ?", userID, month, monthEnd).
			Group("usage_type").
			Scan(&[]map[string]interface{}{})

		// 重新查询
		var records []models.UsageRecord
		c.DB.Where("user_id = ? AND period_start = ?", userID, month).Find(&records)

		for _, record := range records {
			switch record.UsageType {
			case models.UsageTypeAPICall:
				totalAPI += record.UsageValue
			case models.UsageTypeStorage:
				totalStorage += record.UsageValue
			case models.UsageTypeBandwidth:
				totalBandwidth += record.UsageValue
			}
		}

		stats = append(stats, map[string]interface{}{
			"month":           month.Format("2006-01"),
			"api_calls":       totalAPI,
			"storage_gb":      totalStorage,
			"bandwidth_gb":    totalBandwidth,
		})
	}

	// 获取总计
	var totalRecords []models.UsageRecord
	c.DB.Where("user_id = ?", userID).Find(&totalRecords)

	totalAPI := float64(0)
	totalStorage := float64(0)
	totalBandwidth := float64(0)

	for _, record := range totalRecords {
		switch record.UsageType {
		case models.UsageTypeAPICall:
			totalAPI += record.UsageValue
		case models.UsageTypeStorage:
			totalStorage += record.UsageValue
		case models.UsageTypeBandwidth:
			totalBandwidth += record.UsageValue
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"monthly": stats,
			"total": gin.H{
				"api_calls":    totalAPI,
				"storage_gb":   totalStorage,
				"bandwidth_gb": totalBandwidth,
			},
		},
	})
}

// CheckQuota 检查配额是否足够（内部使用）
func (c *UsageController) CheckQuota(userID uint, quotaType string, required float64) (bool, float64, float64) {
	// 获取用户订阅
	var subscription models.UserSubscription
	if err := c.DB.Where("user_id = ? AND status = ?", userID, "active").Order("created_at DESC").First(&subscription).Error; err != nil {
		return false, 0, 0
	}

	// 获取计划配额
	var plan models.SubscriptionPlan
	if err := c.DB.Where("plan_id = ?", subscription.PlanID).First(&plan).Error; err != nil {
		return false, 0, 0
	}

	limit, _ := parseFloat64(plan.Quotas[quotaType])
	if limit == 0 {
		return true, 0, 0 // 无限制
	}

	// 获取当前已使用量
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var usageRecord models.UsageRecord
	if err := c.DB.Where("user_id = ? AND usage_type = ? AND period_start = ?", userID, quotaType, periodStart).First(&usageRecord).Error; err != nil {
		used := float64(0)
		return (used + required) <= limit, used, limit
	}

	return (usageRecord.UsageValue + required) <= limit, usageRecord.UsageValue, limit
}

// RecordUsage 记录用量（内部使用）
func (c *UsageController) RecordUsage(userID uint, quotaType string, value float64) error {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	var record models.UsageRecord
	err := c.DB.Where("user_id = ? AND usage_type = ? AND period_start = ?", userID, quotaType, periodStart).First(&record).Error

	if err == gorm.ErrRecordNotFound {
		record = models.UsageRecord{
			RecordID:    "usage-" + uuid.New().String(),
			UserID:      userID,
			UsageType:   quotaType,
			UsageValue:  value,
			Unit:        getUnitForType(quotaType),
			RecordDate:  now,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}
		return c.DB.Create(&record).Error
	}

	record.UsageValue += value
	return c.DB.Save(&record).Error
}

// ===== 辅助函数 =====

func getUnitForType(usageType string) string {
	switch usageType {
	case models.UsageTypeAPICall:
		return "次"
	case models.UsageTypeDevice:
		return "个"
	case models.UsageTypeStorage:
		return "GB"
	case models.UsageTypeBandwidth:
		return "GB"
	default:
		return "次"
	}
}

func getQuotaStatus(percentage float64) string {
	if percentage >= 100 {
		return "exceeded"
	} else if percentage >= 80 {
		return "warning"
	}
	return "normal"
}

func parseFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case json.Number:
		f, _ := val.Float64()
		return f, true
	default:
		return 0, false
	}
}

func (c *UsageController) isAdmin(ctx *gin.Context) bool {
	role := ctx.GetString("role")
	return role == "admin" || role == "super_admin"
}
