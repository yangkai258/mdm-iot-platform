package controllers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"mdm-backend/middleware"
	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type SecurityController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

var secAesKey = []byte(getSecAESSecret())

func getSecAESSecret() string {
	secret := os.Getenv("APP_SECRET_KEY")
	if secret == "" {
		return "mdm-platform-32-byte-secret-key!"
	}
	if len(secret) < 32 {
		secret = fmt.Sprintf("%-32s", secret)
	}
	return secret[:32]
}

func secEncrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(secAesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

type auditArgs struct {
	eventType     string
	category      string
	severity      int
	targetUserID  uint
	targetUsername string
	ip            string
	userAgent     string
	sessionID     string
	resourceType  string
	resourceID    string
	status        int
	errorMsg      string
	reqMethod     string
	reqPath       string
	responseCode  int
	duration      int
	metadata      string
	tenantID      string
}

// ============ 2FA APIs ============

func (c *SecurityController) Enable2FA(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	var existing2FA models.TwoFactorAuth
	if err := c.DB.Where("user_id = ?", userID).First(&existing2FA).Error; err == nil && existing2FA.IsEnabled {
		ctx.JSON(http.StatusConflict, gin.H{"code": 409, "message": "2FA已启用，无需重复开启"})
		return
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MDM-Platform",
		AccountName: username,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成2FA密钥失败"})
		return
	}

	encryptedSecret, err := secEncrypt(secret.Secret())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "加密2FA密钥失败"})
		return
	}

	recoveryCodes := genRecoveryCodes(8)
	recoveryCodesJSON, _ := json.Marshal(recoveryCodes)
	encryptedRecoveryCodes, _ := secEncrypt(string(recoveryCodesJSON))

	now := time.Now()
	record := models.TwoFactorAuth{
		UserID:          userID,
		Secret:          secret.Secret(),
		SecretEncrypted: encryptedSecret,
		IsEnabled:       false,
		IsVerified:      false,
		RecoveryCodes:   encryptedRecoveryCodes,
		EnabledAt:      &now,
	}

	if existing2FA.ID != 0 {
		record.ID = existing2FA.ID
		record.CreatedAt = existing2FA.CreatedAt
		c.DB.Save(&record)
	} else {
		if err := c.DB.Create(&record).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存2FA配置失败"})
			return
		}
	}

	c.saveAudit(userID, username, auditArgs{
		eventType:    "2fa_enable_request",
		category:    "authentication",
		severity:    1,
		ip:          ctx.ClientIP(),
		userAgent:   ctx.GetHeader("User-Agent"),
		status:      1,
		responseCode: http.StatusOK,
		tenantID:    ctx.GetHeader("X-Tenant-ID"),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"secret":         secret.Secret(),
			"otpauth_uri":    secret.URL(),
			"qr_code_url":    fmt.Sprintf("otpauth://totp/MDM-Platform:%s?secret=%s&issuer=MDM-Platform&digits=6&period=30", username, secret.Secret()),
			"recovery_codes": recoveryCodes,
		},
	})
}

func (c *SecurityController) Verify2FA(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var tfa models.TwoFactorAuth
	if err := c.DB.Where("user_id = ?", userID).First(&tfa).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "请先获取2FA密钥"})
		return
	}

	valid, err := totp.ValidateCustom(req.Code, tfa.Secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil || !valid {
		c.saveAudit(userID, username, auditArgs{
			eventType:    "2fa_verify_fail",
			category:    "authentication",
			severity:    2,
			ip:          ctx.ClientIP(),
			userAgent:   ctx.GetHeader("User-Agent"),
			status:      2,
			errorMsg:    "invalid_totp",
			responseCode: http.StatusOK,
			tenantID:    ctx.GetHeader("X-Tenant-ID"),
		})
		c.recordLoginAttempt(userID, username, ctx.ClientIP(), "verify_2fa", 2, "invalid_totp")
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "验证码错误或已过期"})
		return
	}

	tfa.IsEnabled = true
	tfa.IsVerified = true
	now := time.Now()
	tfa.LastUsedAt = &now
	c.DB.Save(&tfa)

	c.saveAudit(userID, username, auditArgs{
		eventType:    "2fa_verify_success",
		category:    "authentication",
		severity:    1,
		ip:          ctx.ClientIP(),
		userAgent:   ctx.GetHeader("User-Agent"),
		status:      1,
		responseCode: http.StatusOK,
		tenantID:    ctx.GetHeader("X-Tenant-ID"),
	})
	c.recordLoginAttempt(userID, username, ctx.ClientIP(), "verify_2fa", 1, "")

	token, _ := middleware.GenerateToken(userID, username, middleware.GetRoleID(ctx), middleware.GetTenantIDCtx(ctx), middleware.IsSuperAdminCtx(ctx))

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "2FA验证成功，已启用",
		"data": gin.H{
			"is_enabled":  true,
			"is_verified": true,
			"token":       token,
		},
	})
}

