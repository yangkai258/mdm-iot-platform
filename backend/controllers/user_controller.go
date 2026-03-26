package controllers

import (
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserController 用户管理控制器
type UserController struct {
	DB *gorm.DB
}

// List 用户列表
// GET /api/v1/users
func (c *UserController) List(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	var users []models.SysUser
	var total int64

	query := c.DB.Model(&models.SysUser{})

	// 租户隔离
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 关键字搜索
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 角色筛选
	if roleID := ctx.Query("role_id"); roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}

	// 时间范围筛选
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Get 获取单个用户
// GET /api/v1/users/:id
func (c *UserController) Get(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

// Create 创建用户
// POST /api/v1/users
func (c *UserController) Create(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "租户ID不能为空"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		RoleID   uint   `json:"role_id"`
		Status   int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 检查用户名唯一性
	var count int64
	c.DB.Model(&models.SysUser{}).Where("username = ? AND tenant_id = ?", req.Username, tenantID).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	user := models.SysUser{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
		TenantID: tenantID,
		Status:   req.Status,
	}
	if user.Status == 0 {
		user.Status = 1
	}

	if err := c.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

// Update 更新用户
// PUT /api/v1/users/:id
func (c *UserController) Update(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		RoleID   uint   `json:"role_id"`
		Status   int    `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 如果改用户名，检查唯一性
	if req.Username != "" && req.Username != user.Username {
		var cnt int64
		c.DB.Model(&models.SysUser{}).Where("username = ? AND tenant_id = ? AND id != ?",
			req.Username, tenantID, id).Count(&cnt)
		if cnt > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名已存在"})
			return
		}
	}

	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
			return
		}
		updates["password"] = string(hashedPassword)
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RoleID > 0 {
		updates["role_id"] = req.RoleID
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := c.DB.Model(&user).Updates(updates).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
			return
		}
	}

	c.DB.First(&user, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

// Delete 删除用户
// DELETE /api/v1/users/:id
func (c *UserController) Delete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	// 检查是否是当前登录用户
	currentUserID := middleware.GetUserID(ctx)
	if currentUserID > 0 {
		userID, _ := strconv.ParseUint(id, 10, 64)
		if uint(userID) == currentUserID {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不能删除当前登录用户"})
			return
		}
	}

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// UpdateStatus 更新用户状态
// PUT /api/v1/users/:id/status
func (c *UserController) UpdateStatus(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	if err := c.DB.Model(&user).Update("status", req.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新状态失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ResetPassword 重置密码
// POST /api/v1/users/:id/reset-password
func (c *UserController) ResetPassword(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	id := ctx.Param("id")

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	if err := c.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "重置密码失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/users/me
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	tenantID := middleware.GetTenantID(ctx)

	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

// ChangePassword 修改密码
// POST /api/v1/users/change-password
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	tenantID := middleware.GetTenantID(ctx)

	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	var user models.SysUser
	if err := c.DB.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "原密码错误"})
		return
	}

	// 新密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	if err := c.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "修改密码失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// BatchDelete 批量删除用户
// POST /api/v1/users/batch-delete
func (c *UserController) BatchDelete(ctx *gin.Context) {
	tenantID := middleware.GetTenantID(ctx)
	currentUserID := middleware.GetUserID(ctx)

	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要删除的用户"})
		return
	}

	// 排除当前用户
	var idsToDelete []uint
	for _, id := range req.IDs {
		if id != currentUserID {
			idsToDelete = append(idsToDelete, id)
		}
	}

	if len(idsToDelete) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "没有可删除的用户"})
		return
	}

	result := c.DB.Where("id IN ? AND tenant_id = ?", idsToDelete, tenantID).Delete(&models.SysUser{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "批量删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"deleted": result.RowsAffected}})
}
