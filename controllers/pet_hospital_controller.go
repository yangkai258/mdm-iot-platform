package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetHospitalController 宠物医疗控制器
type PetHospitalController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册宠物医疗相关路由
func (p *PetHospitalController) RegisterRoutes(r *gin.RouterGroup) {
	hospital := r.Group("/pet-hospital")
	{
		hospital.GET("/appointments", p.ListAppointments)
		hospital.POST("/appointments", p.CreateAppointment)
		hospital.GET("/appointments/:id", p.GetAppointment)
		hospital.PUT("/appointments/:id/cancel", p.CancelAppointment)
		hospital.PUT("/appointments/:id/checkin", p.CheckInAppointment)
		hospital.POST("/appointments/:id/complete", p.CompleteAppointment)
	}
}

// CreateAppointmentRequest 创建预约请求
type CreateAppointmentRequest struct {
	PetUUID        string `json:"pet_uuid" binding:"required"`
	PetName        string `json:"pet_name"`
	HospitalName   string `json:"hospital_name" binding:"required"`
	HospitalAddress string `json:"hospital_address"`
	HospitalPhone  string `json:"hospital_phone"`
	Department     string `json:"department"`
	DoctorName     string `json:"doctor_name"`
	AppointmentAt  string `json:"appointment_at" binding:"required"` // RFC3339 格式
	Reason         string `json:"reason"`
	HouseholdID    *uint  `json:"household_id"`
}

// CreateAppointment 创建预约
func (p *PetHospitalController) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	appointmentAt, err := time.Parse(time.RFC3339, req.AppointmentAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预约时间格式错误，请使用 RFC3339 格式"})
		return
	}

	appointment := models.HospitalAppointment{
		PetUUID:         req.PetUUID,
		PetName:        req.PetName,
		HospitalName:   req.HospitalName,
		HospitalAddress: req.HospitalAddress,
		HospitalPhone:  req.HospitalPhone,
		Department:     req.Department,
		DoctorName:     req.DoctorName,
		AppointmentAt:  appointmentAt,
		Reason:         req.Reason,
		Status:         "pending",
		HouseholdID:    req.HouseholdID,
		OwnerID:        userID,
		TenantID:       tenantID,
	}

	if err := p.DB.Create(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建预约失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "预约创建成功",
		"data": appointment,
	})
}

// ListAppointments 获取预约列表
func (p *PetHospitalController) ListAppointments(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var appointments []models.HospitalAppointment
	query := p.DB.Where("owner_id = ? AND tenant_id = ?", userID, tenantID)

	// 过滤参数
	if petUUID := c.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if department := c.Query("department"); department != "" {
		query = query.Where("department = ?", department)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("appointment_at >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("appointment_at <= ?", endDate)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.HospitalAppointment{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("appointment_at DESC").Offset(offset).Limit(pageSize).Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询预约列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": appointments,
			"pagination": gin.H{
				"page":       page,
				"page_size":  pageSize,
				"total":      total,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetAppointment 获取预约详情
func (p *PetHospitalController) GetAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var appointment models.HospitalAppointment
	if err := p.DB.Where("(appointment_uuid = ? OR id = ?) AND owner_id = ? AND tenant_id = ?",
		appointmentID, appointmentID, userID, tenantID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询预约失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": appointment})
}

// CancelAppointmentRequest 取消预约请求
type CancelAppointmentRequest struct {
	CancelReason string `json:"cancel_reason"`
}

// CancelAppointment 取消预约
func (p *PetHospitalController) CancelAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var req CancelAppointmentRequest
	c.ShouldBindJSON(&req)

	var appointment models.HospitalAppointment
	if err := p.DB.Where("(appointment_uuid = ? OR id = ?) AND owner_id = ? AND tenant_id = ?",
		appointmentID, appointmentID, userID, tenantID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询预约失败"})
		return
	}

	if appointment.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预约已取消"})
		return
	}
	if appointment.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预约已完成，无法取消"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        "cancelled",
		"cancel_reason": req.CancelReason,
		"cancelled_at":  now,
	}

	if err := p.DB.Model(&appointment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消预约失败"})
		return
	}

	p.DB.First(&appointment, appointment.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "预约已取消", "data": appointment})
}

// CheckInAppointment 签到预约
func (p *PetHospitalController) CheckInAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var appointment models.HospitalAppointment
	if err := p.DB.Where("(appointment_uuid = ? OR id = ?) AND owner_id = ? AND tenant_id = ?",
		appointmentID, appointmentID, userID, tenantID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询预约失败"})
		return
	}

	if appointment.Status != "pending" && appointment.Status != "confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "当前状态无法签到"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":     "confirmed",
		"check_in_at": now,
	}

	if err := p.DB.Model(&appointment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "签到失败"})
		return
	}

	p.DB.First(&appointment, appointment.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "签到成功", "data": appointment})
}

// CompleteAppointmentRequest 完成预约请求
type CompleteAppointmentRequest struct {
	Diagnosis    string  `json:"diagnosis"`
	Prescription string  `json:"prescription"`
	Cost         float64 `json:"cost"`
	Notes        string  `json:"notes"`
}

// CompleteAppointment 完成预约
func (p *PetHospitalController) CompleteAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var req CompleteAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var appointment models.HospitalAppointment
	if err := p.DB.Where("(appointment_uuid = ? OR id = ?) AND owner_id = ? AND tenant_id = ?",
		appointmentID, appointmentID, userID, tenantID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询预约失败"})
		return
	}

	if appointment.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预约已取消"})
		return
	}
	if appointment.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "预约已完成"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        "completed",
		"completed_at":  now,
		"diagnosis":    req.Diagnosis,
		"prescription": req.Prescription,
		"cost":         req.Cost,
		"notes":        req.Notes,
	}

	if err := p.DB.Model(&appointment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "完成预约失败"})
		return
	}

	p.DB.First(&appointment, appointment.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "预约已完成", "data": appointment})
}
