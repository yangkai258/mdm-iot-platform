package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"mdm-backend/models"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

// DefaultPollInterval 默认轮询间隔
const DefaultPollInterval = 5 * time.Minute

// DefaultPauseThreshold 默认暂停阈值（成功率低于此值则自动暂停）
const DefaultPauseThreshold = 0.80

// OTAWorker OTA后台Worker
type OTAWorker struct {
	DB             *gorm.DB
	MQTTClient     pahomqtt.Client
	PollInterval   time.Duration
	PauseThreshold float64 // 0.0-1.0
	stopCh         chan struct{}
}

// NewOTAWorker 创建OTA Worker实例
func NewOTAWorker(db *gorm.DB, mqttClient pahomqtt.Client) *OTAWorker {
	return &OTAWorker{
		DB:             db,
		MQTTClient:     mqttClient,
		PollInterval:   DefaultPollInterval,
		PauseThreshold: DefaultPauseThreshold,
		stopCh:         make(chan struct{}),
	}
}

// Start 启动OTA Worker轮询
func (w *OTAWorker) Start() {
	log.Printf("[OTA-Worker] 启动OTA后台Worker，轮询间隔: %v", w.PollInterval)
	ticker := time.NewTicker(w.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.CheckPendingDeployments()
		case <-w.stopCh:
			log.Printf("[OTA-Worker] 已停止")
			return
		}
	}
}

// Stop 停止OTA Worker
func (w *OTAWorker) Stop() {
	close(w.stopCh)
}

// CheckPendingDeployments 检查并处理所有待下发的部署任务
func (w *OTAWorker) CheckPendingDeployments() {
	var deployments []models.OTADeployment
	if err := w.DB.Model(&models.OTADeployment{}).
		Where("status IN ?", []string{"pending", "running"}).
		Find(&deployments).Error; err != nil {
		log.Printf("[OTA-Worker] 查询部署任务失败: %v", err)
		return
	}

	for _, dep := range deployments {
		w.processDeployment(&dep)
	}
}

// processDeployment 处理单个部署任务
func (w *OTAWorker) processDeployment(dep *models.OTADeployment) {
	log.Printf("[OTA-Worker] 处理部署任务 #%d (状态: %s)", dep.ID, dep.Status)

	// 查询目标固件包
	var pkg models.OTAPackage
	if err := w.DB.First(&pkg, dep.PackageID).Error; err != nil {
		log.Printf("[OTA-Worker] 固件包 #%d 不存在: %v", dep.PackageID, err)
		w.pauseDeployment(dep.ID, "固件包不存在")
		return
	}

	// 查询目标设备列表（根据灰度策略）
	devices, err := w.SelectTargetDevices(dep)
	if err != nil {
		log.Printf("[OTA-Worker] 查询目标设备失败: %v", err)
		return
	}

	if len(devices) == 0 {
		log.Printf("[OTA-Worker] 部署 #%d 没有符合条件的设备", dep.ID)
		w.completeDeployment(dep.ID)
		return
	}

	// 更新部署状态为 running，并记录目标设备数量
	if dep.Status == "pending" {
		w.DB.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Updates(map[string]interface{}{
			"status":              "running",
			"target_device_count": len(devices),
		})
	}

	for _, device := range devices {
		// 跳过已是最新版本的设备
		if device.FirmwareVersion == pkg.Version {
			continue
		}

		// 检查是否已有 progress 记录且已完成
		var existing models.OTAProgress
		err := w.DB.Where("deployment_id = ? AND device_id = ?", dep.ID, device.DeviceID).First(&existing).Error
		if err == nil {
			// 已有记录且已完成，跳过
			if existing.OTAStatus == "success" || existing.OTAStatus == "failed" {
				continue
			}
		}

		// 构建 OTA 指令
		otaCmd := map[string]interface{}{
			"cmd_id":   fmt.Sprintf("ota-%d-%s", dep.ID, device.DeviceID),
			"cmd_type": "ota",
			"ota": map[string]interface{}{
				"version":   pkg.Version,
				"url":       pkg.FileURL,
				"md5":       pkg.FileMD5,
				"size":      pkg.FileSize,
				"mandatory": pkg.IsMandatory,
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}

		// 通过 MQTT 下发到 /device/{device_id}/down/cmd
		if err := w.PublishOTACommand(device.DeviceID, otaCmd); err != nil {
			log.Printf("[OTA-Worker] 设备 %s 下发失败: %v", device.DeviceID, err)
			continue
		}

		// 创建或更新 ota_progress 记录
		now := time.Now()
		if err == nil && existing.ID > 0 {
			// 更新已有记录
			w.DB.Model(&existing).Updates(map[string]interface{}{
				"ota_status":   "pending",
				"to_version":   pkg.Version,
				"from_version": device.FirmwareVersion,
				"progress":     0,
			})
		} else {
			// 创建新记录
			progress := models.OTAProgress{
				DeploymentID: dep.ID,
				DeviceID:     device.DeviceID,
				PackageID:    dep.PackageID,
				FromVersion:  device.FirmwareVersion,
				ToVersion:    pkg.Version,
				OTAStatus:    "pending",
				Progress:     0,
				StartedAt:    &now,
			}
			if err := w.DB.Create(&progress).Error; err != nil {
				log.Printf("[OTA-Worker] 创建进度记录失败: %v", err)
			}
		}
	}

	log.Printf("[OTA-Worker] 部署 #%d 处理完成，目标设备: %d 台", dep.ID, len(devices))

	// 检查是否需要自动暂停
	w.CheckAndAutoPause(dep)
}

// SelectTargetDevices 根据灰度策略选择目标设备
func (w *OTAWorker) SelectTargetDevices(dep *models.OTADeployment) ([]models.Device, error) {
	var devices []models.Device
	query := w.DB.Model(&models.Device{}).
		Where("lifecycle_status = ?", 2). // 2=服役中
		Where("hardware_model = ?", dep.HardwareModel)

	switch dep.StrategyType {
	case "full":
		// 全量：所有目标型号的设备
		if err := query.Find(&devices).Error; err != nil {
			return nil, err
		}

	case "percentage":
		// 百分比灰度：从所有设备中按配置的比例随机选取
		var allDevices []models.Device
		if err := query.Find(&allDevices).Error; err != nil {
			return nil, err
		}
		total := len(allDevices)
		// strategy_config 存储的是百分比数字（如 "30" 表示 30%）
		percentage := 100
		if dep.StrategyConfig != "" {
			if p, err := strconv.Atoi(dep.StrategyConfig); err == nil {
				percentage = p
			}
		}
		target := int(float64(total) * float64(percentage) / 100.0)
		if target > total {
			target = total
		}
		// 简单随机选取（洗牌后取前 target 个）
		shuffled := make([]models.Device, len(allDevices))
		copy(shuffled, allDevices)
		for i := range shuffled {
			j := rand.Intn(i + 1)
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		}
		devices = shuffled[:target]

	case "whitelist":
		// 白名单：解析 strategy_config 中的设备 ID 列表（JSON数组）
		var whitelist []string
		if err := json.Unmarshal([]byte(dep.StrategyConfig), &whitelist); err != nil {
			log.Printf("[OTA-Worker] 解析白名单失败: %v", err)
			return nil, err
		}
		if len(whitelist) > 0 {
			if err := query.Where("device_id IN ?", whitelist).Find(&devices).Error; err != nil {
				return nil, err
			}
		}

	default:
		log.Printf("[OTA-Worker] 未知灰度策略: %s", dep.StrategyType)
		return nil, fmt.Errorf("unknown strategy type: %s", dep.StrategyType)
	}

	return devices, nil
}

// PublishOTACommand 通过 MQTT 下发 OTA 指令
func (w *OTAWorker) PublishOTACommand(deviceID string, cmd map[string]interface{}) error {
	if w.MQTTClient == nil {
		return fmt.Errorf("MQTT 客户端未初始化")
	}

	topic := fmt.Sprintf("/device/%s/down/cmd", deviceID)
	payload, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("序列化指令失败: %w", err)
	}

	token := w.MQTTClient.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT 发布失败: %w", token.Error())
	}

	log.Printf("[OTA-Worker] OTA指令已下发 [设备:%s] -> %s", deviceID, string(payload))
	return nil
}

