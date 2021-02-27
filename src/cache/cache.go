package cache

// Cacher is an interface that represents an object that caches objects.
type Cacher interface {
	Get(key string) (string, error)
	Set(key, value string, expireInSec int) error
}
