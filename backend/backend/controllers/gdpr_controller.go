package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GDPRController GDPR合规控制器
type GDPRController struct {
	DB *gorm.DB
}

// GDPRExportRequest GDPR导出请求
type GDPRExportRequest struct {
	Format string `json:"format"` // json, csv, pdf
}

// GDPRDeleteAccountRequest 账户删除请求
type GDPRDeleteAccountRequest struct {
	ConfirmText string `json:"confirm_text" binding:"required"` // 必须输入 "DELETE"
	Reason      string `json:"reason"`                          // 删除原因（可选）
}

// ExportUserData 导出用户数据（GDPR要求）
// @Summary 导出用户数据
// @Description 导出指定用户的所有个人数据（GDPR合规）
// @Tags GDPR
// @Accept json
// @Produce json
// @Param user_id query int true "用户ID"
// @Param format query string false "导出格式" default(json) enum(json,csv)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/gdpr/export [get]
func (c *GDPRController) ExportUserData(ctx *gin.Context) {
	// 获取目标用户ID
	userIDStr := ctx.Query("user_id")
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少用户ID参数",
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	format := ctx.DefaultQuery("format", "json")

	// 验证权限：只能导出自己的数据，或管理员可以导出
	currentUserID, _ := ctx.Get("user_id")
	isSuperAdmin, _ := ctx.Get("is_super_admin")

	if uint(userID) != currentUserID.(uint) && !isSuperAdmin.(bool) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权导出其他用户的数据",
		})
		return
	}

	// 收集用户数据
	exportData := make(map[string]interface{})

	// 1. 用户基本信息
	var user models.SysUser
	if err := c.DB.First(&user, userID).Error; err == nil {
		userData := map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"phone":      user.Phone,
			"status":     user.Status,
			"tenant_id":  user.TenantID,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		}
		// 对敏感字段进行脱敏
		exportData["user"] = middleware.MaskUserData(userData)
	}

	// 2. 员工扩展信息（包含敏感字段）
	var employee models.Employee
	if err := c.DB.Where("emp_code = (SELECT emp_code FROM sys_users WHERE id = ?)", userID).First(&employee).Error; err == nil {
		engine := middleware.GetMaskingEngine()
		extData := map[string]interface{}{
			"id":          employee.ID,
			"emp_code":    employee.EmpCode,
			"emp_name":    engine.MaskValue("real_name", employee.EmpName), // 脱敏
			"gender":      employee.Gender,
			"birth_date":  employee.BirthDate,
			"phone":       engine.MaskValue("phone", employee.Phone),      // 脱敏
			"email":       engine.MaskValue("email", employee.Email),      // 脱敏
			"id_card":     engine.MaskValue("id_card", employee.IDCard),    // 脱敏
			"province":    employee.Province,
			"city":        employee.City,
			"district":    employee.District,
			"address":     engine.MaskValue("address", employee.Address), // 脱敏
		}
		exportData["employee"] = extData
	}

	// 3. 设备列表
	var devices []models.Device
	c.DB.Where("bind_user_id = ?", fmt.Sprintf("%d", userID)).Find(&devices)
	devicesData := make([]map[string]interface{}, 0)
	for _, d := range devices {
		devicesData = append(devicesData, map[string]interface{}{
			"id":                d.ID,
			"device_id":         d.DeviceID,
			"mac_address":       d.MacAddress,
			"sn_code":           d.SnCode,
			"hardware_model":    d.HardwareModel,
			"firmware_version":  d.FirmwareVersion,
			"lifecycle_status":  d.LifecycleStatus,
			"created_at":        d.CreatedAt,
		})
	}
	exportData["devices"] = devicesData

	// 4. 会员信息
	var members []models.Member
	c.DB.Where("phone = ?", user.Phone).Find(&members)
	membersData := make([]map[string]interface{}, 0)
	for _, m := range members {
		membersData = append(membersData, middleware.MaskMemberData(map[string]interface{}{
			"id":            m.ID,
			"member_code":   m.MemberCode,
			"member_name":   m.MemberName,
			"phone":         m.Phone,
			"gender":       m.Gender,
			"birth_date":   m.BirthDate,
			"email":        m.Email,
			"avatar":       m.Avatar,
			"member_level": m.MemberLevel,
			"points":       m.Points,
			"balance":      m.Balance,
			"created_at":   m.CreatedAt,
		}))
	}
	exportData["members"] = membersData

	// 5. 操作日志
	var opLogs []models.SysOperationLog
	c.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(100).Find(&opLogs)
	opLogsData := make([]map[string]interface{}, 0)
	for _, log := range opLogs {
		opLogsData = append(opLogsData, map[string]interface{}{
			"id":         log.ID,
			"module":     log.Module,
			"operation":  log.Operation,
			"method":     log.Method,
			"path":       log.Path,
			"ip":         log.IP,
			"status":     log.Status,
			"created_at": log.CreatedAt,
		})
	}
	exportData["operation_logs"] = opLogsData

	// 6. 登录日志
	var loginLogs []models.SysLoginLog
	c.DB.Where("user_id = ?", userID).Order("login_time DESC").Limit(100).Find(&loginLogs)
	loginLogsData := make([]map[string]interface{}, 0)
	for _, log := range loginLogs {
		loginLogsData = append(loginLogsData, map[string]interface{}{
			"id":         log.ID,
			"ip":         log.IP,
			"location":   log.Location,
			"browser":    log.Browser,
			"os":         log.OS,
			"status":     log.Status,
			"login_time": log.LoginTime,
		})
	}
	exportData["login_logs"] = loginLogsData

	// 7. 会员积分记录
	var pointsRecords []models.MemberPointsRecord
	c.DB.Joins("JOIN members ON members.phone = ?", user.Phone).
		Where("points_records.member_id = members.id", user.Phone).
		Order("created_at DESC").Limit(100).Find(&pointsRecords)
	pointsData := make([]map[string]interface{}, 0)
	for _, p := range pointsRecords {
		pointsData = append(pointsData, map[string]interface{}{
			"id":            p.ID,
			"member_id":     p.MemberID,
			"points":       p.Points,
			"points_type":  p.PointsType,
			"source_type":  p.SourceType,
			"source_id":    p.SourceID,
			"order_no":     p.OrderNo,
			"created_at":   p.CreatedAt,
		})
	}
	exportData["points_records"] = pointsData

	// 8. 添加导出元数据
	exportData["_meta"] = map[string]interface{}{
		"exported_at":      time.Now().Format(time.RFC3339),
		"exported_by":      currentUserID.(uint),
		"format":           format,
		"gdpr_article":     "Article 15 - Right of access",
		"data_controller": "MDM Platform",
	}

	// 根据格式返回
	switch format {
	case "csv":
		// CSV格式导出（简化版，仅包含基本信息）
		csvData := fmt.Sprintf("Type,Field,Value,Exported_At\n")
		if userData, ok := exportData["user"].(map[string]interface{}); ok {
			for k, v := range userData {
				if k != "_meta" {
					csvData += fmt.Sprintf("user,%s,%v,%s\n", k, v, exportData["_meta"].(map[string]interface{})["exported_at"])
				}
			}
		}
		ctx.Header("Content-Type", "text/csv")
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=gdpr_export_%d_%d.csv", userID, time.Now().Unix()))
		ctx.String(http.StatusOK, csvData)
	default:
		// JSON格式
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "数据导出成功",
			"data":    exportData,
		})
	}
}

