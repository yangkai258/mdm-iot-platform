package services

import (
	"log"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// SubscriptionRenewalWorker 订阅自动续费 Worker
type SubscriptionRenewalWorker struct {
	db        *gorm.DB
	interval  time.Duration
	stopChan  chan struct{}
}

// NewSubscriptionRenewalWorker 创建续费 Worker
func NewSubscriptionRenewalWorker(db *gorm.DB) *SubscriptionRenewalWorker {
	return &SubscriptionRenewalWorker{
		db:       db,
		interval: 1 * time.Hour, // 每小时检查一次
		stopChan: make(chan struct{}),
	}
}

// Start 启动 Worker
func (w *SubscriptionRenewalWorker) Start() {
	log.Println("[SubscriptionRenewalWorker] Started")

	// 立即执行一次
	w.processExpiredSubscriptions()
	w.sendRenewalReminders()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processExpiredSubscriptions()
			w.sendRenewalReminders()
		case <-w.stopChan:
			log.Println("[SubscriptionRenewalWorker] Stopped")
			return
		}
	}
}

// Stop 停止 Worker
func (w *SubscriptionRenewalWorker) Stop() {
	close(w.stopChan)
}

// sendRenewalReminders 发送续费提醒
func (w *SubscriptionRenewalWorker) sendRenewalReminders() {
	// 查找 3 天后到期的订阅
	threeDaysLater := time.Now().AddDate(0, 0, 3)
	sevenDaysLater := threeDaysLater.AddDate(0, 0, 4)

	var subscriptions []models.Subscription
	w.db.Where("status = ? AND auto_renew = ? AND end_date BETWEEN ? AND ? AND (reminder_sent_at IS NULL OR reminder_sent_at < ?)",
		"active", true, threeDaysLater, sevenDaysLater, time.Now().AddDate(0, 0, -1)).
		Find(&subscriptions)

	for _, sub := range subscriptions {
		log.Printf("[SubscriptionRenewalWorker] Sending renewal reminder for subscription %d, user %d", sub.ID, sub.UserID)

		// 创建提醒日志
		now := time.Now()
		renewalLog := models.SubscriptionRenewalLog{
			SubscriptionID: sub.ID,
			Action:        "reminder",
			Amount:        sub.Price,
			Status:        "success",
		}
		w.db.Create(&renewalLog)

		// 更新提醒时间
		w.db.Model(&sub).Update("reminder_sent_at", &now)

		// TODO: 发送通知（邮件/短信/站内信）
	}
}

// processExpiredSubscriptions 处理到期订阅
func (w *SubscriptionRenewalWorker) processExpiredSubscriptions() {
	now := time.Now()

	// 查找已到期且开启自动续费的订阅
	var subscriptions []models.Subscription
	w.db.Where("status = ? AND auto_renew = ? AND end_date < ?",
		"active", true, now).
		Find(&subscriptions)

	for _, sub := range subscriptions {
		log.Printf("[SubscriptionRenewalWorker] Processing renewal for subscription %d, user %d", sub.ID, sub.UserID)

		// 模拟扣费（实际需要调用支付网关）
		success := w.chargeUser(sub.UserID, sub.Price)

		if success {
			// 续费成功
			newEndDate := sub.EndDate.AddDate(0, 0, sub.Duration)
			w.db.Model(&sub).Updates(map[string]interface{}{
				"end_date":      newEndDate,
				"last_renew_at": now,
				"renew_count":   sub.RenewCount + 1,
				"retry_count":  0,
				"status":        "active",
			})

			// 记录日志
			renewalLog := models.SubscriptionRenewalLog{
				SubscriptionID: sub.ID,
				Action:         "renewal_success",
				Amount:         sub.Price,
				Status:         "success",
			}
			w.db.Create(&renewalLog)

			log.Printf("[SubscriptionRenewalWorker] Renewal successful for subscription %d", sub.ID)
		} else {
			// 续费失败
			retryCount := sub.RetryCount + 1

			if retryCount >= 3 {
				// 重试 3 次后仍失败，暂停服务
				log.Printf("[SubscriptionRenewalWorker] Renewal failed 3 times, suspending subscription %d", sub.ID)
				w.db.Model(&sub).Updates(map[string]interface{}{
					"status":           "suspended",
					"suspended_at":     now,
					"retry_count":      retryCount,
					"renew_fail_reason": "payment_failed",
				})

				renewalLog := models.SubscriptionRenewalLog{
					SubscriptionID: sub.ID,
					Action:          "suspended",
					Amount:          sub.Price,
					Status:          "failed",
					FailReason:      "payment_failed",
					RetryCount:      retryCount,
				}
				w.db.Create(&renewalLog)
			} else {
				// 继续重试
				log.Printf("[SubscriptionRenewalWorker] Renewal failed, will retry (attempt %d/3) for subscription %d", retryCount+1, sub.ID)
				w.db.Model(&sub).Update("retry_count", retryCount)

				renewalLog := models.SubscriptionRenewalLog{
					SubscriptionID: sub.ID,
					Action:          "renewal_failed",
					Amount:          sub.Price,
					Status:          "failed",
					FailReason:      "payment_failed",
					RetryCount:      retryCount,
				}
				w.db.Create(&renewalLog)
			}
		}
	}
}

// chargeUser 模拟扣费（实际需要调用支付网关）
func (w *SubscriptionRenewalWorker) chargeUser(userID uint, amount float64) bool {
	// TODO: 实现实际的支付逻辑
	// 这里模拟 90% 成功率
	return time.Now().Unix()%10 != 0
}
