package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OTACompatibilityController OTA兼容性控制器
type OTACompatibilityController struct {
	DB *gorm.DB
}

// NewOTACompatibilityController 创建控制器
func NewOTACompatibilityController(db *gorm.DB) *OTACompatibilityController {
	return &OTACompatibilityController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *OTACompatibilityController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ota/compatibility/matrix", ctrl.ListMatrix)
	rg.GET("/ota/compatibility/matrix/:id", ctrl.GetMatrix)
	rg.POST("/ota/compatibility/matrix", ctrl.CreateMatrix)
	rg.PUT("/ota/compatibility/matrix/:id", ctrl.UpdateMatrix)
	rg.DELETE("/ota/compatibility/matrix/:id", ctrl.DeleteMatrix)
	rg.GET("/ota/compatibility/check", ctrl.CheckCompatibility)
	rg.GET("/ota/compatibility/device/:device_id", ctrl.GetDeviceCompatibility)
	rg.POST("/ota/compatibility/tests", ctrl.RecordTest)
	rg.GET("/ota/compatibility/tests", ctrl.ListTests)
}

// ListMatrix 获取兼容性矩阵列表
func (ctrl *OTACompatibilityController) ListMatrix(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	hardwareModel := c.Query("hardware_model")
	status := c.Query("status")
	fromFirmware := c.Query("from_firmware")
	toFirmware := c.Query("to_firmware")

	query := ctrl.DB.Model(&models.OTACompatibilityMatrix{})

	if hardwareModel != "" {
		query = query.Where("hardware_model = ?", hardwareModel)
	}
	if status != "" {
		query = query.Where("compatibility_status = ?", status)
	}
	if fromFirmware != "" {
		query = query.Where("from_firmware = ?", fromFirmware)
	}
	if toFirmware != "" {
		query = query.Where("to_firmware = ?", toFirmware)
	}

	var total int64
	var list []models.OTACompatibilityMatrix
	query.Count(&total)

	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMatrix 获取单个兼容性记录
func (ctrl *OTACompatibilityController) GetMatrix(c *gin.Context) {
	id := c.Param("id")

	var matrix models.OTACompatibilityMatrix
	if err := ctrl.DB.Where("id = ? OR matrix_id = ?", id, id).First(&matrix).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "兼容性记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": matrix})
}

