package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthController 健康管理控制器
type HealthController struct {
	DB *gorm.DB
}

// ============ HealthWarning ============

// HealthWarningList 健康预警列表
func (c *HealthController) HealthWarningList(ctx *gin.Context) {
	var warnings []models.HealthWarning
	var total int64

	query := c.DB.Model(&models.HealthWarning{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if level := ctx.Query("level"); level != "" {
		query = query.Where("level = ?", level)
	}
	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("priority DESC, start_time DESC").Find(&warnings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": warnings, "total": total, "page": page, "page_size": pageSize,
	}})
}

// HealthWarningGet 获取预警
func (c *HealthController) HealthWarningGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var warning models.HealthWarning
	if err := c.DB.First(&warning, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": warning})
}

// HealthWarningCreate 创建预警
func (c *HealthController) HealthWarningCreate(ctx *gin.Context) {
	var warning models.HealthWarning
	if err := ctx.ShouldBindJSON(&warning); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&warning).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": warning})
}

// HealthWarningUpdate 更新预警
func (c *HealthController) HealthWarningUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var warning models.HealthWarning
	if err := c.DB.First(&warning, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}
	var updateData models.HealthWarning
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

// HealthWarningResolve 解决预警
func (c *HealthController) HealthWarningResolve(ctx *gin.Context) {
	id := ctx.Param("id")
	var warning models.HealthWarning
	if err := c.DB.First(&warning, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预警不存在"})
		return
	}
	now := time.Now()
	warning.Status = "resolved"
	warning.ResolvedAt = &now
	if err := c.DB.Save(&warning).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": warning})
}

