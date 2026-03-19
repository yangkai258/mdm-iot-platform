package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UnbindDevice 解绑设备
func UnbindDevice(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		snCode := ctx.Param("sn_code")

		var device models.Device
		result := db.Where("sn_code = ?", snCode).First(&device)
		
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":      4002,
				"message":   "设备不存在",
				"error_code": "ERR_DEVICE_002",
			})
			return
		}

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "服务器内部错误",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		// 解绑用户
		device.BindUserID = nil
		device.LifecycleStatus = 1 // 改为未绑定状态
		
		if err := db.Save(&device).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "解绑失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id": device.DeviceID,
				"sn_code":   device.SnCode,
				"status":    "unbound",
			},
		})
	}
}

// UpdateDeviceStatus 更新设备状态
func UpdateDeviceStatus(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceID := ctx.Param("device_id")

		var device models.Device
		result := db.Where("device_id = ?", deviceID).First(&device)
		
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":      4002,
				"message":   "设备不存在",
				"error_code": "ERR_DEVICE_002",
			})
			return
		}

		// 解析请求
		var req struct {
			Status int `json:"status"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "参数错误",
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		// 验证状态值 (1:待激活 2:服役中 3:维修中 4:已挂失 5:已报废)
		if req.Status < 1 || req.Status > 5 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "无效的状态值",
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		device.LifecycleStatus = req.Status
		if err := db.Save(&device).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "更新失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":         device.DeviceID,
				"lifecycle_status": device.LifecycleStatus,
			},
		})
	}
}
