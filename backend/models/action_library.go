package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActionLibrary 动作库表
type ActionLibrary struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ActionID         string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"action_id"`
	ActionName       string         `gorm:"type:varchar(64);not null" json:"action_name"`
	ActionNameEn     string         `gorm:"type:varchar(64)" json:"action_name_en"`
	Category         string         `gorm:"type:varchar(32);not null;index" json:"category"`
	Description      string         `gorm:"type:text" json:"description"`
	DurationMs       int            `gorm:"type:int;not null" json:"duration_ms"`
	Priority         int            `gorm:"type:int;default:5" json:"priority"`
	IsEmergency      bool           `gorm:"type:boolean;default:false" json:"is_emergency"`
	CompatibleModels string         `gorm:"type:jsonb;not null" json:"compatible_models"`
	Parameters       string         `gorm:"type:jsonb;default:'{}'" json:"parameters"`
	AnimationData    string         `gorm:"type:jsonb;default:'{}'" json:"animation_data"`
	MotorCommands    string         `gorm:"type:jsonb;default:'{}'" json:"motor_commands"`
	AudioFile        string         `gorm:"type:varchar(256)" json:"audio_file"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ActionLibrary) TableName() string {
	return "embodied_action_library"
}

// BeforeCreate 创建前自动生成 UUID
func (a *ActionLibrary) BeforeCreate(tx *gorm.DB) error {
	if a.ActionID == "" {
		a.ActionID = uuid.New().String()
	}
	return nil
}

// ActionCategory 常量
const (
	ActionCategoryEmotion  = "emotion"  // 情绪动作
	ActionCategoryGreeting = "greeting" // 问候动作
	ActionCategoryPlay     = "play"     // 玩耍动作
	ActionCategoryUtility  = "utility"  // 功能动作
)

// ActionLibraryResponse 动作库响应
type ActionLibraryResponse struct {
	ActionID         string `json:"action_id"`
	ActionName       string `json:"action_name"`
	ActionNameEn     string `json:"action_name_en"`
	Category         string `json:"category"`
	Description      string `json:"description"`
	DurationMs       int    `json:"duration_ms"`
	Priority         int    `json:"priority"`
	IsEmergency      bool   `json:"is_emergency"`
	CompatibleModels string `json:"compatible_models"`
	Parameters       string `json:"parameters"`
	AnimationData    string `json:"animation_data"`
	MotorCommands    string `json:"motor_commands"`
	AudioFile        string `json:"audio_file"`
	CreatedAt        string `json:"created_at"`
}

// ToResponse 转换为响应格式
func (a *ActionLibrary) ToResponse() *ActionLibraryResponse {
	return &ActionLibraryResponse{
		ActionID:         a.ActionID,
		ActionName:       a.ActionName,
		ActionNameEn:     a.ActionNameEn,
		Category:         a.Category,
		Description:      a.Description,
		DurationMs:       a.DurationMs,
		Priority:         a.Priority,
		IsEmergency:      a.IsEmergency,
		CompatibleModels: a.CompatibleModels,
		Parameters:       a.Parameters,
		AnimationData:    a.AnimationData,
		MotorCommands:    a.MotorCommands,
		AudioFile:        a.AudioFile,
		CreatedAt:        a.CreatedAt.Format(time.RFC3339),
	}
}

// ActionExecuteRequest 动作执行请求
type ActionExecuteRequest struct {
	ActionID   string                 `json:"action_id" binding:"required"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ActionExecuteResponse 动作执行响应
type ActionExecuteResponse struct {
	ActionID string `json:"action_id"`
	Status   string `json:"status"` // pending, sent, executed
	Message  string `json:"message"`
}
