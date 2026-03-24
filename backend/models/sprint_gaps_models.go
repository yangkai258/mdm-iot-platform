package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===================== Sprint 26 Phase 4: 家庭相册增强 =====================

// FamilyAlbumContainer 家庭相册容器（相册夹）
type FamilyAlbumContainer struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UUID        string         `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	TenantID    string         `gorm:"type:varchar(36);index" json:"tenant_id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(500)" json:"description"`
	CoverURL    string         `gorm:"type:varchar(500)" json:"cover_url"`
	Visibility  string         `gorm:"type:varchar(20);default:'private'" json:"visibility"` // private, family, public
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	ItemCount   int            `gorm:"default:0" json:"item_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (f *FamilyAlbumContainer) BeforeCreate(tx *gorm.DB) error {
	if f.UUID == "" {
		f.UUID = uuid.New().String()
	}
	return nil
}

// FamilyAlbumItem 家庭相册项（照片/视频）
type FamilyAlbumItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UUID        string         `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	TenantID    string         `gorm:"type:varchar(36);index" json:"tenant_id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	ContainerID uint           `gorm:"not null;index" json:"container_id"`
	PetID       uint           `gorm:"index" json:"pet_id"`
	DeviceID    string         `gorm:"type:varchar(36);index" json:"device_id"`
	MediaType   string         `gorm:"type:varchar(20);not null" json:"media_type"` // photo, video
	Title       string         `gorm:"type:varchar(255)" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	// 存储信息
	URL         string `gorm:"type:varchar(500)" json:"url"`
	ThumbnailURL string `gorm:"type:varchar(500)" json:"thumbnail_url"`
	FileSize    int64  `gorm:"default:0" json:"file_size"`
	Width       int    `gorm:"default:0" json:"width"`
	Height      int    `gorm:"default:0" json:"height"`
	MimeType    string `gorm:"type:varchar(50)" json:"mime_type"`
	Duration    int    `gorm:"default:0" json:"duration"` // 视频时长（秒）
	// 分享
	IsShared    bool   `gorm:"default:false" json:"is_shared"`
	ShareToken  string `gorm:"type:varchar(64);index" json:"share_token"`
	ShareExpiry *time.Time `json:"share_expiry"`
	// 元数据
	Category    string `gorm:"type:varchar(50)" json:"category"` // daily, health, event, milestone, other
	Tags        string `gorm:"type:varchar(255)" json:"tags"`
	TakenAt     *time.Time `json:"taken_at"`
	TakenLat    float64 `gorm:"type:decimal(10,6)" json:"taken_lat"`
	TakenLng    float64 `gorm:"type:decimal(10,6)" json:"taken_lng"`
	TakenAddr   string `gorm:"type:varchar(255)" json:"taken_addr"`
	// 统计
	LikeCount   int `gorm:"default:0" json:"like_count"`
	CommentCount int `gorm:"default:0" json:"comment_count"`
	ViewCount   int `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (f *FamilyAlbumItem) BeforeCreate(tx *gorm.DB) error {
	if f.UUID == "" {
		f.UUID = uuid.New().String()
	}
	if f.ShareToken == "" {
		f.ShareToken = uuid.New().String()
	}
	return nil
}

// ===================== Sprint 26 Phase 4: 儿童模式设置 =====================

