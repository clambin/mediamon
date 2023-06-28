package xxxarr

import (
	"context"
	"fmt"
	"net/http"
)

// RadarrClient calls Radarr endpoints
type RadarrClient struct {
	Client *http.Client
	URL    string
	APIKey string
}

const radarrAPIPrefix = "/api/v3"

func (c RadarrClient) GetURL() string {
	return c.URL
}

// GetSystemStatus calls Radarr's  /api/v3/system/status endpoint. It returns the system status of the Radarr instance
func (c RadarrClient) GetSystemStatus(ctx context.Context) (response RadarrSystemStatusResponse, err error) {
	return call[RadarrSystemStatusResponse](ctx, c.Client, c.URL+radarrAPIPrefix+"/system/status", c.APIKey)
}

// GetHealth calls Radarr's /api/v3/health endpoint. It returns the health of the Radarr instance
func (c RadarrClient) GetHealth(ctx context.Context) (response []RadarrHealthResponse, err error) {
	return call[[]RadarrHealthResponse](ctx, c.Client, c.URL+radarrAPIPrefix+"/health", c.APIKey)
}

// GetCalendar calls Radarr's /api/v3/calendar endpoint. It returns all movies that will become available in the next 24 hours
func (c RadarrClient) GetCalendar(ctx context.Context) (response []RadarrCalendarResponse, err error) {
	return call[[]RadarrCalendarResponse](ctx, c.Client, c.URL+radarrAPIPrefix+"/calendar", c.APIKey)
}

// GetQueuePage calls Radarr's /api/v3/queue/page=:pageNr endpoint. It returns one page of movies currently queued for download
func (c RadarrClient) GetQueuePage(ctx context.Context, pageNr int) (response RadarrQueueResponse, err error) {
	return call[RadarrQueueResponse](ctx, c.Client, fmt.Sprintf(c.URL+radarrAPIPrefix+"/queue?page=%d", pageNr), c.APIKey)
}

// GetQueue calls Radarr's /api/v3/queue endpoint. It returns all movies currently queued for download
func (c RadarrClient) GetQueue(ctx context.Context) (response RadarrQueueResponse, err error) {
	response, err = call[RadarrQueueResponse](ctx, c.Client, c.URL+radarrAPIPrefix+"/queue", c.APIKey)

	for err == nil && len(response.Records) < response.TotalRecords {
		var tmp RadarrQueueResponse
		tmp, err = c.GetQueuePage(ctx, response.Page+1)
		if err == nil {
			response.Page = tmp.Page
			response.Records = append(response.Records, tmp.Records...)
		}
	}

	return
}

// GetMovies calls Radarr's /api/v3/movie endpoint. It returns all movies added to Radarr
func (c RadarrClient) GetMovies(ctx context.Context) (response []RadarrMovieResponse, err error) {
	return call[[]RadarrMovieResponse](ctx, c.Client, c.URL+radarrAPIPrefix+"/movie", c.APIKey)
}

// GetMovieByID calls Radar's "/api/v3/movie/:movieID endpoint. It returns details for the specified movieID
func (c RadarrClient) GetMovieByID(ctx context.Context, movieID int) (response RadarrMovieResponse, err error) {
	return call[RadarrMovieResponse](ctx, c.Client, fmt.Sprintf(c.URL+radarrAPIPrefix+"/movie/%d", movieID), c.APIKey)
}
