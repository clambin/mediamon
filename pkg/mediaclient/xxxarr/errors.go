package xxxarr

var _ error = &ErrParseFailed{}

type ErrParseFailed struct {
	Err  error
	Body []byte
}

func (e *ErrParseFailed) Error() string {
	return "parse: " + e.Err.Error()
}

func (e *ErrParseFailed) Is(other error) bool {
	if _, ok := other.(*ErrParseFailed); ok {
		return true
	}
	if x, ok := e.Err.(interface{ Is(error) bool }); ok {
		return x.Is(other)
	}
	return false
}

func (e *ErrParseFailed) Unwrap() error { return e.Err }
