package controllers

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PairingController 设备配对管理控制器
type PairingController struct {
	DB *gorm.DB
}

// ============ 请求结构 ============

type PairingListRequest struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	DeviceID   string `form:"device_id"`
	UserID     string `form:"user_id"`
	Status     string `form:"status"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
}

type PairingApproveRequest struct {
	Notes string `json:"notes"`
}

type PairingRejectRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type PairingGenerateRequest struct {
	DeviceID string `json:"device_id"`
	ExpireMinutes int `json:"expire_minutes"`
}

type PairingVerifyRequest struct {
	PairingCode string `json:"pairing_code" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceInfo  struct {
		FirmwareVersion string `json:"firmware_version"`
		HardwareVersion string `json:"hardware_version"`
		MACAddress      string `json:"mac_address"`
	} `json:"device_info"`
}

// ============ 配对码管理 ============

// GenerateCode 生成分对码
func (c *PairingController) GenerateCode(ctx *gin.Context) {
	var req PairingGenerateRequest
	ctx.ShouldBindJSON(&req)

	expireMinutes := 5 // 默认5分钟
	if req.ExpireMinutes > 0 {
		expireMinutes = req.ExpireMinutes
	}

	// 生成6位数字配对码
	code := generatePairingCode()

	// 计算过期时间
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)

	// 获取当前用户
	userID := "system"
	userName := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		userID = uid.(string)
	}
	if uname, exists := ctx.Get("user_name"); exists {
		userName = uname.(string)
	}

	// 获取客户端IP
	ipAddress := ctx.ClientIP()

	pairing := models.DevicePairing{
		PairingCode: code,
		DeviceID:    req.DeviceID,
		UserID:      userID,
		UserName:    userName,
		Status:      "pending",
		ExpiresAt:   expiresAt,
		IPAddress:   ipAddress,
	}

	if err := c.DB.Create(&pairing).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成配对码失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"pairing_code":    code,
			"expires_at":      expiresAt.Format(time.RFC3339),
			"expires_in_seconds": expireMinutes * 60,
		},
	})
}

// Verify 设备配对验证
func (c *PairingController) Verify(ctx *gin.Context) {
	var req PairingVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 查找配对码
	var pairing models.DevicePairing
	if err := c.DB.Where("pairing_code = ? AND status = ?", req.PairingCode, "pending").First(&pairing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配对码无效或已过期"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 检查是否过期
	if time.Now().After(pairing.ExpiresAt) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "配对码已过期"})
		return
	}

	// 更新配对信息
	now := time.Now()
	updates := map[string]interface{}{
		"device_id":     req.DeviceID,
		"device_name":   req.DeviceID,
		"firmware_ver":  req.DeviceInfo.FirmwareVersion,
		"hardware_ver":  req.DeviceInfo.HardwareVersion,
		"mac_address":   req.DeviceInfo.MACAddress,
		"status":        "approved",
		"approve_at":    now,
		"approve_by":    "system",
		"paired_at":     now,
	}

	if err := c.DB.Model(&pairing).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "配对失败"})
		return
	}

	c.DB.First(&pairing, pairing.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "配对成功",
		"data": gin.H{
			"pairing_id": pairing.ID,
			"device_id":  pairing.DeviceID,
			"user_id":    pairing.UserID,
			"status":     pairing.Status,
		},
	})
}

// ============ 配对列表 ============

// List 获取配对列表
func (c *PairingController) List(ctx *gin.Context) {
	var req PairingListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var list []models.DevicePairing
	var total int64

	query := c.DB.Model(&models.DevicePairing{})

	if req.DeviceID != "" {
		query = query.Where("device_id = ?", req.DeviceID)
	}
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&list)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": list,
			"pagination": gin.H{"total": total, "page": req.Page, "page_size": req.PageSize},
		},
	})
}

// Get 获取配对详情
func (c *PairingController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var pairing models.DevicePairing
	if err := c.DB.First(&pairing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配对记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取AI授权信息
	var binding models.DeviceOpenClawBinding
	c.DB.Where("device_id = ?", pairing.DeviceID).First(&binding)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"pairing": pairing,
			"binding": binding,
		},
	})
}

// Approve 批准配对
func (c *PairingController) Approve(ctx *gin.Context) {
	id := ctx.Param("id")
	var pairing models.DevicePairing
	if err := c.DB.First(&pairing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配对记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if pairing.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前状态不允许批准"})
		return
	}

	now := time.Now()
	approvedBy := "system"
	if uid, exists := ctx.Get("user_id"); exists {
		approvedBy = uid.(string)
	}

	if err := c.DB.Model(&pairing).Updates(map[string]interface{}{
		"status":     "approved",
		"approve_at": now,
		"approve_by": approvedBy,
		"paired_at":  now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批准失败"})
		return
	}

	// 创建AI授权绑定
	binding := models.DeviceOpenClawBinding{
		DeviceID:   pairing.DeviceID,
		UserID:     pairing.UserID,
		AuthStatus: "authorized",
	}
	c.DB.Create(&binding)

	c.DB.First(&pairing, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "批准成功",
		"data":   pairing,
	})
}

// Reject 拒绝配对
func (c *PairingController) Reject(ctx *gin.Context) {
	id := ctx.Param("id")
	var pairing models.DevicePairing
	if err := c.DB.First(&pairing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配对记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if pairing.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前状态不允许拒绝"})
		return
	}

	var req PairingRejectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&pairing).Updates(map[string]interface{}{
		"status":       "rejected",
		"reject_at":    now,
		"reject_reason": req.Reason,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "拒绝失败"})
		return
	}

	c.DB.First(&pairing, id)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "拒绝成功",
		"data":   pairing,
	})
}

// Unbind 解绑设备
func (c *PairingController) Unbind(ctx *gin.Context) {
	id := ctx.Param("id")
	var pairing models.DevicePairing
	if err := c.DB.First(&pairing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配对记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	ctx.ShouldBindJSON(&req)

	now := time.Now()
	if err := c.DB.Model(&pairing).Updates(map[string]interface{}{
		"status":        "unbound",
		"unbound_at":    now,
		"unbound_reason": req.Reason,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解绑失败"})
		return
	}

	// 取消AI授权
	c.DB.Model(&models.DeviceOpenClawBinding{}).Where("device_id = ?", pairing.DeviceID).
		Update("auth_status", "cancelled")

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "解绑成功"})
}

// ============ 辅助函数 ============

func generatePairingCode() string {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}
