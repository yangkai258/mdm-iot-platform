package ai

import (
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

// RegionalAINode 区域AI节点
type RegionalAINode struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	NodeID      string    `json:"node_id" gorm:"uniqueIndex;size:64"`
	RegionCode  string    `json:"region_code" gorm:"index;size:16"`
	ModelName   string    `json:"model_name" gorm:"size:64"`
	Endpoint    string    `json:"endpoint" gorm:"size:512"`
	Status      string    `json:"status" gorm:"size:20"` // online/offline
	QPSLimit    int       `json:"qps_limit"`
	CurrentQPS  float64   `json:"current_qps"`
	Priority    int       `json:"priority" gorm:"default:50"`
	Config      string    `json:"config" gorm:"type:text"` // JSON配置
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RegionalAINodeService 区域AI节点服务
type RegionalAINodeService struct {
	db   *gorm.DB
	mu   sync.RWMutex
	nodes map[string]*RegionalAINode
}

// NewRegionalAINodeService 创建服务实例
func NewRegionalAINodeService(db *gorm.DB) *RegionalAINodeService {
	return &RegionalAINodeService{
		db:    db,
		nodes: make(map[string]*RegionalAINode),
	}
}

// RegisterNode 注册AI节点
func (s *RegionalAINodeService) RegisterNode(node *RegionalAINode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查节点是否已存在
	var existing RegionalAINode
	if err := s.db.Where("node_id = ?", node.NodeID).First(&existing).Error; err == nil {
		// 更新现有节点
		node.ID = existing.ID
		return s.db.Save(node).Error
	}

	// 创建新节点
	return s.db.Create(node).Error
}

// GetBestNode 获取最优节点（负载最低且在线）
func (s *RegionalAINodeService) GetBestNode(regionCode string) (*RegionalAINode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var node RegionalAINode
	err := s.db.Where("region_code = ? AND status = ?", regionCode, "online").
		Order("current_qps ASC, priority DESC").
		First(&node).Error

	if err != nil {
		return nil, errors.New("no available AI node in region: " + regionCode)
	}

	// 检查是否超载
	if node.CurrentQPS >= float64(node.QPSLimit) {
		// 尝试找备用节点
		var backup RegionalAINode
		if err := s.db.Where("region_code = ? AND status = ? AND current_qps < qps_limit", regionCode, "online").
			Order("current_qps ASC").
			First(&backup).Error; err == nil {
			return &backup, nil
		}
		return nil, errors.New("all AI nodes in region are overloaded: " + regionCode)
	}

	return &node, nil
}

// ReportLoad 上报负载
func (s *RegionalAINodeService) ReportLoad(nodeID string, qps float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Model(&RegionalAINode{}).Where("node_id = ?", nodeID).Updates(map[string]interface{}{
		"current_qps": qps,
		"updated_at":  time.Now(),
	}).Error
}

// ListNodesByRegion 获取区域内所有节点
func (s *RegionalAINodeService) ListNodesByRegion(regionCode string) ([]RegionalAINode, error) {
	var nodes []RegionalAINode
	err := s.db.Where("region_code = ?", regionCode).Order("priority DESC").Find(&nodes).Error
	return nodes, err
}

// UpdateNodeStatus 更新节点状态
func (s *RegionalAINodeService) UpdateNodeStatus(nodeID string, status string) error {
	return s.db.Model(&RegionalAINode{}).Where("node_id = ?", nodeID).Update("status", status).Error
}

// GetAllNodes 获取所有AI节点
func (s *RegionalAINodeService) GetAllNodes() ([]RegionalAINode, error) {
	var nodes []RegionalAINode
	err := s.db.Order("region_code, priority DESC").Find(&nodes).Error
	return nodes, err
}
