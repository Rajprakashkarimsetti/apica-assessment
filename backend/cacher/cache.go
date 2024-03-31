package cacher

import (
	"container/list"
	"sync"
	"time"

	"github.com/Rajprakashkarimsetti/apica-project/models"
)

type Cache struct {
	Capacity int
	Cache    map[string]*list.Element
	LruList  *list.List
	Mutex    sync.Mutex
}

// NewCache creates a new cache with the specified capacity.
// It initializes the cache and starts a goroutine to periodically check for expired cache entries.
func NewCache(capacity int) *Cache {
	cache := &Cache{
		Capacity: capacity,
		Cache:    make(map[string]*list.Element),
		LruList:  list.New(),
	}

	go cache.startExpirationCheck()

	return cache
}

// startExpirationCheck periodically checks for expired cache entries and removes them.
func (c *Cache) startExpirationCheck() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.Mutex.Lock()
		for key, elem := range c.Cache {
			entry := elem.Value.(*models.CacheData)
			if time.Since(entry.Timestamp) > 5*time.Second {
				delete(c.Cache, key)
				c.LruList.Remove(elem)
			}
		}

		c.Mutex.Unlock()
	}
}
