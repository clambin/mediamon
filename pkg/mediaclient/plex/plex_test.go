package plex_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexClient_GetIdentity(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()

	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	client := &plex.Client{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	identity, err := client.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "SomeVersion", identity.MediaContainer.Version)
}

func TestPlexClient_GetStats(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()

	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	client := &plex.Client{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	sessions, err := client.GetSessions(context.Background())
	require.NoError(t, err)

	titles := []string{"pilot", "movie 1", "movie 2", "movie 3"}
	locations := []string{"lan", "wan", "lan", "lan"}
	require.Len(t, sessions.MediaContainer.Metadata, len(titles))

	for index, title := range titles {
		assert.Equal(t, title, sessions.MediaContainer.Metadata[index].Title)
		assert.Equal(t, "Plex Web", sessions.MediaContainer.Metadata[index].Player.Product)
		assert.Equal(t, locations[index], sessions.MediaContainer.Metadata[index].Session.Location)

		if sessions.MediaContainer.Metadata[index].TranscodeSession.VideoDecision == "transcode" {
			assert.NotZero(t, sessions.MediaContainer.Metadata[index].TranscodeSession.Speed)
		} else {
			assert.Zero(t, sessions.MediaContainer.Metadata[index].TranscodeSession.Speed)
		}
	}
}

func TestPlexClient_Authentication(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	client := &plex.Client{
		Client:   &http.Client{},
		URL:      "",
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "badpassword",
	}

	_, err := client.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403 Forbidden")
}

func TestPlexClient_Authentication_Failure(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	authServer.Close()

	client := &plex.Client{
		Client:   &http.Client{},
		URL:      "",
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "badpassword",
	}

	_, err := client.GetIdentity(context.Background())
	require.Error(t, err)
	assert.True(t, errors.Is(err, unix.ECONNREFUSED))
}

func TestPlexClient_WithMetrics(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))

	latencyMetric := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "plex_request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	errorMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "plex_request_errors",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	client := &plex.Client{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
		Options: plex.Options{
			PrometheusMetrics: metrics.APIClientMetrics{
				Latency: latencyMetric,
				Errors:  errorMetric,
			},
		},
	}

	_, err := client.GetIdentity(context.Background())
	require.NoError(t, err)

	// validate that the metrics were recorded
	ch := make(chan prometheus.Metric)
	go latencyMetric.Collect(ch)

	expected := map[string]uint64{
		"auth":      1,
		"/identity": 1,
	}
	for i := 0; i < len(expected); i++ {
		desc := <-ch
		assert.Equal(t, "plex", metrics.MetricLabel(desc, "application"))
		value, ok := expected[metrics.MetricLabel(desc, "request")]
		require.True(t, ok)
		assert.Equal(t, value, metrics.MetricValue(desc).GetSummary().GetSampleCount())
	}

	// shut down the server
	testServer.Close()

	_, err = client.GetIdentity(context.Background())
	require.Error(t, err)

	ch = make(chan prometheus.Metric)
	go errorMetric.Collect(ch)

	expected2 := map[string]float64{
		"auth":      0,
		"/identity": 1,
	}
	for i := 0; i < len(expected2); i++ {
		desc := <-ch
		assert.Equal(t, "plex", metrics.MetricLabel(desc, "application"))
		value, ok := expected2[metrics.MetricLabel(desc, "request")]
		require.True(t, ok)
		assert.Equal(t, value, metrics.MetricValue(desc).GetCounter().GetValue())
	}
}

func TestClient_Failures(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()
	testServer := httptest.NewServer(http.HandlerFunc(plexBadHandler))

	client := &plex.Client{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	_, err := client.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 Internal Server Error", err.Error())

	testServer.Close()
	_, err = client.GetIdentity(context.Background())
	require.Error(t, err)
	assert.True(t, errors.Is(err, unix.ECONNREFUSED))
}

// Server handlers

func plexAuthHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		_ = req.Body.Close()
	}()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if string(body) != `user%5Blogin%5D=user@example.com&user%5Bpassword%5D=somepassword` {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(authResponse))
}

func plexHandler(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("X-Plex-Token")
	if token != "some_token" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	response, ok := plexResponses[req.URL.Path]

	if ok == false {
		http.Error(w, "endpoint not implemented: "+req.URL.Path, http.StatusNotFound)
	} else {
		_, _ = w.Write([]byte(response))
	}
}

func plexBadHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "server's having a hard day", http.StatusInternalServerError)
}

var plexResponses = map[string]string{
	"/identity": `{ "MediaContainer": {
    	"size": 0,
    	"claimed": true,
    	"machineIdentifier": "SomeUUID",
    	"version": "SomeVersion"
  	}}`,
	"/status/sessions": `{ "MediaContainer": {
		"size": 2,
		"Metadata": [
			{ "User": { "title": "foo" },   "Player": { "product": "Plex Web" }, "Session": { "location": "lan"}, "grandparentTitle": "series", "parentTitle": "season 1", "title": "pilot", "type": "episode"},
			{ "User": { "title": "bar" },   "Player": { "product": "Plex Web" }, "Session": { "location": "wan"}, "TranscodeSession": { "throttled": false, "videoDecision": "copy" }, "title": "movie 1" },
			{ "User": { "title": "snafu" }, "Player": { "product": "Plex Web" }, "Session": { "location": "lan"}, "TranscodeSession": { "throttled": true, "speed": 3.1, "videoDecision": "transcode" }, "title": "movie 2" },
			{ "User": { "title": "snafu" }, "Player": { "product": "Plex Web" }, "Session": { "location": "lan"}, "TranscodeSession": { "throttled": true, "speed": 4.1, "videoDecision": "transcode" }, "title": "movie 3" }
		]
	}}`,
}

const (
	authResponse = `<?xml version="1.0" encoding="UTF-8"?>
<user email="user@example.com" id="1" uuid="1" username="user" authenticationToken="some_token" authToken="some_token">
  <subscription active="0" status="Inactive" plan=""></subscription>
  <entitlements all="0"></entitlements>
  <profile_settings/>
  <providers></providers>
  <services></services>
  <username>user</username>
  <email>user@example.com</email>
  <joined-at type="datetime">2000-01-01 00:00:00 UTC</joined-at>
  <authentication-token>some_token</authentication-token>
</user>`
)
