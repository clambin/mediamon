package cache

import (
	"sync"
	"time"
)

// A Cache prevents too many API calls to expensive servers. It stores the latest update, and
// if subsequent calls to Update come before Duration, it returns the cached status update instead.
//
// Set LastStats to an initial value with the type returned by the Updater. Cache panics if the types do not match.
type Cache struct {
	Duration  time.Duration // how long to cache data
	LastStats interface{}   // holds the latest update (and the initial one).
	Updater   UpdateFunc    // returns the next update
	expiry    time.Time
	lock      sync.Mutex
}

// UpdateFunc is the signature of the updater function
type UpdateFunc func() (interface{}, error)

// Update is called to get the next update. If the latest update is more recent than Duration, the cached update is
// returned. Otherwise, the Updater function is called.
func (cache *Cache) Update() interface{} {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if time.Now().After(cache.expiry) {
		stats, err := cache.Updater()

		if err == nil {
			cache.LastStats = stats
			cache.expiry = time.Now().Add(cache.Duration)
		}

	}

	return cache.LastStats
}
