package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Household 家庭/户
type Household struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	HouseholdUUID string       `gorm:"type:varchar(64);uniqueIndex;not null" json:"household_uuid"`
	Name        string         `gorm:"type:varchar(64);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	AvatarURL   string         `gorm:"type:varchar(512)" json:"avatar_url"`
	OwnerID     uint           `gorm:"index;not null" json:"owner_id"`       // 户主用户ID
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`     // 租户ID
	Region      string         `gorm:"type:varchar(64)" json:"region"`       // 地区
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active, archived
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (h *Household) BeforeCreate(tx *gorm.DB) error {
	if h.HouseholdUUID == "" {
		h.HouseholdUUID = uuid.New().String()
	}
	return nil
}

// HouseholdMember 家庭成员
type HouseholdMember struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MemberUUID     string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"member_uuid"`
	HouseholdID    uint           `gorm:"index;not null" json:"household_id"`
	UserID         uint           `gorm:"index;not null" json:"user_id"`
	Nickname       string         `gorm:"type:varchar(32)" json:"nickname"`        // 家庭内昵称
	Role           string         `gorm:"type:varchar(32);default:'member'" json:"role"` // owner, admin, member, guest
	Relationship   string         `gorm:"type:varchar(32)" json:"relationship"`   // 关系：father, mother, child, grandparent, other
	AvatarURL      string         `gorm:"type:varchar(512)" json:"avatar_url"`
	InviteCode     string         `gorm:"type:varchar(64);index" json:"invite_code"` // 邀请码
	InviteStatus   string         `gorm:"type:varchar(20);default:'active'" json:"invite_status"` // pending, active, removed
	JoinedAt       *time.Time     `json:"joined_at"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	Permissions    string         `gorm:"type:varchar(512)" json:"permissions"` // JSON数组字符串
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (h *HouseholdMember) BeforeCreate(tx *gorm.DB) error {
	if h.MemberUUID == "" {
		h.MemberUUID = uuid.New().String()
	}
	return nil
}

// Interaction 用户与宠物交互记录
type Interaction struct {
	ID           uint                   `gorm:"primaryKey" json:"id"`
	UUID         string                 `gorm:"type:varchar(64);uniqueIndex;not null" json:"uuid"`
	PetID        uint                   `gorm:"index;not null" json:"pet_id"`
	DeviceID     string                 `gorm:"type:varchar(64);index" json:"device_id"`
	UserID       uint                   `gorm:"index;not null" json:"user_id"`
	HouseholdID  *uint                  `gorm:"index" json:"household_id"`
	Type         string                 `gorm:"type:varchar(32);not null" json:"type"`            // play, feed, walk, talk, pet, teach, other
	Action       string                 `gorm:"type:varchar(64)" json:"action"`                  // 具体动作
	Content      string                 `gorm:"type:text" json:"content"`                        // 交互内容/对话
	Duration     int                    `gorm:"type:int;default:0" json:"duration"`             // 持续时长（秒）
	Emotion      string                 `gorm:"type:varchar(32)" json:"emotion"`                // 交互后的情绪
	EmotionScore int                    `gorm:"type:smallint" json:"emotion_score"`             // 情绪评分 1-10
	Effect       string                 `gorm:"type:varchar(32)" json:"effect"`                  // 效果：positive, neutral, negative
	Metadata     map[string]interface{} `gorm:"type:jsonb" json:"metadata"`                    // 附加数据
	ImageURLs    string                 `gorm:"type:varchar(1024)" json:"image_urls"`          // 图片URLs，逗号分隔
	TenantID     string                 `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// BeforeCreate 创建前自动生成 UUID
func (i *Interaction) BeforeCreate(tx *gorm.DB) error {
	if i.UUID == "" {
		i.UUID = uuid.New().String()
	}
	return nil
}

// FamilyAlbum 家庭相册
type FamilyAlbum struct {
	ID          uint                   `gorm:"primaryKey" json:"id"`
	UUID        string                 `gorm:"type:varchar(64);uniqueIndex;not null" json:"uuid"`
	HouseholdID uint                   `gorm:"index;not null" json:"household_id"`
	PetID       *uint                  `gorm:"index" json:"pet_id"`
	UploaderID  uint                   `gorm:"index;not null" json:"uploader_id"`
	Title       string                 `gorm:"type:varchar(128)" json:"title"`
	Description string                 `gorm:"type:text" json:"description"`
	ImageURL    string                 `gorm:"type:varchar(512);not null" json:"image_url"`
	ThumbnailURL string                `gorm:"type:varchar(512)" json:"thumbnail_url"`
	Category    string                 `gorm:"type:varchar(32)" json:"category"` // daily, milestone, event, pet, family
	Tags        string                 `gorm:"type:varchar(256)" json:"tags"`   // 逗号分隔
	Metadata    map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	FileSize    int64                  `gorm:"type:bigint" json:"file_size"`  // 文件大小（字节）
	Width       int                    `gorm:"type:int" json:"width"`         // 图片宽度
	Height      int                    `gorm:"type:int" json:"height"`        // 图片高度
	TenantID    string                 `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (f *FamilyAlbum) BeforeCreate(tx *gorm.DB) error {
	if f.UUID == "" {
		f.UUID = uuid.New().String()
	}
	return nil
}

// HouseholdSettings 家庭设置
type HouseholdSettings struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	HouseholdID uint           `gorm:"uniqueIndex;not null" json:"household_id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`

	// 家庭模式
	Mode        string         `gorm:"type:varchar(32);default:'normal'" json:"mode"` // normal, child, elder, pet

	// 儿童模式设置
	ChildModeEnabled bool         `gorm:"default:false" json:"child_mode_enabled"`
	ChildAgeRange    string       `gorm:"type:varchar(16)" json:"child_age_range"` // e.g. "3-6"
	ChildContentFilter string     `gorm:"type:varchar(32)" json:"child_content_filter"` // strict, moderate, light
	ChildScreenTime  int          `gorm:"type:int;default:60" json:"child_screen_time"` // 每日屏幕时间（分钟）
	ChildAllowedActions string    `gorm:"type:varchar(512)" json:"child_allowed_actions"` // JSON数组

	// 老人陪伴模式设置
	ElderModeEnabled bool         `gorm:"default:false" json:"elder_mode_enabled"`
	ElderCareLevel   string       `gorm:"type:varchar(16)" json:"elder_care_level"` // basic, standard, premium
	ElderReminderEnabled bool     `gorm:"default:true" json:"elder_reminder_enabled"`
	ElderMedicationReminders bool `gorm:"default:true" json:"elder_medication_reminders"`
	ElderEmergencyContact string  `gorm:"type:varchar(128)" json:"elder_emergency_contact"`

	// 宠物模式设置
	PetModeEnabled bool           `gorm:"default:true" json:"pet_mode_enabled"`
	PetInteractionLevel string  `gorm:"type:varchar(16)" json:"pet_interaction_level"` // low, medium, high

	// 通知设置
	NotificationEnabled bool      `gorm:"default:true" json:"notification_enabled"`
	NotificationSettings string  `gorm:"type:varchar(512)" json:"notification_settings"` // JSON

	// 隐私设置
	PrivacyLevel string           `gorm:"type:varchar(16);default:'family'" json:"privacy_level"` // family, private

	// 其他设置
	Metadata   map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}
