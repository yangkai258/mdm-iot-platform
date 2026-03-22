package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlertHistoryController 告警历史控制器
type AlertHistoryController struct {
	DB *gorm.DB
}

// GetAlertHistory 获取告警历史列表
func (c *AlertHistoryController) GetAlertHistory(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var alerts []models.DeviceAlert
	query := c.DB.Model(&models.DeviceAlert{})

	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if alertType := ctx.Query("alert_type"); alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&alerts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":     alerts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAlertHistoryByID 获取告警历史详情
func (c *AlertHistoryController) GetAlertHistoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.DeviceAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "告警不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": alert})
}

// ArchiveAlert 归档告警
func (c *AlertHistoryController) ArchiveAlert(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.DeviceAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4041, "message": "告警不存在"})
		return
	}

	now := time.Now()
	alert.Status = 3 // 已归档
	alert.ResolvedAt = &now

	if err := c.DB.Save(&alert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "归档失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "归档成功"})
}
