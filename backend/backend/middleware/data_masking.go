package middleware

import (
	"encoding/json"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
)

// DataMaskingRule 数据脱敏规则
type DataMaskingRule struct {
	ID         uint   `json:"id"`
	Field      string `json:"field"`      // 字段名，如 phone, email, id_card
	Pattern    string `json:"pattern"`    // 正则表达式
	Replacement string `json:"replacement"` // 替换模式，如 $1****$2
	Enabled    bool   `json:"enabled"`    // 是否启用
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// MaskingEngine 脱敏引擎
type MaskingEngine struct {
	rules    map[string]*regexp.Regexp // field -> compiled regex
	replaces map[string]string         // field -> replacement pattern
	mu       sync.RWMutex
}

var (
	defaultEngine     *MaskingEngine
	engineInit        sync.Once
)

// GetMaskingEngine 获取全局脱敏引擎
func GetMaskingEngine() *MaskingEngine {
	engineInit.Do(func() {
		defaultEngine = &MaskingEngine{
			rules:    make(map[string]*regexp.Regexp),
			replaces: make(map[string]string),
		}
		// 注册默认脱敏规则
		defaultEngine.registerDefaultRules()
	})
	return defaultEngine
}

// registerDefaultRules 注册默认脱敏规则
func (e *MaskingEngine) registerDefaultRules() {
	// 手机号：13812345678 -> 138****5678
	e.RegisterRule("phone", `^(\d{3})\d{4}(\d{4})$`, "$1****$2")
	// 身份证：110101199001011234 -> 110101********1234
	e.RegisterRule("id_card", `^(\d{6})\d{8}(\d{4})$`, "$1********$2")
	// 邮箱：user@example.com -> u***@example.com
	e.RegisterRule("email", `^(\w{1})\w+(\w+@\w+\.\w+)$`, "$1***$2")
	// 银行卡：6222021234567890123 -> 622202*********0123
	e.RegisterRule("bank_card", `^(\d{6})\d+(\d{4})$`, "$1*********$2")
	// 地址：只保留省份
	e.RegisterRule("address", `^(.{2}[省市区县]).*$`, "$1***")
	// 真实姓名：显示首尾字符
	e.RegisterRule("real_name", `^(.)(.*)(.)$`, "$1***$3")
}

// RegisterRule 注册脱敏规则
func (e *MaskingEngine) RegisterRule(field, pattern, replacement string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if re, err := regexp.Compile(pattern); err == nil {
		e.rules[field] = re
		e.replaces[field] = replacement
	}
}

// UnregisterRule 移除脱敏规则
func (e *MaskingEngine) UnregisterRule(field string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.rules, field)
	delete(e.replaces, field)
}

// MaskValue 对单个值进行脱敏
func (e *MaskingEngine) MaskValue(field, value string) string {
	if value == "" {
		return value
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	if re, ok := e.rules[field]; ok {
		if repl, ok := e.replaces[field]; ok {
			return re.ReplaceAllString(value, repl)
		}
	}

	// 默认脱敏：保留前3后4位
	if len(value) > 7 {
		return value[:3] + "****" + value[len(value)-4:]
	}
	return value[:1] + "****" + value[len(value)-1:]
}

// MaskMap 对 map 中的敏感字段进行脱敏
func (e *MaskingEngine) MaskMap(data map[string]interface{}, sensitiveFields []string) map[string]interface{} {
	if data == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = v
	}

	for _, field := range sensitiveFields {
		if val, ok := result[field]; ok {
			if strVal, ok := val.(string); ok {
				result[field] = e.MaskValue(field, strVal)
			}
		}
	}

	return result
}

// MaskStructFields 根据字段名列表对结构体进行脱敏（反射实现）
func (e *MaskingEngine) MaskStructFields(data interface{}, fieldNames []string) interface{} {
	if data == nil {
		return nil
	}

	// 转换为 JSON 再处理，保留原始类型
	jsonData, err := json.Marshal(data)
	if err != nil {
		return data
	}

	var m map[string]interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return data
	}

	masked := e.MaskMap(m, fieldNames)

	// 转回原始类型
	result, err := json.Marshal(masked)
	if err != nil {
		return data
	}

	var out interface{}
	if err := json.Unmarshal(result, &out); err != nil {
		return data
	}

	return out
}

