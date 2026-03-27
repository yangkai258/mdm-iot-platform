package models

import (
	"time"

	"gorm.io/gorm"
)

// OfflineCache 本地缓存的设备数据
type OfflineCache struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DeviceID   string         `gorm:"type:varchar(36);index" json:"device_id"`
	DataType   string         `gorm:"type:varchar(50);index" json:"data_type"` // device_status, sensor_data, etc.
	CachedData string         `gorm:"type:text" json:"cached_data"`           // JSON
	CachedAt   time.Time      `json:"cached_at"`
	SyncStatus string         `gorm:"type:varchar(20);default:'pending'" json:"sync_status"` // synced/pending/failed
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// OfflineOperation 离线操作的队列
type OfflineOperation struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DeviceID   string         `gorm:"type:varchar(36);index" json:"device_id"`
	Operation  string         `gorm:"type:varchar(50)" json:"operation"` // control/setting/update
	Payload    string         `gorm:"type:text" json:"payload"`         // JSON
	CreatedAt  time.Time      `json:"created_at"`
	Status     string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending/syncing/completed/failed
	SyncedAt   *time.Time     `json:"synced_at"`
	ErrorMsg   string         `gorm:"type:text" json:"error_msg"`
	RetryCount int            `gorm:"default:0" json:"retry_count"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// SyncConflict 同步冲突记录
type SyncConflict struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeviceID      string    `gorm:"type:varchar(36);index" json:"device_id"`
	OperationID  uint      `gorm:"index" json:"operation_id"`
	ConflictType  string    `gorm:"type:varchar(20)" json:"conflict_type"` // timestamp/server_wins/client_wins
	ClientData    string    `gorm:"type:text" json:"client_data"`
	ServerData    string    `gorm:"type:text" json:"server_data"`
	ResolvedData  string    `gorm:"type:text" json:"resolved_data"`
	Resolution    string    `gorm:"type:varchar(20)" json:"resolution"` // resolved_with_server/resolved_with_client/manual
	CreatedAt     time.Time `json:"created_at"`
	ResolvedAt    *time.Time `json:"resolved_at"`
}

// TableName 设置表名
func (OfflineCache) TableName() string {
	return "offline_caches"
}

func (OfflineOperation) TableName() string {
	return "offline_operations"
}

func (SyncConflict) TableName() string {
	return "sync_conflicts"
}
