package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
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
// GET /api/v1/members/coupons - 获取当前登录用户的优惠券
// GET /api/v1/members/:id/coupons - 获取指定会员的优惠券
func (c *MemberEnhancedController) MemberCouponList(ctx *gin.Context) {
	// 优先从路径参数获取会员ID，否则从当前登录用户获取
	memberIDStr := ctx.Param("id")
	var memberID uint
	var err error
	
	if memberIDStr != "" {
		// 路径有ID，使用路径ID
		id, err := strconv.ParseUint(memberIDStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的会员ID"})
			return
		}
		memberID = uint(id)
	} else {
		// 路径无ID，从当前登录用户获取关联的会员ID
		userID := middleware.GetUserID(ctx)
		if userID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录或无法获取用户信息"})
			return
		}
		// 查找关联的会员
		var member models.Member
		if err := c.DB.Where("user_id = ?", userID).First(&member).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到关联的会员账号"})
			return
		}
		memberID = member.ID
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

	grants, total, err := c.CouponEngine.GetMemberCoupons(memberID, status)
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

// ListMemberPoints 获取会员积分列表（带统计）
// GET /api/v1/members/points
func (c *MemberEnhancedController) ListMemberPoints(ctx *gin.Context) {
	var members []models.Member
	query := c.DB.Model(&models.Member{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("member_name ILIKE ? OR phone ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if level := ctx.Query("level"); level != "" {
		query = query.Where("member_level = ?", level)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&members)

	// 统计
	var stats struct {
		Total       int64 `json:"total"`
		TodayNew    int64 `json:"today_new"`
		TotalPoints int64 `json:"total_points"`
		MonthUsed   int64 `json:"month_used"`
	}
	c.DB.Model(&models.Member{}).Count(&stats.Total)
	c.DB.Model(&models.Member{}).Where("created_at >= ?", time.Now().Truncate(24*time.Hour)).Count(&stats.TodayNew)
	c.DB.Model(&models.Member{}).Select("COALESCE(SUM(points), 0)").Find(&stats.TotalPoints)
	// MonthUsed: 积分抵扣（points_type=2），取绝对值
	c.DB.Model(&models.MemberPointsRecord{}).Where("points_type = 2 AND created_at >= ?", time.Now().AddDate(0, 0, -30)).Select("COALESCE(SUM(ABS(points)), 0)").Find(&stats.MonthUsed)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":  members,
		"total": total,
		"page":  page,
		"stats": stats,
	}})
}

// AdjustMemberPoints 调整会员积分（通用入口）
// POST /api/v1/members/points/adjust
func (c *MemberEnhancedController) AdjustMemberPoints(ctx *gin.Context) {
	var req struct {
		MemberID uint   `json:"member_id" binding:"required"`
		Type     string `json:"type" binding:"required"` // add / deduct
		Points   int64  `json:"points" binding:"required,min=1"`
		Reason   string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Type == "add" {
		record, err := c.PointsEngine.AddPoints(services.AddPointsReq{
			MemberID:   req.MemberID,
			Points:     req.Points,
			Type:       2, // 活动赠送
			SourceType: "admin_adjust",
			Remark:     req.Reason,
			Operator:   ctx.GetString("user_id"),
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
	} else {
		record, err := c.PointsEngine.DeductPoints(services.DeductPointsReq{
			MemberID:   req.MemberID,
			Points:     req.Points,
			SourceType: "admin_adjust",
			Remark:     req.Reason,
			Operator:   ctx.GetString("user_id"),
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
	}
}
