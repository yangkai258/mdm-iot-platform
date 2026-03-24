package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MeshController BLE Mesh 网络控制器
type MeshController struct {
	DB *gorm.DB
}

// RegisterMeshRoutes 注册 Mesh 相关路由
func (c *MeshController) RegisterMeshRoutes(api *gin.RouterGroup) {
	mesh := api.Group("/mesh")
	{
		// Mesh 网络 CRUD
		mesh.GET("/networks", c.ListMeshNetworks)
		mesh.POST("/networks", c.CreateMeshNetwork)
		mesh.GET("/networks/:id", c.GetMeshNetwork)
		mesh.PUT("/networks/:id", c.UpdateMeshNetwork)
		mesh.DELETE("/networks/:id", c.DeleteMeshNetwork)
		mesh.POST("/networks/:id/activate", c.ActivateMeshNetwork)
		mesh.POST("/networks/:id/deactivate", c.DeactivateMeshNetwork)
		mesh.GET("/networks/:id/topology", c.GetMeshTopology)

		// Mesh 设备 CRUD
		mesh.GET("/devices", c.ListMeshDevices)
		mesh.POST("/devices", c.CreateMeshDevice)
		mesh.GET("/devices/:id", c.GetMeshDevice)
		mesh.PUT("/devices/:id", c.UpdateMeshDevice)
		mesh.DELETE("/devices/:id", c.DeleteMeshDevice)
		mesh.POST("/devices/:id/connect", c.ConnectMeshDevice)
		mesh.POST("/devices/:id/disconnect", c.DisconnectMeshDevice)
		mesh.GET("/devices/:id/neighbors", c.GetMeshDeviceNeighbors)
	}
}

// CreateMeshNetworkRequest 创建 Mesh 网络请求
type CreateMeshNetworkRequest struct {
	Name            string `json:"name" binding:"required"`
	NetworkID       string `json:"network_id" binding:"required"`
	SSID            string `json:"ssid"`
	SecurityKey     string `json:"security_key"`
	EncryptionType  string `json:"encryption_type"`
	Channel         int    `json:"channel"`
	GatewayDeviceID string `json:"gateway_device_id"`
}

// CreateMeshNetwork 创建 Mesh 网络
// POST /api/v1/mesh/networks
func (c *MeshController) CreateMeshNetwork(ctx *gin.Context) {
	var req CreateMeshNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	orgID, _ := ctx.Get("org_id")
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 检查 network_id 是否已存在
	var existing models.MeshNetwork
	if err := c.DB.Where("network_id = ?", req.NetworkID).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":    4090,
			"message": "Mesh网络ID已存在",
		})
		return
	}

	network := models.MeshNetwork{
		Name:            req.Name,
		NetworkID:       req.NetworkID,
		SSID:            req.SSID,
		SecurityKey:     req.SecurityKey,
		EncryptionType:  req.EncryptionType,
		Channel:         req.Channel,
		GatewayDeviceID: req.GatewayDeviceID,
		Status:          "inactive",
		OrgID:           orgID.(uint),
		CreateUserID:    uid,
	}
	if network.EncryptionType == "" {
		network.EncryptionType = "aes256"
	}
	if network.Channel == 0 {
		network.Channel = 6
	}

	if err := c.DB.Create(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "创建Mesh网络失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": network,
	})
}

// UpdateMeshNetworkRequest 更新 Mesh 网络请求
type UpdateMeshNetworkRequest struct {
	Name            string `json:"name"`
	SSID            string `json:"ssid"`
	SecurityKey     string `json:"security_key"`
	EncryptionType  string `json:"encryption_type"`
	Channel         int    `json:"channel"`
	GatewayDeviceID string `json:"gateway_device_id"`
}

// UpdateMeshNetwork 更新 Mesh 网络配置
// PUT /api/v1/mesh/networks/:id
func (c *MeshController) UpdateMeshNetwork(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的网络ID",
		})
		return
	}

	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh网络不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络失败: " + err.Error(),
		})
		return
	}

	var req UpdateMeshNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != "" {
		network.Name = req.Name
	}
	if req.SSID != "" {
		network.SSID = req.SSID
	}
	if req.SecurityKey != "" {
		network.SecurityKey = req.SecurityKey
	}
	if req.EncryptionType != "" {
		network.EncryptionType = req.EncryptionType
	}
	if req.Channel != 0 {
		network.Channel = req.Channel
	}
	if req.GatewayDeviceID != "" {
		network.GatewayDeviceID = req.GatewayDeviceID
	}

	if err := c.DB.Save(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "更新Mesh网络失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": network,
	})
}

