package services

import (
	"log"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// RenewalScheduler 续费定时任务调度器
type RenewalScheduler struct {
	DB              *gorm.DB
	RenewalService  *SubscriptionRenewalService
	stopChan        chan struct{}
	isRunning       bool
}

// NewRenewalScheduler 创建续费调度器
func NewRenewalScheduler(db *gorm.DB) *RenewalScheduler {
	return &RenewalScheduler{
		DB:             db,
		RenewalService: NewSubscriptionRenewalService(db),
		stopChan:       make(chan struct{}),
		isRunning:      false,
	}
}

// Start 启动调度器
func (s *RenewalScheduler) Start() {
	if s.isRunning {
		log.Println("[RenewalScheduler] 调度器已在运行中")
		return
	}
	s.isRunning = true
	log.Println("[RenewalScheduler] 续费调度器已启动")

	// 立即执行一次
	go s.runDailyCheck()

	// 每小时执行一次过期检查
	hourlyTicker := time.NewTicker(1 * time.Hour)
	defer hourlyTicker.Stop()

	// 每天凌晨2点执行全面的续费检查
	dailyCheck := time.NewTicker(24 * time.Hour)
	defer dailyCheck.Stop()

	for {
		select {
		case <-hourlyTicker.C:
			go s.processExpiredSubscriptions()
			go s.processExpiringReminders()
		case <-dailyCheck.C:
			go s.runDailyCheck()
		case <-s.stopChan:
			s.isRunning = false
			log.Println("[RenewalScheduler] 续费调度器已停止")
			return
		}
	}
}

// Stop 停止调度器
func (s *RenewalScheduler) Stop() {
	close(s.stopChan)
	s.isRunning = false
}

// runDailyCheck 每天执行一次的全面检查
func (s *RenewalScheduler) runDailyCheck() {
	log.Println("[RenewalScheduler] 执行每日续费检查")

	now := time.Now()

	// 1. 处理已过期且开启自动续费的订阅
	s.processExpiredSubscriptions()

	// 2. 发送即将到期提醒（到期前7天、3天、1天）
	s.processExpiringReminders()

	// 3. 处理暂停的订阅（如果用户已解决支付问题）
	s.processSuspendedSubscriptions()

	// 4. 更新统计数据
	s.updateRenewalStats()

	log.Printf("[RenewalScheduler] 每日检查完成: %s", now.Format("2006-01-02 15:04:05"))
}

// processExpiredSubscriptions 处理到期订阅
func (s *RenewalScheduler) processExpiredSubscriptions() {
	now := time.Now()

	var subs []models.Subscription
	s.DB.Where("status = ? AND auto_renew = ? AND end_date < ?", "active", true, now).
		Find(&subs)

	for _, sub := range subs {
		// 获取自动续费设置
		var setting models.AutoRenewalSetting
		if err := s.DB.Where("subscription_id = ? AND enabled = ?", sub.ID, true).First(&setting).Error; err != nil {
			// 没有启用自动续费设置，跳过
			continue
		}

		// 计算续费后的到期日期
		newExpiredDate := sub.EndDate.AddDate(0, 0, sub.Duration)

		// 创建续费记录
		renewal := models.SubscriptionRenewal{
			SubscriptionID:  sub.ID,
			UserID:          sub.UserID,
			PaymentMethodID: setting.PaymentMethodID,
			Amount:          sub.Price,
			PaymentStatus:   "processing",
			RenewalDate:      now,
			ExpiredDate:     newExpiredDate,
		}

		if setting.PaymentMethodID != nil {
			var pm models.PaymentMethod
			if err := s.DB.First(&pm, *setting.PaymentMethodID).Error; err == nil {
				renewal.PaymentMethod = pm.MethodType
			}
		}

		// 尝试扣费
		success, err := s.RenewalService.ProcessPayment(sub, renewal)

		if success {
			// 扣费成功
			renewal.PaymentStatus = "success"
			s.DB.Create(&renewal)

			s.DB.Model(&sub).Updates(map[string]interface{}{
				"end_date":       newExpiredDate,
				"status":         "active",
				"last_renew_at":  &now,
				"renew_count":    sub.RenewCount + 1,
				"retry_count":   0,
				"renew_fail_reason": "",
			})

			// 发送成功通知
			s.sendRenewalSuccessNotification(sub, newExpiredDate)
			log.Printf("[RenewalScheduler] 自动续费成功: SubscriptionID=%d, UserID=%d", sub.ID, sub.UserID)
		} else {
			// 扣费失败
			renewal.PaymentStatus = "failed"
			renewal.FailReason = err.Error()
			renewal.RetryCount = 1
			s.DB.Create(&renewal)

			newRetryCount := sub.RetryCount + 1
			nextRetry := now.AddDate(0, 0, 1) // 1天后重试

			s.DB.Model(&sub).Updates(map[string]interface{}{
				"retry_count":       newRetryCount,
				"renew_fail_reason": err.Error(),
			})

			// 更新续费记录的下次重试时间
			s.DB.Model(&renewal).Update("next_retry_at", &nextRetry)

			if newRetryCount >= 3 {
				// 重试次数超限，暂停服务
				s.suspendSubscription(sub, err.Error())
			} else {
				// 发送失败通知
				s.sendRenewalFailedNotification(sub, newRetryCount)
			}

			log.Printf("[RenewalScheduler] 自动续费失败: SubscriptionID=%d, RetryCount=%d, Error=%v", sub.ID, newRetryCount, err)
		}
	}
}

// processExpiringReminders 发送即将到期提醒
func (s *RenewalScheduler) processExpiringReminders() {
	now := time.Now()

	var subs []models.Subscription
	s.DB.Where("status = ? AND auto_renew = ?", "active", true).Find(&subs)

	for _, sub := range subs {
		// 获取自动续费设置
		var setting models.AutoRenewalSetting
		if err := s.DB.Where("subscription_id = ?", sub.ID).First(&setting).Error; err != nil {
			continue
		}

		daysUntilExpiry := int(time.Until(sub.EndDate).Hours() / 24)

		// 解析提醒天数设置
		reminderDaysMap := map[int]bool{}
		if setting.ReminderDays != "" {
			for _, d := range strings.Split(setting.ReminderDays, ",") {
				if day, err := strconv.Atoi(strings.TrimSpace(d)); err == nil {
					reminderDaysMap[day] = true
				}
			}
		} else {
			// 默认提醒天数
			reminderDaysMap = map[int]bool{7: true, 3: true, 1: true}
		}

		// 检查是否需要发送提醒
		if days, ok := reminderDaysMap[daysUntilExpiry]; ok && days {
			// 检查是否已发送过提醒
			if setting.LastRemindAt != nil {
				lastRemindDay := int(sub.EndDate.Sub(*setting.LastRemindAt).Hours() / 24)
				if _, sent := reminderDaysMap[lastRemindDay]; sent && lastRemindDay >= daysUntilExpiry {
					continue // 今天已发送过
				}
			}

			// 发送提醒
			s.sendRenewalReminder(sub, daysUntilExpiry)

			// 更新提醒时间
			s.DB.Model(&setting).Update("last_remind_at", now)
		}
	}
}

// processSuspendedSubscriptions 处理暂停的订阅
func (s *RenewalScheduler) processSuspendedSubscriptions() {
	// 查找暂停超过7天的订阅，发送最终提醒
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var subs []models.Subscription
	s.DB.Where("status = ? AND suspended_at < ?", "suspended", sevenDaysAgo).Find(&subs)

	for _, sub := range subs {
		// 创建最终提醒通知
		notification := &models.Notification{
			DeviceID:  "system",
			Title:     "订阅服务即将终止",
			Content:   "您的订阅已暂停超过7天，如继续不续费，服务将被终止。",
			Priority:  3,
			Channel:   "push",
			Status:    "pending",
			CreatedBy: "system",
		}
		s.DB.Create(notification)
		log.Printf("[RenewalScheduler] 发送服务终止预警: SubscriptionID=%d", sub.ID)
	}
}

// suspendSubscription 暂停订阅
func (s *RenewalScheduler) suspendSubscription(sub models.Subscription, reason string) {
	now := time.Now()

	s.DB.Model(&sub).Updates(map[string]interface{}{
		"status":     "suspended",
		"suspended_at": &now,
	})

	// 创建暂停日志
	renewalLog := models.SubscriptionRenewalLog{
		SubscriptionID: sub.ID,
		Action:         "suspended",
		Amount:         sub.Price,
		Status:         "failed",
		FailReason:     reason,
		RetryCount:     sub.RetryCount,
	}
	s.DB.Create(&renewalLog)

	// 发送暂停通知
	s.sendSubscriptionSuspendedNotification(sub)

	log.Printf("[RenewalScheduler] 订阅已暂停: SubscriptionID=%d, Reason=%s", sub.ID, reason)
}

// sendRenewalReminder 发送续费提醒
func (s *RenewalScheduler) sendRenewalReminder(sub models.Subscription, daysLeft int) {
	notification := &models.Notification{
		DeviceID:  "system",
		Title:     "订阅即将到期提醒",
		Content:   "",
		Priority:  2,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}

	if daysLeft == 0 {
		notification.Title = "订阅今日到期"
		notification.Content = "您的" + sub.PlanName + "订阅将于今天到期，为避免服务中断，请及时续费。"
		notification.Priority = 3
	} else {
		notification.Content = "您的" + sub.PlanName + "订阅将于" + strconv.Itoa(daysLeft) + "天后到期，请及时续费以避免服务中断。"
	}

	s.DB.Create(notification)
	log.Printf("[RenewalScheduler] 发送续费提醒: SubscriptionID=%d, DaysLeft=%d", sub.ID, daysLeft)
}

// sendRenewalSuccessNotification 发送续费成功通知
func (s *RenewalScheduler) sendRenewalSuccessNotification(sub models.Subscription, newEndDate time.Time) {
	notification := &models.Notification{
		DeviceID:  "system",
		Title:     "订阅自动续费成功",
		Content:   "您的" + sub.PlanName + "订阅已成功续费，有效期至" + newEndDate.Format("2006-01-02") + "。",
		Priority:  1,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}
	s.DB.Create(notification)
}

// sendRenewalFailedNotification 发送续费失败通知
func (s *RenewalScheduler) sendRenewalFailedNotification(sub models.Subscription, retryCount int) {
	notification := &models.Notification{
		DeviceID:  "system",
		Title:     "订阅续费失败",
		Content:   "您的" + sub.PlanName + "订阅续费失败（第" + strconv.Itoa(retryCount) + "次），明天将自动重试。如多次失败，服务将被暂停。",
		Priority:  2,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}
	s.DB.Create(notification)
}

// sendSubscriptionSuspendedNotification 发送服务暂停通知
func (s *RenewalScheduler) sendSubscriptionSuspendedNotification(sub models.Subscription) {
	notification := &models.Notification{
		DeviceID:  "system",
		Title:     "订阅服务已暂停",
		Content:   "您的" + sub.PlanName + "订阅因续费失败已暂停服务。请登录平台重新开通或联系客服。",
		Priority:  3,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}
	s.DB.Create(notification)
}

// updateRenewalStats 更新续费统计数据
func (s *RenewalScheduler) updateRenewalStats() {
	// 统计今日/本周/本月续费情况
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var todayCount, weekCount, monthCount int64
	var todayAmount, weekAmount, monthAmount float64

	// 今日
	s.DB.Model(&models.SubscriptionRenewal{}).
		Where("payment_status = ? AND created_at >= ?", "success", today).
		Count(&todayCount)
	s.DB.Model(&models.SubscriptionRenewal{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("payment_status = ? AND created_at >= ?", "success", today).
		Scan(&todayAmount)

	// 本周
	s.DB.Model(&models.SubscriptionRenewal{}).
		Where("payment_status = ? AND created_at >= ?", "success", weekStart).
		Count(&weekCount)
	s.DB.Model(&models.SubscriptionRenewal{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("payment_status = ? AND created_at >= ?", "success", weekStart).
		Scan(&weekAmount)

	// 本月
	s.DB.Model(&models.SubscriptionRenewal{}).
		Where("payment_status = ? AND created_at >= ?", "success", monthStart).
		Count(&monthCount)
	s.DB.Model(&models.SubscriptionRenewal{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("payment_status = ? AND created_at >= ?", "success", monthStart).
		Scan(&monthAmount)

	log.Printf("[RenewalScheduler] 续费统计 - 今日: %d笔/%.2f, 本周: %d笔/%.2f, 本月: %d笔/%.2f",
		todayCount, todayAmount, weekCount, weekAmount, monthCount, monthAmount)
}
