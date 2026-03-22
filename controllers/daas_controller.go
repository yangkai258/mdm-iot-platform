package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"
	"mdm-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DaaSController 设备即服务（DaaS）租赁管理控制器
type DaaSController struct {
	DB    *gorm.DB
	Redis *utils.RedisClient
}

// parseTenantID parses tenant_id from context as uint
func parseTenantID(c *gin.Context) uint {
	tenantIDStr := getTenantID(c)
	tid, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return 0
	}
	return uint(tid)
}

// generateBillNo 生成账单编号
func generateBillNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%d%02d%02d%02d%02d%02d", prefix,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
}

// generateContractNo 生成合同编号
func generateContractNo() string {
	return fmt.Sprintf("DaaS%d%02d%02d%02d%02d%02d%d",
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), time.Now().Second(),
		time.Now().Nanosecond()%10000)
}

// generateRentalNo 生成租赁流水号
func generateRentalNo() string {
	return fmt.Sprintf("RTL%d%02d%02d%02d%02d%02d%d",
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), time.Now().Second(),
		time.Now().Nanosecond()%10000)
}

// ===================== 租赁合同管理 =====================

// ListContracts 获取租赁合同列表
// @Summary 获取租赁合同列表
// @Tags DaaS
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "合同状态"
// @Param device_id query int false "设备ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/contracts [get]
func (ctrl *DaaSController) ListContracts(c *gin.Context) {
	var contracts []models.DaaSContract
	var total int64

	tenantID := parseTenantID(c)
	userID := getUserID(c)

	query := ctrl.DB.Model(&models.DaaSContract{}).Where("tenant_id = ?", tenantID)

	// 按状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 按设备ID筛选
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	// 用户只能看自己的合同
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// 统计总数
	query.Count(&total)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  contracts,
			"total": total,
			"page":  page,
		},
	})
}

