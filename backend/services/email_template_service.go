package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// EmailTemplateService 邮件模板服务
type EmailTemplateService struct {
	DB *gorm.DB
}

// NewEmailTemplateService 创建邮件模板服务实例
func NewEmailTemplateService(db *gorm.DB) *EmailTemplateService {
	return &EmailTemplateService{DB: db}
}

// Create 创建邮件模板
func (s *EmailTemplateService) Create(tpl *models.EmailTemplate) error {
	if err := s.DB.Create(tpl).Error; err != nil {
		return fmt.Errorf("创建邮件模板失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取模板
func (s *EmailTemplateService) GetByID(id uint) (*models.EmailTemplate, error) {
	var tpl models.EmailTemplate
	if err := s.DB.First(&tpl, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("模板不存在")
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}
	return &tpl, nil
}

// GetByCode 根据编码获取模板
func (s *EmailTemplateService) GetByCode(code string) (*models.EmailTemplate, error) {
	var tpl models.EmailTemplate
	if err := s.DB.Where("code = ? AND status = ?", code, 1).First(&tpl).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("模板不存在或已禁用: %s", code)
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}
	return &tpl, nil
}

// List 查询模板列表（支持分页和过滤）
func (s *EmailTemplateService) List(tenantID, keyword string, status *int, page, pageSize int) ([]models.EmailTemplate, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	query := s.DB.Model(&models.EmailTemplate{})

	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR subject LIKE ?", keyword, keyword, keyword)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计模板数量失败: %w", err)
	}

	var list []models.EmailTemplate
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询模板列表失败: %w", err)
	}

	return list, total, nil
}

// Update 更新邮件模板
func (s *EmailTemplateService) Update(id uint, updates map[string]interface{}) error {
	result := s.DB.Model(&models.EmailTemplate{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新模板失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("模板不存在")
	}
	return nil
}

// Delete 删除邮件模板（软删除：status=0）
func (s *EmailTemplateService) Delete(id uint) error {
	result := s.DB.Model(&models.EmailTemplate{}).Where("id = ?", id).Update("status", 0)
	if result.Error != nil {
		return fmt.Errorf("删除模板失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("模板不存在")
	}
	return nil
}

// Render 渲染邮件内容，将变量替换为实际值
// vars 格式: map[string]interface{}{"username": "张三", "device_name": "设备A"}
func (s *EmailTemplateService) Render(tpl *models.EmailTemplate, vars map[string]interface{}) (subject, body string, err error) {
	if tpl == nil {
		return "", "", fmt.Errorf("模板为空")
	}

	subject, err = s.replaceVariables(tpl.Subject, vars)
	if err != nil {
		return "", "", fmt.Errorf("渲染主题失败: %w", err)
	}

	body, err = s.replaceVariables(tpl.Body, vars)
	if err != nil {
		return "", "", fmt.Errorf("渲染正文失败: %w", err)
	}

	return subject, body, nil
}

// replaceVariables 将模板字符串中的 {{.varName}} 替换为实际值
func (s *EmailTemplateService) replaceVariables(content string, vars map[string]interface{}) (string, error) {
	if content == "" {
		return "", nil
	}

	// 匹配 {{.variableName}} 格式
	re := regexp.MustCompile(`\{\{\s*\.\w+\s*\}\}`)
	result := re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取变量名: {{.username}} -> username
		reInner := regexp.MustCompile(`\{\{\s*\.(\w+)\s*\}\}`)
		matches := reInner.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		varName := matches[1]

		if val, ok := vars[varName]; ok {
			return fmt.Sprintf("%v", val)
		}
		// 变量不存在，保留原样或返回空字符串
		return match
	})

	return result, nil
}

// GetVariables 解析模板变量列表，返回变量名数组
func (s *EmailTemplateService) GetVariables(tpl *models.EmailTemplate) ([]string, error) {
	if tpl.Variables == "" {
		return nil, nil
	}

	var vars []string
	if err := json.Unmarshal([]byte(tpl.Variables), &vars); err != nil {
		return nil, fmt.Errorf("解析变量列表失败: %w", err)
	}
	return vars, nil
}

