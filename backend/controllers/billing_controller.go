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

type BillingController struct {
	DB *gorm.DB
}

func (ctrl *BillingController) RegisterRoutes(rg *gin.RouterGroup) {
	billing := rg.Group("/billing")
	{
		billing.GET("/records", ctrl.GetBillingRecords)
		billing.GET("/records/:id", ctrl.GetBillingRecord)
		billing.GET("/invoices", ctrl.GetInvoices)
		billing.POST("/invoices", ctrl.CreateInvoice)
		billing.GET("/invoices/:id", ctrl.GetInvoice)
		billing.PUT("/invoices/:id", ctrl.UpdateInvoice)
		billing.POST("/invoices/:id/issue", ctrl.IssueInvoice)
		billing.POST("/invoices/:id/void", ctrl.VoidInvoice)
		billing.GET("/summary", ctrl.GetBillingSummary)
		billing.POST("/pay", ctrl.PayBilling)
	}
}

func (ctrl *BillingController) GetBillingRecords(c *gin.Context) {
	userID := c.Query("user_id")
	billingType := c.Query("type")
	status := c.DefaultQuery("status", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := ctrl.DB.Model(&models.BillingRecord{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if billingType != "" {
		query = query.Where("billing_type = ?", billingType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var records []models.BillingRecord
	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

func (ctrl *BillingController) GetBillingRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.BillingRecord
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}

func (ctrl *BillingController) GetInvoices(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.DefaultQuery("status", "")

	query := ctrl.DB.Model(&models.Invoice{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var invoices []models.Invoice
	query.Order("created_at DESC").Find(&invoices)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoices})
}

func (ctrl *BillingController) CreateInvoice(c *gin.Context) {
	var req struct {
		UserID          uint    `json:"user_id" binding:"required"`
		SubscriptionID  uint    `json:"subscription_id"`
		OrderID         uint    `json:"order_id"`
		Amount          float64 `json:"amount"`
		InvoiceType     string  `json:"invoice_type" binding:"required"`
		Title           string  `json:"title" binding:"required"`
		TaxNo           string  `json:"tax_no"`
		BankName        string  `json:"bank_name"`
		BankAccount     string  `json:"bank_account"`
		Address         string  `json:"address"`
		Email           string  `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	invoice := models.Invoice{
		InvoiceNo:      fmt.Sprintf("INV-%d-%d", req.UserID, time.Now().Unix()),
		UserID:         req.UserID,
		SubscriptionID: req.SubscriptionID,
		OrderID:        req.OrderID,
		Amount:         req.Amount,
		TaxAmount:      req.Amount * 0.06,
		TotalAmount:    req.Amount * 1.06,
		Status:         "pending",
		InvoiceType:    req.InvoiceType,
		Title:          req.Title,
		TaxNo:          req.TaxNo,
		BankName:       req.BankName,
		BankAccount:    req.BankAccount,
		Address:        req.Address,
		Email:          req.Email,
	}
	ctrl.DB.Create(&invoice)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoice})
}

func (ctrl *BillingController) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoice})
}

func (ctrl *BillingController) UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	ctrl.DB.Model(&invoice).Updates(req)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoice})
}

func (ctrl *BillingController) IssueInvoice(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	ctrl.DB.Model(&models.Invoice{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "issued",
		"issue_date": &now,
		"pdf_url":    fmt.Sprintf("/invoices/%s.pdf", id),
	})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "开票成功"})
}

func (ctrl *BillingController) VoidInvoice(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.Invoice{}).Where("id = ?", id).Update("status", "void")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发票已作废"})
}

func (ctrl *BillingController) GetBillingSummary(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_id required"})
		return
	}

	var totalPending float64
	var totalPaid float64
	var overdueCount int64

	ctrl.DB.Model(&models.BillingRecord{}).Where("user_id = ? AND status = 'pending'", userID).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalPending)
	ctrl.DB.Model(&models.BillingRecord{}).Where("user_id = ? AND status = 'paid'", userID).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalPaid)
	ctrl.DB.Model(&models.BillingRecord{}).Where("user_id = ? AND status = 'overdue'", userID).
		Count(&overdueCount)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"total_pending":  totalPending,
		"total_paid":     totalPaid,
		"overdue_count":  overdueCount,
		"currency":       "CNY",
	}})
}

func (ctrl *BillingController) PayBilling(c *gin.Context) {
	var req struct {
		BillingID uint `json:"billing_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now()
	ctrl.DB.Model(&models.BillingRecord{}).Where("id = ?", req.BillingID).Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "支付成功"})
}
