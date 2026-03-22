package models

import (
	"time"
)

// SysUser 系统用户
type SysUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Nickname  string    `gorm:"type:varchar(50)" json:"nickname"`
	Email     string    `gorm:"type:varchar(100)" json:"email"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Status    int       `gorm:"default:1" json:"status"`
	RoleID    uint      `json:"role_id"`
	TenantID  string    `gorm:"index" json:"tenant_id"` // 租户ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SysUser) TableName() string { return "sys_users" }

// SysRole 系统角色
type SysRole struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SysRole) TableName() string { return "sys_roles" }

// SysMenu 系统菜单
type SysMenu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ParentID   uint      `gorm:"default:0" json:"parent_id"`
	Name       string    `gorm:"type:varchar(50);not null" json:"name"`
	Path       string    `gorm:"type:varchar(255)" json:"path"`
	Component  string    `gorm:"type:varchar(255)" json:"component"`
	Icon       string    `gorm:"type:varchar(50)" json:"icon"`
	Sort       int       `gorm:"default:0" json:"sort"`
	Visible    int       `gorm:"default:1" json:"visible"`
	Permission string    `gorm:"type:varchar(100)" json:"permission"`
	Type       int       `gorm:"default:1" json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (SysMenu) TableName() string { return "sys_menus" }

// SysDictionary 字典表
type SysDictionary struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"type"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Label     string    `gorm:"type:varchar(100)" json:"label"`
	Value     string    `gorm:"type:varchar(100)" json:"value"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int       `gorm:"default:1" json:"status"`
	Remark    string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SysDictionary) TableName() string { return "sys_dictionaries" }

// SysOperationLog 操作日志
type SysOperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `gorm:"type:varchar(50)" json:"username"`
	Module     string    `gorm:"type:varchar(50)" json:"module"`
	Operation  string    `gorm:"type:varchar(50)" json:"operation"`
	Method     string    `gorm:"type:varchar(10)" json:"method"`
	Path       string    `gorm:"type:varchar(255)" json:"path"`
	IP         string    `gorm:"type:varchar(50)" json:"ip"`
	Location   string    `gorm:"type:varchar(255)" json:"location"`
	Params     string    `gorm:"type:text" json:"params"`
	Result     string    `gorm:"type:text" json:"result"`
	Status     int       `gorm:"default:1" json:"status"`
	ErrorMsg   string    `gorm:"type:text" json:"error_msg"`
	Duration   int       `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

func (SysOperationLog) TableName() string { return "sys_operation_logs" }

// SysLoginLog 登录日志
type SysLoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `gorm:"type:varchar(50)" json:"username"`
	IP        string    `gorm:"type:varchar(50)" json:"ip"`
	Location  string    `gorm:"type:varchar(255)" json:"location"`
	Browser   string    `gorm:"type:varchar(50)" json:"browser"`
	OS        string    `gorm:"type:varchar(50)" json:"os"`
	Status    int       `json:"status"`
	Msg       string    `gorm:"type:varchar(255)" json:"msg"`
	LoginTime time.Time `json:"login_time"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysLoginLog) TableName() string { return "sys_login_logs" }

// Knowledge 知识库条目
type Knowledge struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Category string    `gorm:"type:varchar(50);index" json:"category"`
	Question string    `gorm:"type:text;not null" json:"question"`
	Answer   string    `gorm:"type:text;not null" json:"answer"`
	Status   int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Knowledge) TableName() string { return "knowledge" }
