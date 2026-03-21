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

		// 修复：从 package_quotas 表查询指定类型的配额记录
		quota, err := models.GetQuota(db, tenantID, string(quotaType))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "DB_ERROR",
				"message": "配额查询失败",
			})
			c.Abort()
			return
		}

		// 租户配额记录不存在，获取套餐默认配额
		limit := getQuotaLimit(db, tenantID, quotaType)
		usedCount := 0

		// 如果配额记录存在，获取已使用量
		if quota != nil && quota.ID != 0 {
			usedCount = quota.QuotaUsed
		}

		if limit == -1 {
			// -1 表示不限
			c.Next()
			return
		}

		if usedCount >= limit {
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

// getQuotaLimit 从 packages 表查询套餐配额上限
// 修复：查询 packages 表（不是 plans 表），从 quota_config JSONB 提取配额
func getQuotaLimit(db *gorm.DB, tenantID string, quotaType QuotaType) int {
	// Directly query package_quotas table for this tenant's quota limit
	var quota models.TenantQuota
	if err := db.Where("tenant_id = ? AND quota_type = ?", tenantID, string(quotaType)).First(&quota).Error; err != nil {
		return getDefaultQuotaLimit("free", quotaType)
	}
	return quota.QuotaLimit
}


// getQuotaKey 将 QuotaType 转换为 quota_config JSONB 的 key
func getQuotaKey(quotaType QuotaType) string {
	keyMap := map[QuotaType]string{
		QuotaUser:   "users",
		QuotaDevice: "devices",
		QuotaDept:   "departments",
		QuotaStore:  "stores",
	}
	if key, ok := keyMap[quotaType]; ok {
		return key
	}
	// 回退：使用 quotaType 本身作为 key
	return string(quotaType)
}

// getDefaultQuotaLimit 默认套餐配额
func getDefaultQuotaLimit(packageCode string, quotaType QuotaType) int {
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

	if d, ok := defaults[packageCode]; ok {
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
