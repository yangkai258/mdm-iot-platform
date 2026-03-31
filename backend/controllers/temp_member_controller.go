package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TempMemberController 临时会员控制器
type TempMemberController struct {
	DB *gorm.DB
}

func NewTempMemberController(db *gorm.DB) *TempMemberController {
	return &TempMemberController{DB: db}
}

// RegisterRoutes 注册临时会员路由
func (ctrl *TempMemberController) RegisterRoutes(rg *gin.RouterGroup) {
	tempMembers := rg.Group("/temp-members")
	{
		tempMembers.GET("", ctrl.List)
		tempMembers.POST("", ctrl.Create)
		tempMembers.GET("/:id", ctrl.GetByID)
		tempMembers.PUT("/:id", ctrl.Update)
	}
}

// List 获取临时会员列表
// GET /api/v1/temp-members
func (ctrl *TempMemberController) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.TempMemberRecord{}).Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if phone := c.Query("phone"); phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("nickname LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var members []models.TempMemberRecord
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&members)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      members,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建临时会员
// POST /api/v1/temp-members
func (ctrl *TempMemberController) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		OpenID     string    `json:"open_id"`
		Nickname   string    `json:"nickname"`
		Avatar     string    `json:"avatar"`
		Phone      string    `json:"phone"`
		StoreID    *uint     `json:"store_id"`
		ExpireTime time.Time `json:"expire_time"`
		Remark     string    `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	member := models.TempMemberRecord{
		OpenID:     req.OpenID,
		Nickname:   req.Nickname,
		Avatar:     req.Avatar,
		Phone:      req.Phone,
		StoreID:    req.StoreID,
		TenantID:   tenantID,
		ExpireTime: req.ExpireTime,
		Remark:     req.Remark,
		Status:     1,
	}

	if err := ctrl.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}

// GetByID 获取临时会员详情
// GET /api/v1/temp-members/:id
func (ctrl *TempMemberController) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var member models.TempMemberRecord
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "临时会员不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}

// Update 更新临时会员
// PUT /api/v1/temp-members/:id
func (ctrl *TempMemberController) Update(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var member models.TempMemberRecord
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "临时会员不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Nickname   string    `json:"nickname"`
		Avatar     string    `json:"avatar"`
		Phone      string    `json:"phone"`
		StoreID    *uint     `json:"store_id"`
		ExpireTime time.Time `json:"expire_time"`
		Status     int       `json:"status"`
		Remark     string    `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.StoreID != nil {
		updates["store_id"] = req.StoreID
	}
	if !req.ExpireTime.IsZero() {
		updates["expire_time"] = req.ExpireTime
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := ctrl.DB.Model(&member).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&member, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}
