package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetProfileController 宠物配置控制器
type PetProfileController struct {
	DB *gorm.DB
}

// GetProfile 获取宠物配置
func (c *PetProfileController) GetProfile(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var profile models.PetProfile
	if err := c.DB.Where("device_id = ?", deviceID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认配置
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"device_id":        deviceID,
					"pet_name":         "Mimi",
					"personality":      "lively",
					"interaction_freq": "medium",
					"dnd_start_time":   "23:00",
					"dnd_end_time":     "08:00",
					"custom_rules":     map[string]interface{}{},
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "服务器内部错误",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    profile,
	})
}

// UpdateProfileRequest 更新宠物配置请求
type UpdateProfileRequest struct {
	PetName         string                 `json:"pet_name"`
	Personality     string                 `json:"personality"`
	InteractionFreq string                 `json:"interaction_freq"`
	DNDStartTime    string                 `json:"dnd_start_time"`
	DNDEndTime      string                 `json:"dnd_end_time"`
	CustomRules     map[string]interface{} `json:"custom_rules"`
}

// UpdateProfile 更新宠物配置
func (c *PetProfileController) UpdateProfile(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var profile models.PetProfile
	result := c.DB.Where("device_id = ?", deviceID).First(&profile)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新配置
		profile = models.PetProfile{
			DeviceID:        deviceID,
			PetName:        req.PetName,
			Personality:    req.Personality,
			InteractionFreq: req.InteractionFreq,
			DNDStartTime:   req.DNDStartTime,
			DNDEndTime:     req.DNDEndTime,
			CustomRules:    req.CustomRules,
		}
		if err := c.DB.Create(&profile).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "创建配置失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "服务器内部错误",
			"error_code": "ERR_INTERNAL",
		})
		return
	} else {
		// 更新现有配置
		if req.PetName != "" {
			profile.PetName = req.PetName
		}
		if req.Personality != "" {
			profile.Personality = req.Personality
		}
		if req.InteractionFreq != "" {
			profile.InteractionFreq = req.InteractionFreq
		}
		if req.DNDStartTime != "" {
			profile.DNDStartTime = req.DNDStartTime
		}
		if req.DNDEndTime != "" {
			profile.DNDEndTime = req.DNDEndTime
		}
		if req.CustomRules != nil {
			profile.CustomRules = req.CustomRules
		}
		profile.UpdatedAt = time.Now()

		if err := c.DB.Save(&profile).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "更新配置失败",
				"error_code": "ERR_INTERNAL",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    profile,
	})
}
