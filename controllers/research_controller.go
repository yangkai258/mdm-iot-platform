package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResearchController 行为研究数据控制器
type ResearchController struct {
	DB *gorm.DB
}

// RegisterRoutes 注册研究平台路由
func (r *ResearchController) RegisterRoutes(api *gin.RouterGroup) {
	// 行为研究数据
	api.GET("/research/data", r.ListAnonymizedData)
	api.POST("/research/data/export", r.ExportData)

	// 数据集管理
	api.GET("/research/datasets", r.ListDatasets)
	api.POST("/research/datasets", r.CreateDataset)
	api.GET("/research/datasets/:id", r.GetDataset)
	api.PUT("/research/datasets/:id", r.UpdateDataset)
	api.DELETE("/research/datasets/:id", r.DeleteDataset)

	// AI行为实验
	api.GET("/research/experiments", r.ListExperiments)
	api.POST("/research/experiments", r.CreateExperiment)
	api.GET("/research/experiments/:id", r.GetExperiment)
	api.PUT("/research/experiments/:id", r.UpdateExperiment)
	api.DELETE("/research/experiments/:id", r.DeleteExperiment)
	api.POST("/research/experiments/:id/start", r.StartExperiment)
	api.POST("/research/experiments/:id/stop", r.StopExperiment)
	api.GET("/research/experiments/:id/runs", r.ListExperimentRuns)
}

// ==================== 行为研究数据 API ====================

