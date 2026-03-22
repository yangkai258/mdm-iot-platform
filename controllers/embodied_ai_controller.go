package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EmbodiedAIController 具身智能控制器
type EmbodiedAIController struct {
	DB *gorm.DB
}

// NewEmbodiedAIController 创建控制器
func NewEmbodiedAIController(db *gorm.DB) *EmbodiedAIController {
	return &EmbodiedAIController{DB: db}
}

// RegisterRoutes 注册具身智能路由
func (c *EmbodiedAIController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/ai/embodied/perceive", c.Perceive)
	rg.POST("/ai/embodied/spatial", c.Spatial)
	rg.POST("/ai/embodied/explore", c.Explore)
	rg.POST("/ai/embodied/interact", c.Interact)
	rg.GET("/ai/embodied/state", c.GetState)
	rg.GET("/ai/embodied/capabilities", c.GetCapabilities)
}

// defaultEmbodiedAIState 获取或创建默认具身AI状态
func (c *EmbodiedAIController) defaultEmbodiedAIState(entityID, entityType string) *models.EmbodiedAIState {
	var state models.EmbodiedAIState
	stateKey := entityID + "_" + entityType

	err := c.DB.Where("state_key = ?", stateKey).First(&state).Error
	if err == gorm.ErrRecordNotFound {
		// 创建默认状态
		state = models.EmbodiedAIState{
			StateKey:        stateKey,
			EntityType:      entityType,
			EntityID:        entityID,
			PositionX:       0,
			PositionY:       0,
			PositionZ:       0,
			Orientation:     0,
			LocationName:    "unknown",
			PerceptionMode:  "active",
			AlertLevel:     "normal",
			EnergyLevel:    100,
			Mood:           "happy",
			MoodIntensity:  0.5,
			CurrentGoal:    "",
			GoalProgress:   0,
			LastUpdatedAt:   time.Now(),
			LastActiveAt:    time.Now(),
		}
		c.DB.Create(&state)
	}
	return &state
}

// getDefaultCapabilities 返回默认能力列表
func getDefaultCapabilities(entityType string) []models.Capability {
	baseCaps := []models.Capability{
		{Key: "vision", Name: "视觉感知", Description: "识别物体、人物、场景", Available: true, Confidence: 0.92},
		{Key: "audio", Name: "听觉感知", Description: "识别声音、语音、方向", Available: true, Confidence: 0.88},
		{Key: "tactile", Name: "触觉感知", Description: "感知触碰、压力、温度", Available: true, Confidence: 0.75},
		{Key: "spatial_mapping", Name: "空间建图", Description: "构建和更新环境地图", Available: true, Confidence: 0.85},
		{Key: "path_planning", Name: "路径规划", Description: "自主规划移动路径", Available: true, Confidence: 0.90},
		{Key: "object_recognition", Name: "物体识别", Description: "识别常见物体并分类", Available: true, Confidence: 0.88},
		{Key: "grasp_planning", Name: "抓取规划", Description: "规划物体抓取策略", Available: true, Confidence: 0.78},
		{Key: "exploration", Name: "自主探索", Description: "未知环境自主探索", Available: true, Confidence: 0.82},
	}

	if entityType == "robot" {
		baseCaps = append(baseCaps, models.Capability{
			Key: "manipulation", Name: "精细操作", Description: "执行精细机械操作", Available: true, Confidence: 0.85,
		})
	}
	return baseCaps
}

// getCapabilitiesMap 返回能力映射
func getCapabilitiesMap(entityType string) string {
	caps := getDefaultCapabilities(entityType)
	data, _ := json.Marshal(caps)
	return string(data)
}

// getSensoryStatusMap 返回传感器状态映射
func getSensoryStatusMap() string {
	status := map[string]interface{}{
		"vision":  map[string]interface{}{"status": "online", "last_calibrated": time.Now().Format(time.RFC3339)},
		"audio":   map[string]interface{}{"status": "online", "last_calibrated": time.Now().Format(time.RFC3339)},
		"tactile": map[string]interface{}{"status": "online", "last_calibrated": time.Now().Format(time.RFC3339)},
		"imu":     map[string]interface{}{"status": "online", "last_calibrated": time.Now().Format(time.RFC3339)},
	}
	data, _ := json.Marshal(status)
	return string(data)
}

// ============ API 实现 ============