// ChildModeSettings 儿童模式设置
type ChildModeSettings struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:varchar(36);index" json:"tenant_id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	DeviceID     string         `gorm:"type:varchar(36);index" json:"device_id"`
	PetID        uint           `gorm:"index" json:"pet_id"`
	IsEnabled    bool           `gorm:"default:false" json:"is_enabled"`
	// 内容过滤
	ContentFilterLevel string `gorm:"type:varchar(20);default:'moderate'" json:"content_filter_level"` // strict, moderate, relaxed, none
	AllowedCategories  string `gorm:"type:varchar(500)" json:"allowed_categories"` // education, entertainment, story, music, game
	BlockedKeywords    string `gorm:"type:text" json:"blocked_keywords"`
	BlockedApps        string `gorm:"type:varchar(500)" json:"blocked_apps"` // 逗号分隔
	// 使用时间限制
	DailyTimeLimit   int    `gorm:"default:0" json:"daily_time_limit"`    // 分钟，0表示不限制
	SessionDuration   int    `gorm:"default:30" json:"session_duration"`   // 单次使用时长（分钟）
	BreakDuration     int    `gorm:"default:10" json:"break_duration"`     // 休息时长（分钟）
	AllowedStartTime  string `gorm:"type:varchar(10);default:'08:00'" json:"allowed_start_time"` // HH:MM
	AllowedEndTime    string `gorm:"type:varchar(10);default:'20:00'" json:"allowed_end_time"`   // HH:MM
	// 使用统计
	TodayUsedMinutes  int  `gorm:"default:0" json:"today_used_minutes"`
	WeekUsedMinutes   int  `gorm:"default:0" json:"week_used_minutes"`
	TotalUsedMinutes  int  `gorm:"default:0" json:"total_used_minutes"`
	TotalSessions     int  `gorm:"default:0" json:"total_sessions"`
	LastSessionAt     *time.Time `json:"last_session_at"`
	// 家长控制
	EmergencyContact string `gorm:"type:varchar(100)" json:"emergency_contact"`
	PinCode         string `gorm:"type:varchar(10)" json:"pin_code"` // 家长验证PIN
	ParentPhone     string `gorm:"type:varchar(50)" json:"parent_phone"`
	// 权限
	CanShareContent  bool `gorm:"default:false" json:"can_share_content"`
	CanDownloadContent bool `gorm:"default:false" json:"can_download_content"`
	AllowCamera      bool `gorm:"default:true" json:"allow_camera"`
	AllowMicrophone  bool `gorm:"default:true" json:"allow_microphone"`
	// 状态
	Status          string `gorm:"type:varchar(20);default:'active'" json:"status"` // active, paused, expired
	ExpiresAt       *time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// ===================== Sprint 17: 语音情绪记录 =====================

// VoiceEmotionRecord 语音情绪记录
type VoiceEmotionRecord struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:varchar(36);index" json:"tenant_id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	PetID        uint           `gorm:"index" json:"pet_id"`
	DeviceID     string         `gorm:"type:varchar(36);index" json:"device_id"`
	SessionID    string         `gorm:"type:varchar(64);index" json:"session_id"`
	// 语音数据
	AudioURL     string `gorm:"type:varchar(500)" json:"audio_url"`
	AudioDuration int   `gorm:"default:0" json:"audio_duration"` // 秒
	FileSize     int64  `gorm:"default:0" json:"file_size"`
	// 情绪分析结果
	EmotionType  string  `gorm:"type:varchar(32)" json:"emotion_type"` // happy, sad, angry, calm, excited, anxious, fearful, neutral
	EmotionScore float64 `gorm:"type:decimal(5,2)" json:"emotion_score"` // 置信度 0-100
	Intensity   int     `gorm:"type:smallint" json:"intensity"` // 1-10
	Valence      float64 `gorm:"type:decimal(5,2)" json:"valence"` // 情绪效价 -1到1
	Arousal      float64 `gorm:"type:decimal(5,2)" json:"arousal"` // 激活度 -1到1
	// 语音特征
	PitchAvg     float64 `gorm:"type:decimal(10,4)" json:"pitch_avg"`
	PitchVar      float64 `gorm:"type:decimal(10,4)" json:"pitch_var"`
	SpeechRate    float64 `gorm:"type:decimal(10,4)" json:"speech_rate"` // 字/秒
	EnergyLevel   float64 `gorm:"type:decimal(10,4)" json:"energy_level"`
	// 语义分析
	Transcript   string `gorm:"type:text" json:"transcript"` // 语音转文字
	Keywords      string `gorm:"type:varchar(500)" json:"keywords"` // 关键词
	Sentiment     string `gorm:"type:varchar(32)" json:"sentiment"` // positive, negative, neutral
	// 上下文
	Context       string `gorm:"type:varchar(255)" json:"context"` // 触发场景
	TriggerPhrase string `gorm:"type:varchar(255)" json:"trigger_phrase"` // 触发词
	// AI响应
	AIResponse    string `gorm:"type:text" json:"ai_response"` // AI的响应内容
	ResponseType  string `gorm:"type:varchar(32)" json:"response_type"` // comfort, encourage, play, quiet, none
	// 状态
	Status        string `gorm:"type:varchar(20);default:'completed'" json:"status"` // processing, completed, failed
	ErrorMessage  string `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
