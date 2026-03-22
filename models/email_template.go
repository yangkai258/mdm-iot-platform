package models

import (
	"time"
)

// EmailTemplate 邮件模板
type EmailTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`                     // 模板名称
	Code      string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`           // 模板编码
	Subject   string    `gorm:"type:varchar(255);not null" json:"subject"`                  // 邮件主题
	Body      string    `gorm:"type:text;not null" json:"body"`                            // 邮件正文
	Variables string    `gorm:"type:text" json:"variables"`                                // 预定义变量（JSON数组格式）
	Remark    string    `gorm:"type:text" json:"remark"`                                    // 备注
	Status    int       `gorm:"type:int;default:1" json:"status"`                           // 状态：0-禁用 1-启用
	TenantID  string    `gorm:"type:varchar(50)" json:"tenant_id"`                          // 租户ID
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EmailTemplate) TableName() string {
	return "email_templates"
}
