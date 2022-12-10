package xxxarr_test

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var sonarrResponses = Responses{
	`/api/v3/system/status`: `{ "version": "1.2.3.4444" }`,
	`/api/v3/series/11`:     `{ "title": "Foo" }`,
	`/api/v3/queue`:         `{ "page": 1, "pageSize": 1, "totalRecords": 2, "records": [ { "title": "foo" } ] }`,
	`/api/v3/queue?page=2`:  `{ "page": 2, "pageSize": 1, "totalRecords": 2, "records": [ { "title": "bar" } ] }`,
	`/api/v3/series`:        `[ { "title": "Foo", "monitored": true }, { "monitored": false } ]`,
	`/api/v3/episode/11`:    `{ "title": "Foo", "seasonNumber": 1, "episodeNumber": 2, "series": { "title": "Bar" } }`,
	`/api/v3/calendar`: `[
  { "seriesId": 11, "title": "Bar", "seasonNumber": 2, "episodeNumber": 9, "hasFile": false, "monitored": true },
  { "hasFile": true, "monitored": true },
  { "hasFile": false, "monitored": false }
]`,
}

func TestNewSonarrClient_GetURL(t *testing.T) {
	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: "foo"}
	assert.Equal(t, "foo", c.GetURL())
}

func TestSonarrClient_GetSystemStatus(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	response, err := c.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1.2.3.4444", response.Version)
}

func TestSonarrClient_GetCalendar(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	calendar, err := c.GetCalendar(context.Background())
	require.NoError(t, err)
	require.Len(t, calendar, 3)
	assert.Equal(t, "Bar", calendar[0].Title)
	assert.True(t, calendar[1].HasFile)
	assert.False(t, calendar[2].Monitored)
}

func TestSonarrClient_GetQueue(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	queue, err := c.GetQueue(context.Background())
	require.NoError(t, err)
	require.Len(t, queue.Records, 2)
	assert.Equal(t, "foo", queue.Records[0].Title)
}

func TestSonarrClient_GetQueuePage(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	queue, err := c.GetQueuePage(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, queue.Records, 1)
	assert.Equal(t, "bar", queue.Records[0].Title)
}

func TestSonarrClient_GetSeries(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	series, err := c.GetSeries(context.Background())
	require.NoError(t, err)
	require.Len(t, series, 2)
	assert.Equal(t, "Foo", series[0].Title)
	assert.True(t, series[0].Monitored)
	assert.False(t, series[1].Monitored)
}

func TestSonarrClient_GetSeriesByID(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	series, err := c.GetSeriesByID(context.Background(), 11)
	require.NoError(t, err)
	assert.Equal(t, "Foo", series.Title)
}

func TestSonarrClient_GetEpisodeByID(t *testing.T) {
	s := NewTestServer(sonarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.SonarrClient{Client: http.DefaultClient, APIKey: "1234", URL: s.server.URL}
	episode, err := c.GetEpisodeByID(context.Background(), 11)
	require.NoError(t, err)
	assert.Equal(t, "Foo", episode.Title)
	assert.Equal(t, "Bar", episode.Series.Title)
	assert.Equal(t, 1, episode.SeasonNumber)
	assert.Equal(t, 2, episode.EpisodeNumber)
}
