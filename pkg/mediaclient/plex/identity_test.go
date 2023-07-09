package plex_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlexClient_GetIdentity(t *testing.T) {
	c, s := makeClientAndServer(nil)
	defer s.Close()

	identity, err := c.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "SomeVersion", identity.Version)
}
