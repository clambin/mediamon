package measurer

import (
	"context"
	"sync"
	"time"
)

// Cached measures a value and caches it for Interval seconds.
type Cached[T any] struct {
	lastCheck time.Time
	lastValue T
	Do        func(context.Context) (T, error)
	Interval  time.Duration
	lock      sync.Mutex
}

// Measure returns the cached value if it's within Interval seconds, otherwise it calls Do to measure a new value.
func (c *Cached[T]) Measure(ctx context.Context) (T, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if time.Since(c.lastCheck) > c.Interval {
		value, err := c.Do(ctx)
		if err != nil {
			return value, err
		}
		c.lastValue = value
		c.lastCheck = time.Now()
	}
	return c.lastValue, nil
}
