package cache_test

import (
	"fmt"
	"github.com/clambin/mediamon/cache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type stats struct {
	next  int
	fail  bool
	value int
}

func (s stats) update() (interface{}, error) {
	if s.fail {
		return nil, fmt.Errorf("failure")
	}
	return stats{value: s.next}, nil
}

func TestCache_Update(t *testing.T) {
	s := stats{next: 1}
	c := cache.Cache{
		Expiry:    time.Hour,
		LastStats: stats{},
		Updater:   s.update,
	}

	updated := c.Update()
	assert.Equal(t, 1, updated.(stats).value)

	s.next = 2
	updated = c.Update()
	assert.Equal(t, 1, updated.(stats).value)

	s.fail = true
	updated = c.Update()
	assert.Equal(t, 1, updated.(stats).value)
}
