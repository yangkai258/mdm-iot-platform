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

type ResearchPlatformController struct {
	DB *gorm.DB
}

func (ctrl *ResearchPlatformController) RegisterRoutes(rg *gin.RouterGroup) {
	research := rg.Group("/research")
	{
		research.GET("/datasets", ctrl.GetDatasets)
		research.POST("/datasets", ctrl.CreateDataset)
		research.GET("/datasets/:id", ctrl.GetDataset)
		research.PUT("/datasets/:id", ctrl.UpdateDataset)
		research.DELETE("/datasets/:id", ctrl.DeleteDataset)
		research.GET("/datasets/:id/versions", ctrl.GetDatasetVersions)
		research.POST("/datasets/:id/versions", ctrl.CreateDatasetVersion)
		research.GET("/datasets/:id/download", ctrl.DownloadDataset)
		research.POST("/datasets/:id/cite", ctrl.CiteDataset)

		research.GET("/projects", ctrl.GetProjects)
		research.POST("/projects", ctrl.CreateProject)
		research.GET("/projects/:id", ctrl.GetProject)
		research.PUT("/projects/:id", ctrl.UpdateProject)
		research.DELETE("/projects/:id", ctrl.DeleteProject)
		research.POST("/projects/:id/submit", ctrl.SubmitProject)
		research.GET("/projects/:id/experiments", ctrl.GetProjectExperiments)

		research.GET("/experiments", ctrl.GetExperiments)
		research.POST("/experiments", ctrl.CreateExperiment)
		research.GET("/experiments/:id", ctrl.GetExperiment)
		research.PUT("/experiments/:id", ctrl.UpdateExperiment)
		research.POST("/experiments/:id/start", ctrl.StartExperiment)
		research.POST("/experiments/:id/stop", ctrl.StopExperiment)
		research.GET("/experiments/:id/results", ctrl.GetExperimentResults)

		research.GET("/collaborations", ctrl.GetCollaborations)
		research.POST("/collaborations", ctrl.CreateCollaboration)
		research.PUT("/collaborations/:id", ctrl.UpdateCollaboration)
		research.DELETE("/collaborations/:id", ctrl.DeleteCollaboration)
	}
}

