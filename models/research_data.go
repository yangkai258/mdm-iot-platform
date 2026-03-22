package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResearchDataset 行为研究数据集
type ResearchDataset struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DatasetID    string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"dataset_id"`
	Name         string         `gorm:"type:varchar(128);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	DataType     string         `gorm:"type:varchar(32);not null" json:"data_type"` // behavior_log, emotion_record, health_metric, sensor_event
	Source       string         `gorm:"type:varchar(64)" json:"source"`              // pet_device, user_app, lab_equipment
	DateFrom     *time.Time     `gorm:"type:date" json:"date_from"`
	DateTo       *time.Time     `gorm:"type:date" json:"date_to"`
	Filters      string         `gorm:"type:jsonb" json:"filters"`                  // 匿名化过滤条件 JSON
	RecordCount  int            `gorm:"default:0" json:"record_count"`
	FileFormat   string         `gorm:"type:varchar(16);default:'json'" json:"file_format"` // json, csv, parquet
	FilePath     string         `gorm:"type:varchar(512)" json:"file_path"`           // 导出文件路径
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"`     // draft, ready, exporting, exported, failed
	Anonymized   bool           `gorm:"default:true" json:"anonymized"`              // 是否已匿名化
	OrgID        uint           `gorm:"index" json:"org_id"`
	CreateUserID uint           `gorm:"index" json:"create_user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (r *ResearchDataset) BeforeCreate(tx *gorm.DB) error {
	if r.DatasetID == "" {
		r.DatasetID = uuid.New().String()
	}
	return nil
}

// ResearchDataRecord 行为研究数据记录（匿名化的原始数据快照引用）
type ResearchDataRecord struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RecordID        string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"record_id"`
	DatasetID       string         `gorm:"type:varchar(36);index" json:"dataset_id"`
	OriginalType    string         `gorm:"type:varchar(32);not null" json:"original_type"`    // 对应源表: behavior_event, emotion_record, etc.
	OriginalID      string         `gorm:"type:varchar(36);index" json:"original_id"`        // 源记录ID（已匿名化）
	AnonymizedID    string         `gorm:"type:varchar(36);index" json:"anonymized_id"`      // 匿名化后的设备/用户ID
	DataSnapshot    string         `gorm:"type:jsonb" json:"data_snapshot"`                  // 匿名化后的数据快照
	Timestamp       time.Time      `gorm:"index" json:"timestamp"`
	DeviceType      string         `gorm:"type:varchar(32)" json:"device_type"`             // 设备类型（无具体设备ID）
	BehaviorType    string         `gorm:"type:varchar(64);index" json:"behavior_type"`     // 行为类型
	OrgID           uint           `gorm:"index" json:"org_id"`
	CreatedAt       time.Time      `json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (r *ResearchDataRecord) BeforeCreate(tx *gorm.DB) error {
	if r.RecordID == "" {
		r.RecordID = uuid.New().String()
	}
	return nil
}

// ResearchExportJob 导出任务
type ResearchExportJob struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	JobID       string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"job_id"`
	DatasetID   string         `gorm:"type:varchar(36);index" json:"dataset_id"`
	Status      string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, processing, completed, failed
	Format      string         `gorm:"type:varchar(16);default:'json'" json:"format"`
	FilePath    string         `gorm:"type:varchar(512)" json:"file_path"`
	FileSize    int64          `gorm:"default:0" json:"file_size"`
	RecordCount int            `gorm:"default:0" json:"record_count"`
	ErrorMsg    string         `gorm:"type:text" json:"error_msg"`
	OrgID       uint           `gorm:"index" json:"org_id"`
	CreateUserID uint          `gorm:"index" json:"create_user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 创建前自动生成 UUID
func (r *ResearchExportJob) BeforeCreate(tx *gorm.DB) error {
	if r.JobID == "" {
		r.JobID = uuid.New().String()
	}
	return nil
}
