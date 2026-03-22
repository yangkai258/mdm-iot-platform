package models

import (
	"time"
)

// Certificate 设备证书表
type Certificate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CertID       string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"cert_id"` // 证书唯一标识
	CertName     string    `gorm:"type:varchar(128)" json:"cert_name"`                    // 证书名称
	CertType     string    `gorm:"type:varchar(32)" json:"cert_type"`                    // device/client/server/ca
	SerialNumber string    `gorm:"type:varchar(128);index" json:"serial_number"`         // 证书序列号
	Subject      string    `gorm:"type:varchar(256)" json:"subject"`                    // 主题 (CN=xxx)
	Issuer       string    `gorm:"type:varchar(256)" json:"issuer"`                     // 颁发者
	Thumbprint   string    `gorm:"type:varchar(64);uniqueIndex" json:"thumbprint"`       // SHA1 指纹
	NotBefore    time.Time `json:"not_before"`                                          // 生效时间
	NotAfter     time.Time `json:"not_after"`                                            // 过期时间
	Status       string    `gorm:"type:varchar(20);default:active;index" json:"status"` // active/expired/revoked/pending
	CertFile     string    `gorm:"type:varchar(512)" json:"cert_file"`                   // 证书文件路径 (公钥)
	KeyFile      string    `gorm:"type:varchar(512)" json:"-"`                          // 私钥文件路径 (敏感，不返回)
	TenantID     string    `gorm:"type:varchar(50);index" json:"tenant_id"`              // 租户ID
	Description  string    `gorm:"type:text" json:"description"`                        // 描述
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Certificate) TableName() string {
	return "certificates"
}

// CertificateRequest 创建/更新证书请求
type CertificateRequest struct {
	CertName    string `json:"cert_name" binding:"required"`
	CertType    string `json:"cert_type" binding:"required"` // device/client/server/ca
	Description string `json:"description"`
}

// CertificateUploadRequest 上传证书文件请求
type CertificateUploadRequest struct {
	CertName     string `json:"cert_name" binding:"required"`
	CertType     string `json:"cert_type" binding:"required"`
	CertFileData string `json:"cert_file_data"` // PEM 格式证书内容 (base64)
	KeyFileData  string `json:"key_file_data"`  // PEM 格式私钥内容 (base64, 可选)
	Description  string `json:"description"`
}

// CertificateResponse 证书列表响应 (隐藏私钥)
type CertificateResponse struct {
	ID           uint      `json:"id"`
	CertID       string    `json:"cert_id"`
	CertName     string    `json:"cert_name"`
	CertType     string    `json:"cert_type"`
	SerialNumber string    `json:"serial_number"`
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	Thumbprint   string    `json:"thumbprint"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	Status       string    `json:"status"`
	CertFile     string    `json:"cert_file"`
	TenantID     string    `json:"tenant_id"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse 转换为响应结构 (不暴露私钥)
func (c *Certificate) ToResponse() CertificateResponse {
	return CertificateResponse{
		ID:           c.ID,
		CertID:       c.CertID,
		CertName:     c.CertName,
		CertType:     c.CertType,
		SerialNumber: c.SerialNumber,
		Subject:      c.Subject,
		Issuer:       c.Issuer,
		Thumbprint:   c.Thumbprint,
		NotBefore:    c.NotBefore,
		NotAfter:     c.NotAfter,
		Status:       c.Status,
		CertFile:     c.CertFile,
		TenantID:     c.TenantID,
		Description:  c.Description,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
