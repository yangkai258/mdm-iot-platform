package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PolicyController 策略配置控制器
type PolicyController struct {
	DB *gorm.DB
}

// ListConfigs 获取策略配置列表
func (c *PolicyController) ListConfigs(ctx *gin.Context) {
	var configs []models.PolicyConfig
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	configType := ctx.Query("config_type")
	enabled := ctx.Query("enabled")
	keyword := ctx.Query("keyword")

	query := c.DB.Model(&models.PolicyConfig{})

	if configType != "" {
		query = query.Where("config_type = ?", configType)
	}
	if enabled != "" {
		query = query.Where("enabled = ?", enabled)
	}
	if keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&configs)

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

// GetConfig 获取单个策略配置
func (c *PolicyController) GetConfig(ctx *gin.Context) {
	id := ctx.Param("id")

	var config models.PolicyConfig
	if err := c.DB.First(&config, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "策略配置不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// CreateConfig 创建策略配置
func (c *PolicyController) CreateConfig(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		ConfigType  string `json:"config_type" binding:"required"`
		SubType     string `json:"sub_type"`
		Description string `json:"description"`
		ConfigData  string `json:"config_data"`
		Enabled     *bool  `json:"enabled"`
		Platform    string `json:"platform"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	username, _ := ctx.Get("username")
	uname, _ := username.(string)

	config := models.PolicyConfig{
		Name:        req.Name,
		ConfigType:  req.ConfigType,
		SubType:     req.SubType,
		Description: req.Description,
		ConfigData:  req.ConfigData,
		Enabled:     enabled,
		Version:     1,
		CreatedBy:  uname,
		UpdatedBy:  uname,
	}

	if err := c.DB.Create(&config).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败: " + err.Error(),
		})
		return
	}

	// record audit
	c.recordPolicyAudit(ctx, "create", "policy_config", config.ID, uid, "创建策略配置")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// UpdateConfig 更新策略配置
func (c *PolicyController) UpdateConfig(ctx *gin.Context) {
	id := ctx.Param("id")

	var config models.PolicyConfig
	if err := c.DB.First(&config, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "策略配置不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	username, _ := ctx.Get("username")
	uname, _ := username.(string)

	updates := map[string]interface{}{
		"version":     config.Version + 1,
		"updated_by": uname,
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ConfigType != "" {
		updates["config_type"] = req.ConfigType
	}
	if req.SubType != "" {
		updates["sub_type"] = req.SubType
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ConfigData != "" {
		updates["config_data"] = req.ConfigData
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := c.DB.Model(&config).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	c.DB.First(&config, id)
	c.recordPolicyAudit(ctx, "update", "policy_config", config.ID, 0, "更新策略配置")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// DeleteConfig 删除策略配置
func (c *PolicyController) DeleteConfig(ctx *gin.Context) {
	id := ctx.Param("id")

	var config models.PolicyConfig
	if err := c.DB.First(&config, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "策略配置不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	if err := c.DB.Delete(&config).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	c.recordPolicyAudit(ctx, "delete", "policy_config", config.ID, uid, "删除策略配置")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ============ Policy 策略主表 CRUD ============

// ListPolicies 获取策略列表
func (c *PolicyController) ListPolicies(ctx *gin.Context) {
	var policies []models.Policy
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	policyType := ctx.Query("policy_type")
	status := ctx.Query("status")
	keyword := ctx.Query("keyword")

	query := c.DB.Model(&models.Policy{})

	if policyType != "" {
		query = query.Where("policy_type = ?", policyType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name ILIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("priority DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&policies)

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

// GetPolicy 获取单个策略
func (c *PolicyController) GetPolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.Policy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// CreatePolicy 创建策略
func (c *PolicyController) CreatePolicy(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		PolicyType  string `json:"policy_type" binding:"required"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		ConfigIDs   string `json:"config_ids"`
		RuleIDs     string `json:"rule_ids"`
		Enabled     *bool  `json:"enabled"`
		Platform    string `json:"platform"`
		Scope       string `json:"scope"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	username, _ := ctx.Get("username")
	uname, _ := username.(string)

	policy := models.Policy{
		Name:        req.Name,
		PolicyType:  req.PolicyType,
		Description: req.Description,
		Priority:    req.Priority,
		ConfigIDs:   req.ConfigIDs,
		RuleIDs:     req.RuleIDs,
		Enabled:     enabled,
		Status:      "active",
		Platform:    req.Platform,
		Scope:       req.Scope,
		CreatedBy:   uname,
		UpdatedBy:   uname,
	}

	if err := c.DB.Create(&policy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// UpdatePolicy 更新策略
func (c *PolicyController) UpdatePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.Policy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		PolicyType  string `json:"policy_type"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		ConfigIDs   string `json:"config_ids"`
		RuleIDs     string `json:"rule_ids"`
		Enabled     *bool  `json:"enabled"`
		Status      string `json:"status"`
		Platform    string `json:"platform"`
		Scope       string `json:"scope"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	username, _ := ctx.Get("username")
	uname, _ := username.(string)
	updates := map[string]interface{}{"updated_by": uname}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.PolicyType != "" {
		updates["policy_type"] = req.PolicyType
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	if req.ConfigIDs != "" {
		updates["config_ids"] = req.ConfigIDs
	}
	if req.RuleIDs != "" {
		updates["rule_ids"] = req.RuleIDs
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Platform != "" {
		updates["platform"] = req.Platform
	}
	if req.Scope != "" {
		updates["scope"] = req.Scope
	}

	if err := c.DB.Model(&policy).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败: " + err.Error()})
		return
	}

	c.DB.First(&policy, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// DeletePolicy 删除策略
func (c *PolicyController) DeletePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.Policy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.DB.Delete(&policy)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PolicyBinding 策略绑定 ============

// BindPolicy 绑定策略到目标
func (c *PolicyController) BindPolicy(ctx *gin.Context) {
	policyID := ctx.Param("id")

	var policy models.Policy
	if err := c.DB.First(&policy, policyID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "策略不存在"})
		return
	}

	var req models.PolicyBindingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	username, _ := ctx.Get("username")
	uname, _ := username.(string)

	for _, targetID := range req.TargetIDs {
		binding := models.PolicyBinding{
			PolicyID:   policy.ID,
			TargetType: req.TargetType,
			TargetID:   targetID,
			BoundBy:    uname,
		}
		c.DB.Create(&binding)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UnbindPolicy 解绑策略
func (c *PolicyController) UnbindPolicy(ctx *gin.Context) {
	policyID := ctx.Param("id")

	var req models.PolicyUnbindRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	username, _ := ctx.Get("username")
	uname, _ := username.(string)

	now := c.DB.NowFunc()
	c.DB.Model(&models.PolicyBinding{}).
		Where("policy_id = ? AND target_type = ? AND target_id IN ?", policyID, req.TargetType, req.TargetIDs).
		Updates(map[string]interface{}{"status": 0, "unbound_by": uname, "unbound_at": now})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetPolicyBindings 获取策略绑定列表
func (c *PolicyController) GetPolicyBindings(ctx *gin.Context) {
	policyID := ctx.Param("id")

	var bindings []models.PolicyBinding
	c.DB.Where("policy_id = ? AND status = 1", policyID).Find(&bindings)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": bindings})
}

// ============ 辅助函数 ============

func (c *PolicyController) recordPolicyAudit(ctx *gin.Context, action, resourceType string, resourceID uint, userID uint, remark string) {
	username, _ := ctx.Get("username")
	log := models.AuditLog{
		Action:        action,
		Module:        "policy",
		ResourceType:  resourceType,
		ResourceID:    strconv.FormatUint(uint64(resourceID), 10),
		UserID:        userID,
		Username:      username.(string),
		IP:            ctx.ClientIP(),
		UserAgent:     ctx.GetHeader("User-Agent"),
		Status:        1,
		RequestMethod: ctx.Request.Method,
		RequestPath:   ctx.Request.URL.Path,
	}
	c.DB.Create(&log)
}
