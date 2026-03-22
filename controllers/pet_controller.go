package controllers

import (
	"net/http"
	"time"

	"mdm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PetCtrl 宠物档案控制器
type PetCtrl struct {
	DB *gorm.DB
}

// RegisterPetRoutes 注册宠物档案路由
func (p *PetCtrl) RegisterPetRoutes(r *gin.RouterGroup) {
	r.GET("/pets", p.ListPets)
	r.POST("/pets", p.CreatePet)
	r.GET("/pets/:pet_id", p.GetPet)
	r.PUT("/pets/:pet_id", p.UpdatePet)
	r.DELETE("/pets/:pet_id", p.DeletePet)
	// 宠物设备绑定
	r.GET("/pets/:pet_id/devices", p.GetPetDevices)
	r.POST("/pets/:pet_id/devices", p.BindDevice)
	r.DELETE("/pets/:pet_id/devices/:device_id", p.UnbindDevice)
}

// getUserID 从上下文获取当前用户ID
func getUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		return id.(uint)
	}
	return 0
}

// getTenantID 从上下文获取租户ID
func getTenantID(c *gin.Context) string {
	if tid, exists := c.Get("tenant_id"); exists {
		return tid.(string)
	}
	return ""
}

// ListPets 获取当前用户的宠物列表
func (p *PetCtrl) ListPets(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var pets []models.Pet
	query := p.DB.Where("owner_id = ? AND tenant_id = ?", userID, tenantID).Where("status != ?", "deceased")

	// 支持按 status 过滤
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	// 支持按 species 过滤
	if species := c.Query("species"); species != "" {
		query = query.Where("species = ?", species)
	}

	if err := query.Order("created_at DESC").Find(&pets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取宠物列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": pets})
}

// CreatePet 创建宠物档案
func (p *PetCtrl) CreatePet(c *gin.Context) {
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		PetName   string  `json:"pet_name" binding:"required"`
		Species   string  `json:"species" binding:"required"`
		Breed     string  `json:"breed"`
		Gender    string  `json:"gender"`
		BirthDate string  `json:"birth_date"` // YYYY-MM-DD
		Weight    float64 `json:"weight"`
		Color     string  `json:"color"`
		Microchip string  `json:"microchip"`
		AvatarURL string  `json:"avatar_url"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	pet := models.Pet{
		PetName:      input.PetName,
		Species:      input.Species,
		Breed:        input.Breed,
		Gender:       input.Gender,
		Weight:       input.Weight,
		Color:        input.Color,
		Microchip:    input.Microchip,
		AvatarURL:    input.AvatarURL,
		OwnerID:      userID,
		Status:       "active",
		Description:  input.Description,
		TenantID:     tenantID,
	}

	if input.BirthDate != "" {
		t, err := time.Parse("2006-01-02", input.BirthDate)
		if err == nil {
			pet.BirthDate = &t
		}
	}

	if err := p.DB.Create(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "创建宠物档案失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": pet})
}

// GetPet 获取宠物详情
func (p *PetCtrl) GetPet(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var pet models.Pet
	err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取宠物详情失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": pet})
}

// UpdatePet 更新宠物档案
func (p *PetCtrl) UpdatePet(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var pet models.Pet
	err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取宠物失败"})
		return
	}

	var input struct {
		PetName      string  `json:"pet_name"`
		Species      string  `json:"species"`
		Breed        string  `json:"breed"`
		Gender       string  `json:"gender"`
		BirthDate    string  `json:"birth_date"`
		Weight       float64 `json:"weight"`
		Color        string  `json:"color"`
		Microchip    string  `json:"microchip"`
		AvatarURL    string  `json:"avatar_url"`
		Status       string  `json:"status"`
		Description  string  `json:"description"`
		VaccinationRecords []string `json:"vaccination_records"`
		HealthRecords      []string `json:"health_records"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if input.PetName != "" { updates["pet_name"] = input.PetName }
	if input.Species != "" { updates["species"] = input.Species }
	if input.Breed != "" { updates["breed"] = input.Breed }
	if input.Gender != "" { updates["gender"] = input.Gender }
	if input.Weight > 0 { updates["weight"] = input.Weight }
	if input.Color != "" { updates["color"] = input.Color }
	if input.Microchip != "" { updates["microchip"] = input.Microchip }
	if input.AvatarURL != "" { updates["avatar_url"] = input.AvatarURL }
	if input.Status != "" { updates["status"] = input.Status }
	if input.Description != "" { updates["description"] = input.Description }
	if len(input.VaccinationRecords) > 0 { updates["vaccination_records"] = input.VaccinationRecords }
	if len(input.HealthRecords) > 0 { updates["health_records"] = input.HealthRecords }

	if input.BirthDate != "" {
		if t, err := time.Parse("2006-01-02", input.BirthDate); err == nil {
			updates["birth_date"] = t
		}
	}

	if err := p.DB.Model(&pet).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "更新失败"})
		return
	}

	p.DB.First(&pet, pet.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": pet})
}

// DeletePet 删除宠物档案（软删除）
func (p *PetCtrl) DeletePet(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var pet models.Pet
	err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取宠物失败"})
		return
	}

	// 软删除
	if err := p.DB.Delete(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// GetPetDevices 获取宠物绑定的设备列表
func (p *PetCtrl) GetPetDevices(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	// 验证宠物归属
	var pet models.Pet
	if err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}

	var bindings []models.PetDeviceBinding
	if err := p.DB.Where("pet_uuid = ? AND is_active = ?", petID, true).Find(&bindings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "获取设备列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": bindings})
}

// BindDevice 绑定设备到宠物
func (p *PetCtrl) BindDevice(c *gin.Context) {
	petID := c.Param("pet_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	var input struct {
		DeviceID    string `json:"device_id" binding:"required"`
		BindingType string `json:"binding_type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "device_id 不能为空"})
		return
	}

	// 验证宠物归属
	var pet models.Pet
	if err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}

	// 检查是否已存在绑定
	var existing models.PetDeviceBinding
	err := p.DB.Where("pet_uuid = ? AND device_id = ? AND is_active = ?", petID, input.DeviceID, true).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "设备已绑定到此宠物"})
		return
	}

	bindingType := input.BindingType
	if bindingType == "" {
		bindingType = "primary"
	}

	binding := models.PetDeviceBinding{
		PetUUID:     petID,
		DeviceID:    input.DeviceID,
		BindingType: bindingType,
		IsActive:    true,
		BoundAt:     time.Now(),
		TenantID:    tenantID,
	}

	if err := p.DB.Create(&binding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "绑定设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "绑定成功", "data": binding})
}

// UnbindDevice 解绑设备
func (p *PetCtrl) UnbindDevice(c *gin.Context) {
	petID := c.Param("pet_id")
	deviceID := c.Param("device_id")
	userID := getUserID(c)
	tenantID := getTenantID(c)

	// 验证宠物归属
	var pet models.Pet
	if err := p.DB.Where("pet_uuid = ? AND owner_id = ? AND tenant_id = ?", petID, userID, tenantID).First(&pet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "宠物不存在"})
		return
	}

	var binding models.PetDeviceBinding
	err := p.DB.Where("pet_uuid = ? AND device_id = ? AND is_active = ?", petID, deviceID, true).First(&binding).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "绑定记录不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "查询失败"})
		return
	}

	now := time.Now()
	binding.IsActive = false
	binding.UnboundAt = &now

	if err := p.DB.Save(&binding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5001, "message": "解绑失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "解绑成功"})
}
