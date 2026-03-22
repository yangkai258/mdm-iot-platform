package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ActionMarket 动作市场（用户自定义动作发布到市场）
type ActionMarket struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	MarketUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"market_uuid"`
	ActionID      uint           `gorm:"index;not null" json:"action_id"`
	OriginalUUID  string         `gorm:"type:varchar(64)" json:"original_uuid"` // 原始动作UUID
	CreatorID     uint           `gorm:"index;not null" json:"creator_id"`
	CreatorName   string         `gorm:"type:varchar(64)" json:"creator_name"`
	ActionName    string         `gorm:"type:varchar(64);not null" json:"action_name"`
	ActionNameEn  string         `gorm:"type:varchar(64)" json:"action_name_en"`
	Category      string         `gorm:"type:varchar(32);index" json:"category"`
	Description   string         `gorm:"type:text" json:"description"`
	PreviewVideo  string         `gorm:"type:varchar(512)" json:"preview_video"`
	ThumbnailURL  string         `gorm:"type:varchar(512)" json:"thumbnail_url"`
	DurationMs    int            `gorm:"type:int" json:"duration_ms"`
	Parameters    string         `gorm:"type:jsonb;default:'{}'" json:"parameters"`
	AnimationData string         `gorm:"type:jsonb;default:'{}'" json:"animation_data"`
	MotorCommands string         `gorm:"type:jsonb;default:'{}'" json:"motor_commands"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	Price         float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	Currency      string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	DownloadCount int            `gorm:"type:int;default:0" json:"download_count"`
	UseCount      int            `gorm:"type:int;default:0" json:"use_count"`
	Rating        float64        `gorm:"type:decimal(3,2);default:0" json:"rating"`
	RatingCount   int            `gorm:"type:int;default:0" json:"rating_count"`
	IsPremium     bool           `gorm:"type:boolean;default:false" json:"is_premium"`
	IsFeatured    bool           `gorm:"type:boolean;default:false" json:"is_featured"`
	Status        string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/approved/rejected/removed
	ReviewNote    string         `gorm:"type:text" json:"review_note"`
	PublishedAt  *time.Time     `json:"published_at"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ActionMarket) TableName() string {
	return "action_market"
}

// BeforeCreate 创建前自动生成 UUID
func (a *ActionMarket) BeforeCreate(tx *gorm.DB) error {
	if a.MarketUUID == "" {
		a.MarketUUID = uuid.New().String()
	}
	return nil
}

// ActionMarketResponse 动作市场响应
type ActionMarketResponse struct {
	MarketUUID    string   `json:"market_uuid"`
	ActionID      uint     `json:"action_id"`
	OriginalUUID  string   `json:"original_uuid"`
	CreatorID     uint     `json:"creator_id"`
	CreatorName   string   `json:"creator_name"`
	ActionName    string   `json:"action_name"`
	ActionNameEn  string   `json:"action_name_en"`
	Category      string   `json:"category"`
	Description   string   `json:"description"`
	PreviewVideo  string   `json:"preview_video"`
	ThumbnailURL  string   `json:"thumbnail_url"`
	DurationMs    int      `json:"duration_ms"`
	Parameters    string   `json:"parameters"`
	AnimationData string   `json:"animation_data"`
	MotorCommands string   `json:"motor_commands"`
	Tags          []string `json:"tags"`
	Price         float64  `json:"price"`
	Currency      string   `json:"currency"`
	DownloadCount int      `json:"download_count"`
	UseCount      int      `json:"use_count"`
	Rating        float64  `json:"rating"`
	RatingCount   int      `json:"rating_count"`
	IsPremium     bool     `json:"is_premium"`
	IsFeatured    bool     `json:"is_featured"`
	Status        string   `json:"status"`
	PublishedAt   string   `json:"published_at,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

