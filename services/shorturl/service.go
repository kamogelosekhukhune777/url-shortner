package shorturl

import (
	"context"
	"errors"
	"fmt"
	"log"
)

// Storer defines the contract for database access (Repository/DAO layer).
type Storer interface {
	Create(ctx context.Context, urlModel *ShortURL) (*ShortURL, error)
	Update(ctx context.Context, shortCode string, urlModel *ShortURL) error
	Delete(ctx context.Context, shortCode string) error
	Exists(ctx context.Context, shortCode string) (bool, error)
	GetByShortCode(ctx context.Context, shortCode string) (*ShortURL, error)
	IncrementAccessCount(ctx context.Context, shortCode string) error
}

// ShortURLService implements the business logic for URL shortening operations.
type ShortURLService struct {
	Logger *log.Logger
	Store  Storer
}

func NewService(store Storer, logger *log.Logger) *ShortURLService {
	return &ShortURLService{
		Store:  store,
		Logger: logger,
	}
}

// --- Service Methods (Business Logic Layer) ---

// Create delegates the creation task to the Storer.
func (s *ShortURLService) Create(ctx context.Context, urlModel *ShortURL) (*ShortURL, error) {
	model, err := s.Store.Create(ctx, urlModel)
	if err != nil {
		s.Logger.Printf("ERROR: failed to create ShortURL: %v", err)
		return nil, fmt.Errorf("service failed to create URL: %w", err)
	}

	return model, nil
}

// Update delegates the update task to the Storer.
func (s *ShortURLService) Update(ctx context.Context, shortCode string, urlModel *ShortURL) error {
	if err := s.Store.Update(ctx, shortCode, urlModel); err != nil {
		s.Logger.Printf("ERROR: failed to update ShortURL with code %s: %v", shortCode, err)
		return fmt.Errorf("service failed to update URL with code '%s': %w", shortCode, err)
	}

	return nil
}

// Delete delegates the deletion task to the Storer.
func (s *ShortURLService) Delete(ctx context.Context, shortCode string) error {
	if err := s.Store.Delete(ctx, shortCode); err != nil {
		s.Logger.Printf("ERROR: failed to delete ShortURL with code %s: %v", shortCode, err)
		return fmt.Errorf("service failed to delete URL with code '%s': %w", shortCode, err)
	}

	return nil
}

// Exists checks if a URL exists via the Storer.
func (s *ShortURLService) Exists(ctx context.Context, shortCode string) (bool, error) {
	exists, err := s.Store.Exists(ctx, shortCode)
	if err != nil {
		s.Logger.Printf("ERROR: failed to check existence for code %s: %v", shortCode, err)
		return false, fmt.Errorf("service failed to check URL existence for code '%s': %w", shortCode, err)
	}

	return exists, nil
}

// GetByShortCode retrieves a URL record via the Storer.
func (s *ShortURLService) GetByShortCode(ctx context.Context, shortCode string) (*ShortURL, error) {
	url, err := s.Store.GetByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, fmt.Errorf("URL with code '%s' not found: %w", shortCode, err)
		}

		s.Logger.Printf("ERROR: failed to get ShortURL by code %s: %v", shortCode, err)
		return nil, fmt.Errorf("service failed to retrieve URL: %w", err)
	}
	return url, nil
}

// IncrementAccessCount delegates the count increment task to the Storer.
func (s *ShortURLService) IncrementAccessCount(ctx context.Context, shortCode string) error {
	if err := s.Store.IncrementAccessCount(ctx, shortCode); err != nil {
		s.Logger.Printf("ERROR: failed to increment access count for code %s: %v", shortCode, err)
		return fmt.Errorf("service failed to increment access count for code '%s': %w", shortCode, err)
	}
	return nil
}
