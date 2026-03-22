package models

import (
	"time"

	"gorm.io/gorm"
)

// AnalyticsRecord 高级分析记录（存储计算好的分析数据）
type AnalyticsRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	AnalysisType string        `gorm:"type:varchar(50);not null" json:"analysis_type"` // device_health, member_activity, ota_performance, usage_pattern
	PeriodStart time.Time     `gorm:"not null" json:"period_start"`
	PeriodEnd   time.Time     `gorm:"not null" json:"period_end"`
	Dimensions  string        `gorm:"type:jsonb" json:"dimensions"`  // JSON: 多维度聚合结果
	Metrics     string        `gorm:"type:jsonb" json:"metrics"`     // JSON: 核心指标
	Summary     string        `gorm:"type:text" json:"summary"`       // 分析摘要文字
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AnalyticsRecord) TableName() string {
	return "analytics_records"
}

// ExportJob 数据导出任务
type ExportJob struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	ExportType   string         `gorm:"type:varchar(30);not null" json:"export_type"` // csv, excel, pdf
	DataSource   string         `gorm:"type:varchar(50);not null" json:"data_source"` // devices, members, alerts, ota, usage, custom
	Filters      string         `gorm:"type:jsonb" json:"filters"`                    // JSON: 筛选条件
	Status       string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, processing, completed, failed
	FilePath     string         `gorm:"type:varchar(512)" json:"file_path"`
	FileSize     int64          `gorm:"default:0" json:"file_size"`
	RecordCount  int64          `gorm:"default:0" json:"record_count"`
	ErrorMsg     string         `gorm:"type:text" json:"error_msg"`
	StartedAt    *time.Time    `json:"started_at"`
	CompletedAt  *time.Time    `json:"completed_at"`
	ExpiresAt    *time.Time    `json:"expires_at"` // 文件过期时间
	CreatedBy    uint           `gorm:"index" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (ExportJob) TableName() string {
	return "export_jobs"
}

// CustomReport 自定义报表定义
type CustomReport struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    string         `gorm:"type:uuid;index" json:"tenant_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(500)" json:"description"`
	ReportConfig string        `gorm:"type:jsonb;not null" json:"report_config"` // JSON: 报表配置（数据源、维度、指标、筛选器、排序）
	ChartType   string         `gorm:"type:varchar(30)" json:"chart_type"`       // bar, line, pie, table, area
	IsScheduled bool           `gorm:"default:false" json:"is_scheduled"`        // 是否定时执行
	CronExpr    string         `gorm:"type:varchar(50)" json:"cron_expr"`         // Cron表达式
	LastRunAt   *time.Time    `json:"last_run_at"`
	NextRunAt   *time.Time    `json:"next_run_at"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`            // 是否公开（其他用户可见）
	CreatedBy   uint           `gorm:"index" json:"created_by"`
	UpdatedBy   uint           `json:"updated_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CustomReport) TableName() string {
	return "custom_reports"
}

// ==================== 请求/响应结构体 ====================

// AdvancedAnalyticsRequest 高级分析请求
type AdvancedAnalyticsRequest struct {
	AnalysisType string `form:"analysis_type" json:"analysis_type"` // device_health, member_activity, ota_performance, usage_pattern
	PeriodType   string `form:"period_type" json:"period_type"`     // daily, weekly, monthly
	StartDate    string `form:"start_date" json:"start_date"`       // YYYY-MM-DD
	EndDate      string `form:"end_date" json:"end_date"`           // YYYY-MM-DD
	DeviceID     string `form:"device_id" json:"device_id"`         // 可选，设备维度
	GroupBy      string `form:"group_by" json:"group_by"`           // day, week, month, device_type, region
}

// TrendRequest 趋势分析请求
type TrendRequest struct {
	Metrics    string `form:"metrics" json:"metrics"`       // comma-separated: device_count, member_count, alert_count, ota_success_rate
	PeriodType string `form:"period_type" json:"period_type"` // daily, weekly, monthly
	StartDate  string `form:"start_date" json:"start_date"`  // YYYY-MM-DD
	EndDate    string `form:"end_date" json:"end_date"`      // YYYY-MM-DD
}

