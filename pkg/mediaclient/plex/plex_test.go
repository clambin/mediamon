package plex_test

import (
	"context"
	"github.com/clambin/mediamon/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_Failures(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()
	testServer := httptest.NewServer(http.HandlerFunc(plexBadHandler))

	c := &plex.Client{
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

	_, err := c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.Equal(t, "500 "+http.StatusText(http.StatusInternalServerError), err.Error())

	testServer.Close()
	_, err = c.GetIdentity(context.Background())
	require.Error(t, err)
	assert.ErrorIs(t, err, unix.ECONNREFUSED)
}

func TestClient_Decode_Failure(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()
	testServer := httptest.NewServer(http.HandlerFunc(plexGarbageHandler))
	defer testServer.Close()

	c := &plex.Client{
		URL:      testServer.URL,
		AuthURL:  authServer.URL,
		UserName: "user@example.com",
		Password: "somepassword",
	}

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

func TestClient_GetAuthToken(t *testing.T) {
	type fields struct {
		AuthToken string
		UserName  string
		Password  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "authToken",
			fields:  fields{AuthToken: "1234"},
			want:    "1234",
			wantErr: assert.NoError,
		},
		{
			name:    "username / password",
			fields:  fields{UserName: "user@example.com", Password: "somepassword"},
			want:    "some_token",
			wantErr: assert.NoError,
		},
		{
			name:    "bad password",
			fields:  fields{UserName: "user@example.com", Password: "bad-password"},
			wantErr: assert.Error,
		},
	}

	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &plex.Client{
				AuthToken: tt.fields.AuthToken,
				AuthURL:   authServer.URL,
				UserName:  tt.fields.UserName,
				Password:  tt.fields.Password,
			}
			got, err := c.GetAuthToken(context.Background())
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
