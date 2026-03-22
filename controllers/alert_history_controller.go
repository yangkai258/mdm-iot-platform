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
// GET /api/v1/alerts/history
func (c *AlertHistoryController) GetAlertHistory(ctx *gin.Context) {
	var history []models.AlertHistory
	query := c.DB.Model(&models.AlertHistory{})

	// 筛选条件
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if alertType := ctx.Query("type"); alertType != "" {
		query = query.Where("alert_type = ?", alertType)
	}
	if severity := ctx.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status := ctx.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if ruleID := ctx.Query("rule_id"); ruleID != "" {
		query = query.Where("rule_id = ?", ruleID)
	}

	// 时间范围
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		} else if t2, err := time.Parse("2006-01-02", startTime); err == nil {
			query = query.Where("created_at >= ?", t2)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z07:00", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		} else if t2, err := time.Parse("2006-01-02", endTime); err == nil {
			query = query.Where("created_at <= ?", t2.Add(24*time.Hour))
		}
	}

	// 分页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&history)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      history,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAlertHistoryByID 获取单条告警历史详情
// GET /api/v1/alerts/history/:id
func (c *AlertHistoryController) GetAlertHistoryByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var history models.AlertHistory
	if err := c.DB.First(&history, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "告警历史不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    history,
	})
}

// ArchiveAlert 将已解决的告警归档到历史表
// POST /api/v1/alerts/history/archive
func (c *AlertHistoryController) ArchiveAlert(ctx *gin.Context) {
	var req struct {
		AlertID uint `json:"alert_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 查询原始告警
	var alert models.DeviceAlert
	if err := c.DB.First(&alert, req.AlertID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "告警不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	// 检查是否已归档
	var existing models.AlertHistory
	if c.DB.Where("original_id = ?", req.AlertID).First(&existing).Error == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该告警已经归档",
		})
		return
	}

	// 创建归档记录
	now := time.Now()
	history := models.AlertHistory{
		OriginalID:     alert.ID,
		RuleID:        alert.RuleID,
		DeviceID:      alert.DeviceID,
		AlertType:     alert.AlertType,
		Severity:      alert.Severity,
		Message:       alert.Message,
		TriggerValue:  strconv.FormatFloat(alert.TriggerVal, 'f', -1, 64),
		Threshold:     strconv.FormatFloat(alert.Threshold, 'f', -1, 64),
		Status:        alert.Status,
		ConfirmedAt:   alert.ConfirmedAt,
		ConfirmedBy:   0,
		ResolvedAt:    alert.ResolvedAt,
		ResolvedBy:    0,
		CreatedAt:    alert.CreatedAt,
		ResolvedAtH:  &now,
		ArchivedAt:   now,
	}

	if err := c.DB.Create(&history).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "归档失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "归档成功",
		"data":    history,
	})
}
