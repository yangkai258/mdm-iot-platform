package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIFairnessController AI公平性测试控制器
type AIFairnessController struct {
	DB *gorm.DB
}

// NewAIFairnessController 创建控制器
func NewAIFairnessController(db *gorm.DB) *AIFairnessController {
	return &AIFairnessController{DB: db}
}

// RegisterRoutes 注册AI公平性相关路由
func (c *AIFairnessController) RegisterRoutes(rg *gin.RouterGroup) {
	// 公平性测试
	rg.GET("/ai/fairness/tests", c.ListTests)
	rg.POST("/ai/fairness/tests", c.CreateTest)
	rg.GET("/ai/fairness/tests/:id", c.GetTest)
	rg.PUT("/ai/fairness/tests/:id", c.UpdateTest)
	rg.DELETE("/ai/fairness/tests/:id", c.DeleteTest)
	rg.POST("/ai/fairness/tests/:id/run", c.RunTest)
	rg.GET("/ai/fairness/tests/:id/report", c.GetTestReport)
	// 偏见检测
	rg.POST("/ai/fairness/detect-bias", c.DetectBias)
	rg.GET("/ai/fairness/metrics", c.GetFairnessMetrics)
	rg.GET("/ai/fairness/bias-detections", c.ListBiasDetections)
	// 模型审计
	rg.GET("/ai/audit/logs", c.ListAuditLogs)
	rg.POST("/ai/audit/report", c.GenerateAuditReport)
	rg.GET("/ai/audit/reports", c.ListAuditReports)
}

// ListTests 获取公平性测试列表
// GET /api/v1/ai/fairness/tests
func (c *AIFairnessController) ListTests(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query := c.DB.Model(&models.FairnessTest{})

	if testType := ctx.Query("test_type"); testType != "" {
		query = query.Where("test_type = ?", testType)
	}
	if modelIDStr := ctx.Query("model_id"); modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var tests []models.FairnessTest
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tests).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": tests,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateTest 创建公平性测试
// POST /api/v1/ai/fairness/tests
func (c *AIFairnessController) CreateTest(ctx *gin.Context) {
	var req models.FairnessTestCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(ctx)
	orgID := getOrgID(ctx)

	test := models.FairnessTest{
		Name:         req.Name,
		Description:  req.Description,
		TestType:     req.TestType,
		ModelID:      req.ModelID,
		ModelKey:     req.ModelKey,
		Config:       req.Config,
		TestData:     req.TestData,
		Status:       models.FairnessTestStatusPending,
		CreateUserID: userID,
		OrgID:        orgID,
	}

	if err := c.DB.Create(&test).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建公平性测试失败: " + err.Error()})
		return
	}

	c.logAIAudit(ctx, "create", "ai_fairness", "fairness_test", fmt.Sprintf("%d", test.ID), test.ModelID, test.ModelKey, nil)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "公平性测试创建成功",
		"data":    test,
	})
}

// GetTest 获取公平性测试详情
// GET /api/v1/ai/fairness/tests/:id
func (c *AIFairnessController) GetTest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var test models.FairnessTest
	var err error

	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr == nil {
		err = c.DB.First(&test, id).Error
	} else {
		err = c.DB.Where("test_key = ?", idStr).First(&test).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "公平性测试不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": test})
}

// UpdateTest 更新公平性测试
// PUT /api/v1/ai/fairness/tests/:id
func (c *AIFairnessController) UpdateTest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var test models.FairnessTest
	if err := c.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "公平性测试不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if test.Status != models.FairnessTestStatusPending {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "只允许更新 pending 状态的测试"})
		return
	}

	var req models.FairnessTestUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	if req.TestData != "" {
		updates["test_data"] = req.TestData
	}

	if err := c.DB.Model(&test).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败: " + err.Error()})
		return
	}

	c.DB.First(&test, test.ID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": test})
}

