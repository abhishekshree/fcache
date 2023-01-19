package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

type Cacher interface {
	Set([]byte, []byte, time.Duration) error
	Has([]byte) bool
	Get([]byte) ([]byte, error)
	Del([]byte) error
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, value []byte, t time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[string(key)] = value

	// if t > 0, set a timer to delete the key
	if t > 0 {
		go func() {
			<-time.After(t)
			c.Del(key)
		}()
	}

	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.data[string(key)]
	return ok
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if value, ok := c.data[string(key)]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("key not found")
}

func (c *Cache) Del(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}
