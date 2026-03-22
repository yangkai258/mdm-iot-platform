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

// InsuranceController 保险理赔控制器
type InsuranceController struct {
	DB *gorm.DB
}

// RegisterInsuranceRoutes 注册保险相关路由
func (ic *InsuranceController) RegisterInsuranceRoutes(r *gin.RouterGroup) {
	// 保险产品
	r.GET("/insurance/products", ic.ListProducts)
	r.POST("/insurance/products", ic.CreateProduct)
	r.GET("/insurance/products/:id", ic.GetProduct)

	// 理赔申请
	r.GET("/insurance/claims", ic.ListClaims)
	r.POST("/insurance/claims", ic.CreateClaim)
	r.GET("/insurance/claims/:id", ic.GetClaim)
	r.PUT("/insurance/claims/:id/status", ic.UpdateClaimStatus)
	r.POST("/insurance/claims/:id/documents", ic.UploadClaimDocument)

	// 宠物健康档案
	r.GET("/insurance/pets/:pet_id/health-records", ic.ListHealthRecords)
	r.POST("/insurance/pets/:pet_id/health-records", ic.CreateHealthRecord)
}

// ============ 保险产品 API ============

// ListProducts 获取保险产品列表
func (ic *InsuranceController) ListProducts(c *gin.Context) {
	tenantID := getTenantID(c)

	var products []models.InsuranceProduct
	query := ic.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	// 过滤：保障类型
	if coverageType := c.Query("coverage_type"); coverageType != "" {
		query = query.Where("coverage_type = ?", coverageType)
	}
	// 过滤：保险公司
	if provider := c.Query("provider"); provider != "" {
		query = query.Where("provider ILIKE ?", "%"+provider+"%")
	}
	// 过滤：物种
	if species := c.Query("species"); species != "" {
		query = query.Where("? = ANY(species_allowed)", species)
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
	offset := (page - 1) * pageSize

	var total int64
	ic.DB.Model(&models.InsuranceProduct{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&total)

	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      products,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateProduct 添加保险产品
func (ic *InsuranceController) CreateProduct(c *gin.Context) {
	tenantID := getTenantID(c)

	var input struct {
		Name           string   `json:"name" binding:"required"`
		Code           string   `json:"code" binding:"required"`
		CoverageType   string   `json:"coverage_type" binding:"required"`
		Provider       string   `json:"provider"`
		CoverageAmount float64  `json:"coverage_amount" binding:"required"`
		Premium        float64  `json:"premium" binding:"required"`
		PremiumPeriod  string   `json:"premium_period"`
		Deductible     float64  `json:"deductible"`
		WaitPeriodDays int      `json:"wait_period_days"`
		CoverAgeMin    int      `json:"cover_age_min"`
		CoverAgeMax    int      `json:"cover_age_max"`
		BreedCodes     []string `json:"breed_codes"`
		SpeciesAllowed []string `json:"species_allowed"`
		Description    string   `json:"description"`
		Terms          string   `json:"terms"`
		Exclusions     string   `json:"exclusions"`
		CoverageItems  []string `json:"coverage_items"`
		MaxClaimAmount float64  `json:"max_claim_amount"`
		AnnualMaxClaim float64  `json:"annual_max_claim"`
		SortOrder      int      `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误: " + err.Error()})
		return
	}

	product := models.InsuranceProduct{
		Name:           input.Name,
		Code:           input.Code,
		CoverageType:   input.CoverageType,
		Provider:       input.Provider,
		CoverageAmount: input.CoverageAmount,
		Premium:        input.Premium,
		PremiumPeriod:  input.PremiumPeriod,
		Deductible:     input.Deductible,
		WaitPeriodDays: input.WaitPeriodDays,
		CoverAgeMin:    input.CoverAgeMin,
		CoverAgeMax:    input.CoverAgeMax,
		BreedCodes:     input.BreedCodes,
		SpeciesAllowed: input.SpeciesAllowed,
		Description:    input.Description,
		Terms:          input.Terms,
		Exclusions:     input.Exclusions,
		CoverageItems:  input.CoverageItems,
		MaxClaimAmount: input.MaxClaimAmount,
		AnnualMaxClaim: input.AnnualMaxClaim,
		SortOrder:      input.SortOrder,
		TenantID:       tenantID,
	}

	if err := ic.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "data": product})
}

// GetProduct 获取产品详情
func (ic *InsuranceController) GetProduct(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	var product models.InsuranceProduct
	query := ic.DB.Where("tenant_id = ?", tenantID)

	// 支持 UUID 或数字 ID
	if len(id) == 36 {
		query = query.Where("product_uuid = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.First(&product).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 4040, "message": "产品不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": product})
}

// ============ 理赔申请 API ============

// ListClaims 获取理赔列表
func (ic *InsuranceController) ListClaims(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)

	var claims []models.InsuranceClaim
	query := ic.DB.Where("tenant_id = ?", tenantID)

	// 权限过滤：普通用户只能看自己的理赔
	if role, _ := c.Get("role"); role != "admin" && role != "super_admin" {
		query = query.Where("owner_id = ?", userID)
	}

	// 过滤：状态
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// 过滤：宠物UUID
	if petUUID := c.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	// 过滤：产品UUID
	if productUUID := c.Query("product_uuid"); productUUID != "" {
		query = query.Where("product_uuid = ?", productUUID)
	}
	// 过滤：出险日期范围
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("incident_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("incident_date <= ?", endDate)
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
	offset := (page - 1) * pageSize

	var total int64
	ic.DB.Model(&models.InsuranceClaim{}).Where("tenant_id = ?", tenantID).Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&claims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      claims,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateClaim 创建理赔
func (ic *InsuranceController) CreateClaim(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)

	var input struct {
		ProductUUID  string  `json:"product_uuid" binding:"required"`
		PetUUID      string  `json:"pet_uuid" binding:"required"`
		IncidentDate string  `json:"incident_date" binding:"required"`
		IncidentType string  `json:"incident_type" binding:"required"`
		IncidentDesc string  `json:"incident_desc"`
		HospitalName string  `json:"hospital_name"`
		Diagnosis    string  `json:"diagnosis"`
		ClaimAmount  float64 `json:"claim_amount" binding:"required"`
		PolicyNo     string  `json:"policy_no"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误: " + err.Error()})
		return
	}

	incidentDate, err := time.Parse("2006-01-02", input.IncidentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "出险日期格式错误，请使用 YYYY-MM-DD"})
		return
	}

	// 验证产品存在
	var product models.InsuranceProduct
	if err := ic.DB.Where("product_uuid = ? AND tenant_id = ? AND is_active = ?", input.ProductUUID, tenantID, true).First(&product).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 4040, "message": "保险产品不存在"})
		return
	}

	claim := models.InsuranceClaim{
		ClaimNo:      generateClaimNo(),
		ProductUUID:  input.ProductUUID,
		PetUUID:      input.PetUUID,
		OwnerID:      userID,
		IncidentDate: incidentDate,
		IncidentType: input.IncidentType,
		IncidentDesc: input.IncidentDesc,
		HospitalName: input.HospitalName,
		Diagnosis:    input.Diagnosis,
		ClaimAmount:  input.ClaimAmount,
		PolicyNo:     input.PolicyNo,
		Status:       "draft",
		TenantID:     tenantID,
	}

	if err := ic.DB.Create(&claim).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建理赔失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "data": claim})
}

