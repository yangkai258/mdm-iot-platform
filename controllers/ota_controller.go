package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OTAController OTA 固件管理控制器
type OTAController struct {
	DB *gorm.DB
}

// OTAWorkerRef 保存全局 OTA Worker 引用
var OTAWorkerRef *services.OTAWorker

// SetOTAWorkerRef 设置 OTA Worker 引用（供 main.go 调用）
func SetOTAWorkerRef(w *services.OTAWorker) {
	OTAWorkerRef = w
}

// CreatePackageRequest 创建固件包请求
type CreatePackageRequest struct {
	Name           string `json:"name" binding:"required"`
	Version        string `json:"version" binding:"required"`
	HardwareModel  string `json:"hardware_model" binding:"required"`
	FileSize       int64  `json:"file_size"`
	FileURL        string `json:"file_url" binding:"required"`
	FileMD5        string `json:"file_md5"`
	IsMandatory    bool   `json:"is_mandatory"`
	AllowDowngrade bool   `json:"allow_downgrade"`
	ReleaseNotes   string `json:"release_notes"`
}

// CreatePackage 创建固件包
func (c *OTAController) CreatePackage(ctx *gin.Context) {
	var req CreatePackageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	pkg := models.OTAPackage{
		Name:           req.Name,
		Version:        req.Version,
		HardwareModel:  req.HardwareModel,
		FileSize:       req.FileSize,
		FileURL:        req.FileURL,
		FileMD5:        req.FileMD5,
		IsMandatory:    req.IsMandatory,
		AllowDowngrade: req.AllowDowngrade,
		ReleaseNotes:   req.ReleaseNotes,
		IsActive:       true,
		CreatedBy:      "system",
	}

	if err := c.DB.Create(&pkg).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建固件包失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    pkg,
	})
}

// ListPackages 固件包列表
func (c *OTAController) ListPackages(ctx *gin.Context) {
	hardwareModel := ctx.Query("hardware_model")
	isActive := ctx.Query("is_active")

	query := c.DB.Model(&models.OTAPackage{})
	if hardwareModel != "" {
		query = query.Where("hardware_model = ?", hardwareModel)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var packages []models.OTAPackage
	if err := query.Order("created_at DESC").Find(&packages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": packages,
		},
	})
}

// CreateDeploymentRequest 创建发布任务请求
type CreateDeploymentRequest struct {
	Name           string     `json:"name" binding:"required"`
	PackageID      uint       `json:"package_id" binding:"required"`
	HardwareModel  string     `json:"hardware_model" binding:"required"`
	StrategyType   string     `json:"strategy_type" binding:"required"` // full / percentage / whitelist
	StrategyConfig string     `json:"strategy_config"`                  // 百分比数字或白名单JSON数组
	ScheduledAt    *time.Time `json:"scheduled_at"`
}

// CreateDeployment 创建发布任务
func (c *OTAController) CreateDeployment(ctx *gin.Context) {
	var req CreateDeploymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查固件包是否存在
	var pkg models.OTAPackage
	if err := c.DB.First(&pkg, req.PackageID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "固件包不存在",
			"error_code": "ERR_OTA_001",
		})
		return
	}

	// 验证硬件型号匹配
	if pkg.HardwareModel != req.HardwareModel {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "固件包硬件型号与目标不匹配",
			"error_code": "ERR_OTA_002",
		})
		return
	}

	deployment := models.OTADeployment{
		Name:                    req.Name,
		PackageID:               req.PackageID,
		HardwareModel:           req.HardwareModel,
		StrategyType:            req.StrategyType,
		StrategyConfig:          req.StrategyConfig,
		Status:                  "pending",
		PauseOnFailureThreshold: 20.0, // 默认 20% 失败率阈值
		ScheduledAt:             req.ScheduledAt,
		CreatedBy:               "system",
	}

	if err := c.DB.Create(&deployment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建发布任务失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    deployment,
	})
}

