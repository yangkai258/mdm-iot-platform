package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecretFromEnv 从环境变量获取 JWT 密钥（延迟读取避免循环依赖）
var jwtSecretUserCtx []byte

func getJWTSecretForUserContext() []byte {
	if jwtSecretUserCtx == nil {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			panic("JWT_SECRET environment variable is not set")
		}
		jwtSecretUserCtx = []byte(secret)
	}
	return jwtSecretUserCtx
}

// UserClaims JWT 中与数据权限相关的用户信息
type UserClaims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	RoleID       uint   `json:"role_id"`
	TenantID     string `json:"tenant_id"`
	OrgID        uint   `json:"org_id"`         // 组织/部门ID
	IsSuperAdmin bool   `json:"is_super_admin"`
}

// ContextKey 用户上下文 key
const (
	ContextKeyUserID       = "user_id"
	ContextKeyUsername     = "username"
	ContextKeyRoleID       = "role_id"
	ContextKeyTenantID     = "tenant_id"
	ContextKeyOrgID        = "org_id"
	ContextKeyIsSuperAdmin = "is_super_admin"
	ContextKeyClaims       = "claims"
)

// UserContext 用户上下文中间件
// 从 JWT 中提取用户ID、OrgID、角色等信息存入 Gin Context
// 供 DataScope 中间件和其他需要用户上下文的逻辑使用
func UserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除不需要用户上下文的路径
		path := c.Request.URL.Path
		if isPublicPath(path) {
			c.Next()
			return
		}

		// 优先使用 JWTAuth 中间件已解析的 claims
		if claimsVal, exists := c.Get(ContextKeyClaims); exists {
			if claimsMap, ok := claimsVal.(map[string]interface{}); ok {
				setUserContextFromMap(c, claimsMap)
				c.Next()
				return
			}
		}

		// 如果 JWTAuth 尚未解析，从 Authorization header 自己解析
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		claims := &UserClaims{}
		tokenObj, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecretForUserContext(), nil
		})

		if err == nil && tokenObj.Valid {
			setUserContextFromClaims(c, claims)
		}

		c.Next()
	}
}

// setUserContextFromClaims 将解析出的 claims 写入 Gin Context
func setUserContextFromClaims(c *gin.Context, claims *UserClaims) {
	c.Set(ContextKeyUserID, claims.UserID)
	c.Set(ContextKeyUsername, claims.Username)
	c.Set(ContextKeyRoleID, claims.RoleID)
	c.Set(ContextKeyTenantID, claims.TenantID)
	c.Set(ContextKeyOrgID, claims.OrgID)
	c.Set(ContextKeyIsSuperAdmin, claims.IsSuperAdmin)
	c.Set(ContextKeyClaims, map[string]interface{}{
		"user_id":       claims.UserID,
		"username":      claims.Username,
		"role_id":       claims.RoleID,
		"tenant_id":     claims.TenantID,
		"org_id":        claims.OrgID,
		"is_super_admin": claims.IsSuperAdmin,
	})
}

// setUserContextFromMap 从 map 设置用户上下文
func setUserContextFromMap(c *gin.Context, m map[string]interface{}) {
	if v, ok := m["user_id"]; ok {
		c.Set(ContextKeyUserID, toUint(v))
	}
	if v, ok := m["username"]; ok {
		if s, ok := v.(string); ok {
			c.Set(ContextKeyUsername, s)
		}
	}
	if v, ok := m["role_id"]; ok {
		c.Set(ContextKeyRoleID, toUint(v))
	}
	if v, ok := m["tenant_id"]; ok {
		c.Set(ContextKeyTenantID, toString(v))
	}
	if v, ok := m["org_id"]; ok {
		c.Set(ContextKeyOrgID, toUint(v))
	}
	if v, ok := m["is_super_admin"]; ok {
		if b, ok := v.(bool); ok {
			c.Set(ContextKeyIsSuperAdmin, b)
		}
	}
}

// GetUserID 获取当前用户ID
func GetUserID(c *gin.Context) uint {
	if v, exists := c.Get(ContextKeyUserID); exists {
		return toUint(v)
	}
	return 0
}

// GetUsername 获取当前用户名
func GetUsername(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyUsername); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetRoleID 获取当前用户角色ID
func GetRoleID(c *gin.Context) uint {
	if v, exists := c.Get(ContextKeyRoleID); exists {
		return toUint(v)
	}
	return 0
}

// GetTenantID 获取当前租户ID
func GetTenantIDCtx(c *gin.Context) string {
	if v, exists := c.Get(ContextKeyTenantID); exists {
		return toString(v)
	}
	return ""
}

// GetOrgID 获取当前用户组织ID
func GetOrgID(c *gin.Context) uint {
	if v, exists := c.Get(ContextKeyOrgID); exists {
		return toUint(v)
	}
	return 0
}

// IsSuperAdmin 判断是否为超级管理员
func IsSuperAdminCtx(c *gin.Context) bool {
	if v, exists := c.Get(ContextKeyIsSuperAdmin); exists {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetUserClaims 获取完整用户上下文
func GetUserClaims(c *gin.Context) *UserClaims {
	return &UserClaims{
		UserID:       GetUserID(c),
		Username:     GetUsername(c),
		RoleID:       GetRoleID(c),
		TenantID:     GetTenantIDCtx(c),
		OrgID:        GetOrgID(c),
		IsSuperAdmin: IsSuperAdminCtx(c),
	}
}

// isPublicPath 判断路径是否公开（不需要认证）
func isPublicPath(path string) bool {
	publicPrefixes := []string{
		"/api/v1/auth",
		"/health",
		"/api/v1/public",
	}
	for _, prefix := range publicPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// toUint 将任意类型转换为 uint
func toUint(v interface{}) uint {
	switch val := v.(type) {
	case uint:
		return val
	case float64:
		return uint(val)
	case int:
		return uint(val)
	case int64:
		return uint(val)
	case string:
		if val == "" {
			return 0
		}
		// 简单解析
		var result uint
		for _, c := range val {
			if c >= '0' && c <= '9' {
				result = result*10 + uint(c-'0')
			} else {
				break
			}
		}
		return result
	default:
		return 0
	}
}

// toString 将任意类型转换为 string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case uint:
		return fmt.Sprintf("%d", val)
	default:
		return ""
	}
}
