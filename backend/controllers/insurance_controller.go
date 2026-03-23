package controllers

import (
	"fmt"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InsuranceController struct {
	DB *gorm.DB
}

func (ctrl *InsuranceController) RegisterRoutes(rg *gin.RouterGroup) {
	insurance := rg.Group("/insurance")
	{
		insurance.GET("/products", ctrl.GetProducts)
		insurance.POST("/products", ctrl.CreateProduct)
		insurance.GET("/products/:id", ctrl.GetProduct)
		insurance.PUT("/products/:id", ctrl.UpdateProduct)
		insurance.DELETE("/products/:id", ctrl.DeleteProduct)

		insurance.GET("/policies", ctrl.GetPolicies)
		insurance.POST("/policies", ctrl.CreatePolicy)
		insurance.GET("/policies/:id", ctrl.GetPolicy)

		insurance.GET("/claims", ctrl.GetClaims)
		insurance.POST("/claims", ctrl.CreateClaim)
		insurance.GET("/claims/:id", ctrl.GetClaim)
		insurance.PUT("/claims/:id/status", ctrl.UpdateClaimStatus)
	}

	medical := rg.Group("/medical")
	{
		medical.GET("/records", ctrl.GetMedicalRecords)
		medical.POST("/records", ctrl.CreateMedicalRecord)
		medical.GET("/records/:id", ctrl.GetMedicalRecord)
		medical.PUT("/records/:id", ctrl.UpdateMedicalRecord)
		medical.DELETE("/records/:id", ctrl.DeleteMedicalRecord)

		medical.GET("/appointments", ctrl.GetAppointments)
		medical.POST("/appointments", ctrl.CreateAppointment)
		medical.GET("/appointments/:id", ctrl.GetAppointment)
		medical.PUT("/appointments/:id", ctrl.UpdateAppointment)
		medical.DELETE("/appointments/:id", ctrl.DeleteAppointment)
	}
}

func (ctrl *InsuranceController) GetProducts(c *gin.Context) {
	var products []models.InsuranceProduct
	query := ctrl.DB.Model(&models.InsuranceProduct{})

	if c.Query("active") == "true" {
		query = query.Where("is_active = ?", true)
	}

	query.Find(&products)
	c.JSON(http.StatusOK, products)
}

func (ctrl *InsuranceController) CreateProduct(c *gin.Context) {
	var product models.InsuranceProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Create(&product)
	c.JSON(http.StatusOK, product)
}

func (ctrl *InsuranceController) GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.InsuranceProduct

	if err := ctrl.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ctrl *InsuranceController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.InsuranceProduct

	if err := ctrl.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req struct {
		Name         string  `json:"name"`
		Provider     string  `json:"provider"`
		CoverageType string  `json:"coverage_type"`
		MonthlyPrice float64 `json:"monthly_price"`
		Description string  `json:"description"`
		IsActive     *bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&product).Updates(req)
	c.JSON(http.StatusOK, product)
}

func (ctrl *InsuranceController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.InsuranceProduct{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *InsuranceController) GetPolicies(c *gin.Context) {
	petID := c.Query("pet_id")
	var policies []models.InsurancePolicy
	query := ctrl.DB.Model(&models.InsurancePolicy{})

	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}

	query.Find(&policies)
	c.JSON(http.StatusOK, policies)
}

