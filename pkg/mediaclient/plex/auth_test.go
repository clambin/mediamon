package plex_test

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlexAuth(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(plexAuthHandler))
	defer authServer.Close()

	server := httptest.NewServer(authenticated("some_token", plexHandler))
	defer server.Client()

	c := plex.New("user@example.com", "somepassword", "", "", server.URL)
	c.HTTPClient.Transport.(*plex.Auth).AuthURL = authServer.URL

	resp, err := c.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, plex.Identity{
		Claimed:           true,
		MachineIdentifier: "SomeUUID",
		Version:           "SomeVersion",
	}, resp)

	c.SetAuthToken("")
	c.HTTPClient.Transport.(*plex.Auth).Password = "badpassword"

	_, err = c.GetIdentity(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plex auth: 403 Forbidden")
}

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
			c := plex.New(tt.fields.UserName, tt.fields.Password, "", "", "")
			c.HTTPClient.Transport.(*plex.Auth).AuthURL = authServer.URL
			if tt.fields.AuthToken != "" {
				c.SetAuthToken(tt.fields.AuthToken)
			}

			got, err := c.GetAuthToken(context.Background())
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func authenticated(token string, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Plex-Token") != token {
			writer.WriteHeader(http.StatusForbidden)
			return
		}
		next(writer, request)
	}
}

func plexHandler(w http.ResponseWriter, req *http.Request) {
	if response, ok := plexResponses[req.URL.Path]; ok {
		_, _ = w.Write([]byte(response))
	} else {
		http.Error(w, "endpoint not implemented: "+req.URL.Path, http.StatusNotFound)
	}
}