// ListDeployments 发布任务列表
func (c *OTAController) ListDeployments(ctx *gin.Context) {
	status := ctx.Query("status")
	hardwareModel := ctx.Query("hardware_model")

	query := c.DB.Model(&models.OTADeployment{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if hardwareModel != "" {
		query = query.Where("hardware_model = ?", hardwareModel)
	}

	var deployments []models.OTADeployment
	if err := query.Order("created_at DESC").Find(&deployments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": deployments,
		},
	})
}

// GetDeployment 获取单个部署详情
func (c *OTAController) GetDeployment(ctx *gin.Context) {
	id := ctx.Param("id")

	var deployment models.OTADeployment
	if err := c.DB.First(&deployment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "部署任务不存在",
			"error_code": "ERR_OTA_003",
		})
		return
	}

	// 关联查询固件包信息
	var pkg models.OTAPackage
	c.DB.First(&pkg, deployment.PackageID)

	// 查询进度统计
	var total, success, failed, pending int64
	c.DB.Model(&models.OTAProgress{}).Where("deployment_id = ?", id).Count(&total)
	c.DB.Model(&models.OTAProgress{}).Where("deployment_id = ? AND ota_status = ?", id, "success").Count(&success)
	c.DB.Model(&models.OTAProgress{}).Where("deployment_id = ? AND ota_status = ?", id, "failed").Count(&failed)
	c.DB.Model(&models.OTAProgress{}).Where("deployment_id = ? AND ota_status NOT IN ?", id, []string{"success", "failed"}).Count(&pending)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"deployment": deployment,
			"package":    pkg,
			"stats": gin.H{
				"total":   total,
				"success": success,
				"failed":  failed,
				"pending": pending,
			},
		},
	})
}

// DeviceOTAReportRequest 设备上报 OTA 升级进度请求
type DeviceOTAReportRequest struct {
	DeploymentID uint   `json:"deployment_id" binding:"required"`
	Status       string `json:"status" binding:"required"` // pending/downloading/verifying/flashing/success/failed
	Progress     int    `json:"progress"`                  // 0-100
	Message      string `json:"message"`
}

// DeviceOTAReport 设备回调：上报 OTA 升级进度
func (c *OTAController) DeviceOTAReport(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var req DeviceOTAReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 验证设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 查找 ota_progress 记录
	var progress models.OTAProgress
	err := c.DB.Where("device_id = ? AND deployment_id = ?", deviceID, req.DeploymentID).First(&progress).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "升级记录不存在",
			"error_code": "ERR_OTA_004",
		})
		return
	}

	// 更新进度
	updates := map[string]interface{}{
		"ota_status":  req.Status,
		"progress":    req.Progress,
		"ota_message": req.Message,
	}
	now := time.Now()
	if req.Status == "success" || req.Status == "failed" {
		updates["completed_at"] = &now
	}
	if req.Status == "pending" && progress.StartedAt == nil {
		updates["started_at"] = &now
	}

	if err := c.DB.Model(&progress).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "更新进度失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 通过全局 Worker 触发自动暂停检查
	if OTAWorkerRef != nil {
		var deployment models.OTADeployment
		if err := c.DB.First(&deployment, req.DeploymentID).Error; err == nil {
			OTAWorkerRef.CheckAndAutoPause(&deployment)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// CheckOTA 检查设备 OTA 状态
func (c *OTAController) CheckOTA(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 获取设备信息
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 查询最新可用固件
	var latestPkg models.OTAPackage
	result := c.DB.Where("hardware_model = ? AND is_active = ?", device.HardwareModel, true).
		Order("created_at DESC").First(&latestPkg)

	if result.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"has_update": false,
				"message":    "当前已是最新版本",
			},
		})
		return
	}

	// 比较版本
	hasUpdate := latestPkg.Version != device.FirmwareVersion

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"has_update":      hasUpdate,
			"current_version":  device.FirmwareVersion,
			"latest_version":   latestPkg.Version,
			"package":          latestPkg,
		},
	})
}

