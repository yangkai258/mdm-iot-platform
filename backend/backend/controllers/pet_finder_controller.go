package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetFinderController 寻宠网络控制器
type PetFinderController struct {
	DB *gorm.DB
}

// ============ PetFinderReport ============

// PetFinderReportList 寻宠报告列表
func (c *PetFinderController) PetFinderReportList(ctx *gin.Context) {
	var reports []models.PetFinderReport
	var total int64

	query := c.DB.Model(&models.PetFinderReport{})

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if userID := ctx.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if reportType := ctx.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Preload("Sightings").Preload("Alerts").Offset(offset).Limit(pageSize).Order("id DESC").Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": reports, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetFinderReportGet 获取报告
func (c *PetFinderController) PetFinderReportGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var report models.PetFinderReport
	if err := c.DB.Preload("Sightings").Preload("Alerts").First(&report, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	// 增加浏览量
	c.DB.Model(&report).Update("view_count", report.ViewCount+1)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// PetFinderReportCreate 创建报告
func (c *PetFinderController) PetFinderReportCreate(ctx *gin.Context) {
	var report models.PetFinderReport
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

// PetFinderReportUpdate 更新报告
func (c *PetFinderController) PetFinderReportUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var report models.PetFinderReport
	if err := c.DB.First(&report, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	var updateData models.PetFinderReport
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

// PetFinderReportResolve 关闭报告
func (c *PetFinderController) PetFinderReportResolve(ctx *gin.Context) {
	id := ctx.Param("id")
	var report models.PetFinderReport
	if err := c.DB.First(&report, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	now := time.Now()
	report.Status = "resolved"
	report.ResolvedAt = &now
	if err := c.DB.Save(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// PetFinderReportDelete 删除报告
func (c *PetFinderController) PetFinderReportDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetFinderReport{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetFinderSighting ============

// PetFinderSightingList 目击记录列表
func (c *PetFinderController) PetFinderSightingList(ctx *gin.Context) {
	var sightings []models.PetFinderSighting
	var total int64

	query := c.DB.Model(&models.PetFinderSighting{})

	if reportID := ctx.Query("report_id"); reportID != "" {
		query = query.Where("report_id = ?", reportID)
	}
	if isVerified := ctx.Query("is_verified"); isVerified != "" {
		query = query.Where("is_verified = ?", isVerified == "true")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sighted_at DESC").Find(&sightings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": sightings, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetFinderSightingCreate 报告目击
func (c *PetFinderController) PetFinderSightingCreate(ctx *gin.Context) {
	var sighting models.PetFinderSighting
	if err := ctx.ShouldBindJSON(&sighting); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&sighting).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": sighting})
}

// PetFinderSightingVerify 验证目击记录
func (c *PetFinderController) PetFinderSightingVerify(ctx *gin.Context) {
	id := ctx.Param("id")
	var sighting models.PetFinderSighting
	if err := c.DB.First(&sighting, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	sighting.IsVerified = true
	if err := c.DB.Save(&sighting).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": sighting})
}

// ============ PetFinderAlert ============

// PetFinderAlertList 警报列表
func (c *PetFinderController) PetFinderAlertList(ctx *gin.Context) {
	var alerts []models.PetFinderAlert
	var total int64

	query := c.DB.Model(&models.PetFinderAlert{})

	if reportID := ctx.Query("report_id"); reportID != "" {
		query = query.Where("report_id = ?", reportID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if err := query.Order("id DESC").Find(&alerts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": alerts, "total": total,
	}})
}

// ============ SightingReport ============

// SightingReportList 目击报告列表
func (c *PetFinderController) SightingReportList(ctx *gin.Context) {
	var reports []models.SightingReport
	var total int64

	query := c.DB.Model(&models.SightingReport{})

	if reportUUID := ctx.Query("report_uuid"); reportUUID != "" {
		query = query.Where("report_uuid = ?", reportUUID)
	}
	if isCredible := ctx.Query("is_credible"); isCredible != "" {
		query = query.Where("is_credible = ?", isCredible == "true")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": reports, "total": total, "page": page, "page_size": pageSize,
	}})
}

// SightingReportCreate 创建目击报告
func (c *PetFinderController) SightingReportCreate(ctx *gin.Context) {
	var report models.SightingReport
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

// SightingReportDelete 删除目击报告
func (c *PetFinderController) SightingReportDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.SightingReport{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
