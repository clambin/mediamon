package xxxarr

import (
	"context"
	"fmt"
	"github.com/clambin/go-metrics/caller"
	"net/http"
)

// RadarrAPI contains all supported Radarr APIs
//go:generate mockery --name RadarrAPI
type RadarrAPI interface {
	GetURL() (url string)
	GetSystemStatus(ctx context.Context) (response RadarrSystemStatusResponse, err error)
	GetCalendar(ctx context.Context) (response []RadarrCalendarResponse, err error)
	GetQueuePage(ctx context.Context, pageNr int) (response RadarrQueueResponse, err error)
	GetQueue(ctx context.Context) (response RadarrQueueResponse, err error)
	GetMovies(ctx context.Context) (response []RadarrMovieResponse, err error)
	GetMovieByID(ctx context.Context, movieID int) (response RadarrMovieResponse, err error)
}

// RadarrClient calls Radarr endpoints
type RadarrClient struct {
	APICaller
}

var _ RadarrAPI = &RadarrClient{}

// NewRadarrClient creates a new RadarrClient, using http.DefaultClient as http.InstrumentedClient
func NewRadarrClient(apiKey, url string, options caller.Options) *RadarrClient {
	return NewRadarrClientWithCaller(apiKey, url, &caller.InstrumentedClient{
		BaseClient:  caller.BaseClient{HTTPClient: http.DefaultClient},
		Options:     options,
		Application: "radarr",
	})
}

// NewRadarrClientWithCaller creates a new RadarrClient with a specified caller
func NewRadarrClientWithCaller(apiKey, url string, caller caller.Caller) *RadarrClient {
	return &RadarrClient{APICaller: &APIClient{
		Caller: caller,
		URL:    url,
		APIKey: apiKey,
	}}
}

// GetSystemStatus calls Radarr's  /api/v3/system/status endpoint. It returns the system status of the Radarr instance
func (rc RadarrClient) GetSystemStatus(ctx context.Context) (response RadarrSystemStatusResponse, err error) {
	err = rc.Get(ctx, "/api/v3/system/status", &response)
	return
}

// GetCalendar calls Radarr's /api/v3/calendar endpoint. It returns all movies that will become available in the next 24 hours
func (rc RadarrClient) GetCalendar(ctx context.Context) (response []RadarrCalendarResponse, err error) {
	err = rc.Get(ctx, "/api/v3/calendar", &response)
	return
}

// GetQueuePage calls Radarr's /api/v3/queue/page=:pageNr endpoint. It returns one page of movies currently queued for download
func (rc RadarrClient) GetQueuePage(ctx context.Context, pageNr int) (response RadarrQueueResponse, err error) {
	err = rc.Get(ctx, fmt.Sprintf("/api/v3/queue?page=%d", pageNr), &response)
	return
}

// GetQueue calls Radarr's /api/v3/queue endpoint. It returns all movies currently queued for download
func (rc RadarrClient) GetQueue(ctx context.Context) (response RadarrQueueResponse, err error) {
	err = rc.Get(ctx, "/api/v3/queue", &response)

	for err == nil && len(response.Records) < response.TotalRecords {
		var tmp RadarrQueueResponse
		tmp, err = rc.GetQueuePage(ctx, response.Page+1)
		if err == nil {
			response.Page = tmp.Page
			response.Records = append(response.Records, tmp.Records...)
		}
	}

	return
}

// GetMovies calls Radarr's /api/v3/movie endpoint. It returns all movies added to Radarr
func (rc RadarrClient) GetMovies(ctx context.Context) (response []RadarrMovieResponse, err error) {
	err = rc.Get(ctx, "/api/v3/movie", &response)
	return
}

// GetMovieByID calls Radar's "/api/v3/movie/:movieID endpoint. It returns details for the specified movieID
func (rc RadarrClient) GetMovieByID(ctx context.Context, movieID int) (response RadarrMovieResponse, err error) {
	err = rc.Get(ctx, fmt.Sprintf("/api/v3/movie/%d", movieID), &response)
	return
}
