package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============ 儿童模式 ============

// ChildrenProfile 儿童档案
type ChildrenProfile struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ProfileID    string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"profile_id"`
	DeviceID     string         `gorm:"type:varchar(36);index" json:"device_id"`      // 绑定的设备ID
	Name         string         `gorm:"type:varchar(64);not null" json:"name"`        // 儿童姓名
	Nickname     string         `gorm:"type:varchar(64)" json:"nickname"`             // 昵称
	Avatar       string         `gorm:"type:varchar(500)" json:"avatar"`              // 头像URL
	BirthDate    *time.Time     `json:"birth_date"`                                   // 出生日期
	Gender       string         `gorm:"type:varchar(10)" json:"gender"`               // 性别
	Age          int            `json:"age"`                                          // 年龄（计算字段）
	ParentUserID string         `gorm:"type:varchar(36);index" json:"parent_user_id"` // 家长用户ID
	OrgID        uint           `gorm:"index" json:"org_id"`                          // 组织ID
	Status       int            `gorm:"type:smallint;default:1" json:"status"`        // 状态：1正常 0禁用
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (c *ChildrenProfile) BeforeCreate(tx *gorm.DB) error {
	if c.ProfileID == "" {
		c.ProfileID = uuid.New().String()
	}
	return nil
}

// ContentFilterRule 内容过滤规则
type ContentFilterRule struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	ProfileID         string         `gorm:"type:varchar(36);index;not null" json:"profile_id"` // 儿童档案ID
	RuleID            string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"rule_id"`
	RuleName          string         `gorm:"type:varchar(128)" json:"rule_name"`          // 规则名称
	FilterLevel       int            `gorm:"type:smallint;default:2" json:"filter_level"` // 过滤级别：1严格 2普通 3宽松
	BlockAdult        bool           `gorm:"default:true" json:"block_adult"`             // 屏蔽成人内容
	BlockViolence     bool           `gorm:"default:true" json:"block_violence"`          // 屏蔽暴力内容
	BlockGambling     bool           `gorm:"default:true" json:"block_gambling"`          // 屏蔽赌博内容
	BlockAds          bool           `gorm:"default:true" json:"block_ads"`               // 屏蔽广告
	BlockGames        bool           `gorm:"default:false" json:"block_games"`            // 屏蔽游戏
	AllowedCategories []string       `gorm:"type:jsonb" json:"allowed_categories"`        // 允许的分类
	BlockedKeywords   []string       `gorm:"type:jsonb" json:"blocked_keywords"`          // 屏蔽关键词
	AllowedApps       []string       `gorm:"type:jsonb" json:"allowed_apps"`              // 允许的应用
	BlockedApps       []string       `gorm:"type:jsonb" json:"blocked_apps"`              // 屏蔽的应用
	WhitelistMode     bool           `gorm:"default:false" json:"whitelist_mode"`         // 白名单模式
	Enabled           bool           `gorm:"default:true" json:"enabled"`                 // 是否启用
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (c *ContentFilterRule) BeforeCreate(tx *gorm.DB) error {
	if c.RuleID == "" {
		c.RuleID = uuid.New().String()
	}
	return nil
}

