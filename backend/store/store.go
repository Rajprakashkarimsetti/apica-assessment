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
		c.cache.MoveToFront(elem)
		return elem.Value
	}

	return ""
}

// Set stores the key-value pair in the cache, updating the value if the key already exists.
// It also maintains the LRU list by moving accessed elements to the front and removing the least recently used element when the cache exceeds its capacity.
func (c store) Set(cache *models.CacheData) {
	c.cache.Mutex.Lock()
	defer c.cache.Mutex.Unlock()

	if elem, ok := c.cache.Cache[cache.Key]; ok {
		elem.Value = cache.Value
		c.cache.MoveToFront(elem)
	}

	if len(c.cache.Cache) >= c.cache.Capacity {
		delete(c.cache.Cache, c.cache.Tail.Key)
		c.cache.RemoveTail()

	}

	newEntry := &models.CacheData{
		Key:        cache.Key,
		Value:      cache.Value,
		Expiration: cache.Expiration,
		TimeStamp:  time.Now(),
		Next:       c.cache.Head,
	}

	if c.cache.Head != nil {
		c.cache.Head.Prev = newEntry
	}

	c.cache.Head = newEntry
	if c.cache.Tail == nil {
		c.cache.Tail = newEntry
	}

	c.cache.Cache[cache.Key] = newEntry
}
