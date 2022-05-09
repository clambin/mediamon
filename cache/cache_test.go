package cache_test

import (
	"fmt"
	"github.com/clambin/mediamon/cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type stats struct {
	next int
	fail bool
}

func (s *stats) update() (int, error) {
	if s.fail {
		return -1, fmt.Errorf("failure")
	}
	return s.next, nil
}

func TestCache_Update(t *testing.T) {
	s := &stats{next: 1}
	c := cache.Cache[int]{
		Duration: time.Hour,
		Updater:  s.update,
	}

	updated := c.Update()
	assert.Equal(t, 1, updated)

	s.next = 2
	updated = c.Update()
	assert.Equal(t, 1, updated)

	s.fail = true
	updated = c.Update()
	assert.Equal(t, 1, updated)
}

func TestCache_Expiry(t *testing.T) {
	s := &stats{next: 1}
	c := cache.Cache[int]{
		Duration: 0,
		Updater:  s.update,
	}

	updated := c.Update()
	assert.Equal(t, 1, updated)

	s.next = 2
	updated = c.Update()
	assert.Equal(t, 2, updated)
}
