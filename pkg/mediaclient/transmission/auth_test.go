package transmission

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthenticator_RoundTrip(t *testing.T) {
	h := server{sessionID: "1234"}
	s := httptest.NewServer(&h)
	defer s.Close()

	c := http.Client{Transport: &authenticator{next: http.DefaultTransport}}
	resp, err := c.Post(s.URL, "application/json", io.NopCloser(strings.NewReader("foo")))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "1234", resp.Header.Get(transmissionSessionIDHeader))

	// simulate an expired session ID
	h.sessionID = "4321"
	resp, err = c.Post(s.URL, "application/json", io.NopCloser(strings.NewReader("foo")))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "4321", resp.Header.Get(transmissionSessionIDHeader))
}

type server struct {
	sessionID string
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(transmissionSessionIDHeader, s.sessionID)

	if r.Header.Get(transmissionSessionIDHeader) != s.sessionID {
		w.WriteHeader(http.StatusConflict)
		return
	}
	body, _ := io.ReadAll(r.Body)
	if string(body) != "foo" {
		w.WriteHeader(http.StatusBadRequest)
	}
}
