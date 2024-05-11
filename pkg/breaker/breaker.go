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

var stateStrings = map[State]string{
	StateClosed:   "closed",
	StateHalfOpen: "half-open",
	StateOpen:     "open",
}

func (s State) String() string {
	if value, ok := stateStrings[s]; ok {
		return value
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

func New(failureThreshold int, openDuration time.Duration, successThreshold int, logger *slog.Logger) *CircuitBreaker {
	if logger == nil {
		logger = slog.Default()
	}
	return &CircuitBreaker{
		FailureThreshold: failureThreshold,
		OpenDuration:     openDuration,
		SuccessThreshold: successThreshold,
		Logger:           logger,
	}
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
			// one error during half-open state opens the circuit again
			// alternatively: one move after c.FailureThreshold errors?
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
	c.Logger.Debug("circuit breaker", "state", c.state)
}
