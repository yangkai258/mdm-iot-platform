package services

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// PetWebSocketMessage 宠物WebSocket消息结构
type PetWebSocketMessage struct {
	Type      string                 `json:"type"`           // message, status_update, action_executed, mood_change
	DeviceID  string                 `json:"device_id"`      // 设备ID
	Content   string                 `json:"content"`        // 消息内容
	Data      map[string]interface{} `json:"data,omitempty"` // 附加数据
	Timestamp string                 `json:"timestamp"`      // 时间戳
}

// PetWebSocketClient 宠物WebSocket客户端
type PetWebSocketClient struct {
	hub      *PetWebSocketHub
	conn     *websocket.Conn
	send     chan []byte
	deviceID string
	userID   string
}

// PetWebSocketHub 宠物WebSocket连接管理Hub
type PetWebSocketHub struct {
	// 设备级别的订阅管理
	deviceClients map[string]map[*PetWebSocketClient]bool
	// 用户级别的订阅管理
	userClients map[string]map[*PetWebSocketClient]bool
	// 全局广播
	broadcast  chan *PetWebSocketMessage
	register   chan *PetWebSocketClient
	unregister chan *PetWebSocketClient
	mu         sync.RWMutex
}

var petHub *PetWebSocketHub
var petHubOnce sync.Once

// GetPetHub 获取宠物WebSocket Hub实例
func GetPetHub() *PetWebSocketHub {
	petHubOnce.Do(func() {
		petHub = &PetWebSocketHub{
			deviceClients: make(map[string]map[*PetWebSocketClient]bool),
			userClients:   make(map[string]map[*PetWebSocketClient]bool),
			broadcast:     make(chan *PetWebSocketMessage, 256),
			register:      make(chan *PetWebSocketClient),
			unregister:    make(chan *PetWebSocketClient),
		}
		go petHub.run()
	})
	return petHub
}

func (h *PetWebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// 按设备注册
			if _, ok := h.deviceClients[client.deviceID]; !ok {
				h.deviceClients[client.deviceID] = make(map[*PetWebSocketClient]bool)
			}
			h.deviceClients[client.deviceID][client] = true
			// 按用户注册
			if _, ok := h.userClients[client.userID]; !ok {
				h.userClients[client.userID] = make(map[*PetWebSocketClient]bool)
			}
			h.userClients[client.userID][client] = true
			h.mu.Unlock()
			log.Printf("[PetWS] Client registered for device=%s, user=%s", client.deviceID, client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			// 从设备订阅中移除
			if clients, ok := h.deviceClients[client.deviceID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
				}
				if len(clients) == 0 {
					delete(h.deviceClients, client.deviceID)
				}
			}
			// 从用户订阅中移除
			if clients, ok := h.userClients[client.userID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.userClients, client.userID)
				}
			}
			h.mu.Unlock()
			log.Printf("[PetWS] Client unregistered for device=%s", client.deviceID)

		case message := <-h.broadcast:
			h.mu.RLock()
			data, _ := json.Marshal(message)
			// 按设备推送
			if message.DeviceID != "" {
				if clients, ok := h.deviceClients[message.DeviceID]; ok {
					for client := range clients {
						select {
						case client.send <- data:
						default:
							h.mu.RUnlock()
							h.unregister <- client
							h.mu.RLock()
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register 注册客户端
func (c *PetWebSocketClient) Register() {
	GetPetHub().register <- c
}

// Unregister 注销客户端
func (c *PetWebSocketClient) Unregister() {
	GetPetHub().unregister <- c
}

// SendToDevice 向指定设备的所有连接发送消息
func SendPetMessageToDevice(deviceID string, msgType string, content string, data map[string]interface{}) {
	message := &PetWebSocketMessage{
		Type:      msgType,
		DeviceID:  deviceID,
		Content:   content,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	GetPetHub().broadcast <- message
}

// SendToUser 向指定用户的所有连接发送消息
func SendPetMessageToUser(userID string, msgType string, deviceID string, content string, data map[string]interface{}) {
	message := &PetWebSocketMessage{
		Type:      msgType,
		DeviceID:  deviceID,
		Content:   content,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	GetPetHub().broadcast <- message
}

// ReadPump 处理读消息
func (c *PetWebSocketClient) ReadPump() {
	defer func() {
		c.Unregister()
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[PetWS] Read error: %v", err)
			}
			break
		}

		// 解析客户端消息
		var msg PetWebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[PetWS] Invalid message format: %v", err)
			continue
		}

		// 处理客户端消息
		c.handleMessage(&msg)
	}
}

// WritePump 处理写消息
func (c *PetWebSocketClient) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (c *PetWebSocketClient) handleMessage(msg *PetWebSocketMessage) {
	switch msg.Type {
	case "ping":
		// 心跳响应
		response := &PetWebSocketMessage{
			Type:      "pong",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		data, _ := json.Marshal(response)
		c.send <- data
	case "subscribe":
		// 订阅设备更新
		log.Printf("[PetWS] Client subscribed to device: %s", msg.DeviceID)
	case "unsubscribe":
		// 取消订阅
		log.Printf("[PetWS] Client unsubscribed from device: %s", msg.DeviceID)
	}
}

// PetWebSocketHandler WebSocket处理函数
type PetWebSocketHandler struct{}

func NewPetWebSocketHandler() *PetWebSocketHandler {
	return &PetWebSocketHandler{}
}

// HandleWebSocket 处理WebSocket连接
func (h *PetWebSocketHandler) HandleWebSocket(c *websocket.Conn, deviceID, userID string) {
	client := &PetWebSocketClient{
		hub:      GetPetHub(),
		conn:     c,
		send:     make(chan []byte, 256),
		deviceID: deviceID,
		userID:   userID,
	}
	client.Register()

	go client.WritePump()
	client.ReadPump()
}

// BroadcastPetStatusUpdate 广播宠物状态更新
func BroadcastPetStatusUpdate(deviceID string, mood, energy, hunger int, expression string) {
	data := map[string]interface{}{
		"mood":       mood,
		"energy":     energy,
		"hunger":     hunger,
		"expression": expression,
	}
	SendPetMessageToDevice(deviceID, "status_update", "", data)
}

// BroadcastNewMessage 广播新消息
func BroadcastNewMessage(deviceID, messageID, content string, senderType int) {
	data := map[string]interface{}{
		"message_id":  messageID,
		"content":     content,
		"sender_type": senderType,
	}
	SendPetMessageToDevice(deviceID, "new_message", content, data)
}

// BroadcastActionExecuted 广播动作执行结果
func BroadcastActionExecuted(deviceID, actionID, result string) {
	data := map[string]interface{}{
		"action_id": actionID,
		"result":    result,
	}
	SendPetMessageToDevice(deviceID, "action_executed", "", data)
}

// BroadcastMoodChange 广播心情变化
func BroadcastMoodChange(deviceID string, mood int, reason string) {
	data := map[string]interface{}{
		"mood":   mood,
		"reason": reason,
	}
	SendPetMessageToDevice(deviceID, "mood_change", "", data)
}
