package controllers

import (
	"net/http"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictTypeController 字典类型控制器
type DictTypeController struct {
	DB *gorm.DB
}

func (c *DictTypeController) List(ctx *gin.Context) {
	var items []models.SysDictType
	query := c.DB.Model(&models.SysDictType{})

	if typeCode := ctx.Query("type_code"); typeCode != "" {
		query = query.Where("type_code LIKE ?", "%"+typeCode+"%")
	}
	if typeName := ctx.Query("type_name"); typeName != "" {
		query = query.Where("type_name LIKE ?", "%"+typeName+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": items, "pagination": gin.H{"total": total, "current": page, "pageSize": pageSize}}})
}

func (c *DictTypeController) Get(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysDictType
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典类型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item})
}

func (c *DictTypeController) Create(ctx *gin.Context) {
	var req struct {
		TypeCode string `json:"type_code" binding:"required"`
		TypeName string `json:"type_name" binding:"required"`
		Status   int    `json:"status"`
		Remark   string `json:"remark"`
		Sort     int    `json:"sort"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var exist models.SysDictType
	if err := c.DB.Where("type_code = ?", req.TypeCode).First(&exist).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "字典编码已存在"})
		return
	}

	item := models.SysDictType{TypeCode: req.TypeCode, TypeName: req.TypeName, Status: req.Status, Remark: req.Remark, Sort: req.Sort}
	if item.Status == 0 {
		item.Status = 1
	}
	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "创建成功"})
}

func (c *DictTypeController) Update(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysDictType
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典类型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		TypeCode string `json:"type_code"`
		TypeName string `json:"type_name"`
		Status   int    `json:"status"`
		Remark   string `json:"remark"`
		Sort     int    `json:"sort"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.TypeCode != "" && req.TypeCode != item.TypeCode {
		var exist models.SysDictType
		if err := c.DB.Where("type_code = ? AND id != ?", req.TypeCode, id).First(&exist).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "字典编码已存在"})
			return
		}
		item.TypeCode = req.TypeCode
	}
	if req.TypeName != "" {
		item.TypeName = req.TypeName
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

func (c *DictTypeController) Delete(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysDictType
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "字典类型不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	tx := c.DB.Begin()
	if err := tx.Where("dict_type_id = ?", id).Delete(&models.SysDictItem{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除字典项失败"})
		return
	}
	if err := tx.Delete(&item).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}
