package models

import (
	"time"

	"gorm.io/gorm"
)

// Member 会员信息
type Member struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	MemberCode   string         `gorm:"size:50;uniqueIndex" json:"member_code"`   // 会员编号
	MemberName  string         `gorm:"size:100" json:"member_name"`               // 会员名称
	Phone       string         `gorm:"size:20;index" json:"phone"`                // 手机号
	Gender      string         `gorm:"size:10" json:"gender"`                     // 性别
	BirthDate   *time.Time     `json:"birth_date"`                               // 生日
	Email       string         `gorm:"size:100" json:"email"`                    // 邮箱
	Avatar      string         `gorm:"size:500" json:"avatar"`                   // 头像
	MemberLevel int            `gorm:"default:1" json:"member_level"`           // 会员等级
	Points      int64          `gorm:"default:0" json:"points"`                  // 积分余额
	Balance     float64        `gorm:"default:0" json:"balance"`                 // 储值余额
	CardID      *uint          `gorm:"index" json:"card_id"`                     // 会员卡ID
	Card        MemberCard     `gorm:"foreignKey:CardID" json:"card"`            // 会员卡关联
	StoreID     *uint          `gorm:"index" json:"store_id"`                   // 所属店铺
	Status      int            `gorm:"default:1" json:"status"`                 // 状态
	Source      string         `gorm:"size:50" json:"source"`                   // 会员来源
	Remark      string         `gorm:"size:500" json:"remark"`                 // 备注
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// MemberCard 会员卡信息
type MemberCard struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CardCode    string    `gorm:"size:50;uniqueIndex" json:"card_code"`    // 卡号
	CardName    string    `gorm:"size:100" json:"card_name"`              // 卡名称
	CardType    int       `gorm:"default:1" json:"card_type"`            // 卡类型: 1储值卡 2积分卡 3打折卡
	GroupID     *uint     `gorm:"index" json:"group_id"`                  // 分组ID
	Discount    float64   `gorm:"default:100" json:"discount"`           // 折扣率
	PointsRate  float64   `gorm:"default:1" json:"points_rate"`          // 积分倍率
	InitPoints  int       `gorm:"default:0" json:"init_points"`          // 开卡赠送积分
	InitBalance float64   `gorm:"default:0" json:"init_balance"`         // 开卡储值
	ValidDays   int       `gorm:"default:0" json:"valid_days"`           // 有效期(天)
	Status      int        `gorm:"default:1" json:"status"`                // 状态
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MemberCardGroup 会员卡分组
type MemberCardGroup struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	GroupName string   `gorm:"size:100" json:"group_name"`    // 分组名称
	ParentID *uint    `gorm:"index" json:"parent_id"`        // 上级分组
	Sort     int       `gorm:"default:0" json:"sort"`        // 排序
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MemberLevel 会员等级
type MemberLevel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LevelName  string    `gorm:"size:50" json:"level_name"`    // 等级名称
	LevelCode  string    `gorm:"size:20" json:"level_code"`  // 等级编码
	MinPoints   int64     `gorm:"default:0" json:"min_points"` // 最低积分
	MaxPoints   int64     `gorm:"default:0" json:"max_points"` // 最高积分
	Discount    float64   `gorm:"default:100" json:"discount"` // 享受折扣
	PointsRate  float64   `gorm:"default:1" json:"points_rate"` // 积分倍率
	Sort        int       `gorm:"default:0" json:"sort"`       // 排序
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MemberUpgradeRule 会员升级规则
type MemberUpgradeRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FromLevel   uint      `gorm:"index" json:"from_level"`    // 原等级
	ToLevel     uint      `gorm:"index" json:"to_level"`      // 目标等级
	PointsThreshold int64  `json:"points_threshold"`           // 积分阈值
	AmountThreshold float64 `json:"amount_threshold"`         // 消费金额阈值
	PointsReward int64     `json:"points_reward"`             // 升级赠送积分
	Status      int        `gorm:"default:1" json:"status"`   // 状态
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MemberTag 会员标签
type MemberTag struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	TagCode  string    `gorm:"size:50;uniqueIndex" json:"tag_code"`   // 标签编码
	TagName  string    `gorm:"size:100" json:"tag_name"`             // 标签名称
	TagType  int       `gorm:"default:1" json:"tag_type"`           // 标签类型: 1行为标签 2属性标签 3自定义
	Category string    `gorm:"size:50" json:"category"`             // 分类
	Color    string    `gorm:"size:20" json:"color"`              // 颜色
	Icon     string    `gorm:"size:50" json:"icon"`               // 图标
	Sort     int        `gorm:"default:0" json:"sort"`             // 排序
	Status   int        `gorm:"default:1" json:"status"`           // 状态
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MemberTagRecord 会员标签流水
type MemberTagRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MemberID  uint      `gorm:"index" json:"member_id"`     // 会员ID
	TagID     uint      `gorm:"index" json:"tag_id"`        // 标签ID
	Action    string    `gorm:"size:20" json:"action"`    // 动作: add/remove
	Operator  string    `gorm:"size:50" json:"operator"`   // 操作人
	CreatedAt time.Time `json:"created_at"`
}

