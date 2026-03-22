package services

import (
	"errors"
	"strconv"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// PointsEngine 积分引擎
type PointsEngine struct {
	DB *gorm.DB
}

// NewPointsEngine 创建积分引擎实例
func NewPointsEngine(db *gorm.DB) *PointsEngine {
	return &PointsEngine{DB: db}
}

// PointsChangeType 积分变动类型
const (
	PointsTypeEarn   = 1 // 消费获得
	PointsTypeDeduct = 2 // 兑换抵扣
	PointsTypeGift   = 3 // 活动赠送
	PointsTypeRefund = 4 // 退款返还
)

// AddPointsReq 增加积分请求
type AddPointsReq struct {
	MemberID   uint   `json:"member_id" binding:"required"`
	Points     int64  `json:"points" binding:"required,min=1"`
	Type       int    `json:"type" binding:"required"`       // 1消费获得 2活动赠送 3退款返还
	SourceType string `json:"source_type"`                   // 来源类型
	SourceID   string `json:"source_id"`                     // 来源ID
	OrderNo    string `json:"order_no"`                      // 关联订单号
	Operator   string `json:"operator"`                      // 操作人
	Remark     string `json:"remark"`                       // 备注
}

// DeductPointsReq 抵扣积分请求
type DeductPointsReq struct {
	MemberID   uint   `json:"member_id" binding:"required"`
	Points     int64  `json:"points" binding:"required,min=1"`
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	OrderNo    string `json:"order_no"`
	Operator   string `json:"operator"`
	Remark     string `json:"remark"`
}

// AddPoints 增加积分
func (e *PointsEngine) AddPoints(req AddPointsReq) (*models.MemberPointsRecord, error) {
	var member models.Member
	if err := e.DB.First(&member, req.MemberID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会员不存在")
		}
		return nil, err
	}

	if member.Status != 1 {
		return nil, errors.New("会员状态异常，无法操作积分")
	}

	beforeBalance := member.Points
	afterBalance := beforeBalance + req.Points

	// 更新会员积分
	if err := e.DB.Model(&member).Update("points", afterBalance).Error; err != nil {
		return nil, err
	}

	// 记录流水
	record := &models.MemberPointsRecord{
		MemberID:      req.MemberID,
		Points:        req.Points,
		PointsType:    req.Type,
		SourceType:    req.SourceType,
		SourceID:      req.SourceID,
		OrderNo:       req.OrderNo,
		BeforeBalance: beforeBalance,
		AfterBalance:  afterBalance,
		Operator:      req.Operator,
		Remark:        req.Remark,
	}
	if err := e.DB.Create(record).Error; err != nil {
		// 回滚积分
		e.DB.Model(&member).Update("points", beforeBalance)
		return nil, err
	}

	// 检查是否需要升级
	e.checkAndUpgrade(req.MemberID)

	return record, nil
}

// DeductPoints 抵扣积分
func (e *PointsEngine) DeductPoints(req DeductPointsReq) (*models.MemberPointsRecord, error) {
	var member models.Member
	if err := e.DB.First(&member, req.MemberID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会员不存在")
		}
		return nil, err
	}

	if member.Status != 1 {
		return nil, errors.New("会员状态异常，无法操作积分")
	}

	if member.Points < req.Points {
		return nil, errors.New("积分余额不足")
	}

	beforeBalance := member.Points
	afterBalance := beforeBalance - req.Points

	// 更新会员积分
	if err := e.DB.Model(&member).Update("points", afterBalance).Error; err != nil {
		return nil, err
	}

	// 记录流水
	record := &models.MemberPointsRecord{
		MemberID:      req.MemberID,
		Points:        -req.Points,
		PointsType:    PointsTypeDeduct,
		SourceType:    req.SourceType,
		SourceID:      req.SourceID,
		OrderNo:       req.OrderNo,
		BeforeBalance: beforeBalance,
		AfterBalance:  afterBalance,
		Operator:      req.Operator,
		Remark:        req.Remark,
	}
	if err := e.DB.Create(record).Error; err != nil {
		// 回滚积分
		e.DB.Model(&member).Update("points", beforeBalance)
		return nil, err
	}

	// 检查是否需要降级
	e.checkAndDowngrade(req.MemberID)

	return record, nil
}

