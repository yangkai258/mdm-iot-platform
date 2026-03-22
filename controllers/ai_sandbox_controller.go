package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AISandboxController AI沙箱测试控制器
type AISandboxController struct {
	DB *gorm.DB
}

// NewAISandboxController 创建控制器
func NewAISandboxController(db *gorm.DB) *AISandboxController {
	return &AISandboxController{DB: db}
}

// CreateTest 创建沙箱测试任务
// POST /api/v1/ai/sandbox/test
func (c *AISandboxController) CreateTest(ctx *gin.Context) {
	var req models.SandboxTestCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": err.Error()})
		return
	}

	// 校验测试类型
	validTypes := map[string]bool{"unit": true, "integration": true, "stress": true}
	if !validTypes[req.TestType] {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4005, "message": "test_type 必须为 unit/integration/stress 之一"})
		return
	}

	var userID uint
	if uid, exists := ctx.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = id
		}
	}

	taskID := "test-" + uuid.New().String()
	test := models.SandboxTest{
		TaskID:     taskID,
		ModelID:    req.ModelID,
		TestDataID: req.TestDataID,
		TestType:   req.TestType,
		TestName:   req.TestName,
		TestCases:  req.TestCases,
		Status:     "pending",
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
	}

	if err := c.DB.Create(&test).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"task_id": taskID,
			"status":  "pending",
		},
	})
}

// GetTest 获取测试结果
// GET /api/v1/ai/sandbox/test/:task_id
func (c *AISandboxController) GetTest(ctx *gin.Context) {
	taskID := ctx.Param("task_id")

	var test models.SandboxTest
	if err := c.DB.Where("task_id = ?", taskID).First(&test).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "测试任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": test})
}

// ListTests 测试任务列表
// GET /api/v1/ai/sandbox/tests
func (c *AISandboxController) ListTests(ctx *gin.Context) {
	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)

	query := c.DB.Model(&models.SandboxTest{})

	if modelID := ctx.Query("model_id"); modelID != "" {
		query = query.Where("model_id = ?", modelID)
	}
	if testType := ctx.Query("test_type"); testType != "" {
		query = query.Where("test_type = ?", testType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var tests []models.SandboxTest
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

// CancelTest 取消测试
// POST /api/v1/ai/sandbox/test/:task_id/cancel
func (c *AISandboxController) CancelTest(ctx *gin.Context) {
	taskID := ctx.Param("task_id")

	var test models.SandboxTest
	if err := c.DB.Where("task_id = ?", taskID).First(&test).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "测试任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if test.Status != "pending" && test.Status != "running" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "当前状态不允许取消"})
		return
	}

	if err := c.DB.Model(&test).Updates(map[string]interface{}{
		"status":       "failed",
		"completed_at": time.Now(),
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	c.DB.First(&test, "task_id = ?", taskID)
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": test})
}

// RunTest 执行沙箱测试（触发异步执行）
// POST /api/v1/ai/sandbox/test/:task_id/run
func (c *AISandboxController) RunTest(ctx *gin.Context) {
	taskID := ctx.Param("task_id")

	var test models.SandboxTest
	if err := c.DB.Where("task_id = ?", taskID).First(&test).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 4004, "message": "测试任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	if test.Status != "pending" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 4003, "message": "只有 pending 状态的任务可以执行"})
		return
	}

	// 更新状态为 running
	now := time.Now()
	if err := c.DB.Model(&test).Updates(map[string]interface{}{
		"status":     "running",
		"started_at": now,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": err.Error()})
		return
	}

	// 异步执行测试（实际沙箱逻辑在 sandbox/test_runner.go 中）
	go c.runSandboxTest(taskID)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"task_id": taskID,
			"status":  "running",
		},
	})
}

// runSandboxTest 异步执行沙箱测试
func (c *AISandboxController) runSandboxTest(taskID string) {
	var test models.SandboxTest
	if err := c.DB.Where("task_id = ?", taskID).First(&test).Error; err != nil {
		return
	}

	// 模拟测试执行（实际应调用 sandbox/test_runner.go）
	// 这里生成一个模拟结果，实际项目中应调用真实的沙箱测试逻辑
	result := `{"passed":true,"total":10,"passed_count":9,"failed_count":1,"avg_latency_ms":45,"max_latency_ms":120,"memory_mb":256}`

	completedAt := time.Now()
	c.DB.Model(&test).Updates(map[string]interface{}{
		"status":       "completed",
		"result":       result,
		"completed_at": completedAt,
	})
}
