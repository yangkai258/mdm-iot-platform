package models

import (
	"time"
)

// LDAPConfig LDAP 配置表
type LDAPConfig struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ConfigName   string     `gorm:"type:varchar(100)" json:"config_name"`      // 配置名称
	Host         string     `gorm:"type:varchar(255);not null" json:"host"`      // LDAP 服务器地址
	Port         int        `gorm:"default:389" json:"port"`                    // 端口
	BaseDN       string     `gorm:"type:varchar(255);not null" json:"base_dn"` // 基准 DN
	BindDN       string     `gorm:"type:varchar(255)" json:"bind_dn"`           // 管理员 DN
	BindPassword string     `gorm:"type:varchar(255)" json:"-"`                 // 加密后的密码
	UseSSL       bool       `gorm:"default:false" json:"use_ssl"`              // 使用 SSL
	UseTLS       bool       `gorm:"default:false" json:"use_tls"`              // 使用 STARTTLS
	UserFilter   string     `gorm:"type:varchar(500)" json:"user_filter"`       // 用户搜索过滤器
	GroupFilter  string     `gorm:"type:varchar(500)" json:"group_filter"`      // 分组搜索过滤器
	SyncInterval int        `gorm:"default:3600" json:"sync_interval"`          // 同步间隔(秒)
	IsEnabled    bool       `gorm:"default:false" json:"is_enabled"`            // 是否启用
	LastSyncAt   *time.Time `json:"last_sync_at"`                               // 最后同步时间
	Status       string     `gorm:"type:varchar(20);default:inactive" json:"status"` // active/inactive/error
	TenantID     string     `gorm:"type:varchar(50);index" json:"tenant_id"`   // 租户ID
	Description  string     `gorm:"type:text" json:"description"`               // 描述
	CreatedBy    uint       `json:"created_by"`                                 // 创建人
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (LDAPConfig) TableName() string {
	return "ldap_configs"
}

// LDAPUserMapping LDAP 用户映射表
type LDAPUserMapping struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	LDAPDN       string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"ldap_dn"` // LDAP DN
	LocalUserID  uint      `gorm:"index" json:"local_user_id"`                            // 本地用户ID
	Username     string    `gorm:"type:varchar(100)" json:"username"`
	Email        string    `gorm:"type:varchar(255)" json:"email"`
	DisplayName  string    `gorm:"type:varchar(100)" json:"display_name"`
	LDAPGroups   string    `gorm:"type:text" json:"ldap_groups"`           // JSON 数组
	SyncStatus   string    `gorm:"type:varchar(20);default:synced" json:"sync_status"` // synced/pending/removed
	LastSyncedAt *time.Time `json:"last_synced_at"`                        // 最后同步时间
	TenantID     string    `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (LDAPUserMapping) TableName() string {
	return "ldap_user_mappings"
}

// LDAPGroupRoleMapping LDAP 分组-角色映射表
type LDAPGroupRoleMapping struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	LDAPGroupDN   string    `gorm:"type:varchar(255);not null" json:"ldap_group_dn"`
	LDAPGroupName string    `gorm:"type:varchar(100)" json:"ldap_group_name"`
	RoleID        uint      `gorm:"not null" json:"role_id"`
	RoleName      string    `gorm:"type:varchar(100)" json:"role_name"` // 冗余存储
	TenantID      string    `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func (LDAPGroupRoleMapping) TableName() string {
	return "ldap_group_role_mappings"
}

// LDAPConfigRequest LDAP 配置请求
type LDAPConfigRequest struct {
	ConfigName   string `json:"config_name"`
	Host         string `json:"host" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	BaseDN       string `json:"base_dn" binding:"required"`
	BindDN       string `json:"bind_dn" binding:"required"`
	BindPassword string `json:"bind_password"` // 新密码时才传
	UseSSL       bool   `json:"use_ssl"`
	UseTLS       bool   `json:"use_tls"`
	UserFilter   string `json:"user_filter"`
	GroupFilter  string `json:"group_filter"`
	SyncInterval int    `json:"sync_interval"`
	IsEnabled    bool   `json:"is_enabled"`
	Description  string `json:"description"`
}

// LDAPGroupMappingRequest 分组-角色映射请求
type LDAPGroupMappingRequest struct {
	LDAPGroupDN   string `json:"ldap_group_dn" binding:"required"`
	LDAPGroupName string `json:"ldap_group_name"`
	RoleID        uint   `json:"role_id" binding:"required"`
}