// CreateContract 创建租赁合同
// @Summary 创建租赁合同
// @Tags DaaS
// @Accept json
// @Produce json
// @Param body body models.DaaSContract true "合同信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/contracts [post]
func (ctrl *DaaSController) CreateContract(c *gin.Context) {
	var input struct {
		DeviceID       uint       `json:"device_id" binding:"required"`
		DeviceSN      string     `json:"device_sn"`
		DeviceName    string     `json:"device_name"`
		PlanName      string     `json:"plan_name"`
		DailyRate     float64    `json:"daily_rate" binding:"required"`
		MonthlyRate   float64    `json:"monthly_rate"`
		DepositAmount float64    `json:"deposit_amount"`
		ContractPeriod int        `json:"contract_period"` // 0表示不定期
		StartDate     *time.Time `json:"start_date"`
		EndDate       *time.Time `json:"end_date"`
		BillingCycle  string     `json:"billing_cycle"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误"})
		return
	}

	tenantID := parseTenantID(c)
	userID := getUserID(c)

	// 检查设备是否存在
	var device models.Device
	if err := ctrl.DB.First(&device, input.DeviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "设备不存在"})
		return
	}

	// 检查是否已有生效中的合同
	var existingCount int64
	ctrl.DB.Model(&models.DaaSContract{}).
		Where("device_id = ? AND status IN ?", input.DeviceID, []string{models.DaaSContractStatusActive, models.DaaSContractStatusPaused}).
		Count(&existingCount)
	if existingCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "error": "该设备已有生效中的租赁合同"})
		return
	}

	contract := models.DaaSContract{
		ContractNo:     generateContractNo(),
		TenantID:       tenantID,
		UserID:         userID,
		DeviceID:       input.DeviceID,
		DeviceSN:       input.DeviceSN,
		DeviceName:    input.DeviceName,
		PlanName:      input.PlanName,
		DailyRate:     input.DailyRate,
		MonthlyRate:   input.MonthlyRate,
		DepositAmount: input.DepositAmount,
		ContractPeriod: input.ContractPeriod,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		Status:        models.DaaSContractStatusActive,
		BillingCycle:  input.BillingCycle,
	}

	if contract.BillingCycle == "" {
		contract.BillingCycle = "monthly"
	}

	// 如果没有指定开始日期，使用当前时间
	if contract.StartDate == nil {
		now := time.Now()
		contract.StartDate = &now
	}

	// 如果有周期限制，计算结束日期
	if input.ContractPeriod > 0 && contract.EndDate == nil {
		endDate := contract.StartDate.AddDate(0, 0, input.ContractPeriod)
		contract.EndDate = &endDate
	}

	if err := ctrl.DB.Create(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "创建合同失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": contract,
	})
}

// GetContract 获取合同详情
// @Summary 获取合同详情
// @Tags DaaS
// @Accept json
// @Produce json
// @Param id path int true "合同ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/contracts/{id} [get]
func (ctrl *DaaSController) GetContract(c *gin.Context) {
	id := c.Param("id")
	var contract models.DaaSContract

	if err := ctrl.DB.First(&contract, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "合同不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询失败"})
		return
	}

	tenantID := parseTenantID(c)
	if contract.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "无权访问"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": contract,
	})
}

// UpdateContract 更新合同
// @Summary 更新合同
// @Tags DaaS
// @Accept json
// @Produce json
// @Param id path int true "合同ID"
// @Param body body map[string]interface{} true "更新字段"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/contracts/{id} [put]
func (ctrl *DaaSController) UpdateContract(c *gin.Context) {
	id := c.Param("id")
	var contract models.DaaSContract

	if err := ctrl.DB.First(&contract, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "合同不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询失败"})
		return
	}

	tenantID := parseTenantID(c)
	if contract.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "无权访问"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误"})
		return
	}

	// 只允许更新部分字段
	allowedFields := []string{"plan_name", "daily_rate", "monthly_rate", "deposit_amount",
		"contract_period", "end_date", "billing_cycle"}
	updateFields := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := input[field]; ok {
			updateFields[field] = val
		}
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "没有可更新的字段"})
		return
	}

	if err := ctrl.DB.Model(&contract).Updates(updateFields).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "更新失败"})
		return
	}

	// 重新加载
	ctrl.DB.First(&contract, id)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
		"data": contract,
	})
}

// TerminateContract 终止合同
// @Summary 终止合同
// @Tags DaaS
// @Accept json
// @Produce json
// @Param id path int true "合同ID"
// @Param body body map[string]interface{} true "终止原因"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/contracts/{id}/terminate [post]
func (ctrl *DaaSController) TerminateContract(c *gin.Context) {
	id := c.Param("id")
	var contract models.DaaSContract

	if err := ctrl.DB.First(&contract, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "合同不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询失败"})
		return
	}

	tenantID := parseTenantID(c)
	if contract.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "无权访问"})
		return
	}

	if contract.Status == models.DaaSContractStatusTerminated {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "合同已终止"})
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&input)

	userID := getUserID(c)
	now := time.Now()

	updates := map[string]interface{}{
		"status":           models.DaaSContractStatusTerminated,
		"terminate_reason": input.Reason,
		"terminated_at":    now,
		"terminated_by":    userID,
	}

	if err := ctrl.DB.Model(&contract).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "终止失败"})
		return
	}

	// 归还设备（创建归还记录）
	rental := models.DaaSDeviceRental{
		RentalNo:       generateRentalNo(),
		ContractID:     contract.ID,
		TenantID:       tenantID,
		UserID:         userID,
		DeviceID:       contract.DeviceID,
		DeviceSN:       contract.DeviceSN,
		DeviceName:     contract.DeviceName,
		Action:         models.DaaSRentalActionReturn,
		ReturnTime:     &now,
		Status:         models.DaaSRentalStatusReturned,
		DepositAmount:  contract.DepositAmount,
	}
	ctrl.DB.Create(&rental)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "合同已终止",
	})
}

// ===================== 设备租赁管理 =====================

// ListDaasDevices 获取可租赁的设备列表
// @Summary 获取可租赁的设备列表
// @Tags DaaS
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "租赁状态: rented/available"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/devices [get]
func (ctrl *DaaSController) ListDaasDevices(c *gin.Context) {
	tenantID := parseTenantID(c)
	status := c.Query("status")

	// 获取已租出的设备
	var rentals []models.DaaSDeviceRental
	var total int64
	var result []map[string]interface{}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if status == "available" {
		// 可用设备（未被租用的设备）
		subQuery := ctrl.DB.Model(&models.DaaSDeviceRental{}).
			Select("device_id").
			Where("tenant_id = ? AND status = ?", tenantID, models.DaaSRentalStatusRented)

		ctrl.DB.Model(&models.Device{}).Where("tenant_id = ?", tenantID).Count(&total)

		var devices []models.Device
		ctrl.DB.Where("tenant_id = ? AND id NOT IN ?", tenantID, subQuery).
			Offset(offset).Limit(pageSize).Find(&devices)

		for _, d := range devices {
			result = append(result, map[string]interface{}{
				"device_id":   d.ID,
				"device_sn":   d.SnCode,
				"device_name": d.DeviceID,
				"status":      "available",
			})
		}
	} else {
		// 已租用的设备
		query := ctrl.DB.Model(&models.DaaSDeviceRental{}).Where("tenant_id = ?", tenantID)
		if status == "rented" {
			query = query.Where("status = ?", models.DaaSRentalStatusRented)
		}
		query.Count(&total)
		query.Offset(offset).Limit(pageSize).Find(&rentals)

		for _, r := range rentals {
			result = append(result, map[string]interface{}{
				"rental_id":   r.ID,
				"device_id":   r.DeviceID,
				"device_sn":   r.DeviceSN,
				"device_name": r.DeviceName,
				"contract_id": r.ContractID,
				"rent_time":   r.RentTime,
				"status":      r.Status,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  result,
			"total": total,
			"page":  page,
		},
	})
}

// RentDevice 租用设备
// @Summary 租用设备
// @Tags DaaS
// @Accept json
// @Produce json
// @Param device_id path int true "设备ID"
// @Param body body map[string]interface{} true "租赁信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/devices/{device_id}/rent [post]
func (ctrl *DaaSController) RentDevice(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的设备ID"})
		return
	}

	tenantID := parseTenantID(c)
	userID := getUserID(c)

	var input struct {
		ContractID     uint       `json:"contract_id"`
		ExpectedReturn *time.Time `json:"expected_return"`
		Notes         string     `json:"notes"`
	}
	c.ShouldBindJSON(&input)

	// 检查设备是否存在
	var device models.Device
	if err := ctrl.DB.First(&device, deviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "设备不存在"})
		return
	}

	// 检查是否已被租用
	var activeRental models.DaaSDeviceRental
	notReturned := ctrl.DB.Where("device_id = ? AND status = ?", deviceID, models.DaaSRentalStatusRented).
		First(&activeRental)
	if notReturned.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "error": "设备已被租用"})
		return
	}

	// 如果指定了合同，检查合同有效性
	var contract models.DaaSContract
	if input.ContractID > 0 {
		if err := ctrl.DB.First(&contract, input.ContractID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "合同不存在"})
			return
		}
		if contract.Status != models.DaaSContractStatusActive {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "合同状态无效"})
			return
		}
	}

	now := time.Now()
	rental := models.DaaSDeviceRental{
		RentalNo:        generateRentalNo(),
		ContractID:      input.ContractID,
		TenantID:       tenantID,
		UserID:         userID,
		DeviceID:       uint(deviceID),
		DeviceSN:       device.SnCode,
		DeviceName:     device.DeviceID,
		Action:         models.DaaSRentalActionRent,
		RentTime:       &now,
		ExpectedReturn: input.ExpectedReturn,
		Status:         models.DaaSRentalStatusRented,
		Notes:          input.Notes,
	}

	if contract.ID > 0 {
		rental.DepositAmount = contract.DepositAmount
	}

	if err := ctrl.DB.Create(&rental).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "租用失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "租用成功",
		"data": rental,
	})
}

// ReturnDevice 归还设备
// @Summary 归还设备
// @Tags DaaS
// @Accept json
// @Produce json
// @Param device_id path int true "设备ID"
// @Param body body map[string]interface{} true "归还信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/devices/{device_id}/return [post]
func (ctrl *DaaSController) ReturnDevice(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "无效的设备ID"})
		return
	}

	tenantID := parseTenantID(c)
	userID := getUserID(c)

	var input struct {
		Notes string `json:"notes"`
	}
	c.ShouldBindJSON(&input)

	// 查找当前租用记录
	var rental models.DaaSDeviceRental
	if err := ctrl.DB.Where("device_id = ? AND status = ?", deviceID, models.DaaSRentalStatusRented).
		First(&rental).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "没有找到租用记录"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询失败"})
		return
	}

	now := time.Now()
	rental.ReturnTime = &now
	rental.Status = models.DaaSRentalStatusReturned
	if input.Notes != "" {
		rental.Notes = input.Notes
	}

	if err := ctrl.DB.Save(&rental).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "归还失败"})
		return
	}

	// 如果有押金，创建押金退还账单
	if rental.DepositAmount > 0 {
		billing := models.DaaSBilling{
			BillNo:         generateBillNo("DaaSDEP"),
			TenantID:       tenantID,
			UserID:         userID,
			ContractID:     rental.ContractID,
			RentalID:       rental.ID,
			BillType:       models.DaaSBillTypeDepositRefund,
			DepositAmount:  rental.DepositAmount,
			Amount:         rental.DepositAmount,
			Currency:       "CNY",
			Status:         models.DaaSBillStatusPending,
			Description:    "设备押金退还",
		}
		ctrl.DB.Create(&billing)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "归还成功",
		"data": rental,
	})
}

// ===================== 租赁账单管理 =====================

// ListBillings 获取租赁账单列表
// @Summary 获取租赁账单列表
// @Tags DaaS
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param status query string false "账单状态"
// @Param contract_id query int false "合同ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/billing [get]
func (ctrl *DaaSController) ListBillings(c *gin.Context) {
	var billings []models.DaaSBilling
	var total int64

	tenantID := parseTenantID(c)
	userID := getUserID(c)

	query := ctrl.DB.Model(&models.DaaSBilling{}).Where("tenant_id = ?", tenantID)

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if contractID := c.Query("contract_id"); contractID != "" {
		query = query.Where("contract_id = ?", contractID)
	}

	query.Count(&total)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&billings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  billings,
			"total": total,
			"page":  page,
		},
	})
}

// CalculateBilling 计算租赁费用
// @Summary 计算租赁费用
// @Tags DaaS
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "计费参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/daas/billing/calculate [post]
func (ctrl *DaaSController) CalculateBilling(c *gin.Context) {
	var input struct {
		ContractID  uint      `json:"contract_id" binding:"required"`
		RentalID   uint      `json:"rental_id"`
		PeriodStart time.Time `json:"period_start" binding:"required"`
		PeriodEnd   time.Time `json:"period_end" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "参数错误"})
		return
	}

	tenantID := parseTenantID(c)

	// 获取合同信息
	var contract models.DaaSContract
	if err := ctrl.DB.First(&contract, input.ContractID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "合同不存在"})
		return
	}

	if contract.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "error": "无权访问"})
		return
	}

	// 计算天数
	days := int(input.PeriodEnd.Sub(input.PeriodStart).Hours()/24) + 1

	// 计算租金
	var amount float64
	if contract.BillingCycle == "daily" {
		amount = float64(days) * contract.DailyRate
	} else {
		// 按月计费：计算完整月数 + 剩余天数
		months := days / 30
		remainingDays := days % 30
		amount = float64(months)*contract.MonthlyRate + float64(remainingDays)*contract.DailyRate
	}

	// 检查是否有逾期
	overdueDays := 0
	if input.RentalID > 0 {
		var rental models.DaaSDeviceRental
		if err := ctrl.DB.First(&rental, input.RentalID).Error; err == nil {
			if rental.ExpectedReturn != nil && rental.ExpectedReturn.Before(time.Now()) && rental.Status == models.DaaSRentalStatusRented {
				overdueDays = int(time.Since(*rental.ExpectedReturn).Hours() / 24)
			}
		}
	}

	penalty := 0.0
	if overdueDays > 0 {
		// 逾期每天收取日租金的 50% 作为违约金
		penalty = float64(overdueDays) * contract.DailyRate * 0.5
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"contract_id":   contract.ID,
			"rental_id":    input.RentalID,
			"period_start": input.PeriodStart,
			"period_end":   input.PeriodEnd,
			"days":         days,
			"daily_rate":   contract.DailyRate,
			"rental_amount": amount,
			"overdue_days": overdueDays,
			"penalty":      penalty,
			"deposit":      contract.DepositAmount,
			"total_amount": amount + penalty,
			"currency":     "CNY",
		},
	})
}
