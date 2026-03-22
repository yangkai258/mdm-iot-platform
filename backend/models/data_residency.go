package models

import (
	"time"
)

// DataResidencyRule 数据驻留规则
type DataResidencyRule struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RuleID      string    `json:"rule_id" gorm:"uniqueIndex;size:64"`
	DataType    string    `json:"data_type" gorm:"size:32"`    // user/device/pet/alert/log
	RegionCode  string    `json:"region_code" gorm:"index;size:16"` // 必须存储的区域
	TenantID    uint      `json:"tenant_id" gorm:"index"`
	Description string    `json:"description" gorm:"size:256"`
	Status      string    `json:"status" gorm:"size:20"` // active/inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
