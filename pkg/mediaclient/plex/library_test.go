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

func TestClient_GetLibraries(t *testing.T) {
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

	libraries, err := c.GetLibraries(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []plex.LibrariesDirectory{
		{Key: "1", Type: "movie", Title: "Movies"},
		{Key: "2", Type: "show", Title: "Shows"},
	}, libraries.Directory)
}

func TestClient_GetMovieLibrary(t *testing.T) {
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

	movies, err := c.GetMovieLibrary(context.Background(), "1")
	require.NoError(t, err)
	assert.Equal(t, plex.MovieLibrary{Metadata: []plex.MovieLibraryEntry{{Guid: "1", Title: "foo"}}}, movies)
}

func TestClient_GetShowLibrary(t *testing.T) {
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

	movies, err := c.GetShowLibrary(context.Background(), "2")
	require.NoError(t, err)
	assert.Equal(t, plex.ShowLibrary{Metadata: []plex.ShowLibraryEntry{{Guid: "2", Title: "bar"}}}, movies)
}
