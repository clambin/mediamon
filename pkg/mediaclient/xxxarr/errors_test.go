package xxxarr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestErrInvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		other    error
		expectIs bool
		expectAs bool
	}{
		{
			name:     "direct",
			err:      &ErrInvalidJSON{Err: errors.New("bad content"), Body: []byte("hello")},
			other:    &ErrInvalidJSON{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "wrapped",
			err:      fmt.Errorf("error: %w", &ErrInvalidJSON{Err: errors.New("bad content"), Body: []byte("hello")}),
			other:    &ErrInvalidJSON{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "is not",
			err:      &ErrInvalidJSON{},
			other:    errors.New("error"),
			expectIs: false,
			expectAs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.err)
			assert.Equal(t, tt.expectIs, errors.Is(tt.err, tt.other))
			var newErr *ErrInvalidJSON
			assert.Equal(t, tt.expectAs, errors.As(tt.err, &newErr))
			if tt.expectAs {
				assert.True(t, errors.Is(newErr, &ErrInvalidJSON{}))
			}
		})
	}
}

func TestErrInvalidJSON_Unwrap(t *testing.T) {
	e := &ErrInvalidJSON{
		Err: errors.New("error"),
	}
	e2 := errors.Unwrap(e)
	assert.Equal(t, "error", e2.Error())
}

func TestErrHTTPFailed(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		errorString string
		other       error
		expectIs    bool
		expectAs    bool
	}{
		{
			name:        "simple",
			err:         &ErrHTTPFailed{Status: "error"},
			errorString: "error",
			other:       &ErrHTTPFailed{},
			expectIs:    true,
			expectAs:    true,
		},
		{
			name:        "httpCode",
			err:         &ErrHTTPFailed{StatusCode: http.StatusForbidden},
			errorString: fmt.Sprintf("%d - %s", http.StatusForbidden, http.StatusText(http.StatusForbidden)),
			other:       &ErrHTTPFailed{},
			expectIs:    true,
			expectAs:    true,
		},
		{
			name:     "statuscode comparison",
			err:      &ErrHTTPFailed{StatusCode: http.StatusInternalServerError},
			other:    &ErrHTTPFailed{StatusCode: http.StatusInternalServerError},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "statuscode mismatch",
			err:      &ErrHTTPFailed{StatusCode: http.StatusForbidden},
			other:    &ErrHTTPFailed{StatusCode: http.StatusInternalServerError},
			expectIs: false,
			expectAs: true,
		},
		{
			name:     "status comparison",
			err:      &ErrHTTPFailed{Status: "error 1"},
			other:    &ErrHTTPFailed{Status: "error 1"},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "status mismatch",
			err:      &ErrHTTPFailed{Status: "error 1"},
			other:    &ErrHTTPFailed{Status: "error 2"},
			expectIs: false,
			expectAs: true,
		},
		{
			name:     "wrapped",
			err:      fmt.Errorf("error: %w", &ErrHTTPFailed{Status: "error"}),
			other:    &ErrHTTPFailed{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "is not",
			err:      &ErrHTTPFailed{StatusCode: http.StatusForbidden},
			other:    errors.New("error"),
			expectIs: false,
			expectAs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.err)
			assert.Equal(t, tt.expectIs, errors.Is(tt.err, tt.other))
			var newErr *ErrHTTPFailed
			assert.Equal(t, tt.expectAs, errors.As(tt.err, &newErr))
			if tt.expectAs {
				assert.True(t, errors.Is(newErr, &ErrHTTPFailed{}))
			}
		})
	}

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

	e3 = new(ErrHTTPFailed)
	assert.True(t, errors.Is(e2, e3))
	assert.True(t, errors.As(e2, &e3))
	assert.Equal(t, "foo", e3.Error())

	//e3.StatusCode = http.StatusOK
	assert.True(t, errors.Is(e2, e3))
}
