package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BatchController 批量操作控制器
type BatchController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterBatchRoutes 注册批量操作路由
func RegisterBatchRoutes(api *gin.RouterGroup, db *gorm.DB, redisClient *utils.RedisClient) {
	ctrl := &BatchController{DB: db, Redis: redisClient}
	api.POST("/devices/batch/bind", ctrl.BatchBind)
	api.POST("/devices/batch/unbind", ctrl.BatchUnbind)
	api.POST("/devices/batch/transfer", ctrl.BatchTransfer)
	api.POST("/devices/batch/status", ctrl.BatchUpdateStatus)
}

// ============ 批量绑定设备 ============

// BatchBindRequest 批量绑定请求
type BatchBindRequest struct {
	DeviceIDs  []string `json:"device_ids" binding:"required,min=1"`
	BindUserID string   `json:"bind_user_id" binding:"required"`
}

// BatchBind 批量绑定设备
func (c *BatchController) BatchBind(ctx *gin.Context) {
	var req BatchBindRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if len(req.DeviceIDs) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "单次批量操作最多支持100台设备",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	type Result struct {
		DeviceID string `json:"device_id"`
		Success  bool   `json:"success"`
		Message  string `json:"message"`
	}

	results := make([]Result, 0, len(req.DeviceIDs))
	successCount := 0
	failCount := 0

	for _, deviceID := range req.DeviceIDs {
		var device models.Device
		if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "设备不存在",
			})
			failCount++
			continue
		}

		// 检查状态：只有 1(待激活) 或 2(服役中) 可以绑定
		if device.LifecycleStatus != 1 && device.LifecycleStatus != 2 {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "设备状态不允许绑定",
			})
			failCount++
			continue
		}

		device.BindUserID = &req.BindUserID
		device.LifecycleStatus = 2 // 服役中
		device.UpdatedAt = time.Now()

		if err := c.DB.Save(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "绑定失败: " + err.Error(),
			})
			failCount++
			continue
		}

		results = append(results, Result{
			DeviceID: deviceID,
			Success:  true,
			Message:  "绑定成功",
		})
		successCount++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量绑定完成",
		"data": gin.H{
			"total":        len(req.DeviceIDs),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		},
	})
}

// ============ 批量解绑设备 ============

// BatchUnbindRequest 批量解绑请求
type BatchUnbindRequest struct {
	DeviceIDs []string `json:"device_ids" binding:"required,min=1"`
}

// BatchUnbind 批量解绑设备
func (c *BatchController) BatchUnbind(ctx *gin.Context) {
	var req BatchUnbindRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if len(req.DeviceIDs) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "单次批量操作最多支持100台设备",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	type Result struct {
		DeviceID string `json:"device_id"`
		Success  bool   `json:"success"`
		Message  string `json:"message"`
	}

	results := make([]Result, 0, len(req.DeviceIDs))
	successCount := 0
	failCount := 0

	for _, deviceID := range req.DeviceIDs {
		var device models.Device
		if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "设备不存在",
			})
			failCount++
			continue
		}

		// 只有 2(服役中) 状态的设备可以解绑
		if device.LifecycleStatus != 2 {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "只有服役中的设备才能解绑",
			})
			failCount++
			continue
		}

		// 解绑：清空绑定用户，状态改回待激活(1)
		device.BindUserID = nil
		device.LifecycleStatus = 1
		device.UpdatedAt = time.Now()

		if err := c.DB.Save(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "解绑失败: " + err.Error(),
			})
			failCount++
			continue
		}

		// 从 Redis 清除设备影子（模拟设备下线）
		c.Redis.DelDeviceShadow(deviceID)

		results = append(results, Result{
			DeviceID: deviceID,
			Success:  true,
			Message:  "解绑成功",
		})
		successCount++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量解绑完成",
		"data": gin.H{
			"total":        len(req.DeviceIDs),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		},
	})
}

// ============ 批量转移设备 ============

// BatchTransferRequest 批量转移请求
type BatchTransferRequest struct {
	DeviceIDs       []string `json:"device_ids" binding:"required,min=1"`
	SourceUserID    string   `json:"source_user_id"` // 可选：指定源用户（为空则不校验）
	TargetUserID    string   `json:"target_user_id" binding:"required"`
	TargetTenantID  string   `json:"target_tenant_id"` // 可选：跨租户转移
}

// BatchTransfer 批量转移设备
func (c *BatchController) BatchTransfer(ctx *gin.Context) {
	var req BatchTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if len(req.DeviceIDs) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "单次批量操作最多支持100台设备",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if req.SourceUserID == req.TargetUserID && req.SourceUserID != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "源用户和目标用户不能相同",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	type Result struct {
		DeviceID string `json:"device_id"`
		Success  bool   `json:"success"`
		Message  string `json:"message"`
	}

	results := make([]Result, 0, len(req.DeviceIDs))
	successCount := 0
	failCount := 0

	for _, deviceID := range req.DeviceIDs {
		var device models.Device
		if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "设备不存在",
			})
			failCount++
			continue
		}

		// 检查源用户（如果指定）
		if req.SourceUserID != "" && device.BindUserID != nil && *device.BindUserID != req.SourceUserID {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "设备不属于指定源用户",
			})
			failCount++
			continue
		}

		// 跨租户转移：清空绑定关系
		if req.TargetTenantID != "" {
			device.BindUserID = nil
		} else {
			device.BindUserID = &req.TargetUserID
		}
		device.UpdatedAt = time.Now()

		if err := c.DB.Save(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "转移失败: " + err.Error(),
			})
			failCount++
			continue
		}

		results = append(results, Result{
			DeviceID: deviceID,
			Success:  true,
			Message:  "转移成功",
		})
		successCount++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量转移完成",
		"data": gin.H{
			"total":        len(req.DeviceIDs),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		},
	})
}

// ============ 批量更新设备状态 ============

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	DeviceIDs []string `json:"device_ids" binding:"required,min=1"`
	Status    int      `json:"status" binding:"required,min=1,max=4"` // 1:待激活 2:服役中 3:维修 4:报废
}

// BatchUpdateStatus 批量更新设备状态
func (c *BatchController) BatchUpdateStatus(ctx *gin.Context) {
	var req BatchUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	if len(req.DeviceIDs) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "单次批量操作最多支持100台设备",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	type Result struct {
		DeviceID string `json:"device_id"`
		Success  bool   `json:"success"`
		Message  string `json:"message"`
	}

	results := make([]Result, 0, len(req.DeviceIDs))
	successCount := 0
	failCount := 0

	now := time.Now()
	for _, deviceID := range req.DeviceIDs {
		var device models.Device
		if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "设备不存在",
			})
			failCount++
			continue
		}

		device.LifecycleStatus = req.Status
		device.UpdatedAt = now

		if err := c.DB.Save(&device).Error; err != nil {
			results = append(results, Result{
				DeviceID: deviceID,
				Success:  false,
				Message:  "更新失败: " + err.Error(),
			})
			failCount++
			continue
		}

		results = append(results, Result{
			DeviceID: deviceID,
			Success:  true,
			Message:  "状态更新成功",
		})
		successCount++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批量状态更新完成",
		"data": gin.H{
			"total":        len(req.DeviceIDs),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		},
	})
}