// DeleteAccount 删除账户（GDPR要求 - 被遗忘权）
// @Summary 删除账户
// @Description 删除用户账户及其所有相关数据（GDPR合规 - 被遗忘权）
// @Tags GDPR
// @Accept json
// @Produce json
// @Param request body GDPRDeleteAccountRequest true "删除确认"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/gdpr/delete-account [delete]
func (c *GDPRController) DeleteAccount(ctx *gin.Context) {
	var req GDPRDeleteAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误，请确认已输入 DELETE 进行二次确认",
			"error":   err.Error(),
		})
		return
	}

	// 验证确认文本
	if req.ConfirmText != "DELETE" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "必须输入 DELETE 进行二次确认",
		})
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 如果是管理员操作，支持删除其他用户
	targetUserID := userID.(uint)
	if adminIDStr := ctx.Query("target_user_id"); adminIDStr != "" {
		isSuperAdmin, _ := ctx.Get("is_super_admin")
		if isSuperAdmin.(bool) {
			if parsed, err := strconv.ParseUint(adminIDStr, 10, 32); err == nil {
				targetUserID = uint(parsed)
			}
		}
	}

	// 开启事务
	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 收集删除前数据（用于审计）
	var user models.SysUser
	if err := tx.First(&user, targetUserID).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 2. 删除用户扩展信息
	tx.Where("user_id = ?", targetUserID).Delete(&models.SysUserExt{})

	// 3. 解除设备绑定
	tx.Model(&models.Device{}).Where("user_id = ?", targetUserID).Updates(map[string]interface{}{
		"user_id":    0,
		"bind_time":  nil,
		"unbind_time": time.Now(),
	})

	// 4. 会员数据处理（根据业务需求，可以是软删除或匿名化）
	var member models.Member
	if tx.Where("phone = ?", user.Phone).First(&member).Error == nil {
		// 匿名化会员数据（保留积分等业务数据）
		member.Phone = fmt.Sprintf("deleted_%d@anonymized.local", member.ID)
		member.Email = ""
		member.MemberName = "已注销用户"
		member.Avatar = ""
		tx.Save(&member)
	}

	// 5. 删除操作日志
	tx.Where("user_id = ?", targetUserID).Delete(&models.SysOperationLog{})

	// 6. 删除登录日志
	tx.Where("user_id = ?", targetUserID).Delete(&models.SysLoginLog{})

	// 7. 删除用户角色关联
	tx.Where("user_id = ?", targetUserID).Delete(&models.SysUserRole{})

	// 8. 记录账户删除审计日志
	auditLog := map[string]interface{}{
		"user_id":       targetUserID,
		"action":        "gdpr_account_deletion",
		"reason":        req.Reason,
		"deleted_at":    time.Now(),
		"deleted_by":    userID,
		"anonymized_data": map[string]string{
			"phone": fmt.Sprintf("deleted_%d@anonymized.local", member.ID),
			"email": "",
		},
	}
	auditJSON, _ := json.Marshal(auditLog)
	tx.Create(&models.SysOperationLog{
		UserID:    userID.(uint),
		Module:    "GDPR",
		Operation: "Account Deletion",
		Path:      "/api/v1/gdpr/delete-account",
		Method:    "DELETE",
		Params:    string(auditJSON),
		Status:    2, // 成功
	})

	// 9. 软删除用户（保留记录用于审计）
	user.Status = 0 // 禁用状态
	user.Username = fmt.Sprintf("deleted_%d_%d", targetUserID, time.Now().Unix())
	tx.Save(&user)

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除账户失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "账户删除成功",
		"data": map[string]interface{}{
			"deleted_user_id": targetUserID,
			"deleted_at":      time.Now().Format(time.RFC3339),
			"gdpr_article":    "Article 17 - Right to erasure ('right to be forgotten')",
		},
	})
}

