package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthController 认证控制器
type AuthController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// 登录限流配置
const (
	loginMaxAttempts   = 3       // 5分钟内最大失败次数
	loginBlockMinutes = 5       // 锁定分钟数
	loginWindowMinutes = 5      // 时间窗口分钟数
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 登录限流检查（基于 IP + 用户名）
	if c.Redis != nil {
		loginKey := fmt.Sprintf("login:ratelimit:%s:%s", ctx.ClientIP(), req.Username)
		ctx2 := context.Background()

		// 获取当前失败次数
		failedAttempts, err := c.Redis.Client().Get(ctx2, loginKey).Int()
		if err == nil && failedAttempts >= loginMaxAttempts {
			// 获取锁定剩余时间
			ttl, _ := c.Redis.Client().TTL(ctx2, loginKey).Result()
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": fmt.Sprintf("登录失败次数过多，请 %d 分钟后再试", int(ttl.Minutes())+1),
			})
			return
		}
	}

	// 查询用户
	var user models.SysUser
	result := c.DB.Where("username = ?", req.Username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// 登录限流：记录失败次数（防止暴力猜测用户名）
		if c.Redis != nil {
			loginKey := fmt.Sprintf("login:ratelimit:%s:%s", ctx.ClientIP(), req.Username)
			ctx2 := context.Background()
			c.Redis.Client().Incr(ctx2, loginKey)
			c.Redis.Client().Expire(ctx2, loginKey, loginWindowMinutes*time.Minute)
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 登录限流：记录失败次数
		if c.Redis != nil {
			loginKey := fmt.Sprintf("login:ratelimit:%s:%s", ctx.ClientIP(), req.Username)
			ctx2 := context.Background()
			// 增加失败计数，并设置过期时间
			c.Redis.Client().Incr(ctx2, loginKey)
			c.Redis.Client().Expire(ctx2, loginKey, loginWindowMinutes*time.Minute)
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "账号已被禁用",
		})
		return
	}

	// 生成 Token（携带 tenant_id）
	token, err := middleware.GenerateToken(user.ID, user.Username, user.RoleID, "", false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成 Token 失败",
		})
		return
	}

	// 记录登录日志
	loginLog := models.SysLoginLog{
		UserID:   user.ID,
		Username: user.Username,
		IP:       ctx.ClientIP(),
		Status:   1,
		Msg:      "登录成功",
	}
	c.DB.Create(&loginLog)

	// 记录活动日志（ActivityLog）
	RecordActivity(c.DB, user.ID, user.Username, "login", "auth", user.ID, user.Username, map[string]interface{}{
		"browser": ctx.GetHeader("User-Agent"),
	}, ctx.ClientIP())

	// 记录到独立 login_logs 表
	loginLog2 := models.LoginLog{
		UserID:    user.ID,
		Username:  user.Username,
		IP:        ctx.ClientIP(),
		Status:    1,
		Msg:       "登录成功",
		TenantID:  user.TenantID,
		LoginTime: time.Now(),
	}
	c.DB.Create(&loginLog2)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.Nickname,
				"email":    user.Email,
				"role_id":  user.RoleID,
			},
		},
	})
}

// Logout 登出
func (c *AuthController) Logout(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("user_id")
	usernameVal, _ := ctx.Get("username")
	uid := uint(userIDVal.(int))
	uname := usernameVal.(string)

	// 记录登出日志
	loginLog := models.SysLoginLog{
		UserID:   uid,
		Username: uname,
		IP:       ctx.ClientIP(),
		Status:   1,
		Msg:      "退出登录",
	}
	c.DB.Create(&loginLog)

	// 记录活动日志
	RecordActivity(c.DB, uid, uname, "logout", "auth", uid, uname, nil, ctx.ClientIP())

	// 记录到独立 login_logs 表
	loginLog2 := models.LoginLog{
		UserID:    uid,
		Username:  uname,
		IP:        ctx.ClientIP(),
		Status:    1,
		Msg:       "退出登录",
		LoginTime: time.Now(),
	}
	c.DB.Create(&loginLog2)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// GetUserInfo 获取用户信息
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var user models.SysUser
	if err := c.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"nickname":  user.Nickname,
			"email":     user.Email,
			"phone":     user.Phone,
			"role_id":   user.RoleID,
			"created_at": user.CreatedAt,
		},
	})
}
