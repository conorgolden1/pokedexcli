package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	v  map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	ticker := time.NewTicker(interval)
	cache := Cache{

		v: map[string]CacheEntry{},
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.reapLoop(interval)
			}

		}
	}()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheEntry := CacheEntry{time.Now(), val}
	c.v[key] = cacheEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.v[key]; ok {
		return v.val, ok
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range c.v {
		if time.Since(value.createdAt) > interval {
			delete(c.v, key)
		}
	}
}