// ValidateVariables 检查必填变量是否都在 vars 中
func (s *EmailTemplateService) ValidateVariables(tpl *models.EmailTemplate, vars map[string]interface{}) (missing []string) {
	tplVars, err := s.GetVariables(tpl)
	if err != nil || tplVars == nil {
		return
	}
	for _, v := range tplVars {
		if _, ok := vars[v]; !ok {
			missing = append(missing, v)
		}
	}
	return
}

// SendEmail 渲染并发送邮件（需要配置 SMTP）
// 如需真实发送，可在此集成 email 包（如 gopkg.in/gomail.v2）
func (s *EmailTemplateService) SendEmail(code string, to []string, vars map[string]interface{}) error {
	tpl, err := s.GetByCode(code)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}

	subject, body, err := s.Render(tpl, vars)
	if err != nil {
		return fmt.Errorf("渲染邮件失败: %w", err)
	}

	// TODO: 集成真实邮件发送（gomail 等）
	// 这里仅打印，实际发送需配合 SMTP 配置
	fmt.Printf("[EmailTemplate] To: %s | Subject: %s | Body: %s\n",
		strings.Join(to, ","), subject, body)
	return nil
}

// InitDefaultTemplates 初始化默认邮件模板（仅当表为空时）
func (s *EmailTemplateService) InitDefaultTemplates() error {
	var count int64
	s.DB.Model(&models.EmailTemplate{}).Count(&count)
	if count > 0 {
		return nil
	}

	defaults := []models.EmailTemplate{
		{
			Name:      "设备告警通知",
			Code:      "device_alert",
			Subject:   "【{{.system_name}}】设备 {{.device_name}} 发生告警",
			Body:      "尊敬的 {{.username}}：\n\n您的设备 \"{{.device_name}}\" (MAC: {{.mac_address}}) 于 {{.alert_time}} 发生告警。\n\n告警类型：{{.alert_type}}\n告警级别：{{.alert_level}}\n告警内容：{{.alert_message}}\n\n请及时登录 MDM 管控中心查看详情。\n\n-- {{.system_name}}",
			Variables: `["username","system_name","device_name","mac_address","alert_time","alert_type","alert_level","alert_message"]`,
			Status:    1,
			Remark:    "设备告警通知邮件模板",
		},
		{
			Name:      "OTA升级通知",
			Code:      "ota_upgrade",
			Subject:   "【{{.system_name}}】设备 {{.device_name}} 可升级至 {{.target_version}}",
			Body:      "尊敬的 {{.username}}：\n\n设备 \"{{.device_name}}\" 有新的固件版本可升级。\n\n当前版本：{{.current_version}}\n目标版本：{{.target_version}}\n升级说明：{{.upgrade_note}}\n\n请登录 MDM 管控中心确认升级。\n\n-- {{.system_name}}",
			Variables: `["username","system_name","device_name","current_version","target_version","upgrade_note"]`,
			Status:    1,
			Remark:    "OTA升级通知邮件模板",
		},
		{
			Name:      "新设备绑定通知",
			Code:      "device_bind",
			Subject:   "【{{.system_name}}】新设备 \"{{.device_name}}\" 已绑定",
			Body:      "尊敬的 {{.username}}：\n\n您的新设备已成功绑定至 MDM 管控中心。\n\n设备名称：{{.device_name}}\n设备型号：{{.hardware_model}}\n绑定时间：{{.bind_time}}\n\n-- {{.system_name}}",
			Variables: `["username","system_name","device_name","hardware_model","bind_time"]`,
			Status:    1,
			Remark:    "新设备绑定通知邮件模板",
		},
	}

	for i := range defaults {
		if err := s.DB.Create(&defaults[i]).Error; err != nil {
			return fmt.Errorf("初始化默认邮件模板失败: %w", err)
		}
	}
	return nil
}
