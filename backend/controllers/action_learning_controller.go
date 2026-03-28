package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActionLearningController 动作模仿学习进度控制器
type ActionLearningController struct {
	DB *gorm.DB
}

// NewActionLearningController 创建控制器
func NewActionLearningController(db *gorm.DB) *ActionLearningController {
	return &ActionLearningController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *ActionLearningController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/pets/:id/action-progress", ctrl.GetProgressList)
	rg.GET("/pets/:id/action-progress/:action_id", ctrl.GetProgress)
	rg.POST("/pets/:id/action-progress/:action_id/practice", ctrl.RecordPractice)
	rg.GET("/pets/:id/action-progress/:action_id/sessions", ctrl.GetSessions)
	rg.GET("/members/:id/action-summary", ctrl.GetMemberSummary)
}

// GetProgressList 获取宠物动作学习进度列表
func (ctrl *ActionLearningController) GetProgressList(c *gin.Context) {
	petID := c.Param("id")

	var list []models.ActionLearningProgress
	ctrl.DB.Where("pet_id = ?", petID).Order("mastery_rate DESC, level DESC").Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":        list,
			"total":       len(list),
			"mastered":    countByStatus(list, "mastered"),
			"learning":    countByStatus(list, "learning"),
		},
	})
}

func countByStatus(list []models.ActionLearningProgress, status string) int {
	count := 0
	for _, p := range list {
		if p.Status == status {
			count++
		}
	}
	return count
}

// GetProgress 获取单个动作学习进度
func (ctrl *ActionLearningController) GetProgress(c *gin.Context) {
	petID := c.Param("id")
	actionID := c.Param("action_id")

	var progress models.ActionLearningProgress
	if err := ctrl.DB.Where("pet_id = ? AND action_id = ?", petID, actionID).First(&progress).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "学习进度不存在",
		})
		return
	}

	// 获取最近练习记录
	var sessions []models.ActionLearningSession
	ctrl.DB.Where("progress_id = ?", progress.ID).Order("session_date DESC").Limit(10).Find(&sessions)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"progress": progress,
			"recent_sessions": sessions,
		},
	})
}

// RecordPractice 记录练习
func (ctrl *ActionLearningController) RecordPractice(c *gin.Context) {
	petID := c.Param("id")
	actionID := c.Param("action_id")

	var req struct {
		Duration    int     `json:"duration"`
		Accuracy    float64 `json:"accuracy"`
		IsSuccess   bool    `json:"is_success"`
		Score       int     `json:"score"`
		Feedback    string  `json:"feedback"`
		VideoURL    string  `json:"video_url"`
		ScreenshotURL string `json:"screenshot_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 获取或创建进度
	var progress models.ActionLearningProgress
	exists := ctrl.DB.Where("pet_id = ? AND action_id = ?", petID, actionID).First(&progress).Error

	if exists == gorm.ErrRecordNotFound {
		progress = models.ActionLearningProgress{
			PetID:      petID,
			ActionID:   parseUintAction(actionID),
			Level:      1,
			Exp:        0,
			ExpToNextLevel: 100,
			Status:     "learning",
		}
		ctrl.DB.Create(&progress)
	}

	// 计算获得经验
	expGained := calculateExp(req.Duration, req.Accuracy, req.IsSuccess)

	// 创建练习记录
	session := models.ActionLearningSession{
		ProgressID:   progress.ID,
		PetID:        petID,
		ActionID:     parseUintAction(actionID),
		Duration:     req.Duration,
		Accuracy:     req.Accuracy,
		IsSuccess:    req.IsSuccess,
		Score:        req.Score,
		ExpGained:    expGained,
		Feedback:     req.Feedback,
		VideoURL:     req.VideoURL,
		ScreenshotURL: req.ScreenshotURL,
		SessionDate:  time.Now(),
	}
	ctrl.DB.Create(&session)

	// 更新进度
	progress.Exp += expGained
	progress.PracticeCount++
	if req.IsSuccess {
		progress.SuccessCount++
	}
	progress.TotalDuration += req.Duration
	progress.AvgAccuracy = (progress.AvgAccuracy*float64(progress.PracticeCount-1) + req.Accuracy) / float64(progress.PracticeCount)
	now := time.Now()
	progress.LastPracticeAt = &now

	// 检查升级
	if progress.Exp >= progress.ExpToNextLevel && progress.Level < 10 {
		progress.Level++
		progress.Exp = progress.Exp - progress.ExpToNextLevel
		progress.ExpToNextLevel = int(float64(progress.ExpToNextLevel) * 1.2) // 下一级需要更多经验
	}

	// 检查精通
	if progress.MasteryRate >= 90 && progress.Status != "mastered" {
		progress.Status = "mastered"
		progress.MasteredAt = &now
	} else if progress.MasteryRate > 0 {
		progress.Status = "practicing"
	}

	// 更新掌握度
	progress.MasteryRate = calculateMasteryRate(progress.SuccessCount, progress.PracticeCount)

	ctrl.DB.Save(&progress)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "练习记录成功",
		"data": gin.H{
			"exp_gained":  expGained,
			"level":       progress.Level,
			"exp":         progress.Exp,
			"mastery_rate": progress.MasteryRate,
			"status":      progress.Status,
		},
	})
}

// calculateExp 计算获得经验
func calculateExp(duration int, accuracy float64, success bool) int {
	base := duration / 10
	if base < 1 {
		base = 1
	}
	accBonus := int(accuracy / 10)
	if success {
		return base + accBonus + 5
	}
	return base
}

// calculateMasteryRate 计算掌握度
func calculateMasteryRate(success, total int) float64 {
	if total == 0 {
		return 0
	}
	rate := float64(success) / float64(total) * 100
	return rate
}

// GetSessions 获取练习记录
func (ctrl *ActionLearningController) GetSessions(c *gin.Context) {
	petID := c.Param("id")
	actionID := c.Param("action_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var progress models.ActionLearningProgress
	if err := ctrl.DB.Where("pet_id = ? AND action_id = ?", petID, actionID).First(&progress).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "进度不存在"})
		return
	}

	var total int64
	var sessions []models.ActionLearningSession
	ctrl.DB.Model(&models.ActionLearningSession{}).Where("progress_id = ?", progress.ID).Count(&total)

	ctrl.DB.Where("progress_id = ?", progress.ID).Order("session_date DESC").
		Offset((page-1)*pageSize).Limit(pageSize).Find(&sessions)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      sessions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMemberSummary 获取会员动作学习汇总
func (ctrl *ActionLearningController) GetMemberSummary(c *gin.Context) {
	memberID := c.Param("id")

	var totalPets int64
	ctrl.DB.Model(&models.Pet{}).Where("owner_id = ?", memberID).Count(&totalPets)

	var progress []models.ActionLearningProgress
	ctrl.DB.Joins("JOIN pets ON pets.id = action_learning_progress.pet_id").
		Where("pets.owner_id = ?", memberID).Find(&progress)

	mastered := 0
	learning := 0
	totalExp := 0
	for _, p := range progress {
		if p.Status == "mastered" {
			mastered++
		} else {
			learning++
		}
		totalExp += p.Exp
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total_pets":    totalPets,
			"mastered_count": mastered,
			"learning_count": learning,
			"total_exp":     totalExp,
		},
	})
}

func parseUintAction(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint(v)
}
