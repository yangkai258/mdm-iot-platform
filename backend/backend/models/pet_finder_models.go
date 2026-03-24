package models

import (
	"time"
)

// PetFinderAlert 寻宠警报
type PetFinderAlert struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ReportID     int64          `gorm:"not null;index" json:"report_id"`
	AlertType    string         `gorm:"size:20" json:"alert_type"` // sms, push, email
	RecipientID  int64          `gorm:"index" json:"recipient_id"`
	SentAt       *time.Time     `json:"sent_at"`
	Status       string         `gorm:"size:20" json:"status"` // sent, failed, pending
	ErrorMsg     string         `gorm:"type:text" json:"error_msg"`
	CreatedAt    time.Time      `json:"created_at"`
}

// PetFinderReport 寻宠报告
type PetFinderReport struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PetID          int64          `gorm:"not null;index" json:"pet_id"`
	UserID         int64          `gorm:"not null;index" json:"user_id"`
	ReportType     string         `gorm:"size:20;not null" json:"report_type"` // lost, found
	Title          string         `gorm:"size:255;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	LastSeenAt     *time.Time     `json:"last_seen_at"`
	LastSeenLat    float64        `gorm:"type:numeric(10,6)" json:"last_seen_lat"`
	LastSeenLng    float64        `gorm:"type:numeric(10,6)" json:"last_seen_lng"`
	LastSeenAddr   string         `gorm:"size:500" json:"last_seen_addr"`
	ContactName    string         `gorm:"size:100" json:"contact_name"`
	ContactPhone   string         `gorm:"size:50" json:"contact_phone"`
	Reward         float64        `gorm:"type:numeric(10,2)" json:"reward"`
	RewardMemo     string         `gorm:"size:255" json:"reward_memo"`
	Photos         []string       `gorm:"type:character varying(500)[]" json:"photos"`
	Status         string         `gorm:"size:20;default:'active'" json:"status"`
	ResolvedAt     *time.Time     `json:"resolved_at"`
	ViewCount      int64          `gorm:"default:0" json:"view_count"`
	AlertRadius    float64        `gorm:"type:numeric(10,2);default:5" json:"alert_radius"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Sightings      []PetFinderSighting `gorm:"foreignKey:ReportID" json:"sightings,omitempty"`
	Alerts         []PetFinderAlert    `gorm:"foreignKey:ReportID" json:"alerts,omitempty"`
}

// PetFinderSighting 寻宠目击记录
type PetFinderSighting struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ReportID       int64          `gorm:"not null;index" json:"report_id"`
	ReporterID     int64          `gorm:"not null;index" json:"reporter_id"`
	SightedAt      *time.Time     `json:"sighted_at"`
	SightedLat     float64        `gorm:"type:numeric(10,6)" json:"sighted_lat"`
	SightedLng     float64        `gorm:"type:numeric(10,6)" json:"sighted_lng"`
	SightedAddr    string         `gorm:"size:500" json:"sighted_addr"`
	Description    string         `gorm:"type:text" json:"description"`
	PetStatus      string         `gorm:"size:50" json:"pet_status"` // alone, with_owner, rescued
	PhotoURL       string         `gorm:"size:500" json:"photo_url"`
	ContactName    string         `gorm:"size:100" json:"contact_name"`
	ContactPhone   string         `gorm:"size:50" json:"contact_phone"`
	IsVerified     bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt      time.Time      `json:"created_at"`
}

// SightingReport 目击报告（另一套）
type SightingReport struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SightingUUID   string         `gorm:"size:64;uniqueIndex;not null" json:"sighting_uuid"`
	ReportUUID     string         `gorm:"size:64;index" json:"report_uuid"`
	Location       string         `gorm:"type:jsonb" json:"location"` // {"lat": 0, "lng": 0}
	SightingTime   *time.Time     `json:"sighting_time"`
	Description    string         `gorm:"type:text" json:"description"`
	PhotoURL       string         `gorm:"size:512" json:"photo_url"`
	ReporterName   string         `gorm:"size:64" json:"reporter_name"`
	ContactPhone   string         `gorm:"size:32" json:"contact_phone"`
	IsCredible     bool           `gorm:"default:true" json:"is_credible"`
	ReporterID     int64          `gorm:"index" json:"reporter_id"`
	TenantID       string         `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
}
