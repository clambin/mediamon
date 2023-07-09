package plex_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Failures(t *testing.T) {
	c, s := makeClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "server's having a hard day", http.StatusInternalServerError)
	}))

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 "+http.StatusText(http.StatusInternalServerError), err.Error())

	s.Close()
	_, err = c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, unix.ECONNREFUSED)
}

func TestClient_Decode_Failure(t *testing.T) {
	c, s := makeClientAndServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("this is definitely not json"))
	}))
	defer s.Close()

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "decode: invalid character 'h' in literal true (expecting 'r')", err.Error())
}

func makeClientAndServer(h http.Handler) (*plex.Client, *httptest.Server) {
	if h == nil {
		h = http.HandlerFunc(testutil.Handler)
	}
	s := httptest.NewServer(h)
	c := plex.New("user@example.com", "somepassword", "", "", s.URL, nil)
	// cut out the authenticator
	c.HTTPClient.Transport = http.DefaultTransport
	return c, s
}
