package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberCardController 会员卡管理控制器
type MemberCardController struct {
	DB *gorm.DB
}

func NewMemberCardController(db *gorm.DB) *MemberCardController {
	return &MemberCardController{DB: db}
}

// RegisterRoutes 注册会员卡路由
func (ctrl *MemberCardController) RegisterRoutes(rg *gin.RouterGroup) {
	cards := rg.Group("/cards")
	{
		cards.GET("", ctrl.List)
		cards.POST("", ctrl.Create)
		cards.GET("/:id", ctrl.GetByID)
		cards.PUT("/:id", ctrl.Update)
		cards.DELETE("/:id", ctrl.Delete)
	}
}

// List 获取会员卡列表
// GET /api/v1/cards
func (ctrl *MemberCardController) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := ctrl.DB.Model(&models.MemberCardRecord{}).Where("tenant_id = ?", tenantID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if cardType := c.Query("card_type"); cardType != "" {
		query = query.Where("card_type = ?", cardType)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("card_number LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var cards []models.MemberCardRecord
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&cards)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      cards,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Create 创建会员卡
// POST /api/v1/cards
func (ctrl *MemberCardController) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req struct {
		CardNumber string     `json:"card_number"`
		CardType   int        `json:"card_type"`
		MemberID   *uint      `json:"member_id"`
		Status     int        `json:"status"`
		ExpiredAt  *time.Time `json:"expired_at"`
		Remark     string     `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.CardNumber == "" {
		// 自动生成卡号
		var maxID int64
		ctrl.DB.Model(&models.MemberCardRecord{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID)
		req.CardNumber = fmt.Sprintf("CARD%06d", maxID+1)
	}

	card := models.MemberCardRecord{
		CardNumber: req.CardNumber,
		CardType:   req.CardType,
		MemberID:   req.MemberID,
		TenantID:   tenantID,
		Status:     req.Status,
		IssuedAt:   time.Now(),
		ExpiredAt:  req.ExpiredAt,
		Remark:     req.Remark,
	}

	if err := ctrl.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": card})
}

// GetByID 获取会员卡详情
// GET /api/v1/cards/:id
func (ctrl *MemberCardController) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var card models.MemberCardRecord
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&card).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员卡不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": card})
}

// Update 更新会员卡
// PUT /api/v1/cards/:id
func (ctrl *MemberCardController) Update(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	var card models.MemberCardRecord
	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&card).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员卡不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		CardType  int        `json:"card_type"`
		MemberID  *uint      `json:"member_id"`
		Status    int        `json:"status"`
		ExpiredAt *time.Time `json:"expired_at"`
		Balance   float64    `json:"balance"`
		Points    int64      `json:"points"`
		Remark    string     `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.CardType > 0 {
		updates["card_type"] = req.CardType
	}
	if req.MemberID != nil {
		updates["member_id"] = req.MemberID
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	if req.ExpiredAt != nil {
		updates["expired_at"] = req.ExpiredAt
	}
	if req.Balance > 0 {
		updates["balance"] = req.Balance
	}
	if req.Points > 0 {
		updates["points"] = req.Points
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if err := ctrl.DB.Model(&card).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctrl.DB.First(&card, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": card})
}

// Delete 删除会员卡
// DELETE /api/v1/cards/:id
func (ctrl *MemberCardController) Delete(c *gin.Context) {
	id := c.Param("id")
	tenantID := middleware.GetTenantID(c)

	if err := ctrl.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.MemberCardRecord{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
