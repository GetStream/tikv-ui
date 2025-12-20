package utils

import "sync"

type Cache struct {
	mu   sync.RWMutex
	data map[string]map[string]any
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]map[string]any),
	}
}

func (c *Cache) Set(category, key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[category]; !ok {
		c.data[category] = make(map[string]any)
	}
	c.data[category][key] = value
}

func (c *Cache) Get(category, key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m, ok := c.data[category]
	if !ok {
		return nil, false
	}
	v, ok := m[key]
	return v, ok
}
