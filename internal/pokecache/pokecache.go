package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	data		map[string]cacheEntry
	sync.Mutex
	interval	time.Duration
}

type cacheEntry struct {
	CreatedAt	time.Time
	Val			[]byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		data: make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.Lock()
	defer c.Unlock()
	entry := cacheEntry{
		CreatedAt: time.Now(),
		Val: val,
	}
	c.data[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()
	value, ok := c.data[key]
	if ok {
		return value.Val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.Lock()
		for key, value := range c.data {
			if (time.Now().Sub(value.CreatedAt)) > c.interval {
				delete(c.data, key)
			}
		}
		c.Unlock()
	}
}
