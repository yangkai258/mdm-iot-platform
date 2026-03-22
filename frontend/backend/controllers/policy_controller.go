package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PolicyController 策略管理控制器
type PolicyController struct {
	DB *gorm.DB
}

// ===== PolicyConfig CRUD =====

// ListConfigs 获取配置文件列表
func (c *PolicyController) ListConfigs(ctx *gin.Context) {
	var configs []models.PolicyConfig
	query := c.DB.Model(&models.PolicyConfig{})

	// 关键词过滤
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	// 类型过滤
	if configType := ctx.Query("config_type"); configType != "" {
		query = query.Where("config_type = ?", configType)
	}
	// 启用状态过滤
	if enabled := ctx.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch configs"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      configs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateConfig 创建配置文件
func (c *PolicyController) CreateConfig(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		ConfigType  string `json:"config_type" binding:"required"`
		SubType     string `json:"sub_type"`
		Description string `json:"description"`
		ConfigData  string `json:"config_data"`
		Enabled     *bool  `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	// 验证 config_data 是否为有效 JSON
	if req.ConfigData != "" {
		var dummy interface{}
		if err := json.Unmarshal([]byte(req.ConfigData), &dummy); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "config_data must be valid JSON"})
			return
		}
	}

	config := models.PolicyConfig{
		Name:        req.Name,
		ConfigType:  req.ConfigType,
		SubType:     req.SubType,
		Description: req.Description,
		ConfigData:  req.ConfigData,
		Enabled:     enabled,
		CreatedBy:   ctx.GetString("user_id"),
	}

	if err := c.DB.Create(&config).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create config"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// UpdateConfig 更新配置文件
func (c *PolicyController) UpdateConfig(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var config models.PolicyConfig
	if err := c.DB.First(&config, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		ConfigType  string `json:"config_type"`
		SubType     string `json:"sub_type"`
		Description string `json:"description"`
		ConfigData  string `json:"config_data"`
		Enabled     *bool  `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ConfigData != "" {
		var dummy interface{}
		if err := json.Unmarshal([]byte(req.ConfigData), &dummy); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "config_data must be valid JSON"})
			return
		}
		config.ConfigData = req.ConfigData
	}

	if req.Name != "" {
		config.Name = req.Name
	}
	if req.ConfigType != "" {
		config.ConfigType = req.ConfigType
	}
	if req.SubType != "" {
		config.SubType = req.SubType
	}
	if req.Description != "" {
		config.Description = req.Description
	}
	if req.Enabled != nil {
		config.Enabled = *req.Enabled
	}
	config.UpdatedBy = ctx.GetString("user_id")
	config.Version++

	if err := c.DB.Save(&config).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// DeleteConfig 删除配置文件
