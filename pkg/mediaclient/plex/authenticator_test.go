package plex

import (
	"context"
	"github.com/clambin/mediamon/v2/pkg/mediaclient/plex/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticator_RoundTrip(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(testutil.AuthHandler))

	server := httptest.NewServer(testutil.WithToken("some_token", testutil.Handler))
	defer server.Client()

	c := New("user@example.com", "somepassword", "", "", server.URL, nil)
	c.authenticator.authURL = authServer.URL

	resp, err := c.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, Identity{
		Claimed:           true,
		MachineIdentifier: "SomeUUID",
		Version:           "SomeVersion",
	}, resp)

	c.SetAuthToken("")
	c.HTTPClient.Transport.(*authenticator).password = "badpassword"

	_, err = c.GetIdentity(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plex auth: 403 Forbidden")

	authServer.Close()
	_, err = c.GetIdentity(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connect: connection refused")

}

func TestAuthenticator_Custom_RoundTripper(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(testutil.AuthHandler))
	defer authServer.Close()

	server := httptest.NewServer(testutil.WithToken("some_token", testutil.Handler))
	defer server.Client()

	c := New("user@example.com", "somepassword", "", "", server.URL, &dummyRoundTripper{next: http.DefaultTransport})
	c.authenticator.authURL = authServer.URL

	resp, err := c.GetIdentity(context.Background())
	require.NoError(t, err)
	assert.Equal(t, Identity{
		Claimed:           true,
		MachineIdentifier: "SomeUUID",
		Version:           "SomeVersion",
	}, resp)

	c.SetAuthToken("")
	c.HTTPClient.Transport.(*authenticator).password = "badpassword"

	_, err = c.GetIdentity(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "plex auth: 403 Forbidden")
}

var _ http.RoundTripper = &dummyRoundTripper{}

type dummyRoundTripper struct {
	next http.RoundTripper
}

func (d *dummyRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("X-Dummy", "foo")
	return d.next.RoundTrip(request)
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

	authServer := httptest.NewServer(http.HandlerFunc(testutil.AuthHandler))
	defer authServer.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.fields.UserName, tt.fields.Password, "", "", "", nil)
			c.authenticator.authURL = authServer.URL
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
