package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Config 配置结构体
type Config struct {
	JWTSecret string `json:"jwt_secret"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

var jwtSecret = []byte(getJWTSecret())

// getJWTSecret 从配置文件或环境变量获取 JWT 密钥
func getJWTSecret() string {
	// 优先从配置文件读取
	configPaths := []string{
		"config.json",
		"../config.json",
		"../../config.json",
		os.Getenv("CONFIG_PATH"),
	}

	for _, path := range configPaths {
		if path == "" {
			continue
		}
		if cfg, err := LoadConfig(path); err == nil && cfg.JWTSecret != "" {
			return cfg.JWTSecret
		}
	}

	// 降级到环境变量
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set and no config file found")
	}
	return secret
}

// JWTClaims JWT 载荷
type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	RoleID       uint   `json:"role_id"`
	TenantID     string `json:"tenant_id"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	jwt.RegisteredClaims
}

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除登录接口和健康检查
		if c.Request.URL.Path == "/api/v1/auth/login" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "未授权，请先登录",
			})
			c.Abort()
			return
		}

		// 去掉 Bearer 前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// 解析 Token
		claims := &JWTClaims{}
		tokenObj, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !tokenObj.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"message": "Token 已过期，请重新登录",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("claims", map[string]interface{}{
			"user_id":       claims.UserID,
			"username":      claims.Username,
			"role_id":       claims.RoleID,
			"tenant_id":     claims.TenantID,
			"is_super_admin": claims.IsSuperAdmin,
		})

		c.Next()
	}
}

// GenerateToken 生成 Token（支持 tenant_id 和 is_super_admin）
func GenerateToken(userID uint, username string, roleID uint, tenantID string, isSuperAdmin bool) (string, error) {
	claims := JWTClaims{
		UserID:       userID,
		Username:     username,
		RoleID:       roleID,
		TenantID:     tenantID,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
