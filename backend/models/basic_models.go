package models

import (
	"time"
)

// SysDictType 字典类型表
type SysDictType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TypeCode  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"type_code"`
	TypeName  string    `gorm:"type:varchar(100);not null" json:"type_name"`
	Status    int       `gorm:"default:1" json:"status"`
	Remark    string    `gorm:"type:varchar(255)" json:"remark"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SysDictType) TableName() string { return "sys_dict_types" }

// SysDictItem 字典项表
type SysDictItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DictTypeID uint      `gorm:"not null;index" json:"dict_type_id"`
	ItemText   string    `gorm:"type:varchar(100);not null" json:"item_text"`
	ItemValue  string    `gorm:"type:varchar(100);not null" json:"item_value"`
	Sort       int       `gorm:"default:0" json:"sort"`
	Status     int       `gorm:"default:1" json:"status"`
	Remark     string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (SysDictItem) TableName() string { return "sys_dict_items" }

// SysNumberRule 编号规则表
type SysNumberRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RuleName     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"rule_name"`
	Prefix       string    `gorm:"type:varchar(20)" json:"prefix"`
	DateFormat   string    `gorm:"type:varchar(50)" json:"date_format"`
	SeqFormat    string    `gorm:"type:varchar(20)" json:"seq_format"`
	CurrentValue uint      `gorm:"default:0" json:"current_value"`
	Increment    int       `gorm:"default:1" json:"increment"`
	Status       int       `gorm:"default:1" json:"status"`
	Remark       string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SysNumberRule) TableName() string { return "sys_number_rules" }

// SysScheduleJob 调度任务表
type SysScheduleJob struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	JobName     string     `gorm:"type:varchar(100);not null" json:"job_name"`
	JobType     string     `gorm:"type:varchar(50);not null" json:"job_type"`
	CronExpr    string     `gorm:"type:varchar(100)" json:"cron_expr"`
	ApiUrl      string     `gorm:"type:varchar(500)" json:"api_url"`
	HttpMethod  string     `gorm:"type:varchar(10)" json:"http_method"`
	Headers     string     `gorm:"type:text" json:"headers"`
	BodyTpl     string     `gorm:"type:text" json:"body_tpl"`
	Enabled     int        `gorm:"default:1" json:"enabled"`
	LastRunAt   *time.Time `json:"last_run_at"`
	NextRunAt   *time.Time `json:"next_run_at"`
	RunCount    int        `gorm:"default:0" json:"run_count"`
	Status      string     `gorm:"type:varchar(20);default:'idle'" json:"status"`
	LastResult  string     `gorm:"type:text" json:"last_result"`
	StatusMsg   string     `gorm:"type:varchar(255)" json:"status_msg"`
	Remark      string     `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (SysScheduleJob) TableName() string { return "sys_schedule_jobs" }
