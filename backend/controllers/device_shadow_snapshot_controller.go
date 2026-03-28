package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeviceShadowSnapshotController 设备影子快照控制器
type DeviceShadowSnapshotController struct {
	DB *gorm.DB
}

// NewDeviceShadowSnapshotController 创建控制器
func NewDeviceShadowSnapshotController(db *gorm.DB) *DeviceShadowSnapshotController {
	return &DeviceShadowSnapshotController{DB: db}
}

// RegisterRoutes 注册路由
func (ctrl *DeviceShadowSnapshotController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/devices/:device_id/shadows/snapshots", ctrl.ListSnapshots)
	rg.GET("/devices/:device_id/shadows/snapshots/:snapshot_id", ctrl.GetSnapshot)
	rg.POST("/devices/:device_id/shadows/snapshots", ctrl.CreateSnapshot)
	rg.DELETE("/devices/:device_id/shadows/snapshots/:snapshot_id", ctrl.DeleteSnapshot)
	rg.GET("/devices/:device_id/shadows/snapshots/:snapshot_id/export", ctrl.ExportSnapshot)
	rg.POST("/devices/:device_id/shadows/snapshots/compare", ctrl.CompareSnapshots)
}

// ListSnapshots 获取快照列表
func (ctrl *DeviceShadowSnapshotController) ListSnapshots(c *gin.Context) {
	deviceID := c.Param("device_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	snapshotType := c.Query("type")

	query := ctrl.DB.Model(&models.DeviceShadowSnapshot{}).Where("device_id = ?", deviceID)
	if snapshotType != "" {
		query = query.Where("snapshot_type = ?", snapshotType)
	}

	var total int64
	var list []models.DeviceShadowSnapshot
	query.Count(&total)

	query.Order("created_at DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list)

	// 解析JSON字段
	for i := range list {
		if list[i].DesiredState != "" {
			var state map[string]interface{}
			json.Unmarshal([]byte(list[i].DesiredState), &state)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSnapshot 获取单个快照
func (ctrl *DeviceShadowSnapshotController) GetSnapshot(c *gin.Context) {
	deviceID := c.Param("device_id")
	snapshotID := c.Param("snapshot_id")

	var snapshot models.DeviceShadowSnapshot
	if err := ctrl.DB.Where("device_id = ? AND snapshot_id = ?", deviceID, snapshotID).First(&snapshot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "快照不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": snapshot})
}

// CreateSnapshot 创建快照
func (ctrl *DeviceShadowSnapshotController) CreateSnapshot(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req struct {
		SnapshotType string   `json:"snapshot_type"` // manual/auto/scheduled
		Reason       string   `json:"reason"`
		Tags         []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数可选
		req.SnapshotType = "manual"
	}

	// 获取当前设备影子 - 直接查询数据库
	var shadow struct {
		DeviceID        string `gorm:"column:device_id"`
		DesiredConfig   string `gorm:"column:desired_config"`
		DesiredVersion  string `gorm:"column:desired_version"`
	}
	if err := ctrl.DB.Table("device_shadows").Where("device_id = ?", deviceID).First(&shadow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "设备影子不存在"})
		return
	}

	// 生成快照ID
	snapshotID := fmt.Sprintf("SNAP-%s-%d", deviceID[:8], time.Now().Unix())

	// 计算版本号
	var maxVersion int
	ctrl.DB.Model(&models.DeviceShadowSnapshot{}).
		Where("device_id = ?", deviceID).
		Select("COALESCE(MAX(version),0)").
		Row().Scan(&maxVersion)

	// 计算状态差异（与上一版本）
	var delta string
	var stateDiff int
	var lastSnapshot models.DeviceShadowSnapshot
	if err := ctrl.DB.Where("device_id = ?", deviceID).Order("version DESC").First(&lastSnapshot).Error; err == nil {
		diff := computeDelta(lastSnapshot.DesiredState, shadow.DesiredConfig)
		delta = diff
		stateDiff = countDeltaFields(diff)
	}

	// 计算校验和
	checksum := sha256.Sum256([]byte(shadow.DesiredConfig))

	tags := ""
	if len(req.Tags) > 0 {
		data, _ := json.Marshal(req.Tags)
		tags = string(data)
	}

	snapshot := models.DeviceShadowSnapshot{
		DeviceID:        deviceID,
		SnapshotID:      snapshotID,
		Version:         maxVersion + 1,
		SnapshotType:    req.SnapshotType,
		Reason:         req.Reason,
		DesiredState:   shadow.DesiredConfig,
		DesiredVersion:  0, // device_shadows 没有 reported version
		Metadata:       fmt.Sprintf(`{"device_id":"%s","created_at":"%s"}`, deviceID, time.Now().Format(time.RFC3339)),
		Delta:          delta,
		Tags:           tags,
		StateDiff:      stateDiff,
		IsHealthy:      stateDiff == 0,
		Checksum:        hex.EncodeToString(checksum[:]),
		CreatedBy:       c.GetString("username"),
	}

	if err := ctrl.DB.Create(&snapshot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建快照失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "快照创建成功",
		"data": gin.H{
			"snapshot_id": snapshotID,
			"version":     snapshot.Version,
			"state_diff":  stateDiff,
		},
	})
}

// DeleteSnapshot 删除快照
func (ctrl *DeviceShadowSnapshotController) DeleteSnapshot(c *gin.Context) {
	deviceID := c.Param("device_id")
	snapshotID := c.Param("snapshot_id")

	var snapshot models.DeviceShadowSnapshot
	if err := ctrl.DB.Where("device_id = ? AND snapshot_id = ?", deviceID, snapshotID).First(&snapshot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "快照不存在"})
		return
	}

	ctrl.DB.Delete(&snapshot)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// ExportSnapshot 导出快照
