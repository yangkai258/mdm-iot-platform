package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingsController 系统设置控制器
type SettingsController struct {
	DB *gorm.DB
}

// GetSettings 获取系统设置
// GET /api/v1/settings
func (c *SettingsController) GetSettings(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	var configs []models.SystemConfig
	query := c.DB.Model(&models.SystemConfig{})

	if tenantID != "" {
		query = query.Where("tenant_id = ? OR tenant_id = ''", tenantID)
	} else {
		query = query.Where("tenant_id = ''")
	}

	if group := ctx.Query("group"); group != "" {
		query = query.Where("`group` = ?", group)
	}

	if err := query.Where("status = 1").Order("sort ASC, id ASC").Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	settings := make(map[string]interface{})
	for _, cfg := range configs {
		switch cfg.Type {
		case "number":
			var num int
			for _, ch := range cfg.Value {
				if ch >= '0' && ch <= '9' {
					num = num*10 + int(ch-'0')
				}
			}
			settings[cfg.ConfigKey] = num
		case "boolean":
			settings[cfg.ConfigKey] = cfg.Value == "true" || cfg.Value == "1"
		case "json":
			settings[cfg.ConfigKey] = cfg.Value
		default:
			settings[cfg.ConfigKey] = cfg.Value
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    settings,
	})
}

// GetSetting 获取单个设置
// GET /api/v1/settings/:key
func (c *SettingsController) GetSetting(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	key := ctx.Param("key")

	var config models.SystemConfig
	query := c.DB.Where("config_key = ?", key).Where("status = 1")

	if tenantID != "" {
		query = query.Where("tenant_id = ? OR tenant_id = ''", tenantID)
	} else {
		query = query.Where("tenant_id = ''")
	}

	if err := query.First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配置不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"key":    config.ConfigKey,
			"value":  config.Value,
			"type":   config.Type,
			"label":  config.Label,
			"remark": config.Remark,
		},
	})
}

// UpdateSetting 更新设置
// PUT /api/v1/settings/:key
func (c *SettingsController) UpdateSetting(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	key := ctx.Param("key")

	var req struct {
		Value string `json:"value" binding:"required"`
		Type  string `json:"type"`
		Label string `json:"label"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var config models.SystemConfig
	query := c.DB.Where("config_key = ?", key)

	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	} else {
		query = query.Where("tenant_id = ''")
	}

	if err := query.First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			config = models.SystemConfig{
				ConfigKey: key,
				Value:     req.Value,
				Type:      req.Type,
				Label:     req.Label,
				TenantID:  tenantID,
				Status:    1,
			}
			if config.Type == "" {
				config.Type = "string"
			}
			if err := c.DB.Create(&config).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
			return
		}
	} else {
		updates := map[string]interface{}{
			"value": req.Value,
		}
		if req.Type != "" {
			updates["type"] = req.Type
		}
		if err := c.DB.Model(&config).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"key":   config.ConfigKey,
			"value": req.Value,
		},
	})
}

// BatchUpdateSettings 批量更新设置
// PUT /api/v1/settings
func (c *SettingsController) BatchUpdateSettings(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	var req map[string]interface{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	for key, value := range req {
		strValue := ""
		switch v := value.(type) {
		case string:
			strValue = v
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			if v {
				strValue = "true"
			} else {
				strValue = "false"
			}
		default:
			strValue = ""
		}

		var config models.SystemConfig
		query := c.DB.Where("config_key = ?", key)
		if tenantID != "" {
			query = query.Where("tenant_id = ?", tenantID)
		} else {
			query = query.Where("tenant_id = ''")
		}

		if err := query.First(&config).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				config = models.SystemConfig{
					ConfigKey: key,
					Value:     strValue,
					TenantID:  tenantID,
					Status:    1,
				}
				c.DB.Create(&config)
			}
		} else {
			c.DB.Model(&config).Update("value", strValue)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ListSettingsByGroup 按分组获取设置列表
// GET /api/v1/settings/groups
func (c *SettingsController) ListSettingsByGroup(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	var configs []models.SystemConfig
	query := c.DB.Model(&models.SystemConfig{}).Where("status = 1")

	if tenantID != "" {
		query = query.Where("tenant_id = ? OR tenant_id = ''", tenantID)
	} else {
		query = query.Where("tenant_id = ''")
	}

	if err := query.Order("`group` ASC, sort ASC").Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	groups := make(map[string][]gin.H)
	for _, cfg := range configs {
		groups[cfg.Group] = append(groups[cfg.Group], gin.H{
			"key":    cfg.ConfigKey,
			"value":  cfg.Value,
			"type":   cfg.Type,
			"label":  cfg.Label,
			"remark": cfg.Remark,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    groups,
	})
}
