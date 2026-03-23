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

// InvoiceStatus 发票状态常量
const (
	InvoiceStatusPending  = "pending"  // 待审核
	InvoiceStatusApproved = "approved"  // 已审核
	InvoiceStatusRejected = "rejected" // 已拒绝
	InvoiceStatusIssued   = "issued"   // 已开票
	InvoiceStatusVoid     = "void"     // 已作废
	InvoiceStatusSent     = "sent"     // 已寄送
	InvoiceStatusDelivered = "delivered" // 已送达
)

// InvoiceReviewRequest 发票审核请求
type InvoiceReviewRequest struct {
	Approved  bool   `json:"approved" binding:"required"`
	Comment   string `json:"comment"`
}

// InvoiceDeliveryRequest 发票寄送请求
type InvoiceDeliveryRequest struct {
	Carrier       string `json:"carrier" binding:"required"`        // 快递公司
	TrackingNo    string `json:"tracking_no" binding:"required"`    // 快递单号
	RecipientName string `json:"recipient_name" binding:"required"`  // 收件人
	RecipientPhone string `json:"recipient_phone" binding:"required"` // 收件人电话
	RecipientAddr string `json:"recipient_addr" binding:"required"`  // 收件地址
}

// InvoiceShipping 发票寄送记录
type InvoiceShipping struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	InvoiceID       uint       `gorm:"index;not null" json:"invoice_id"`
	Carrier         string     `gorm:"type:varchar(50)" json:"carrier"`           // 快递公司
	TrackingNo      string     `gorm:"type:varchar(100)" json:"tracking_no"`     // 快递单号
	RecipientName   string     `gorm:"type:varchar(100)" json:"recipient_name"`  // 收件人
	RecipientPhone  string     `gorm:"type:varchar(20)" json:"recipient_phone"`  // 收件人电话
	RecipientAddr   string     `gorm:"type:varchar(500)" json:"recipient_addr"`  // 收件地址
	ShippedAt       *time.Time `json:"shipped_at"`                               // 寄出时间
	DeliveredAt     *time.Time `json:"delivered_at"`                             // 送达时间
	DeliveryStatus  string     `gorm:"type:varchar(20);default:'pending'" json:"delivery_status"` // pending, in_transit, delivered, failed
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName 设置表名
func (InvoiceShipping) TableName() string {
	return "invoice_shippings"
}

// BillingController 账单控制器
type BillingController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册路由
func (ctrl *BillingController) RegisterRoutes(rg *gin.RouterGroup) {
	billing := rg.Group("/billing")
	{
		// 计费记录
		billing.GET("/records", ctrl.GetBillingRecords)
		billing.GET("/records/:id", ctrl.GetBillingRecord)
		billing.POST("/pay", ctrl.PayBilling)

		// 发票
		billing.GET("/invoices", ctrl.GetInvoices)
		billing.POST("/invoices", ctrl.CreateInvoice)
		billing.GET("/invoices/:id", ctrl.GetInvoice)
		billing.PUT("/invoices/:id", ctrl.UpdateInvoice)

		// 发票审核
		billing.POST("/invoices/:id/review", ctrl.ReviewInvoice)
		billing.POST("/invoices/:id/approve", ctrl.ApproveInvoice)
		billing.POST("/invoices/:id/reject", ctrl.RejectInvoice)

		// 开票/作废
		billing.POST("/invoices/:id/issue", ctrl.IssueInvoice)
		billing.POST("/invoices/:id/void", ctrl.VoidInvoice)

		// 发票寄送
		billing.POST("/invoices/:id/ship", ctrl.ShipInvoice)
		billing.GET("/invoices/:id/tracking", ctrl.GetInvoiceTracking)
		billing.POST("/invoices/:id/track", ctrl.UpdateTrackingStatus)

		// 发票PDF
		billing.GET("/invoices/:id/pdf", ctrl.GeneratePDF)

		// 汇总
		billing.GET("/summary", ctrl.GetBillingSummary)
	}
}

// GetBillingRecords 获取计费记录列表
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
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// GetBillingRecord 获取单条计费记录
func (ctrl *BillingController) GetBillingRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.BillingRecord
	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": record})
}

// GetInvoices 获取发票列表
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

