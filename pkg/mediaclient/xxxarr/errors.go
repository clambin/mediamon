package xxxarr

import (
	"fmt"
	"net/http"
)

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

var _ error = &ErrHTTPFailed{}

type ErrHTTPFailed struct {
	StatusCode int
	Status     string
}

func (e *ErrHTTPFailed) Error() string {
	if e.Status == "" {
		return fmt.Sprintf("%d - %s", e.StatusCode, http.StatusText(e.StatusCode))
	}
	return e.Status
}

func (e *ErrHTTPFailed) Is(target error) bool {
	if t, ok := target.(*ErrHTTPFailed); ok {
		if t == nil {
			return true
		}
		return (e.StatusCode == t.StatusCode || t.StatusCode == 0) &&
			(e.Status == t.Status || t.Status == "")
	}
	return false
}
