package cache

import "time"

// A Cache prevents too many API calls to expensive servers. It stores the latest update, and
// if subsequent calls to Update come before Expiry, it returns the cached status update instead.
type Cache struct {
	// Expiry indicates how long the last  update will be cached
	Expiry time.Duration
	// LastStats holds the latest update (and the initial one). Cache panics if the type of LastStats
	// does not match the return type of the Updater,
	LastStats interface{}
	// Updater provides the next update
	Updater    UpdateFunc
	lastUpdate time.Time
}

// UpdateFunc is the signature of the Updater function
type UpdateFunc func() (interface{}, error)

// Update is called to get the next update. If the latest update is more recent than Expiry, the cached update is
// returned. Otherwise, the Updater function is called.
func (cache *Cache) Update() interface{} {
	if time.Now().After(cache.lastUpdate.Add(cache.Expiry)) {
		stats, err := cache.Updater()

		if err == nil {
			cache.LastStats = stats
			cache.lastUpdate = time.Now()
		}
	}
	return cache.LastStats
}
