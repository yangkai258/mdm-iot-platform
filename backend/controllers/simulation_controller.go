package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SimulationController 仿真测试控制器
type SimulationController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// RegisterSimulationRoutes 注册仿真测试相关路由
func (sc *SimulationController) RegisterSimulationRoutes(r *gin.RouterGroup) {
	// 仿真场景 API
	r.GET("/simulation/scenes", sc.ListScenes)
	r.POST("/simulation/scenes", sc.CreateScene)
	r.GET("/simulation/scenes/:id", sc.GetScene)
	r.PUT("/simulation/scenes/:id", sc.UpdateScene)
	r.DELETE("/simulation/scenes/:id", sc.DeleteScene)

	// 仿真会话 API
	r.POST("/simulation/sessions", sc.CreateSession)
	r.GET("/simulation/sessions/:id", sc.GetSession)
	r.POST("/simulation/sessions/:id/start", sc.StartSession)
	r.POST("/simulation/sessions/:id/stop", sc.StopSession)
	r.GET("/simulation/sessions/:id/results", sc.GetSessionResults)

	// 虚拟宠物仿真 API
	r.POST("/simulation/pet/start", sc.StartPetSimulation)
	r.POST("/simulation/pet/stop", sc.StopPetSimulation)
	r.GET("/simulation/pet/:session_id", sc.GetPetStatus)

	// 压力测试 API
	r.POST("/simulation/stress-test", sc.CreateStressTest)
	r.GET("/simulation/stress-test/:id", sc.GetStressTest)
	r.POST("/simulation/stress-test/:id/run", sc.RunStressTest)
	r.POST("/simulation/stress-test/:id/stop", sc.StopStressTest)
	r.GET("/simulation/stress-test/:id/report", sc.GetStressTestReport)

	// A/B实验仿真 API
	r.POST("/simulation/ab-experiments", sc.CreateABExperiment)
	r.GET("/simulation/ab-experiments/:id", sc.GetABExperiment)
	r.POST("/simulation/ab-experiments/:id/run", sc.RunABExperiment)
	r.GET("/simulation/ab-experiments/:id/results", sc.GetABExperimentResults)

	// 回放系统 API
	r.POST("/simulation/recordings", sc.CreateRecording)
	r.GET("/simulation/recordings/:id", sc.GetRecording)
	r.POST("/simulation/recordings/:id/replay", sc.ReplayRecording)
	r.POST("/simulation/recordings/:id/stop", sc.StopRecording)
	r.GET("/simulation/recordings", sc.ListRecordings)
}

// ============================================================
// 仿真场景 API
// ============================================================

// ListScenes 获取场景列表
func (sc *SimulationController) ListScenes(c *gin.Context) {
	var scenes []models.SimulationScene
	query := sc.DB.Model(&models.SimulationScene{})

	// 筛选参数
	if sceneType := c.Query("scene_type"); sceneType != "" {
		query = query.Where("scene_type = ?", sceneType)
	}
	if isPublic := c.Query("is_public"); isPublic != "" {
		query = query.Where("is_public = ?", isPublic == "true")
	}
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags @> ?", fmt.Sprintf(`["%s"]`, strings.Join(strings.Split(tags, ","), "\",\"")))
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&scenes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取场景列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items":    scenes,
			"total":    total,
			"page":     page,
			"page_size": pageSize,
		},
	})
}

