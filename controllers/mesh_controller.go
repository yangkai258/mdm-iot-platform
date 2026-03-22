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
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize)

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

// ConnectMeshDevice 连接 Mesh 设备
// POST /api/v1/mesh/devices/:id/connect
func (c *MeshController) ConnectMeshDevice(ctx *gin.Context) {
	deviceID := ctx.Param("id")
	if deviceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "设备ID不能为空",
		})
		return
	}

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

	// 更新连接状态为 connecting
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

	// TODO: 通过 MQTT/BLE 协议向设备下发连接指令
	// mqttHandler.PublishMeshConnect(device)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device_id":        device.DeviceID,
			"connection_status": "connecting",
			"connected_at":    now,
		},
	})
}

// SetupMeshNetworkRequest 设置 Mesh 网络请求
type SetupMeshNetworkRequest struct {
	Name           string   `json:"name" binding:"required"`
	NetworkID      string   `json:"network_id" binding:"required"`
	SSID           string   `json:"ssid"`
	EncryptionType string   `json:"encryption_type"`
	Channel        int      `json:"channel"`
	GatewayDeviceID string  `json:"gateway_device_id"`
	DeviceIDs      []string `json:"device_ids"`
}

// SetupMeshNetwork 设置 Mesh 网络
// POST /api/v1/mesh/network/setup
func (c *MeshController) SetupMeshNetwork(ctx *gin.Context) {
	var req SetupMeshNetworkRequest
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

	// 检查网络是否已存在
	var existing models.MeshNetwork
	if err := c.DB.Where("network_id = ?", req.NetworkID).First(&existing).Error; err == nil {
		// 网络已存在，更新配置
		existing.Name = req.Name
		existing.SSID = req.SSID
		existing.EncryptionType = req.EncryptionType
		existing.Channel = req.Channel
		existing.GatewayDeviceID = req.GatewayDeviceID
		existing.Status = "configuring"
		if err := c.DB.Save(&existing).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    5000,
				"message": "更新Mesh网络失败: " + err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": existing,
			"message": "Mesh网络已更新",
		})
		return
	}

	// 创建新网络
	network := models.MeshNetwork{
		Name:            req.Name,
		NetworkID:       req.NetworkID,
		SSID:            req.SSID,
		EncryptionType:  req.EncryptionType,
		Channel:         req.Channel,
		GatewayDeviceID: req.GatewayDeviceID,
		Status:          "configuring",
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

	// 批量创建设备成员
	if len(req.DeviceIDs) > 0 {
		now := time.Now()
		members := make([]models.MeshNetworkMember, 0, len(req.DeviceIDs))
		for i, devID := range req.DeviceIDs {
			members = append(members, models.MeshNetworkMember{
				NetworkID:   network.ID,
				DeviceID:    devID,
				MeshAddress: fmt.Sprintf("%04x", i+1), // 自动分配 Mesh 地址
				Role:        "node",
				JoinedAt:    &now,
				IsActive:    true,
			})
		}
		c.DB.Create(&members)
		network.DeviceCount = len(req.DeviceIDs)
		c.DB.Save(&network)
	}

	// TODO: 通过 MQTT 下发 Mesh 网络配置到网关设备

	ctx.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"data": gin.H{
			"network":     network,
			"members":     len(req.DeviceIDs),
		},
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
	c.DB.Where("device_id = ?", deviceID).Preload("MeshNetwork").First(&member)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"device":  device,
			"network": member.MeshNetwork,
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
