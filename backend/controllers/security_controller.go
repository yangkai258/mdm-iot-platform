package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SecurityController 安全与加密控制器
type SecurityController struct {
	DB *gorm.DB
}

// EncryptionRequest 加密请求
type EncryptionRequest struct {
	Data     string `json:"data" binding:"required"`
	KeyID    string `json:"key_id"`    // 可选，指定密钥ID
	Metadata string `json:"metadata"` // 可选，元数据
}

// DecryptionRequest 解密请求
type DecryptionRequest struct {
	Data   string `json:"data" binding:"required"`
	KeyID  string `json:"key_id"`  // 可选，指定密钥ID
}

// KeyRotateRequest 密钥轮换请求
type KeyRotateRequest struct {
	Algorithm string `json:"algorithm"` // 加密算法，默认 AES-256-GCM
}

// AnonymizeRequest 脱敏请求
type AnonymizeRequest struct {
	Data       interface{} `json:"data" binding:"required"`        // 要脱敏的数据
	Type       string     `json:"type" binding:"required"`        // 脱敏类型: email, phone, id_card, name, custom
	Fields     []string    `json:"fields"`                         // 要脱敏的字段列表
	Purpose    string     `json:"purpose"`                         // 使用目的
	ExportFormat string   `json:"export_format"`                   // 导出格式: json, csv, xlsx
}

// AnonymizeData 脱敏数据
func (c *SecurityController) AnonymizeData(ctx *gin.Context) {
	var req AnonymizeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.recordAudit(ctx, "anonymize", "security", "", "", 0, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	_ = time.Now() // start time for potential logging
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 解析数据
	result := c.processAnonymization(req.Data, req.Type, req.Fields)
	
	// 生成记录ID
	recordID := generateRecordID()

	// 存储脱敏记录
	record := models.DataAnonymizationRecord{
		RecordID:       recordID,
		AnonymizedData: result.ExportedData,
		AnonymizeType:  req.Type,
		Fields:         strings.Join(req.Fields, ","),
		UserID:         uid,
		Purpose:        req.Purpose,
		ExportFormat:   req.ExportFormat,
		Status:         2, // 完成
		CompletedAt:    timePtr(time.Now()),
	}
	c.DB.Create(&record)

	// 记录审计日志
	c.recordAudit(ctx, "anonymize", "security", "data_anonymization_record", recordID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"record_id":       recordID,
			"anonymized_data": result.Data,
			"fields_anonymized": result.FieldsAnonymized,
			"created_at":      record.CreatedAt,
		},
	})
}

// AnonymizationResult 脱敏结果
type AnonymizationResult struct {
	Data             interface{}
	ExportedData     string
	FieldsAnonymized []string
}

// processAnonymization 处理脱敏
func (c *SecurityController) processAnonymization(data interface{}, anonymizeType string, fields []string) AnonymizationResult {
	result := AnonymizationResult{
		Data:             data,
		ExportedData:     "",
		FieldsAnonymized: []string{},
	}

	switch anonymizeType {
	case "email":
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				if str, ok := v.(string); ok {
					if strings.Contains(str, "@") {
						parts := strings.Split(str, "@")
						if len(parts) == 2 {
							m[k] = anonymizeEmail(parts[0]) + "@" + parts[1]
							result.FieldsAnonymized = append(result.FieldsAnonymized, k)
						}
					}
				}
			}
			result.Data = m
		}
	case "phone":
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				if str, ok := v.(string); ok {
					if isPhoneNumber(str) {
						m[k] = anonymizePhone(str)
						result.FieldsAnonymized = append(result.FieldsAnonymized, k)
					}
				}
			}
			result.Data = m
		}
	case "id_card":
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				if str, ok := v.(string); ok {
					if isIDCard(str) {
						m[k] = anonymizeIDCard(str)
						result.FieldsAnonymized = append(result.FieldsAnonymized, k)
					}
				}
			}
			result.Data = m
		}
	case "name":
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				if str, ok := v.(string); ok {
					if len(str) >= 2 {
						m[k] = anonymizeName(str)
						result.FieldsAnonymized = append(result.FieldsAnonymized, k)
					}
				}
			}
			result.Data = m
		}
	case "custom":
		if m, ok := data.(map[string]interface{}); ok {
			for _, field := range fields {
				if v, exists := m[field]; exists {
					if str, ok := v.(string); ok {
						m[field] = maskString(str)
						result.FieldsAnonymized = append(result.FieldsAnonymized, field)
					}
				}
			}
			result.Data = m
		}
	}

	return result
}