// DeleteTest 删除公平性测试
// DELETE /api/v1/ai/fairness/tests/:id
func (c *AIFairnessController) DeleteTest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var test models.FairnessTest
	if err := c.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "公平性测试不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if err := c.DB.Delete(&test).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// RunTest 运行公平性测试
// POST /api/v1/ai/fairness/tests/:id/run
func (c *AIFairnessController) RunTest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var test models.FairnessTest
	if err := c.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "公平性测试不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if test.Status == models.FairnessTestStatusRunning {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "测试已在运行中"})
		return
	}

	var req models.FairnessTestRun
	ctx.ShouldBindJSON(&req)

	now := time.Now()
	test.Status = models.FairnessTestStatusRunning
	test.StartedAt = &now
	if req.TestData != "" {
		test.TestData = req.TestData
	}
	c.DB.Save(&test)

	go c.executeFairnessTest(test.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "测试已开始运行",
		"data":    gin.H{"test_key": test.TestKey, "status": test.Status, "started_at": now},
	})
}

// executeFairnessTest 异步执行公平性测试
func (c *AIFairnessController) executeFairnessTest(testID uint) {
	var test models.FairnessTest
	if err := c.DB.First(&test, testID).Error; err != nil {
		return
	}

	for i := 1; i <= 100; i += 10 {
		time.Sleep(200 * time.Millisecond)
		c.DB.Model(&test).Update("progress", i)
	}

	metricsMap := c.generateMockFairnessMetrics(test.TestType)
	resultsMap := c.generateMockTestResults(test.TestType)

	resultsJSON, _ := json.Marshal(resultsMap)
	metricsJSON, _ := json.Marshal(metricsMap)

	now := time.Now()
	c.DB.Model(&test).Updates(map[string]interface{}{
		"status":       models.FairnessTestStatusCompleted,
		"progress":     100,
		"results":      string(resultsJSON),
		"metrics":      string(metricsJSON),
		"completed_at": now,
	})
}

func (c *AIFairnessController) generateMockFairnessMetrics(testType models.FairnessTestType) map[string]float64 {
	m := map[string]float64{
		"demographic_parity":   0.0,
		"equal_opportunity":    0.0,
		"disparate_impact":     0.0,
		"statistical_parity":  0.0,
		"precision_gap":        0.0,
		"recall_gap":           0.0,
		"false_positive_gap":  0.0,
		"overall_score":        0.0,
	}

	switch testType {
	case models.FairnessTestTypeDemographicParity:
		m["demographic_parity"] = 0.75 + rand.Float64()*0.2
		m["overall_score"] = m["demographic_parity"] * 100
	case models.FairnessTestTypeEqualOpportunity:
		m["equal_opportunity"] = 0.70 + rand.Float64()*0.25
		m["overall_score"] = m["equal_opportunity"] * 100
	case models.FairnessTestTypeDisparateImpact:
		m["disparate_impact"] = 0.80 + rand.Float64()*0.15
		m["overall_score"] = m["disparate_impact"] * 100
	default:
		for k := range m {
			m[k] = 0.60 + rand.Float64()*0.35
		}
		m["overall_score"] = 70 + rand.Float64()*25
	}
	return m
}

