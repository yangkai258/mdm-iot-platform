package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

// BroadcastManager 广播管理器
type BroadcastManager struct {
	hub *ClientHub
}

var broadcastManager *BroadcastManager
var bmOnce sync.Once

// GetBroadcastManager 获取广播管理器单例
func GetBroadcastManager() *BroadcastManager {
	bmOnce.Do(func() {
		broadcastManager = &BroadcastManager{
			hub: GetHub(),
		}
	})
	return broadcastManager
}

// BroadcastMessage 广播普通消息
func (bm *BroadcastManager) BroadcastMessage(msgType, title, content string, data interface{}) {
	notification := &Notification{
		Type:      NotificationType(msgType),
		Title:     title,
		Content:   content,
		Data:      data,
		Timestamp: time.Now(),
	}
	bm.hub.Broadcast(notification)
}

// BroadcastAlert 广播告警消息
func (bm *BroadcastManager) BroadcastAlert(title, content string, data interface{}) {
	bm.BroadcastMessage(string(NotificationTypeAlert), title, content, data)
}

// BroadcastOTA 广播OTA升级消息
func (bm *BroadcastManager) BroadcastOTA(title, content string, data interface{}) {
	bm.BroadcastMessage(string(NotificationTypeOTA), title, content, data)
}

// BroadcastSystem 广播系统消息
func (bm *BroadcastManager) BroadcastSystem(title, content string, data interface{}) {
	bm.BroadcastMessage(string(NotificationTypeSystem), title, content, data)
}

// MarshalNotification 将通知转换为JSON
func MarshalNotification(n *Notification) ([]byte, error) {
	return json.Marshal(n)
}
