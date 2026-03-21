package models

import (
	"time"

	"gorm.io/gorm"
)

// PolicyConfig 策略配置文件（Wi-Fi/VPN/Email/证书/限制策略配置）
type PolicyConfig struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	ConfigType  string         `gorm:"type:varchar(32);not null;index" json:"config_type"` // wifi / vpn / email / certificate / restrictions
	SubType     string         `gorm:"type:varchar(64)" json:"sub_type"`                   // 子类型，如 WPA2/WPA3, OpenVPN/IPSec 等
	Description string         `gorm:"type:text" json:"description"`
	ConfigData  string         `gorm:"type:jsonb" json:"config_data"` // JSON 配置内容
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	Version     int            `gorm:"default:1" json:"version"` // 配置版本
	CreatedBy   string         `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedBy   string         `gorm:"type:varchar(64)" json:"updated_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PolicyConfig) TableName() string {
	return "policy_configs"
}

// Policy 策略主表
type Policy struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	Name        string         `gorm:"type:varchar(128);not null" json:"name"`
	PolicyType  string         `gorm:"type:varchar(32);not null;index" json:"policy_type"` // compliance / security / network / app / device
	Description string         `gorm:"type:text" json:"description"`
	Priority    int            `gorm:"default:0" json:"priority"`                 // 优先级，数字越大优先级越高
	ConfigIDs   string         `gorm:"type:jsonb;default:'[]'" json:"config_ids"` // 引用的配置文件 IDs
	RuleIDs     string         `gorm:"type:jsonb;default:'[]'" json:"rule_ids"`   // 引用的合规规则 IDs
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active / draft / archived
	Platform    string         `gorm:"type:varchar(32)" json:"platform"`                // 适用平台：ios / android / windows / mac / all
	Scope       string         `gorm:"type:varchar(32);default:'all'" json:"scope"`     // 适用范围：all / group / individual
	CreatedBy   string         `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedBy   string         `gorm:"type:varchar(64)" json:"updated_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Policy) TableName() string {
	return "policies"
}

// PolicyBinding 策略绑定表
type PolicyBinding struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TenantID   string     `gorm:"type:uuid;index" json:"tenant_id"` // 租户ID
	PolicyID   uint       `gorm:"not null;index" json:"policy_id"`
	TargetType string     `gorm:"type:varchar(32);not null;index" json:"target_type"` // device / user / group / org_unit
	TargetID   string     `gorm:"type:varchar(64);not null;index" json:"target_id"`   // 绑定目标 ID
	TargetName string     `gorm:"type:varchar(128)" json:"target_name"`               // 目标名称（冗余存储便于展示）
	BoundBy    string     `gorm:"type:varchar(64)" json:"bound_by"`                   // 绑定操作人
	BoundAt    time.Time  `json:"bound_at"`                                           // 绑定时间
	UnboundBy  string     `gorm:"type:varchar(64)" json:"unbound_by"`                 // 解绑操作人
	UnboundAt  *time.Time `json:"unbound_at"`                                         // 解绑时间
	Status     int        `gorm:"default:1" json:"status"`                            // 1:有效 0:已解绑
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (PolicyBinding) TableName() string {
	return "policy_bindings"
}

// PolicyBindingRequest 策略绑定请求
type PolicyBindingRequest struct {
	TargetType string   `json:"target_type" binding:"required"` // device / user / group / org_unit
	TargetIDs  []string `json:"target_ids" binding:"required"`  // 目标 ID 列表
}

// PolicyUnbindRequest 策略解绑请求
type PolicyUnbindRequest struct {
	TargetType string   `json:"target_type" binding:"required"`
	TargetIDs  []string `json:"target_ids" binding:"required"`
}
