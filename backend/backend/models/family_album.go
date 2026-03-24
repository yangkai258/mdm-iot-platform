package models

import (
	"time"

	"gorm.io/gorm"
)

// FamilyAlbumContainer 家庭相册（相册容器）
type FamilyAlbumContainer struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FamilyID    uint           `gorm:"index;not null" json:"family_id"`
	AlbumName   string         `gorm:"type:varchar(128);not null" json:"album_name"`
	Description string         `gorm:"type:text" json:"description"`
	CoverURL    string         `gorm:"type:varchar(500)" json:"cover_url"`
	PhotoCount  int            `gorm:"default:0" json:"photo_count"`
	CreatedBy   uint           `gorm:"index" json:"created_by"`
	Privacy     string         `gorm:"type:varchar(20);default:'family'" json:"privacy"` // family/private/public
	Tags        StringArray    `gorm:"type:text" json:"tags"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (FamilyAlbumContainer) TableName() string {
	return "family_album_containers"
}

// FamilyAlbumItem 家庭相册照片项
type FamilyAlbumItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	AlbumID       uint           `gorm:"index;not null" json:"album_id"`
	UploaderID    uint           `gorm:"index;not null" json:"uploader_id"`
	PhotoURL      string         `gorm:"type:varchar(500);not null" json:"photo_url"`
	ThumbnailURL  string         `gorm:"type:varchar(500)" json:"thumbnail_url"`
	Caption       string         `gorm:"type:varchar(256)" json:"caption"`
	Width         int            `gorm:"type:int" json:"width"`
	Height        int            `gorm:"type:int" json:"height"`
	FileSize      int64          `gorm:"type:bigint" json:"file_size"`
	TakenAt       *time.Time     `json:"taken_at"`
	LocationLat   float64        `gorm:"type:decimal(10,7)" json:"location_lat"`
	LocationLng   float64        `gorm:"type:decimal(10,7)" json:"location_lng"`
	LocationName  string         `gorm:"type:varchar(256)" json:"location_name"`
	Tags          StringArray    `gorm:"type:text" json:"tags"`
	IsFeatured    bool           `gorm:"default:false" json:"is_featured"`
	LikeCount     int            `gorm:"default:0" json:"like_count"`
	CommentCount  int            `gorm:"default:0" json:"comment_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns table name
func (FamilyAlbumItem) TableName() string {
	return "family_album_items"
}
