package controllers

import (
	"fmt"
	"net/http"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BillingController 计费控制器
type BillingController struct {
	DB *gorm.DB
}

// InvoiceListRequest 发票列表请求
type InvoiceListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Status   string `form:"status"`
	Period   string `form:"billing_period"` // e.g. "2026-03"
}

// InvoiceCreateRequest 创建发票请求
type InvoiceCreateRequest struct {
	Plan          string  `json:"plan" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	BillingPeriod string `json:"billing_period" binding:"required"`
	DueDate       string `json:"due_date"` // RFC3339格式
	Remark        string `json:"remark"`
}

// ListInvoices 获取发票列表
// GET /api/v1/billing/invoices
func (c *BillingController) ListInvoices(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	var req InvoiceListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 2, "message": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	query := c.DB.Model(&models.Invoice{}).Where("tenant_id = ?", tenantID)

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Period != "" {
		query = query.Where("billing_period = ?", req.Period)
	}

	var list []models.Invoice
	var total int64
	query.Count(&total)
	query.Offset((req.Page-1)*req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&list)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": list,
			"pagination": gin.H{
				"total":    total,
				"page":     req.Page,
				"page_size": req.PageSize,
			},
		},
	})
}

// GetInvoice 获取发票详情
// GET /api/v1/billing/invoices/:id
func (c *BillingController) GetInvoice(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	id := ctx.Param("id")
	var invoice models.Invoice
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 2, "message": "发票不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": invoice,
	})
}

// CreateInvoice 创建发票（测试/演示用）
// POST /api/v1/billing/invoices
func (c *BillingController) CreateInvoice(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	var req InvoiceCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 2, "message": "参数错误: " + err.Error()})
		return
	}

	// 查询租户名称
	var tenant models.Tenant
	if err := c.DB.Where("id = ?", tenantID).First(&tenant).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": "获取租户信息失败"})
		return
	}

	// 生成发票号
	invoiceNo := fmt.Sprintf("INV-%s-%d", time.Now().Format("200601"), time.Now().UnixNano()%1000000)

	dueDate := time.Now().AddDate(0, 0, 30)
	if req.DueDate != "" {
		if parsed, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			dueDate = parsed
		}
	}

	invoice := models.Invoice{
		InvoiceNo:     invoiceNo,
		TenantID:      tenantID,
		TenantName:    tenant.Name,
		Plan:          req.Plan,
		Amount:        req.Amount,
		Currency:      "CNY",
		Status:        "pending",
		BillingPeriod: req.BillingPeriod,
		DueDate:       dueDate,
		Remark:        req.Remark,
	}

	if err := c.DB.Create(&invoice).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 4, "message": "创建发票失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "发票创建成功",
		"data": invoice,
	})
}

// PayInvoice 支付发票
// POST /api/v1/billing/invoices/:id/pay
func (c *BillingController) PayInvoice(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "租户ID不能为空"})
		return
	}

	id := ctx.Param("id")
	var invoice models.Invoice
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 2, "message": "发票不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 3, "message": err.Error()})
		return
	}

	if invoice.Status == "paid" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4, "message": "发票已支付"})
		return
	}

	now := time.Now()
	invoice.Status = "paid"
	invoice.PaidAt = &now

	if err := c.DB.Save(&invoice).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5, "message": "支付失败: " + err.Error()})
		return
	}

	// 支付成功后自动升级租户套餐
	if err := c.DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Update("plan", invoice.Plan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 6, "message": "更新套餐失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "支付成功，套餐已升级",
		"data": gin.H{
			"invoice_no": invoice.InvoiceNo,
			"status":     invoice.Status,
			"paid_at":    invoice.PaidAt,
			"plan":       invoice.Plan,
		},
	})
}
