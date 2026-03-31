package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberServiceController 会员服务控制器
type MemberServiceController struct {
	DB *gorm.DB
}

func NewMemberServiceController(db *gorm.DB) *MemberServiceController {
	return &MemberServiceController{DB: db}
}

// RegisterRoutes 注册会员服务路由
func (ctrl *MemberServiceController) RegisterRoutes(rg *gin.RouterGroup) {
	services := rg.Group("/services")
	{
		services.GET("", ctrl.List)
		services.POST("", ctrl.Create)
		services.GET("/:id", ctrl.GetByID)
	}
}

// List 获取服务记录列表
// GET /api/v1/services
func (ctrl *MemberServiceController) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.MemberService{}).Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if serviceType := c.Query("service_type"); serviceType != "" {
		query = query.Where("service_type = ?", serviceType)
	}
	if memberID := c.Query("member_id"); memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}
	if operatorID := c.Query("operator_id"); operatorID != "" {
		query = query.Where("operator_id = ?", operatorID)
	}

	var total int64
	query.Count(&total)

	var records []models.MemberService
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建服务记录
// POST /api/v1/services
func (ctrl *MemberServiceController) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	var req struct {
		MemberID    uint   `json:"member_id" binding:"required"`
		ServiceType int    `json:"service_type"`
		Content     string `json:"content"`
		Operator    string `json:"operator"`
		Rating      int    `json:"rating"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	service := models.MemberService{
		MemberID:    req.MemberID,
		TenantID:    tenantID,
		ServiceType: req.ServiceType,
		Content:     req.Content,
		OperatorID:  &userID,
		Operator:    req.Operator,
		Rating:      req.Rating,
		Status:      1,
	}

	if err := ctrl.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": service})
}

// GetByID 获取服务详情
// GET /api/v1/services/:id
func (ctrl *MemberServiceController) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var service models.MemberService
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "服务记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": service})
}
