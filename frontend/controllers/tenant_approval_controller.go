package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TenantApprovalController 租户申请审批控制器
type TenantApprovalController struct {
	DB *gorm.DB
}

// RegisterTenantApprovalRoutes 注册租户申请审批路由
func (tc *TenantApprovalController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tenant-approvals", tc.ListApplications)
	rg.GET("/tenant-approvals/:id", tc.GetApplication)
	rg.POST("/tenant-approvals/:id/approve", tc.ApproveApplication)
	rg.POST("/tenant-approvals/:id/reject", tc.RejectApplication)
	rg.GET("/tenant-approvals/:id/history", tc.GetApprovalHistory)
	rg.DELETE("/tenant-approvals/:id", tc.DeleteApplication)
	rg.POST("/tenant-approvals/batch-delete", tc.BatchDeleteApplications)
}

// ListApplications 获取租户申请列表
func (tc *TenantApprovalController) ListApplications(c *gin.Context) {
	var applications []models.TenantApplication
	query := tc.DB.Model(&models.TenantApplication{})

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("company_name ILIKE ? OR contact_name ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applications)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":  applications,
			"total": total,
			"page":  page,
		},
	})
}

// GetApplication 获取申请详情
func (tc *TenantApprovalController) GetApplication(c *gin.Context) {
	id := c.Param("id")

	var app models.TenantApplication
	if err := tc.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "申请不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data":    app,
	})
}

// ApproveApplication 审批通过
func (tc *TenantApprovalController) ApproveApplication(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetString("user_id")

	var app models.TenantApplication
	if err := tc.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "申请不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if app.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只能审批待处理的申请"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// 更新申请状态
	app.Status = "approved"
	app.ApprovedBy = operator
	app.ApprovedAt = now

	if err := tc.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审批失败"})
		return
	}

	// 记录审批历史
	history := models.ApprovalHistory{
		ApplicationID: app.ID,
		Action:        "approve",
		ActionText:    "审批通过",
		Operator:      operator,
		CreatedAt:     now,
	}
	tc.DB.Create(&history)

	// 创建正式租户
	tenantID := fmt.Sprintf("%d", app.ID)
	tenant := models.Tenant{
		ID:           tenantID,
		TenantCode:   app.ApplicationCode,
		Name:         app.CompanyName,
		ContactName:  app.ContactName,
		ContactPhone: app.ContactPhone,
		ContactEmail: app.ContactEmail,
		Plan:         app.PlanName,
		Status:       "active",
	}
	tc.DB.Create(&tenant)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "审批成功",
		"data":    app,
	})
}

// RejectApplication 审批拒绝
func (tc *TenantApprovalController) RejectApplication(c *gin.Context) {
	id := c.Param("id")
	operator := c.GetString("user_id")

	var req struct {
		RejectReason string `json:"reject_reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供拒绝原因"})
		return
	}

	var app models.TenantApplication
	if err := tc.DB.First(&app, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "申请不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if app.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只能审批待处理的申请"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	app.Status = "rejected"
	app.RejectReason = req.RejectReason
	app.ApprovedBy = operator
	app.ApprovedAt = now

	if err := tc.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审批失败"})
		return
	}

	// 记录审批历史
	history := models.ApprovalHistory{
		ApplicationID: app.ID,
		Action:        "reject",
		ActionText:    "审批拒绝",
		Operator:      operator,
		Comment:       req.RejectReason,
		CreatedAt:     now,
	}
	tc.DB.Create(&history)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已拒绝",
		"data":    app,
	})
}

// GetApprovalHistory 获取审批历史
func (tc *TenantApprovalController) GetApprovalHistory(c *gin.Context) {
	id := c.Param("id")

	var history []models.ApprovalHistory
	tc.DB.Where("application_id = ?", id).Order("created_at DESC").Find(&history)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    history,
	})
}

// DeleteApplication 删除申请
func (tc *TenantApprovalController) DeleteApplication(c *gin.Context) {
	id := c.Param("id")

	if err := tc.DB.Delete(&models.TenantApplication{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// BatchDeleteApplications 批量删除申请
func (tc *TenantApprovalController) BatchDeleteApplications(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供要删除的ID列表"})
		return
	}

	if err := tc.DB.Delete(&models.TenantApplication{}, req.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}