// DeleteMeshNetwork 删除 Mesh 网络
// DELETE /api/v1/mesh/networks/:id
func (c *MeshController) DeleteMeshNetwork(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的网络ID",
		})
		return
	}

	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh网络不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络失败: " + err.Error(),
		})
		return
	}

	// 删除网络成员
	c.DB.Where("network_id = ?", id).Delete(&models.MeshNetworkMember{})

	// 软删除网络
	if err := c.DB.Delete(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "删除Mesh网络失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "Mesh网络已删除",
	})
}

// ListMeshNetworks 获取 Mesh 网络列表
// GET /api/v1/mesh/networks
func (c *MeshController) ListMeshNetworks(ctx *gin.Context) {
	var networks []models.MeshNetwork
	query := c.DB.Model(&models.MeshNetwork{})

	name := ctx.Query("name")
	status := ctx.Query("status")
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize)

	if err := query.Find(&networks).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      networks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMeshNetwork 获取 Mesh 网络详情
// GET /api/v1/mesh/networks/:id
func (c *MeshController) GetMeshNetwork(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的网络ID",
		})
		return
	}

	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh网络不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络失败: " + err.Error(),
		})
		return
	}

	// 获取成员列表
	var members []models.MeshNetworkMember
	c.DB.Where("network_id = ?", id).Find(&members)

	// 获取设备列表
	var devices []models.MeshDevice
	deviceIDs := make([]string, len(members))
	for i, m := range members {
		deviceIDs[i] = m.DeviceID
	}
	if len(deviceIDs) > 0 {
		c.DB.Where("device_id IN ?", deviceIDs).Find(&devices)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"network":  network,
			"members":  members,
			"devices":  devices,
		},
	})
}

// ActivateMeshNetwork 激活 Mesh 网络
// POST /api/v1/mesh/networks/:id/activate
func (c *MeshController) ActivateMeshNetwork(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的网络ID",
		})
		return
	}

	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh网络不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络失败: " + err.Error(),
		})
		return
	}

	network.Status = "active"
	if err := c.DB.Save(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "激活Mesh网络失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": network,
	})
}

// DeactivateMeshNetwork 停用 Mesh 网络
// POST /api/v1/mesh/networks/:id/deactivate
func (c *MeshController) DeactivateMeshNetwork(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的网络ID",
		})
		return
	}

	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh网络不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络失败: " + err.Error(),
		})
		return
	}

	network.Status = "inactive"
	if err := c.DB.Save(&network).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "停用Mesh网络失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": network,
	})
}

