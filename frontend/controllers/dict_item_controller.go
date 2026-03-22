package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictItemController 字典项控制器
type DictItemController struct {
	DB *gorm.DB
}

func (c *DictItemController) List(ctx *gin.Context) {
	var items []models.SysDictItem
	query := c.DB.Model(&models.SysDictItem{})

	if dictTypeID := ctx.Query("dict_type_id"); dictTypeID != "" {
		query = query.Where("dict_type_id = ?", dictTypeID)
	}
	if typeCode := ctx.Query("type_code"); typeCode != "" {
		var dictType models.SysDictType
		if err := c.DB.Where("type_code = ?", typeCode).First(&dictType).Error; err == nil {
			query = query.Where("dict_type_id = ?", dictType.ID)
		}
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query.Order("sort ASC, id ASC").Offset((page-1)*pageSize).Limit(pageSize).Find(&items)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": items, "pagination": gin.H{"total": total, "current": page, "pageSize": pageSize}}})
}

func (c *DictItemController) Get(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysDictItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典项不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item})
}

func (c *DictItemController) Create(ctx *gin.Context) {
	var req struct {
		DictTypeID uint   `json:"dict_type_id" binding:"required"`
		ItemText   string `json:"item_text" binding:"required"`
		ItemValue  string `json:"item_value" binding:"required"`
		Sort       int    `json:"sort"`
		Status     int    `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var dictType models.SysDictType
	if err := c.DB.First(&dictType, req.DictTypeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "字典类型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询字典类型失败"})
		return
	}

	item := models.SysDictItem{DictTypeID: req.DictTypeID, ItemText: req.ItemText, ItemValue: req.ItemValue, Sort: req.Sort, Status: req.Status, Remark: req.Remark}
	if item.Status == 0 {
		item.Status = 1
	}
	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "创建成功"})
}

func (c *DictItemController) Update(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysDictItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典项不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		DictTypeID uint   `json:"dict_type_id"`
		ItemText   string `json:"item_text"`
		ItemValue  string `json:"item_value"`
		Sort       int    `json:"sort"`
		Status     int    `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.DictTypeID != 0 {
		item.DictTypeID = req.DictTypeID
	}
	if req.ItemText != "" {
		item.ItemText = req.ItemText
	}
	if req.ItemValue != "" {
		item.ItemValue = req.ItemValue
	}
	if req.Remark != "" {
		item.Remark = req.Remark
	}
	item.Sort = req.Sort
	if req.Status != 0 {
		item.Status = req.Status
	}

	if err := c.DB.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "更新成功"})
}

func (c *DictItemController) Delete(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysDictItem
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典项不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	if err := c.DB.Delete(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

func (c *DictItemController) GetByType(ctx *gin.Context) {
	typeCode := ctx.Param("type")
	var dictType models.SysDictType
	if err := c.DB.Where("type_code = ?", typeCode).First(&dictType).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典类型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var items []models.SysDictItem
	c.DB.Where("dict_type_id = ? AND status = 1", dictType.ID).Order("sort ASC, id ASC").Find(&items)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": items})
}