func (c *SecurityController) Disable2FA(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var tfa models.TwoFactorAuth
	if err := c.DB.Where("user_id = ?", userID).First(&tfa).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "2FA未启用"})
		return
	}

	valid, err := totp.ValidateCustom(req.Code, tfa.Secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil || !valid {
		c.saveAudit(userID, username, auditArgs{
			eventType:    "2fa_disable_fail",
			category:    "authentication",
			severity:    2,
			ip:          ctx.ClientIP(),
			userAgent:   ctx.GetHeader("User-Agent"),
			status:      2,
			errorMsg:    "invalid_totp",
			responseCode: http.StatusOK,
			tenantID:    ctx.GetHeader("X-Tenant-ID"),
		})
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "验证码错误"})
		return
	}

	now := time.Now()
	tfa.IsEnabled = false
	tfa.IsVerified = false
	tfa.DisabledAt = &now
	tfa.Secret = ""
	tfa.SecretEncrypted = ""
	tfa.RecoveryCodes = ""
	c.DB.Save(&tfa)

	c.saveAudit(userID, username, auditArgs{
		eventType:    "2fa_disable",
		category:    "authentication",
		severity:    1,
		ip:          ctx.ClientIP(),
		userAgent:   ctx.GetHeader("User-Agent"),
		status:      1,
		responseCode: http.StatusOK,
		tenantID:    ctx.GetHeader("X-Tenant-ID"),
	})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "2FA已禁用"})
}

// ============ Session Management APIs ============

func (c *SecurityController) GetSessions(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)

	page := 1
	pageSize := 20
	if p := ctx.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if ps := ctx.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			pageSize = v
		}
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var sessions []models.SecuritySession
	var total int64

	c.DB.Model(&models.SecuritySession{}).Where("user_id = ? AND status = ?", userID, 1).Count(&total)

	offset := (page - 1) * pageSize
	if err := c.DB.Where("user_id = ? AND status = ?", userID, 1).Order("last_active_at DESC").Offset(offset).Limit(pageSize).Find(&sessions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询会话失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      sessions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (c *SecurityController) DeleteSession(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	sessionID := ctx.Param("id")

	var session models.SecuritySession
	if err := c.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "会话不存在"})
		return
	}

	session.Status = 2
	c.DB.Save(&session)

	if c.Redis != nil {
		key := fmt.Sprintf("session:%s:%d", sessionID, userID)
		c.Redis.Client().Del(context.Background(), key)
	}

	c.saveAudit(userID, middleware.GetUsername(ctx), auditArgs{
		eventType:    "session_terminate",
		category:    "authentication",
		severity:    1,
		ip:          ctx.ClientIP(),
		userAgent:   ctx.GetHeader("User-Agent"),
		sessionID:   sessionID,
		status:      1,
		responseCode: http.StatusOK,
		tenantID:    ctx.GetHeader("X-Tenant-ID"),
	})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "会话已终止"})
}

func (c *SecurityController) DeleteAllSessions(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	currentToken := ctx.GetHeader("Authorization")
	if strings.HasPrefix(currentToken, "Bearer ") {
		currentToken = strings.TrimPrefix(currentToken, "Bearer ")
	}

	if err := c.DB.Model(&models.SecuritySession{}).
		Where("user_id = ? AND status = ? AND token != ?", userID, 1, currentToken).
		Updates(map[string]interface{}{"status": 2, "updated_at": time.Now()}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "终止会话失败"})
		return
	}

	c.saveAudit(userID, username, auditArgs{
		eventType:    "session_terminate_all",
		category:    "authentication",
		severity:    1,
		ip:          ctx.ClientIP(),
		userAgent:   ctx.GetHeader("User-Agent"),
		status:      1,
		responseCode: http.StatusOK,
		tenantID:    ctx.GetHeader("X-Tenant-ID"),
	})

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "所有其他会话已终止"})
}

