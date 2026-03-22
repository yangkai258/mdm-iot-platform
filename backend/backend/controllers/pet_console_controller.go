package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// PetConversation 对话
type PetConversation struct {
	ID        uint   `json:"id"`
	DeviceID  string `json:"device_id"`
	Title     string `json:"title"`
	Status    string `json:"status"` // active/closed
}

// PetMessage 消息
type PetMessage struct {
	ID             uint   `json:"id"`
	ConversationID uint   `json:"conversation_id"`
	Role           string `json:"role"` // user/assistant
	Content        string `json:"content"`
}

// PetConsoleController 宠物控制台
type PetConsoleController struct{}

// RegisterRoutes 注册路由
func (c *PetConsoleController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/pet/conversations", c.ListConversations)
	r.POST("/pet/conversations", c.CreateConversation)
	r.GET("/pet/conversations/:id", c.GetConversation)
	r.POST("/pet/messages", c.SendMessage)
	r.GET("/pet/status/:device_id", c.GetPetStatus)
	r.PUT("/pet/status/:device_id", c.UpdatePetStatus)
}

func (c *PetConsoleController) ListConversations(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": []PetConversation{}})
}

func (c *PetConsoleController) CreateConversation(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": PetConversation{}})
}

func (c *PetConsoleController) GetConversation(g *gin.Context) {
	id := g.Param("id")
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": PetConversation{ID: 1, DeviceID: id}})
}

func (c *PetConsoleController) SendMessage(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": PetMessage{}})
}

func (c *PetConsoleController) GetPetStatus(g *gin.Context) {
	deviceID := g.Param("device_id")
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"device_id": deviceID, "status": "happy"}})
}

func (c *PetConsoleController) UpdatePetStatus(g *gin.Context) {
	deviceID := g.Param("device_id")
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"device_id": deviceID}})
}
