package xxxarr

import (
	"context"

	"github.com/clambin/mediaclients/radarr"
	"github.com/clambin/mediaclients/sonarr"
)

var _ Client = fakeClient{}

type fakeClient struct {
	version  string
	health   map[string]int
	calendar []string
	queue    []QueuedItem
	library  Library
}

func (f fakeClient) GetVersion(_ context.Context) (string, error) {
	return f.version, nil
}

func (f fakeClient) GetHealth(_ context.Context) (map[string]int, error) {
	return f.health, nil
}

func (f fakeClient) GetCalendar(_ context.Context, _ int) ([]string, error) {
	return f.calendar, nil
}

func (f fakeClient) GetQueue(_ context.Context) ([]QueuedItem, error) {
	return f.queue, nil
}

func (f fakeClient) GetLibrary(_ context.Context) (Library, error) {
	return f.library, nil
}

var _ SonarrClient = fakeSonarrClient{}

type fakeSonarrClient struct {
	systemStatus *sonarr.GetApiV3SystemStatusResponse
	health       *sonarr.GetApiV3HealthResponse
	calendar     *sonarr.GetApiV3CalendarResponse
	queue        *sonarr.GetApiV3QueueResponse
	series       *sonarr.GetApiV3SeriesResponse
}

func (f fakeSonarrClient) GetApiV3SystemStatusWithResponse(_ context.Context, _ ...sonarr.RequestEditorFn) (*sonarr.GetApiV3SystemStatusResponse, error) {
	return f.systemStatus, nil
}

func (f fakeSonarrClient) GetApiV3HealthWithResponse(_ context.Context, _ ...sonarr.RequestEditorFn) (*sonarr.GetApiV3HealthResponse, error) {
	return f.health, nil
}

func (f fakeSonarrClient) GetApiV3CalendarWithResponse(_ context.Context, _ *sonarr.GetApiV3CalendarParams, _ ...sonarr.RequestEditorFn) (*sonarr.GetApiV3CalendarResponse, error) {
	return f.calendar, nil
}

func (f fakeSonarrClient) GetApiV3QueueWithResponse(_ context.Context, _ *sonarr.GetApiV3QueueParams, _ ...sonarr.RequestEditorFn) (*sonarr.GetApiV3QueueResponse, error) {
	return f.queue, nil
}

func (f fakeSonarrClient) GetApiV3SeriesWithResponse(_ context.Context, _ *sonarr.GetApiV3SeriesParams, _ ...sonarr.RequestEditorFn) (*sonarr.GetApiV3SeriesResponse, error) {
	return f.series, nil
}

var _ RadarrClient = fakeRadarrClient{}

type fakeRadarrClient struct {
	systemStatus *radarr.GetApiV3SystemStatusResponse
	health       *radarr.GetApiV3HealthResponse
	calendar     *radarr.GetApiV3CalendarResponse
	queue        *radarr.GetApiV3QueueResponse
	movies       *radarr.GetApiV3MovieResponse
}

func (f fakeRadarrClient) GetApiV3SystemStatusWithResponse(_ context.Context, _ ...radarr.RequestEditorFn) (*radarr.GetApiV3SystemStatusResponse, error) {
	return f.systemStatus, nil
}

func (f fakeRadarrClient) GetApiV3HealthWithResponse(_ context.Context, _ ...radarr.RequestEditorFn) (*radarr.GetApiV3HealthResponse, error) {
	return f.health, nil
}

func (f fakeRadarrClient) GetApiV3CalendarWithResponse(_ context.Context, _ *radarr.GetApiV3CalendarParams, _ ...radarr.RequestEditorFn) (*radarr.GetApiV3CalendarResponse, error) {
	return f.calendar, nil
}

func (f fakeRadarrClient) GetApiV3QueueWithResponse(_ context.Context, _ *radarr.GetApiV3QueueParams, _ ...radarr.RequestEditorFn) (*radarr.GetApiV3QueueResponse, error) {
	return f.queue, nil
}

func (f fakeRadarrClient) GetApiV3MovieWithResponse(_ context.Context, _ *radarr.GetApiV3MovieParams, _ ...radarr.RequestEditorFn) (*radarr.GetApiV3MovieResponse, error) {
	return f.movies, nil
}
