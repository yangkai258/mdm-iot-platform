package controllers

import (
	"net/http"
	"strings"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes 注册路由 - 合并所有设备相关路由
func RegisterRoutes(r *gin.Engine, db *gorm.DB, redisClient *utils.RedisClient) {
	deviceCtrl := &DeviceController{
		DB:    db,
		Redis: redisClient,
	}

	otaCtrl := &OTAController{DB: db}
	profileCtrl := &PetProfileController{DB: db}
	cmdCtrl := &CommandController{
		DB:    db,
		Redis: redisClient,
		MQTT:  nil,
	}
	orgCtrl := &OrgController{DB: db}
	permCtrl := &PermissionController{DB: db}
	roleCtrl := &RoleController{DB: db}
	memberCtrl := &MemberController{DB: db}

	api := r.Group("/api/v1")
	{
		// ============ 设备管理 ============
		api.POST("/devices/register", deviceCtrl.Register)
		api.GET("/devices", deviceCtrl.List)
		
		// 设备绑定/解绑
		api.POST("/devices/bind/:sn_code", deviceCtrl.Bind)
		api.POST("/devices/unbind/:sn_code", UnbindDevice(db))
		
		// 设备详情及子路由 - 统一用 device_id
		api.GET("/devices/:device_id", deviceCtrl.Get)
		api.PUT("/devices/:device_id", UpdateDevice(db))
		api.DELETE("/devices/:device_id", DeleteDevice(db))
		api.PUT("/devices/:device_id/status", UpdateDeviceStatus(db))
		api.GET("/devices/:device_id/profile", profileCtrl.GetProfile)
		api.PUT("/devices/:device_id/profile", profileCtrl.UpdateProfile)
		api.POST("/devices/:device_id/commands", cmdCtrl.SendCommand)
		api.GET("/devices/:device_id/commands", cmdCtrl.GetCommandHistory)
		
		// OTA 路由
		api.POST("/ota/packages", otaCtrl.CreatePackage)
		api.GET("/ota/packages", otaCtrl.ListPackages)
		api.POST("/ota/deployments", otaCtrl.CreateDeployment)
		api.GET("/ota/devices/:device_id/check", otaCtrl.CheckOTA)

		// ============ 组织管理 ============
		// 公司管理
		api.GET("/org/companies", orgCtrl.CompanyList)
		api.POST("/org/companies", orgCtrl.CompanyCreate)
		api.PUT("/org/companies/:id", orgCtrl.CompanyUpdate)
		api.DELETE("/org/companies/:id", orgCtrl.CompanyDelete)
		
		// 部门管理
		api.GET("/org/departments", orgCtrl.DepartmentList)
		api.GET("/org/departments/tree", orgCtrl.DepartmentTree)
		api.POST("/org/departments", orgCtrl.DepartmentCreate)
		api.PUT("/org/departments/:id", orgCtrl.DepartmentUpdate)
		api.DELETE("/org/departments/:id", orgCtrl.DepartmentDelete)
		
		// 岗位管理
		api.GET("/org/positions", orgCtrl.PositionList)
		api.POST("/org/positions", orgCtrl.PositionCreate)
		api.PUT("/org/positions/:id", orgCtrl.PositionUpdate)
		api.DELETE("/org/positions/:id", orgCtrl.PositionDelete)
		
		// 员工管理
		api.GET("/org/employees", orgCtrl.EmployeeList)
		api.POST("/org/employees", orgCtrl.EmployeeCreate)
		api.PUT("/org/employees/:id", orgCtrl.EmployeeUpdate)
		api.DELETE("/org/employees/:id", orgCtrl.EmployeeDelete)
		
		// 基准岗位管理
		api.GET("/org/standard-positions", orgCtrl.StandardPositionList)
		api.POST("/org/standard-positions", orgCtrl.StandardPositionCreate)
		api.PUT("/org/standard-positions/:id", orgCtrl.StandardPositionUpdate)
		api.DELETE("/org/standard-positions/:id", orgCtrl.StandardPositionDelete)

		// ============ 权限管理 ============
		api.GET("/permissions", permCtrl.List)
		api.POST("/permissions", permCtrl.Create)
		api.PUT("/permissions/:id", permCtrl.Update)
		api.DELETE("/permissions/:id", permCtrl.Delete)

		// ============ 角色管理 ============
		api.GET("/roles", roleCtrl.List)
		api.POST("/roles", roleCtrl.Create)
		api.PUT("/roles/:id", roleCtrl.Update)
		api.DELETE("/roles/:id", roleCtrl.Delete)
		api.GET("/roles/:id/perms", roleCtrl.GetPerms)
		api.PUT("/roles/:id/perms", roleCtrl.SetPerms)

		// ============ 会员管理 ============
		// 会员信息
		api.GET("/members", memberCtrl.MemberList)
		api.POST("/members", memberCtrl.MemberCreate)
		api.PUT("/members/:id", memberCtrl.MemberUpdate)
		api.DELETE("/members/:id", memberCtrl.MemberDelete)
		api.GET("/members/:id", memberCtrl.MemberDetail)
		
		// 会员卡
		api.GET("/member/cards", memberCtrl.CardList)
		api.POST("/member/cards", memberCtrl.CardCreate)
		api.PUT("/member/cards/:id", memberCtrl.CardUpdate)
		api.DELETE("/member/cards/:id", memberCtrl.CardDelete)
		
		// 优惠券
		api.GET("/member/coupons", memberCtrl.CouponList)
		api.POST("/member/coupons", memberCtrl.CouponCreate)
		api.PUT("/member/coupons/:id", memberCtrl.CouponUpdate)
		api.DELETE("/member/coupons/:id", memberCtrl.CouponDelete)
		
		// 店铺管理
		api.GET("/member/stores", memberCtrl.StoreList)
		api.POST("/member/stores", memberCtrl.StoreCreate)
		api.PUT("/member/stores/:id", memberCtrl.StoreUpdate)
		api.DELETE("/member/stores/:id", memberCtrl.StoreDelete)
		
		// 会员标签
		api.GET("/member/tags", memberCtrl.TagList)
		api.POST("/member/tags", memberCtrl.TagCreate)
		api.PUT("/member/tags/:id", memberCtrl.TagUpdate)
		api.DELETE("/member/tags/:id", memberCtrl.TagDelete)
		
		// 促销活动
		api.GET("/member/promotions", memberCtrl.PromotionList)
		api.POST("/member/promotions", memberCtrl.PromotionCreate)
		api.PUT("/member/promotions/:id", memberCtrl.PromotionUpdate)
		api.DELETE("/member/promotions/:id", memberCtrl.PromotionDelete)
		
		// 会员等级
		api.GET("/member/levels", memberCtrl.LevelList)
		
		// 积分规则
		api.GET("/member/points/rules", memberCtrl.PointsRuleList)
		api.POST("/member/points/rules", memberCtrl.PointsRuleCreate)
		api.PUT("/member/points/rules/:id", memberCtrl.PointsRuleUpdate)
		api.DELETE("/member/points/rules/:id", memberCtrl.PointsRuleDelete)
		
		// 积分流水
		api.GET("/member/points/records", memberCtrl.PointsRecordList)
	}
}

// DeviceController 设备管理控制器
type DeviceController struct {
	DB  *gorm.DB
	Redis *utils.RedisClient
}

// RegisterRequest 设备注册请求
type RegisterRequest struct {
	DeviceID        string `json:"device_id" binding:"omitempty"`
	MacAddress      string `json:"mac_address" binding:"required"`
	SnCode         string `json:"sn_code" binding:"required"`
	HardwareModel   string `json:"hardware_model" binding:"required"`
	FirmwareVersion string `json:"firmware_version" binding:"required"`
}

// BindRequest 绑定请求
type BindRequest struct {
	BindUserID string `json:"bind_user_id" binding:"required"`
}

// Register 设备注册/绑定
func (c *DeviceController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败: " + err.Error(),
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	// 校验 MAC 地址格式
	if !isValidMAC(req.MacAddress) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "无效的MAC地址格式",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var device models.Device
	result := c.DB.Where("mac_address = ?", req.MacAddress).First(&device)

	if result.Error == nil {
		// 设备已存在，更新绑定
		device.FirmwareVersion = req.FirmwareVersion
		device.UpdatedAt = time.Now()
		if err := c.DB.Save(&device).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":      5001,
				"message":   "服务器内部错误",
				"error_code": "ERR_INTERNAL",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"device_id":        device.DeviceID,
				"lifecycle_status": device.LifecycleStatus,
				"created_at":       device.CreatedAt,
			},
		})
		return
	}

	// 新建设备
	device = models.Device{
		MacAddress:      req.MacAddress,
		SnCode:         req.SnCode,
		HardwareModel:   req.HardwareModel,
		FirmwareVersion: req.FirmwareVersion,
		LifecycleStatus: 1, // 待激活
	}

	if err := c.DB.Create(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "设备注册失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id":        device.DeviceID,
			"lifecycle_status": device.LifecycleStatus,
			"created_at":       device.CreatedAt,
		},
	})
}

