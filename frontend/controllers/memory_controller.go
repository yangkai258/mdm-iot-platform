package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemoryController 记忆库控制器
type MemoryController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册记忆库路由
func (m *MemoryController) RegisterRoutes(r *gin.RouterGroup) {
	// 短期记忆
	r.GET("/pets/:device_id/memory/short", m.ListShortTermMemory)
	r.POST("/pets/:device_id/memory/short", m.CreateShortTermMemory)
	r.DELETE("/pets/:device_id/memory/short/:id", m.DeleteShortTermMemory)
	r.POST("/pets/:device_id/memory/short/cleanup", m.CleanupExpiredMemory)

	// 长期记忆
	r.GET("/pets/:device_id/memory/long", m.ListLongTermMemory)
	r.POST("/pets/:device_id/memory/long", m.CreateLongTermMemory)
	r.PUT("/pets/:device_id/memory/long/:id", m.UpdateLongTermMemory)
	r.DELETE("/pets/:device_id/memory/long/:id", m.DeleteLongTermMemory)
	r.POST("/pets/:device_id/memory/long/:id/reinforce", m.ReinforceMemory)
}

// ListShortTermMemory 获取短期记忆列表
func (m *MemoryController) ListShortTermMemory(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var query models.MemoryListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "查询参数错误"})
		return
	}

	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	userID := int64(1) // TODO: 从认证中间件获取

	var memories []models.PetShortTermMemory
	db := m.DB.Where("device_id = ? AND user_id = ?", deviceID, userID)
	if query.SessionID != "" {
		db = db.Where("session_id = ?", query.SessionID)
	}
	if query.Type != "" {
		db = db.Where("memory_type = ?", query.Type)
	}

	if err := db.Order("created_at DESC").Offset(query.Offset).Limit(limit).Find(&memories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询记忆失败"})
		return
	}

	var responses []*models.ShortTermMemoryResponse
	for i := range memories {
		responses = append(responses, memories[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// CreateShortTermMemory 创建短期记忆
func (m *MemoryController) CreateShortTermMemory(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req models.CreateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证中间件获取

	contentJSON, err := json.Marshal(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容序列化失败"})
		return
	}

	importance := req.Importance
	if importance == 0 {
		importance = 0.5
	}

	memory := models.PetShortTermMemory{
		DeviceID:    deviceID,
		UserID:      userID,
		SessionID:   req.SessionID,
		MessageID:   req.MessageID,
		MemoryType:  req.MemoryType,
		Content:     string(contentJSON),
		Importance:  importance,
		AccessCount: 0,
	}

	// 设置过期时间
	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		memory.ExpiresAt = &expiresAt
	}

	if err := m.DB.Create(&memory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建记忆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": memory.ToResponse(),
	})
}

// DeleteShortTermMemory 删除短期记忆
func (m *MemoryController) DeleteShortTermMemory(c *gin.Context) {
	memoryID := c.Param("id")
	if memoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "记忆ID不能为空"})
		return
	}

	if err := m.DB.Where("memory_id = ?", memoryID).Delete(&models.PetShortTermMemory{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除记忆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// CleanupExpiredMemory 清理过期记忆
func (m *MemoryController) CleanupExpiredMemory(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	now := time.Now()
	result := m.DB.Where("device_id = ? AND expires_at IS NOT NULL AND expires_at < ?", deviceID, now).
		Delete(&models.PetShortTermMemory{})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"deleted_count": result.RowsAffected,
		},
	})
}

// ListLongTermMemory 获取长期记忆列表
func (m *MemoryController) ListLongTermMemory(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var query models.LongTermMemoryListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "查询参数错误"})
		return
	}

	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	userID := int64(1) // TODO: 从认证中间件获取

	var memories []models.PetLongTermMemory
	db := m.DB.Where("device_id = ? AND user_id = ?", deviceID, userID)
	if query.Category != "" {
		db = db.Where("memory_category = ?", query.Category)
	}

	if err := db.Order("created_at DESC").
		Offset(query.Offset).
		Limit(limit).
		Find(&memories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询记忆失败"})
		return
	}

	var responses []*models.LongTermMemoryResponse
	for i := range memories {
		responses = append(responses, memories[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// CreateLongTermMemory 创建长期记忆
func (m *MemoryController) CreateLongTermMemory(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var req models.CreateLongTermMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	userID := int64(1) // TODO: 从认证中间件获取

	contentJSON, err := json.Marshal(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "内容序列化失败"})
		return
	}

	keywordsJSON, _ := json.Marshal(req.Keywords)
	if req.Keywords == nil {
		keywordsJSON = []byte("[]")
	}

	confidence := req.Confidence
	if confidence == 0 {
		confidence = 0.8
	}

	memory := models.PetLongTermMemory{
		DeviceID:       deviceID,
		UserID:         userID,
		MemoryCategory: req.MemoryCategory,
		Content:        string(contentJSON),
		Keywords:       string(keywordsJSON),
		Confidence:     confidence,
		IsLocked:       req.IsLocked,
	}

	if err := m.DB.Create(&memory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建记忆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": memory.ToResponse(),
	})
}

// UpdateLongTermMemory 更新长期记忆
func (m *MemoryController) UpdateLongTermMemory(c *gin.Context) {
	memoryID := c.Param("id")
	if memoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "记忆ID不能为空"})
		return
	}

	var memory models.PetLongTermMemory
	if err := m.DB.Where("memory_id = ?", memoryID).First(&memory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记忆不存在"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	updates["updated_at"] = time.Now()
	if err := m.DB.Model(&memory).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新记忆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": memory.ToResponse(),
	})
}

// DeleteLongTermMemory 删除长期记忆
func (m *MemoryController) DeleteLongTermMemory(c *gin.Context) {
	memoryID := c.Param("id")
	if memoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "记忆ID不能为空"})
		return
	}

	if err := m.DB.Where("memory_id = ?", memoryID).Delete(&models.PetLongTermMemory{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除记忆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ReinforceMemory 强化记忆
func (m *MemoryController) ReinforceMemory(c *gin.Context) {
	memoryID := c.Param("id")
	if memoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "记忆ID不能为空"})
		return
	}

	var memory models.PetLongTermMemory
	if err := m.DB.Where("memory_id = ?", memoryID).First(&memory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "记忆不存在"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"reinforcement_count": gorm.Expr("reinforcement_count + 1"),
		"last_reinforced_at":  now,
		"decay_score":         gorm.Expr("decay_score * 1.1"),
		"updated_at":          now,
	}

	m.DB.Model(&memory).Updates(updates)
	m.DB.Where("memory_id = ?", memoryID).First(&memory)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": memory.ToResponse(),
	})
}
