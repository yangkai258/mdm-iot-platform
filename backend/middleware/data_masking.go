package middleware

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// DataMasking 数据脱敏中间件
type DataMasking struct {
	rules map[string]MaskRule
}

// MaskRule 脱敏规则
type MaskRule struct {
	Pattern *regexp.Regexp
	Replace string // 如 "***" 或 "1***9"
}

// 默认脱敏规则
var defaultMasking = &DataMasking{
	rules: map[string]MaskRule{
		"phone": {
			Pattern: regexp.MustCompile(`(\d{3})\d{4}(\d{4})`),
			Replace: "$1****$2",
		},
		"email": {
			Pattern: regexp.MustCompile(`(\w)[\w.-]*@([\w.-]+\.\w+)`),
			Replace: "$1***@$2",
		},
		"id_card": {
			Pattern: regexp.MustCompile(`(\d{6})\d{8}(\d{4})`),
			Replace: "$1********$2",
		},
		"bank_card": {
			Pattern: regexp.MustCompile(`(\d{4})\d+(\d{4})`),
			Replace: "$1****$2",
		},
		"name": {
			Pattern: regexp.MustCompile(`(.)\S+(\S*)`),
			Replace: "$1*$2",
		},
	},
}

// MaskData 根据字段名脱敏
func (m *DataMasking) MaskData(data interface{}, fields []string) interface{} {
	if m == nil {
		m = defaultMasking
	}

	result := make(map[string]interface{})
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	for _, field := range fields {
		if val, ok := result[field]; ok {
			result[field] = m.maskByType(val, field)
		}
	}
	return result
}

func (m *DataMasking) maskByType(value interface{}, fieldType string) string {
	str := fmt.Sprintf("%v", value)
	if rule, ok := m.rules[fieldType]; ok {
		return rule.Pattern.ReplaceAllString(str, rule.Replace)
	}
	// 默认脱敏：显示首尾各1个字符
	if len(str) > 2 {
		return str[:1] + strings.Repeat("*", len(str)-2) + str[len(str)-1:]
	}
	return strings.Repeat("*", len(str))
}

// MaskSensitiveFields 通用脱敏方法
func MaskSensitiveFields(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	sensitiveFields := []string{"phone", "email", "id_card", "bank_card", "real_name", "address", "contact", "tax_no"}

	for _, field := range sensitiveFields {
		if val, ok := result[field]; ok {
			result[field] = defaultMasking.maskByType(val, field)
		}
	}
	return result
}