func (ctrl *ResearchPlatformController) GetDatasets(c *gin.Context) {
	category := c.Query("category")
	accessLevel := c.Query("access_level")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var datasets []models.Dataset
	query := ctrl.DB.Model(&models.Dataset{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if accessLevel != "" {
		query = query.Where("access_level = ?", accessLevel)
	} else {
		query = query.Where("access_level = ?", "public")
	}

	var total int64
	query.Count(&total)

	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&datasets)

	c.JSON(http.StatusOK, gin.H{
		"datasets":  datasets,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *ResearchPlatformController) CreateDataset(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Category    string   `json:"category" binding:"required"`
		Tags        []string `json:"tags"`
		DataFormat  string   `json:"data_format" binding:"required"`
		DataSize    int64    `json:"data_size"`
		RecordCount int      `json:"record_count"`
		FileURL     string   `json:"file_url"`
		License     string   `json:"license"`
		AccessLevel string   `json:"access_level"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagsJSON, _ := json.Marshal(req.Tags)

	dataset := models.Dataset{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Tags:        string(tagsJSON),
		DataFormat:  req.DataFormat,
		DataSize:    req.DataSize,
		RecordCount: req.RecordCount,
		FileURL:     req.FileURL,
		License:     req.License,
		AccessLevel: req.AccessLevel,
	}

	ctrl.DB.Create(&dataset)
	c.JSON(http.StatusOK, dataset)
}

func (ctrl *ResearchPlatformController) GetDataset(c *gin.Context) {
	id := c.Param("id")
	var dataset models.Dataset

	if err := ctrl.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	c.JSON(http.StatusOK, dataset)
}

func (ctrl *ResearchPlatformController) UpdateDataset(c *gin.Context) {
	id := c.Param("id")
	var dataset models.Dataset

	if err := ctrl.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&dataset).Updates(req)
	c.JSON(http.StatusOK, dataset)
}

func (ctrl *ResearchPlatformController) DeleteDataset(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.Dataset{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *ResearchPlatformController) GetDatasetVersions(c *gin.Context) {
	id := c.Param("id")
	var versions []models.ResearchDatasetVersion

	ctrl.DB.Where("dataset_id = ?", id).Order("published_at DESC").Find(&versions)
	c.JSON(http.StatusOK, versions)
}

func (ctrl *ResearchPlatformController) CreateDatasetVersion(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Version     string `json:"version" binding:"required"`
		Changes    string `json:"changes"`
		FileURL    string `json:"file_url"`
		RecordCount int   `json:"record_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	version := models.ResearchDatasetVersion{
		DatasetID:   toUintResearch(id),
		Version:     req.Version,
		Changes:     req.Changes,
		FileURL:     req.FileURL,
		RecordCount: req.RecordCount,
	}

	ctrl.DB.Create(&version)
	c.JSON(http.StatusOK, version)
}

func (ctrl *ResearchPlatformController) DownloadDataset(c *gin.Context) {
	id := c.Param("id")
	var dataset models.Dataset

	if err := ctrl.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	ctrl.DB.Model(&dataset).UpdateColumn("download_count", gorm.Expr("download_count + 1"))

	c.JSON(http.StatusOK, gin.H{
		"download_url": dataset.FileURL,
		"name":        dataset.Name,
		"doi":         dataset.DOI,
	})
}

func (ctrl *ResearchPlatformController) CiteDataset(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.Dataset{}).Where("id = ?", id).UpdateColumn("citation_count", gorm.Expr("citation_count + 1"))
	c.JSON(http.StatusOK, gin.H{"message": "cited"})
}

func (ctrl *ResearchPlatformController) GetProjects(c *gin.Context) {
	ownerID := c.Query("owner_id")
	status := c.Query("status")
	var projects []models.ResearchProject
	query := ctrl.DB.Model(&models.ResearchProject{})

	if ownerID != "" {
		query = query.Where("owner_id = ?", ownerID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("created_at DESC").Find(&projects)
	c.JSON(http.StatusOK, projects)
}

func (ctrl *ResearchPlatformController) CreateProject(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		OwnerID     uint     `json:"owner_id" binding:"required"`
		DatasetIDs  []uint   `json:"dataset_ids"`
		StartDate   time.Time `json:"start_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	datasetIDsJSON, _ := json.Marshal(req.DatasetIDs)

	project := models.ResearchProject{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		DatasetIDs:  string(datasetIDsJSON),
		StartDate:   req.StartDate,
		Status:      "draft",
	}

	ctrl.DB.Create(&project)
	c.JSON(http.StatusOK, project)
}

func (ctrl *ResearchPlatformController) GetProject(c *gin.Context) {
	id := c.Param("id")
	var project models.ResearchProject

	if err := ctrl.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (ctrl *ResearchPlatformController) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.ResearchProject

	if err := ctrl.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&project).Updates(req)
	c.JSON(http.StatusOK, project)
}

func (ctrl *ResearchPlatformController) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.ResearchProject{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *ResearchPlatformController) SubmitProject(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Model(&models.ResearchProject{}).Where("id = ?", id).Update("status", "active")
	c.JSON(http.StatusOK, gin.H{"message": "submitted"})
}

func (ctrl *ResearchPlatformController) GetProjectExperiments(c *gin.Context) {
	id := c.Param("id")
	var experiments []models.ExperimentRun

	ctrl.DB.Where("project_id = ?", id).Order("started_at DESC").Find(&experiments)
	c.JSON(http.StatusOK, experiments)
}

func (ctrl *ResearchPlatformController) GetExperiments(c *gin.Context) {
	projectID := c.Query("project_id")
	status := c.Query("status")
	var experiments []models.ExperimentRun
	query := ctrl.DB.Model(&models.ExperimentRun{})

	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("started_at DESC").Find(&experiments)
	c.JSON(http.StatusOK, experiments)
}

func (ctrl *ResearchPlatformController) CreateExperiment(c *gin.Context) {
	var req struct {
		ProjectID uint                 `json:"project_id" binding:"required"`
		Name      string                `json:"name" binding:"required"`
		Config    map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configJSON, _ := json.Marshal(req.Config)

	exp := models.ExperimentRun{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Config:    string(configJSON),
		Status:    "draft",
	}

	ctrl.DB.Create(&exp)
	c.JSON(http.StatusOK, exp)
}

func (ctrl *ResearchPlatformController) GetExperiment(c *gin.Context) {
	id := c.Param("id")
	var exp models.ExperimentRun

	if err := ctrl.DB.First(&exp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	c.JSON(http.StatusOK, exp)
}

func (ctrl *ResearchPlatformController) UpdateExperiment(c *gin.Context) {
	id := c.Param("id")
	var exp models.ExperimentRun

	if err := ctrl.DB.First(&exp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Model(&exp).Updates(req)
	c.JSON(http.StatusOK, exp)
}

func (ctrl *ResearchPlatformController) StartExperiment(c *gin.Context) {
	id := c.Param("id")
	exp := models.ExperimentRun{
		Status:    "running",
		StartedAt: time.Now(),
	}

	ctrl.DB.Model(&models.ExperimentRun{}).Where("id = ?", id).Updates(exp)
	c.JSON(http.StatusOK, gin.H{"message": "started"})
}

func (ctrl *ResearchPlatformController) StopExperiment(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":       req.Status,
		"completed_at": &now,
	}

	ctrl.DB.Model(&models.ExperimentRun{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "stopped"})
}

func (ctrl *ResearchPlatformController) GetExperimentResults(c *gin.Context) {
	id := c.Param("id")
	var exp models.ExperimentRun

	if err := ctrl.DB.First(&exp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Experiment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": exp.Results,
		"metrics": exp.Metrics,
		"status":  exp.Status,
	})
}

func (ctrl *ResearchPlatformController) GetCollaborations(c *gin.Context) {
	projectID := c.Query("project_id")
	status := c.Query("status")
	var collabs []models.ResearchCollaboration
	query := ctrl.DB.Model(&models.ResearchCollaboration{})

	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&collabs)
	c.JSON(http.StatusOK, collabs)
}

func (ctrl *ResearchPlatformController) CreateCollaboration(c *gin.Context) {
	var req struct {
		ProjectID      uint   `json:"project_id" binding:"required"`
		CollaboratorID uint  `json:"collaborator_id" binding:"required"`
		Role          string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collab := models.ResearchCollaboration{
		ProjectID:      req.ProjectID,
		CollaboratorID: req.CollaboratorID,
		Role:          req.Role,
		Status:        "pending",
		InvitedAt:     time.Now(),
	}

	ctrl.DB.Create(&collab)
	c.JSON(http.StatusOK, collab)
}

func (ctrl *ResearchPlatformController) UpdateCollaboration(c *gin.Context) {
	id := c.Param("id")
	var collab models.ResearchCollaboration

	if err := ctrl.DB.First(&collab, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collaboration not found"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":       req.Status,
		"responded_at": &now,
	}

	ctrl.DB.Model(&collab).Updates(updates)
	c.JSON(http.StatusOK, collab)
}

func (ctrl *ResearchPlatformController) DeleteCollaboration(c *gin.Context) {
	id := c.Param("id")
	ctrl.DB.Delete(&models.ResearchCollaboration{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func toUintResearch(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 64)
	return uint(v)
}
