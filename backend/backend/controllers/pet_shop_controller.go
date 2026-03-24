package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetShopController 宠物商店控制器
type PetShopController struct {
	DB *gorm.DB
}

// ============ PetShopProduct ============

// PetShopProductList 产品列表
func (c *PetShopController) PetShopProductList(ctx *gin.Context) {
	var products []models.PetShopProduct
	var total int64

	query := c.DB.Model(&models.PetShopProduct{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR brand LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if isActive := ctx.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if isFeatured := ctx.Query("is_featured"); isFeatured != "" {
		query = query.Where("is_featured = ?", isFeatured == "true")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("is_featured DESC, sold_count DESC, id DESC").Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": products, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetShopProductGet 获取产品
func (c *PetShopController) PetShopProductGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.PetShopProduct
	if err := c.DB.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

// PetShopProductCreate 创建产品
func (c *PetShopController) PetShopProductCreate(ctx *gin.Context) {
	var product models.PetShopProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

// PetShopProductUpdate 更新产品
func (c *PetShopController) PetShopProductUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.PetShopProduct
	if err := c.DB.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	var updateData models.PetShopProduct
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updateData})
}

// PetShopProductDelete 删除产品
func (c *PetShopController) PetShopProductDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetShopProduct{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetShopOrder ============

// PetShopOrderList 订单列表
func (c *PetShopController) PetShopOrderList(ctx *gin.Context) {
	var orders []models.PetShopOrder
	var total int64

	query := c.DB.Model(&models.PetShopOrder{})

	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if orderStatus := ctx.Query("order_status"); orderStatus != "" {
		query = query.Where("order_status = ?", orderStatus)
	}
	if payStatus := ctx.Query("pay_status"); payStatus != "" {
		query = query.Where("pay_status = ?", payStatus)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Preload("Items").Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": orders, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetShopOrderGet 获取订单
func (c *PetShopController) PetShopOrderGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.PetShopOrder
	if err := c.DB.Preload("Items").First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// PetShopOrderCreate 创建订单
func (c *PetShopController) PetShopOrderCreate(ctx *gin.Context) {
	var order models.PetShopOrder
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// PetShopOrderUpdate 更新订单
func (c *PetShopController) PetShopOrderUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.PetShopOrder
	if err := c.DB.First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
		return
	}
	var updateData models.PetShopOrder
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updateData})
}

// PetShopOrderDelete 删除订单
func (c *PetShopController) PetShopOrderDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetShopOrder{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetShopCart ============

// PetShopCartList 购物车列表
func (c *PetShopController) PetShopCartList(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	var carts []models.PetShopCart
	var total int64

	query := c.DB.Model(&models.PetShopCart{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)

	if err := query.Order("id DESC").Find(&carts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": carts, "total": total,
	}})
}

// PetShopCartAdd 加入购物车
func (c *PetShopController) PetShopCartAdd(ctx *gin.Context) {
	var cart models.PetShopCart
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	// 检查是否已存在
	var existing models.PetShopCart
	err := c.DB.Where("user_id = ? AND product_uuid = ?", cart.UserID, cart.ProductUUID).First(&existing).Error
	if err == nil {
		existing.Quantity += cart.Quantity
		if err := c.DB.Save(&existing).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": existing})
		return
	}
	if err := c.DB.Create(&cart).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": cart})
}

// PetShopCartRemove 移除购物车商品
func (c *PetShopController) PetShopCartRemove(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetShopCart{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetShopFavorite ============

// PetShopFavoriteList 收藏列表
func (c *PetShopController) PetShopFavoriteList(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	var favorites []models.PetShopFavorite
	var total int64

	query := c.DB.Model(&models.PetShopFavorite{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)

	if err := query.Order("id DESC").Find(&favorites).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": favorites, "total": total,
	}})
}

// PetShopFavoriteAdd 添加收藏
func (c *PetShopController) PetShopFavoriteAdd(ctx *gin.Context) {
	var favorite models.PetShopFavorite
	if err := ctx.ShouldBindJSON(&favorite); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&favorite).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": favorite})
}

// PetShopFavoriteRemove 移除收藏
func (c *PetShopController) PetShopFavoriteRemove(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetShopFavorite{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
