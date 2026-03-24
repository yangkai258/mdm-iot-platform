package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmotionController 情感管理控制器
type EmotionController struct {
	DB *gorm.DB
}

// ============ EmotionRecord ============

// EmotionRecordList 情感记录列表
func (c *EmotionController) EmotionRecordList(ctx *gin.Context) {
	var records []models.EmotionRecord
	var total int64

	query := c.DB.Model(&models.EmotionRecord{})

	if subjectType := ctx.Query("subject_type"); subjectType != "" {
		query = query.Where("subject_type = ?", subjectType)
	}
	if subjectID := ctx.Query("subject_id"); subjectID != "" {
		query = query.Where("subject_id = ?", subjectID)
	}
	if emotionType := ctx.Query("emotion_type"); emotionType != "" {
		query = query.Where("emotion_type = ?", emotionType)
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("recorded_at >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("recorded_at <= ?", endDate)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("recorded_at DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

// EmotionRecordGet 获取记录
func (c *EmotionController) EmotionRecordGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.EmotionRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// EmotionRecordCreate 创建记录
func (c *EmotionController) EmotionRecordCreate(ctx *gin.Context) {
	var record models.EmotionRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if record.RecordedAt.IsZero() {
		record.RecordedAt = time.Now()
	}
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// EmotionRecordDelete 删除记录
func (c *EmotionController) EmotionRecordDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.EmotionRecord{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ EmotionReport ============

// EmotionReportList 情感报告列表
func (c *EmotionController) EmotionReportList(ctx *gin.Context) {
	var reports []models.EmotionReport
	var total int64

	query := c.DB.Model(&models.EmotionReport{})

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if reportType := ctx.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("generated_at DESC").Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": reports, "total": total, "page": page, "page_size": pageSize,
	}})
}

// EmotionReportGet 获取报告
func (c *EmotionController) EmotionReportGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var report models.EmotionReport
	if err := c.DB.First(&report, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// EmotionReportCreate 创建报告
func (c *EmotionController) EmotionReportCreate(ctx *gin.Context) {
	var report models.EmotionReport
	if err := ctx.ShouldBindJSON(&report); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// EmotionReportDelete 删除报告
func (c *EmotionController) EmotionReportDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.EmotionReport{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ EmotionResponseConfig ============

// EmotionResponseConfigList 情感响应配置列表
func (c *EmotionController) EmotionResponseConfigList(ctx *gin.Context) {
	var configs []models.EmotionResponseConfig
	var total int64

	query := c.DB.Model(&models.EmotionResponseConfig{})

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if enabled := ctx.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	query.Count(&total)

	if err := query.Order("id DESC").Find(&configs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": configs, "total": total,
	}})
}

// EmotionResponseConfigGet 获取配置
func (c *EmotionController) EmotionResponseConfigGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var config models.EmotionResponseConfig
	if err := c.DB.First(&config, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配置不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// EmotionResponseConfigCreate 创建配置
func (c *EmotionController) EmotionResponseConfigCreate(ctx *gin.Context) {
	var config models.EmotionResponseConfig
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&config).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": config})
}

// EmotionResponseConfigUpdate 更新配置
func (c *EmotionController) EmotionResponseConfigUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var config models.EmotionResponseConfig
	if err := c.DB.First(&config, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "配置不存在"})
		return
	}
	var updateData models.EmotionResponseConfig
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updateData})
}

// EmotionResponseConfigDelete 删除配置
func (c *EmotionController) EmotionResponseConfigDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.EmotionResponseConfig{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ 情感统计 ============

// EmotionStats 获取情感统计数据
func (c *EmotionController) EmotionStats(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	period := ctx.DefaultQuery("period", "7d") // 7d, 30d

	var startDate string
	now := time.Now()
	switch period {
	case "24h":
		startDate = now.AddDate(0, 0, -1).Format("2006-01-02")
	case "7d":
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
	case "30d":
		startDate = now.AddDate(0, 0, -30).Format("2006-01-02")
	default:
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	// 按情绪类型统计
	type EmotionStat struct {
		EmotionType string  `json:"emotion_type"`
		Count       int64   `json:"count"`
		AvgIntensity float64 `json:"avg_intensity"`
	}
	var emotionStats []EmotionStat

	query := c.DB.Model(&models.EmotionRecord{}).Where("recorded_at >= ?", startDate)
	if petID != "" {
		query = query.Where("subject_id = ? AND subject_type = 'pet'", petID)
	}

	query.Select("emotion_type, COUNT(*) as count, AVG(intensity) as avg_intensity").
		Group("emotion_type").
		Scan(&emotionStats)

	// 情绪趋势（按天）
	type EmotionTrend struct {
		Date       string `json:"date"`
		EmotionType string `json:"emotion_type"`
		Count     int64  `json:"count"`
	}
	var trend []EmotionTrend
	rows, err := c.DB.Model(&models.EmotionRecord{}).
		Select("DATE(recorded_at) as date, emotion_type, COUNT(*) as count").
		Where("recorded_at >= ?", startDate).
		Group("DATE(recorded_at), emotion_type").
		Order("date ASC").
		Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t EmotionTrend
			rows.Scan(&t.Date, &t.EmotionType, &t.Count)
			trend = append(trend, t)
		}
	}

	// 总体统计
	var totalStats struct {
		TotalRecords  int64
		AvgIntensity float64
		AvgConfidence float64
	}
	c.DB.Model(&models.EmotionRecord{}).
		Select("COUNT(*) as total_records, AVG(intensity) as avg_intensity, AVG(confidence) as avg_confidence").
		Where("recorded_at >= ?", startDate).
		Scan(&totalStats)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"period":     period,
			"start_date": startDate,
			"end_date":   now.Format("2006-01-02"),
			"total":      totalStats,
			"by_emotion_type": emotionStats,
			"trend":       trend,
		},
	})
}
