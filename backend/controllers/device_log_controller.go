package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceLogController 设备日志控制器
type DeviceLogController struct {
	DB *gorm.DB
}

func NewDeviceLogController(db *gorm.DB) *DeviceLogController {
	return &DeviceLogController{DB: db}
}

// RegisterRoutes 注册设备日志路由
func (ctrl *DeviceLogController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/device/logs", ctrl.List)
	rg.GET("/devices/:device_id/logs", ctrl.ListByDevice)
}

// List 获取设备日志列表
// GET /api/v1/device/logs
func (ctrl *DeviceLogController) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.DeviceLog{}).Where("tenant_id = ?", tenantID)

	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if logType := c.Query("log_type"); logType != "" {
		query = query.Where("log_type = ?", logType)
	}
	if action := c.Query("action"); action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var logs []models.DeviceLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ListByDevice 获取单设备日志列表
// GET /api/v1/devices/:device_id/logs
func (ctrl *DeviceLogController) ListByDevice(c *gin.Context) {
	deviceID := c.Param("device_id")
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.DeviceLog{}).Where("device_id = ? AND tenant_id = ?", deviceID, tenantID)

	if logType := c.Query("log_type"); logType != "" {
		query = query.Where("log_type = ?", logType)
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var total int64
	query.Count(&total)

	var logs []models.DeviceLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
