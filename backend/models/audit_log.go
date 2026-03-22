package models

import (
	"time"
)

// AuditLog 审计日志
type AuditLog struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Action        string         `gorm:"type:varchar(50);not null;index" json:"action"`          // 操作类型: encrypt, decrypt, anonymize, export, gdpr_access, gdpr_delete
	Module        string         `gorm:"type:varchar(50);not null;index" json:"module"`         // 模块: security, compliance, audit
	ResourceType  string         `gorm:"type:varchar(50)" json:"resource_type"`               // 资源类型: device, member, user, data
	ResourceID    string         `gorm:"type:varchar(100);index" json:"resource_id"`          // 资源ID
	UserID        uint           `gorm:"index" json:"user_id"`                                 // 操作人ID
	Username      string         `gorm:"type:varchar(100)" json:"username"`                   // 操作人用户名
	IP            string         `gorm:"type:varchar(45)" json:"ip"`                          // IP地址
	UserAgent     string         `gorm:"type:varchar(500)" json:"user_agent"`                 // User-Agent
	Status        int            `gorm:"default:1" json:"status"`                             // 1:成功 2:失败
	ErrorMsg      string         `gorm:"type:text" json:"error_msg"`                          // 错误信息
	RequestMethod string         `gorm:"type:varchar(10)" json:"request_method"`              // HTTP方法
	RequestPath   string         `gorm:"type:varchar(500)" json:"request_path"`                // 请求路径
	RequestBody   string         `gorm:"type:text" json:"request_body"`                        // 请求体（脱敏后）
	ResponseCode  int            `json:"response_code"`                                       // 响应码
	Duration      int            `json:"duration"`                                           // 操作耗时（毫秒）
	Metadata      string         `gorm:"type:json" json:"metadata"`                          // 额外元数据（JSON）
	TenantID      uint           `gorm:"index" json:"tenant_id"`                             // 租户ID
	CreatedAt     time.Time      `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// EncryptionKey 加密密钥管理
type EncryptionKey struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	KeyID          string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key_id"` // 密钥唯一标识
	KeyVersion     int       `gorm:"not null" json:"key_version"`                         // 密钥版本
	EncryptedKey   string    `gorm:"type:text;not null" json:"encrypted_key"`              // 加密后的密钥（用于密钥轮换）
	IsPrimary      bool      `gorm:"default:false" json:"is_primary"`                     // 是否为主密钥
	Status         int       `gorm:"default:1" json:"status"`                            // 1:激活 2:禁用 3:已轮换
	Algorithm      string    `gorm:"type:varchar(20);default:'AES-256-GCM'" json:"algorithm"`
	RotatedAt      *time.Time `json:"rotated_at"`                                         // 轮换时间
	ExpiresAt      *time.Time `json:"expires_at"`                                         // 过期时间
	RotatedBy      uint      `json:"rotated_by"`                                          // 轮换操作人
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (EncryptionKey) TableName() string {
	return "encryption_keys"
}

// DataAnonymizationRecord 数据脱敏记录
type DataAnonymizationRecord struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RecordID        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"record_id"` // 记录唯一标识
	OriginalData    string    `gorm:"type:text" json:"original_data"`                          // 原始数据（加密存储）
	AnonymizedData  string    `gorm:"type:text" json:"anonymized_data"`                        // 脱敏后数据
	AnonymizeType   string    `gorm:"type:varchar(50);not null" json:"anonymize_type"`         // 脱敏类型: email, phone, id_card, name, custom
	Fields          string    `gorm:"type:varchar(500)" json:"fields"`                         // 脱敏字段列表
	UserID          uint      `gorm:"index" json:"user_id"`                                    // 操作人ID
	Username        string    `gorm:"type:varchar(100)" json:"username"`
	Purpose         string    `gorm:"type:varchar(100)" json:"purpose"`                        // 使用目的
	ExportFormat    string    `gorm:"type:varchar(20)" json:"export_format"`                   // 导出格式: json, csv, xlsx
	Status          int       `gorm:"default:1" json:"status"`                               // 1:处理中 2:完成 3:失败
	CreatedAt       time.Time `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at"`
}

func (DataAnonymizationRecord) TableName() string {
	return "data_anonymization_records"
}

// GDPRRequest GDPR请求记录
type GDPRRequest struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RequestID       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"request_id"` // 请求唯一标识
	RequestType     string    `gorm:"type:varchar(50);not null" json:"request_type"`            // 请求类型: data_access, data_deletion, data_portability
	RequesterEmail  string    `gorm:"type:varchar(255);not null;index" json:"requester_email"`  // 请求者邮箱
	RequesterName   string    `gorm:"type:varchar(100)" json:"requester_name"`                   // 请求者姓名
	UserID          uint      `gorm:"index" json:"user_id"`                                     // 关联用户ID（如果有）
	Status          int       `gorm:"default:1" json:"status"`                                  // 1:待处理 2:处理中 3:已完成 4:已拒绝
	RequestReason   string    `gorm:"type:text" json:"request_reason"`                           // 请求原因
	ProcessedBy     uint      `json:"processed_by"`                                             // 处理人
	ProcessedAt     *time.Time `json:"processed_at"`                                            // 处理时间
	CompletedAt     *time.Time `json:"completed_at"`                                            // 完成时间
	RejectedReason  string    `gorm:"type:text" json:"rejected_reason"`                        // 拒绝原因
	ExportPath      string    `gorm:"type:varchar(500)" json:"export_path"`                     // 导出文件路径
	ResponseData    string    `gorm:"type:text" json:"response_data"`                           // 响应数据
	TenantID        uint      `gorm:"index" json:"tenant_id"`                                   // 租户ID
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (GDPRRequest) TableName() string {
	return "gdpr_requests"
}