// CreateScene 创建仿真场景
func (sc *SimulationController) CreateScene(c *gin.Context) {
	var input struct {
		SceneName   string                  `json:"scene_name" binding:"required"`
		SceneType   string                  `json:"scene_type" binding:"required"`
		Environment *models.JSONB           `json:"environment"`
		Objects     []string                 `json:"objects"`
		Events      *models.JSONB           `json:"events"`
		Config      *models.JSONB           `json:"config"`
		IsPublic    bool                    `json:"is_public"`
		Tags        []string                 `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	scene := models.SimulationScene{
		SceneName: input.SceneName,
		SceneType: input.SceneType,
		Environment: func() models.JSONB {
			if input.Environment != nil {
				return *input.Environment
			}
			return models.JSONB{}
		}(),
		Objects: input.Objects,
		Events: func() models.JSONB {
			if input.Events != nil {
				return *input.Events
			}
			return models.JSONB{}
		}(),
		Config: func() models.JSONB {
			if input.Config != nil {
				return *input.Config
			}
			return models.JSONB{}
		}(),
		IsPublic: input.IsPublic,
		Tags:     input.Tags,
		Status:   "idle",
	}

	// 获取当前用户ID
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		scene.CreatedBy = &uid
	}

	if err := sc.DB.Create(&scene).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建场景失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"message": "success",
		"data": scene,
	})
}

// GetScene 获取场景详情
func (sc *SimulationController) GetScene(c *gin.Context) {
	id := c.Param("id")

	var scene models.SimulationScene
	if err := sc.DB.First(&scene, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "场景不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取场景失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": scene})
}

// UpdateScene 更新场景
func (sc *SimulationController) UpdateScene(c *gin.Context) {
	id := c.Param("id")

	var scene models.SimulationScene
	if err := sc.DB.First(&scene, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "场景不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取场景失败", "error": err.Error()})
		return
	}

	var input struct {
		SceneName   *string                  `json:"scene_name"`
		SceneType   *string                  `json:"scene_type"`
		Environment *models.JSONB           `json:"environment"`
		Objects     []string                 `json:"objects"`
		Events      *models.JSONB           `json:"events"`
		Config      *models.JSONB           `json:"config"`
		IsPublic    *bool                    `json:"is_public"`
		Tags        []string                 `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if input.SceneName != nil {
		updates["scene_name"] = *input.SceneName
	}
	if input.SceneType != nil {
		updates["scene_type"] = *input.SceneType
	}
	if input.Environment != nil {
		updates["environment"] = *input.Environment
	}
	if input.Objects != nil {
		updates["objects"] = models.StringSlice(input.Objects)
	}
	if input.Events != nil {
		updates["events"] = *input.Events
	}
	if input.Config != nil {
		updates["config"] = *input.Config
	}
	if input.IsPublic != nil {
		updates["is_public"] = *input.IsPublic
	}
	if input.Tags != nil {
		updates["tags"] = models.StringSlice(input.Tags)
	}

	if err := sc.DB.Model(&scene).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "更新场景失败", "error": err.Error()})
		return
	}

	sc.DB.First(&scene, id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": scene})
}

// DeleteScene 删除场景
func (sc *SimulationController) DeleteScene(c *gin.Context) {
	id := c.Param("id")

	var scene models.SimulationScene
	if err := sc.DB.First(&scene, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "场景不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取场景失败", "error": err.Error()})
		return
	}

	if err := sc.DB.Delete(&scene).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "删除场景失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============================================================
// 仿真会话 API
// ============================================================

