package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"mdm-backend/middleware"
	"mdm-backend/models"
	"gorm.io/gorm"
)

// DevicePairingController 设备配对控制器
type DevicePairingController struct {
	DB *gorm.DB
}

func NewDevicePairingController(db *gorm.DB) *DevicePairingController {
	return &DevicePairingController{DB: db}
}

// RegisterRoutes 注册设备配对路由
func (ctrl *DevicePairingController) RegisterRoutes(rg *gin.RouterGroup) {
	pairing := rg.Group("/devices/pairing")
	{
		pairing.POST("/code", ctrl.GeneratePairingCode)
		pairing.POST("/verify", ctrl.VerifyPairingCode)
		pairing.GET("/history", ctrl.GetPairingHistory)
	}
	rg.POST("/devices/:device_id/unbind", ctrl.UnbindDevice)
}

// GeneratePairingCode 生成配对码
// POST /api/v1/devices/pairing/code
func (ctrl *DevicePairingController) GeneratePairingCode(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req models.PairingCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 生成6位随机配对码
	code, err := generateRandomCode(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成配对码失败"})
		return
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	record := models.PairingRecord{
		PairingCode: code,
		UserID:      userID,
		Status:      models.PairingStatusPending,
		ExpiresAt:   expiresAt,
	}

	if err := ctrl.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建配对记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"pairing_code":      code,
			"expires_at":          expiresAt.Format(time.RFC3339),
			"expires_in_seconds": 300,
		},
	})
}

// VerifyPairingCode 验证配对码并绑定设备
// POST /api/v1/devices/pairing/verify
func (ctrl *DevicePairingController) VerifyPairingCode(c *gin.Context) {
	var req models.PairingVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 查询配对码记录
	var record models.PairingRecord
	if err := ctrl.DB.Where("pairing_code = ? AND status = ?", req.PairingCode, models.PairingStatusPending).First(&record).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "配对码无效或已使用"})
		return
	}

	// 检查是否过期
	if time.Now().After(record.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "配对码已过期"})
		return
	}

	// 更新配对记录
	now := time.Now()
	record.DeviceID = req.DeviceID
	record.Status = models.PairingStatusPaired
	record.PairedAt = &now

	if err := ctrl.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新配对记录失败"})
		return
	}

	// 创建设备AI授权
	binding := models.DeviceOpenClawBinding{
		DeviceID:    req.DeviceID,
		UserID:      record.UserID,
		AuthStatus:  models.AuthStatusOK,
		AuthToken:   generateAuthToken(),
	}
	ctrl.DB.Create(&binding)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "配对成功",
		"data": gin.H{
			"device_id":   req.DeviceID,
			"paired_at":   now.Format(time.RFC3339),
		},
	})
}

// GetPairingHistory 查询配对历史
// GET /api/v1/devices/pairing/history
func (ctrl *DevicePairingController) GetPairingHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	deviceID := c.Query("device_id")

	query := ctrl.DB.Model(&models.PairingRecord{}).Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	var total int64
	query.Count(&total)

	var records []models.PairingRecord
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": records,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// UnbindDevice 解绑设备
// POST /api/v1/devices/:device_id/unbind
func (ctrl *DevicePairingController) UnbindDevice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	deviceID := c.Param("device_id")

	// 更新配对记录
	var record models.PairingRecord
	if err := ctrl.DB.Where("device_id = ? AND user_id = ? AND status = ?", deviceID, userID, models.PairingStatusPaired).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到配对记录"})
		return
	}

	now := time.Now()
	record.Status = models.PairingStatusUnbound
	record.UnboundAt = &now
	record.UnboundReason = "用户主动解绑"
	ctrl.DB.Save(&record)

	// 删除AI授权
	ctrl.DB.Where("device_id = ? AND user_id = ?", deviceID, userID).Delete(&models.DeviceOpenClawBinding{})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "解绑成功",
		"data": gin.H{
			"device_id": deviceID,
			"unbound_at": now.Format(time.RFC3339),
		},
	})
}

// ============ 辅助函数 ============

func generateRandomCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}

func generateAuthToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
