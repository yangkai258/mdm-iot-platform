package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmbodiedController 具身智能控制器
type EmbodiedController struct {
	DB *gorm.DB
}

// RegisterEmbodiedRoutes 注册具身智能相关路由
func (ctrl *EmbodiedController) RegisterEmbodiedRoutes(api *gin.RouterGroup) {
	embodied := api.Group("/embodied")
	embodiedCtrl := &EmbodiedController{DB: ctrl.DB}

	// 环境感知 API
	embodied.GET("/:device_id/perception", embodiedCtrl.GetPerception)
	embodied.POST("/:device_id/perception/visual", embodiedCtrl.ReportVisualPerception)
	embodied.POST("/:device_id/perception/depth", embodiedCtrl.ReportDepthPerception)
	embodied.POST("/:device_id/perception/touch", embodiedCtrl.ReportTouchPerception)

	// 空间认知 API
	embodied.GET("/:device_id/map", embodiedCtrl.GetMap)
	embodied.POST("/:device_id/map/update", embodiedCtrl.UpdateMap)
	embodied.GET("/:device_id/localization", embodiedCtrl.GetLocalization)
	embodied.POST("/:device_id/localization/calibrate", embodiedCtrl.CalibrateLocalization)

	// 自主探索 API
	embodied.POST("/:device_id/navigate", embodiedCtrl.Navigate)
	embodied.POST("/:device_id/stop", embodiedCtrl.StopMovement)
	embodied.POST("/:device_id/follow", embodiedCtrl.FollowTarget)
	embodied.GET("/:device_id/explore/status", embodiedCtrl.GetExploreStatus)
	embodied.POST("/:device_id/explore/start", embodiedCtrl.StartExplore)

	// 动作模仿 API（设备相关）
	embodied.POST("/:device_id/action/execute", embodiedCtrl.ExecuteAction)
	embodied.POST("/:device_id/action/stop", embodiedCtrl.StopAction)

	// 决策引擎 API
	embodied.GET("/:device_id/decision/context", embodiedCtrl.GetDecisionContext)
	embodied.POST("/:device_id/decision/strategy", embodiedCtrl.SetDecisionStrategy)
	embodied.GET("/:device_id/decision/logs", embodiedCtrl.GetDecisionLogs)

	// 安全边界 API
	embodied.GET("/:device_id/safety/zones", embodiedCtrl.GetSafetyZones)
	embodied.POST("/:device_id/safety/zones", embodiedCtrl.CreateSafetyZone)
	embodied.DELETE("/:device_id/safety/zones/:id", embodiedCtrl.DeleteSafetyZone)
	embodied.POST("/:device_id/safety/emergency-stop", embodiedCtrl.EmergencyStop)
	embodied.GET("/:device_id/safety/logs", embodiedCtrl.GetSafetyLogs)

	// 动作模仿 API（全局，动作库管理）
	library := embodied.Group("")
	library.GET("/action-library", embodiedCtrl.ListActionLibrary)
	library.POST("/action-library/record", embodiedCtrl.RecordAction)
	library.POST("/action-library/:id/learn", embodiedCtrl.LearnAction)
	library.POST("/action-library/:id/share", embodiedCtrl.ShareAction)
}

// ===================== 环境感知 API =====================

// GetPerception 获取当前感知结果
func (ctrl *EmbodiedController) GetPerception(c *gin.Context) {
	deviceID := c.Param("device_id")
	perceptionType := c.DefaultQuery("type", "") // visual/depth/touch

	query := ctrl.DB.Model(&models.PerceptionData{}).Where("device_id = ?", deviceID)
	if perceptionType != "" {
		query = query.Where("type = ?", perceptionType)
	}

	var perception models.PerceptionData
	result := query.Order("created_at DESC").First(&perception)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "no perception data",
			"data":    nil,
		})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get perception"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": perception})
}

