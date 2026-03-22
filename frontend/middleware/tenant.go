package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenant_id"

// TenantContext 租户上下文中间件
// 优先级：1. X-Tenant-ID Header  2. JWT Claims tenant_id  3. jwt_claims tenant_id
func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先级 1: 从 X-Tenant-ID Header 读取（适用于 Service-to-Service 调用或特殊场景）
		if tenantID := c.GetHeader("X-Tenant-ID"); tenantID != "" {
			c.Set(TenantIDKey, strings.TrimSpace(tenantID))
			c.Next()
			return
		}

		// 优先级 2: 从 JWTAuth 中间件注入的 claims 中读取
		if claims, exists := c.Get("claims"); exists {
			if m, ok := claims.(map[string]interface{}); ok {
				if tenantID, ok := m["tenant_id"]; ok {
					c.Set(TenantIDKey, formatTenantID(tenantID))
					c.Next()
					return
				}
			}
		}

		// 优先级 3: 从 jwt_claims map 中读取（兜底兼容）
		if claimsMap, ok := c.Get("jwt_claims"); ok {
			if m, ok := claimsMap.(map[string]interface{}); ok {
				if tenantID, ok := m["tenant_id"]; ok {
					c.Set(TenantIDKey, formatTenantID(tenantID))
					c.Next()
					return
				}
			}
		}

		c.Next()
	}
}

// formatTenantID 将各种类型转换为 string
func formatTenantID(v interface{}) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return strings.TrimSpace(s)
	case float64:
		return fmt.Sprintf("%.0f", s)
	case int:
		return fmt.Sprintf("%d", s)
	case int64:
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
