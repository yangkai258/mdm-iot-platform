package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HouseholdPetCtrl 家庭多宠物管理控制器
type HouseholdPetCtrl struct {
	DB *gorm.DB
}

// RegisterHouseholdRoutes 注册家庭宠物管理路由
func (h *HouseholdPetCtrl) RegisterHouseholdRoutes(r *gin.RouterGroup) {
	r.GET("/household/pets", h.ListHouseholdPets)
	r.POST("/household/pets/invite", h.InviteMember)
	r.GET("/household/members", h.ListMembers)
	r.POST("/household/members", h.AddMember)
	r.PUT("/household/members/:id", h.UpdateMember)
	r.DELETE("/household/members/:id", h.RemoveMember)
}

// ListHouseholdPets 获取当前用户家庭的所有宠物
func (h *HouseholdPetCtrl) ListHouseholdPets(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	// 查找该用户所属的所有家庭
	var memberRecords []models.HouseholdMember
	if err := h.DB.Where("user_id = ? AND tenant_id = ? AND status = ?", userID, tenantID, "active").Find(&memberRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if len(memberRecords) == 0 {
		// 用户不属于任何家庭，只返回自己拥有的宠物
		var pets []models.Pet
		h.DB.Where("owner_id = ? AND tenant_id = ? AND status != ?", userID, tenantID, "deceased").
			Order("created_at DESC").Find(&pets)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": pets})
		return
	}

	// 收集所有家庭ID
	householdIDs := make([]uint, len(memberRecords))
	for i, m := range memberRecords {
		householdIDs[i] = m.HouseholdID
	}

	// 查询家庭所有宠物 + 用户自己的宠物
	var pets []models.Pet
	h.DB.Where("(household_id IN ? OR owner_id = ?) AND tenant_id = ? AND status != ?", householdIDs, userID, tenantID, "deceased").
		Order("created_at DESC").Find(&pets)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": pets})
}

// InviteMember 邀请家庭成员添加宠物（生成邀请码）
func (h *HouseholdPetCtrl) InviteMember(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		HouseholdID uint   `json:"household_id" binding:"required"`
		Email       string `json:"email"`
		Role        string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "household_id 不能为空"})
		return
	}

	// 验证用户是该家庭的所有者或成员
	var member models.HouseholdMember
	err := h.DB.Where("household_id = ? AND user_id = ? AND tenant_id = ?", input.HouseholdID, userID, tenantID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "您不是该家庭成员"})
		return
	}

	role := input.Role
	if role == "" {
		role = "member"
	}

	inviteCode := generateInviteCode()

	invite := models.HouseholdMember{
		HouseholdID:  input.HouseholdID,
		UserID:       0, // 尚未加入
		Role:         role,
		Status:       "pending",
		TenantID:     tenantID,
		InviteCode:   inviteCode,
		InvitedEmail: input.Email,
		InvitedBy:    &userID,
	}

	if err := h.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "生成邀请失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "邀请成功",
		"data": gin.H{
			"invite_code": inviteCode,
			"invite":      invite,
		},
	})
}

// ListMembers 获取家庭成员列表
func (h *HouseholdPetCtrl) ListMembers(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	// 获取用户所属的所有家庭
	var memberRecords []models.HouseholdMember
	if err := h.DB.Where("user_id = ? AND tenant_id = ? AND status = ?", userID, tenantID, "active").Find(&memberRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if len(memberRecords) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": []interface{}{}})
		return
	}

	householdIDs := make([]uint, len(memberRecords))
	for i, m := range memberRecords {
		householdIDs[i] = m.HouseholdID
	}

	var members []models.HouseholdMember
	h.DB.Where("household_id IN ? AND tenant_id = ?", householdIDs, tenantID).
		Order("created_at ASC").Find(&members)

	// 填充用户信息（如果有的话）
	type MemberWithUser struct {
		models.HouseholdMember
		Username string `json:"username,omitempty"`
	}

	result := make([]MemberWithUser, len(members))
	for i, m := range members {
		result[i] = MemberWithUser{HouseholdMember: m}
		if m.UserID > 0 {
			var user models.SysUser
			if err := h.DB.First(&user, m.UserID).Error; err == nil {
				result[i].Username = user.Username
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// AddMember 添加家庭成员（通过邀请码加入或直接添加）
func (h *HouseholdPetCtrl) AddMember(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		InviteCode string `json:"invite_code"` // 优先通过邀请码加入
		HouseholdID uint  `json:"household_id"`
		Email       string `json:"email"`
		Role        string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	role := input.Role
	if role == "" {
		role = "member"
	}

	// 如果有邀请码，验证并加入
	if input.InviteCode != "" {
		var invite models.HouseholdMember
		err := h.DB.Where("invite_code = ? AND tenant_id = ? AND status = ?", input.InviteCode, tenantID, "pending").First(&invite).Error
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "邀请码无效或已过期"})
			return
		}

		// 更新为正式成员
		invite.UserID = userID
		invite.Status = "active"
		invite.InviteCode = ""
		joinedAt := invite.CreatedAt
		invite.JoinedAt = &joinedAt

		if err := h.DB.Save(&invite).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "加入失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "加入成功", "data": invite})
		return
	}

	// 直接添加成员（需要是家庭成员才能操作）
	if input.HouseholdID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "household_id 不能为空"})
		return
	}

	// 验证操作者是家庭成员
	var self models.HouseholdMember
	err := h.DB.Where("household_id = ? AND user_id = ? AND tenant_id = ? AND status = ?", input.HouseholdID, userID, tenantID, "active").First(&self).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "您不是该家庭成员"})
		return
	}

	member := models.HouseholdMember{
		HouseholdID:  input.HouseholdID,
		UserID:       userID,
		Role:         role,
		Status:       "active",
		TenantID:     tenantID,
		InvitedEmail: input.Email,
		InvitedBy:    &userID,
	}

	if err := h.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "添加成员失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "添加成功", "data": member})
}

// UpdateMember 更新成员权限
func (h *HouseholdPetCtrl) UpdateMember(c *gin.Context) {
	memberIDStr := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的成员ID"})
		return
	}

	var member models.HouseholdMember
	err = h.DB.Where("id = ? AND tenant_id = ?", uint(memberID), tenantID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "成员不存在"})
		return
	}

	// 验证操作者权限（必须是 owner）
	var self models.HouseholdMember
	err = h.DB.Where("household_id = ? AND user_id = ? AND tenant_id = ? AND role = ?", member.HouseholdID, userID, tenantID, "owner").First(&self).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "只有家庭所有者可以修改成员权限"})
		return
	}

	var input struct {
		Role   string `json:"role"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if input.Role != "" {
		updates["role"] = input.Role
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}

	if len(updates) > 0 {
		h.DB.Model(&member).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": member})
}

// RemoveMember 移除家庭成员
func (h *HouseholdPetCtrl) RemoveMember(c *gin.Context) {
	memberIDStr := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的成员ID"})
		return
	}

	var member models.HouseholdMember
	err = h.DB.Where("id = ? AND tenant_id = ?", uint(memberID), tenantID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "成员不存在"})
		return
	}

	// 验证权限：本人可以退出，owner 可以移除其他成员
	if member.UserID != userID {
		var self models.HouseholdMember
		err = h.DB.Where("household_id = ? AND user_id = ? AND tenant_id = ? AND role = ?", member.HouseholdID, userID, tenantID, "owner").First(&self).Error
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权移除该成员"})
			return
		}
	}

	// 软删除
	if err := h.DB.Delete(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "移除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "成员已移除"})
}
