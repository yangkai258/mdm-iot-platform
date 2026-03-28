package models

import (
	"time"

	"gorm.io/gorm"
)

// Member360Profile 会员360度画像
type Member360Profile struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	MemberID        string         `gorm:"type:varchar(64);uniqueIndex" json:"member_id"`
	MemberName      string         `gorm:"type:varchar(100)" json:"member_name"`

	// 基础画像
	Age             int            `gorm:"default:0" json:"age"`
	Gender          string         `gorm:"type:varchar(10)" json:"gender"`
	Location        string         `gorm:"type:varchar(200)" json:"location"`
	Occupation      string         `gorm:"type:varchar(100)" json:"occupation"`
	Interests       string         `gorm:"type:text" json:"interests"` // JSON数组

	// 消费画像
	TotalSpend      float64        `gorm:"type:decimal(12,2);default:0" json:"total_spend"`
	AvgOrderValue   float64        `gorm:"type:decimal(12,2);default:0" json:"avg_order_value"`
	TotalOrders     int            `gorm:"default:0" json:"total_orders"`
	LastOrderDate   *time.Time     `json:"last_order_date"`
	PreferredPayment string        `gorm:"type:varchar(50)" json:"preferred_payment"`
	PriceSensitivity float64       `gorm:"type:decimal(5,2);default:0" json:"price_sensitivity"` // 0-1

	// 行为画像
	LoginFrequency  int            `gorm:"default:0" json:"login_frequency"`  // 近30天
	AvgSessionDuration int         `gorm:"default:0" json:"avg_session_duration"` // 秒
	ActiveDaysPerWeek int          `gorm:"default:0" json:"active_days_per_week"`
	FeatureUsageMap string         `gorm:"type:text" json:"feature_usage_map"` // JSON
	PetEngagement   float64        `gorm:"type:decimal(5,2);default:0" json:"pet_engagement"` // 0-100

	// 偏好画像
	PreferredPetType string        `gorm:"type:varchar(50)" json:"preferred_pet_type"`
	PreferredAIStyle string        `gorm:"type:varchar(50)" json:"preferred_ai_style"` // 活泼/温柔/幽默
	NotificationPref string        `gorm:"type:varchar(20)" json:"notification_pref"` // push/email/sms/none

	// 生命周期
	MemberLevel     string         `gorm:"type:varchar(20)" json:"member_level"`
	MemberSince     *time.Time     `json:"member_since"`
	ChurnRisk       float64        `gorm:"type:decimal(5,2);default:0" json:"churn_risk"` // 0-1
	LTV             float64        `gorm:"type:decimal(12,2);default:0" json:"ltv"` // Lifetime Value

	// 标签
	Tags            string         `gorm:"type:text" json:"tags"` // JSON数组
	CustomLabels    string         `gorm:"type:text" json:"custom_labels"` // JSON

	GeneratedAt     time.Time      `json:"generated_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Member360Profile) TableName() string {
	return "member_360_profiles"
}

// MemberPortraitInsight 画像洞察
type MemberPortraitInsight struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MemberID    string    `gorm:"type:varchar(64);index" json:"member_id"`
	InsightType string   `gorm:"type:varchar(50)" json:"insight_type"` // consumption/habit/preference/risk
	Title       string   `gorm:"type:varchar(200)" json:"title"`
	Description string   `gorm:"type:text" json:"description"`
	Confidence  float64  `gorm:"type:decimal(5,2)" json:"confidence"` // 0-1
	Action      string   `gorm:"type:text" json:"action"` // 建议的行动
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 表名
func (MemberPortraitInsight) TableName() string {
	return "member_portrait_insights"
}
