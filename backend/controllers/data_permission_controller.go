package controllers

import (
	"net/http"
	"strings"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DataPermissionController 数据权限控制器
type DataPermissionController struct {
	DB *gorm.DB
}

// GetRoleDataPermissions 获取角色数据权限
func (c *DataPermissionController) GetRoleDataPermissions(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	roleID := ctx.Param("role_id")

	var permissions []models.UserDataPermission
	if err := c.DB.Where("role_id = ? AND tenant_id = ?", roleID, tenantID).Find(&permissions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  permissions,
			"total": len(permissions),
		},
		"message": "success",
	})
}

// UpdateRoleDataPermissions 更新角色数据权限
func (c *DataPermissionController) UpdateRoleDataPermissions(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	roleID := ctx.Param("role_id")
	userID := middleware.GetUserID(ctx)

	var req struct {
		Permissions []models.UserDataPermissionRequest `json:"permissions" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证角色存在
	var role models.Role
	if err := c.DB.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "角色不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 删除旧权限
	c.DB.Where("role_id = ? AND tenant_id = ?", roleID, tenantID).Delete(&models.UserDataPermission{})

	// 创建新权限
	for _, permReq := range req.Permissions {
		columnFields := models.StringArray(permReq.ColumnFields)
		dataScope := models.JSONMap(permReq.DataScope)
		if permReq.DataScope == nil {
			dataScope = models.JSONMap{}
		}

		perm := models.UserDataPermission{
			UserID:       0, // 角色级别为0
			RoleID:       uint(parseIntStr(roleID, 10)),
			ResourceType: permReq.ResourceType,
			RuleType:     permReq.RuleType,
			ColumnFields: columnFields,
			DataScope:    dataScope,
			FilterExpr:   permReq.FilterExpr,
			Priority:     permReq.Priority,
			IsActive:     permReq.IsActive,
			TenantID:     tenantID,
			CreatedBy:    userID,
		}
		if err := c.DB.Create(&perm).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建权限失败: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// GetUserDataPermissions 获取用户数据权限
func (c *DataPermissionController) GetUserDataPermissions(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	userID := ctx.Param("user_id")

	var permissions []models.UserDataPermission
	if err := c.DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&permissions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  permissions,
			"total": len(permissions),
		},
		"message": "success",
	})
}

// UpdateUserDataPermissions 更新用户数据权限
func (c *DataPermissionController) UpdateUserDataPermissions(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	userIDParam := ctx.Param("user_id")
	userID := middleware.GetUserID(ctx)

	userIDNum := parseIntStr(userIDParam, 10)

	var req struct {
		Permissions []models.UserDataPermissionRequest `json:"permissions" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 删除旧权限
	c.DB.Where("user_id = ? AND tenant_id = ?", userIDNum, tenantID).Delete(&models.UserDataPermission{})

	// 创建新权限
	for _, permReq := range req.Permissions {
		columnFields := models.StringArray(permReq.ColumnFields)
		dataScope := models.JSONMap(permReq.DataScope)
		if permReq.DataScope == nil {
			dataScope = models.JSONMap{}
		}

		perm := models.UserDataPermission{
			UserID:       uint(userIDNum),
			RoleID:       0, // 用户级别为0
			ResourceType: permReq.ResourceType,
			RuleType:     permReq.RuleType,
			ColumnFields: columnFields,
			DataScope:    dataScope,
			FilterExpr:   permReq.FilterExpr,
			Priority:     permReq.Priority,
			IsActive:     permReq.IsActive,
			TenantID:     tenantID,
			CreatedBy:    userID,
		}
		if err := c.DB.Create(&perm).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建权限失败: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// GetColumnPermissions 获取可配置列级权限字段列表
func (c *DataPermissionController) GetColumnPermissions(ctx *gin.Context) {
	// 定义各资源类型的可配置列级权限字段
	columns := []models.ColumnPermission{
		// 设备相关
		{Field: "device_id", Label: "设备ID", TableName: "devices", DataType: "string", Sensitive: true},
		{Field: "mac_address", Label: "MAC地址", TableName: "devices", DataType: "string", Sensitive: true},
		{Field: "sn_code", Label: "序列号", TableName: "devices", DataType: "string", Sensitive: true},
		{Field: "firmware_version", Label: "固件版本", TableName: "devices", DataType: "string", Sensitive: false},
		{Field: "lifecycle_status", Label: "生命周期状态", TableName: "devices", DataType: "int", Sensitive: false},
		{Field: "created_at", Label: "创建时间", TableName: "devices", DataType: "datetime", Sensitive: false},
		// 会员相关
		{Field: "phone", Label: "手机号", TableName: "members", DataType: "string", Sensitive: true},
		{Field: "email", Label: "邮箱", TableName: "members", DataType: "string", Sensitive: true},
		{Field: "balance", Label: "余额", TableName: "members", DataType: "decimal", Sensitive: true},
		{Field: "points", Label: "积分", TableName: "members", DataType: "int", Sensitive: false},
		// 组织相关
		{Field: "company_name", Label: "公司名称", TableName: "companies", DataType: "string", Sensitive: false},
		{Field: "legal_person", Label: "法人", TableName: "companies", DataType: "string", Sensitive: true},
		{Field: "license_no", Label: "营业执照号", TableName: "companies", DataType: "string", Sensitive: true},
	}

	resourceType := ctx.Query("resource_type")
	if resourceType != "" {
		// 过滤指定资源类型的字段
		filtered := make([]models.ColumnPermission, 0)
		for _, col := range columns {
			if strings.HasPrefix(col.TableName, resourceType) {
				filtered = append(filtered, col)
			}
		}
		columns = filtered
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  columns,
			"total": len(columns),
		},
		"message": "success",
	})
}

// ValidatePermissionExpression 验证权限表达式
func (c *DataPermissionController) ValidatePermissionExpression(ctx *gin.Context) {
	var req models.DataPermissionValidateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	response := models.DataPermissionValidateResponse{
		Valid:   true,
		Result:  true,
		Details: []string{},
	}

	// 简单表达式验证
	// 实际生产应使用表达式引擎（如 goja、gval 等）
	expression := strings.TrimSpace(req.Expression)

	if expression == "" {
		response.Valid = false
		response.Error = "表达式不能为空"
	} else if strings.Contains(expression, "DROP") || strings.Contains(expression, "DELETE") || strings.Contains(expression, "TRUNCATE") {
		response.Valid = false
		response.Error = "表达式包含禁止的关键字"
	} else if strings.Contains(expression, "rm -rf") || strings.Contains(expression, "format") {
		response.Valid = false
		response.Error = "表达式包含危险命令"
	} else {
		// 检查语法（简单括号匹配）
		openCount := 0
		for _, ch := range expression {
			if ch == '(' {
				openCount++
			} else if ch == ')' {
				openCount--
			}
		}
		if openCount != 0 {
			response.Valid = false
			response.Error = "括号不匹配"
		} else {
			response.Details = append(response.Details, "语法检查通过")
			response.Details = append(response.Details, "表达式解析成功")
		}
	}

	if !response.Valid {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    400,
			"data":    response,
			"message": response.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    response,
		"message": "验证通过",
	})
}

