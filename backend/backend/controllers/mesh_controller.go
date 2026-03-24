package controllers

import (
	"net/http"
	"strconv"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MeshController BLE Mesh 网络控制器
type MeshController struct {
	DB *gorm.DB
}

// ============ MeshNetwork ============

// MeshNetworkList 网络列表
func (c *MeshController) MeshNetworkList(ctx *gin.Context) {
	var networks []models.MeshNetwork
	var total int64

	query := c.DB.Model(&models.MeshNetwork{})

	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&networks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": networks, "total": total, "page": page, "page_size": pageSize,
	}})
}

// MeshNetworkGet 获取网络
func (c *MeshController) MeshNetworkGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "网络不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": network})
}

// MeshNetworkCreate 创建网络
func (c *MeshController) MeshNetworkCreate(ctx *gin.Context) {
	var network models.MeshNetwork
	if err := ctx.ShouldBindJSON(&network); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": network})
}

// MeshNetworkUpdate 更新网络
func (c *MeshController) MeshNetworkUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "网络不存在"})
		return
	}
	var updateData models.MeshNetwork
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updateData})
}

// MeshNetworkDelete 删除网络
func (c *MeshController) MeshNetworkDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.MeshNetwork{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ MeshNode ============

// MeshNodeList 节点列表
func (c *MeshController) MeshNodeList(ctx *gin.Context) {
	var nodes []models.MeshNode
	var total int64

	query := c.DB.Model(&models.MeshNode{})

	if networkUUID := ctx.Query("network_uuid"); networkUUID != "" {
		query = query.Where("network_uuid = ?", networkUUID)
	}
	if isOnline := ctx.Query("is_online"); isOnline != "" {
		query = query.Where("is_online = ?", isOnline == "true")
	}
	if isProvisioned := ctx.Query("is_provisioned"); isProvisioned != "" {
		query = query.Where("is_provisioned = ?", isProvisioned == "true")
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("is_online DESC, last_seen_at DESC").Find(&nodes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": nodes, "total": total, "page": page, "page_size": pageSize,
	}})
}

// MeshNodeGet 获取节点
func (c *MeshController) MeshNodeGet(ctx *gin.Context) {
	id := ctx.Param("id")
	var node models.MeshNode
	if err := c.DB.First(&node, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "节点不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": node})
}

// MeshNodeCreate 创建节点
func (c *MeshController) MeshNodeCreate(ctx *gin.Context) {
	var node models.MeshNode
	if err := ctx.ShouldBindJSON(&node); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&node).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": node})
}

// MeshNodeUpdate 更新节点
func (c *MeshController) MeshNodeUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var node models.MeshNode
	if err := c.DB.First(&node, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "节点不存在"})
		return
	}
	var updateData models.MeshNode
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Save(&updateData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": updateData})
}

// MeshNodeDelete 删除节点
func (c *MeshController) MeshNodeDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.MeshNode{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ MeshGroup ============

// MeshGroupList 组列表
func (c *MeshController) MeshGroupList(ctx *gin.Context) {
	var groups []models.MeshGroup
	var total int64

	query := c.DB.Model(&models.MeshGroup{})

	if networkUUID := ctx.Query("network_uuid"); networkUUID != "" {
		query = query.Where("network_uuid = ?", networkUUID)
	}

	query.Count(&total)

	if err := query.Order("id DESC").Find(&groups).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": groups, "total": total,
	}})
}

// MeshGroupCreate 创建组
func (c *MeshController) MeshGroupCreate(ctx *gin.Context) {
	var group models.MeshGroup
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&group).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": group})
}

// MeshGroupDelete 删除组
func (c *MeshController) MeshGroupDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.DB.Delete(&models.MeshGroup{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ MeshNodeGroup ============

// MeshNodeGroupList 节点组成员
func (c *MeshController) MeshNodeGroupList(ctx *gin.Context) {
	groupUUID := ctx.Query("group_uuid")
	var nodeGroups []models.MeshNodeGroup
	var total int64

	query := c.DB.Model(&models.MeshNodeGroup{})
	if groupUUID != "" {
		query = query.Where("group_uuid = ?", groupUUID)
	}

	query.Count(&total)

	if err := query.Find(&nodeGroups).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": nodeGroups, "total": total,
	}})
}

// MeshNodeGroupAdd 添加节点到组
func (c *MeshController) MeshNodeGroupAdd(ctx *gin.Context) {
	var nodeGroup models.MeshNodeGroup
	if err := ctx.ShouldBindJSON(&nodeGroup); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	// 检查是否已存在
	var existing models.MeshNodeGroup
	err := c.DB.Where("node_uuid = ? AND group_uuid = ?", nodeGroup.NodeUUID, nodeGroup.GroupUUID).First(&existing).Error
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": existing})
		return
	}
	if err := c.DB.Create(&nodeGroup).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": nodeGroup})
}

// MeshNodeGroupRemove 从组移除节点
func (c *MeshController) MeshNodeGroupRemove(ctx *gin.Context) {
	nodeUUID := ctx.Query("node_uuid")
	groupUUID := ctx.Query("group_uuid")
	if err := c.DB.Where("node_uuid = ? AND group_uuid = ?", nodeUUID, groupUUID).Delete(&models.MeshNodeGroup{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ============ MeshTelemetry ============

// MeshTelemetryList 遥测数据列表
func (c *MeshController) MeshTelemetryList(ctx *gin.Context) {
	var telemetry []models.MeshTelemetry
	var total int64

	query := c.DB.Model(&models.MeshTelemetry{})

	if nodeUUID := ctx.Query("node_uuid"); nodeUUID != "" {
		query = query.Where("node_uuid = ?", nodeUUID)
	}
	if telemetryType := ctx.Query("telemetry_type"); telemetryType != "" {
		query = query.Where("telemetry_type = ?", telemetryType)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("recorded_at DESC").Find(&telemetry).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{
		"list": telemetry, "total": total, "page": page, "page_size": pageSize,
	}})
}

// MeshTelemetryCreate 创建遥测数据
func (c *MeshController) MeshTelemetryCreate(ctx *gin.Context) {
	var telemetry models.MeshTelemetry
	if err := ctx.ShouldBindJSON(&telemetry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := c.DB.Create(&telemetry).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": telemetry})
}