// Bind 扫码绑定设备
func (c *DeviceController) Bind(ctx *gin.Context) {
	snCode := ctx.Param("sn_code")
	var req BindRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4005,
			"message":   "参数校验失败",
			"error_code": "ERR_VALIDATION",
		})
		return
	}

	var device models.Device
	if err := c.DB.Where("sn_code = ?", snCode).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 检查状态
	if device.LifecycleStatus != 1 && device.LifecycleStatus != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":      4003,
			"message":   "非法设备状态，无法绑定",
			"error_code": "ERR_DEVICE_003",
		})
		return
	}

	// 更新绑定
	device.BindUserID = &req.BindUserID
	device.LifecycleStatus = 2 // 服役中
	device.UpdatedAt = time.Now()

	if err := c.DB.Save(&device).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "绑定失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id":        device.DeviceID,
			"bind_user_id":     *device.BindUserID,
			"lifecycle_status": device.LifecycleStatus,
			"message":          "绑定成功",
		},
	})
}

// ListRequest 设备列表请求
type ListRequest struct {
	Page            int    `form:"page"`
	PageSize        int    `form:"page_size"`
	Status          string `form:"status"`
	LifecycleStatus int    `form:"lifecycle_status"`
	HardwareModel   string `form:"hardware_model"`
	Search          string `form:"search"`
}

