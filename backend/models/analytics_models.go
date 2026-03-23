package models

import (
	"time"

	"gorm.io/gorm"
)

// AnalyticsEvent 用户行为分析事件
type AnalyticsEvent struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	EventType    string         `gorm:"type:varchar(50);not null;index" json:"event_type"` // click, view, action, error, custom
	EventName    string         `gorm:"type:varchar(100);not null;index" json:"event_name"`
	DeviceID     string         `gorm:"type:varchar(36);index" json:"device_id"`
	UserID       string         `gorm:"type:varchar(36);index" json:"user_id"`
	SessionID    string         `gorm:"type:varchar(64);index" json:"session_id"`
	PageURL      string         `gorm:"type:varchar(512)" json:"page_url"`
	PageName     string         `gorm:"type:varchar(100)" json:"page_name"`
	ElementID    string         `gorm:"type:varchar(100)" json:"element_id"`
	ElementName  string         `gorm:"type:varchar(200)" json:"element_name"`
	EventData    string         `gorm:"type:jsonb" json:"event_data"` // 扩展属性 JSON
	Duration     int            `gorm:"type:int;default:0" json:"duration"` // 事件持续时长(ms)
	ScreenWidth  int            `gorm:"type:int" json:"screen_width"`
	ScreenHeight int            `gorm:"type:int" json:"screen_height"`
	Browser      string         `gorm:"type:varchar(50)" json:"browser"`
	OS           string         `gorm:"type:varchar(50)" json:"os"`
	AppVersion   string         `gorm:"type:varchar(20)" json:"app_version"`
	IPAddress    string         `gorm:"type:varchar(45)" json:"ip_address"`
	Country      string         `gorm:"type:varchar(50)" json:"country"`
	Province     string         `gorm:"type:varchar(50)" json:"province"`
	City         string         `gorm:"type:varchar(50)" json:"city"`
	Latitude     float64        `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude    float64        `gorm:"type:decimal(10,7)" json:"longitude"`
	OccurredAt   time.Time      `gorm:"index" json:"occurred_at"` // 事件发生时间
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AnalyticsEvent) TableName() string {
	return "analytics_events"
}

// AnalyticsReport 分析报告
type AnalyticsReport struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	ReportName   string         `gorm:"type:varchar(200);not null" json:"report_name"`
	ReportType   string         `gorm:"type:varchar(50);not null;index" json:"report_type"` // funnel, cohort, retention, custom, dashboard
	Description  string         `gorm:"type:text" json:"description"`
	Config       string         `gorm:"type:jsonb" json:"config"` // 报告配置参数 JSON
	ResultData   string         `gorm:"type:jsonb" json:"result_data"` // 报告结果数据 JSON
	FilePath     string         `gorm:"type:varchar(512)" json:"file_path"`
	Status       string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, generating, completed, failed
	StartDate    string         `gorm:"type:varchar(10)" json:"start_date"` // YYYY-MM-DD
	EndDate      string         `gorm:"type:varchar(10)" json:"end_date"` // YYYY-MM-DD
	GeneratedBy  string         `gorm:"type:varchar(36)" json:"generated_by"`
	GeneratedAt  *time.Time     `json:"generated_at"`
	Scheduled    bool           `gorm:"default:false" json:"scheduled"` // 是否定期生成
	ScheduleCron string         `gorm:"type:varchar(50)" json:"schedule_cron"` // Cron表达式
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (AnalyticsReport) TableName() string {
	return "analytics_reports"
}

// FunnelAnalysis 漏斗分析
type FunnelAnalysis struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TenantID     string         `gorm:"type:uuid;index" json:"tenant_id"`
	Name         string         `gorm:"type:varchar(200);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Steps        string         `gorm:"type:jsonb;not null" json:"steps"` // 漏斗步骤 JSON: [{"step":1,"name":"步骤名称","event_name":"事件名","event_type":"event|page"}]
	ConversionWindow int        `gorm:"type:int;default:86400" json:"conversion_window"` // 转化时间窗口(秒)，默认1天
	StartDate    string         `gorm:"type:varchar(10)" json:"start_date"` // YYYY-MM-DD
	EndDate      string         `gorm:"type:varchar(10)" json:"end_date"` // YYYY-MM-DD
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, active, archived
	CreatedBy    string         `gorm:"type:varchar(36)" json:"created_by"`
	ResultData   string         `gorm:"type:jsonb" json:"result_data"` // 漏斗计算结果 JSON
	LastRunAt    *time.Time     `json:"last_run_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (FunnelAnalysis) TableName() string {
	return "funnel_analyses"
}

// CohortAnalysis 群组分析
type CohortAnalysis struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TenantID       string         `gorm:"type:uuid;index" json:"tenant_id"`
	Name           string         `gorm:"type:varchar(200);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	CohortType     string         `gorm:"type:varchar(50);not null" json:"cohort_type"` // daily, weekly, monthly
	EntryEvent     string         `gorm:"type:varchar(100);not null" json:"entry_event"` // 入口事件
	RetentionEvent string         `gorm:"type:varchar(100);not null" json:"retention_event"` // 留存事件
	Periods        int            `gorm:"type:int;default:12" json:"periods"` // 分析周期数
	StartDate      string         `gorm:"type:varchar(10)" json:"start_date"` // YYYY-MM-DD
	EndDate        string         `gorm:"type:varchar(10)" json:"end_date"` // YYYY-MM-DD
	SegmentFilter  string         `gorm:"type:jsonb" json:"segment_filter"` // 分群过滤条件 JSON
	ResultData     string         `gorm:"type:jsonb" json:"result_data"` // 群组分析结果 JSON
	Status         string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, active, archived
	CreatedBy      string         `gorm:"type:varchar(36)" json:"created_by"`
	LastRunAt      *time.Time     `json:"last_run_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CohortAnalysis) TableName() string {
	return "cohort_analyses"
}

// RetentionReport 留存报告
type RetentionReport struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	TenantID      string         `gorm:"type:uuid;index" json:"tenant_id"`
	Name          string         `gorm:"type:varchar(200);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	AnalysisType  string         `gorm:"type:varchar(50);not null" json:"analysis_type"` // user, device, event
	TargetEvent   string         `gorm:"type:varchar(100);not null" json:"target_event"` // 分析目标事件
	ReturnEvents  string         `gorm:"type:varchar(500)" json:"return_events"` // 留存判断事件(逗号分隔)
	PeriodType    string         `gorm:"type:varchar(20);default:'day'" json:"period_type"` // day, week, month
	PeriodCount   int            `gorm:"type:int;default:30" json:"period_count"` // 分析周期数
	SegmentFilter string         `gorm:"type:jsonb" json:"segment_filter"` // 分群条件 JSON
	StartDate     string         `gorm:"type:varchar(10)" json:"start_date"` // YYYY-MM-DD
	EndDate       string         `gorm:"type:varchar(10)" json:"end_date"` // YYYY-MM-DD
	ResultData    string         `gorm:"type:jsonb" json:"result_data"` // 留存分析结果 JSON
	Status        string         `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, active, archived
	CreatedBy     string         `gorm:"type:varchar(36)" json:"created_by"`
	LastRunAt     *time.Time     `json:"last_run_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RetentionReport) TableName() string {
	return "retention_reports"
}