func (ctrl *InsuranceController) CreatePolicy(c *gin.Context) {
	var req struct {
		PetID       uint      `json:"pet_id" binding:"required"`
		ProductID   uint      `json:"product_id" binding:"required"`
		Premium     float64   `json:"premium"`
		StartDate   time.Time `json:"start_date"`
		EndDate     time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	policy := models.InsurancePolicy{
		PetID:       req.PetID,
		ProductID:   req.ProductID,
		PolicyNumber: fmt.Sprintf("POL-%d-%d", req.PetID, time.Now().Unix()),
		Status:      "active",
		Premium:     req.Premium,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	ctrl.DB.Create(&policy)
	c.JSON(http.StatusOK, policy)
}

func (ctrl *InsuranceController) GetPolicy(c *gin.Context) {
	id := c.Param("id")
	var policy models.InsurancePolicy

	if err := ctrl.DB.First(&policy, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
		return
	}

	c.JSON(http.StatusOK, policy)
}

func (ctrl *InsuranceController) GetClaims(c *gin.Context) {
	policyID := c.Query("policy_id")
	var claims []models.InsuranceClaim
	query := ctrl.DB.Model(&models.InsuranceClaim{})

	if policyID != "" {
		query = query.Where("policy_id = ?", policyID)
	}

	query.Find(&claims)
	c.JSON(http.StatusOK, claims)
}

func (ctrl *InsuranceController) CreateClaim(c *gin.Context) {
	var req struct {
		PolicyID    uint    `json:"policy_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		Reason     string  `json:"reason" binding:"required"`
		Description string  `json:"description"`
		Diagnosis  string  `json:"diagnosis"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claim := models.InsuranceClaim{
		PolicyID:    req.PolicyID,
		ClaimNumber: fmt.Sprintf("CLM-%d-%d", req.PolicyID, time.Now().Unix()),
		Amount:      req.Amount,
		Reason:      req.Reason,
		Description: req.Description,
		Diagnosis:   req.Diagnosis,
		Status:      "pending",
		SubmittedAt: time.Now(),
	}

	ctrl.DB.Create(&claim)
	c.JSON(http.StatusOK, claim)
}

func (ctrl *InsuranceController) GetClaim(c *gin.Context) {
	id := c.Param("id")
	var claim models.InsuranceClaim

	if err := ctrl.DB.First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	c.JSON(http.StatusOK, claim)
}

func (ctrl *InsuranceController) UpdateClaimStatus(c *gin.Context) {
	id := c.Param("id")
	var claim models.InsuranceClaim

	if err := ctrl.DB.First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == "approved" || req.Status == "rejected" {
		updates["processed_at"] = &now
	}

	ctrl.DB.Model(&claim).Updates(updates)
	c.JSON(http.StatusOK, claim)
}

func (ctrl *InsuranceController) GetMedicalRecords(c *gin.Context) {
	petID := c.Query("pet_id")
	recordType := c.Query("type")
	var records []models.PetMedicalRecord
	query := ctrl.DB.Model(&models.PetMedicalRecord{})

	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	query.Order("record_date DESC").Find(&records)
	c.JSON(http.StatusOK, records)
}

func (ctrl *InsuranceController) CreateMedicalRecord(c *gin.Context) {
	var record models.PetMedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Create(&record)
	c.JSON(http.StatusOK, record)
}

func (ctrl *InsuranceController) GetMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.PetMedicalRecord

	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (ctrl *InsuranceController) UpdateMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.PetMedicalRecord

	if err := ctrl.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&record).Updates(req)
	c.JSON(http.StatusOK, record)
}

func (ctrl *InsuranceController) DeleteMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.PetMedicalRecord{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *InsuranceController) GetAppointments(c *gin.Context) {
	petID := c.Query("pet_id")
	status := c.Query("status")
	var appointments []models.VetAppointment
	query := ctrl.DB.Model(&models.VetAppointment{})

	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("appointment_at ASC").Find(&appointments)
	c.JSON(http.StatusOK, appointments)
}

func (ctrl *InsuranceController) CreateAppointment(c *gin.Context) {
	var apt models.VetAppointment
	if err := c.ShouldBindJSON(&apt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apt.Status = "scheduled"
	ctrl.DB.Create(&apt)
	c.JSON(http.StatusOK, apt)
}

func (ctrl *InsuranceController) GetAppointment(c *gin.Context) {
	id := c.Param("id")
	var apt models.VetAppointment

	if err := ctrl.DB.First(&apt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, apt)
}

func (ctrl *InsuranceController) UpdateAppointment(c *gin.Context) {
	id := c.Param("id")
	var apt models.VetAppointment

	if err := ctrl.DB.First(&apt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&apt).Updates(req)
	c.JSON(http.StatusOK, apt)
}

func (ctrl *InsuranceController) DeleteAppointment(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.VetAppointment{}).Where("id = ?", id).Update("status", "cancelled")
	c.JSON(http.StatusOK, gin.H{"message": "cancelled"})
}