// GetClaim 获取理赔详情
func (ic *InsuranceController) GetClaim(c *gin.Context) {
	tenantID := getTenantID(c)
	userID := getUserID(c)
	id := c.Param("id")

	var claim models.InsuranceClaim
	query := ic.DB.Where("tenant_id = ?", tenantID)

	// 支持 UUID 或数字 ID
	if len(id) == 36 {
		query = query.Where("claim_uuid = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	// 权限过滤：普通用户只能看自己的
	if role, _ := c.Get("role"); role != "admin" && role != "super_admin" {
		query = query.Where("owner_id = ?", userID)
	}

	if err := query.Preload("ClaimDocuments").First(&claim).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 4040, "message": "理赔记录不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": claim})
}

// UpdateClaimStatus 更新理赔状态
func (ic *InsuranceController) UpdateClaimStatus(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	var input struct {
		Status   string `json:"status" binding:"required"`
		ApprovedAmount float64 `json:"approved_amount"`
		RejectionReason string `json:"rejection_reason"`
		ReviewNotes string `json:"review_notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误: " + err.Error()})
		return
	}

	// 验证状态转换合法
	validTransitions := map[string][]string{
		"draft":      {"submitted"},
		"submitted":  {"under_review", "rejected"},
		"under_review": {"approved", "rejected"},
		"approved":  {"paid", "closed"},
		"paid":       {"closed"},
	}

	var claim models.InsuranceClaim
	query := ic.DB.Where("tenant_id = ?", tenantID)
	if len(id) == 36 {
		query = query.Where("claim_uuid = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.First(&claim).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 4040, "message": "理赔记录不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	allowed, ok := validTransitions[claim.Status]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": fmt.Sprintf("当前状态 %s 不允许更新", claim.Status)})
		return
	}

	validNext := false
	for _, s := range allowed {
		if s == input.Status {
			validNext = true
			break
		}
	}
	if !validNext {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": fmt.Sprintf("状态 %s 不能转换为 %s", claim.Status, input.Status)})
		return
	}

	updates := map[string]interface{}{
		"status": input.Status,
	}
	if input.ReviewNotes != "" {
		updates["review_notes"] = input.ReviewNotes
	}
	if input.ApprovedAmount > 0 {
		updates["approved_amount"] = input.ApprovedAmount
	}
	if input.RejectionReason != "" {
		updates["rejection_reason"] = input.RejectionReason
	}
	if input.Status == "paid" {
		now := time.Now()
		updates["paid_at"] = &now
	}

	if err := ic.DB.Model(&claim).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新状态失败"})
		return
	}

	ic.DB.Where("claim_uuid = ?", claim.ClaimUUID).First(&claim)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": claim})
}

// UploadClaimDocument 上传理赔文档
func (ic *InsuranceController) UploadClaimDocument(c *gin.Context) {
	tenantID := getTenantID(c)
	id := c.Param("id")

	// 查找理赔记录
	var claim models.InsuranceClaim
	query := ic.DB.Where("tenant_id = ?", tenantID)
	if len(id) == 36 {
		query = query.Where("claim_uuid = ?", id)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.First(&claim).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 4040, "message": "理赔记录不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	// 限制只能为草稿或已提交状态的理赔添加文档
	if claim.Status != "draft" && claim.Status != "submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "只能在草稿或已提交状态下上传文档"})
		return
	}

	var input struct {
		DocType     string `json:"doc_type" binding:"required"`
		FileName    string `json:"file_name" binding:"required"`
		FileURL     string `json:"file_url" binding:"required"`
		FileSize    int64  `json:"file_size"`
		MimeType    string `json:"mime_type"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误: " + err.Error()})
		return
	}

	doc := models.InsuranceClaimDocument{
		ClaimUUID:   claim.ClaimUUID,
		DocType:     input.DocType,
		FileName:    input.FileName,
		FileURL:     input.FileURL,
		FileSize:    input.FileSize,
		MimeType:    input.MimeType,
		Description: input.Description,
		TenantID:    tenantID,
	}

	if err := ic.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "上传文档失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "data": doc})
}

