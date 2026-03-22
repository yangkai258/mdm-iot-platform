package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/timezone"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TimezoneController 时区控制器
type TimezoneController struct {
	DB  *gorm.DB
	svc *timezone.TimezoneService
}

// NewTimezoneController 创建时区控制器
func NewTimezoneController(db *gorm.DB) *TimezoneController {
	return &TimezoneController{
		DB:  db,
		svc: timezone.NewTimezoneService(db),
	}
}

// GetCurrentUserTimezone 获取当前用户时区
// @Summary 获取当前用户时区
// @Tags timezone
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/timezone [get]
func (ctrl *TimezoneController) GetCurrentUserTimezone(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	tz, err := ctrl.svc.GetUserTimezone(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回租户时区或默认
			tz, _ = ctrl.svc.GetDefaultTimezone()
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"timezone": tz,
			"entity_type": "user",
			"entity_id": userID,
		},
	})
}

// UpdateCurrentUserTimezone 更新当前用户时区
// @Summary 更新当前用户时区
// @Tags timezone
// @Accept json
// @Produce json
// @Param timezone body map[string]string true "时区信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/timezone [put]
func (ctrl *TimezoneController) UpdateCurrentUserTimezone(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var req struct {
		Timezone string `json:"timezone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证时区有效性
	if !timezone.ValidateTimezone(req.Timezone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timezone"})
		return
	}

	if err := ctrl.svc.SetUserTimezone(userID, req.Timezone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"timezone": req.Timezone,
			"entity_type": "user",
			"entity_id": userID,
		},
	})
}

// GetSupportedTimezones 获取支持的时区列表
// @Summary 获取支持的时区列表
// @Tags timezone
// @Produce json
// @Success 200 {array} map[string]string
// @Router /api/v1/timezone/list [get]
func (ctrl *TimezoneController) GetSupportedTimezones(c *gin.Context) {
	timezones := timezone.GetSupportedTimezones()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": timezones,
	})
}

// GetTenantTimezone 获取租户时区
// @Summary 获取租户时区
// @Tags timezone
// @Produce json
// @Param id path int true "租户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/timezone/tenant/{id} [get]
func (ctrl *TimezoneController) GetTenantTimezone(c *gin.Context) {
	id := c.Param("id")
	tenantID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return
	}

	tz, err := ctrl.svc.GetTenantTimezone(uint(tenantID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tz = "UTC"
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"timezone": tz,
			"entity_type": "tenant",
			"entity_id": tenantID,
		},
	})
}

// UpdateTenantTimezone 更新租户时区
// @Summary 更新租户时区
// @Tags timezone
// @Accept json
// @Produce json
// @Param id path int true "租户ID"
// @Param timezone body map[string]string true "时区信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/timezone/tenant/{id} [put]
func (ctrl *TimezoneController) UpdateTenantTimezone(c *gin.Context) {
	id := c.Param("id")
	tenantID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant id"})
		return
	}

	var req struct {
		Timezone string `json:"timezone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证时区有效性
	if !timezone.ValidateTimezone(req.Timezone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timezone"})
		return
	}

	if err := ctrl.svc.SetTenantTimezone(uint(tenantID), req.Timezone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"timezone": req.Timezone,
			"entity_type": "tenant",
			"entity_id": tenantID,
		},
	})
}

// GetTimezoneConfig 获取时区配置
// @Summary 获取时区配置
// @Tags timezone
// @Produce json
// @Param entity_type query string true "实体类型 (user/tenant/system)"
// @Param entity_id query int true "实体ID"
// @Success 200 {object} models.TimezoneConfig
// @Router /api/v1/timezone/config [get]
func (ctrl *TimezoneController) GetTimezoneConfig(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")
	entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity id"})
		return
	}

	config, err := ctrl.svc.GetTimezoneConfig(entityType, uint(entityID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "timezone config not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config,
	})
}

// ListTenantTimezones 列出所有租户时区配置
// @Summary 列出所有租户时区配置
// @Tags timezone
// @Produce json
// @Success 200 {array} models.TimezoneConfig
// @Router /api/v1/timezone/tenant-configs [get]
func (ctrl *TimezoneController) ListTenantTimezones(c *gin.Context) {
	configs, err := ctrl.svc.ListTimezoneConfigs("tenant")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": configs,
	})
}

// ConvertTime 时区转换
// @Summary 时区转换
// @Tags timezone
// @Accept json
// @Produce json
// @Param time body map[string]interface{} true "时间转换请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/timezone/convert [post]
func (ctrl *TimezoneController) ConvertTime(c *gin.Context) {
	var req struct {
		Time     string `json:"time" binding:"required"`      // 原始时间
		FromTZ   string `json:"from_tz" binding:"required"`  // 源时区
		ToTZ     string `json:"to_tz" binding:"required"`     // 目标时区
		Format   string `json:"format"`                       // 格式化模板
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析时间
	t, err := parseTime(req.Time, req.FromTZ)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	// 转换到目标时区
	result := timezone.FormatTimeWithTimezone(t, req.ToTZ)
	if req.Format != "" {
		result["formatted"] = timezone.FormatInTimezone(t, req.ToTZ, req.Format)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// parseTime 解析时间字符串
func parseTime(timeStr, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}

	// 尝试多种格式
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unable to parse time")
}
