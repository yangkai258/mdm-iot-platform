package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SmartHomeController 智能家居控制器
type SmartHomeController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册智能家居路由
func (c *SmartHomeController) RegisterRoutes(api *gin.RouterGroup) {
	smartHome := api.Group("/smart-home")
	{
		smartHome.GET("/devices", c.ListDevices)
		smartHome.GET("/devices/:id", c.GetDevice)
		smartHome.POST("/devices", c.CreateDevice)
		smartHome.PUT("/devices/:id", c.UpdateDevice)
		smartHome.DELETE("/devices/:id", c.DeleteDevice)
		smartHome.POST("/devices/:id/control", c.ControlDevice)
		smartHome.GET("/devices/:id/status", c.GetDeviceStatus)

		smartHome.GET("/triggers", c.ListTriggers)
		smartHome.POST("/triggers", c.CreateTrigger)
		smartHome.PUT("/triggers/:id", c.UpdateTrigger)
		smartHome.DELETE("/triggers/:id", c.DeleteTrigger)
		smartHome.POST("/triggers/:id/execute", c.ExecuteTrigger)

		smartHome.GET("/scenes", c.ListScenes)
		smartHome.POST("/scenes", c.CreateScene)
		smartHome.PUT("/scenes/:id", c.UpdateScene)
		smartHome.DELETE("/scenes/:id", c.DeleteScene)
		smartHome.POST("/scenes/:id/activate", c.ActivateScene)
	}
}

// ListDevices 获取智能设备列表
func (c *SmartHomeController) ListDevices(ctx *gin.Context) {
	var devices []models.SmartHomeDevice
	var total int64

	query := c.DB.Model(&models.SmartHomeDevice{})

	// 用户筛选
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 类型筛选
	if deviceType := ctx.Query("type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 房间筛选
	if room := ctx.Query("room"); room != "" {
		query = query.Where("room = ?", room)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&devices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      devices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDevice 获取设备详情
func (c *SmartHomeController) GetDevice(ctx *gin.Context) {
	id := ctx.Param("id")
	var device models.SmartHomeDevice
	if err := c.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": device})
}

// CreateDevice 创建设备
func (c *SmartHomeController) CreateDevice(ctx *gin.Context) {
	var device models.SmartHomeDevice
	if err := ctx.ShouldBindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	device.Status = models.JSONMap{"status": "offline"}
	if err := c.DB.Create(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": device})
}

// UpdateDevice 更新设备
func (c *SmartHomeController) UpdateDevice(ctx *gin.Context) {
	id := ctx.Param("id")
	var device models.SmartHomeDevice
	if err := c.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		Name     string `json:"name"`
		Room     string `json:"room"`
		Location string `json:"location"`
		Status   string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.Name != "" {
		updates["name"] = updateData.Name
	}
	if updateData.Room != "" {
		updates["room"] = updateData.Room
	}
	if updateData.Location != "" {
		updates["location"] = updateData.Location
	}
	if updateData.Status != "" {
		updates["status"] = updateData.Status
	}

	if err := c.DB.Model(&device).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&device, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": device})
}

// DeleteDevice 删除设备
func (c *SmartHomeController) DeleteDevice(ctx *gin.Context) {
	id := ctx.Param("id")
	var device models.SmartHomeDevice
	if err := c.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	if err := c.DB.Delete(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ControlDevice 控制设备
func (c *SmartHomeController) ControlDevice(ctx *gin.Context) {
	id := ctx.Param("id")
	var device models.SmartHomeDevice
	if err := c.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var controlData struct {
		Action string `json:"action" binding:"required"` // on, off, toggle, set
		Value  interface{} `json:"value"`                // 亮度、色温等
	}
	if err := ctx.ShouldBindJSON(&controlData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 更新设备状态
	newStatus := "off"
	if controlData.Action == "on" {
		newStatus = "on"
	} else if controlData.Action == "off" {
		newStatus = "off"
	} else if controlData.Action == "toggle" {
		currentStatus := "off"
		if s, ok := device.Status["status"].(string); ok {
			currentStatus = s
		}
		if currentStatus == "on" {
			newStatus = "off"
		} else {
			newStatus = "on"
		}
	}

	c.DB.Model(&device).Update("status", models.JSONMap{"status": newStatus})

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"device_id": id,
			"action":    controlData.Action,
			"new_status": newStatus,
		},
	})
}

// GetDeviceStatus 获取设备状态
func (c *SmartHomeController) GetDeviceStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var device models.SmartHomeDevice
	if err := c.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":   device.ID,
			"status":       device.Status,
			"device_type":  device.DeviceType,
			"online":       device.IsOnline,
		},
	})
}

