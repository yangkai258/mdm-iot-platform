package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB is an alias for JSON for clarity in simulation context
type JSONB = JSON

// StringSlice is an alias for StringArray for clarity in simulation context
type StringSlice = StringArray

// IntSlice is a custom type for int arrays
type IntSlice []int64

// Value implements driver.Valuer for IntSlice
func (i IntSlice) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}
	return json.Marshal(i)
}

// Scan implements sql.Scanner for IntSlice
func (i *IntSlice) Scan(value interface{}) error {
	if value == nil {
		*i = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, i)
}

// ============================================================
// 仿真场景 (SimulationScene)
// ============================================================

// SimulationScene 仿真场景表
type SimulationScene struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	SceneName     string     `gorm:"type:varchar(255);not null" json:"scene_name"`
	SceneType     string     `gorm:"type:varchar(50);not null" json:"scene_type"` // preset/custom
	Environment   JSONB      `gorm:"type:jsonb" json:"environment"`
	Objects       StringSlice `gorm:"type:text[]" json:"objects"`
	Events        JSONB      `gorm:"type:jsonb" json:"events"`
	Config        JSONB      `gorm:"type:jsonb" json:"config"`
	IsPublic      bool       `gorm:"default:false" json:"is_public"`
	Score         float64    `gorm:"type:decimal(5,2)" json:"score"`
	Downloads     int        `gorm:"default:0" json:"downloads"`
	Tags          StringSlice `gorm:"type:text[]" json:"tags"`
	Status        string     `gorm:"type:varchar(20);default:'idle'" json:"status"` // idle/running/completed
	CreatedBy     *uint      `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (SimulationScene) TableName() string { return "simulation_scenes" }

// ============================================================
// 仿真会话 (SimulationSession)
// ============================================================

// SimulationSession 仿真会话表
type SimulationSession struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SceneID      uint       `gorm:"not null;index" json:"scene_id"`
	PetID        *uint      `json:"pet_id"`
	SessionName  string     `gorm:"type:varchar(255)" json:"session_name"`
	Status       string     `gorm:"type:varchar(20);default:'created'" json:"status"` // created/running/paused/stopped/completed
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	DurationMs   int        `json:"duration_ms"`
	Parameters   JSONB      `gorm:"type:jsonb" json:"parameters"`
	Environment  JSONB      `gorm:"type:jsonb" json:"environment"`
	CurrentState JSONB      `gorm:"type:jsonb" json:"current_state"`
	ResultID     *uint      `json:"result_id"`
	CreatedBy    *uint      `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (SimulationSession) TableName() string { return "simulation_sessions" }

// ============================================================
// 仿真任务 (SimulationTask)
// ============================================================

// SimulationTask 仿真任务表
type SimulationTask struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	SessionID     uint       `gorm:"not null;index" json:"session_id"`
	TaskName      string     `gorm:"type:varchar(255);not null" json:"task_name"`
	TaskType      string     `gorm:"type:varchar(50);not null" json:"task_type"` // test_case/stress_test/ab_experiment/playback
	Status        string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/running/completed/failed/cancelled
	Progress      int        `gorm:"default:0" json:"progress"`
	Priority      int        `gorm:"default:0" json:"priority"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	DurationMs    int        `json:"duration_ms"`
	Config        JSONB      `gorm:"type:jsonb" json:"config"`
	ResultData    JSONB      `gorm:"type:jsonb" json:"result_data"`
	ErrorMessage  string     `gorm:"type:text" json:"error_message"`
	RetryCount    int        `gorm:"default:0" json:"retry_count"`
	CreatedBy     *uint      `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (SimulationTask) TableName() string { return "simulation_tasks" }

// ============================================================
// 仿真结果 (SimulationResult)
// ============================================================