// DataMaskingMiddleware 数据脱敏中间件
// 对指定路由的响应进行敏感字段脱敏
func DataMaskingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 继续处理请求
		c.Next()

		// 只对成功响应进行处理
		if c.Writer.Status() >= 400 {
			return
		}

		// 检查是否需要脱敏
		if !shouldMask(c.Request.URL.Path) {
			return
		}

		// 获取响应Writer的原始内容
		// 注意：这里需要在写入响应之前拦截，比较复杂
		// 实际应用中建议在 Controller 层显式调用 MaskMap
	}
}

// shouldMask 判断路径是否需要脱敏
func shouldMask(path string) bool {
	// 定义需要脱敏的路径模式
	maskingPaths := []string{
		"/api/v1/members",
		"/api/v1/users",
		"/api/v1/employees",
		"/api/v1/auth/profile",
	}

	for _, p := range maskingPaths {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}
	return false
}

// CommonSensitiveFields 常用敏感字段列表
var CommonSensitiveFields = []string{
	"phone", "mobile", "telephone",      // 电话
	"id_card", "identity", "identity_card", // 身份证
	"email",                              // 邮箱
	"bank_card", "bank_account", "account", // 银行账号
	"address", "home_address", "work_address", // 地址
	"real_name", "name", "full_name",    // 真实姓名
	"password", "pwd",                   // 密码（通常已用 json:"-" 忽略）
}

// MaskSensitiveData 通用敏感数据脱敏函数
func MaskSensitiveData(data map[string]interface{}) map[string]interface{} {
	engine := GetMaskingEngine()
	return engine.MaskMap(data, CommonSensitiveFields)
}

// MaskMemberData 会员数据脱敏
func MaskMemberData(member map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":            member["id"],
		"member_code":   member["member_code"],
		"member_name":   maskField(member, "member_name", "name"),
		"phone":         maskField(member, "phone", "phone"),
		"gender":        member["gender"],
		"birth_date":    member["birth_date"],
		"email":         maskField(member, "email", "email"),
		"avatar":        member["avatar"],
		"member_level":  member["member_level"],
		"points":        member["points"],
		"balance":       member["balance"],
		"status":        member["status"],
		"source":        member["source"],
		"created_at":    member["created_at"],
	}
}

// maskField 对特定字段进行脱敏
func maskField(data map[string]interface{}, field, maskType string) interface{} {
	engine := GetMaskingEngine()
	if val, ok := data[field].(string); ok {
		return engine.MaskValue(maskType, val)
	}
	return data[field]
}

// MaskUserData 用户数据脱敏
func MaskUserData(user map[string]interface{}) map[string]interface{} {
	engine := GetMaskingEngine()
	result := make(map[string]interface{})

	// 复制所有字段
	for k, v := range user {
		result[k] = v
	}

	// 脱敏敏感字段
	if phone, ok := result["phone"].(string); ok {
		result["phone"] = engine.MaskValue("phone", phone)
	}
	if email, ok := result["email"].(string); ok {
		result["email"] = engine.MaskValue("email", email)
	}
	if idCard, ok := result["id_card"].(string); ok {
		result["id_card"] = engine.MaskValue("id_card", idCard)
	}

	return result
}

// GDPRExportRequest GDPR数据导出请求
type GDPRExportRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	DataType string `json:"data_type"` // all, personal, activity, devices
}

// GDPRDeleteRequest 账户删除请求
type GDPRDeleteRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	ConfirmText string `json:"confirm_text" binding:"required"` // 需要输入 "DELETE" 确认
	Reason      string `json:"reason"`                          // 删除原因
}

// GDPRExportResponse GDPR数据导出响应
type GDPRExportResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
