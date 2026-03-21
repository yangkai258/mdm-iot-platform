package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Knowledge 知识条目
type Knowledge struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// KnowledgeController 知识库控制器
type KnowledgeController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册路由（不使用DB，路由在 main.go 中单独注册）
func (c *KnowledgeController) RegisterRoutes(rg *gin.RouterGroup) {
}

// RegisterRoutesWithDB 注册路由（传入DB）
func (c *KnowledgeController) RegisterRoutesWithDB(rg *gin.RouterGroup, db *gorm.DB) {
	c.DB = db
	rg.GET("/knowledge", c.ListKnowledge)
	rg.POST("/knowledge", c.CreateKnowledge)
	rg.PUT("/knowledge/:id", c.UpdateKnowledge)
	rg.DELETE("/knowledge/:id", c.DeleteKnowledge)
	rg.GET("/knowledge/query", c.QueryKnowledge)
	rg.POST("/knowledge/batch-delete", c.BatchDeleteKnowledge)
}

// InitDB 初始化数据库连接（由 main.go 注入）
func (c *KnowledgeController) InitDB(db *gorm.DB) {
	c.DB = db
}

// ListKnowledge 获取知识库列表
func (c *KnowledgeController) ListKnowledge(g *gin.Context) {
	if c.DB == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	var knowledgeList []models.Knowledge
	query := c.DB.Model(&models.Knowledge{})

	// 关键词搜索
	if keyword := g.Query("keyword"); keyword != "" {
		query = query.Where("question ILIKE ? OR answer ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	// 分类筛选
	if category := g.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(g.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&knowledgeList)

	g.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":  knowledgeList,
			"total": total,
			"page":  page,
		},
	})
}

// CreateKnowledge 创建知识条目
func (c *KnowledgeController) CreateKnowledge(g *gin.Context) {
	if c.DB == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	var req models.Knowledge
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if err := c.DB.Create(&req).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	g.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": req})
}

// UpdateKnowledge 更新知识条目
func (c *KnowledgeController) UpdateKnowledge(g *gin.Context) {
	if c.DB == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	id := g.Param("id")
	var req models.Knowledge
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var existing models.Knowledge
	if err := c.DB.First(&existing, id).Error; err != nil {
		g.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "知识不存在"})
		return
	}

	existing.Category = req.Category
	existing.Question = req.Question
	existing.Answer = req.Answer

	if err := c.DB.Save(&existing).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	g.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": existing})
}

// DeleteKnowledge 删除知识条目
func (c *KnowledgeController) DeleteKnowledge(g *gin.Context) {
	if c.DB == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	id := g.Param("id")
	if err := c.DB.Delete(&models.Knowledge{}, id).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	g.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// QueryKnowledge 搜索知识
func (c *KnowledgeController) QueryKnowledge(g *gin.Context) {
	if c.DB == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	q := g.Query("q")
	var results []models.Knowledge
	query := c.DB.Model(&models.Knowledge{})

	if q != "" {
		query = query.Where("question ILIKE ? OR answer ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	query.Order("id DESC").Limit(20).Find(&results)

	g.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list":  results,
		"query": q,
	}})
}

// GetWeather 获取天气（占位接口）
func (c *KnowledgeController) GetWeather(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"temp": 25, "weather": "sunny"}})
}

// BatchDeleteKnowledge 批量删除知识
func (c *KnowledgeController) BatchDeleteKnowledge(g *gin.Context) {
	if c.DB == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := g.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		g.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供要删除的ID列表"})
		return
	}

	if err := c.DB.Delete(&models.Knowledge{}, req.IDs).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	g.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}
