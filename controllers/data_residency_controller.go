package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DataResidencyController 数据驻留控制器
type DataResidencyController struct {
	DB *gorm.DB
}

// NewDataResidencyController 创建数据驻留控制器
func NewDataResidencyController(db *gorm.DB) *DataResidencyController {
	return &DataResidencyController{DB: db}
}

// ListDataResidencyRules 获取数据驻留规则列表
// @Summary 获取数据驻留规则列表
// @Tags data-residency
// @Produce json
// @Success 200 {array} models.DataResidencyRule
// @Router /api/v1/data-residency/rules [get]
func (ctrl *DataResidencyController) ListDataResidencyRules(c *gin.Context) {
	var rules []models.DataResidencyRule
	
	query := ctrl.DB.Model(&models.DataResidencyRule{})
	
	// 按租户过滤
	tenantID := c.GetUint("tenant_id")
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 按数据类型过滤
	if dataType := c.Query("data_type"); dataType != "" {
		query = query.Where("data_type = ?", dataType)
	}

	// 按状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rules,
	})
}

// GetDataResidencyRule 获取数据驻留规则详情
// @Summary 获取数据驻留规则详情
// @Tags data-residency
// @Produce json
// @Param id path int true "规则ID"
// @Success 200 {object} models.DataResidencyRule
// @Router /api/v1/data-residency/rules/{id} [get]
func (ctrl *DataResidencyController) GetDataResidencyRule(c *gin.Context) {
	id := c.Param("id")
	var ruleID uint
	if _, err := parseUint(id, &ruleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var rule models.DataResidencyRule
	if err := ctrl.DB.First(&rule, ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 租户权限检查
	tenantID := c.GetUint("tenant_id")
	if tenantID > 0 && rule.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rule,
	})
}

// CreateDataResidencyRule 创建数据驻留规则
// @Summary 创建数据驻留规则
// @Tags data-residency
// @Accept json
// @Produce json
// @Param rule body models.DataResidencyRule true "规则信息"
// @Success 201 {object} models.DataResidencyRule
// @Router /api/v1/data-residency/rules [post]
func (ctrl *DataResidencyController) CreateDataResidencyRule(c *gin.Context) {
	var rule models.DataResidencyRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成规则ID
	rule.RuleID = uuid.New().String()

	// 设置租户
	tenantID := c.GetUint("tenant_id")
	if tenantID > 0 {
		rule.TenantID = tenantID
	}

	// 设置默认状态
	if rule.Status == "" {
		rule.Status = "active"
	}

	// 检查区域是否存在
	var region models.Region
	if err := ctrl.DB.Where("region_code = ?", rule.RegionCode).First(&region).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "region not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": rule,
	})
}

// UpdateDataResidencyRule 更新数据驻留规则
// @Summary 更新数据驻留规则
// @Tags data-residency
// @Accept json
// @Produce json
// @Param id path int true "规则ID"
// @Param rule body models.DataResidencyRule true "规则信息"
// @Success 200 {object} models.DataResidencyRule
// @Router /api/v1/data-residency/rules/{id} [put]
func (ctrl *DataResidencyController) UpdateDataResidencyRule(c *gin.Context) {
	id := c.Param("id")
	var ruleID uint
	if _, err := parseUint(id, &ruleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var existingRule models.DataResidencyRule
	if err := ctrl.DB.First(&existingRule, ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 租户权限检查
	tenantID := c.GetUint("tenant_id")
	if tenantID > 0 && existingRule.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 删除不允许更新的字段
	delete(updates, "id")
	delete(updates, "rule_id")
	delete(updates, "tenant_id")
	delete(updates, "created_at")

	if err := ctrl.DB.Model(&existingRule).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.First(&existingRule, ruleID)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": existingRule,
	})
}

// DeleteDataResidencyRule 删除数据驻留规则
// @Summary 删除数据驻留规则
// @Tags data-residency
// @Produce json
// @Param id path int true "规则ID"
// @Success 200 {object} gin.H
// @Router /api/v1/data-residency/rules/{id} [delete]
func (ctrl *DataResidencyController) DeleteDataResidencyRule(c *gin.Context) {
	id := c.Param("id")
	var ruleID uint
	if _, err := parseUint(id, &ruleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}

	var rule models.DataResidencyRule
	if err := ctrl.DB.First(&rule, ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 租户权限检查
	tenantID := c.GetUint("tenant_id")
	if tenantID > 0 && rule.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := ctrl.DB.Delete(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "rule deleted",
	})
}

// ValidateDataResidency 验证数据驻留合规性
// @Summary 验证数据驻留合规性
// @Tags data-residency
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "待验证数据"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/data-residency/rules/validate [post]
func (ctrl *DataResidencyController) ValidateDataResidency(c *gin.Context) {
	var req struct {
		DataType  string `json:"data_type" binding:"required"`
		RegionCode string `json:"region_code" binding:"required"`
		DataID    uint   `json:"data_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")

	// 查询匹配的规则
	var rule models.DataResidencyRule
	err := ctrl.DB.Where("data_type = ? AND tenant_id = ? AND status = ?", 
		req.DataType, tenantID, "active").First(&rule).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": map[string]interface{}{
					"compliant":  true,
					"message":    "no matching rule found, data is compliant by default",
					"matched_rule": nil,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查数据是否在正确的区域
	compliant := rule.RegionCode == req.RegionCode

	result := map[string]interface{}{
		"compliant":   compliant,
		"data_type":   req.DataType,
		"data_region": req.RegionCode,
		"required_region": rule.RegionCode,
		"matched_rule": rule.RuleID,
	}

	if compliant {
		result["message"] = "data is stored in the correct region"
	} else {
		result["message"] = "data residency violation: data must be stored in " + rule.RegionCode
		result["violation"] = map[string]interface{}{
			"type":        "data_residency",
			"expected":    rule.RegionCode,
			"actual":      req.RegionCode,
			"description": rule.Description,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// BatchCreateDataResidencyRules 批量创建数据驻留规则
// @Summary 批量创建数据驻留规则
// @Tags data-residency
// @Accept json
// @Produce json
// @Param rules body []models.DataResidencyRule true "规则列表"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/data-residency/rules/batch [post]
func (ctrl *DataResidencyController) BatchCreateDataResidencyRules(c *gin.Context) {
	var rules []models.DataResidencyRule
	if err := c.ShouldBindJSON(&rules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetUint("tenant_id")
	now := time.Now()

	for i := range rules {
		rules[i].RuleID = uuid.New().String()
		rules[i].TenantID = tenantID
		rules[i].CreatedAt = now
		rules[i].UpdatedAt = now
		if rules[i].Status == "" {
			rules[i].Status = "active"
		}
	}

	if err := ctrl.DB.Create(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"created": len(rules),
			"rules":   rules,
		},
	})
}
