package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CompanyController 公司管理控制器
type CompanyController struct {
	DB *gorm.DB
}

// CompanyList 公司列表
func (c *CompanyController) CompanyList(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var companies []models.Company
	var total int64

	query := c.DB.Model(&models.Company{})

	// 按租户筛选
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 搜索条件
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("company_name LIKE ? OR company_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&companies).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      companies,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CompanyCreate 创建公司
func (c *CompanyController) CompanyCreate(ctx *gin.Context) {
	var company models.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tenantID := middleware.GetTenantID(ctx)
	if tenantID != "" {
		company.TenantID = tenantID
	}

	if err := c.DB.Create(&company).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": company})
}

// CompanyGet 获取公司详情
func (c *CompanyController) CompanyGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var company models.Company
	if err := c.DB.First(&company, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "公司不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": company})
}

// CompanyUpdate 更新公司
func (c *CompanyController) CompanyUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var company models.Company
	if err := c.DB.First(&company, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "公司不存在"})
		return
	}

	var updateData models.Company
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 更新字段
	updateFields := map[string]interface{}{}
	if updateData.CompanyName != "" {
		updateFields["company_name"] = updateData.CompanyName
	}
	if updateData.ShortName != "" {
		updateFields["short_name"] = updateData.ShortName
	}
	if updateData.Logo != "" {
		updateFields["logo"] = updateData.Logo
	}
	if updateData.Province != "" {
		updateFields["province"] = updateData.Province
	}
	if updateData.City != "" {
		updateFields["city"] = updateData.City
	}
	if updateData.District != "" {
		updateFields["district"] = updateData.District
	}
	if updateData.Address != "" {
		updateFields["address"] = updateData.Address
	}
	if updateData.LegalPerson != "" {
		updateFields["legal_person"] = updateData.LegalPerson
	}
	if updateData.Contact != "" {
		updateFields["contact"] = updateData.Contact
	}
	if updateData.Phone != "" {
		updateFields["phone"] = updateData.Phone
	}
	if updateData.Email != "" {
		updateFields["email"] = updateData.Email
	}
	updateFields["status"] = updateData.Status
	updateFields["sort"] = updateData.Sort
	if updateData.Remark != "" {
		updateFields["remark"] = updateData.Remark
	}

	if err := c.DB.Model(&company).Updates(updateFields).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&company, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": company})
}

// CompanyDelete 删除公司
func (c *CompanyController) CompanyDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	// 检查是否有部门关联
	var deptCount int64
	c.DB.Model(&models.Department{}).Where("company_id = ?", id).Count(&deptCount)
	if deptCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该公司下还有部门，无法删除"})
		return
	}

	if err := c.DB.Delete(&models.Company{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
