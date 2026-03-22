package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BehaviorController 行为引擎控制器
type BehaviorController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册行为引擎路由
func (b *BehaviorController) RegisterRoutes(r *gin.RouterGroup) {
	// 动作库管理
	r.GET("/actions", b.ListActions)
	r.GET("/actions/:id", b.GetAction)
	r.POST("/actions", b.CreateAction)
	r.PUT("/actions/:id", b.UpdateAction)
	r.DELETE("/actions/:id", b.DeleteAction)

	// 行为记录
	r.GET("/pets/:device_id/behaviors", b.ListBehaviors)
	r.POST("/pets/:device_id/behaviors/trigger", b.TriggerBehavior)
}

// ListActions 获取动作列表
// @Summary 获取动作列表
// @Description 获取所有可用动作
// @Tags behavior
// @Accept json
// @Produce json
// @Param category query string false "动作分类"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/actions [GET]
func (b *BehaviorController) ListActions(c *gin.Context) {
	category := c.Query("category")

	var actions []models.ActionLibrary
	db := b.DB.Model(&models.ActionLibrary{})
	if category != "" {
		db = db.Where("category = ?", category)
	}

	if err := db.Order("priority DESC, created_at DESC").Find(&actions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询动作列表失败"})
		return
	}

	var responses []*models.ActionLibraryResponse
	for i := range actions {
		responses = append(responses, actions[i].ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// GetAction 获取单个动作
// @Summary 获取单个动作
// @Description 根据ID获取动作详情
// @Tags behavior
// @Accept json
// @Produce json
// @Param id path string true "动作ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/actions/{id} [GET]
func (b *BehaviorController) GetAction(c *gin.Context) {
	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "动作ID不能为空"})
		return
	}

	var action models.ActionLibrary
	if err := b.DB.Where("action_id = ?", actionID).First(&action).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "动作不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询动作失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// CreateAction 创建动作
// @Summary 创建动作
// @Description 创建新的动作
// @Tags behavior
// @Accept json
// @Produce json
// @Param body body models.ActionLibrary true "动作信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/actions [POST]
func (b *BehaviorController) CreateAction(c *gin.Context) {
	var action models.ActionLibrary
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误: " + err.Error()})
		return
	}

	if err := b.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建动作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// UpdateAction 更新动作
// @Summary 更新动作
// @Description 更新动作信息
// @Tags behavior
// @Accept json
// @Produce json
// @Param id path string true "动作ID"
// @Param body body models.ActionLibrary true "动作信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/actions/{id} [PUT]
func (b *BehaviorController) UpdateAction(c *gin.Context) {
	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "动作ID不能为空"})
		return
	}

	var action models.ActionLibrary
	if err := b.DB.Where("action_id = ?", actionID).First(&action).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "动作不存在"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	updates["updated_at"] = time.Now()
	if err := b.DB.Model(&action).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新动作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": action.ToResponse(),
	})
}

// DeleteAction 删除动作
// @Summary 删除动作
// @Description 删除动作
// @Tags behavior
// @Accept json
// @Produce json
// @Param id path string true "动作ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/actions/{id} [DELETE]
func (b *BehaviorController) DeleteAction(c *gin.Context) {
	actionID := c.Param("id")
	if actionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "动作ID不能为空"})
		return
	}

	if err := b.DB.Where("action_id = ?", actionID).Delete(&models.ActionLibrary{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除动作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ListBehaviors 获取行为记录列表
// @Summary 获取行为记录列表
// @Description 获取宠物的行为记录
// @Tags behavior
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Param limit query int false "数量限制"
// @Param offset query int false "偏移量"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/behaviors [GET]
func (b *BehaviorController) ListBehaviors(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v := parseInt(l); v > 0 {
			limit = min(v, 100)
		}
	}
	if o := c.Query("offset"); o != "" {
		if v := parseInt(o); v >= 0 {
			offset = v
		}
	}

	type BehaviorRecord struct {
		models.PetBehaviorAction
		CreatedAtStr string `json:"created_at_str"`
	}

	var behaviors []models.PetBehaviorAction
	if err := b.DB.Where("device_id = ?", deviceID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&behaviors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询行为记录失败"})
		return
	}

	var responses []map[string]interface{}
	for _, b := range behaviors {
		responses = append(responses, map[string]interface{}{
			"id":            b.ID,
			"device_id":     b.DeviceID,
			"action_type":   b.ActionType,
			"action_name":   b.ActionName,
			"sequence":      b.Sequence,
			"duration":      b.Duration,
			"trigger":       b.Trigger,
			"decision_path": b.DecisionPath,
			"parameters":    b.Parameters,
			"result":        b.Result,
			"created_at":    b.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": responses,
	})
}

// TriggerBehavior 触发行为
// @Summary 触发行为
// @Description 触发宠物的特定行为
// @Tags behavior
// @Accept json
// @Produce json
// @Param device_id path string true "设备ID"
// @Param body body map[string]interface{} true "行为参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/pets/{device_id}/behaviors/trigger [POST]
func (b *BehaviorController) TriggerBehavior(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "设备ID不能为空"})
		return
	}

	var params map[string]interface{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	actionType, _ := params["action_type"].(string)
	actionName, _ := params["action_name"].(string)

	// 下发动作到设备
	if mqtt.GlobalMQTTClient != nil {
		payload := map[string]interface{}{
			"type":        "behavior",
			"action_type": actionType,
			"action_name": actionName,
			"parameters":  params,
			"timestamp":   time.Now().Format(time.RFC3339),
		}
		if err := mqtt.PublishAction(mqtt.GlobalMQTTClient, deviceID, payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "下发行为失败"})
			return
		}
	}

	// 记录行为
	behavior := models.PetBehaviorAction{
		DeviceID:     deviceID,
		ActionType:   actionType,
		ActionName:   actionName,
		Trigger:      "manual",
		DecisionPath: "user_triggered",
		Parameters:   params,
		Result:       "sent",
	}
	b.DB.Create(&behavior)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"behavior_id": behavior.ID,
			"status":      "triggered",
		},
	})
}

func parseIntWithError(s string) (int, error) {
	var v int
	_, err := parseStrToInt(s, &v)
	return v, err
}

func parseStrToInt(s string, v *int) (bool, error) {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		*v = *v*10 + int(c-'0')
	}
	return true, nil
}
