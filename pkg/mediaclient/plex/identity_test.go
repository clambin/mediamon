package plex_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexClient_GetIdentity(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(testutil.Handler))
	defer testServer.Close()

	c := plex.New("user@example.com", "somepassword", "", "", testServer.URL, nil)
	c.HTTPClient.Transport = http.DefaultTransport

	identity, err := c.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "SomeVersion", identity.Version)
}
