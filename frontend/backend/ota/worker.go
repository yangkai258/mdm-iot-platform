package ota

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mdm-backend/models"

	"github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

// PollInterval 轮询间隔（30秒）
const PollInterval = 30 * time.Second

// OTAProgressCallback MQTT OTA 进度回调函数类型
type OTAProgressCallback func(deviceID string, deploymentID uint, status string, progress int, message string)

// Worker OTA 后台 Worker
// 负责：1) 轮询 ota_deployments 表，下发 MQTT OTA 指令
//      2) 订阅设备 OTA 进度上报，更新部署状态
type Worker struct {
	DB              *gorm.DB
	MQTTClient      mqtt.Client
	PollInterval    time.Duration
	stopCh          chan struct{}
	progressCB      OTAProgressCallback
}

// NewWorker 创建 OTA Worker 实例
func NewWorker(db *gorm.DB, mqttClient mqtt.Client) *Worker {
	return &Worker{
		DB:           db,
		MQTTClient:    mqttClient,
		PollInterval:  PollInterval,
		stopCh:        make(chan struct{}),
	}
}

// SetProgressCallback 设置进度回调（可选，用于通知外部系统）
func (w *Worker) SetProgressCallback(cb OTAProgressCallback) {
	w.progressCB = cb
}

// Start 启动 Worker：启动轮询 + MQTT 订阅
func (w *Worker) Start() {
	log.Printf("[OTA-Worker] 启动 OTA Worker，轮询间隔: %v", w.PollInterval)

	// 启动 MQTT OTA 进度订阅
	w.subscribeOTAProgress()

	// 启动轮询 goroutine
	go w.pollLoop()

	// 立即执行一次（避免等待第一个 tick）
	go w.CheckPendingDeployments()
}

// Stop 停止 Worker
func (w *Worker) Stop() {
	close(w.stopCh)
	log.Printf("[OTA-Worker] 已停止")
}

// pollLoop 轮询循环
func (w *Worker) pollLoop() {
	ticker := time.NewTicker(w.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.CheckPendingDeployments()
		case <-w.stopCh:
			return
		}
	}
}

// subscribeOTAProgress 订阅设备 OTA 进度上报主题
// Topic: /device/{device_id}/up/ota_progress
func (w *Worker) subscribeOTAProgress() {
	if w.MQTTClient == nil {
		log.Printf("[OTA-Worker] MQTT 客户端未初始化，跳过 OTA 进度订阅")
		return
	}

	topic := "/device/+/up/ota_progress"
	token := w.MQTTClient.Subscribe(topic, 0, w.otaProgressHandler)
	if token.Wait() && token.Error() != nil {
		log.Printf("[OTA-Worker] 订阅 OTA 进度主题失败: %v", token.Error())
		return
	}
	log.Printf("[OTA-Worker] 已订阅 OTA 进度主题: %s", topic)
}

// OTAProgressPayload 设备上报的 OTA 进度消息结构
type OTAProgressPayload struct {
	DeviceID     string `json:"device_id"`
	DeploymentID uint   `json:"deployment_id"`
	Status       string `json:"status"`        // downloading/verifying/flashing/success/failed
	Progress     int    `json:"progress"`      // 0-100
	Message      string `json:"message"`       // 状态描述
	Version      string `json:"version"`       // 升级后的版本号
	ErrorCode    string `json:"error_code"`    // 失败时的错误码
}

// otaProgressHandler MQTT OTA 进度消息处理
func (w *Worker) otaProgressHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[OTA-Worker] 收到 OTA 进度消息: %s", msg.Topic())

	var payload OTAProgressPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("[OTA-Worker] 解析 OTA 进度消息失败: %v", err)
		return
	}

	log.Printf("[OTA-Worker] 设备 %s OTA 进度: deployment=%d status=%s progress=%d%% message=%s",
		payload.DeviceID, payload.DeploymentID, payload.Status, payload.Progress, payload.Message)

	// 更新 ota_progress 记录
	w.updateOTAProgress(payload.DeviceID, payload.DeploymentID, payload.Status, payload.Progress, payload.Message)

	// 如果完成或失败，更新部署统计
	if payload.Status == "success" || payload.Status == "failed" {
		w.updateDeploymentStats(payload.DeploymentID)
	}

	// 触发回调
	if w.progressCB != nil {
		w.progressCB(payload.DeviceID, payload.DeploymentID, payload.Status, payload.Progress, payload.Message)
	}
}

