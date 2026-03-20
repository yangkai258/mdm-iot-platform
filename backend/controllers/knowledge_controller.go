package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Knowledge 知识条目
type Knowledge struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// KnowledgeController 知识库控制器
type KnowledgeController struct{}

// RegisterRoutes 注册路由
func (c *KnowledgeController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/knowledge", c.ListKnowledge)
	r.POST("/knowledge", c.CreateKnowledge)
	r.PUT("/knowledge/:id", c.UpdateKnowledge)
	r.DELETE("/knowledge/:id", c.DeleteKnowledge)
	r.GET("/knowledge/query", c.QueryKnowledge)
	r.GET("/knowledge/weather", c.GetWeather)
}

func (c *KnowledgeController) ListKnowledge(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": []Knowledge{}})
}

func (c *KnowledgeController) CreateKnowledge(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": Knowledge{}})
}

func (c *KnowledgeController) UpdateKnowledge(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0})
}

func (c *KnowledgeController) DeleteKnowledge(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0})
}

func (c *KnowledgeController) QueryKnowledge(g *gin.Context) {
	q := g.Query("q")
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": []Knowledge{}, "query": q})
}

func (c *KnowledgeController) GetWeather(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"temp": 25, "weather": "sunny"}})
}
