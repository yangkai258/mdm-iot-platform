package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberProfileController 会员画像控制器
type MemberProfileController struct {
	DB *gorm.DB
}

// NewMemberProfileController 创建控制器
func NewMemberProfileController(db *gorm.DB) *MemberProfileController {
	return &MemberProfileController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *MemberProfileController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/members/:id/profile", ctrl.GetProfile)
	rg.POST("/members/:id/profile/generate", ctrl.GenerateProfile)
	rg.GET("/members/:id/insights", ctrl.GetInsights)
	rg.GET("/members/profiles/search", ctrl.SearchProfiles)
}

// GetProfile 获取会员360度画像
func (ctrl *MemberProfileController) GetProfile(c *gin.Context) {
	memberID := c.Param("id")

	var profile models.Member360Profile
	if err := ctrl.DB.Where("member_id = ?", memberID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员画像不存在，请先生成"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取画像失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": profile})
}

// GenerateProfile 生成会员画像
func (ctrl *MemberProfileController) GenerateProfile(c *gin.Context) {
	memberID := c.Param("id")

	// 获取会员信息
	var member models.Member
	if err := ctrl.DB.Where("id = ?", memberID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会员不存在"})
		return
	}

	// 计算消费画像 - 使用 member_orders 表
	var totalSpend, avgOrder float64
	var totalOrders int64
	ctrl.DB.Model(&models.MemberOrder{}).Where("member_id = ?", memberID).
		Select("COALESCE(SUM(total_amount),0), COUNT(*), COALESCE(AVG(total_amount),0)").
		Row().Scan(&totalSpend, &totalOrders, &avgOrder)

	// 计算行为数据 - 使用 activity_logs
	var loginCount int64
	ctrl.DB.Model(&models.ActivityLog{}).
		Where("member_id = ? AND created_at > ?", memberID, time.Now().AddDate(0,0,-30)).
		Count(&loginCount)

	// 获取宠物互动数据
	var petEngagement float64
	var petCount int64
	ctrl.DB.Model(&models.Pet{}).Where("owner_id = ?", memberID).Count(&petCount)
	if petCount > 0 {
		// 直接查询 interaction_records 表
		var interactionCount int64
		ctrl.DB.Raw("SELECT COUNT(*) FROM interaction_records WHERE pet_id IN (SELECT id FROM pets WHERE owner_id = ?)", memberID).Scan(&interactionCount)
		petEngagement = float64(interactionCount) / float64(petCount) / 10.0
		if petEngagement > 100 {
			petEngagement = 100
		}
	}

	// 计算流失风险
	churnRisk := 0.3
	if totalOrders == 0 {
		churnRisk = 0.7
	} else {
		// 检查最近订单时间
		var lastOrderTime time.Time
		ctrl.DB.Model(&models.MemberOrder{}).
			Where("member_id = ?", memberID).
			Select("COALESCE(MAX(created_at), '1970-01-01')").
			Row().Scan(&lastOrderTime)
		if !lastOrderTime.IsZero() {
			daysSinceLastOrder := time.Since(lastOrderTime).Hours() / 24
			if daysSinceLastOrder > 90 {
				churnRisk = 0.8
			} else if daysSinceLastOrder > 60 {
				churnRisk = 0.5
			}
		}
	}

	profile := models.Member360Profile{
		MemberID:          memberID,
		MemberName:        member.MemberName,
		Gender:            member.Gender,
		Location:          "", // Member 表无此字段
		TotalSpend:        totalSpend,
		AvgOrderValue:     avgOrder,
		TotalOrders:       int(totalOrders),
		LoginFrequency:    int(loginCount),
		PetEngagement:     petEngagement,
		MemberLevel:      strconv.Itoa(member.MemberLevel),
		MemberSince:       &member.CreatedAt,
		ChurnRisk:         churnRisk,
		LTV:               totalSpend * 1.5,
		GeneratedAt:        time.Now(),
	}

	// 保存或更新
	existing := models.Member360Profile{}
	if err := ctrl.DB.Where("member_id = ?", memberID).First(&existing).Error; err == nil {
		profile.ID = existing.ID
		ctrl.DB.Save(&profile)
	} else {
		ctrl.DB.Create(&profile)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "画像生成成功", "data": profile})
}

// GetInsights 获取会员洞察
func (ctrl *MemberProfileController) GetInsights(c *gin.Context) {
	memberID := c.Param("id")

	var insights []models.MemberPortraitInsight
	ctrl.DB.Where("member_id = ? AND (expires_at IS NULL OR expires_at > ?)", memberID, time.Now()).
		Order("confidence DESC").Find(&insights)

	// 生成默认洞察
	if len(insights) == 0 {
		insights = ctrl.generateDefaultInsights(memberID)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": insights})
}

func (ctrl *MemberProfileController) generateDefaultInsights(memberID string) []models.MemberPortraitInsight {
	var profile models.Member360Profile
	if err := ctrl.DB.Where("member_id = ?", memberID).First(&profile).Error; err != nil {
		return []models.MemberPortraitInsight{}
	}

	insights := []models.MemberPortraitInsight{}

	if profile.TotalOrders > 10 {
		insights = append(insights, models.MemberPortraitInsight{
			MemberID:    memberID,
			InsightType: "consumption",
			Title:       "高价值会员",
			Description: "历史订单数超过10单，是核心用户群体",
			Confidence:  0.95,
			Action:      "建议提供专属VIP服务和优先体验新功能",
		})
	}

	if profile.ChurnRisk > 0.6 {
		insights = append(insights, models.MemberPortraitInsight{
			MemberID:    memberID,
			InsightType: "risk",
			Title:       "流失风险预警",
			Description: "该会员已有较长时间未互动，流失风险较高",
			Confidence:  profile.ChurnRisk,
			Action:      "建议推送专属优惠券或召回活动",
		})
	}

	if profile.PetEngagement > 70 {
		insights = append(insights, models.MemberPortraitInsight{
			MemberID:    memberID,
			InsightType: "habit",
			Title:       "宠物互动达人",
			Description: "宠物互动频率很高，对宠物功能非常活跃",
			Confidence:  0.88,
			Action:      "推荐宠物进阶功能和新宠物配件",
		})
	}

	// 保存洞察
	for i := range insights {
		expireAt := time.Now().AddDate(0, 0, 7)
		insights[i].ExpiresAt = &expireAt
		ctrl.DB.Create(&insights[i])
	}

	return insights
}

// SearchProfiles 搜索会员画像
func (ctrl *MemberProfileController) SearchProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
	churnRisk := c.Query("churn_risk")
	keyword := c.Query("keyword")

	query := ctrl.DB.Model(&models.Member360Profile{})

	if level != "" {
		query = query.Where("member_level = ?", level)
	}
	if churnRisk != "" {
		query = query.Where("churn_risk >= ?", churnRisk)
	}
	if keyword != "" {
		query = query.Where("member_name LIKE ? OR member_id LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	var list []models.Member360Profile
	query.Count(&total)

	query.Order("ltv DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