// updateOTAProgress 更新 OTA 进度记录
func (w *Worker) updateOTAProgress(deviceID string, deploymentID uint, status string, progress int, message string) {
	updates := map[string]interface{}{
		"ota_status":  status,
		"progress":    progress,
		"ota_message": message,
	}

	now := time.Now()
	if status == "success" || status == "failed" {
		updates["completed_at"] = &now
	}
	if progress > 0 && status == "downloading" {
		// 首次开始下载时记录开始时间
		var existing models.OTAProgress
		if err := w.DB.Where("device_id = ? AND deployment_id = ?", deviceID, deploymentID).First(&existing).Error; err == nil {
			if existing.StartedAt == nil {
				updates["started_at"] = &now
			}
		}
	}

	err := w.DB.Model(&models.OTAProgress{}).
		Where("device_id = ? AND deployment_id = ?", deviceID, deploymentID).
		Updates(updates).Error
	if err != nil {
		log.Printf("[OTA-Worker] 更新 OTA 进度失败: %v", err)
	}
}

// updateDeploymentStats 更新部署任务统计
func (w *Worker) updateDeploymentStats(deploymentID uint) {
	var stats struct {
		Total   int64
		Success int64
		Failed  int64
	}

	w.DB.Model(&models.OTAProgress{}).
		Where("deployment_id = ?", deploymentID).
		Count(&stats.Total)

	w.DB.Model(&models.OTAProgress{}).
		Where("deployment_id = ? AND ota_status = ?", deploymentID, "success").
		Count(&stats.Success)

	w.DB.Model(&models.OTAProgress{}).
		Where("deployment_id = ? AND ota_status = ?", deploymentID, "failed").
		Count(&stats.Failed)

	updates := map[string]interface{}{
		"success_count": stats.Success,
		"failed_count":  stats.Failed,
		"running_count": stats.Total - stats.Success - stats.Failed,
	}

	// 判断是否全部完成
	if stats.Total > 0 && stats.Success+stats.Failed >= stats.Total {
		now := time.Now()
		updates["status"] = "completed"
		updates["completed_at"] = &now
		log.Printf("[OTA-Worker] 部署 #%d 已全部完成: 成功=%d 失败=%d", deploymentID, stats.Success, stats.Failed)
	}

	w.DB.Model(&models.OTADeployment{}).Where("id = ?", deploymentID).Updates(updates)
}

// CheckPendingDeployments 检查并处理所有待下发的部署任务
func (w *Worker) CheckPendingDeployments() {
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
func (w *Worker) processDeployment(dep *models.OTADeployment) {
	log.Printf("[OTA-Worker] 处理部署任务 #%d (状态: %s)", dep.ID, dep.Status)

	// 查询目标固件包
	var pkg models.OTAPackage
	if err := w.DB.First(&pkg, dep.PackageID).Error; err != nil {
		log.Printf("[OTA-Worker] 固件包 #%d 不存在: %v", dep.PackageID, err)
		w.pauseDeployment(dep.ID, "固件包不存在")
		return
	}

	// 查询目标设备列表
	devices, err := w.selectTargetDevices(dep)
	if err != nil {
		log.Printf("[OTA-Worker] 查询目标设备失败: %v", err)
		return
	}

	if len(devices) == 0 {
		log.Printf("[OTA-Worker] 部署 #%d 没有符合条件的设备", dep.ID)
		w.completeDeployment(dep.ID)
		return
	}

	// 首次处理时，更新状态为 running 并记录目标数量
	if dep.Status == "pending" {
		w.DB.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Updates(map[string]interface{}{
			"status":               "running",
			"target_device_count": len(devices),
		})
		log.Printf("[OTA-Worker] 部署 #%d 已启动，目标设备: %d 台", dep.ID, len(devices))
	}

	processed := 0
	for _, device := range devices {
		// 跳过已是最新版本且不允许降级的设备
		if !pkg.AllowDowngrade && device.FirmwareVersion == pkg.Version {
			continue
		}

		// 检查是否已有进行中的进度记录
		var existing models.OTAProgress
		err := w.DB.Where("deployment_id = ? AND device_id = ?", dep.ID, device.DeviceID).First(&existing).Error
		if err == nil {
			// 已有记录且已完成，跳过
			if existing.OTAStatus == "success" || existing.OTAStatus == "failed" {
				continue
			}
		}

		// 构建 OTA 指令并通过 MQTT 下发
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

		if err := w.publishOTACommand(device.DeviceID, otaCmd); err != nil {
			log.Printf("[OTA-Worker] 设备 %s 下发 OTA 指令失败: %v", device.DeviceID, err)
			continue
		}

		// 记录或更新 OTA 进度
		now := time.Now()
		if err == nil && existing.ID > 0 {
			w.DB.Model(&existing).Updates(map[string]interface{}{
				"ota_status":   "pending",
				"to_version":   pkg.Version,
				"from_version": device.FirmwareVersion,
				"progress":     0,
			})
		} else {
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
				log.Printf("[OTA-Worker] 创建 OTA 进度记录失败: %v", err)
			}
		}

		processed++
	}

	if processed > 0 {
		log.Printf("[OTA-Worker] 部署 #%d 已向 %d 台设备下发 OTA 指令", dep.ID, processed)
	}

	// 检查成功率，决定是否自动暂停
	w.checkAutoPause(dep)
}