// CheckAndAutoPause 检查成功率并自动暂停部署
func (w *OTAWorker) CheckAndAutoPause(dep *models.OTADeployment) {
	// 只对 running 状态检查
	if dep.Status != "running" {
		return
	}

	var stats struct {
		Total   int64
		Success int64
		Failed  int64
	}

	w.DB.Model(&models.OTAProgress{}).
		Where("deployment_id = ?", dep.ID).
		Count(&stats.Total)

	w.DB.Model(&models.OTAProgress{}).
		Where("deployment_id = ? AND ota_status = ?", dep.ID, "success").
		Count(&stats.Success)

	w.DB.Model(&models.OTAProgress{}).
		Where("deployment_id = ? AND ota_status = ?", dep.ID, "failed").
		Count(&stats.Failed)

	if stats.Total == 0 {
		return
	}

	successRate := float64(stats.Success) / float64(stats.Total)
	failureRate := 1.0 - successRate

	// 更新部署统计
	w.DB.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Updates(map[string]interface{}{
		"success_count": stats.Success,
		"failed_count":  stats.Failed,
		"running_count": stats.Total - stats.Success - stats.Failed,
	})

	log.Printf("[OTA-Worker] 部署 #%d 进度: 总数=%d 成功=%d 失败=%d 成功率=%.2f%%",
		dep.ID, stats.Total, stats.Success, stats.Failed, successRate*100)

	// 失败率阈值判断（PauseOnFailureThreshold 存储的是百分比，如 20 表示 20%）
	threshold := dep.PauseOnFailureThreshold / 100.0
	if failureRate > threshold && !dep.AutoPaused {
		log.Printf("[OTA-Worker] 部署 #%d 失败率 %.2f%% 超过阈值 %.2f%%，自动暂停",
			dep.ID, failureRate*100, threshold*100)
		w.DB.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Updates(map[string]interface{}{
			"status":       "paused",
			"auto_paused":  true,
			"pause_reason": fmt.Sprintf("失败率 %.2f%% 超过阈值 %.2f%%", failureRate*100, threshold*100),
		})
	}
}

// pauseDeployment 暂停部署任务
func (w *OTAWorker) pauseDeployment(deploymentID uint, reason string) {
	w.DB.Model(&models.OTADeployment{}).Where("id = ?", deploymentID).Updates(map[string]interface{}{
		"status":       "paused",
		"pause_reason": reason,
	})
}

// completeDeployment 完成部署任务
func (w *OTAWorker) completeDeployment(deploymentID uint) {
	now := time.Now()
	w.DB.Model(&models.OTADeployment{}).Where("id = ?", deploymentID).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
	})
}

// UpdateProgress 更新设备升级进度（由设备回调调用）
func (w *OTAWorker) UpdateProgress(deviceID string, deploymentID uint, status string, progress int, message string) error {
	updates := map[string]interface{}{
		"ota_status":  status,
		"progress":    progress,
		"ota_message": message,
	}

	if status == "success" || status == "failed" {
		now := time.Now()
		updates["completed_at"] = &now
	}

	return w.DB.Model(&models.OTAProgress{}).
		Where("device_id = ? AND deployment_id = ?", deviceID, deploymentID).
		Updates(updates).Error
}