// ReportVisualPerception 上报视觉感知
func (ctrl *EmbodiedController) ReportVisualPerception(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		Data      string    `json:"data" binding:"required"`
		Timestamp time.Time `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	perception := models.PerceptionData{
		DeviceID:  deviceID,
		Type:      "visual",
		Data:      req.Data,
		Timestamp: req.Timestamp,
	}

	if err := ctrl.DB.Create(&perception).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to save perception"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "visual perception recorded", "data": perception})
}

// ReportDepthPerception 上报深度感知
func (ctrl *EmbodiedController) ReportDepthPerception(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		Data      string    `json:"data" binding:"required"`
		Timestamp time.Time `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	perception := models.PerceptionData{
		DeviceID:  deviceID,
		Type:      "depth",
		Data:      req.Data,
		Timestamp: req.Timestamp,
	}

	if err := ctrl.DB.Create(&perception).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to save perception"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "depth perception recorded", "data": perception})
}

// ReportTouchPerception 上报触觉感知
func (ctrl *EmbodiedController) ReportTouchPerception(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		Data      string    `json:"data" binding:"required"`
		Timestamp time.Time `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	perception := models.PerceptionData{
		DeviceID:  deviceID,
		Type:      "touch",
		Data:      req.Data,
		Timestamp: req.Timestamp,
	}

	if err := ctrl.DB.Create(&perception).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to save perception"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "touch perception recorded", "data": perception})
}

// ===================== 空间认知 API =====================

// GetMap 获取地图
func (ctrl *EmbodiedController) GetMap(c *gin.Context) {
	deviceID := c.Param("device_id")

	var m models.EmbodiedMap
	result := ctrl.DB.Where("device_id = ? AND is_active = ?", deviceID, true).First(&m)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "no map found", "data": nil})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get map"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": m})
}

// UpdateMap 更新地图
func (ctrl *EmbodiedController) UpdateMap(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		MapType    string  `json:"map_type" binding:"required"`
		MapData    string  `json:"map_data" binding:"required"`
		Resolution float64 `json:"resolution"`
		Size       string  `json:"size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	// 查找现有激活地图
	var existing models.EmbodiedMap
	findResult := ctrl.DB.Where("device_id = ? AND is_active = ?", deviceID, true).First(&existing)

	var m models.EmbodiedMap
	if findResult.Error == gorm.ErrRecordNotFound {
		// 创建新地图
		m = models.EmbodiedMap{
			DeviceID:   deviceID,
			MapType:    req.MapType,
			MapData:    req.MapData,
			Resolution: req.Resolution,
			Size:       req.Size,
			Version:    1,
			IsActive:   true,
		}
	} else {
		// 更新现有地图
		m = existing
		m.MapData = req.MapData
		m.MapType = req.MapType
		m.Resolution = req.Resolution
		m.Size = req.Size
		m.Version = existing.Version + 1
	}

	ops := []func(tx *gorm.DB) error{
		func(tx *gorm.DB) error {
			if findResult.Error == gorm.ErrRecordNotFound {
				return tx.Create(&m).Error
			}
			return tx.Save(&m).Error
		},
	}

	if err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
		return ops[0](tx)
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to update map"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "map updated", "data": m})
}

// GetLocalization 获取位置
func (ctrl *EmbodiedController) GetLocalization(c *gin.Context) {
	deviceID := c.Param("device_id")

	var pos models.SpatialPosition
	result := ctrl.DB.Where("device_id = ?", deviceID).Order("recorded_at DESC").First(&pos)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "no position found", "data": nil})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get localization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": pos})
}

// CalibrateLocalization 定位校准
func (ctrl *EmbodiedController) CalibrateLocalization(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		PositionX  float64 `json:"position_x" binding:"required"`
		PositionY  float64 `json:"position_y" binding:"required"`
		PositionZ  float64 `json:"position_z"`
		Orientation float64 `json:"orientation"`
		Confidence float64 `json:"confidence"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	pos := models.SpatialPosition{
		DeviceID:    deviceID,
		PositionX:   req.PositionX,
		PositionY:   req.PositionY,
		PositionZ:   req.PositionZ,
		Orientation: req.Orientation,
		Confidence:  req.Confidence,
		RecordedAt:  time.Now(),
	}

	if err := ctrl.DB.Create(&pos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to calibrate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "localization calibrated", "data": pos})
}

// ===================== 自主探索 API =====================

// Navigate 导航到目标
func (ctrl *EmbodiedController) Navigate(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		TargetX float64 `json:"target_x" binding:"required"`
		TargetY float64 `json:"target_y" binding:"required"`
		TargetZ float64 `json:"target_z"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	now := time.Now()
	task := models.NavigationTask{
		DeviceID:  deviceID,
		TargetX:  req.TargetX,
		TargetY:  req.TargetY,
		TargetZ:  req.TargetZ,
		Status:   "running",
		StartedAt: &now,
	}

	if err := ctrl.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create navigation task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "navigation started", "data": task})
}

