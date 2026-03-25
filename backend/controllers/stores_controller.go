package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StoresController 门店管理控制器
type StoresController struct {
	DB *gorm.DB
}

// StoreList 门店列表
func (c *StoresController) StoreList(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var stores []models.Store
	var total int64

	query := c.DB.Model(&models.Store{})

	// 按租户筛选
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 搜索条件 - 按门店编码或名称
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("store_code LIKE ? OR store_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 门店类型筛选
	if storeType := ctx.Query("store_type"); storeType != "" {
		query = query.Where("store_type = ?", storeType)
	}

	// 省份筛选
	if province := ctx.Query("province"); province != "" {
		query = query.Where("province = ?", province)
	}

	// 城市筛选
	if city := ctx.Query("city"); city != "" {
		query = query.Where("city = ?", city)
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

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&stores).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      stores,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// StoreCreate 创建门店
func (c *StoresController) StoreCreate(ctx *gin.Context) {
	var store models.Store
	if err := ctx.ShouldBindJSON(&store); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	tenantID := middleware.GetTenantID(ctx)
	if tenantID != "" {
		store.TenantID = tenantID
	}

	// 检查门店编码唯一性（租户内唯一）
	var count int64
	query := c.DB.Model(&models.Store{}).Where("store_code = ?", store.StoreCode)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	query.Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "门店编码已存在"})
		return
	}

	if err := c.DB.Create(&store).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": store})
}

// StoreGet 获取门店详情
func (c *StoresController) StoreGet(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")
	
	query := c.DB.Model(&models.Store{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	var store models.Store
	if err := query.First(&store, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "门店不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": store})
}

// StoreUpdate 更新门店
func (c *StoresController) StoreUpdate(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")
	var store models.Store
	
	// 应用租户过滤
	query := c.DB.Model(&models.Store{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.First(&store, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "门店不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	var updateData struct {
		StoreCode  string `json:"store_code"`
		StoreName  string `json:"store_name"`
		StoreType  int    `json:"store_type"`
		Province   string `json:"province"`
		City       string `json:"city"`
		District   string `json:"district"`
		Address    string `json:"address"`
		Contact    string `json:"contact"`
		Phone      string `json:"phone"`
		Status     int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查门店编码唯一性（租户内唯一，排除自身）
	if updateData.StoreCode != "" && updateData.StoreCode != store.StoreCode {
		var count int64
		query := c.DB.Model(&models.Store{}).Where("store_code = ? AND id != ?", updateData.StoreCode, id)
		if tenantID != "" {
			query = query.Where("tenant_id = ?", tenantID)
		}
		query.Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "门店编码已存在"})
			return
		}
	}

	// 更新字段 - 只更新非零值和显式提供的字段
	updateFields := map[string]interface{}{}
	if updateData.StoreCode != "" {
		updateFields["store_code"] = updateData.StoreCode
	}
	if updateData.StoreName != "" {
		updateFields["store_name"] = updateData.StoreName
	}
	if updateData.StoreType > 0 {
		updateFields["store_type"] = updateData.StoreType
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
	if updateData.Contact != "" {
		updateFields["contact"] = updateData.Contact
	}
	if updateData.Phone != "" {
		updateFields["phone"] = updateData.Phone
	}
	if updateData.Status > 0 {
		updateFields["status"] = updateData.Status
	}

	if err := c.DB.Model(&store).Updates(updateFields).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&store, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": store})
}

// StoreDelete 删除门店
func (c *StoresController) StoreDelete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	query := c.DB.Model(&models.Store{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	var store models.Store
	if err := query.First(&store, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "门店不存在"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		}
		return
	}

	// 检查是否有会员关联
	var memberCount int64
	c.DB.Model(&models.Member{}).Where("store_id = ?", id).Count(&memberCount)
	if memberCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该门店下还有会员，无法删除"})
		return
	}

	if err := c.DB.Delete(&models.Store{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// StoreUpdateStatus 更新门店状态
func (c *StoresController) StoreUpdateStatus(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var updateData struct {
		Status int `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	query := c.DB.Model(&models.Store{}).Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	result := query.Update("status", updateData.Status)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新状态失败"})
		return
	}
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "门店不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// StoreBatchDelete 批量删除门店
func (c *StoresController) StoreBatchDelete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var ids struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(ids.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要删除的门店"})
		return
	}

	// 检查是否有会员关联
	var memberCount int64
	c.DB.Model(&models.Member{}).Where("store_id IN ?", ids.IDs).Count(&memberCount)
	if memberCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "部分门店下还有会员，无法删除"})
		return
	}

	query := c.DB.Where("id IN ?", ids.IDs)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.Delete(&models.Store{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// StoreStatistics 获取门店统计
func (c *StoresController) StoreStatistics(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)

	query := c.DB.Model(&models.Store{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	var total int64
	var activeCount int64
	var directStoreCount int64
	var franchiseCount int64

	query.Count(&total)
	query.Where("status = 1").Count(&activeCount)
	query.Where("store_type = 1").Count(&directStoreCount)
	query.Where("store_type = 2").Count(&franchiseCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"total":           total,
			"active":          activeCount,
			"inactive":        total - activeCount,
			"direct_stores":   directStoreCount,
			"franchises":      franchiseCount,
		},
	})
}
