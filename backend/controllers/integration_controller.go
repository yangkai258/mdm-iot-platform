package controllers

import (
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
	Description string    `json:"description"`
	Documents   []string  `json:"documents"`
	ClaimDate   time.Time `json:"claim_date"`
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

	// 宠物医疗
	medical := rg.Group("/integrations/medical")
	medical.GET("/records", c.ListMedicalRecords)
	medical.POST("/records", c.CreateMedicalRecord)
	medical.GET("/records/:id", c.GetMedicalRecord)
	medical.PUT("/records/:id", c.UpdateMedicalRecord)

	// 宠物保险
	insurance := rg.Group("/integrations/insurance")
	insurance.GET("/claims", c.ListInsuranceClaims)
	insurance.POST("/claims", c.CreateInsuranceClaim)
	insurance.GET("/claims/:id", c.GetInsuranceClaim)
	insurance.PUT("/claims/:id", c.UpdateInsuranceClaim)

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
		Attachments:  req.Attachments,
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
		PolicyID:    req.PolicyID,
		ClaimNumber: claimNumber,
		ClaimType:   req.ClaimType,
		Amount:      req.Amount,
		Status:      "pending",
		Description: req.Description,
		Documents:   req.Documents,
		ClaimDate:   &req.ClaimDate,
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

// ============ 寻回网络 ============

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
