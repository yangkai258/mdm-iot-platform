package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============ Family/Household Controller ============

// FamilyController 家庭管理控制器
type FamilyController struct {
	DB *gorm.DB
}

// RegisterFamilyRoutes 注册家庭相关路由
func (c *FamilyController) RegisterFamilyRoutes(api *gin.RouterGroup) {
	family := api.Group("/family")
	{
		family.GET("", c.ListFamilies)
		family.POST("", c.CreateFamily)
		family.GET("/:family_id", c.GetFamily)
		family.PUT("/:family_id", c.UpdateFamily)
		family.DELETE("/:family_id", c.DeleteFamily)

		// 家庭成员管理
		family.GET("/members", c.ListMembers)
		family.POST("/members/invite", c.InviteMember)
		family.PUT("/members/:id/role", c.UpdateMemberRole)
		family.DELETE("/members/:id", c.RemoveMember)
		family.GET("/members/pending", c.ListPendingMembers)

		// 家庭相册
		family.GET("/albums", c.ListAlbums)
		family.POST("/albums", c.CreateAlbum)
		family.DELETE("/albums/:id", c.DeleteAlbum)

		// 家庭设置
		family.GET("/:family_id/settings", c.GetSettings)
		family.PUT("/:family_id/settings", c.UpdateSettings)
	}
}

// ============ 请求结构体 ============

// CreateFamilyRequest 创建家庭请求
type CreateFamilyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	Region      string `json:"region"`
}

// UpdateFamilyRequest 更新家庭请求
type UpdateFamilyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	Region      string `json:"region"`
	Status      string `json:"status"`
}