// Perceive 环境感知
// POST /api/v1/ai/embodied/perceive
func (c *EmbodiedAIController) Perceive(ctx *gin.Context) {
	var req models.PerceiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 获取或创建状态
	state := c.defaultEmbodiedAIState(req.EntityID, req.EntityType)

	// 模拟感知处理
	var objects []string
	var events []string
	var entities []string
	var envTags []string

	switch req.PerceiveMode {
	case "visual", "all":
		objects = []string{"table", "chair", "person", "cup"}
		entities = []string{"person_001", "pet_002"}
		envTags = []string{"indoor", "living_room", "well_lit"}
	case "audio":
		events = []string{"footstep", "door_open", "voice"}
		envTags = []string{"indoor", "moderate_noise"}
	case "tactile":
		objects = []string{"surface_smooth", "object_soft"}
		envTags = []string{"indoor", "room_temperature"}
	default:
		objects = []string{"table", "chair", "person"}
		entities = []string{"person_001"}
		envTags = []string{"indoor", "living_room"}
	}

	// 更新状态中的感知缓存
	state.LastSeenObjects = objects
	state.NearbyEntities = entities
	state.EnvironmentTags = envTags
	state.LastActiveAt = time.Now()
	state.LastUpdatedAt = time.Now()
	c.DB.Save(state)

	resp := models.PerceiveResponse{
		EntityID:        req.EntityID,
		Timestamp:       time.Now().Format(time.RFC3339),
		Objects:         objects,
		Events:          events,
		Entities:        entities,
		EnvironmentTags: envTags,
		Confidence:      0.89,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "感知成功",
		"data":    resp,
	})
}

// Spatial 空间认知
// POST /api/v1/ai/embodied/spatial
func (c *EmbodiedAIController) Spatial(ctx *gin.Context) {
	var req models.SpatialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 获取或创建状态
	state := c.defaultEmbodiedAIState(req.EntityID, "pet")

	resp := models.SpatialResponse{
		EntityID:  req.EntityID,
		QueryType: req.QueryType,
	}

	switch req.QueryType {
	case "position":
		resp.Position = models.SpatialPosition{
			X: state.PositionX,
			Y: state.PositionY,
			Z: state.PositionZ,
		}
		resp.Orientation = state.Orientation

	case "direction":
		if req.TargetX != 0 || req.TargetY != 0 {
			dx := req.TargetX - state.PositionX
			dy := req.TargetY - state.PositionY
			resp.Direction = calculateDirection(dx, dy)
		}
		resp.Position = models.SpatialPosition{
			X: state.PositionX,
			Y: state.PositionY,
		}

	case "distance":
		if req.TargetX != 0 || req.TargetY != 0 {
			resp.Distance = calculateDistance(state.PositionX, state.PositionY, req.TargetX, req.TargetY)
		}

	case "map":
		var envMap models.EnvironmentMap
		mapKey := req.EntityID + "_default"
		if req.MapKey != "" {
			mapKey = req.MapKey
		}
		if err := c.DB.Where("map_key = ?", mapKey).First(&envMap).Error; err == gorm.ErrRecordNotFound {
			// 创建默认地图
			envMap = models.EnvironmentMap{
				MapKey:         mapKey,
				EntityID:       req.EntityID,
				EntityType:     "pet",
				MapName:        "Default Environment Map",
				MapType:        "2d_grid",
				Resolution:     0.1,
				Scale:          10.0,
				Width:          10.0,
				Height:         10.0,
				Confidence:     85.0,
				ExploredRatio:  45.0,
				IsActive:       true,
				Version:        1,
				LastUpdatedAt:  time.Now(),
			}
			c.DB.Create(&envMap)
		}
		resp.MapData = string(formatMapData(envMap))

	default:
		resp.Position = models.SpatialPosition{
			X: state.PositionX,
			Y: state.PositionY,
			Z: state.PositionZ,
		}
		resp.Orientation = state.Orientation
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "空间查询成功",
		"data":    resp,
	})
}

// Explore 自主探索
// POST /api/v1/ai/embodied/explore
func (c *EmbodiedAIController) Explore(ctx *gin.Context) {
	var req models.ExploreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 创建探索会话
	session := models.ExplorationSession{
		EntityID:          req.EntityID,
		EntityType:        req.EntityType,
		Strategy:          req.Strategy,
		ExplorationGoal:   req.ExplorationGoal,
		MaxDuration:       req.MaxDuration,
		MaxDistance:       req.MaxDistance,
		Status:            "active",
		Progress:          0,
		Waypoints:         []string{},
		VisitedAreas:      []string{},
		DiscoveredObjects: []string{},
		CoverageRate:      0,
		PathLength:        0,
		NewDiscoveryCount: 0,
		StartedAt:         time.Now(),
	}

	// 设置边界
	if req.Boundary != "" {
		var boundary map[string]float64
		if err := json.Unmarshal([]byte(req.Boundary), &boundary); err == nil {
			session.BoundaryMinX = boundary["min_x"]
			session.BoundaryMaxX = boundary["max_x"]
			session.BoundaryMinY = boundary["min_y"]
			session.BoundaryMaxY = boundary["max_y"]
		}
	}

	if err := c.DB.Create(&session).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5001,
			"message": "创建探索会话失败: " + err.Error(),
		})
		return
	}

	// 异步模拟探索进度
	go c.simulateExploration(session.ID, req)

	// 计算预估时间
	estimatedTime := 60
	if req.MaxDuration > 0 {
		estimatedTime = req.MaxDuration
	}

	resp := models.ExploreResponse{
		SessionKey:          session.SessionKey,
		Strategy:            session.Strategy,
		Status:              session.Status,
		Progress:            0,
		Waypoints:           []string{},
		DiscoveredObjects:   []string{},
		EstimatedTime:       estimatedTime,
		CoverageRate:        0,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "探索会话已创建",
		"data":    resp,
	})
}

