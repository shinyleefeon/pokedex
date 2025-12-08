package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	data      []byte
	timestamp time.Time
}

type Cache struct {
	val       map[string]CacheEntry
	mu        sync.RWMutex
}



func NewCache(val map[string]CacheEntry, interval time.Duration) *Cache {
	newCache := &Cache{
		val: val,
		mu:  sync.RWMutex{},
	}
	newCache.val["0"] = CacheEntry{[]byte("Initialized"), time.Now()}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.val[key] = CacheEntry{
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
		time.Sleep(time.Second * 10)
	}
}