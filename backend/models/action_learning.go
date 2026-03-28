package models

import (
	"time"

	"gorm.io/gorm"
)

// ActionLearningProgress 动作模仿学习进度
type ActionLearningProgress struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PetID        string         `gorm:"type:varchar(64);index" json:"pet_id"`
	MemberID     string         `gorm:"type:varchar(64);index" json:"member_id"`
	ActionID     uint           `gorm:"not null" json:"action_id"`
	ActionName   string         `gorm:"type:varchar(100)" json:"action_name"`

	// 学习进度
	Level         int            `gorm:"default:1" json:"level"`            // 当前等级 1-10
	Exp           int            `gorm:"default:0" json:"exp"`            // 当前经验值
	ExpToNextLevel int           `gorm:"default:100" json:"exp_to_next"`   // 升级所需经验
	MasteryRate   float64        `gorm:"type:decimal(5,2);default:0" json:"mastery_rate"` // 掌握度 0-100%

	// 练习数据
	PracticeCount  int            `gorm:"default:0" json:"practice_count"`  // 练习次数
	SuccessCount   int            `gorm:"default:0" json:"success_count"`   // 成功次数
	TotalDuration  int            `gorm:"default:0" json:"total_duration"` // 总练习时长(秒)
	AvgAccuracy    float64        `gorm:"type:decimal(5,2);default:0" json:"avg_accuracy"` // 平均准确率

	// 时间追踪
	LastPracticeAt *time.Time     `json:"last_practice_at"`
	NextPracticeAt *time.Time     `json:"next_practice_at"` // 推荐下次练习时间
	LearnedAt      *time.Time     `json:"learned_at"`       // 学会时间
	MasteredAt     *time.Time     `json:"mastered_at"`      // 精通时间

	// 状态
	Status        string         `gorm:"type:varchar(20);default:'learning'" json:"status"` // learning/practicing/mastered
	StarRating    int            `gorm:"default:0" json:"star_rating"`      // 星级评分 1-5

	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (ActionLearningProgress) TableName() string {
	return "action_learning_progress"
}

// ActionLearningSession 单次练习记录
type ActionLearningSession struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ProgressID   uint           `gorm:"not null;index" json:"progress_id"`
	PetID        string         `gorm:"type:varchar(64)" json:"pet_id"`
	ActionID     uint           `gorm:"not null" json:"action_id"`
	ActionName   string         `gorm:"type:varchar(100)" json:"action_name"`

	Duration      int            `gorm:"default:0" json:"duration"`       // 练习时长(秒)
	Accuracy      float64        `gorm:"type:decimal(5,2)" json:"accuracy"` // 准确率 0-100
	IsSuccess     bool           `gorm:"default:false" json:"is_success"`
	Score         int            `gorm:"default:0" json:"score"`         // 得分
	ExpGained     int            `gorm:"default:0" json:"exp_gained"`    // 获得经验

	Feedback      string         `gorm:"type:text" json:"feedback"`      // AI反馈
	VideoURL      string         `gorm:"type:varchar(512)" json:"video_url"` // 练习视频URL
	ScreenshotURL string         `gorm:"type:varchar(512)" json:"screenshot_url"`

	SessionDate   time.Time      `json:"session_date"`
	CreatedAt     time.Time      `json:"created_at"`
}

// TableName 表名
func (ActionLearningSession) TableName() string {
	return "action_learning_sessions"
}
