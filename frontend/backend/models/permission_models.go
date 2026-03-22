package models

import (
	"time"
)

// Menu 菜单表
type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ParentID   *uint     `gorm:"index" json:"parent_id"`
	MenuName   string    `gorm:"size:100;not null" json:"menu_name"`
	MenuCode   string    `gorm:"size:100" json:"menu_code"`
	Icon       string    `gorm:"size:100" json:"icon"`
	RoutePath  string    `gorm:"size:255" json:"route_path"`
	Component  string    `gorm:"size:255" json:"component"`
	Permission string    `gorm:"size:100" json:"permission"`
	MenuType   int       `gorm:"default:1" json:"menu_type"` // 1=目录 2=菜单 3=按钮
	Sort       int       `gorm:"default:0" json:"sort"`
	Status     int       `gorm:"default:1" json:"status"`
	TenantID   string    `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Menu) TableName() string { return "menus" }

// ApiPermission API权限表
type ApiPermission struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ApiPath        string    `gorm:"size:255;not null" json:"api_path"`
	ApiName        string    `gorm:"size:100;not null" json:"api_name"`
	Method         string    `gorm:"size:20" json:"method"` // GET/POST/PUT/DELETE
	PermissionCode string    `gorm:"size:100" json:"permission_code"`
	MenuID         *uint     `gorm:"index" json:"menu_id"`
	Status         int       `gorm:"default:1" json:"status"`
	TenantID       string    `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ApiPermission) TableName() string { return "api_permissions" }

// Role 角色表（多租户版本）
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoleName    string    `gorm:"size:100;not null" json:"role_name"`
	RoleCode    string    `gorm:"size:100;not null" json:"role_code"`
	Description string    `gorm:"size:500" json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	TenantID    string    `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Role) TableName() string { return "roles" }

// PermGroup 权限组表（多租户版本）
type PermGroup struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GroupName   string    `gorm:"size:100;not null" json:"group_name"`
	GroupCode   string    `gorm:"size:100;not null" json:"group_code"`
	Description string    `gorm:"size:500" json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	TenantID    string    `gorm:"size:50;index" json:"tenant_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (PermGroup) TableName() string { return "permission_groups" }

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        uint      `gorm:"primaryKey"`
	RoleID    uint      `gorm:"index"`
	MenuID    uint      `gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

func (RoleMenu) TableName() string { return "role_menus" }

// RoleApiPermission 角色API权限关联表
type RoleApiPermission struct {
	ID              uint      `gorm:"primaryKey"`
	RoleID          uint      `gorm:"index"`
	ApiPermissionID uint      `gorm:"index"`
	CreatedAt       time.Time `json:"created_at"`
}

func (RoleApiPermission) TableName() string { return "role_api_permissions" }

// RolePermissionGroup 角色权限组关联表
type RolePermissionGroup struct {
	ID                 uint      `gorm:"primaryKey"`
	RoleID             uint      `gorm:"index"`
	PermissionGroupID  uint      `gorm:"index"`
	CreatedAt          time.Time `json:"created_at"`
}

func (RolePermissionGroup) TableName() string { return "role_permission_groups" }
