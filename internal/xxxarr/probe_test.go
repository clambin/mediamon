package xxxarr_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"mediamon/internal/metrics"
	"mediamon/internal/xxxarr"

	"mediamon/internal/httpstub"
)

func TestProbe_InvalidProbe(t *testing.T) {
	assert.Panics(t, func() { xxxarr.NewProbeWithHTTPClient(&http.Client{}, "", "", "invalid") })
}

func TestProbe_Run(t *testing.T) {
	for _, application := range []string{"sonarr", "radarr"} {
		probe := xxxarr.NewProbeWithHTTPClient(httpstub.NewTestClient(loopback), "http://example.com", "1234", application)

		log.SetLevel(log.DebugLevel)

		probe.Run()

		value, ok := metrics.LoadValue("version", application, "1.2.3.4444")
		assert.True(t, ok)
		assert.Equal(t, float64(1), value)

		count, ok := metrics.LoadValue("xxxarr_calendar", application)
		assert.True(t, ok)
		assert.Equal(t, float64(1), count)

		count, ok = metrics.LoadValue("xxxarr_queue", application)
		assert.True(t, ok)
		assert.Equal(t, float64(2), count)

		count, ok = metrics.LoadValue("xxxarr_monitored", application)
		assert.True(t, ok)
		assert.Equal(t, float64(2), count)

		count, ok = metrics.LoadValue("xxxarr_unmonitored", application)
		assert.True(t, ok)
		assert.Equal(t, float64(1), count)
	}
}

// Server loopback function

func loopback(req *http.Request) *http.Response {
	if req.Header.Get("X-Api-Key") != "1234" {
		return &http.Response{
			StatusCode: 409,
			Status:     "No/invalid Application Key",
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}
	}
	switch req.URL.Path {
	case "/api/system/status":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(systemStatus)),
		}
	case "/api/calendar":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(calendar)),
		}
	case "/api/queue":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(queued)),
		}
	case "/api/series":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(monitored)),
		}
	case "/api/movie":
		return &http.Response{
			StatusCode: 200,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(monitored)),
		}
	default:
		return &http.Response{
			StatusCode: 404,
			Header:     nil,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
		}
	}
}

// Responses

const (
	systemStatus = `{
  "version": "1.2.3.4444"
}`

	calendar = `[
  {
    "hasFile": false
  },
  {
    "hasFile": true
  }
]`
	queued = `[ {}, {} ]`

	monitored = `[ { "monitored": true }, { "monitored": false }, { "monitored": true } ]`
)
