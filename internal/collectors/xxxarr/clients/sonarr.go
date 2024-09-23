package clients

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/sonarr"
	"net/http"
	"time"
)

type SonarrClient interface {
	GetApiV3SystemStatusWithResponse(ctx context.Context, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3SystemStatusResponse, error)
	GetApiV3HealthWithResponse(ctx context.Context, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3HealthResponse, error)
	GetApiV3CalendarWithResponse(ctx context.Context, params *sonarr.GetApiV3CalendarParams, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3CalendarResponse, error)
	GetApiV3EpisodeWithResponse(ctx context.Context, params *sonarr.GetApiV3EpisodeParams, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3EpisodeResponse, error)
	GetApiV3SeriesIdWithResponse(ctx context.Context, id int32, params *sonarr.GetApiV3SeriesIdParams, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3SeriesIdResponse, error)
	GetApiV3QueueWithResponse(ctx context.Context, params *sonarr.GetApiV3QueueParams, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3QueueResponse, error)
	GetApiV3SeriesWithResponse(ctx context.Context, params *sonarr.GetApiV3SeriesParams, reqEditors ...sonarr.RequestEditorFn) (*sonarr.GetApiV3SeriesResponse, error)
}

type Sonarr struct {
	Client SonarrClient
}

func NewSonarrClient(url, token string, httpClient *http.Client) (*Sonarr, error) {
	var s Sonarr
	var err error
	s.Client, err = sonarr.NewClientWithResponses(url, sonarr.WithRequestEditorFn(WithToken(token)), sonarr.WithHTTPClient(httpClient))
	return &s, err
}

func (s Sonarr) GetVersion(ctx context.Context) (string, error) {
	resp, err := s.Client.GetApiV3SystemStatusWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("GetApiV3SystemStatusWithResponse: %w", err)
	}
	return *resp.JSON200.Version, err
}

func (s Sonarr) GetHealth(ctx context.Context) (map[string]int, error) {
	resp, err := s.Client.GetApiV3HealthWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetApiV3HealthWithResponse: %w", err)
	}
	health := make(map[string]int, len(*resp.JSON200))
	for _, healthItem := range *resp.JSON200 {
		healthType := string(*healthItem.Type)
		value := health[healthType]
		health[healthType] = value + 1
	}
	return health, err
}

func (s Sonarr) GetCalendar(ctx context.Context, days int) ([]string, error) {
	from := time.Now()
	to := from.AddDate(0, 0, days)
	params := sonarr.GetApiV3CalendarParams{
		Start: &from,
		End:   &to,
	}
	resp, err := s.Client.GetApiV3CalendarWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("GetApiV3CalendarWithResponse: %w", err)
	}
	calendar := make([]string, len(*resp.JSON200))
	for i, episode := range *resp.JSON200 {
		name, err := s.getEpisodeNameFromEpisodeResource(ctx, episode)
		if err != nil {
			return nil, fmt.Errorf("getEpisodeNameFromEpisodeResource: %w", err)
		}
		calendar[i] = name
	}
	return calendar, err
}

func (s Sonarr) getEpisodeNameFromEpisodeResource(ctx context.Context, episode sonarr.EpisodeResource) (string, error) {
	seriesIdResponse, err := s.Client.GetApiV3SeriesIdWithResponse(ctx, *episode.SeriesId, &sonarr.GetApiV3SeriesIdParams{})
	if err != nil {
		return "", fmt.Errorf("GetApiV3SeriesIdWithResponse: %w", err)
	}
	return fmt.Sprintf("%s - S%02dE%02d - %s",
		*seriesIdResponse.JSON200.Title,
		*episode.SeasonNumber,
		*episode.EpisodeNumber,
		*episode.Title,
	), nil
}

func (s Sonarr) GetQueue(ctx context.Context) ([]QueuedItem, error) {
	var page int32
	pageSize := int32(100)
	var entries []QueuedItem
	for {
		params := sonarr.GetApiV3QueueParams{
			Page:     &page,
			PageSize: &pageSize,
		}
		resp, err := s.Client.GetApiV3QueueWithResponse(ctx, &params)
		if err != nil {
			return nil, fmt.Errorf("GetApiV3QueueWithResponse: %w", err)
		}
		for _, record := range *resp.JSON200.Records {
			name, err := s.getEpisodeNameFromQueueResource(ctx, record)
			if err != nil {
				return nil, fmt.Errorf("getEpisodeNameFromQueueResource: %w", err)
			}
			entries = append(entries, QueuedItem{
				Name:            name,
				TotalBytes:      int64(*record.Size),
				DownloadedBytes: int64(*record.Size - *record.Sizeleft),
			})
		}
		if len(entries) == int(*resp.JSON200.TotalRecords) {
			break
		}
		page++
	}
	return entries, nil
}

func (s Sonarr) getEpisodeNameFromQueueResource(ctx context.Context, episode sonarr.QueueResource) (string, error) {
	seriesIdResponse, err := s.Client.GetApiV3SeriesIdWithResponse(ctx, *episode.SeriesId, &sonarr.GetApiV3SeriesIdParams{})
	if err != nil {
		return "", fmt.Errorf("GetApiV3SeriesIdWithResponse: %w", err)
	}
	episodeDetails, err := s.Client.GetApiV3EpisodeWithResponse(ctx, &sonarr.GetApiV3EpisodeParams{EpisodeIds: &[]int32{*episode.EpisodeId}})
	if err != nil {
		return "", fmt.Errorf("GetApiV3EpisodeWithResponse: %w", err)
	}
	return fmt.Sprintf("%s - S%02dE%02d - %s",
		*seriesIdResponse.JSON200.Title,
		*episode.SeasonNumber,
		*(*episodeDetails.JSON200)[0].EpisodeNumber,
		*episode.Title,
	), nil
}

func (s Sonarr) GetLibrary(ctx context.Context) (Library, error) {
	resp, err := s.Client.GetApiV3SeriesWithResponse(ctx, &sonarr.GetApiV3SeriesParams{})
	if err != nil {
		return Library{}, fmt.Errorf("GetApiV3SeriesWithResponse: %w", err)
	}
	var library Library
	for _, entry := range *resp.JSON200 {
		if *entry.Monitored {
			library.Monitored++
		} else {
			library.Unmonitored++
		}
	}
	return library, err
}