// PauseDeployment 暂停部署任务
func (c *OTAController) PauseDeployment(ctx *gin.Context) {
	id := ctx.Param("id")

	var deployment models.OTADeployment
	if err := c.DB.First(&deployment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "部署任务不存在",
			"error_code": "ERR_OTA_003",
		})
		return
	}

	if deployment.Status != "pending" && deployment.Status != "running" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "当前状态不允许暂停",
			"error_code": "ERR_OTA_005",
		})
		return
	}

	c.DB.Model(&deployment).Updates(map[string]interface{}{
		"status":       "paused",
		"pause_reason": "手动暂停",
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ResumeDeployment 恢复部署任务
func (c *OTAController) ResumeDeployment(ctx *gin.Context) {
	id := ctx.Param("id")

	var deployment models.OTADeployment
	if err := c.DB.First(&deployment, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":       4002,
			"message":    "部署任务不存在",
			"error_code": "ERR_OTA_003",
		})
		return
	}

	if deployment.Status != "paused" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "只有暂停状态可以恢复",
			"error_code": "ERR_OTA_006",
		})
		return
	}

	c.DB.Model(&deployment).Updates(map[string]interface{}{
		"status":      "running",
		"auto_paused": false,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// CancelDeployment 取消部署任务
func (c *OTAController) CancelDeployment(ctx *gin.Context) {
	id := ctx.Param("id")
	now := time.Now()

	result := c.DB.Model(&models.OTADeployment{}).
		Where("id = ? AND status IN ?", id, []string{"pending", "running", "paused"}).
		Updates(map[string]interface{}{
			"status":       "cancelled",
			"cancelled_by": "system",
			"cancelled_at": &now,
		})

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "无法取消该部署任务",
			"error_code": "ERR_OTA_007",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// GetDeploymentProgress 获取部署进度详情
func (c *OTAController) GetDeploymentProgress(ctx *gin.Context) {
	id := ctx.Param("id")

	var progressList []models.OTAProgress
	if err := c.DB.Where("deployment_id = ?", id).Order("created_at DESC").Find(&progressList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": progressList,
		},
	})
}

// PartialUpgradeRequest 分片升级请求
type PartialUpgradeRequest struct {
	PackageID   uint   `json:"package_id" binding:"required"`
	TotalShards int    `json:"total_shards" binding:"required"`
}

// StartPartialUpgrade 发起分片升级
// POST /api/v1/device/:id/ota/partial
func (c *OTAController) StartPartialUpgrade(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "设备ID不能为空",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var req PartialUpgradeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查固件包是否存在
	var pkg models.OTAPackage
	if err := c.DB.First(&pkg, req.PackageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4002,
				"message":    "固件包不存在",
				"error_code": "ERR_OTA_001",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询固件包失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 检查设备是否存在
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4002,
				"message":    "设备不存在",
				"error_code": "ERR_DEVICE_001",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询设备失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 检查是否已有进行中的分片升级
	var existing models.OTAPartialUpgrade
	if err := c.DB.Where("device_id = ? AND shard_status NOT IN ?", deviceID, []string{"done", "failed"}).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "设备存在进行中的升级任务",
			"error_code": "ERR_OTA_007",
		})
		return
	}

	// 计算分片大小
	chunkSize := pkg.FileSize / int64(req.TotalShards)
	if chunkSize <= 0 {
		chunkSize = 1024 * 1024 // 默认 1MB per chunk
	}

	now := time.Now()
	upgrade := models.OTAPartialUpgrade{
		DeviceID:        deviceID,
		PackageID:       req.PackageID,
		FromVersion:     device.FirmwareVersion,
		ToVersion:       pkg.Version,
		TotalShards:     req.TotalShards,
		ShardIndex:      0,
		TotalChunks:     int(pkg.FileSize / chunkSize),
		ByteOffset:      0,
		TotalBytes:      pkg.FileSize,
		ShardStatus:     "pending",
		ChunkStatus:     "idle",
		RetryCount:      0,
		StartedAt:       &now,
	}

	if err := c.DB.Create(&upgrade).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "创建升级记录失败: " + err.Error(),
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// TODO: 通过 MQTT 下发分片升级指令到设备

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "分片升级任务已创建",
		"data": gin.H{
			"upgrade_id":   upgrade.ID,
			"device_id":    deviceID,
			"total_shards": req.TotalShards,
			"total_bytes":  pkg.FileSize,
			"version":      pkg.Version,
		},
	})
}