// List 获取设备列表（带分页和筛选）
func (c *DeviceController) List(ctx *gin.Context) {
	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	query := c.DB.Model(&models.Device{})

	// 筛选条件
	if req.LifecycleStatus > 0 {
		query = query.Where("lifecycle_status = ?", req.LifecycleStatus)
	}
	if req.HardwareModel != "" {
		query = query.Where("hardware_model = ?", req.HardwareModel)
	}
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("device_id LIKE ? OR sn_code LIKE ?", search, search)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 查询列表
	var devices []models.Device
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&devices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":      5001,
			"message":   "查询失败",
			"error_code": "ERR_INTERNAL",
		})
		return
	}

	// 获取在线状态（从 Redis）
	type DeviceWithShadow struct {
		models.Device
		IsOnline      bool `json:"is_online"`
		BatteryLevel  int  `json:"battery_level"`
	}

	result := make([]DeviceWithShadow, len(devices))
	for i, d := range devices {
		result[i].Device = d
		// 从 Redis 获取设备影子
		shadow, err := c.Redis.GetDeviceShadow(d.DeviceID)
		if err == nil && shadow != nil {
			result[i].IsOnline = shadow.IsOnline
			result[i].BatteryLevel = shadow.BatteryLevel
		}
	}

	// 如果有 status 筛选，内存过滤
	if req.Status != "" {
		filtered := make([]DeviceWithShadow, 0)
		for _, d := range result {
			if req.Status == "online" && d.IsOnline {
				filtered = append(filtered, d)
			} else if req.Status == "offline" && !d.IsOnline {
				filtered = append(filtered, d)
			}
		}
		result = filtered
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list": result,
			"pagination": gin.H{
				"page":        req.Page,
				"page_size":   req.PageSize,
				"total":       total,
				"total_pages": (int(total) + req.PageSize - 1) / req.PageSize,
			},
		},
	})
}

// Get 获取单个设备详情
func (c *DeviceController) Get(ctx *gin.Context) {
	deviceID := ctx.Param("device_id")

	var device models.Device
	if err := c.DB.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":      4002,
			"message":   "设备不存在",
			"error_code": "ERR_DEVICE_002",
		})
		return
	}

	// 获取设备影子
	shadow, _ := c.Redis.GetDeviceShadow(deviceID)
	
	// 获取宠物配置
	var petProfile models.PetProfile
	c.DB.Where("device_id = ?", deviceID).First(&petProfile)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"device_id":         device.DeviceID,
			"mac_address":      device.MacAddress,
			"sn_code":          device.SnCode,
			"hardware_model":   device.HardwareModel,
			"firmware_version": device.FirmwareVersion,
			"bind_user_id":     device.BindUserID,
			"lifecycle_status": device.LifecycleStatus,
			"shadow":          shadow,
			"pet_profile":     petProfile,
			"created_at":      device.CreatedAt,
			"updated_at":      device.UpdatedAt,
		},
	})
}

// 校验 MAC 地址格式
func isValidMAC(mac string) bool {
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		return false
	}
	for _, p := range parts {
		if len(p) != 2 {
			return false
		}
	}
	return true
}
