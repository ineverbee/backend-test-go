package cache

import "time"

// Cache interface
type Cache interface {
	Get(key interface{}) (int, error)
	Set(key, val interface{}, ttl time.Duration) error
}
