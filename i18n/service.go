package i18n

import (
	"errors"

	"mdm-backend/models"

	"gorm.io/gorm"
)

// TranslationService 翻译服务
type TranslationService struct {
	db *gorm.DB
}

// NewTranslationService 创建翻译服务
func NewTranslationService(db *gorm.DB) *TranslationService {
	return &TranslationService{db: db}
}

// ListTranslations 获取翻译列表（支持分页和过滤）
func (s *TranslationService) ListTranslations(filter models.TranslationFilter) ([]models.Translation, int64, error) {
	var translations []models.Translation
	var total int64

	query := s.db.Model(&models.Translation{})

	// 应用过滤条件
	if filter.Locale != "" {
		query = query.Where("locale = ?", filter.Locale)
	}
	if filter.Namespace != "" {
		query = query.Where("namespace = ?", filter.Namespace)
	}
	if filter.Key != "" {
		query = query.Where("key LIKE ?", "%"+filter.Key+"%")
	}
	if filter.Tags != "" {
		query = query.Where("tags LIKE ?", "%"+filter.Tags+"%")
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	if offset < 0 {
		offset = 0
	}

	// 获取数据
	if err := query.Order("locale, namespace, key").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&translations).Error; err != nil {
		return nil, 0, err
	}

	return translations, total, nil
}

// GetTranslationByID 根据ID获取翻译
func (s *TranslationService) GetTranslationByID(id uint) (*models.Translation, error) {
	var translation models.Translation
	err := s.db.First(&translation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &translation, nil
}

// GetTranslationByKey 根据locale和key获取翻译
func (s *TranslationService) GetTranslationByKey(locale, key string) (*models.Translation, error) {
	var translation models.Translation
	err := s.db.Where("locale = ? AND key = ?", locale, key).First(&translation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &translation, nil
}

// CreateTranslation 创建翻译
func (s *TranslationService) CreateTranslation(req *models.TranslationRequest) (*models.Translation, error) {
	// 检查是否已存在
	existing, _ := s.GetTranslationByKey(req.Locale, req.Key)
	if existing != nil {
		return nil, errors.New("translation already exists for this locale and key")
	}

	translation := &models.Translation{
		Locale:    req.Locale,
		Key:       req.Key,
		Namespace: req.Namespace,
		Value:     req.Value,
		Context:   req.Context,
		Tags:      req.Tags,
		IsActive:  true,
	}
	if req.IsActive != nil {
		translation.IsActive = *req.IsActive
	}

	if err := s.db.Create(translation).Error; err != nil {
		return nil, err
	}
	return translation, nil
}

// UpdateTranslation 更新翻译
func (s *TranslationService) UpdateTranslation(id uint, req *models.TranslationRequest) (*models.Translation, error) {
	translation, err := s.GetTranslationByID(id)
	if err != nil {
		return nil, err
	}

	// 如果更改了locale或key，检查是否冲突
	if req.Locale != translation.Locale || req.Key != translation.Key {
		existing, _ := s.GetTranslationByKey(req.Locale, req.Key)
		if existing != nil && existing.ID != id {
			return nil, errors.New("translation already exists for this locale and key")
		}
	}

	// 更新字段
	translation.Locale = req.Locale
	translation.Key = req.Key
	translation.Namespace = req.Namespace
	translation.Value = req.Value
	translation.Context = req.Context
	translation.Tags = req.Tags
	if req.IsActive != nil {
		translation.IsActive = *req.IsActive
	}

	if err := s.db.Save(translation).Error; err != nil {
		return nil, err
	}
	return translation, nil
}

// DeleteTranslation 删除翻译
func (s *TranslationService) DeleteTranslation(id uint) error {
	result := s.db.Delete(&models.Translation{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// BatchCreateTranslations 批量创建翻译
func (s *TranslationService) BatchCreateTranslations(translations []*models.Translation) (int, error) {
	if len(translations) == 0 {
		return 0, nil
	}

	created := 0
	for _, t := range translations {
		if err := s.db.Create(t).Error; err == nil {
			created++
		}
	}
	return created, nil
}

// GetSupportedLocales 获取支持的语言列表
func (s *TranslationService) GetSupportedLocales() ([]string, error) {
	var locales []string
	err := s.db.Model(&models.Translation{}).
		Distinct("locale").
		Pluck("locale", &locales).Error
	return locales, err
}

// GetNamespaces 获取所有命名空间
func (s *TranslationService) GetNamespaces() ([]string, error) {
	var namespaces []string
	err := s.db.Model(&models.Translation{}).
		Distinct("namespace").
		Where("namespace != ''").
		Pluck("namespace", &namespaces).Error
	return namespaces, err
}
