package controllers

import (
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"mdm-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataMaskingRule 数据脱敏规则
type DataMaskingRule struct {
	ID          uint      `json:"id"`
	Field       string    `json:"field"`        // 字段名
	Pattern     string    `json:"pattern"`      // 正则表达式
	Replacement string    `json:"replacement"` // 替换模式
	Enabled     bool      `json:"enabled"`      // 是否启用
	Description string    `json:"description"`  // 规则描述
	TenantID    string    `json:"tenant_id"`    // 租户ID（空表示全局规则）
	CreatedBy  uint      `json:"created_by"`   // 创建人
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MaskingRuleDB 规则存储模型（用于数据库持久化）
type MaskingRuleDB struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Field       string    `gorm:"type:varchar(50);not null" json:"field"`
	Pattern     string    `gorm:"type:varchar(255);not null" json:"pattern"`
	Replacement string    `gorm:"type:varchar(255);not null" json:"replacement"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	TenantID    string    `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedBy  uint      `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (MaskingRuleDB) TableName() string {
	return "data_masking_rules"
}

// DataMaskingController 数据脱敏控制器
type DataMaskingController struct {
	DB *gorm.DB
}

// 内存中的规则缓存（生产环境建议使用 Redis）
var (
	rulesCache      = make(map[string]*middleware.DataMaskingRule)
	rulesCacheMutex sync.RWMutex
)

// GetMaskingRules 获取脱敏规则列表
// @Summary 获取脱敏规则列表
// @Description 获取所有数据脱敏规则
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param field query string false "字段名筛选"
// @Param enabled query bool false "是否启用"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/rules [get]
func (c *DataMaskingController) GetMaskingRules(ctx *gin.Context) {
	page := parseIntDefault(ctx.Query("page"), 1)
	pageSize := parseIntDefault(ctx.Query("page_size"), 20)
	field := ctx.Query("field")
	enabledStr := ctx.Query("enabled")

	var rules []MaskingRuleDB
	query := c.DB.Model(&MaskingRuleDB{})

	if field != "" {
		query = query.Where("field LIKE ?", "%"+field+"%")
	}
	if enabledStr != "" {
		enabled := enabledStr == "true"
		query = query.Where("enabled = ?", enabled)
	}

	var total int64
	query.Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rules).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取脱敏规则失败",
			"error":   err.Error(),
		})
		return
	}

	// 转换为API响应格式
	result := make([]map[string]interface{}, 0, len(rules))
	for _, r := range rules {
		result = append(result, map[string]interface{}{
			"id":          r.ID,
			"field":       r.Field,
			"pattern":     r.Pattern,
			"replacement": r.Replacement,
			"enabled":     r.Enabled,
			"description": r.Description,
			"tenant_id":   r.TenantID,
			"created_at":  r.CreatedAt,
			"updated_at":  r.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items": result,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreateMaskingRule 创建脱敏规则
// @Summary 创建脱敏规则
// @Description 创建新的数据脱敏规则
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Param rule body DataMaskingRuleInput true "脱敏规则"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/rules [post]
func (c *DataMaskingController) CreateMaskingRule(ctx *gin.Context) {
	var input struct {
		Field       string `json:"field" binding:"required"`
		Pattern     string `json:"pattern" binding:"required"`
		Replacement string `json:"replacement" binding:"required"`
		Enabled     bool   `json:"enabled"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证正则表达式是否有效
	if _, err := regexp.Compile(input.Pattern); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的正则表达式",
			"error":   err.Error(),
		})
		return
	}

	// 获取当前用户
	userID, _ := ctx.Get("user_id")
	tenantID, _ := ctx.Get("tenant_id")

	rule := MaskingRuleDB{
		Field:       input.Field,
		Pattern:     input.Pattern,
		Replacement: input.Replacement,
		Enabled:     input.Enabled,
		Description: input.Description,
	}

	// 设置租户ID
	if tenantIDStr, ok := tenantID.(string); ok && tenantIDStr != "" {
		rule.TenantID = tenantIDStr
	}

	// 设置创建人
	if uid, ok := userID.(uint); ok {
		rule.CreatedBy = uid
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建脱敏规则失败",
			"error":   err.Error(),
		})
		return
	}

	// 更新内存缓存
	engine := middleware.GetMaskingEngine()
	if rule.Enabled {
		engine.RegisterRule(rule.Field, rule.Pattern, rule.Replacement)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "创建脱敏规则成功",
		"data": map[string]interface{}{
			"id":          rule.ID,
			"field":       rule.Field,
			"pattern":     rule.Pattern,
			"replacement": rule.Replacement,
			"enabled":     rule.Enabled,
			"description": rule.Description,
			"created_at":  rule.CreatedAt,
		},
	})
}