// PredictionRequest 预测分析请求
type PredictionRequest struct {
	Metric     string `form:"metric" json:"metric"`        // device_count, member_count, alert_count
	Method     string `form:"method" json:"method"`         // linear_regression, moving_average, arima
	Periods    int    `form:"periods" json:"periods"`      // 预测周期数
	StartDate  string `form:"start_date" json:"start_date"` // YYYY-MM-DD 训练数据起始
	EndDate    string `form:"end_date" json:"end_date"`     // YYYY-MM-DD 训练数据截止
}

// ExportRequest 数据导出请求
type ExportRequest struct {
	DataSource string            `json:"data_source" binding:"required"` // devices, members, alerts, ota, usage, custom
	ExportType string            `json:"export_type" binding:"required"`   // csv, excel, pdf
	Filters    map[string]string `json:"filters"`                         // 筛选条件
	Columns    []string          `json:"columns"`                        // 导出的列
	StartDate  string            `json:"start_date"`                     // YYYY-MM-DD
	EndDate    string            `json:"end_date"`                       // YYYY-MM-DD
}

// CreateCustomReportRequest 创建自定义报表请求
type CreateCustomReportRequest struct {
	Name        string            `json:"name" binding:"required,min=1,max=100"`
	Description string            `json:"description" binding:"max=500"`
	ReportConfig map[string]interface{} `json:"report_config" binding:"required"`
	ChartType   string            `json:"chart_type"` // bar, line, pie, table, area
	IsScheduled bool              `json:"is_scheduled"`
	CronExpr    string            `json:"cron_expr"`
	IsPublic    bool              `json:"is_public"`
}

// UpdateCustomReportRequest 更新自定义报表请求
type UpdateCustomReportRequest struct {
	Name        string            `json:"name" binding:"omitempty,min=1,max=100"`
	Description string            `json:"description" binding:"max=500"`
	ReportConfig map[string]interface{} `json:"report_config"`
	ChartType   string            `json:"chart_type"`
	IsScheduled *bool             `json:"is_scheduled"`
	CronExpr    string            `json:"cron_expr"`
	IsPublic    *bool             `json:"is_public"`
}

// ==================== 响应结构体 ====================

// AdvancedAnalyticsResponse 高级分析响应
type AdvancedAnalyticsResponse struct {
	AnalysisType string                 `json:"analysis_type"`
	PeriodStart  string                 `json:"period_start"`
	PeriodEnd    string                 `json:"period_end"`
	Dimensions   map[string]interface{} `json:"dimensions"`
	Metrics      map[string]interface{} `json:"metrics"`
	Summary      string                 `json:"summary"`
}

// TrendDataPoint 趋势数据点
type TrendDataPoint struct {
	Date    string  `json:"date"`
	Metrics map[string]float64 `json:"metrics"`
}

// TrendsResponse 趋势分析响应
type TrendsResponse struct {
	Metrics    []string          `json:"metrics"`
	PeriodType string            `json:"period_type"`
	DataPoints []TrendDataPoint  `json:"data_points"`
}

// PredictionDataPoint 预测数据点
type PredictionDataPoint struct {
	Date   string  `json:"date"`
	Value  float64 `json:"value"`
	Lower  float64 `json:"lower"`  // 置信区间下限
	Upper  float64 `json:"upper"`  // 置信区间上限
}

// PredictionsResponse 预测分析响应
type PredictionsResponse struct {
	Metric     string               `json:"metric"`
	Method     string               `json:"method"`
	Historical []PredictionDataPoint `json:"historical"` // 历史数据
	Forecast   []PredictionDataPoint `json:"forecast"`   // 预测数据
	Accuracy   float64              `json:"accuracy"`   // 模型准确度（R2分数）
}

// ExportJobResponse 导出任务响应
type ExportJobResponse struct {
	ID          uint      `json:"id"`
	Status      string    `json:"status"`
	DataSource  string    `json:"data_source"`
	ExportType  string    `json:"export_type"`
	RecordCount int64     `json:"record_count"`
	FileSize    int64     `json:"file_size"`
	DownloadURL string    `json:"download_url,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
