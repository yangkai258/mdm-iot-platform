package controllers

import (
	"net/http"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ctrl *AuthController) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", ctrl.Login)
		auth.POST("/refresh", ctrl.RefreshToken)
		auth.POST("/logout", ctrl.Logout)
		auth.GET("/me", ctrl.Me)
		auth.GET("/menu", ctrl.Menu)
	}
}

// Login 用户登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 查找用户
	var user models.SysUser
	if err := ctrl.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 密码校验（bcrypt）
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 生成 Token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID, user.TenantID, false, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	refreshToken, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID, user.TenantID, false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成refresh token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"token":         token,
			"refresh_token": refreshToken,
			"user_id":       user.ID,
			"username":      user.Username,
			"expires_in":    3600,
		},
	})
}

// Logout 用户登出
func (ctrl *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 20000, "message": "登出成功"})
}

// Me 获取当前用户信息
func (ctrl *AuthController) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.SysUser
	if err := ctrl.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
			"phone":    user.Phone,
			"role_id":  user.RoleID,
			"tenant_id": user.TenantID,
		},
	})
}

// Menu 获取用户菜单
func (ctrl *AuthController) Menu(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": []gin.H{
			{"path": "/dashboard", "name": "Dashboard", "icon": "dashboard"},
			{"path": "/device", "name": "设备管理", "icon": "device"},
			{"path": "/member", "name": "会员管理", "icon": "member"},
			{"path": "/pet", "name": "宠物管理", "icon": "pet"},
			{"path": "/alerts", "name": "告警中心", "icon": "alert"},
			{"path": "/ota", "name": "OTA升级", "icon": "upgrade"},
		},
	})
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
		"code": 20000,
		"data": gin.H{
			"access_token":  newToken,
			"refresh_token": newRefreshToken,
			"expires_in":    3600,
		},
	})
}