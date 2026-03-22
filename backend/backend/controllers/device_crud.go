package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeleteDevice 删除设备
func DeleteDevice(db *gorm.DB) gin.HandlerFunc {
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

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "服务器内部错误",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		// 软删除
		if err := db.Delete(&device).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "删除失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id": deviceID,
			},
		})
	}
}

// UpdateDevice 更新设备
func UpdateDevice(db *gorm.DB) gin.HandlerFunc {
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
			MacAddress       string `json:"mac_address"`
			HardwareModel    string `json:"hardware_model"`
			FirmwareVersion  string `json:"firmware_version"`
			LifecycleStatus  int    `json:"lifecycle_status"`
		}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":      4005,
				"message":   "参数错误",
				"error_code": "ERR_VALIDATION",
			})
			return
		}

		// 更新字段
		if req.MacAddress != "" {
			device.MacAddress = req.MacAddress
		}
		if req.HardwareModel != "" {
			device.HardwareModel = req.HardwareModel
		}
		if req.FirmwareVersion != "" {
			device.FirmwareVersion = req.FirmwareVersion
		}
		if req.LifecycleStatus > 0 {
			device.LifecycleStatus = req.LifecycleStatus
		}

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
			"data":    device,
		})
	}
}
