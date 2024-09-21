package clients

import (
	"cmp"
	"context"
	"errors"
	"github.com/clambin/mediaclients/sonarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/maps"
	"slices"
	"testing"
)

func TestSonarrGetVersion(t *testing.T) {
	s := mocks.NewSonarrClient(t)
	c := Sonarr{Client: s}

	ctx := context.Background()
	s.EXPECT().
		GetApiV3SystemStatusWithResponse(ctx).
		Return(&sonarr.GetApiV3SystemStatusResponse{JSON200: &sonarr.SystemResource{Version: constP("1.0")}}, nil).
		Once()
	version, err := c.GetVersion(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "1.0", version)

	s.EXPECT().
		GetApiV3SystemStatusWithResponse(ctx).
		Return(nil, errors.New("blah"))
	_, err = c.GetVersion(ctx)
	assert.Error(t, err)
}

func TestSonarrGetHealth(t *testing.T) {
	s := mocks.NewSonarrClient(t)
	c := Sonarr{Client: s}

	ctx := context.Background()
	healthResources := []sonarr.HealthResource{
		{Type: constP[sonarr.HealthCheckResult]("foo")},
		{Type: constP[sonarr.HealthCheckResult]("bar")},
		{Type: constP[sonarr.HealthCheckResult]("bar")},
	}
	s.EXPECT().
		GetApiV3HealthWithResponse(ctx).
		Return(&sonarr.GetApiV3HealthResponse{JSON200: &healthResources}, nil).
		Once()
	health, err := c.GetHealth(ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"foo": 1, "bar": 2}, health)

	s.EXPECT().GetApiV3HealthWithResponse(ctx).Return(nil, errors.New("blah"))
	_, err = c.GetHealth(ctx)
	assert.Error(t, err)
}

func TestSonarrGetCalendar(t *testing.T) {
	s := mocks.NewSonarrClient(t)
	c := Sonarr{Client: s}

	ctx := context.Background()
	episodes := map[int32]sonarr.EpisodeResource{
		11: {Title: constP("series 1 season 5 episode 1"), Id: constP[int32](11), SeriesId: constP[int32](1), SeasonNumber: constP[int32](5), EpisodeNumber: constP[int32](1)},
		21: {Title: constP("series 2 season 9 episode 1"), Id: constP[int32](21), SeriesId: constP[int32](2), SeasonNumber: constP[int32](9), EpisodeNumber: constP[int32](1)},
	}
	allEpisodes := maps.Values(episodes)
	s.EXPECT().
		GetApiV3CalendarWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3CalendarParams")).
		Return(&sonarr.GetApiV3CalendarResponse{JSON200: &allEpisodes}, nil).
		Once()
	s.EXPECT().
		GetApiV3SeriesIdWithResponse(ctx, int32(1), mock.AnythingOfType("*sonarr.GetApiV3SeriesIdParams")).
		Return(&sonarr.GetApiV3SeriesIdResponse{
			JSON200: &sonarr.SeriesResource{Title: constP("series 1")},
		}, nil).
		Once()
	s.EXPECT().
		GetApiV3SeriesIdWithResponse(ctx, int32(2), mock.AnythingOfType("*sonarr.GetApiV3SeriesIdParams")).
		Return(&sonarr.GetApiV3SeriesIdResponse{
			JSON200: &sonarr.SeriesResource{Title: constP("series 2")},
		}, nil).
		Once()
	resp, err := c.GetCalendar(ctx, 1)
	assert.NoError(t, err)
	slices.Sort(resp)
	assert.Equal(t, []string{
		"series 1 - S05E01 - series 1 season 5 episode 1",
		"series 2 - S09E01 - series 2 season 9 episode 1",
	}, resp)

	s.EXPECT().
		GetApiV3CalendarWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3CalendarParams")).
		Return(nil, errors.New("blah"))
	_, err = c.GetCalendar(ctx, 1)
	assert.Error(t, err)
}

