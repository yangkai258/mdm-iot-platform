package controllers

import (
	"fmt"
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

// DeptRow 部门查询结果行
type DeptRow struct {
	ID        int
	DeptID    *string
	ParentID  *int
	DeptCode  *string
	DeptName  *string
	ManagerID *int
	Status    *string
	SortOrder *int
	TenantID  *string
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

	// 自动生成公司编号
	if company.CompanyCode == "" {
		var maxID int64
		c.DB.Model(&models.Company{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
		company.CompanyCode = fmt.Sprintf("CMP%05d", maxID+1)
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
	var rows []DeptRow
	result := c.DB.Raw(`
		SELECT id, dept_id, parent_id, dept_code, dept_name,
		       manager_id, status, sort_order, tenant_id
		FROM departments
		ORDER BY sort_order ASC, id DESC
		LIMIT 100
	`).Scan(&rows)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	list := make([]map[string]interface{}, 0)
	for _, r := range rows {
		list = append(list, map[string]interface{}{
			"id":         r.ID,
			"dept_id":    r.DeptID,
			"parent_id":  r.ParentID,
			"dept_code":  r.DeptCode,
			"dept_name":  r.DeptName,
			"manager_id":  r.ManagerID,
			"status":     r.Status,
			"sort_order": r.SortOrder,
			"tenant_id":  r.TenantID,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  list,
			"total": len(list),
		},
	})
}

// DepartmentTree 部门树
func (c *OrgController) DepartmentTree(ctx *gin.Context) {
	var rows []DeptRow
	c.DB.Raw(`
		SELECT id, dept_id, parent_id, dept_code, dept_name,
		       manager_id, status, sort_order, tenant_id
		FROM departments
		ORDER BY sort_order ASC, id DESC
	`).Scan(&rows)

	// 构建树形结构
	tree := buildDeptTreeFromRows(rows, 0)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tree})
}

// buildDeptTreeFromRows 从行数据构建树
func buildDeptTreeFromRows(rows []DeptRow, parentID int) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, r := range rows {
		isRoot := r.ParentID == nil || *r.ParentID == 0
		if (parentID == 0 && isRoot) || (parentID != 0 && r.ParentID != nil && *r.ParentID == parentID) {
			node := map[string]interface{}{
				"id":          r.ID,
				"dept_id":     r.DeptID,
				"dept_code":   r.DeptCode,
				"dept_name":   r.DeptName,
				"manager_id":   r.ManagerID,
				"status":      r.Status,
				"sort_order":  r.SortOrder,
				"children":    buildDeptTreeFromRows(rows, r.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
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
	var req struct {
		DeptCode  string  `json:"dept_code" binding:"required"`
		DeptName  string  `json:"dept_name" binding:"required"`
		ParentID  *int    `json:"parent_id"`
		ManagerID *int    `json:"manager_id"`
		Status    string  `json:"status"`
		SortOrder *int    `json:"sort_order"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.Status == "" {
		req.Status = "active"
	}

	// 使用原始 SQL 插入
	sql := `INSERT INTO departments (dept_id, dept_code, dept_name, parent_id, manager_id, status, sort_order)
	        VALUES (gen_random_uuid()::text, $1, $2, $3, $4, $5, $6) RETURNING id`

	var newID int
	err := c.DB.Raw(sql, req.DeptCode, req.DeptName, req.ParentID, req.ManagerID, req.Status, req.SortOrder).Scan(&newID).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id":        newID,
			"dept_code": req.DeptCode,
			"dept_name": req.DeptName,
			"parent_id": req.ParentID,
			"manager_id": req.ManagerID,
			"status":    req.Status,
			"sort_order": req.SortOrder,
		},
	})
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

	// 自动生成员工工号
	if employee.EmpCode == "" {
		var maxID int64
		c.DB.Model(&models.Employee{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
		employee.EmpCode = fmt.Sprintf("EMP%05d", maxID+1)
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
