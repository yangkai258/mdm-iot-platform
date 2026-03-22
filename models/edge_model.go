package models

import (
	"time"

	"gorm.io/gorm"
)

// EdgeModel 端侧推理模型
type EdgeModel struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(128);not null" json:"name"`
	Version         string         `gorm:"type:varchar(32);not null" json:"version"`
	Description     string         `gorm:"type:text" json:"description"`
	ModelType       string         `gorm:"type:varchar(32);not null" json:"model_type"` // yolo, mobilenet, bert, whisper, tflite, onnx
	ModelFormat     string         `gorm:"type:varchar(16);not null" json:"model_format"` // tflite, onnx, pt, h5
	FileSize        int64          `gorm:"default:0" json:"file_size"`
	FileURL         string         `gorm:"type:varchar(512)" json:"file_url"`
	FileMD5         string         `gorm:"type:varchar(32)" json:"file_md5"`
	InputShape      string         `gorm:"type:varchar(128)" json:"input_shape"`      // e.g., "1,3,224,224"
	OutputShape     string         `gorm:"type:varchar(128)" json:"output_shape"`    // e.g., "1,1000"
	Input_dtype     string         `gorm:"type:varchar(16);default:'float32'" json:"input_dtype"`
	Output_dtype    string         `gorm:"type:varchar(16);default:'float32'" json:"output_dtype"`
	Framework       string         `gorm:"type:varchar(32)" json:"framework"`        // TensorFlow, PyTorch, ONNXRuntime
	HardwareTarget  string         `gorm:"type:varchar(32)" json:"hardware_target"`  // cortex_m, esp32, nrf52, generic
	MaxBatchSize    int            `gorm:"default:1" json:"max_batch_size"`
	AvgLatencyMs    float64        `gorm:"default:0" json:"avg_latency_ms"`   // 平均推理延迟(ms)
	Accuracy        float64        `gorm:"default:0" json:"accuracy"`         // 精度指标
	MemoryUsageKB   int            `gorm:"default:0" json:"memory_usage_kb"`  // 内存占用(KB)
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	OrgID           uint           `gorm:"index" json:"org_id"`
	CreateUserID    uint           `gorm:"index" json:"create_user_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (EdgeModel) TableName() string {
	return "edge_models"
}

// EdgeModelDeployment 端侧模型部署记录
type EdgeModelDeployment struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	ModelID        uint             `gorm:"not null;index" json:"model_id"`
	DeviceID       string           `gorm:"type:varchar(64);not null;index" json:"device_id"`
	Status         string           `gorm:"type:varchar(16);default:'pending'" json:"status"` // pending, deploying, running, stopped, failed
	DeployedAt     *time.Time       `json:"deployed_at"`
	StoppedAt      *time.Time       `json:"stopped_at"`
	ErrorMessage   string           `gorm:"type:text" json:"error_message"`
	RuntimeVersion string           `gorm:"type:varchar(32)" json:"runtime_version"` // 运行时版本
	OrgID          uint             `gorm:"index" json:"org_id"`
	CreateUserID   uint             `gorm:"index" json:"create_user_id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
	// GORM 关联
	EdgeModel      EdgeModel        `gorm:"foreignKey:ModelID" json:"edge_model,omitempty"`
}

// TableName 指定表名
func (EdgeModelDeployment) TableName() string {
	return "edge_model_deployments"
}

// EdgeInferenceLog 端侧推理日志
type EdgeInferenceLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	DeploymentID uint          `gorm:"index" json:"deployment_id"`
	DeviceID    string         `gorm:"type:varchar(64);not null;index" json:"device_id"`
	ModelID     uint           `gorm:"not null" json:"model_id"`
	InputSize   int64          `gorm:"default:0" json:"input_size"`   // 输入数据大小(bytes)
	OutputSize  int64          `gorm:"default:0" json:"output_size"`  // 输出数据大小(bytes)
	LatencyMs   float64        `gorm:"default:0" json:"latency_ms"`   // 推理耗时(ms)
	ResultCode  string         `gorm:"type:varchar(16)" json:"result_code"` // success, error, timeout
	ErrorMsg    string         `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 指定表名
func (EdgeInferenceLog) TableName() string {
	return "edge_inference_logs"
}
