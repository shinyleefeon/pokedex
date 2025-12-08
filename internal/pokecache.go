package pokecache

import "time"


type cacheEntry struct {
	data      []byte
	timestamp time.Time
}

type Cache struct {
	val       map[string]cacheEntry
	mu        sync.RWMutex
}



func NewCache(val map[string]cacheEntry, interval time.Duration) *Cache {
	newCache := &Cache{
		val: val,
		mu:  sync.RWMutex{},
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val[key] = cacheEntry{
		data:      val,
		timestamp: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.val[key]
	if !exists {
		return nil, false
	}
	return entry.data, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	for len(c.val) > 0 {
		for key, entry := range c.val {
			if time.Since(entry.timestamp) > interval {
				delete(c.val, key)
			}
		}
	}
}