package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification 通知记录
type Notification struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeviceID    string         `gorm:"type:varchar(36);index;not null" json:"device_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Priority    int            `gorm:"type:smallint;default:0" json:"priority"` // 0:普通 1:重要 2:紧急
	Channel     string         `gorm:"type:varchar(20);default:'push'" json:"channel"` // push, sms, email
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, sent, delivered, failed
	SentAt      *time.Time     `json:"sent_at"`
	DeliveredAt *time.Time     `json:"delivered_at"`
	CreatedBy   string         `gorm:"type:varchar(36)" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.DeviceID == "" {
		n.DeviceID = uuid.New().String()
	}
	return nil
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // 模板编码，如 "low_battery_alert"
	TitleTpl  string         `gorm:"type:varchar(255);not null" json:"title_tpl"`       // 标题模板，支持 {{variable}} 替换
	ContentTpl string         `gorm:"type:text;not null" json:"content_tpl"`            // 内容模板，支持 {{variable}} 替换
	Channel   string         `gorm:"type:varchar(20);default:'push'" json:"channel"`   // push, sms, email
	Priority  int            `gorm:"type:smallint;default:0" json:"priority"`           // 默认优先级
	Variables string         `gorm:"type:text" json:"variables"`                       // 变量列表，JSON 格式，如 ["device_name", "battery_level"]
	Enabled   bool           `gorm:"type:boolean;default:true" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Announcement 企业公告
type Announcement struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Type       string         `gorm:"type:varchar(20);default:'info'" json:"type"` // info, warning, critical
	Priority   int            `gorm:"type:smallint;default:0" json:"priority"`
	TargetType string         `gorm:"type:varchar(20);default:'all'" json:"target_type"` // all, company, device_group
	TargetID   string         `gorm:"type:varchar(36)" json:"target_id"`                // 公司ID或设备组ID
	Status     string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, published, archived
	StartTime  *time.Time     `json:"start_time"`
	EndTime    *time.Time     `json:"end_time"`
	CreatedBy  string         `gorm:"type:varchar(36)" json:"created_by"`
	PublishedBy string        `gorm:"type:varchar(36)" json:"published_by"`
	PublishedAt *time.Time    `json:"published_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
