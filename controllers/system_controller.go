package controllers

import (
	"net/http"
	"strings"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SystemConfigController 系统参数控制器
type SystemConfigController struct {
	DB *gorm.DB
}

// ConfigCategory 系统参数分类
const (
	ConfigCategoryBasic  = "basic"  // 基础配置
	ConfigCategoryEmail  = "email"  // 邮件配置
	ConfigCategorySMS    = "sms"    // SMS配置
	ConfigCategoryPush   = "push"   // 推送配置
)

// GetConfigs 获取系统参数（按分类）
func (c *SystemConfigController) GetConfigs(ctx *gin.Context) {
	category := strings.TrimSpace(ctx.Query("category"))

	query := c.DB.Model(&models.SysConfig{}).Where("status = ?", 1)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var configs []models.SysConfig
	if err := query.Order("id ASC").Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 按分类聚合
	grouped := make(map[string][]models.SysConfig)
	for _, cfg := range configs {
		grouped[cfg.Category] = append(grouped[cfg.Category], cfg)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":     configs,
			"grouped":  grouped,
			"category": category,
		},
	})
}

// GetConfigByKey 根据 key 获取单个配置
func (c *SystemConfigController) GetConfigByKey(ctx *gin.Context) {
	key := ctx.Param("key")

	var cfg models.SysConfig
	if err := c.DB.Where("config_key = ? AND status = ?", key, 1).First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "配置项不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    cfg,
	})
}