// ListDataPermissionRules 列出数据权限规则
func (c *DataPermissionController) ListDataPermissionRules(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	page := parseIntDefault(ctx.Query("page"), 1)
	pageSize := parseIntDefault(ctx.Query("page_size"), 20)
	resourceType := ctx.Query("resource_type")
	ruleType := ctx.Query("rule_type")
	isActive := ctx.Query("is_active")

	query := c.DB.Model(&models.DataPermissionRule{}).Where("tenant_id = ?", tenantID)

	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if ruleType != "" {
		query = query.Where("rule_type = ?", ruleType)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var total int64
	query.Count(&total)

	var rules []models.DataPermissionRule
	offset := (page - 1) * pageSize
	if err := query.Order("priority DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      rules,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// CreateDataPermissionRule 创建数据权限规则
func (c *DataPermissionController) CreateDataPermissionRule(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	userID := middleware.GetUserID(ctx)

	var req models.DataPermissionRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if req.RuleType != "row" && req.RuleType != "column" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "规则类型必须是 row 或 column",
		})
		return
	}

	permExpr := models.JSONMap(req.PermissionExpr)
	if req.PermissionExpr == nil {
		permExpr = models.JSONMap{}
	}

	rule := models.DataPermissionRule{
		RuleName:       req.RuleName,
		ResourceType:   req.ResourceType,
		RuleType:       req.RuleType,
		ResourceIDs:    req.ResourceIDs,
		PermissionExpr: permExpr,
		Priority:       req.Priority,
		IsActive:       req.IsActive,
		Description:    req.Description,
		TenantID:       tenantID,
		CreatedBy:      userID,
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建规则失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rule,
		"message": "创建成功",
	})
}

// UpdateDataPermissionRule 更新数据权限规则
func (c *DataPermissionController) UpdateDataPermissionRule(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var rule models.DataPermissionRule
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "规则不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	var req models.DataPermissionRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	rule.RuleName = req.RuleName
	rule.ResourceType = req.ResourceType
	rule.RuleType = req.RuleType
	rule.ResourceIDs = req.ResourceIDs
	if req.PermissionExpr != nil {
		rule.PermissionExpr = models.JSONMap(req.PermissionExpr)
	}
	rule.Priority = req.Priority
	rule.IsActive = req.IsActive
	rule.Description = req.Description

	if err := c.DB.Save(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    rule,
		"message": "更新成功",
	})
}

// DeleteDataPermissionRule 删除数据权限规则
func (c *DataPermissionController) DeleteDataPermissionRule(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var rule models.DataPermissionRule
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "规则不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	if err := c.DB.Delete(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// parseIntStr 简单字符串转int
func parseIntStr(s string, base int) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*base + int(c-'0')
		}
	}
	return n
}
