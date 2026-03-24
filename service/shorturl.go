package service

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"shorturl/models"

	"gorm.io/gorm"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	baseURL = "http://localhost:8080/"
)

var shortURLService *ShortURLService

type ShortURLService struct{}

func NewShortURLService() *ShortURLService {
	if shortURLService == nil {
		shortURLService = &ShortURLService{}
	}
	return shortURLService
}

// generateShortCode 生成随机短码
func (s *ShortURLService) generateShortCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// generateUniqueShortCode 生成唯一的短码
func (s *ShortURLService) generateUniqueShortCode() (string, error) {
	for i := 0; i < 10; i++ {
		code := s.generateShortCode(6)
		// 检查是否已存在
		var count int64
		if err := models.DB.Model(&models.ShortURL{}).Where("short_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
	return "", errors.New("无法生成唯一的短码")
}

// CreateShortURL 创建短链
func (s *ShortURLService) CreateShortURL(originalURL string, customCode string, expiresAt *time.Time) (*models.ShortURL, error) {
	// 验证 URL 格式
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		return nil, errors.New("URL 必须以 http:// 或 https:// 开头")
	}

	var shortCode string
	var err error

	if customCode != "" {
		// 检查自定义短码是否已存在
		var count int64
		if err := models.DB.Model(&models.ShortURL{}).Where("short_code = ?", customCode).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("自定义短码已存在")
		}
		shortCode = customCode
	} else {
		// 生成随机短码
		shortCode, err = s.generateUniqueShortCode()
		if err != nil {
			return nil, err
		}
	}

	shortURL := baseURL + shortCode

	url := &models.ShortURL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		ShortURL:    shortURL,
		Clicks:      0,
		IsActive:    true,
		ExpiresAt:   expiresAt,
	}

	if err := models.DB.Create(url).Error; err != nil {
		return nil, err
	}

	return url, nil
}

// UpdateShortURL 更新短链
func (s *ShortURLService) UpdateShortURL(id uint, originalURL string, customCode string, expiresAt *time.Time, isActive *bool) (*models.ShortURL, error) {
	var url models.ShortURL
	if err := models.DB.First(&url, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("短链不存在")
		}
		return nil, err
	}

	// 更新字段
	if originalURL != "" {
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			return nil, errors.New("URL 必须以 http:// 或 https:// 开头")
		}
		url.OriginalURL = originalURL
	}

	if customCode != "" && customCode != url.ShortCode {
		// 检查自定义短码是否已存在
		var count int64
		if err := models.DB.Model(&models.ShortURL{}).Where("short_code = ? AND id != ?", customCode, id).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("自定义短码已存在")
		}
		url.ShortCode = customCode
		url.ShortURL = baseURL + customCode
	}

	if expiresAt != nil {
		url.ExpiresAt = expiresAt
	}

	if isActive != nil {
		url.IsActive = *isActive
	}

	if err := models.DB.Save(&url).Error; err != nil {
		return nil, err
	}

	return &url, nil
}

// DeleteShortURL 删除短链
func (s *ShortURLService) DeleteShortURL(id uint) error {
	var url models.ShortURL
	if err := models.DB.First(&url, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("短链不存在")
		}
		return err
	}

	if err := models.DB.Delete(&models.ShortURL{}, id).Error; err != nil {
		return err
	}

	return nil
}

// GetShortURL 获取短链详情
func (s *ShortURLService) GetShortURL(id uint) (*models.ShortURL, error) {
	var url models.ShortURL
	if err := models.DB.First(&url, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("短链不存在")
		}
		return nil, err
	}

	return &url, nil
}

// ListShortURLs 获取短链列表
func (s *ShortURLService) ListShortURLs(page, pageSize int) ([]models.ShortURL, int64, error) {
	var urls []models.ShortURL
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	if err := models.DB.Model(&models.ShortURL{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := models.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&urls).Error; err != nil {
		return nil, 0, err
	}

	return urls, total, nil
}

// GetByShortCode 根据短码获取短链
func (s *ShortURLService) GetByShortCode(shortCode string) (*models.ShortURL, error) {
	var url models.ShortURL
	if err := models.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("短链不存在")
		}
		return nil, err
	}

	// 检查是否激活
	if !url.IsActive {
		return nil, errors.New("短链已禁用")
	}

	// 检查是否过期
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return nil, errors.New("短链已过期")
	}

	// 增加点击次数
	models.DB.Model(&url).UpdateColumn("clicks", url.Clicks+1)

	return &url, nil
}

// StatsShortURL 获取短链统计信息
func (s *ShortURLService) StatsShortURL(id uint) (map[string]interface{}, error) {
	var url models.ShortURL
	if err := models.DB.First(&url, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("短链不存在")
		}
		return nil, err
	}

	stats := map[string]interface{}{
		"id":           url.ID,
		"short_code":   url.ShortCode,
		"short_url":    url.ShortURL,
		"original_url": url.OriginalURL,
		"clicks":       url.Clicks,
		"is_active":    url.IsActive,
		"created_at":   url.CreatedAt,
		"expires_at":   url.ExpiresAt,
	}

	return stats, nil
}
