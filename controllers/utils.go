package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"mdm-backend/middleware"
)

// getUserIDFromContext 从 Context 获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	if uid, exists := c.Get("user_id"); exists {
		switch v := uid.(type) {
		case uint:
			return v
		case int:
			return uint(v)
		case int64:
			return uint(v)
		case float64:
			return uint(v)
		case string:
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				return uint(id)
			}
		}
	}
	return 0
}

// getUserID 获取用户ID
func getUserID(c *gin.Context) uint {
	return getUserIDFromContext(c)
}

// getTenantID 获取租户ID
func getTenantID(c *gin.Context) string {
	return middleware.GetTenantID(c)
}

// parseInt 字符串转 int
func parseInt(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return 0
}
