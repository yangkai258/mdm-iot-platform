package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// I18nController 国际化控制器
type I18nController struct {
	DB *gorm.DB
}

// TranslationList 翻译列表
func (c *I18nController) TranslationList(ctx *gin.Context) {
	var translations []models.Translation
	var total int64

	query := c.DB.Model(&models.Translation{})

	if locale := ctx.Query("locale"); locale != "" {
		query = query.Where("locale = ?", locale)
	}
	if namespace := ctx.Query("namespace"); namespace != "" {
		query = query.Where("namespace = ?", namespace)
	}
	if key := ctx.Query("key"); key != "" {
		query = query.Where("key LIKE ?", "%"+key+"%")
	}
	if isActive := ctx.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "50"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("locale ASC, namespace ASC, `key` ASC").Find(&translations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": translations, "total": total, "page": page, "page_size": pageSize,
	}})
}

// TranslationGet 获取翻译
func (c *I18nController) TranslationGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var translation models.Translation
	if err := c.DB.First(&translation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "翻译不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": translation})
}

// TranslationCreate 创建翻译
func (c *I18nController) TranslationCreate(ctx *gin.Context) {
	var translation models.Translation
	if err := ctx.ShouldBindJSON(&translation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&translation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": translation})
}

// TranslationUpdate 更新翻译
func (c *I18nController) TranslationUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var translation models.Translation
	if err := c.DB.First(&translation, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "翻译不存在"})
		return
	}
	var updateData models.Translation
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

// TranslationDelete 删除翻译
func (c *I18nController) TranslationDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.Translation{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// TranslationBatchCreate 批量创建翻译
func (c *I18nController) TranslationBatchCreate(ctx *gin.Context) {
	var translations []models.Translation
	if err := ctx.ShouldBindJSON(&translations); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&translations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"count": len(translations)}})
}

// TranslationByLocale 获取指定语言的翻译（前端加载用）
func (c *I18nController) TranslationByLocale(ctx *gin.Context) {
	locale := ctx.Param("locale")
	var translations []models.Translation

	if err := c.DB.Where("locale = ? AND is_active = ?", locale, true).Find(&translations).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 转换为 key-value map
	result := make(map[string]string)
	for _, t := range translations {
		result[t.Key] = t.Value
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// TranslationNamespaceList 获取所有命名空间
func (c *I18nController) TranslationNamespaceList(ctx *gin.Context) {
	var namespaces []string
	c.DB.Model(&models.Translation{}).Distinct("namespace").Pluck("namespace", &namespaces)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": namespaces})
}

// TranslationLocaleList 获取所有语言
func (c *I18nController) TranslationLocaleList(ctx *gin.Context) {
	var locales []string
	c.DB.Model(&models.Translation{}).Distinct("locale").Pluck("locale", &locales)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": locales})
}
