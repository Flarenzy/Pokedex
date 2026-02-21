package pokecache

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key already exists")
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Cached   map[string]cacheEntry
	interval time.Duration
	mu       sync.RWMutex
	done     chan bool
}

func NewCache(duration time.Duration) *Cache {
	c := &Cache{
		Cached:   make(map[string]cacheEntry),
		mu:       sync.RWMutex{},
		interval: duration,
		done:     make(chan bool),
	}

	ticker := time.NewTicker(duration)
	c.reapLoop(ticker)
	return c
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.Cached[key]
	if !ok {
		//slog.Error("Key not found: ", key)
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return entry.val, nil
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	slog.Debug("Adding key: ", key, string(val))
	if _, ok := c.Cached[key]; ok {
		//slog.Debug("Key already exists: ", key)
		return fmt.Errorf("%w: %s", ErrKeyExists, key)
	}
	c.Cached[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	//slog.Info("Added key: ", "key", key)
	return nil
}

func (c *Cache) Done() {
	c.done <- true
	close(c.done)
}

func (c *Cache) reapLoop(ticker *time.Ticker) {
	go func() {
		for {
			select {
			case <-ticker.C:
				for k, v := range c.Cached {
					c.mu.Lock()
					if v.createdAt.Add(c.interval).Before(time.Now()) {
						delete(c.Cached, k)
						//slog.Info("Cache reaped: ", k)
					}
					c.mu.Unlock()
				}
			case <-c.done:
				ticker.Stop()
				//slog.Info("Closing cache...")
				return
			}
		}
	}()

}
