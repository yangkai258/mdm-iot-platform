package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetStatusV2 宠物状态表 (Sprint 9)
type PetStatusV2 struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	DeviceID          string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"device_id"`
	PetName           string         `gorm:"type:varchar(32);not null;default:'小爪'" json:"pet_name"`
	PetType           string         `gorm:"type:varchar(16);default:'cat'" json:"pet_type"`
	Personality       string         `gorm:"type:jsonb;default:'{}'" json:"personality"`
	Appearance        string         `gorm:"type:jsonb;default:'{}'" json:"appearance"`
	Mood              int            `gorm:"type:int;default:50" json:"mood"`
	Energy            int            `gorm:"type:int;default:100" json:"energy"`
	Hunger            int            `gorm:"type:int;default:0" json:"hunger"`
	PositionX         float64        `gorm:"type:float;default:0" json:"position_x"`
	PositionY         float64        `gorm:"type:float;default:0" json:"position_y"`
	CurrentExpression string         `gorm:"type:varchar(32);default:'happy'" json:"current_expression"`
	CurrentAction     string         `gorm:"type:varchar(32)" json:"current_action"`
	IsOnline          bool           `gorm:"type:boolean;default:false" json:"is_online"`
	LastSeenAt        *time.Time     `json:"last_seen_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetStatusV2) TableName() string {
	return "pet_status"
}

// BeforeCreate 创建前自动生成 UUID
func (p *PetStatusV2) BeforeCreate(tx *gorm.DB) error {
	if p.DeviceID == "" {
		p.DeviceID = uuid.New().String()
	}
	return nil
}

// PetStatusV2Response 宠物状态响应
type PetStatusV2Response struct {
	DeviceID          string  `json:"device_id"`
	PetName           string  `json:"pet_name"`
	PetType           string  `json:"pet_type"`
	Personality       string  `json:"personality"`
	Appearance        string  `json:"appearance"`
	Mood              int     `json:"mood"`
	Energy            int     `json:"energy"`
	Hunger            int     `json:"hunger"`
	PositionX         float64 `json:"position_x"`
	PositionY         float64 `json:"position_y"`
	CurrentExpression string  `json:"current_expression"`
	CurrentAction     string  `json:"current_action"`
	IsOnline          bool    `json:"is_online"`
	LastSeenAt        *string `json:"last_seen_at"`
}

// ToResponse 转换为响应格式
func (p *PetStatusV2) ToResponse() *PetStatusV2Response {
	resp := &PetStatusV2Response{
		DeviceID:          p.DeviceID,
		PetName:           p.PetName,
		PetType:           p.PetType,
		Personality:       p.Personality,
		Appearance:        p.Appearance,
		Mood:              p.Mood,
		Energy:            p.Energy,
		Hunger:            p.Hunger,
		PositionX:         p.PositionX,
		PositionY:         p.PositionY,
		CurrentExpression: p.CurrentExpression,
		CurrentAction:     p.CurrentAction,
		IsOnline:          p.IsOnline,
	}
	if p.LastSeenAt != nil {
		t := p.LastSeenAt.Format(time.RFC3339)
		resp.LastSeenAt = &t
	}
	return resp
}

// PetSettingsUpdate 宠物设置更新请求
type PetSettingsUpdate struct {
	PetName     string `json:"pet_name"`
	PetType     string `json:"pet_type"`
	Personality string `json:"personality"`
	Appearance  string `json:"appearance"`
}

// MoodBoost 心情激励请求
type MoodBoost struct {
	BoostType string `json:"boost_type"` // food, play, praise, music
	Amount    int    `json:"amount"`
}
