package service

import "github.com/Rajprakashkarimsetti/apica-project/models"

type LruCacher interface {
	Get(key string) string
	Set(cache *models.CacheData)
}