// Coupon 优惠券信息
type Coupon struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CouponCode   string    `gorm:"size:50;uniqueIndex" json:"coupon_code"`   // 优惠券编码
	CouponName   string    `gorm:"size:200" json:"coupon_name"`               // 优惠券名称
	CouponType   int       `gorm:"default:1" json:"coupon_type"`             // 类型: 1满减 2折扣 3兑换
	FaceValue    float64   `gorm:"default:0" json:"face_value"`              // 面值
	MinAmount    float64   `gorm:"default:0" json:"min_amount"`              // 最低消费
	DiscountRate float64   `gorm:"default:100" json:"discount_rate"`        // 折扣率
	TotalStock   int       `gorm:"default:0" json:"total_stock"`            // 总库存
	RemainStock  int       `gorm:"default:0" json:"remain_stock"`           // 剩余库存
	ValidDays    int        `gorm:"default:0" json:"valid_days"`              // 有效期(天)
	StartTime    *time.Time `json:"start_time"`                             // 开始时间
	EndTime     *time.Time `json:"end_time"`                               // 结束时间
	Status       int        `gorm:"default:1" json:"status"`                // 状态
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CouponGrant 优惠券发放记录
type CouponGrant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CouponID  uint      `gorm:"index" json:"coupon_id"`      // 优惠券ID
	MemberID  uint      `gorm:"index" json:"member_id"`      // 会员ID
	GrantTime time.Time `json:"grant_time"`                 // 发放时间
	UseTime   *time.Time `json:"use_time"`                 // 使用时间
	OrderID   *uint     `gorm:"index" json:"order_id"`     // 关联订单
	Status    int        `gorm:"default:1" json:"status"`   // 状态: 1未使用 2已使用 3已过期
	CreatedAt time.Time `json:"created_at"`
}