// CreateInvoice 创建发票申请
func (ctrl *BillingController) CreateInvoice(c *gin.Context) {
	var req struct {
		UserID          uint    `json:"user_id" binding:"required"`
		SubscriptionID  uint    `json:"subscription_id"`
		OrderID         uint    `json:"order_id"`
		Amount          float64 `json:"amount" binding:"required"`
		InvoiceType     string  `json:"invoice_type" binding:"required"`
		Title           string  `json:"title" binding:"required"`
		TaxNo           string  `json:"tax_no"`
		BankName        string  `json:"bank_name"`
		BankAccount     string  `json:"bank_account"`
		Address         string  `json:"address"`
		Email           string  `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 验证增值税发票必填字段
	if req.InvoiceType == "VAT" {
		if req.TaxNo == "" || req.BankName == "" || req.BankAccount == "" || req.Address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "增值税发票需要填写税号、开户行、账号、地址电话"})
			return
		}
	}

	invoice := models.Invoice{
		InvoiceNo:      fmt.Sprintf("INV-%d-%d", req.UserID, time.Now().Unix()),
		UserID:         req.UserID,
		SubscriptionID: req.SubscriptionID,
		OrderID:        req.OrderID,
		Amount:         req.Amount,
		TaxAmount:      req.Amount * 0.06,
		TotalAmount:    req.Amount * 1.06,
		Status:         InvoiceStatusPending, // 默认待审核
		InvoiceType:    req.InvoiceType,
		Title:          req.Title,
		TaxNo:          req.TaxNo,
		BankName:       req.BankName,
		BankAccount:    req.BankAccount,
		Address:        req.Address,
		Email:          req.Email,
	}

	if err := ctrl.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建发票失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoice, "message": "发票申请已提交，等待审核"})
}

// GetInvoice 获取发票详情
func (ctrl *BillingController) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	// 获取寄送信息
	var shipping InvoiceShipping
	ctrl.DB.Where("invoice_id = ?", id).First(&shipping)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"invoice":  invoice,
		"shipping": shipping,
	}})
}

// UpdateInvoice 更新发票
func (ctrl *BillingController) UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	// 只有待审核状态可修改
	if invoice.Status != InvoiceStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有待审核的发票可以修改"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 过滤不可修改的字段
	delete(req, "invoice_no")
	delete(req, "status")
	delete(req, "created_at")

	if err := ctrl.DB.Model(&invoice).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoice})
}

// ReviewInvoice 审核发票（通用）
func (ctrl *BillingController) ReviewInvoice(c *gin.Context) {
	id := c.Param("id")
	var req InvoiceReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	if invoice.Status != InvoiceStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有待审核的发票可以审核"})
		return
	}

	now := time.Now()
	if req.Approved {
		invoice.Status = InvoiceStatusApproved
	} else {
		invoice.Status = InvoiceStatusRejected
	}

	if err := ctrl.DB.Model(&invoice).Updates(map[string]interface{}{
		"status":    invoice.Status,
		"issue_date": &now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "审核失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "审核完成", "data": gin.H{
		"status":  invoice.Status,
		"comment": req.Comment,
	}})
}

// ApproveInvoice 批准发票
func (ctrl *BillingController) ApproveInvoice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Comment string `json:"comment"`
	}
	c.ShouldBindJSON(&req)

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	if invoice.Status != InvoiceStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有待审核的发票可以批准"})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&invoice).Updates(map[string]interface{}{
		"status":     InvoiceStatusApproved,
		"issue_date": &now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批准失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发票已批准", "data": gin.H{"status": InvoiceStatusApproved}})
}

// RejectInvoice 拒绝发票
func (ctrl *BillingController) RejectInvoice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请填写拒绝原因"})
		return
	}

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	if invoice.Status != InvoiceStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有待审核的发票可以拒绝"})
		return
	}

	if err := ctrl.DB.Model(&invoice).Update("status", InvoiceStatusRejected).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "拒绝失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发票已拒绝", "data": gin.H{
		"status": InvoiceStatusRejected,
		"reason": req.Reason,
	}})
}

// IssueInvoice 开票
func (ctrl *BillingController) IssueInvoice(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	// 只有已审核的发票可以开票
	if invoice.Status != InvoiceStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有已审核的发票可以开票"})
		return
	}

	pdfURL := fmt.Sprintf("/invoices/%d.pdf", id)

	if err := ctrl.DB.Model(&invoice).Updates(map[string]interface{}{
		"status":     InvoiceStatusIssued,
		"issue_date": &now,
		"pdf_url":    pdfURL,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "开票失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "开票成功", "data": gin.H{
		"pdf_url": pdfURL,
	}})
}

// VoidInvoice 作废发票
func (ctrl *BillingController) VoidInvoice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请填写作废原因"})
		return
	}

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	// 已开票的发票可以作废
	if invoice.Status != InvoiceStatusIssued && invoice.Status != InvoiceStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有已开票的发票可以作废"})
		return
	}

	if err := ctrl.DB.Model(&invoice).Update("status", InvoiceStatusVoid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "作废失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发票已作废", "data": gin.H{
		"reason": req.Reason,
	}})
}

// ShipInvoice 寄送发票
func (ctrl *BillingController) ShipInvoice(c *gin.Context) {
	id := c.Param("id")
	var req InvoiceDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	// 只有已开票的可以寄送
	if invoice.Status != InvoiceStatusIssued {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有已开票的发票可以寄送"})
		return
	}

	now := time.Now()
	shipping := InvoiceShipping{
		InvoiceID:       invoice.ID,
		Carrier:         req.Carrier,
		TrackingNo:      req.TrackingNo,
		RecipientName:   req.RecipientName,
		RecipientPhone:  req.RecipientPhone,
		RecipientAddr:   req.RecipientAddr,
		ShippedAt:       &now,
		DeliveryStatus:  "in_transit",
	}

	if err := ctrl.DB.Create(&shipping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "寄送记录创建失败"})
		return
	}

	// 更新发票状态
	ctrl.DB.Model(&invoice).Update("status", InvoiceStatusSent)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发票已寄出", "data": shipping})
}

// GetInvoiceTracking 获取发票寄送跟踪信息
func (ctrl *BillingController) GetInvoiceTracking(c *gin.Context) {
	id := c.Param("id")

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	var shipping InvoiceShipping
	if err := ctrl.DB.Where("invoice_id = ?", id).First(&shipping).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
			"has_shipping": false,
			"message":      "暂无寄送信息",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"has_shipping":     true,
		"carrier":          shipping.Carrier,
		"tracking_no":      shipping.TrackingNo,
		"recipient_name":   shipping.RecipientName,
		"recipient_phone":  shipping.RecipientPhone,
		"recipient_addr":   shipping.RecipientAddr,
		"shipped_at":      shipping.ShippedAt,
		"delivered_at":    shipping.DeliveredAt,
		"delivery_status":  shipping.DeliveryStatus,
	}})
}

// UpdateTrackingStatus 更新物流状态（模拟快递公司回调）
func (ctrl *BillingController) UpdateTrackingStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"` // in_transit, delivered, failed
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var shipping InvoiceShipping
	if err := ctrl.DB.Where("invoice_id = ?", id).First(&shipping).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "寄送记录不存在"})
		return
	}

	updates := map[string]interface{}{
		"delivery_status": req.Status,
	}

	if req.Status == "delivered" {
		now := time.Now()
		updates["delivered_at"] = &now

		// 更新发票状态为已送达
		ctrl.DB.Model(&models.Invoice{}).Where("id = ?", id).Update("status", InvoiceStatusDelivered)
	}

	if err := ctrl.DB.Model(&shipping).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "物流状态已更新"})
}

