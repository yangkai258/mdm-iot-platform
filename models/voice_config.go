package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VoiceGender 声音性别
const (
	VoiceGenderMale   = "male"
	VoiceGenderFemale = "female"
	VoiceGenderNeutral = "neutral"
)

// VoiceAgeGroup 年龄组
const (
	VoiceAgeChild   = "child"
	VoiceAgeTeen    = "teen"
	VoiceAgeAdult   = "adult"
	VoiceAgeSenior  = "senior"
)

// VoiceConfig 声音配置
type VoiceConfig struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ConfigUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"config_uuid"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	Name          string         `gorm:"type:varchar(64);not null" json:"name"`
	Provider      string         `gorm:"type:varchar(32);not null" json:"provider"` // elevenlabs/azure/gtts/custom
	VoiceID       string         `gorm:"type:varchar(128)" json:"voice_id"`        // 第三方声音ID
	Language      string         `gorm:"type:varchar(16);default:'zh-CN'" json:"language"`
	Gender        string         `gorm:"type:varchar(16)" json:"gender"`          // male/female/neutral
	AgeGroup      string         `gorm:"type:varchar(16)" json:"age_group"`       // child/teen/adult/senior
	Pitch         float64        `gorm:"type:decimal(5,2);default:1.0" json:"pitch"` // 音调 0.5-2.0
	Speed         float64        `gorm:"type:decimal(5,2);default:1.0" json:"speed"` // 语速 0.5-2.0
	Volume        float64        `gorm:"type:decimal(5,2);default:1.0" json:"volume"` // 音量 0-1.0
	Emotion       string         `gorm:"type:varchar(32)" json:"emotion"`         // happy/sad/angry/calm/neutral
	Style         string         `gorm:"type:varchar(64)" json:"style"`           // narration/casual/formal/broadcast
	IsDefault     bool           `gorm:"type:boolean;default:false" json:"is_default"`
	IsPublic      bool           `gorm:"type:boolean;default:false" json:"is_public"` // 是否公开（可市场共享）
	PreviewURL    string         `gorm:"type:varchar(512)" json:"preview_url"`     // 预览音频URL
	Settings      string         `gorm:"type:jsonb;default:'{}'" json:"settings"` // 额外设置
	UseCount      int            `gorm:"type:int;default:0" json:"use_count"`
	Status        string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (VoiceConfig) TableName() string {
	return "voice_configs"
}

// BeforeCreate 创建前自动生成 UUID
func (v *VoiceConfig) BeforeCreate(tx *gorm.DB) error {
	if v.ConfigUUID == "" {
		v.ConfigUUID = uuid.New().String()
	}
	return nil
}

// VoiceConfigResponse 声音配置响应
type VoiceConfigResponse struct {
	ConfigUUID string  `json:"config_uuid"`
	UserID     uint    `json:"user_id"`
	Name       string  `json:"name"`
	Provider   string  `json:"provider"`
	VoiceID    string  `json:"voice_id"`
	Language   string  `json:"language"`
	Gender     string  `json:"gender"`
	AgeGroup   string  `json:"age_group"`
	Pitch      float64 `json:"pitch"`
	Speed      float64 `json:"speed"`
	Volume     float64 `json:"volume"`
	Emotion    string  `json:"emotion"`
	Style      string  `json:"style"`
	IsDefault  bool    `json:"is_default"`
	IsPublic   bool    `json:"is_public"`
	PreviewURL string  `json:"preview_url"`
	Settings   string  `json:"settings"`
	UseCount   int     `json:"use_count"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
}

// ToResponse 转换为响应格式
func (v *VoiceConfig) ToResponse() *VoiceConfigResponse {
	return &VoiceConfigResponse{
		ConfigUUID: v.ConfigUUID,
		UserID:     v.UserID,
		Name:       v.Name,
		Provider:   v.Provider,
		VoiceID:    v.VoiceID,
		Language:   v.Language,
		Gender:     v.Gender,
		AgeGroup:   v.AgeGroup,
		Pitch:      v.Pitch,
		Speed:      v.Speed,
		Volume:     v.Volume,
		Emotion:    v.Emotion,
		Style:      v.Style,
		IsDefault:  v.IsDefault,
		IsPublic:   v.IsPublic,
		PreviewURL: v.PreviewURL,
		Settings:   v.Settings,
		UseCount:   v.UseCount,
		Status:     v.Status,
		CreatedAt:  v.CreatedAt.Format(time.RFC3339),
	}
}

// VoicePreviewRequest 声音预览请求
type VoicePreviewRequest struct {
	Text     string  `json:"text" binding:"required"`
	Provider string  `json:"provider" binding:"required"`
	VoiceID  string  `json:"voice_id"`
	Language string  `json:"language"`
	Pitch    float64 `json:"pitch"`
	Speed    float64 `json:"speed"`
	Volume   float64 `json:"volume"`
	Emotion  string  `json:"emotion"`
	Style    string  `json:"style"`
}

// VoicePreviewResponse 声音预览响应
type VoicePreviewResponse struct {
	PreviewURL string `json:"preview_url"`
	DurationMs int    `json:"duration_ms"`
	Format     string `json:"format"` // mp3/wav/ogg
}

// VoiceMarketItem 声音市场项
type VoiceMarketItem struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ItemUUID   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"item_uuid"`
	ConfigID   uint           `gorm:"index;not null" json:"config_id"`
	CreatorID  uint           `gorm:"index;not null" json:"creator_id"`
	CreatorName string        `gorm:"type:varchar(64)" json:"creator_name"`
	Name       string         `gorm:"type:varchar(64);not null" json:"name"`
	Provider   string         `gorm:"type:varchar(32);not null" json:"provider"`
	VoiceID    string         `gorm:"type:varchar(128)" json:"voice_id"`
	Language   string         `gorm:"type:varchar(16)" json:"language"`
	Gender     string         `gorm:"type:varchar(16)" json:"gender"`
	AgeGroup   string         `gorm:"type:varchar(16)" json:"age_group"`
	Emotion    string         `gorm:"type:varchar(32)" json:"emotion"`
	Style      string         `gorm:"type:varchar(64)" json:"style"`
	PreviewURL string         `gorm:"type:varchar(512)" json:"preview_url"`
	UseCount   int            `gorm:"type:int;default:0" json:"use_count"`
	Rating     float64        `gorm:"type:decimal(3,2);default:0" json:"rating"`
	Price      float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	IsPremium  bool           `gorm:"type:boolean;default:false" json:"is_premium"`
	Status     string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	TenantID   string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (VoiceMarketItem) TableName() string {
	return "voice_market_items"
}

// BeforeCreate 创建前自动生成 UUID
func (v *VoiceMarketItem) BeforeCreate(tx *gorm.DB) error {
	if v.ItemUUID == "" {
		v.ItemUUID = uuid.New().String()
	}
	return nil
}
