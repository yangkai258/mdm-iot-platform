package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmployeeController 员工管理控制器
type EmployeeController struct {
	DB *gorm.DB
}

// EmployeeList 员工列表
func (c *EmployeeController) EmployeeList(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var employees []models.Employee
	var total int64

	query := c.DB.Preload("Department").Preload("Position").Model(&models.Employee{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("emp_name LIKE ? OR emp_code LIKE ? OR phone LIKE ? OR id_card LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if deptID := ctx.Query("dept_id"); deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}

	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if empStatus := ctx.Query("emp_status"); empStatus != "" {
		query = query.Where("emp_status = ?", empStatus)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&employees).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      employees,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// EmployeeOnboard 办理入职
func (c *EmployeeController) EmployeeOnboard(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var input struct {
		EmpCode    string  `json:"emp_code" required:"true"`
		EmpName    string  `json:"emp_name" required:"true"`
		Gender     string  `json:"gender"`
		BirthDate  string  `json:"birth_date"`
		Phone      string  `json:"phone"`
		Email      string  `json:"email"`
		IDCard     string  `json:"id_card"`
		Province   string  `json:"province"`
		City       string  `json:"city"`
		District   string  `json:"district"`
		Address    string  `json:"address"`
		DeptID     uint    `json:"dept_id"`
		PositionID uint    `json:"position_id"`
		CompanyID  uint    `json:"company_id" required:"true"`
		EntryDate  string  `json:"entry_date"`
		Remark     string  `json:"remark"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	employee := models.Employee{
		EmpCode:   input.EmpCode,
		EmpName:   input.EmpName,
		Gender:    input.Gender,
		Phone:     input.Phone,
		Email:     input.Email,
		IDCard:    input.IDCard,
		Province:  input.Province,
		City:      input.City,
		District:  input.District,
		Address:   input.Address,
		DeptID:    input.DeptID,
		PositionID: input.PositionID,
		CompanyID: input.CompanyID,
		EmpStatus: 1, // 在职
		Status:    1,
		Remark:    input.Remark,
	}

	if tenantID != "" {
		employee.TenantID = tenantID
	}

	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil {
			employee.BirthDate = &t
		}
	}

	if input.EntryDate != "" {
		t, err := time.Parse("2006-01-02", input.EntryDate)
		if err == nil {
			employee.EntryDate = &t
		}
	} else {
		now := time.Now()
		employee.EntryDate = &now
	}

	if err := c.DB.Create(&employee).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "入职办理失败"})
		return
	}

	c.DB.Preload("Department").Preload("Position").First(&employee, employee.ID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": employee})
}

// EmployeeGet 获取员工详情
func (c *EmployeeController) EmployeeGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var employee models.Employee
	if err := c.DB.Preload("Department").Preload("Position").First(&employee, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "员工不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": employee})
}

// EmployeeLeave 办理离职
func (c *EmployeeController) EmployeeLeave(ctx *gin.Context) {
	id := ctx.Param("id")
	var employee models.Employee
	if err := c.DB.First(&employee, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "员工不存在"})
		return
	}

	if employee.EmpStatus == 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该员工已经离职"})
		return
	}

	var input struct {
		LeaveDate string `json:"leave_date"`
		Reason    string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		// 允许空body
		input = struct {
			LeaveDate string `json:"leave_date"`
			Reason    string `json:"reason"`
		}{}
	}

	updates := map[string]interface{}{
		"emp_status": 2, // 离职
		"status":     0, // 账号禁用
	}
	if input.Reason != "" {
		updates["remark"] = employee.Remark + "\n离职原因: " + input.Reason
	}

	if input.LeaveDate != "" {
		t, err := time.Parse("2006-01-02", input.LeaveDate)
		if err == nil {
			updates["leave_date"] = t
		}
	} else {
		updates["leave_date"] = time.Now()
	}

	if err := c.DB.Model(&employee).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "离职办理失败"})
		return
	}

	c.DB.Preload("Department").Preload("Position").First(&employee, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": employee})
}

// EmployeeUpdate 更新员工
func (c *EmployeeController) EmployeeUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var employee models.Employee
	if err := c.DB.First(&employee, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "员工不存在"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	allowedFields := []string{
		"emp_name", "gender", "birth_date", "phone", "email", "id_card",
		"photo", "province", "city", "district", "address",
		"dept_id", "position_id", "company_id", "status", "remark",
	}
	updates := map[string]interface{}{}
	for _, field := range allowedFields {
		if val, ok := updateData[field]; ok {
			updates[field] = val
		}
	}

	if err := c.DB.Model(&employee).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.Preload("Department").Preload("Position").First(&employee, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": employee})
}

// EmployeeDelete 删除员工
func (c *EmployeeController) EmployeeDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Employee{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