// StopMovement 停止移动
func (ctrl *EmbodiedController) StopMovement(c *gin.Context) {
	deviceID := c.Param("device_id")

	// 停止所有运行中的导航任务
	result := ctrl.DB.Model(&models.NavigationTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{
			"status": "interrupted",
		})

	// 停止探索任务
	ctrl.DB.Model(&models.ExploreTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{
			"status": "paused",
		})

	// 停止跟随任务
	ctrl.DB.Model(&models.FollowTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{
			"status": "stopped",
		})

	c.JSON(http.StatusOK, gin.H{
		"code":           200,
		"message":        "movement stopped",
		"nav_tasks_stopped": result.RowsAffected,
	})
}

// FollowTarget 跟随目标
func (ctrl *EmbodiedController) FollowTarget(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		TargetID string  `json:"target_id" binding:"required"`
		Distance float64 `json:"distance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	if req.Distance == 0 {
		req.Distance = 1.5
	}

	// 停止现有跟随
	ctrl.DB.Model(&models.FollowTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{"status": "stopped"})

	now := time.Now()
	task := models.FollowTask{
		DeviceID:  deviceID,
		TargetID:  req.TargetID,
		Status:    "running",
		Distance:  req.Distance,
		StartedAt: &now,
	}

	if err := ctrl.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to start follow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "follow started", "data": task})
}

// GetExploreStatus 获取探索状态
func (ctrl *EmbodiedController) GetExploreStatus(c *gin.Context) {
	deviceID := c.Param("device_id")

	var task models.ExploreTask
	result := ctrl.DB.Where("device_id = ?", deviceID).Order("created_at DESC").First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "no explore task", "data": gin.H{
			"status":   "idle",
			"coverage": 0,
		}})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get explore status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": task})
}

// StartExplore 开始探索
func (ctrl *EmbodiedController) StartExplore(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		Strategy string `json:"strategy"`
	}
	c.ShouldBindJSON(&req) // optional

	now := time.Now()
	task := models.ExploreTask{
		DeviceID:  deviceID,
		Status:    "running",
		Strategy:  req.Strategy,
		StartedAt: &now,
	}

	if err := ctrl.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to start explore"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "explore started", "data": task})
}

// ===================== 动作模仿 API =====================

// ListActionLibrary 动作库列表
func (ctrl *EmbodiedController) ListActionLibrary(c *gin.Context) {
	page := defaultPage(c)
	pageSize := defaultPageSize(c)
	keyword := c.Query("keyword")
	category := c.Query("category")

	query := ctrl.DB.Model(&models.ActionLibrary{})
	if keyword != "" {
		query = query.Where("action_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to list action library", "detail": err.Error()})
		return
	}

	var items []models.ActionLibrary
	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to list action library", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"items":      items,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
		},
	})
}

// RecordAction 录制动作
func (ctrl *EmbodiedController) RecordAction(c *gin.Context) {
	var req struct {
		ActionName       string `json:"action_name" binding:"required"`
		Description      string `json:"description"`
		Category         string `json:"category"`
		DurationMs       int    `json:"duration_ms"`
		CompatibleModels string `json:"compatible_models"`
		Parameters       string `json:"parameters"`
		AnimationData    string `json:"animation_data"`
		MotorCommands    string `json:"motor_commands"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	action := models.ActionLibrary{
		ActionName:       req.ActionName,
		Description:      req.Description,
		Category:         req.Category,
		DurationMs:       req.DurationMs,
		CompatibleModels: req.CompatibleModels,
		Parameters:       req.Parameters,
		AnimationData:    req.AnimationData,
		MotorCommands:    req.MotorCommands,
		IsEmergency:      false,
	}

	if err := ctrl.DB.Create(&action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to record action"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "action recorded", "data": action})
}

// LearnAction 学习动作
func (ctrl *EmbodiedController) LearnAction(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		DeviceID string `json:"device_id" binding:"required"`
		Refine   string `json:"refine"` // optional refined sequence
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	var action models.ActionLibrary
	if err := ctrl.DB.First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "action not found"})
		return
	}

	// 记录执行作为学习过程
	now := time.Now()
	exec := models.ActionExecution{
		DeviceID:      req.DeviceID,
		ExecutionType: "learning",
		StartTime:     now,
		Status:        "running",
		Parameters:    req.Refine,
	}

	if err := ctrl.DB.Create(&exec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to record learning"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "learning started", "data": exec})
}

