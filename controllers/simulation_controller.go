package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SimulationController 仿真测试平台控制器
type SimulationController struct {
	DB     *gorm.DB
	Service *services.SimulationService
}

// NewSimulationController 创建仿真控制器
func NewSimulationController(db *gorm.DB) *SimulationController {
	return &SimulationController{
		DB:     db,
		Service: services.NewSimulationService(db),
	}
}

// RegisterRoutes 注册仿真测试路由
func (c *SimulationController) RegisterRoutes(r *gin.RouterGroup) {
	// 虚拟宠物
	r.POST("/simulation/virtual-pets", c.CreateVirtualPet)
	r.GET("/simulation/virtual-pets", c.ListVirtualPets)
	r.GET("/simulation/virtual-pets/:id", c.GetVirtualPet)
	r.DELETE("/simulation/virtual-pets/:id", c.DeleteVirtualPet)
	r.POST("/simulation/virtual-pets/:id/behavior", c.SimulateBehavior)
	r.PUT("/simulation/virtual-pets/:id/attrs", c.UpdatePetAttrs)

	// 测试环境
	r.POST("/simulation/environments", c.CreateEnvironment)
	r.GET("/simulation/environments", c.ListEnvironments)
	r.GET("/simulation/environments/:id", c.GetEnvironment)
	r.PUT("/simulation/environments/:id/status", c.UpdateEnvStatus)

	// 测试运行
	r.POST("/simulation/run", c.CreateRun)
	r.POST("/simulation/run/:id/execute", c.ExecuteRun)
	r.GET("/simulation/run/:id", c.GetRun)
	r.GET("/simulation/run", c.ListRuns)

	// 性能指标
	r.GET("/simulation/metrics", c.GetMetrics)
	r.GET("/simulation/metrics/aggregate", c.GetAggregatedMetrics)
}

// getPageParams 解析分页参数
func getPageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

// ============ 虚拟宠物接口 ============

// CreateVirtualPet 创建虚拟宠物
// POST /api/v1/simulation/virtual-pets
func (c *SimulationController) CreateVirtualPet(ctx *gin.Context) {
	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Species     string                 `json:"species" binding:"required"`
		Personality string                 `json:"personality"`
		Mood        string                 `json:"mood"`
		Health      int                    `json:"health"`
		Hunger      int                    `json:"hunger"`
		Energy      int                    `json:"energy"`
		Happiness   int                    `json:"happiness"`
		Age         int                    `json:"age"`
		Weight      float64                `json:"weight"`
		AvatarURL   string                 `json:"avatar_url"`
		CustomAttrs map[string]interface{} `json:"custom_attrs"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request: " + err.Error()})
		return
	}

	pet := &models.VirtualPet{
		Name:         req.Name,
		Species:      req.Species,
		Personality:  req.Personality,
		Mood:         req.Mood,
		Health:       req.Health,
		Hunger:       req.Hunger,
		Energy:       req.Energy,
		Happiness:    req.Happiness,
		Age:          req.Age,
		Weight:       req.Weight,
		AvatarURL:    req.AvatarURL,
		OrgID:        getOrgID(ctx),
		CreateUserID: getUserID(ctx),
	}

	// 设置默认值
	if pet.Personality == "" {
		pet.Personality = "lively"
	}
	if pet.Mood == "" {
		pet.Mood = "happy"
	}
	if pet.Health == 0 {
		pet.Health = 100
	}
	if pet.Hunger == 0 {
		pet.Hunger = 80
	}
	if pet.Energy == 0 {
		pet.Energy = 100
	}
	if pet.Happiness == 0 {
		pet.Happiness = 80
	}

	if err := c.Service.CreateVirtualPet(pet); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create pet failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": pet})
}

// ListVirtualPets 获取虚拟宠物列表
// GET /api/v1/simulation/virtual-pets
func (c *SimulationController) ListVirtualPets(ctx *gin.Context) {
	page, pageSize := getPageParams(ctx)
	species := ctx.Query("species")
	personality := ctx.Query("personality")

	pets, total, err := c.Service.ListVirtualPets(getOrgID(ctx), page, pageSize, species, personality)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list pets failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      pets,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetVirtualPet 获取虚拟宠物详情
// GET /api/v1/simulation/virtual-pets/:id
func (c *SimulationController) GetVirtualPet(ctx *gin.Context) {
	petID := ctx.Param("id")

	pet, err := c.Service.GetVirtualPet(petID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "pet not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get pet failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": pet})
}

// DeleteVirtualPet 删除虚拟宠物
// DELETE /api/v1/simulation/virtual-pets/:id
func (c *SimulationController) DeleteVirtualPet(ctx *gin.Context) {
	petID := ctx.Param("id")

	if err := c.Service.DeleteVirtualPet(petID); err != nil {
		if err.Error() == "pet not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "pet not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete pet failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// SimulateBehavior 模拟宠物行为
// POST /api/v1/simulation/virtual-pets/:id/behavior
func (c *SimulationController) SimulateBehavior(ctx *gin.Context) {
	petID := ctx.Param("id")

	behavior, err := c.Service.SimulatePetBehavior(petID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "pet not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "simulate behavior failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": behavior})
}

// UpdatePetAttrs 更新宠物属性
// PUT /api/v1/simulation/virtual-pets/:id/attrs
func (c *SimulationController) UpdatePetAttrs(ctx *gin.Context) {
	petID := ctx.Param("id")

	var attrs map[string]interface{}
	if err := ctx.ShouldBindJSON(&attrs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request: " + err.Error()})
		return
	}

	if err := c.Service.UpdatePetAttributes(petID, attrs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update attrs failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "updated"})
}

// ============ 测试环境接口 ============

// CreateEnvironment 创建测试环境
// POST /api/v1/simulation/environments
func (c *SimulationController) CreateEnvironment(ctx *gin.Context) {
	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		SceneType   string                 `json:"scene_type" binding:"required"`
		SceneConfig map[string]interface{} `json:"scene_config"`
		Parameters  map[string]interface{} `json:"parameters"`
		Tags        string                 `json:"tags"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request: " + err.Error()})
		return
	}

	sceneConfigJSON := jsonMarshal(req.SceneConfig)
	paramsJSON := jsonMarshal(req.Parameters)

	env := &models.SimulationEnvironment{
		Name:         req.Name,
		Description:  req.Description,
		SceneType:    req.SceneType,
		SceneConfig:  sceneConfigJSON,
		Parameters:   paramsJSON,
		Tags:         req.Tags,
		Status:       "idle",
		OrgID:        getOrgID(ctx),
		CreateUserID: getUserID(ctx),
	}

	if err := c.Service.CreateEnvironment(env); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create environment failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": env})
}

