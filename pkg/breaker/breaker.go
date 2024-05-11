package breaker

import (
	"log/slog"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	}
	return "unknown"
}

type CircuitBreaker struct {
	FailureThreshold int
	OpenDuration     time.Duration
	SuccessThreshold int
	Logger           *slog.Logger
	lock             sync.Mutex
	failures         int
	successes        int
	state            State
	openExpiration   time.Time
}

func (c *CircuitBreaker) Do(f func() error) {
	state := c.getState()
	if state == StateOpen {
		return
	}

	err := f()

	c.lock.Lock()
	defer c.lock.Unlock()

	switch state {
	case StateClosed:
		if err != nil {
			c.failures++
			if c.failures >= c.FailureThreshold {
				c.setState(StateOpen)
			}
		} else {
			c.failures = 0
		}
	case StateHalfOpen:
		if err != nil {
			c.setState(StateOpen)
		} else {
			c.successes++
			if c.successes >= c.SuccessThreshold {
				c.setState(StateClosed)
			}
		}
	default:
		// never called
	}
}

func (c *CircuitBreaker) getState() State {
	c.lock.Lock()
	defer c.lock.Unlock()
	state := c.state
	if state == StateOpen && time.Until(c.openExpiration) < 0 {
		state = StateHalfOpen
		c.setState(state)
	}
	return state
}

func (c *CircuitBreaker) setState(state State) {
	c.state = state
	c.successes = 0
	c.failures = 0
	if state == StateOpen {
		c.openExpiration = time.Now().Add(c.OpenDuration)
	}
	if c.Logger != nil {
		c.Logger.Debug("circuit breaker state set", "state", c)
	}
}