// simulateExploration 异步模拟探索过程
func (c *EmbodiedAIController) simulateExploration(sessionID uint, req models.ExploreRequest) {
	steps := 10
	stepDuration := 100 * time.Millisecond
	if req.MaxDuration > 0 {
		stepDuration = time.Duration(req.MaxDuration*1000/steps) * time.Millisecond
	}

	for i := 1; i <= steps; i++ {
		time.Sleep(stepDuration)

		progress := float64(i) / float64(steps) * 100
		coverage := progress * 0.9

		// 生成随机路径点
		waypoint := map[string]interface{}{
			"x": float64(i) * 0.5,
			"y": float64(i%3) * 0.3,
			"t": time.Now().Format(time.RFC3339),
		}
		waypointJSON, _ := json.Marshal(waypoint)

		// 更新会话
		c.DB.Model(&models.ExplorationSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
			"progress":          progress,
			"coverage_rate":     coverage,
			"path_length":       float64(i) * 1.2,
			"new_discovery_count": i % 3,
		})
		c.DB.Exec("UPDATE exploration_sessions SET waypoints = array_append(waypoints, ?) WHERE id = ?", string(waypointJSON), sessionID)
	}

	// 探索完成
	now := time.Now()
	c.DB.Model(&models.ExplorationSession{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
		"status":      "completed",
		"progress":    100,
		"coverage_rate": 95.0,
		"ended_at":    now,
	})
}

// Interact 环境交互
// POST /api/v1/ai/embodied/interact
func (c *EmbodiedAIController) Interact(ctx *gin.Context) {
	var req models.InteractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 获取状态
	stateKey := req.EntityID + "_pet"
	var state models.EmbodiedAIState
	if err := c.DB.Where("state_key = ?", stateKey).First(&state).Error; err != nil {
		state = *c.defaultEmbodiedAIState(req.EntityID, "pet")
	}

	resp := models.InteractResponse{
		EntityID:   req.EntityID,
		ActionType: req.ActionType,
	}

	// 模拟交互处理
	energyUsed := 0.0
	success := true
	result := "动作执行成功"
	feedback := "已完成目标操作"

	switch req.ActionType {
	case "move":
		state.PositionX = req.TargetX
		state.PositionY = req.TargetY
		state.PositionZ = req.TargetZ
		energyUsed = calculateDistance(state.PositionX, state.PositionY, req.TargetX, req.TargetY) * 0.5
		result = "移动完成"

	case "grasp":
		energyUsed = 5.0
		if req.TargetObject == "" {
			success = false
			result = "抓取失败：未指定目标物体"
			feedback = "请指定要抓取的物体"
		} else {
			result = "成功抓取: " + req.TargetObject
			feedback = "抓取力度适中，物体稳定"
		}

	case "place":
		energyUsed = 3.0
		result = "已将物体放置到目标位置"
		feedback = "放置位置准确"

	case "push", "pull":
		energyUsed = 4.0
		result = "推拉动作完成"
		feedback = "力度控制良好"

	case "activate":
		energyUsed = 2.0
		result = "激活操作完成"
		feedback = "目标设备已响应"

	default:
		energyUsed = 1.0
		result = "动作执行完成"
		feedback = "动作已执行"
	}

	// 更新状态
	state.EnergyLevel -= energyUsed
	if state.EnergyLevel < 0 {
		state.EnergyLevel = 0
	}
	state.LastActiveAt = time.Now()
	state.LastUpdatedAt = time.Now()
	c.DB.Save(&state)

	resp.Success = success
	resp.Result = result
	resp.Feedback = feedback
	resp.EnergyUsed = energyUsed
	resp.NewPosition = models.SpatialPosition{
		X: state.PositionX,
		Y: state.PositionY,
		Z: state.PositionZ,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": result,
		"data":    resp,
	})
}

