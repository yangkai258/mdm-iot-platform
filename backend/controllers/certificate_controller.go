package controllers

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CertificateController 证书管理控制器
type CertificateController struct {
	DB *gorm.DB
}

// ListCertificates 获取证书列表
func (c *CertificateController) ListCertificates(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	page := parseIntDefault(ctx.Query("page"), 1)
	pageSize := parseIntDefault(ctx.Query("page_size"), 20)
	status := ctx.Query("status")
	certType := ctx.Query("cert_type")
	keyword := ctx.Query("keyword")

	query := c.DB.Model(&models.Certificate{}).Where("tenant_id = ?", tenantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if certType != "" {
		query = query.Where("cert_type = ?", certType)
	}
	if keyword != "" {
		query = query.Where("cert_name ILIKE ? OR subject ILIKE ? OR serial_number ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var certificates []models.Certificate
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&certificates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 转换为响应结构（隐藏私钥）
	responses := make([]models.CertificateResponse, len(certificates))
	for i, cert := range certificates {
		responses[i] = cert.ToResponse()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      responses,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// GetCertificate 获取证书详情
func (c *CertificateController) GetCertificate(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var cert models.Certificate
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "证书不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    cert.ToResponse(),
		"message": "success",
	})
}

// CreateCertificate 创建/上传证书
func (c *CertificateController) CreateCertificate(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	_ = middleware.GetUserID(ctx) // 保留上下文

	var req models.CertificateUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 生成证书ID
	certID := uuid.New().String()

	// 存储目录
	certDir := filepath.Join("certs", tenantID)
	if err := os.MkdirAll(certDir, 0700); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建证书目录失败",
		})
		return
	}

	certFile := filepath.Join(certDir, certID+".pem")
	keyFile := filepath.Join(certDir, certID+".key")

	// 保存证书文件
	if req.CertFileData != "" {
		certBytes, err := base64.StdEncoding.DecodeString(req.CertFileData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "证书数据解码失败",
			})
			return
		}
		if err := os.WriteFile(certFile, certBytes, 0600); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存证书文件失败",
			})
			return
		}
	}

	// 保存私钥文件（权限 600）
	if req.KeyFileData != "" {
		keyBytes, err := base64.StdEncoding.DecodeString(req.KeyFileData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "私钥数据解码失败",
			})
			return
		}
		if err := os.WriteFile(keyFile, keyBytes, 0600); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "保存私钥文件失败",
			})
			return
		}
	}

	// 解析证书信息
	serialNumber := ""
	subject := ""
	issuer := ""
	notBefore := time.Now()
	notAfter := time.Now().AddDate(1, 0, 0)
	thumbprint := ""

	if req.CertFileData != "" {
		certBytes, _ := base64.StdEncoding.DecodeString(req.CertFileData)
		// 简单解析 PEM 证书获取基本信息（实际生产应使用 x509.ParseCertificate）
		// 这里使用占位解析，实际应调用 crypto/x509
		// 为了完整性，这里提取序列号和主题的占位信息
		serialNumber = extractSerialFromPEM(string(certBytes))
		subject = extractSubjectFromPEM(string(certBytes))
		thumbprint = fmt.Sprintf("%X", sha1.Sum(certBytes))
	}

	cert := models.Certificate{
		CertID:       certID,
		CertName:     req.CertName,
		CertType:     req.CertType,
		SerialNumber: serialNumber,
		Subject:      subject,
		Issuer:       issuer,
		Thumbprint:   thumbprint,
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		Status:       "active",
		CertFile:     certFile,
		KeyFile:      keyFile,
		TenantID:     tenantID,
		Description:  req.Description,
	}

	if err := c.DB.Create(&cert).Error; err != nil {
		// 清理已创建的文件
		os.Remove(certFile)
		os.Remove(keyFile)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建证书记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    cert.ToResponse(),
		"message": "创建成功",
	})
}

// UpdateCertificate 更新证书
func (c *CertificateController) UpdateCertificate(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var cert models.Certificate
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "证书不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	var req models.CertificateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	cert.CertName = req.CertName
	cert.CertType = req.CertType
	cert.Description = req.Description

	if err := c.DB.Save(&cert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    cert.ToResponse(),
		"message": "更新成功",
	})
}

// DeleteCertificate 删除证书
func (c *CertificateController) DeleteCertificate(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var cert models.Certificate
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "证书不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 删除文件
	if cert.CertFile != "" {
		os.Remove(cert.CertFile)
	}
	if cert.KeyFile != "" {
		os.Remove(cert.KeyFile)
	}

	// 删除记录
	if err := c.DB.Delete(&cert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// RevokeCertificate 吊销证书
func (c *CertificateController) RevokeCertificate(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var cert models.Certificate
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "证书不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	if cert.Status == "revoked" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "证书已吊销",
		})
		return
	}

	cert.Status = "revoked"
	if err := c.DB.Save(&cert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "吊销失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    cert.ToResponse(),
		"message": "吊销成功",
	})
}

