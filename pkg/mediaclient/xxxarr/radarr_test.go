package xxxarr_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var radarrResponses = Responses{
	`/api/v3/system/status`: `{ "version": "1.2.3.4444" }`,
	`/api/v3/queue`:         `{ "page": 1, "pageSize": 1, "totalRecords": 2, "records": [ { "title": "foo" } ] }`,
	`/api/v3/queue?page=2`:  `{ "page": 2, "pageSize": 1, "totalRecords": 2, "records": [ { "title": "bar" } ] }`,
	`/api/v3/movie`:         `[ { "monitored": true }, { "monitored": false }, { "monitored": true } ]`,
	`/api/v3/movie/11`:      `{ "title": "foo", "monitored": true }`,
	`/api/v3/calendar`: `[
  { "title": "Bar", "hasFile": false, "monitored": true },
  { "hasFile": true, "monitored": true },
  { "hasFile": false, "monitored": false }
]`,
}

func TestNewRadarrClient_GetURL(t *testing.T) {
	c := xxxarr.NewRadarrClient("1234", "foo", client.Options{})
	assert.Equal(t, "foo", c.GetURL())
}

func TestRadarrClient_SystemStatus(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient("1234", s.server.URL, client.Options{})

	response, err := c.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1.2.3.4444", response.Version)

}

func TestRadarrClient_GetCalendar(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient("1234", s.server.URL, client.Options{})

	_, err := c.GetCalendar(context.Background())
	require.NoError(t, err)

}

func TestRadarrClient_GetQueuePage(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient("1234", s.server.URL, client.Options{})
	queue, err := c.GetQueuePage(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, queue.Records, 1)
	assert.Equal(t, "bar", queue.Records[0].Title)
}

func TestRadarrClient_GetQueue(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient("1234", s.server.URL, client.Options{})
	queue, err := c.GetQueue(context.Background())
	require.NoError(t, err)
	require.Len(t, queue.Records, 2)
	assert.Equal(t, "foo", queue.Records[0].Title)
	assert.Equal(t, "bar", queue.Records[1].Title)
}

func TestRadarrClient_GetMovies(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient("1234", s.server.URL, client.Options{})
	movies, err := c.GetMovies(context.Background())
	require.NoError(t, err)
	require.Len(t, movies, 3)
}

func TestRadarrClient_GetMovieByID(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient("1234", s.server.URL, client.Options{})
	movie, err := c.GetMovieByID(context.Background(), 11)
	require.NoError(t, err)
	assert.Equal(t, "foo", movie.Title)
}