// GetPartialUpgradeStatus 获取分片升级状态
// GET /api/v1/device/:id/ota/status
func (c *OTAController) GetPartialUpgradeStatus(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "设备ID不能为空",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var upgrade models.OTAPartialUpgrade
	if err := c.DB.Where("device_id = ?", deviceID).Order("created_at DESC").First(&upgrade).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4002,
				"message":    "未找到升级记录",
				"error_code": "ERR_OTA_004",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询升级状态失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 获取固件包信息
	var pkg models.OTAPackage
	c.DB.First(&pkg, upgrade.PackageID)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"upgrade": gin.H{
				"id":                  upgrade.ID,
				"device_id":           upgrade.DeviceID,
				"from_version":        upgrade.FromVersion,
				"to_version":          upgrade.ToVersion,
				"total_shards":        upgrade.TotalShards,
				"shard_index":         upgrade.ShardIndex,
				"total_chunks":        upgrade.TotalChunks,
				"chunk_index":         upgrade.ChunkIndex,
				"byte_offset":         upgrade.ByteOffset,
				"transferred_bytes":   upgrade.TransferredBytes,
				"total_bytes":         upgrade.TotalBytes,
				"progress":            upgrade.Progress,
				"shard_status":        upgrade.ShardStatus,
				"chunk_status":        upgrade.ChunkStatus,
				"retry_count":         upgrade.RetryCount,
				"error_message":       upgrade.ErrorMessage,
				"started_at":          upgrade.StartedAt,
				"completed_at":        upgrade.CompletedAt,
			},
			"package": gin.H{
				"id":       pkg.ID,
				"name":     pkg.Name,
				"version":  pkg.Version,
				"file_url": pkg.FileURL,
				"file_md5": pkg.FileMD5,
			},
		},
	})
}

// UpdatePartialUpgradeProgress 更新分片升级进度（设备回调）
// POST /api/v1/device/:id/ota/partial/progress
func (c *OTAController) UpdatePartialUpgradeProgress(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	type ProgressUpdate struct {
		ShardIndex       int    `json:"shard_index"`
		ChunkIndex       int    `json:"chunk_index"`
		ByteOffset       int64  `json:"byte_offset"`
		TransferredBytes int64  `json:"transferred_bytes"`
		Progress         int    `json:"progress"`
		ShardStatus      string `json:"shard_status"`
		ChunkStatus      string `json:"chunk_status"`
		ErrorMessage     string `json:"error_message"`
	}

	var req ProgressUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":       4005,
			"message":    "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var upgrade models.OTAPartialUpgrade
	if err := c.DB.Where("device_id = ?", deviceID).Order("created_at DESC").First(&upgrade).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":       4002,
				"message":    "未找到升级记录",
				"error_code": "ERR_OTA_004",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":       5001,
			"message":    "查询升级状态失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	updates := map[string]interface{}{
		"shard_index":        req.ShardIndex,
		"chunk_index":        req.ChunkIndex,
		"byte_offset":        req.ByteOffset,
		"transferred_bytes":  req.TransferredBytes,
		"progress":           req.Progress,
		"shard_status":       req.ShardStatus,
		"chunk_status":       req.ChunkStatus,
		"updated_at":         time.Now(),
	}
	if req.ErrorMessage != "" {
		updates["error_message"] = req.ErrorMessage
	}
	if req.ShardStatus == "done" {
		now := time.Now()
		updates["completed_at"] = &now
	}

	c.DB.Model(&upgrade).Updates(updates)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "进度已更新",
	})
}
