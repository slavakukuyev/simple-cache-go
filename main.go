package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	ttl    time.Duration
	data   map[string]interface{}
	mutex  sync.RWMutex
	stopCh chan struct{} // Channel to signal goroutine to stop
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl:    ttl,
		data:   make(map[string]interface{}),
		stopCh: make(chan struct{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	val, found := c.data[key]
	if !found {
		return nil, false
	}

	return val, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value

	// Set a timer to delete the item after the TTL duration
	go func() {
		select {
		case <-time.After(c.ttl):
			c.mutex.Lock()
			delete(c.data, key)
			c.mutex.Unlock()
		case <-c.stopCh:
			return // Stop the goroutine
		}
	}()
}

// Stop the goroutine that expires cache entries
func (c *Cache) StopExpiration() {
	c.stopCh <- struct{}{}
}

func main() {
	cache := NewCache(5 * time.Second) // Set cache with a TTL of 5 seconds

	// Test the cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if val, found := cache.Get("key1"); found {
		fmt.Println("Value for key1:", val)
	} else {
		fmt.Println("Key1 not found in cache.")
	}

	cache.StopExpiration() // Stop the goroutine that expires cache entries

	time.Sleep(6 * time.Second) // Wait for cache entries to expire

	if val, found := cache.Get("key1"); found {
		fmt.Printf("Key1 still found in cache after TTL, because goroutine stopped: %v\n", val)
	} else {
		fmt.Println("Key1 expired from cache.")
	}

	if val, found := cache.Get("key2"); found {
		fmt.Printf("Key2 still found in cache after TTL: %v\n", val)
	} else {
		fmt.Println("Key2 expired from cache.")
	}
}
