package cache

import (
	"sync"
	"time"
)

// A Cache prevents too many API calls to expensive servers. It stores the latest update, and
// if subsequent calls to Update come before Duration, it returns the cached status update instead.
type Cache[T any] struct {
	Duration  time.Duration // how long to cache data
	Updater   UpdateFunc[T] // returns the next update
	lastStats T             // holds the latest update (and the initial one).
	expiry    time.Time
	lock      sync.Mutex
}

// UpdateFunc is the signature of the updater function
type UpdateFunc[T any] func() (T, error)

// Update is called to get the next update. If the latest update is more recent than Duration, the cached update is
// returned. Otherwise, the Updater function is called.
func (cache *Cache[T]) Update() T {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if time.Now().After(cache.expiry) {
		stats, err := cache.Updater()

		if err == nil {
			cache.lastStats = stats
			cache.expiry = time.Now().Add(cache.Duration)
		}

	}

	return cache.lastStats
}
