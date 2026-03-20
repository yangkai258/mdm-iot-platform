package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberEnhancedController 会员增强功能控制器（积分引擎、优惠券）
type MemberEnhancedController struct {
	DB             *gorm.DB
	PointsEngine   *services.PointsEngine
	CouponEngine   *services.CouponEngine
}

// NewMemberEnhancedController 创建增强控制器
func NewMemberEnhancedController(db *gorm.DB) *MemberEnhancedController {
	return &MemberEnhancedController{
		DB:           db,
		PointsEngine: services.NewPointsEngine(db),
		CouponEngine: services.NewCouponEngine(db),
	}
}

// ==================== 积分接口 ====================

// AddPoints 增加积分
// POST /api/v1/members/:id/points/add
func (c *MemberEnhancedController) AddPoints(ctx *gin.Context) {
	id := ctx.Param("id")
	memberID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的会员ID"})
		return
	}

	var req services.AddPointsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	req.MemberID = uint(memberID)

	record, err := c.PointsEngine.AddPoints(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeductPoints 抵扣积分
// POST /api/v1/members/:id/points/deduct
func (c *MemberEnhancedController) DeductPoints(ctx *gin.Context) {
	id := ctx.Param("id")
	memberID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的会员ID"})
		return
	}

	var req services.DeductPointsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	req.MemberID = uint(memberID)

	record, err := c.PointsEngine.DeductPoints(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// GetBalance 查询积分余额
// GET /api/v1/members/:id/points/balance
func (c *MemberEnhancedController) GetBalance(ctx *gin.Context) {
	id := ctx.Param("id")
	memberID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的会员ID"})
		return
	}

	balance, err := c.PointsEngine.GetBalance(uint(memberID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 获取会员等级信息
	var member models.Member
	c.DB.First(&member, memberID)
	var level models.MemberLevel
	levelName := ""
	if c.DB.First(&level, member.MemberLevel).Error == nil {
		levelName = level.LevelName
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"member_id":    memberID,
		"balance":      balance,
		"member_level": member.MemberLevel,
		"level_name":   levelName,
	}})
}

// GetPointsLogs 查询积分流水
// GET /api/v1/members/:id/points/logs
func (c *MemberEnhancedController) GetPointsLogs(ctx *gin.Context) {
	id := ctx.Param("id")
	memberID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的会员ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := c.PointsEngine.GetLogs(uint(memberID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// ==================== 优惠券接口 ====================

// CouponIssue 发放优惠券
// POST /api/v1/coupons/:id/issue
func (c *MemberEnhancedController) CouponIssue(ctx *gin.Context) {
	id := ctx.Param("id")
	couponID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的优惠券ID"})
		return
	}

	var req services.IssueCouponReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	req.CouponID = uint(couponID)

	grant, err := c.CouponEngine.IssueCoupon(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": grant})
}

// CouponUse 核销优惠券
// POST /api/v1/coupons/:id/use
func (c *MemberEnhancedController) CouponUse(ctx *gin.Context) {
	var req services.UseCouponReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	grant, err := c.CouponEngine.UseCoupon(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": grant})
}

// MemberCouponList 获取会员优惠券列表
// GET /api/v1/members/:id/coupons
func (c *MemberEnhancedController) MemberCouponList(ctx *gin.Context) {
	id := ctx.Param("id")
	memberID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的会员ID"})
		return
	}

	status, _ := strconv.Atoi(ctx.DefaultQuery("status", "0"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	grants, total, err := c.CouponEngine.GetMemberCoupons(uint(memberID), status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":      grants,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// CouponListNew 优惠券列表（新路径）
// GET /api/v1/coupons
func (c *MemberEnhancedController) CouponListNew(ctx *gin.Context) {
	var coupons []models.Coupon
	var total int64

	query := c.DB.Model(&models.Coupon{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("coupon_name LIKE ? OR coupon_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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
		"list":      coupons,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// CouponCreateNew 创建优惠券（新路径）
// POST /api/v1/coupons
func (c *MemberEnhancedController) CouponCreateNew(ctx *gin.Context) {
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

// ==================== 促销活动接口（新路径）====================

// PromotionListNew 促销列表（新路径）
// GET /api/v1/promotions
func (c *MemberEnhancedController) PromotionListNew(ctx *gin.Context) {
	var promos []models.Promotion
	var total int64

	query := c.DB.Model(&models.Promotion{})
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("promo_name LIKE ? OR promo_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
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
		"list":      promos,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// PromotionCreateNew 创建促销（新路径）
// POST /api/v1/promotions
func (c *MemberEnhancedController) PromotionCreateNew(ctx *gin.Context) {
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

// PromotionUpdateNew 更新促销（新路径）
// PUT /api/v1/promotions/:id
func (c *MemberEnhancedController) PromotionUpdateNew(ctx *gin.Context) {
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

// PromotionDeleteNew 删除促销（新路径）
// DELETE /api/v1/promotions/:id
func (c *MemberEnhancedController) PromotionDeleteNew(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Promotion{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PromotionDetailNew 促销详情
// GET /api/v1/promotions/:id
func (c *MemberEnhancedController) PromotionDetailNew(ctx *gin.Context) {
	id := ctx.Param("id")
	var promo models.Promotion
	if err := c.DB.First(&promo, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "促销不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": promo})
}