// InviteMemberRequest 邀请成员请求
type InviteMemberRequest struct {
	HouseholdID  uint   `json:"household_id" binding:"required"`
	Nickname     string `json:"nickname"`
	Relationship string `json:"relationship" binding:"required"`
	Role         string `json:"role"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Role   string `json:"role" binding:"required"`
	UserID uint   `json:"user_id"`
}

// CreateAlbumRequest 创建相册请求
type CreateAlbumRequest struct {
	HouseholdID   uint                   `json:"household_id" binding:"required"`
	PetID         *uint                  `json:"pet_id"`
	Title         string                 `json:"title" binding:"required"`
	Description   string                 `json:"description"`
	ImageURL      string                 `json:"image_url" binding:"required"`
	ThumbnailURL  string                 `json:"thumbnail_url"`
	Category      string                 `json:"category"`
	Tags          string                 `json:"tags"`
	Metadata      map[string]interface{} `json:"metadata"`
	FileSize      int64                  `json:"file_size"`
	Width         int                    `json:"width"`
	Height        int                    `json:"height"`
}

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	Mode                     string `json:"mode"`
	ChildModeEnabled         *bool  `json:"child_mode_enabled"`
	ChildAgeRange            string `json:"child_age_range"`
	ChildContentFilter       string `json:"child_content_filter"`
	ChildScreenTime          int    `json:"child_screen_time"`
	ChildAllowedActions      string `json:"child_allowed_actions"`
	ElderModeEnabled         *bool  `json:"elder_mode_enabled"`
	ElderCareLevel           string `json:"elder_care_level"`
	ElderReminderEnabled     *bool  `json:"elder_reminder_enabled"`
	ElderMedicationReminders *bool  `json:"elder_medication_reminders"`
	ElderEmergencyContact    string `json:"elder_emergency_contact"`
	PetModeEnabled           *bool  `json:"pet_mode_enabled"`
	PetInteractionLevel      string `json:"pet_interaction_level"`
	NotificationEnabled     *bool  `json:"notification_enabled"`
	NotificationSettings     string `json:"notification_settings"`
	PrivacyLevel             string `json:"privacy_level"`
}

// InteractionRequest 记录交互请求
type InteractionRequest struct {
	PetID        uint                   `json:"pet_id" binding:"required"`
	DeviceID     string                 `json:"device_id"`
	HouseholdID  *uint                  `json:"household_id"`
	Type         string                 `json:"type" binding:"required"`
	Action       string                 `json:"action"`
	Content      string                 `json:"content"`
	Duration     int                    `json:"duration"`
	Emotion      string                 `json:"emotion"`
	EmotionScore int                    `json:"emotion_score"`
	Effect       string                 `json:"effect"`
	Metadata     map[string]interface{} `json:"metadata"`
	ImageURLs    string                 `json:"image_urls"`
}

// ============ 家庭 CRUD ============

// ListFamilies 获取当前用户的家庭列表
func (c *FamilyController) ListFamilies(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")

	var households []models.Household
	query := c.DB.Where("tenant_id = ?", tenantID).
		Where("owner_id = ? OR id IN (SELECT household_id FROM household_members WHERE user_id = ? AND invite_status = 'active')", userID, userID)

	if err := query.Order("created_at DESC").Find(&households).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "查询失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": households,
		},
	})
}

// CreateFamily 创建家庭
func (c *FamilyController) CreateFamily(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")

	var req CreateFamilyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	household := models.Household{
		Name:        req.Name,
		Description: req.Description,
		AvatarURL:   req.AvatarURL,
		OwnerID:     userID,
		TenantID:    tenantID,
		Region:      req.Region,
		Status:      "active",
	}

	if err := c.DB.Create(&household).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "创建失败",
		})
		return
	}

	// 自动将创建者设为 owner 角色成员
	now := time.Now()
	member := models.HouseholdMember{
		HouseholdID:  household.ID,
		UserID:       userID,
		Nickname:     "我",
		Role:         "owner",
		Relationship: "self",
		InviteStatus: "active",
		JoinedAt:     &now,
		TenantID:     tenantID,
	}
	c.DB.Create(&member)

	// 创建默认设置
	settings := models.HouseholdSettings{
		HouseholdID: household.ID,
		TenantID:    tenantID,
		Mode:        "normal",
	}
	c.DB.Create(&settings)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    household,
	})
}

// GetFamily 获取家庭详情
func (c *FamilyController) GetFamily(ctx *gin.Context) {
	tenantID := ctx.GetString("tenant_id")
	familyID, err := strconv.ParseUint(ctx.Param("family_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "无效的家庭ID",
		})
		return
	}

	var household models.Household
	if err := c.DB.Where("id = ? AND tenant_id = ?", familyID, tenantID).First(&household).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4002,
				"message": "家庭不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "查询失败",
		})
		return
	}

	var memberCount int64
	c.DB.Model(&models.HouseholdMember{}).Where("household_id = ? AND invite_status = 'active'", familyID).Count(&memberCount)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"household":    household,
			"member_count": memberCount,
		},
	})
}

// UpdateFamily 更新家庭
func (c *FamilyController) UpdateFamily(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")
	familyID, err := strconv.ParseUint(ctx.Param("family_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "无效的家庭ID",
		})
		return
	}

	var household models.Household
	if err := c.DB.Where("id = ? AND tenant_id = ?", familyID, tenantID).First(&household).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    4002,
			"message": "家庭不存在",
		})
		return
	}

	if !c.hasHouseholdAdminPermission(household.ID, userID) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    4003,
			"message": "无权限操作",
		})
		return
	}

	var req UpdateFamilyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "参数校验失败",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Region != "" {
		updates["region"] = req.Region
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := c.DB.Model(&household).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "更新失败",
		})
		return
	}

	c.DB.First(&household, familyID)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    household,
	})
}

// DeleteFamily 删除家庭
func (c *FamilyController) DeleteFamily(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")
	familyID, err := strconv.ParseUint(ctx.Param("family_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "无效的家庭ID",
		})
		return
	}

	var household models.Household
	if err := c.DB.Where("id = ? AND tenant_id = ?", familyID, tenantID).First(&household).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    4002,
			"message": "家庭不存在",
		})
		return
	}

	if household.OwnerID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    4003,
			"message": "仅家庭所有者可删除",
		})
		return
	}

	c.DB.Delete(&household)
	c.DB.Where("household_id = ?", familyID).Delete(&models.HouseholdMember{})
	c.DB.Where("household_id = ?", familyID).Delete(&models.HouseholdSettings{})

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ============ 成员管理 ============

// ListMembers 获取家庭成员列表
func (c *FamilyController) ListMembers(ctx *gin.Context) {
	tenantID := ctx.GetString("tenant_id")

	householdID := ctx.Query("household_id")
	if householdID == "" {
		userID := ctx.GetUint("user_id")
		var member models.HouseholdMember
		if err := c.DB.Where("user_id = ? AND invite_status = 'active'", userID).First(&member).Error; err == nil {
			householdID = strconv.FormatUint(uint64(member.HouseholdID), 10)
		}
	}

	familyID, err := strconv.ParseUint(householdID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4005,
			"message": "无效的家庭ID",
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var members []models.HouseholdMember
	query := c.DB.Where("household_id = ? AND tenant_id = ? AND invite_status = 'active'", familyID, tenantID)

	var total int64
	query.Model(&models.HouseholdMember{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at ASC").Find(&members).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": members,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// InviteMember 邀请成员
func (c *FamilyController) InviteMember(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")

	var req InviteMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "参数校验失败"})
		return
	}

	var household models.Household
	if err := c.DB.Where("id = ? AND tenant_id = ?", req.HouseholdID, tenantID).First(&household).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4002, "message": "家庭不存在"})
		return
	}

	if !c.hasHouseholdAdminPermission(household.ID, userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "无邀请权限"})
		return
	}

	inviteCode := generateInviteCode()
	role := req.Role
	if role == "" {
		role = "member"
	}

	member := models.HouseholdMember{
		HouseholdID:  req.HouseholdID,
		UserID:       0,
		Nickname:     req.Nickname,
		Role:         role,
		Relationship: req.Relationship,
		InviteCode:   inviteCode,
		InviteStatus: "pending",
		TenantID:     tenantID,
	}

	if err := c.DB.Create(&member).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "邀请失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"member":      member,
			"invite_code": inviteCode,
		},
	})
}

// UpdateMemberRole 更新成员角色
func (c *FamilyController) UpdateMemberRole(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")
	memberID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的成员ID"})
		return
	}

	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "参数校验失败"})
		return
	}

	var member models.HouseholdMember
	if err := c.DB.Where("id = ? AND tenant_id = ?", memberID, tenantID).First(&member).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4002, "message": "成员不存在"})
		return
	}

	if member.UserID != userID && !c.hasHouseholdAdminPermission(member.HouseholdID, userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "无权限操作"})
		return
	}

	if member.Role == "owner" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "不能修改所有者角色"})
		return
	}

	updates := map[string]interface{}{}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.UserID > 0 {
		updates["user_id"] = req.UserID
		updates["invite_status"] = "active"
		now := time.Now()
		updates["joined_at"] = &now
	}

	if err := c.DB.Model(&member).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	c.DB.First(&member, memberID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": member})
}

// RemoveMember 移除成员
func (c *FamilyController) RemoveMember(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")
	memberID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的成员ID"})
		return
	}

	var member models.HouseholdMember
	if err := c.DB.Where("id = ? AND tenant_id = ?", memberID, tenantID).First(&member).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4002, "message": "成员不存在"})
		return
	}

	if member.Role == "owner" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "不能移除家庭所有者"})
		return
	}

	if member.UserID != userID && !c.hasHouseholdAdminPermission(member.HouseholdID, userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "无权限操作"})
		return
	}

	member.InviteStatus = "removed"
	c.DB.Save(&member)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ListPendingMembers 获取待处理邀请
func (c *FamilyController) ListPendingMembers(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")

	var households []models.Household
	c.DB.Where("owner_id = ? OR id IN (SELECT household_id FROM household_members WHERE user_id = ? AND role IN ('owner','admin') AND invite_status = 'active')", userID, userID).Find(&households)
	var householdIDs []uint
	for _, h := range households {
		householdIDs = append(householdIDs, h.ID)
	}

	var members []models.HouseholdMember
	if len(householdIDs) > 0 {
		c.DB.Where("tenant_id = ? AND invite_status = 'pending' AND household_id IN ?", tenantID, householdIDs).Find(&members)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    gin.H{"list": members},
	})
}

// ============ 家庭相册 ============

// ListAlbums 获取相册列表
func (c *FamilyController) ListAlbums(ctx *gin.Context) {
	tenantID := ctx.GetString("tenant_id")

	householdID := ctx.Query("household_id")
	familyID, err := strconv.ParseUint(householdID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的家庭ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := c.DB.Where("household_id = ? AND tenant_id = ?", familyID, tenantID)

	if category := ctx.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if petID := ctx.Query("pet_id"); petID != "" {
		query = query.Where("pet_id = ?", petID)
	}

	var total int64
	query.Model(&models.FamilyAlbum{}).Count(&total)

	var albums []models.FamilyAlbum
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&albums).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": albums,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// CreateAlbum 上传照片到相册
func (c *FamilyController) CreateAlbum(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")

	var req CreateAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "参数校验失败"})
		return
	}

	var household models.Household
	if err := c.DB.Where("id = ? AND tenant_id = ?", req.HouseholdID, tenantID).First(&household).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4002, "message": "家庭不存在"})
		return
	}

	var member models.HouseholdMember
	if err := c.DB.Where("household_id = ? AND user_id = ? AND invite_status = 'active'", req.HouseholdID, userID).First(&member).Error; err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "非家庭成员无法上传"})
		return
	}

	category := req.Category
	if category == "" {
		category = "daily"
	}

	album := models.FamilyAlbum{
		HouseholdID:  req.HouseholdID,
		PetID:        req.PetID,
		UploaderID:   userID,
		Title:        req.Title,
		Description:  req.Description,
		ImageURL:     req.ImageURL,
		ThumbnailURL: req.ThumbnailURL,
		Category:     category,
		Tags:         req.Tags,
		Metadata:     req.Metadata,
		FileSize:     req.FileSize,
		Width:        req.Width,
		Height:       req.Height,
		TenantID:     tenantID,
	}

	if err := c.DB.Create(&album).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "上传失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": album})
}

// DeleteAlbum 删除照片
func (c *FamilyController) DeleteAlbum(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")
	albumID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的照片ID"})
		return
	}

	var album models.FamilyAlbum
	if err := c.DB.Where("id = ? AND tenant_id = ?", albumID, tenantID).First(&album).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 4002, "message": "照片不存在"})
		return
	}

	if album.UploaderID != userID && !c.hasHouseholdAdminPermission(album.HouseholdID, userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "无权限删除"})
		return
	}

	c.DB.Delete(&album)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ 家庭设置 ============

// GetSettings 获取家庭设置
func (c *FamilyController) GetSettings(ctx *gin.Context) {
	tenantID := ctx.GetString("tenant_id")
	familyID, err := strconv.ParseUint(ctx.Param("family_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的家庭ID"})
		return
	}

	var settings models.HouseholdSettings
	if err := c.DB.Where("household_id = ? AND tenant_id = ?", familyID, tenantID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.HouseholdSettings{
				HouseholdID: uint(familyID),
				TenantID:    tenantID,
				Mode:        "normal",
			}
			c.DB.Create(&settings)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": settings})
}

// UpdateSettings 更新家庭设置
func (c *FamilyController) UpdateSettings(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")
	familyID, err := strconv.ParseUint(ctx.Param("family_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的家庭ID"})
		return
	}

	if !c.hasHouseholdAdminPermission(uint(familyID), userID) {
		ctx.JSON(http.StatusForbidden, gin.H{"code": 4003, "message": "无权限操作"})
		return
	}

	var settings models.HouseholdSettings
	if err := c.DB.Where("household_id = ? AND tenant_id = ?", familyID, tenantID).First(&settings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.HouseholdSettings{HouseholdID: uint(familyID), TenantID: tenantID}
			c.DB.Create(&settings)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
			return
		}
	}

	var req UpdateSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "参数校验失败"})
		return
	}

	updates := map[string]interface{}{}
	if req.Mode != "" {
		updates["mode"] = req.Mode
	}
	if req.ChildModeEnabled != nil {
		updates["child_mode_enabled"] = *req.ChildModeEnabled
	}
	if req.ChildAgeRange != "" {
		updates["child_age_range"] = req.ChildAgeRange
	}
	if req.ChildContentFilter != "" {
		updates["child_content_filter"] = req.ChildContentFilter
	}
	if req.ChildScreenTime > 0 {
		updates["child_screen_time"] = req.ChildScreenTime
	}
	if req.ChildAllowedActions != "" {
		updates["child_allowed_actions"] = req.ChildAllowedActions
	}
	if req.ElderModeEnabled != nil {
		updates["elder_mode_enabled"] = *req.ElderModeEnabled
	}
	if req.ElderCareLevel != "" {
		updates["elder_care_level"] = req.ElderCareLevel
	}
	if req.ElderReminderEnabled != nil {
		updates["elder_reminder_enabled"] = *req.ElderReminderEnabled
	}
	if req.ElderMedicationReminders != nil {
		updates["elder_medication_reminders"] = *req.ElderMedicationReminders
	}
	if req.ElderEmergencyContact != "" {
		updates["elder_emergency_contact"] = req.ElderEmergencyContact
	}
	if req.PetModeEnabled != nil {
		updates["pet_mode_enabled"] = *req.PetModeEnabled
	}
	if req.PetInteractionLevel != "" {
		updates["pet_interaction_level"] = req.PetInteractionLevel
	}
	if req.NotificationEnabled != nil {
		updates["notification_enabled"] = *req.NotificationEnabled
	}
	if req.NotificationSettings != "" {
		updates["notification_settings"] = req.NotificationSettings
	}
	if req.PrivacyLevel != "" {
		updates["privacy_level"] = req.PrivacyLevel
	}

	if err := c.DB.Model(&settings).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	c.DB.First(&settings, settings.ID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": settings})
}

// ============ 多用户交互 ============

// InteractionController 交互记录控制器
type InteractionController struct {
	DB *gorm.DB
}

// RegisterInteractionRoutes 注册交互路由
func (c *InteractionController) RegisterInteractionRoutes(api *gin.RouterGroup) {
	api.POST("/interactions", c.RecordInteraction)
	api.GET("/interactions/:pet_id", c.GetPetInteractions)
}

// RecordInteraction 记录交互
func (c *InteractionController) RecordInteraction(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetString("tenant_id")

	var req InteractionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "参数校验失败"})
		return
	}

	interaction := models.Interaction{
		PetID:        req.PetID,
		DeviceID:     req.DeviceID,
		UserID:       userID,
		HouseholdID:  req.HouseholdID,
		Type:         req.Type,
		Action:       req.Action,
		Content:      req.Content,
		Duration:     req.Duration,
		Emotion:      req.Emotion,
		EmotionScore: req.EmotionScore,
		Effect:       req.Effect,
		Metadata:     req.Metadata,
		ImageURLs:    req.ImageURLs,
		TenantID:     tenantID,
	}

	if err := c.DB.Create(&interaction).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "记录失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": interaction})
}

// GetPetInteractions 获取宠物交互记录
func (c *InteractionController) GetPetInteractions(ctx *gin.Context) {
	tenantID := ctx.GetString("tenant_id")
	petID, err := strconv.ParseUint(ctx.Param("pet_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "无效的宠物ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := c.DB.Where("pet_id = ? AND tenant_id = ?", petID, tenantID)

	if interactionType := ctx.Query("type"); interactionType != "" {
		query = query.Where("type = ?", interactionType)
	}

	var total int64
	query.Model(&models.Interaction{}).Count(&total)

	var interactions []models.Interaction
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&interactions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": interactions,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (int(total) + pageSize - 1) / pageSize,
			},
		},
	})
}

// ============ 辅助函数 ============

// hasHouseholdAdminPermission 检查用户是否有家庭管理权限
func (c *FamilyController) hasHouseholdAdminPermission(householdID, userID uint) bool {
	var member models.HouseholdMember
	err := c.DB.Where("household_id = ? AND user_id = ? AND invite_status = 'active' AND role IN ('owner', 'admin')", householdID, userID).First(&member).Error
	return err == nil
}

// generateInviteCode 生成随机邀请码
func generateInviteCode() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
