package breaker

import "errors"

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
