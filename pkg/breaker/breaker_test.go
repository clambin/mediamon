package breaker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"math/rand"
	"testing"
	"time"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		name string
		s    State
		want string
	}{
		{
			name: "closed",
			s:    StateClosed,
			want: "closed",
		},
		{
			name: "open",
			s:    StateOpen,
			want: "open",
		},
		{
			name: "half-open",
			s:    StateHalfOpen,
			want: "half-open",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.String())
		})
	}
}

func TestCircuitBreaker_Do(t *testing.T) {
	cfg := Configuration{
		FailureThreshold: 5,
		OpenDuration:     500 * time.Millisecond,
		SuccessThreshold: 10,
	}
	newCB := func(state State) *CircuitBreaker {
		cb := New(cfg, slog.Default())
		cb.setState(state)
		return cb
	}

	t.Run("closed cb remains closed on success", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateClosed)
		cb.Do(func() error { return nil })
		assert.Equal(t, StateClosed, cb.getState())
	})

	t.Run("closed cb closes after FailureThreshold errors", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateClosed)
		var calls int
		for cb.getState() == StateClosed {
			cb.Do(func() error { return errors.New("error") })
			calls++
		}
		assert.Equal(t, cb.FailureThreshold, calls)
	})

	t.Run("closed cb only closes after FailureThreshold consecutive errors", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateClosed)
		for range cb.FailureThreshold {
			cb.Do(func() error { return errors.New("error") })
			cb.Do(func() error { return nil })
		}
		assert.Equal(t, StateClosed, cb.getState())
	})

	t.Run("open cb doesn't perform the call", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateOpen)
		cb.OpenDuration = time.Hour
		var calls int
		for range 10 {
			cb.Do(func() error {
				calls++
				return nil
			})
		}
		assert.Zero(t, calls)
	})

	t.Run("open cb goes half-open after OpenDuration", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateOpen)
		start := time.Now()
		for cb.getState() == StateOpen {
			time.Sleep(100 * time.Millisecond)
		}
		assert.GreaterOrEqual(t, time.Since(start), cb.OpenDuration)
	})

	t.Run("half-open cb opens on error", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateHalfOpen)
		cb.Do(func() error { return errors.New("error") })
		assert.Equal(t, StateOpen, cb.getState())
	})

	t.Run("half-open cb closes after SuccessThreshold successes", func(t *testing.T) {
		t.Parallel()
		cb := newCB(StateHalfOpen)
		var calls int
		for cb.getState() == StateHalfOpen {
			cb.Do(func() error { return nil })
			calls++
		}
		assert.Equal(t, StateClosed, cb.getState())
		assert.Equal(t, cb.SuccessThreshold, calls)
	})
}

func BenchmarkCircuitBreaker_Do(b *testing.B) {
	cb := CircuitBreaker{
		Configuration: Configuration{
			FailureThreshold: 5,
			OpenDuration:     time.Millisecond,
			SuccessThreshold: 10,
		},
		Logger: slog.Default(),
	}

	for range b.N {
		cb.Do(func() error {
			if rand.Intn(10) < 3 {
				return errors.New("error")
			}
			return nil
		})
	}
}
