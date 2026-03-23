package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceShadowController 设备影子控制器
type DeviceShadowController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册设备影子路由
func (c *DeviceShadowController) RegisterRoutes(api *gin.RouterGroup) {
	shadow := api.Group("/device-shadow")
	{
		shadow.GET("/:device_id", c.GetShadow)
		shadow.PUT("/:device_id", c.UpdateShadow)
		shadow.POST("/:device_id/desired", c.UpdateDesired)
		shadow.POST("/:device_id/reported", c.UpdateReported)
	}
}

// GetShadow 获取设备影子
func (c *DeviceShadowController) GetShadow(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var shadow models.DeviceShadow
	if err := c.DB.Where("device_id = ?", deviceID).First(&shadow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果影子不存在，返回空数据
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"device_id":     deviceID,
					"is_online":     false,
					"desired_state": gin.H{},
					"reported_state": gin.H{},
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    shadow,
	})
}

// UpdateShadow 更新设备影子
func (c *DeviceShadowController) UpdateShadow(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var input struct {
		BatteryLevel  *int                    `json:"battery_level"`
		CurrentMode  *string                 `json:"current_mode"`
		IsOnline     *bool                    `json:"is_online"`
		DesiredState map[string]interface{}   `json:"desired_state"`
		ReportedState map[string]interface{}  `json:"reported_state"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var shadow models.DeviceShadow
	err := c.DB.Where("device_id = ?", deviceID).First(&shadow).Error
	if err == gorm.ErrRecordNotFound {
		// 创建新影子
		shadow = models.DeviceShadow{DeviceID: deviceID}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	// 更新字段
	if input.BatteryLevel != nil {
		shadow.BatteryLevel = *input.BatteryLevel
	}
	if input.CurrentMode != nil {
		shadow.CurrentMode = *input.CurrentMode
	}
	if input.IsOnline != nil {
		shadow.IsOnline = *input.IsOnline
	}

	if err := c.DB.Save(&shadow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": shadow, "message": "success"})
}

// UpdateDesired 更新期望状态
func (c *DeviceShadowController) UpdateDesired(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var input struct {
		NRDEnabled *bool   `json:"nrd_enabled"`
		NRDStart   *string `json:"nrd_start"`
		NRDEnd     *string `json:"nrd_end"`
		DNDEnabled *bool   `json:"dnd_enabled"`
		Config     map[string]interface{} `json:"config"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var shadow models.DeviceShadow
	err := c.DB.Where("device_id = ?", deviceID).First(&shadow).Error
	if err == gorm.ErrRecordNotFound {
		shadow = models.DeviceShadow{DeviceID: deviceID}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if input.NRDEnabled != nil {
		shadow.DesiredNRDEnabled = *input.NRDEnabled
	}
	if input.NRDStart != nil {
		shadow.DesiredNRDStart = *input.NRDStart
	}
	if input.NRDEnd != nil {
		shadow.DesiredNRDEnd = *input.NRDEnd
	}
	if input.DNDEnabled != nil {
		shadow.DesiredDNDEnabled = *input.DNDEnabled
	}

	if err := c.DB.Save(&shadow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": shadow, "message": "success"})
}

// UpdateReported 更新报告状态
func (c *DeviceShadowController) UpdateReported(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var input struct {
		BatteryLevel *int                    `json:"battery_level"`
		CurrentMode  *string                 `json:"current_mode"`
		IsOnline     *bool                    `json:"is_online"`
		LastIP       *string                  `json:"last_ip"`
		IsJailbroken *bool                    `json:"is_jailbroken"`
		Latitude     *float64                `json:"latitude"`
		Longitude    *float64                `json:"longitude"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var shadow models.DeviceShadow
	err := c.DB.Where("device_id = ?", deviceID).First(&shadow).Error
	if err == gorm.ErrRecordNotFound {
		shadow = models.DeviceShadow{DeviceID: deviceID}
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if input.BatteryLevel != nil {
		shadow.BatteryLevel = *input.BatteryLevel
	}
	if input.CurrentMode != nil {
		shadow.CurrentMode = *input.CurrentMode
	}
	if input.IsOnline != nil {
		shadow.IsOnline = *input.IsOnline
	}
	if input.LastIP != nil {
		shadow.LastIP = *input.LastIP
	}
	if input.IsJailbroken != nil {
		shadow.IsJailbroken = *input.IsJailbroken
	}
	if input.Latitude != nil {
		shadow.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		shadow.Longitude = *input.Longitude
	}

	if err := c.DB.Save(&shadow).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": shadow, "message": "success"})
}
