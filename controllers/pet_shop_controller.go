package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetShopController 宠物用品电商控制器
type PetShopController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册宠物用品电商相关路由
func (p *PetShopController) RegisterRoutes(r *gin.RouterGroup) {
	shop := r.Group("/pet-shop")
	{
		shop.GET("/products", p.ListProducts)
		shop.GET("/products/:id", p.GetProduct)
		shop.GET("/categories", p.ListCategories)
		shop.POST("/recommendations", p.GetRecommendations)
		shop.POST("/orders", p.CreateOrder)
		shop.GET("/orders", p.ListOrders)
		shop.GET("/orders/:id", p.GetOrder)
		shop.PUT("/orders/:id/cancel", p.CancelOrder)
	}
}

// ListProducts 获取商品列表
func (p *PetShopController) ListProducts(c *gin.Context) {
	var products []models.Product
	query := p.DB.Where("is_active = ?", true)

	// 过滤参数
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if categoryName := c.Query("category_name"); categoryName != "" {
		query = query.Where("category_name LIKE ?", "%"+categoryName+"%")
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR brand LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if brand := c.Query("brand"); brand != "" {
		query = query.Where("brand = ?", brand)
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if isFeatured := c.Query("is_featured"); isFeatured == "true" {
		query = query.Where("is_featured = ?", true)
	}
	if petUUID := c.Query("pet_uuid"); petUUID != "" {
		// 根据宠物类型过滤
		query = query.Where("tags LIKE ? OR tags = ''", "%"+petUUID+"%")
	}

	// 排序
	sortField := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "desc")
	allowedSorts := map[string]bool{"price": true, "sales_count": true, "rating": true, "created_at": true}
	if !allowedSorts[sortField] {
		sortField = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.Product{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询商品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": products,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetProduct 获取商品详情
func (p *PetShopController) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	var product models.Product
	if err := p.DB.Where("product_uuid = ? OR id = ?", productID, productID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "商品不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询商品失败"})
		return
	}

	// 增加浏览次数
	p.DB.Model(&product).Update("view_count", product.ViewCount+1)

	// 解析 images 和 specs JSON
	var images, specs []string
	json.Unmarshal([]byte(product.Images), &images)
	json.Unmarshal([]byte(product.Specs), &specs)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"product":        product,
			"images":         images,
			"specs":          specs,
		},
	})
}

// ListCategories 获取商品分类
func (p *PetShopController) ListCategories(c *gin.Context) {
	var categories []models.ProductCategory
	if err := p.DB.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": categories})
}

// GetRecommendationsRequest 获取推荐请求
type GetRecommendationsRequest struct {
	PetUUID    string  `json:"pet_uuid"`
	PetType    string  `json:"pet_type"`
	CategoryID *uint   `json:"category_id"`
	Limit      int     `json:"limit"`
	PriceRange string  `json:"price_range"` // low/mid/high
}

// GetRecommendations 获取商品推荐
func (p *PetShopController) GetRecommendations(c *gin.Context) {
	var req GetRecommendationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	query := p.DB.Where("is_active = ? AND stock > 0", true)

	// 价格区间过滤
	if req.PriceRange != "" {
		switch req.PriceRange {
		case "low":
			query = query.Where("price < 50")
		case "mid":
			query = query.Where("price >= 50 AND price < 200")
		case "high":
			query = query.Where("price >= 200")
		}
	}

	// 分类过滤
	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	// 如果有宠物类型，根据标签推荐
	if req.PetType != "" {
		query = query.Where("tags LIKE ? OR tags = ''", "%"+req.PetType+"%")
	}

	var products []models.Product
	// 优先推荐高评分、高销量、有特色的商品
	if err := query.Order("is_featured DESC, rating DESC, sales_count DESC").Limit(limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取推荐失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"recommendations": products,
			"count":           len(products),
		},
	})
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	PetUUID          string           `json:"pet_uuid"`
	Items            []OrderItemInput `json:"items" binding:"required,min=1"`
	ShippingAddress  string           `json:"shipping_address" binding:"required"`
	ShippingName     string           `json:"shipping_name" binding:"required"`
	ShippingPhone    string           `json:"shipping_phone" binding:"required"`
	PayMethod        string           `json:"pay_method"`
	Remarks          string           `json:"remarks"`
	HouseholdID      *uint            `json:"household_id"`
	CouponCode       string           `json:"coupon_code"`
}

