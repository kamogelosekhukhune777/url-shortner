package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kamogelosekhukhune777/url-shortner/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DSNKey is a context key for passing the DSN without exposing it directly. (Optional but good for context usage)
type DSNKey struct{}

// Config holds the configuration for the PostgreSQL database connection.
type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	URLModel        interface{}
}

// Repository implements the database operations.
type Repository struct {
	db *gorm.DB
}

// NewDB initializes and connects to the database, returning a *Repository.
func NewDB(cfg *Config) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	sqlDB, err := dbClient.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database with DSN (%s): %w", dsn, err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// Assume cfg.ConnMaxLifetime is already the correct time.Duration value
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Optional: Auto-migrate the model if the config provides it
	if cfg.URLModel != nil {
		if err := dbClient.AutoMigrate(cfg.URLModel); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate model: %w", err)
		}
	}

	return &Repository{db: dbClient}, nil
}

// Close gracefully closes the database connection.
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// ----- CRUD Operations -----

// Create inserts a new ShortURL record.
func (r *Repository) Create(ctx context.Context, urlModel *service.ShortURL) (*service.ShortURL, error) {
	if err := r.db.WithContext(ctx).Create(urlModel).Error; err != nil {
		return nil, fmt.Errorf("failed to create ShortURL: %w", err)
	}

	return urlModel, nil
}

// Update updates an existing ShortURL record by its short code.
func (r *Repository) Update(ctx context.Context, shortCode string, urlModel *service.ShortURL) error {
	result := r.db.WithContext(ctx).
		Model(urlModel).
		Where("short_code = ?", shortCode).
		Updates(urlModel)

	if result.Error != nil {
		return fmt.Errorf("failed to update ShortURL with code '%s': %w", shortCode, result.Error)
	}
	// Best Practice: Check if any rows were affected for a successful update
	if result.RowsAffected == 0 {
		return errors.New("update failed: short code not found or no changes made")
	}

	return nil
}

// Delete removes a ShortURL record by its short code.
func (r *Repository) Delete(ctx context.Context, shortCode string) error {
	result := r.db.WithContext(ctx).
		Where("short_code = ?", shortCode).
		Delete(&service.ShortURL{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete ShortURL with code '%s': %w", shortCode, result.Error)
	}

	// Check for existence before deletion, or check RowsAffected
	if result.RowsAffected == 0 {
		return errors.New("delete failed: short code not found")
	}

	return nil
}

// Exists checks if a ShortURL record exists by its short code.
func (r *Repository) Exists(ctx context.Context, shortCode string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&service.ShortURL{}).
		Where("short_code = ?", shortCode).
		Count(&count).
		Error

	if err != nil {
		return false, fmt.Errorf("failed to check existence for code '%s': %w", shortCode, err)
	}

	return count > 0, nil
}

// GetByShortCode retrieves a ShortURL record by its short code.
func (r *Repository) GetByShortCode(ctx context.Context, shortCode string) (*service.ShortURL, error) {
	model := &service.ShortURL{}
	err := r.db.WithContext(ctx).
		Where("short_code = ?", shortCode).
		First(model).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ShortURL by code '%s': %w", shortCode, err)
	}

	return model, nil
}

// IncrementAccessCount atomically increments the access_count field.
func (r *Repository) IncrementAccessCount(ctx context.Context, shortCode string) error {
	result := r.db.WithContext(ctx).
		Model(&service.ShortURL{}).
		Where("short_code = ?", shortCode).
		Update("access_count", gorm.Expr("access_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("failed to increment access count for code '%s': %w", shortCode, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("increment failed: short code not found")
	}

	return nil
}
