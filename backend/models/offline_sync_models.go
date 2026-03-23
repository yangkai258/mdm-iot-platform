package models

import "time"

// OfflineCache 本地缓存
type OfflineCache struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	DeviceID   uint       `gorm:"index" json:"device_id"`
	UserID     uint       `gorm:"index" json:"user_id"`
	CacheKey   string     `gorm:"index" json:"cache_key"`
	CacheData  string     `json:"cache_data"`
	SyncStatus string     `gorm:"default:'pending'" json:"sync_status"`
	CreatedAt  time.Time  `json:"created_at"`
	SyncedAt   *time.Time `json:"synced_at"`
	RetryCount int        `gorm:"default:0" json:"retry_count"`
}

// OfflineQueue 离线操作队列
type OfflineQueue struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	DeviceID   uint       `gorm:"index" json:"device_id"`
	UserID     uint       `gorm:"index" json:"user_id"`
	ActionType string     `json:"action_type"`
	ActionData string     `json:"action_data"`
	Status     string     `gorm:"default:'pending'" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	SentAt     *time.Time `json:"sent_at"`
	AckedAt    *time.Time `json:"acked_at"`
	ErrorMsg   string     `json:"error_msg"`
}
