package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ Sprint 32: 高级安全功能模型 ============

// AuditReport 审计报告
type AuditReport struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ReportName  string         `gorm:"type:varchar(200);not null" json:"report_name"` // 报告名称
	ReportType  string         `gorm:"type:varchar(50);not null" json:"report_type"`   // 报告类型: security, compliance, access, operations
	PeriodStart time.Time      `gorm:"not null" json:"period_start"`                   // 报告周期开始
	PeriodEnd   time.Time      `gorm:"not null" json:"period_end"`                   // 报告周期结束
	FilePath    string         `gorm:"type:varchar(500)" json:"file_path"`           // 报告文件路径
	FileSize    int64          `json:"file_size"`                                    // 文件大小（字节）
	Format      string         `gorm:"type:varchar(20);default:'pdf'" json:"format"` // 报告格式: pdf, xlsx, csv, json
	Status      int            `gorm:"default:1" json:"status"`                      // 1:生成中 2:已完成 3:失败
	GeneratedBy uint           `json:"generated_by"`                                   // 生成人
	Summary     string         `gorm:"type:text" json:"summary"`                     // 报告摘要
	Detail      string         `gorm:"type:json" json:"detail"`                      // 详细数据（JSON）
	TenantID    uint           `gorm:"index" json:"tenant_id"`                      // 租户ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuditReport) TableName() string {
	return "audit_reports"
}

// ComplianceReport 合规报告
type ComplianceReport struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ReportName     string         `gorm:"type:varchar(200);not null" json:"report_name"`      // 报告名称
	RegulationType string         `gorm:"type:varchar(50);not null" json:"regulation_type"` // 法规类型: gdpr, ccpa, hipaa, sox, iso27001
	Scope          string         `gorm:"type:varchar(100)" json:"scope"`                  // 审计范围
	PeriodStart    time.Time      `gorm:"not null" json:"period_start"`                      // 报告周期开始
	PeriodEnd      time.Time      `gorm:"not null" json:"period_end"`                      // 报告周期结束
	Status         int            `gorm:"default:1" json:"status"`                         // 1:草稿 2:审核中 3:已通过 4:已拒绝
	Score          float64        `json:"score"`                                           // 合规评分（0-100）
	ViolationsCount int           `json:"violations_count"`                                // 违规项数量
	FilePath       string         `gorm:"type:varchar(500)" json:"file_path"`              // 报告文件路径
	Findings       string         `gorm:"type:text" json:"findings"`                        // 调查发现
	Recommendations string       `gorm:"type:text" json:"recommendations"`                // 改进建议
	ApprovedBy     uint           `json:"approved_by"`                                     // 审批人
	ApprovedAt     *time.Time     `json:"approved_at"`                                     // 审批时间
	TenantID       uint           `gorm:"index" json:"tenant_id"`                         // 租户ID
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ComplianceReport) TableName() string {
	return "compliance_reports"
}

// DataExport 数据导出任务
type DataExport struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ExportID    string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"export_id"` // 导出任务唯一标识
	ExportType  string         `gorm:"type:varchar(50);not null" json:"export_type"`           // 导出类型: full, partial, audit, compliance
	DataTypes   string         `gorm:"type:varchar(200)" json:"data_types"`                    // 数据类型: devices, members, activity, all
	Format      string         `gorm:"type:varchar(20);default:'json'" json:"format"`         // 导出格式: json, csv, xlsx, zip
	FilePath    string         `gorm:"type:varchar(500)" json:"file_path"`                     // 导出文件路径
	FileSize    int64          `json:"file_size"`                                            // 文件大小（字节）
	Filters     string         `gorm:"type:json" json:"filters"`                             // 导出过滤条件（JSON）
	Status      int            `gorm:"default:1" json:"status"`                              // 1:等待中 2:处理中 3:完成 4:失败 5:已过期
	RequesterID uint           `json:"requester_id"`                                         // 请求人ID
	RequesterEmail string      `gorm:"type:varchar(255)" json:"requester_email"`              // 请求人邮箱
	DownloadCount int          `gorm:"default:0" json:"download_count"`                      // 下载次数
	ExpiresAt   *time.Time     `json:"expires_at"`                                           // 过期时间
	ErrorMsg    string         `gorm:"type:text" json:"error_msg"`                          // 错误信息
	TenantID    uint           `gorm:"index" json:"tenant_id"`                               // 租户ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DataExport) TableName() string {
	return "data_exports"
}

// GDPRRequest GDPR请求（扩展自 audit_log.go 中的定义，增加补充字段）
// 注意：基础字段已在 audit_log.go 的 GDPRRequest 中定义
// 此处仅扩展额外字段供关联查询使用
type GDPRRequestExtra struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	GDPRRequestID   uint      `gorm:"uniqueIndex;not null" json:"gdpr_request_id"` // 关联 audit_log.go 中的 GDPRRequest
	VerifiedAt      *time.Time `json:"verified_at"`                                 // 身份验证时间
	IdentityDoc     string    `gorm:"type:varchar(100)" json:"identity_doc"`       // 身份验证文档类型
	ProcessingNotes string    `gorm:"type:text" json:"processing_notes"`           // 处理备注
	CompletionNotes string    `gorm:"type:text" json:"completion_notes"`           // 完成备注
	Signature       string    `gorm:"type:text" json:"signature"`                 // 电子签名（可选）
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (GDPRRequestExtra) TableName() string {
	return "gdpr_request_extras"
}

// ConsentRecord 同意记录
type ConsentRecord struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ConsentType   string         `gorm:"type:varchar(50);not null" json:"consent_type"`    // 同意类型: privacy, terms, marketing, data_processing, third_party_sharing
	ConsentAction string         `gorm:"type:varchar(20);not null" json:"consent_action"` // 同意动作: granted, withdrawn, updated
	UserID        uint           `gorm:"index" json:"user_id"`                            // 用户ID
	UserEmail     string         `gorm:"type:varchar(255);index" json:"user_email"`      // 用户邮箱
	DeviceID      string         `gorm:"type:varchar(36);index" json:"device_id"`       // 设备ID
	IPAddress     string         `gorm:"type:varchar(45)" json:"ip_address"`             // IP地址
	UserAgent     string         `gorm:"type:varchar(500)" json:"user_agent"`            // User-Agent
	Version       string         `gorm:"type:varchar(50)" json:"version"`               // 协议版本
	PolicyURL     string         `gorm:"type:varchar(500)" json:"policy_url"`            // 政策URL
	ConsentProof  string         `gorm:"type:text" json:"consent_proof"`                 // 同意证明（JSON格式）
	WithdrawnAt   *time.Time     `json:"withdrawn_at"`                                   // 撤回时间
	TenantID      uint           `gorm:"index" json:"tenant_id"`                         // 租户ID
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ConsentRecord) TableName() string {
	return "consent_records"
}
