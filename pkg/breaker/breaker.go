package breaker

import (
	"sync"
	"time"
)

type CircuitBreaker struct {
	Configuration
	lock               sync.Mutex
	counters           Counters
	state              State
	openExpiration     time.Time
	halfOpenExpiration time.Time
}

type Configuration struct {
	FailureThreshold int
	OpenDuration     time.Duration
	SuccessThreshold int
	HalfOpenDuration time.Duration
	ShouldOpen       func(Counters) bool
	ShouldClose      func(Counters) bool
}

func New(configuration Configuration) *CircuitBreaker {
	if configuration.ShouldOpen == nil {
		configuration.ShouldOpen = defaultShouldOpen(configuration)
	}
	if configuration.ShouldClose == nil {
		configuration.ShouldClose = defaultShouldClose(configuration)
	}
	return &CircuitBreaker{
		Configuration: configuration,
	}
}

func defaultShouldOpen(configuration Configuration) func(counters Counters) bool {
	return func(counters Counters) bool {
		return counters.ConsecutiveFailures >= configuration.FailureThreshold
	}
}

func defaultShouldClose(configuration Configuration) func(counters Counters) bool {
	return func(counters Counters) bool {
		return counters.ConsecutiveSuccesses >= configuration.SuccessThreshold
	}
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
	if c.state == StateHalfOpen && c.ShouldClose(c.counters) {
		c.close()
	}
}

func (c *CircuitBreaker) onFailure() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counters.fail()
	// any failure during half-open state immediately opens the circuit again. Too harsh?
	if c.state == StateHalfOpen || c.ShouldOpen(c.counters) {
		c.open()
	}
}

func (c *CircuitBreaker) open() {
	c.setState(StateOpen)
	c.openExpiration = time.Now().Add(c.OpenDuration)
}
func (c *CircuitBreaker) halfOpen() {
	c.setState(StateHalfOpen)
	c.halfOpenExpiration = time.Now().Add(c.HalfOpenDuration)
}
func (c *CircuitBreaker) close() {
	c.setState(StateClosed)
}

func (c *CircuitBreaker) setState(state State) {
	c.state = state
	c.counters.reset()
}

type Counters struct {
	Calls                int
	Successes            int
	ConsecutiveSuccesses int
	Failures             int
	ConsecutiveFailures  int
}

func (c *Counters) pass() {
	c.Calls++
	c.Successes++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counters) fail() {
	c.Calls++
	c.Failures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counters) reset() {
	c.Calls = 0
	c.Successes = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
	c.Failures = 0
	c.Successes = 0
}
