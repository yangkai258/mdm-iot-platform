package middleware

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QuotaType 定义配额类型
type QuotaType string

const (
	QuotaUser   QuotaType = "user"
	QuotaDevice QuotaType = "device"
	QuotaDept   QuotaType = "dept"
	QuotaStore  QuotaType = "store"
)

// QuotaCheck 创建配额校验中间件
// usage: api.PUT("/users", QuotaCheck(db, QuotaUser), CreateUser)
func QuotaCheck(db *gorm.DB, quotaType QuotaType) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID == "" {
			// 租户ID为空，说明是非租户请求（如超管内部接口），跳过校验
			c.Next()
			return
		}

		quota, err := models.GetQuota(db, tenantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "DB_ERROR",
				"message": "配额查询失败",
			})
			c.Abort()
			return
		}

		// 租户配额记录不存在，初始化默认值
		if quota.ID == 0 {
			quota = &models.TenantQuota{
				TenantID:   tenantID,
				UserCount:  5,
				DeviceCount: 10,
				DeptCount:  1,
				StoreCount: 1,
			}
		}

		usedCount := getUsedCount(quota, quotaType)
		limit := getQuotaLimit(db, tenantID, quotaType)

		if limit == -1 {
			// -1 表示不限
			c.Next()
			return
		}

		if int(usedCount) >= limit {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "QUOTA_EXCEEDED",
				"message": getQuotaExceededMessage(quotaType, usedCount, limit),
				"data": gin.H{
					"quota_type":  quotaType,
					"used":        usedCount,
					"limit":       limit,
					"upgrade_url": "/admin/tenants/" + tenantID + "/change-plan",
				},
			})
			c.Abort()
			return
		}

		// 将配额信息存入 context，供后续 handler 使用
		c.Set("quota_info", map[string]interface{}{
			"tenant_id":  tenantID,
			"quota_type": quotaType,
			"used_count": usedCount,
			"limit":      limit,
		})

		c.Next()
	}
}

func getUsedCount(quota *models.TenantQuota, quotaType QuotaType) int {
	switch quotaType {
	case QuotaUser:
		return quota.UserCount
	case QuotaDevice:
		return quota.DeviceCount
	case QuotaDept:
		return quota.DeptCount
	case QuotaStore:
		return quota.StoreCount
	default:
		return 0
	}
}

// getQuotaLimit 从 plans 表查询套餐配额上限
func getQuotaLimit(db *gorm.DB, tenantID string, quotaType QuotaType) int {
	var tenant models.Tenant
	if err := db.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		return 5 // 默认免费配额
	}

	planCode := tenant.Plan
	if planCode == "" {
		planCode = "free"
	}

	var plan models.Plan
	if err := db.Where("plan_code = ?", planCode).First(&plan).Error; err != nil {
		return getDefaultQuotaLimit(planCode, quotaType)
	}

	switch quotaType {
	case QuotaUser:
		return plan.UserQuota
	case QuotaDevice:
		return plan.DeviceQuota
	case QuotaDept:
		return plan.DeptQuota
	case QuotaStore:
		return plan.StoreQuota
	default:
		return -1
	}
}

// getDefaultQuotaLimit 默认套餐配额
func getDefaultQuotaLimit(planCode string, quotaType QuotaType) int {
	defaults := map[string]map[QuotaType]int{
		"free": {
			QuotaUser:   5,
			QuotaDevice: 10,
			QuotaDept:   1,
			QuotaStore:  1,
		},
		"basic": {
			QuotaUser:   50,
			QuotaDevice: 100,
			QuotaDept:   5,
			QuotaStore:  10,
		},
		"professional": {
			QuotaUser:   200,
			QuotaDevice: 500,
			QuotaDept:   20,
			QuotaStore:  50,
		},
		"enterprise": {
			QuotaUser:   -1,
			QuotaDevice: -1,
			QuotaDept:   -1,
			QuotaStore:  -1,
		},
	}

	if d, ok := defaults[planCode]; ok {
		if v, ok := d[quotaType]; ok {
			return v
		}
	}
	return 5
}

func getQuotaExceededMessage(quotaType QuotaType, used, limit int) string {
	typeNames := map[QuotaType]string{
		QuotaUser:   "用户数",
		QuotaDevice: "设备数",
		QuotaDept:   "部门数",
		QuotaStore:  "门店数",
	}
	typeName := typeNames[quotaType]
	if typeName == "" {
		typeName = "配额"
	}
	return typeName + "已达上限（" + itoa(used) + "/" + itoa(limit) + "），请联系管理员升级套餐"
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	negative := false
	if i < 0 {
		negative = true
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
