package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImportExportController 数据导入导出控制器
type ImportExportController struct {
	DB *gorm.DB
}

// ============ 设备导入导出 ============

// ExportDevices 导出设备列表
func (c *ImportExportController) ExportDevices(ctx *gin.Context) {
	var devices []models.Device
	query := c.DB.Model(&models.Device{})

	// 应用现有筛选条件
	if deviceType := ctx.Query("device_type"); deviceType != "" {
		query = query.Where("hardware_model = ?", deviceType)
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("lifecycle_status = ?", status)
	}
	if search := ctx.Query("search"); search != "" {
		search = "%" + search + "%"
		query = query.Where("device_id LIKE ? OR sn_code LIKE ? OR mac_address LIKE ?", search, search, search)
	}
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	if err := query.Order("created_at DESC").Find(&devices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出失败"})
		return
	}

	// 设置CSV头
	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=devices_%s.csv", time.Now().Format("20060102150405")))

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	// 写入表头
	header := []string{"设备ID", "MAC地址", "序列号", "硬件型号", "固件版本", "生命周期状态", "绑定用户ID", "组织ID", "创建时间"}
	writer.Write(header)

	// 写入数据
	for _, d := range devices {
		statusMap := map[int]string{1: "待激活", 2: "服役中", 3: "维修", 4: "报废"}
		row := []string{
			d.DeviceID,
			d.MacAddress,
			d.SnCode,
			d.HardwareModel,
			d.FirmwareVersion,
			statusMap[d.LifecycleStatus],
			stringPtrToString(d.BindUserID),
			strconv.FormatUint(uint64(d.OrgID), 10),
			d.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}
}

// ImportDevices 导入设备
func (c *ImportExportController) ImportDevices(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请上传文件"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// 读取表头
	header, err := reader.Read()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件格式错误"})
		return
	}

	// 建立表头索引
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	var successCount, failCount int
	var errors []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			failCount++
			continue
		}

		device := models.Device{}
		
		if idx, ok := colIndex["设备ID"]; ok && idx < len(record) {
			device.DeviceID = record[idx]
		}
		if idx, ok := colIndex["MAC地址"]; ok && idx < len(record) {
			device.MacAddress = record[idx]
		}
		if idx, ok := colIndex["序列号"]; ok && idx < len(record) {
			device.SnCode = record[idx]
		}
		if idx, ok := colIndex["硬件型号"]; ok && idx < len(record) {
			device.HardwareModel = record[idx]
		}
		if idx, ok := colIndex["固件版本"]; ok && idx < len(record) {
			device.FirmwareVersion = record[idx]
		}

		if err := c.DB.Create(&device).Error; err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("行%d: %v", successCount+failCount+1, err))
		} else {
			successCount++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "导入完成",
		"data": gin.H{
			"success": successCount,
			"fail":    failCount,
			"errors":  errors,
		},
	})
}

// ============ 会员导入导出 ============