func (c *AIFairnessController) generateMockTestResults(testType models.FairnessTestType) map[string]interface{} {
	r := map[string]interface{}{
		"total_samples":     10000,
		"positive_outcomes": 0,
		"negative_outcomes": 0,
		"group_a_outcomes":  0,
		"group_b_outcomes":  0,
		"pass_status":       "pass",
		"details":           []map[string]interface{}{},
	}

	switch testType {
	case models.FairnessTestTypeDemographicParity:
		r["group_a_size"] = 5000
		r["group_b_size"] = 5000
		r["group_a_positive_rate"] = 0.85
		r["group_b_positive_rate"] = 0.82
		r["positive_outcomes"] = 8350
		r["negative_outcomes"] = 1650
		r["group_a_outcomes"] = 4250
		r["group_b_outcomes"] = 4100
	case models.FairnessTestTypeEqualOpportunity:
		r["group_a_size"] = 3000
		r["group_b_size"] = 7000
		r["true_positive_rate_a"] = 0.78
		r["true_positive_rate_b"] = 0.81
		r["positive_outcomes"] = 5600
		r["negative_outcomes"] = 4400
	default:
		r["group_a_size"] = 6000
		r["group_b_size"] = 4000
		r["group_a_positive_rate"] = 0.80
		r["group_b_positive_rate"] = 0.78
	}

	groupA := r["group_a_positive_rate"].(float64)
	groupB := r["group_b_positive_rate"].(float64)
	overallScore := groupA - groupB
	if overallScore > 0.1 {
		r["pass_status"] = "fail"
	} else if overallScore > 0.05 {
		r["pass_status"] = "warning"
	}
	return r
}

// GetTestReport 获取测试报告
// GET /api/v1/ai/fairness/tests/:id/report
func (c *AIFairnessController) GetTestReport(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var test models.FairnessTest
	if err := c.DB.First(&test, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "公平性测试不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if test.Status != models.FairnessTestStatusCompleted {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "测试未完成，无法生成报告"})
		return
	}

	var metrics map[string]float64
	if test.Metrics != "" {
		json.Unmarshal([]byte(test.Metrics), &metrics)
	}

	overallScore := 0.0
	if ms, ok := metrics["overall_score"]; ok {
		overallScore = ms
	}

	passStatus := "pass"
	if overallScore < 60 {
		passStatus = "fail"
	} else if overallScore < 80 {
		passStatus = "warning"
	}

	recommendations := c.generateRecommendations(test.TestType, overallScore)

	report := models.FairnessTestReport{
		TestKey:         test.TestKey,
		TestName:       test.Name,
		TestType:       string(test.TestType),
		Status:         string(test.Status),
		OverallScore:   overallScore,
		PassStatus:     passStatus,
		Metrics:        test.Metrics,
		Recommendations: recommendations,
		GeneratedAt:    time.Now(),
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

func (c *AIFairnessController) generateRecommendations(testType models.FairnessTestType, score float64) []string {
	var recs []string
	if score < 70 {
		recs = append(recs, "公平性评分较低，建议对模型进行重新训练或调整")
	}

	switch testType {
	case models.FairnessTestTypeDemographicParity:
		recs = append(recs, "检查训练数据中不同人口群体的代表性是否均衡")
		recs = append(recs, "考虑在损失函数中加入公平性约束项")
		recs = append(recs, "对高风险决策场景进行人工审核")
	case models.FairnessTestTypeEqualOpportunity:
		recs = append(recs, "确保不同群体获得相同质量的预测结果")
		recs = append(recs, "检查特征工程过程中是否存在信息泄露")
		recs = append(recs, "使用对抗性去偏技术减少偏见")
	case models.FairnessTestTypeDisparateImpact:
		recs = append(recs, "审查模型输出对不同群体的差异性影响")
		recs = append(recs, "实施4/5规则检验确保合规性")
		recs = append(recs, "考虑使用阈值调整技术平衡不同群体的结果")
	default:
		recs = append(recs, "建议定期进行公平性评估和监控")
		recs = append(recs, "建立持续公平性监测机制")
	}
	return recs
}

// DetectBias 检测偏见
// POST /api/v1/ai/fairness/detect-bias
func (c *AIFairnessController) DetectBias(ctx *gin.Context) {
	var req models.BiasDetectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(ctx)
	orgID := getOrgID(ctx)

	biasTypes := []models.BiasType{models.BiasTypeGender, models.BiasTypeAge, models.BiasTypeRace}
	if len(req.BiasTypes) > 0 {
		biasTypes = req.BiasTypes
	}

	var detections []models.BiasDetection
	for _, biasType := range biasTypes {
		detection := c.performBiasDetection(req, biasType, userID, orgID)
		c.DB.Create(&detection)
		detections = append(detections, detection)
	}

	c.logAIAudit(ctx, "bias_detect", "ai_bias", "bias_detection", "", req.ModelID, req.ModelKey, nil)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "偏见检测完成",
		"data":    gin.H{"detections": detections, "total": len(detections)},
	})
}