// FunnelStep 漏斗步骤(用于JSON序列化)
type FunnelStep struct {
	Step      int    `json:"step"`
	Name      string `json:"name"`
	EventName string `json:"event_name"`
	EventType string `json:"event_type"` // event, page
}

// FunnelResult 漏斗计算结果
type FunnelResult struct {
	FunnelID        uint           `json:"funnel_id"`
	TotalUsers      int64          `json:"total_users"`
	Steps           []FunnelStepResult `json:"steps"`
	OverallRate     float64        `json:"overall_rate"` // 总体转化率
	GeneratedAt     time.Time      `json:"generated_at"`
}

// FunnelStepResult 单步漏斗结果
type FunnelStepResult struct {
	Step       int    `json:"step"`
	Name       string `json:"name"`
	EventName  string `json:"event_name"`
	Users      int64  `json:"users"`
	Rate       float64 `json:"rate"` // 相对上一步转化率
	DropRate   float64 `json:"drop_rate"` // 流失率
}

// CohortResult 群组分析结果
type CohortResult struct {
	CohortID       uint              `json:"cohort_id"`
	CohortType     string            `json:"cohort_type"`
	Cohorts        []CohortData      `json:"cohorts"`
	GeneratedAt    time.Time         `json:"generated_at"`
}

// CohortData 单个群组数据
type CohortData struct {
	CohortPeriod  string          `json:"cohort_period"` // 群组周期，如 "2026-01"
	CohortSize    int64           `json:"cohort_size"`   // 群组大小
	RetentionData []RetentionPoint `json:"retention_data"` // 各周期留存率
}

// RetentionPoint 留存点
type RetentionPoint struct {
	Period    int     `json:"period"` // 周期序号 0,1,2...
	Retention float64 `json:"retention"` // 留存率 0.0-1.0
	Count     int64   `json:"count"` // 绝对人数
}

// EventSummary 事件汇总统计
type EventSummary struct {
	TotalEvents   int64             `json:"total_events"`
	TodayEvents   int64             `json:"today_events"`
	WeekEvents    int64             `json:"week_events"`
	TopEvents     []EventCount      `json:"top_events"` // Top10 事件
	ByType        []EventTypeCount  `json:"by_type"`
	ByDevice      int64             `json:"by_device"`
	ByUser        int64             `json:"by_user"`
	AvgDuration   float64           `json:"avg_duration"`
	GeneratedAt   time.Time         `json:"generated_at"`
}

// EventCount 事件计数
type EventCount struct {
	EventName string `json:"event_name"`
	Count     int64  `json:"count"`
}

// EventTypeCount 事件类型计数
type EventTypeCount struct {
	EventType string `json:"event_type"`
	Count     int64  `json:"count"`
}

// AnalyticsDashboard 分析仪表板数据
type AnalyticsDashboard struct {
	TotalEvents   int64           `json:"total_events"`
	TodayEvents    int64          `json:"today_events"`
	ActiveFunnels  int64          `json:"active_funnels"`
	ActiveCohorts  int64          `json:"active_cohorts"`
	RetentionRate  float64        `json:"retention_rate"`
	TopEvents      []EventCount   `json:"top_events"`
	FunnelSummary  []FunnelSummary `json:"funnel_summary"`
	CohortSummary  []CohortSummary `json:"cohort_summary"`
	GeneratedAt    time.Time      `json:"generated_at"`
}

// FunnelSummary 漏斗摘要
type FunnelSummary struct {
	FunnelID   uint    `json:"funnel_id"`
	Name       string  `json:"name"`
	TotalUsers int64   `json:"total_users"`
	FinalRate  float64 `json:"final_rate"`
}

// CohortSummary 群组摘要
type CohortSummary struct {
	CohortID      uint    `json:"cohort_id"`
	Name          string  `json:"name"`
	CohortType    string  `json:"cohort_type"`
	CohortSize    int64   `json:"cohort_size"`
	LatestRetention float64 `json:"latest_retention"`
}
