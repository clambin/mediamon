package clients

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/radarr"
	"net/http"
	"time"
)

type RadarrClient interface {
	GetApiV3SystemStatusWithResponse(ctx context.Context, reqEditors ...radarr.RequestEditorFn) (*radarr.GetApiV3SystemStatusResponse, error)
	GetApiV3HealthWithResponse(ctx context.Context, reqEditors ...radarr.RequestEditorFn) (*radarr.GetApiV3HealthResponse, error)
	GetApiV3CalendarWithResponse(ctx context.Context, params *radarr.GetApiV3CalendarParams, reqEditors ...radarr.RequestEditorFn) (*radarr.GetApiV3CalendarResponse, error)
	GetApiV3QueueWithResponse(ctx context.Context, params *radarr.GetApiV3QueueParams, reqEditors ...radarr.RequestEditorFn) (*radarr.GetApiV3QueueResponse, error)
	GetApiV3MovieWithResponse(ctx context.Context, params *radarr.GetApiV3MovieParams, reqEditors ...radarr.RequestEditorFn) (*radarr.GetApiV3MovieResponse, error)
}

type Radarr struct {
	Client RadarrClient
}

func NewRadarrClient(url, token string, httpClient *http.Client) (*Radarr, error) {
	var r Radarr
	var err error
	r.Client, err = radarr.NewClientWithResponses(url, radarr.WithRequestEditorFn(WithToken(token)), radarr.WithHTTPClient(httpClient))
	return &r, err
}

func (r Radarr) GetVersion(ctx context.Context) (string, error) {
	resp, err := r.Client.GetApiV3SystemStatusWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("GetApiV3SystemStatusWithResponse: %w", err)
	}
	return *resp.JSON200.Version, err
}

func (r Radarr) GetHealth(ctx context.Context) (map[string]int, error) {
	resp, err := r.Client.GetApiV3HealthWithResponse(ctx)
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

func (r Radarr) GetCalendar(ctx context.Context, days int) ([]string, error) {
	from := time.Now()
	to := from.AddDate(0, 0, days)
	params := radarr.GetApiV3CalendarParams{
		Start: &from,
		End:   &to,
	}
	resp, err := r.Client.GetApiV3CalendarWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("GetApiV3CalendarWithResponse: %w", err)
	}
	calendar := make([]string, len(*resp.JSON200))
	for i, movie := range *resp.JSON200 {
		calendar[i] = *movie.Title
	}
	return calendar, err
}

func (r Radarr) GetQueue(ctx context.Context) ([]QueuedItem, error) {
	var page int32
	pageSize := int32(100)

	var entries []QueuedItem
	for {
		params := radarr.GetApiV3QueueParams{
			Page:     &page,
			PageSize: &pageSize,
		}
		resp, err := r.Client.GetApiV3QueueWithResponse(ctx, &params)
		if err != nil {
			return nil, fmt.Errorf("GetApiV3QueueWithResponse: %w", err)
		}
		for _, record := range *resp.JSON200.Records {
			entries = append(entries, QueuedItem{
				Name:            *record.Title,
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

func (r Radarr) GetLibrary(ctx context.Context) (Library, error) {
	resp, err := r.Client.GetApiV3MovieWithResponse(ctx, &radarr.GetApiV3MovieParams{})
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
