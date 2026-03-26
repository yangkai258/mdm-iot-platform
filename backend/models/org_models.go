package models

import (
	"time"

	"gorm.io/gorm"
)

// Company 公司表
type Company struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`                    // 租户ID
	CompanyCode  string         `gorm:"size:50;uniqueIndex;not null" json:"company_code"`  // 公司编码
	CompanyName  string         `gorm:"size:200;not null" json:"company_name"`             // 公司名称
	ShortName    string         `gorm:"size:100" json:"short_name"`                       // 简称
	Logo         string         `gorm:"size:500" json:"logo"`                             // Logo
	Province     string         `gorm:"size:50" json:"province"`                         // 省
	City         string         `gorm:"size:50" json:"city"`                             // 市
	District     string         `gorm:"size:50" json:"district"`                         // 区
	Address      string         `gorm:"size:500" json:"address"`                          // 详细地址
	LegalPerson  string         `gorm:"size:50" json:"legal_person"`                      // 法人代表
	Contact      string         `gorm:"size:50" json:"contact"`                          // 联系人
	Phone        string         `gorm:"size:20" json:"phone"`                             // 联系电话
	Email        string         `gorm:"size:100" json:"email"`                            // 邮箱
	Status       int            `gorm:"default:1" json:"status"`                           // 状态: 1正常 0禁用
	Sort         int            `gorm:"default:0" json:"sort"`                             // 排序
	Remark       string         `gorm:"size:500" json:"remark"`                           // 备注
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Department 部门表
type Department struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     *string        `gorm:"column:tenant_id;type:uuid;index" json:"tenant_id"`        // 租户ID
	DeptCode     string         `gorm:"column:dept_code;size:50;not null" json:"dept_code"`      // 部门编码
	DeptName     string         `gorm:"column:dept_name;size:100;not null" json:"dept_name"`     // 部门名称
	ParentID     *uint          `gorm:"column:parent_id;index" json:"parent_id"`                  // 上级部门
	Level        int            `gorm:"default:1" json:"level"`                  // 层级
	Path         string         `gorm:"size:500" json:"path"`                    // 路径
	ManagerID    uint           `gorm:"column:manager_id" json:"manager_id"`     // 负责人ID
	Phone        string         `gorm:"size:20" json:"phone"`                    // 联系电话
	Email        string         `gorm:"size:100" json:"email"`                   // 邮箱
	Status       string         `gorm:"column:status;default:'active'" json:"status"`                 // 状态
	SortOrder    int            `gorm:"column:sort_order;default:0" json:"sort_order"`                   // 排序
	CompanyID    uint           `gorm:"index" json:"company_id"`                 // 所属公司
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Children     []Department   `gorm:"-" json:"children,omitempty"`             // 子部门
}

// Position 岗位表
type Position struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`       // 租户ID
	PosCode      string         `gorm:"size:50;not null" json:"pos_code"`       // 岗位编码
	PosName      string         `gorm:"size:100;not null" json:"pos_name"`      // 岗位名称
	Category     string         `gorm:"size:50" json:"category"`               // 岗位类别
	Level        int            `gorm:"default:1" json:"level"`                 // 级别
	DeptID       *uint          `gorm:"index" json:"dept_id"`                   // 所属部门
	CompanyID    uint           `gorm:"index" json:"company_id"`                 // 所属公司
	Description  string         `gorm:"size:500" json:"description"`            // 岗位描述
	Status       int            `gorm:"default:1" json:"status"`                // 状态
	Sort         int            `gorm:"default:0" json:"sort"`                   // 排序
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Employee 员工表
type Employee struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"size:50;index" json:"tenant_id"`                 // 租户ID
	EmpCode      string         `gorm:"size:50;uniqueIndex;not null" json:"emp_code"`    // 工号
	EmpName      string         `gorm:"size:50;not null" json:"emp_name"`               // 姓名
	Gender       string         `gorm:"size:10" json:"gender"`                          // 性别
	BirthDate    *time.Time     `json:"birth_date"`                                    // 出生日期
	Phone        string         `gorm:"size:20" json:"phone"`                           // 手机号
	Email        string         `gorm:"size:100" json:"email"`                          // 邮箱
	IDCard       string         `gorm:"size:20" json:"id_card"`                         // 身份证号
	Photo        string         `gorm:"size:500" json:"photo"`                         // 照片
	Province     string         `gorm:"size:50" json:"province"`                       // 省
	City         string         `gorm:"size:50" json:"city"`                           // 市
	District     string         `gorm:"size:50" json:"district"`                       // 区
	Address      string         `gorm:"size:500" json:"address"`                        // 详细地址
	DeptID       uint           `gorm:"index" json:"dept_id"`                          // 部门
	PositionID   uint           `gorm:"index" json:"position_id"`                      // 岗位
	CompanyID    uint           `gorm:"index" json:"company_id"`                       // 公司
	EntryDate    *time.Time     `json:"entry_date"`                                   // 入职日期
	EmpStatus    int            `gorm:"default:1" json:"emp_status"`                   // 员工状态
	Status       int            `gorm:"default:1" json:"status"`                        // 账号状态
	Remark       string         `gorm:"size:500" json:"remark"`                       // 备注
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Department   *Department    `gorm:"foreignKey:DeptID" json:"department,omitempty"`
	Position     *Position      `gorm:"foreignKey:PositionID" json:"position,omitempty"`
}

// StandardPosition 基准岗位表
type StandardPosition struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SpCode       string         `gorm:"size:50;not null" json:"sp_code"`      // 编码
	SpName       string         `gorm:"size:100;not null" json:"sp_name"`     // 名称
	Category     string         `gorm:"size:50" json:"category"`             // 类别
	Level        int            `gorm:"default:1" json:"level"`              // 级别
	Description  string         `gorm:"size:1000" json:"description"`       // 描述
	Responsibility string       `gorm:"type:text" json:"responsibility"`     // 职责
	Requirement  string         `gorm:"type:text" json:"requirement"`       // 任职要求
	Skills       string         `gorm:"size:500" json:"skills"`             // 技能要求
	Status       int            `gorm:"default:1" json:"status"`             // 状态
	Sort         int            `gorm:"default:0" json:"sort"`               // 排序
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// PositionTemplate 基准岗位模板表
type PositionTemplate struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TenantID      string         `gorm:"size:50;index" json:"tenant_id"`     // 租户ID
	Name          string         `gorm:"size:100;not null" json:"name"`     // 模板名称
	Code          string         `gorm:"size:50;not null" json:"code"`        // 模板编码
	Description   string         `gorm:"size:500" json:"description"`          // 描述
	Permissions   string         `gorm:"type:text" json:"permissions"`        // 权限code列表，JSON数组
	Status        int            `gorm:"default:1" json:"status"`           // 1启用 0禁用
	InheritedFrom *uint          `gorm:"index" json:"inherited_from"`        // 继承自模板ID（复制/继承时设置）
	ParentID      *uint          `gorm:"index" json:"parent_id"`              // 父模板ID（用于模板层级）
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// SysUser 系统用户表(扩展)
type SysUserExt struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"uniqueIndex;not null" json:"user_id"`    // 用户ID
	EmployeeID   uint           `gorm:"index" json:"employee_id"`                 // 员工ID
	DeptID       uint           `gorm:"index" json:"dept_id"`                    // 部门ID
	CompanyID    uint           `gorm:"index" json:"company_id"`                 // 公司ID
	RoleIDs      string         `gorm:"size:500" json:"role_ids"`               // 角色IDs
	DataScope    int            `gorm:"default:1" json:"data_scope"`             // 数据权限范围
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
