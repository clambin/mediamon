package xxxarr

import (
	"context"
	"fmt"
	"net/http"
)

// SonarrClient calls Sonarr endpoints
type SonarrClient struct {
	Client *http.Client
	URL    string
	APIKey string
}

const sonarrAPIPrefix = "/api/v3"

func (c SonarrClient) GetURL() string {
	return c.URL
}

// GetSystemStatus calls Sonarr's /api/v3/system/status endpoint. It returns the system status of the Sonarr instance
func (c SonarrClient) GetSystemStatus(ctx context.Context) (response SonarrSystemStatusResponse, err error) {
	return call[SonarrSystemStatusResponse](ctx, c.Client, c.URL+sonarrAPIPrefix+"/system/status", c.APIKey)
}

// GetHealth calls Sonarr's /api/v3/health endpoint. It returns the health of the Radarr instance
func (c SonarrClient) GetHealth(ctx context.Context) (response []SonarrHealthResponse, err error) {
	return call[[]SonarrHealthResponse](ctx, c.Client, c.URL+sonarrAPIPrefix+"/health", c.APIKey)
}

// GetCalendar calls Sonarr's /api/v3/calendar endpoint. It returns all episodes that will become available in the next 24 hours
func (c SonarrClient) GetCalendar(ctx context.Context) (response []SonarrCalendarResponse, err error) {
	return call[[]SonarrCalendarResponse](ctx, c.Client, c.URL+sonarrAPIPrefix+"/calendar", c.APIKey)
}

// GetQueuePage calls Sonarr's /api/v3/queue/page=:pageNr endpoint. It returns one page of episodes currently queued for download
func (c SonarrClient) GetQueuePage(ctx context.Context, pageNr int) (response SonarrQueueResponse, err error) {
	return call[SonarrQueueResponse](ctx, c.Client, fmt.Sprintf(c.URL+sonarrAPIPrefix+"/queue?page=%d", pageNr), c.APIKey)
}

// GetQueue calls Sonarr's /api/v3/queue endpoint. It returns all episodes currently queued for download
func (c SonarrClient) GetQueue(ctx context.Context) (response SonarrQueueResponse, err error) {
	response, err = call[SonarrQueueResponse](ctx, c.Client, c.URL+sonarrAPIPrefix+"/queue", c.APIKey)

	for err == nil && len(response.Records) < response.TotalRecords {
		var tmp SonarrQueueResponse
		tmp, err = c.GetQueuePage(ctx, response.Page+1)
		if err == nil {
			response.Page = tmp.Page
			response.Records = append(response.Records, tmp.Records...)
		}
	}

	return
}

// GetSeries calls Sonarr's /api/v3/series endpoint. It returns all series added to Sonarr
func (c SonarrClient) GetSeries(ctx context.Context) (response []SonarrSeriesResponse, err error) {
	return call[[]SonarrSeriesResponse](ctx, c.Client, c.URL+sonarrAPIPrefix+"/series", c.APIKey)
}

// GetSeriesByID calls Sonarr's /api/v3/series/:seriesID endpoint. It returns details for the specified seriesID
func (c SonarrClient) GetSeriesByID(ctx context.Context, seriesID int) (response SonarrSeriesResponse, err error) {
	return call[SonarrSeriesResponse](ctx, c.Client, fmt.Sprintf(c.URL+sonarrAPIPrefix+"/series/%d", seriesID), c.APIKey)
}

// GetEpisodeByID calls Sonarr's /api/v3/episode/:episodeID endpoint. It returns details for the specified responseID
func (c SonarrClient) GetEpisodeByID(ctx context.Context, episodeID int) (response SonarrEpisodeResponse, err error) {
	return call[SonarrEpisodeResponse](ctx, c.Client, fmt.Sprintf(c.URL+sonarrAPIPrefix+"/episode/%d", episodeID), c.APIKey)
}
