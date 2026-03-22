package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// DeviceSecurityController 设备安全控制器
type DeviceSecurityController struct{}

func NewDeviceSecurityController() *DeviceSecurityController {
	return &DeviceSecurityController{}
}

// GetSecurityStatus 获取设备安全状态
// GET /api/v1/devices/:device_id/security
func (c *DeviceSecurityController) GetSecurityStatus(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	var security models.DeviceSecurity
	if err := models.DB.Where("device_id = ?", deviceID).First(&security).Error; err != nil {
		// 返回默认安全状态
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
			"device_id":         deviceID,
			"encryption_enabled": false,
			"firewall_enabled":   false,
			"last_security_scan": nil,
			"security_level":    "basic",
		}})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": security})
}

// UpdateSecurityConfig 更新设备安全配置
// PUT /api/v1/devices/:device_id/security
func (c *DeviceSecurityController) UpdateSecurityConfig(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	var input struct {
		EncryptionEnabled bool   `json:"encryption_enabled"`
		FirewallEnabled   bool   `json:"firewall_enabled"`
		SecurityLevel     string `json:"security_level"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	var security models.DeviceSecurity
	if err := models.DB.Where("device_id = ?", deviceID).FirstOrCreate(&security, models.DeviceSecurity{DeviceID: deviceID}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	security.EncryptionEnabled = input.EncryptionEnabled
	security.FirewallEnabled = input.FirewallEnabled
	security.SecurityLevel = input.SecurityLevel
	models.DB.Save(&security)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "security config updated", "data": security})
}

// ScanSecurity 执行安全扫描
// POST /api/v1/devices/:device_id/security/scan
func (c *DeviceSecurityController) ScanSecurity(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "security scan triggered", "device_id": deviceID})
}

// GetAuditLog 获取安全审计日志
// GET /api/v1/devices/:device_id/security/audit
func (c *DeviceSecurityController) GetAuditLog(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	var logs []models.SecurityAuditLog
	if err := models.DB.Where("device_id = ?", deviceID).Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": logs})
}
