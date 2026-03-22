package controllers

import (
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LostFoundCtrl 寻回网络控制器
type LostFoundCtrl struct {
	DB *gorm.DB
}

// RegisterLostFoundRoutes 注册寻回网络路由
func (l *LostFoundCtrl) RegisterLostFoundRoutes(r *gin.RouterGroup) {
	r.GET("/lost-found/reports", l.ListReports)
	r.POST("/lost-found/reports", l.CreateReport)
	r.GET("/lost-found/reports/:id", l.GetReport)
	r.PUT("/lost-found/reports/:id", l.UpdateReport)
	r.DELETE("/lost-found/reports/:id", l.CloseReport)
	r.POST("/lost-found/reports/:id/sighting", l.AddSighting)
	r.GET("/lost-found/reports/:id/sightings", l.ListSightings)
	r.GET("/lost-found/nearby", l.NearbyReports)
}

// ListReports 获取失宠报告列表
func (l *LostFoundCtrl) ListReports(c *gin.Context) {
	tenantID := getTenantID(c)

	var reports []models.LostPet
	query := l.DB.Where("tenant_id = ?", tenantID)

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	} else {
		// 默认只显示 searching 状态
		query = query.Where("status = ?", "searching")
	}

	// 物种过滤
	if species := c.Query("species"); species != "" {
		query = query.Where("species = ?", species)
	}

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

	var total int64
	l.DB.Model(&models.LostPet{}).Where("tenant_id = ? AND status = ?", tenantID, "searching").Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取列表失败"})
		return
	}

	// 脱敏处理：隐藏精确位置，只对发布者可见
	userID := getUserID(c)
	sanitized := make([]models.LostPet, len(reports))
	for i, r := range reports {
		sanitized[i] = r
		if r.ReporterID != userID {
			// 模糊化位置（保留大概区域）
			sanitized[i].LastLocation = blurLocation(r.LastLocation)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": sanitized,
			"total": total,
			"page": page,
			"page_size": pageSize,
		},
	})
}

// blurLocation 模糊化位置信息（大概区域）
func blurLocation(loc models.JSON) models.JSON {
	if loc == nil {
		return nil
	}
	// 只保留一位小数，降低精度到约10km范围
	if lat, ok := loc["lat"].(float64); ok {
		loc["lat"] = math.Round(lat*10) / 10
	}
	if lng, ok := loc["lng"].(float64); ok {
		loc["lng"] = math.Round(lng*10) / 10
	}
	return loc
}

// CreateReport 创建失宠报告
func (l *LostFoundCtrl) CreateReport(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		PetUUID      string   `json:"pet_uuid"`
		PetName      string   `json:"pet_name" binding:"required"`
		Species      string   `json:"species" binding:"required"`
		LastLocation models.JSON `json:"last_location" binding:"required"`
		LostTime     string   `json:"lost_time" binding:"required"`
		Reward       string   `json:"reward"`
		Contact      string   `json:"contact"`
		Description  string   `json:"description"`
		PhotoURLs    []string `json:"photo_urls"`
		SpreadRadius float64  `json:"spread_radius_km"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	lostTime, err := time.Parse(time.RFC3339, input.LostTime)
	if err != nil {
		// 尝试日期格式
		lostTime, _ = time.Parse("2006-01-02T15:04:05Z", input.LostTime)
	}

	spreadRadius := input.SpreadRadius
	if spreadRadius <= 0 {
		spreadRadius = 10 // 默认10km
	}
	// 限制扩散范围不超过50km
	if spreadRadius > 50 {
		spreadRadius = 50
	}

	report := models.LostPet{
		PetUUID:      input.PetUUID,
		PetName:      input.PetName,
		Species:      input.Species,
		LastLocation: input.LastLocation,
		LostTime:     lostTime,
		Status:       "searching",
		Reward:       input.Reward,
		Contact:      input.Contact,
		Description:  input.Description,
		PhotoURLs:     input.PhotoURLs,
		ReporterID:   userID,
		SpreadRadius: spreadRadius,
		TenantID:     tenantID,
	}

	if err := l.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发布成功", "data": report})
}

// GetReport 获取报告详情
func (l *LostFoundCtrl) GetReport(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var report models.LostPet
	err := l.DB.Where("report_uuid = ? AND tenant_id = ?", reportID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取详情失败"})
		return
	}

	// 位置隐私：精确位置只对发布者可见
	if report.ReporterID != userID {
		report.LastLocation = blurLocation(report.LastLocation)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// UpdateReport 更新报告（仅发布者可更新）
func (l *LostFoundCtrl) UpdateReport(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var report models.LostPet
	err := l.DB.Where("report_uuid = ? AND reporter_id = ? AND tenant_id = ?", reportID, userID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在或无权修改"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	var input struct {
		PetName      string   `json:"pet_name"`
		Status       string   `json:"status"`
		Reward       string   `json:"reward"`
		Contact      string   `json:"contact"`
		Description  string   `json:"description"`
		LastLocation models.JSON `json:"last_location"`
		PhotoURLs    []string `json:"photo_urls"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if input.PetName != "" { updates["pet_name"] = input.PetName }
	if input.Status != "" { updates["status"] = input.Status }
	if input.Reward != "" { updates["reward"] = input.Reward }
	if input.Contact != "" { updates["contact"] = input.Contact }
	if input.Description != "" { updates["description"] = input.Description }
	if input.LastLocation != nil { updates["last_location"] = input.LastLocation }
	if len(input.PhotoURLs) > 0 { updates["photo_urls"] = input.PhotoURLs }

	if err := l.DB.Model(&report).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	l.DB.First(&report, report.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": report})
}

// CloseReport 关闭报告（软删除/标记关闭）
func (l *LostFoundCtrl) CloseReport(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var report models.LostPet
	err := l.DB.Where("report_uuid = ? AND reporter_id = ? AND tenant_id = ?", reportID, userID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在或无权操作"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if err := l.DB.Model(&report).Update("status", "closed").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "关闭报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "报告已关闭"})
}