// CreateSession 创建仿真会话
func (sc *SimulationController) CreateSession(c *gin.Context) {
	var input struct {
		SceneID     uint                  `json:"scene_id" binding:"required"`
		PetID       *uint                 `json:"pet_id"`
		SessionName string                `json:"session_name"`
		Parameters  *models.JSONB         `json:"parameters"`
		Environment *models.JSONB         `json:"environment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	session := models.SimulationSession{
		SceneID: input.SceneID,
		PetID:   input.PetID,
		SessionName: input.SessionName,
		Status:  "created",
		Parameters: func() models.JSONB {
			if input.Parameters != nil {
				return *input.Parameters
			}
			return models.JSONB{}
		}(),
		Environment: func() models.JSONB {
			if input.Environment != nil {
				return *input.Environment
			}
			return models.JSONB{}
		}(),
		CurrentState: models.JSONB{"step": 0, "progress": 0},
	}

	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		session.CreatedBy = &uid
	}

	if err := sc.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建会话失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": session})
}

// GetSession 获取会话详情
func (sc *SimulationController) GetSession(c *gin.Context) {
	id := c.Param("id")

	var session models.SimulationSession
	if err := sc.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取会话失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": session})
}

// StartSession 开始仿真
func (sc *SimulationController) StartSession(c *gin.Context) {
	id := c.Param("id")

	var session models.SimulationSession
	if err := sc.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取会话失败", "error": err.Error()})
		return
	}

	now := time.Now()
	if err := sc.DB.Model(&session).Updates(map[string]interface{}{
		"status":     "running",
		"start_time": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "启动仿真失败", "error": err.Error()})
		return
	}

	// 更新场景状态
	sc.DB.Model(&models.SimulationScene{}).Where("id = ?", session.SceneID).Update("status", "running")

	session.Status = "running"
	session.StartTime = &now
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": session})
}

// StopSession 停止仿真
func (sc *SimulationController) StopSession(c *gin.Context) {
	id := c.Param("id")

	var session models.SimulationSession
	if err := sc.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取会话失败", "error": err.Error()})
		return
	}

	now := time.Now()
	var durationMs int
	if session.StartTime != nil {
		durationMs = int(now.Sub(*session.StartTime).Milliseconds())
	}

	if err := sc.DB.Model(&session).Updates(map[string]interface{}{
		"status":      "stopped",
		"end_time":     now,
		"duration_ms":  durationMs,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "停止仿真失败", "error": err.Error()})
		return
	}

	// 更新场景状态
	sc.DB.Model(&models.SimulationScene{}).Where("id = ?", session.SceneID).Update("status", "idle")

	session.Status = "stopped"
	session.EndTime = &now
	session.DurationMs = durationMs
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": session})
}

// GetSessionResults 获取会话仿真结果
func (sc *SimulationController) GetSessionResults(c *gin.Context) {
	id := c.Param("id")

	var session models.SimulationSession
	if err := sc.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取会话失败", "error": err.Error()})
		return
	}

	var results []models.SimulationResult
	if err := sc.DB.Where("session_id = ?", id).Order("created_at DESC").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取仿真结果失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": results})
}

// ============================================================
// 虚拟宠物仿真 API
// ============================================================

// StartPetSimulation 启动虚拟宠物仿真
func (sc *SimulationController) StartPetSimulation(c *gin.Context) {
	var input struct {
		PetID       uint                  `json:"pet_id"`
		PetName     string                `json:"pet_name"`
		PetType     string                `json:"pet_type"`
		Environment *models.JSONB         `json:"environment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果没有输入，创建一个默认的虚拟宠物会话
		input = struct {
			PetID       uint                  `json:"pet_id"`
			PetName     string                `json:"pet_name"`
			PetType     string                `json:"pet_type"`
			Environment *models.JSONB         `json:"environment"`
		}{PetName: "默认宠物", PetType: "cat"}
	}

	// 创建仿真会话
	session := models.SimulationSession{
		SessionName: fmt.Sprintf("pet_simulation_%d", time.Now().Unix()),
		Status:      "running",
		Environment: func() models.JSONB {
			if input.Environment != nil {
				return *input.Environment
			}
			return models.JSONB{
				"pet_name": input.PetName,
				"pet_type": input.PetType,
				"status":   "running",
				"emotion":  "happy",
			}
		}(),
		CurrentState: models.JSONB{
			"emotion":     "happy",
			"position":    models.JSONB{"x": 0, "y": 0, "z": 0},
			"battery":     1.0,
			"action":      "idle",
		},
	}

	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		session.CreatedBy = &uid
	}
	if input.PetID > 0 {
		session.PetID = &input.PetID
	}

	now := time.Now()
	session.StartTime = &now

	if err := sc.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "启动虚拟宠物失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"session_id":  session.ID,
			"status":      session.Status,
			"pet_name":    input.PetName,
			"pet_type":    input.PetType,
			"emotion":     "happy",
			"started_at":  now,
		},
	})
}

