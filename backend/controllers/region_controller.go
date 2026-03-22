package controllers

import (
	"net/http"

	"mdm-backend/multi_region"
	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegionController 区域控制器
type RegionController struct {
	DB       *gorm.DB
	regionSvc *multi_region.RegionService
	nodeSvc   *multi_region.RegionalNodeService
}

// NewRegionController 创建区域控制器
func NewRegionController(db *gorm.DB) *RegionController {
	return &RegionController{
		DB:        db,
		regionSvc: multi_region.NewRegionService(db),
		nodeSvc:   multi_region.NewRegionalNodeService(db),
	}
}

// ListRegions 获取区域列表
// @Summary 获取区域列表
// @Tags regions
// @Produce json
// @Success 200 {array} models.Region
// @Router /api/v1/regions [get]
func (ctrl *RegionController) ListRegions(c *gin.Context) {
	regions, err := ctrl.regionSvc.ListRegions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": regions,
	})
}

// GetRegion 获取区域详情
// @Summary 获取区域详情
// @Tags regions
// @Produce json
// @Param id path int true "区域ID"
// @Success 200 {object} models.Region
// @Router /api/v1/regions/{id} [get]
func (ctrl *RegionController) GetRegion(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	region, err := ctrl.regionSvc.GetRegionByID(regionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": region,
	})
}

// CreateRegion 创建区域
// @Summary 创建区域
// @Tags regions
// @Accept json
// @Produce json
// @Param region body models.Region true "区域信息"
// @Success 201 {object} models.Region
// @Router /api/v1/regions [post]
func (ctrl *RegionController) CreateRegion(c *gin.Context) {
	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成区域代码
	if region.RegionCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "region_code is required"})
		return
	}

	if err := ctrl.regionSvc.CreateRegion(&region); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": region,
	})
}

// UpdateRegion 更新区域
// @Summary 更新区域
// @Tags regions
// @Accept json
// @Produce json
// @Param id path int true "区域ID"
// @Param region body models.Region true "区域信息"
// @Success 200 {object} models.Region
// @Router /api/v1/regions/{id} [put]
func (ctrl *RegionController) UpdateRegion(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.regionSvc.UpdateRegion(regionID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	region, _ := ctrl.regionSvc.GetRegionByID(regionID)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": region,
	})
}

// DeleteRegion 删除区域
// @Summary 删除区域
// @Tags regions
// @Produce json
// @Param id path int true "区域ID"
// @Success 200 {object} gin.H
// @Router /api/v1/regions/{id} [delete]
func (ctrl *RegionController) DeleteRegion(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	if err := ctrl.regionSvc.DeleteRegion(regionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "region deleted",
	})
}

// ListRegionNodes 获取区域节点列表
// @Summary 获取区域节点列表
// @Tags region-nodes
// @Produce json
// @Param id path int true "区域ID"
// @Success 200 {array} models.RegionalNode
// @Router /api/v1/regions/{id}/nodes [get]
func (ctrl *RegionController) ListRegionNodes(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	region, err := ctrl.regionSvc.GetRegionByID(regionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	nodes, err := ctrl.nodeSvc.ListNodesByRegion(region.RegionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": nodes,
	})
}

// CreateRegionNode 创建区域节点
// @Summary 创建区域节点
// @Tags region-nodes
// @Accept json
// @Produce json
// @Param id path int true "区域ID"
// @Param node body models.RegionalNode true "节点信息"
// @Success 201 {object} models.RegionalNode
// @Router /api/v1/regions/{id}/nodes [post]
func (ctrl *RegionController) CreateRegionNode(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	region, err := ctrl.regionSvc.GetRegionByID(regionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	var node models.RegionalNode
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node.RegionCode = region.RegionCode
	if err := ctrl.nodeSvc.CreateNode(&node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": node,
	})
}

// GetRegionNode 获取节点详情
// @Summary 获取节点详情
// @Tags region-nodes
// @Produce json
// @Param id path int true "区域ID"
// @Param node_id path string true "节点ID"
// @Success 200 {object} models.RegionalNode
// @Router /api/v1/regions/{id}/nodes/{node_id} [get]
func (ctrl *RegionController) GetRegionNode(c *gin.Context) {
	nodeID := c.Param("node_id")
	node, err := ctrl.nodeSvc.GetNodeByNodeID(nodeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": node,
	})
}

// UpdateRegionNode 更新节点
// @Summary 更新节点
// @Tags region-nodes
// @Accept json
// @Produce json
// @Param id path int true "区域ID"
// @Param node_id path string true "节点ID"
// @Param node body models.RegionalNode true "节点信息"
// @Success 200 {object} models.RegionalNode
// @Router /api/v1/regions/{id}/nodes/{node_id} [put]
func (ctrl *RegionController) UpdateRegionNode(c *gin.Context) {
	nodeID := c.Param("node_id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.nodeSvc.UpdateNode(nodeID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	node, _ := ctrl.nodeSvc.GetNodeByNodeID(nodeID)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": node,
	})
}

// DeleteRegionNode 删除节点
// @Summary 删除节点
// @Tags region-nodes
// @Produce json
// @Param id path int true "区域ID"
// @Param node_id path string true "节点ID"
// @Success 200 {object} gin.H
// @Router /api/v1/regions/{id}/nodes/{node_id} [delete]
func (ctrl *RegionController) DeleteRegionNode(c *gin.Context) {
	nodeID := c.Param("node_id")

	if err := ctrl.nodeSvc.DeleteNode(nodeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "node deleted",
	})
}

// RegionHealthCheck 区域健康检查
// @Summary 区域健康检查
// @Tags regions
// @Produce json
// @Param id path int true "区域ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/regions/{id}/health [get]
func (ctrl *RegionController) RegionHealthCheck(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	health, err := ctrl.regionSvc.HealthCheck(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": health,
	})
}

// RegionFailover 区域故障切换
// @Summary 区域故障切换
// @Tags regions
// @Produce json
// @Param id path int true "区域ID"
// @Success 200 {object} gin.H
// @Router /api/v1/regions/{id}/failover [post]
func (ctrl *RegionController) RegionFailover(c *gin.Context) {
	id := c.Param("id")
	var regionID uint
	if _, err := parseUint(id, &regionID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region id"})
		return
	}

	if err := ctrl.regionSvc.Failover(regionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "failover initiated",
	})
}