// GetState 获取AI状态
// GET /api/v1/ai/embodied/state
func (c *EmbodiedAIController) GetState(ctx *gin.Context) {
	entityID := ctx.Query("entity_id")
	entityType := ctx.DefaultQuery("entity_type", "pet")

	if entityID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": "entity_id 为必填参数",
		})
		return
	}

	state := c.defaultEmbodiedAIState(entityID, entityType)

	// 确保能力和传感器状态已更新
	if state.Capabilities == "" {
		state.Capabilities = getCapabilitiesMap(entityType)
		c.DB.Save(state)
	}

	resp := models.EmbodiedAIStateResponse{
		StateKey:       state.StateKey,
		EntityID:       state.EntityID,
		EntityType:     state.EntityType,
		Position: models.SpatialPosition{
			X: state.PositionX,
			Y: state.PositionY,
			Z: state.PositionZ,
		},
		Orientation:    state.Orientation,
		LocationName:   state.LocationName,
		PerceptionMode: state.PerceptionMode,
		AlertLevel:     state.AlertLevel,
		EnergyLevel:    state.EnergyLevel,
		Mood:           state.Mood,
		MoodIntensity:  state.MoodIntensity,
		CurrentGoal:    state.CurrentGoal,
		GoalProgress:   state.GoalProgress,
		Capabilities:   state.Capabilities,
		LastActiveAt:   state.LastActiveAt.Format(time.RFC3339),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    resp,
	})
}

// GetCapabilities 获取能力列表
// GET /api/v1/ai/embodied/capabilities
func (c *EmbodiedAIController) GetCapabilities(ctx *gin.Context) {
	entityID := ctx.Query("entity_id")
	entityType := ctx.DefaultQuery("entity_type", "pet")

	if entityID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": "entity_id 为必填参数",
		})
		return
	}

	caps := getDefaultCapabilities(entityType)

	resp := models.CapabilitiesResponse{
		EntityID:     entityID,
		EntityType:   entityType,
		Capabilities: caps,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    resp,
	})
}

// ============ 辅助函数 ============

// calculateDirection 计算朝向角度
func calculateDirection(dx, dy float64) float64 {
	if dx == 0 && dy == 0 {
		return 0
	}
	angle := 0.0
	if dx == 0 {
		if dy > 0 {
			angle = 90.0
		} else {
			angle = 270.0
		}
	} else {
		angle = 180.0 * (dx / (dx + dy))
	}
	return angle
}

// calculateDistance 计算两点距离
func calculateDistance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return (dx*dx + dy*dy) * 100 // 简化计算
}

// formatMapData 格式化地图数据
func formatMapData(envMap models.EnvironmentMap) []byte {
	data := map[string]interface{}{
		"map_key":       envMap.MapKey,
		"map_name":      envMap.MapName,
		"map_type":      envMap.MapType,
		"resolution":    envMap.Resolution,
		"width":         envMap.Width,
		"height":        envMap.Height,
		"confidence":    envMap.Confidence,
		"explored_ratio": envMap.ExploredRatio,
		"version":       envMap.Version,
		"last_updated":  envMap.LastUpdatedAt.Format(time.RFC3339),
	}
	result, _ := json.Marshal(data)
	return result
}

// ============ 探索会话管理（额外CRUD） ============

// ListExplorationSessions 获取探索会话列表
// GET /api/v1/ai/embodied/explore/sessions
func (c *EmbodiedAIController) ListExplorationSessions(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.ExplorationSession{})

	if entityID := ctx.Query("entity_id"); entityID != "" {
		query = query.Where("entity_id = ?", entityID)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if strategy := ctx.Query("strategy"); strategy != "" {
		query = query.Where("strategy = ?", strategy)
	}

	var total int64
	query.Count(&total)

	var sessions []models.ExplorationSession
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&sessions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": sessions,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetExplorationSession 获取单个探索会话
// GET /api/v1/ai/embodied/explore/sessions/:id
func (c *EmbodiedAIController) GetExplorationSession(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的会话ID"})
		return
	}

	var session models.ExplorationSession
	if err := c.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "探索会话不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    session,
	})
}

// StopExplorationSession 停止探索会话
// POST /api/v1/ai/embodied/explore/sessions/:id/stop
func (c *EmbodiedAIController) StopExplorationSession(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "无效的会话ID"})
		return
	}

	var session models.ExplorationSession
	if err := c.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "探索会话不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if session.Status != "active" && session.Status != "paused" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "当前状态不允许停止"})
		return
	}

	now := time.Now()
	if err := c.DB.Model(&session).Updates(map[string]interface{}{
		"status":     "aborted",
		"ended_at":   now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "探索会话已停止",
		"data": gin.H{
			"session_key": session.SessionKey,
			"status":      "aborted",
		},
	})
}