// ListAnonymizedData 获取匿名数据
// @Summary 获取匿名行为数据
// @Description 分页查询已匿名化的行为数据
// @Tags research
// @Accept json
// @Produce json
// @Param data_type query string false "数据类型: behavior_event, emotion_record, sensor_event"
// @Param device_type query string false "设备类型"
// @Param behavior_type query string false "行为类型"
// @Param date_from query string false "开始日期 YYYY-MM-DD"
// @Param date_to query string false "结束日期 YYYY-MM-DD"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/research/data [GET]
func (r *ResearchController) ListAnonymizedData(c *gin.Context) {
	dataType := c.Query("data_type")
	deviceType := c.Query("device_type")
	behaviorType := c.Query("behavior_type")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orgID, _ := c.Get("org_id")
	db := r.DB.Model(&models.ResearchDataRecord{}).Where("org_id = ?", orgID)

	if dataType != "" {
		db = db.Where("original_type = ?", dataType)
	}
	if deviceType != "" {
		db = db.Where("device_type = ?", deviceType)
	}
	if behaviorType != "" {
		db = db.Where("behavior_type = ?", behaviorType)
	}
	if dateFrom != "" {
		t, err := time.Parse("2006-01-02", dateFrom)
		if err == nil {
			db = db.Where("timestamp >= ?", t)
		}
	}
	if dateTo != "" {
		t, err := time.Parse("2006-01-02", dateTo)
		if err == nil {
			db = db.Where("timestamp <= ?", t.Add(24*time.Hour-time.Second))
		}
	}

	var total int64
	db.Count(&total)

	var records []models.ResearchDataRecord
	offset := (page - 1) * pageSize
	if err := db.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询数据失败"})
		return
	}

	// 解析 data_snapshot
	type DataSnapshot map[string]interface{}
	type ResponseRecord struct {
		models.ResearchDataRecord
		Data DataSnapshot `json:"data"`
	}
	var responses []ResponseRecord
	for _, rec := range records {
		var snap DataSnapshot
		json.Unmarshal([]byte(rec.DataSnapshot), &snap)
		responses = append(responses, ResponseRecord{
			ResearchDataRecord: rec,
			Data:               snap,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items":      responses,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ExportData 导出研究数据
// @Summary 导出研究数据
// @Description 创建异步导出任务
// @Tags research
// @Accept json
// @Produce json
// @Param body body models.ExportDataRequest true "导出请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/research/data/export [POST]
func (r *ResearchController) ExportData(c *gin.Context) {
	var req struct {
		DatasetID string `json:"dataset_id" binding:"required"`
		Format    string `json:"format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}
	if req.Format == "" {
		req.Format = "json"
	}
	if req.Format != "json" && req.Format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的导出格式，仅支持 json/csv"})
		return
	}

	orgID, _ := c.Get("org_id")
	userID, _ := c.Get("user_id")

	job := models.ResearchExportJob{
		DatasetID:   req.DatasetID,
		Status:      "pending",
		Format:      req.Format,
		OrgID:       orgID.(uint),
		CreateUserID: userID.(uint),
	}

	if err := r.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建导出任务失败"})
		return
	}

	// 异步执行导出（简化版：同步执行）
	go func(jobID uint) {
		var j models.ResearchExportJob
		r.DB.First(&j, jobID)
		j.Status = "processing"
		r.DB.Save(&j)

		// 查询数据集
		var ds models.ResearchDataset
		if err := r.DB.Where("dataset_id = ? AND org_id = ?", req.DatasetID, j.OrgID).First(&ds).Error; err != nil {
			j.Status = "failed"
			j.ErrorMsg = "数据集不存在"
			r.DB.Save(&j)
			return
		}

		// 查询数据记录
		var records []models.ResearchDataRecord
		query := r.DB.Where("dataset_id = ?", req.DatasetID)
		if ds.DateFrom != nil {
			query = query.Where("timestamp >= ?", *ds.DateFrom)
		}
		if ds.DateTo != nil {
			query = query.Where("timestamp <= ?", *ds.DateTo)
		}
		query.Find(&records)

		// 生成导出文件
		cwd, _ := os.Getwd()
		exportDir := filepath.Join(cwd, "exports", "research")
		os.MkdirAll(exportDir, 0755)
		filename := fmt.Sprintf("export_%s_%d.%s", req.DatasetID, time.Now().Unix(), req.Format)
		filePath := filepath.Join(exportDir, filename)

		if req.Format == "json" {
			data, _ := json.MarshalIndent(records, "", "  ")
			os.WriteFile(filePath, data, 0644)
		} else {
			// CSV
			var sb strings.Builder
			sb.WriteString("record_id,original_type,anonymized_id,behavior_type,timestamp,device_type\n")
			for _, rec := range records {
				sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s\n",
					rec.RecordID, rec.OriginalType, rec.AnonymizedID,
					rec.BehaviorType, rec.Timestamp.Format(time.RFC3339), rec.DeviceType))
			}
			os.WriteFile(filePath, []byte(sb.String()), 0644)
		}

		// 更新任务状态
		j.Status = "completed"
		j.FilePath = filePath
		j.RecordCount = len(records)
		if fi, err := os.Stat(filePath); err == nil {
			j.FileSize = fi.Size()
		}
		r.DB.Save(&j)
	}(job.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"job_id": job.JobID,
			"status": "pending",
		},
		"message": "导出任务已创建",
	})
}

// ==================== 数据集管理 API ====================

// ListDatasets 数据集列表
func (r *ResearchController) ListDatasets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	orgID, _ := c.Get("org_id")

	var total int64
	r.DB.Model(&models.ResearchDataset{}).Where("org_id = ?", orgID).Count(&total)

	var datasets []models.ResearchDataset
	offset := (page - 1) * pageSize
	r.DB.Where("org_id = ?", orgID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&datasets)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items":      datasets,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateDataset 创建数据集
func (r *ResearchController) CreateDataset(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		DataType    string `json:"data_type" binding:"required"`
		Source      string `json:"source"`
		DateFrom    string `json:"date_from"`
		DateTo      string `json:"date_to"`
		FileFormat  string `json:"file_format"`
		Anonymized  *bool  `json:"anonymized"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}
	orgID, _ := c.Get("org_id")
	userID, _ := c.Get("user_id")

	ds := models.ResearchDataset{
		Name:        req.Name,
		Description: req.Description,
		DataType:    req.DataType,
		Source:      req.Source,
		FileFormat:  req.FileFormat,
		Status:      "draft",
		OrgID:       orgID.(uint),
		CreateUserID: userID.(uint),
	}
	if req.Anonymized != nil {
		ds.Anonymized = *req.Anonymized
	}
	if req.DateFrom != "" {
		t, err := time.Parse("2006-01-02", req.DateFrom)
		if err == nil {
			ds.DateFrom = &t
		}
	}
	if req.DateTo != "" {
		t, err := time.Parse("2006-01-02", req.DateTo)
		if err == nil {
			ds.DateTo = &t
		}
	}

	if err := r.DB.Create(&ds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建数据集失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": ds, "message": "创建成功"})
}

// GetDataset 获取数据集详情
func (r *ResearchController) GetDataset(c *gin.Context) {
	datasetID := c.Param("id")
	orgID, _ := c.Get("org_id")

	var ds models.ResearchDataset
	if err := r.DB.Where("dataset_id = ? AND org_id = ?", datasetID, orgID).First(&ds).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "数据集不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": ds})
}

// UpdateDataset 更新数据集
func (r *ResearchController) UpdateDataset(c *gin.Context) {
	datasetID := c.Param("id")
	orgID, _ := c.Get("org_id")

	var ds models.ResearchDataset
	if err := r.DB.Where("dataset_id = ? AND org_id = ?", datasetID, orgID).First(&ds).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "数据集不存在"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		DateFrom    string `json:"date_from"`
		DateTo      string `json:"date_to"`
		FileFormat  string `json:"file_format"`
		Filters     string `json:"filters"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			updates["date_from"] = t
		}
	}
	if req.DateTo != "" {
		if t, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			updates["date_to"] = t
		}
	}
	if req.FileFormat != "" {
		updates["file_format"] = req.FileFormat
	}
	if req.Filters != "" {
		updates["filters"] = req.Filters
	}

	if len(updates) > 0 {
		r.DB.Model(&ds).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": ds, "message": "更新成功"})
}

// DeleteDataset 删除数据集
func (r *ResearchController) DeleteDataset(c *gin.Context) {
	datasetID := c.Param("id")
	orgID, _ := c.Get("org_id")

	if err := r.DB.Where("dataset_id = ? AND org_id = ?", datasetID, orgID).Delete(&models.ResearchDataset{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	// 同时删除关联的数据记录
	r.DB.Where("dataset_id = ?", datasetID).Delete(&models.ResearchDataRecord{})

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// ==================== AI 行为实验 API ====================

// ListExperiments 实验列表
func (r *ResearchController) ListExperiments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	orgID, _ := c.Get("org_id")

	db := r.DB.Model(&models.Experiment{}).Where("org_id = ?", orgID)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var experiments []models.Experiment
	offset := (page - 1) * pageSize
	db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&experiments)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"items":      experiments,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateExperiment 创建实验
func (r *ResearchController) CreateExperiment(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Hypothesis  string `json:"hypothesis"`
		TargetModel string `json:"target_model"`
		Variables   string `json:"variables"`
		ControlGroup string `json:"control_group"`
		TestGroup   string `json:"test_group"`
		SampleSize  int    `json:"sample_size"`
		Metrics     string `json:"metrics"`
		Tags        string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}
	orgID, _ := c.Get("org_id")
	userID, _ := c.Get("user_id")

	exp := models.Experiment{
		Name:          req.Name,
		Description:   req.Description,
		Hypothesis:    req.Hypothesis,
		TargetModel:   req.TargetModel,
		Variables:     req.Variables,
		ControlGroup:  req.ControlGroup,
		TestGroup:     req.TestGroup,
		SampleSize:    req.SampleSize,
		Metrics:       req.Metrics,
		Tags:          req.Tags,
		Status:        "draft",
		OrgID:         orgID.(uint),
		CreateUserID:  userID.(uint),
	}

	if err := r.DB.Create(&exp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建实验失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": exp, "message": "创建成功"})
}

// GetExperiment 实验详情
func (r *ResearchController) GetExperiment(c *gin.Context) {
	expID := c.Param("id")
	orgID, _ := c.Get("org_id")

	var exp models.Experiment
	if err := r.DB.Where("exp_id = ? AND org_id = ?", expID, orgID).First(&exp).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实验不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 查询关联的运行记录
	var runs []models.ExperimentRun
	r.DB.Where("exp_id = ?", expID).Order("created_at DESC").Limit(10).Find(&runs)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"experiment": exp,
			"runs":       runs,
		},
	})
}

// UpdateExperiment 更新实验
func (r *ResearchController) UpdateExperiment(c *gin.Context) {
	expID := c.Param("id")
	orgID, _ := c.Get("org_id")

	var exp models.Experiment
	if err := r.DB.Where("exp_id = ? AND org_id = ?", expID, orgID).First(&exp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实验不存在"})
		return
	}
	if exp.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "运行中的实验不能编辑"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Hypothesis  string `json:"hypothesis"`
		Variables   string `json:"variables"`
		Metrics     string `json:"metrics"`
		Tags        string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Hypothesis != "" {
		updates["hypothesis"] = req.Hypothesis
	}
	if req.Variables != "" {
		updates["variables"] = req.Variables
	}
	if req.Metrics != "" {
		updates["metrics"] = req.Metrics
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}

	if len(updates) > 0 {
		r.DB.Model(&exp).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": exp, "message": "更新成功"})
}

// DeleteExperiment 删除实验
func (r *ResearchController) DeleteExperiment(c *gin.Context) {
	expID := c.Param("id")
	orgID, _ := c.Get("org_id")

	var exp models.Experiment
	if err := r.DB.Where("exp_id = ? AND org_id = ?", expID, orgID).First(&exp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实验不存在"})
		return
	}
	if exp.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "运行中的实验不能删除"})
		return
	}

	r.DB.Where("exp_id = ?", expID).Delete(&models.ExperimentRun{})
	r.DB.Where("exp_id = ?", expID).Delete(&models.ExperimentParticipant{})
	r.DB.Delete(&exp)

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// StartExperiment 开始实验
func (r *ResearchController) StartExperiment(c *gin.Context) {
	expID := c.Param("id")
	orgID, _ := c.Get("org_id")
	userID, _ := c.Get("user_id")

	var exp models.Experiment
	if err := r.DB.Where("exp_id = ? AND org_id = ?", expID, orgID).First(&exp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实验不存在"})
		return
	}
	if exp.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "实验已在运行中"})
		return
	}
	if exp.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "已完成的实验不能重新开始，请创建新实验"})
		return
	}

	now := time.Now()
	exp.Status = "running"
	exp.StartTime = &now
	r.DB.Save(&exp)

	// 创建实验运行记录
	run := models.ExperimentRun{
		ExpID:         expID,
		Status:        "running",
		StartTime:     &now,
		OrgID:         orgID.(uint),
		CreateUserID:  userID.(uint),
	}
	r.DB.Create(&run)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"experiment": exp,
			"run_id":     run.RunID,
		},
		"message": "实验已启动",
	})
}

// StopExperiment 停止实验
func (r *ResearchController) StopExperiment(c *gin.Context) {
	expID := c.Param("id")
	orgID, _ := c.Get("org_id")

	var exp models.Experiment
	if err := r.DB.Where("exp_id = ? AND org_id = ?", expID, orgID).First(&exp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实验不存在"})
		return
	}
	if exp.Status != "running" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "实验不在运行中"})
		return
	}

	now := time.Now()
	exp.Status = "completed"
	exp.EndTime = &now
	if exp.StartTime != nil {
		exp.Duration = int(now.Sub(*exp.StartTime).Seconds())
	}
	r.DB.Save(&exp)

	// 更新运行记录
	var run models.ExperimentRun
	r.DB.Where("exp_id = ? AND status = ?", expID, "running").Order("created_at DESC").First(&run)
	if run.ID != 0 {
		run.Status = "completed"
		run.EndTime = &now
		if run.StartTime != nil {
			run.Duration = int(now.Sub(*run.StartTime).Seconds())
		}
		r.DB.Save(&run)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"experiment": exp,
			"run_id":     run.RunID,
		},
		"message": "实验已停止",
	})
}

// ListExperimentRuns 实验运行记录列表
func (r *ResearchController) ListExperimentRuns(c *gin.Context) {
	expID := c.Param("id")
	orgID, _ := c.Get("org_id")

	// 验证实验存在
	var exp models.Experiment
	if err := r.DB.Where("exp_id = ? AND org_id = ?", expID, orgID).First(&exp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实验不存在"})
		return
	}

	var runs []models.ExperimentRun
	r.DB.Where("exp_id = ?", expID).Order("created_at DESC").Find(&runs)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": runs})
}
