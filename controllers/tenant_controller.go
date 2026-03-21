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

// ==================== 请求/响应结构 ====================

type TenantListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	List     []models.Tenant  `json:"list"`
}

type UpdateTenantRequest struct {
	Name         string `json:"name"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	LogoURL      string `json:"logo_url"`
	Domain       string `json:"domain"`
}

type ExtendRequest struct {
	ExtendDays int `json:"extend_days" binding:"required,min=1"`
}

type CreateTenantRequest struct {
	TenantCode    string `json:"tenant_code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	Plan         string `json:"plan"` // default: "free"
	ExpiresAt    string `json:"expires_at"` // RFC3339 format, optional
}

type ChangePlanRequest struct {
	PlanID        uint   `json:"plan_id" binding:"required"`
	EffectiveType string `json:"effective_type"` // immediate | end_of_cycle
}

// ==================== TenantController ====================

type TenantController struct {
	DB *gorm.DB
}

// RegisterTenantRoutes 注册租户管理路由（超管）
func (tc *TenantController) RegisterTenantRoutes(r *gin.RouterGroup) {
	// 租户 CRUD
	r.POST("/tenants", tc.CreateTenant)   // POST /api/v1/admin/tenants
	r.GET("/tenants", tc.ListTenants)
	r.GET("/tenants/:id", tc.GetTenant)
	r.PUT("/tenants/:id", tc.UpdateTenant)
	r.DELETE("/tenants/:id", tc.DeleteTenant)
	r.PUT("/tenants/:id/suspend", tc.SuspendTenant)
	r.PUT("/tenants/:id/activate", tc.ActivateTenant)
	r.PUT("/tenants/:id/extend", tc.ExtendTenant)
	r.PUT("/tenants/:id/upgrade", tc.ChangePlan)       // 套餐升级/降级别名
	r.POST("/tenants/:id/change-plan", tc.ChangePlan) // 原 change-plan 路由

	// 套餐管理
	r.GET("/plans", tc.ListPlans)
}

// RegisterTenantAPIRoutes 注册租户 CRUD API（/api/v1/tenants 路径，超管可用）
func (tc *TenantController) RegisterTenantAPIRoutes(r *gin.RouterGroup) {
	r.POST("/tenants", tc.CreateTenant)
	r.GET("/tenants", tc.ListTenants)
	r.GET("/tenants/:id", tc.GetTenant)
	r.PUT("/tenants/:id", tc.UpdateTenant)
	r.DELETE("/tenants/:id", tc.DeleteTenant)
	r.PUT("/tenants/:id/upgrade", tc.ChangePlan) // PUT /api/v1/tenants/:id/upgrade
}

// CreateTenant 创建租户（超管）
func (tc *TenantController) CreateTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PARAMS", "message": "参数错误: " + err.Error()})
		return
	}

	// 检查 tenant_code 是否已存在
	var existing models.Tenant
	if err := tc.DB.Where("tenant_code = ?", req.TenantCode).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": "DUPLICATE_CODE", "message": "租户编码已存在"})
		return
	}

	planCode := req.Plan
	if planCode == "" {
		planCode = "free"
	}

	// 验证套餐是否存在
	var plan models.Plan
	if err := tc.DB.Where("plan_code = ? AND is_active = ?", planCode, true).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PLAN", "message": "套餐不存在或未激活"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	// 初始化租户配额
	defaultQuota := models.TenantQuota{
		UserCount:   plan.UserQuota,
		DeviceCount: plan.DeviceQuota,
		DeptCount:   plan.DeptQuota,
		StoreCount:  plan.StoreQuota,
	}

	// 解析过期时间
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, req.ExpiresAt); err == nil {
			expiresAt = &t
		}
	} else {
		// 默认 30 天试用期
		t := time.Now().AddDate(0, 0, 30)
		expiresAt = &t
	}

	tenant := models.Tenant{
		TenantCode:    req.TenantCode,
		Name:          req.Name,
		ContactName:   req.ContactName,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
		Plan:          planCode,
		Status:        "active",
		ExpiresAt:     expiresAt,
	}

	if err := tc.DB.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "创建租户失败: " + err.Error()})
		return
	}

	// 创建租户配额记录
	defaultQuota.TenantID = tenant.ID
	if err := tc.DB.Create(&defaultQuota).Error; err != nil {
		// 配额创建失败不影响租户创建，只记录日志
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "租户创建成功但配额初始化失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "租户创建成功",
		"data": gin.H{
			"tenant": tenant,
			"quota":  defaultQuota,
		},
	})
}