// GetMeshTopology 获取网络拓扑
// GET /api/v1/mesh/networks/:id/topology
func (c *MeshController) GetMeshTopology(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "无效的网络ID",
		})
		return
	}

	var network models.MeshNetwork
	if err := c.DB.First(&network, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh网络不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh网络失败: " + err.Error(),
		})
		return
	}

	// 获取所有成员
	var members []models.MeshNetworkMember
	c.DB.Where("network_id = ?", id).Find(&members)

	// 获取所有设备详情
	deviceIDs := make([]string, len(members))
	for i, m := range members {
		deviceIDs[i] = m.DeviceID
	}

	var devices []models.MeshDevice
	if len(deviceIDs) > 0 {
		c.DB.Where("device_id IN ?", deviceIDs).Find(&devices)
	}

	// 构建拓扑数据
	type TopologyNode struct {
		DeviceID        string  `json:"device_id"`
		MeshAddress     string  `json:"mesh_address"`
		ParentDeviceID  string  `json:"parent_device_id"`
		Role            string  `json:"role"`
		HopCount        int     `json:"hop_count"`
		SignalStrength  int     `json:"signal_strength"`
		ConnectionStatus string `json:"connection_status"`
		LatencyMs       float64 `json:"latency_ms"`
		IsActive        bool    `json:"is_active"`
	}

	type TopologyLink struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}

	nodes := make([]TopologyNode, 0, len(devices))
	links := make([]TopologyLink, 0)
	deviceMap := make(map[string]models.MeshDevice)

	for _, d := range devices {
		deviceMap[d.DeviceID] = d
	}

	for _, m := range members {
		var latency float64
		for _, d := range devices {
			if d.DeviceID == m.DeviceID {
				latency = float64(d.SignalStrength) // 暂用信号强度代替延迟
				break
			}
		}
		node := TopologyNode{
			DeviceID:        m.DeviceID,
			MeshAddress:     m.MeshAddress,
			Role:            m.Role,
			LatencyMs:       latency,
			IsActive:        m.IsActive,
		}
		for _, d := range devices {
			if d.DeviceID == m.DeviceID {
				node.ParentDeviceID = d.ParentDeviceID
				node.HopCount = d.HopCount
				node.SignalStrength = d.SignalStrength
				node.ConnectionStatus = d.ConnectionStatus
				break
			}
		}
		nodes = append(nodes, node)

		// 构建链路关系
		for _, d := range devices {
			if d.DeviceID == m.DeviceID && d.ParentDeviceID != "" {
				links = append(links, TopologyLink{
					Source: d.ParentDeviceID,
					Target: d.DeviceID,
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"network_id": network.NetworkID,
			"network_name": network.Name,
			"status":     network.Status,
			"nodes":      nodes,
			"links":      links,
		},
	})
}

// ============ Mesh 设备 API ============

// CreateMeshDeviceRequest 添加设备到 Mesh 请求
type CreateMeshDeviceRequest struct {
	DeviceID       string `json:"device_id" binding:"required"`
	MeshUUID       string `json:"mesh_uuid"`
	NetworkID      uint   `json:"network_id" binding:"required"`
	ParentDeviceID string `json:"parent_device_id"`
	MeshAddress    string `json:"mesh_address"`
	Role           string `json:"role"`
	HardwareModel  string `json:"hardware_model"`
	FirmwareVersion string `json:"firmware_version"`
}

// CreateMeshDevice 添加设备到 Mesh
// POST /api/v1/mesh/devices
func (c *MeshController) CreateMeshDevice(ctx *gin.Context) {
	var req CreateMeshDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	orgID, _ := ctx.Get("org_id")
	userID, _ := ctx.Get("user_id")
	uid, _ := userID.(uint)

	// 检查设备是否已存在
	var existing models.MeshDevice
	if err := c.DB.Where("device_id = ?", req.DeviceID).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":    4090,
			"message": "Mesh设备已存在",
		})
		return
	}

	// 检查网络是否存在
	var network models.MeshNetwork
	if err := c.DB.First(&network, req.NetworkID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    4000,
				"message": "Mesh网络不存在",
			})
			return
		}
	}

	device := models.MeshDevice{
		DeviceID:        req.DeviceID,
		MeshUUID:        req.MeshUUID,
		ParentDeviceID:  req.ParentDeviceID,
		MeshAddress:     req.MeshAddress,
		Role:            req.Role,
		HardwareModel:   req.HardwareModel,
		FirmwareVersion: req.FirmwareVersion,
		ConnectionStatus: "disconnected",
		OrgID:           orgID.(uint),
		CreateUserID:    uid,
	}
	if device.Role == "" {
		device.Role = "node"
	}
	if device.MeshAddress == "" {
		device.MeshAddress = fmt.Sprintf("%04x", time.Now().UnixNano()&0xFFFF)
	}

	if err := c.DB.Create(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "创建Mesh设备失败: " + err.Error(),
		})
		return
	}

	// 添加到网络成员
	now := time.Now()
	member := models.MeshNetworkMember{
		NetworkID:   req.NetworkID,
		DeviceID:    device.DeviceID,
		MeshAddress: device.MeshAddress,
		Role:        device.Role,
		JoinedAt:    &now,
		IsActive:    true,
	}
	c.DB.Create(&member)

	// 更新网络设备数量
	c.DB.Model(&network).Update("device_count", gorm.Expr("device_count + 1"))

	ctx.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": gin.H{
			"device": device,
			"member": member,
		},
	})
}

// ListMeshDevices 获取 Mesh 设备列表
// GET /api/v1/mesh/devices
func (c *MeshController) ListMeshDevices(ctx *gin.Context) {
	var devices []models.MeshDevice
	query := c.DB.Model(&models.MeshDevice{})

	// 过滤条件
	networkID := ctx.Query("network_id")
	status := ctx.Query("connection_status")
	role := ctx.Query("role")

	if networkID != "" {
		// 通过 MeshNetworkMember 关联过滤
		query = query.Joins("JOIN mesh_network_members mnm ON mnm.device_id = mesh_devices.device_id").
			Where("mnm.network_id = ?", networkID)
	}
	if status != "" {
		query = query.Where("connection_status = ?", status)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	page := defaultPage(ctx)
	pageSize := defaultPageSize(ctx)
	var total int64

	query.Count(&total)
	query.Order("mesh_devices.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize)

	if err := query.Find(&devices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询Mesh设备列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      devices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMeshDevice 获取 Mesh 设备详情
// GET /api/v1/mesh/devices/:id
func (c *MeshController) GetMeshDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	var device models.MeshDevice
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	// 获取所属网络
	var member models.MeshNetworkMember
	c.DB.Where("device_id = ?", deviceID).First(&member)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device":  device,
			"member":  member,
		},
	})
}

