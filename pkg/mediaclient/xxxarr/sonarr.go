package xxxarr

import (
	"context"
	"fmt"
	"github.com/clambin/go-metrics/client"
	"net/http"
)

// SonarrAPI contains all supported Sonarr APIs
//go:generate mockery --name SonarrAPI
type SonarrAPI interface {
	GetURL() string
	GetSystemStatus(ctx context.Context) (response SonarrSystemStatusResponse, err error)
	GetCalendar(ctx context.Context) (response []SonarrCalendarResponse, err error)
	GetQueuePage(ctx context.Context, pageNr int) (response SonarrQueueResponse, err error)
	GetQueue(ctx context.Context) (response SonarrQueueResponse, err error)
	GetSeries(ctx context.Context) (response []SonarrSeriesResponse, err error)
	GetSeriesByID(ctx context.Context, seriesID int) (response SonarrSeriesResponse, err error)
	GetEpisodeByID(ctx context.Context, episodeID int) (response SonarrEpisodeResponse, err error)
}

// SonarrClient calls Sonarr endpoints
type SonarrClient struct {
	APICaller
}

var _ SonarrAPI = &SonarrClient{}

// NewSonarrClient creates a new SonarrClient, using http.DefaultClient as http.InstrumentedClient
func NewSonarrClient(apiKey, url string, options client.Options) *SonarrClient {
	return NewSonarrClientWithCaller(apiKey, url, &client.InstrumentedClient{
		BaseClient:  client.BaseClient{HTTPClient: http.DefaultClient},
		Options:     options,
		Application: "sonarr",
	})
}

// NewSonarrClientWithCaller creates a new SonarrClient with a specified Caller
func NewSonarrClientWithCaller(apiKey, url string, caller client.Caller) *SonarrClient {
	return &SonarrClient{APICaller: &APIClient{
		Caller: caller,
		URL:    url,
		APIKey: apiKey,
	}}
}

// GetSystemStatus calls Sonarr's /api/v3/system/status endpoint. It returns the system status of the Sonarr instance
func (sc SonarrClient) GetSystemStatus(ctx context.Context) (response SonarrSystemStatusResponse, err error) {
	err = sc.Get(ctx, "/api/v3/system/status", &response)
	return
}

// GetCalendar calls Sonarr's /api/v3/calendar endpoint. It returns all episodes that will become available in the next 24 hours
func (sc SonarrClient) GetCalendar(ctx context.Context) (response []SonarrCalendarResponse, err error) {
	err = sc.Get(ctx, "/api/v3/calendar", &response)
	return
}

// GetQueuePage calls Sonarr's /api/v3/queue/page=:pageNr endpoint. It returns one page of episodes currently queued for download
func (sc SonarrClient) GetQueuePage(ctx context.Context, pageNr int) (response SonarrQueueResponse, err error) {
	err = sc.Get(ctx, fmt.Sprintf("/api/v3/queue?page=%d", pageNr), &response)
	return
}

// GetQueue calls Sonarr's /api/v3/queue endpoint. It returns all episodes currently queued for download
func (sc SonarrClient) GetQueue(ctx context.Context) (response SonarrQueueResponse, err error) {
	err = sc.Get(ctx, "/api/v3/queue", &response)

	for err == nil && len(response.Records) < response.TotalRecords {
		var tmp SonarrQueueResponse
		tmp, err = sc.GetQueuePage(ctx, response.Page+1)
		if err == nil {
			response.Page = tmp.Page
			response.Records = append(response.Records, tmp.Records...)
		}
	}

	return
}

// GetSeries calls Sonarr's /api/v3/series endpoint. It returns all series added to Sonarr
func (sc SonarrClient) GetSeries(ctx context.Context) (response []SonarrSeriesResponse, err error) {
	err = sc.Get(ctx, "/api/v3/series", &response)
	return
}

// GetSeriesByID calls Sonarr's /api/v3/series/:seriesID endpoint. It returns details for the specified seriesID
func (sc SonarrClient) GetSeriesByID(ctx context.Context, seriesID int) (response SonarrSeriesResponse, err error) {
	err = sc.Get(ctx, fmt.Sprintf("/api/v3/series/%d", seriesID), &response)
	return
}

// GetEpisodeByID calls Sonarr's /api/v3/episode/:episodeID endpoint. It returns details for the specified responseID
func (sc SonarrClient) GetEpisodeByID(ctx context.Context, episodeID int) (response SonarrEpisodeResponse, err error) {
	err = sc.Get(ctx, fmt.Sprintf("/api/v3/episode/%d", episodeID), &response)
	return
}