// GetDataProcessingInfo 获取数据处理信息
// @Summary 获取数据处理信息
// @Description 获取用户数据的处理方式和保留期限
// @Tags GDPR
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/gdpr/info [get]
func (c *GDPRController) GetDataProcessingInfo(ctx *gin.Context) {
	info := map[string]interface{}{
		"data_controller": "MDM Platform",
		"privacy_policy_url": "/privacy-policy",
		"contact_email":    "privacy@mdm-platform.com",
		"data_categories": []map[string]string{
			{"category": "个人基本信息", "retention": "账户存续期间", "purpose": "用户身份识别和服务提供"},
			{"category": "设备数据", "retention": "账户存续期间 + 2年", "purpose": "设备管理和支持"},
			{"category": "操作日志", "retention": "180天", "purpose": "安全审计"},
			{"category": "会员信息", "retention": "账户存续期间", "purpose": "会员服务"},
		},
		"data_subjects": []string{
			"用户账户信息",
			"设备绑定信息",
			"会员数据",
			"操作和登录日志",
			"积分和余额信息",
		},
		"rights": []map[string]string{
			{"right": "访问权", "article": "Article 15", "description": "获取您的个人数据副本"},
			{"right": "更正权", "article": "Article 16", "description": "更正不准确的个人数据"},
			{"right": "删除权", "article": "Article 17", "description": "要求删除您的个人数据"},
			{"right": "限制处理权", "article": "Article 18", "description": "限制我们处理您的数据"},
			{"right": "数据可携权", "article": "Article 20", "description": "获取您的数据格式"},
			{"right": "反对权", "article": "Article 21", "description": "反对某些类型的处理"},
		},
		"international_transfers": map[string]string{
			"status": "无",
			"description": "您的数据仅在本地服务器处理",
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": info,
	})
}

// AnonymizeData 匿名化数据（用于数据保留策略）
// @Summary 匿名化历史数据
// @Description 将超过保留期限的数据匿名化处理
// @Tags GDPR
// @Accept json
// @Produce json
// @Param days query int false "超过多少天的数据" default(365)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/gdpr/anonymize [post]
func (c *GDPRController) AnonymizeData(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "365"))
	
	// 验证管理员权限
	isSuperAdmin, _ := ctx.Get("is_super_admin")
	if !isSuperAdmin.(bool) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "需要管理员权限",
		})
		return
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)

	// 匿名化超过期限的登录日志
	result := c.DB.Model(&models.SysLoginLog{}).
		Where("login_time < ?", cutoffDate).
		Updates(map[string]interface{}{
			"ip":       "anonymized",
			"location": "anonymized",
		})

	anonymizedCount := result.RowsAffected

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "数据匿名化完成",
		"data": map[string]interface{}{
			"anonymized_records": anonymizedCount,
			"cutoff_date":         cutoffDate.Format(time.RFC3339),
			"processed_at":        time.Now().Format(time.RFC3339),
		},
	})
}
