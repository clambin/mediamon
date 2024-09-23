package client

import (
	"context"
	"fmt"
	"github.com/clambin/mediaclients/prowlarr"
	"github.com/clambin/mediamon/v2/internal/collectors/xxxarr/clients"
	"net/http"
)

type Prowlarr struct {
	Client ProwlarrClient
}

type ProwlarrClient interface {
	GetApiV1IndexerstatsWithResponse(ctx context.Context, params *prowlarr.GetApiV1IndexerstatsParams, reqEditors ...prowlarr.RequestEditorFn) (*prowlarr.GetApiV1IndexerstatsResponse, error)
}

func NewProwlarrClient(url, token string, httpClient *http.Client) (*Prowlarr, error) {
	c, err := prowlarr.NewClientWithResponses(url, prowlarr.WithRequestEditorFn(clients.WithToken(token)), prowlarr.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("NewClientWithResponses: %w", err)
	}
	return &Prowlarr{Client: c}, nil
}

func (p Prowlarr) GetIndexStats(ctx context.Context) (*prowlarr.IndexerStatsResource, error) {
	rsp, err := p.Client.GetApiV1IndexerstatsWithResponse(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAPIV1Indexerstats: %w", err)
	}
	return rsp.JSON200, nil
}
