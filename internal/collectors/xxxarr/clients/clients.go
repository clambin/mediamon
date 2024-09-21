// Package clients collects mediamon metrics, using the openapi-generated API.
package clients

import (
	"context"
	"fmt"
	"net/http"
)

type QueuedItem struct {
	Name            string
	TotalBytes      int64
	DownloadedBytes int64
}

type Library struct {
	Monitored   int
	Unmonitored int
}

func WithToken(token string) func(ctx context.Context, req *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		if token == "" {
			return fmt.Errorf("no token provided")
		}
		req.Header.Set("X-Api-Key", token)
		return nil
	}
}
