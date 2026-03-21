package models

import (
	"encoding/json"
	"time"
)

// ActivityLog 活动日志/审计日志
type ActivityLog struct {
	ID           uint                 `gorm:"primaryKey" json:"id"`
	UserID       uint                 `gorm:"index" json:"user_id"`
	Username     string               `gorm:"type:varchar(64);index" json:"username"`
	Action       string               `gorm:"type:varchar(64);index" json:"action"`               // create/update/delete/login/logout
	ResourceType string               `gorm:"type:varchar(64);index" json:"resource_type"`        // device/member/role/config
	ResourceID   uint                 `gorm:"index" json:"resource_id"`                            // 资源ID
	ResourceName string               `gorm:"type:varchar(255)" json:"resource_name"`            // 资源名称（冗余，便于展示）
	Details      map[string]interface{} `gorm:"type:jsonb" json:"details"`                      // 详细信息（JSON）
	IP           string               `gorm:"type:varchar(32)" json:"ip"`
	UserAgent    string               `gorm:"type:varchar(255)" json:"user_agent"`
	TenantID     string               `gorm:"index" json:"tenant_id"`
	CreatedAt    time.Time            `gorm:"index" json:"created_at"`
}

func (ActivityLog) TableName() string { return "activity_logs" }

// LoginLog 登录日志（独立表，与 SysLoginLog 分开）
type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"type:varchar(64);index" json:"username"`
	IP        string    `gorm:"type:varchar(32)" json:"ip"`
	Location  string    `gorm:"type:varchar(255)" json:"location"`
	Browser   string    `gorm:"type:varchar(50)" json:"browser"`
	OS        string    `gorm:"type:varchar(50)" json:"os"`
	Status    int       `gorm:"default:1" json:"status"` // 1:成功 0:失败
	Msg       string    `gorm:"type:varchar(255)" json:"msg"`
	TenantID  string    `gorm:"index" json:"tenant_id"`
	LoginTime time.Time `gorm:"index" json:"login_time"`
	CreatedAt time.Time `json:"created_at"`
}

func (LoginLog) TableName() string { return "login_logs" }

// SetDetails 将 Details 字段序列化为 JSON 字符串
func (a *ActivityLog) SetDetails(data map[string]interface{}) {
	a.Details = data
}

// GetDetailsJSON 返回 Details 的 JSON 字符串
func (a *ActivityLog) GetDetailsJSON() string {
	if a.Details == nil {
		return "{}"
	}
	b, _ := json.Marshal(a.Details)
	return string(b)
}
