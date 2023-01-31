package xxxarr

var _ error = &ErrInvalidJSON{}

type ErrInvalidJSON struct {
	Err  error
	Body []byte
}

func (e *ErrInvalidJSON) Error() string {
	return "parse: " + e.Err.Error()
}

func (e *ErrInvalidJSON) Is(target error) bool {
	_, ok := target.(*ErrInvalidJSON)
	return ok
}

func (e *ErrInvalidJSON) Unwrap() error {
	return e.Err
}
