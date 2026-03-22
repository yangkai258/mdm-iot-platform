package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mdm-backend/models"
	"mdm-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FlowController 流程管理控制器
type FlowController struct {
	DB    *gorm.DB
	Engine *services.FlowEngine
}

// NewFlowController 创建流程控制器
func NewFlowController(db *gorm.DB) *FlowController {
	return &FlowController{
		DB:     db,
		Engine: services.NewFlowEngine(db),
	}
}

// CreateDefinitionRequest 创建流程定义请求
type CreateDefinitionRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description"`
	Nodes       []models.FlowNode    `json:"nodes" binding:"required,min=3"`
}

// CreateInstanceRequest 发起流程实例请求
type CreateInstanceRequest struct {
	FlowDefinitionID uint                  `json:"flow_definition_id" binding:"required"`
	BusinessKey      string                `json:"business_key"`
	BusinessType    string                `json:"business_type"`
	FormData        map[string]interface{} `json:"form_data"`
}

// ApproveTaskRequest 审批通过请求
type ApproveTaskRequest struct {
	Remark string `json:"remark"`
}

// RejectTaskRequest 审批拒绝请求
type RejectTaskRequest struct {
	Remark string `json:"remark" binding:"required"`
}

// CreateDefinition 创建流程定义
// POST /api/v1/flow/definitions
func (fc *FlowController) CreateDefinition(c *gin.Context) {
	var req CreateDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	tenantID := getTenantIDFromContext(c)
	def, err := fc.Engine.CreateFlowDefinition(req.Name, req.Description, req.Nodes, fmt.Sprintf("%d", tenantID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": def})
}

// ListDefinitions 获取流程定义列表
// GET /api/v1/flow/definitions
func (fc *FlowController) ListDefinitions(c *gin.Context) {
	tenantID := getTenantIDFromContext(c)
	defs, err := fc.Engine.GetFlowDefinitions(fmt.Sprintf("%d", tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": defs})
}

// StartInstance 发起流程实例
// POST /api/v1/flow/instances
func (fc *FlowController) StartInstance(c *gin.Context) {
	var req CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	var formData json.RawMessage
	if req.FormData != nil {
		b, _ := json.Marshal(req.FormData)
		formData = b
	}

	instance, err := fc.Engine.StartInstance(req.FlowDefinitionID, userID, fmt.Sprintf("%d", tenantID), req.BusinessKey, req.BusinessType, formData)
	if err != nil {
		if err == services.ErrFlowNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "流程定义不存在"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "发起失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "success", "data": instance})
}

// ListInstances 查询流程实例列表
// GET /api/v1/flow/instances
func (fc *FlowController) ListInstances(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	status := c.Query("status")
	
	query := fc.DB.Where("initiator_id = ? AND tenant_id = ?", userID, tenantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var instances []models.FlowInstance
	if err := query.Order("created_at DESC").Find(&instances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": instances})
}

// GetInstance 获取流程实例详情
// GET /api/v1/flow/instances/:id
func (fc *FlowController) GetInstance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的实例ID"})
		return
	}

	var instance models.FlowInstance
	if err := fc.DB.First(&instance, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实例不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	// 获取实例的任务列表
	tasks, _ := fc.Engine.GetInstanceTasks(instance.ID)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"instance": instance,
		"tasks":    tasks,
	}})
}

// ListPendingTasks 获取待审批任务列表
// GET /api/v1/flow/tasks/pending
func (fc *FlowController) ListPendingTasks(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tenantID := getTenantIDFromContext(c)

	tasks, err := fc.Engine.GetPendingTasks(userID, fmt.Sprintf("%d", tenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": tasks})
}

// ApproveTask 审批通过任务
// POST /api/v1/flow/tasks/:id/approve
func (fc *FlowController) ApproveTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的任务ID"})
		return
	}

	var req ApproveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// remark is optional for approval
		req.Remark = ""
	}

	userID := getUserIDFromContext(c)

	instance, err := fc.Engine.ApproveTask(uint(id), userID, req.Remark)
	if err != nil {
		switch err {
		case services.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务不存在"})
		case services.ErrNotCurrentApprover:
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "您不是当前审批人"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "审批成功", "data": instance})
}

// RejectTask 审批拒绝任务
// POST /api/v1/flow/tasks/:id/reject
func (fc *FlowController) RejectTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的任务ID"})
		return
	}

	var req RejectTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请提供审批意见"})
		return
	}

	userID := getUserIDFromContext(c)

	instance, err := fc.Engine.RejectTask(uint(id), userID, req.Remark)
	if err != nil {
		switch err {
		case services.ErrTaskNotFound:
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "任务不存在"})
		case services.ErrNotCurrentApprover:
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "您不是当前审批人"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已拒绝", "data": instance})
}

// CancelInstance 取消流程实例
// POST /api/v1/flow/instances/:id/cancel
func (fc *FlowController) CancelInstance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的实例ID"})
		return
	}

	userID := getUserIDFromContext(c)

	if err := fc.Engine.CancelInstance(uint(id), userID); err != nil {
		switch err {
		case services.ErrInstanceNotFound:
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "实例不存在"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "已取消"})
}

// RegisterFlowRoutes 注册流程相关路由
func (fc *FlowController) RegisterRoutes(r *gin.RouterGroup) {
	definitions := r.Group("/flow/definitions")
	{
		definitions.POST("", fc.CreateDefinition)
		definitions.GET("", fc.ListDefinitions)
	}

	instances := r.Group("/flow/instances")
	{
		instances.POST("", fc.StartInstance)
		instances.GET("", fc.ListInstances)
		instances.GET("/:id", fc.GetInstance)
		instances.POST("/:id/cancel", fc.CancelInstance)
	}

	tasks := r.Group("/flow/tasks")
	{
		tasks.GET("/pending", fc.ListPendingTasks)
		tasks.POST("/:id/approve", fc.ApproveTask)
		tasks.POST("/:id/reject", fc.RejectTask)
	}
}
