package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderController 会员订单控制器
type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{DB: db}
}

// RegisterRoutes 注册订单路由
func (ctrl *OrderController) RegisterRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	{
		orders.GET("", ctrl.List)
		orders.POST("", ctrl.Create)
		orders.GET("/:id", ctrl.GetByID)
	}
}

// List 获取订单列表
// GET /api/v1/orders
func (ctrl *OrderController) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.Order{}).Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if orderType := c.Query("order_type"); orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}
	if memberID := c.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建订单
// POST /api/v1/orders
func (ctrl *OrderController) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		MemberID     uint    `json:"member_id" binding:"required"`
		TotalAmount  float64 `json:"total_amount" binding:"required"`
		OrderType    int     `json:"order_type"`
		PayType      int     `json:"pay_type"`
		Discount     float64 `json:"discount"`
		PointsEarned int64   `json:"points_earned"`
		Remark       string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 自动生成订单号
	var maxID int64
	ctrl.DB.Model(&models.Order{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
	orderNo := fmt.Sprintf("ORD%d%06d", time.Now().Unix(), maxID+1)

	order := models.Order{
		OrderNo:      orderNo,
		MemberID:    req.MemberID,
		TenantID:    tenantID,
		TotalAmount: req.TotalAmount,
		OrderType:   req.OrderType,
		PayType:     req.PayType,
		Discount:    req.Discount,
		PointsEarned: req.PointsEarned,
		Status:      1,
		Remark:      req.Remark,
	}

	if err := ctrl.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// GetByID 获取订单详情
// GET /api/v1/orders/:id
func (ctrl *OrderController) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var order models.Order
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}
