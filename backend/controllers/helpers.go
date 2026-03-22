package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// getUserIDFromContext 从 gin.Context 获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		if uid, ok := id.(uint); ok {
			return uid
		}
	}
	if idStr := c.Query("user_id"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			return uint(id)
		}
	}
	if idStr := c.GetHeader("X-User-ID"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			return uint(id)
		}
	}
	return 1 // 默认管理员ID
}

// getTenantIDFromContext 从 gin.Context 获取租户ID
func getTenantIDFromContext(c *gin.Context) string {
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}
	if tid := c.Query("tenant_id"); tid != "" {
		return tid
	}
	return "default"
}
