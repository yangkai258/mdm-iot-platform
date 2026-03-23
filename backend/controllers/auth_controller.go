package controllers

import (
	"net/http"
	"time"

	"mdm-backend/middleware"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ctrl *AuthController) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/refresh", ctrl.RefreshToken)
	}
}

// RefreshToken 刷新 Access Token
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证 refresh token
	claims, err := middleware.ParseToken(req.RefreshToken)
	if err != nil || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "refresh token无效"})
		return
	}

	// 检查是否过期（refresh token通常7天）
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "refresh token已过期，请重新登录"})
		return
	}

	// 生成新的 access token
	newToken, err := middleware.GenerateToken(claims.UserID, claims.Username, claims.RoleID, claims.TenantID, claims.IsSuperAdmin, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	// 生成新的 refresh token
	newRefreshToken, err := middleware.GenerateToken(claims.UserID, claims.Username, claims.RoleID, claims.TenantID, claims.IsSuperAdmin, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成refresh token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"access_token":  newToken,
			"refresh_token": newRefreshToken,
			"expires_in":   3600,
		},
	})
}
