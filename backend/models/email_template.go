package models

import (
	"time"
)

// EmailTemplate 邮件模板
type EmailTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`            // 模板名称
	Code      string    `gorm:"type:varchar(50);uniqueIndex" json:"code"`          // 模板编码
	Subject   string    `gorm:"type:varchar(255);not null" json:"subject"`        // 邮件主题
	Body      string    `gorm:"type:text;not null" json:"body"`                   // 邮件正文
	Variables string    `gorm:"type:text" json:"variables"`                       // 变量列表，JSON格式
	Type      string    `gorm:"type:varchar(50);default:'general'" json:"type"`  // 模板类型
	Status    int       `gorm:"default:1" json:"status"`                          // 状态: 1=active, 0=inactive
	Remark    string    `gorm:"type:text" json:"remark"`                           // 备注
	TenantID  string    `gorm:"type:uuid;index" json:"tenant_id"`                 // 租户ID
	CreatedBy uint      `json:"created_by"`                                       // 创建者
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (EmailTemplate) TableName() string {
	return "email_templates"
}
