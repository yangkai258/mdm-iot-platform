package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"mdm-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SimulationService 仿真测试服务
type SimulationService struct {
	db *gorm.DB
}

// NewSimulationService 创建仿真服务实例
func NewSimulationService(db *gorm.DB) *SimulationService {
	return &SimulationService{db: db}
}

// ============ 虚拟宠物相关 ============

// CreateVirtualPet 创建虚拟宠物
func (s *SimulationService) CreateVirtualPet(pet *models.VirtualPet) error {
	if pet.PetID == "" {
		pet.PetID = uuid.New().String()
	}
	return s.db.Create(pet).Error
}

// ListVirtualPets 列出虚拟宠物
func (s *SimulationService) ListVirtualPets(orgID uint, page, pageSize int, species, personality string) ([]models.VirtualPet, int64, error) {
	var pets []models.VirtualPet
	var total int64

	query := s.db.Model(&models.VirtualPet{})
	if orgID > 0 {
		query = query.Where("org_id = ?", orgID)
	}
	if species != "" {
		query = query.Where("species = ?", species)
	}
	if personality != "" {
		query = query.Where("personality = ?", personality)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&pets).Error; err != nil {
		return nil, 0, err
	}

	return pets, total, nil
}

// GetVirtualPet 获取虚拟宠物详情
func (s *SimulationService) GetVirtualPet(petID string) (*models.VirtualPet, error) {
	var pet models.VirtualPet
	if err := s.db.Where("pet_id = ?", petID).First(&pet).Error; err != nil {
		return nil, err
	}
	return &pet, nil
}

