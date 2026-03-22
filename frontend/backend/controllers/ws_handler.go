package controllers

import (
	"net/http"

	ws "mdm-backend/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	// 从上下文获取用户ID（如果已登录）
	userID := ""
	if uid, exists := ctx.Get("user_id"); exists {
		userID = uid.(string)
	}

	hub := ws.GetHub()
	client := ws.NewClient(hub, conn, userID)
	hub.Register(client)

	// 启动读写goroutine
	go client.WritePump()
	go client.ReadPump()
}
