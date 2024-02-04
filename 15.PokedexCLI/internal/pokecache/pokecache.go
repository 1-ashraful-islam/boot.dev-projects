package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mutex *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cache: make(map[string]cacheEntry),
		mutex: &sync.RWMutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) reap(interval time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, entry := range c.cache {
		if time.Since(entry.createdAt) > interval*time.Second {
			delete(c.cache, key)
		}
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.reap(interval)
	}
}
