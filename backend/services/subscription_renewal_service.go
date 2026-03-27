package services

import (
	"fmt"
	"log"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

const (
	// RenewalRetryMax 续费失败最大重试次数
	RenewalRetryMax = 3
	// ReminderDaysBefore 到期前几天发送提醒
	ReminderDaysBefore = 3
)

// SubscriptionRenewalService 订阅续费服务
type SubscriptionRenewalService struct {
	DB *gorm.DB
}

// NewSubscriptionRenewalService 创建续费服务
func NewSubscriptionRenewalService(db *gorm.DB) *SubscriptionRenewalService {
	return &SubscriptionRenewalService{DB: db}
}

// StartRenewalJob 启动续费定时任务
func (s *SubscriptionRenewalService) StartRenewalJob() {
	log.Println("[SubscriptionRenewal] 续费服务已启动，每小时检查一次到期订阅")

	// 立即执行一次
	go s.processExpiredAndExpiringSubscriptions()

	// 每小时执行一次
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		go s.processExpiredAndExpiringSubscriptions()
	}
}

// processExpiredAndExpiringSubscriptions 处理到期和即将到期的订阅
func (s *SubscriptionRenewalService) processExpiredAndExpiringSubscriptions() {
	now := time.Now()

	// 1. 处理已过期且开启自动续费的订阅
	var expiredSubs []models.Subscription
	s.DB.Where("status = ? AND auto_renew = ? AND end_date < ?", "active", true, now).
		Find(&expiredSubs)

	for _, sub := range expiredSubs {
		s.processRenewal(sub)
	}

	// 2. 发送即将到期提醒（到期前3天）
	var expiringSubs []models.Subscription
	reminderTime := now.AddDate(0, 0, ReminderDaysBefore)
	s.DB.Where("status = ? AND auto_renew = ? AND end_date BETWEEN ? AND ? AND (reminder_sent_at IS NULL OR reminder_sent_at < ?)",
		"active", true, now, reminderTime, now).
		Find(&expiringSubs)

	for _, sub := range expiringSubs {
		s.sendRenewalReminder(sub)
	}

	log.Printf("[SubscriptionRenewal] 处理完成: 过期续费 %d 条, 提醒 %d 条", len(expiredSubs), len(expiringSubs))
}

// processRenewal 处理单个订阅续费
func (s *SubscriptionRenewalService) processRenewal(sub models.Subscription) {
	log.Printf("[SubscriptionRenewal] 开始处理订阅续费: SubscriptionID=%d, UserID=%d, RetryCount=%d", sub.ID, sub.UserID, sub.RetryCount)

	// 尝试扣费续期
	success, err := s.attemptRenewal(sub)

	if success {
		// 续费成功
		s.handleRenewalSuccess(sub)
	} else {
		// 续费失败
		s.handleRenewalFailure(sub, err)
	}
}

// attemptRenewal 尝试扣费续期（模拟支付网关调用）
func (s *SubscriptionRenewalService) attemptRenewal(sub models.Subscription) (bool, error) {
	// TODO: 集成真实支付网关（如微信支付、支付宝）
	// 这里模拟支付逻辑：95%成功率用于测试
	// 实际项目中调用 paymentGateway.Charge(sub.UserID, sub.Price)

	// 模拟支付延迟
	time.Sleep(100 * time.Millisecond)

	// 模拟支付结果（实际项目中由支付网关返回）
	// 这里随机模拟成功/失败，用于测试重试逻辑
	// return true, nil // 成功
	return false, fmt.Errorf("支付网关连接超时") // 测试重试逻辑时使用
}

// ProcessPayment 处理支付（供手动续费调用）
func (s *SubscriptionRenewalService) ProcessPayment(sub models.Subscription, renewal models.SubscriptionRenewal) (bool, error) {
	// TODO: 集成真实支付网关（如微信支付、支付宝）
	// 实际项目中调用 paymentGateway.Charge(sub.UserID, renewal.Amount, renewal.PaymentMethod)

	// 模拟支付延迟
	time.Sleep(100 * time.Millisecond)

	// 模拟支付结果：90%成功率
	// 实际项目中根据支付网关返回结果判断
	if time.Now().Unix()%10 != 0 {
		return true, nil // 成功
	}
	return false, fmt.Errorf("支付失败: 余额不足或支付方式无效")
}

