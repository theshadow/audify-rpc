package api

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cacher interface {
	Get(key string) (interface{}, bool, error)
	Set(key string, x interface{}, d time.Duration) error
}

type GoCache struct {
	cache.Cache
}

func (c *GoCache) Get(key string) (interface{}, bool, error) {
	data, exists := c.Cache.Get(key)
	if !exists {
		return nil, false, nil
	}
	return data, true, nil
}

func (c *GoCache) Set(key string, data interface{}, d time.Duration) error {
	c.Cache.Set(key, data, d)
	return nil
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *GoCache {
	return &GoCache{
		Cache: *cache.New(defaultExpiration, cleanupInterval),
	}
}