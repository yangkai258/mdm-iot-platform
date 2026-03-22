package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/i18n"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// I18nController 国际化控制器
type I18nController struct {
	DB  *gorm.DB
	svc *i18n.TranslationService
}

// NewI18nController 创建国际化控制器
func NewI18nController(db *gorm.DB) *I18nController {
	return &I18nController{
		DB:  db,
		svc: i18n.NewTranslationService(db),
	}
}

// ListTranslations 获取翻译列表
// @Summary 获取翻译列表
// @Tags i18n
// @Produce json
// @Param locale query string false "语言代码 (e.g. zh-CN, en-US)"
// @Param namespace query string false "命名空间"
// @Param key query string false "翻译Key (模糊匹配)"
// @Param tags query string false "标签 (模糊匹配)"
// @Param is_active query bool false "是否激活"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/i18n/translations [get]
func (ctrl *I18nController) ListTranslations(c *gin.Context) {
	var filter models.TranslationFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	translations, total, err := ctrl.svc.ListTranslations(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"items":       translations,
			"total":       total,
			"page":        filter.Page,
			"page_size":   filter.PageSize,
			"total_pages": (total + int64(filter.PageSize) - 1) / int64(filter.PageSize),
		},
	})
}

// GetTranslation 获取翻译详情
// @Summary 获取翻译详情
// @Tags i18n
// @Produce json
// @Param id path int true "翻译ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/i18n/translations/{id} [get]
func (ctrl *I18nController) GetTranslation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	translation, err := ctrl.svc.GetTranslationByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "translation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": translation,
	})
}

// CreateTranslation 创建翻译
// @Summary 创建翻译
// @Tags i18n
// @Accept json
// @Produce json
// @Param translation body models.TranslationRequest true "翻译信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/i18n/translations [post]
func (ctrl *I18nController) CreateTranslation(c *gin.Context) {
	var req models.TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证locale格式
	if !isValidLocale(req.Locale) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid locale format"})
		return
	}

	translation, err := ctrl.svc.CreateTranslation(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": translation,
	})
}

// UpdateTranslation 更新翻译
// @Summary 更新翻译
// @Tags i18n
// @Accept json
// @Produce json
// @Param id path int true "翻译ID"
// @Param translation body models.TranslationRequest true "翻译信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/i18n/translations/{id} [put]
func (ctrl *I18nController) UpdateTranslation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证locale格式
	if !isValidLocale(req.Locale) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid locale format"})
		return
	}

	translation, err := ctrl.svc.UpdateTranslation(uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "translation not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": translation,
	})
}

// DeleteTranslation 删除翻译
// @Summary 删除翻译
// @Tags i18n
// @Produce json
// @Param id path int true "翻译ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/i18n/translations/{id} [delete]
func (ctrl *I18nController) DeleteTranslation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := ctrl.svc.DeleteTranslation(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "translation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "translation deleted",
	})
}

// GetSupportedLocales 获取支持的语言列表
// @Summary 获取支持的语言列表
// @Tags i18n
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/i18n/locales [get]
func (ctrl *I18nController) GetSupportedLocales(c *gin.Context) {
	locales, err := ctrl.svc.GetSupportedLocales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": locales,
	})
}

// GetNamespaces 获取所有命名空间
// @Summary 获取所有命名空间
// @Tags i18n
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/i18n/namespaces [get]
func (ctrl *I18nController) GetNamespaces(c *gin.Context) {
	namespaces, err := ctrl.svc.GetNamespaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": namespaces,
	})
}

// isValidLocale 验证locale格式
func isValidLocale(locale string) bool {
	if len(locale) < 2 {
		return false
	}
	// 基本格式检查：language-COUNTRY 或 language
	for _, c := range locale {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-' || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}
