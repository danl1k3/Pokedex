package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	entries map[string]cacheEntry
	mu	sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	value []byte
}



func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		entries: map[string]cacheEntry{},
	}
	go reapLoop(interval, newCache)
	return newCache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := c.entries[key]
	entry.value = value
	entry.createdAt = time.Now()
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
	    return nil, false
	}

	return entry.value, true
}

func reapLoop(interval time.Duration, cache *Cache) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		cache.mu.Lock()
		for key, _ := range cache.entries {
			entry := cache.entries[key]
			if time.Now().Sub(entry.createdAt) > interval {
				delete(cache.entries, key)
			}
		}
		cache.mu.Unlock()
	}
}

