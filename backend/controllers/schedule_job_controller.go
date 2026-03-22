package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ScheduleJobController 调度任务控制器
type ScheduleJobController struct {
	DB *gorm.DB
}

func (c *ScheduleJobController) List(ctx *gin.Context) {
	var items []models.SysScheduleJob
	query := c.DB.Model(&models.SysScheduleJob{})

	if jobName := ctx.Query("job_name"); jobName != "" {
		query = query.Where("job_name LIKE ?", "%"+jobName+"%")
	}
	if jobType := ctx.Query("job_type"); jobType != "" {
		query = query.Where("job_type = ?", jobType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if enabled := ctx.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled)
	}

	var total int64
	query.Count(&total)

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	query.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&items)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": items, "pagination": gin.H{"total": total, "current": page, "pageSize": pageSize}}})
}

func (c *ScheduleJobController) Get(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysScheduleJob
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item})
}

func (c *ScheduleJobController) Create(ctx *gin.Context) {
	var req struct {
		JobName    string `json:"job_name" binding:"required"`
		JobType    string `json:"job_type" binding:"required"`
		CronExpr   string `json:"cron_expr"`
		ApiUrl     string `json:"api_url"`
		HttpMethod string `json:"http_method"`
		Headers    string `json:"headers"`
		BodyTpl    string `json:"body_tpl"`
		Enabled    int    `json:"enabled"`
		Remark     string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.JobType == "http" && req.ApiUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "HTTP类型任务必须填写API地址"})
		return
	}
	if req.HttpMethod == "" {
		req.HttpMethod = "POST"
	}
	if req.Enabled == 0 {
		req.Enabled = 1
	}

	item := models.SysScheduleJob{
		JobName: req.JobName, JobType: req.JobType, CronExpr: req.CronExpr,
		ApiUrl: req.ApiUrl, HttpMethod: req.HttpMethod, Headers: req.Headers,
		BodyTpl: req.BodyTpl, Enabled: req.Enabled, Remark: req.Remark, Status: "idle",
	}

	if err := c.DB.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "创建成功"})
}

func (c *ScheduleJobController) Update(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysScheduleJob
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var req struct {
		JobName    string `json:"job_name"`
		JobType    string `json:"job_type"`
		CronExpr   string `json:"cron_expr"`
		ApiUrl     string `json:"api_url"`
		HttpMethod string `json:"http_method"`
		Headers    string `json:"headers"`
		BodyTpl    string `json:"body_tpl"`
		Enabled    int    `json:"enabled"`
		Status     string `json:"status"`
		StatusMsg  string `json:"status_msg"`
		Remark     string `json:"remark"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.JobName != "" {
		item.JobName = req.JobName
	}
	if req.JobType != "" {
		item.JobType = req.JobType
	}
	if req.CronExpr != "" {
		item.CronExpr = req.CronExpr
	}
	if req.ApiUrl != "" {
		item.ApiUrl = req.ApiUrl
	}
	if req.HttpMethod != "" {
		item.HttpMethod = req.HttpMethod
	}
	if req.Headers != "" {
		item.Headers = req.Headers
	}
	if req.BodyTpl != "" {
		item.BodyTpl = req.BodyTpl
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.StatusMsg != "" {
		item.StatusMsg = req.StatusMsg
	}
	if req.Remark != "" {
		item.Remark = req.Remark
	}
	if req.Enabled != 0 {
		item.Enabled = req.Enabled
	}

	if err := c.DB.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": item, "message": "更新成功"})
}

func (c *ScheduleJobController) Delete(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var item models.SysScheduleJob
	if err := c.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	if err := c.DB.Delete(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

func (c *ScheduleJobController) Execute(ctx *gin.Context) {
	id := parseUint(ctx.Param("id"))
	var job models.SysScheduleJob
	if err := c.DB.First(&job, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.DB.Model(&job).Updates(map[string]interface{}{"status": "running", "status_msg": "执行中..."})

	go func() {
		result, err := executeJob(&job)
		status := "success"
		statusMsg := "执行成功"
		if err != nil {
			status = "failed"
			statusMsg = err.Error()
		}

		finishTime := time.Now()
		c.DB.Model(&models.SysScheduleJob{}).Where("id = ?", job.ID).Updates(map[string]interface{}{
			"status":      status,
			"status_msg":  statusMsg,
			"last_run_at": finishTime,
			"run_count":   gorm.Expr("run_count + 1"),
			"last_result": truncateStr(result, 2000),
		})
		log.Printf("[ScheduleJob] id=%d name=%s status=%s", job.ID, job.JobName, status)
	}()

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "任务已触发执行", "data": gin.H{"job_id": job.ID, "job_name": job.JobName, "triggered": true}})
}

func executeJob(job *models.SysScheduleJob) (string, error) {
	if job.JobType != "http" {
		return "", fmt.Errorf("unsupported job type: %s", job.JobType)
	}

	var headers map[string]string
	if job.Headers != "" {
		_ = json.Unmarshal([]byte(job.Headers), &headers)
	}

	var body io.Reader
	if job.BodyTpl != "" {
		body = bytes.NewBufferString(job.BodyTpl)
	}

	req, err := http.NewRequest(job.HttpMethod, job.ApiUrl, body)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	result := fmt.Sprintf("status=%d body=%s", resp.StatusCode, string(respBody))
	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return result, nil
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
