package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FlowProcessController BPMN流程控制器
type FlowProcessController struct {
	DB *gorm.DB
}

func NewFlowProcessController(db *gorm.DB) *FlowProcessController {
	return &FlowProcessController{DB: db}
}

// RegisterRoutes 注册流程路由
func (ctrl *FlowProcessController) RegisterRoutes(rg *gin.RouterGroup) {
	processes := rg.Group("/flow/processes")
	{
		processes.GET("", ctrl.List)
		processes.POST("", ctrl.Create)
		processes.GET("/:id", ctrl.GetByID)
	}
}

// List 获取流程列表
// GET /api/v1/flow/processes
func (ctrl *FlowProcessController) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.FlowProcess{}).Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("process_name LIKE ? OR process_key LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var processes []models.FlowProcess
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&processes)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      processes,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建流程
// POST /api/v1/flow/processes
func (ctrl *FlowProcessController) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		ProcessName string `json:"process_name" binding:"required"`
		ProcessKey  string `json:"process_key" binding:"required"`
		BpmnXML     string `json:"bpmn_xml"`
		Status      int    `json:"status"`
		Description string `json:"description"`
		Category    string `json:"category"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 检查 process_key 是否已存在
	var existing models.FlowProcess
	if err := ctrl.DB.Where("process_key = ? AND tenant_id = ?", req.ProcessKey, tenantID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "流程标识已存在"})
		return
	}

	process := models.FlowProcess{
		ProcessName: req.ProcessName,
		ProcessKey:  req.ProcessKey,
		BpmnXML:     req.BpmnXML,
		TenantID:    tenantID,
		Status:      req.Status,
		Description: req.Description,
		Category:    req.Category,
		Version:     1,
	}

	if err := ctrl.DB.Create(&process).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": process})
}

// GetByID 获取流程详情
// GET /api/v1/flow/processes/:id
func (ctrl *FlowProcessController) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var process models.FlowProcess
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&process).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "流程不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": process})
}
