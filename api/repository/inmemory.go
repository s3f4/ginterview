package repository

import (
	"sync"
)

// InMemoryRepository interface
type InMemoryRepository interface {
	Create(string, interface{})
	Get(string) interface{}
	Exist(string) bool
}

type inMemoryRepository struct {
	// lock is used to thread safety of items map
	lock sync.RWMutex
	// items holds key value storage in-memory
	items map[string]interface{}
}

// NewInMemoryRepository creates an items map and returns inMemoryRepository pointer
func NewInMemoryRepository() InMemoryRepository {
	return &inMemoryRepository{
		items: make(map[string]interface{}),
	}
}

// Create gets key and value and stores
func (p *inMemoryRepository) Create(key string, value interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.items[key] = value
}

// Get returns value of key
func (p *inMemoryRepository) Get(key string) interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.items[key]
}

// Exists checks if key is exists in items map
func (p *inMemoryRepository) Exist(key string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	_, ok := p.items[key]
	return ok
}