// GetBalance 查询积分余额
func (e *PointsEngine) GetBalance(memberID uint) (int64, error) {
	var member models.Member
	if err := e.DB.First(&member, memberID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("会员不存在")
		}
		return 0, err
	}
	return member.Points, nil
}

// GetLogs 查询积分流水
func (e *PointsEngine) GetLogs(memberID uint, page, pageSize int) ([]models.MemberPointsRecord, int64, error) {
	var records []models.MemberPointsRecord
	var total int64

	query := e.DB.Model(&models.MemberPointsRecord{}).Where("member_id = ?", memberID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// checkAndUpgrade 检查并执行升级
func (e *PointsEngine) checkAndUpgrade(memberID uint) {
	var member models.Member
	if err := e.DB.First(&member, memberID).Error; err != nil {
		return
	}

	// 查询所有等级，按最低积分降序
	var levels []models.MemberLevel
	e.DB.Where("status = 1").Order("min_points DESC").Find(&levels)

	for _, level := range levels {
		if member.Points >= level.MinPoints && member.MemberLevel < int(level.ID) {
			// 升级
			e.DB.Model(&member).Update("member_level", level.ID)

			// 记录升级流水
			upgradeRecord := &models.MemberUpgradeRecord{
				MemberID:  memberID,
				FromLevel: uint(member.MemberLevel),
				ToLevel:   level.ID,
				Reason:    "积分达到升级条件",
				Operator:  "system",
			}
			e.DB.Create(upgradeRecord)

			// 检查是否有升级奖励规则
			var rule models.MemberUpgradeRule
			if err := e.DB.Where("from_level = ? AND to_level = ? AND status = 1", member.MemberLevel, level.ID).First(&rule).Error; err == nil {
				if rule.PointsReward > 0 {
					// 发放升级奖励积分
					rewardReq := AddPointsReq{
						MemberID:   memberID,
						Points:     rule.PointsReward,
						Type:       PointsTypeGift,
						SourceType: "upgrade_reward",
						SourceID:   strconv.FormatUint(uint64(rule.ID), 10),
						Operator:   "system",
						Remark:     "升级奖励",
					}
					e.AddPoints(rewardReq)
				}
			}
			break
		}
	}
}

// checkAndDowngrade 检查并执行降级
func (e *PointsEngine) checkAndDowngrade(memberID uint) {
	var member models.Member
	if err := e.DB.First(&member, memberID).Error; err != nil {
		return
	}

	// 查询所有等级，按最低积分降序
	var levels []models.MemberLevel
	e.DB.Where("status = 1").Order("min_points DESC").Find(&levels)

	for i := len(levels) - 1; i >= 0; i-- {
		level := levels[i]
		if member.Points < level.MinPoints && member.MemberLevel > int(level.ID) {
			// 降级（只降一级）
			e.DB.Model(&member).Update("member_level", level.ID)

			upgradeRecord := &models.MemberUpgradeRecord{
				MemberID:  memberID,
				FromLevel: uint(member.MemberLevel),
				ToLevel:   level.ID,
				Reason:    "积分不足降级",
				Operator:  "system",
			}
			e.DB.Create(upgradeRecord)
			break
		}
	}
}

// CalculatePoints 计算积分（根据消费金额和规则）
func (e *PointsEngine) CalculatePoints(memberID uint, amount float64, ruleCode string) (int64, error) {
	var member models.Member
	if err := e.DB.First(&member, memberID).Error; err != nil {
		return 0, errors.New("会员不存在")
	}

	// 查询积分规则
	var rule models.PointsRule
	query := e.DB.Model(&models.PointsRule{}).Where("status = 1").Where("rule_type = 1")
	if ruleCode != "" {
		query = query.Where("rule_code = ?", ruleCode)
	}
	if err := query.First(&rule).Error; err != nil {
		// 无规则时，默认每消费1元得1分
		return int64(amount), nil
	}

	// 获取会员等级倍率
	var level models.MemberLevel
	if err := e.DB.First(&level, member.MemberLevel).Error; err == nil {
		return int64(float64(rule.Points) * (amount / rule.Amount) * level.PointsRate), nil
	}

	return int64(float64(rule.Points) * (amount / rule.Amount)), nil
}

// CouponEngine 优惠券引擎
type CouponEngine struct {
	DB *gorm.DB
}

// NewCouponEngine 创建优惠券引擎
func NewCouponEngine(db *gorm.DB) *CouponEngine {
	return &CouponEngine{DB: db}
}

// IssueCouponReq 发放优惠券请求
type IssueCouponReq struct {
	CouponID  uint   `json:"coupon_id" binding:"required"`
	MemberID  uint   `json:"member_id" binding:"required"`
	Operator  string `json:"operator"`
	Remark    string `json:"remark"`
}

// UseCouponReq 核销优惠券请求
type UseCouponReq struct {
	CouponGrantID uint   `json:"coupon_grant_id" binding:"required"`
	OrderID       uint   `json:"order_id"`
	Operator      string `json:"operator"`
	Remark        string `json:"remark"`
}

// IssueCoupon 发放优惠券
func (e *CouponEngine) IssueCoupon(req IssueCouponReq) (*models.CouponGrant, error) {
	// 检查优惠券是否存在
	var coupon models.Coupon
	if err := e.DB.First(&coupon, req.CouponID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("优惠券不存在")
		}
		return nil, err
	}

	if coupon.Status != 1 {
		return nil, errors.New("优惠券已下架")
	}

	if coupon.RemainStock <= 0 {
		return nil, errors.New("优惠券库存不足")
	}

	// 检查会员是否存在
	var member models.Member
	if err := e.DB.First(&member, req.MemberID).Error; err != nil {
		return nil, errors.New("会员不存在")
	}

	// 扣减库存
	if err := e.DB.Model(&coupon).Update("remain_stock", coupon.RemainStock-1).Error; err != nil {
		return nil, err
	}

	// 创建发放记录
	now := time.Now()
	grant := &models.CouponGrant{
		CouponID:  req.CouponID,
		MemberID:  req.MemberID,
		GrantTime: now,
		Status:    1, // 未使用
	}
	if err := e.DB.Create(grant).Error; err != nil {
		// 回滚库存
		e.DB.Model(&coupon).Update("remain_stock", coupon.RemainStock)
		return nil, err
	}

	return grant, nil
}

// UseCoupon 核销优惠券
func (e *CouponEngine) UseCoupon(req UseCouponReq) (*models.CouponGrant, error) {
	var grant models.CouponGrant
	if err := e.DB.First(&grant, req.CouponGrantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("优惠券发放记录不存在")
		}
		return nil, err
	}

	if grant.Status == 2 {
		return nil, errors.New("优惠券已核销")
	}
	if grant.Status == 3 {
		return nil, errors.New("优惠券已过期")
	}

	// 检查有效期
	var coupon models.Coupon
	if err := e.DB.First(&coupon, grant.CouponID).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	// 检查固定有效期
	if coupon.StartTime != nil && now.Before(*coupon.StartTime) {
		return nil, errors.New("优惠券未到使用期")
	}
	if coupon.EndTime != nil && now.After(*coupon.EndTime) {
		// 更新为已过期状态
		e.DB.Model(&grant).Update("status", 3)
		return nil, errors.New("优惠券已过期")
	}
	// 检查相对有效期
	if coupon.ValidDays > 0 {
		expireTime := grant.GrantTime.AddDate(0, 0, coupon.ValidDays)
		if now.After(expireTime) {
			e.DB.Model(&grant).Update("status", 3)
			return nil, errors.New("优惠券已过期")
		}
	}

	// 核销
	nowTime := time.Now()
	grant.UseTime = &nowTime
	grant.OrderID = &req.OrderID
	grant.Status = 2 // 已使用

	if err := e.DB.Save(&grant).Error; err != nil {
		return nil, err
	}

	return &grant, nil
}

// GetMemberCoupons 获取会员的优惠券列表
func (e *CouponEngine) GetMemberCoupons(memberID uint, status int) ([]models.CouponGrant, int64, error) {
	var grants []models.CouponGrant
	var total int64

	query := e.DB.Model(&models.CouponGrant{}).Where("member_id = ?", memberID)
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)

	if err := query.Preload("Coupon").Order("id DESC").Find(&grants).Error; err != nil {
		return nil, 0, err
	}

	return grants, total, nil
}
