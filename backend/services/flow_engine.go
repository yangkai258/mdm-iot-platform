package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// FlowEngine 简单的顺序流程引擎
type FlowEngine struct {
	db *gorm.DB
}

// NewFlowEngine 创建流程引擎实例
func NewFlowEngine(db *gorm.DB) *FlowEngine {
	return &FlowEngine{db: db}
}

// ErrFlowNotFound 流程定义未找到
var ErrFlowNotFound = errors.New("flow definition not found")

// ErrInstanceNotFound 流程实例未找到
var ErrInstanceNotFound = errors.New("flow instance not found")

// ErrTaskNotFound 流程任务未找到
var ErrTaskNotFound = errors.New("flow task not found")

// ErrInvalidTransition 无效的流程转移
var ErrInvalidTransition = errors.New("invalid flow transition")

// ErrNotCurrentApprover 不是当前审批人
var ErrNotCurrentApprover = errors.New("user is not the current approver")

// StartInstance 发起流程实例
func (e *FlowEngine) StartInstance(flowDefID uint, initiatorID uint, tenantID string, businessKey string, businessType string, formData json.RawMessage) (*models.FlowInstance, error) {
	var flowDef models.FlowDefinition
	if err := e.db.First(&flowDef, flowDefID).Error; err != nil {
		return nil, ErrFlowNotFound
	}

	var nodes []models.FlowNode
	if err := json.Unmarshal(flowDef.Nodes, &nodes); err != nil {
		return nil, err
	}

	// 找到start节点
	var startNode *models.FlowNode
	for i := range nodes {
		if nodes[i].NodeType == models.FlowNodeStart {
			startNode = &nodes[i]
			break
		}
	}
	if startNode == nil {
		return nil, errors.New("flow has no start node")
	}

	// 找到start后的第一个节点（通常是审批节点）
	var firstApprovalNode *models.FlowNode
	if startNode.NextNodeID != "" {
		for i := range nodes {
			if nodes[i].NodeID == startNode.NextNodeID {
				firstApprovalNode = &nodes[i]
				break
			}
		}
	}

	instanceUUID := generateUUID()

	instance := &models.FlowInstance{
		InstanceUUID:     instanceUUID,
		FlowDefinitionID: flowDefID,
		FlowName:         flowDef.Name,
		BusinessKey:      businessKey,
		BusinessType:     businessType,
		Status:           models.FlowStatusRunning,
		CurrentNodeID:    firstApprovalNode.NodeID,
		InitiatorID:     initiatorID,
		TenantID:        tenantID,
		FormData:        formData,
	}

	tx := e.db.Begin()
	if err := tx.Create(instance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果第一个审批节点存在，则创建待审批任务
	if firstApprovalNode != nil && firstApprovalNode.ApproverID != nil {
		task := &models.FlowTask{
			InstanceID:   instance.ID,
			InstanceUUID: instanceUUID,
			NodeID:       firstApprovalNode.NodeID,
			NodeName:     firstApprovalNode.NodeName,
			ApproverID:   *firstApprovalNode.ApproverID,
			Status:       models.TaskStatusPending,
		}
		if err := tx.Create(task).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	log.Printf("[FlowEngine] Started instance %s, flow=%s, current_node=%s",
		instanceUUID, flowDef.Name, instance.CurrentNodeID)
	return instance, nil
}

// ApproveTask 审批通过任务
func (e *FlowEngine) ApproveTask(taskID uint, approverID uint, remark string) (*models.FlowInstance, error) {
	var task models.FlowTask
	if err := e.db.First(&task, taskID).Error; err != nil {
		return nil, ErrTaskNotFound
	}

	if task.ApproverID != approverID {
		return nil, ErrNotCurrentApprover
	}

	if task.Status != models.TaskStatusPending {
		return nil, errors.New("task is not pending")
	}

	var instance models.FlowInstance
	if err := e.db.First(&instance, task.InstanceID).Error; err != nil {
		return nil, ErrInstanceNotFound
	}

	if instance.Status != models.FlowStatusRunning {
		return nil, errors.New("instance is not running")
	}

	// 获取流程定义
	var flowDef models.FlowDefinition
	if err := e.db.First(&flowDef, instance.FlowDefinitionID).Error; err != nil {
		return nil, ErrFlowNotFound
	}

	var nodes []models.FlowNode
	if err := json.Unmarshal(flowDef.Nodes, &nodes); err != nil {
		return nil, err
	}

	// 找到当前节点
	var currentNode *models.FlowNode
	for i := range nodes {
		if nodes[i].NodeID == task.NodeID {
			currentNode = &nodes[i]
			break
		}
	}
	if currentNode == nil {
		return nil, ErrInvalidTransition
	}

	now := time.Now()
	tx := e.db.Begin()

	// 更新任务为已审批
	task.Status = models.TaskStatusApproved
	task.Remark = remark
	task.CompletedAt = &now
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 检查下一个节点
	var nextNode *models.FlowNode
	if currentNode.NextNodeID != "" {
		for i := range nodes {
			if nodes[i].NodeID == currentNode.NextNodeID {
				nextNode = &nodes[i]
				break
			}
		}
	}

	if nextNode == nil || nextNode.NodeType == models.FlowNodeEnd {
		// 流程结束
		instance.Status = models.FlowStatusApproved
		instance.CurrentNodeID = ""
	} else {
		// 流转到下一个审批节点
		instance.CurrentNodeID = nextNode.NodeID

		// 创建下一个任务（如果需要审批）
		if nextNode.ApproverID != nil {
			nextTask := &models.FlowTask{
				InstanceID:   instance.ID,
				InstanceUUID: instance.InstanceUUID,
				NodeID:       nextNode.NodeID,
				NodeName:     nextNode.NodeName,
				ApproverID:   *nextNode.ApproverID,
				Status:       models.TaskStatusPending,
			}
			if err := tx.Create(nextTask).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Save(&instance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	log.Printf("[FlowEngine] Approved task %d, instance %s -> %s",
		taskID, instance.InstanceUUID, instance.Status)
	return &instance, nil
}

// RejectTask 审批拒绝任务
func (e *FlowEngine) RejectTask(taskID uint, approverID uint, remark string) (*models.FlowInstance, error) {
	var task models.FlowTask
	if err := e.db.First(&task, taskID).Error; err != nil {
		return nil, ErrTaskNotFound
	}

	if task.ApproverID != approverID {
		return nil, ErrNotCurrentApprover
	}

	if task.Status != models.TaskStatusPending {
		return nil, errors.New("task is not pending")
	}

	var instance models.FlowInstance
	if err := e.db.First(&instance, task.InstanceID).Error; err != nil {
		return nil, ErrInstanceNotFound
	}

	if instance.Status != models.FlowStatusRunning {
		return nil, errors.New("instance is not running")
	}

	now := time.Now()
	tx := e.db.Begin()

	task.Status = models.TaskStatusRejected
	task.Remark = remark
	task.CompletedAt = &now
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	instance.Status = models.FlowStatusRejected
	instance.CurrentNodeID = ""
	if err := tx.Save(&instance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	log.Printf("[FlowEngine] Rejected task %d, instance %s", taskID, instance.InstanceUUID)
	return &instance, nil
}

// CancelInstance 取消流程实例（发起人可取消）
func (e *FlowEngine) CancelInstance(instanceID uint, userID uint) error {
	var instance models.FlowInstance
	if err := e.db.First(&instance, instanceID).Error; err != nil {
		return ErrInstanceNotFound
	}

	if instance.InitiatorID != userID {
		return errors.New("only initiator can cancel")
	}

	if instance.Status != models.FlowStatusRunning {
		return errors.New("instance is not running")
	}

	instance.Status = models.FlowStatusCanceled
	instance.CurrentNodeID = ""
	return e.db.Save(&instance).Error
}

// GetPendingTasks 获取用户的待审批任务列表
func (e *FlowEngine) GetPendingTasks(approverID uint, tenantID string) ([]models.FlowTask, error) {
	var tasks []models.FlowTask
	err := e.db.Where("approver_id = ? AND status = ? AND tenant_id = ?",
		approverID, models.TaskStatusPending, tenantID).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// GetMyInstances 获取我发起的流程实例
func (e *FlowEngine) GetMyInstances(initiatorID uint, tenantID string) ([]models.FlowInstance, error) {
	var instances []models.FlowInstance
	err := e.db.Where("initiator_id = ? AND tenant_id = ?", initiatorID, tenantID).
		Order("created_at DESC").
		Find(&instances).Error
	return instances, err
}

// GetInstanceTasks 获取流程实例的所有任务
func (e *FlowEngine) GetInstanceTasks(instanceID uint) ([]models.FlowTask, error) {
	var tasks []models.FlowTask
	err := e.db.Where("instance_id = ?", instanceID).
		Order("created_at ASC").
		Find(&tasks).Error
	return tasks, err
}

// GetFlowDefinitions 获取流程定义列表
func (e *FlowEngine) GetFlowDefinitions(tenantID string) ([]models.FlowDefinition, error) {
	var defs []models.FlowDefinition
	err := e.db.Where("tenant_id = ? AND is_active = ?", tenantID, true).
		Order("created_at DESC").
		Find(&defs).Error
	return defs, err
}

// CreateFlowDefinition 创建流程定义
func (e *FlowEngine) CreateFlowDefinition(name, description string, nodes []models.FlowNode, tenantID string) (*models.FlowDefinition, error) {
	// 验证节点合法性：必须有且仅有一个start和一个end
	var hasStart, hasEnd bool
	for _, n := range nodes {
		if n.NodeType == models.FlowNodeStart {
			hasStart = true
		}
		if n.NodeType == models.FlowNodeEnd {
			hasEnd = true
		}
	}
	if !hasStart || !hasEnd {
		return nil, errors.New("flow must have both start and end nodes")
	}

	nodesJSON, err := json.Marshal(nodes)
	if err != nil {
		return nil, err
	}

	def := &models.FlowDefinition{
		Name:        name,
		Description: description,
		Version:     1,
		Nodes:       nodesJSON,
		TenantID:    tenantID,
		IsActive:    true,
	}
	if err := e.db.Create(def).Error; err != nil {
		return nil, err
	}
	return def, nil
}

// generateUUID 生成简单的UUID（使用github.com/google/uuid兼容格式）
func generateUUID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
