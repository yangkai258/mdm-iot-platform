package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
)

// EmailTemplateController 邮件模板控制器
type EmailTemplateController struct {
	Service *services.EmailTemplateService
}

// NewEmailTemplateController 创建邮件模板控制器
func NewEmailTemplateController(svc *services.EmailTemplateService) *EmailTemplateController {
	return &EmailTemplateController{Service: svc}
}

// List 获取邮件模板列表
func (c *EmailTemplateController) List(ctx *gin.Context) {
	tenantID := ctx.Query("tenant_id")
	keyword := ctx.Query("keyword")
	statusStr := ctx.Query("status")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	var status *int
	if statusStr != "" {
		s := 1
		if statusStr == "0" {
			s = 0
		}
		status = &s
	}

	list, total, err := c.Service.List(tenantID, keyword, status, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": list,
			"pagination": gin.H{
				"total":    total,
				"current":  page,
				"pageSize": pageSize,
			},
		},
	})
}

// Get 获取单个邮件模板
func (c *EmailTemplateController) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	tpl, err := c.Service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tpl,
	})
}

// Create 创建邮件模板
func (c *EmailTemplateController) Create(ctx *gin.Context) {
	var tpl models.EmailTemplate
	if err := ctx.ShouldBindJSON(&tpl); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if tpl.Code == "" || tpl.Name == "" || tpl.Subject == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "code、name、subject 为必填字段",
		})
		return
	}

	// 提取 tenant_id（如果系统使用多租户）
	if tenantID, exists := ctx.Get("tenant_id"); exists {
		tpl.TenantID = tenantID.(string)
	}

	if err := c.Service.Create(&tpl); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    tpl,
	})
}

// Update 更新邮件模板
func (c *EmailTemplateController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := c.Service.Update(uint(id), updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// Delete 删除邮件模板（软删除）
func (c *EmailTemplateController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}

	if err := c.Service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// RenderPreview 预览渲染结果（用于模板编辑时的实时预览）
func (c *EmailTemplateController) RenderPreview(ctx *gin.Context) {
	var req struct {
		Code string                 `json:"code" binding:"required"`
		Vars map[string]interface{} `json:"vars"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	tpl, err := c.Service.GetByCode(req.Code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	subject, body, err := c.Service.Render(tpl, req.Vars)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"subject": subject,
			"body":    body,
		},
	})
}