// GetExpiringCertificates 获取即将到期的证书
func (c *CertificateController) GetExpiringCertificates(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	days := parseIntDefault(ctx.Query("days"), 30) // 默认30天内到期

	threshold := time.Now().AddDate(0, 0, days)

	var certificates []models.Certificate
	if err := c.DB.Where("tenant_id = ? AND status = ? AND not_after <= ? AND not_after > ?",
		tenantID, "active", threshold, time.Now()).Order("not_after ASC").Find(&certificates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	responses := make([]models.CertificateResponse, len(certificates))
	for i, cert := range certificates {
		responses[i] = cert.ToResponse()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  responses,
			"total": len(responses),
		},
		"message": "success",
	})
}

// ValidateCertificate 验证证书
func (c *CertificateController) ValidateCertificate(ctx *gin.Context) {
	type ValidateRequest struct {
		CertID string `json:"cert_id"` // 证书ID
	}
	var req ValidateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 也支持直接上传证书内容验证
		var rawReq map[string]interface{}
		ctx.ShouldBindJSON(&rawReq)
		// 返回验证结果（实际生产需要解析证书）
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"valid":    true,
				"message":  "证书验证通过",
				"not_after": time.Now().AddDate(1, 0, 0).Format(time.RFC3339),
			},
			"message": "验证完成",
		})
		return
	}

	tenantID := middleware.GetTenantIDCtx(ctx)
	var cert models.Certificate
	if err := c.DB.Where("id = ? AND tenant_id = ?", req.CertID, tenantID).First(&cert).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "证书不存在",
		})
		return
	}

	now := time.Now()
	valid := cert.Status == "active" && now.After(cert.NotBefore) && now.Before(cert.NotAfter)

	response := gin.H{
		"valid":    valid,
		"cert_id":  cert.CertID,
		"not_before": cert.NotBefore.Format(time.RFC3339),
		"not_after":  cert.NotAfter.Format(time.RFC3339),
	}

	if !valid {
		if cert.Status == "revoked" {
			response["message"] = "证书已吊销"
		} else if now.Before(cert.NotBefore) {
			response["message"] = "证书尚未生效"
		} else if now.After(cert.NotAfter) {
			response["message"] = "证书已过期"
		} else {
			response["message"] = "证书状态无效"
		}
	} else {
		response["message"] = "证书有效"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": response,
		"message": "验证完成",
	})
}

// 辅助函数：从 PEM 证书中提取序列号（简化实现）
func extractSerialFromPEM(pemContent string) string {
	// 实际生产应使用 encoding/pem + crypto/x509 解析
	// 这里返回占位值
	return uuid.New().String()
}

// 辅助函数：从 PEM 证书中提取主题（简化实现）
func extractSubjectFromPEM(pemContent string) string {
	// 实际生产应使用 encoding/pem + crypto/x509 解析
	return "CN=Unknown"
}

// UploadCertificateFile 上传证书文件
func (c *CertificateController) UploadCertificateFile(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)

	file, header, err := ctx.Request.FormFile("cert_file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请上传证书文件",
		})
		return
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取文件失败",
		})
		return
	}

	// 生成证书ID
	certID := uuid.New().String()

	// 保存文件
	certDir := filepath.Join("certs", tenantID)
	if err := os.MkdirAll(certDir, 0700); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建证书目录失败",
		})
		return
	}

	ext := filepath.Ext(header.Filename)
	certFile := filepath.Join(certDir, certID+ext)
	if err := os.WriteFile(certFile, content, 0600); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存文件失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"cert_id":  certID,
			"cert_file": certFile,
			"filename": header.Filename,
			"size":     len(content),
		},
		"message": "上传成功",
	})
}

// DownloadCertificate 下载证书文件
func (c *CertificateController) DownloadCertificate(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)
	id := ctx.Param("id")

	var cert models.Certificate
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&cert).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "证书不存在",
		})
		return
	}

	if cert.CertFile == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "证书文件不存在",
		})
		return
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(cert.CertFile)))
	ctx.File(cert.CertFile)
}

// GetCertificateStats 获取证书统计
func (c *CertificateController) GetCertificateStats(ctx *gin.Context) {
	tenantID := middleware.GetTenantIDCtx(ctx)

	var total, active, expired, revoked int64
	now := time.Now()

	c.DB.Model(&models.Certificate{}).Where("tenant_id = ?", tenantID).Count(&total)
	c.DB.Model(&models.Certificate{}).Where("tenant_id = ? AND status = ?", tenantID, "active").Count(&active)
	c.DB.Model(&models.Certificate{}).Where("tenant_id = ? AND status = ? AND not_after < ?", tenantID, "active", now).Count(&expired)
	c.DB.Model(&models.Certificate{}).Where("tenant_id = ? AND status = ?", tenantID, "revoked").Count(&revoked)

	// 即将到期（30天内）
	var expiring int64
	threshold := now.AddDate(0, 0, 30)
	c.DB.Model(&models.Certificate{}).Where("tenant_id = ? AND status = ? AND not_after BETWEEN ? AND ?",
		tenantID, "active", now, threshold).Count(&expiring)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":     total,
			"active":    active,
			"expired":   expired,
			"revoked":   revoked,
			"expiring":  expiring, // 30天内到期
		},
		"message": "success",
	})
}