// DeleteVirtualPet 删除虚拟宠物
func (s *SimulationService) DeleteVirtualPet(petID string) error {
	result := s.db.Where("pet_id = ?", petID).Delete(&models.VirtualPet{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("pet not found")
	}
	return nil
}

// SimulatePetBehavior 模拟宠物行为
func (s *SimulationService) SimulatePetBehavior(petID string) (map[string]interface{}, error) {
	pet, err := s.GetVirtualPet(petID)
	if err != nil {
		return nil, err
	}

	// 模拟行为变化
	behavior := s.generateBehavior(pet)
	return behavior, nil
}

func (s *SimulationService) generateBehavior(pet *models.VirtualPet) map[string]interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 根据性格生成不同行为倾向
	var action string
	var moodChange string

	switch pet.Personality {
	case "lively":
		actions := []string{"跑跳", "撒娇", "探索", "玩耍"}
		action = actions[r.Intn(len(actions))]
		moodChange = "excited"
	case "calm":
		actions := []string{"休息", "发呆", "晒太阳", "打盹"}
		action = actions[r.Intn(len(actions))]
		moodChange = "neutral"
	case "shy":
		actions := []string{"躲藏", "观望", "蹭脚", "依偎"}
		action = actions[r.Intn(len(actions))]
		moodChange = "neutral"
	case "curious":
		actions := []string{"探索", "嗅闻", "追逐", "研究"}
		action = actions[r.Intn(len(actions))]
		moodChange = "excited"
	default:
		actions := []string{"走动", "观望", "叫唤"}
		action = actions[r.Intn(len(actions))]
		moodChange = pet.Mood
	}

	// 随机属性变化
	deltaHealth := r.Intn(3) - 1 // -1 ~ +1
	deltaHunger := r.Intn(5) + 2 // +2 ~ +6
	deltaEnergy := r.Intn(7) - 3 // -3 ~ +3
	deltaHappiness := r.Intn(5) - 2 // -2 ~ +2

	newHealth := clamp(pet.Health+deltaHealth, 0, 100)
	newHunger := clamp(pet.Hunger+deltaHunger, 0, 100)
	newEnergy := clamp(pet.Energy+deltaEnergy, 0, 100)
	newHappiness := clamp(pet.Happiness+deltaHappiness, 0, 100)

	return map[string]interface{}{
		"pet_id":         pet.PetID,
		"action":         action,
		"mood_before":    pet.Mood,
		"mood_after":      moodChange,
		"health":         newHealth,
		"hunger":         newHunger,
		"energy":         newEnergy,
		"happiness":     newHappiness,
		"behavior_time":  time.Now().Format(time.RFC3339),
	}
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// UpdatePetAttributes 更新宠物属性
func (s *SimulationService) UpdatePetAttributes(petID string, attrs map[string]interface{}) error {
	updates := map[string]interface{}{}

	if v, ok := attrs["name"].(string); ok && v != "" {
		updates["name"] = v
	}
	if v, ok := attrs["personality"].(string); ok && v != "" {
		updates["personality"] = v
	}
	if v, ok := attrs["mood"].(string); ok && v != "" {
		updates["mood"] = v
	}
	if v, ok := attrs["health"].(float64); ok {
		updates["health"] = int(v)
	}
	if v, ok := attrs["hunger"].(float64); ok {
		updates["hunger"] = int(v)
	}
	if v, ok := attrs["energy"].(float64); ok {
		updates["energy"] = int(v)
	}
	if v, ok := attrs["happiness"].(float64); ok {
		updates["happiness"] = int(v)
	}
	if v, ok := attrs["avatar_url"].(string); ok {
		updates["avatar_url"] = v
	}

	if len(updates) == 0 {
		return errors.New("no valid fields to update")
	}

	return s.db.Model(&models.VirtualPet{}).Where("pet_id = ?", petID).Updates(updates).Error
}

// ============ 测试环境相关 ============

// CreateEnvironment 创建测试环境
func (s *SimulationService) CreateEnvironment(env *models.SimulationEnvironment) error {
	if env.EnvID == "" {
		env.EnvID = uuid.New().String()
	}
	return s.db.Create(env).Error
}

// ListEnvironments 列出测试环境
func (s *SimulationService) ListEnvironments(orgID uint, page, pageSize int, sceneType, status string) ([]models.SimulationEnvironment, int64, error) {
	var envs []models.SimulationEnvironment
	var total int64

	query := s.db.Model(&models.SimulationEnvironment{})
	if orgID > 0 {
		query = query.Where("org_id = ?", orgID)
	}
	if sceneType != "" {
		query = query.Where("scene_type = ?", sceneType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&envs).Error; err != nil {
		return nil, 0, err
	}

	return envs, total, nil
}

// GetEnvironment 获取环境详情
func (s *SimulationService) GetEnvironment(envID string) (*models.SimulationEnvironment, error) {
	var env models.SimulationEnvironment
	if err := s.db.Where("env_id = ?", envID).First(&env).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

// UpdateEnvironmentStatus 更新环境状态
func (s *SimulationService) UpdateEnvironmentStatus(envID, status string) error {
	return s.db.Model(&models.SimulationEnvironment{}).Where("env_id = ?", envID).Update("status", status).Error
}

// ============ 测试运行相关 ============

// CreateRun 创建测试运行
func (s *SimulationService) CreateRun(run *models.SimulationRun) error {
	if run.RunID == "" {
		run.RunID = uuid.New().String()
	}
	now := time.Now()
	run.StartedAt = &now
	run.Status = "running"
	return s.db.Create(run).Error
}

// RunSimulation 执行仿真测试
func (s *SimulationService) RunSimulation(runID string) (*models.SimulationRun, error) {
	var run models.SimulationRun
	if err := s.db.Where("run_id = ?", runID).First(&run).Error; err != nil {
		return nil, err
	}

	// 解析场景配置
	var scenarioConfig map[string]interface{}
	if run.ScenarioConfig != "" {
		json.Unmarshal([]byte(run.ScenarioConfig), &scenarioConfig)
	}

	// 模拟执行测试
	result := s.executeSimulation(&run, scenarioConfig)

	// 更新运行结果
	now := time.Now()
	run.CompletedAt = &now
	run.Status = "success"
	run.ResultData = result
	run.Duration = int(now.Sub(*run.StartedAt).Seconds())

	if err := s.db.Save(&run).Error; err != nil {
		return nil, err
	}

	// 更新环境状态
	if run.EnvID != "" {
		s.UpdateEnvironmentStatus(run.EnvID, "idle")
	}

	return &run, nil
}

func (s *SimulationService) executeSimulation(run *models.SimulationRun, config map[string]interface{}) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 模拟测试结果
	passRate := 0.85 + r.Float64()*0.15 // 85% ~ 100%
	responseTime := 50 + r.Intn(200)    // 50ms ~ 250ms
	throughput := 100 + r.Intn(400)     // 100 ~ 500 tps

	result := map[string]interface{}{
		"pass_rate":     fmt.Sprintf("%.2f%%", passRate*100),
		"response_time": fmt.Sprintf("%dms", responseTime),
		"throughput":    fmt.Sprintf("%d tps", throughput),
		"total_cases":   100,
		"passed_cases":  int(passRate * 100),
		"failed_cases":  100 - int(passRate*100),
		"warnings":      r.Intn(5),
		"simulated_at":  time.Now().Format(time.RFC3339),
	}

	// 如果配置了特定场景，添加相应指标
	if scenario, ok := config["scenario"].(string); ok {
		result["scenario"] = scenario
	}

	resultJSON, _ := json.Marshal(result)
	return string(resultJSON)
}

// GetRun 获取测试运行详情
func (s *SimulationService) GetRun(runID string) (*models.SimulationRun, error) {
	var run models.SimulationRun
	if err := s.db.Where("run_id = ?", runID).First(&run).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

// ListRuns 列出测试运行
func (s *SimulationService) ListRuns(orgID uint, page, pageSize int, status, petID, envID string) ([]models.SimulationRun, int64, error) {
	var runs []models.SimulationRun
	var total int64

	query := s.db.Model(&models.SimulationRun{})
	if orgID > 0 {
		query = query.Where("org_id = ?", orgID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if envID != "" {
		query = query.Where("env_id = ?", envID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&runs).Error; err != nil {
		return nil, 0, err
	}

	return runs, total, nil
}

// ============ 性能指标相关 ============

// RecordMetrics 记录性能指标
func (s *SimulationService) RecordMetrics(metric *models.SimulationMetrics) error {
	if metric.MetricID == "" {
		metric.MetricID = uuid.New().String()
	}
	return s.db.Create(metric).Error
}

// GetMetrics 获取性能指标
func (s *SimulationService) GetMetrics(orgID uint, page, pageSize int, metricType, petID, envID string) ([]models.SimulationMetrics, int64, error) {
	var metrics []models.SimulationMetrics
	var total int64

	query := s.db.Model(&models.SimulationMetrics{})
	if orgID > 0 {
		query = query.Where("org_id = ?", orgID)
	}
	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}
	if petID != "" {
		query = query.Where("pet_id = ?", petID)
	}
	if envID != "" {
		query = query.Where("env_id = ?", envID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&metrics).Error; err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}

// GetAggregatedMetrics 获取聚合指标
func (s *SimulationService) GetAggregatedMetrics(orgID uint, metricType string) (map[string]interface{}, error) {
	type AggResult struct {
		MetricType string
		AvgValue   float64
		MaxValue   float64
		MinValue   float64
		Count      int64
	}

	var result AggResult
	query := s.db.Model(&models.SimulationMetrics{}).Select(
		"metric_type, AVG(metric_value) as avg_value, MAX(metric_value) as max_value, MIN(metric_value) as min_value, COUNT(*) as count",
	)
	if orgID > 0 {
		query = query.Where("org_id = ?", orgID)
	}
	if metricType != "" {
		query = query.Where("metric_type = ?", metricType)
	}
	query = query.Group("metric_type").Scan(&result)

	return map[string]interface{}{
		"metric_type": result.MetricType,
		"avg_value":   result.AvgValue,
		"max_value":   result.MaxValue,
		"min_value":   result.MinValue,
		"total_count": result.Count,
		"generated_at": time.Now().Format(time.RFC3339),
	}, nil
}
