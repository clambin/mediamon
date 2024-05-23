package breaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

var ErrCircuitOpen = errors.New("circuit is open")

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
	Configuration
	lock           sync.Mutex
	failures       int
	successes      int
	state          State
	openExpiration time.Time
}

type Configuration struct {
	FailureThreshold int
	OpenDuration     time.Duration
	SuccessThreshold int
}

func (c *CircuitBreaker) Do(f func() error) error {
	if state := c.getState(); state == StateOpen {
		return ErrCircuitOpen
	}

	err := f()
	if err == nil {
		c.onSuccess()
	} else {
		c.onFailure()
	}
	return err
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

func (c *CircuitBreaker) onSuccess() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.successes++
	c.failures = 0
	if c.state == StateHalfOpen && c.successes >= c.SuccessThreshold {
		c.setState(StateClosed)
	}
}

func (c *CircuitBreaker) onFailure() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.successes = 0
	c.failures++
	if c.state == StateHalfOpen || c.failures >= c.FailureThreshold {
		c.setState(StateOpen)
	}
}

func (c *CircuitBreaker) setState(state State) {
	c.state = state
	c.successes = 0
	c.failures = 0
	if state == StateOpen {
		c.openExpiration = time.Now().Add(c.OpenDuration)
	}
}
