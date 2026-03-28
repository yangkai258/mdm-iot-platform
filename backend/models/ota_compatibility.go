package models

import (
	"time"

	"gorm.io/gorm"
)

// OTACompatibilityMatrix OTA固件兼容性矩阵
type OTACompatibilityMatrix struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	MatrixID     string         `gorm:"type:varchar(64);uniqueIndex" json:"matrix_id"` // 矩阵唯一标识

	// 硬件信息
	HardwareModel string         `gorm:"type:varchar(64);index" json:"hardware_model"`   // 硬件型号
	HardwareVersion string       `gorm:"type:varchar(50)" json:"hardware_version"`      // 硬件版本

	// 当前固件
	FromFirmware string         `gorm:"type:varchar(50);index" json:"from_firmware"`    // 当前固件版本

	// 目标固件
	ToFirmware   string         `gorm:"type:varchar(50);index" json:"to_firmware"`      // 目标固件版本

	// 兼容性状态
	CompatibilityStatus string    `gorm:"type:varchar(20);default:'compatible'" json:"status"` // compatible/incompatible/conditional/unknown
	CompatibilityScore float64    `gorm:"type:decimal(5,2);default:100" json:"score"` // 兼容性评分 0-100

	// 限制条件
	MinBatteryLevel int          `gorm:"default:30" json:"min_battery_level"`       // 最低电量要求
	MinStorageKB   int          `gorm:"default:1024" json:"min_storage_kb"`       // 最低存储KB
	MinMemoryKB    int          `gorm:"default:512" json:"min_memory_kb"`         // 最低内存KB
	NetworkRequired bool         `gorm:"default:false" json:"network_required"`     // 是否需要网络

	// 限制说明
	Constraints   string         `gorm:"type:text" json:"constraints"`             // 限制条件说明
	Warning      string         `gorm:"type:text" json:"warning"`               // 升级警告
	BreakingChanges string      `gorm:"type:text" json:"breaking_changes"`      // 破坏性变更

	// 回滚信息
	RollbackSupported bool       `gorm:"default:true" json:"rollback_supported"`   // 是否支持回滚
	RollbackMinVersion string   `gorm:"type:varchar(50)" json:"rollback_min_version"` // 支持回滚的最小版本

	// 统计
	SuccessCount  int           `gorm:"default:0" json:"success_count"`          // 成功次数
	FailureCount  int           `gorm:"default:0" json:"failure_count"`         // 失败次数
	SuccessRate   float64       `gorm:"type:decimal(5,2);default:0" json:"success_rate"` // 成功率

	// 审核状态
	IsVerified    bool          `gorm:"default:false" json:"is_verified"`        // 是否已验证
	VerifiedAt    *time.Time   `json:"verified_at"`
	VerifiedBy    string        `gorm:"type:varchar(64)" json:"verified_by"`

	IsActive      bool          `gorm:"default:true" json:"is_active"`          // 是否启用

	CreatedBy     string        `gorm:"type:varchar(64)" json:"created_by"`
	UpdatedAt     time.Time     `json:"updated_at"`
	CreatedAt     time.Time     `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (OTACompatibilityMatrix) TableName() string {
	return "ota_compatibility_matrix"
}

// OTACompatibilityTest 兼容性测试记录
type OTACompatibilityTest struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	MatrixID     string         `gorm:"type:varchar(64);index" json:"matrix_id"`
	DeviceID     string         `gorm:"type:varchar(64);index" json:"device_id"`

	FromFirmware string         `gorm:"type:varchar(50)" json:"from_firmware"`
	ToFirmware   string         `gorm:"type:varchar(50)" json:"to_firmware"`

	TestResult   string         `gorm:"type:varchar(20)" json:"result"` // passed/failed/skipped
	ErrorCode    string         `gorm:"type:varchar(50)" json:"error_code"`
	ErrorMessage string         `gorm:"type:text" json:"error_message"`

	Duration     int            `gorm:"default:0" json:"duration"`        // 测试耗时(秒)

	Environment  string         `gorm:"type:text" json:"environment"`   // 测试环境信息

	TesterID     string         `gorm:"type:varchar(64)" json:"tester_id"`
	TestType     string         `gorm:"type:varchar(20)" json:"test_type"` // auto/manual
	CreatedAt    time.Time      `json:"created_at"`
}

// TableName 表名
func (OTACompatibilityTest) TableName() string {
	return "ota_compatibility_tests"
}
