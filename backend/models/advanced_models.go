package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===================== 儿童模式配置 =====================

// ChildModeConfig 儿童模式配置
type ChildModeConfig struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	DeviceID     string         `gorm:"type:varchar(36);index" json:"device_id"`
	IsEnabled    bool           `gorm:"default:false" json:"is_enabled"`
	// 内容过滤
	ContentFilterLevel string `gorm:"type:varchar(20);default:'moderate'" json:"content_filter_level"` // strict/moderate/relaxed/none
	AllowedCategories  string `gorm:"type:varchar(500)" json:"allowed_categories"`                     // 逗号分隔：education,entertainment,story,music
	BlockedKeywords     string `gorm:"type:text" json:"blocked_keywords"`                                // 逗号分隔
	// 使用限制
	DailyTimeLimit  int `gorm:"default:0" json:"daily_time_limit"`  // 分钟，0表示不限制
	SessionDuration int `gorm:"default:30" json:"session_duration"`  // 单次使用时长（分钟）
	BreakDuration   int `gorm:"default:10" json:"break_duration"`    // 休息时长（分钟）
	AllowedStartTime string `gorm:"type:varchar(10);default:'08:00'" json:"allowed_start_time"` // HH:MM
	AllowedEndTime   string `gorm:"type:varchar(10);default:'20:00'" json:"allowed_end_time"`   // HH:MM
	// 使用统计
	TodayUsedMinutes  int `gorm:"default:0" json:"today_used_minutes"`
	WeekUsedMinutes   int `gorm:"default:0" json:"week_used_minutes"`
	TotalUsedMinutes  int `gorm:"default:0" json:"total_used_minutes"`
	TotalSessions     int `gorm:"default:0" json:"total_sessions"`
	LastSessionAt     *time.Time `json:"last_session_at"`
	// 家长配置
	EmergencyContact string `gorm:"type:varchar(100)" json:"emergency_contact"`
	PinCode          string `gorm:"type:varchar(10)" json:"pin_code"` // 家长验证PIN
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (ChildModeConfig) TableName() string { return "child_mode_configs" }

// ===================== 老人陪伴配置 =====================

// ElderlyCareConfig 老人陪伴配置
type ElderlyCareConfig struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	DeviceID     string    `gorm:"type:varchar(36);index" json:"device_id"`
	IsEnabled    bool      `gorm:"default:false" json:"is_enabled"`
	// 健康监测设置
	HealthMonitorEnabled bool `gorm:"default:true" json:"health_monitor_enabled"`
	HeartRateAlertHigh   int  `gorm:"default:100" json:"heart_rate_alert_high"`   // 心率上限
	HeartRateAlertLow    int  `gorm:"default:50" json:"heart_rate_alert_low"`      // 心率下限
	ActivityGoal         int  `gorm:"default:6000" json:"activity_goal"`          // 每日活动目标（步数）
	SleepMonitoring      bool `gorm:"default:true" json:"sleep_monitoring"`       // 睡眠监测
	// 用药提醒
	MedicationReminders  bool `gorm:"default:true" json:"medication_reminders"`
	MedicationTimes      string `gorm:"type:varchar(500)" json:"medication_times"`    // JSON数组，如 ["08:00","12:00","20:00"]
	MedicationNames      string `gorm:"type:varchar(500)" json:"medication_names"`     // 对应药品名称
	// 情感陪伴设置
	CompanionEnabled     bool   `gorm:"default:true" json:"companion_enabled"`
	InteractionFrequency string `gorm:"type:varchar(20);default:'normal'" json:"interaction_frequency"` // frequent/normal/sparse
	VoiceCallEnabled      bool   `gorm:"default:true" json:"voice_call_enabled"`
	FallDetectionEnabled  bool   `gorm:"default:true" json:"fall_detection_enabled"`
	// 紧急联系人
	EmergencyContactName  string `gorm:"type:varchar(100)" json:"emergency_contact_name"`
	EmergencyContactPhone string `gorm:"type:varchar(50)" json:"emergency_contact_phone"`
	EmergencyMessage      string `gorm:"type:varchar(500)" json:"emergency_message"`
	// 统计
	TotalInteractions int       `gorm:"default:0" json:"total_interactions"`
	LastInteractionAt *time.Time `json:"last_interaction_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (ElderlyCareConfig) TableName() string { return "elderly_care_configs" }

// ===================== 家庭相册 =====================

// FamilyAlbum 家庭相册照片
type FamilyAlbum struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UUID        string         `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	DeviceID    string         `gorm:"type:varchar(36);index" json:"device_id"`
	PetID       uint           `gorm:"index" json:"pet_id"`
	Title       string         `gorm:"type:varchar(255)" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	// 照片存储
	PhotoURL    string `gorm:"type:varchar(500)" json:"photo_url"`
	ThumbnailURL string `gorm:"type:varchar(500)" json:"thumbnail_url"`
	FileSize    int64  `gorm:"default:0" json:"file_size"`    // 字节
	Width       int    `gorm:"default:0" json:"width"`         // 像素
	Height      int    `gorm:"default:0" json:"height"`        // 像素
	MimeType    string `gorm:"type:varchar(50)" json:"mime_type"`
	// 分享
	IsShared    bool   `gorm:"default:false" json:"is_shared"`
	ShareToken  string `gorm:"type:varchar(64);index" json:"share_token"`
	ShareExpiry *time.Time `json:"share_expiry"`
	// 标签与分类
	Category    string `gorm:"type:varchar(50)" json:"category"`     // daily/health/event/other
	Tags        string `gorm:"type:varchar(255)" json:"tags"`        // 逗号分隔
	// 元数据
	TakenAt     *time.Time `json:"taken_at"`    // 拍摄时间
	TakenLat    float64   `gorm:"type:decimal(10,6)" json:"taken_lat"`
	TakenLng    float64   `gorm:"type:decimal(10,6)" json:"taken_lng"`
	TakenAddr   string    `gorm:"type:varchar(255)" json:"taken_addr"`
	// 点赞/评论统计
	LikeCount   int `gorm:"default:0" json:"like_count"`
	CommentCount int `gorm:"default:0" json:"comment_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (f *FamilyAlbum) BeforeCreate(tx *gorm.DB) error {
	if f.UUID == "" {
		f.UUID = uuid.New().String()
	}
	if f.ShareToken == "" {
		f.ShareToken = uuid.New().String()
	}
	return nil
}