// StopPetSimulation 停止虚拟宠物仿真
func (sc *SimulationController) StopPetSimulation(c *gin.Context) {
	var input struct {
		SessionID uint `json:"session_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	var session models.SimulationSession
	if err := sc.DB.First(&session, input.SessionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取会话失败", "error": err.Error()})
		return
	}

	now := time.Now()
	var durationMs int
	if session.StartTime != nil {
		durationMs = int(now.Sub(*session.StartTime).Milliseconds())
	}

	if err := sc.DB.Model(&session).Updates(map[string]interface{}{
		"status":     "stopped",
		"end_time":    now,
		"duration_ms": durationMs,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "停止虚拟宠物失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"session_id":  session.ID,
			"status":      "stopped",
			"duration_ms": durationMs,
			"stopped_at":  now,
		},
	})
}

// GetPetStatus 获取宠物状态
func (sc *SimulationController) GetPetStatus(c *gin.Context) {
	sessionID := c.Param("session_id")
	id, err := strconv.ParseUint(sessionID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "无效的会话ID"})
		return
	}

	var session models.SimulationSession
	if err := sc.DB.First(&session, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "会话不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取宠物状态失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"session_id":    session.ID,
			"status":        session.Status,
			"current_state": session.CurrentState,
			"start_time":    session.StartTime,
			"duration_ms":   session.DurationMs,
		},
	})
}

// ============================================================
// 压力测试 API
// ============================================================

// CreateStressTest 创建压力测试
func (sc *SimulationController) CreateStressTest(c *gin.Context) {
	var input struct {
		TestName string         `json:"test_name" binding:"required"`
		TestType string         `json:"test_type" binding:"required"`
		Config   *models.JSONB `json:"config" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	test := models.StressTestConfig{
		TestName: input.TestName,
		TestType: input.TestType,
		Config:   *input.Config,
		Status:   "draft",
	}

	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		test.CreatedBy = &uid
	}

	if err := sc.DB.Create(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建压力测试失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": test})
}

// GetStressTest 获取压力测试详情
func (sc *SimulationController) GetStressTest(c *gin.Context) {
	id := c.Param("id")

	var test models.StressTestConfig
	if err := sc.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "压力测试不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取压力测试失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": test})
}

// RunStressTest 运行压力测试
func (sc *SimulationController) RunStressTest(c *gin.Context) {
	id := c.Param("id")

	var test models.StressTestConfig
	if err := sc.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "压力测试不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取压力测试失败", "error": err.Error()})
		return
	}

	if test.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "压力测试运行中"})
		return
	}

	now := time.Now()
	if err := sc.DB.Model(&test).Updates(map[string]interface{}{
		"status":     "running",
		"start_time": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "启动压力测试失败", "error": err.Error()})
		return
	}

	// TODO: 实际的压力测试引擎会在后台运行，这里先模拟启动
	// 在实际实现中，应该启动一个 goroutine 来执行压力测试

	test.Status = "running"
	test.StartTime = &now

	c.JSON(http.StatusAccepted, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"test_id":    test.ID,
			"status":     "running",
			"started_at": now,
		},
	})
}

// StopStressTest 停止压力测试
func (sc *SimulationController) StopStressTest(c *gin.Context) {
	id := c.Param("id")

	var test models.StressTestConfig
	if err := sc.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "压力测试不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取压力测试失败", "error": err.Error()})
		return
	}

	if test.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 4002, "message": "压力测试未运行"})
		return
	}

	now := time.Now()
	var durationSeconds int
	if test.StartTime != nil {
		durationSeconds = int(now.Sub(*test.StartTime).Seconds())
	}

	// 模拟生成测试报告
	summary := models.JSONB{
		"total_requests":    10000,
		"failed_requests":  50,
		"success_rate":      99.5,
		"avg_response_time": 85.0,
		"p50_response_time": 72.0,
		"p95_response_time": 150.0,
		"p99_response_time": 200.0,
		"max_response_time": 500.0,
		"requests_per_second": 120.0,
	}

	if err := sc.DB.Model(&test).Updates(map[string]interface{}{
		"status":            "completed",
		"end_time":          now,
		"duration_seconds":  durationSeconds,
		"summary":           summary,
		"thresholds_passed": true,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "停止压力测试失败", "error": err.Error()})
		return
	}

	test.Status = "completed"
	test.EndTime = &now
	test.DurationSeconds = durationSeconds
	test.Summary = summary

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"test_id":         test.ID,
			"status":          "completed",
			"duration_seconds": durationSeconds,
			"summary":         summary,
		},
	})
}

