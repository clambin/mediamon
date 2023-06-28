package xxxarr_test

import (
	"context"
	"errors"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/xxxarr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var radarrResponses = Responses{
	`/api/v3/system/status`: xxxarr.RadarrSystemStatusResponse{Version: "1.2.3.4444"},
	`/api/v3/queue`:         xxxarr.RadarrQueueResponse{Page: 1, PageSize: 1, TotalRecords: 2, Records: []xxxarr.RadarrQueueResponseRecord{{Title: "foo"}}},
	`/api/v3/queue?page=2`:  xxxarr.RadarrQueueResponse{Page: 2, PageSize: 1, TotalRecords: 2, Records: []xxxarr.RadarrQueueResponseRecord{{Title: "bar"}}},
	`/api/v3/movie`:         []xxxarr.RadarrMovieResponse{{Monitored: true}, {Monitored: false}, {Monitored: true}},
	`/api/v3/movie/11`:      xxxarr.RadarrMovieResponse{Title: "foo", Monitored: true},
	`/api/v3/calendar`: []xxxarr.RadarrCalendarResponse{
		{Title: "foo", HasFile: false, Monitored: true},
		{Title: "bar", HasFile: true, Monitored: true},
		{Title: "snafu", HasFile: false, Monitored: false},
	},
}

func TestNewRadarrClient_GetURL(t *testing.T) {
	c := xxxarr.NewRadarrClient("foo", "1234", nil)
	assert.Equal(t, "foo", c.GetURL())
}

func TestRadarrClient_SystemStatus(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "1234", nil)
	response, err := c.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1.2.3.4444", response.Version)
}

func TestRadarrClient_GetCalendar(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "1234", nil)
	_, err := c.GetCalendar(context.Background())
	require.NoError(t, err)

}

func TestRadarrClient_GetQueuePage(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "1234", nil)
	queue, err := c.GetQueuePage(context.Background(), 2)
	require.NoError(t, err)
	require.Len(t, queue.Records, 1)
	assert.Equal(t, "bar", queue.Records[0].Title)
}

func TestRadarrClient_GetQueue(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "1234", nil)
	queue, err := c.GetQueue(context.Background())
	require.NoError(t, err)
	require.Len(t, queue.Records, 2)
	assert.Equal(t, "foo", queue.Records[0].Title)
	assert.Equal(t, "bar", queue.Records[1].Title)
}

func TestRadarrClient_GetMovies(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "1234", nil)
	movies, err := c.GetMovies(context.Background())
	require.NoError(t, err)
	require.Len(t, movies, 3)
}

func TestRadarrClient_GetMovieByID(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "1234", nil)
	movie, err := c.GetMovieByID(context.Background(), 11)
	require.NoError(t, err)
	assert.Equal(t, "foo", movie.Title)
}

func TestRadarrClient_BadKey(t *testing.T) {
	s := NewTestServer(radarrResponses, "1234")
	defer s.server.Close()

	c := xxxarr.NewRadarrClient(s.server.URL, "4321", nil)
	_, err := c.GetSystemStatus(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403 Forbidden")
}

func TestRadarrClient_BadOutput(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("bad output"))
	}))
	defer s.Close()

	c := xxxarr.NewRadarrClient(s.URL, "1234", nil)
	_, err := c.GetHealth(context.Background())
	assert.Error(t, err)
	var err2 *xxxarr.ErrInvalidJSON
	assert.True(t, errors.As(err, &err2))
	assert.Equal(t, "bad output", string(err2.Body))
}
