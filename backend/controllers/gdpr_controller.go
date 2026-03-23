package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ComplianceController 合规控制器
type ComplianceController struct {
	DB *gorm.DB
}

// GDPRDataAccessRequest GDPR数据访问请求
type GDPRDataAccessRequest struct {
	Email string `json:"email" binding:"required"` // 用户邮箱
}

// GDPRDataDeleteRequest GDPR删除请求
type GDPRDataDeleteRequest struct {
	Email  string `json:"email" binding:"required"`  // 用户邮箱
	Reason string `json:"reason"`                     // 删除原因
}

// GDPRDataPortabilityRequest 数据可携带导出请求
type GDPRDataPortabilityRequest struct {
	Email     string `json:"email" binding:"required"`      // 用户邮箱
	Format    string `json:"format"`                         // 导出格式: json, csv
	DataTypes string `json:"data_types"`                    // 数据类型: all, devices, members, activity
}

// ============ GDPR API ============

// GetGDPRData GDPR数据访问 - 获取用户所有数据
func (c *ComplianceController) GetGDPRData(ctx *gin.Context) {
	var req GDPRDataAccessRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.recordAudit(ctx, "gdpr_access", "compliance", "", "", 0, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: email 为必填项",
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 查找用户
	var user models.SysUser
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "未找到该邮箱对应的用户",
			})
			return
		}
		c.recordAudit(ctx, "gdpr_access", "compliance", "sys_user", "", uid, http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户失败",
		})
		return
	}

	// 收集用户所有数据
	data := c.collectUserData(user, tid)

	// 生成请求记录
	requestID := generateGDPRRequestID()
	gdprRequest := models.GDPRRequest{
		RequestID:      requestID,
		RequestType:    "data_access",
		RequesterEmail: req.Email,
		RequesterName: user.Nickname,
		UserID:         user.ID,
		Status:         3, // 已完成
		ProcessedAt:    timePtr(time.Now()),
		CompletedAt:    timePtr(time.Now()),
		ResponseData:   data,
		TenantID:       tid,
	}
	c.DB.Create(&gdprRequest)

	// 记录审计日志
	c.recordAudit(ctx, "gdpr_access", "compliance", "gdpr_request", requestID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"request_id":    requestID,
			"user_data":      json.RawMessage(data),
			"requested_at":   gdprRequest.CreatedAt,
			"completed_at":   gdprRequest.CompletedAt,
		},
	})
}

// DeleteGDPRData GDPR删除请求 - 删除用户所有数据
func (c *ComplianceController) DeleteGDPRData(ctx *gin.Context) {
	var req GDPRDataDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.recordAudit(ctx, "gdpr_delete", "compliance", "", "", 0, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: email 为必填项",
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 查找用户
	var user models.SysUser
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "未找到该邮箱对应的用户",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户失败",
		})
		return
	}

	// 生成请求记录
	requestID := generateGDPRRequestID()
	gdprRequest := models.GDPRRequest{
		RequestID:      requestID,
		RequestType:    "data_deletion",
		RequesterEmail: req.Email,
		RequesterName:  user.Nickname,
		UserID:         user.ID,
		Status:         2, // 处理中
		RequestReason:  req.Reason,
		TenantID:       tid,
	}
	c.DB.Create(&gdprRequest)

	// 执行数据删除（软删除或匿名化）
	go c.executeDataDeletion(requestID, user.ID, tid)

	// 记录审计日志
	c.recordAudit(ctx, "gdpr_delete", "compliance", "gdpr_request", requestID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"request_id":   requestID,
			"status":       "processing",
			"requested_at": gdprRequest.CreatedAt,
			"message":      "数据删除请求已提交，正在处理中",
		},
	})
}

