package service

import (
	"github.com/Rajprakashkarimsetti/apica-project/models"
	"github.com/Rajprakashkarimsetti/apica-project/store"
)

type Service struct {
	lruCacherStore store.LruCacher
}

func New(lruCacherStore store.LruCacher) Service {
	return Service{lruCacherStore: lruCacherStore}
}

// Get retrieves the value associated with the given key from the store's LRU cache.
func (s Service) Get(key string) string {
	return s.lruCacherStore.Get(key)
}

// Set stores the key-value pair in the store's LRU cache.
func (s Service) Set(cache *models.CacheData) {
	s.lruCacherStore.Set(cache)
}