// ExportMembers 导出会员列表
func (c *ImportExportController) ExportMembers(ctx *gin.Context) {
	var members []models.Member
	query := c.DB.Model(&models.Member{})

	// 应用筛选条件
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("member_name LIKE ? OR member_code LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if level := ctx.Query("member_level"); level != "" {
		query = query.Where("member_level = ?", level)
	}
	if pointsMin := ctx.Query("points_min"); pointsMin != "" {
		if p, err := strconv.ParseInt(pointsMin, 10, 64); err == nil {
			query = query.Where("points >= ?", p)
		}
	}
	if pointsMax := ctx.Query("points_max"); pointsMax != "" {
		if p, err := strconv.ParseInt(pointsMax, 10, 64); err == nil {
			query = query.Where("points <= ?", p)
		}
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("Card").Order("created_at DESC").Find(&members).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出失败"})
		return
	}

	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=members_%s.csv", time.Now().Format("20060102150405")))

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	header := []string{"会员编号", "会员名称", "手机号", "性别", "邮箱", "会员等级", "积分", "储值余额", "状态", "来源", "创建时间"}
	writer.Write(header)

	for _, m := range members {
		statusMap := map[int]string{1: "正常", 2: "禁用", 0: "已删除"}
		row := []string{
			m.MemberCode,
			m.MemberName,
			m.Phone,
			m.Gender,
			m.Email,
			strconv.Itoa(m.MemberLevel),
			strconv.FormatInt(m.Points, 10),
			strconv.FormatFloat(m.Balance, 'f', 2, 64),
			statusMap[m.Status],
			m.Source,
			m.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}
}

// ImportMembers 导入会员
func (c *ImportExportController) ImportMembers(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请上传文件"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文件格式错误"})
		return
	}

	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[col] = i
	}

	var successCount, failCount int
	var errors []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			failCount++
			continue
		}

		member := models.Member{}

		if idx, ok := colIndex["会员编号"]; ok && idx < len(record) {
			member.MemberCode = record[idx]
		}
		if idx, ok := colIndex["会员名称"]; ok && idx < len(record) {
			member.MemberName = record[idx]
		}
		if idx, ok := colIndex["手机号"]; ok && idx < len(record) {
			member.Phone = record[idx]
		}
		if idx, ok := colIndex["性别"]; ok && idx < len(record) {
			member.Gender = record[idx]
		}
		if idx, ok := colIndex["邮箱"]; ok && idx < len(record) {
			member.Email = record[idx]
		}
		if idx, ok := colIndex["会员等级"]; ok && idx < len(record) {
			if level, err := strconv.Atoi(record[idx]); err == nil {
				member.MemberLevel = level
			}
		}
		if idx, ok := colIndex["积分"]; ok && idx < len(record) {
			if points, err := strconv.ParseInt(record[idx], 10, 64); err == nil {
				member.Points = points
			}
		}
		if idx, ok := colIndex["储值余额"]; ok && idx < len(record) {
			if balance, err := strconv.ParseFloat(record[idx], 64); err == nil {
				member.Balance = balance
			}
		}

		if err := c.DB.Create(&member).Error; err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("行%d: %v", successCount+failCount+1, err))
		} else {
			successCount++
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "导入完成",
		"data": gin.H{
			"success": successCount,
			"fail":    failCount,
			"errors":  errors,
		},
	})
}

// ============ 活动日志导出 ============

// ExportActivityLogs 导出活动日志
func (c *ImportExportController) ExportActivityLogs(ctx *gin.Context) {
	var logs []models.ActivityLog
	query := c.DB.Model(&models.ActivityLog{})

	// 筛选条件
	if username := ctx.Query("username"); username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action := ctx.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType := ctx.Query("resource_type"); resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if startTime := ctx.Query("start_time"); startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime := ctx.Query("end_time"); endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	if err := query.Order("created_at DESC").Limit(10000).Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "导出失败"})
		return
	}

	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=activity_logs_%s.csv", time.Now().Format("20060102150405")))

	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	header := []string{"ID", "用户ID", "用户名", "动作", "资源类型", "资源ID", "资源名称", "IP地址", "时间"}
	writer.Write(header)

	for _, l := range logs {
		row := []string{
			strconv.FormatUint(uint64(l.ID), 10),
			strconv.FormatUint(uint64(l.UserID), 10),
			l.Username,
			l.Action,
			l.ResourceType,
			strconv.FormatUint(uint64(l.ResourceID), 10),
			l.ResourceName,
			l.IP,
			l.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}
}

// ============ 辅助函数 ============

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// BatchStatusRequest 批量更新状态请求
type BatchStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required"`
	Status int    `json:"status"`
}

// ============ 批量操作实现 ============

// BatchDeleteDevices 批量删除设备
func (c *ImportExportController) BatchDeleteDevices(ctx *gin.Context) {
	var req BatchDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要删除的设备"})
		return
	}

	result := c.DB.Where("id IN ?", req.IDs).Delete(&models.Device{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data": gin.H{
			"deleted": result.RowsAffected,
		},
	})
}

// BatchUpdateDeviceStatus 批量更新设备状态
func (c *ImportExportController) BatchUpdateDeviceStatus(ctx *gin.Context) {
	var req BatchStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要更新的设备"})
		return
	}

	if req.Status < 1 || req.Status > 4 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "状态值无效（1:待激活 2:服役中 3:维修 4:报废）"})
		return
	}

	result := c.DB.Model(&models.Device{}).Where("id IN ?", req.IDs).Update("lifecycle_status", req.Status)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data": gin.H{
			"updated": result.RowsAffected,
		},
	})
}

// BatchDeleteMembers 批量删除会员
func (c *ImportExportController) BatchDeleteMembers(ctx *gin.Context) {
	var req BatchDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请选择要删除的会员"})
		return
	}

	result := c.DB.Where("id IN ?", req.IDs).Delete(&models.Member{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data": gin.H{
			"deleted": result.RowsAffected,
		},
	})
}
