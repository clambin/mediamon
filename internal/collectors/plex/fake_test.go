package plex

import (
	"context"

	"github.com/clambin/mediaclients/plex"
	"github.com/clambin/mediamon/v2/iplocator"
)

var (
	_ libraryGetter  = fakeGetter{}
	_ sessionGetter  = fakeGetter{}
	_ identityGetter = fakeGetter{}
)

type fakeGetter struct {
	libraries []plex.Library
	movies    []plex.Movie
	shows     []plex.Show
	seasons   map[string][]plex.Season
	episodes  map[string][]plex.Episode
	sessions  []plex.Session
	identity  plex.Identity
}

func (f fakeGetter) GetLibraries(_ context.Context) ([]plex.Library, error) {
	return f.libraries, nil
}

func (f fakeGetter) GetMovies(_ context.Context, _ string) ([]plex.Movie, error) {
	return f.movies, nil
}

func (f fakeGetter) GetShows(_ context.Context, _ string) ([]plex.Show, error) {
	return f.shows, nil
}

func (f fakeGetter) GetSeasons(_ context.Context, key string) ([]plex.Season, error) {
	return f.seasons[key], nil
}

func (f fakeGetter) GetEpisodes(_ context.Context, key string) ([]plex.Episode, error) {
	return f.episodes[key], nil
}

func (f fakeGetter) GetSessions(_ context.Context) ([]plex.Session, error) {
	return f.sessions, nil
}

func (f fakeGetter) GetIdentity(_ context.Context) (plex.Identity, error) {
	return f.identity, nil
}

var _ IPLocator = fakeIPLocator{}

type fakeIPLocator struct {
	ips map[string]iplocator.Location
}

func (f fakeIPLocator) Locate(s string) (iplocator.Location, error) {
	return f.ips[s], nil
}