// ShareAction 分享动作
func (ctrl *EmbodiedController) ShareAction(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ShareTo string `json:"share_to"` // device_id or "public"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	var action models.ActionLibrary
	if err := ctrl.DB.First(&action, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "action not found"})
		return
	}

	// 简单标记分享（实际分享逻辑可扩展）
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "action shared", "data": gin.H{
		"action_id": action.ActionID,
		"share_to":  req.ShareTo,
		"shared_at": time.Now(),
	}})
}

// ExecuteAction 执行动作
func (ctrl *EmbodiedController) ExecuteAction(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		ActionID   uint    `json:"action_id" binding:"required"`
		Parameters string  `json:"parameters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	now := time.Now()
	exec := models.ActionExecution{
		DeviceID:      deviceID,
		ActionID:      req.ActionID,
		ExecutionType: "manual",
		StartTime:     now,
		Status:        "running",
		Parameters:    req.Parameters,
	}

	if err := ctrl.DB.Create(&exec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to execute action"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "action executed", "data": exec})
}

// StopAction 停止动作
func (ctrl *EmbodiedController) StopAction(c *gin.Context) {
	deviceID := c.Param("device_id")

	now := time.Now()
	result := ctrl.DB.Model(&models.ActionExecution{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{
			"status":              "interrupted",
			"end_time":            &now,
			"interruption_reason": "user_stopped",
		})

	c.JSON(http.StatusOK, gin.H{
		"code":             200,
		"message":          "action stopped",
		"actions_stopped": result.RowsAffected,
	})
}

// ===================== 决策引擎 API =====================

// GetDecisionContext 获取决策上下文
func (ctrl *EmbodiedController) GetDecisionContext(c *gin.Context) {
	deviceID := c.Param("device_id")

	// 获取当前定位
	var pos models.SpatialPosition
	ctrl.DB.Where("device_id = ?", deviceID).Order("recorded_at DESC").First(&pos)

	// 获取当前策略
	var strategy models.DecisionStrategy
	ctrl.DB.Where("device_id = ?", deviceID).First(&strategy)

	// 获取最新感知
	var perception models.PerceptionData
	ctrl.DB.Where("device_id = ?", deviceID).Order("created_at DESC").First(&perception)

	// 获取激活地图
	var m models.EmbodiedMap
	ctrl.DB.Where("device_id = ? AND is_active = ?", deviceID, true).First(&m)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": gin.H{
		"position":  pos,
		"strategy":  strategy,
		"perception": perception,
		"map":       m,
		"timestamp": time.Now(),
	}})
}

// SetDecisionStrategy 设置决策策略
func (ctrl *EmbodiedController) SetDecisionStrategy(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		Strategy string `json:"strategy" binding:"required"`
		Config   string `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	var s models.DecisionStrategy
	result := ctrl.DB.Where("device_id = ?", deviceID).First(&s)

	if result.Error == gorm.ErrRecordNotFound {
		s = models.DecisionStrategy{
			DeviceID: deviceID,
			Strategy: req.Strategy,
			Config:   req.Config,
		}
		if err := ctrl.DB.Create(&s).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to set strategy"})
			return
		}
	} else {
		s.Strategy = req.Strategy
		s.Config = req.Config
		if err := ctrl.DB.Save(&s).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to update strategy"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "strategy updated", "data": s})
}

