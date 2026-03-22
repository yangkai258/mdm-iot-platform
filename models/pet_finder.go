package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// PetLostReport 走失报告
type PetLostReport struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ReportUUID    string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"report_uuid"`
	PetUUID       string         `gorm:"type:varchar(64);index" json:"pet_uuid"`
	PetName       string         `gorm:"type:varchar(32);not null" json:"pet_name"`
	Species       string         `gorm:"type:varchar(32)" json:"species"` // dog/cat/bird/rabbit/other
	Breed         string         `gorm:"type:varchar(64)" json:"breed"`
	Color         string         `gorm:"type:varchar(32)" json:"color"`
	Gender        string         `gorm:"type:varchar(16)" json:"gender"` // male/female/unknown
	Age           string         `gorm:"type:varchar(32)" json:"age"`   // 如 "2岁"
	LastLocation  JSON           `gorm:"type:jsonb" json:"last_location"` // {"lat": xx, "lng": xx, "address": "..."}
	LostTime      time.Time      `json:"lost_time"`
	Status        string         `gorm:"type:varchar(20);default:'searching'" json:"status"` // searching/found/closed/abandoned
	Reward        string         `gorm:"type:varchar(256)" json:"reward"`
	ContactName   string         `gorm:"type:varchar(64)" json:"contact_name"`
	ContactPhone  string         `gorm:"type:varchar(32)" json:"contact_phone"`
	ContactWechat string         `gorm:"type:varchar(64)" json:"contact_wechat"`
	Description   string         `gorm:"type:text" json:"description"`
	PhotoURLs     pq.StringArray `gorm:"type:text[]" json:"photo_urls"`
	ReporterID    uint           `gorm:"index;not null" json:"reporter_id"`
	SpreadRadius  float64        `gorm:"type:decimal(5,2);default:10" json:"spread_radius_km"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (PetLostReport) TableName() string {
	return "pet_lost_reports"
}

// BeforeCreate 创建前自动生成 UUID
func (p *PetLostReport) BeforeCreate(tx *gorm.DB) error {
	if p.ReportUUID == "" {
		p.ReportUUID = uuid.New().String()
	}
	return nil
}

// SightingReport 目击报告
type SightingReport struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SightingUUID string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"sighting_uuid"`
	ReportUUID   string         `gorm:"type:varchar(64);index;not null" json:"report_uuid"`
	Location     JSON           `gorm:"type:jsonb" json:"location"` // {"lat": xx, "lng": xx, "address": "..."}
	SightingTime time.Time      `json:"sighting_time"`
	Description  string         `gorm:"type:text" json:"description"`
	PhotoURL     string         `gorm:"type:varchar(512)" json:"photo_url"`
	ReporterName string         `gorm:"type:varchar(64)" json:"reporter_name"`
	ContactPhone string         `gorm:"type:varchar(32)" json:"contact_phone"`
	IsCredible   bool           `gorm:"type:boolean;default:true" json:"is_credible"`
	ReporterID   uint           `gorm:"index" json:"reporter_id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt    time.Time      `json:"created_at"`
}

// TableName 指定表名
func (SightingReport) TableName() string {
	return "sighting_reports"
}

// BeforeCreate 创建前自动生成 UUID
func (s *SightingReport) BeforeCreate(tx *gorm.DB) error {
	if s.SightingUUID == "" {
		s.SightingUUID = uuid.New().String()
	}
	return nil
}

// FinderAlert 寻宠警报订阅
type FinderAlert struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	AlertUUID   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"alert_uuid"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	// 订阅条件
	Species     string         `gorm:"type:varchar(32)" json:"species"`   // 订阅的物种（空=全部）
	Latitude    float64        `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude   float64        `gorm:"type:decimal(10,7)" json:"longitude"`
	RadiusKM    float64        `gorm:"type:decimal(5,2);default:10" json:"radius_km"` // 半径km
	// 通知设置
	NotifyEmail  bool          `gorm:"type:boolean;default:true" json:"notify_email"`
	NotifySMS    bool          `gorm:"type:boolean;default:false" json:"notify_sms"`
	NotifyApp    bool          `gorm:"type:boolean;default:true" json:"notify_app"`
	// 状态
	IsActive    bool           `gorm:"type:boolean;default:true" json:"is_active"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (FinderAlert) TableName() string {
	return "finder_alerts"
}

// BeforeCreate 创建前自动生成 UUID
func (f *FinderAlert) BeforeCreate(tx *gorm.DB) error {
	if f.AlertUUID == "" {
		f.AlertUUID = uuid.New().String()
	}
	return nil
}