func TestSonarrGetQueue(t *testing.T) {
	s := mocks.NewSonarrClient(t)
	c := Sonarr{Client: s}

	ctx := context.Background()
	s.EXPECT().
		GetApiV3QueueWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3QueueParams")).
		RunAndReturn(func(ctx context.Context, params *sonarr.GetApiV3QueueParams, fn ...sonarr.RequestEditorFn) (*sonarr.GetApiV3QueueResponse, error) {
			var resp = []sonarr.GetApiV3QueueResponse{
				{
					JSON200: &sonarr.QueueResourcePagingResource{
						Page:         constP[int32](0),
						PageSize:     constP[int32](100),
						TotalRecords: constP[int32](2),
						Records: &[]sonarr.QueueResource{
							{Title: constP("show 1 season 1 episode 1"), SeriesId: constP[int32](1), SeasonNumber: constP[int32](1), EpisodeId: constP[int32](111), Size: constP(100.0), Sizeleft: constP(25.0)},
						},
					},
				},
				{
					JSON200: &sonarr.QueueResourcePagingResource{
						Page:         constP[int32](1),
						PageSize:     constP[int32](100),
						TotalRecords: constP[int32](2),
						Records: &[]sonarr.QueueResource{
							{Title: constP("show 2 season 2 episode 2"), SeriesId: constP[int32](2), SeasonNumber: constP[int32](2), EpisodeId: constP[int32](222), Size: constP(100.0), Sizeleft: constP(50.0)},
						},
					},
				},
			}
			return &resp[*params.Page], nil
		}).
		Twice()
	s.EXPECT().
		GetApiV3SeriesIdWithResponse(ctx, int32(1), mock.AnythingOfType("*sonarr.GetApiV3SeriesIdParams")).
		Return(&sonarr.GetApiV3SeriesIdResponse{
			JSON200: &sonarr.SeriesResource{Title: constP("series 1")},
		}, nil).
		Once()
	s.EXPECT().
		GetApiV3SeriesIdWithResponse(ctx, int32(2), mock.AnythingOfType("*sonarr.GetApiV3SeriesIdParams")).
		Return(&sonarr.GetApiV3SeriesIdResponse{
			JSON200: &sonarr.SeriesResource{Title: constP("series 2")},
		}, nil).
		Once()
	s.EXPECT().
		GetApiV3EpisodeWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3EpisodeParams")).
		RunAndReturn(func(ctx context.Context, params *sonarr.GetApiV3EpisodeParams, fn ...sonarr.RequestEditorFn) (*sonarr.GetApiV3EpisodeResponse, error) {
			switch (*params.EpisodeIds)[0] {
			case int32(111):
				return &sonarr.GetApiV3EpisodeResponse{JSON200: &[]sonarr.EpisodeResource{{EpisodeNumber: constP[int32](1)}}}, nil
			case int32(222):
				return &sonarr.GetApiV3EpisodeResponse{JSON200: &[]sonarr.EpisodeResource{{EpisodeNumber: constP[int32](2)}}}, nil
			default:
				return nil, errors.New("blah")
			}
		}).
		Twice()
	resp, err := c.GetQueue(ctx)
	assert.NoError(t, err)
	slices.SortFunc(resp, func(a, b QueuedItem) int { return cmp.Compare(a.Name, b.Name) })
	assert.Equal(t, []QueuedItem{
		{Name: "series 1 - S01E01 - show 1 season 1 episode 1", TotalBytes: 100, DownloadedBytes: 75},
		{Name: "series 2 - S02E02 - show 2 season 2 episode 2", TotalBytes: 100, DownloadedBytes: 50},
	}, resp)

	s.EXPECT().GetApiV3QueueWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3QueueParams")).Return(nil, errors.New("blah"))
	_, err = c.GetQueue(ctx)
	assert.Error(t, err)
}

func TestSonarrGetLibrary(t *testing.T) {
	s := mocks.NewSonarrClient(t)
	c := Sonarr{Client: s}

	ctx := context.Background()
	s.EXPECT().
		GetApiV3SeriesWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3SeriesParams")).
		Return(&sonarr.GetApiV3SeriesResponse{
			JSON200: &[]sonarr.SeriesResource{
				{Title: constP("series 1"), Monitored: constP(false)},
				{Title: constP("series 2"), Monitored: constP(true)},
				{Title: constP("series 3"), Monitored: constP(true)},
			},
		}, nil).
		Once()
	resp, err := c.GetLibrary(ctx)
	assert.NoError(t, err)
	assert.Equal(t, Library{Monitored: 2, Unmonitored: 1}, resp)

	s.EXPECT().
		GetApiV3SeriesWithResponse(ctx, mock.AnythingOfType("*sonarr.GetApiV3SeriesParams")).
		Return(nil, errors.New("blah")).
		Once()
	_, err = c.GetLibrary(ctx)
	assert.Error(t, err)
}
