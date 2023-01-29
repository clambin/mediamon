package xxxarr

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrParseFailed(t *testing.T) {
	e := &ErrParseFailed{
		Err:  errors.New("bad content"),
		Body: []byte("hello"),
	}

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
