package timezone

import (
	"fmt"
	"time"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// SupportedTimezones 支持的时区列表
var SupportedTimezones = []map[string]string{
	{"value": "Pacific/Honolulu", "label": "夏威夷 (UTC-10)"},
	{"value": "America/Anchorage", "label": "阿拉斯加 (UTC-9)"},
	{"value": "America/Los_Angeles", "label": "太平洋时间 (UTC-8)"},
	{"value": "America/Denver", "label": "山地时间 (UTC-7)"},
	{"value": "America/Chicago", "label": "中部时间 (UTC-6)"},
	{"value": "America/New_York", "label": "东部时间 (UTC-5)"},
	{"value": "America/Sao_Paulo", "label": "巴西利亚 (UTC-3)"},
	{"value": "Atlantic/Azores", "label": "亚速尔 (UTC-1)"},
	{"value": "Europe/London", "label": "伦敦 (UTC+0)"},
	{"value": "Europe/Paris", "label": "巴黎 (UTC+1)"},
	{"value": "Europe/Berlin", "label": "柏林 (UTC+1)"},
	{"value": "Europe/Moscow", "label": "莫斯科 (UTC+3)"},
	{"value": "Asia/Dubai", "label": "迪拜 (UTC+4)"},
	{"value": "Asia/Karachi", "label": "卡拉奇 (UTC+5)"},
	{"value": "Asia/Kolkata", "label": "印度 (UTC+5:30)"},
	{"value": "Asia/Dhaka", "label": "达卡 (UTC+6)"},
	{"value": "Asia/Bangkok", "label": "曼谷 (UTC+7)"},
	{"value": "Asia/Shanghai", "label": "中国 (UTC+8)"},
	{"value": "Asia/Singapore", "label": "新加坡 (UTC+8)"},
	{"value": "Asia/Tokyo", "label": "东京 (UTC+9)"},
	{"value": "Australia/Sydney", "label": "悉尼 (UTC+10)"},
	{"value": "Pacific/Auckland", "label": "奥克兰 (UTC+12)"},
}

// TimezoneService 时区服务
type TimezoneService struct {
	db *gorm.DB
}

// NewTimezoneService 创建时区服务
func NewTimezoneService(db *gorm.DB) *TimezoneService {
	return &TimezoneService{db: db}
}

// GetUserTimezone 获取用户时区
func (s *TimezoneService) GetUserTimezone(userID uint) (string, error) {
	var config models.TimezoneConfig
	err := s.db.Where("entity_type = ? AND entity_id = ?", "user", userID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 尝试获取租户时区
			return s.GetDefaultTimezone()
		}
		return "", err
	}
	return config.Timezone, nil
}

// GetTenantTimezone 获取租户时区
func (s *TimezoneService) GetTenantTimezone(tenantID uint) (string, error) {
	var config models.TimezoneConfig
	err := s.db.Where("entity_type = ? AND entity_id = ?", "tenant", tenantID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "UTC", nil
		}
		return "", err
	}
	return config.Timezone, nil
}

// GetDefaultTimezone 获取系统默认时区
func (s *TimezoneService) GetDefaultTimezone() (string, error) {
	var config models.TimezoneConfig
	err := s.db.Where("entity_type = ? AND entity_id = ?", "system", 0).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "UTC", nil
		}
		return "", err
	}
	return config.Timezone, nil
}

// SetUserTimezone 设置用户时区
func (s *TimezoneService) SetUserTimezone(userID uint, timezone string) error {
	return s.upsertTimezone("user", userID, timezone)
}

// SetTenantTimezone 设置租户时区
func (s *TimezoneService) SetTenantTimezone(tenantID uint, timezone string) error {
	return s.upsertTimezone("tenant", tenantID, timezone)
}

// SetSystemTimezone 设置系统默认时区
func (s *TimezoneService) SetSystemTimezone(timezone string) error {
	return s.upsertTimezone("system", 0, timezone)
}

// upsertTimezone 插入或更新时区配置
func (s *TimezoneService) upsertTimezone(entityType string, entityID uint, timezone string) error {
	var config models.TimezoneConfig
	err := s.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&config).Error

	if err == gorm.ErrRecordNotFound {
		config = models.TimezoneConfig{
			EntityType: entityType,
			EntityID:   entityID,
			Timezone:   timezone,
			IsActive:   true,
		}
		return s.db.Create(&config).Error
	}

	if err != nil {
		return err
	}

	return s.db.Model(&config).Update("timezone", timezone).Error
}

// GetTimezoneConfig 获取时区配置
func (s *TimezoneService) GetTimezoneConfig(entityType string, entityID uint) (*models.TimezoneConfig, error) {
	var config models.TimezoneConfig
	err := s.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetTimezoneConfigByID 根据ID获取时区配置
func (s *TimezoneService) GetTimezoneConfigByID(id uint) (*models.TimezoneConfig, error) {
	var config models.TimezoneConfig
	err := s.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ListTimezoneConfigs 列出所有时区配置
func (s *TimezoneService) ListTimezoneConfigs(entityType string) ([]models.TimezoneConfig, error) {
	var configs []models.TimezoneConfig
	err := s.db.Where("entity_type = ?", entityType).Find(&configs).Error
	return configs, err
}

// ConvertToTimezone 转换时间到指定时区
func ConvertToTimezone(t time.Time, tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	return t.In(loc)
}

// GetUserTimezone 获取用户时区（独立函数）
func GetUserTimezone(db *gorm.DB, userID uint) string {
	svc := NewTimezoneService(db)
	tz, err := svc.GetUserTimezone(userID)
	if err != nil {
		return "UTC"
	}
	return tz
}

// FormatInTimezone 在指定时区格式化时间
func FormatInTimezone(t time.Time, tz string, format string) string {
	converted := ConvertToTimezone(t, tz)
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return converted.Format(format)
}

// GetTimezoneOffset 获取时区偏移量（小时）
func GetTimezoneOffset(tz string) (float64, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	_, offset := now.In(loc).Zone()
	return float64(offset) / 3600.0, nil
}

// ValidateTimezone 验证时区是否有效
func ValidateTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

// GetSupportedTimezones 获取支持的时区列表
func GetSupportedTimezones() []map[string]string {
	return SupportedTimezones
}

// FormatTimeWithTimezone 格式化时间并附带时区信息
func FormatTimeWithTimezone(t time.Time, tz string) map[string]interface{} {
	converted := ConvertToTimezone(t, tz)
	offset, _ := GetTimezoneOffset(tz)

	return map[string]interface{}{
		"utc":      t.UTC().Format("2006-01-02T15:04:05Z"),
		"local":    converted.Format("2006-01-02T15:04:05"),
		"timezone": tz,
		"offset":   fmt.Sprintf("UTC%+d", int(offset)),
	}
}
