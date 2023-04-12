package main

import (
	"context"
	"fmt"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/http"
	"sort"
)

func views(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	serverTokens, err := getServerTokens(ctx)
	if err != nil {
		slog.Error("failed to get tokens", err)
		return
	}

	totalViewCount, err := getViewCount(ctx, serverTokens)
	if err != nil {
		slog.Error("failed to get view count", err)
		return
	}

	flattened := totalViewCount.flatten()
	sort.Slice(flattened, func(i, j int) bool {
		return flattened[i].Views < flattened[j].Views
	})

	for _, entry := range flattened {
		//if entry.Views == 0 {
		fmt.Printf("%s: %s - %d\n", entry.Library, entry.Title, entry.Views)
		//}
	}
}

type viewCount map[guid]viewCountEntry

func (vc *viewCount) merge(s viewCount) {
	for id, entry := range s {
		if existing, ok := (*vc)[id]; ok {
			entry.Views += existing.Views
		}
		(*vc)[id] = entry
	}
}

func (vc *viewCount) flatten() []viewCountEntry {
	viewCounts := make([]viewCountEntry, 0, len(*vc))
	for _, entry := range *vc {
		viewCounts = append(viewCounts, entry)
	}
	return viewCounts
}

type guid string
type viewCountEntry struct {
	Library string
	Title   string
	Views   int
}

func getViewCount(ctx context.Context, authTokens []string) (viewCount, error) {
	totalViewCount := make(viewCount)

	for _, token := range authTokens {
		vc, err := getUserViewCount(ctx, token)
		if err != nil {
			return totalViewCount, err
		}
		totalViewCount.merge(vc)
	}
	return totalViewCount, nil
}

func getUserViewCount(ctx context.Context, authToken string) (viewCount, error) {
	result := make(viewCount)
	c := plex.Client{AuthToken: authToken, URL: viper.GetString("url")}
	libraries, err := c.GetLibraries(ctx)
	if err != nil {
		return result, err
	}
	for _, library := range libraries.Directory {
		var libraryViewCount viewCount
		switch library.Type {
		case "movie":
			libraryViewCount, err = getMovieViewCount(ctx, &c, library)
		case "show":
			libraryViewCount, err = getShowViewCount(ctx, &c, library)
		}

		result.merge(libraryViewCount)
	}
	return result, err
}

func getMovieViewCount(ctx context.Context, client *plex.Client, library plex.LibrariesDirectory) (viewCount, error) {
	result := make(viewCount)
	entries, err := client.GetMovieLibrary(ctx, library.Key)
	if err != nil {
		return result, err
	}
	for _, entry := range entries.Metadata {
		existing, ok := result[guid(entry.Guid)]
		if !ok {
			existing = viewCountEntry{
				Library: library.Title,
				Title:   entry.Title,
			}
		}
		existing.Views += entry.ViewCount
		result[guid(entry.Guid)] = existing
	}
	return result, nil
}

func getShowViewCount(ctx context.Context, client *plex.Client, library plex.LibrariesDirectory) (viewCount, error) {
	result := make(viewCount)
	entries, err := client.GetShowLibrary(ctx, library.Key)
	if err != nil {
		return result, err
	}
	for _, entry := range entries.Metadata {
		existing, ok := result[guid(entry.Guid)]
		if !ok {
			existing = viewCountEntry{
				Library: library.Title,
				Title:   entry.Title,
			}
		}
		existing.Views += entry.ViewCount
		result[guid(entry.Guid)] = existing
	}
	return result, nil
}

func getServerTokens(ctx context.Context) ([]string, error) {
	serverToken := viper.GetString("auth.serverToken")
	if serverToken == "" {
		return nil, fmt.Errorf("no server token configured")
	}

	c := plex.Client{HTTPClient: http.DefaultClient, AuthToken: serverToken}

	tokens, err := c.GetAccessTokens(ctx)
	if err != nil {
		return nil, err
	}

	serverTokens := []string{serverToken}
	for _, token := range tokens {
		if token.Type == "server" {
			slog.Debug("token found", "user", token.Invited.Title)
			serverTokens = append(serverTokens, token.Token)
		}
	}

	return serverTokens, nil
}
