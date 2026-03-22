package models

import (
	"time"
)

// AlertRule 告警规则 (DeviceAlertRule 的别名)
type AlertRule = DeviceAlertRule

// AlertRuleResponse 告警规则响应
type AlertRuleResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	DeviceID   string    `json:"device_id"`
	AlertType  string    `json:"alert_type"`
	Condition  string    `json:"condition"`
	Threshold  float64   `json:"threshold"`
	Severity   int       `json:"severity"`
	Enabled    bool      `json:"enabled"`
	NotifyWays string    `json:"notify_ways"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AlertRuleCreateRequest 创建告警规则请求
type AlertRuleCreateRequest struct {
	Name       string                 `json:"name" binding:"required"`
	RuleName   string                 `json:"rule_name"`
	RuleType   string                 `json:"rule_type"`
	DeviceID   string                 `json:"device_id"`
	AlertType  string                 `json:"alert_type" binding:"required"`
	Condition  string                 `json:"condition" binding:"required"`
	Conditions interface{}            `json:"conditions"` // JSON conditions
	Actions    interface{}            `json:"actions"`
	Threshold  float64                `json:"threshold" binding:"required"`
	Priority   int                    `json:"priority"`
	Severity   int                    `json:"severity"`
	Enabled    *bool                  `json:"enabled"`
	NotifyWays string                 `json:"notify_ways"`
	Remark     string                 `json:"remark"`
}

// AlertRuleUpdateRequest 更新告警规则请求
type AlertRuleUpdateRequest struct {
	RuleName   string                 `json:"rule_name"`
	RuleType   string                 `json:"rule_type"`
	DeviceID   string                 `json:"device_id"`
	Conditions interface{}            `json:"conditions"`
	Actions    interface{}            `json:"actions"`
	Enabled    *bool                  `json:"enabled"`
	Priority   int                    `json:"priority"`
}

// ToResponse 转换为响应结构
func (r *DeviceAlertRule) ToResponse() *AlertRuleResponse {
	return &AlertRuleResponse{
		ID:         r.ID,
		Name:       r.Name,
		DeviceID:   r.DeviceID,
		AlertType:  r.AlertType,
		Condition:  r.Condition,
		Threshold:  r.Threshold,
		Severity:   r.Severity,
		Enabled:    r.Enabled,
		NotifyWays: r.NotifyWays,
		Remark:     r.Remark,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}
