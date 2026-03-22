package models

import (
	"time"
)

type SecurityAudit struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EventType      string    `gorm:"type:varchar(50);not null;index" json:"event_type"`
	EventCategory  string    `gorm:"type:varchar(50);not null;index" json:"event_category"`
	Severity       int       `gorm:"default:1" json:"severity"`
	UserID         uint      `gorm:"index" json:"user_id"`
	Username       string    `gorm:"type:varchar(100);index" json:"username"`
	TargetUserID   uint      `gorm:"index" json:"target_user_id"`
	TargetUsername string    `gorm:"type:varchar(100)" json:"target_username"`
	IP             string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent      string    `gorm:"type:varchar(500)" json:"user_agent"`
	SessionID      string    `gorm:"type:varchar(100);index" json:"session_id"`
	ResourceType   string    `gorm:"type:varchar(50)" json:"resource_type"`
	ResourceID     string    `gorm:"type:varchar(100);index" json:"resource_id"`
	Action         string    `gorm:"type:varchar(100)" json:"action"`
	Status         int       `gorm:"default:1" json:"status"`
	ErrorMsg       string    `gorm:"type:text" json:"error_msg"`
	RequestMethod  string    `gorm:"type:varchar(10)" json:"request_method"`
	RequestPath    string    `gorm:"type:varchar(500)" json:"request_path"`
	ResponseCode   int       `json:"response_code"`
	Duration       int       `json:"duration"`
	Metadata       string    `gorm:"type:json" json:"metadata"`
	Country        string    `gorm:"type:varchar(50)" json:"country"`
	City           string    `gorm:"type:varchar(100)" json:"city"`
	TenantID       string    `gorm:"index" json:"tenant_id"`
	CreatedAt      time.Time `gorm:"index" json:"created_at"`
}

func (SecurityAudit) TableName() string {
	return "security_audits"
}

type SecurityReport struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ReportType      string    `gorm:"type:varchar(50);not null" json:"report_type"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	Summary         string    `gorm:"type:text" json:"summary"`
	Stats           string    `gorm:"type:json" json:"stats"`
	RiskLevel       int       `gorm:"default:1" json:"risk_level"`
	Findings        string    `gorm:"type:json" json:"findings"`
	Recommendations string    `gorm:"type:text" json:"recommendations"`
	Status          int       `gorm:"default:1" json:"status"`
	GeneratedBy     uint      `json:"generated_by"`
	GeneratedAt     time.Time `json:"generated_at"`
	TenantID        string    `gorm:"index" json:"tenant_id"`
	CreatedAt       time.Time `json:"created_at"`
}

func (SecurityReport) TableName() string {
	return "security_reports"
}