// SimulationResult 仿真结果表
type SimulationResult struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	SessionID       uint       `gorm:"not null;index" json:"session_id"`
	TaskID          *uint      `json:"task_id"`
	ResultName      string     `gorm:"type:varchar(255)" json:"result_name"`
	ResultType      string     `gorm:"type:varchar(50)" json:"result_type"` // test_case/stress_test/ab_experiment/playback
	Status          string     `gorm:"type:varchar(20)" json:"status"` // passed/failed/skipped
	Summary         JSONB      `gorm:"type:jsonb" json:"summary"`
	Metrics         JSONB      `gorm:"type:jsonb" json:"metrics"`
	DurationMs      int        `json:"duration_ms"`
	TestStepsPassed int        `json:"test_steps_passed"`
	TestStepsFailed int        `json:"test_steps_failed"`
	TotalRequests   int        `json:"total_requests"`
	FailedRequests  int        `json:"failed_requests"`
	SuccessRate     float64    `gorm:"type:decimal(5,2)" json:"success_rate"`
	AvgResponseTime float64    `gorm:"type:decimal(10,2)" json:"avg_response_time"`
	P50ResponseTime float64    `gorm:"type:decimal(10,2)" json:"p50_response_time"`
	P95ResponseTime float64    `gorm:"type:decimal(10,2)" json:"p95_response_time"`
	P99ResponseTime float64    `gorm:"type:decimal(10,2)" json:"p99_response_time"`
	Logs            string     `gorm:"type:text" json:"logs"`
	Screenshots     StringSlice `gorm:"type:text[]" json:"screenshots"`
	Attachments     JSONB      `gorm:"type:jsonb" json:"attachments"`
	GeneratedAt     time.Time  `json:"generated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (SimulationResult) TableName() string { return "simulation_results" }

// ============================================================
// 仿真数据集 (SimulationDataset)
// ============================================================

// SimulationDataset 仿真数据集表
type SimulationDataset struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	DatasetName    string     `gorm:"type:varchar(255);not null" json:"dataset_name"`
	DatasetType    string     `gorm:"type:varchar(50);not null" json:"dataset_type"` // sensor/voice/image/behavior/environment
	Description    string     `gorm:"type:text" json:"description"`
	Schema         JSONB      `gorm:"type:jsonb" json:"schema"`
	CurrentVersion string     `gorm:"type:varchar(20)" json:"current_version"`
	RecordCount    int        `gorm:"default:0" json:"record_count"`
	FileSize       int64      `gorm:"default:0" json:"file_size"`
	Tags           StringSlice `gorm:"type:text[]" json:"tags"`
	IsPublic       bool       `gorm:"default:false" json:"is_public"`
	CreatedBy      *uint      `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (SimulationDataset) TableName() string { return "simulation_datasets" }

