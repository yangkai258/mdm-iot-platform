package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mdm-backend/ldap"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LDAPController LDAP 管理控制器
type LDAPController struct {
	DB *gorm.DB
}

// GetLDAPConfig 获取 LDAP 配置
func (c *LDAPController) GetLDAPConfig(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	var config models.LDAPConfig
	if err := c.DB.Where("tenant_id = ?", tenantID).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认空配置
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"data":    nil,
				"message": "暂无配置",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取配置失败: " + err.Error(),
		})
		return
	}

	// 不返回加密密码
	config.BindPassword = ""
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    config,
		"message": "success",
	})
}

// UpdateLDAPConfig 更新 LDAP 配置
func (c *LDAPController) UpdateLDAPConfig(ctx *gin.Context) {
	var req models.LDAPConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	userID := middleware.GetUserID(ctx)

	// 查找现有配置
	var config models.LDAPConfig
	err := c.DB.Where("tenant_id = ?", tenantID).First(&config).Error
	isCreate := false
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			isCreate = true
			config = models.LDAPConfig{}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询配置失败: " + err.Error(),
			})
			return
		}
	}

	// 更新字段
	config.ConfigName = req.ConfigName
	config.Host = req.Host
	config.Port = req.Port
	config.BaseDN = req.BaseDN
	config.BindDN = req.BindDN
	config.UseSSL = req.UseSSL
	config.UseTLS = req.UseTLS
	config.UserFilter = req.UserFilter
	if config.UserFilter == "" {
		config.UserFilter = "(objectClass=user)"
	}
	config.GroupFilter = req.GroupFilter
	if config.GroupFilter == "" {
		config.GroupFilter = "(objectClass=group)"
	}
	config.SyncInterval = req.SyncInterval
	config.IsEnabled = req.IsEnabled
	config.Description = req.Description
	config.TenantID = tenantID
	config.CreatedBy = userID
	config.Status = "inactive"

	// 如果传了新密码，需要加密存储
	if req.BindPassword != "" {
		encrypted, err := utils.EncryptAES(req.BindPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码加密失败: " + err.Error(),
			})
			return
		}
		config.BindPassword = encrypted
	}

	if isCreate {
		if err := c.DB.Create(&config).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建配置失败: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.DB.Save(&config).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新配置失败: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    config,
		"message": "保存成功",
	})
}

// TestLDAPConnection 测试 LDAP 连接
func (c *LDAPController) TestLDAPConnection(ctx *gin.Context) {
	var req models.LDAPConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 临时创建 LDAPService 测试连接
	ldapService := &ldap.LDAPService{
		Host:         req.Host,
		Port:         req.Port,
		BaseDN:       req.BaseDN,
		BindDN:       req.BindDN,
		BindPassword: req.BindPassword,
		UseSSL:       req.UseSSL,
		UseTLS:       req.UseTLS,
		UserFilter:   req.UserFilter,
		GroupFilter:  req.GroupFilter,
	}

	result, err := ldapService.TestConnection()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"data":    result,
			"message": "测试完成",
		})
		return
	}

	if !result.Success {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    result,
			"message": result.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    result,
		"message": "连接成功",
	})
}

