package multi_region

import (
	"errors"
	"sync"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// RegionService 区域服务
type RegionService struct {
	db *gorm.DB
	mu sync.RWMutex
}

// NewRegionService 创建区域服务
func NewRegionService(db *gorm.DB) *RegionService {
	return &RegionService{db: db}
}

// CreateRegion 创建区域
func (s *RegionService) CreateRegion(region *models.Region) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果设为默认，先取消其他默认
	if region.IsDefault {
		s.db.Model(&models.Region{}).Where("is_default = ?", true).Update("is_default", false)
	}

	return s.db.Create(region).Error
}

// GetRegionByID 获取区域
func (s *RegionService) GetRegionByID(id uint) (*models.Region, error) {
	var region models.Region
	err := s.db.First(&region, id).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetRegionByCode 通过代码获取区域
func (s *RegionService) GetRegionByCode(code string) (*models.Region, error) {
	var region models.Region
	err := s.db.Where("region_code = ?", code).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// ListRegions 获取所有区域
func (s *RegionService) ListRegions() ([]models.Region, error) {
	var regions []models.Region
	err := s.db.Order("is_default DESC, created_at DESC").Find(&regions).Error
	return regions, err
}

// UpdateRegion 更新区域
func (s *RegionService) UpdateRegion(id uint, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果设为默认，先取消其他默认
	if isDefault, ok := updates["is_default"].(bool); ok && isDefault {
		s.db.Model(&models.Region{}).Where("is_default = ? AND id != ?", true, id).Update("is_default", false)
	}

	return s.db.Model(&models.Region{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRegion 删除区域
func (s *RegionService) DeleteRegion(id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 删除区域下的节点
	s.db.Where("region_code IN (?)", s.db.Model(&models.Region{}).Select("region_code").Where("id = ?", id)).Delete(&models.RegionalNode{})

	return s.db.Delete(&models.Region{}, id).Error
}

// GetDefaultRegion 获取默认区域
func (s *RegionService) GetDefaultRegion() (*models.Region, error) {
	var region models.Region
	err := s.db.Where("is_default = ?", true).First(&region).Error
	if err != nil {
		return nil, errors.New("no default region configured")
	}
	return &region, nil
}

// HealthCheck 健康检查
func (s *RegionService) HealthCheck(id uint) (map[string]interface{}, error) {
	region, err := s.GetRegionByID(id)
	if err != nil {
		return nil, err
	}

	// 获取区域节点
	var nodes []models.RegionalNode
	s.db.Where("region_code = ?", region.RegionCode).Find(&nodes)

	onlineCount := 0
	totalLoad := 0.0
	nodeStatus := make([]map[string]interface{}, 0)

	for _, node := range nodes {
		status := map[string]interface{}{
			"node_id":     node.NodeID,
			"node_type":   node.NodeType,
			"node_status": node.NodeStatus,
			"load":        node.Load,
		}
		nodeStatus = append(nodeStatus, status)

		if node.NodeStatus == "online" {
			onlineCount++
			totalLoad += node.Load
		}
	}

	avgLoad := 0.0
	if onlineCount > 0 {
		avgLoad = totalLoad / float64(onlineCount)
	}

	return map[string]interface{}{
		"region_id":    region.ID,
		"region_code":  region.RegionCode,
		"status":       region.Status,
		"node_count":   len(nodes),
		"online_count": onlineCount,
		"avg_load":     avgLoad,
		"nodes":        nodeStatus,
	}, nil
}

// Failover 故障切换
func (s *RegionService) Failover(id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	region, err := s.GetRegionByID(id)
	if err != nil {
		return err
	}

	// 查找备用区域
	var backupRegion models.Region
	err = s.db.Where("region_type = ? AND status = ? AND id != ?", "replica", "active", id).
		First(&backupRegion).Error
	if err != nil {
		return errors.New("no backup region available")
	}

	// 更新原区域状态
	s.db.Model(&region).Update("status", "failover")

	// 激活备用区域
	return s.db.Model(&backupRegion).Updates(map[string]interface{}{
		"status": "active",
	}).Error
}

// RegionalNodeService 区域节点服务
type RegionalNodeService struct {
	db *gorm.DB
	mu sync.RWMutex
}

// NewRegionalNodeService 创建节点服务
func NewRegionalNodeService(db *gorm.DB) *RegionalNodeService {
	return &RegionalNodeService{db: db}
}

// CreateNode 创建节点
func (s *RegionalNodeService) CreateNode(node *models.RegionalNode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Create(node).Error
}

// GetNodeByID 获取节点
func (s *RegionalNodeService) GetNodeByID(id uint) (*models.RegionalNode, error) {
	var node models.RegionalNode
	err := s.db.First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// GetNodeByNodeID 通过NodeID获取节点
func (s *RegionalNodeService) GetNodeByNodeID(nodeID string) (*models.RegionalNode, error) {
	var node models.RegionalNode
	err := s.db.Where("node_id = ?", nodeID).First(&node).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// ListNodesByRegion 获取区域节点列表
func (s *RegionalNodeService) ListNodesByRegion(regionCode string) ([]models.RegionalNode, error) {
	var nodes []models.RegionalNode
	err := s.db.Where("region_code = ?", regionCode).Order("node_type, created_at DESC").Find(&nodes).Error
	return nodes, err
}

// UpdateNode 更新节点
func (s *RegionalNodeService) UpdateNode(nodeID string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Model(&models.RegionalNode{}).Where("node_id = ?", nodeID).Updates(updates).Error
}

// DeleteNode 删除节点
func (s *RegionalNodeService) DeleteNode(nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Where("node_id = ?", nodeID).Delete(&models.RegionalNode{}).Error
}

// UpdateNodeStatus 更新节点状态
func (s *RegionalNodeService) UpdateNodeStatus(nodeID string, status string) error {
	return s.db.Model(&models.RegionalNode{}).Where("node_id = ?", nodeID).Update("node_status", status).Error
}

// ReportLoad 上报负载
func (s *RegionalNodeService) ReportLoad(nodeID string, load float64) error {
	return s.db.Model(&models.RegionalNode{}).Where("node_id = ?", nodeID).Updates(map[string]interface{}{
		"load":       load,
		"updated_at": time.Now(),
	}).Error
}

// Heartbeat 节点心跳
func (s *RegionalNodeService) Heartbeat(nodeID string) error {
	return s.db.Model(&models.RegionalNode{}).Where("node_id = ?", nodeID).Updates(map[string]interface{}{
		"node_status": "online",
		"updated_at":  time.Now(),
	}).Error
}
