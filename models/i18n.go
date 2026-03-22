package models

import (
	"time"
)

// Translation 翻译条目
type Translation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Locale    string    `json:"locale" gorm:"size:16;index:idx_locale_key,unique"` // zh-CN/en-US/ja-JP
	Key       string    `json:"key" gorm:"size:255;index:idx_locale_key,unique"`   // translation key
	Namespace string    `json:"namespace" gorm:"size:64;index"`                    // messages/common/button
	Value     string    `json:"value" gorm:"type:text"`                            // translated value
	Context   string    `json:"context" gorm:"size:512"`                           // optional context/comment
	Tags      string    `json:"tags" gorm:"size:255"`                              // comma-separated tags
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Translation model
func (Translation) TableName() string {
	return "translations"
}

// TranslationRequest 创建/更新翻译请求
type TranslationRequest struct {
	Locale    string `json:"locale" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Namespace string `json:"namespace"`
	Value     string `json:"value" binding:"required"`
	Context   string `json:"context"`
	Tags      string `json:"tags"`
	IsActive  *bool  `json:"is_active"`
}

// TranslationFilter 翻译查询过滤
type TranslationFilter struct {
	Locale    string `form:"locale"`
	Namespace string `form:"namespace"`
	Key       string `form:"key"`
	Tags      string `form:"tags"`
	IsActive  *bool  `form:"is_active"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}