// ============ 宠物健康档案 API ============

// ListHealthRecords 获取健康档案列表
func (ic *InsuranceController) ListHealthRecords(c *gin.Context) {
	tenantID := getTenantID(c)
	petID := c.Param("pet_id")

	var records []models.PetHealthRecord
	query := ic.DB.Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID)

	// 过滤：记录类型
	if recordType := c.Query("record_type"); recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}
	// 过滤：日期范围
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("record_date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("record_date <= ?", endDate)
	}
	// 过滤：是否关联保险
	if isInsured := c.Query("is_insured"); isInsured != "" {
		query = query.Where("is_insured = ?", isInsured == "true")
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
	offset := (page - 1) * pageSize

	var total int64
	ic.DB.Model(&models.PetHealthRecord{}).Where("pet_uuid = ? AND tenant_id = ?", petID, tenantID).Count(&total)

	if err := query.Order("record_date DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"items":      records,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateHealthRecord 添加健康记录
func (ic *InsuranceController) CreateHealthRecord(c *gin.Context) {
	tenantID := getTenantID(c)
	petID := c.Param("pet_id")

	var input struct {
		RecordType   string    `json:"record_type" binding:"required"`
		RecordDate   string    `json:"record_date" binding:"required"`
		Title        string    `json:"title" binding:"required"`
		Hospital     string    `json:"hospital"`
		Doctor       string    `json:"doctor"`
		VetName      string    `json:"vet_name"`
		Diagnosis    string    `json:"diagnosis"`
		Treatment    string    `json:"treatment"`
		Prescription string    `json:"prescription"`
		Medications  []string  `json:"medications"`
		Cost         float64   `json:"cost"`
		VaccineName  string    `json:"vaccine_name"`
		NextDueDate  string    `json:"next_due_date"`
		Weight       float64   `json:"weight"`
		Notes        string    `json:"notes"`
		Attachments  []string  `json:"attachments"`
		IsInsured    bool      `json:"is_insured"`
		InsuranceClaimUUID string `json:"insurance_claim_uuid"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数错误: " + err.Error()})
		return
	}

	recordDate, err := time.Parse("2006-01-02", input.RecordDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "记录日期格式错误，请使用 YYYY-MM-DD"})
		return
	}

	var nextDueDate *time.Time
	if input.NextDueDate != "" {
		t, err := time.Parse("2006-01-02", input.NextDueDate)
		if err == nil {
			nextDueDate = &t
		}
	}

	record := models.PetHealthRecord{
		PetUUID:              petID,
		RecordType:           input.RecordType,
		RecordDate:           recordDate,
		Title:                input.Title,
		Hospital:             input.Hospital,
		Doctor:               input.Doctor,
		VetName:              input.VetName,
		Diagnosis:            input.Diagnosis,
		Treatment:            input.Treatment,
		Prescription:         input.Prescription,
		Medications:          input.Medications,
		Cost:                 input.Cost,
		VaccineName:          input.VaccineName,
		NextDueDate:          nextDueDate,
		Weight:               input.Weight,
		Notes:                input.Notes,
		Attachments:         input.Attachments,
		IsInsured:           input.IsInsured,
		InsuranceClaimUUID:   input.InsuranceClaimUUID,
		TenantID:             tenantID,
	}

	if err := ic.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建健康记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "data": record})
}

// generateClaimNo 理赔单号生成
func generateClaimNo() string {
	now := time.Now()
	return fmt.Sprintf("CL%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
}