// FamilyAlbumComment 照片评论
type FamilyAlbumComment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PhotoUUID  string         `gorm:"type:varchar(36);index" json:"photo_uuid"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	UserName   string         `gorm:"type:varchar(100)" json:"user_name"`
	Content    string         `gorm:"type:text" json:"content"`
	ParentID   uint           `gorm:"default:0" json:"parent_id"`   // 回复ID
	LikeCount  int            `gorm:"default:0" json:"like_count"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FamilyAlbumComment) TableName() string { return "family_album_comments" }

// FamilyAlbumLike 照片点赞
type FamilyAlbumLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PhotoUUID string    `gorm:"type:varchar(36);index" json:"photo_uuid"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (FamilyAlbumLike) TableName() string { return "family_album_likes" }

// ===================== 宠物疫苗接种 =====================

// PetVaccination 宠物疫苗接种记录
type PetVaccination struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PetID       uint      `gorm:"not null;index" json:"pet_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	VaccineName string    `gorm:"type:varchar(100);not null" json:"vaccine_name"` // 疫苗名称：rabies/distemper/parvo/combo/etc
	VaccineType string    `gorm:"type:varchar(50)" json:"vaccine_type"`           // 类别：core/non-core
	LotNumber   string    `gorm:"type:varchar(100)" json:"lot_number"`            // 批号
	Manufacturer string   `gorm:"type:varchar(200)" json:"manufacturer"`          // 生产商
	// 接种信息
	InoculationDate time.Time `gorm:"not null" json:"inoculation_date"`
	InoculationAge  string    `gorm:"type:varchar(50)" json:"inoculation_age"`     // 接种时月龄
	InoculationSite string    `gorm:"type:varchar(100)" json:"inoculation_site"`  // 接种部位
	Inoculator      string    `gorm:"type:varchar(100)" json:"inoculator"`         // 接种人
	VetClinic        string    `gorm:"type:varchar(200)" json:"vet_clinic"`        // 接种诊所
	// 下次接种
	NextDoseDate *time.Time `json:"next_dose_date"`
	NextDoseMemo string     `gorm:"type:varchar(255)" json:"next_dose_memo"`
	// 副作用记录
	AdverseReactions string `gorm:"type:text" json:"adverse_reactions"` // 无/轻微/严重
	AdverseDetail     string `gorm:"type:text" json:"adverse_detail"`
	// 文件
	CertificateURL string `gorm:"type:varchar(500)" json:"certificate_url"` // 证书图片URL
	Remark         string `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (PetVaccination) TableName() string { return "pet_vaccinations" }

// VaccinationReminder 疫苗提醒
type VaccinationReminder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VaccinationID uint    `gorm:"not null;index" json:"vaccination_id"`
	PetID       uint      `gorm:"not null;index" json:"pet_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	RemindAt    time.Time `gorm:"not null" json:"remind_at"`    // 提醒时间
	RemindType  string    `gorm:"type:varchar(20)" json:"remind_type"` // push/sms/email/in_app
	IsSent      bool      `gorm:"default:false" json:"is_sent"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	Memo        string    `gorm:"type:varchar(255)" json:"memo"`
	CreatedAt   time.Time `json:"created_at"`
}

func (VaccinationReminder) TableName() string { return "vaccination_reminders" }

// ===================== 宠物饮食记录 =====================

// PetDietRecord 宠物饮食记录
type PetDietRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PetID       uint      `gorm:"not null;index" json:"pet_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	DeviceID    string    `gorm:"type:varchar(36);index" json:"device_id"`
	MealType    string    `gorm:"type:varchar(20);not null" json:"meal_type"`    // breakfast/lunch/dinner/snack/supplement
	FoodName    string    `gorm:"type:varchar(200)" json:"food_name"`             // 食物名称
	FoodBrand   string    `gorm:"type:varchar(100)" json:"food_brand"`             // 品牌
	FoodType    string    `gorm:"type:varchar(50)" json:"food_type"`               // dry/wet/fresh/raw/treat
	// 食量
	Amount      float64   `gorm:"type:decimal(10,2)" json:"amount"`               // 数量
	AmountUnit  string    `gorm:"type:varchar(20);default:'g'" json:"amount_unit"` // g/kg/cup/piece
	Calories    int       `gorm:"default:0" json:"calories"`                     // 估算卡路里
	// 喂食方式
	FeedingMethod string  `gorm:"type:varchar(50)" json:"feeding_method"`         // 自动/手动/混合
	AutoFeederID  string  `gorm:"type:varchar(36)" json:"auto_feeder_id"`          // 自动喂食器设备ID
	// 时间
	EatTime     *time.Time `json:"eat_time"`    // 实际进食时间
	PlannedTime *time.Time `json:"planned_time"` // 计划时间
	Duration    int        `gorm:"default:0" json:"duration"`  // 进食时长（分钟）
	// 状态
	Status      string    `gorm:"type:varchar(20);default:'completed'" json:"status"` // planned/completed/partial/skipped
	LeftOver    float64   `gorm:"type:decimal(10,2);default:0" json:"left_over"`      // 剩余量
	// 健康关联
	AppetiteScore int     `gorm:"default:5" json:"appetite_score"` // 食欲评分 1-10
	HealthNote    string  `gorm:"type:text" json:"health_note"`
	// 照片
	PhotoURL    string    `gorm:"type:varchar(500)" json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PetDietRecord) TableName() string { return "pet_diet_records" }

// ===================== 寻宠报告增强 =====================

// PetFinderReport 寻宠报告（增强版，对应 PetLostReport）
type PetFinderReport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PetID          uint      `gorm:"not null;index" json:"pet_id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	ReportType     string    `gorm:"type:varchar(20);not null" json:"report_type"` // lost/found/theft
	Title          string    `gorm:"type:varchar(255);not null" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	LastSeenAt     *time.Time `json:"last_seen_at"`
	LastSeenLat    float64   `gorm:"type:decimal(10,6)" json:"last_seen_lat"`
	LastSeenLng    float64   `gorm:"type:decimal(10,6)" json:"last_seen_lng"`
	LastSeenAddr   string    `gorm:"type:varchar(500)" json:"last_seen_addr"`
	ContactName    string    `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone   string    `gorm:"type:varchar(50)" json:"contact_phone"`
	Reward         float64   `gorm:"type:decimal(10,2)" json:"reward"`
	RewardMemo     string    `gorm:"type:varchar(255)" json:"reward_memo"`
	Photos         StringArray `gorm:"type:varchar(500)[]" json:"photos"`
	Status         string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active/found/closed
	ResolvedAt     *time.Time `json:"resolved_at"`
	ViewCount      int       `gorm:"default:0" json:"view_count"`
	AlertRadius    float64   `gorm:"type:decimal(10,2);default:5.0" json:"alert_radius"` // km
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (PetFinderReport) TableName() string { return "pet_finder_reports" }

