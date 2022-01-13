package plex_test

import (
	"context"
	"errors"
	"github.com/clambin/mediamon/pkg/mediaclient/metrics"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexClient_GetVersion(t *testing.T) {
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

	version, err := client.GetVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "SomeVersion", version)

	version, err = client.GetVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "SomeVersion", version)
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
	assert.Equal(t, []plex.Session{
		{Title: "", User: "foo", Local: true, Transcode: false, Throttled: false, Speed: 0},
		{Title: "", User: "bar", Local: false, Transcode: false, Throttled: false, Speed: 2.1},
		{Title: "", User: "snafu", Local: true, Transcode: true, Throttled: true, Speed: 3.1},
		{Title: "", User: "snafu", Local: true, Transcode: true, Throttled: true, Speed: 4.1},
	}, sessions)
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

	_, err := client.GetVersion(context.Background())
	require.Error(t, err)
	assert.Equal(t, "403 Forbidden", err.Error())
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
			PrometheusMetrics: metrics.PrometheusMetrics{
				Latency: latencyMetric,
				Errors:  errorMetric,
			},
		},
	}

	_, err := client.GetVersion(context.Background())
	require.NoError(t, err)

	// validate that the metrics were recorded
	ch := make(chan prometheus.Metric)
	go latencyMetric.Collect(ch)

	desc := <-ch
	var m io_prometheus_client.Metric
	err = desc.Write(&m)
	require.NoError(t, err)
	// TODO: why isn't this 2 (one for auth, one for API call)?
	assert.Equal(t, uint64(1), m.Summary.GetSampleCount())

	// shut down the server
	testServer.Close()

	_, err = client.GetVersion(context.Background())
	require.Error(t, err)

	ch = make(chan prometheus.Metric)
	go errorMetric.Collect(ch)

	desc = <-ch
	err = desc.Write(&m)
	require.NoError(t, err)
	// TODO: why isn't this 2 (one for auth, one for API call)?
	assert.Equal(t, float64(1), m.Counter.GetValue())
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

	_, err := client.GetVersion(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 Internal Server Error", err.Error())

	testServer.Close()
	_, err = client.GetVersion(context.Background())
	require.Error(t, err)
	assert.True(t, errors.Is(err, unix.ECONNREFUSED))
}

// Server handlers

func plexAuthHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		_ = req.Body.Close()
	}()
	body, err := ioutil.ReadAll(req.Body)

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
			{ "User": { "title": "foo" }, "Player": { "local": true }},
			{ "User": { "title": "bar" },  "Player": { "local": false }, "TranscodeSession": { "throttled": false, "speed": 2.1, "videoDecision": "copy" } },
			{ "User": { "title": "snafu" }, "Player": { "local": true }, "TranscodeSession": { "throttled": true, "speed": 3.1, "videoDecision": "transcode" } },
			{ "User": { "title": "snafu" }, "Player": { "local": true }, "TranscodeSession": { "throttled": true, "speed": 4.1, "videoDecision": "transcode" } }
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
