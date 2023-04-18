package plex_test

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexClient_Authentication(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))

	c := plex.Client{
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "badpassword",
	}

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403 Forbidden")

	c.SetAuthToken("some_token")
	_, err = c.GetIdentity(context.Background())
	require.NoError(t, err)

	authServer.Close()
	c.SetAuthToken("")
	_, err = c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, unix.ECONNREFUSED)
}
