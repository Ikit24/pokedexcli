package pokecache

import (
	"time"
)

type Cache struct {
	data map[string]cacheEntry
	sync.Mutex
}

type cacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		data: make(map[string]cacheEntry),
	}
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