// DatasetVersion 数据集版本表
type DatasetVersion struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	DatasetID   uint       `gorm:"not null;index" json:"dataset_id"`
	Version     string     `gorm:"type:varchar(20);not null" json:"version"`
	Description string     `gorm:"type:text" json:"description"`
	RecordCount int        `gorm:"default:0" json:"record_count"`
	FileSize    int64      `gorm:"default:0" json:"file_size"`
	FilePath    string     `gorm:"type:varchar(500)" json:"file_path"`
	SchemaHash  string     `gorm:"type:varchar(64)" json:"schema_hash"`
	Stats       JSONB      `gorm:"type:jsonb" json:"stats"`
	IsPublished bool       `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedBy   *uint      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (DatasetVersion) TableName() string { return "simulation_dataset_versions" }

// ============================================================
// 压力测试配置 (StressTestConfig)
// ============================================================

// StressTestConfig 压力测试表
type StressTestConfig struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TestName        string     `gorm:"type:varchar(255);not null" json:"test_name"`
	TestType        string     `gorm:"type:varchar(20);not null" json:"test_type"` // concurrent/performance/stability
	Config          JSONB      `gorm:"type:jsonb;not null" json:"config"`
	Status          string     `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/running/paused/stopped/completed/failed
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	DurationSeconds int        `json:"duration_seconds"`
	Metrics         JSONB      `gorm:"type:jsonb" json:"metrics"`
	Summary         JSONB      `gorm:"type:jsonb" json:"summary"`
	ThresholdsPassed *bool     `json:"thresholds_passed"`
	ReportURL       string     `gorm:"type:varchar(500)" json:"report_url"`
	CreatedBy       *uint      `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (StressTestConfig) TableName() string { return "stress_tests" }

// ============================================================
// A/B实验 (ABExperiment)
// ============================================================

// ABExperiment A/B实验表
type ABExperiment struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ExperimentName string     `gorm:"type:varchar(255);not null" json:"experiment_name"`
	ExperimentKey  string     `gorm:"type:varchar(100);uniqueIndex" json:"experiment_key"`
	Description    string     `gorm:"type:text" json:"description"`
	Hypothesis     string     `gorm:"type:text" json:"hypothesis"`
	TrafficPercent int        `gorm:"default:100" json:"traffic_percent"` // 流量百分比 0-100
	Status         string     `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/running/paused/completed/archived
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	VariantAConfig JSONB      `gorm:"type:jsonb" json:"variant_a_config"` // 对照组配置
	VariantBConfig JSONB      `gorm:"type:jsonb" json:"variant_b_config"` // 实验组配置
	TargetMetrics  StringSlice `gorm:"type:text[]" json:"target_metrics"`
	ResultSummary  JSONB      `gorm:"type:jsonb" json:"result_summary"`
	Winner         string     `gorm:"type:varchar(10)" json:"winner"` // A/B/none
	Confidence     float64    `gorm:"type:decimal(5,2)" json:"confidence"`
	Tags           StringSlice `gorm:"type:text[]" json:"tags"`
	CreatedBy      *uint      `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (ABExperiment) TableName() string { return "ab_experiments" }

// ABExperimentResult A/B实验结果表
type ABExperimentResult struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ExperimentID  uint       `gorm:"not null;index" json:"experiment_id"`
	Variant       string     `gorm:"type:varchar(10);not null" json:"variant"` // A/B
	SampleSize    int        `gorm:"default:0" json:"sample_size"`
	Conversions   int        `gorm:"default:0" json:"conversions"`
	ConversionRate float64   `gorm:"type:decimal(10,4)" json:"conversion_rate"`
	Metrics       JSONB      `gorm:"type:jsonb" json:"metrics"`
	AvgValue      float64    `gorm:"type:decimal(10,4)" json:"avg_value"`
	Variance      float64    `gorm:"type:decimal(10,4)" json:"variance"`
	PValue        float64    `gorm:"type:decimal(10,6)" json:"p_value"`
	StatSignificance bool    `gorm:"default:false" json:"stat_significance"`
	RecordedAt    time.Time  `json:"recorded_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (ABExperimentResult) TableName() string { return "ab_experiment_results" }

// ============================================================
// 回放录制 (Recording)
// ============================================================