// AddSighting 添加目击记录
func (l *LostFoundCtrl) AddSighting(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	// 验证报告存在且正在搜索
	var report models.LostPet
	err := l.DB.Where("report_uuid = ? AND tenant_id = ? AND status = ?", reportID, tenantID, "searching").First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在或已关闭"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	var input struct {
		Location      models.JSON `json:"location" binding:"required"`
		SightingTime  string      `json:"sighting_time" binding:"required"`
		Description   string      `json:"description"`
		PhotoURL      string      `json:"photo_url"`
		ReporterName  string      `json:"reporter_name"`
		Contact       string      `json:"contact"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sightingTime, err := time.Parse(time.RFC3339, input.SightingTime)
	if err != nil {
		sightingTime, _ = time.Parse("2006-01-02T15:04:05Z", input.SightingTime)
	}

	sighting := models.PetSighting{
		ReportUUID:   reportID,
		Location:     input.Location,
		SightingTime: sightingTime,
		Description:  input.Description,
		PhotoURL:     input.PhotoURL,
		ReporterName: input.ReporterName,
		Contact:      input.Contact,
		ReporterID:   userID,
		TenantID:     tenantID,
	}

	if err := l.DB.Create(&sighting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "添加目击记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "目击记录已提交", "data": sighting})
}

// ListSightings 获取目击记录列表
func (l *LostFoundCtrl) ListSightings(c *gin.Context) {
	reportID := c.Param("id")
	tenantID := getTenantID(c)

	// 验证报告存在
	var report models.LostPet
	err := l.DB.Where("report_uuid = ? AND tenant_id = ?", reportID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}

	var sightings []models.PetSighting
	query := l.DB.Where("report_uuid = ?", reportID)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if err := query.Order("sighting_time DESC").Offset(offset).Limit(pageSize).Find(&sightings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取目击记录失败"})
		return
	}

	// 脱敏：目击位置也做模糊处理
	userID := getUserID(c)
	sanitized := make([]models.PetSighting, len(sightings))
	for i, s := range sightings {
		sanitized[i] = s
		if report.ReporterID != userID {
			sanitized[i].Location = blurLocation(s.Location)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list": sanitized,
			"page": page,
			"page_size": pageSize,
		},
	})
}

// NearbyReports 附近失宠报告（按位置半径查询）
func (l *LostFoundCtrl) NearbyReports(c *gin.Context) {
	tenantID := getTenantID(c)

	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "10") // 默认10km

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的纬度"})
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的经度"})
		return
	}
	radius, _ := strconv.ParseFloat(radiusStr, 64)
	if radius <= 0 || radius > 50 {
		radius = 10
	}

	var reports []models.LostPet
	if err := l.DB.Where("tenant_id = ? AND status = ?", tenantID, "searching").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	// 简单的距离过滤（使用 Haversine 公式）
	userID := getUserID(c)
	var filtered []models.LostPet
	for _, r := range reports {
		if r.LastLocation == nil {
			continue
		}
		repLat, ok1 := r.LastLocation["lat"].(float64)
		repLng, ok2 := r.LastLocation["lng"].(float64)
		if !ok1 || !ok2 {
			continue
		}

		dist := haversineDistance(lat, lng, repLat, repLng)
		if dist <= radius {
			reportCopy := r
			if reportCopy.ReporterID != userID {
				reportCopy.LastLocation = blurLocation(reportCopy.LastLocation)
			}
			reportCopy.LastLocation["distance_km"] = math.Round(dist*100) / 100
			filtered = append(filtered, reportCopy)
		}
	}

	// 按距离排序
	sort.Slice(filtered, func(i, j int) bool {
		di := 100.0
		if d, ok := filtered[i].LastLocation["distance_km"].(float64); ok {
			di = d
		}
		dj := 100.0
		if d, ok := filtered[j].LastLocation["distance_km"].(float64); ok {
			dj = d
		}
		return di < dj
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":          filtered,
			"search_center": gin.H{"lat": lat, "lng": lng},
			"radius_km":     radius,
		},
	})
}

// haversineDistance 计算两点间的球面距离（km）
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // 地球半径 km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
