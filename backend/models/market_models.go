package models

import (
	"time"

	"gorm.io/gorm"
)

// ===== Emoticon Pack =====

// EmoticonPack emoticon pack model
type EmoticonPack struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	PackName     string         `gorm:"type:varchar(255);not null" json:"pack_name"`
	PackType     string         `gorm:"type:varchar(20);not null" json:"pack_type"` // built-in/official/creator/seasonal
	Description  string         `gorm:"type:text" json:"description"`
	ThumbnailURL string         `gorm:"type:varchar(500)" json:"thumbnail_url"`
	PreviewURL   string         `gorm:"type:varchar(500)" json:"preview_url"`
	Emoticons    JSON           `gorm:"type:jsonb" json:"emoticons"`
	Price        float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	IsFree       bool           `gorm:"default:true" json:"is_free"`
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/pending/approved/rejected/published/removed
	Downloads    int            `gorm:"default:0" json:"downloads"`
	RatingAvg    float64        `gorm:"type:decimal(3,2);default:0" json:"rating_avg"`
	RatingCount  int            `gorm:"default:0" json:"rating_count"`
	Tags         StringArray    `gorm:"type:text" json:"tags"`
	ReviewedAt   *time.Time     `json:"reviewed_at"`
	ReviewedBy   *uint          `gorm:"index" json:"reviewed_by"`
	PublishedAt  *time.Time     `json:"published_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (EmoticonPack) TableName() string {
	return "emoticon_packs"
}

// ===== Single Emoticon =====

// Emoticon single emoticon model
type Emoticon struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PackID      uint           `gorm:"index;not null" json:"pack_id"`
	EmoticonKey string         `gorm:"type:varchar(64);not null" json:"emoticon_key"`
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url"`
	GifURL      string         `gorm:"type:varchar(500)" json:"gif_url"`
	Width       int            `gorm:"type:int" json:"width"`
	Height      int            `gorm:"type:int" json:"height"`
	SortOrder   int            `gorm:"type:int;default:0" json:"sort_order"`
	Tags        StringArray    `gorm:"type:text" json:"tags"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (Emoticon) TableName() string {
	return "emoticons"
}

// ===== Action Resource =====

// ActionResource action resource model
type ActionResource struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	ActionName    string         `gorm:"type:varchar(255);not null" json:"action_name"`
	Description   string         `gorm:"type:text" json:"description"`
	ActionType    string         `gorm:"type:varchar(20);not null" json:"action_type"` // built-in/official/creator
	Difficulty    string         `gorm:"type:varchar(20)" json:"difficulty"` // easy/medium/hard
	ThumbnailURL  string         `gorm:"type:varchar(500)" json:"thumbnail_url"`
	VideoURL      string         `gorm:"type:varchar(500)" json:"video_url"`
	MotionData    JSON           `gorm:"type:jsonb" json:"motion_data"`
	Price         float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	IsFree        bool           `gorm:"default:true" json:"is_free"`
	Status        string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/pending/approved/rejected/published/removed
	Downloads     int            `gorm:"default:0" json:"downloads"`
	RatingAvg     float64        `gorm:"type:decimal(3,2);default:0" json:"rating_avg"`
	RatingCount   int            `gorm:"default:0" json:"rating_count"`
	Tags          StringArray    `gorm:"type:text" json:"tags"`
	DurationSec   int            `gorm:"type:int;default:0" json:"duration_sec"`
	ReviewedAt    *time.Time     `json:"reviewed_at"`
	ReviewedBy    *uint          `gorm:"index" json:"reviewed_by"`
	PublishedAt   *time.Time     `json:"published_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (ActionResource) TableName() string {
	return "action_resources"
}

// ===== Voice Config =====

// VoiceConfig voice configuration model
type VoiceConfig struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	VoiceName    string         `gorm:"type:varchar(255);not null" json:"voice_name"`
	VoiceType    string         `gorm:"type:varchar(20);not null" json:"voice_type"` // built-in/official/custom
	Description  string         `gorm:"type:text" json:"description"`
	PreviewURL   string         `gorm:"type:varchar(500)" json:"preview_url"`
	AudioSamples JSON           `gorm:"type:jsonb" json:"audio_samples"`
	VoiceParams  JSON           `gorm:"type:jsonb" json:"voice_params"`
	Price        float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	IsFree       bool           `gorm:"default:true" json:"is_free"`
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/pending/approved/rejected/published/removed
	Downloads    int            `gorm:"default:0" json:"downloads"`
	RatingAvg    float64        `gorm:"type:decimal(3,2);default:0" json:"rating_avg"`
	RatingCount  int            `gorm:"default:0" json:"rating_count"`
	Tags         StringArray    `gorm:"type:text" json:"tags"`
	CloneTaskID  string         `gorm:"type:varchar(64)" json:"clone_task_id"`
	CloneStatus  string         `gorm:"type:varchar(20)" json:"clone_status"` // pending/processing/completed/failed
	ReviewedAt   *time.Time     `json:"reviewed_at"`
	ReviewedBy   *uint          `gorm:"index" json:"reviewed_by"`
	PublishedAt  *time.Time     `json:"published_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (VoiceConfig) TableName() string {
	return "voice_configs"
}

// ===== Content Review =====

// ContentReview content review record
type ContentReview struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ItemType   string         `gorm:"type:varchar(20);not null" json:"item_type"` // emoticon/action/voice
	ItemID     uint           `gorm:"index;not null" json:"item_id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ReviewerID *uint          `gorm:"index" json:"reviewer_id"`
	Status     string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/approved/rejected
	Result     string         `gorm:"type:varchar(20)" json:"result"`
	Reason     string         `gorm:"type:text" json:"reason"`
	SubmittedAt time.Time     `json:"submitted_at"`
	ReviewedAt *time.Time     `json:"reviewed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (ContentReview) TableName() string {
	return "content_reviews"
}

// ===== User Purchase =====

// UserPurchase user purchase record
type UserPurchase struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	ItemType      string         `gorm:"type:varchar(20);not null" json:"item_type"` // plugin/emoticon/action/voice
	ItemID        uint           `gorm:"index;not null" json:"item_id"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	PaymentMethod string         `gorm:"type:varchar(20)" json:"payment_method"`
	PaymentID     string         `gorm:"type:varchar(100)" json:"payment_id"`
	Status        string         `gorm:"type:varchar(20);default:'completed'" json:"status"` // completed/pending/refunded
	PurchasedAt   time.Time      `json:"purchased_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (UserPurchase) TableName() string {
	return "user_purchases"
}

// ===== Rating =====

// Rating rating record
type Rating struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	ItemType  string         `gorm:"type:varchar(20);not null" json:"item_type"` // plugin/emoticon/action/voice
	ItemID    uint           `gorm:"index;not null" json:"item_id"`
	Rating    int            `gorm:"type:int;not null" json:"rating"` // 1-5
	Review    string         `gorm:"type:text" json:"review"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (Rating) TableName() string {
	return "ratings"
}
