package clients

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithToken(t *testing.T) {
	const token = "1234"
	ctx := context.Background()

	f := WithToken(token)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, f(ctx, r))
	assert.Equal(t, token, r.Header.Get("X-API-KEY"))

	f = WithToken("")
	r, _ = http.NewRequest(http.MethodGet, "/", nil)
	assert.Error(t, f(ctx, r))
}

// constP is used by tests to convert a constant to a pointer to that value
func constP[T any](t T) *T {
	return &t
}