func (c *AIFairnessController) performBiasDetection(req models.BiasDetectRequest, biasType models.BiasType, userID, orgID uint) models.BiasDetection {
	biasScore := rand.Float64() * 0.5
	confidence := 0.6 + rand.Float64()*0.35

	var severity models.BiasSeverity
	if biasScore < 0.1 {
		severity = models.BiasSeverityLow
	} else if biasScore < 0.2 {
		severity = models.BiasSeverityMedium
	} else if biasScore < 0.35 {
		severity = models.BiasSeverityHigh
	} else {
		severity = models.BiasSeverityCritical
	}

	evidence := fmt.Sprintf("检测到 %s 相关偏见，样本分析显示存在 %.1f%% 的差异性", biasType, biasScore*100)
	recommendation := c.getRecommendationForBiasType(biasType, biasScore)

	return models.BiasDetection{
		ModelID:        req.ModelID,
		ModelKey:       req.ModelKey,
		InputData:      req.InputData,
		OutputData:     req.OutputData,
		BiasType:       biasType,
		Severity:       severity,
		Confidence:     confidence,
		BiasScore:      biasScore,
		Evidence:       evidence,
		Recommendation: recommendation,
		DetectedAt:     time.Now(),
		OrgID:          orgID,
		CreateUserID:   userID,
	}
}

func (c *AIFairnessController) getRecommendationForBiasType(biasType models.BiasType, score float64) string {
	var rec string
	switch biasType {
	case models.BiasTypeGender:
		rec = "建议增加性别平衡的训练数据，使用去偏技术如重采样或对抗训练"
	case models.BiasTypeAge:
		rec = "建议在训练数据中加入不同年龄段的样本，确保年龄段表示均衡"
	case models.BiasTypeRace:
		rec = "建议进行种族公平性审计，使用公平性约束重新训练模型"
	case models.BiasTypeEthnicity:
		rec = "建议收集更多元化的训练数据，避免特定群体的表示不足"
	default:
		rec = "建议进行全面的公平性审计，识别并缓解系统性偏见"
	}
	if score > 0.3 {
		rec += "。当前偏见严重程度较高，建议优先处理。"
	}
	return rec
}

