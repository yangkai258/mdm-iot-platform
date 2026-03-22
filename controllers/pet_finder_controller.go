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

// PetFinderCtrl 宠物寻回网络控制器
type PetFinderCtrl struct {
	DB *gorm.DB
}

// RegisterPetFinderRoutes 注册宠物寻回网络路由
func (p *PetFinderCtrl) RegisterPetFinderRoutes(r *gin.RouterGroup) {
	r.POST("/pet-finder/reports", p.CreateReport)
	r.GET("/pet-finder/reports", p.ListReports)
	r.GET("/pet-finder/reports/:id", p.GetReport)
	r.PUT("/pet-finder/reports/:id/status", p.UpdateStatus)
	r.POST("/pet-finder/reports/:id/sightings", p.AddSighting)
	r.GET("/pet-finder/reports/:id/sightings", p.ListSightings)
	r.GET("/pet-finder/nearby", p.NearbyReports)
	r.POST("/pet-finder/alerts", p.ManageAlert)
	r.GET("/pet-finder/alerts", p.ListAlerts)
	r.DELETE("/pet-finder/alerts/:id", p.DeleteAlert)
}

// CreateReport 创建走失报告
func (p *PetFinderCtrl) CreateReport(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		PetUUID       string   `json:"pet_uuid"`
		PetName       string   `json:"pet_name" binding:"required"`
		Species       string   `json:"species" binding:"required"`
		Breed         string   `json:"breed"`
		Color         string   `json:"color"`
		Gender        string   `json:"gender"`
		Age           string   `json:"age"`
		LastLocation  models.JSON `json:"last_location" binding:"required"`
		LostTime      string   `json:"lost_time" binding:"required"`
		Reward        string   `json:"reward"`
		ContactName   string   `json:"contact_name"`
		ContactPhone  string   `json:"contact_phone"`
		ContactWechat string   `json:"contact_wechat"`
		Description   string   `json:"description"`
		PhotoURLs     []string `json:"photo_urls"`
		SpreadRadius  float64  `json:"spread_radius_km"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	lostTime, err := time.Parse(time.RFC3339, input.LostTime)
	if err != nil {
		lostTime, _ = time.Parse("2006-01-02T15:04:05Z", input.LostTime)
	}

	spreadRadius := input.SpreadRadius
	if spreadRadius <= 0 {
		spreadRadius = 10
	}
	if spreadRadius > 50 {
		spreadRadius = 50
	}

	report := models.PetLostReport{
		PetUUID:       input.PetUUID,
		PetName:       input.PetName,
		Species:       input.Species,
		Breed:         input.Breed,
		Color:         input.Color,
		Gender:        input.Gender,
		Age:           input.Age,
		LastLocation:  input.LastLocation,
		LostTime:      lostTime,
		Status:        "searching",
		Reward:        input.Reward,
		ContactName:   input.ContactName,
		ContactPhone:  input.ContactPhone,
		ContactWechat: input.ContactWechat,
		Description:   input.Description,
		PhotoURLs:     input.PhotoURLs,
		ReporterID:    userID,
		SpreadRadius:  spreadRadius,
		TenantID:      tenantID,
	}

	if err := p.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建报告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发布成功", "data": report})
}

// ListReports 获取走失报告列表
func (p *PetFinderCtrl) ListReports(c *gin.Context) {
	tenantID := getTenantID(c)

	var reports []models.PetLostReport
	query := p.DB.Where("tenant_id = ?", tenantID)

	// 状态过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	} else {
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
	p.DB.Model(&models.PetLostReport{}).Where("tenant_id = ?", tenantID).Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取列表失败"})
		return
	}

	// 脱敏处理
	userID := getUserID(c)
	sanitized := make([]models.PetLostReport, len(reports))
	for i, r := range reports {
		sanitized[i] = r
		if r.ReporterID != userID {
			sanitized[i].LastLocation = blurPetFinderLocation(r.LastLocation)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      sanitized,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetReport 获取报告详情
func (p *PetFinderCtrl) GetReport(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var report models.PetLostReport
	err := p.DB.Where("report_uuid = ? AND tenant_id = ?", reportID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取详情失败"})
		return
	}

	// 位置隐私
	if report.ReporterID != userID {
		report.LastLocation = blurPetFinderLocation(report.LastLocation)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": report})
}

// UpdateStatus 更新报告状态
func (p *PetFinderCtrl) UpdateStatus(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var report models.PetLostReport
	err := p.DB.Where("report_uuid = ? AND tenant_id = ?", reportID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	// 权限检查：只有发布者可以更新状态
	if report.ReporterID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作此报告"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 状态值校验
	validStatuses := map[string]bool{"searching": true, "found": true, "closed": true, "abandoned": true}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效的状态值"})
		return
	}

	if err := p.DB.Model(&report).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新状态失败"})
		return
	}

	report.Status = input.Status
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "状态已更新", "data": report})
}

// AddSighting 报告目击
func (p *PetFinderCtrl) AddSighting(c *gin.Context) {
	reportID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	// 验证报告存在且正在搜索
	var report models.PetLostReport
	err := p.DB.Where("report_uuid = ? AND tenant_id = ? AND status = ?", reportID, tenantID, "searching").First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在或已关闭"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	var input struct {
		Location     models.JSON `json:"location" binding:"required"`
		SightingTime string      `json:"sighting_time" binding:"required"`
		Description  string      `json:"description"`
		PhotoURL     string      `json:"photo_url"`
		ReporterName string      `json:"reporter_name"`
		ContactPhone string      `json:"contact_phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	sightingTime, err := time.Parse(time.RFC3339, input.SightingTime)
	if err != nil {
		sightingTime, _ = time.Parse("2006-01-02T15:04:05Z", input.SightingTime)
	}

	sighting := models.SightingReport{
		ReportUUID:   reportID,
		Location:     input.Location,
		SightingTime: sightingTime,
		Description:  input.Description,
		PhotoURL:     input.PhotoURL,
		ReporterName: input.ReporterName,
		ContactPhone: input.ContactPhone,
		ReporterID:   userID,
		TenantID:     tenantID,
	}

	if err := p.DB.Create(&sighting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "添加目击记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "目击记录已提交", "data": sighting})
}

// ListSightings 获取目击记录列表
func (p *PetFinderCtrl) ListSightings(c *gin.Context) {
	reportID := c.Param("id")
	tenantID := getTenantID(c)

	// 验证报告存在
	var report models.PetLostReport
	err := p.DB.Where("report_uuid = ? AND tenant_id = ?", reportID, tenantID).First(&report).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "报告不存在"})
		return
	}

	var sightings []models.SightingReport
	query := p.DB.Where("report_uuid = ?", reportID)

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var total int64
	p.DB.Model(&models.SightingReport{}).Where("report_uuid = ?", reportID).Count(&total)

	if err := query.Order("sighting_time DESC").Offset(offset).Limit(pageSize).Find(&sightings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取目击记录失败"})
		return
	}

	// 脱敏处理
	userID := getUserID(c)
	sanitized := make([]models.SightingReport, len(sightings))
	for i, s := range sightings {
		sanitized[i] = s
		if report.ReporterID != userID {
			sanitized[i].Location = blurPetFinderLocation(s.Location)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      sanitized,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// NearbyReports 附近走失宠物
func (p *PetFinderCtrl) NearbyReports(c *gin.Context) {
	tenantID := getTenantID(c)

	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "10")

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

	var reports []models.PetLostReport
	if err := p.DB.Where("tenant_id = ? AND status = ?", tenantID, "searching").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	userID := getUserID(c)
	var filtered []models.PetLostReport
	for _, r := range reports {
		if r.LastLocation == nil {
			continue
		}
		repLat, ok1 := r.LastLocation["lat"].(float64)
		repLng, ok2 := r.LastLocation["lng"].(float64)
		if !ok1 || !ok2 {
			continue
		}

		dist := haversineDistancePetFinder(lat, lng, repLat, repLng)
		if dist <= radius {
			reportCopy := r
			if reportCopy.ReporterID != userID {
				reportCopy.LastLocation = blurPetFinderLocation(reportCopy.LastLocation)
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

// ManageAlert 订阅/取消警报
func (p *PetFinderCtrl) ManageAlert(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		AlertUUID    string   `json:"alert_uuid"` // 可选，有值表示更新
		Species      string   `json:"species"`
		Latitude     float64  `json:"latitude" binding:"required"`
		Longitude    float64  `json:"longitude" binding:"required"`
		RadiusKM     float64  `json:"radius_km"`
		NotifyEmail  bool     `json:"notify_email"`
		NotifySMS    bool     `json:"notify_sms"`
		NotifyApp    bool     `json:"notify_app"`
		IsActive     *bool    `json:"is_active"` // 用于取消订阅
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 如果传入了 alert_uuid，则更新现有订阅
	if input.AlertUUID != "" {
		var alert models.FinderAlert
		err := p.DB.Where("alert_uuid = ? AND user_id = ? AND tenant_id = ?", input.AlertUUID, userID, tenantID).First(&alert).Error
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
			return
		}

		updates := map[string]interface{}{}
		if input.Species != "" {
			updates["species"] = input.Species
		}
		updates["latitude"] = input.Latitude
		updates["longitude"] = input.Longitude
		if input.RadiusKM > 0 {
			if input.RadiusKM > 50 {
				input.RadiusKM = 50
			}
			updates["radius_km"] = input.RadiusKM
		}
		updates["notify_email"] = input.NotifyEmail
		updates["notify_sms"] = input.NotifySMS
		updates["notify_app"] = input.NotifyApp
		if input.IsActive != nil {
			updates["is_active"] = *input.IsActive
		}

		if err := p.DB.Model(&alert).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新订阅失败"})
			return
		}

		p.DB.First(&alert, alert.ID)
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订阅已更新", "data": alert})
		return
	}

	// 新建订阅
	radius := input.RadiusKM
	if radius <= 0 {
		radius = 10
	}
	if radius > 50 {
		radius = 50
	}

	alert := models.FinderAlert{
		UserID:      userID,
		Species:     input.Species,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		RadiusKM:    radius,
		NotifyEmail: input.NotifyEmail,
		NotifySMS:   input.NotifySMS,
		NotifyApp:   input.NotifyApp,
		IsActive:   true,
		TenantID:    tenantID,
	}

	if err := p.DB.Create(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建订阅失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订阅成功", "data": alert})
}

// ListAlerts 获取用户的警报订阅列表
func (p *PetFinderCtrl) ListAlerts(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var alerts []models.FinderAlert
	query := p.DB.Where("user_id = ? AND tenant_id = ?", userID, tenantID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var total int64
	p.DB.Model(&models.FinderAlert{}).Where("user_id = ? AND tenant_id = ?", userID, tenantID).Count(&total)

	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取订阅列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "success",
		"data": gin.H{
			"list":      alerts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DeleteAlert 删除警报订阅
func (p *PetFinderCtrl) DeleteAlert(c *gin.Context) {
	alertID := c.Param("id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var alert models.FinderAlert
	err := p.DB.Where("alert_uuid = ? AND user_id = ? AND tenant_id = ?", alertID, userID, tenantID).First(&alert).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "订阅不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	if err := p.DB.Delete(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除订阅失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "订阅已删除"})
}

// blurPetFinderLocation 模糊化位置信息
func blurPetFinderLocation(loc models.JSON) models.JSON {
	if loc == nil {
		return nil
	}
	if lat, ok := loc["lat"].(float64); ok {
		loc["lat"] = math.Round(lat*10) / 10
	}
	if lng, ok := loc["lng"].(float64); ok {
		loc["lng"] = math.Round(lng*10) / 10
	}
	return loc
}

// haversineDistancePetFinder 计算两点间的球面距离（km）
func haversineDistancePetFinder(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
