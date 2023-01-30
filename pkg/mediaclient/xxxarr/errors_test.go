package xxxarr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestErrParseFailed(t *testing.T) {
	e := &ErrParseFailed{
		Err:  errors.New("bad content"),
		Body: []byte("hello"),
	}
	assert.Error(t, e)
	assert.Equal(t, "parse: bad content", e.Error())

	e2 := fmt.Errorf("error: %w", e)
	assert.Equal(t, "error: parse: bad content", e2.Error())

	var e3 *ErrParseFailed
	assert.True(t, errors.Is(e, e3))
	assert.True(t, errors.Is(e2, e3))

	var e4 *ErrParseFailed
	assert.True(t, errors.As(e2, &e4))
	assert.Equal(t, "parse: bad content", e4.Error())
	assert.Equal(t, "hello", string(e4.Body))
}

func TestErrHTTPFailed(t *testing.T) {
	code := http.StatusInternalServerError
	e := &ErrHTTPFailed{
		StatusCode: code,
	}
	assert.Error(t, e)
	assert.Equal(t, fmt.Sprintf("%d - %s", code, http.StatusText(code)), e.Error())

	e.Status = "foo"
	assert.Equal(t, "foo", e.Error())

	e2 := fmt.Errorf("wrapper: %w", e)

	var e3 *ErrHTTPFailed
	assert.True(t, errors.Is(e2, e3))
	assert.True(t, errors.As(e2, &e3))
	assert.Equal(t, "foo", e3.Error())
}
