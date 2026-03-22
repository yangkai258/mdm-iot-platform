package models

import (
	"time"
)

// PushStatus 推送状态
const (
	PushStatusPending   = "pending"   // 待发送
	PushStatusSent      = "sent"      // 已发送
	PushStatusDelivered = "delivered" // 已送达
	PushStatusClicked   = "clicked"   // 已点击
	PushStatusFailed    = "failed"    // 发送失败
)

// PushType 推送类型
const (
	PushTypeAlert     = "alert"     // 告警通知
	PushTypeDevice    = "device"    // 设备通知
	PushTypeMember    = "member"    // 会员通知
	PushTypeSystem    = "system"    // 系统通知
	PushTypeOTA       = "ota"       // OTA升级通知
	PushTypeMarketing = "marketing" // 营销通知
)

// AppPush App推送记录表
type AppPush struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	PushID      string     `gorm:"type:varchar(36);uniqueIndex;not null" json:"push_id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	DeviceID    string     `gorm:"type:varchar(36);index" json:"device_id"` // 关联设备ID（可选）
	Platform    string     `gorm:"type:varchar(20);index" json:"platform"`  // ios, android, miniapp
	ClientID    string     `gorm:"type:varchar(128)" json:"client_id"`      // 客户端标识
	Title       string     `gorm:"type:varchar(128)" json:"title"`
	Body        string     `gorm:"type:varchar(512)" json:"body"`
	PushType    string     `gorm:"type:varchar(32);index" json:"push_type"`                // alert, device, member, system, ota, marketing
	Data        string     `gorm:"type:text" json:"data"`                                  // 透传数据 JSON
	Badge       int        `gorm:"default:1" json:"badge"`                                 // 角标数字
	Sound       string     `gorm:"type:varchar(64)" json:"sound"`                          // 提示音
	ChannelID   string     `gorm:"type:varchar(64)" json:"channel_id"`                     // 推送渠道ID
	Tag         string     `gorm:"type:varchar(64)" json:"tag"`                            // 标签（用于分组/撤回）
	Status      string     `gorm:"type:varchar(20);default:'pending';index" json:"status"` // pending, sent, delivered, clicked, failed
	ErrorMsg    string     `gorm:"type:text" json:"error_msg"`                             // 错误信息
	SentAt      *time.Time `json:"sent_at"`                                                // 实际发送时间
	DeliveredAt *time.Time `json:"delivered_at"`                                           // 送达时间
	ClickedAt   *time.Time `json:"clicked_at"`                                             // 点击时间
	ExpiresAt   *time.Time `json:"expires_at"`                                             // 过期时间（用于撤回）
	RetryCount  int        `gorm:"default:0" json:"retry_count"`                           // 重试次数
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (AppPush) TableName() string {
	return "app_pushes"
}

// MiniAppDevice 小程序设备绑定表
type MiniAppDevice struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index" json:"user_id"`
	OpenID       string     `gorm:"type:varchar(64);index" json:"open_id"`   // 微信小程序OpenID
	DeviceID     string     `gorm:"type:varchar(36);index" json:"device_id"` // 绑定的设备ID
	Nickname     string     `gorm:"type:varchar(64)" json:"nickname"`        // 设备昵称（如"客厅的Momo"）
	BindTime     time.Time  `json:"bind_time"`                               // 绑定时间
	UnbindTime   *time.Time `json:"unbind_time"`                             // 解绑时间（nil表示有效绑定）
	IsActive     bool       `gorm:"default:true;index" json:"is_active"`     // 是否有效
	BindToken    string     `gorm:"type:varchar(128)" json:"bind_token"`     // 绑定Token（用于二维码扫码绑定）
	BindTokenExp *time.Time `json:"bind_token_exp"`                          // 绑定Token过期时间
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (MiniAppDevice) TableName() string {
	return "miniapp_devices"
}

// MiniAppQRCodeBind 二维码扫码绑定记录
type MiniAppQRCodeBind struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	QRCodeID      string     `gorm:"type:varchar(36);uniqueIndex;not null" json:"qr_code_id"`
	DeviceID      string     `gorm:"type:varchar(36);index" json:"device_id"`
	Scene         string     `gorm:"type:varchar(32)" json:"scene"`         // 扫码场景
	OpenID        string     `gorm:"type:varchar(64);index" json:"open_id"` // 扫码用户OpenID
	ExpireMinutes int        `gorm:"default:30" json:"expire_minutes"`      // 过期分钟数
	ExpiresAt     time.Time  `json:"expires_at"`                            // 过期时间
	BoundAt       *time.Time `json:"bound_at"`                              // 绑定成功时间
	IsBound       bool       `gorm:"default:false" json:"is_bound"`         // 是否已绑定
	BindUserID    uint       `gorm:"index" json:"bind_user_id"`             // 绑定用户ID
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (MiniAppQRCodeBind) TableName() string {
	return "miniapp_qrcode_binds"
}
