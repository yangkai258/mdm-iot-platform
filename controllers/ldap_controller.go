package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// LDAPController LDAP 认证控制器
type LDAPController struct{}

func NewLDAPController() *LDAPController {
	return &LDAPController{}
}

// GetConfig 获取 LDAP 配置
// GET /api/v1/ldap/config
func (c *LDAPController) GetConfig(ctx *gin.Context) {
	var config models.LDAPConfig
	if err := models.DB.First(&config).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
			"host":     "",
			"port":     389,
			"base_dn":  "",
			"enabled":  false,
		}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
}

// UpdateConfig 更新 LDAP 配置
// PUT /api/v1/ldap/config
func (c *LDAPController) UpdateConfig(ctx *gin.Context) {
	var input struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		BaseDN       string `json:"base_dn"`
		BindDN       string `json:"bind_dn"`
		BindPassword string `json:"bind_password"`
		UserFilter   string `json:"user_filter"`
		Enabled      bool   `json:"enabled"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	var config models.LDAPConfig
	if err := models.DB.FirstOrCreate(&config, models.LDAPConfig{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	config.Host = input.Host
	config.Port = input.Port
	config.BaseDN = input.BaseDN
	config.BindDN = input.BindDN
	config.UserFilter = input.UserFilter
	config.Enabled = input.Enabled
	models.DB.Save(&config)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "LDAP config updated"})
}

// TestConnection 测试 LDAP 连接
// POST /api/v1/ldap/test
func (c *LDAPController) TestConnection(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "LDAP connection test passed (simulated)"})
}

// GetUsers 获取 LDAP 用户列表
// GET /api/v1/ldap/users
func (c *LDAPController) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": []gin.H{
		{"username": "test_user", "email": "test@example.com", "display_name": "Test User"},
	}})
}
