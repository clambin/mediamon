package clients

import (
	"context"
	"errors"
	"github.com/clambin/mediaclients/radarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRadarrGetVersion(t *testing.T) {
	r := mocks.NewRadarrClient(t)
	c := Radarr{Client: r}

	ctx := context.Background()
	r.EXPECT().
		GetApiV3SystemStatusWithResponse(ctx).
		Return(&radarr.GetApiV3SystemStatusResponse{JSON200: &radarr.SystemResource{Version: constP("1.0")}}, nil).
		Once()
	version, err := c.GetVersion(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "1.0", version)

	r.EXPECT().
		GetApiV3SystemStatusWithResponse(ctx).
		Return(nil, errors.New("blah"))
	_, err = c.GetVersion(ctx)
	assert.Error(t, err)
}

func TestRadarrGetHealth(t *testing.T) {
	r := mocks.NewRadarrClient(t)
	c := Radarr{Client: r}

	ctx := context.Background()
	healthResources := []radarr.HealthResource{
		{Type: constP[radarr.HealthCheckResult]("foo")},
		{Type: constP[radarr.HealthCheckResult]("bar")},
		{Type: constP[radarr.HealthCheckResult]("bar")},
	}
	r.EXPECT().
		GetApiV3HealthWithResponse(ctx).
		Return(&radarr.GetApiV3HealthResponse{JSON200: &healthResources}, nil).
		Once()
	health, err := c.GetHealth(ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"foo": 1, "bar": 2}, health)

	r.EXPECT().GetApiV3HealthWithResponse(ctx).Return(nil, errors.New("blah"))
	_, err = c.GetHealth(ctx)
	assert.Error(t, err)
}

func TestRadarrGetCalendar(t *testing.T) {
	r := mocks.NewRadarrClient(t)
	c := Radarr{Client: r}

	ctx := context.Background()
	movies := []radarr.MovieResource{
		{Title: constP("movie 1")},
		{Title: constP("movie 2")},
	}
	r.EXPECT().
		GetApiV3CalendarWithResponse(ctx, mock.AnythingOfType("*radarr.GetApiV3CalendarParams")).
		Return(&radarr.GetApiV3CalendarResponse{JSON200: &movies}, nil).
		Once()
	resp, err := c.GetCalendar(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"movie 1", "movie 2"}, resp)

	r.EXPECT().
		GetApiV3CalendarWithResponse(ctx, mock.AnythingOfType("*radarr.GetApiV3CalendarParams")).
		Return(nil, errors.New("blah"))
	_, err = c.GetCalendar(ctx, 1)
	assert.Error(t, err)
}

func TestRadarrGetQueue(t *testing.T) {
	r := mocks.NewRadarrClient(t)
	c := Radarr{Client: r}

	ctx := context.Background()
	r.EXPECT().
		GetApiV3QueueWithResponse(ctx, mock.AnythingOfType("*radarr.GetApiV3QueueParams")).
		RunAndReturn(func(ctx context.Context, params *radarr.GetApiV3QueueParams, fn ...radarr.RequestEditorFn) (*radarr.GetApiV3QueueResponse, error) {
			var resp = []radarr.GetApiV3QueueResponse{
				{
					JSON200: &radarr.QueueResourcePagingResource{
						Page:         constP[int32](0),
						PageSize:     constP[int32](100),
						TotalRecords: constP[int32](2),
						Records:      &[]radarr.QueueResource{{Title: constP("movie 1"), Size: constP(100.0), Sizeleft: constP(25.0)}},
					},
				},
				{
					JSON200: &radarr.QueueResourcePagingResource{
						Page:         constP[int32](1),
						PageSize:     constP[int32](100),
						TotalRecords: constP[int32](2),
						Records:      &[]radarr.QueueResource{{Title: constP("movie 2"), Size: constP(100.0), Sizeleft: constP(50.0)}},
					},
				},
			}
			return &resp[*params.Page], nil
		}).
		Twice()
	resp, err := c.GetQueue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []QueuedItem{{Name: "movie 1", TotalBytes: 100, DownloadedBytes: 75}, {Name: "movie 2", TotalBytes: 100, DownloadedBytes: 50}}, resp)

	r.EXPECT().GetApiV3QueueWithResponse(ctx, mock.AnythingOfType("*radarr.GetApiV3QueueParams")).Return(nil, errors.New("blah"))
	_, err = c.GetQueue(ctx)
	assert.Error(t, err)
}

func TestRadarrGetLibrary(t *testing.T) {
	r := mocks.NewRadarrClient(t)
	c := Radarr{Client: r}

	ctx := context.Background()
	r.EXPECT().
		GetApiV3MovieWithResponse(ctx, mock.AnythingOfType("*radarr.GetApiV3MovieParams")).
		Return(&radarr.GetApiV3MovieResponse{
			JSON200: &[]radarr.MovieResource{
				{Title: constP("movie 1"), Monitored: constP(false)},
				{Title: constP("movie 2"), Monitored: constP(true)},
				{Title: constP("movie 3"), Monitored: constP(true)},
			},
		}, nil).
		Once()
	resp, err := c.GetLibrary(ctx)
	assert.NoError(t, err)
	assert.Equal(t, Library{Monitored: 2, Unmonitored: 1}, resp)

	r.EXPECT().
		GetApiV3MovieWithResponse(ctx, mock.AnythingOfType("*radarr.GetApiV3MovieParams")).
		Return(nil, errors.New("blah")).
		Once()
	_, err = c.GetLibrary(ctx)
	assert.Error(t, err)
}