// GetStressTestReport 获取压力测试报告
func (sc *SimulationController) GetStressTestReport(c *gin.Context) {
	id := c.Param("id")

	var test models.StressTestConfig
	if err := sc.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "压力测试不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取压力测试失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"id":               test.ID,
			"test_name":        test.TestName,
			"test_type":        test.TestType,
			"status":           test.Status,
			"duration_seconds": test.DurationSeconds,
			"summary":          test.Summary,
			"metrics":          test.Metrics,
			"thresholds_passed": test.ThresholdsPassed,
			"generated_at":      time.Now(),
		},
	})
}

// ============================================================
// A/B实验仿真 API
// ============================================================

// CreateABExperiment 创建A/B实验
func (sc *SimulationController) CreateABExperiment(c *gin.Context) {
	var input struct {
		ExperimentName string         `json:"experiment_name" binding:"required"`
		ExperimentKey  string         `json:"experiment_key" binding:"required"`
		Description    string         `json:"description"`
		Hypothesis     string         `json:"hypothesis"`
		TrafficPercent int            `json:"traffic_percent"`
		VariantAConfig  *models.JSONB `json:"variant_a_config"`
		VariantBConfig  *models.JSONB `json:"variant_b_config"`
		TargetMetrics   []string      `json:"target_metrics"`
		Tags           []string       `json:"tags"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	experiment := models.ABExperiment{
		ExperimentName: input.ExperimentName,
		ExperimentKey:  input.ExperimentKey,
		Description:   input.Description,
		Hypothesis:    input.Hypothesis,
		TrafficPercent: 100,
		Status:        "draft",
		VariantAConfig: func() models.JSONB {
			if input.VariantAConfig != nil {
				return *input.VariantAConfig
			}
			return models.JSONB{}
		}(),
		VariantBConfig: func() models.JSONB {
			if input.VariantBConfig != nil {
				return *input.VariantBConfig
			}
			return models.JSONB{}
		}(),
		TargetMetrics: input.TargetMetrics,
		Tags:          input.Tags,
	}

	if input.TrafficPercent > 0 && input.TrafficPercent <= 100 {
		experiment.TrafficPercent = input.TrafficPercent
	}

	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		experiment.CreatedBy = &uid
	}

	if err := sc.DB.Create(&experiment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建实验失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": experiment})
}

// GetABExperiment 获取实验详情
func (sc *SimulationController) GetABExperiment(c *gin.Context) {
	id := c.Param("id")

	var experiment models.ABExperiment
	if err := sc.DB.First(&experiment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "实验不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取实验失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": experiment})
}

// RunABExperiment 运行实验
func (sc *SimulationController) RunABExperiment(c *gin.Context) {
	id := c.Param("id")

	var experiment models.ABExperiment
	if err := sc.DB.First(&experiment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "实验不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取实验失败", "error": err.Error()})
		return
	}

	now := time.Now()
	if err := sc.DB.Model(&experiment).Updates(map[string]interface{}{
		"status":     "running",
		"start_time": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "启动实验失败", "error": err.Error()})
		return
	}

	experiment.Status = "running"
	experiment.StartTime = &now

	c.JSON(http.StatusAccepted, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"experiment_id": experiment.ID,
			"status":         "running",
			"started_at":     now,
		},
	})
}

// GetABExperimentResults 获取实验结果
func (sc *SimulationController) GetABExperimentResults(c *gin.Context) {
	id := c.Param("id")

	var experiment models.ABExperiment
	if err := sc.DB.First(&experiment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "实验不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取实验失败", "error": err.Error()})
		return
	}

	var results []models.ABExperimentResult
	if err := sc.DB.Where("experiment_id = ?", id).Order("created_at DESC").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取实验结果失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"experiment":  experiment,
			"results":     results,
			"result_summary": experiment.ResultSummary,
		},
	})
}

// ============================================================
// 回放系统 API
// ============================================================

// CreateRecording 创建录制
func (sc *SimulationController) CreateRecording(c *gin.Context) {
	var input struct {
		DeviceID   string          `json:"device_id"`
		PetID      *uint           `json:"pet_id"`
		RecordType string          `json:"record_type" binding:"required"`
		Metadata   *models.JSONB   `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1002, "message": "参数无效", "error": err.Error()})
		return
	}

	recording := models.Recording{
		DeviceID:   input.DeviceID,
		PetID:      input.PetID,
		RecordType: input.RecordType,
		StartTime:  time.Now(),
		Status:     "recording",
		Metadata: func() models.JSONB {
			if input.Metadata != nil {
				return *input.Metadata
			}
			return models.JSONB{}
		}(),
	}

	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		recording.CreatedBy = &uid
	}

	if err := sc.DB.Create(&recording).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "创建录制失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": recording})
}