func (c *PolicyController) DeleteConfig(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 检查是否被策略引用
	var count int64
	c.DB.Model(&models.Policy{}).Where("config_ids LIKE ?", "%"+strconv.Itoa(int(id))+"%").Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Config is referenced by policies, cannot delete"})
		return
	}

	if err := c.DB.Delete(&models.PolicyConfig{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete config"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ===== Policy CRUD =====

// ListPolicies 获取策略列表
func (c *PolicyController) ListPolicies(ctx *gin.Context) {
	var policies []models.Policy
	query := c.DB.Model(&models.Policy{})

	// 关键词过滤
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	// 类型过滤
	if policyType := ctx.Query("policy_type"); policyType != "" {
		query = query.Where("policy_type = ?", policyType)
	}
	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// 平台过滤
	if platform := ctx.Query("platform"); platform != "" {
		query = query.Where("platform = ? OR platform = 'all'", platform)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("priority DESC, id DESC").Offset(offset).Limit(pageSize).Find(&policies).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch policies"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      policies,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreatePolicy 创建策略
func (c *PolicyController) CreatePolicy(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		PolicyType  string `json:"policy_type" binding:"required"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		ConfigIDs   []uint `json:"config_ids"`
		RuleIDs     []uint `json:"rule_ids"`
		Enabled     *bool  `json:"enabled"`
		Status      string `json:"status"`
		Platform    string `json:"platform"`
		Scope       string `json:"scope"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	status := "draft"
	if req.Status != "" {
		status = req.Status
	}
	platform := "all"
	if req.Platform != "" {
		platform = req.Platform
	}
	scope := "all"
	if req.Scope != "" {
		scope = req.Scope
	}

	configIDsJSON, _ := json.Marshal(req.ConfigIDs)
	ruleIDsJSON, _ := json.Marshal(req.RuleIDs)

	policy := models.Policy{
		Name:        req.Name,
		PolicyType:  req.PolicyType,
		Description: req.Description,
		Priority:    req.Priority,
		ConfigIDs:   string(configIDsJSON),
		RuleIDs:     string(ruleIDsJSON),
		Enabled:     enabled,
		Status:      status,
		Platform:    platform,
		Scope:       scope,
		CreatedBy:   ctx.GetString("user_id"),
	}

	if err := c.DB.Create(&policy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create policy"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// UpdatePolicy 更新策略
func (c *PolicyController) UpdatePolicy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var policy models.Policy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		PolicyType  string `json:"policy_type"`
		Description string `json:"description"`
		Priority    *int   `json:"priority"`
		ConfigIDs   []uint `json:"config_ids"`
		RuleIDs     []uint `json:"rule_ids"`
		Enabled     *bool  `json:"enabled"`
		Status      string `json:"status"`
		Platform    string `json:"platform"`
		Scope       string `json:"scope"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		policy.Name = req.Name
	}
	if req.PolicyType != "" {
		policy.PolicyType = req.PolicyType
	}
	if req.Description != "" {
		policy.Description = req.Description
	}
	if req.Priority != nil {
		policy.Priority = *req.Priority
	}
	if req.ConfigIDs != nil {
		configIDsJSON, _ := json.Marshal(req.ConfigIDs)
		policy.ConfigIDs = string(configIDsJSON)
	}
	if req.RuleIDs != nil {
		ruleIDsJSON, _ := json.Marshal(req.RuleIDs)
		policy.RuleIDs = string(ruleIDsJSON)
	}
	if req.Enabled != nil {
		policy.Enabled = *req.Enabled
	}
	if req.Status != "" {
		policy.Status = req.Status
	}
	if req.Platform != "" {
		policy.Platform = req.Platform
	}
	if req.Scope != "" {
		policy.Scope = req.Scope
	}
	policy.UpdatedBy = ctx.GetString("user_id")

	if err := c.DB.Save(&policy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update policy"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// DeletePolicy 删除策略
func (c *PolicyController) DeletePolicy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 删除关联绑定
	c.DB.Where("policy_id = ?", id).Delete(&models.PolicyBinding{})

	if err := c.DB.Delete(&models.Policy{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete policy"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// BindPolicy 绑定策略
func (c *PolicyController) BindPolicy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 检查策略是否存在
	var policy models.Policy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req models.PolicyBindingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	boundBy := ctx.GetString("user_id")
	boundAt := time.Now()
	created := 0

	for _, targetID := range req.TargetIDs {
		// 检查是否已存在有效绑定
		var existing models.PolicyBinding
		err := c.DB.Where("policy_id = ? AND target_type = ? AND target_id = ? AND status = 1", id, req.TargetType, targetID).First(&existing).Error
		if err == nil {
			continue // 已存在，跳过
		}

		binding := models.PolicyBinding{
			PolicyID:   uint(id),
			TargetType: req.TargetType,
			TargetID:   targetID,
			BoundBy:    boundBy,
			BoundAt:    boundAt,
			Status:     1,
		}

		// 尝试获取目标名称
		binding.TargetName = targetID // 默认使用 ID

		if err := c.DB.Create(&binding).Error; err == nil {
			created++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"created": created}})
}

// UnbindPolicy 解绑策略
func (c *PolicyController) UnbindPolicy(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.PolicyUnbindRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	unboundBy := ctx.GetString("user_id")

	result := c.DB.Model(&models.PolicyBinding{}).
		Where("policy_id = ? AND target_type = ? AND target_id IN ? AND status = 1", id, req.TargetType, req.TargetIDs).
		Updates(map[string]interface{}{
			"status":     0,
			"unbound_by": unboundBy,
			"unbound_at": now,
		})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"updated": result.RowsAffected}})
}

// GetPolicyBindings 获取策略绑定列表
func (c *PolicyController) GetPolicyBindings(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var bindings []models.PolicyBinding
	query := c.DB.Where("policy_id = ? AND status = 1", id)

	if targetType := ctx.Query("target_type"); targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	if err := query.Order("bound_at DESC").Find(&bindings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bindings"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": bindings})
}

// ===== Compliance Rules =====

// ComplianceController 合规模型控制器（复用已有 CompliancePolicy/ComplianceViolation 模型）
type ComplianceController struct {
	DB *gorm.DB
}

// ListRules 获取合规规则列表
func (c *ComplianceController) ListRules(ctx *gin.Context) {
	var rules []models.CompliancePolicy
	query := c.DB.Model(&models.CompliancePolicy{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}
	if policyType := ctx.Query("policy_type"); policyType != "" {
		query = query.Where("policy_type = ?", policyType)
	}
	if enabled := ctx.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      rules,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateRule 创建合规规则
func (c *ComplianceController) CreateRule(ctx *gin.Context) {
	var req struct {
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	enforceScope := "all"
	if req.EnforceScope != "" {
		enforceScope = req.EnforceScope
	}
	severity := 2
	if req.Severity > 0 {
		severity = req.Severity
	}

	rule := models.CompliancePolicy{
		Name:              req.Name,
		Description:       req.Description,
		PolicyType:        req.PolicyType,
		TargetValue:       req.TargetValue,
		Condition:         req.Condition,
		Severity:          severity,
		RemediationAction: req.RemediationAction,
		Enabled:           enabled,
		EnforceScope:      enforceScope,
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rule"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// UpdateRule 更新合规规则
func (c *ComplianceController) UpdateRule(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var rule models.CompliancePolicy
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Rule not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req struct {
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		rule.Name = req.Name
	}
	if req.Description != "" {
		rule.Description = req.Description
	}
	if req.PolicyType != "" {
		rule.PolicyType = req.PolicyType
	}
	if req.TargetValue != "" {
		rule.TargetValue = req.TargetValue
	}
	if req.Condition != "" {
		rule.Condition = req.Condition
	}
	if req.Severity != nil {
		rule.Severity = *req.Severity
	}
	if req.RemediationAction != "" {
		rule.RemediationAction = req.RemediationAction
	}
	if req.Enabled != nil {
		rule.Enabled = *req.Enabled
	}
	if req.EnforceScope != "" {
		rule.EnforceScope = req.EnforceScope
	}

	if err := c.DB.Save(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rule"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// DeleteRule 删除合规规则
func (c *ComplianceController) DeleteRule(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.DB.Delete(&models.CompliancePolicy{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rule"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ListViolations 获取违规记录
func (c *ComplianceController) ListViolations(ctx *gin.Context) {
	var violations []models.ComplianceViolation
	query := c.DB.Model(&models.ComplianceViolation{})

	if policyID := ctx.Query("policy_id"); policyID != "" {
		query = query.Where("policy_id = ?", policyID)
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id ILIKE ?", "%"+deviceID+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if severity := ctx.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&violations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch violations"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      violations,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ResolveViolation 处理违规记录
func (c *ComplianceController) ResolveViolation(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var violation models.ComplianceViolation
	if err := c.DB.First(&violation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Violation not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"` // 3:已解决 4:已忽略
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusInt, err := strconv.Atoi(req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	now := time.Now()
	violation.Status = statusInt
	violation.ResolvedAt = &now
	violation.ResolvedBy = ctx.GetString("user_id")

	if err := c.DB.Save(&violation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update violation"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": violation})
}