// UpdateConfigs 批量更新系统参数
func (c *SystemConfigController) UpdateConfigs(ctx *gin.Context) {
	var req struct {
		Items []struct {
			ConfigKey   string `json:"config_key" binding:"required"`
			ConfigVal   string `json:"config_val"`
			ConfigType  string `json:"config_type"`
			Category    string `json:"category"`
			Remark      string `json:"remark"`
		} `json:"items" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	updated := 0
	created := 0

	for _, item := range req.Items {
		var cfg models.SysConfig
		err := c.DB.Where("config_key = ?", item.ConfigKey).First(&cfg).Error

		if err == gorm.ErrRecordNotFound {
			// 新建
			cfg = models.SysConfig{
				ConfigKey:  item.ConfigKey,
				ConfigVal:  item.ConfigVal,
				ConfigType: item.ConfigType,
				Category:   item.Category,
				Remark:     item.Remark,
				Status:     1,
			}
			if cfg.Category == "" {
				cfg.Category = ConfigCategoryBasic
			}
			if err := c.DB.Create(&cfg).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "创建配置失败: " + err.Error(),
				})
				return
			}
			created++
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询配置失败: " + err.Error(),
			})
			return
		} else {
			// 更新
			cfg.ConfigVal = item.ConfigVal
			if item.ConfigType != "" {
				cfg.ConfigType = item.ConfigType
			}
			if item.Category != "" {
				cfg.Category = item.Category
			}
			if item.Remark != "" {
				cfg.Remark = item.Remark
			}
			if err := c.DB.Save(&cfg).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "更新配置失败: " + err.Error(),
				})
				return
			}
			updated++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"updated": updated,
			"created": created,
		},
	})
}

// DeleteConfig 删除系统参数（软删除标记 status=0）
func (c *SystemConfigController) DeleteConfig(ctx *gin.Context) {
	key := ctx.Param("key")

	result := c.DB.Model(&models.SysConfig{}).
		Where("config_key = ?", key).
		Update("status", 0)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "配置项不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// InitDefaultConfigs 初始化默认系统参数（仅当表为空时）
func (c *SystemConfigController) InitDefaultConfigs() error {
	var count int64
	c.DB.Model(&models.SysConfig{}).Count(&count)
	if count > 0 {
		return nil
	}

	defaults := []models.SysConfig{
		// 基础配置
		{ConfigKey: "system_name", ConfigVal: "MDM管控中心", ConfigType: "string", Category: ConfigCategoryBasic, Status: 1, Remark: "系统名称"},
		{ConfigKey: "system_logo", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryBasic, Status: 1, Remark: "系统Logo"},
		{ConfigKey: "copyright", ConfigVal: "© 2026 MDM", ConfigType: "string", Category: ConfigCategoryBasic, Status: 1, Remark: "版权信息"},
		{ConfigKey: "login_timeout", ConfigVal: "86400", ConfigType: "int", Category: ConfigCategoryBasic, Status: 1, Remark: "登录超时时间(秒)"},
		{ConfigKey: "max_devices_per_user", ConfigVal: "100", ConfigType: "int", Category: ConfigCategoryBasic, Status: 1, Remark: "每个用户最大设备数"},
		// 邮件配置
		{ConfigKey: "smtp_host", ConfigVal: "smtp.example.com", ConfigType: "string", Category: ConfigCategoryEmail, Status: 1, Remark: "SMTP服务器地址"},
		{ConfigKey: "smtp_port", ConfigVal: "587", ConfigType: "int", Category: ConfigCategoryEmail, Status: 1, Remark: "SMTP端口"},
		{ConfigKey: "smtp_user", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryEmail, Status: 1, Remark: "SMTP用户名"},
		{ConfigKey: "smtp_password", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryEmail, Status: 1, Remark: "SMTP密码"},
		{ConfigKey: "smtp_from", ConfigVal: "noreply@example.com", ConfigType: "string", Category: ConfigCategoryEmail, Status: 1, Remark: "发件人地址"},
		{ConfigKey: "smtp_from_name", ConfigVal: "MDM系统", ConfigType: "string", Category: ConfigCategoryEmail, Status: 1, Remark: "发件人名称"},
		// SMS配置
		{ConfigKey: "sms_provider", ConfigVal: "aliyun", ConfigType: "string", Category: ConfigCategorySMS, Status: 1, Remark: "SMS服务商: aliyun/tencent"},
		{ConfigKey: "sms_access_key", ConfigVal: "", ConfigType: "string", Category: ConfigCategorySMS, Status: 1, Remark: "SMS AccessKey"},
		{ConfigKey: "sms_secret_key", ConfigVal: "", ConfigType: "string", Category: ConfigCategorySMS, Status: 1, Remark: "SMS SecretKey"},
		{ConfigKey: "sms_sign_name", ConfigVal: "MDM管控", ConfigType: "string", Category: ConfigCategorySMS, Status: 1, Remark: "SMS签名"},
		// 推送配置
		{ConfigKey: "push_provider", ConfigVal: "fcm", ConfigType: "string", Category: ConfigCategoryPush, Status: 1, Remark: "推送服务商: fcm/apns"},
		{ConfigKey: "fcm_server_key", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryPush, Status: 1, Remark: "FCM Server Key"},
		{ConfigKey: "apns_key_path", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryPush, Status: 1, Remark: "APNs密钥文件路径"},
		{ConfigKey: "apns_key_id", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryPush, Status: 1, Remark: "APNs密钥ID"},
		{ConfigKey: "apns_team_id", ConfigVal: "", ConfigType: "string", Category: ConfigCategoryPush, Status: 1, Remark: "APNs Team ID"},
	}

	for i := range defaults {
		if err := c.DB.Create(&defaults[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

// MenuController 菜单控制器
type MenuController struct {
	DB *gorm.DB
}

// GetMenuTree 获取菜单树
func (c *MenuController) GetMenuTree(ctx *gin.Context) {
	var menus []models.SysMenu
	if err := c.DB.Order("sort ASC").Find(&menus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	tree := buildMenuTree(menus, 0)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": tree,
	})
}

func buildMenuTree(menus []models.SysMenu, parentID uint) []interface{} {
	var result []interface{}
	for _, m := range menus {
		if m.ParentID == parentID {
			item := map[string]interface{}{
				"id":         m.ID,
				"name":       m.Name,
				"path":       m.Path,
				"component":  m.Component,
				"icon":       m.Icon,
				"sort":       m.Sort,
				"visible":    m.Visible,
				"permission": m.Permission,
				"type":       m.Type,
			}
			children := buildMenuTree(menus, m.ID)
			if len(children) > 0 {
				item["children"] = children
			}
			result = append(result, item)
		}
	}
	return result
}

// DictController 字典控制器
type DictController struct {
	DB *gorm.DB
}

// GetDictByType 根据类型获取字典
func (c *DictController) GetDictByType(ctx *gin.Context) {
	dictType := ctx.Param("type")

	var dicts []models.SysDictionary
	if err := c.DB.Where("type = ? AND status = 1", dictType).Order("sort ASC").Find(&dicts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dicts,
	})
}

// LogController 日志控制器
type LogController struct {
	DB *gorm.DB
}

// GetOperationLogs 获取操作日志
func (c *LogController) GetOperationLogs(ctx *gin.Context) {
	query := c.DB.Model(&models.SysOperationLog{})

	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if module := ctx.Query("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	if startTime := ctx.Query("start_time"); startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var logs []models.SysOperationLog
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "20")

	query.Order("id DESC").Offset((systemParseInt(page) - 1) * systemParseInt(pageSize)).Limit(systemParseInt(pageSize)).Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": logs,
			"pagination": gin.H{
				"total":    total,
				"current":  systemParseInt(page),
				"pageSize": systemParseInt(pageSize),
			},
		},
	})
}

// GetLoginLogs 获取登录日志
func (c *LogController) GetLoginLogs(ctx *gin.Context) {
	var logs []models.SysLoginLog
	c.DB.Order("id DESC").Limit(100).Find(&logs)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": logs,
		},
	})
}

func systemParseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
