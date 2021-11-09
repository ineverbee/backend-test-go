package rediscache

import (
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

// Cache ...
type Cache struct {
	password string
	pool     *redigo.Pool
}

// New redis cache
func New(pool *redigo.Pool, password string) *Cache {
	return &Cache{
		password: password,
		pool:     pool,
	}
}

func (s *Cache) getClient() (redigo.Conn, error) {
	c := s.pool.Get()
	if _, err := c.Do("AUTH", s.password); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

// Get ...
func (s *Cache) Get(key interface{}) (int, error) {
	c, err := s.getClient()
	if err != nil {
		return 0, err
	}
	defer c.Close()
	data, err := redigo.Int(c.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return data, nil
}

// Set ...
func (s *Cache) Set(key, val interface{}, ttl time.Duration) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}
	defer c.Close()
	_, err = c.Do("SET", key, val)
	if err != nil {
		return fmt.Errorf("error setting key %v to %v: %v", key, val, err)
	}
	return nil
}
