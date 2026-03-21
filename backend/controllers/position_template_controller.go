package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PositionTemplateController 基准岗位模板控制器
type PositionTemplateController struct {
	DB *gorm.DB
}

// PositionTemplateList 基准岗位模板列表
// GET /api/v1/position-templates
func (c *PositionTemplateController) PositionTemplateList(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var templates []models.PositionTemplate
	var total int64

	query := c.DB.Model(&models.PositionTemplate{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&templates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 解析 permissions JSON 字段为数组
	type TemplateResponse struct {
		models.PositionTemplate
		Permissions []string `json:"permissions"`
	}
	resp := make([]TemplateResponse, 0, len(templates))
	for _, t := range templates {
		var perms []string
		json.Unmarshal([]byte(t.Permissions), &perms)
		resp = append(resp, TemplateResponse{PositionTemplate: t, Permissions: perms})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      resp,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// PositionTemplateCreate 创建基准岗位模板
// POST /api/v1/position-templates
func (c *PositionTemplateController) PositionTemplateCreate(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var input struct {
		Name        string   `json:"name" binding:"required"`
		Code        string   `json:"code" binding:"required"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
		Status      int      `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查 code 是否重复
	var count int64
	c.DB.Model(&models.PositionTemplate{}).Where("code = ?", input.Code).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "模板编码已存在"})
		return
	}

	permissionsJSON, _ := json.Marshal(input.Permissions)
	template := models.PositionTemplate{
		TenantID:    tenantID,
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		Permissions: string(permissionsJSON),
		Status:      input.Status,
	}
	if template.Status == 0 {
		template.Status = 1
	}

	if err := c.DB.Create(&template).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	var perms []string
	json.Unmarshal([]byte(template.Permissions), &perms)
	ctx.JSON(http.StatusOK, gin.H{
		"code":        0,
		"message":     "success",
		"data":        gin.H{"id": template.ID, "name": template.Name, "code": template.Code, "description": template.Description, "permissions": perms, "status": template.Status},
	})
}

// PositionTemplateGet 获取基准岗位模板详情
// GET /api/v1/position-templates/:id
func (c *PositionTemplateController) PositionTemplateGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var template models.PositionTemplate
	if err := c.DB.First(&template, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	var perms []string
	json.Unmarshal([]byte(template.Permissions), &perms)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":          template.ID,
			"name":        template.Name,
			"code":        template.Code,
			"description": template.Description,
			"permissions": perms,
			"status":      template.Status,
			"created_at":  template.CreatedAt,
			"updated_at":  template.UpdatedAt,
		},
	})
}

// PositionTemplateUpdate 更新基准岗位模板
// PUT /api/v1/position-templates/:id
func (c *PositionTemplateController) PositionTemplateUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var template models.PositionTemplate
	if err := c.DB.First(&template, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	var input map[string]interface{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	allowedFields := []string{"name", "code", "description", "status"}
	updates := map[string]interface{}{}
	for _, field := range allowedFields {
		if val, ok := input[field]; ok {
			updates[field] = val
		}
	}

	if perms, ok := input["permissions"]; ok {
		if permsList, ok := perms.([]interface{}); ok {
			strPerms := make([]string, 0, len(permsList))
			for _, p := range permsList {
				if s, ok := p.(string); ok {
					strPerms = append(strPerms, s)
				}
			}
			permissionsJSON, _ := json.Marshal(strPerms)
			updates["permissions"] = string(permissionsJSON)
		}
	}

	if err := c.DB.Model(&template).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&template, id)
	var finalPerms []string
	json.Unmarshal([]byte(template.Permissions), &finalPerms)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":          template.ID,
			"name":        template.Name,
			"code":        template.Code,
			"description": template.Description,
			"permissions": finalPerms,
			"status":      template.Status,
			"created_at":  template.CreatedAt,
			"updated_at":  template.UpdatedAt,
		},
	})
}

// PositionTemplateDelete 删除基准岗位模板
// DELETE /api/v1/position-templates/:id
func (c *PositionTemplateController) PositionTemplateDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	var template models.PositionTemplate
	if err := c.DB.First(&template, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	if err := c.DB.Delete(&template).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PositionTemplateEnable 启用岗位模板
// POST /api/v1/position-templates/:id/enable
func (c *PositionTemplateController) PositionTemplateEnable(ctx *gin.Context) {
	id := ctx.Param("id")
	var template models.PositionTemplate
	if err := c.DB.First(&template, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	if err := c.DB.Model(&template).Update("status", 1).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "启用失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"id": template.ID, "status": 1}})
}

// PositionTemplateDisable 禁用岗位模板
// POST /api/v1/position-templates/:id/disable
func (c *PositionTemplateController) PositionTemplateDisable(ctx *gin.Context) {
	id := ctx.Param("id")
	var template models.PositionTemplate
	if err := c.DB.First(&template, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	if err := c.DB.Model(&template).Update("status", 0).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "禁用失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"id": template.ID, "status": 0}})
}

// PositionTemplateClone 复制/继承岗位模板（创建新模板，继承原模板数据）
// POST /api/v1/position-templates/:id/clone
func (c *PositionTemplateController) PositionTemplateClone(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	templateID := ctx.Param("id")

	var sourceTemplate models.PositionTemplate
	if err := c.DB.First(&sourceTemplate, templateID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
		Code string `json:"code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查新编码是否重复
	var count int64
	c.DB.Model(&models.PositionTemplate{}).Where("code = ?", input.Code).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "模板编码已存在"})
		return
	}

	newTemplate := models.PositionTemplate{
		TenantID:      tenantID,
		Name:          input.Name,
		Code:          input.Code,
		Description:   sourceTemplate.Description,
		Permissions:  sourceTemplate.Permissions,
		Status:        sourceTemplate.Status,
		InheritedFrom: &sourceTemplate.ID,
	}

	if err := c.DB.Create(&newTemplate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "复制失败"})
		return
	}

	var perms []string
	json.Unmarshal([]byte(newTemplate.Permissions), &perms)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":             newTemplate.ID,
			"name":           newTemplate.Name,
			"code":           newTemplate.Code,
			"description":    newTemplate.Description,
			"permissions":    perms,
			"status":         newTemplate.Status,
			"inherited_from": sourceTemplate.ID,
		},
	})
}

// PositionTemplateCopy 复制岗位模板到目标部门（应用模板创建岗位）
// POST /api/v1/position-templates/:id/copy
func (c *PositionTemplateController) PositionTemplateCopy(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	templateID := ctx.Param("id")

	var template models.PositionTemplate
	if err := c.DB.First(&template, templateID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "模板不存在"})
		return
	}

	var input struct {
		TargetDepartmentID uint `json:"target_department_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证目标部门存在
	var dept models.Department
	if err := c.DB.First(&dept, input.TargetDepartmentID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "目标部门不存在"})
		return
	}

	// 生成新岗位编码（模板编码 + 部门编码后缀）
	newPosCode := template.Code + "_" + dept.DeptCode

	// 创建岗位
	position := models.Position{
		TenantID:    tenantID,
		PosCode:     newPosCode,
		PosName:     template.Name,
		DeptID:      &input.TargetDepartmentID,
		CompanyID:   dept.CompanyID,
		Description: template.Description,
		Status:      template.Status,
	}

	if err := c.DB.Create(&position).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "复制失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"position_id":   position.ID,
			"position_code": position.PosCode,
			"position_name": position.PosName,
			"dept_id":       position.DeptID,
		},
	})
}
