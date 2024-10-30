package keeper

import (
	"sync"
	"time"
)

type keeperMCacheItem struct {
	value      interface{}
	expiration int64
}

type keeperMCache struct {
	data   map[string]keeperMCacheItem
	mu     sync.RWMutex
	ttl    time.Duration
	ticker *time.Ticker
	stop   chan bool
}

func NewKeeperMCache(ttl time.Duration) *keeperMCache {
	cache := &keeperMCache{
		data:   make(map[string]keeperMCacheItem),
		ttl:    ttl,
		ticker: time.NewTicker(ttl),
		stop:   make(chan bool),
	}
	go cache.cleanup()
	return cache
}

func (c *keeperMCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = keeperMCacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl).UnixNano(),
	}
}

func (c *keeperMCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.data[key]
	if !exists || time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}

func (c *keeperMCache) cleanup() {
	for {
		select {
		case <-c.ticker.C:
			c.mu.Lock()
			for key, item := range c.data {
				if time.Now().UnixNano() > item.expiration {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			c.ticker.Stop()
			return
		}
	}
}

func (c *keeperMCache) StopCleanup() {
	c.stop <- true
}