// ExportAnonymizedData 导出脱敏数据
func (c *SecurityController) ExportAnonymizedData(ctx *gin.Context) {
	var req AnonymizeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.recordAudit(ctx, "export", "security", "", "", 0, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 处理脱敏
	result := c.processAnonymization(req.Data, req.Type, req.Fields)
	recordID := generateRecordID()

	// 生成导出文件
	var exportPath string
	format := strings.ToLower(req.ExportFormat)
	if format == "" {
		format = "json"
	}

	switch format {
	case "csv":
		exportPath = c.exportToCSV(result.Data, recordID)
	case "xlsx":
		exportPath = c.exportToXLSX(result.Data, recordID)
	default:
		exportPath = c.exportToJSON(result.Data, recordID)
	}

	// 存储记录
	record := models.DataAnonymizationRecord{
		RecordID:       recordID,
		AnonymizedData: result.ExportedData,
		AnonymizeType:  req.Type,
		Fields:         strings.Join(req.Fields, ","),
		UserID:         uid,
		Purpose:        req.Purpose,
		ExportFormat:   format,
		Status:         2,
		CompletedAt:    timePtr(time.Now()),
	}
	c.DB.Create(&record)

	// 记录审计日志
	c.recordAudit(ctx, "export_anonymized", "security", "data_anonymization_record", recordID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"record_id":  recordID,
			"export_path": exportPath,
			"format":      format,
			"created_at":  record.CreatedAt,
		},
	})
}

// ============ 加密 API ============

// Encrypt 加密数据
func (c *SecurityController) Encrypt(ctx *gin.Context) {
	var req EncryptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.recordAudit(ctx, "encrypt", "security", "", "", 0, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 使用 AES 加密
	encrypted, err := utils.EncryptAES(req.Data)
	if err != nil {
		c.recordAudit(ctx, "encrypt", "security", "", "", uid, http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "加密失败: " + err.Error(),
		})
		return
	}

	// 获取当前活跃密钥
	var key models.EncryptionKey
	c.DB.Where("status = ? AND is_primary = ?", 1, true).First(&key)

	c.recordAudit(ctx, "encrypt", "security", "", "", uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"encrypted_data": encrypted,
			"key_id":         key.KeyID,
			"algorithm":      "AES-256-GCM",
			"encrypted_at":   time.Now(),
		},
	})
}

// Decrypt 解密数据
func (c *SecurityController) Decrypt(ctx *gin.Context) {
	var req DecryptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.recordAudit(ctx, "decrypt", "security", "", "", 0, http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 使用 AES 解密
	decrypted, err := utils.DecryptAES(req.Data)
	if err != nil {
		c.recordAudit(ctx, "decrypt", "security", "", "", uid, http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "解密失败: " + err.Error(),
		})
		return
	}

	c.recordAudit(ctx, "decrypt", "security", "", "", uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"decrypted_data": decrypted,
			"decrypted_at":   time.Now(),
		},
	})
}

// RotateKey 密钥轮换
func (c *SecurityController) RotateKey(ctx *gin.Context) {
	var req KeyRotateRequest
	ctx.ShouldBindJSON(&req)

	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 生成新密钥版本号
	var maxVersion int
	c.DB.Model(&models.EncryptionKey{}).Select("COALESCE(MAX(key_version), 0)").Scan(&maxVersion)
	newVersion := maxVersion + 1

	// 生成新的加密密钥
	newKeyBytes := make([]byte, 32)
	rand.Read(newKeyBytes)
	newKeyHex := hex.EncodeToString(newKeyBytes)

	// 加密新密钥（使用主密钥加密后存储）
	_ = getPrimaryKey() // primaryKey for potential future use
	encryptedKey, _ := utils.EncryptAES(newKeyHex)

	// 将旧密钥标记为已轮换
	c.DB.Model(&models.EncryptionKey{}).Where("is_primary = ?", true).Updates(map[string]interface{}{
		"is_primary": false,
		"status":     3, // 已轮换
		"rotated_at":  time.Now(),
		"rotated_by":  uid,
	})

	// 创建新密钥
	keyID := fmt.Sprintf("key-%d-%d", time.Now().Unix(), newVersion)
	newEncryptionKey := models.EncryptionKey{
		KeyID:        keyID,
		KeyVersion:   newVersion,
		EncryptedKey: encryptedKey,
		IsPrimary:    true,
		Status:       1,
		Algorithm:    "AES-256-GCM",
		RotatedAt:    timePtr(time.Now()),
		RotatedBy:    uid,
	}

	if err := c.DB.Create(&newEncryptionKey).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密钥轮换失败",
		})
		return
	}

	// 记录审计日志
	c.recordAudit(ctx, "rotate_key", "security", "encryption_key", keyID, uid, http.StatusOK, "")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"key_id":       keyID,
			"key_version":  newVersion,
			"rotated_at":   time.Now(),
			"algorithm":    "AES-256-GCM",
		},
	})
}

