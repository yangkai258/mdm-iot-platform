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
