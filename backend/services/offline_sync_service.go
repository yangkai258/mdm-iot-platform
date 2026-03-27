package services

import (
	"encoding/json"
	"errors"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// ConflictResolutionMode 冲突解决模式
type ConflictResolutionMode string

const (
	ConflictTimestamp      ConflictResolutionMode = "timestamp"      // 时间戳优先
	ConflictServerWins     ConflictResolutionMode = "server_wins"     // 服务端优先
	ConflictClientWins     ConflictResolutionMode = "client_wins"     // 客户端优先
)

// OfflineSyncService 离线同步服务
type OfflineSyncService struct {
	DB *gorm.DB
}

// NewOfflineSyncService 创建离线同步服务
func NewOfflineSyncService(db *gorm.DB) *OfflineSyncService {
	return &OfflineSyncService{DB: db}
}

// ProcessOfflineQueue 处理离线队列
func (s *OfflineSyncService) ProcessOfflineQueue(deviceID string) (map[string]interface{}, error) {
	var operations []models.OfflineOperation
	if err := s.DB.Where("device_id = ? AND status = ?", deviceID, "pending").
		Order("created_at ASC").
		Find(&operations).Error; err != nil {
		return nil, err
	}

	results := map[string]interface{}{
		"total":     len(operations),
		"succeeded": 0,
		"failed":    0,
		"items":     []map[string]interface{}{},
	}

	for _, op := range operations {
		// 更新状态为 syncing
		s.DB.Model(&op).Updates(map[string]interface{}{
			"status": "syncing",
		})

		// 执行同步操作
		err := s.executeOperation(op)
		if err != nil {
			s.DB.Model(&op).Updates(map[string]interface{}{
				"status":     "failed",
				"error_msg":  err.Error(),
				"retry_count": op.RetryCount + 1,
			})
			results["failed"] = results["failed"].(int) + 1
		} else {
			now := time.Now()
			s.DB.Model(&op).Updates(map[string]interface{}{
				"status":    "completed",
				"synced_at": &now,
			})
			results["succeeded"] = results["succeeded"].(int) + 1
		}

		results["items"] = append(results["items"].([]map[string]interface{}), map[string]interface{}{
			"id":     op.ID,
			"status": "syncing",
		})
	}

	return results, nil
}

// executeOperation 执行单个离线操作
func (s *OfflineSyncService) executeOperation(op models.OfflineOperation) error {
	// 根据操作类型执行不同逻辑
	switch op.Operation {
	case "control":
		return s.executeControl(op)
	case "setting":
		return s.executeSetting(op)
	case "update":
		return s.executeUpdate(op)
	default:
		return errors.New("unknown operation type: " + op.Operation)
	}
}

// executeControl 执行控制指令
func (s *OfflineSyncService) executeControl(op models.OfflineOperation) error {
	// TODO: 根据 Payload 发送 MQTT 控制指令到设备
	// 这里仅做演示，实际应调用 mqtt 服务发送指令
	return nil
}

// executeSetting 执行设置操作
func (s *OfflineSyncService) executeSetting(op models.OfflineOperation) error {
	// TODO: 更新设备配置
	return nil
}

// executeUpdate 执行数据更新
func (s *OfflineSyncService) executeUpdate(op models.OfflineOperation) error {
	// TODO: 更新设备状态到数据库
	return nil
}

// ConflictResolution 冲突解决
func (s *OfflineSyncService) ConflictResolution(opID uint, mode ConflictResolutionMode, resolvedData string) error {
	var conflict models.SyncConflict
	if err := s.DB.Where("operation_id = ?", opID).First(&conflict).Error; err != nil {
		return err
	}

	now := time.Now()
	switch mode {
	case ConflictServerWins:
		// 服务端数据优先，无需处理
		conflict.Resolution = "resolved_with_server"
	case ConflictClientWins:
		// 客户端数据优先，使用 resolvedData
		conflict.ResolvedData = resolvedData
		conflict.Resolution = "resolved_with_client"
	case ConflictTimestamp:
		// 时间戳优先：比较时间后决定
		var serverData, clientData map[string]interface{}
		json.Unmarshal([]byte(conflict.ServerData), &serverData)
		json.Unmarshal([]byte(conflict.ClientData), &clientData)

		serverTime, _ := time.Parse(time.RFC3339, serverData["updated_at"].(string))
		clientTime, _ := time.Parse(time.RFC3339, clientData["updated_at"].(string))

		if serverTime.After(clientTime) {
			conflict.Resolution = "resolved_with_server"
		} else {
			conflict.ResolvedData = resolvedData
			conflict.Resolution = "resolved_with_client"
		}
	}

	conflict.ResolvedAt = &now
	return s.DB.Save(&conflict).Error
}

// GetPendingOperations 获取待同步操作
func (s *OfflineSyncService) GetPendingOperations(deviceID string) ([]models.OfflineOperation, error) {
	var operations []models.OfflineOperation
	err := s.DB.Where("device_id = ? AND status = ?", deviceID, "pending").
		Order("created_at ASC").
		Find(&operations).Error
	return operations, err
}

// GetCachedData 获取设备缓存数据
func (s *OfflineSyncService) GetCachedData(deviceID, dataType string) ([]models.OfflineCache, error) {
	var caches []models.OfflineCache
	query := s.DB.Where("device_id = ?", deviceID)
	if dataType != "" {
		query = query.Where("data_type = ?", dataType)
	}
	err := query.Order("cached_at DESC").Find(&caches).Error
	return caches, err
}

// SaveCache 保存设备缓存数据
func (s *OfflineSyncService) SaveCache(deviceID, dataType, cachedData string) (*models.OfflineCache, error) {
	cache := models.OfflineCache{
		DeviceID:   deviceID,
		DataType:   dataType,
		CachedData: cachedData,
		CachedAt:   time.Now(),
		SyncStatus: "pending",
	}
	if err := s.DB.Create(&cache).Error; err != nil {
		return nil, err
	}
	return &cache, nil
}

// MarkCacheSynced 标记缓存已同步
func (s *OfflineSyncService) MarkCacheSynced(cacheID uint) error {
	now := time.Now()
	return s.DB.Model(&models.OfflineCache{}).Where("id = ?", cacheID).Updates(map[string]interface{}{
		"sync_status": "synced",
		"cached_at":   &now,
	}).Error
}

// AutoMigrateOfflineModels 自动迁移离线相关表
func AutoMigrateOfflineModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.OfflineCache{},
		&models.OfflineOperation{},
		&models.SyncConflict{},
	)
}
