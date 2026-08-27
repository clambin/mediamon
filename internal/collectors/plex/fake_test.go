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
	movies    []plex.MediaMetadata
	episodes  []plex.MediaMetadata
	sessions  []plex.Session
	identity  plex.Identity
}

func (f fakeGetter) GetLibraries(_ context.Context) ([]plex.Library, error) {
	return f.libraries, nil
}

func (f fakeGetter) GetAllLibraryMedia(_ context.Context, _ string) ([]plex.MediaMetadata, error) {
	return append(f.movies, f.episodes...), nil
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