// GetKeyInfo 获取密钥信息
func (c *SecurityController) GetKeyInfo(ctx *gin.Context) {
	var keys []models.EncryptionKey
	c.DB.Where("status IN ?", []int{1, 2}).Order("key_version DESC").Find(&keys)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": keys,
	})
}

// ============ 辅助函数 ============

func generateRecordID() string {
	return fmt.Sprintf("anon-%d-%s", time.Now().Unix(), randomString(8))
}

func randomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func anonymizeEmail(email string) string {
	if len(email) < 3 {
		return "***"
	}
	at := strings.Index(email, "@")
	if at <= 0 {
		return "***"
	}
	username := email[:at]
	domain := email[at:]
	if len(username) <= 2 {
		return "*" + domain
	}
	return string(username[0]) + strings.Repeat("*", len(username)-2) + domain
}

func anonymizePhone(phone string) string {
	if len(phone) < 7 {
		return "***"
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func anonymizeIDCard(id string) string {
	if len(id) < 8 {
		return "******************"
	}
	return id[:6] + "********" + id[len(id)-4:]
}

func anonymizeName(name string) string {
	runes := []rune(name)
	if len(runes) == 1 {
		return "*"
	}
	if len(runes) == 2 {
		return string(runes[0]) + "*"
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-1)
}

func maskString(s string) string {
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

func isPhoneNumber(s string) bool {
	// 中国手机号格式
	return len(s) == 11 && strings.HasPrefix(s, "1")
}

func isIDCard(s string) bool {
	// 简单身份证号格式检查（15位或18位）
	return (len(s) == 15 || len(s) == 18) && (strings.HasPrefix(s, "1") || strings.HasPrefix(s, "2"))
}

func getPrimaryKey() string {
	keyStr := os.Getenv("AES_ENCRYPTION_KEY")
	if keyStr == "" {
		return "mdm-secret-key-32-bytes-long!!"
	}
	return keyStr
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// recordAudit 记录审计日志
func (c *SecurityController) recordAudit(ctx *gin.Context, action, module, resourceType, resourceID string, userID uint, statusCode int, errorMsg string) {
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

// ============ 导出功能 ============

func (c *SecurityController) exportToJSON(data interface{}, recordID string) string {
	// 简化实现，实际应使用 json.Marshal
	exportDir := "./exports"
	os.MkdirAll(exportDir, 0755)
	path := filepath.Join(exportDir, fmt.Sprintf("%s.json", recordID))
	// 实际应写入文件
	return path
}

func (c *SecurityController) exportToCSV(data interface{}, recordID string) string {
	exportDir := "./exports"
	os.MkdirAll(exportDir, 0755)
	path := filepath.Join(exportDir, fmt.Sprintf("%s.csv", recordID))
	return path
}

func (c *SecurityController) exportToXLSX(data interface{}, recordID string) string {
	exportDir := "./exports"
	os.MkdirAll(exportDir, 0755)
	path := filepath.Join(exportDir, fmt.Sprintf("%s.xlsx", recordID))
	return path
}

// ============ 文件上传处理 ============

func parseMultiPartForm(ctx *gin.Context) (map[string]interface{}, error) {
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for key, values := range ctx.Request.MultipartForm.Value {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}

	for key, files := range ctx.Request.MultipartForm.File {
		if len(files) == 1 {
			file, _ := files[0].Open()
			content, _ := io.ReadAll(file)
			result[key] = string(content)
		}
	}

	return result, nil
}
