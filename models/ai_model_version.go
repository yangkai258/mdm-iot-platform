package models

import (
	"time"
)

// AIModelVersion AI模型版本
type AIModelVersion struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ModelID      string    `json:"model_id" gorm:"uniqueIndex:idx_model_version;size:64"`
	ModelName    string    `json:"model_name" gorm:"size:64"`
	Version      string    `json:"version" gorm:"uniqueIndex:idx_model_version;size:32"`
	Status       string    `json:"status" gorm:"size:20"` // testing/staging/production/deprecated
	ModelPath    string    `json:"model_path" gorm:"size:512"`
	Config       string    `json:"config" gorm:"type:jsonb"`
	Metrics      string    `json:"metrics" gorm:"type:jsonb"`
	RollbackFrom string    `json:"rollback_from" gorm:"size:32"`
	PublishedBy  uint      `json:"published_by"`
	PublishedAt  time.Time `json:"published_at"`
	DeprecatedAt time.Time `json:"deprecated_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 表名
func (AIModelVersion) TableName() string {
	return "ai_model_versions"
}

// AIModel AI模型（注册新模型时使用）
type AIModel struct {
	ModelID   string `json:"model_id" binding:"required" gorm:"size:64"`
	ModelName string `json:"model_name" binding:"required" gorm:"size:64"`
	ModelType string `json:"model_type" gorm:"size:30"` // openai/anthropic/local
}

// AIModelVersionCreate 创建版本请求
type AIModelVersionCreate struct {
	Version   string `json:"version" binding:"required" gorm:"size:32"`
	ModelPath string `json:"model_path" gorm:"size:512"`
	Config    string `json:"config"`
	Metrics   string `json:"metrics"`
}

// AIModelVersionUpdate 更新版本请求
type AIModelVersionUpdate struct {
	ModelName string `json:"model_name" gorm:"size:64"`
	Status    string `json:"status" gorm:"size:20"`
	ModelPath string `json:"model_path" gorm:"size:512"`
	Config    string `json:"config" gorm:"type:jsonb"`
	Metrics   string `json:"metrics" gorm:"type:jsonb"`
}

// AIModelRollback 回滚请求
type AIModelRollback struct {
	TargetVersion string `json:"target_version" binding:"required"`
	Reason        string `json:"reason" gorm:"type:text"`
}

// SandboxTest 沙箱测试任务
type SandboxTest struct {
	TaskID      string    `json:"task_id" gorm:"uniqueIndex;size:64"`
	ModelID     string    `json:"model_id" gorm:"index;size:64"`
	TestDataID  string    `json:"test_data_id" gorm:"size:64"`
	TestType    string    `json:"test_type" gorm:"size:30"` // unit/integration/stress
	TestName    string    `json:"test_name" gorm:"size:100"`
	TestCases   string    `json:"test_cases" gorm:"type:jsonb"`
	Status      string    `json:"status" gorm:"size:20"` // pending/running/completed/failed
	Result      string    `json:"result" gorm:"type:jsonb"`
	ReportPath  string    `json:"report_path" gorm:"size:512"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 表名
func (SandboxTest) TableName() string {
	return "ai_sandbox_tests"
}

// SandboxTestCreate 创建沙箱测试请求
type SandboxTestCreate struct {
	ModelID    string `json:"model_id" binding:"required" gorm:"size:64"`
	TestDataID string `json:"test_data_id" gorm:"size:64"`
	TestType   string `json:"test_type" binding:"required" gorm:"size:30"`
	TestName   string `json:"test_name" binding:"required" gorm:"size:100"`
	TestCases  string `json:"test_cases" binding:"required"`
}
