package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===================== Fairness Test Models =====================

// FairnessTestStatus 公平性测试状态
type FairnessTestStatus string

const (
	FairnessTestStatusPending   FairnessTestStatus = "pending"
	FairnessTestStatusRunning   FairnessTestStatus = "running"
	FairnessTestStatusCompleted FairnessTestStatus = "completed"
	FairnessTestStatusFailed    FairnessTestStatus = "failed"
)

// FairnessTestType 公平性测试类型
type FairnessTestType string

const (
	FairnessTestTypeDemographicParity FairnessTestType = "demographic_parity"
	FairnessTestTypeEqualOpportunity  FairnessTestType = "equal_opportunity"
	FairnessTestTypeDisparateImpact   FairnessTestType = "disparate_impact"
	FairnessTestTypeBiasDetection     FairnessTestType = "bias_detection"
	FairnessTestTypeRepresentation    FairnessTestType = "representation"
)

// FairnessTest 公平性测试
type FairnessTest struct {
	ID           uint               `json:"id" gorm:"primaryKey"`
	TestKey      string             `json:"test_key" gorm:"uniqueIndex;size:64;not null"`
	Name         string             `json:"name" gorm:"size:128;not null"`
	Description  string             `json:"description" gorm:"type:text"`
	TestType     FairnessTestType   `json:"test_type" gorm:"size:32;not null;index"`
	ModelID      uint               `json:"model_id" gorm:"index"`
	ModelKey     string             `json:"model_key" gorm:"size:128"`
	Status       FairnessTestStatus `json:"status" gorm:"size:20;default:'pending';index"`
	Progress     int                `json:"progress" gorm:"type:smallint;default:0"`
	Config       string             `json:"config" gorm:"type:jsonb"`         // 测试配置：敏感属性定义、阈值、样本数等
	TestData     string             `json:"test_data" gorm:"type:jsonb"`       // 测试数据集
	Results      string             `json:"results" gorm:"type:jsonb"`         // 测试结果
	Metrics      string             `json:"metrics" gorm:"type:jsonb"`         // 公平性指标
	ReportPath   string             `json:"report_path" gorm:"size:512"`       // 报告文件路径
	ErrorMessage string             `json:"error_message" gorm:"type:text"`
	StartedAt    *time.Time         `json:"started_at"`
	CompletedAt  *time.Time         `json:"completed_at"`
	OrgID        uint               `gorm:"index" json:"org_id"`
	CreateUserID uint               `gorm:"index" json:"create_user_id"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	DeletedAt    gorm.DeletedAt     `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (f *FairnessTest) BeforeCreate(tx *gorm.DB) error {
	if f.TestKey == "" {
		f.TestKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (FairnessTest) TableName() string {
	return "ai_fairness_tests"
}

// FairnessTestCreate 创建公平性测试请求
type FairnessTestCreate struct {
	Name        string             `json:"name" binding:"required" gorm:"size:128"`
	Description string             `json:"description" gorm:"type:text"`
	TestType    FairnessTestType   `json:"test_type" binding:"required" gorm:"size:32"`
	ModelID     uint               `json:"model_id" binding:"required"`
	ModelKey    string              `json:"model_key" gorm:"size:128"`
	Config      string             `json:"config" gorm:"type:jsonb"`
	TestData    string             `json:"test_data" gorm:"type:jsonb"`
}

// FairnessTestUpdate 更新公平性测试请求
type FairnessTestUpdate struct {
	Name        string `json:"name" gorm:"size:128"`
	Description string `json:"description" gorm:"type:text"`
	Config      string `json:"config" gorm:"type:jsonb"`
	TestData    string `json:"test_data" gorm:"type:jsonb"`
}

// FairnessTestRun 运行测试请求
type FairnessTestRun struct {
	TestData string `json:"test_data" gorm:"type:jsonb"` // 可选：覆盖测试数据集
}

// FairnessTestReport 公平性测试报告
type FairnessTestReport struct {
	TestKey       string    `json:"test_key"`
	TestName      string    `json:"test_name"`
	TestType      string    `json:"test_type"`
	Status        string    `json:"status"`
	OverallScore  float64   `json:"overall_score"`  // 0-100 总体公平性评分
	PassStatus    string    `json:"pass_status"`    // pass, fail, warning
	Metrics       string    `json:"metrics"`         // 详细指标 JSON
	Recommendations []string `json:"recommendations"` // 改进建议
	GeneratedAt   time.Time `json:"generated_at"`
}

// ===================== Bias Detection Models =====================

// BiasType 偏见类型
type BiasType string

const (
	BiasTypeGender        BiasType = "gender"
	BiasTypeAge          BiasType = "age"
	BiasTypeRace         BiasType = "race"
	BiasTypeEthnicity    BiasType = "ethnicity"
	BiasTypeReligion      BiasType = "religion"
	BiasTypeSexualOrient  BiasType = "sexual_orientation"
	BiasTypeDisability   BiasType = "disability"
	BiasTypeSocioeconomic BiasType = "socioeconomic"
	BiasTypeGeographic    BiasType = "geographic"
	BiasTypeLanguage      BiasType = "language"
)

// BiasSeverity 偏见严重程度
type BiasSeverity string

const (
	BiasSeverityLow      BiasSeverity = "low"
	BiasSeverityMedium   BiasSeverity = "medium"
	BiasSeverityHigh     BiasSeverity = "high"
	BiasSeverityCritical BiasSeverity = "critical"
)

// BiasDetection 偏见检测记录
type BiasDetection struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	DetectionKey    string         `json:"detection_key" gorm:"uniqueIndex;size:64;not null"`
	ModelID         uint           `json:"model_id" gorm:"index"`
	ModelKey        string         `json:"model_key" gorm:"size:128"`
	InputData       string         `json:"input_data" gorm:"type:text"`            // 输入数据
	OutputData      string         `json:"output_data" gorm:"type:text"`           // 输出数据
	BiasType        BiasType       `json:"bias_type" gorm:"size:32;not null;index"`
	Severity        BiasSeverity  `json:"severity" gorm:"size:20;default:'low'"`
	Confidence      float64       `json:"confidence" gorm:"type:decimal(5,4)"`    // 置信度 0-1
	BiasScore       float64       `json:"bias_score" gorm:"type:decimal(5,4)"`   // 偏见得分 0-1
	Evidence        string        `json:"evidence" gorm:"type:text"`              // 偏见证据
	Context         string        `json:"context" gorm:"type:text"`               // 上下文信息
	Recommendation  string        `json:"recommendation" gorm:"type:text"`       // 建议
	DetectedAt      time.Time     `json:"detected_at"`
	OrgID           uint          `gorm:"index" json:"org_id"`
	CreateUserID    uint          `gorm:"index" json:"create_user_id"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (b *BiasDetection) BeforeCreate(tx *gorm.DB) error {
	if b.DetectionKey == "" {
		b.DetectionKey = uuid.New().String()
	}
	if b.DetectedAt.IsZero() {
		b.DetectedAt = time.Now()
	}
	return nil
}

// TableName 表名
func (BiasDetection) TableName() string {
	return "ai_bias_detections"
}

// BiasDetectRequest 偏见检测请求
type BiasDetectRequest struct {
	ModelID    uint      `json:"model_id" binding:"required"`
	ModelKey   string    `json:"model_key" gorm:"size:128"`
	InputData  string    `json:"input_data" binding:"required" gorm:"type:text"`
	OutputData string    `json:"output_data" gorm:"type:text"`
	BiasTypes  []BiasType `json:"bias_types"` // 可选：指定要检测的偏见类型
}

// FairnessMetrics 公平性指标
type FairnessMetrics struct {
	ModelID              uint    `json:"model_id" gorm:"primaryKey"`
	ModelKey             string  `json:"model_key" gorm:"size:128"`
	DemographicParity    float64 `json:"demographic_parity" gorm:"type:decimal(5,4)"`    // 人口统计均等
	EqualOpportunity     float64 `json:"equal_opportunity" gorm:"type:decimal(5,4)"`     // 机会均等
	DisparateImpact       float64 `json:"disparate_impact" gorm:"type:decimal(5,4)"`      // 不同影响
	StatisticalParity    float64 `json:"statistical_parity" gorm:"type:decimal(5,4)"`   // 统计均等
	PrecisionGap         float64 `json:"precision_gap" gorm:"type:decimal(5,4)"`         // 精确度差距
	RecallGap            float64 `json:"recall_gap" gorm:"type:decimal(5,4)"`           // 召回率差距
	FalsePositiveGap     float64 `json:"false_positive_gap" gorm:"type:decimal(5,4)"`   // 假阳性差距
	OverallScore         float64 `json:"overall_score" gorm:"type:decimal(5,2)"`         // 总体公平性评分 0-100
	TotalTestsRun        int     `json:"total_tests_run" gorm:"type:integer"`            // 总测试次数
	TotalBiasDetections  int     `json:"total_bias_detections" gorm:"type:integer"`      // 总偏见检测次数
	CriticalBiasCount    int     `json:"critical_bias_count" gorm:"type:smallint"`       // 严重偏见数
	HighBiasCount        int     `json:"high_bias_count" gorm:"type:smallint"`           // 高度偏见数
	LastUpdated          time.Time `json:"last_updated"`
}

// TableName 表名
func (FairnessMetrics) TableName() string {
	return "ai_fairness_metrics"
}

// ===================== Model Audit Models =====================

// AuditLog AI模型操作审计日志
type AIAuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Action        string    `gorm:"type:varchar(50);not null;index" json:"action"`         // 操作类型: inference, training, deploy, fairness_test, bias_detect, model_create, model_update, model_delete
	Module        string    `gorm:"type:varchar(50);not null;index" json:"module"`         // 模块: ai_inference, ai_training, ai_fairness, ai_bias, ai_model
	ResourceType  string    `gorm:"type:varchar(50)" json:"resource_type"`               // 资源类型: model, inference, training, fairness_test, bias_detection
	ResourceID    string    `gorm:"type:varchar(100);index" json:"resource_id"`           // 资源ID
	ModelID       uint      `gorm:"index" json:"model_id"`                              // 关联模型ID
	ModelKey      string    `gorm:"size:128" json:"model_key"`                           // 模型key
	UserID        uint      `gorm:"index" json:"user_id"`                               // 操作人ID
	Username      string    `gorm:"type:varchar(100)" json:"username"`                   // 操作人用户名
	IP            string    `gorm:"type:varchar(45)" json:"ip"`                          // IP地址
	UserAgent     string    `gorm:"type:varchar(500)" json:"user_agent"`                 // User-Agent
	Status        int       `gorm:"default:1" json:"status"`                            // 1:成功 2:失败
	ErrorMsg      string    `gorm:"type:text" json:"error_msg"`                         // 错误信息
	RequestMethod string    `gorm:"type:varchar(10)" json:"request_method"`             // HTTP方法
	RequestPath   string    `gorm:"type:varchar(500)" json:"request_path"`              // 请求路径
	RequestBody   string    `gorm:"type:text" json:"request_body"`                      // 请求体（脱敏后）
	ResponseCode  int       `json:"response_code"`                                       // 响应码
	Duration      int       `json:"duration"`                                           // 操作耗时（毫秒）
	Metadata      string    `gorm:"type:json" json:"metadata"`                          // 额外元数据（JSON）
	FairnessScore float64   `gorm:"type:decimal(5,2)" json:"fairness_score"`            // 公平性评分（相关操作时）
	BiasDetected  bool      `gorm:"default:false" json:"bias_detected"`                 // 是否检测到偏见
	TenantID      uint      `gorm:"index" json:"tenant_id"`                            // 租户ID
	OrgID         uint      `gorm:"index" json:"org_id"`                               // 组织ID
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 表名
func (AIAuditLog) TableName() string {
	return "ai_audit_logs"
}

// AIAuditReport AI审计报告
type AIAuditReport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ReportKey      string    `gorm:"uniqueIndex;size:64;not null" json:"report_key"`
	ReportType     string    `gorm:"type:varchar(50);not null" json:"report_type"` // fairness, bias, comprehensive, model_performance
	ModelID        uint      `gorm:"index" json:"model_id"`
	ModelKey       string    `gorm:"size:128" json:"model_key"`
	PeriodStart    time.Time `json:"period_start" gorm:"type:date"`
	PeriodEnd      time.Time `json:"period_end" gorm:"type:date"`
	Summary        string    `gorm:"type:text" json:"summary"`                      // 报告摘要
	Findings       string    `gorm:"type:json" json:"findings"`                      // 主要发现 JSON数组
	Metrics        string    `gorm:"type:json" json:"metrics"`                      // 详细指标 JSON
	Recommendations string   `gorm:"type:text" json:"recommendations"`              // 建议
	RiskLevel      string    `gorm:"type:varchar:20" json:"risk_level"`             // low, medium, high, critical
	GeneratedBy    uint      `gorm:"index" json:"generated_by"`                     // 生成人
	GeneratedAt    time.Time `json:"generated_at"`
	TenantID       uint      `gorm:"index" json:"tenant_id"`
	OrgID          uint      `gorm:"index" json:"org_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// BeforeCreate 创建前自动生成 UUID
func (a *AIAuditReport) BeforeCreate(tx *gorm.DB) error {
	if a.ReportKey == "" {
		a.ReportKey = uuid.New().String()
	}
	return nil
}

// TableName 表名
func (AIAuditReport) TableName() string {
	return "ai_audit_reports"
}

// AIAuditReportRequest 生成审计报告请求
type AIAuditReportRequest struct {
	ReportType  string `json:"report_type" binding:"required" gorm:"size:50"` // fairness, bias, comprehensive
	ModelID     uint   `json:"model_id"`
	ModelKey    string `json:"model_key" gorm:"size:128"`
	PeriodStart string `json:"period_start"` // YYYY-MM-DD
	PeriodEnd   string `json:"period_end"`   // YYYY-MM-DD
}

// AIAuditLogQuery 审计日志查询参数
type AIAuditLogQuery struct {
	ModelID    uint   `form:"model_id"`
	ModelKey   string `form:"model_key"`
	Action     string `form:"action"`
	Module     string `form:"module"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Status     int    `form:"status"`
}
