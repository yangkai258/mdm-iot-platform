package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertRuleController 告警规则控制器
type AlertRuleController struct {
	DB *gorm.DB
}

// RegisterAlertRuleRoutes 注册告警规则路由
func (ctrl *AlertRuleController) RegisterAlertRuleRoutes(api *gin.RouterGroup) {
	api.GET("/alert-rules", ctrl.List)
	api.POST("/alert-rules", ctrl.Create)
	api.GET("/alert-rules/:id", ctrl.Get)
	api.PUT("/alert-rules/:id", ctrl.Update)
	api.DELETE("/alert-rules/:id", ctrl.Delete)
	api.POST("/alert-rules/:id/test", ctrl.Test)
	api.PUT("/alert-rules/:id/toggle", ctrl.Toggle)
}

// List 获取告警规则列表
// GET /api/v1/alert-rules
func (ctrl *AlertRuleController) List(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	ruleType := c.Query("rule_type")
	enabled := c.Query("enabled")
	deviceID := c.Query("device_id")
	sortBy := c.DefaultQuery("sort_by", "priority")
	order := c.DefaultQuery("order", "desc")

	query := ctrl.DB.Model(&models.AlertRule{})

	if ruleType != "" {
		query = query.Where("rule_type = ?", ruleType)
	}
	if enabled != "" {
		if enabled == "true" {
			query = query.Where("enabled = ?", true)
		} else if enabled == "false" {
			query = query.Where("enabled = ?", false)
		}
	}
	if deviceID != "" {
		query = query.Where("device_id = ? OR device_id = ''", deviceID)
	}

	orderMap := map[string]string{"asc": "asc", "desc": "desc"}
	if orderMap[order] == "" {
		order = "desc"
	}
	query = query.Order(sortBy + " " + order)

	var total int64
	query.Count(&total)

	var rules []models.AlertRule
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to query alert rules: " + err.Error(),
		})
		return
	}

	responses := make([]*models.AlertRuleResponse, len(rules))
	for i := range rules {
		responses[i] = rules[i].ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      responses,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Get 获取告警规则详情
// GET /api/v1/alert-rules/:id
func (ctrl *AlertRuleController) Get(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Alert rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rule.ToResponse(),
	})
}

// Create 创建告警规则
// POST /api/v1/alert-rules
func (ctrl *AlertRuleController) Create(c *gin.Context) {
	var req models.AlertRuleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	priority := 50
	if req.Priority > 0 {
		priority = req.Priority
	}

	conditions, _ := json.Marshal(req.Conditions)
	actions, _ := json.Marshal(req.Actions)

	rule := models.AlertRule{
		RuleName:   req.RuleName,
		RuleType:   req.RuleType,
		Conditions: models.JSON{},
		Actions:    models.JSON{},
		Enabled:    enabled,
		Priority:   priority,
		DeviceID:   req.DeviceID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	json.Unmarshal(conditions, &rule.Conditions)
	json.Unmarshal(actions, &rule.Actions)

	if err := ctrl.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create alert rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": rule.ToResponse(),
	})
}

// Update 更新告警规则
// PUT /api/v1/alert-rules/:id
func (ctrl *AlertRuleController) Update(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Alert rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	var req models.AlertRuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request: " + err.Error()})
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.RuleName != "" {
		updates["rule_name"] = req.RuleName
	}
	if req.RuleType != "" {
		updates["rule_type"] = req.RuleType
	}
	if req.Conditions != nil {
		conds, _ := json.Marshal(req.Conditions)
		var cJSON models.JSON
		json.Unmarshal(conds, &cJSON)
		updates["conditions"] = cJSON
	}
	if req.Actions != nil {
		acts, _ := json.Marshal(req.Actions)
		var aJSON models.JSON
		json.Unmarshal(acts, &aJSON)
		updates["actions"] = aJSON
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	if req.DeviceID != "" {
		updates["device_id"] = req.DeviceID
	}

	if err := ctrl.DB.Model(&rule).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update alert rule: " + err.Error()})
		return
	}

	ctrl.DB.First(&rule, rule.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": rule.ToResponse(),
	})
}

// Delete 删除告警规则
// DELETE /api/v1/alert-rules/:id
func (ctrl *AlertRuleController) Delete(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Alert rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if err := ctrl.DB.Delete(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete alert rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Alert rule deleted successfully",
	})
}

