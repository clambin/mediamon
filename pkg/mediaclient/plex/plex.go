package plex

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clambin/httpclient"
	"github.com/clambin/mediamon/version"
	"net/http"
)

// API interface
//
//go:generate mockery --name API
type API interface {
	GetIdentity(context.Context) (identity IdentityResponse, err error)
	GetSessions(ctx context.Context) (sessions SessionsResponse, err error)
}

// Client calls the Plex APIs
type Client struct {
	httpclient.Caller
	URL       string
	AuthURL   string
	UserName  string
	Password  string
	authToken string
}

var _ API = &Client{}

// GetIdentity calls Plex' /identity endpoint. Mainly useful to get the server's version.
func (client *Client) GetIdentity(ctx context.Context) (identity IdentityResponse, err error) {
	err = client.call(ctx, "/identity", &identity)
	return
}

// GetSessions retrieves session information from the server.
func (client *Client) GetSessions(ctx context.Context) (sessions SessionsResponse, err error) {
	err = client.call(ctx, "/status/sessions", &sessions)
	return
}

// call the specified Plex API endpoint
func (client *Client) call(ctx context.Context, endpoint string, response interface{}) (err error) {
	err = client.authenticate(ctx)

	if err != nil {
		return
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, client.URL+endpoint, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Plex-Token", client.authToken)

	var resp *http.Response
	resp, err = client.Caller.Do(req)

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&response)
}

// authenticate logs in to plex.tv and gets an authentication token
// to be used for calls to the Plex server APIs
func (client *Client) authenticate(ctx context.Context) (err error) {
	if client.authToken != "" {
		return
	}

	authBody := fmt.Sprintf("user%%5Blogin%%5D=%s&user%%5Bpassword%%5D=%s",
		client.UserName,
		client.Password,
	)
	authURL := "https://plex.tv/users/sign_in.xml"
	if client.AuthURL != "" {
		authURL = client.AuthURL
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, authURL, bytes.NewBufferString(authBody))
	req.Header.Add("X-Plex-Product", "github.com/clambin/mediamon")
	req.Header.Add("X-Plex-Version", version.BuildVersion)
	req.Header.Add("X-Plex-Client-Identifier", "github.com/clambin/mediamon-v"+version.BuildVersion)

	var resp *http.Response
	resp, err = client.Caller.Do(req)

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("plex auth failed: %s", resp.Status)
	}

	// TODO: there's three different places in the response where the authToken appears.
	// Which is the officially supported version?
	var authResponse struct {
		XMLName             xml.Name `xml:"user"`
		AuthenticationToken string   `xml:"authenticationToken,attr"`
	}

	if err = xml.NewDecoder(resp.Body).Decode(&authResponse); err == nil {
		client.authToken = authResponse.AuthenticationToken
	}

	return
}
