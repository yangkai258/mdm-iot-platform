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