// HealthWarningDelete 删除预警
func (c *HealthController) HealthWarningDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.HealthWarning{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ VitalRecord ============

// VitalRecordList 体征记录列表
func (c *HealthController) VitalRecordList(ctx *gin.Context) {
	var records []models.VitalRecord
	var total int64

	query := c.DB.Model(&models.VitalRecord{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if vitalType := ctx.Query("vital_type"); vitalType != "" {
		query = query.Where("vital_type = ?", vitalType)
	}
	if isAbnormal := ctx.Query("is_abnormal"); isAbnormal != "" {
		query = query.Where("is_abnormal = ?", isAbnormal == "true")
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

// VitalRecordGet 获取体征记录
func (c *HealthController) VitalRecordGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.VitalRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// VitalRecordCreate 创建体征记录
func (c *HealthController) VitalRecordCreate(ctx *gin.Context) {
	var record models.VitalRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// VitalRecordDelete 删除体征记录
func (c *HealthController) VitalRecordDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.VitalRecord{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// VitalTrendList 体征趋势
func (c *HealthController) VitalTrendList(ctx *gin.Context) {
	petUUID := ctx.Query("pet_uuid")
	vitalType := ctx.Query("vital_type")
	var trends []models.VitalTrend

	query := c.DB.Model(&models.VitalTrend{})
	if petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if vitalType != "" {
		query = query.Where("vital_type = ?", vitalType)
	}

	if err := query.Order("last_updated DESC").Find(&trends).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": trends})
}

// ============ ExerciseGoal ============

// ExerciseGoalList 运动目标列表
func (c *HealthController) ExerciseGoalList(ctx *gin.Context) {
	var goals []models.ExerciseGoal
	var total int64

	query := c.DB.Model(&models.ExerciseGoal{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("priority DESC, id DESC").Find(&goals).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": goals, "total": total, "page": page, "page_size": pageSize,
	}})
}

// ExerciseGoalGet 获取目标
func (c *HealthController) ExerciseGoalGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var goal models.ExerciseGoal
	if err := c.DB.First(&goal, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "目标不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": goal})
}

// ExerciseGoalCreate 创建目标
func (c *HealthController) ExerciseGoalCreate(ctx *gin.Context) {
	var goal models.ExerciseGoal
	if err := ctx.ShouldBindJSON(&goal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&goal).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": goal})
}

// ExerciseGoalUpdate 更新目标
func (c *HealthController) ExerciseGoalUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var goal models.ExerciseGoal
	if err := c.DB.First(&goal, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "目标不存在"})
		return
	}
	var updateData models.ExerciseGoal
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

// ExerciseGoalDelete 删除目标
func (c *HealthController) ExerciseGoalDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.ExerciseGoal{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ ExerciseRecord ============

// ExerciseRecordList 运动记录列表
func (c *HealthController) ExerciseRecordList(ctx *gin.Context) {
	var records []models.ExerciseRecord
	var total int64

	query := c.DB.Model(&models.ExerciseRecord{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if exerciseType := ctx.Query("exercise_type"); exerciseType != "" {
		query = query.Where("exercise_type = ?", exerciseType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("start_time DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

// ExerciseRecordGet 获取运动记录
func (c *HealthController) ExerciseRecordGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.ExerciseRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// ExerciseRecordCreate 创建运动记录
func (c *HealthController) ExerciseRecordCreate(ctx *gin.Context) {
	var record models.ExerciseRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// ExerciseRecordDelete 删除运动记录
func (c *HealthController) ExerciseRecordDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.ExerciseRecord{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ExerciseSummaryList 运动摘要列表
func (c *HealthController) ExerciseSummaryList(ctx *gin.Context) {
	var summaries []models.ExerciseSummary
	var total int64

	query := c.DB.Model(&models.ExerciseSummary{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("summary_date DESC").Find(&summaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": summaries, "total": total, "page": page, "page_size": pageSize,
	}})
}

// ExerciseTrendList 运动趋势列表
func (c *HealthController) ExerciseTrendList(ctx *gin.Context) {
	var trends []models.ExerciseTrend
	if err := c.DB.Order("date DESC").Limit(30).Find(&trends).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": trends})
}

// ============ HealthAlertRule ============

// HealthAlertRuleList 健康告警规则列表
func (c *HealthController) HealthAlertRuleList(ctx *gin.Context) {
	var rules []models.HealthAlertRule
	var total int64

	query := c.DB.Model(&models.HealthAlertRule{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if isEnabled := ctx.Query("is_enabled"); isEnabled != "" {
		query = query.Where("is_enabled = ?", isEnabled == "true")
	}

	query.Count(&total)

	if err := query.Order("priority DESC, id DESC").Find(&rules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": rules, "total": total,
	}})
}

// HealthAlertRuleGet 获取规则
func (c *HealthController) HealthAlertRuleGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.HealthAlertRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// HealthAlertRuleCreate 创建规则
func (c *HealthController) HealthAlertRuleCreate(ctx *gin.Context) {
	var rule models.HealthAlertRule
	if err := ctx.ShouldBindJSON(&rule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// HealthAlertRuleUpdate 更新规则
func (c *HealthController) HealthAlertRuleUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var rule models.HealthAlertRule
	if err := c.DB.First(&rule, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "规则不存在"})
		return
	}
	var updateData models.HealthAlertRule
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

// HealthAlertRuleDelete 删除规则
func (c *HealthController) HealthAlertRuleDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.HealthAlertRule{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ HealthMonitorSetting ============

// HealthMonitorSettingList 监控设置列表
func (c *HealthController) HealthMonitorSettingList(ctx *gin.Context) {
	var settings []models.HealthMonitorSetting
	var total int64

	query := c.DB.Model(&models.HealthMonitorSetting{})
	query.Count(&total)

	if err := query.Order("id DESC").Find(&settings).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": settings, "total": total,
	}})
}

// HealthMonitorSettingGet 获取监控设置
func (c *HealthController) HealthMonitorSettingGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var setting models.HealthMonitorSetting
	if err := c.DB.First(&setting, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": setting})
}

// HealthMonitorSettingCreate 创建监控设置
func (c *HealthController) HealthMonitorSettingCreate(ctx *gin.Context) {
	var setting models.HealthMonitorSetting
	if err := ctx.ShouldBindJSON(&setting); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&setting).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": setting})
}

// HealthMonitorSettingUpdate 更新监控设置
func (c *HealthController) HealthMonitorSettingUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var setting models.HealthMonitorSetting
	if err := c.DB.First(&setting, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设置不存在"})
		return
	}
	var updateData models.HealthMonitorSetting
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

// ============ HealthAlert ============

// HealthAlertList 健康告警列表
func (c *HealthController) HealthAlertList(ctx *gin.Context) {
	var alerts []models.HealthAlert
	var total int64

	query := c.DB.Model(&models.HealthAlert{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if alertLevel := ctx.Query("alert_level"); alertLevel != "" {
		query = query.Where("alert_level = ?", alertLevel)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("urgency DESC, occurred_at DESC").Find(&alerts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": alerts, "total": total, "page": page, "page_size": pageSize,
	}})
}

// HealthAlertGet 获取告警
func (c *HealthController) HealthAlertGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.HealthAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": alert})
}

// HealthAlertCreate 创建告警
func (c *HealthController) HealthAlertCreate(ctx *gin.Context) {
	var alert models.HealthAlert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&alert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": alert})
}

// HealthAlertResolve 解决告警
func (c *HealthController) HealthAlertResolve(ctx *gin.Context) {
	id := ctx.Param("id")
	var alert models.HealthAlert
	if err := c.DB.First(&alert, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "告警不存在"})
		return
	}
	now := time.Now()
	alert.Status = "resolved"
	alert.ResolvedAt = &now
	if err := c.DB.Save(&alert).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": alert})
}

// HealthAlertDelete 删除告警
func (c *HealthController) HealthAlertDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.HealthAlert{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetVaccination ============

// PetVaccinationList 疫苗接种列表
func (c *HealthController) PetVaccinationList(ctx *gin.Context) {
	var records []models.PetVaccination
	var total int64

	query := c.DB.Model(&models.PetVaccination{})

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("inoculation_date DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetVaccinationCreate 创建疫苗记录
func (c *HealthController) PetVaccinationCreate(ctx *gin.Context) {
	var record models.PetVaccination
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// PetVaccinationDelete 删除疫苗记录
func (c *HealthController) PetVaccinationDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetVaccination{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetDietRecord ============

// PetDietRecordList 饮食记录列表
func (c *HealthController) PetDietRecordList(ctx *gin.Context) {
	var records []models.PetDietRecord
	var total int64

	query := c.DB.Model(&models.PetDietRecord{})

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if mealType := ctx.Query("meal_type"); mealType != "" {
		query = query.Where("meal_type = ?", mealType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("eat_time DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetDietRecordCreate 创建饮食记录
func (c *HealthController) PetDietRecordCreate(ctx *gin.Context) {
	var record models.PetDietRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// PetDietRecordDelete 删除饮食记录
func (c *HealthController) PetDietRecordDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetDietRecord{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ PetHealthRecord ============

// PetHealthRecordList 健康记录列表
func (c *HealthController) PetHealthRecordList(ctx *gin.Context) {
	var records []models.PetHealthRecord
	var total int64

	query := c.DB.Model(&models.PetHealthRecord{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if recordType := ctx.Query("record_type"); recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("record_date DESC").Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": records, "total": total, "page": page, "page_size": pageSize,
	}})
}

// PetHealthRecordGet 获取健康记录
func (c *HealthController) PetHealthRecordGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.PetHealthRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// PetHealthRecordCreate 创建健康记录
func (c *HealthController) PetHealthRecordCreate(ctx *gin.Context) {
	var record models.PetHealthRecord
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// PetHealthRecordUpdate 更新健康记录
func (c *HealthController) PetHealthRecordUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var record models.PetHealthRecord
	if err := c.DB.First(&record, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	var updateData models.PetHealthRecord
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

// PetHealthRecordDelete 删除健康记录
func (c *HealthController) PetHealthRecordDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.PetHealthRecord{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