// UpdateMeshDeviceRequest 更新设备请求
type UpdateMeshDeviceRequest struct {
	MeshUUID        string `json:"mesh_uuid"`
	ParentDeviceID  string `json:"parent_device_id"`
	MeshAddress     string `json:"mesh_address"`
	Role            string `json:"role"`
	HopCount        int    `json:"hop_count"`
	SignalStrength  int    `json:"signal_strength"`
	HardwareModel   string `json:"hardware_model"`
	FirmwareVersion string `json:"firmware_version"`
}

// UpdateMeshDevice 更新设备
// PUT /api/v1/mesh/devices/:id
func (c *MeshController) UpdateMeshDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	var device models.MeshDevice
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	var req UpdateMeshDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "参数校验失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.MeshUUID != "" {
		device.MeshUUID = req.MeshUUID
	}
	if req.ParentDeviceID != "" {
		device.ParentDeviceID = req.ParentDeviceID
	}
	if req.MeshAddress != "" {
		device.MeshAddress = req.MeshAddress
	}
	if req.Role != "" {
		device.Role = req.Role
	}
	if req.HopCount >= 0 {
		device.HopCount = req.HopCount
	}
	if req.SignalStrength != 0 {
		device.SignalStrength = req.SignalStrength
	}
	if req.HardwareModel != "" {
		device.HardwareModel = req.HardwareModel
	}
	if req.FirmwareVersion != "" {
		device.FirmwareVersion = req.FirmwareVersion
	}

	if err := c.DB.Save(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "更新设备失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": device,
	})
}

// DeleteMeshDevice 从 Mesh 移除设备
// DELETE /api/v1/mesh/devices/:id
func (c *MeshController) DeleteMeshDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	var device models.MeshDevice
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	// 获取设备所在网络，更新成员计数
	var member models.MeshNetworkMember
	if err := c.DB.Where("device_id = ?", deviceID).First(&member).Error; err == nil {
		c.DB.Model(&models.MeshNetwork{}).Where("id = ?", member.NetworkID).Update("device_count", gorm.Expr("device_count - 1"))
		// 删除网络成员
		c.DB.Where("device_id = ?", deviceID).Delete(&models.MeshNetworkMember{})
	}

	// 软删除设备
	if err := c.DB.Delete(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "删除设备失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "设备已从Mesh网络移除",
	})
}

// ConnectMeshDevice 连接设备
// POST /api/v1/mesh/devices/:id/connect
func (c *MeshController) ConnectMeshDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	var device models.MeshDevice
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	now := time.Now()
	device.ConnectionStatus = "connecting"
	device.ConnectedAt = &now
	device.LastSeenAt = &now

	if err := c.DB.Save(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "更新设备状态失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":         device.DeviceID,
			"connection_status": "connecting",
			"connected_at":      now,
		},
	})
}

// DisconnectMeshDevice 断开设备
// POST /api/v1/mesh/devices/:id/disconnect
func (c *MeshController) DisconnectMeshDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	var device models.MeshDevice
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	now := time.Now()
	device.ConnectionStatus = "disconnected"
	device.LastSeenAt = &now

	if err := c.DB.Save(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "更新设备状态失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":         device.DeviceID,
			"connection_status": "disconnected",
			"last_seen_at":      now,
		},
	})
}

// GetMeshDeviceNeighbors 获取邻居设备
// GET /api/v1/mesh/devices/:id/neighbors
func (c *MeshController) GetMeshDeviceNeighbors(ctx *gin.Context) {
	deviceID := ctx.Param("id")

	var device models.MeshDevice
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    4040,
				"message": "Mesh设备不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": "查询设备失败: " + err.Error(),
		})
		return
	}

	// 查找同一网络中的其他设备作为邻居
	var members []models.MeshNetworkMember
	c.DB.Where("device_id = ?", deviceID).Find(&members)

	var neighbors []models.MeshDevice
	if len(members) > 0 {
		networkID := members[0].NetworkID
		var neighborMembers []models.MeshNetworkMember
		c.DB.Where("network_id = ? AND device_id != ?", networkID, deviceID).Find(&neighborMembers)

		neighborIDs := make([]string, len(neighborMembers))
		for i, nm := range neighborMembers {
			neighborIDs[i] = nm.DeviceID
		}
		if len(neighborIDs) > 0 {
			c.DB.Where("device_id IN ?", neighborIDs).Find(&neighbors)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id": deviceID,
			"neighbors": neighbors,
			"count":     len(neighbors),
		},
	})
}
