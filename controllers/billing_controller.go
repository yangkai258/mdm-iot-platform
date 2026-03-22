package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BillingController 账单控制器
type BillingController struct {
	DB *gorm.DB
}

// ===== 请求结构 =====

// CreateInvoiceRequest 申请发票请求
type CreateInvoiceRequest struct {
	BillID        string  `json:"bill_id" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	TaxNo         string  `json:"tax_no"`
	Amount        float64 `json:"amount" binding:"required"`
	InvoiceType   string  `json:"invoice_type"` // normal/special/electronic
	ReceiverEmail string  `json:"receiver_email"`
	ReceiverPhone string  `json:"receiver_phone"`
}

// ===== 账单管理 =====

// ListRecords 获取账单记录
func (c *BillingController) ListRecords(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	var records []models.BillingRecord
	query := c.DB.Where("user_id = ?", userID)

	// 类型过滤
	if billType := ctx.Query("type"); billType != "" {
		query = query.Where("type = ?", billType)
	}

	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 日期范围
	if startDate := ctx.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour-time.Second))
		}
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": records,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// GetRecord 获取账单详情
func (c *BillingController) GetRecord(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	id := ctx.Param("id")
	var record models.BillingRecord

	query := c.DB.Where("bill_id = ?", id)
	// 非管理员只能查看自己的
	if !c.isAdmin(ctx) {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "账单不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    record,
	})
}

// ===== 发票管理 =====

// ListInvoices 获取发票列表
func (c *BillingController) ListInvoices(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	var invoices []models.Invoice
	query := c.DB.Where("user_id = ?", userID)

	// 状态过滤
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 类型过滤
	if invoiceType := ctx.Query("invoice_type"); invoiceType != "" {
		query = query.Where("invoice_type = ?", invoiceType)
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&invoices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": invoices,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// CreateInvoice 申请发票
func (c *BillingController) CreateInvoice(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	var req CreateInvoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	// 验证账单存在且属于该用户
	var bill models.BillingRecord
	if err := c.DB.Where("bill_id = ? AND user_id = ?", req.BillID, userID).First(&bill).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "账单不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 账单必须已支付
	if bill.Status != models.BillingStatusPaid {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "账单未支付，无法申请发票", "error_code": "ERR_BILL_NOT_PAID"})
		return
	}

	// 检查是否已申请过发票
	var existing models.Invoice
	if err := c.DB.Where("bill_id = ? AND user_id = ?", req.BillID, userID).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "该账单已申请过发票", "error_code": "ERR_INVOICE_EXISTS"})
		return
	}

	// 设置默认值
	invoiceType := req.InvoiceType
	if invoiceType == "" {
		invoiceType = "electronic"
	}

	// 计算税额（默认6%）
	taxRate := 6.0
	taxAmount := req.Amount * taxRate / 100

	invoice := models.Invoice{
		InvoiceID:     "inv-" + uuid.New().String(),
		BillID:       req.BillID,
		UserID:       userID,
		Title:        req.Title,
		TaxNo:        req.TaxNo,
		Amount:       req.Amount,
		TaxAmount:    taxAmount,
		TaxRate:      taxRate,
		Status:       models.InvoiceStatusPending,
		InvoiceType:  invoiceType,
		ReceiverEmail: req.ReceiverEmail,
		ReceiverPhone: req.ReceiverPhone,
	}

	if err := c.DB.Create(&invoice).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "申请失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "发票申请成功",
		"data":    invoice,
	})
}

// GetInvoice 获取发票详情
func (c *BillingController) GetInvoice(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	id := ctx.Param("id")
	var invoice models.Invoice

	query := c.DB.Where("invoice_id = ?", id)
	if !c.isAdmin(ctx) {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "发票不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	// 获取关联账单
	var bill models.BillingRecord
	c.DB.Where("bill_id = ?", invoice.BillID).First(&bill)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"invoice": invoice,
			"bill":    bill,
		},
	})
}

// ===== 账单汇总 =====

// GetSummary 获取账单汇总
func (c *BillingController) GetSummary(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 4011, "message": "未登录", "error_code": "ERR_UNAUTHORIZED"})
		return
	}

	// 获取本月统计
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

	// 本月消费
	var monthlyTotal float64
	c.DB.Model(&models.BillingRecord{}).
		Where("user_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", userID, models.BillingStatusPaid, monthStart, monthEnd).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlyTotal)

	// 本月待支付
	var monthlyPending float64
	c.DB.Model(&models.BillingRecord{}).
		Where("user_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", userID, models.BillingStatusPending, monthStart, monthEnd).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlyPending)

	// 累计消费
	var totalSpent float64
	c.DB.Model(&models.BillingRecord{}).
		Where("user_id = ? AND status = ?", userID, models.BillingStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpent)

	// 累计发票
	var totalInvoices int64
	c.DB.Model(&models.Invoice{}).
		Where("user_id = ? AND status = ?", userID, models.InvoiceStatusIssued).
		Count(&totalInvoices)

	// 待开票金额
	var pendingInvoiceAmount float64
	c.DB.Model(&models.Invoice{}).
		Where("user_id = ? AND status = ?", userID, models.InvoiceStatusPending).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&pendingInvoiceAmount)

	// 按类型统计本月消费
	typeSummary := make(map[string]float64)
	var typeRecords []map[string]interface{}
	c.DB.Model(&models.BillingRecord{}).
		Select("type, SUM(amount) as total").
		Where("user_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", userID, models.BillingStatusPaid, monthStart, monthEnd).
		Group("type").
		Scan(&typeRecords)

	for _, r := range typeRecords {
		if t, ok := r["type"].(string); ok {
			if total, ok := r["total"].(float64); ok {
				typeSummary[t] = total
			}
		}
	}

	// 最近消费记录
	var recentBills []models.BillingRecord
	c.DB.Where("user_id = ? AND status = ?", userID, models.BillingStatusPaid).
		Order("created_at DESC").
		Limit(5).
		Find(&recentBills)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"current_period": gin.H{
				"start":       monthStart,
				"end":         monthEnd,
				"spent":       monthlyTotal,
				"pending":     monthlyPending,
				"by_type":     typeSummary,
			},
			"all_time": gin.H{
				"total_spent":          totalSpent,
				"total_invoices":       totalInvoices,
				"pending_invoice_amount": pendingInvoiceAmount,
			},
			"recent_bills": recentBills,
		},
	})
}

// ===== 管理员功能 =====

// IssueInvoice 开票（管理员）
func (c *BillingController) IssueInvoice(ctx *gin.Context) {
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4031, "message": "权限不足", "error_code": "ERR_FORBIDDEN"})
		return
	}

	id := ctx.Param("id")
	var invoice models.Invoice
	if err := c.DB.Where("invoice_id = ?", id).First(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "发票不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	if invoice.Status != models.InvoiceStatusPending {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "发票状态无法开票", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	now := time.Now()
	invoice.Status = models.InvoiceStatusIssued
	invoice.IssuedAt = &now

	if err := c.DB.Save(&invoice).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "开票成功",
		"data":    invoice,
	})
}

// RejectInvoice 拒绝开票（管理员）
func (c *BillingController) RejectInvoice(ctx *gin.Context) {
	if !c.isAdmin(ctx) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4031, "message": "权限不足", "error_code": "ERR_FORBIDDEN"})
		return
	}

	id := ctx.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误", "error_code": "ERR_BAD_REQUEST"})
		return
	}

	var invoice models.Invoice
	if err := c.DB.Where("invoice_id = ?", id).First(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "发票不存在", "error_code": "ERR_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败", "error_code": "ERR_INTERNAL"})
		return
	}

	invoice.Status = models.InvoiceStatusRejected
	if err := c.DB.Save(&invoice).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败", "error_code": "ERR_INTERNAL"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已拒绝开票",
		"data":    invoice,
	})
}

// ===== 辅助函数 =====

func (c *BillingController) isAdmin(ctx *gin.Context) bool {
	role := ctx.GetString("role")
	return role == "admin" || role == "super_admin"
}