// ToResponse 转换为响应格式
func (a *ActionMarket) ToResponse() *ActionMarketResponse {
	resp := &ActionMarketResponse{
		MarketUUID:    a.MarketUUID,
		ActionID:      a.ActionID,
		OriginalUUID:  a.OriginalUUID,
		CreatorID:     a.CreatorID,
		CreatorName:   a.CreatorName,
		ActionName:    a.ActionName,
		ActionNameEn:  a.ActionNameEn,
		Category:      a.Category,
		Description:   a.Description,
		PreviewVideo:  a.PreviewVideo,
		ThumbnailURL:  a.ThumbnailURL,
		DurationMs:    a.DurationMs,
		Parameters:    a.Parameters,
		AnimationData: a.AnimationData,
		MotorCommands: a.MotorCommands,
		Tags:          a.Tags,
		Price:         a.Price,
		Currency:      a.Currency,
		DownloadCount: a.DownloadCount,
		UseCount:      a.UseCount,
		Rating:        a.Rating,
		RatingCount:   a.RatingCount,
		IsPremium:     a.IsPremium,
		IsFeatured:    a.IsFeatured,
		Status:        a.Status,
		CreatedAt:     a.CreatedAt.Format(time.RFC3339),
	}
	if a.PublishedAt != nil {
		resp.PublishedAt = a.PublishedAt.Format(time.RFC3339)
	}
	return resp
}

// CustomAction 用户自定义动作（与动作库区分，用户创建的私人动作）
type CustomAction struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ActionUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"action_uuid"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	ActionName    string         `gorm:"type:varchar(64);not null" json:"action_name"`
	ActionNameEn  string         `gorm:"type:varchar(64)" json:"action_name_en"`
	Category      string         `gorm:"type:varchar(32);index" json:"category"`
	Description   string         `gorm:"type:text" json:"description"`
	DurationMs    int            `gorm:"type:int" json:"duration_ms"`
	Priority      int            `gorm:"type:int;default:5" json:"priority"`
	IsEmergency   bool           `gorm:"type:boolean;default:false" json:"is_emergency"`
	Parameters    string         `gorm:"type:jsonb;default:'{}'" json:"parameters"`
	AnimationData string         `gorm:"type:jsonb;default:'{}'" json:"animation_data"`
	MotorCommands string         `gorm:"type:jsonb;default:'{}'" json:"motor_commands"`
	AudioFile     string         `gorm:"type:varchar(256)" json:"audio_file"`
	IsPublished   bool           `gorm:"type:boolean;default:false" json:"is_published"`
	PublishedAt   *time.Time     `json:"published_at"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CustomAction) TableName() string {
	return "custom_actions"
}

// BeforeCreate 创建前自动生成 UUID
func (c *CustomAction) BeforeCreate(tx *gorm.DB) error {
	if c.ActionUUID == "" {
		c.ActionUUID = uuid.New().String()
	}
	return nil
}

// CustomActionResponse 自定义动作响应
type CustomActionResponse struct {
	ActionUUID    string   `json:"action_uuid"`
	UserID        uint     `json:"user_id"`
	ActionName    string   `json:"action_name"`
	ActionNameEn  string   `json:"action_name_en"`
	Category      string   `json:"category"`
	Description   string   `json:"description"`
	DurationMs    int      `json:"duration_ms"`
	Priority      int      `json:"priority"`
	IsEmergency   bool     `json:"is_emergency"`
	Parameters    string   `json:"parameters"`
	AnimationData string   `json:"animation_data"`
	MotorCommands string   `json:"motor_commands"`
	AudioFile     string   `json:"audio_file"`
	IsPublished   bool     `json:"is_published"`
	PublishedAt   string   `json:"published_at,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

// ToResponse 转换为响应格式
func (c *CustomAction) ToResponse() *CustomActionResponse {
	resp := &CustomActionResponse{
		ActionUUID:    c.ActionUUID,
		UserID:        c.UserID,
		ActionName:    c.ActionName,
		ActionNameEn:  c.ActionNameEn,
		Category:      c.Category,
		Description:   c.Description,
		DurationMs:    c.DurationMs,
		Priority:      c.Priority,
		IsEmergency:   c.IsEmergency,
		Parameters:    c.Parameters,
		AnimationData: c.AnimationData,
		MotorCommands: c.MotorCommands,
		AudioFile:     c.AudioFile,
		IsPublished:   c.IsPublished,
		CreatedAt:     c.CreatedAt.Format(time.RFC3339),
	}
	if c.PublishedAt != nil {
		resp.PublishedAt = c.PublishedAt.Format(time.RFC3339)
	}
	return resp
}