// ListTenants 获取租户列表（超管）
func (tc *TenantController) ListTenants(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	var tenants []models.Tenant
	var total int64

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	query := tc.DB.Model(&models.Tenant{})

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// 套餐筛选
	if plan := c.Query("plan"); plan != "" {
		query = query.Where("plan = ?", plan)
	}
	// 关键词搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name ILIKE ? OR contact_name ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tenants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"list":      tenants,
		},
	})
}

// GetTenant 获取租户详情
func (tc *TenantController) GetTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	var quota models.TenantQuota
	tc.DB.Where("tenant_id = ?", id).First(&quota)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"tenant": tenant,
			"quota":  quota,
		},
	})
}

// UpdateTenant 更新租户信息
func (tc *TenantController) UpdateTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PARAMS", "message": "参数错误"})
		return
	}

	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.ContactEmail != "" {
		updates["contact_email"] = req.ContactEmail
	}
	if req.LogoURL != "" {
		updates["logo_url"] = req.LogoURL
	}
	if req.Domain != "" {
		updates["domain"] = req.Domain
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PARAMS", "message": "无更新字段"})
		return
	}

	if err := tc.DB.Model(&tenant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

// DeleteTenant 删除租户（仅 expired 状态可删除）
func (tc *TenantController) DeleteTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	if tenant.Status != "expired" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_STATUS", "message": "仅已到期的租户可删除"})
		return
	}

	if err := tc.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("tenant_id = ?", id).Delete(&models.TenantQuota{})
		return tx.Where("id = ?", id).Delete(&models.Tenant{}).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// SuspendTenant 禁用租户
func (tc *TenantController) SuspendTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	if tenant.Status == "suspended" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "ALREADY_SUSPENDED", "message": "租户已处于禁用状态"})
		return
	}

	if err := tc.DB.Model(&tenant).Update("status", "suspended").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "租户已禁用"})
}

// ActivateTenant 启用租户
func (tc *TenantController) ActivateTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	if tenant.Status == "active" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "ALREADY_ACTIVE", "message": "租户已是启用状态"})
		return
	}

	if err := tc.DB.Model(&tenant).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "租户已启用"})
}

// ExtendTenant 延长租期
func (tc *TenantController) ExtendTenant(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var req ExtendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PARAMS", "message": "请提供 extend_days 参数"})
		return
	}

	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	now := time.Now()
	newExpiry := now.AddDate(0, 0, req.ExtendDays)
	if tenant.ExpiresAt != nil {
		newExpiry = tenant.ExpiresAt.AddDate(0, 0, req.ExtendDays)
	}

	if err := tc.DB.Model(&tenant).Update("expires_at", newExpiry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "租期已延长",
		"data": gin.H{
			"expires_at": newExpiry,
		},
	})
}

// ChangePlan 变更套餐
func (tc *TenantController) ChangePlan(c *gin.Context) {
	if !middleware.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "需要超级管理员权限"})
		return
	}

	id := c.Param("id")
	var req ChangePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_PARAMS", "message": "参数错误"})
		return
	}

	var plan models.Plan
	if err := tc.DB.Where("id = ?", req.PlanID).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "PLAN_NOT_FOUND", "message": "套餐不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	var tenant models.Tenant
	if err := tc.DB.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "租户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}

	effectiveType := req.EffectiveType
	if effectiveType == "" {
		effectiveType = "immediate"
	}

	if effectiveType == "immediate" {
		if err := tc.DB.Model(&tenant).Update("plan", plan.PlanCode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "变更失败"})
			return
		}
	}

	settings := tenant.Settings
	if settings == nil {
		settings = make(map[string]interface{})
	}
	settings["pending_plan_id"] = req.PlanID
	settings["pending_plan_code"] = plan.PlanCode
	settings["effective_type"] = effectiveType
	tc.DB.Model(&tenant).Update("settings", settings)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "套餐变更已提交",
		"data": gin.H{
			"plan_code":      plan.PlanCode,
			"effective_type": effectiveType,
		},
	})
}

// ListPlans 获取套餐列表
func (tc *TenantController) ListPlans(c *gin.Context) {
	var plans []models.Plan
	if err := tc.DB.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "DB_ERROR", "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": plans,
	})
}

// ==================== 辅助函数 ====================