// handleRenewalSuccess 处理续费成功
func (s *SubscriptionRenewalService) handleRenewalSuccess(sub models.Subscription) {
	now := time.Now()
	newEndDate := sub.EndDate.AddDate(0, 0, sub.Duration)

	// 更新订阅状态
	updates := map[string]interface{}{
		"end_date":        newEndDate,
		"status":          "active",
		"last_renew_at":   now,
		"renew_count":     sub.RenewCount + 1,
		"retry_count":     0,          // 重置重试次数
		"renew_fail_reason": "",       // 清除失败原因
	}
	s.DB.Model(&sub).Updates(updates)

	// 记录续费日志
	s.createRenewalLog(sub.ID, "renewal_success", sub.Price, "success", "", 0)

	log.Printf("[SubscriptionRenewal] 续费成功: SubscriptionID=%d, NewEndDate=%s", sub.ID, newEndDate.Format("2006-01-02"))

	// 发送续费成功通知
	s.sendRenewalSuccessNotification(sub, newEndDate)
}

// handleRenewalFailure 处理续费失败
func (s *SubscriptionRenewalService) handleRenewalFailure(sub models.Subscription, err error) {
	newRetryCount := sub.RenewCount + 1

	log.Printf("[SubscriptionRenewal] 续费失败: SubscriptionID=%d, RetryCount=%d, Error=%v", sub.ID, newRetryCount, err)

	if newRetryCount >= RenewalRetryMax {
		// 重试次数超限，暂停服务
		s.suspendSubscription(sub, err.Error())
	} else {
		// 更新重试次数
		s.DB.Model(&sub).Updates(map[string]interface{}{
			"retry_count":       newRetryCount,
			"renew_fail_reason": err.Error(),
		})

		// 记录失败日志
		s.createRenewalLog(sub.ID, "renewal_failed", sub.Price, "failed", err.Error(), newRetryCount)

		// 发送续费失败通知
		s.sendRenewalFailedNotification(sub, newRetryCount)
	}
}

// suspendSubscription 暂停订阅服务
func (s *SubscriptionRenewalService) suspendSubscription(sub models.Subscription, reason string) {
	now := time.Now()

	s.DB.Model(&sub).Updates(map[string]interface{}{
		"status":        "suspended",
		"suspended_at": now,
	})

	// 记录暂停日志
	s.createRenewalLog(sub.ID, "suspended", 0, "success", reason, RenewalRetryMax)

	log.Printf("[SubscriptionRenewal] 服务已暂停: SubscriptionID=%d, Reason=%s", sub.ID, reason)

	// 发送服务暂停通知
	s.sendSubscriptionSuspendedNotification(sub)
}

// sendRenewalReminder 发送续费提醒
func (s *SubscriptionRenewalService) sendRenewalReminder(sub models.Subscription) {
	daysUntilExpiry := int(time.Until(sub.EndDate).Hours() / 24)

	// 创建站内通知
	notification := &models.Notification{
		DeviceID:  fmt.Sprintf("user_%d", sub.UserID),
		Title:     "订阅即将到期提醒",
		Content:   fmt.Sprintf("您的%s订阅将于%d天后到期，请及时续费以避免服务中断。", sub.PlanName, daysUntilExpiry),
		Priority:  2, // 重要
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}

	if err := s.DB.Create(notification).Error; err != nil {
		log.Printf("[SubscriptionRenewal] 提醒通知创建失败: SubscriptionID=%d, Error=%v", sub.ID, err)
		return
	}

	// 更新提醒发送时间
	now := time.Now()
	s.DB.Model(&sub).Update("reminder_sent_at", now)

	// 记录提醒日志
	s.createRenewalLog(sub.ID, "reminder", 0, "success", "", 0)

	log.Printf("[SubscriptionRenewal] 续费提醒已发送: SubscriptionID=%d, DaysLeft=%d", sub.ID, daysUntilExpiry)
}

