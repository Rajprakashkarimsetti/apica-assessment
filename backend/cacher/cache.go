package cacher

import (
	"sync"
	"time"

	"github.com/Rajprakashkarimsetti/apica-project/models"
)

type Cache struct {
	Capacity int
	Cache    map[string]*models.CacheData
	Head     *models.CacheData
	Tail     *models.CacheData
	Mutex    sync.Mutex
}

// NewCache creates a new cache with the specified capacity.
// It initializes the cache and starts a goroutine to periodically check for expired cache entries.
func NewCache(capacity int) *Cache {
	cache := &Cache{
		Capacity: capacity,
		Cache:    make(map[string]*models.CacheData),
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

		curr := c.Head
		for curr != nil {
			next := curr.Next
			expirationTime := time.Duration(curr.Expiration) * time.Second
			if time.Since(curr.TimeStamp) > expirationTime {
				delete(c.Cache, curr.Key)
				c.RemoveTail()
			}

			curr = next
		}

		c.Mutex.Unlock()
	}
}

func (c *Cache) MoveToFront(entry *models.CacheData) {
	if entry == c.Head {
		return
	}

	if entry == c.Tail {
		c.Tail = entry.Prev
	}

	if entry.Prev != nil {
		entry.Prev.Next = entry.Next
	}

	if entry.Next != nil {
		entry.Next.Prev = entry.Prev
	}

	entry.Prev = nil
	entry.Next = c.Head
	c.Head.Prev = entry
	c.Head = entry
}

func (c *Cache) RemoveTail() {
	if c.Tail == nil {
		return
	}

	prev := c.Tail.Prev
	if prev != nil {
		prev.Next = nil
	} else {
		c.Head = nil
	}
	c.Tail = prev
}