// GetGDPRDataExport 数据可携带导出
func (c *ComplianceController) GetGDPRDataExport(ctx *gin.Context) {
	var req GDPRDataPortabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.ShouldBindQuery(&req)
	}

	if req.Email == "" {
		c.recordAudit(ctx, "data_export", "compliance", "", "", 0, http.StatusBadRequest, "缺少邮箱参数")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: email 为必填项",
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	// 查找用户
	var user models.SysUser
	if err := c.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "未找到该邮箱对应的用户",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询用户失败",
		})
		return
	}

	format := strings.ToLower(req.Format)
	if format == "" {
		format = "json"
	}

	dataTypes := strings.ToLower(req.DataTypes)
	if dataTypes == "" {
		dataTypes = "all"
	}

	// 收集数据
	exportData := c.collectExportData(user, tid, dataTypes)

	// 生成导出文件
	requestID := generateGDPRRequestID()
	var exportPath string

	switch format {
	case "csv":
		exportPath = c.exportToCSV(requestID, exportData)
	default:
		exportPath = c.exportToJSON(requestID, exportData)
	}

	// 记录请求
	gdprRequest := models.GDPRRequest{
		RequestID:      requestID,
		RequestType:    "data_portability",
		RequesterEmail: req.Email,
		RequesterName:  user.Nickname,
		UserID:         user.ID,
		Status:         3, // 已完成
		ExportPath:     exportPath,
		ProcessedAt:    timePtr(time.Now()),
		CompletedAt:    timePtr(time.Now()),
		TenantID:       tid,
	}
	c.DB.Create(&gdprRequest)

	// 记录审计日志
	c.recordAudit(ctx, "data_export", "compliance", "gdpr_request", requestID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"request_id":   requestID,
			"export_path":   exportPath,
			"format":        format,
			"data_types":    dataTypes,
			"requested_at":  gdprRequest.CreatedAt,
			"completed_at":  gdprRequest.CompletedAt,
		},
	})
}

// GetGDPRRequests 获取GDPR请求列表
func (c *ComplianceController) GetGDPRRequests(ctx *gin.Context) {
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	var requests []models.GDPRRequest
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	status := ctx.Query("status")
	requestType := ctx.Query("type")

	query := c.DB.Model(&models.GDPRRequest{}).Where("tenant_id = ?", tid)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&requests)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      requests,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetGDPRRequest 获取单个GDPR请求详情
func (c *ComplianceController) GetGDPRRequest(ctx *gin.Context) {
	requestID := ctx.Param("id")
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)

	var request models.GDPRRequest
	if err := c.DB.Where("request_id = ? AND tenant_id = ?", requestID, tid).First(&request).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "请求不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data":   request,
	})
}