// UsageLimit 使用限制
type UsageLimit struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	ProfileID              string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"profile_id"` // 儿童档案ID
	LimitID                string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"limit_id"`
	DailyTimeLimit         int            `gorm:"type:int;default:120" json:"daily_time_limit"`              // 每日使用时长限制（分钟）
	WeeklyTimeLimit        int            `gorm:"type:int;default:480" json:"weekly_time_limit"`             // 每周使用时长限制（分钟）
	MaxSingleSession       int            `gorm:"type:int;default:30" json:"max_single_session"`             // 单次最大使用时长（分钟）
	AllowedStartTime       string         `gorm:"type:varchar(8);default:'08:00'" json:"allowed_start_time"` // 允许使用开始时间
	AllowedEndTime         string         `gorm:"type:varchar(8);default:'22:00'" json:"allowed_end_time"`   // 允许使用结束时间
	AllowedDays            []int          `gorm:"type:jsonb" json:"allowed_days"`                            // 允许使用的星期：1-7
	BreakInterval          int            `gorm:"type:int;default:10" json:"break_interval"`                 // 休息间隔（分钟）
	BreakDuration          int            `gorm:"type:int;default:5" json:"break_duration"`                  // 休息时长（分钟）
	EyeProtectionEnabled   bool           `gorm:"default:true" json:"eye_protection_enabled"`                // 护眼模式
	EyeProtectionInterval  int            `gorm:"type:int;default:30" json:"eye_protection_interval"`        // 护眼提醒间隔（分钟）
	PostureReminderEnabled bool           `gorm:"default:true" json:"posture_reminder_enabled"`              // 坐姿提醒
	Enabled                bool           `gorm:"default:true" json:"enabled"`                               // 是否启用
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (u *UsageLimit) BeforeCreate(tx *gorm.DB) error {
	if u.LimitID == "" {
		u.LimitID = uuid.New().String()
	}
	return nil
}

// ============ 老人陪伴模式 ============