// Promotion 促销活动
type Promotion struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PromoCode    string    `gorm:"size:50;uniqueIndex" json:"promo_code"`   // 促销编码
	PromoName    string    `gorm:"size:200" json:"promo_name"`               // 促销名称
	PromoType    int       `gorm:"default:1" json:"promo_type"`             // 类型: 1买赠 2直减 3满减 4打折
	StartTime    time.Time `json:"start_time"`                               // 开始时间
	EndTime      time.Time `json:"end_time"`                                 // 结束时间
	RuleConfig   string    `gorm:"type:jsonb" json:"rule_config"`          // 规则配置
	Status       int        `gorm:"default:1" json:"status"`                 // 状态
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Store 店铺信息
type Store struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"size:50;index" json:"tenant_id"`          // 租户ID
	StoreCode  string    `gorm:"size:50;uniqueIndex:idx_tenant_store" json:"store_code"`    // 店铺编码
	StoreName  string    `gorm:"size:200" json:"store_name"`              // 店铺名称
	StoreType  int       `gorm:"default:1" json:"store_type"`             // 类型: 1直营店 2加盟店
	Province   string    `gorm:"size:50" json:"province"`                 // 省
	City       string    `gorm:"size:50" json:"city"`                     // 市
	District   string    `gorm:"size:50" json:"district"`                 // 区
	Address    string    `gorm:"size:500" json:"address"`                 // 详细地址
	Contact    string    `gorm:"size:50" json:"contact"`                 // 联系人
	Phone      string    `gorm:"size:20" json:"phone"`                   // 联系电话
	Status     int        `gorm:"default:1" json:"status"`               // 状态
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PointsRule 积分规则设置
type PointsRule struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	RuleCode string    `gorm:"size:50;uniqueIndex" json:"rule_code"`  // 规则编码
	RuleName string    `gorm:"size:100" json:"rule_name"`             // 规则名称
	RuleType int       `gorm:"default:1" json:"rule_type"`           // 类型: 1获取积分 2抵扣积分 3不积分
	Points   int       `json:"points"`                                // 积分值
	Amount   float64   `json:"amount"`                               // 对应金额
	Remark   string    `gorm:"size:500" json:"remark"`               // 备注
	Status   int        `gorm:"default:1" json:"status"`             // 状态
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MemberPointsRecord 会员积分流水
type MemberPointsRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MemberID   uint      `gorm:"index" json:"member_id"`     // 会员ID
	Points     int64      `json:"points"`                    // 积分变动
	PointsType int       `gorm:"index" json:"points_type"` // 类型: 1消费获得 2兑换抵扣 3活动赠送 4退款返还
	SourceType string    `gorm:"size:20" json:"source_type"` // 来源类型
	SourceID   string    `gorm:"size:50" json:"source_id"`   // 来源ID
	OrderNo    string    `gorm:"size:50" json:"order_no"`    // 关联订单号
	BeforeBalance int64   `json:"before_balance"`            // 变动前余额
	AfterBalance  int64   `json:"after_balance"`             // 变动后余额
	Operator    string    `gorm:"size:50" json:"operator"`  // 操作人
	Remark      string    `gorm:"size:500" json:"remark"`   // 备注
	CreatedAt   time.Time `json:"created_at"`
}

// MemberOperationRecord 会员操作流水
type MemberOperationRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MemberID  uint      `gorm:"index" json:"member_id"`    // 会员ID
	Operation string    `gorm:"size:50" json:"operation"` // 操作类型
	Detail    string    `gorm:"type:text" json:"detail"`  // 操作详情
	Operator  string    `gorm:"size:50" json:"operator"` // 操作人
	IP        string    `gorm:"size:50" json:"ip"`       // IP地址
	CreatedAt time.Time `json:"created_at"`
}

// TempMember 临时会员
type TempMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OpenID    string    `gorm:"size:100;index" json:"open_id"`     // 微信OpenID
	Nickname  string    `gorm:"size:100" json:"nickname"`        // 昵称
	Avatar    string    `gorm:"size:500" json:"avatar"`           // 头像
	Phone     string    `gorm:"size:20" json:"phone"`             // 手机号
	StoreID   *uint    `gorm:"index" json:"store_id"`           // 店铺ID
	ExpireTime time.Time `json:"expire_time"`                   // 过期时间
	CreatedAt time.Time `json:"created_at"`
}

// MemberOrder 会员订单信息
type MemberOrder struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderNo      string    `gorm:"size:50;uniqueIndex" json:"order_no"`
	MemberID     uint      `gorm:"index" json:"member_id"`
	OrderType    int       `gorm:"default:1" json:"order_type"`    // 1消费 2充值 3退款
	Amount       float64   `json:"amount"`                          // 订单金额
	PointsEarned int64     `json:"points_earned"`                   // 获得积分
	PayType      int       `json:"pay_type"`                       // 支付方式
	Status       int       `gorm:"default:1" json:"status"`        // 订单状态
	Remark       string    `gorm:"size:500" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MemberUpgradeRecord 会员等级调整流水
type MemberUpgradeRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MemberID  uint      `gorm:"index" json:"member_id"`
	FromLevel uint      `json:"from_level"`
	ToLevel   uint      `json:"to_level"`
	Reason    string    `gorm:"size:200" json:"reason"`
	Operator  string    `gorm:"size:50" json:"operator"`
	CreatedAt time.Time `json:"created_at"`
}