// CreateMatrix 创建兼容性记录
func (ctrl *OTACompatibilityController) CreateMatrix(c *gin.Context) {
	var req struct {
		HardwareModel       string  `json:"hardware_model" binding:"required"`
		HardwareVersion     string  `json:"hardware_version"`
		FromFirmware       string  `json:"from_firmware" binding:"required"`
		ToFirmware         string  `json:"to_firmware" binding:"required"`
		CompatibilityStatus string  `json:"status"`
		CompatibilityScore float64 `json:"score"`
		MinBatteryLevel    int     `json:"min_battery_level"`
		MinStorageKB       int     `json:"min_storage_kb"`
		MinMemoryKB        int     `json:"min_memory_kb"`
		NetworkRequired     bool    `json:"network_required"`
		Constraints        string  `json:"constraints"`
		Warning            string  `json:"warning"`
		BreakingChanges    string  `json:"breaking_changes"`
		RollbackSupported  bool    `json:"rollback_supported"`
		RollbackMinVersion string  `json:"rollback_min_version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	status := "compatible"
	if req.CompatibilityStatus != "" {
		status = req.CompatibilityStatus
	}

	matrixID := fmt.Sprintf("CM-%s-%s-%s", req.HardwareModel, req.FromFirmware, req.ToFirmware)

	matrix := models.OTACompatibilityMatrix{
		MatrixID:              matrixID,
		HardwareModel:         req.HardwareModel,
		HardwareVersion:       req.HardwareVersion,
		FromFirmware:         req.FromFirmware,
		ToFirmware:           req.ToFirmware,
		CompatibilityStatus:  status,
		CompatibilityScore:   req.CompatibilityScore,
		MinBatteryLevel:      req.MinBatteryLevel,
		MinStorageKB:         req.MinStorageKB,
		MinMemoryKB:          req.MinMemoryKB,
		NetworkRequired:      req.NetworkRequired,
		Constraints:          req.Constraints,
		Warning:              req.Warning,
		BreakingChanges:       req.BreakingChanges,
		RollbackSupported:    req.RollbackSupported,
		RollbackMinVersion:   req.RollbackMinVersion,
		IsActive:             true,
		CreatedBy:            c.GetString("username"),
	}

	if err := ctrl.DB.Create(&matrix).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": matrix})
}

// UpdateMatrix 更新兼容性记录
func (ctrl *OTACompatibilityController) UpdateMatrix(c *gin.Context) {
	id := c.Param("id")

	var matrix models.OTACompatibilityMatrix
	if err := ctrl.DB.Where("id = ? OR matrix_id = ?", id, id).First(&matrix).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	var req struct {
		CompatibilityStatus string  `json:"status"`
		CompatibilityScore float64 `json:"score"`
		MinBatteryLevel    int    `json:"min_battery_level"`
		MinStorageKB       int    `json:"min_storage_kb"`
		MinMemoryKB        int    `json:"min_memory_kb"`
		Constraints        string `json:"constraints"`
		Warning            string `json:"warning"`
		IsActive           *bool  `json:"is_active"`
		IsVerified         *bool  `json:"is_verified"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}

	if req.CompatibilityStatus != "" {
		updates["compatibility_status"] = req.CompatibilityStatus
	}
	if req.CompatibilityScore > 0 {
		updates["compatibility_score"] = req.CompatibilityScore
	}
	if req.MinBatteryLevel > 0 {
		updates["min_battery_level"] = req.MinBatteryLevel
	}
	if req.MinStorageKB > 0 {
		updates["min_storage_kb"] = req.MinStorageKB
	}
	if req.MinMemoryKB > 0 {
		updates["min_memory_kb"] = req.MinMemoryKB
	}
	if req.Constraints != "" {
		updates["constraints"] = req.Constraints
	}
	if req.Warning != "" {
		updates["warning"] = req.Warning
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsVerified != nil {
		updates["is_verified"] = *req.IsVerified
		if *req.IsVerified {
			now := time.Now()
			updates["verified_at"] = now
			updates["verified_by"] = c.GetString("username")
		}
	}

	ctrl.DB.Model(&matrix).Updates(updates)
	ctrl.DB.Where("id = ?", matrix.ID).First(&matrix)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": matrix})
}

// DeleteMatrix 删除记录
func (ctrl *OTACompatibilityController) DeleteMatrix(c *gin.Context) {
	id := c.Param("id")

	var matrix models.OTACompatibilityMatrix
	if err := ctrl.DB.Where("id = ? OR matrix_id = ?", id, id).First(&matrix).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	ctrl.DB.Delete(&matrix)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// CheckCompatibility 检查兼容性
func (ctrl *OTACompatibilityController) CheckCompatibility(c *gin.Context) {
	hardwareModel := c.Query("hardware_model")
	fromFirmware := c.Query("from_firmware")
	toFirmware := c.Query("to_firmware")

	if hardwareModel == "" || fromFirmware == "" || toFirmware == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少必要参数"})
		return
	}

	var matrix models.OTACompatibilityMatrix
	err := ctrl.DB.Where("hardware_model = ? AND from_firmware = ? AND to_firmware = ? AND is_active = ?",
		hardwareModel, fromFirmware, toFirmware, true).First(&matrix).Error

	if err != nil {
		// 没有找到精确匹配，尝试查找条件更宽的记录
		var partial []models.OTACompatibilityMatrix
		ctrl.DB.Where("hardware_model = ? AND is_active = ?", hardwareModel, true).
			Where("from_firmware = ? AND to_firmware = ?", fromFirmware, toFirmware).
			Find(&partial)

		if len(partial) > 0 {
			matrix = partial[0]
			err = nil
		}
	}

	result := gin.H{
		"hardware_model": hardwareModel,
		"from_firmware": fromFirmware,
		"to_firmware":   toFirmware,
	}

	if err == nil && matrix.ID > 0 {
		result["compatible"] = matrix.CompatibilityStatus == "compatible" || matrix.CompatibilityStatus == "conditional"
		result["status"] = matrix.CompatibilityStatus
		result["score"] = matrix.CompatibilityScore
		result["constraints"] = gin.H{
			"min_battery_level": matrix.MinBatteryLevel,
			"min_storage_kb":   matrix.MinStorageKB,
			"min_memory_kb":    matrix.MinMemoryKB,
			"network_required": matrix.NetworkRequired,
		}
		result["warning"] = matrix.Warning
		result["breaking_changes"] = matrix.BreakingChanges
		result["rollback_supported"] = matrix.RollbackSupported
	} else {
		result["compatible"] = false
		result["status"] = "unknown"
		result["message"] = "未找到兼容性数据，请在管理后台配置"
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// GetDeviceCompatibility 获取设备可用的OTA升级列表
func (ctrl *OTACompatibilityController) GetDeviceCompatibility(c *gin.Context) {
	deviceID := c.Param("device_id")

	// 获取设备信息
	var device models.Device
	if err := ctrl.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}

	// 获取当前固件版本对应的所有可升级版本
	var upgrades []models.OTACompatibilityMatrix
	query := ctrl.DB.Model(&models.OTACompatibilityMatrix{}).
		Where("hardware_model = ? AND from_firmware = ? AND is_active = ?",
			device.HardwareModel, device.FirmwareVersion, true).
		Where("compatibility_status IN ('compatible','conditional')")

	query.Order("compatibility_score DESC, success_rate DESC").Find(&upgrades)

	result := gin.H{
		"device_id":        deviceID,
		"hardware_model":   device.HardwareModel,
		"current_version":  device.FirmwareVersion,
		"upgrades":         upgrades,
		"upgrade_count":    len(upgrades),
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// RecordTest 记录兼容性测试结果
func (ctrl *OTACompatibilityController) RecordTest(c *gin.Context) {
	var req struct {
		MatrixID     string `json:"matrix_id" binding:"required"`
		DeviceID    string `json:"device_id"`
		FromFirmware string `json:"from_firmware" binding:"required"`
		ToFirmware  string `json:"to_firmware" binding:"required"`
		TestResult  string `json:"result" binding:"required"`
		ErrorCode   string `json:"error_code"`
		ErrorMessage string `json:"error_message"`
		Duration    int    `json:"duration"`
		TestType    string `json:"test_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	test := models.OTACompatibilityTest{
		MatrixID:     req.MatrixID,
		DeviceID:     req.DeviceID,
		FromFirmware: req.FromFirmware,
		ToFirmware:  req.ToFirmware,
		TestResult:  req.TestResult,
		ErrorCode:   req.ErrorCode,
		ErrorMessage: req.ErrorMessage,
		Duration:    req.Duration,
		TestType:    req.TestType,
		TesterID:    c.GetString("user_id"),
	}

	if err := ctrl.DB.Create(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "记录失败"})
		return
	}

	// 更新矩阵的统计数据
	if req.MatrixID != "" {
		var matrix models.OTACompatibilityMatrix
		if err := ctrl.DB.Where("matrix_id = ?", req.MatrixID).First(&matrix).Error; err == nil {
			updates := map[string]interface{}{}
			if req.TestResult == "passed" {
				updates["success_count"] = matrix.SuccessCount + 1
			} else {
				updates["failure_count"] = matrix.FailureCount + 1
			}
			total := matrix.SuccessCount + matrix.FailureCount + 1
			updates["success_rate"] = float64(matrix.SuccessCount) / float64(total) * 100
			ctrl.DB.Model(&matrix).Updates(updates)
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "记录成功", "data": test})
}

// ListTests 获取测试记录列表
func (ctrl *OTACompatibilityController) ListTests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	matrixID := c.Query("matrix_id")
	result := c.Query("result")

	query := ctrl.DB.Model(&models.OTACompatibilityTest{})

	if matrixID != "" {
		query = query.Where("matrix_id = ?", matrixID)
	}
	if result != "" {
		query = query.Where("test_result = ?", result)
	}

	var total int64
	var list []models.OTACompatibilityTest
	query.Count(&total)

	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