// ElderlyProfile 老人档案
type ElderlyProfile struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ProfileID        string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"profile_id"`
	DeviceID         string         `gorm:"type:varchar(36);index" json:"device_id"`         // 绑定的设备ID
	Name             string         `gorm:"type:varchar(64);not null" json:"name"`           // 老人姓名
	Nickname         string         `gorm:"type:varchar(64)" json:"nickname"`                // 昵称
	Avatar           string         `gorm:"type:varchar(500)" json:"avatar"`                 // 头像URL
	BirthDate        *time.Time     `json:"birth_date"`                                      // 出生日期
	Gender           string         `gorm:"type:varchar(10)" json:"gender"`                  // 性别
	Age              int            `json:"age"`                                             // 年龄
	Phone            string         `gorm:"type:varchar(20)" json:"phone"`                   // 联系电话
	EmergencyContact string         `gorm:"type:varchar(64)" json:"emergency_contact"`       // 紧急联系人
	EmergencyPhone   string         `gorm:"type:varchar(20)" json:"emergency_phone"`         // 紧急联系电话
	Address          string         `gorm:"type:varchar(500)" json:"address"`                // 住址
	MedicalHistory   string         `gorm:"type:text" json:"medical_history"`                // 病史
	Allergies        string         `gorm:"type:text" json:"allergies"`                      // 过敏信息
	BloodType        string         `gorm:"type:varchar(10)" json:"blood_type"`              // 血型
	CaregiverUserID  string         `gorm:"type:varchar(36);index" json:"caregiver_user_id"` // 照护者用户ID
	OrgID            uint           `gorm:"index" json:"org_id"`                             // 组织ID
	Status           int            `gorm:"type:smallint;default:1" json:"status"`           // 状态：1正常 0禁用
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (e *ElderlyProfile) BeforeCreate(tx *gorm.DB) error {
	if e.ProfileID == "" {
		e.ProfileID = uuid.New().String()
	}
	return nil
}

// HealthMonitorSetting 健康监控设置
type HealthMonitorSetting struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	ProfileID             string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"profile_id"` // 老人档案ID
	SettingID             string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"setting_id"`
	HeartRateMonitor      bool           `gorm:"default:true" json:"heart_rate_monitor"`                   // 心率监控
	HeartRateMin          int            `gorm:"type:smallint;default:50" json:"heart_rate_min"`           // 心率下限
	HeartRateMax          int            `gorm:"type:smallint;default:120" json:"heart_rate_max"`          // 心率上限
	BPMonitor             bool           `gorm:"default:true" json:"bp_monitor"`                           // 血压监控
	BPSystolicMin         int            `gorm:"type:smallint;default:90" json:"bp_systolic_min"`          // 收缩压下限
	BPSystolicMax         int            `gorm:"type:smallint;default:140" json:"bp_systolic_max"`         // 收缩压上限
	BPDiastolicMin        int            `gorm:"type:smallint;default:60" json:"bp_diastolic_min"`         // 舒张压下限
	BPDiastolicMax        int            `gorm:"type:smallint;default:90" json:"bp_diastolic_max"`         // 舒张压上限
	SleepMonitor          bool           `gorm:"default:true" json:"sleep_monitor"`                        // 睡眠监控
	ActivityMonitor       bool           `gorm:"default:true" json:"activity_monitor"`                     // 活动监控
	StepGoal              int            `gorm:"type:int;default:6000" json:"step_goal"`                   // 每日步数目标
	FallDetection         bool           `gorm:"default:true" json:"fall_detection"`                       // 跌倒检测
	FallAlertEnabled      bool           `gorm:"default:true" json:"fall_alert_enabled"`                   // 跌倒告警
	EmergencyAlertEnabled bool           `gorm:"default:true" json:"emergency_alert_enabled"`              // 紧急告警
	AlertThreshold        int            `gorm:"type:int;default:3" json:"alert_threshold"`                // 告警阈值（连续异常次数）
	CheckInterval         int            `gorm:"type:int;default:60" json:"check_interval"`                // 健康检查间隔（分钟）
	ReportFrequency       string         `gorm:"type:varchar(20);default:'daily'" json:"report_frequency"` // 报告频率：daily, weekly, monthly
	Enabled               bool           `gorm:"default:true" json:"enabled"`                              // 是否启用
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (h *HealthMonitorSetting) BeforeCreate(tx *gorm.DB) error {
	if h.SettingID == "" {
		h.SettingID = uuid.New().String()
	}
	return nil
}

// ElderlyReminder 老人提醒
type ElderlyReminder struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ReminderID     string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"reminder_id"`
	ProfileID      string         `gorm:"type:varchar(36);index;not null" json:"profile_id"`     // 老人档案ID
	Title          string         `gorm:"type:varchar(128);not null" json:"title"`               // 提醒标题
	Content        string         `gorm:"type:varchar(500)" json:"content"`                      // 提醒内容
	ReminderType   string         `gorm:"type:varchar(32);not null" json:"reminder_type"`        // 类型：medication, exercise, appointment, custom
	MedicineName   string         `gorm:"type:varchar(128)" json:"medicine_name"`                // 药品名称（用药提醒时）
	MedicineDosage string         `gorm:"type:varchar(64)" json:"medicine_dosage"`               // 药品剂量
	ScheduleType   string         `gorm:"type:varchar(20);default:'daily'" json:"schedule_type"` // 调度类型：once, daily, weekly, monthly
	ScheduleTime   string         `gorm:"type:varchar(8);not null" json:"schedule_time"`         // 提醒时间（HH:MM）
	ScheduleDays   []int          `gorm:"type:jsonb" json:"schedule_days"`                       // 周期：星期几（weekly时使用）
	ScheduleDates  []int          `gorm:"type:jsonb" json:"schedule_dates"`                      // 日期（monthly时使用）
	AdvanceNotice  int            `gorm:"type:int;default:0" json:"advance_notice"`              // 提前通知分钟数
	RepeatCount    int            `gorm:"type:int;default:1" json:"repeat_count"`                // 重复次数
	RepeatInterval int            `gorm:"type:int;default:0" json:"repeat_interval"`             // 重复间隔（分钟）
	IsCompleted    bool           `gorm:"default:false" json:"is_completed"`                     // 是否完成
	CompletedAt    *time.Time     `json:"completed_at"`                                          // 完成时间
	IsEnabled      bool           `gorm:"default:true" json:"is_enabled"`                        // 是否启用
	DeviceID       string         `gorm:"type:varchar(36)" json:"device_id"`                     // 关联设备ID
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (e *ElderlyReminder) BeforeCreate(tx *gorm.DB) error {
	if e.ReminderID == "" {
		e.ReminderID = uuid.New().String()
	}
	return nil
}
