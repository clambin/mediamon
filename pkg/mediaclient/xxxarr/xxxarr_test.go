package xxxarr_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/xxxarr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXXXArrClient_GetApplication(t *testing.T) {
	client := xxxarr.Client{
		Client:      &http.Client{},
		URL:         "",
		APIKey:      "",
		Application: "sonarr",
	}

	assert.Equal(t, "sonarr", client.GetApplication())
}

func TestXXXArrClient_GetURL(t *testing.T) {
	client := xxxarr.Client{
		Client:      &http.Client{},
		URL:         "https://localhost:4321",
		APIKey:      "",
		Application: "sonarr",
	}

	assert.Equal(t, "https://localhost:4321", client.GetURL())
}

func TestXXXArrClient_GetVersion(t *testing.T) {
	type fields struct {
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    "1.2.3.4444",
			wantErr: false,
		},
		{
			name: "sonarr",
			fields: fields{
				APIKey:      "",
				Application: "sonarr",
			},
			wantErr: true,
		},
		{
			name: "radarr",
			fields: fields{
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    "1.2.3.4444",
			wantErr: false,
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(xxxArrHandler))
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &xxxarr.Client{
				Client:      &http.Client{},
				URL:         testServer.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, err := client.GetVersion(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXXXArrClient_GetCalendar(t *testing.T) {
	type fields struct {
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    1,
			wantErr: false,
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(xxxArrHandler))
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &xxxarr.Client{
				Client:      &http.Client{},
				URL:         testServer.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, err := client.GetCalendar(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCalendar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCalendar() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXXXArrClient_GetQueue(t *testing.T) {
	type fields struct {
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    3,
			wantErr: false,
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(xxxArrHandler))
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &xxxarr.Client{
				Client:      &http.Client{},
				URL:         testServer.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, err := client.GetQueue(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetQueue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXXXArrClient_GetMonitored(t *testing.T) {
	type fields struct {
		APIKey      string
		Application string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		want1   int
		wantErr bool
	}{
		{
			name: "sonarr",
			fields: fields{
				APIKey:      "1234",
				Application: "sonarr",
			},
			want:    2,
			want1:   1,
			wantErr: false,
		},
		{
			name: "radarr",
			fields: fields{
				APIKey:      "1234",
				Application: "radarr",
			},
			want:    2,
			want1:   1,
			wantErr: false,
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(xxxArrHandler))
	defer testServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &xxxarr.Client{
				Client:      &http.Client{},
				URL:         testServer.URL,
				APIKey:      tt.fields.APIKey,
				Application: tt.fields.Application,
			}
			got, got1, err := client.GetMonitored(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMonitored() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMonitored() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetMonitored() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestXXXArrClient_GetMonitored_Panic(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(xxxArrHandler))
	defer testServer.Close()

	client := &xxxarr.Client{
		Client:      &http.Client{},
		URL:         testServer.URL,
		APIKey:      "",
		Application: "invalid",
	}

	assert.Panics(t, func() { _, _, _ = client.GetMonitored(context.Background()) })
}

func TestXXXArrClient_ServerDown(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(xxxArrDownHandler))

	client := &xxxarr.Client{
		Client:      &http.Client{},
		URL:         testServer.URL,
		APIKey:      "",
		Application: "sonarr",
	}

	_, err := client.GetVersion(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 Internal Server Error", err.Error())

	testServer.Close()
	_, err = client.GetVersion(context.Background())
	require.Error(t, err)
	assert.True(t, errors.Is(err, unix.ECONNREFUSED))
}

func TestXXXArrClient_WithMetrics(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(xxxArrHandler))
	defer testServer.Close()

	latencyMetric := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "xxxarr_request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	errorMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "xxxarr_request_errors",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	client := &xxxarr.Client{
		Client:      &http.Client{},
		URL:         testServer.URL,
		APIKey:      "1234",
		Application: "sonarr",
		Options: xxxarr.Options{
			PrometheusMetrics: metrics.APIClientMetrics{
				Latency: latencyMetric,
				Errors:  errorMetric,
			},
		},
	}

	_, err := client.GetVersion(context.Background())
	require.NoError(t, err)

	// validate that a metric was recorded
	ch := make(chan prometheus.Metric)
	go latencyMetric.Collect(ch)

	desc := <-ch
	assert.Equal(t, uint64(1), metrics.MetricValue(desc).GetSummary().GetSampleCount())

	// shut down the server
	testServer.Close()

	_, err = client.GetVersion(context.Background())
	require.Error(t, err)

	ch = make(chan prometheus.Metric)
	go errorMetric.Collect(ch)

	desc = <-ch
	assert.Equal(t, 1.0, metrics.MetricValue(desc).GetCounter().GetValue())
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
	queued = `{ "totalRecords": 3 }`

	monitored = `[ { "monitored": true }, { "monitored": false }, { "monitored": true } ]`
)

// Server handler
func xxxArrHandler(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Api-Key") != "1234" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var response string
	switch req.URL.Path {
	case "/api/v3/system/status":
		response = systemStatus
	case "/api/v3/calendar":
		response = calendar
	case "/api/v3/queue":
		response = queued
	case "/api/v3/series":
		response = monitored
	case "/api/v3/movie":
		response = monitored
	default:
		http.Error(w, "endpoint not implemented", http.StatusNotFound)
		return
	}

	_, _ = w.Write([]byte(response))
}

func xxxArrDownHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "the software, it no workie", http.StatusInternalServerError)
}
