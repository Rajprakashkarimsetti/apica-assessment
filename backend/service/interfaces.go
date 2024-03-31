package service

type LruCacher interface {
	Get(key string) string
	Set(key, val string)
}
