package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GeofenceController 设备地理围栏控制器
type GeofenceController struct {
	DB *gorm.DB
}

// ============ 请求结构 ============

type GeofenceListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	DeviceID string `form:"device_id"`
	Enabled  *bool  `form:"enabled"`
}

type GeofenceCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	DeviceID    string  `json:"device_id"`
	CenterLat   float64 `json:"center_lat" binding:"required"`
	CenterLng   float64 `json:"center_lng" binding:"required"`
	RadiusMeters float64 `json:"radius_meters" binding:"required"`
	AlertOn     string  `json:"alert_on"`
	Severity    int     `json:"severity"`
	Enabled     *bool   `json:"enabled"`
	NotifyWays  string  `json:"notify_ways"`
	Remark      string  `json:"remark"`
}

type GeofenceUpdateRequest struct {
	Name        string  `json:"name"`
	DeviceID    string  `json:"device_id"`
	CenterLat   float64 `json:"center_lat"`
	CenterLng   float64 `json:"center_lng"`
	RadiusMeters float64 `json:"radius_meters"`
	AlertOn     string  `json:"alert_on"`
	Severity    int     `json:"severity"`
	Enabled     *bool   `json:"enabled"`
	NotifyWays  string  `json:"notify_ways"`
	Remark      string  `json:"remark"`
}

type BindDeviceRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

// ============ 围栏CRUD ============

// List 获取围栏列表
func (c *GeofenceController) List(ctx *gin.Context) {
	var req GeofenceListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	var list []models.GeofenceRule
	var total int64

	query := c.DB.Model(&models.GeofenceRule{})

	if req.Keyword != "" {
		query = query.Where("name ILIKE ?", "%"+req.Keyword+"%")
	}
	if req.DeviceID != "" {
		query = query.Where("device_id = ?", req.DeviceID)
	}
	if req.Enabled != nil {
		query = query.Where("enabled = ?", *req.Enabled)
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("created_at DESC").Find(&list)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": list,
			"pagination": gin.H{"total": total, "page": req.Page, "page_size": req.PageSize},
		},
	})
}

// Get 获取围栏详情
func (c *GeofenceController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.GeofenceRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "围栏不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取关联的设备列表（通过device_id字段，这里简化为单设备）
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"rule":      rule,
			"bound_devices": []string{rule.DeviceID},
		},
	})
}

// Create 创建围栏
func (c *GeofenceController) Create(ctx *gin.Context) {
	var req GeofenceCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.AlertOn == "" {
		req.AlertOn = "enter"
	}
	if req.Severity == 0 {
		req.Severity = 2
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	rule := models.GeofenceRule{
		Name:         req.Name,
		DeviceID:     req.DeviceID,
		CenterLat:    req.CenterLat,
		CenterLng:    req.CenterLng,
		RadiusMeters: req.RadiusMeters,
		AlertOn:      req.AlertOn,
		Severity:     req.Severity,
		Enabled:      enabled,
		NotifyWays:   req.NotifyWays,
		Remark:       req.Remark,
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// Update 更新围栏
func (c *GeofenceController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.GeofenceRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "围栏不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req GeofenceUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.DeviceID != "" {
		updates["device_id"] = req.DeviceID
	}
	if req.CenterLat != 0 {
		updates["center_lat"] = req.CenterLat
	}
	if req.CenterLng != 0 {
		updates["center_lng"] = req.CenterLng
	}
	if req.RadiusMeters > 0 {
		updates["radius_meters"] = req.RadiusMeters
	}
	if req.AlertOn != "" {
		updates["alert_on"] = req.AlertOn
	}
	if req.Severity > 0 {
		updates["severity"] = req.Severity
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.NotifyWays != "" {
		updates["notify_ways"] = req.NotifyWays
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := c.DB.Model(&rule).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&rule, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": rule})
}

// Delete 删除围栏
func (c *GeofenceController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.GeofenceRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "围栏不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Delete(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// BindDevice 绑定设备到围栏
func (c *GeofenceController) BindDevice(ctx *gin.Context) {
	ruleID := ctx.Param("id")
	var rule models.GeofenceRule
	if err := c.DB.First(&rule, ruleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "围栏不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req BindDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 更新围栏的device_id（简化版：单设备绑定）
	if err := c.DB.Model(&rule).Update("device_id", req.DeviceID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "绑定设备失败"})
		return
	}

	c.DB.First(&rule, ruleID)
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "设备绑定成功",
		"data": gin.H{
			"rule_id":   rule.ID,
			"device_id": rule.DeviceID,
		},
	})
}

// ============ 围栏告警 ============

// ListAlerts 获取围栏告警列表
func (c *GeofenceController) ListAlerts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var alerts []models.GeofenceAlert
	var total int64

	query := c.DB.Model(&models.GeofenceAlert{})

	if ruleID := ctx.Query("rule_id"); ruleID != "" {
		query = query.Where("rule_id = ?", ruleID)
	}
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if alertType := ctx.Query("alert_type"); alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&alerts)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list": alerts,
			"pagination": gin.H{"total": total, "page": page, "page_size": pageSize},
		},
	})
}

// UpdateAlertStatus 更新告警状态
func (c *GeofenceController) UpdateAlertStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.GeofenceAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Status int `json:"status"` // 1:未处理 2:已确认 3:已解决 4:已忽略
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Model(&alert).Update("status", req.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}