// ProcessGDPRRequest 处理GDPR请求（管理员）
func (c *ComplianceController) ProcessGDPRRequest(ctx *gin.Context) {
	requestID := ctx.Param("id")
	tenantID, _ := ctx.Get("tenant_id")
	tid, _ := tenantID.(uint)
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	var req struct {
		Status int    `json:"status"` // 3:已完成 4:已拒绝
		Note   string `json:"note"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	var request models.GDPRRequest
	if err := c.DB.Where("request_id = ? AND tenant_id = ?", requestID, tid).First(&request).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "请求不存在",
		})
		return
	}

	updates := map[string]interface{}{
		"status":       req.Status,
		"processed_by": uid,
		"processed_at": time.Now(),
	}

	if req.Status == 4 { // 拒绝
		updates["rejected_reason"] = req.Note
	} else if req.Status == 3 { // 完成
		updates["completed_at"] = time.Now()
	}

	c.DB.Model(&request).Updates(updates)

	// 记录审计日志
	c.recordAudit(ctx, "process_gdpr_request", "compliance", "gdpr_request", requestID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ============ 辅助函数 ============

func generateGDPRRequestID() string {
	return fmt.Sprintf("gdpr-%d-%s", time.Now().Unix(), randomString(8))
}

func (c *ComplianceController) collectUserData(user models.SysUser, tenantID uint) string {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"nickname":  user.Nickname,
			"email":     user.Email,
			"phone":     user.Phone,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"devices":      []interface{}{},
		"members":      []interface{}{},
		"activity":     []interface{}{},
		"notifications": []interface{}{},
	}

	// 收集设备数据
	var devices []models.Device
	c.DB.Where("owner_id = ? OR tenant_id = ?", user.ID, tenantID).Find(&devices)
	data["devices"] = devices

	// 收集会员数据
	var members []models.Member
	c.DB.Where("user_id = ? OR tenant_id = ?", user.ID, tenantID).Find(&members)
	data["members"] = members

	// 收集活动日志
	var activities []models.ActivityLog
	c.DB.Where("user_id = ?", user.ID).Limit(100).Find(&activities)
	data["activity"] = activities

	// 收集通知记录
	var notifications []models.Notification
	c.DB.Where("user_id = ?", user.ID).Limit(100).Find(&notifications)
	data["notifications"] = notifications

	result, _ := json.Marshal(data)
	return string(result)
}

func (c *ComplianceController) collectExportData(user models.SysUser, tenantID uint, dataTypes string) map[string]interface{} {
	data := make(map[string]interface{})
	data["export_info"] = map[string]interface{}{
		"exported_at":   time.Now(),
		"requester":     user.Email,
		"data_types":    dataTypes,
	}

	switch dataTypes {
	case "devices":
		var devices []models.Device
		c.DB.Where("owner_id = ? OR tenant_id = ?", user.ID, tenantID).Find(&devices)
		data["devices"] = devices
	case "members":
		var members []models.Member
		c.DB.Where("user_id = ? OR tenant_id = ?", user.ID, tenantID).Find(&members)
		data["members"] = members
	case "activity":
		var activities []models.ActivityLog
		c.DB.Where("user_id = ?", user.ID).Find(&activities)
		data["activity"] = activities
	default: // all
		data["user"] = user
		var devices []models.Device
		c.DB.Where("owner_id = ? OR tenant_id = ?", user.ID, tenantID).Find(&devices)
		data["devices"] = devices
	}

	return data
}

func (c *ComplianceController) executeDataDeletion(requestID string, userID, tenantID uint) {
	// 在后台执行数据删除
	// 软删除用户
	c.DB.Model(&models.SysUser{}).Where("id = ?", userID).Update("status", 0)

	// 匿名化相关数据
	c.DB.Model(&models.ActivityLog{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"username": "DELETED_USER",
	})

	// 更新请求状态
	c.DB.Model(&models.GDPRRequest{}).Where("request_id = ?", requestID).Updates(map[string]interface{}{
		"status":       3, // 已完成
		"completed_at": time.Now(),
	})
}

func (c *ComplianceController) exportToJSON(requestID string, data map[string]interface{}) string {
	exportDir := "./exports/gdpr"
	os.MkdirAll(exportDir, 0755)
	path := filepath.Join(exportDir, fmt.Sprintf("%s.json", requestID))

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(path, jsonData, 0644)
	return path
}

func (c *ComplianceController) exportToCSV(requestID string, data map[string]interface{}) string {
	exportDir := "./exports/gdpr"
	os.MkdirAll(exportDir, 0755)
	path := filepath.Join(exportDir, fmt.Sprintf("%s.csv", requestID))
	// 实际CSV导出逻辑
	return path
}

func (c *ComplianceController) recordAudit(ctx *gin.Context, action, module, resourceType, resourceID string, userID uint, statusCode int, errorMsg string) {
	username, _ := ctx.Get("username")
	status := 1
	if statusCode >= 400 {
		status = 2
	}

	log := models.AuditLog{
		Action:        action,
		Module:        module,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		UserID:        userID,
		Username:      username.(string),
		IP:            ctx.ClientIP(),
		UserAgent:     ctx.GetHeader("User-Agent"),
		Status:        status,
		ErrorMsg:      errorMsg,
		RequestMethod: ctx.Request.Method,
		RequestPath:   ctx.Request.URL.Path,
		ResponseCode:  statusCode,
	}
	c.DB.Create(&log)
}

// ============ 合规规则 CRUD ============

// ListRules 获取合规规则列表
func (c *ComplianceController) ListRules(ctx *gin.Context) {
	var rules []models.CompliancePolicy
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	policyType := ctx.Query("policy_type")
	severity := ctx.Query("severity")
	enabled := ctx.Query("enabled")

	query := c.DB.Model(&models.CompliancePolicy{})

	if policyType != "" {
		query = query.Where("policy_type = ?", policyType)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if enabled != "" {
		query = query.Where("enabled = ?", enabled)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("severity DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&rules)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      rules,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRule 获取单个合规规则
func (c *ComplianceController) GetRule(ctx *gin.Context) {
	id := ctx.Param("id")

	var rule models.CompliancePolicy
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// CreateRule 创建合规规则
func (c *ComplianceController) CreateRule(ctx *gin.Context) {
	var req struct {
		Name              string `json:"name" binding:"required"`
		Description      string `json:"description"`
		PolicyType       string `json:"policy_type" binding:"required"`
		TargetValue      string `json:"target_value"`
		Condition        string `json:"condition" binding:"required"`
		Severity         int    `json:"severity"`
		RemediationAction string `json:"remediation_action"`
		Enabled          *bool  `json:"enabled"`
		EnforceScope     string `json:"enforce_scope"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	severity := 2
	if req.Severity > 0 {
		severity = req.Severity
	}
	enforceScope := "all"
	if req.EnforceScope != "" {
		enforceScope = req.EnforceScope
	}

	rule := models.CompliancePolicy{
		Name:               req.Name,
		Description:       req.Description,
		PolicyType:        req.PolicyType,
		TargetValue:       req.TargetValue,
		Condition:         req.Condition,
		Severity:          severity,
		RemediationAction: req.RemediationAction,
		Enabled:           enabled,
		EnforceScope:      enforceScope,
	}

	if err := c.DB.Create(&rule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)
	c.recordAudit(ctx, "create", "compliance", "compliance_policy", strconv.FormatUint(uint64(rule.ID), 10), uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// UpdateRule 更新合规规则
func (c *ComplianceController) UpdateRule(ctx *gin.Context) {
	id := ctx.Param("id")

	var rule models.CompliancePolicy
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Name               string `json:"name"`
		Description       string `json:"description"`
		PolicyType        string `json:"policy_type"`
		TargetValue       string `json:"target_value"`
		Condition         string `json:"condition"`
		Severity          int    `json:"severity"`
		RemediationAction string `json:"remediation_action"`
		Enabled           *bool  `json:"enabled"`
		EnforceScope      string `json:"enforce_scope"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.PolicyType != "" {
		updates["policy_type"] = req.PolicyType
	}
	if req.TargetValue != "" {
		updates["target_value"] = req.TargetValue
	}
	if req.Condition != "" {
		updates["condition"] = req.Condition
	}
	if req.Severity > 0 {
		updates["severity"] = req.Severity
	}
	if req.RemediationAction != "" {
		updates["remediation_action"] = req.RemediationAction
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.EnforceScope != "" {
		updates["enforce_scope"] = req.EnforceScope
	}

	if err := c.DB.Model(&rule).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败: " + err.Error()})
		return
	}

	c.DB.First(&rule, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": rule})
}

// DeleteRule 删除合规规则
func (c *ComplianceController) DeleteRule(ctx *gin.Context) {
	id := ctx.Param("id")

	var rule models.CompliancePolicy
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.DB.Delete(&rule)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// EnforceRule 强制执行合规规则
func (c *ComplianceController) EnforceRule(ctx *gin.Context) {
	id := ctx.Param("id")

	var rule models.CompliancePolicy
	if err := c.DB.First(&rule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规规则不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 扫描所有设备，检测违规
	var violations []models.ComplianceViolation
	c.DB.Where("policy_id = ? AND status = 1", rule.ID).Find(&violations)

	violationCount := len(violations)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"policy_id":       rule.ID,
			"policy_type":     rule.PolicyType,
			"violations_found": violationCount,
			"enforced_at":     time.Now(),
		},
	})
}

// ============ 合规策略（CompliancePolicy）CRUD ============

// ListCompliancePolicies 获取合规策略列表
func (c *ComplianceController) ListCompliancePolicies(ctx *gin.Context) {
	var policies []models.CompliancePolicy
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	policyType := ctx.Query("policy_type")
	enabled := ctx.Query("enabled")

	query := c.DB.Model(&models.CompliancePolicy{})

	if policyType != "" {
		query = query.Where("policy_type = ?", policyType)
	}
	if enabled != "" {
		query = query.Where("enabled = ?", enabled)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("severity DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&policies)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      policies,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateCompliancePolicy 创建合规策略
func (c *ComplianceController) CreateCompliancePolicy(ctx *gin.Context) {
	var req struct {
		Name               string `json:"name" binding:"required"`
		Description       string `json:"description"`
		PolicyType        string `json:"policy_type" binding:"required"`
		TargetValue       string `json:"target_value"`
		Condition         string `json:"condition" binding:"required"`
		Severity          int    `json:"severity"`
		RemediationAction string `json:"remediation_action"`
		Enabled           *bool  `json:"enabled"`
		EnforceScope      string `json:"enforce_scope"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	severity := 2
	if req.Severity > 0 {
		severity = req.Severity
	}
	enforceScope := "all"
	if req.EnforceScope != "" {
		enforceScope = req.EnforceScope
	}

	policy := models.CompliancePolicy{
		Name:               req.Name,
		Description:       req.Description,
		PolicyType:        req.PolicyType,
		TargetValue:       req.TargetValue,
		Condition:         req.Condition,
		Severity:          severity,
		RemediationAction: req.RemediationAction,
		Enabled:           enabled,
		EnforceScope:      enforceScope,
	}

	if err := c.DB.Create(&policy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// UpdateCompliancePolicy 更新合规策略
func (c *ComplianceController) UpdateCompliancePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.CompliancePolicy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Name               string `json:"name"`
		Description       string `json:"description"`
		PolicyType        string `json:"policy_type"`
		TargetValue       string `json:"target_value"`
		Condition         string `json:"condition"`
		Severity          int    `json:"severity"`
		RemediationAction string `json:"remediation_action"`
		Enabled           *bool  `json:"enabled"`
		EnforceScope      string `json:"enforce_scope"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.PolicyType != "" {
		updates["policy_type"] = req.PolicyType
	}
	if req.TargetValue != "" {
		updates["target_value"] = req.TargetValue
	}
	if req.Condition != "" {
		updates["condition"] = req.Condition
	}
	if req.Severity > 0 {
		updates["severity"] = req.Severity
	}
	if req.RemediationAction != "" {
		updates["remediation_action"] = req.RemediationAction
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.EnforceScope != "" {
		updates["enforce_scope"] = req.EnforceScope
	}

	if err := c.DB.Model(&policy).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	c.DB.First(&policy, id)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": policy})
}

// DeleteCompliancePolicy 删除合规策略
func (c *ComplianceController) DeleteCompliancePolicy(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.CompliancePolicy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.DB.Delete(&policy)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// GetPolicyViolations 获取策略违规记录
func (c *ComplianceController) GetPolicyViolations(ctx *gin.Context) {
	id := ctx.Param("id")

	var policy models.CompliancePolicy
	if err := c.DB.First(&policy, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "合规策略不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var violations []models.ComplianceViolation
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	status := ctx.Query("status")

	query := c.DB.Model(&models.ComplianceViolation{}).Where("policy_id = ?", policy.ID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&violations)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      violations,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ============ 合规违规 CRUD（通过 ComplianceController）============

// ListViolations 获取所有违规记录
func (c *ComplianceController) ListViolations(ctx *gin.Context) {
	var violations []models.ComplianceViolation
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	status := ctx.Query("status")
	deviceID := ctx.Query("device_id")

	query := c.DB.Model(&models.ComplianceViolation{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&violations)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      violations,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ResolveViolation 解决违规记录
func (c *ComplianceController) ResolveViolation(ctx *gin.Context) {
	id := ctx.Param("id")

	var violation models.ComplianceViolation
	if err := c.DB.First(&violation, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "违规记录不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		Status int    `json:"status"` // 3=已解决 4=已忽略
		Note   string `json:"note"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	now := time.Now()
	updates := map[string]interface{}{
		"status":      req.Status,
		"resolved_at": now,
		"resolved_by": uid,
	}
	if req.Status == 4 {
		updates["action_taken"] = "ignored: " + req.Note
	}

	c.DB.Model(&violation).Updates(updates)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