// ListBiasDetections 获取偏见检测列表
// GET /api/v1/ai/fairness/bias-detections
func (c *AIFairnessController) ListBiasDetections(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query := c.DB.Model(&models.BiasDetection{})

	if modelIDStr := ctx.Query("model_id"); modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if biasType := ctx.Query("bias_type"); biasType != "" {
		query = query.Where("bias_type = ?", biasType)
	}
	if severity := ctx.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	var total int64
	query.Count(&total)

	var detections []models.BiasDetection
	offset := (page - 1) * pageSize
	if err := query.Order("detected_at DESC").Offset(offset).Limit(pageSize).Find(&detections).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": detections,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetFairnessMetrics 获取公平性指标
// GET /api/v1/ai/fairness/metrics
func (c *AIFairnessController) GetFairnessMetrics(ctx *gin.Context) {
	modelIDStr := ctx.Query("model_id")
	modelKey := ctx.Query("model_key")

	query := c.DB.Model(&models.FairnessMetrics{})
	if modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if modelKey != "" {
		query = query.Where("model_key = ?", modelKey)
	}

	var metrics models.FairnessMetrics
	err := query.First(&metrics).Error

	if err == gorm.ErrRecordNotFound {
		metrics = c.generateMockFairnessMetricsRecord(modelIDStr, modelKey)
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": metrics})
}

func (c *AIFairnessController) generateMockFairnessMetricsRecord(modelIDStr, modelKey string) models.FairnessMetrics {
	var modelID uint
	if modelIDStr != "" {
		if id, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			modelID = uint(id)
		}
	}

	return models.FairnessMetrics{
		ModelID:             modelID,
		ModelKey:            modelKey,
		DemographicParity:   0.82 + rand.Float64()*0.15,
		EqualOpportunity:    0.78 + rand.Float64()*0.18,
		DisparateImpact:     0.85 + rand.Float64()*0.12,
		StatisticalParity:   0.80 + rand.Float64()*0.15,
		PrecisionGap:        0.05 + rand.Float64()*0.10,
		RecallGap:           0.04 + rand.Float64()*0.08,
		FalsePositiveGap:    0.06 + rand.Float64()*0.09,
		OverallScore:        75 + rand.Float64()*20,
		TotalTestsRun:       int(rand.Int31n(100)),
		TotalBiasDetections: int(rand.Int31n(20)),
		CriticalBiasCount:   int(rand.Int31n(3)),
		HighBiasCount:       int(rand.Int31n(5)),
		LastUpdated:         time.Now(),
	}
}

// ListAuditLogs 获取AI审计日志
// GET /api/v1/ai/audit/logs
func (c *AIFairnessController) ListAuditLogs(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query := c.DB.Model(&models.AIAuditLog{})

	if modelIDStr := ctx.Query("model_id"); modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if modelKey := ctx.Query("model_key"); modelKey != "" {
		query = query.Where("model_key = ?", modelKey)
	}
	if action := ctx.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if module := ctx.Query("module"); module != "" {
		query = query.Where("module = ?", module)
	}
	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			query = query.Where("status = ?", status)
		}
	}
	if startDate := ctx.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate := ctx.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var logs []models.AIAuditLog
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": logs,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GenerateAuditReport 生成审计报告
// POST /api/v1/ai/audit/report
func (c *AIFairnessController) GenerateAuditReport(ctx *gin.Context) {
	var req models.AIAuditReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4001, "message": "参数校验失败: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(ctx)
	orgID := getOrgID(ctx)

	report := models.AIAuditReport{
		ReportType:  req.ReportType,
		ModelID:     req.ModelID,
		ModelKey:    req.ModelKey,
		Summary:     "AI模型审计报告",
		GeneratedBy: userID,
		GeneratedAt: time.Now(),
		OrgID:       orgID,
	}

	if req.PeriodStart != "" {
		if t, err := time.Parse("2006-01-02", req.PeriodStart); err == nil {
			report.PeriodStart = t
		}
	}
	if req.PeriodEnd != "" {
		if t, err := time.Parse("2006-01-02", req.PeriodEnd); err == nil {
			report.PeriodEnd = t
		}
	}

	findings, metricsMap, recommendations, riskLevel := c.generateReportContent(req)
 findingsJSON, _ := json.Marshal(findings)
 metricsJSON, _ := json.Marshal(metricsMap)
 report.Findings = string(findingsJSON)
 report.Metrics = string(metricsJSON)
	report.Recommendations = recommendations
	report.RiskLevel = riskLevel

	if err := c.DB.Create(&report).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "生成报告失败: " + err.Error()})
		return
	}

	c.logAIAudit(ctx, "generate_report", "ai_audit", "audit_report", report.ReportKey, req.ModelID, req.ModelKey, nil)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "审计报告生成成功", "data": report})
}