// Test 测试告警规则
// POST /api/v1/alert-rules/:id/test
func (ctrl *AlertRuleController) Test(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Alert rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 获取测试数据
	var testData struct {
		Data map[string]interface{} `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&testData); err != nil {
		// 如果没有提供测试数据，使用默认测试数据
		testData.Data = map[string]interface{}{
			"temperature": 35.0,
			"battery":     20.0,
			"humidity":    80.0,
		}
	}

	// 模拟规则匹配
	matched, err := evaluateAlertRule(rule.Conditions, testData.Data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"data":    gin.H{"matched": false, "error": err.Error()},
			"message": "Rule evaluation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"matched":      matched,
			"test_data":    testData.Data,
			"conditions":   rule.Conditions,
			"actions":      rule.Actions,
			"message":      "Rule test completed",
		},
	})
}

// Toggle 启用/禁用告警规则
// PUT /api/v1/alert-rules/:id/toggle
func (ctrl *AlertRuleController) Toggle(c *gin.Context) {
	id := c.Param("id")

	var rule models.AlertRule
	if err := ctrl.DB.Where("id = ? OR rule_id = ?", id, id).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Alert rule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	rule.Enabled = !rule.Enabled
	rule.UpdatedAt = time.Now()

	if err := ctrl.DB.Save(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to toggle alert rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"enabled":  rule.Enabled,
			"rule_id":  rule.RuleID,
			"message":  "Alert rule toggled successfully",
		},
	})
}

// evaluateAlertRule 评估告警规则条件是否匹配
func evaluateAlertRule(conditions models.JSON, data map[string]interface{}) (bool, error) {
	if conditions == nil || len(conditions) == 0 {
		return false, nil
	}

	// 支持数组格式的条件列表
	items, ok := conditions["items"].([]interface{})
	if !ok {
		// 尝试直接作为数组解析
		items = make([]interface{}, 0)
		for _, v := range conditions {
			if vmap, ok := v.(map[string]interface{}); ok {
				items = append(items, vmap)
			}
		}
	}

	if len(items) == 0 {
		// 单条件格式：{"field": "temperature", "operator": ">", "value": 30}
		field, ok := conditions["field"].(string)
		if !ok {
			return false, nil
		}
		op, _ := conditions["operator"].(string)
		threshold, _ := toFloat64(conditions["value"])

		return alertRuleEvaluateCondition(field, op, threshold, data), nil
	}

	// 多条件 AND 逻辑
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		field, ok := itemMap["field"].(string)
		if !ok {
			continue
		}
		op, _ := itemMap["operator"].(string)
		threshold, _ := toFloat64(itemMap["value"])
		if !alertRuleEvaluateCondition(field, op, threshold, data) {
			return false, nil
		}
	}

	return true, nil
}

// evaluateCondition 评估单个条件
func alertRuleEvaluateCondition(field, operator string, threshold float64, data map[string]interface{}) bool {
	val, exists := data[field]
	if !exists {
		return false
	}

	var numVal float64
	switch v := val.(type) {
	case float64:
		numVal = v
	case int:
		numVal = float64(v)
	case int64:
		numVal = float64(v)
	case string:
		// 尝试解析为数字
		var parsed float64
		if _, err := json.Marshal(v); err == nil {
			// 简单解析
			for _, c := range v {
				if c >= '0' && c <= '9' || c == '.' {
					continue
				}
				return false
				}
			}
		_ = parsed
	}

	switch operator {
	case ">":
		return numVal > threshold
	case "<":
		return numVal < threshold
	case ">=":
		return numVal >= threshold
	case "<=":
		return numVal <= threshold
	case "=":
		return numVal == threshold
	case "!=":
		return numVal != threshold
	default:
		return false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		var f float64
		for _, c := range val {
			if c == '.' {
				continue
			}
			if c < '0' || c > '9' {
				return 0, false
			}
		}
		// simple parse
		result := 0.0
		dot := false
		div := 1.0
		for _, c := range val {
			if c == '.' {
				dot = true
				continue
			}
			if c < '0' || c > '9' {
				return 0, false
			}
			result = result*10 + float64(c-'0')
			if dot {
				div *= 10
			}
		}
		return result / div * result / result, true
	}
	return 0, false
}
