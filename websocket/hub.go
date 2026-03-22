package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// NotificationType 通知类型
type NotificationType string

const (
	NotificationTypeAlert    NotificationType = "alert"     // 新告警
	NotificationTypeOTA     NotificationType = "ota"       // OTA升级进度
	NotificationTypeSystem  NotificationType = "system"   // 系统消息
)

// Notification 通知消息结构
type Notification struct {
	Type      NotificationType     `json:"type"`                // 通知类型
	Title     string               `json:"title"`               // 通知标题
	Content   string               `json:"content"`             // 通知内容
	Data      interface{}          `json:"data,omitempty"`       // 附加数据
	Timestamp time.Time            `json:"timestamp"`            // 时间戳
}

// ClientHub 所有客户端的集中管理
type ClientHub struct {
	clients    map[*Client]bool
	broadcast  chan *Notification
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var hub *ClientHub
var hubOnce sync.Once

// GetHub 获取全局Hub实例
func GetHub() *ClientHub {
	hubOnce.Do(func() {
		hub = &ClientHub{
			clients:    make(map[*Client]bool),
			broadcast:  make(chan *Notification, 256),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
		go hub.run()
	})
	return hub
}

func (h *ClientHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[WS] Client registered, total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("[WS] Client unregistered, total: %d", len(h.clients))

		case notification := <-h.broadcast:
			h.mu.RLock()
			msg, _ := json.Marshal(notification)
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register 注册客户端
func (h *ClientHub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端
func (h *ClientHub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast 广播通知到所有客户端
func (h *ClientHub) Broadcast(notification *Notification) {
	h.broadcast <- notification
}

// BroadcastToAll 广播通知（便捷方法）
func BroadcastToAll(notification *Notification) {
	GetHub().Broadcast(notification)
}

// SendAlertNotification 发送告警通知
func SendAlertNotification(title, content string, data interface{}) {
	BroadcastToAll(&Notification{
		Type:      NotificationTypeAlert,
		Title:     title,
		Content:   content,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// SendOTANotification 发送OTA升级进度通知
func SendOTANotification(deviceID, status string, progress int, data interface{}) {
	BroadcastToAll(&Notification{
		Type:    NotificationTypeOTA,
		Title:   "OTA升级进度",
		Content: "设备 " + deviceID + " OTA升级" + status,
		Data: map[string]interface{}{
			"device_id": deviceID,
			"status":    status,
			"progress":  progress,
			"extra":     data,
		},
		Timestamp: time.Now(),
	})
}

// SendSystemNotification 发送系统消息通知
func SendSystemNotification(title, content string, data interface{}) {
	BroadcastToAll(&Notification{
		Type:      NotificationTypeSystem,
		Title:     title,
		Content:   content,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// ClientCount 获取当前客户端数量
func (h *ClientHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