// ListEnvironments 获取环境列表
// GET /api/v1/simulation/environments
func (c *SimulationController) ListEnvironments(ctx *gin.Context) {
	page, pageSize := getPageParams(ctx)
	sceneType := ctx.Query("scene_type")
	status := ctx.Query("status")

	envs, total, err := c.Service.ListEnvironments(getOrgID(ctx), page, pageSize, sceneType, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list environments failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      envs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetEnvironment 获取环境详情
// GET /api/v1/simulation/environments/:id
func (c *SimulationController) GetEnvironment(ctx *gin.Context) {
	envID := ctx.Param("id")

	env, err := c.Service.GetEnvironment(envID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "environment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get environment failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": env})
}

// UpdateEnvStatus 更新环境状态
// PUT /api/v1/simulation/environments/:id/status
func (c *SimulationController) UpdateEnvStatus(ctx *gin.Context) {
	envID := ctx.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{"idle": true, "running": true, "paused": true, "stopped": true}
	if !validStatuses[req.Status] {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid status value"})
		return
	}

	if err := c.Service.UpdateEnvironmentStatus(envID, req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update status failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "updated"})
}

// ============ 测试运行接口 ============

// CreateRun 创建测试运行
// POST /api/v1/simulation/run
func (c *SimulationController) CreateRun(ctx *gin.Context) {
	var req struct {
		Name           string                 `json:"name" binding:"required"`
		EnvID          string                 `json:"env_id"`
		PetID          string                 `json:"pet_id"`
		ScenarioConfig map[string]interface{} `json:"scenario_config"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request: " + err.Error()})
		return
	}

	scenarioJSON := jsonMarshal(req.ScenarioConfig)

	run := &models.SimulationRun{
		Name:           req.Name,
		EnvID:          req.EnvID,
		PetID:          req.PetID,
		ScenarioConfig: scenarioJSON,
		Status:         "pending",
		OrgID:          getOrgID(ctx),
		CreateUserID:   getUserID(ctx),
	}

	if err := c.Service.CreateRun(run); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create run failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": run})
}

// ExecuteRun 执行仿真测试
// POST /api/v1/simulation/run/:id/execute
func (c *SimulationController) ExecuteRun(ctx *gin.Context) {
	runID := ctx.Param("id")

	run, err := c.Service.RunSimulation(runID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "run not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "execute run failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": run})
}

// GetRun 获取测试运行详情
// GET /api/v1/simulation/run/:id
func (c *SimulationController) GetRun(ctx *gin.Context) {
	runID := ctx.Param("id")

	run, err := c.Service.GetRun(runID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "run not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get run failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": run})
}

// ListRuns 获取测试运行列表
// GET /api/v1/simulation/run
func (c *SimulationController) ListRuns(ctx *gin.Context) {
	page, pageSize := getPageParams(ctx)
	status := ctx.Query("status")
	petID := ctx.Query("pet_id")
	envID := ctx.Query("env_id")

	runs, total, err := c.Service.ListRuns(getOrgID(ctx), page, pageSize, status, petID, envID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list runs failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      runs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ============ 性能指标接口 ============

// GetMetrics 获取性能指标列表
// GET /api/v1/simulation/metrics
func (c *SimulationController) GetMetrics(ctx *gin.Context) {
	page, pageSize := getPageParams(ctx)
	metricType := ctx.Query("metric_type")
	petID := ctx.Query("pet_id")
	envID := ctx.Query("env_id")

	metrics, total, err := c.Service.GetMetrics(getOrgID(ctx), page, pageSize, metricType, petID, envID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get metrics failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      metrics,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAggregatedMetrics 获取聚合指标
// GET /api/v1/simulation/metrics/aggregate
func (c *SimulationController) GetAggregatedMetrics(ctx *gin.Context) {
	metricType := ctx.Query("metric_type")

	result, err := c.Service.GetAggregatedMetrics(getOrgID(ctx), metricType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "get aggregated metrics failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": result})
}

// jsonMarshal 封装 json.Marshal，错误返回空字符串
func jsonMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
