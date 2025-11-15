package measurer_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/clambin/mediamon/v2/internal/measurer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCached_Measure(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var calls atomic.Int64
		c := measurer.Cached[int64]{
			Interval: time.Second,
			Do: func(ctx context.Context) (int64, error) {
				return calls.Add(1), nil
			},
		}
		ctx := t.Context()

		// measure a value. should be called.
		value, err := c.Measure(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(1), value)

		// measure again. should be cached.
		value, err = c.Measure(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(1), value)

		// wait for the cache to expire.
		time.Sleep(time.Second * 5)

		// measure again. should be called again.
		value, err = c.Measure(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(2), value)

		// measure again. should be cached again.
		value, err = c.Measure(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(2), value)

		// wait for the cache to expire.
		time.Sleep(time.Second * 5)

		// measure again. errors should be returned.
		c.Do = func(_ context.Context) (int64, error) { return 0, errors.New("failed") }
		_, err = c.Measure(ctx)
		require.Error(t, err)
	})
}