// GetDecisionLogs 决策日志
func (ctrl *EmbodiedController) GetDecisionLogs(c *gin.Context) {
	deviceID := c.Param("device_id")
	page := defaultPage(c)
	pageSize := defaultPageSize(c)

	query := ctrl.DB.Model(&models.EmbodiedDecisionLog{}).Where("device_id = ?", deviceID)

	var total int64
	query.Count(&total)

	var logs []models.EmbodiedDecisionLog
	if err := query.Order("decided_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"items":     logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ===================== 安全边界 API =====================

// GetSafetyZones 获取禁区
func (ctrl *EmbodiedController) GetSafetyZones(c *gin.Context) {
	deviceID := c.Param("device_id")

	var zones []models.SafetyZone
	if err := ctrl.DB.Where("device_id = ?", deviceID).Order("created_at DESC").Find(&zones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get safety zones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": zones})
}

// CreateSafetyZone 设置禁区
func (ctrl *EmbodiedController) CreateSafetyZone(c *gin.Context) {
	deviceID := c.Param("device_id")
	var req struct {
		ZoneType  string `json:"zone_type" binding:"required"`
		ZoneShape string `json:"zone_shape" binding:"required"`
		ZoneData  string `json:"zone_data" binding:"required"`
		ZoneName  string `json:"zone_name"`
		IsEnabled *bool  `json:"is_enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}

	zone := models.SafetyZone{
		DeviceID:  deviceID,
		ZoneType:  req.ZoneType,
		ZoneShape: req.ZoneShape,
		ZoneData:  req.ZoneData,
		ZoneName:  req.ZoneName,
		IsEnabled: enabled,
	}

	if err := ctrl.DB.Create(&zone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to create safety zone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "safety zone created", "data": zone})
}

// DeleteSafetyZone 删除禁区
func (ctrl *EmbodiedController) DeleteSafetyZone(c *gin.Context) {
	deviceID := c.Param("device_id")
	id := c.Param("id")

	result := ctrl.DB.Where("id = ? AND device_id = ?", id, deviceID).Delete(&models.SafetyZone{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to delete safety zone"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "safety zone not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "safety zone deleted"})
}

// EmergencyStop 紧急停止
func (ctrl *EmbodiedController) EmergencyStop(c *gin.Context) {
	deviceID := c.Param("device_id")

	// 停止所有运行中的任务
	now := time.Now()
	ctrl.DB.Model(&models.NavigationTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{"status": "failed", "error_message": "emergency_stop", "completed_at": &now})

	ctrl.DB.Model(&models.ExploreTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{"status": "paused"})

	ctrl.DB.Model(&models.FollowTask{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{"status": "stopped", "stopped_at": &now})

	ctrl.DB.Model(&models.ActionExecution{}).
		Where("device_id = ? AND status = ?", deviceID, "running").
		Updates(map[string]interface{}{"status": "interrupted", "interruption_reason": "emergency_stop", "end_time": &now})

	// 记录安全日志
	safetyLog := models.SafetyLog{
		DeviceID:  deviceID,
		EventType: "emergency_stop",
		Severity:  "critical",
		Details:   `{"reason":"manual_emergency_stop"}`,
		Resolved:  false,
	}
	ctrl.DB.Create(&safetyLog)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "emergency stop executed"})
}

// GetSafetyLogs 安全日志
func (ctrl *EmbodiedController) GetSafetyLogs(c *gin.Context) {
	deviceID := c.Param("device_id")
	page := defaultPage(c)
	pageSize := defaultPageSize(c)
	eventType := c.Query("event_type")

	query := ctrl.DB.Model(&models.SafetyLog{}).Where("device_id = ?", deviceID)
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	var total int64
	query.Count(&total)

	var logs []models.SafetyLog
	if err := query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to get safety logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": gin.H{
			"items":     logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
