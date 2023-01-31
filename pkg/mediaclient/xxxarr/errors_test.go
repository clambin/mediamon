package xxxarr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
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
				assert.ErrorIs(t, newErr, &ErrInvalidJSON{})
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