// GetRecording 获取录制详情
func (sc *SimulationController) GetRecording(c *gin.Context) {
	id := c.Param("id")

	var recording models.Recording
	if err := sc.DB.First(&recording, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "录制不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取录制失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": recording})
}

// ListRecordings 获取录制列表
func (sc *SimulationController) ListRecordings(c *gin.Context) {
	var recordings []models.Recording
	query := sc.DB.Model(&models.Recording{})

	// 筛选参数
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}
	if recordType := c.Query("record_type"); recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&recordings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取录制列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"items":     recordings,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ReplayRecording 回放录制
func (sc *SimulationController) ReplayRecording(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		StartPositionMs int64           `json:"start_position_ms"`
		Speed          float64         `json:"speed"`
		Variables      *models.JSONB   `json:"variables"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// 使用默认值
		input = struct {
			StartPositionMs int64           `json:"start_position_ms"`
			Speed          float64         `json:"speed"`
			Variables      *models.JSONB   `json:"variables"`
		}{Speed: 1.0}
	}

	var recording models.Recording
	if err := sc.DB.First(&recording, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "录制不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取录制失败", "error": err.Error()})
		return
	}

	if recording.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "录制未完成，无法回放"})
		return
	}

	// TODO: 实际回放逻辑会在后台运行
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"recording_id":      recording.ID,
			"status":            "playing",
			"current_position_ms": input.StartPositionMs,
			"speed":             input.Speed,
		},
	})
}

// StopRecording 停止录制/回放
func (sc *SimulationController) StopRecording(c *gin.Context) {
	id := c.Param("id")

	var recording models.Recording
	if err := sc.DB.First(&recording, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 1001, "message": "录制不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "获取录制失败", "error": err.Error()})
		return
	}

	now := time.Now()
	var durationMs int
	if recording.StartTime.Unix() > 0 {
		durationMs = int(now.Sub(recording.StartTime).Milliseconds())
	}

	updates := map[string]interface{}{
		"status":     "completed",
		"end_time":   now,
		"duration_ms": durationMs,
	}

	if err := sc.DB.Model(&recording).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "停止录制失败", "error": err.Error()})
		return
	}

	recording.Status = "completed"
	recording.EndTime = &now
	recording.DurationMs = durationMs

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"recording_id": recording.ID,
			"status":       "completed",
			"duration_ms":  durationMs,
		},
	})
}