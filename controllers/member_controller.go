package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberController 会员管理控制器
type MemberController struct {
	DB *gorm.DB
}

// NewMemberController 创建会员控制器
func NewMemberController(db *gorm.DB) *MemberController {
	return &MemberController{DB: db}
}

// MemberList 会员列表
func (c *MemberController) MemberList(ctx *gin.Context) {
	var members []models.Member
	var total int64

	query := c.DB.Preload("Card").Model(&models.Member{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("member_name LIKE ? OR member_code LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	// 会员等级筛选（支持 level 和 member_level 两种参数名）
	if level := ctx.Query("level"); level != "" {
		query = query.Where("member_level = ?", level)
	}
	if memberLevel := ctx.Query("member_level"); memberLevel != "" {
		query = query.Where("member_level = ?", memberLevel)
	}
	// 积分范围筛选
	if pointsMin := ctx.Query("points_min"); pointsMin != "" {
		if p, err := strconv.ParseInt(pointsMin, 10, 64); err == nil {
			query = query.Where("points >= ?", p)
		}
	}
	if pointsMax := ctx.Query("points_max"); pointsMax != "" {
		if p, err := strconv.ParseInt(pointsMax, 10, 64); err == nil {
			query = query.Where("points <= ?", p)
		}
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// 时间范围筛选
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Preload("Card").Offset(offset).Limit(pageSize).Order("id DESC").Find(&members).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   members,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// MemberCreate 创建会员
func (c *MemberController) MemberCreate(ctx *gin.Context) {
	var member models.Member
	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&member).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}

// MemberUpdate 更新会员
func (c *MemberController) MemberUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var member models.Member
	if err := c.DB.First(&member, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员不存在"})
		return
	}

	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Save(&member).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}

// MemberDelete 删除会员
func (c *MemberController) MemberDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Member{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// MemberDetail 会员详情
func (c *MemberController) MemberDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	var member models.Member
	if err := c.DB.Preload("Card").First(&member, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}

// CardList 会员卡列表
func (c *MemberController) CardList(ctx *gin.Context) {
	var cards []models.MemberCard
	var total int64

	query := c.DB.Model(&models.MemberCard{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("card_name LIKE ? OR card_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&cards).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   cards,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// CardCreate 创建会员卡
func (c *MemberController) CardCreate(ctx *gin.Context) {
	var card models.MemberCard
	if err := ctx.ShouldBindJSON(&card); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&card).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": card})
}

// CardUpdate 更新会员卡
func (c *MemberController) CardUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var card models.MemberCard
	if err := c.DB.First(&card, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员卡不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&card); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&card).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": card})
}

// CardDelete 删除会员卡
func (c *MemberController) CardDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.MemberCard{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// CouponList 优惠券列表
func (c *MemberController) CouponList(ctx *gin.Context) {
	var coupons []models.Coupon
	var total int64

	query := c.DB.Model(&models.Coupon{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("coupon_name LIKE ? OR coupon_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&coupons).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   coupons,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// CouponCreate 创建优惠券
func (c *MemberController) CouponCreate(ctx *gin.Context) {
	var coupon models.Coupon
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	coupon.RemainStock = coupon.TotalStock
	if err := c.DB.Create(&coupon).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": coupon})
}

// CouponUpdate 更新优惠券
func (c *MemberController) CouponUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var coupon models.Coupon
	if err := c.DB.First(&coupon, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "优惠券不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&coupon).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": coupon})
}

// CouponDelete 删除优惠券
func (c *MemberController) CouponDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Coupon{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// StoreList 店铺列表
func (c *MemberController) StoreList(ctx *gin.Context) {
	var stores []models.Store
	var total int64

	query := c.DB.Model(&models.Store{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("store_name LIKE ? OR store_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&stores).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   stores,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// StoreCreate 创建店铺
func (c *MemberController) StoreCreate(ctx *gin.Context) {
	var store models.Store
	if err := ctx.ShouldBindJSON(&store); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&store).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": store})
}

// StoreUpdate 更新店铺
func (c *MemberController) StoreUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var store models.Store
	if err := c.DB.First(&store, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "店铺不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&store); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&store).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": store})
}

// StoreDelete 删除店铺
func (c *MemberController) StoreDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Store{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PointsRecordList 积分流水列表
func (c *MemberController) PointsRecordList(ctx *gin.Context) {
	var records []models.MemberPointsRecord
	var total int64

	query := c.DB.Model(&models.MemberPointsRecord{})
	if memberID := ctx.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   records,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// TagList 标签列表
func (c *MemberController) TagList(ctx *gin.Context) {
	var tags []models.MemberTag
	var total int64

	query := c.DB.Model(&models.MemberTag{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("tag_name LIKE ? OR tag_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tags).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   tags,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// TagCreate 创建标签
func (c *MemberController) TagCreate(ctx *gin.Context) {
	var tag models.MemberTag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tag})
}

// TagUpdate 更新标签
func (c *MemberController) TagUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var tag models.MemberTag
	if err := c.DB.First(&tag, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "标签不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tag})
}

// TagDelete 删除标签
func (c *MemberController) TagDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.MemberTag{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PromotionList 促销列表
func (c *MemberController) PromotionList(ctx *gin.Context) {
	var promos []models.Promotion
	var total int64

	query := c.DB.Model(&models.Promotion{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("promo_name LIKE ? OR promo_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&promos).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   promos,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// PromotionCreate 创建促销
func (c *MemberController) PromotionCreate(ctx *gin.Context) {
	var promo models.Promotion
	if err := ctx.ShouldBindJSON(&promo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&promo).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": promo})
}

// PromotionUpdate 更新促销
func (c *MemberController) PromotionUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var promo models.Promotion
	if err := c.DB.First(&promo, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "促销不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&promo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&promo).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": promo})
}

// PromotionDelete 删除促销
func (c *MemberController) PromotionDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Promotion{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// LevelList 会员等级列表
func (c *MemberController) LevelList(ctx *gin.Context) {
	var levels []models.MemberLevel
	if err := c.DB.Order("sort ASC, id ASC").Find(&levels).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": levels})
}

// LevelCreate 创建会员等级
func (c *MemberController) LevelCreate(ctx *gin.Context) {
	var level models.MemberLevel
	if err := ctx.ShouldBindJSON(&level); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&level).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": level})
}

// LevelUpdate 更新会员等级
func (c *MemberController) LevelUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var level models.MemberLevel
	if err := c.DB.First(&level, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "等级不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&level); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&level).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": level})
}

// LevelDelete 删除会员等级
func (c *MemberController) LevelDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.MemberLevel{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// OrderList 会员订单列表
func (c *MemberController) OrderList(ctx *gin.Context) {
	var orders []models.MemberOrder
	var total int64

	query := c.DB.Model(&models.MemberOrder{})
	if memberID := ctx.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	if orderType := ctx.Query("order_type"); orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// OrderCreate 创建会员订单
func (c *MemberController) OrderCreate(ctx *gin.Context) {
	var order models.MemberOrder
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if order.OrderNo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "订单号不能为空"})
		return
	}
	if err := c.DB.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// OrderDetail 会员订单详情
func (c *MemberController) OrderDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.MemberOrder
	if err := c.DB.First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订单不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": order})
}

// UpgradeRecordList 等级调整流水列表
func (c *MemberController) UpgradeRecordList(ctx *gin.Context) {
	var records []models.MemberUpgradeRecord
	var total int64

	query := c.DB.Model(&models.MemberUpgradeRecord{})
	if memberID := ctx.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// PointsRuleList 积分规则列表
func (c *MemberController) PointsRuleList(ctx *gin.Context) {
	var rules []models.PointsRule
	var total int64

	query := c.DB.Model(&models.PointsRule{})
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&rules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":   rules,
		"total":  total,
		"page":   page,
		"page_size": pageSize,
	}})
}

// PointsRuleCreate 创建积分规则
func (c *MemberController) PointsRuleCreate(ctx *gin.Context) {
	var rule models.PointsRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// PointsRuleUpdate 更新积分规则
func (c *MemberController) PointsRuleUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.PointsRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// PointsRuleDelete 删除积分规则
func (c *MemberController) PointsRuleDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PointsRule{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// OperationRecordList 会员操作记录列表
func (c *MemberController) OperationRecordList(ctx *gin.Context) {
	var records []models.MemberOperationRecord
	var total int64

	query := c.DB.Model(&models.MemberOperationRecord{})
	if memberID := ctx.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	if operation := ctx.Query("operation"); operation != "" {
		query = query.Where("operation = ?", operation)
	}
	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}
