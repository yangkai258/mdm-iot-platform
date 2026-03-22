package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenant_id"

// TenantContext 从 JWT claims 中解析 tenant_id 并注入到 Gin Context
func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			// 尝试从 map 格式 claims（如通用 jwt 解析）
			if claimsMap, ok := c.Get("jwt_claims"); ok {
				if m, ok := claimsMap.(map[string]interface{}); ok {
					if tenantID, ok := m["tenant_id"]; ok {
						c.Set(TenantIDKey, formatTenantID(tenantID))
					}
				}
			}
			c.Next()
			return
		}

		// 尝试 map 格式
		if m, ok := claims.(map[string]interface{}); ok {
			if tenantID, ok := m["tenant_id"]; ok {
				c.Set(TenantIDKey, formatTenantID(tenantID))
			}
		}

		c.Next()
	}
}

// formatTenantID 将各种类型转换为 string
func formatTenantID(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case float64:
		return fmt.Sprintf("%.0f", s)
	case int:
		return fmt.Sprintf("%d", s)
	default:
		return ""
	}
}

// GetTenantID 获取当前请求的租户ID
func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get(TenantIDKey); exists {
		return formatTenantID(tenantID)
	}
	return ""
}

// IsSuperAdmin 判断是否为超级管理员
func IsSuperAdmin(c *gin.Context) bool {
	if claims, exists := c.Get("claims"); exists {
		if m, ok := claims.(map[string]interface{}); ok {
			if isSuper, ok := m["is_super_admin"].(bool); ok {
				return isSuper
			}
		}
	}
	return false
}