// UpdateMaskingRule 更新脱敏规则
// @Summary 更新脱敏规则
// @Description 更新指定ID的脱敏规则
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Param id path int true "规则ID"
// @Param rule body DataMaskingRuleInput true "脱敏规则"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/rules/{id} [put]
func (c *DataMaskingController) UpdateMaskingRule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的规则ID",
		})
		return
	}

	var input struct {
		Field       string `json:"field"`
		Pattern     string `json:"pattern"`
		Replacement string `json:"replacement"`
		Enabled     *bool  `json:"enabled"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查找现有规则
	var rule MaskingRuleDB
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "规则不存在",
		})
		return
	}

	// 验证正则表达式（如果提供）
	if input.Pattern != "" {
		if _, err := regexp.Compile(input.Pattern); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的正则表达式",
				"error":   err.Error(),
			})
			return
		}
		rule.Pattern = input.Pattern
	}

	// 更新字段
	if input.Field != "" {
		rule.Field = input.Field
	}
	if input.Replacement != "" {
		rule.Replacement = input.Replacement
	}
	if input.Enabled != nil {
		rule.Enabled = *input.Enabled
	}
	if input.Description != "" {
		rule.Description = input.Description
	}

	if err := c.DB.Save(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新脱敏规则失败",
			"error":   err.Error(),
		})
		return
	}

	// 更新内存缓存
	engine := middleware.GetMaskingEngine()
	if rule.Enabled {
		engine.RegisterRule(rule.Field, rule.Pattern, rule.Replacement)
	} else {
		engine.UnregisterRule(rule.Field)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "更新脱敏规则成功",
		"data": map[string]interface{}{
			"id":          rule.ID,
			"field":       rule.Field,
			"pattern":     rule.Pattern,
			"replacement": rule.Replacement,
			"enabled":     rule.Enabled,
			"description": rule.Description,
			"updated_at":  rule.UpdatedAt,
		},
	})
}

// DeleteMaskingRule 删除脱敏规则
// @Summary 删除脱敏规则
// @Description 删除指定ID的脱敏规则
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Param id path int true "规则ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/rules/{id} [delete]
func (c *DataMaskingController) DeleteMaskingRule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的规则ID",
		})
		return
	}

	// 查找规则
	var rule MaskingRuleDB
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "规则不存在",
		})
		return
	}

	// 删除规则
	if err := c.DB.Delete(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除脱敏规则失败",
			"error":   err.Error(),
		})
		return
	}

	// 从缓存移除
	engine := middleware.GetMaskingEngine()
	engine.UnregisterRule(rule.Field)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除脱敏规则成功",
	})
}

// GetDefaultRules 获取默认脱敏规则
// @Summary 获取默认脱敏规则
// @Description 获取系统内置的默认脱敏规则
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/defaults [get]
func (c *DataMaskingController) GetDefaultRules(ctx *gin.Context) {
	defaults := []map[string]interface{}{
		{
			"field":        "phone",
			"pattern":      `^(\d{3})\d{4}(\d{4})$`,
			"replacement": "$1****$2",
			"description": "手机号脱敏：13812345678 -> 138****5678",
		},
		{
			"field":        "id_card",
			"pattern":      `^(\d{6})\d{8}(\d{4})$`,
			"replacement": "$1********$2",
			"description": "身份证脱敏：显示前6位和后4位",
		},
		{
			"field":        "email",
			"pattern":      `^(\w{1})\w+(\w+@\w+\.\w+)$`,
			"replacement": "$1***$2",
			"description": "邮箱脱敏：user@example.com -> u***@example.com",
		},
		{
			"field":        "bank_card",
			"pattern":      `^(\d{6})\d+(\d{4})$`,
			"replacement": "$1*********$2",
			"description": "银行卡脱敏：显示前6位和后4位",
		},
		{
			"field":        "address",
			"pattern":      `^(.{2}[省市区县]).*$`,
			"replacement": "$1***",
			"description": "地址脱敏：只显示省市區",
		},
		{
			"field":        "real_name",
			"pattern":      `^(.)(.*)(.)$`,
			"replacement": "$1***$3",
			"description": "姓名脱敏：显示首尾字符",
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": defaults,
	})
}

// TestMaskingRule 测试脱敏规则
// @Summary 测试脱敏规则
// @Description 测试指定的脱敏规则效果
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Param test body MaskingRuleTestInput true "测试数据"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/test [post]
func (c *DataMaskingController) TestMaskingRule(ctx *gin.Context) {
	var input struct {
		Field       string `json:"field" binding:"required"`
		Pattern     string `json:"pattern" binding:"required"`
		Replacement string `json:"replacement" binding:"required"`
		TestValue   string `json:"test_value" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证正则表达式
	re, err := regexp.Compile(input.Pattern)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的正则表达式",
			"error":   err.Error(),
		})
		return
	}

	// 执行脱敏
	maskedValue := re.ReplaceAllString(input.TestValue, input.Replacement)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"original": input.TestValue,
			"masked":   maskedValue,
			"valid":     true,
		},
	})
}

// MaskDataRequest 脱敏请求
type MaskDataRequest struct {
	Data           map[string]interface{} `json:"data" binding:"required"`
	SensitiveFields []string             `json:"sensitive_fields" binding:"required"`
}

// MaskData 手动脱敏数据
// @Summary 手动脱敏数据
// @Description 对指定数据进行脱敏处理
// @Tags 数据脱敏
// @Accept json
// @Produce json
// @Param request body MaskDataRequest true "脱敏请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-masking/mask [post]
func (c *DataMaskingController) MaskData(ctx *gin.Context) {
	var input MaskDataRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	engine := middleware.GetMaskingEngine()
	masked := engine.MaskMap(input.Data, input.SensitiveFields)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": masked,
	})
}

// InitMaskingRulesFromDB 从数据库加载脱敏规则到内存
func (c *DataMaskingController) InitMaskingRulesFromDB() error {
	var rules []MaskingRuleDB
	if err := c.DB.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return err
	}

	engine := middleware.GetMaskingEngine()
	for _, r := range rules {
		engine.RegisterRule(r.Field, r.Pattern, r.Replacement)
	}

	return nil
}
