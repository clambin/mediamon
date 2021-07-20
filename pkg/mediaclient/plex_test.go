package mediaclient_test

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
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

	client := &mediaclient.PlexClient{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	version, err := client.GetVersion(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "SomeVersion", version)

	version, err = client.GetVersion(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "SomeVersion", version)
}

func TestPlexClient_GetStats(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()

	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	client := &mediaclient.PlexClient{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	users, modes, transcoding, speed, err := client.GetSessions(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 3)
	assert.Len(t, modes, 3)
	assert.Equal(t, 2, transcoding)
	assert.Equal(t, 3.1, speed)

	// User count
	for _, testCase := range []struct {
		user  string
		ok    bool
		value int
	}{
		{"foo", true, 1},
		{"bar", true, 1},
		{"snafu", true, 2},
		{"ufans", false, 0},
	} {
		userCount, ok := users[testCase.user]
		assert.Equal(t, testCase.ok, ok)
		if testCase.ok {
			assert.Equal(t, testCase.value, userCount, testCase.user)
		}
	}
	// Mode count
	for _, testCase := range []struct {
		mode  string
		ok    bool
		value int
	}{
		{"direct", true, 1},
		{"copy", true, 1},
		{"transcode", true, 2},
		{"snafu", false, 0},
	} {
		modeCount, ok := modes[testCase.mode]
		assert.Equal(t, testCase.ok, ok)
		if testCase.ok {
			assert.Equal(t, testCase.value, modeCount, testCase.mode)
		}
	}
}

func TestPlexClient_Authentication(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	client := &mediaclient.PlexClient{
		Client:   &http.Client{},
		URL:      "",
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "badpassword",
	}

	_, err := client.GetVersion(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "403 Forbidden", err.Error())
}

func TestPlexClient_WithMetrics(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()
	testServer := httptest.NewServer(http.HandlerFunc(plexHandler))
	defer testServer.Close()

	requestDuration := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "plex_request_duration_seconds",
		Help: "Duration of API requests.",
	}, []string{"application", "request"})

	client := &mediaclient.PlexClient{
		Client:   &http.Client{},
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
		Options: mediaclient.PlexOpts{
			PrometheusSummary: requestDuration,
		},
	}

	_, err := client.GetVersion(context.Background())
	assert.NoError(t, err)

	// validate that the metrics were recorded
	ch := make(chan prometheus.Metric)
	go requestDuration.Collect(ch)

	desc := <-ch
	var m io_prometheus_client.Metric
	err = desc.Write(&m)
	assert.NoError(t, err)
	// TODO: why isn't this 2 (one for auth, one for API call)?
	assert.Equal(t, uint64(1), *m.Summary.SampleCount)
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
			{ "User": { "title": "foo" } },
			{ "User": { "title": "bar" },  "TranscodeSession": { "throttled": false, "speed": "3.1", "videoDecision": "copy" } },
			{ "User": { "title": "snafu" }, "TranscodeSession": { "throttled": true, "speed": "3.1", "videoDecision": "transcode" } },
			{ "User": { "title": "snafu" }, "TranscodeSession": { "throttled": true, "speed": "3.1", "videoDecision": "transcode" } }
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
