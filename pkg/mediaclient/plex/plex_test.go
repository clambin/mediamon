package plex

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_WithRoundTripper(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(testutil.AuthHandler))
	defer authServer.Close()

	server := httptest.NewServer(testutil.WithToken("some_token", func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Dummy") != "foo" {
			http.Error(writer, "missing X-Dummy header", http.StatusBadRequest)
			return
		}
		testutil.Handler(writer, request)
	}))
	defer server.Client()

	c := New("user@example.com", "somepassword", "", "", server.URL, &dummyRoundTripper{next: http.DefaultTransport})
	c.HTTPClient.Transport.(*authenticator).authURL = authServer.URL

	_, err := c.GetSessions(context.Background())
	assert.NoError(t, err)
}

var _ http.RoundTripper = &dummyRoundTripper{}

type dummyRoundTripper struct {
	next http.RoundTripper
}

func (d *dummyRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("X-Dummy", "foo")
	return d.next.RoundTrip(request)
}

func TestClient_Failures(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "server's having a hard day", http.StatusInternalServerError)
	}))

	c := New("user@example.com", "somepassword", "", "", testServer.URL, nil)
	c.HTTPClient.Transport = http.DefaultTransport

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 "+http.StatusText(http.StatusInternalServerError), err.Error())

	testServer.Close()
	_, err = c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, unix.ECONNREFUSED)
}

func TestClient_Decode_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("this is definitely not json"))
	}))
	defer testServer.Close()

	c := New("user@example.com", "somepassword", "", "", testServer.URL, nil)
	c.HTTPClient.Transport = http.DefaultTransport

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "decode: invalid character 'h' in literal true (expecting 'r')", err.Error())
}
