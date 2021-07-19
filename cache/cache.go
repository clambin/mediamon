package cache

import "time"

type Cache struct {
	interval   time.Duration
	lastUpdate time.Time
	lastStats  interface{}
	updater    UpdateFunc
}

type UpdateFunc func() (interface{}, error)

func New(interval time.Duration, stats interface{}, updater UpdateFunc) *Cache {
	return &Cache{
		interval:  interval,
		lastStats: stats,
		updater:   updater,
	}
}

func (cache *Cache) Update() interface{} {
	if time.Now().After(cache.lastUpdate.Add(cache.interval)) {
		stats, err := cache.updater()

		if err == nil {
			cache.lastStats = stats
			cache.lastUpdate = time.Now()
		}
	}
	return cache.lastStats
}
