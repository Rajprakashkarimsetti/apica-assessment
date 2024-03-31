package store

import (
	"time"

	"github.com/Rajprakashkarimsetti/apica-project/cacher"
	"github.com/Rajprakashkarimsetti/apica-project/models"
)

type store struct {
	cache *cacher.Cache
}

func New(cache *cacher.Cache) store {
	return store{cache: cache}
}

// Get retrieves the value associated with the given key from the store's cache.
// It moves the accessed element to the front of the LRU list to indicate recent usage.
func (c store) Get(key string) string {
	c.cache.Mutex.Lock()
	defer c.cache.Mutex.Unlock()

	if elem, ok := c.cache.Cache[key]; ok {
		c.cache.LruList.MoveToFront(elem)
		return elem.Value.(*models.CacheData).Value
	}

	return ""
}

// Set stores the key-value pair in the cache, updating the value if the key already exists.
// It also maintains the LRU list by moving accessed elements to the front and removing the least recently used element when the cache exceeds its capacity.
func (c store) Set(key string, value string) {
	c.cache.Mutex.Lock()
	defer c.cache.Mutex.Unlock()

	if elem, ok := c.cache.Cache[key]; ok {
		elem.Value.(*models.CacheData).Value = value
		c.cache.LruList.MoveToFront(elem)
	}

	if len(c.cache.Cache) >= c.cache.Capacity {
		back := c.cache.LruList.Back()
		if back != nil {
			delete(c.cache.Cache, back.Value.(*models.CacheData).Key)
			c.cache.LruList.Remove(back)
		}
	}

	elem := c.cache.LruList.PushFront(&models.CacheData{Key: key, Value: value, Timestamp: time.Now()})
	c.cache.Cache[key] = elem
}