// PetFinderSighting 寻宠目击报告
type PetFinderSighting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReportID    uint      `gorm:"not null;index" json:"report_id"`
	ReporterID  uint      `gorm:"not null;index" json:"reporter_id"`
	SightedAt   *time.Time `json:"sighted_at"`
	SightedLat  float64   `gorm:"type:decimal(10,6)" json:"sighted_lat"`
	SightedLng  float64   `gorm:"type:decimal(10,6)" json:"sighted_lng"`
	SightedAddr string    `gorm:"type:varchar(500)" json:"sighted_addr"`
	Description string    `gorm:"type:text" json:"description"`
	PetStatus   string    `gorm:"type:varchar(50)" json:"pet_status"`  // alone/scared/hungry/with_owner
	PhotoURL    string    `gorm:"type:varchar(500)" json:"photo_url"`
	ContactName string    `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string   `gorm:"type:varchar(50)" json:"contact_phone"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
}

func (PetFinderSighting) TableName() string { return "pet_finder_sightings" }

// PetFinderAlert 寻宠网络推送
type PetFinderAlert struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ReportID   uint      `gorm:"not null;index" json:"report_id"`
	AlertType  string    `gorm:"type:varchar(20)" json:"alert_type"`  // push/email/sms
	RecipientID uint     `gorm:"index" json:"recipient_id"`
	SentAt     *time.Time `json:"sent_at"`
	Status     string    `gorm:"type:varchar(20)" json:"status"`  // pending/sent/failed
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PetFinderAlert) TableName() string { return "pet_finder_alerts" }
