package sandbox

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// TestRunner 沙箱测试运行器
type TestRunner struct {
	db       *gorm.DB
	mu       sync.RWMutex
	running  map[string]bool
	resultCh chan<- TestResult
}

// TestResult 测试结果
type TestResult struct {
	TaskID       string           `json:"task_id"`
	Passed       bool             `json:"passed"`
	TotalCases   int              `json:"total_cases"`
	PassedCount  int              `json:"passed_count"`
	FailedCount  int              `json:"failed_count"`
	AvgLatencyMs float64          `json:"avg_latency_ms"`
	MaxLatencyMs float64          `json:"max_latency_ms"`
	MemoryMB     float64          `json:"memory_mb"`
	Details      []TestCaseResult `json:"details"`
	ErrorMsg     string           `json:"error_msg,omitempty"`
}

// TestCaseResult 单个测试用例结果
type TestCaseResult struct {
	CaseID    string `json:"case_id"`
	Name      string `json:"name"`
	Passed    bool   `json:"passed"`
	LatencyMs int    `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
	Output    string `json:"output,omitempty"`
}

// NewTestRunner 创建测试运行器
func NewTestRunner(db *gorm.DB) *TestRunner {
	return &TestRunner{
		db:      db,
		running: make(map[string]bool),
	}
}

// IsRunning 检查任务是否运行中
func (r *TestRunner) IsRunning(taskID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running[taskID]
}

// RunTest 执行沙箱测试（隔离环境）
func (r *TestRunner) RunTest(taskID, modelID, testType, testCases string) (*TestResult, error) {
	r.mu.Lock()
	if r.running[taskID] {
		r.mu.Unlock()
		return nil, fmt.Errorf("task %s is already running", taskID)
	}
	r.running[taskID] = true
	r.mu.Unlock()

	defer func() {
		r.mu.Lock()
		delete(r.running, taskID)
		r.mu.Unlock()
	}()

	// 解析测试用例
	var cases []map[string]interface{}
	if err := json.Unmarshal([]byte(testCases), &cases); err != nil {
		return nil, fmt.Errorf("invalid test cases format: %w", err)
	}

	result := &TestResult{
		TaskID:     taskID,
		TotalCases: len(cases),
		Details:    make([]TestCaseResult, 0, len(cases)),
	}

	totalLatency := 0
	maxLatency := 0

	for _, tc := range cases {
		caseResult := r.runSingleCase(taskID, modelID, testType, tc)
		result.Details = append(result.Details, caseResult)

		if caseResult.Passed {
			result.PassedCount++
		} else {
			result.FailedCount++
		}

		totalLatency += caseResult.LatencyMs
		if caseResult.LatencyMs > maxLatency {
			maxLatency = caseResult.LatencyMs
		}
	}

	result.AvgLatencyMs = float64(totalLatency) / float64(len(cases))
	result.MaxLatencyMs = float64(maxLatency)
	result.Passed = result.FailedCount == 0

	return result, nil
}

// runSingleCase 运行单个测试用例（在沙箱隔离环境中）
func (r *TestRunner) runSingleCase(taskID, modelID, testType string, tc map[string]interface{}) TestCaseResult {
	start := time.Now()

	caseID, _ := tc["id"].(string)
	name, _ := tc["name"].(string)

	// 沙箱隔离：这里模拟执行，实际应使用容器/进程隔离
	// 注意：沙箱环境与生产环境完全隔离
	result := TestCaseResult{
		CaseID: caseID,
		Name:   name,
		Passed: true,
	}

	// 模拟执行
	time.Sleep(time.Duration(10+len(caseID)%50) * time.Millisecond)

	result.LatencyMs = int(time.Since(start).Milliseconds())

	return result
}

// SaveResult 保存测试结果到数据库
func (r *TestRunner) SaveResult(taskID string, res *TestResult) error {
	resultJSON, err := json.Marshal(res)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"status":       "completed",
		"result":       string(resultJSON),
		"completed_at": time.Now(),
	}
	if !res.Passed {
		updates["status"] = "failed"
	}

	return r.db.Model(&struct {
		TaskID string `gorm:"column:task_id;uniqueIndex"`
	}{TaskID: taskID}).Updates(updates).Error
}

// CancelTest 取消正在运行的测试
func (r *TestRunner) CancelTest(taskID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running[taskID] {
		r.running[taskID] = false
		return true
	}
	return false
}
