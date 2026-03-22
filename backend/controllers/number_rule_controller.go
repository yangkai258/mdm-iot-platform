package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NumberRuleController 编号规则控制器
type NumberRuleController struct {
	DB *gorm.DB
}

func (c *NumberRuleController) List(ctx *gin.Context) {
	var items []models.SysNumberRule
	query := c.DB.Model(&models.SysNumberRule{})

	if ruleName := ctx.Query("rule_name"); ruleName != "" {
		query = query.Where("rule_name LIKE ?", "%"+ruleName+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&items)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": items, "pagination": gin.H{"total": total, "current": page, "pageSize": pageSize}}})
}

func (c *NumberRuleController) Get(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysNumberRule
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "编号规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item})
}

func (c *NumberRuleController) Create(ctx *gin.Context) {
	var req struct {
		RuleName   string `json:"rule_name" binding:"required"`
		Prefix     string `json:"prefix"`
		DateFormat string `json:"date_format"`
		SeqFormat  string `json:"seq_format"`
		Increment  int    `json:"increment"`
		Status     int    `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var exist models.SysNumberRule
	if err := c.DB.Where("rule_name = ?", req.RuleName).First(&exist).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "规则名称已存在"})
		return
	}

	item := models.SysNumberRule{RuleName: req.RuleName, Prefix: req.Prefix, DateFormat: req.DateFormat, SeqFormat: req.SeqFormat, Increment: req.Increment, Status: req.Status, Remark: req.Remark}
	if item.Increment == 0 {
		item.Increment = 1
	}
	if item.Status == 0 {
		item.Status = 1
	}
	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "创建成功"})
}

func (c *NumberRuleController) Update(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysNumberRule
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "编号规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		RuleName   string `json:"rule_name"`
		Prefix     string `json:"prefix"`
		DateFormat string `json:"date_format"`
		SeqFormat  string `json:"seq_format"`
		Increment  int    `json:"increment"`
		Status     int    `json:"status"`
		Remark     string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.RuleName != "" {
		item.RuleName = req.RuleName
	}
	if req.Prefix != "" {
		item.Prefix = req.Prefix
	}
	if req.DateFormat != "" {
		item.DateFormat = req.DateFormat
	}
	if req.SeqFormat != "" {
		item.SeqFormat = req.SeqFormat
	}
	if req.Increment > 0 {
		item.Increment = req.Increment
	}
	if req.Remark != "" {
		item.Remark = req.Remark
	}
	if req.Status != 0 {
		item.Status = req.Status
	}

	if err := c.DB.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "更新成功"})
}

func (c *NumberRuleController) Delete(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysNumberRule
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "编号规则不存在"})
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

func (c *NumberRuleController) Generate(ctx *gin.Context) {
	var req struct {
		RuleName string `json:"rule_name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "规则名称不能为空"})
		return
	}

	var rule models.SysNumberRule
	if err := c.DB.Where("rule_name = ? AND status = 1", req.RuleName).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "编号规则不存在或已停用"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Model(&rule).Update("current_value", gorm.Expr("current_value + ?", rule.Increment)).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成编号失败"})
		return
	}

	var fresh models.SysNumberRule
	if err := c.DB.First(&fresh, rule.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成编号失败"})
		return
	}

	number := buildNumber(&fresh)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"rule_name": fresh.RuleName, "number": number, "current_value": fresh.CurrentValue}})
}

func buildNumber(rule *models.SysNumberRule) string {
	var sb strings.Builder
	if rule.Prefix != "" {
		sb.WriteString(rule.Prefix)
	}
	if rule.DateFormat != "" {
		sb.WriteString(time.Now().Format(rule.DateFormat))
	}
	if rule.SeqFormat != "" {
		seqLen := len(rule.SeqFormat)
		if seqLen == 0 {
			seqLen = 4
		}
		seq := fmt.Sprintf("%0"+fmt.Sprintf("%d", seqLen)+"d", rule.CurrentValue)
		sb.WriteString(seq)
	}
	return sb.String()
}