// OrderItemInput 订单项输入
type OrderItemInput struct {
	ProductUUID string `json:"product_uuid" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	Specs       string `json:"specs"`
}

// CreateOrder 创建订单
func (p *PetShopController) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	// 查询商品并计算总价
	var totalAmount float64
	orderItems := make([]models.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		var product models.Product
		if err := p.DB.Where("product_uuid = ? OR id = ?", item.ProductUUID, item.ProductUUID).First(&product).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": fmt.Sprintf("商品不存在: %s", item.ProductUUID)})
			return
		}

		if !product.IsActive {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": fmt.Sprintf("商品已下架: %s", product.Name)})
			return
		}

		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": fmt.Sprintf("库存不足: %s, 剩余 %d", product.Name, product.Stock)})
			return
		}

		orderItem := models.OrderItem{
			ProductUUID:  product.ProductUUID,
			ProductName:  product.Name,
			ProductImage: product.Images,
			Price:        product.Price,
			Quantity:     item.Quantity,
			Specs:        item.Specs,
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += product.Price * float64(item.Quantity)

		// 扣减库存
		p.DB.Model(&product).Update("stock", product.Stock-item.Quantity)
		// 增加销量
		p.DB.Model(&product).Update("sales_count", product.SalesCount+item.Quantity)
	}

	// 序列化订单项
	itemsJSON, _ := json.Marshal(orderItems)

	order := models.Order{
		PetUUID:         req.PetUUID,
		TotalAmount:    totalAmount,
		DiscountAmount: 0,
		PayAmount:      totalAmount,
		PayMethod:      req.PayMethod,
		PayStatus:      "pending",
		Status:         "pending",
		ShippingAddress: req.ShippingAddress,
		ShippingName:   req.ShippingName,
		ShippingPhone:  req.ShippingPhone,
		Remarks:        req.Remarks,
		Items:          string(itemsJSON),
		HouseholdID:    req.HouseholdID,
		OwnerID:        userID,
		TenantID:       fmt.Sprintf("%d", tenantID),
	}

	if order.PayMethod == "" {
		order.PayMethod = "wechat"
	}

	if err := p.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建订单失败: " + err.Error()})
		return
	}

	// 解析返回订单项
	var orderItemsResult []models.OrderItem
	json.Unmarshal([]byte(order.Items), &orderItemsResult)

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "订单创建成功",
		"data": gin.H{
			"order":       order,
			"order_items": orderItemsResult,
		},
	})
}

// ListOrders 获取订单列表
func (p *PetShopController) ListOrders(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var orders []models.Order
	query := p.DB.Where("owner_id = ? AND tenant_id = ?", userID, tenantID)

	// 过滤参数
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if payStatus := c.Query("pay_status"); payStatus != "" {
		query = query.Where("pay_status = ?", payStatus)
	}
	if petUUID := c.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询订单列表失败"})
		return
	}

	// 解析订单项
	type OrderWithItems struct {
		models.Order
		Items []models.OrderItem `json:"items"`
	}
	result := make([]OrderWithItems, len(orders))
	for i, o := range orders {
		result[i] = OrderWithItems{Order: o}
		json.Unmarshal([]byte(o.Items), &result[i].Items)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": result,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetOrder 获取订单详情
func (p *PetShopController) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var order models.Order
	if err := p.DB.Where("(order_uuid = ? OR id = ? OR order_no = ?) AND owner_id = ? AND tenant_id = ?",
		orderID, orderID, orderID, userID, tenantID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询订单失败"})
		return
	}

	var items []models.OrderItem
	json.Unmarshal([]byte(order.Items), &items)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"order":  order,
			"items":  items,
		},
	})
}

// CancelOrder 取消订单
func (p *PetShopController) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var order models.Order
	if err := p.DB.Where("(order_uuid = ? OR id = ? OR order_no = ?) AND owner_id = ? AND tenant_id = ?",
		orderID, orderID, orderID, userID, tenantID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询订单失败"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前状态无法取消订单"})
		return
	}

	// 恢复库存
	var items []models.OrderItem
	json.Unmarshal([]byte(order.Items), &items)
	for _, item := range items {
		var product models.Product
		if err := p.DB.Where("product_uuid = ?", item.ProductUUID).First(&product).Error; err == nil {
			p.DB.Model(&product).Update("stock", product.Stock+item.Quantity)
			p.DB.Model(&product).Update("sales_count", product.SalesCount-item.Quantity)
		}
	}

	updates := map[string]interface{}{
		"status":     "cancelled",
		"pay_status": "refunded",
	}
	if err := p.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消订单失败"})
		return
	}

	p.DB.First(&order, order.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订单已取消", "data": order})
}