// GetLDAPUsers 获取 LDAP 用户列表
func (c *LDAPController) GetLDAPUsers(ctx *gin.Context) {
	query := ctx.Query("query")
	page := parseIntDefault(ctx.Query("page"), 1)
	pageSize := parseIntDefault(ctx.Query("page_size"), 20)

	// 获取 LDAP 配置
	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	var config models.LDAPConfig
	if err := c.DB.Where("tenant_id = ? AND is_enabled = ?", tenantID, true).First(&config).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "LDAP 未配置或未启用",
		})
		return
	}

	// 解密密码
	decryptedPass, err := utils.DecryptAES(config.BindPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码解密失败",
		})
		return
	}

	ldapService := &ldap.LDAPService{
		Host:         config.Host,
		Port:         config.Port,
		BaseDN:       config.BaseDN,
		BindDN:       config.BindDN,
		BindPassword: decryptedPass,
		UseSSL:       config.UseSSL,
		UseTLS:       config.UseTLS,
		UserFilter:   config.UserFilter,
		GroupFilter:  config.GroupFilter,
	}

	users, err := ldapService.SearchUsers(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索用户失败: " + err.Error(),
		})
		return
	}

	// 分页
	total := len(users)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		users = []ldap.LDAPUser{}
	} else if end > total {
		users = users[start:]
	} else {
		users = users[start:end]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  users,
			"total": total,
			"page":  page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// SyncLDAPUsers 同步 LDAP 用户
func (c *LDAPController) SyncLDAPUsers(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}
	_ = middleware.GetUserID(ctx)     // 保留上下文
	_ = middleware.GetUsername(ctx)    // 保留上下文

	// 获取 LDAP 配置
	var config models.LDAPConfig
	if err := c.DB.Where("tenant_id = ? AND is_enabled = ?", tenantID, true).First(&config).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "LDAP 未配置或未启用",
		})
		return
	}

	// 解密密码
	decryptedPass, err := utils.DecryptAES(config.BindPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码解密失败",
		})
		return
	}

	ldapService := &ldap.LDAPService{
		Host:         config.Host,
		Port:         config.Port,
		BaseDN:       config.BaseDN,
		BindDN:       config.BindDN,
		BindPassword: decryptedPass,
		UseSSL:       config.UseSSL,
		UseTLS:       config.UseTLS,
		UserFilter:   config.UserFilter,
		GroupFilter:  config.GroupFilter,
	}

	users, err := ldapService.SearchUsers("")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "同步用户失败: " + err.Error(),
		})
		return
	}

	// 同步到本地数据库
	var added, updated, skipped int
	var errors []string

	for _, user := range users {
		var mapping models.LDAPUserMapping
		err := c.DB.Where("ldap_dn = ? AND tenant_id = ?", user.DN, tenantID).First(&mapping).Error

		groupsJSON, _ := json.Marshal(user.Groups)

		if err == gorm.ErrRecordNotFound {
			// 新增
			mapping = models.LDAPUserMapping{
				LDAPDN:       user.DN,
				Username:     user.Username,
				Email:        user.Email,
				DisplayName:  user.DisplayName,
				LDAPGroups:   string(groupsJSON),
				SyncStatus:   "synced",
				LastSyncedAt: nil,
				TenantID:     tenantID,
			}
			if err := c.DB.Create(&mapping).Error; err != nil {
				errors = append(errors, "创建用户 "+user.Username+" 失败: "+err.Error())
				skipped++
			} else {
				added++
			}
		} else if err != nil {
			errors = append(errors, "查询用户 "+user.Username+" 失败: "+err.Error())
			skipped++
		} else {
			// 更新
			mapping.Username = user.Username
			mapping.Email = user.Email
			mapping.DisplayName = user.DisplayName
			mapping.LDAPGroups = string(groupsJSON)
			mapping.SyncStatus = "synced"
			if err := c.DB.Save(&mapping).Error; err != nil {
				errors = append(errors, "更新用户 "+user.Username+" 失败: "+err.Error())
				skipped++
			} else {
				updated++
			}
		}
	}

	// 更新最后同步时间
	now := time.Now()
	config.LastSyncAt = &now
	config.Status = "active"
	c.DB.Save(&config)

	result := ldap.SyncResult{
		TotalUsers: len(users),
		Added:      added,
		Updated:    updated,
		Skipped:    skipped,
		Errors:     errors,
		SyncedAt:   now,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    result,
		"message": "同步完成",
	})
}

// GetLDAPGroups 获取 LDAP 分组列表
func (c *LDAPController) GetLDAPGroups(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	var config models.LDAPConfig
	if err := c.DB.Where("tenant_id = ? AND is_enabled = ?", tenantID, true).First(&config).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "LDAP 未配置或未启用",
		})
		return
	}

	decryptedPass, err := utils.DecryptAES(config.BindPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码解密失败",
		})
		return
	}

	ldapService := &ldap.LDAPService{
		Host:         config.Host,
		Port:         config.Port,
		BaseDN:       config.BaseDN,
		BindDN:       config.BindDN,
		BindPassword: decryptedPass,
		UseSSL:       config.UseSSL,
		UseTLS:       config.UseTLS,
		UserFilter:   config.UserFilter,
		GroupFilter:  config.GroupFilter,
	}

	groups, err := ldapService.SearchGroups()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索分组失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  groups,
			"total": len(groups),
		},
		"message": "success",
	})
}

// SetGroupRoleMapping 设置分组-角色映射
func (c *LDAPController) SetGroupRoleMapping(ctx *gin.Context) {
	var req models.LDAPGroupMappingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	// 验证角色存在
	var role models.Role
	if err := c.DB.Where("id = ? AND tenant_id = ?", req.RoleID, tenantID).First(&role).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "角色不存在",
		})
		return
	}

	// 查找现有映射或创建
	var mapping models.LDAPGroupRoleMapping
	err := c.DB.Where("ldap_group_dn = ? AND tenant_id = ?", req.LDAPGroupDN, tenantID).First(&mapping).Error
	isCreate := false
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			isCreate = true
			mapping = models.LDAPGroupRoleMapping{}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询映射失败: " + err.Error(),
			})
			return
		}
	}

	mapping.LDAPGroupDN = req.LDAPGroupDN
	mapping.LDAPGroupName = req.LDAPGroupName
	mapping.RoleID = req.RoleID
	mapping.RoleName = role.RoleName
	mapping.TenantID = tenantID

	if isCreate {
		if err := c.DB.Create(&mapping).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建映射失败: " + err.Error(),
			})
			return
		}
	} else {
		if err := c.DB.Save(&mapping).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新映射失败: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    mapping,
		"message": "设置成功",
	})
}

// GetGroupRoleMappings 获取所有分组-角色映射
func (c *LDAPController) GetGroupRoleMappings(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	var mappings []models.LDAPGroupRoleMapping
	if err := c.DB.Where("tenant_id = ?", tenantID).Find(&mappings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询映射失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  mappings,
			"total": len(mappings),
		},
		"message": "success",
	})
}

// parseIntDefault 解析整数，带默认值
func parseIntDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			break
		}
	}
	return result
}
