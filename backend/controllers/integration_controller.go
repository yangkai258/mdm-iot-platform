package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IntegrationController 第三方集成控制器
type IntegrationController struct {
	DB *gorm.DB
}

// ============ 请求/响应结构 ============

// DeviceControlRequest 设备控制请求
type DeviceControlRequest struct {
	Action string                 `json:"action" binding:"required"` // on/off/set
	Params map[string]interface{} `json:"params"`
}

// MedicalRecordRequest 创建/更新医疗记录请求
type MedicalRecordRequest struct {
	PetID        uint      `json:"pet_id" binding:"required"`
	RecordType   string    `json:"record_type" binding:"required"`
	HospitalName string    `json:"hospital_name"`
	DoctorName   string    `json:"doctor_name"`
	Diagnosis    string    `json:"diagnosis"`
	Treatment    string    `json:"treatment"`
	Prescription string    `json:"prescription"`
	RecordDate   time.Time `json:"record_date" binding:"required"`
	Attachments  []string  `json:"attachments"`
}

// InsuranceClaimRequest 创建/更新理赔请求
type InsuranceClaimRequest struct {
	PolicyID    uint      `json:"policy_id" binding:"required"`
	ClaimType   string    `json:"claim_type"`
	Amount      float64   `json:"amount" binding:"required"`
	Reason      string    `json:"reason"`
	Description string    `json:"description"`
	Diagnosis   string    `json:"diagnosis"`
	Documents   []string  `json:"documents"`
	ClaimDate   time.Time `json:"claim_date"`
}

// InsuranceProductRequest 创建/更新保险产品请求
type InsuranceProductRequest struct {
	Name         string  `json:"name" binding:"required"`
	Provider     string  `json:"provider" binding:"required"`
	CoverageType string  `json:"coverage_type" binding:"required"`
	MonthlyPrice float64 `json:"monthly_price" binding:"required"`
	CoverageLimit float64 `json:"coverage_limit" binding:"required"`
	Description  string  `json:"description"`
	Terms        string  `json:"terms"`
	IsActive     bool    `json:"is_active"`
}

