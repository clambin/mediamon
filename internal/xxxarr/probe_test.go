package xxxarr_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/clambin/httpstub"
	"github.com/stretchr/testify/assert"

	"mediamon/internal/xxxarr"
	"mediamon/pkg/metrics"
)

func TestProbe_InvalidProbe(t *testing.T) {
	assert.Panics(t, func() { xxxarr.NewProbe("", "", "invalid") })
}

func TestFailingProbe(t *testing.T) {
	probe := xxxarr.NewProbe("", "1234", "sonarr")
	probe.Client.Client = httpstub.NewTestClient(httpstub.Failing)

	assert.NotNil(t, probe)
	assert.NotPanics(t, func() { probe.Run() })
}

func TestProbe_Run(t *testing.T) {
	for _, application := range []string{"sonarr", "radarr"} {
		probe := xxxarr.NewProbe("", "1234", application)
		probe.Client.Client = httpstub.NewTestClient(loopback)

		probe.Run()

		value, _ := metrics.LoadValue("mediaserver_server_info", application, "1.2.3.4444")
		assert.Equal(t, float64(1), value)
		count, _ := metrics.LoadValue("mediaserver_calendar_count", application)
		assert.Equal(t, float64(1), count)
		count, _ = metrics.LoadValue("mediaserver_queued_count", application)
		assert.Equal(t, float64(2), count)
		count, _ = metrics.LoadValue("mediaserver_monitored_count", application)
		assert.Equal(t, float64(2), count)
		count, _ = metrics.LoadValue("mediaserver_unmonitored_count", application)
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
