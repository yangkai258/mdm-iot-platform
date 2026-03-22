package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"
	plugins "mdm-backend/plugins"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApiPermissionController API权限控制器
type ApiPermissionController struct {
	DB *gorm.DB
}

// List API权限列表
// GET /api/v1/api-permissions
func (c *ApiPermissionController) List(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var apis []models.ApiPermission
	var total int64

	query := c.DB.Model(&models.ApiPermission{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("api_name LIKE ? OR api_path LIKE ? OR permission_code LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if method := ctx.Query("method"); method != "" {
		query = query.Where("method = ?", method)
	}

	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&apis).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      apis,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建API权限
// POST /api/v1/api-permissions
func (c *ApiPermissionController) Create(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "租户ID不能为空"})
		return
	}

	var req struct {
		ApiPath        string `json:"api_path" binding:"required"`
		ApiName        string `json:"api_name" binding:"required"`
		Method         string `json:"method"`
		PermissionCode string `json:"permission_code"`
		MenuID         *uint  `json:"menu_id"`
		Status         int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	api := models.ApiPermission{
		ApiPath:        req.ApiPath,
		ApiName:        req.ApiName,
		Method:         req.Method,
		PermissionCode: req.PermissionCode,
		MenuID:         req.MenuID,
		Status:         req.Status,
		TenantID:       tenantID,
	}
	if api.Status == 0 {
		api.Status = 1
	}

	db := plugins.WithTenantID(c.DB, tenantID)
	if err := db.Create(&api).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": api})
}

// Update 更新API权限
// PUT /api/v1/api-permissions/:id
func (c *ApiPermissionController) Update(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var api models.ApiPermission
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&api).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "API权限不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		ApiPath        string `json:"api_path"`
		ApiName        string `json:"api_name"`
		Method         string `json:"method"`
		PermissionCode string `json:"permission_code"`
		MenuID         *uint  `json:"menu_id"`
		Status         int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.ApiPath != "" {
		updates["api_path"] = req.ApiPath
	}
	if req.ApiName != "" {
		updates["api_name"] = req.ApiName
	}
	if req.Method != "" {
		updates["method"] = req.Method
	}
	if req.PermissionCode != "" {
		updates["permission_code"] = req.PermissionCode
	}
	if req.MenuID != nil {
		updates["menu_id"] = req.MenuID
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := c.DB.Model(&api).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.DB.First(&api, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": api})
}

// Delete 删除API权限
// DELETE /api/v1/api-permissions/:id
func (c *ApiPermissionController) Delete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.ApiPermission{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Import 批量导入API权限
// POST /api/v1/api-permissions/import
func (c *ApiPermissionController) Import(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "租户ID不能为空"})
		return
	}

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请上传文件"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// 跳过表头
	if _, err := reader.Read(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件格式错误"})
		return
	}

	var imported, skipped int64
	db := plugins.WithTenantID(c.DB, tenantID)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// 格式: api_path,api_name,method,permission_code,menu_id
		if len(record) < 3 {
			skipped++
			continue
		}

		api := models.ApiPermission{
			ApiPath:        record[0],
			ApiName:        record[1],
			Method:         "",
			PermissionCode: "",
			Status:         1,
			TenantID:       tenantID,
		}

		if len(record) >= 3 && record[2] != "" {
			api.Method = record[2]
		}
		if len(record) >= 4 && record[3] != "" {
			api.PermissionCode = record[3]
		}
		if len(record) >= 5 && record[4] != "" {
			if menuID, err := strconv.ParseUint(record[4], 10, 32); err == nil {
				id := uint(menuID)
				api.MenuID = &id
			}
		}

		if err := db.Create(&api).Error; err != nil {
			skipped++
			continue
		}
		imported++
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"imported": imported,
			"skipped":  skipped,
		},
	})
}

// Export 批量导出API权限
// GET /api/v1/api-permissions/export
func (c *ApiPermissionController) Export(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	var apis []models.ApiPermission
	query := c.DB.Model(&models.ApiPermission{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if err := query.Order("id ASC").Find(&apis).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出失败"})
		return
	}

	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=api_permissions_%s.csv", tenantID))

	writer := csv.NewWriter(ctx.Writer)
	// 写表头
	writer.Write([]string{"api_path", "api_name", "method", "permission_code", "menu_id", "status"})

	for _, api := range apis {
		menuIDStr := ""
		if api.MenuID != nil {
			menuIDStr = strconv.FormatUint(uint64(*api.MenuID), 10)
		}
		writer.Write([]string{
			api.ApiPath,
			api.ApiName,
			api.Method,
			api.PermissionCode,
			menuIDStr,
			strconv.Itoa(api.Status),
		})
	}
	writer.Flush()
}
