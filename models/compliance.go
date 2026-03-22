package models

import (
	"time"
)

// CompliancePolicy 合规策略
type CompliancePolicy struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"type:varchar(100);not null" json:"name"`
	Description  string        `gorm:"type:varchar(255)" json:"description"`
	PolicyType   string        `gorm:"type:varchar(50);not null" json:"policy_type"` // firmware_version, battery_level, region_lock, encryption_required
	TargetValue  string        `gorm:"type:varchar(100)" json:"target_value"`       // 目标值，如版本号、最低电量等
	Condition    string        `gorm:"type:varchar(20);not null" json:"condition"`  // =, !=, >=, <=, <, >
	Severity     int           `gorm:"default:2" json:"severity"`                   // 1:低 2:中 3:高 4:严重
	RemediationAction string   `gorm:"type:varchar(50)" json:"remediation_action"`   // isolate, wipe, notify, block
	Enabled      bool          `gorm:"default:true" json:"enabled"`
	EnforceScope string        `gorm:"type:varchar(50);default:'all'" json:"enforce_scope"` // all, group, individual
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func (CompliancePolicy) TableName() string {
	return "compliance_policies"
}

// ComplianceViolation 合规违规记录
type ComplianceViolation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PolicyID      uint           `gorm:"not null;index" json:"policy_id"`
	DeviceID      string         `gorm:"type:varchar(36);index" json:"device_id"`
	PolicyType    string         `gorm:"type:varchar(50)" json:"policy_type"`
	ExpectedValue string         `gorm:"type:varchar(100)" json:"expected_value"`
	ActualValue   string         `gorm:"type:varchar(100)" json:"actual_value"`
	Severity      int            `json:"severity"`
	ActionTaken   string         `gorm:"type:varchar(50)" json:"action_taken"` // isolated, wiped, notified, blocked
	Status        int            `gorm:"default:1" json:"status"`              // 1:待处理 2:处理中 3:已解决 4:已忽略
	ResolvedAt    *time.Time     `json:"resolved_at"`
	ResolvedBy    string         `gorm:"type:varchar(36)" json:"resolved_by"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (ComplianceViolation) TableName() string {
	return "compliance_violations"
}
