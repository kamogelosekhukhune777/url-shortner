package shorturl

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Base is the base model for all database entities.
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:TIMESTAMP with time zone;not null"`
	UpdatedAt sql.NullTime   `json:"updated_at" gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ShortURL represents a short URL record in the database.
type ShortURL struct {
	Base
	OriginalURL string `json:"original_url" gorm:"not null"`
	ShortCode   string `json:"short_code" gorm:"unique;not null"`
	AccessCount int    `json:"access_count" gorm:"default:0"`
}
