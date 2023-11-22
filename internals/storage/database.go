package storage

import (
	"errors"
	"sync"

	"github.com/kamogelosekhukhune777/url-shortner/pkg/models"
)

var ErrMappingNotFound = errors.New("mapping not found")

type InMemoryDB struct {
	sync.RWMutex
	data map[string]*models.URLMapping
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]*models.URLMapping),
	}
}

func (db *InMemoryDB) Save(mapping *models.URLMapping) error {
	db.Lock()
	defer db.Unlock()
	db.data[mapping.ShortURL] = mapping

	return nil
}

func (db *InMemoryDB) Get(shortURL string) (*models.URLMapping, error) {
	db.RLock()
	defer db.RUnlock()
	if mapping, ok := db.data[shortURL]; ok {
		return mapping, nil
	}

	return nil, ErrMappingNotFound
}

func (db *InMemoryDB) Delete(shortURL string) error {

	db.Lock()
	defer db.Unlock()
	if _, exists := db.data[shortURL]; !exists {
		return ErrMappingNotFound
	}

	return nil
}

func (db *InMemoryDB) Update(mapping *models.URLMapping) error {
	db.Lock()
	defer db.Unlock()
	if _, exists := db.data[mapping.ShortURL]; !exists {
		return ErrMappingNotFound
	}
	db.data[mapping.ShortURL] = mapping

	return nil
}
