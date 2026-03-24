package models

import (
	"time"
)

// Translation 国际化翻译
type Translation struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Locale     string         `gorm:"size:16;index" json:"locale"` // en, zh-CN, zh-TW
	Key        string         `gorm:"size:255;index" json:"key"`
	Namespace  string         `gorm:"size:64;index" json:"namespace"` // common, button, message
	Value      string         `gorm:"type:text" json:"value"`
	Context    string         `gorm:"size:512" json:"context"`
	Tags       string         `gorm:"size:255" json:"tags"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
