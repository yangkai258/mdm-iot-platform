package notification

import (
	"context"
	"log"
	"time"
)

// NotificationJob 通知任务
type NotificationJob struct {
	ID          uint                  `json:"id"`
	ChannelType string                `json:"channel_type"` // email/sms/webhook
	Recipient   string                `json:"recipient"`
	Subject     string                `json:"subject"`
	Content     string                `json:"content"`
	Payload     []byte                `json:"payload,omitempty"`
	Secret      string                `json:"secret,omitempty"`
	RetryCount  int                   `json:"retry_count"`
	MaxRetries  int                   `json:"max_retries"`
	ChannelID   uint                  `json:"channel_id"`
	AlertID     uint                  `json:"alert_id"`
	Status      string                `json:"status"` // pending/success/failed
	ErrorMsg    string               `json:"error_msg"`
	CreatedAt   time.Time            `json:"created_at"`
}

// NotificationHandler 处理通知的函数类型
type NotificationHandler func(job *NotificationJob) error

// RetryWorker 重试机制 worker
type RetryWorker struct {
	queue      chan *NotificationJob
	handler    NotificationHandler
	maxRetries int
	backoffs   []time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewRetryWorker 创建重试 worker
func NewRetryWorker(handler NotificationHandler, queueSize int) *RetryWorker {
	ctx, cancel := context.WithCancel(context.Background())
	worker := &RetryWorker{
		queue:      make(chan *NotificationJob, queueSize),
		handler:    handler,
		maxRetries: 5,
		// 指数退避: 1s, 2s, 4s, 8s, 16s
		backoffs: []time.Duration{
			1 * time.Second,
			2 * time.Second,
			4 * time.Second,
			8 * time.Second,
			16 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
	return worker
}

// Start 启动 worker
func (w *RetryWorker) Start() {
	log.Printf("[RetryWorker] 启动重试 worker，队列容量: %d, 最大重试次数: %d", cap(w.queue), w.maxRetries)
	go w.run()
}

// Stop 停止 worker
func (w *RetryWorker) Stop() {
	w.cancel()
	close(w.queue)
	log.Printf("[RetryWorker] 已停止")
}

// Enqueue 添加任务到队列
func (w *RetryWorker) Enqueue(job *NotificationJob) {
	if job.MaxRetries == 0 {
		job.MaxRetries = w.maxRetries
	}
	select {
	case w.queue <- job:
		log.Printf("[RetryWorker] 任务入队: JobID=%d, Type=%s, Recipient=%s",
			job.ID, job.ChannelType, job.Recipient)
	default:
		log.Printf("[RetryWorker] 队列已满，任务被丢弃: JobID=%d", job.ID)
	}
}

func (w *RetryWorker) run() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case job := <-w.queue:
			w.processJob(job)
		}
	}
}

func (w *RetryWorker) processJob(job *NotificationJob) {
	err := w.handler(job)
	if err == nil {
		job.Status = "success"
		job.ErrorMsg = ""
		log.Printf("[RetryWorker] 通知发送成功: JobID=%d, Type=%s", job.ID, job.ChannelType)
		return
	}

	job.RetryCount++
	job.ErrorMsg = err.Error()

	if job.RetryCount > job.MaxRetries {
		job.Status = "failed"
		log.Printf("[RetryWorker] 通知发送失败（已达最大重试次数）: JobID=%d, Type=%s, Error=%v",
			job.ID, job.ChannelType, err)
		return
	}

	// 计算退避时间
	backoffIdx := job.RetryCount - 1
	if backoffIdx >= len(w.backoffs) {
		backoffIdx = len(w.backoffs) - 1
	}
	backoff := w.backoffs[backoffIdx]

	job.Status = "pending"
	log.Printf("[RetryWorker] 通知发送失败，准备重试 (#%d/%d, %.0fs后): JobID=%d, Type=%s, Error=%v",
		job.RetryCount, job.MaxRetries, backoff.Seconds(), job.ID, job.ChannelType, err)

	// 延迟后重新入队
	go func(j *NotificationJob, delay time.Duration) {
		time.Sleep(delay)
		w.Enqueue(j)
	}(job, backoff)
}

// Stats 获取 worker 统计信息
func (w *RetryWorker) Stats() map[string]interface{} {
	return map[string]interface{}{
		"queue_size":    len(w.queue),
		"max_retries":   w.maxRetries,
		"backoffs":      w.backoffs,
		"running":       w.ctx.Err() == nil,
	}
}

// DefaultRetryWorker 全局默认 worker 实例
var DefaultRetryWorker *RetryWorker

// InitDefaultRetryWorker 初始化全局默认 worker
func InitDefaultRetryWorker(handler NotificationHandler) {
	DefaultRetryWorker = NewRetryWorker(handler, 1000)
	DefaultRetryWorker.Start()
	log.Printf("[RetryWorker] 全局默认 worker 已初始化")
}

// EnqueueNotificationJob 全局入队方法
func EnqueueNotificationJob(job *NotificationJob) {
	if DefaultRetryWorker != nil {
		DefaultRetryWorker.Enqueue(job)
	} else {
		log.Printf("[RetryWorker] 全局 worker 未初始化，直接处理: %v", job)
	}
}
