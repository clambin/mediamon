package breaker

import (
	"sync"
	"time"
)

type CircuitBreaker struct {
	Configuration
	lock           sync.Mutex
	counters       Counters
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
	if c.state == StateOpen && time.Until(c.openExpiration) < 0 {
		c.halfOpen()
	}
	return c.state
}

func (c *CircuitBreaker) onSuccess() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counters.pass()
	if c.state == StateHalfOpen && c.counters.ConsecutiveSuccesses >= c.SuccessThreshold {
		c.close()
	}
}

func (c *CircuitBreaker) onFailure() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counters.fail()
	if c.state == StateHalfOpen || c.counters.ConsecutiveFailures >= c.FailureThreshold {
		c.open()
	}
}

func (c *CircuitBreaker) open() {
	c.setState(StateOpen)
	c.openExpiration = time.Now().Add(c.OpenDuration)
}
func (c *CircuitBreaker) halfOpen() {
	c.setState(StateHalfOpen)
}
func (c *CircuitBreaker) close() {
	c.setState(StateClosed)
}

func (c *CircuitBreaker) setState(state State) {
	c.state = state
	c.counters.reset()
}
