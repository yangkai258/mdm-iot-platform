package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InsuranceController 保险管理控制器
type InsuranceController struct {
	DB *gorm.DB
}

// ============ InsuranceProduct 产品管理 ============

// InsuranceProductList 产品列表
func (c *InsuranceController) InsuranceProductList(ctx *gin.Context) {
	var products []models.InsuranceProduct
	var total int64

	query := c.DB.Model(&models.InsuranceProduct{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if coverageType := ctx.Query("coverage_type"); coverageType != "" {
		query = query.Where("coverage_type = ?", coverageType)
	}
	if isActive := ctx.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sort_order ASC, id DESC").Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": products, "total": total, "page": page, "page_size": pageSize,
	}})
}

// InsuranceProductGet 获取产品
func (c *InsuranceController) InsuranceProductGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.InsuranceProduct
	if err := c.DB.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

// InsuranceProductCreate 创建产品
func (c *InsuranceController) InsuranceProductCreate(ctx *gin.Context) {
	var product models.InsuranceProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}

// InsuranceProductUpdate 更新产品
func (c *InsuranceController) InsuranceProductUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.InsuranceProduct
	if err := c.DB.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}
	var updateData models.InsuranceProduct
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

// InsuranceProductDelete 删除产品
func (c *InsuranceController) InsuranceProductDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.InsuranceProduct{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ InsuranceClaim 理赔管理 ============

// InsuranceClaimList 理赔列表
func (c *InsuranceController) InsuranceClaimList(ctx *gin.Context) {
	var claims []models.InsuranceClaim
	var total int64

	query := c.DB.Model(&models.InsuranceClaim{})

	if petUUID := ctx.Query("pet_uuid"); petUUID != "" {
		query = query.Where("pet_uuid = ?", petUUID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if productUUID := ctx.Query("product_uuid"); productUUID != "" {
		query = query.Where("product_uuid = ?", productUUID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Preload("Documents").Offset(offset).Limit(pageSize).Order("id DESC").Find(&claims).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": claims, "total": total, "page": page, "page_size": pageSize,
	}})
}

// InsuranceClaimGet 获取理赔
func (c *InsuranceController) InsuranceClaimGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var claim models.InsuranceClaim
	if err := c.DB.Preload("Documents").First(&claim, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// InsuranceClaimCreate 创建理赔
func (c *InsuranceController) InsuranceClaimCreate(ctx *gin.Context) {
	var claim models.InsuranceClaim
	if err := ctx.ShouldBindJSON(&claim); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&claim).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// InsuranceClaimUpdate 更新理赔
func (c *InsuranceController) InsuranceClaimUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var claim models.InsuranceClaim
	if err := c.DB.First(&claim, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}
	var updateData models.InsuranceClaim
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

// InsuranceClaimDelete 删除理赔
func (c *InsuranceController) InsuranceClaimDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.InsuranceClaim{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// InsuranceClaimReview 审核理赔
func (c *InsuranceController) InsuranceClaimReview(ctx *gin.Context) {
	id := ctx.Param("id")
	var claim models.InsuranceClaim
	if err := c.DB.First(&claim, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "理赔不存在"})
		return
	}

	var req struct {
		Status         string  `json:"status" binding:"required"` // approved, rejected
		ApprovedAmount float64 `json:"approved_amount"`
		ReviewNotes    string  `json:"review_notes"`
		RejectionReason string `json:"rejection_reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	claim.Status = req.Status
	claim.ApprovedAmount = req.ApprovedAmount
	claim.ReviewNotes = req.ReviewNotes
	if req.RejectionReason != "" {
		claim.RejectionReason = req.RejectionReason
	}

	if err := c.DB.Save(&claim).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": claim})
}

// ============ InsuranceClaimDocument 文档管理 ============

// InsuranceClaimDocumentList 文档列表
func (c *InsuranceController) InsuranceClaimDocumentList(ctx *gin.Context) {
	claimUUID := ctx.Query("claim_uuid")
	var documents []models.InsuranceClaimDocument
	var total int64

	query := c.DB.Model(&models.InsuranceClaimDocument{})
	if claimUUID != "" {
		query = query.Where("claim_uuid = ?", claimUUID)
	}

	query.Count(&total)
	if err := query.Order("id DESC").Find(&documents).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": documents, "total": total,
	}})
}

// InsuranceClaimDocumentCreate 创建文档
func (c *InsuranceController) InsuranceClaimDocumentCreate(ctx *gin.Context) {
	var doc models.InsuranceClaimDocument
	if err := ctx.ShouldBindJSON(&doc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&doc).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": doc})
}

// InsuranceClaimDocumentDelete 删除文档
func (c *InsuranceController) InsuranceClaimDocumentDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.InsuranceClaimDocument{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
