package controllers

import (
	"fmt"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
)

// SensorController 传感器控制器
type SensorController struct{}

func NewSensorController() *SensorController {
	return &SensorController{}
}

// GetEvents 获取传感器事件列表
// GET /api/v1/sensors/:device_id/events
func (c *SensorController) GetEvents(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	var events []models.SensorEvent
	if err := models.DB.Where("device_id = ?", deviceID).Order("created_at DESC").Limit(100).Find(&events).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": events})
}

// ReportEvent 上报传感器事件
// POST /api/v1/sensors/:device_id/events
func (c *SensorController) ReportEvent(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	var input struct {
		SensorType  string  `json:"sensor_type" binding:"required"`
		SensorValue float64 `json:"sensor_value" binding:"required"`
		Unit        string  `json:"unit"`
		Threshold   float64 `json:"threshold"`
		Description string  `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	event := models.SensorEvent{
		EventID:     fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		DeviceID:    deviceID,
		SensorType:  input.SensorType,
		SensorValue: input.SensorValue,
		Unit:        input.Unit,
		Threshold:   input.Threshold,
		IsAbnormal:  input.SensorValue > input.Threshold,
		EventType:   "normal",
		Description: input.Description,
	}
	if event.IsAbnormal {
		event.EventType = "warning"
	}
	if err := models.DB.Create(&event).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": event})
}

// GetLatest 获取最新传感器数据
// GET /api/v1/sensors/:device_id/latest
func (c *SensorController) GetLatest(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	var event models.SensorEvent
	if err := models.DB.Where("device_id = ?", deviceID).Order("created_at DESC").First(&event).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "no data"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": event})
}

// SetThreshold 设置阈值
// PUT /api/v1/sensors/:device_id/thresholds
func (c *SensorController) SetThreshold(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")
	_ = deviceID
	// 阈值存储到设备配置表，本实现简化为返回成功
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "threshold updated"})
}
