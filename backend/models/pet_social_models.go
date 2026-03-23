package models

import "time"

// Post 宠物动态
type Post struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    PetID        uint      `gorm:"index" json:"pet_id"`
    AuthorID     uint      `gorm:"index" json:"author_id"`
    Content      string    `json:"content"`
    MediaURLs    string    `json:"media_urls"` // JSON array ["url1","url2"]
    PostType     string    `json:"post_type"`  // photo, video, milestone, achievement
    Location     string    `json:"location"`
    LikeCount    int       `gorm:"default:0" json:"like_count"`
    CommentCount int       `gorm:"default:0" json:"comment_count"`
    ShareCount   int       `gorm:"default:0" json:"share_count"`
    IsPublic     bool      `gorm:"default:true" json:"is_public"`
    CreatedAt    time.Time `json:"created_at"`
}

// PostComment 动态评论
type PostComment struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    PostID    uint      `gorm:"index" json:"post_id"`
    UserID    uint      `gorm:"index" json:"user_id"`
    Content   string    `json:"content"`
    LikeCount int       `gorm:"default:0" json:"like_count"`
    CreatedAt time.Time `json:"created_at"`
}

// PostLike 动态点赞
type PostLike struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    PostID    uint      `gorm:"uniqueIndex:idx_post_user" json:"post_id"`
    UserID    uint      `gorm:"uniqueIndex:idx_post_user" json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}

// Follow 关注关系
type Follow struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    FollowerID uint      `gorm:"index" json:"follower_id"`
    FolloweeID uint      `gorm:"index" json:"followee_id"`
    FollowType string    `json:"follow_type"` // pet, user
    CreatedAt  time.Time `json:"created_at"`
}

// PetPlaydate 宠物约会
type PetPlaydate struct {
    ID          uint       `gorm:"primaryKey" json:"id"`
    OrganizerID uint       `gorm:"index" json:"organizer_id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Location    string     `json:"location"`
    PetIDs      string     `json:"pet_ids"` // JSON array of invited pet IDs
    StartTime  time.Time  `json:"start_time"`
    EndTime    time.Time  `json:"end_time"`
    Status      string     `json:"status"` // pending, confirmed, cancelled
    MaxPets     int        `json:"max_pets"`
    CreatedAt   time.Time  `json:"created_at"`
}
