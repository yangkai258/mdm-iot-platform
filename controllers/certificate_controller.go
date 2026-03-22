package controllers

import (
	"fmt"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// CertificateController 设备证书控制器
type CertificateController struct{}

func NewCertificateController() *CertificateController {
	return &CertificateController{}
}

// List 获取证书列表
// GET /api/v1/certificates
func (c *CertificateController) List(ctx *gin.Context) {
	var certs []models.DeviceCertificate
	query := models.DB.Model(&models.DeviceCertificate{})

	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if certType := ctx.Query("type"); certType != "" {
		query = query.Where("type = ?", certType)
	}

	if err := query.Order("created_at DESC").Limit(100).Find(&certs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": certs})
}

// Create 创建证书
// POST /api/v1/certificates
func (c *CertificateController) Create(ctx *gin.Context) {
	var input struct {
		DeviceID    string `json:"device_id" binding:"required"`
		Type       string `json:"type" binding:"required"`
		CommonName string `json:"common_name"`
		NotBefore  string `json:"not_before"`
		NotAfter   string `json:"not_after"`
		Serial     string `json:"serial"`
		Fingerprint string `json:"fingerprint"`
		Status     string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	cert := models.DeviceCertificate{
		CertID:       fmt.Sprintf("cert_%d", time.Now().UnixNano()),
		DeviceID:     input.DeviceID,
		Type:         input.Type,
		CommonName:   input.CommonName,
		Serial:       input.Serial,
		Fingerprint:  input.Fingerprint,
		Status:       "active",
	}
	if input.NotBefore != "" {
		cert.NotBefore = &[]time.Time{time.Now()}[0]
	}
	if input.NotAfter != "" {
		notAfter, _ := time.Parse("2006-01-02", input.NotAfter)
		cert.NotAfter = &notAfter
	}

	if err := models.DB.Create(&cert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": cert})
}

// Revoke 吊销证书
// DELETE /api/v1/certificates/:id
func (c *CertificateController) Revoke(ctx *gin.Context) {
	certID := ctx.Param("id")
	var cert models.DeviceCertificate
	if err := models.DB.Where("cert_id = ?", certID).First(&cert).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "certificate not found"})
		return
	}
	cert.Status = "revoked"
	models.DB.Save(&cert)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "certificate revoked"})
}

// Get 获取单个证书详情
// GET /api/v1/certificates/:id
func (c *CertificateController) Get(ctx *gin.Context) {
	certID := ctx.Param("id")
	var cert models.DeviceCertificate
	if err := models.DB.Where("cert_id = ?", certID).First(&cert).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "certificate not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": cert})
}
