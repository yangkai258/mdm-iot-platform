package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// EmoticonCategory 表情包分类
type EmoticonCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CategoryID  string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"category_id"`
	Name        string         `gorm:"type:varchar(64);not null" json:"name"`
	NameEn      string         `gorm:"type:varchar(64)" json:"name_en"`
	IconURL     string         `gorm:"type:varchar(512)" json:"icon_url"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"type:int;default:0" json:"sort_order"`
	IsActive    bool           `gorm:"type:boolean;default:true" json:"is_active"`
	EmoticonCount int          `gorm:"-" json:"emoticon_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (EmoticonCategory) TableName() string {
	return "emoticon_categories"
}

// BeforeCreate 创建前自动生成 UUID
func (e *EmoticonCategory) BeforeCreate(tx *gorm.DB) error {
	if e.CategoryID == "" {
		e.CategoryID = uuid.New().String()
	}
	return nil
}

// Emoticon 表情包
type Emoticon struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	EmoticonUUID  string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"emoticon_uuid"`
	Name          string         `gorm:"type:varchar(128);not null" json:"name"`
	CategoryID    string         `gorm:"type:varchar(64);index;not null" json:"category_id"`
	CreatorID     uint           `gorm:"index;not null" json:"creator_id"`
	CreatorType   string         `gorm:"type:varchar(20);default:'user'" json:"creator_type"` // user/system/market
	ThumbnailURL  string         `gorm:"type:varchar(512)" json:"thumbnail_url"`
	ImageURL      string         `gorm:"type:varchar(512)" json:"image_url"`
	GifURL        string         `gorm:"type:varchar(512)" json:"gif_url"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2);default:0" json:"price"` // 0 = 免费
	Currency      string         `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	DownloadCount int            `gorm:"type:int;default:0" json:"download_count"`
	UseCount      int            `gorm:"type:int;default:0" json:"use_count"`
	IsPremium     bool           `gorm:"type:boolean;default:false" json:"is_premium"`
	IsFeatured    bool           `gorm:"type:boolean;default:false" json:"is_featured"`
	Status        string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive/pending
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Emoticon) TableName() string {
	return "emoticons"
}

// BeforeCreate 创建前自动生成 UUID
func (e *Emoticon) BeforeCreate(tx *gorm.DB) error {
	if e.EmoticonUUID == "" {
		e.EmoticonUUID = uuid.New().String()
	}
	return nil
}

// EmoticonPurchase 表情包购买记录
type EmoticonPurchase struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PurchaseUUID  string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"purchase_uuid"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	EmoticonID    uint      `gorm:"index;not null" json:"emoticon_id"`
	EmoticonUUID  string    `gorm:"type:varchar(64)" json:"emoticon_uuid"`
	Price         float64   `gorm:"type:decimal(10,2)" json:"price"`
	Currency      string    `gorm:"type:varchar(8);default:'CNY'" json:"currency"`
	PaymentMethod string    `gorm:"type:varchar(32)" json:"payment_method"` // coin/cash/free
	TransactionID string    `gorm:"type:varchar(128)" json:"transaction_id"`
	TenantID      string    `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 指定表名
func (EmoticonPurchase) TableName() string {
	return "emoticon_purchases"
}

// BeforeCreate 创建前自动生成 UUID
func (e *EmoticonPurchase) BeforeCreate(tx *gorm.DB) error {
	if e.PurchaseUUID == "" {
		e.PurchaseUUID = uuid.New().String()
	}
	return nil
}

// EmoticonResponse 表情包响应
type EmoticonResponse struct {
	EmoticonUUID  string   `json:"emoticon_uuid"`
	Name         string   `json:"name"`
	CategoryID   string   `json:"category_id"`
	CategoryName string   `json:"category_name,omitempty"`
	ThumbnailURL string   `json:"thumbnail_url"`
	ImageURL     string   `json:"image_url"`
	GifURL       string   `json:"gif_url"`
	Tags         []string `json:"tags"`
	Description  string   `json:"description"`
	Price        float64  `json:"price"`
	Currency     string   `json:"currency"`
	IsPremium    bool     `json:"is_premium"`
	IsFeatured   bool     `json:"is_featured"`
	DownloadCount int     `json:"download_count"`
	UseCount     int      `json:"use_count"`
	Status       string   `json:"status"`
	CreatedAt    string   `json:"created_at"`
}

// ToResponse 转换为响应格式
func (e *Emoticon) ToResponse() *EmoticonResponse {
	return &EmoticonResponse{
		EmoticonUUID:  e.EmoticonUUID,
		Name:          e.Name,
		CategoryID:    e.CategoryID,
		ThumbnailURL:  e.ThumbnailURL,
		ImageURL:     e.ImageURL,
		GifURL:       e.GifURL,
		Tags:         e.Tags,
		Description:  e.Description,
		Price:        e.Price,
		Currency:     e.Currency,
		IsPremium:    e.IsPremium,
		IsFeatured:   e.IsFeatured,
		DownloadCount: e.DownloadCount,
		UseCount:     e.UseCount,
		Status:       e.Status,
		CreatedAt:    e.CreatedAt.Format(time.RFC3339),
	}
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
	IconURL     string `json:"icon_url"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	EmoticonCount int  `json:"emoticon_count"`
	IsActive    bool   `json:"is_active"`
}

// ToCategoryResponse 转换为响应格式
func (e *EmoticonCategory) ToCategoryResponse() *CategoryResponse {
	return &CategoryResponse{
		CategoryID:    e.CategoryID,
		Name:          e.Name,
		NameEn:        e.NameEn,
		IconURL:       e.IconURL,
		Description:   e.Description,
		SortOrder:     e.SortOrder,
		EmoticonCount: e.EmoticonCount,
		IsActive:      e.IsActive,
	}
}