// GeneratePDF 生成发票PDF（stub实现）
func (ctrl *BillingController) GeneratePDF(c *gin.Context) {
	id := c.Param("id")

	var invoice models.Invoice
	if err := ctrl.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	if invoice.Status != InvoiceStatusIssued {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只有已开票的发票才能生成PDF"})
		return
	}

	// Stub: 实际项目中应使用 PDF 库（如 github.com/jung-kurt/gofpdf）生成
	// 这里返回模拟的PDF内容
	pdfContent := fmt.Sprintf("Invoice PDF Content - InvoiceNo: %s, Amount: %.2f",
		invoice.InvoiceNo, invoice.TotalAmount)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"invoice_id":  invoice.ID,
		"invoice_no":  invoice.InvoiceNo,
		"pdf_url":     invoice.PDFURL,
		"pdf_content": pdfContent, // 实际应为PDF二进制流
	}})
}

// GetBillingSummary 获取账单汇总
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

	// 发票统计
	var invoicePending, invoiceIssued, invoiceVoid int64
	ctrl.DB.Model(&models.Invoice{}).Where("user_id = ? AND status = 'pending'", userID).Count(&invoicePending)
	ctrl.DB.Model(&models.Invoice{}).Where("user_id = ? AND status = 'issued'", userID).Count(&invoiceIssued)
	ctrl.DB.Model(&models.Invoice{}).Where("user_id = ? AND status = 'void'", userID).Count(&invoiceVoid)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"total_pending":   totalPending,
		"total_paid":      totalPaid,
		"overdue_count":   overdueCount,
		"currency":        "CNY",
		"invoice_pending": invoicePending,
		"invoice_issued":  invoiceIssued,
		"invoice_void":   invoiceVoid,
	}})
}

// PayBilling 支付账单
func (ctrl *BillingController) PayBilling(c *gin.Context) {
	var req struct {
		BillingID uint `json:"billing_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now()
	if err := ctrl.DB.Model(&models.BillingRecord{}).Where("id = ?", req.BillingID).Updates(map[string]interface{}{
		"status":  "paid",
		"paid_at": &now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "支付失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "支付成功"})
}
