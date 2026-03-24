package models


import (
	"time"

	"gorm.io/gorm"
)

type ShortURL struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OriginalURL  string         `gorm:"type:text;not null" json:"original_url"`
	ShortCode    string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"short_code"`
	ShortURL     string         `gorm:"type:text" json:"short_url"`
	Clicks       int            `gorm:"default:0" json:"clicks"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ShortURL) TableName() string {
	return "short_urls"
}