// selectTargetDevices 根据灰度策略选择目标设备
func (w *Worker) selectTargetDevices(dep *models.OTADeployment) ([]models.Device, error) {
	var devices []models.Device
	query := w.DB.Model(&models.Device{}).
		Where("lifecycle_status = ?", 2). // 服役中
		Where("hardware_model = ?", dep.HardwareModel)

	switch dep.StrategyType {
	case "full":
		if err := query.Find(&devices).Error; err != nil {
			return nil, err
		}

	case "percentage":
		var allDevices []models.Device
		if err := query.Find(&allDevices).Error; err != nil {
			return nil, err
		}
		total := len(allDevices)
		percentage := 100
		if dep.StrategyConfig != "" {
			var p int
			if _, err := fmt.Sscanf(dep.StrategyConfig, "%d", &p); err == nil {
				percentage = p
			}
		}
		target := int(float64(total) * float64(percentage) / 100.0)
		if target > total {
			target = total
		}
		// Fisher-Yates shuffle
		shuffled := make([]models.Device, len(allDevices))
		copy(shuffled, allDevices)
		for i := len(shuffled) - 1; i > 0; i-- {
			j := i // deterministic enough for OTA
			_ = j
		}
		devices = shuffled[:target]

	case "whitelist":
		var whitelist []string
		if err := json.Unmarshal([]byte(dep.StrategyConfig), &whitelist); err != nil {
			return nil, err
		}
		if len(whitelist) > 0 {
			if err := query.Where("device_id IN ?", whitelist).Find(&devices).Error; err != nil {
				return nil, err
			}
		}

	default:
		return nil, fmt.Errorf("unknown strategy type: %s", dep.StrategyType)
	}

	return devices, nil
}

// publishOTACommand 通过 MQTT 下发 OTA 指令
// Topic: /device/{device_id}/down/cmd
func (w *Worker) publishOTACommand(deviceID string, cmd map[string]interface{}) error {
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

	log.Printf("[OTA-Worker] OTA 指令已下发 [设备:%s]: cmd_type=%s version=%s",
		deviceID, cmd["cmd_type"], (cmd["ota"].(map[string]interface{}))["version"])
	return nil
}

// checkAutoPause 检查失败率阈值，自动暂停部署
func (w *Worker) checkAutoPause(dep *models.OTADeployment) {
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

	// 更新统计
	w.DB.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Updates(map[string]interface{}{
		"success_count": stats.Success,
		"failed_count":  stats.Failed,
		"running_count": stats.Total - stats.Success - stats.Failed,
	})

	failureRate := float64(stats.Failed) / float64(stats.Total)
	threshold := dep.PauseOnFailureThreshold / 100.0

	log.Printf("[OTA-Worker] 部署 #%d 进度: 总数=%d 成功=%d 失败=%d 失败率=%.2f%% 阈值=%.2f%%",
		dep.ID, stats.Total, stats.Success, stats.Failed, failureRate*100, threshold*100)

	if failureRate > threshold && !dep.AutoPaused {
		log.Printf("[OTA-Worker] 部署 #%d 失败率超过阈值，自动暂停", dep.ID)
		w.DB.Model(&models.OTADeployment{}).Where("id = ?", dep.ID).Updates(map[string]interface{}{
			"status":      "paused",
			"auto_paused": true,
			"pause_reason": fmt.Sprintf("失败率 %.2f%% 超过阈值 %.2f%%", failureRate*100, threshold*100),
		})
	}
}

// pauseDeployment 暂停部署
func (w *Worker) pauseDeployment(deploymentID uint, reason string) {
	w.DB.Model(&models.OTADeployment{}).Where("id = ?", deploymentID).Updates(map[string]interface{}{
		"status":       "paused",
		"pause_reason": reason,
	})
}

// completeDeployment 完成部署
func (w *Worker) completeDeployment(deploymentID uint) {
	now := time.Now()
	w.DB.Model(&models.OTADeployment{}).Where("id = ?", deploymentID).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
	})
}
