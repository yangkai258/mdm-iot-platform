package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ComplianceController 合规策略控制器
type ComplianceController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册合规路由
func (c *ComplianceController) RegisterRoutes(api *gin.RouterGroup) {
	compliance := api.Group("/compliance")
	{
		compliance.GET("/policies", c.ListPolicies)
		compliance.GET("/policies/:id", c.GetPolicy)
		compliance.POST("/policies", c.CreatePolicy)
		compliance.PUT("/policies/:id", c.UpdatePolicy)
		compliance.DELETE("/policies/:id", c.DeletePolicy)
	}

	// 合规规则
	api.GET("/compliance-rules", c.ListRules)
	api.POST("/compliance-rules", c.CreateRule)
	api.PUT("/compliance-rules/:id", c.UpdateRule)
	api.DELETE("/compliance-rules/:id", c.DeleteRule)
}

// ListPolicies 获取合规策略列表
func (c *ComplianceController) ListPolicies(ctx *gin.Context) {
	var policies []models.CompliancePolicy
	var total int64

	query := c.DB.Model(&models.CompliancePolicy{})

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&policies).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":       policies,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
		},
		"message": "success",
	})
}

// GetPolicy 获取策略详情
func (c *ComplianceController) GetPolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.CompliancePolicy
	if err := c.DB.Where("id = ?", id).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": policy, "message": "success"})
}

// CreatePolicy 创建策略
func (c *ComplianceController) CreatePolicy(ctx *gin.Context) {
	var input struct {
		Name              string `json:"name" binding:"required"`
		Description       string `json:"description"`
		PolicyType        string `json:"policy_type" binding:"required"`
		TargetValue       string `json:"target_value"`
		Condition         string `json:"condition" binding:"required"`
		Severity          int    `json:"severity"`
		RemediationAction string `json:"remediation_action"`
		Enabled           *bool  `json:"enabled"`
		EnforceScope      string `json:"enforce_scope"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	policy := models.CompliancePolicy{
		Name:              input.Name,
		Description:       input.Description,
		PolicyType:        input.PolicyType,
		TargetValue:       input.TargetValue,
		Condition:         input.Condition,
		RemediationAction: input.RemediationAction,
		EnforceScope:      input.EnforceScope,
	}

	if input.Severity > 0 {
		policy.Severity = input.Severity
	}
	if input.Enabled != nil {
		policy.Enabled = *input.Enabled
	} else {
		policy.Enabled = true
	}

	if err := c.DB.Create(&policy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "data": policy, "message": "success"})
}

// UpdatePolicy 更新策略
func (c *ComplianceController) UpdatePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.CompliancePolicy
	if err := c.DB.Where("id = ?", id).First(&policy).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "策略不存在"})
		return
	}

	var input struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		PolicyType        string `json:"policy_type"`
		TargetValue       string `json:"target_value"`
		Condition         string `json:"condition"`
		Severity          *int   `json:"severity"`
		RemediationAction string `json:"remediation_action"`
		Enabled           *bool  `json:"enabled"`
		EnforceScope      string `json:"enforce_scope"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.PolicyType != "" {
		updates["policy_type"] = input.PolicyType
	}
	if input.TargetValue != "" {
		updates["target_value"] = input.TargetValue
	}
	if input.Condition != "" {
		updates["condition"] = input.Condition
	}
	if input.Severity != nil {
		updates["severity"] = *input.Severity
	}
	if input.RemediationAction != "" {
		updates["remediation_action"] = input.RemediationAction
	}
	if input.Enabled != nil {
		updates["enabled"] = *input.Enabled
	}
	if input.EnforceScope != "" {
		updates["enforce_scope"] = input.EnforceScope
	}

	if err := c.DB.Model(&policy).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": policy, "message": "success"})
}

// DeletePolicy 删除策略
func (c *ComplianceController) DeletePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.DB.Where("id = ?", id).Delete(&models.CompliancePolicy{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ListRules 获取合规规则列表
func (c *ComplianceController) ListRules(ctx *gin.Context) {
	var rules []models.PolicyConfig
	var total int64

	query := c.DB.Model(&models.PolicyConfig{})

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":     rules,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// CreateRule 创建规则
func (c *ComplianceController) CreateRule(ctx *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		ConfigType  string `json:"config_type" binding:"required"`
		SubType     string `json:"sub_type"`
		Description string `json:"description"`
		ConfigData  string `json:"config_data"`
		Enabled     *bool  `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	rule := models.PolicyConfig{
		Name:        input.Name,
		ConfigType:  input.ConfigType,
		SubType:     input.SubType,
		Description: input.Description,
		ConfigData:  input.ConfigData,
	}

	if input.Enabled != nil {
		rule.Enabled = *input.Enabled
	} else {
		rule.Enabled = true
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "data": rule, "message": "success"})
}

// UpdateRule 更新规则
func (c *ComplianceController) UpdateRule(ctx *gin.Context) {
	id := ctx.Param("id")

	var rule models.PolicyConfig
	if err := c.DB.Where("id = ?", id).First(&rule).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		ConfigType  string `json:"config_type"`
		SubType     string `json:"sub_type"`
		Description string `json:"description"`
		ConfigData  string `json:"config_data"`
		Enabled     *bool  `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.ConfigType != "" {
		updates["config_type"] = input.ConfigType
	}
	if input.SubType != "" {
		updates["sub_type"] = input.SubType
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.ConfigData != "" {
		updates["config_data"] = input.ConfigData
	}
	if input.Enabled != nil {
		updates["enabled"] = *input.Enabled
	}

	if err := c.DB.Model(&rule).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule, "message": "success"})
}

// DeleteRule 删除规则
func (c *ComplianceController) DeleteRule(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.DB.Where("id = ?", id).Delete(&models.PolicyConfig{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
