package plex_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_WithRoundTripper(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	server := httptest.NewServer(authenticated("some_token", func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Dummy") != "foo" {
			http.Error(writer, "missing X-Dummy header", http.StatusBadRequest)
			return
		}
		plexHandler(writer, request)
	}))
	defer server.Client()

	c := plex.New("user@example.com", "somepassword", "", "", server.URL, &dummyRoundTripper{next: http.DefaultTransport})
	c.HTTPClient.Transport.(*plex.Authenticator).AuthURL = authServer.URL

	_, err := c.GetSessions(context.Background())
	assert.NoError(t, err)
}

var _ http.RoundTripper = &dummyRoundTripper{}

type dummyRoundTripper struct {
	next http.RoundTripper
}

func (d *dummyRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("X-Dummy", "foo")
	return d.next.RoundTrip(request)
}

func TestClient_Failures(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexBadHandler))

	c := plex.New("user@example.com", "somepassword", "", "", testServer.URL, nil)
	c.HTTPClient.Transport = http.DefaultTransport

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 "+http.StatusText(http.StatusInternalServerError), err.Error())

	testServer.Close()
	_, err = c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, unix.ECONNREFUSED)
}

func TestClient_Decode_Failure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(plexGarbageHandler))
	defer testServer.Close()

	c := plex.New("user@example.com", "somepassword", "", "", testServer.URL, nil)
	c.HTTPClient.Transport = http.DefaultTransport

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "decode: invalid character 'h' in literal true (expecting 'r')", err.Error())
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

	auth, err := url.PathUnescape(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if auth != `user[login]=user@example.com&user[password]=somepassword` {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(authResponse))
}

func plexBadHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "server's having a hard day", http.StatusInternalServerError)
}

func plexGarbageHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("this is definitely not json"))
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
	"/library/sections": `{ "MediaContainer": {
		"size": 2,
        "Directory": [
           { "Key": "1", "Type": "movie", "Title": "Movies" },
           { "Key": "2", "Type": "show", "Title": "Shows" }
        ]
    }}`,
	"/library/sections/1/all": `{ "MediaContainer" : {
        "Metadata": [
           { "guid": "1", "title": "foo" }
        ]
    }}`,
	"/library/sections/2/all": `{ "MediaContainer" : {
        "Metadata": [
           { "guid": "2", "title": "bar" }
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
