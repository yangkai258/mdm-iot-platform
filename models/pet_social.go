package models

import (
	"time"
)

// PetSocialPost 宠物社交动态
type PetSocialPost struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PetID       string    `gorm:"type:varchar(64);not null;index" json:"pet_id"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Images      string    `gorm:"type:text" json:"images"` // JSON array of image URLs
	LikeCount   int       `gorm:"default:0" json:"like_count"`
	CommentCount int       `gorm:"default:0" json:"comment_count"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PetSocialPost) TableName() string {
	return "pet_social_posts"
}

// PetSocialComment 动态评论
type PetSocialComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	PetID     string    `gorm:"type:varchar(64);not null;index" json:"pet_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PetSocialComment) TableName() string {
	return "pet_social_comments"
}

// PetSocialFollow 宠物关注关系
type PetSocialFollow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FollowerID  string    `gorm:"type:varchar(64);not null;index" json:"follower_id"`  // 关注者
	FollowingID string    `gorm:"type:varchar(64);not null;index" json:"following_id"` // 被关注者
	CreatedAt   time.Time `json:"created_at"`
}

func (PetSocialFollow) TableName() string {
	return "pet_social_follows"
}

// PetSocialLike 动态点赞
type PetSocialLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	PetID     string    `gorm:"type:varchar(64);not null;index" json:"pet_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (PetSocialLike) TableName() string {
	return "pet_social_likes"
}
