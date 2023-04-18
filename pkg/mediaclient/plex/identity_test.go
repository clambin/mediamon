package plex_test

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexClient_GetIdentity(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()

	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	c := plex.Client{
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	identity, err := c.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "SomeVersion", identity.Version)
}