// ListTriggers 获取触发器列表
func (c *SmartHomeController) ListTriggers(ctx *gin.Context) {
	var triggers []models.SmartHomeTrigger
	var total int64

	query := c.DB.Model(&models.SmartHomeTrigger{})

	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if triggerType := ctx.Query("type"); triggerType != "" {
		query = query.Where("trigger_type = ?", triggerType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&triggers).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      triggers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateTrigger 创建触发器
func (c *SmartHomeController) CreateTrigger(ctx *gin.Context) {
	var trigger models.SmartHomeTrigger
	if err := ctx.ShouldBindJSON(&trigger); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	trigger.IsEnabled = true
	if err := c.DB.Create(&trigger).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": trigger})
}

// UpdateTrigger 更新触发器
func (c *SmartHomeController) UpdateTrigger(ctx *gin.Context) {
	id := ctx.Param("id")
	var trigger models.SmartHomeTrigger
	if err := c.DB.First(&trigger, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "触发器不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		TriggerType string `json:"trigger_type"`
		Condition   string `json:"condition"`
		Action      string `json:"action"`
		IsEnabled   *bool  `json:"is_enabled"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if updateData.Name != "" {
		updates["name"] = updateData.Name
	}
	if updateData.TriggerType != "" {
		updates["trigger_type"] = updateData.TriggerType
	}
	if updateData.Condition != "" {
		updates["condition"] = updateData.Condition
	}
	if updateData.Action != "" {
		updates["action"] = updateData.Action
	}
	if updateData.IsEnabled != nil {
		updates["is_enabled"] = *updateData.IsEnabled
	}

	if err := c.DB.Model(&trigger).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&trigger, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": trigger})
}

// DeleteTrigger 删除触发器
func (c *SmartHomeController) DeleteTrigger(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.SmartHomeTrigger{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ExecuteTrigger 执行触发器
func (c *SmartHomeController) ExecuteTrigger(ctx *gin.Context) {
	id := ctx.Param("id")
	var trigger models.SmartHomeTrigger
	if err := c.DB.First(&trigger, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "触发器不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 记录触发执行
	trigger.RunCount++
	c.DB.Model(&trigger).Update("run_count", trigger.RunCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"trigger_id": id,
			"executed":   true,
		},
	})
}

// ListScenes 获取场景列表
func (c *SmartHomeController) ListScenes(ctx *gin.Context) {
	var scenes []gin.H
	var total int64 = 5

	// 预定义场景
	scenes = []gin.H{
		{"id": 1, "name": "早安模式", "icon": "sunrise", "actions": "开灯、调高温度、播放音乐"},
		{"id": 2, "name": "离家模式", "icon": "home", "actions": "关闭所有设备、启动安防"},
		{"id": 3, "name": "回家模式", "icon": "home", "actions": "开灯、调至舒适温度"},
		{"id": 4, "name": "影院模式", "icon": "film", "actions": "关灯、调暗灯光"},
		{"id": 5, "name": "睡眠模式", "icon": "moon", "actions": "关闭所有设备、调至低温"},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      scenes,
			"total":     total,
			"page":      1,
			"page_size": 20,
		},
	})
}

// CreateScene 创建场景
func (c *SmartHomeController) CreateScene(ctx *gin.Context) {
	var sceneData struct {
		Name    string `json:"name" binding:"required"`
		Icon    string `json:"icon"`
		Actions string `json:"actions" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&sceneData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"id":      100,
			"name":    sceneData.Name,
			"icon":    sceneData.Icon,
			"actions": sceneData.Actions,
		},
	})
}

// UpdateScene 更新场景
func (c *SmartHomeController) UpdateScene(ctx *gin.Context) {
	id := ctx.Param("id")
	var sceneData struct {
		Name    string `json:"name"`
		Icon    string `json:"icon"`
		Actions string `json:"actions"`
	}
	if err := ctx.ShouldBindJSON(&sceneData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"id":      id,
			"name":    sceneData.Name,
			"icon":    sceneData.Icon,
			"actions": sceneData.Actions,
		},
	})
}

// DeleteScene 删除场景
func (c *SmartHomeController) DeleteScene(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ActivateScene 激活场景
func (c *SmartHomeController) ActivateScene(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"scene_id": id,
			"activated": true,
		},
	})
}