// ============ Security Audit APIs ============

func (c *SecurityController) GetSecurityAudits(ctx *gin.Context) {
	eventType := ctx.Query("event_type")
	severity := ctx.Query("severity")
	username := ctx.Query("username")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	query := c.DB.Model(&models.SecurityAudit{}).Order("created_at DESC")

	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if severity != "" {
		if sev, err := strconv.Atoi(severity); err == nil && sev > 0 {
			query = query.Where("severity = ?", sev)
		}
	}
	if username != "" {
		query = query.Where("username = ?", username)
	}
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02 23:59:59", endDate+" 23:59:59"); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	page := 1
	pageSize := 20
	if p := ctx.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if ps := ctx.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			pageSize = v
		}
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var total int64
	query.Model(&models.SecurityAudit{}).Count(&total)

	offset := (page - 1) * pageSize
	var audits []models.SecurityAudit
	if err := query.Offset(offset).Limit(pageSize).Find(&audits).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询审计日志失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      audits,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (c *SecurityController) GenerateSecurityReport(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	var req struct {
		ReportType  string `json:"report_type" binding:"required"`
		PeriodStart string `json:"period_start"`
		PeriodEnd   string `json:"period_end"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	now := time.Now()
	var periodStart, periodEnd time.Time

	switch req.ReportType {
	case "daily":
		periodStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.Add(24 * time.Hour)
	case "weekly":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		periodStart = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.Add(7 * 24 * time.Hour)
	case "monthly":
		periodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
	case "custom":
		var err1, err2 error
		periodStart, err1 = time.Parse("2006-01-02", req.PeriodStart)
		periodEnd, err2 = time.Parse("2006-01-02 23:59:59", req.PeriodEnd+" 23:59:59")
		if err1 != nil || err2 != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "自定义周期日期格式错误，请使用YYYY-MM-DD"})
			return
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的报告类型"})
		return
	}

	tenantID := middleware.GetTenantIDCtx(ctx)

	var totalEvents, loginSuccess, loginFail, faEnable, faDisable int64
	var faVerifySuccess, faVerifyFail, sessionTerminate, sessionTerminateAll int64
	var sev1, sev2, sev3, sev4 int64
	var criticalEvents []models.SecurityAudit

	baseQuery := c.DB.Model(&models.SecurityAudit{}).Where("created_at >= ? AND created_at <= ?", periodStart, periodEnd)
	if tenantID != "" {
		baseQuery = baseQuery.Where("tenant_id = ?", tenantID)
	}

	baseQuery.Count(&totalEvents)
	baseQuery.Where("event_type = 'login_success'").Count(&loginSuccess)
	baseQuery.Where("event_type = 'login_fail'").Count(&loginFail)
	baseQuery.Where("event_type = '2fa_enable_request'").Count(&faEnable)
	baseQuery.Where("event_type = '2fa_disable'").Count(&faDisable)
	baseQuery.Where("event_type = '2fa_verify_success'").Count(&faVerifySuccess)
	baseQuery.Where("event_type = '2fa_verify_fail'").Count(&faVerifyFail)
	baseQuery.Where("event_type = 'session_terminate'").Count(&sessionTerminate)
	baseQuery.Where("event_type = 'session_terminate_all'").Count(&sessionTerminateAll)
	baseQuery.Where("severity = 1").Count(&sev1)
	baseQuery.Where("severity = 2").Count(&sev2)
	baseQuery.Where("severity = 3").Count(&sev3)
	baseQuery.Where("severity = 4").Count(&sev4)

	c.DB.Where("created_at >= ? AND created_at <= ? AND severity = 4", periodStart, periodEnd).
		Order("created_at DESC").Limit(10).Find(&criticalEvents)

	riskLevel := 1
	if sev4 > 0 || loginFail > 50 {
		riskLevel = 4
	} else if sev3 > 10 || loginFail > 20 {
		riskLevel = 3
	} else if sev2 > 20 || loginFail > 5 {
		riskLevel = 2
	}

	statsMap := map[string]interface{}{
		"total_events":            totalEvents,
		"login_success":          loginSuccess,
		"login_fail":            loginFail,
		"fa_enable":              faEnable,
		"fa_disable":             faDisable,
		"fa_verify_success":      faVerifySuccess,
		"fa_verify_fail":         faVerifyFail,
		"session_terminate":      sessionTerminate,
		"session_terminate_all":   sessionTerminateAll,
		"severity_1_count":      sev1,
		"severity_2_count":      sev2,
		"severity_3_count":      sev3,
		"severity_4_count":      sev4,
		"critical_events_count": len(criticalEvents),
	}
	statsJSON, _ := json.Marshal(statsMap)

	findings := []string{}
	if loginFail > 10 {
		findings = append(findings, fmt.Sprintf("检测到%d次登录失败，存在暴力破解风险", loginFail))
	}
	if faVerifyFail > 5 {
		findings = append(findings, fmt.Sprintf("检测到%d次2FA验证失败", faVerifyFail))
	}
	if sev4 > 0 {
		findings = append(findings, fmt.Sprintf("存在%d个高危安全事件需要立即处理", sev4))
	}
	if len(findings) == 0 {
		findings = append(findings, "未检测到明显安全风险")
	}
	findingsJSON, _ := json.Marshal(findings)

	recommendations := []string{}
	if faEnable == 0 {
		recommendations = append(recommendations, "建议所有用户启用双因素认证(2FA)以提升账户安全")
	}
	if loginFail > 5 {
		recommendations = append(recommendations, "建议启用登录失败锁定策略，防止暴力破解攻击")
	}
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "继续保持良好的安全实践，定期审查会话和审计日志")
	}
	recommendationsJSON, _ := json.Marshal(recommendations)

	report := models.SecurityReport{
		ReportType:      req.ReportType,
		Title:           fmt.Sprintf("安全报告 - %s - %s", req.ReportType, now.Format("2006-01-02")),
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
		Summary:         fmt.Sprintf("统计周期: %s | 总事件数: %d | 登录成功: %d | 登录失败: %d | 高危事件: %d", req.ReportType, totalEvents, loginSuccess, loginFail, sev4),
		Stats:           string(statsJSON),
		RiskLevel:       riskLevel,
		Findings:        string(findingsJSON),
		Recommendations: string(recommendationsJSON),
		Status:          2,
		GeneratedBy:     userID,
		GeneratedAt:     now,
		TenantID:        tenantID,
	}

	if err := c.DB.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成报告失败"})
		return
	}

	c.saveAudit(userID, username, auditArgs{
		eventType:    "security_report_generate",
		category:    "data",
		severity:    1,
		ip:          ctx.ClientIP(),
		userAgent:   ctx.GetHeader("User-Agent"),
		resourceID:  fmt.Sprintf("report_type=%s,period=%s~%s", req.ReportType, periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02")),
		status:      1,
		responseCode: http.StatusOK,
		tenantID:    tenantID,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"report": report,
			"stats":  statsMap,
		},
	})
}

// ============ Helper Functions ============

func (c *SecurityController) saveAudit(userID uint, username string, args auditArgs) {
	audit := models.SecurityAudit{
		EventType:     args.eventType,
		EventCategory: args.category,
		Severity:      args.severity,
		UserID:        userID,
		Username:      username,
		TargetUserID:  args.targetUserID,
		IP:            args.ip,
		UserAgent:     args.userAgent,
		SessionID:     args.sessionID,
		ResourceType:  args.resourceType,
		ResourceID:    args.resourceID,
		Status:        args.status,
		ErrorMsg:      args.errorMsg,
		RequestMethod: args.reqMethod,
		RequestPath:   args.reqPath,
		ResponseCode:   args.responseCode,
		Duration:      args.duration,
		Metadata:      args.metadata,
		TenantID:      args.tenantID,
	}
	c.DB.Create(&audit)
}

func (c *SecurityController) recordLoginAttempt(userID uint, username, ip, action string, status int, reason string) {
	attempt := models.LoginAttempt{
		UserID:   userID,
		Username: username,
		IP:       ip,
		Action:   action,
		Status:   status,
		Reason:   reason,
	}
	c.DB.Create(&attempt)
}

func genRecoveryCodes(n int) []string {
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		bytes := make([]byte, 8)
		rand.Read(bytes)
		codes[i] = strings.ToUpper(hex.EncodeToString(bytes))
	}
	return codes
}
