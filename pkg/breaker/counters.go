package breaker

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
