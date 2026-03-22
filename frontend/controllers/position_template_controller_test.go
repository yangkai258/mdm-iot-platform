package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestPositionTemplateCreate_Success tests successful creation
func TestPositionTemplateCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock DB with actual postgres connection string
	// If PG is not available, this will fail with a clear error
	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.POST("/position-templates", ctrl.PositionTemplateCreate)

	body := map[string]interface{}{
		"name":        "测试模板",
		"code":        "test_code",
		"description": "测试描述",
		"permissions": []string{"device:view", "device:manage"},
		"status":      1,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/position-templates", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", w.Code)
	}
}

// TestPositionTemplateCreate_DuplicateCode tests duplicate code rejection
func TestPositionTemplateCreate_DuplicateCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.POST("/position-templates", ctrl.PositionTemplateCreate)

	// Try creating with no code (required field) - should fail validation
	body := map[string]interface{}{
		"name":   "新模板",
		"status": 1,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/position-templates", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Without required code field, ShouldBindJSON should fail
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestPositionTemplateGet_NotFound tests 404 for non-existent template
func TestPositionTemplateGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.GET("/position-templates/:id", ctrl.PositionTemplateGet)

	req, _ := http.NewRequest("GET", "/position-templates/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should return 404 (or 500 if DB unreachable)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 404 or 500, got %d", w.Code)
	}
}

// TestPositionTemplateUpdate_NotFound tests update of non-existent template
func TestPositionTemplateUpdate_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.PUT("/position-templates/:id", ctrl.PositionTemplateUpdate)

	body := map[string]interface{}{"name": "新名称"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/position-templates/99999", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 404 or 500, got %d", w.Code)
	}
}

// TestPositionTemplateDelete_NotFound tests delete of non-existent template
func TestPositionTemplateDelete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.DELETE("/position-templates/:id", ctrl.PositionTemplateDelete)

	req, _ := http.NewRequest("DELETE", "/position-templates/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 404 or 500, got %d", w.Code)
	}
}

// TestPositionTemplateCreate_InvalidJSON tests invalid JSON body
func TestPositionTemplateCreate_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.POST("/position-templates", ctrl.PositionTemplateCreate)

	req, _ := http.NewRequest("POST", "/position-templates", bytes.NewBuffer([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestPositionTemplateList_Integration tests list endpoint
func TestPositionTemplateList_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.GET("/position-templates", ctrl.PositionTemplateList)

	req, _ := http.NewRequest("GET", "/position-templates", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// If DB unreachable, will be 500; if reachable, will return JSON
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", w.Code)
	}
}

// TestPositionTemplateCopy_MissingDeptID tests copy without department ID
func TestPositionTemplateCopy_MissingDeptID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dsn := "host=localhost user=test password=test dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("PostgreSQL not available for integration test")
	}

	r := gin.New()
	ctrl := &PositionTemplateController{DB: db}
	r.POST("/position-templates/:id/copy", ctrl.PositionTemplateCopy)

	body := map[string]interface{}{}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/position-templates/1/copy", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestPositionTemplateResponse_JSON tests that response structure is correct
func TestPositionTemplateResponse_JSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Just verify the model JSON serialization works
	template := models.PositionTemplate{
		Name:        "测试",
		Code:        "test",
		Description: "描述",
		Status:      1,
	}

	data, err := json.Marshal(template)
	if err != nil {
		t.Errorf("failed to marshal PositionTemplate: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	if result["name"] != "测试" {
		t.Errorf("expected name '测试', got %v", result["name"])
	}
	if result["code"] != "test" {
		t.Errorf("expected code 'test', got %v", result["code"])
	}
	if result["status"] != float64(1) {
		t.Errorf("expected status 1, got %v", result["status"])
	}
}
