package controllers

import (
	"net/http"
	"strings"
	"time"

	"mdm-backend/models"
	"mdm-backend/mqtt"
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
	}
	orgCtrl := &OrgController{DB: db}
	companyCtrl := &CompanyController{DB: db}
	deptCtrl := &DepartmentController{DB: db}
	postCtrl := &PostController{DB: db}
	empCtrl := &EmployeeController{DB: db}
	// roleCtrl 已迁移到 NewRoleController (main.go)
	memberCtrl := &MemberController{DB: db}
	memberEnhancedCtrl := NewMemberEnhancedController(db)
	positionTemplateCtrl := &PositionTemplateController{DB: db}

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
		api.GET("/devices/:device_id/desired-state", GetDesiredState(db))
		api.PUT("/devices/:device_id/desired-state", SetDesiredState(db, mqtt.GlobalMQTTClient))
		api.GET("/devices/:device_id/commands", cmdCtrl.GetCommandHistory)
		
		// ============ OTA 路由 ============
		// 固件包管理
		api.POST("/ota/packages", otaCtrl.CreatePackage)
		api.GET("/ota/packages", otaCtrl.ListPackages)
		// 部署任务管理
		api.POST("/ota/deployments", otaCtrl.CreateDeployment)
		api.GET("/ota/deployments", otaCtrl.ListDeployments)
		api.GET("/ota/deployments/:id", otaCtrl.GetDeployment)
		api.POST("/ota/deployments/:id/pause", otaCtrl.PauseDeployment)
		api.POST("/ota/deployments/:id/resume", otaCtrl.ResumeDeployment)
		api.POST("/ota/deployments/:id/cancel", otaCtrl.CancelDeployment)
		api.GET("/ota/deployments/:id/progress", otaCtrl.GetDeploymentProgress)
		// 设备 OTA 回调（设备主动上报进度）
		api.POST("/ota/devices/:device_id/report", otaCtrl.DeviceOTAReport)
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
		// posts 别名（与 /org/positions 等效）
		api.GET("/org/posts", orgCtrl.PositionList)
		api.POST("/org/posts", orgCtrl.PositionCreate)
		api.PUT("/org/posts/:id", orgCtrl.PositionUpdate)
		api.DELETE("/org/posts/:id", orgCtrl.PositionDelete)

		// 员工管理
		api.GET("/org/employees", orgCtrl.EmployeeList)
		api.POST("/org/employees", orgCtrl.EmployeeCreate)
		api.PUT("/org/employees/:id", orgCtrl.EmployeeUpdate)
		api.DELETE("/org/employees/:id", orgCtrl.EmployeeDelete)

		// ============ Sprint 5: 租户隔离的组织管理 API ============
		// 公司管理（租户隔离）
		tenantGroup := api.Group("/tenants/:tenant_id")
		{
			tenantGroup.GET("/companies", companyCtrl.CompanyList)
			tenantGroup.POST("/companies", companyCtrl.CompanyCreate)
			tenantGroup.GET("/companies/:id", companyCtrl.CompanyGet)
			tenantGroup.PUT("/companies/:id", companyCtrl.CompanyUpdate)
			tenantGroup.DELETE("/companies/:id", companyCtrl.CompanyDelete)

			// 部门管理（租户隔离）
			tenantGroup.GET("/departments", deptCtrl.DepartmentList)
			tenantGroup.GET("/departments/tree", deptCtrl.DepartmentTree)
			tenantGroup.POST("/departments", deptCtrl.DepartmentCreate)
			tenantGroup.GET("/departments/:id", deptCtrl.DepartmentGet)
			tenantGroup.PUT("/departments/:id", deptCtrl.DepartmentUpdate)
			tenantGroup.DELETE("/departments/:id", deptCtrl.DepartmentDelete)

			// 岗位管理（租户隔离）
			tenantGroup.GET("/posts", postCtrl.PostList)
			tenantGroup.POST("/posts", postCtrl.PostCreate)
			tenantGroup.GET("/posts/:id", postCtrl.PostGet)
			tenantGroup.PUT("/posts/:id", postCtrl.PostUpdate)
			tenantGroup.DELETE("/posts/:id", postCtrl.PostDelete)

			// 员工管理（租户隔离）
			tenantGroup.GET("/employees", empCtrl.EmployeeList)
			tenantGroup.POST("/employees/onboard", empCtrl.EmployeeOnboard)
			tenantGroup.GET("/employees/:id", empCtrl.EmployeeGet)
			tenantGroup.PUT("/employees/:id/leave", empCtrl.EmployeeLeave)
		}
		
		// 基准岗位管理
		api.GET("/org/standard-positions", orgCtrl.StandardPositionList)
		api.POST("/org/standard-positions", orgCtrl.StandardPositionCreate)
		api.PUT("/org/standard-positions/:id", orgCtrl.StandardPositionUpdate)
		api.DELETE("/org/standard-positions/:id", orgCtrl.StandardPositionDelete)
		api.POST("/org/standard-positions/:id/enable", orgCtrl.StandardPositionEnable)
		api.POST("/org/standard-positions/:id/disable", orgCtrl.StandardPositionDisable)

		// ============ 基准岗位模板 ============
		api.GET("/position-templates", positionTemplateCtrl.PositionTemplateList)
		api.POST("/position-templates", positionTemplateCtrl.PositionTemplateCreate)
		api.GET("/position-templates/:id", positionTemplateCtrl.PositionTemplateGet)
		api.PUT("/position-templates/:id", positionTemplateCtrl.PositionTemplateUpdate)
		api.DELETE("/position-templates/:id", positionTemplateCtrl.PositionTemplateDelete)
		api.POST("/position-templates/:id/enable", positionTemplateCtrl.PositionTemplateEnable)
		api.POST("/position-templates/:id/disable", positionTemplateCtrl.PositionTemplateDisable)
		api.POST("/position-templates/:id/clone", positionTemplateCtrl.PositionTemplateClone)
		api.POST("/position-templates/:id/copy", positionTemplateCtrl.PositionTemplateCopy)

		// ============ 角色管理 ============
		// (已迁移到 main.go NewRoleController)
		// api.GET("/roles", roleCtrl.List)
		// api.POST("/roles", roleCtrl.Create)
		// api.GET("/roles/:id", roleCtrl.Get)
		// api.PUT("/roles/:id", roleCtrl.Update)
		// api.DELETE("/roles/:id", roleCtrl.Delete)
		// api.GET("/roles/:id/permissions", roleCtrl.GetPermissions)
		// api.POST("/roles/:id/permissions", roleCtrl.AssignPermissions)
		// api.GET("/permissions", roleCtrl.ListPermissions)

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
		api.POST("/member/levels", memberCtrl.LevelCreate)
		api.PUT("/member/levels/:id", memberCtrl.LevelUpdate)
		api.DELETE("/member/levels/:id", memberCtrl.LevelDelete)

		// 会员订单
		api.GET("/member/orders", memberCtrl.OrderList)
		api.POST("/member/orders", memberCtrl.OrderCreate)
		api.GET("/member/orders/:id", memberCtrl.OrderDetail)

		// 等级调整流水
		api.GET("/member/upgrade-records", memberCtrl.UpgradeRecordList)
		
		// 积分规则
		api.GET("/member/points/rules", memberCtrl.PointsRuleList)
		api.POST("/member/points/rules", memberCtrl.PointsRuleCreate)
		api.PUT("/member/points/rules/:id", memberCtrl.PointsRuleUpdate)
		api.DELETE("/member/points/rules/:id", memberCtrl.PointsRuleDelete)
		
		// 积分流水
		api.GET("/member/points/records", memberCtrl.PointsRecordList)

		// ============ Sprint 3.2 会员管理增强 ============
		// 积分引擎
		api.POST("/members/:id/points/add", memberEnhancedCtrl.AddPoints)
		api.POST("/members/:id/points/deduct", memberEnhancedCtrl.DeductPoints)
		api.GET("/members/:id/points/balance", memberEnhancedCtrl.GetBalance)
		api.GET("/members/:id/points/logs", memberEnhancedCtrl.GetPointsLogs)
		// 会员积分列表（/api/v1/members/points）
		api.GET("/members/points", memberEnhancedCtrl.ListMemberPoints)
		api.POST("/members/points/adjust", memberEnhancedCtrl.AdjustMemberPoints)
		api.GET("/members/:id/coupons", memberEnhancedCtrl.MemberCouponList)

		// 优惠券（新路径 /api/v1/coupons）
		api.GET("/coupons", memberEnhancedCtrl.CouponListNew)
		api.POST("/coupons", memberEnhancedCtrl.CouponCreateNew)
		api.POST("/coupons/:id/issue", memberEnhancedCtrl.CouponIssue)
		api.POST("/coupons/:id/use", memberEnhancedCtrl.CouponUse)

		// 促销活动（新路径 /api/v1/promotions）
		api.GET("/promotions", memberEnhancedCtrl.PromotionListNew)
		api.POST("/promotions", memberEnhancedCtrl.PromotionCreateNew)
		api.PUT("/promotions/:id", memberEnhancedCtrl.PromotionUpdateNew)
		api.DELETE("/promotions/:id", memberEnhancedCtrl.PromotionDeleteNew)
		api.GET("/promotions/:id", memberEnhancedCtrl.PromotionDetailNew)

		// ============ 应用管理 ============
		appCtrl := &AppController{DB: db}
		// 应用 CRUD
		api.GET("/apps", appCtrl.List)
		api.POST("/apps", appCtrl.Create)
		api.GET("/apps/:id", appCtrl.Get)
		api.PUT("/apps/:id", appCtrl.Update)
		api.DELETE("/apps/:id", appCtrl.Delete)
		// 版本管理
		api.GET("/apps/:id/versions", appCtrl.ListVersions)
		api.POST("/apps/:id/versions", appCtrl.CreateVersion)
		api.DELETE("/apps/:id/versions/:version_id", appCtrl.DeleteVersion)
		// 分发任务
		api.GET("/app/distributions", appCtrl.ListDistributions)
		api.POST("/app/distributions", appCtrl.CreateDistribution)
		api.GET("/app/distributions/:id", appCtrl.GetDistribution)
		api.POST("/app/distributions/:id/cancel", appCtrl.CancelDistribution)
		// 分发任务（/apps/distributions 路径别名，供前端使用）
		api.GET("/apps/distributions", appCtrl.ListDistributions)
		// 统计
		api.GET("/apps/:id/stats", appCtrl.GetStats)

		// ============ 策略管理 ============
		policyCtrl := &PolicyController{DB: db}
		complianceCtrl := &ComplianceController{DB: db}

		// 配置文件 CRUD
		api.GET("/policies/configs", policyCtrl.ListConfigs)
		api.POST("/policies/configs", policyCtrl.CreateConfig)
		api.PUT("/policies/configs/:id", policyCtrl.UpdateConfig)
		api.DELETE("/policies/configs/:id", policyCtrl.DeleteConfig)

		// 策略 CRUD
		api.GET("/policies", policyCtrl.ListPolicies)
		api.POST("/policies", policyCtrl.CreatePolicy)
		api.PUT("/policies/:id", policyCtrl.UpdatePolicy)
		api.DELETE("/policies/:id", policyCtrl.DeletePolicy)
		// 策略绑定
		api.POST("/policies/:id/bind", policyCtrl.BindPolicy)
		api.DELETE("/policies/:id/unbind", policyCtrl.UnbindPolicy)
		api.GET("/policies/:id/bindings", policyCtrl.GetPolicyBindings)

		// 合规规则 CRUD
		api.GET("/compliance/rules", complianceCtrl.ListRules)
		api.POST("/compliance/rules", complianceCtrl.CreateRule)
		api.PUT("/compliance/rules/:id", complianceCtrl.UpdateRule)
		api.DELETE("/compliance/rules/:id", complianceCtrl.DeleteRule)
		// 违规记录
		api.GET("/compliance/violations", complianceCtrl.ListViolations)
		api.PUT("/compliance/violations/:id/resolve", complianceCtrl.ResolveViolation)
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
	Status          string `form:"status"`           // online/offline
	LifecycleStatus int    `form:"lifecycle_status"` // 1:待激活 2:服役中 3:维修 4:报废
	HardwareModel   string `form:"hardware_model"`
	DeviceType      string `form:"device_type"`     // 设备类型筛选（M5Stack, ESP32等）
	OnlineStatus    string `form:"online_status"`   // 在线状态筛选: online/offline
	TenantID        string `form:"tenant_id"`       // 租户筛选
	Search          string `form:"search"`          // 关键词搜索（device_id/sn_code/mac_address）
	StartTime       string `form:"start_time"`     // 创建时间范围-开始
	EndTime         string `form:"end_time"`        // 创建时间范围-结束
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
	// device_type 是 hardware_model 的别名
	if req.DeviceType != "" {
		query = query.Where("hardware_model = ?", req.DeviceType)
	}
	if req.TenantID != "" {
		query = query.Where("tenant_id = ?", req.TenantID)
	}
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("device_id LIKE ? OR sn_code LIKE ? OR mac_address LIKE ?", search, search, search)
	}

	// 时间范围筛选
	if req.StartTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.StartTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if req.EndTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.EndTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
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

	// 如果有 status 或 online_status 筛选，内存过滤
	onlineStatusFilter := req.Status
	if req.OnlineStatus != "" {
		onlineStatusFilter = req.OnlineStatus
	}
	if onlineStatusFilter != "" {
		filtered := make([]DeviceWithShadow, 0)
		for _, d := range result {
			if onlineStatusFilter == "online" && d.IsOnline {
				filtered = append(filtered, d)
			} else if onlineStatusFilter == "offline" && !d.IsOnline {
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
