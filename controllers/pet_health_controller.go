package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetHealthCtrl 宠物健康提醒控制器
type PetHealthCtrl struct {
	DB *gorm.DB
}

// RegisterPetHealthRoutes 注册宠物健康路由
func (p *PetHealthCtrl) RegisterPetHealthRoutes(r *gin.RouterGroup) {
	r.GET("/pets/:pet_id/health/reminders", p.ListReminders)
	r.POST("/pets/:pet_id/health/reminders", p.CreateReminder)
	r.PUT("/pets/:pet_id/health/reminders/:id", p.UpdateReminder)
	r.DELETE("/pets/:pet_id/health/reminders/:id", p.DeleteReminder)
	r.POST("/pets/:pet_id/health/checkup", p.RecordCheckup)
	r.GET("/pets/:pet_id/health/checkups", p.ListCheckups)
}

// verifyPetOwnership 验证宠物归属
func (p *PetHealthCtrl) verifyPetOwnership(petID string, userID uint, tenantID string) (*models.Pet, error) {
	var pet models.Pet
	err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error
	return &pet, err
}

// ListReminders 获取健康提醒列表
func (p *PetHealthCtrl) ListReminders(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	_, err := p.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	var reminders []models.PetHealthReminder
	query := p.DB.Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)

	// 过滤：已完成/未完成
	if completed := c.Query("completed"); completed != "" {
		if completed == "true" {
			query = query.Where("is_completed = ?", true)
		} else {
			query = query.Where("is_completed = ?", false)
		}
	}

	// 过滤：提醒类型
	if reminderType := c.Query("type"); reminderType != "" {
		query = query.Where("reminder_type = ?", reminderType)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var total int64
	p.DB.Model(&models.PetHealthReminder{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID).Count(&total)

	if err := query.Order("scheduled_at ASC").Offset(offset).Limit(pageSize).Find(&reminders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取提醒列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      reminders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateReminder 创建健康提醒
func (p *PetHealthCtrl) CreateReminder(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	_, err := p.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "验证失败"})
		return
	}

	var input struct {
		ReminderType   string `json:"reminder_type" binding:"required"`
		Title          string `json:"title" binding:"required"`
		Description    string `json:"description"`
		ScheduledAt    string `json:"scheduled_at" binding:"required"`
		RepeatInterval string `json:"repeat_interval"` // none/daily/weekly/monthly/yearly
		Notes          string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	scheduledAt, err := time.Parse(time.RFC3339, input.ScheduledAt)
	if err != nil {
		scheduledAt, _ = time.Parse("2006-01-02T15:04:05Z", input.ScheduledAt)
	}

	repeatInterval := input.RepeatInterval
	if repeatInterval == "" {
		repeatInterval = "none"
	}

	reminder := models.PetHealthReminder{
		PetUUID:      petID,
		ReminderType: input.ReminderType,
		Title:        input.Title,
		Description:  input.Description,
		ScheduledAt:  scheduledAt,
		RepeatInterval: repeatInterval,
		IsCompleted: false,
		Notes:        input.Notes,
		TenantID:     tenantID,
	}

	if err := p.DB.Create(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建提醒失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": reminder})
}

// UpdateReminder 更新提醒
func (p *PetHealthCtrl) UpdateReminder(c *gin.Context) {
	petID := c.Param("pet_id")
	reminderIDStr := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	reminderID, err := strconv.ParseUint(reminderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	_, err = p.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}

	var reminder models.PetHealthReminder
	err = p.DB.Where("id = ? AND pet_uuid = ? AND tenant_id = ?", uint(reminderID), petID, tenantID).First(&reminder).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "提醒不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	var input struct {
		ReminderType   string `json:"reminder_type"`
		Title           string `json:"title"`
		Description     string `json:"description"`
		ScheduledAt     string `json:"scheduled_at"`
		RepeatInterval string `json:"repeat_interval"`
		IsCompleted     *bool  `json:"is_completed"`
		Notes           string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if input.ReminderType != "" { updates["reminder_type"] = input.ReminderType }
	if input.Title != "" { updates["title"] = input.Title }
	if input.Description != "" { updates["description"] = input.Description }
	if input.RepeatInterval != "" { updates["repeat_interval"] = input.RepeatInterval }
	if input.Notes != "" { updates["notes"] = input.Notes }
	if input.IsCompleted != nil {
		updates["is_completed"] = *input.IsCompleted
		if *input.IsCompleted {
			now := time.Now()
			updates["completed_at"] = &now
		}
	}

	if input.ScheduledAt != "" {
		if t, err := time.Parse(time.RFC3339, input.ScheduledAt); err == nil {
			updates["scheduled_at"] = t
		}
	}

	if len(updates) > 0 {
		p.DB.Model(&reminder).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": reminder})
}

// DeleteReminder 删除提醒
func (p *PetHealthCtrl) DeleteReminder(c *gin.Context) {
	petID := c.Param("pet_id")
	reminderIDStr := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	reminderID, err := strconv.ParseUint(reminderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的ID"})
		return
	}

	_, err = p.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}

	var reminder models.PetHealthReminder
	err = p.DB.Where("id = ? AND pet_uuid = ? AND tenant_id = ?", uint(reminderID), petID, tenantID).First(&reminder).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "提醒不存在"})
		return
	}

	if err := p.DB.Delete(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// RecordCheckup 记录宠物体检
func (p *PetHealthCtrl) RecordCheckup(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	_, err := p.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "验证失败"})
		return
	}

	var input struct {
		CheckupDate   string   `json:"checkup_date" binding:"required"`
		Hospital      string   `json:"hospital"`
		Doctor        string   `json:"doctor"`
		Weight        float64  `json:"weight"`
		Symptoms      string   `json:"symptoms"`
		Diagnosis     string   `json:"diagnosis"`
		Treatment     string   `json:"treatment"`
		Prescription  string   `json:"prescription"`
		Cost          float64  `json:"cost"`
		NextDate      string   `json:"next_date"`
		Attachments   []string `json:"attachments"`
		Notes         string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	checkupDate, err := time.Parse("2006-01-02", input.CheckupDate)
	if err != nil {
		checkupDate, _ = time.Parse(time.RFC3339[:10], input.CheckupDate)
	}

	checkup := models.PetCheckup{
		PetUUID:      petID,
		CheckupDate:  checkupDate,
		Hospital:     input.Hospital,
		Doctor:       input.Doctor,
		Weight:       input.Weight,
		Symptoms:     input.Symptoms,
		Diagnosis:    input.Diagnosis,
		Treatment:    input.Treatment,
		Prescription: input.Prescription,
		Cost:         input.Cost,
		Attachments:  input.Attachments,
		Notes:        input.Notes,
		TenantID:     tenantID,
	}

	if input.NextDate != "" {
		t, err := time.Parse("2006-01-02", input.NextDate)
		if err == nil {
			checkup.NextDate = &t
		}
	}

	if err := p.DB.Create(&checkup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "记录体检失败"})
		return
	}

	// 如果有体重更新，同步到宠物档案
	if input.Weight > 0 {
		p.DB.Model(&models.Pet{}).Where("pet_uuid = ?", petID).Update("weight", input.Weight)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "体检记录已保存", "data": checkup})
}

// ListCheckups 获取体检记录列表
func (p *PetHealthCtrl) ListCheckups(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	_, err := p.verifyPetOwnership(petID, userID, tenantID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var checkups []models.PetCheckup
	var total int64
	p.DB.Model(&models.PetCheckup{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID).Count(&total)

	if err := p.DB.Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID).
		Order("checkup_date DESC").Offset(offset).Limit(pageSize).Find(&checkups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取体检记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      checkups,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
