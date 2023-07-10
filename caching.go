package caching

import (
	"errors"
	"sync"
	"time"
)

type valueGetter[T any] func() T

type expiringValue[T any] struct {
	value    T
	expireAt time.Time
}

// Cache is a simple in-memory cache with a TTL.
type Cache[T any] struct {
	ttl    time.Duration
	memory sync.Map
}

// New creates a new Cache with a given TTL.
func New[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		ttl:    ttl,
		memory: sync.Map{},
	}
}

func (c *Cache[T]) getNewExpiringValue(key string, getValue valueGetter[T]) expiringValue[T] {
	v := expiringValue[T]{
		value:    getValue(),
		expireAt: time.Now().Add(c.ttl),
	}

	return v
}

// Remember returns the value for a given key, or computes it if it is not in the cache (or expired).
func (c *Cache[T]) Remember(key string, getValue valueGetter[T]) T {
	v, found := c.memory.Load(key)
	if !found {
		newExpiringValue := c.getNewExpiringValue(key, getValue)
		c.memory.Store(key, newExpiringValue)

		return newExpiringValue.value
	}

	expiringValueInCache := v.(expiringValue[T])
	if expiringValueInCache.expireAt.Before(time.Now()) {
		newExpiringValue := c.getNewExpiringValue(key, getValue)
		c.memory.Store(key, newExpiringValue)

		return newExpiringValue.value
	}

	return expiringValueInCache.value
}

// GetStaleThenRecompute returns the value for a given key, and recomputes it in background if it is expired.
func (c *Cache[T]) GetStaleThenRecompute(key string, getValue valueGetter[T]) T {
	v, found := c.memory.Load(key)
	if !found {
		newExpiringValue := c.getNewExpiringValue(key, getValue)
		c.memory.Store(key, newExpiringValue)

		return newExpiringValue.value
	}

	expiringValueInCache := v.(expiringValue[T])
	if expiringValueInCache.expireAt.Before(time.Now()) {
		go func() {
			newExpiringValue := c.getNewExpiringValue(key, getValue)
			c.memory.Store(key, newExpiringValue)
		}()
	}

	return expiringValueInCache.value
}

// Get returns the value for a given key with boolean indicating if it is expired, or an error if it is not in the cache.
func (c *Cache[T]) Get(key string) (value T, expired bool, err error) {
	var defaultValue T
	v, found := c.memory.Load(key)
	if !found {
		return defaultValue, false, errors.New("value not found")
	}

	expiringValueInCache := v.(expiringValue[T])
	if v.(expiringValue[T]).expireAt.Before(time.Now()) {
		return expiringValueInCache.value, true, nil
	}

	return expiringValueInCache.value, false, nil
}
