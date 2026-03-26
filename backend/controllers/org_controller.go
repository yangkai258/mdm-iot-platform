package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrgController 组织管理控制器
type OrgController struct {
	DB *gorm.DB
}

// CompanyList 公司列表
func (c *OrgController) CompanyList(ctx *gin.Context) {
	var companies []models.Company
	var total int64

	query := c.DB.Model(&models.Company{})

	// 搜索条件
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("company_name LIKE ? OR company_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	query.Count(&total)

	// 分页
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
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":  companies,
			"total": total,
			"page":  page,
			"page_size": pageSize,
		},
	})
}

// CompanyCreate 创建公司
func (c *OrgController) CompanyCreate(ctx *gin.Context) {
	var company models.Company
	if err := ctx.ShouldBindJSON(&company); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&company).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": company})
}

// CompanyUpdate 更新公司
func (c *OrgController) CompanyUpdate(ctx *gin.Context) {
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

	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updateData})
}

// CompanyDelete 删除公司
func (c *OrgController) CompanyDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Company{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// DepartmentList 部门列表
func (c *OrgController) DepartmentList(ctx *gin.Context) {
	var departments []models.Department
	var total int64

	query := c.DB.Model(&models.Department{})

	// 公司筛选
	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	// 上级部门筛选
	if parentID := ctx.Query("parent_id"); parentID != "" {
		query = query.Where("parent_id = ?", parentID)
	} else {
		// 不加 parent_id 条件，直接查所有（让前端传 parent_id 来筛选）
		// 避免 GORM IS NULL 查询在不同版本的行为差异
	}

	query.Count(&total)

	if err := query.Order("sort ASC, id DESC").Find(&departments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":  departments,
			"total": total,
		},
	})
}

// DepartmentTree 部门树
func (c *OrgController) DepartmentTree(ctx *gin.Context) {
	var departments []models.Department
	c.DB.Order("sort ASC, id DESC").Find(&departments)

	// 构建树形结构
	tree := buildDeptTree(departments, 0)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tree})
}

func buildDeptTree(depts []models.Department, parentID uint) []models.Department {
	var tree []models.Department
	for _, dept := range depts {
		if (parentID == 0 && dept.ParentID == nil) || (dept.ParentID != nil && *dept.ParentID == parentID) {
			children := buildDeptTree(depts, dept.ID)
			if len(children) > 0 {
				dept.Children = children
			}
			tree = append(tree, dept)
		}
	}
	return tree
}

// DepartmentCreate 创建部门
func (c *OrgController) DepartmentCreate(ctx *gin.Context) {
	var dept models.Department
	if err := ctx.ShouldBindJSON(&dept); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 生成路径（Level 在插入前计算，Path 在插入后用真实 dept.ID 更新）
	if dept.ParentID != nil {
		var parent models.Department
		if err := c.DB.First(&parent, *dept.ParentID).Error; err == nil {
			dept.Level = parent.Level + 1
		}
	}

	if err := c.DB.Create(&dept).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	// 插入后用真实的 dept.ID 更新 path
	if dept.ParentID != nil {
		var parent models.Department
		if err := c.DB.First(&parent, *dept.ParentID).Error; err == nil {
			newPath := parent.Path + "/" + strconv.Itoa(int(dept.ID))
			c.DB.Model(&dept).Update("path", newPath)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": dept})
}

// DepartmentUpdate 更新部门
func (c *OrgController) DepartmentUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var dept models.Department
	if err := c.DB.First(&dept, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}

	if err := ctx.ShouldBindJSON(&dept); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Save(&dept).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": dept})
}

// DepartmentDelete 删除部门
func (c *OrgController) DepartmentDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	// 检查是否有子部门
	var count int64
	c.DB.Model(&models.Department{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先删除子部门"})
		return
	}

	if err := c.DB.Delete(&models.Department{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// PositionList 岗位列表
func (c *OrgController) PositionList(ctx *gin.Context) {
	var positions []models.Position
	var total int64

	query := c.DB.Model(&models.Position{})

	if deptID := ctx.Query("dept_id"); deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}

	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&positions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":   positions,
			"total":  total,
			"page":   page,
			"page_size": pageSize,
		},
	})
}

// PositionCreate 创建岗位
func (c *OrgController) PositionCreate(ctx *gin.Context) {
	var position models.Position
	if err := ctx.ShouldBindJSON(&position); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&position).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// PositionUpdate 更新岗位
func (c *OrgController) PositionUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var position models.Position
	if err := c.DB.First(&position, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "岗位不存在"})
		return
	}

	if err := ctx.ShouldBindJSON(&position); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Save(&position).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// PositionDelete 删除岗位
func (c *OrgController) PositionDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Position{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// EmployeeList 员工列表
func (c *OrgController) EmployeeList(ctx *gin.Context) {
	var employees []models.Employee
	var total int64

	query := c.DB.Preload("Department").Preload("Position").Model(&models.Employee{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("emp_name LIKE ? OR emp_code LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if deptID := ctx.Query("dept_id"); deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}

	if companyID := ctx.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&employees).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":   employees,
			"total":  total,
			"page":   page,
			"page_size": pageSize,
		},
	})
}

// EmployeeCreate 创建员工
func (c *OrgController) EmployeeCreate(ctx *gin.Context) {
	var employee models.Employee
	if err := ctx.ShouldBindJSON(&employee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&employee).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": employee})
}

// EmployeeUpdate 更新员工
func (c *OrgController) EmployeeUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var employee models.Employee
	if err := c.DB.First(&employee, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "员工不存在"})
		return
	}

	if err := ctx.ShouldBindJSON(&employee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Save(&employee).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": employee})
}

// EmployeeDelete 删除员工
func (c *OrgController) EmployeeDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Employee{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// StandardPositionList 基准岗位列表
func (c *OrgController) StandardPositionList(ctx *gin.Context) {
	var positions []models.StandardPosition
	var total int64

	query := c.DB.Model(&models.StandardPosition{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("sp_name LIKE ? OR sp_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("sort ASC, id DESC").Find(&positions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":   positions,
			"total":  total,
			"page":   page,
			"page_size": pageSize,
		},
	})
}

// StandardPositionCreate 创建基准岗位
func (c *OrgController) StandardPositionCreate(ctx *gin.Context) {
	var position models.StandardPosition
	if err := ctx.ShouldBindJSON(&position); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&position).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// StandardPositionUpdate 更新基准岗位
func (c *OrgController) StandardPositionUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var position models.StandardPosition
	if err := c.DB.First(&position, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "基准岗位不存在"})
		return
	}

	if err := ctx.ShouldBindJSON(&position); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Save(&position).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": position})
}

// StandardPositionDelete 删除基准岗位
func (c *OrgController) StandardPositionDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.StandardPosition{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// StandardPositionEnable 启用基准岗位
// POST /api/v1/org/standard-positions/:id/enable
func (c *OrgController) StandardPositionEnable(ctx *gin.Context) {
	id := ctx.Param("id")
	var position models.StandardPosition
	if err := c.DB.First(&position, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "基准岗位不存在"})
		return
	}
	if err := c.DB.Model(&position).Update("status", 1).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "启用失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"id": position.ID, "status": 1}})
}

// StandardPositionDisable 禁用基准岗位
// POST /api/v1/org/standard-positions/:id/disable
func (c *OrgController) StandardPositionDisable(ctx *gin.Context) {
	id := ctx.Param("id")
	var position models.StandardPosition
	if err := c.DB.First(&position, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "基准岗位不存在"})
		return
	}
	if err := c.DB.Model(&position).Update("status", 0).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "禁用失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"id": position.ID, "status": 0}})
}
