package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MapIntegrationController 地图服务集成控制器
type MapIntegrationController struct {
	DB *gorm.DB
}

// NewMapIntegrationController 创建控制器
func NewMapIntegrationController(db *gorm.DB) *MapIntegrationController {
	return &MapIntegrationController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *MapIntegrationController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/map/config", ctrl.GetConfig)
	rg.POST("/map/config", ctrl.CreateConfig)
	rg.PUT("/map/config/:provider", ctrl.UpdateConfig)
	rg.GET("/map/location/:pet_id", ctrl.GetPetLocation)
	rg.POST("/map/location", ctrl.RecordLocation)
	rg.GET("/map/location/history/:pet_id", ctrl.GetLocationHistory)
	rg.GET("/map/geocode", ctrl.Geocode)
	rg.GET("/map/reverse-geocode", ctrl.ReverseGeocode)
	rg.GET("/map/route", ctrl.GetRoute)
	rg.GET("/map/logs", ctrl.GetServiceLogs)
}

// GetConfig 获取地图配置
func (ctrl *MapIntegrationController) GetConfig(c *gin.Context) {
	var configs []models.MapIntegrationConfig
	ctrl.DB.Find(&configs)

	// 隐藏敏感信息
	for i := range configs {
		if len(configs[i].APIKey) > 8 {
			configs[i].APIKey = configs[i].APIKey[:4] + "****" + configs[i].APIKey[len(configs[i].APIKey)-4:]
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": configs})
}

// CreateConfig 创建地图配置
func (ctrl *MapIntegrationController) CreateConfig(c *gin.Context) {
	var req struct {
		Provider   string `json:"provider" binding:"required"`
		APIKey    string `json:"api_key" binding:"required"`
		APISecret string `json:"api_secret"`
		Services  string `json:"services"`
		QuotaLimit int   `json:"quota_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	config := models.MapIntegrationConfig{
		Provider:    req.Provider,
		APIKey:     req.APIKey,
		APISecret:  req.APISecret,
		Services:   req.Services,
		QuotaLimit: req.QuotaLimit,
		IsActive:   false,
		Status:     "inactive",
		CreatedBy:  c.GetString("username"),
	}

	if err := ctrl.DB.Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	// 隐藏敏感信息
	config.APIKey = req.Provider + "****"
	config.APISecret = ""

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": config})
}

// UpdateConfig 更新配置
func (ctrl *MapIntegrationController) UpdateConfig(c *gin.Context) {
	provider := c.Param("provider")

	var config models.MapIntegrationConfig
	if err := ctrl.DB.Where("provider = ?", provider).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配置不存在"})
		return
	}

	var req struct {
		APIKey     string `json:"api_key"`
		APISecret  string `json:"api_secret"`
		IsActive   *bool  `json:"is_active"`
		QuotaLimit int    `json:"quota_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}

	if req.APIKey != "" {
		updates["api_key"] = req.APIKey
	}
	if req.APISecret != "" {
		updates["api_secret"] = req.APISecret
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
		if *req.IsActive {
			updates["status"] = "active"
		} else {
			updates["status"] = "inactive"
		}
	}
	if req.QuotaLimit > 0 {
		updates["quota_limit"] = req.QuotaLimit
	}

	ctrl.DB.Model(&config).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// GetPetLocation 获取宠物最新位置
func (ctrl *MapIntegrationController) GetPetLocation(c *gin.Context) {
	petID := c.Param("pet_id")

	var location models.PetLocation
	if err := ctrl.DB.Where("pet_id = ?", petID).Order("created_at DESC").First(&location).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "暂无位置数据"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": location})
}

// RecordLocation 记录位置
func (ctrl *MapIntegrationController) RecordLocation(c *gin.Context) {
	var req struct {
		PetID      string  `json:"pet_id" binding:"required"`
		DeviceID   string  `json:"device_id"`
		Latitude   float64 `json:"latitude" binding:"required"`
		Longitude  float64 `json:"longitude" binding:"required"`
		Altitude   float64 `json:"altitude"`
		Accuracy   float64 `json:"accuracy"`
		BatteryLevel int   `json:"battery_level"`
		LocationType string `json:"location_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	location := models.PetLocation{
		PetID:        req.PetID,
		DeviceID:    req.DeviceID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Altitude:   req.Altitude,
		Accuracy:   req.Accuracy,
		BatteryLevel: req.BatteryLevel,
		LocationType: req.LocationType,
	}

	if err := ctrl.DB.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "位置记录成功", "data": location})
}

// GetLocationHistory 获取位置历史
func (ctrl *MapIntegrationController) GetLocationHistory(c *gin.Context) {
	petID := c.Param("pet_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := ctrl.DB.Model(&models.PetLocation{}).Where("pet_id = ?", petID)

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	var total int64
	var list []models.PetLocation
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

// Geocode 地址转坐标
func (ctrl *MapIntegrationController) Geocode(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "地址不能为空"})
		return
	}

	// 获取active配置
	var config models.MapIntegrationConfig
	if err := ctrl.DB.Where("is_active = ?", true).First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先配置地图服务"})
		return
	}

	// 记录调用日志
	log := models.MapServiceLog{
		Provider:    config.Provider,
		ServiceType: "geocoding",
		Endpoint:   "/geocode",
		RequestData: address,
	}
	ctrl.DB.Create(&log)

	// TODO: 调用实际地图API
	// 这里返回mock数据，实际使用时替换为真实API调用
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"provider":  config.Provider,
			"address":   address,
			"latitude":  31.2304,
			"longitude": 121.4737,
			"message":   "Geocoding API待集成，使用模拟数据",
		},
	})
}

// ReverseGeocode 坐标转地址
func (ctrl *MapIntegrationController) ReverseGeocode(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("latitude"), 64)
	lng, _ := strconv.ParseFloat(c.Query("longitude"), 64)

	if lat == 0 || lng == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "坐标不能为空"})
		return
	}

	var config models.MapIntegrationConfig
	if err := ctrl.DB.Where("is_active = ?", true).First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先配置地图服务"})
		return
	}

	log := models.MapServiceLog{
		Provider:    config.Provider,
		ServiceType: "reverse_geocoding",
		Endpoint:   "/reverse-geocode",
		RequestData: strconv.FormatFloat(lat, 'f', 6, 64) + "," + strconv.FormatFloat(lng, 'f', 6, 64),
	}
	ctrl.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"provider":  config.Provider,
			"latitude":  lat,
			"longitude": lng,
			"address":   "上海市黄浦区人民大道",
			"city":     "上海市",
			"district": "黄浦区",
			"message":   "Reverse Geocoding API待集成，使用模拟数据",
		},
	})
}

// GetRoute 路径规划
func (ctrl *MapIntegrationController) GetRoute(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "起点和终点不能为空"})
		return
	}

	var config models.MapIntegrationConfig
	if err := ctrl.DB.Where("is_active = ?", true).First(&config).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先配置地图服务"})
		return
	}

	log := models.MapServiceLog{
		Provider:    config.Provider,
		ServiceType: "routing",
		Endpoint:   "/route",
		RequestData: origin + " -> " + destination,
	}
	ctrl.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"provider":     config.Provider,
			"origin":      origin,
			"destination": destination,
			"distance":     5000,
			"duration":     1200,
			"message":      "Routing API待集成，使用模拟数据",
		},
	})
}

// GetServiceLogs 获取服务调用日志
func (ctrl *MapIntegrationController) GetServiceLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	provider := c.Query("provider")

	query := ctrl.DB.Model(&models.MapServiceLog{})

	if provider != "" {
		query = query.Where("provider = ?", provider)
	}

	var total int64
	var list []models.MapServiceLog
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
