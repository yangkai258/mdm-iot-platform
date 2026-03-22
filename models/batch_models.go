package models

import (
	"time"
)

// BatchTaskCreateRequest 批量任务创建请求
type BatchTaskCreateRequest struct {
	DeviceIDs []string `json:"device_ids" binding:"required"`
	Action    string   `json:"action" binding:"required"`
	Params    JSON     `json:"params"`
}

// BatchTask 批量任务
type BatchTask struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    string    `json:"task_id" gorm:"uniqueIndex;size:64"`
	TaskType  string    `json:"task_type" gorm:"size:32"` // upgrade/restart/config
	Total     int       `json:"total"`
	Success   int       `json:"success" gorm:"default:0"`
	Failed    int       `json:"failed" gorm:"default:0"`
	Pending   int       `json:"pending" gorm:"default:0"`
	Status    string    `json:"status" gorm:"size:20"` // pending/running/completed
	Results   string    `json:"results" gorm:"type:text"` // JSON array as string
	CreatorID uint      `json:"creator_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (BatchTask) TableName() string {
	return "batch_tasks"
}

// BatchTaskDeviceResult 批量任务中单个设备的结果
type BatchTaskDeviceResult struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"` // success/failed
	Error    string `json:"error,omitempty"`
}
