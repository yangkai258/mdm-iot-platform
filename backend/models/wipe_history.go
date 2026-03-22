package models

import (
	"time"
)

// WipeHistory 设备擦除历史记录表
type WipeHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DeviceID     string    `gorm:"type:varchar(64);index;not null" json:"device_id"`
	OperatorID   uint      `gorm:"not null" json:"operator_id"`                       // 操作人ID
	OperatorName string    `gorm:"type:varchar(100)" json:"operator_name"`             // 操作人名称 (冗余存储)
	WipeType     string    `gorm:"type:varchar(32)" json:"wipe_type"`                 // full/selective
	Status       string    `gorm:"type:varchar(20);index;default:pending" json:"status"` // pending/executing/completed/failed
	ConfirmToken string    `gorm:"type:varchar(64)" json:"confirm_token"`              // 二次确认token
	ConfirmedAt  *time.Time `json:"confirmed_at"`                                      // 确认时间
	ExecutedAt   *time.Time `json:"executed_at"`                                       // 执行时间
	CompletedAt  *time.Time `json:"completed_at"`                                      // 完成时间
	Result       string    `gorm:"type:text" json:"result"`                            // 操作结果/错误信息
	Reason       string    `gorm:"type:text" json:"reason"`                            // 擦除原因
	TenantID     string    `gorm:"type:varchar(50);index" json:"tenant_id"`            // 租户ID
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (WipeHistory) TableName() string {
	return "wipe_history"
}

// WipeRequest 擦除请求
type WipeRequest struct {
	WipeType string `json:"wipe_type" binding:"required"` // full/selective
	Reason   string `json:"reason"`                        // 擦除原因
}

// WipeConfirmRequest 擦除确认请求
type WipeConfirmRequest struct {
	ConfirmToken string `json:"confirm_token" binding:"required"` // 确认token
}

// WipeHistoryResponse 擦除历史响应
type WipeHistoryResponse struct {
	ID           uint       `json:"id"`
	DeviceID     string     `json:"device_id"`
	OperatorID   uint       `json:"operator_id"`
	OperatorName string     `json:"operator_name"`
	WipeType     string     `json:"wipe_type"`
	Status       string     `json:"status"`
	ConfirmedAt  *time.Time `json:"confirmed_at,omitempty"`
	ExecutedAt   *time.Time `json:"executed_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Result       string     `json:"result"`
	Reason       string     `json:"reason"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ToResponse 转换为响应结构
func (w *WipeHistory) ToResponse() WipeHistoryResponse {
	return WipeHistoryResponse{
		ID:           w.ID,
		DeviceID:     w.DeviceID,
		OperatorID:   w.OperatorID,
		OperatorName: w.OperatorName,
		WipeType:     w.WipeType,
		Status:       w.Status,
		ConfirmedAt:  w.ConfirmedAt,
		ExecutedAt:   w.ExecutedAt,
		CompletedAt:  w.CompletedAt,
		Result:       w.Result,
		Reason:       w.Reason,
		CreatedAt:    w.CreatedAt,
	}
}