// sendRenewalSuccessNotification 发送续费成功通知
func (s *SubscriptionRenewalService) sendRenewalSuccessNotification(sub models.Subscription, newEndDate time.Time) {
	notification := &models.Notification{
		DeviceID:  fmt.Sprintf("user_%d", sub.UserID),
		Title:     "订阅续费成功",
		Content:   fmt.Sprintf("您的%s订阅已成功续费，有效期至%s。", sub.PlanName, newEndDate.Format("2006-01-02")),
		Priority:  1,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}
	s.DB.Create(notification)
}

// sendRenewalFailedNotification 发送续费失败通知
func (s *SubscriptionRenewalService) sendRenewalFailedNotification(sub models.Subscription, retryCount int) {
	notification := &models.Notification{
		DeviceID:  fmt.Sprintf("user_%d", sub.UserID),
		Title:     "订阅续费失败",
		Content:   fmt.Sprintf("您的%s订阅续费失败（第%d次），即将进行第%d次重试。如续费仍失败，服务将被暂停。请及时检查支付方式。",
			sub.PlanName, retryCount, retryCount+1),
		Priority:  2,
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}
	s.DB.Create(notification)
}

// sendSubscriptionSuspendedNotification 发送服务暂停通知
func (s *SubscriptionRenewalService) sendSubscriptionSuspendedNotification(sub models.Subscription) {
	notification := &models.Notification{
		DeviceID:  fmt.Sprintf("user_%d", sub.UserID),
		Title:     "订阅服务已暂停",
		Content:   fmt.Sprintf("您的%s订阅因续费失败已暂停服务。请联系客服或前往个人中心重新开通服务。", sub.PlanName),
		Priority:  3, // 紧急
		Channel:   "push",
		Status:    "sent",
		CreatedBy: "system",
	}
	s.DB.Create(notification)
}

// createRenewalLog 创建续费日志
func (s *SubscriptionRenewalService) createRenewalLog(subscriptionID uint, action string, amount float64, status, failReason string, retryCount int) {
	log := &models.SubscriptionRenewalLog{
		SubscriptionID: subscriptionID,
		Action:         action,
		Amount:         amount,
		Status:         status,
		FailReason:     failReason,
		RetryCount:     retryCount,
	}
	s.DB.Create(log)
}

// ManualRenew 手动触发续费（API调用）
func (s *SubscriptionRenewalService) ManualRenew(subID uint) error {
	var sub models.Subscription
	if err := s.DB.First(&sub, subID).Error; err != nil {
		return fmt.Errorf("订阅不存在: %v", err)
	}

	if sub.Status == "suspended" {
		// 暂停状态的订阅也可以手动续费恢复
		s.processRenewal(sub)
		return nil
	}

	if sub.Status == "cancelled" {
		return fmt.Errorf("已取消的订阅无法续费")
	}

	s.processRenewal(sub)
	return nil
}

// GetRenewalLogs 获取续费日志
func (s *SubscriptionRenewalService) GetRenewalLogs(subID uint) ([]models.SubscriptionRenewalLog, error) {
	var logs []models.SubscriptionRenewalLog
	err := s.DB.Where("subscription_id = ?", subID).Order("created_at DESC").Limit(50).Find(&logs).Error
	return logs, err
}

// ResumeSubscription 恢复已暂停的订阅（手动处理支付问题后）
func (s *SubscriptionRenewalService) ResumeSubscription(subID uint) error {
	var sub models.Subscription
	if err := s.DB.First(&sub, subID).Error; err != nil {
		return fmt.Errorf("订阅不存在: %v", err)
	}

	if sub.Status != "suspended" {
		return fmt.Errorf("只有暂停状态的订阅才能恢复")
	}

	// 重置状态
	updates := map[string]interface{}{
		"status":              "active",
		"suspended_at":        nil,
		"retry_count":         0,
		"renew_fail_reason":   "",
	}
	s.DB.Model(&sub).Updates(updates)

	// 记录恢复日志
	s.createRenewalLog(subID, "resumed", 0, "success", "手动恢复", 0)

	log.Printf("[SubscriptionRenewal] 订阅已恢复: SubscriptionID=%d", subID)
	return nil
}
