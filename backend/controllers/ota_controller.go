package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OTAController OTA 固件管理控制器
type OTAController struct {
	DB *gorm.DB
}

// CreatePackageRequest 创建固件包请求
type CreatePackageRequest struct {
	VersionCode   string `json:"version_code" binding:"required"`
	HardwareModel string `json:"hardware_model" binding:"required"`
	BinURL        string `json:"bin_url" binding:"required"`
	Md5Hash       string `json:"md5_hash" binding:"required"`
	IsMandatory   bool   `json:"is_mandatory"`
}

// CreatePackage 创建固件包
func (c *OTAController) CreatePackage(ctx *gin.Context) {
	var req CreatePackageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	pkg := models.OTAPackage{
		VersionCode:    req.VersionCode,
		HardwareModel:  req.HardwareModel,
		BinURL:        req.BinURL,
		Md5Hash:       req.Md5Hash,
		IsMandatory:    req.IsMandatory,
		ReleaseStatus: 0, // 默认测试版
	}

	if err := c.DB.Create(&pkg).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "创建固件包失败",
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

	query := c.DB.Model(&models.OTAPackage{})
	if hardwareModel != "" {
		query = query.Where("hardware_model = ?", hardwareModel)
	}

	var packages []models.OTAPackage
	if err := query.Order("created_at DESC").Find(&packages).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "查询失败",
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
	PackageID       uint   `json:"package_id" binding:"required"`
	TargetHardware  string `json:"target_hardware" binding:"required"`
	RolloutStrategy string `json:"rollout_strategy" binding:"required"` // full, percentage, whitelist
	Percentage      int    `json:"percentage"`
}

// CreateDeployment 创建发布任务
func (c *OTAController) CreateDeployment(ctx *gin.Context) {
	var req CreateDeploymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 检查固件包是否存在
	var pkg models.OTAPackage
	if err := c.DB.First(&pkg, req.PackageID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "固件包不存在",
			"error_code": "ERR_OTA_001",
		})
		return
	}

	deployment := models.OTADeployment{
		PackageID:       req.PackageID,
		TargetHardware:  req.TargetHardware,
		RolloutStrategy: req.RolloutStrategy,
		Percentage:      req.Percentage,
		Status:          "pending",
	}

	if err := c.DB.Create(&deployment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "创建发布任务失败",
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

// CheckOTA 检查设备 OTA 状态
func (c *OTAController) CheckOTA(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	// 获取设备信息
	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 查询最新可用固件
	var latestPkg models.OTAPackage
	result := c.DB.Where("hardware_model = ? AND release_status > 0", device.HardwareModel).
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
	hasUpdate := latestPkg.VersionCode != device.FirmwareVersion

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"has_update":     hasUpdate,
			"current_version": device.FirmwareVersion,
			"latest_version": latestPkg.VersionCode,
			"package":        latestPkg,
		},
	})
}