func (ctrl *DeviceShadowSnapshotController) ExportSnapshot(c *gin.Context) {
	deviceID := c.Param("device_id")
	snapshotID := c.Param("snapshot_id")
	format := c.DefaultQuery("format", "json")

	var snapshot models.DeviceShadowSnapshot
	if err := ctrl.DB.Where("device_id = ? AND snapshot_id = ?", deviceID, snapshotID).First(&snapshot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "快照不存在"})
		return
	}

	var exportData []byte
	var contentType string

	switch format {
	case "json":
		data := gin.H{
			"snapshot_id":    snapshot.SnapshotID,
			"device_id":      snapshot.DeviceID,
			"version":        snapshot.Version,
			"desired_state":  snapshot.DesiredState,
			"reported_state": snapshot.ReportedState,
			"created_at":     snapshot.CreatedAt,
			"checksum":       snapshot.Checksum,
		}
		exportData, _ = json.MarshalIndent(data, "", "  ")
		contentType = "application/json"
	case "csv":
		exportData = []byte(fmt.Sprintf("field,value\ndesired_state,%s\nreported_state,%s\nversion,%d\n",
			snapshot.DesiredState, snapshot.ReportedState, snapshot.Version))
		contentType = "text/csv"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的格式"})
		return
	}

	// 记录导出
	export := models.DeviceShadowSnapshotExport{
		SnapshotID: snapshotID,
		DeviceID:   deviceID,
		Format:     format,
		FileSize:   int64(len(exportData)),
		CreatedBy:  c.GetString("username"),
		ExpiresAt:  func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
	}
	ctrl.DB.Create(&export)

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="snapshot_%s.%s"`, snapshotID, format))
	c.Data(http.StatusOK, contentType, exportData)
}

// CompareSnapshots 比较两个快照
func (ctrl *DeviceShadowSnapshotController) CompareSnapshots(c *gin.Context) {
	deviceID := c.Param("device_id")

	var req struct {
		SnapshotID1 string `json:"snapshot_id_1" binding:"required"`
		SnapshotID2 string `json:"snapshot_id_2" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var snap1, snap2 models.DeviceShadowSnapshot
	if err := ctrl.DB.Where("device_id = ? AND snapshot_id = ?", deviceID, req.SnapshotID1).First(&snap1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "快照1不存在"})
		return
	}
	if err := ctrl.DB.Where("device_id = ? AND snapshot_id = ?", deviceID, req.SnapshotID2).First(&snap2).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "快照2不存在"})
		return
	}

	compare := gin.H{
		"snapshot_1": gin.H{
			"snapshot_id":      snap1.SnapshotID,
			"version":          snap1.Version,
			"created_at":       snap1.CreatedAt,
			"desired_version":   snap1.DesiredVersion,
			"reported_version": snap1.ReportedVersion,
		},
		"snapshot_2": gin.H{
			"snapshot_id":      snap2.SnapshotID,
			"version":          snap2.Version,
			"created_at":       snap2.CreatedAt,
			"desired_version":   snap2.DesiredVersion,
			"reported_version": snap2.ReportedVersion,
		},
		"delta":        computeDelta(snap1.DesiredState, snap2.DesiredState),
		"desired_diff_count": countDeltaFields(computeDelta(snap1.DesiredState, snap2.DesiredState)),
		"reported_diff_count": countDeltaFields(computeDelta(snap1.ReportedState, snap2.ReportedState)),
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": compare})
}

// computeDelta 计算两个JSON状态的差异
func computeDelta(state1, state2 string) string {
	if state1 == "" && state2 == "" {
		return "{}"
	}
	var m1, m2 map[string]interface{}
	json.Unmarshal([]byte(state1), &m1)
	json.Unmarshal([]byte(state2), &m2)
	if m1 == nil {
		m1 = map[string]interface{}{}
	}
	if m2 == nil {
		m2 = map[string]interface{}{}
	}

	delta := map[string]interface{}{}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2) {
			delta[k] = map[string]interface{}{
				"from": v1,
				"to":   v2,
			}
		}
	}
	for k, v2 := range m2 {
		if _, ok := m1[k]; !ok {
			delta[k] = map[string]interface{}{
				"from": nil,
				"to":   v2,
			}
		}
	}

	result, _ := json.Marshal(delta)
	return string(result)
}

// countDeltaFields 统计差异字段数量
func countDeltaFields(delta string) int {
	var d map[string]interface{}
	json.Unmarshal([]byte(delta), &d)
	return len(d)
}
