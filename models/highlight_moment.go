package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HighlightType 高光时刻类型
const (
	HighlightTypeCute         = "cute"
	HighlightTypePlay         = "play"
	HighlightTypeMilestone    = "milestone"
	HighlightTypeAchieve      = "achieve"
	HighlightTypeHeartwarming = "heartwarming"
	HighlightTypeFunny        = "funny"
	HighlightTypeRare         = "rare"
	HighlightTypeSocial       = "social"
)

// HighlightMoment 高光时刻 (Sprint 18)
type HighlightMoment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	MomentUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"moment_uuid"`
	PetUUID       string         `gorm:"type:varchar(64);index;not null" json:"pet_uuid"`
	DeviceID      string         `gorm:"type:varchar(64);index" json:"device_id"`
	HighlightType string         `gorm:"type:varchar(32);not null;index" json:"highlight_type"`
	Title         string         `gorm:"type:varchar(128);not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	OccurredAt    time.Time      `gorm:"type:timestamp;not null;index" json:"occurred_at"`
	MediaURLs     StringArray    `gorm:"type:text[]" json:"media_urls"`
	ThumbnailURL  string         `gorm:"type:varchar(512)" json:"thumbnail_url"`
	VideoURL      string         `gorm:"type:varchar(512)" json:"video_url"`
	Duration      int            `gorm:"type:int" json:"duration"`
	CapturedBy    string         `gorm:"type:varchar(32);default:'auto'" json:"captured_by"`
	AIConfidence  float64        `gorm:"type:decimal(5,4)" json:"ai_confidence"`
	EmotionScore  float64        `gorm:"type:decimal(5,2)" json:"emotion_score"`
	Shareable     bool           `gorm:"type:boolean;default:true" json:"shareable"`
	IsPublic      bool           `gorm:"type:boolean;default:false" json:"is_public"`
	LikeCount     int            `gorm:"type:int;default:0" json:"like_count"`
	ViewCount     int            `gorm:"type:int;default:0" json:"view_count"`
	Tags          StringArray    `gorm:"type:text[]" json:"tags"`
	Metadata      JSON           `gorm:"type:jsonb" json:"metadata"`
	RelatedBehavior string        `gorm:"type:varchar(32)" json:"related_behavior"`
	RelatedVitals JSON           `gorm:"type:jsonb" json:"related_vitals"`
	OwnerID       uint           `gorm:"index" json:"owner_id"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (HighlightMoment) TableName() string {
	return "highlight_moments"
}

func (h *HighlightMoment) BeforeCreate(tx *gorm.DB) error {
	if h.MomentUUID == "" {
		h.MomentUUID = uuid.New().String()
	}
	return nil
}

// TimelineEvent 时间轴事件（统一格式）
type TimelineEvent struct {
	ID           uint      `json:"id"`
	EventUUID    string    `json:"event_uuid"`
	EventType    string    `json:"event_type"`
	EventSubType string    `json:"event_sub_type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	OccurredAt   time.Time `json:"occurred_at"`
	Duration     int       `json:"duration"`
	Intensity    float64   `json:"intensity"`
	IsAnomaly    bool      `json:"is_anomaly"`
	AlertLevel   string    `json:"alert_level,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	VideoURL     string    `json:"video_url,omitempty"`
	MediaURLs    []string  `json:"media_urls,omitempty"`
	RelatedData  JSON      `json:"related_data,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
}

// ReqTimelineQuery 时间轴查询请求
type ReqTimelineQuery struct {
	EventTypes    string `form:"event_types" json:"event_types"`
	StartTime     string `form:"start_time" json:"start_time"`
	EndTime       string `form:"end_time" json:"end_time"`
	IsAnomaly     *bool  `form:"is_anomaly" json:"is_anomaly"`
	HighlightOnly bool   `form:"highlight_only" json:"highlight_only"`
	Page          int    `form:"page" json:"page"`
	PageSize      int    `form:"page_size" json:"page_size"`
}

// RespTimeline 时间轴响应
type RespTimeline struct {
	PetUUID    string           `json:"pet_uuid"`
	Events     []TimelineEvent  `json:"events"`
	Total      int64            `json:"total"`
	Highlights []HighlightMoment `json:"highlights"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

// ReqHighlightQuery 高光时刻查询请求
type ReqHighlightQuery struct {
	HighlightType string `form:"highlight_type" json:"highlight_type"`
	StartTime     string `form:"start_time" json:"start_time"`
	EndTime       string `form:"end_time" json:"end_time"`
	IsPublic      *bool  `form:"is_public" json:"is_public"`
	Page          int    `form:"page" json:"page"`
	PageSize      int    `form:"page_size" json:"page_size"`
}

// ReqCreateHighlight 创建高光时刻请求
type ReqCreateHighlight struct {
	HighlightType string   `json:"highlight_type" binding:"required"`
	Title         string   `json:"title" binding:"required"`
	Description   string   `json:"description"`
	OccurredAt    string   `json:"occurred_at"`
	MediaURLs     []string `json:"media_urls"`
	Shareable     bool     `json:"shareable"`
	IsPublic      bool     `json:"is_public"`
	Tags          []string `json:"tags"`
}

// EventDetail 事件详情
type EventDetail struct {
	EventType     string           `json:"event_type"`
	EventUUID     string           `json:"event_uuid"`
	Data          interface{}      `json:"data"`
	RelatedEvents []TimelineEvent  `json:"related_events,omitempty"`
	Analysis      JSON             `json:"analysis,omitempty"`
}