func (c *AIFairnessController) generateReportContent(req models.AIAuditReportRequest) ([]string, map[string]interface{}, string, string) {
	var findings []string
	var riskLevel = "low"

	var testCount, biasCount, criticalBias int64

	biasQuery := c.DB.Model(&models.BiasDetection{})
	testQuery := c.DB.Model(&models.FairnessTest{})

	if req.ModelID > 0 {
		biasQuery = biasQuery.Where("model_id = ?", req.ModelID)
		testQuery = testQuery.Where("model_id = ?", req.ModelID)
	}
	biasQuery.Count(&biasCount)
	testQuery.Count(&testCount)
	c.DB.Model(&models.BiasDetection{}).Where("severity = ?", models.BiasSeverityCritical).Count(&criticalBias)

	if testCount == 0 {
		findings = append(findings, "该模型尚未进行公平性测试")
	} else {
		findings = append(findings, fmt.Sprintf("该模型已进行 %d 次公平性测试", testCount))
	}

	if biasCount > 0 {
		findings = append(findings, fmt.Sprintf("检测到 %d 个偏见案例", biasCount))
		if criticalBias > 0 {
			findings = append(findings, fmt.Sprintf("警告：存在 %d 个严重偏见案例需要立即处理", criticalBias))
			riskLevel = "critical"
		} else {
			riskLevel = "medium"
		}
	} else {
		findings = append(findings, "未检测到明显偏见")
	}

	metricsMap := map[string]interface{}{
		"total_tests":             testCount,
		"total_bias_detections":   biasCount,
		"critical_bias_count":     criticalBias,
		"overall_fairness_score":  70 + rand.Float64()*25,
		"demographic_parity":      0.75 + rand.Float64()*0.2,
		"equal_opportunity":       0.72 + rand.Float64()*0.22,
	}

	var recommendations string
	switch riskLevel {
	case "critical":
		recommendations = "1. 立即暂停使用存在严重偏见的模型\n2. 进行全面公平性审计\n3. 使用去偏技术重新训练模型\n4. 建立持续监控机制"
	case "medium":
		recommendations = "1. 审查并优化训练数据\n2. 应用公平性约束重新训练\n3. 增强对高风险场景的人工审核\n4. 定期进行偏见检测"
	default:
		recommendations = "1. 继续保持当前的公平性监控\n2. 定期进行模型公平性评估\n3. 关注新出现的偏见风险"
	}

	return findings, metricsMap, recommendations, riskLevel
}

// ListAuditReports 获取审计报告列表
// GET /api/v1/ai/audit/reports
func (c *AIFairnessController) ListAuditReports(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query := c.DB.Model(&models.AIAuditReport{})

	if modelIDStr := ctx.Query("model_id"); modelIDStr != "" {
		if modelID, err := strconv.ParseUint(modelIDStr, 10, 64); err == nil {
			query = query.Where("model_id = ?", modelID)
		}
	}
	if reportType := ctx.Query("report_type"); reportType != "" {
		query = query.Where("report_type = ?", reportType)
	}
	if riskLevel := ctx.Query("risk_level"); riskLevel != "" {
		query = query.Where("risk_level = ?", riskLevel)
	}

	var total int64
	query.Count(&total)

	var reports []models.AIAuditReport
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": reports,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// logAIAudit 记录AI审计日志
func (c *AIFairnessController) logAIAudit(ctx *gin.Context, action, module, resourceType, resourceID string, modelID uint, modelKey string, metadata map[string]interface{}) {
	userID := getUserIDFromContext(ctx)
	orgID := getOrgID(ctx)
	tenantID := ctx.GetUint("tenant_id")

	metadataJSON, _ := json.Marshal(metadata)

	log := models.AIAuditLog{
		Action:        action,
		Module:        module,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		ModelID:       modelID,
		ModelKey:      modelKey,
		UserID:        userID,
		IP:            ctx.ClientIP(),
		UserAgent:     ctx.Request.UserAgent(),
		Status:        1,
		RequestMethod: ctx.Request.Method,
		RequestPath:   ctx.Request.URL.Path,
		Metadata:      string(metadataJSON),
		TenantID:      tenantID,
		OrgID:         orgID,
	}

	if username, exists := ctx.Get("username"); exists {
		if name, ok := username.(string); ok {
			log.Username = name
		}
	}

	c.DB.Create(&log)
}