package models

import (
	"time"

	"gorm.io/gorm"
)

// ChildModeSettings 儿童模式设置
type ChildModeSettings struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	IsEnabled     bool           `gorm:"default:false" json:"is_enabled"`
	// 使用时间限制
	DailyTimeLimitMinutes int       `gorm:"default:60" json:"daily_time_limit_minutes"` // 每日使用时长限制（分钟），0表示不限制
	SingleSessionMinutes  int       `gorm:"default:30" json:"single_session_minutes"`  // 单次使用时长限制
	// 允许使用的功能
	AllowVoiceChat    bool `gorm:"default:true" json:"allow_voice_chat"`
	AllowVideoChat    bool `gorm:"default:false" json:"allow_video_chat"`
	AllowAppInstall   bool `gorm:"default:false" json:"allow_app_install"`
	AllowAppUninstall bool `gorm:"default:false" json:"allow_app_uninstall"`
	AllowGameMode     bool `gorm:"default:true" json:"allow_game_mode"`
	AllowSocialFeatures bool `gorm:"default:false" json:"allow_social_features"`
	AllowInAppPurchases bool `gorm:"default:false" json:"allow_in_app_purchases"`
	// 内容过滤
	ContentFilterLevel string `gorm:"type:varchar(20);default:'moderate'" json:"content_filter_level"` // strict/moderate/lenient/none
	// 使用时间段限制
	AllowedStartTime string `gorm:"type:varchar(8);default:'08:00'" json:"allowed_start_time"` // HH:MM
	AllowedEndTime   string `gorm:"type:varchar(8);default:'20:00'" json:"allowed_end_time"`   // HH:MM
	// 休息提醒
	RestReminderMinutes int `gorm:"default:30" json:"rest_reminder_minutes"` // 每隔多少分钟提醒休息
	// 家长监督
	ParentPin        string `gorm:"type:varchar(8)" json:"parent_pin"`         // 家长PIN码
	MonitorPetID     uint   `json:"monitor_pet_id"`                            // 监督的宠物ID
	WeeklyReportEnabled bool `gorm:"default:true" json:"weekly_report_enabled"`
	// 统计
	TodayUsedMinutes int            `gorm:"-" json:"today_used_minutes"` // 不存储，计算得出
	WeekUsedMinutes  int            `gorm:"-" json:"week_used_minutes"`  // 不存储，计算得出
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (ChildModeSettings) TableName() string {
	return "child_mode_settings"
}