// Recording 回放录制表
type Recording struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	DeviceID    string     `gorm:"type:varchar(100);index" json:"device_id"`
	PetID       *uint      `json:"pet_id"`
	RecordType  string     `gorm:"type:varchar(20);not null" json:"record_type"` // auto/manual
	StartTime   time.Time  `gorm:"not null" json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	DurationMs  int        `json:"duration_ms"`
	SensorData  JSONB      `gorm:"type:jsonb" json:"sensor_data"`
	UserActions JSONB      `gorm:"type:jsonb" json:"user_actions"`
	Events      JSONB      `gorm:"type:jsonb" json:"events"`
	PlaybackURL string     `gorm:"type:varchar(500)" json:"playback_url"`
	FilePath    string     `gorm:"type:varchar(500)" json:"file_path"`
	FileSize    int64      `json:"file_size"`
	Status      string     `gorm:"type:varchar(20);default:'recording'" json:"status"` // recording/paused/completed
	Metadata    JSONB      `gorm:"type:jsonb" json:"metadata"`
	CreatedBy   *uint      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (Recording) TableName() string { return "simulation_recordings" }

// ============================================================
// 测试用例 (TestCase) - 自动化测试框架
// ============================================================

// TestCase 测试用例表
type TestCase struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	CaseName        string     `gorm:"type:varchar(255);not null" json:"case_name"`
	CaseType        string     `gorm:"type:varchar(30);not null" json:"case_type"` // functional/performance/stress/regression/smoke
	Module          string     `gorm:"type:varchar(50)" json:"module"`
	Priority        string     `gorm:"type:varchar(20);default:'medium'" json:"priority"` // high/medium/low
	Status          string     `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft/active/archived
	Description     string     `gorm:"type:text" json:"description"`
	Preconditions   string     `gorm:"type:text" json:"preconditions"`
	TestSteps       JSONB      `gorm:"type:jsonb" json:"test_steps"`
	ExpectedResult  string     `gorm:"type:text" json:"expected_result"`
	Tags            StringSlice `gorm:"type:text[]" json:"tags"`
	Dependencies    IntSlice    `gorm:"type:text[]" json:"dependencies"`
	Version         int        `gorm:"default:1" json:"version"`
	CreatedBy       *uint      `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (TestCase) TableName() string { return "simulation_testcases" }

// ============================================================
// 测试执行记录 (TestExecution)
// ============================================================

// TestExecution 测试执行记录表
type TestExecution struct {
	ID              string     `gorm:"type:varchar(50);primaryKey" json:"id"`
	TestCaseID      uint       `gorm:"not null;index" json:"testcase_id"`
	BatchID         string     `gorm:"type:varchar(50);index" json:"batch_id"`
	ExecutionType   string     `gorm:"type:varchar(20)" json:"execution_type"` // single/batch
	TriggerParams   JSONB      `gorm:"type:jsonb" json:"trigger_params"`
	Status          string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/running/passed/failed/skipped/cancelled
	Progress        int        `gorm:"default:0" json:"progress"`
	CurrentStep     int        `gorm:"default:0" json:"current_step"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	DurationMs      int        `json:"duration_ms"`
	Environment     JSONB      `gorm:"type:jsonb" json:"environment"`
	ResultDetails   JSONB      `gorm:"type:jsonb" json:"result_details"`
	Screenshots     StringSlice `gorm:"type:text[]" json:"screenshots"`
	Logs            string     `gorm:"type:text" json:"logs"`
	ErrorMessage    string     `gorm:"type:text" json:"error_message"`
	RetryCount      int        `gorm:"default:0" json:"retry_count"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (TestExecution) TableName() string { return "test_executions" }

// ============================================================
// 测试报告 (TestReport)
// ============================================================

// TestReport 测试报告表
type TestReport struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	ReportName   string     `gorm:"type:varchar(255)" json:"report_name"`
	ReportType   string     `gorm:"type:varchar(20);default:'execution'" json:"report_type"`
	ExecutionIDs StringSlice `gorm:"type:text[]" json:"execution_ids"`
	BatchID      string     `gorm:"type:varchar(50)" json:"batch_id"`
	Summary      JSONB      `gorm:"type:jsonb" json:"summary"`
	PassCount    int        `gorm:"default:0" json:"pass_count"`
	FailCount    int        `gorm:"default:0" json:"fail_count"`
	SkipCount    int        `gorm:"default:0" json:"skip_count"`
	TotalCount   int        `gorm:"default:0" json:"total_count"`
	PassRate     float64    `gorm:"type:decimal(5,2)" json:"pass_rate"`
	AvgDurationMs int       `json:"avg_duration_ms"`
	Coverage     JSONB      `gorm:"type:jsonb" json:"coverage"`
	TrendData    JSONB      `gorm:"type:jsonb" json:"trend_data"`
	FailedCases  JSONB      `gorm:"type:jsonb" json:"failed_cases"`
	CreatedBy    *uint      `json:"created_by"`
	GeneratedAt  time.Time  `json:"generated_at"`
}

func (TestReport) TableName() string { return "test_reports" }
