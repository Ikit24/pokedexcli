package pokecache

type Cache struct {
	data map[string]cacheEntry
	sync.Mutex
}

type cacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

func NewCache() {}