// InsurancePolicyRequest 创建/购买保单请求
type InsurancePolicyRequest struct {
	ProductID   uint      `json:"product_id" binding:"required"`
	PetID       uint      `json:"pet_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

// UpdatePolicyStatusRequest 更新保单状态请求
type UpdatePolicyStatusRequest struct {
	Status string `json:"status" binding:"required"` // active, expired, cancelled
}

// ClaimStatusRequest 更新理赔状态请求
type ClaimStatusRequest struct {
	Status        string  `json:"status" binding:"required"` // pending, approved, rejected, paid
	ApprovedAmount float64 `json:"approved_amount"`
	Reason        string  `json:"reason"`
}

// VetAppointmentRequest 创建/更新预约请求
type VetAppointmentRequest struct {
	PetID         uint      `json:"pet_id" binding:"required"`
	VetName       string    `json:"vet_name" binding:"required"`
	ClinicName    string    `json:"clinic_name"`
	AppointmentAt time.Time `json:"appointment_at" binding:"required"`
	Reason        string    `json:"reason"`
	Notes         string    `json:"notes"`
	Status        string    `json:"status"`
}

// LostReportRequest 创建/更新寻宠报告请求
type LostReportRequest struct {
	PetID        uint      `json:"pet_id" binding:"required"`
	ReportType   string    `json:"report_type" binding:"required"` // lost/found
	Title        string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	LastSeenAt   time.Time `json:"last_seen_at"`
	LastSeenLat  float64   `json:"last_seen_lat"`
	LastSeenLng  float64   `json:"last_seen_lng"`
	LastSeenAddr string    `json:"last_seen_addr"`
	ContactName  string    `json:"contact_name"`
	ContactPhone string    `json:"contact_phone"`
	Reward       float64   `json:"reward"`
	Photos       []string  `json:"photos"`
	Status       string    `json:"status"`
}

// MapConfigRequest 更新地图配置请求
type MapConfigRequest struct {
	Provider string `json:"provider" binding:"required"`
	APIKey   string `json:"api_key"`
	IsActive bool   `json:"is_active"`
}

// RegisterIntegrationRoutes 注册第三方集成路由
func (c *IntegrationController) RegisterIntegrationRoutes(rg *gin.RouterGroup) {
	// 智能家居
	smarthome := rg.Group("/integrations/smarthome")
	smarthome.GET("/devices", c.ListSmartHomeDevices)
	smarthome.POST("/devices", c.CreateSmartHomeDevice)
	smarthome.DELETE("/devices/:id", c.DeleteSmartHomeDevice)
	smarthome.POST("/devices/:id/control", c.ControlSmartHomeDevice)
	smarthome.GET("/status", c.GetSmartHomeStatus)

	// 宠物医疗 (旧路由 - 兼容)
	medical := rg.Group("/integrations/medical")
	medical.GET("/records", c.ListMedicalRecords)
	medical.POST("/records", c.CreateMedicalRecord)
	medical.GET("/records/:id", c.GetMedicalRecord)
	medical.PUT("/records/:id", c.UpdateMedicalRecord)
	medical.DELETE("/records/:id", c.DeleteMedicalRecord)
	medical.GET("/records/:id/attachments", c.GetMedicalRecordAttachments)

	// 宠物保险 (旧路由 - 兼容)
	insurance := rg.Group("/integrations/insurance")
	insurance.GET("/claims", c.ListInsuranceClaims)
	insurance.POST("/claims", c.CreateInsuranceClaim)
	insurance.GET("/claims/:id", c.GetInsuranceClaim)
	insurance.PUT("/claims/:id", c.UpdateInsuranceClaim)

	// ============ 新版宠物保险 API (/api/v1/insurance) ============
	insuranceV1 := rg.Group("/insurance")
	// 产品
	insuranceV1.GET("/products", c.ListInsuranceProducts)
	insuranceV1.POST("/products", c.CreateInsuranceProduct)
	insuranceV1.GET("/products/:id", c.GetInsuranceProduct)
	insuranceV1.PUT("/products/:id", c.UpdateInsuranceProduct)
	insuranceV1.DELETE("/products/:id", c.DeleteInsuranceProduct)
	// 保单
	insuranceV1.GET("/policies", c.ListInsurancePolicies)
	insuranceV1.POST("/policies", c.PurchaseInsurancePolicy)
	insuranceV1.GET("/policies/:id", c.GetInsurancePolicy)
	insuranceV1.PUT("/policies/:id", c.UpdateInsurancePolicy)
	insuranceV1.PUT("/policies/:id/status", c.UpdatePolicyStatus)
	// 理赔
	insuranceV1.GET("/claims", c.ListInsuranceClaimsV1)
	insuranceV1.POST("/claims", c.CreateInsuranceClaimV1)
	insuranceV1.GET("/claims/:id", c.GetInsuranceClaimV1)
	insuranceV1.PUT("/claims/:id", c.UpdateInsuranceClaimV1)
	insuranceV1.PUT("/claims/:id/status", c.UpdateClaimStatus)
	insuranceV1.POST("/claims/:id/approve", c.ApproveClaim)
	insuranceV1.POST("/claims/:id/reject", c.RejectClaim)

	// ============ 新版宠物医疗 API (/api/v1/medical) ============
	medicalV1 := rg.Group("/medical")
	// 病历
	medicalV1.GET("/records", c.ListMedicalRecordsV1)
	medicalV1.POST("/records", c.CreateMedicalRecordV1)
	medicalV1.GET("/records/:id", c.GetMedicalRecordV1)
	medicalV1.PUT("/records/:id", c.UpdateMedicalRecordV1)
	medicalV1.DELETE("/records/:id", c.DeleteMedicalRecord)
	medicalV1.GET("/records/:id/attachments", c.GetMedicalRecordAttachmentsV1)
	// 预约
	medicalV1.GET("/appointments", c.ListVetAppointments)
	medicalV1.POST("/appointments", c.CreateVetAppointment)
	medicalV1.GET("/appointments/:id", c.GetVetAppointment)
	medicalV1.PUT("/appointments/:id", c.UpdateVetAppointment)
	medicalV1.DELETE("/appointments/:id", c.CancelVetAppointment)

	// 寻回网络
	lostfound := rg.Group("/integrations/lost-found")
	lostfound.GET("/reports", c.ListLostReports)
	lostfound.POST("/reports", c.CreateLostReport)
	lostfound.GET("/reports/:id", c.GetLostReport)
	lostfound.PUT("/reports/:id", c.UpdateLostReport)
	lostfound.GET("/nearby", c.SearchNearbyPets)

	// 地图服务
	maps := rg.Group("/integrations/maps")
	maps.GET("/config", c.GetMapConfig)
	maps.PUT("/config", c.UpdateMapConfig)
}

// ============ 智能家居 ============

// ListSmartHomeDevices 获取智能家居设备列表
func (c *IntegrationController) ListSmartHomeDevices(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var devices []models.SmartHomeDevice
	query := c.DB.Where("user_id = ?", userID)

	if platform := ctx.Query("platform"); platform != "" {
		query = query.Where("platform = ?", platform)
	}
	if deviceType := ctx.Query("device_type"); deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}

	if err := query.Order("updated_at DESC").Find(&devices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": devices})
}

// CreateSmartHomeDevice 添加智能家居设备
func (c *IntegrationController) CreateSmartHomeDevice(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var device models.SmartHomeDevice
	if err := ctx.ShouldBindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}
	device.UserID = userID.(uint)

	if err := c.DB.Create(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建设备失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": device})
}

// DeleteSmartHomeDevice 删除智能家居设备
func (c *IntegrationController) DeleteSmartHomeDevice(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	deviceID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的设备ID"})
		return
	}

	result := c.DB.Where("id = ? AND user_id = ?", deviceID, userID).Delete(&models.SmartHomeDevice{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除设备失败", "error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ControlSmartHomeDevice 控制智能家居设备
func (c *IntegrationController) ControlSmartHomeDevice(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	deviceID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的设备ID"})
		return
	}

	var device models.SmartHomeDevice
	if err := c.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询设备失败", "error": err.Error()})
		}
		return
	}

	var req DeviceControlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 更新设备状态
	now := time.Now()
	device.LastControlAt = &now
	if device.Status == nil {
		device.Status = make(map[string]interface{})
	}
	device.Status["last_action"] = req.Action
	device.Status["last_params"] = req.Params
	device.Status["last_control_at"] = now.Unix()

	if err := c.DB.Save(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "控制设备失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id":     device.ID,
			"action":       req.Action,
			"status":       device.Status,
			"controlled_at": now,
		},
	})
}

// GetSmartHomeStatus 获取智能家居集成状态
func (c *IntegrationController) GetSmartHomeStatus(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var integrations []models.Integration
	c.DB.Where("user_id = ? AND integration_type IN ?", userID,
		[]string{"mi_home", "tmall_genie", "homekit", "google_home"}).Find(&integrations)

	var onlineCount int64
	c.DB.Model(&models.SmartHomeDevice{}).Where("user_id = ? AND is_online = ?", userID, true).Count(&onlineCount)

	var totalCount int64
	c.DB.Model(&models.SmartHomeDevice{}).Where("user_id = ?", userID).Count(&totalCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"integrations": integrations,
			"online_count": onlineCount,
			"total_count":  totalCount,
		},
	})
}

// ============ 宠物医疗 ============

// ListMedicalRecords 获取医疗记录列表
func (c *IntegrationController) ListMedicalRecords(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	query := c.DB.Where("pet_id IN (SELECT id FROM pets WHERE owner_id = ?)", userID)

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if recordType := ctx.Query("record_type"); recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("record_date >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("record_date <= ?", endDate)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.PetMedicalRecord{}).Count(&total)

	var records []models.PetMedicalRecord
	offset := (page - 1) * pageSize
	if err := query.Order("record_date DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"records":  records,
			"total":    total,
			"page":     page,
			"page_size": pageSize,
		},
	})
}

// CreateMedicalRecord 创建医疗记录
func (c *IntegrationController) CreateMedicalRecord(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req MedicalRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证宠物归属
	var pet models.Pet
	if err := c.DB.Where("id = ? AND owner_id = ?", req.PetID, userID).First(&pet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该宠物"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "验证失败", "error": err.Error()})
		}
		return
	}

	record := models.PetMedicalRecord{
		PetID:        req.PetID,
		RecordType:   req.RecordType,
		HospitalName: req.HospitalName,
		DoctorName:   req.DoctorName,
		Diagnosis:    req.Diagnosis,
		Treatment:    req.Treatment,
		Prescription: req.Prescription,
		RecordDate:   req.RecordDate,
		Attachments:  func() string { b, _ := json.Marshal(req.Attachments); return string(b) }(),
		SyncStatus:   "synced",
	}

	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建记录失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": record})
}

// GetMedicalRecord 获取医疗记录详情
func (c *IntegrationController) GetMedicalRecord(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	recordID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的记录ID"})
		return
	}

	var record models.PetMedicalRecord
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		recordID, userID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateMedicalRecord 更新医疗记录
func (c *IntegrationController) UpdateMedicalRecord(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	recordID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的记录ID"})
		return
	}

	var record models.PetMedicalRecord
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		recordID, userID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		}
		return
	}

	var req MedicalRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"record_type":   req.RecordType,
		"hospital_name": req.HospitalName,
		"doctor_name":   req.DoctorName,
		"diagnosis":     req.Diagnosis,
		"treatment":     req.Treatment,
		"prescription":  req.Prescription,
		"record_date":   req.RecordDate,
		"attachments":   req.Attachments,
	}

	if err := c.DB.Model(&record).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// ============ 宠物保险 ============

// ListInsuranceClaims 获取理赔列表
func (c *IntegrationController) ListInsuranceClaims(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	query := c.DB.Where("policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)", userID)

	if policyID := ctx.Query("policy_id"); policyID != "" {
		query = query.Where("policy_id = ?", policyID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.InsuranceClaim{}).Count(&total)

	var claims []models.InsuranceClaim
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&claims).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"claims":    claims,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateInsuranceClaim 创建理赔
func (c *IntegrationController) CreateInsuranceClaim(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req InsuranceClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证保单归属
	var policy models.InsurancePolicy
	if err := c.DB.Where("id = ? AND user_id = ?", req.PolicyID, userID).First(&policy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该保单"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "验证失败", "error": err.Error()})
		}
		return
	}

	// 生成理赔编号
	claimNumber := fmt.Sprintf("CLM%d%d", time.Now().Unix(), time.Now().Nanosecond()%10000)

	claim := models.InsuranceClaim{
		PolicyID:     req.PolicyID,
		ClaimNumber:  claimNumber,
		ClaimType:    req.ClaimType,
		Amount:       req.Amount,
		Reason:       req.Reason,
		Status:       "pending",
		Description:  req.Description,
		Diagnosis:    req.Diagnosis,
		Attachments:  func() string { b, _ := json.Marshal(req.Documents); return string(b) }(),
		SubmittedAt:  time.Now(),
	}

	if err := c.DB.Create(&claim).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建理赔失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": claim})
}

// GetInsuranceClaim 获取理赔详情
func (c *IntegrationController) GetInsuranceClaim(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.Where("id = ? AND policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)",
		claimID, userID).First(&claim).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// UpdateInsuranceClaim 更新理赔
func (c *IntegrationController) UpdateInsuranceClaim(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.Where("id = ? AND policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)",
		claimID, userID).First(&claim).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		}
		return
	}

	var req InsuranceClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"claim_type":  req.ClaimType,
		"amount":      req.Amount,
		"description": req.Description,
		"documents":   req.Documents,
		"claim_date":  req.ClaimDate,
	}

	if err := c.DB.Model(&claim).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// ============ 保险产品 (/api/v1/insurance/products) ============

// ListInsuranceProducts 获取保险产品列表
func (c *IntegrationController) ListInsuranceProducts(ctx *gin.Context) {
	query := c.DB.Where("is_active = ?", true)

	if coverageType := ctx.Query("coverage_type"); coverageType != "" {
		query = query.Where("coverage_type = ?", coverageType)
	}
	if provider := ctx.Query("provider"); provider != "" {
		query = query.Where("provider = ?", provider)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Model(&models.InsuranceProduct{}).Count(&total)

	var products []models.InsuranceProduct
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// CreateInsuranceProduct 创建保险产品
func (c *IntegrationController) CreateInsuranceProduct(ctx *gin.Context) {
	var req InsuranceProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	product := models.InsuranceProduct{
		Name:         req.Name,
		Provider:     req.Provider,
		CoverageType: req.CoverageType,
		MonthlyPrice: req.MonthlyPrice,
		CoverageLimit: req.CoverageLimit,
		Description:  req.Description,
		Terms:        req.Terms,
		IsActive:     req.IsActive,
	}

	if err := c.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": product})
}

// GetInsuranceProduct 获取保险产品详情
func (c *IntegrationController) GetInsuranceProduct(ctx *gin.Context) {
	productID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的产品ID"})
		return
	}

	var product models.InsuranceProduct
	if err := c.DB.First(&product, productID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

// UpdateInsuranceProduct 更新保险产品
func (c *IntegrationController) UpdateInsuranceProduct(ctx *gin.Context) {
	productID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的产品ID"})
		return
	}

	var product models.InsuranceProduct
	if err := c.DB.First(&product, productID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	var req InsuranceProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"name":          req.Name,
		"provider":      req.Provider,
		"coverage_type": req.CoverageType,
		"monthly_price": req.MonthlyPrice,
		"coverage_limit": req.CoverageLimit,
		"description":   req.Description,
		"terms":         req.Terms,
		"is_active":     req.IsActive,
	}

	if err := c.DB.Model(&product).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

// DeleteInsuranceProduct 删除保险产品
func (c *IntegrationController) DeleteInsuranceProduct(ctx *gin.Context) {
	productID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的产品ID"})
		return
	}

	if err := c.DB.Delete(&models.InsuranceProduct{}, productID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ 保险保单 (/api/v1/insurance/policies) ============

// ListInsurancePolicies 获取保单列表
func (c *IntegrationController) ListInsurancePolicies(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	query := c.DB.Where("user_id = ?", userID)

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Model(&models.InsurancePolicy{}).Count(&total)

	var policies []models.InsurancePolicy
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&policies).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"policies":  policies,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	}})
}

// PurchaseInsurancePolicy 购买保单
func (c *IntegrationController) PurchaseInsurancePolicy(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req InsurancePolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证产品存在
	var product models.InsuranceProduct
	if err := c.DB.First(&product, req.ProductID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "保险产品不存在"})
		return
	}

	// 验证宠物归属
	var pet models.Pet
	if err := c.DB.Where("id = ? AND owner_id = ?", req.PetID, userID).First(&pet).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权为该宠物购买保险"})
		return
	}

	// 生成保单号
	policyNumber := fmt.Sprintf("POL%d%d", time.Now().Unix(), time.Now().Nanosecond()%10000)

	// 计算保费
	duration := req.EndDate.Sub(req.StartDate).Hours() / 24 / 30 // 月数
	premium := product.MonthlyPrice * duration
	if premium < product.MonthlyPrice {
		premium = product.MonthlyPrice
	}

	policy := models.InsurancePolicy{
		UserID:         userID.(uint),
		PetID:          req.PetID,
		ProductID:     product.ID,
		PlanName:       product.Name,
		Insurer:        product.Provider,
		PolicyNumber:   policyNumber,
		Premium:        premium,
		CoverageAmount: product.CoverageLimit,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         "active",
	}

	if err := c.DB.Create(&policy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "购买失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": policy})
}

// GetInsurancePolicy 获取保单详情
func (c *IntegrationController) GetInsurancePolicy(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	policyID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的保单ID"})
		return
	}

	var policy models.InsurancePolicy
	if err := c.DB.Where("id = ? AND user_id = ?", policyID, userID).First(&policy).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "保单不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// UpdateInsurancePolicy 更新保单
func (c *IntegrationController) UpdateInsurancePolicy(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	policyID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的保单ID"})
		return
	}

	var policy models.InsurancePolicy
	if err := c.DB.Where("id = ? AND user_id = ?", policyID, userID).First(&policy).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "保单不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// UpdatePolicyStatus 更新保单状态
func (c *IntegrationController) UpdatePolicyStatus(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	policyID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的保单ID"})
		return
	}

	var policy models.InsurancePolicy
	if err := c.DB.Where("id = ? AND user_id = ?", policyID, userID).First(&policy).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "保单不存在"})
		return
	}

	var req UpdatePolicyStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	if err := c.DB.Model(&policy).Update("status", req.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// ============ 理赔 (/api/v1/insurance/claims) - V1 版本 ============

// ListInsuranceClaimsV1 获取理赔列表
func (c *IntegrationController) ListInsuranceClaimsV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	query := c.DB.Where("policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)", userID)

	if policyID := ctx.Query("policy_id"); policyID != "" {
		query = query.Where("policy_id = ?", policyID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Model(&models.InsuranceClaim{}).Count(&total)

	var claims []models.InsuranceClaim
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&claims).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"claims":    claims,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// CreateInsuranceClaimV1 创建理赔
func (c *IntegrationController) CreateInsuranceClaimV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req InsuranceClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证保单归属
	var policy models.InsurancePolicy
	if err := c.DB.Where("id = ? AND user_id = ?", req.PolicyID, userID).First(&policy).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该保单"})
		return
	}

	// 生成理赔编号
	claimNumber := fmt.Sprintf("CLM%d%d", time.Now().Unix(), time.Now().Nanosecond()%10000)

	claim := models.InsuranceClaim{
		PolicyID:     req.PolicyID,
		ClaimNumber:  claimNumber,
		ClaimType:    req.ClaimType,
		Amount:       req.Amount,
		Reason:       req.Reason,
		Description:  req.Description,
		Diagnosis:    req.Diagnosis,
		Status:       "pending",
		Attachments:  func() string { b, _ := json.Marshal(req.Documents); return string(b) }(),
		SubmittedAt:  time.Now(),
	}

	if err := c.DB.Create(&claim).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建理赔失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": claim})
}

// GetInsuranceClaimV1 获取理赔详情
func (c *IntegrationController) GetInsuranceClaimV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.Where("id = ? AND policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)",
		claimID, userID).First(&claim).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// UpdateInsuranceClaimV1 更新理赔
func (c *IntegrationController) UpdateInsuranceClaimV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.Where("id = ? AND policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)",
		claimID, userID).First(&claim).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}

	var req InsuranceClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"claim_type":  req.ClaimType,
		"amount":      req.Amount,
		"reason":      req.Reason,
		"description": req.Description,
		"diagnosis":   req.Diagnosis,
		"documents":   req.Documents,
		"claim_date":  req.ClaimDate,
	}

	if err := c.DB.Model(&claim).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// UpdateClaimStatus 更新理赔状态
func (c *IntegrationController) UpdateClaimStatus(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.Where("id = ? AND policy_id IN (SELECT id FROM insurance_policies WHERE user_id = ?)",
		claimID, userID).First(&claim).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}

	var req ClaimStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"status":          req.Status,
		"approved_amount": req.ApprovedAmount,
	}
	if req.Status == "approved" || req.Status == "rejected" || req.Status == "paid" {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if err := c.DB.Model(&claim).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// ApproveClaim 批准理赔
func (c *IntegrationController) ApproveClaim(ctx *gin.Context) {
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.First(&claim, claimID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}

	var req struct {
		ApprovedAmount float64 `json:"approved_amount"`
	}
	ctx.ShouldBindJSON(&req)

	now := time.Now()
	if err := c.DB.Model(&claim).Updates(map[string]interface{}{
		"status":          "approved",
		"approved_amount": req.ApprovedAmount,
		"processed_at":    &now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批准失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// RejectClaim 拒绝理赔
func (c *IntegrationController) RejectClaim(ctx *gin.Context) {
	claimID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的理赔ID"})
		return
	}

	var claim models.InsuranceClaim
	if err := c.DB.First(&claim, claimID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	ctx.ShouldBindJSON(&req)

	now := time.Now()
	if err := c.DB.Model(&claim).Updates(map[string]interface{}{
		"status":       "rejected",
		"reason":       req.Reason,
		"processed_at": &now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "拒绝失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// ============ 病历 (/api/v1/medical/records) ============

// ListMedicalRecordsV1 获取病历列表
func (c *IntegrationController) ListMedicalRecordsV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	query := c.DB.Where("pet_id IN (SELECT id FROM pets WHERE owner_id = ?)", userID)

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if recordType := ctx.Query("record_type"); recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Model(&models.PetMedicalRecord{}).Count(&total)

	var records []models.PetMedicalRecord
	offset := (page - 1) * pageSize
	if err := query.Order("record_date DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"records":   records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}})
}

// CreateMedicalRecordV1 创建病历
func (c *IntegrationController) CreateMedicalRecordV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req MedicalRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证宠物归属
	var pet models.Pet
	if err := c.DB.Where("id = ? AND owner_id = ?", req.PetID, userID).First(&pet).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作该宠物"})
		return
	}

	record := models.PetMedicalRecord{
		PetID:        req.PetID,
		RecordType:   req.RecordType,
		HospitalName: req.HospitalName,
		DoctorName:   req.DoctorName,
		Diagnosis:    req.Diagnosis,
		Treatment:    req.Treatment,
		Prescription: req.Prescription,
		RecordDate:   req.RecordDate,
		Attachments:  func() string { b, _ := json.Marshal(req.Attachments); return string(b) }(),
		SyncStatus:   "synced",
	}

	if err := c.DB.Create(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": record})
}

// GetMedicalRecordV1 获取病历详情
func (c *IntegrationController) GetMedicalRecordV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	recordID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的记录ID"})
		return
	}

	var record models.PetMedicalRecord
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		recordID, userID).First(&record).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// UpdateMedicalRecordV1 更新病历
func (c *IntegrationController) UpdateMedicalRecordV1(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	recordID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的记录ID"})
		return
	}

	var record models.PetMedicalRecord
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		recordID, userID).First(&record).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	var req MedicalRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"record_type":   req.RecordType,
		"hospital_name": req.HospitalName,
		"doctor_name":   req.DoctorName,
		"diagnosis":     req.Diagnosis,
		"treatment":     req.Treatment,
		"prescription":  req.Prescription,
		"record_date":   req.RecordDate,
		"attachments":   req.Attachments,
	}

	if err := c.DB.Model(&record).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record})
}

// DeleteMedicalRecord 删除病历
func (c *IntegrationController) DeleteMedicalRecord(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	recordID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的记录ID"})
		return
	}

	var record models.PetMedicalRecord
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		recordID, userID).First(&record).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}

	if err := c.DB.Delete(&record).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetMedicalRecordAttachments 获取病历附件
func (c *IntegrationController) GetMedicalRecordAttachments(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	recordID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的记录ID"})
		return
	}

	var record models.PetMedicalRecord
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		recordID, userID).First(&record).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记录不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": record.Attachments})
}

// GetMedicalRecordAttachmentsV1 获取病历附件 (V1)
func (c *IntegrationController) GetMedicalRecordAttachmentsV1(ctx *gin.Context) {
	c.GetMedicalRecordAttachments(ctx)
}

// ============ 预约兽医 (/api/v1/medical/appointments) ============

// ListVetAppointments 获取预约列表
func (c *IntegrationController) ListVetAppointments(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	query := c.DB.Where("pet_id IN (SELECT id FROM pets WHERE owner_id = ?)", userID)

	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var total int64
	query.Model(&models.VetAppointment{}).Count(&total)

	var appointments []models.VetAppointment
	offset := (page - 1) * pageSize
	if err := query.Order("appointment_at DESC").Offset(offset).Limit(pageSize).Find(&appointments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"appointments": appointments,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	}})
}

// CreateVetAppointment 创建预约
func (c *IntegrationController) CreateVetAppointment(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req VetAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	// 验证宠物归属
	var pet models.Pet
	if err := c.DB.Where("id = ? AND owner_id = ?", req.PetID, userID).First(&pet).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权为该宠物预约"})
		return
	}

	status := req.Status
	if status == "" {
		status = "scheduled"
	}

	appointment := models.VetAppointment{
		PetID:         req.PetID,
		VetName:       req.VetName,
		ClinicName:    req.ClinicName,
		AppointmentAt: req.AppointmentAt,
		Reason:        req.Reason,
		Notes:         req.Notes,
		Status:        status,
	}

	if err := c.DB.Create(&appointment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建预约失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": appointment})
}

// GetVetAppointment 获取预约详情
func (c *IntegrationController) GetVetAppointment(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	appointmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的预约ID"})
		return
	}

	var appointment models.VetAppointment
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		appointmentID, userID).First(&appointment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": appointment})
}

// UpdateVetAppointment 更新预约
func (c *IntegrationController) UpdateVetAppointment(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	appointmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的预约ID"})
		return
	}

	var appointment models.VetAppointment
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		appointmentID, userID).First(&appointment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
		return
	}

	var req VetAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"vet_name":       req.VetName,
		"clinic_name":    req.ClinicName,
		"appointment_at": req.AppointmentAt,
		"reason":         req.Reason,
		"notes":          req.Notes,
		"status":         req.Status,
	}

	if err := c.DB.Model(&appointment).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": appointment})
}

// CancelVetAppointment 取消预约
func (c *IntegrationController) CancelVetAppointment(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	appointmentID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的预约ID"})
		return
	}

	var appointment models.VetAppointment
	if err := c.DB.Where("id = ? AND pet_id IN (SELECT id FROM pets WHERE owner_id = ?)",
		appointmentID, userID).First(&appointment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "预约不存在"})
		return
	}

	if err := c.DB.Model(&appointment).Update("status", "cancelled").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": appointment})
}

// ListLostReports 获取寻宠报告列表
func (c *IntegrationController) ListLostReports(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	query := c.DB.Where("user_id = ?", userID)

	if reportType := ctx.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.PetLostReport{}).Count(&total)

	var reports []models.PetLostReport
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"reports":   reports,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateLostReport 创建寻宠报告
func (c *IntegrationController) CreateLostReport(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req LostReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	report := models.PetLostReport{
		PetID:        req.PetID,
		UserID:       userID.(uint),
		ReportType:   req.ReportType,
		Title:        req.Title,
		Description:  req.Description,
		LastSeenAt:   &req.LastSeenAt,
		LastSeenLat:  req.LastSeenLat,
		LastSeenLng:  req.LastSeenLng,
		LastSeenAddr: req.LastSeenAddr,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Reward:       req.Reward,
		Photos:       req.Photos,
		Status:       "active",
	}

	if err := c.DB.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建报告失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": report})
}

// GetLostReport 获取寻宠报告详情
func (c *IntegrationController) GetLostReport(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	reportID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的报告ID"})
		return
	}

	var report models.PetLostReport
	if err := c.DB.Where("id = ? AND user_id = ?", reportID, userID).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// UpdateLostReport 更新寻宠报告
func (c *IntegrationController) UpdateLostReport(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	reportID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的报告ID"})
		return
	}

	var report models.PetLostReport
	if err := c.DB.Where("id = ? AND user_id = ?", reportID, userID).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		}
		return
	}

	var req LostReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"report_type":    req.ReportType,
		"title":          req.Title,
		"description":    req.Description,
		"last_seen_at":   req.LastSeenAt,
		"last_seen_lat":  req.LastSeenLat,
		"last_seen_lng":  req.LastSeenLng,
		"last_seen_addr": req.LastSeenAddr,
		"contact_name":   req.ContactName,
		"contact_phone":  req.ContactPhone,
		"reward":         req.Reward,
		"photos":         req.Photos,
	}
	if req.Status != "" {
		updates["status"] = req.Status
		if req.Status == "found" || req.Status == "closed" {
			now := time.Now()
			updates["resolved_at"] = &now
		}
	}

	if err := c.DB.Model(&report).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// SearchNearbyPets 搜索附近宠物
func (c *IntegrationController) SearchNearbyPets(ctx *gin.Context) {
	lat, _ := strconv.ParseFloat(ctx.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(ctx.Query("lng"), 64)
	radiusKm, _ := strconv.ParseFloat(ctx.DefaultQuery("radius_km", "10"), 64)

	if lat == 0 && lng == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "lat 和 lng 参数必填"})
		return
	}

	// 使用矩形过滤（简化版），实际距离计算在应用层
	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / (111.0 * 0.7)

	minLat := lat - latDelta
	maxLat := lat + latDelta
	minLng := lng - lngDelta
	maxLng := lng + lngDelta

	var reports []models.PetLostReport
	if err := c.DB.Where("status = 'active'").
		Where("last_seen_lat BETWEEN ? AND ? AND last_seen_lng BETWEEN ? AND ?",
			minLat, maxLat, minLng, maxLng).
		Order("created_at DESC").
		Limit(50).
		Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "搜索失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"reports":    reports,
			"total":      len(reports),
			"center_lat": lat,
			"center_lng": lng,
			"radius_km":  radiusKm,
		},
	})
}

// ============ 地图服务 ============

// GetMapConfig 获取地图配置
func (c *IntegrationController) GetMapConfig(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var config models.ThirdPartyMapConfig
	if err := c.DB.Where("user_id = ?", userID).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"provider":    "",
					"is_active":   false,
					"quota_used":  0,
					"quota_limit": 10000,
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
		return
	}

	// 不返回 API Key 明文
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":          config.ID,
			"provider":    config.Provider,
			"api_key_set": config.APIKey != "",
			"is_active":   config.IsActive,
			"quota_used":  config.QuotaUsed,
			"quota_limit": config.QuotaLimit,
		},
	})
}

// UpdateMapConfig 更新地图配置
func (c *IntegrationController) UpdateMapConfig(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req MapConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}

	var config models.ThirdPartyMapConfig
	if err := c.DB.Where("user_id = ?", userID).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新配置
			config = models.ThirdPartyMapConfig{
				UserID:     userID.(uint),
				Provider:   req.Provider,
				APIKey:     req.APIKey,
				IsActive:   req.IsActive,
				QuotaLimit: 10000,
			}
			if err := c.DB.Create(&config).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建配置失败", "error": err.Error()})
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败", "error": err.Error()})
			return
		}
	} else {
		// 更新现有配置
		updates := map[string]interface{}{
			"provider":  req.Provider,
			"is_active": req.IsActive,
		}
		if req.APIKey != "" {
			updates["api_key"] = req.APIKey
		}
		if err := c.DB.Model(&config).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新配置失败", "error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":          config.ID,
			"provider":    config.Provider,
			"api_key_set": config.APIKey != "",
			"is_active":   config.IsActive,
			"quota_used":  config.QuotaUsed,
			"quota_limit": config.QuotaLimit,
		},
	})
}
