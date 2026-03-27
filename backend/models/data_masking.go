package models

import (
	"time"
)

// DataMaskingRule 数据脱敏规则
type DataMaskingRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Field       string    `gorm:"type:varchar(50);not null" json:"field"`
	Pattern     string    `gorm:"type:varchar(255);not null" json:"pattern"`
	Replacement string    `gorm:"type:varchar(255);not null" json:"replacement"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	TenantID    string    `gorm:"type:varchar(50);index" json:"tenant_id"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DataMaskingRule) TableName() string {
	return "data_masking_rules"
}
